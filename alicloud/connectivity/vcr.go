package connectivity

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// vcrServer is the singleton VCR proxy server, started once per process.
var (
	vcrOnce      sync.Once
	vcrLocalAddr string // "127.0.0.1:PORT"
	vcrRecorder  *recorder.Recorder
	vcrCleanup   func()
)

// VCRRandSeed returns a deterministic seed derived from VCR_PATH.
// Returns 0 when VCR is not enabled.
func VCRRandSeed() int64 {
	p := os.Getenv("VCR_PATH")
	if p == "" {
		return 0
	}
	h := sha256.Sum256([]byte(p))
	return int64(binary.LittleEndian.Uint64(h[:8]))
}

// VCRRandIntRange is a drop-in replacement for acctest.RandIntRange.
// In VCR mode it returns deterministic values; otherwise falls through
// to a time-seeded source.
var vcrRandOnce sync.Once
var vcrRandSrc *rand.Rand

func VCRRandIntRange(min, max int) int {
	if os.Getenv("VCR_PATH") == "" {
		// Not in VCR mode — behave like acctest.RandIntRange
		return min + rand.New(rand.NewSource(rand.Int63())).Intn(max-min)
	}
	vcrRandOnce.Do(func() {
		vcrRandSrc = rand.New(rand.NewSource(VCRRandSeed()))
	})
	return min + vcrRandSrc.Intn(max-min)
}

// VCRLocalAddr returns the local VCR proxy address if VCR is enabled, or "".
// On first call it lazily starts the proxy server.
func VCRLocalAddr() string {
	if os.Getenv("VCR_PATH") == "" {
		return ""
	}
	vcrOnce.Do(startVCR)
	return vcrLocalAddr
}

// StopVCR stops the VCR recorder and saves the cassette. Call in TestMain or defer.
func StopVCR() {
	if vcrCleanup != nil {
		vcrCleanup()
	}
}

func startVCR() {
	vcrPath := os.Getenv("VCR_PATH")
	if vcrPath == "" {
		return
	}

	mode := os.Getenv("VCR_MODE")

	var transport http.RoundTripper
	if mode == "replay" {
		// Custom replay transport with sequential consumption + fallback
		rt, err := newVcrReplayTransport(vcrPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[VCR] Failed to load cassette: %v\n", err)
			return
		}
		transport = rt
	} else {
		// Use go-vcr recorder for recording
		rec, err := newVCRRecorder(vcrPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[VCR] Failed to create recorder: %v\n", err)
			return
		}
		vcrRecorder = rec
		transport = rec
	}

	// Start local HTTP reverse proxy backed by VCR
	server := httptest.NewServer(&vcrProxyHandler{rec: transport})
	vcrLocalAddr = strings.TrimPrefix(server.URL, "http://")

	// Bypass system HTTP proxy for VCR server (tea SDK does exact host:port match)
	noProxy := os.Getenv("NO_PROXY")
	if noProxy != "" {
		os.Setenv("NO_PROXY", noProxy+","+vcrLocalAddr)
	} else {
		os.Setenv("NO_PROXY", vcrLocalAddr)
	}

	fmt.Fprintf(os.Stderr, "[VCR] Proxy at %s  mode=%s  cassette=%s\n", vcrLocalAddr, mode, vcrPath)

	vcrCleanup = func() {
		server.Close()
		if vcrRecorder != nil {
			if err := vcrRecorder.Stop(); err != nil {
				fmt.Fprintf(os.Stderr, "[VCR] Stop error: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "[VCR] Cassette saved\n")
			}
		} else {
			fmt.Fprintf(os.Stderr, "[VCR] Cassette saved\n")
		}
	}
}

// newVCRRecorder creates a go-vcr recorder for recording mode.
func newVCRRecorder(cassettePath string) (*recorder.Recorder, error) {
	return recorder.New(cassettePath, []recorder.Option{
		recorder.WithMode(recorder.ModeRecordOnly),
		recorder.WithSkipRequestLatency(true),
		recorder.WithMatcher(vcrMatcher),
		recorder.WithHook(sanitizeCassette, recorder.BeforeSaveHook),
	}...)
}

// ---------- Custom replay transport ----------

// vcrReplayTransport replays cassette interactions with sequential consumption.
// When a request matches multiple recorded interactions (e.g. DescribeVpcAttribute
// called before and after ModifyVpcAttribute), it returns them in order.
// If all matching interactions are consumed, it falls back to the last consumed
// match — handling extra calls that weren't in the recording.
type vcrReplayTransport struct {
	interactions []*cassette.Interaction
	mu           sync.Mutex
	consumed     map[int]bool // interaction index -> consumed
}

func newVcrReplayTransport(cassettePath string) (*vcrReplayTransport, error) {
	cas, err := cassette.Load(cassettePath)
	if err != nil {
		return nil, err
	}
	t := &vcrReplayTransport{
		interactions: cas.Interactions,
		consumed:     make(map[int]bool),
	}
	return t, nil
}

func (t *vcrReplayTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Read body so we can match against it (and restore for potential retries)
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// First pass: find first unconsumed matching interaction
	lastConsumedMatch := -1
	for idx, interaction := range t.interactions {
		// Restore body for each matcher call
		if bodyBytes != nil {
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		if !vcrMatcher(r, interaction.Request) {
			continue
		}
		if t.consumed[idx] {
			lastConsumedMatch = idx
			continue
		}
		// Found unconsumed match — consume and return
		t.consumed[idx] = true
		return buildHTTPResponse(interaction), nil
	}

	// Fallback: all matching interactions consumed, reuse the last consumed match
	if lastConsumedMatch >= 0 {
		return buildHTTPResponse(t.interactions[lastConsumedMatch]), nil
	}

	return nil, fmt.Errorf("requested interaction not found")
}

// buildHTTPResponse converts a cassette interaction to an http.Response.
func buildHTTPResponse(i *cassette.Interaction) *http.Response {
	header := make(http.Header)
	for k, vals := range i.Response.Headers {
		for _, v := range vals {
			header.Add(k, v)
		}
	}
	return &http.Response{
		Status:        i.Response.Status,
		StatusCode:    i.Response.Code,
		Proto:         i.Response.Proto,
		ProtoMajor:    i.Response.ProtoMajor,
		ProtoMinor:    i.Response.ProtoMinor,
		Header:        header,
		Body:          io.NopCloser(strings.NewReader(i.Response.Body)),
		ContentLength: i.Response.ContentLength,
	}
}

// ---------- VCR proxy handler ----------

type vcrProxyHandler struct {
	rec http.RoundTripper
}

func (h *vcrProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read body
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body.Close()
	}

	// Read the original host from the per-request query parameter.
	originalHost := r.URL.Query().Get("__vcr_host")
	if originalHost == "" {
		originalHost = r.Host // fallback
	}

	outURL := fmt.Sprintf("https://%s%s", originalHost, r.URL.RequestURI())
	outReq, err := http.NewRequest(r.Method, outURL, bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for k, vv := range r.Header {
		// Skip host headers — they contain the local proxy address.
		if strings.EqualFold(k, "Host") {
			continue
		}
		for _, v := range vv {
			outReq.Header.Add(k, v)
		}
	}
	outReq.Host = originalHost

	resp, err := h.rec.RoundTrip(outReq)
	if err != nil {
		// Extract Action for better error logging
		action := r.URL.Query().Get("Action")
		if action == "" {
			action = originalHost + r.URL.Path
		}
		fmt.Fprintf(os.Stderr, "[VCR] %s %s → error: %v\n", r.Method, action, err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	n, _ := io.Copy(w, resp.Body)

	// Extract Action from query for concise logging
	action := r.URL.Query().Get("Action")
	if action == "" {
		action = r.URL.Path
	}
	fmt.Fprintf(os.Stderr, "[VCR] %s %s → %d (%d bytes)\n", r.Method, action, resp.StatusCode, n)
}

// ---------- Matcher ----------

func vcrMatcher(r *http.Request, i cassette.Request) bool {
	if r.Method != i.Method {
		return false
	}
	reqURL, err := url.Parse(i.URL)
	if err != nil {
		return r.URL.String() == i.URL
	}
	// Match host + path
	if r.URL.Host != reqURL.Host || r.URL.Path != reqURL.Path {
		return false
	}
	// Match query params (ignoring signature/auth params that change per request)
	if filterSignatureParams(r.URL.Query()).Encode() !=
		filterSignatureParams(reqURL.Query()).Encode() {
		return false
	}
	// Match request body (form-encoded API parameters like VpcId, VpcName, etc.)
	// This distinguishes CreateVpc(name=A) from CreateVpc(name=B) and
	// DescribeVpc(id=X) from DescribeVpc(id=Y).
	// Normalize numbered params (Tag.N.Key etc.) to handle Go map iteration
	// non-determinism between recording and replay.
	return normalizeNumberedParams(filterSignatureParams(requestFormValues(r))).Encode() ==
		normalizeNumberedParams(filterSignatureParams(i.Form)).Encode()
}

// requestFormValues extracts all form values (query string + body) from an
// http.Request, matching go-vcr's cassette.Request.Form behavior.
func requestFormValues(r *http.Request) url.Values {
	// Merge query params + body params, same as http.Request.ParseForm()
	vals := make(url.Values)
	for k, v := range r.URL.Query() {
		vals[k] = v
	}
	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err == nil {
			r.Body = io.NopCloser(bytes.NewReader(body))
			if bodyVals, err := url.ParseQuery(string(body)); err == nil {
				for k, v := range bodyVals {
					vals[k] = append(vals[k], v...)
				}
			}
		}
	}
	return vals
}

// sanitizeCassette removes sensitive data from recorded interactions.
// Follows AWS provider convention: strip auth headers and redact credentials
// in URLs/forms. Response bodies are left intact (cassettes are not committed).
func sanitizeCassette(i *cassette.Interaction) error {
	// Request headers
	delete(i.Request.Headers, "Authorization")
	delete(i.Request.Headers, "X-Acs-Security-Token")

	// Request URL: redact credential parameters
	i.Request.URL = redactURLParams(i.Request.URL)

	// Request Form fields
	for _, k := range []string{"AccessKeyId", "SecurityToken", "Signature", "SignatureNonce"} {
		if _, ok := i.Request.Form[k]; ok {
			i.Request.Form[k] = []string{"REDACTED"}
		}
	}
	// Strip VCR-internal parameter
	delete(i.Request.Form, "__vcr_host")

	return nil
}

// redactURLParams replaces sensitive query parameters in a URL string.
func redactURLParams(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	for _, k := range []string{"AccessKeyId", "SecurityToken", "Signature", "SignatureNonce"} {
		if q.Get(k) != "" {
			q.Set(k, "REDACTED")
		}
	}
	// Strip VCR-internal parameter from recorded URL
	q.Del("__vcr_host")
	u.RawQuery = q.Encode()
	return u.String()
}

// numberedParamRe matches Alibaba Cloud API indexed params like Tag.1.Key,
// Tag.2.Value, Filter.1.Key, etc. Go map iteration is non-deterministic, so
// the same map{"A":"1","B":"2"} might produce Tag.1.Key=A or Tag.1.Key=B
// across runs. This function renumbers them canonically (sorted by values).
var numberedParamRe = regexp.MustCompile(`^(.+)\.(\d+)\.(.+)$`)

func normalizeNumberedParams(vals url.Values) url.Values {
	result := make(url.Values)

	// prefix -> index -> field -> value
	type entry struct {
		fields map[string]string
	}
	groups := make(map[string]map[int]*entry)

	for k, v := range vals {
		m := numberedParamRe.FindStringSubmatch(k)
		if m == nil {
			result[k] = v
			continue
		}
		prefix := m[1]
		idx, _ := strconv.Atoi(m[2])
		field := m[3]
		if groups[prefix] == nil {
			groups[prefix] = make(map[int]*entry)
		}
		if groups[prefix][idx] == nil {
			groups[prefix][idx] = &entry{fields: make(map[string]string)}
		}
		groups[prefix][idx].fields[field] = v[0]
	}

	// Sort each group's entries by their field values and renumber
	for prefix, items := range groups {
		entries := make([]*entry, 0, len(items))
		for _, e := range items {
			entries = append(entries, e)
		}
		sort.Slice(entries, func(i, j int) bool {
			return entrySortKey(entries[i].fields) < entrySortKey(entries[j].fields)
		})
		for i, e := range entries {
			for field, value := range e.fields {
				result[fmt.Sprintf("%s.%d.%s", prefix, i+1, field)] = []string{value}
			}
		}
	}

	return result
}

// entrySortKey builds a canonical string from a map for sorting.
func entrySortKey(fields map[string]string) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(fields[k])
		sb.WriteByte('&')
	}
	return sb.String()
}

func filterSignatureParams(q url.Values) url.Values {
	filtered := make(url.Values)
	skip := map[string]bool{
		// Auth/signature params (change every request)
		"Signature": true, "SignatureNonce": true,
		"Timestamp": true, "SecurityToken": true,
		"AccessKeyId": true, "SignatureMethod": true,
		"SignatureVersion": true,
		// Idempotency token (random UUID per apply)
		"ClientToken": true,
		// VCR internal: original API host passed per-request
		"__vcr_host": true,
	}
	for k, v := range q {
		if !skip[k] {
			filtered[k] = v
		}
	}
	return filtered
}
