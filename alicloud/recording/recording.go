// Package recording records and replays acceptance-test HTTP interactions.
// It is opt-in, process-scoped, and never falls back to the network in replay mode.
package recording

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ModeEnv = "TF_ACC_RECORDING_MODE"
	PathEnv = "TF_ACC_RECORDING_PATH"
	// CaseEnv is an optional legacy assertion. New callers select the case via -run.
	CaseEnv = "TF_ACC_RECORDING_CASE"
	maxBody = 8 << 20
)

type request struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}
type interaction struct {
	Request request     `json:"request"`
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}
type cassette struct {
	FormatVersion    int           `json:"formatVersion"`
	MatcherVersion   string        `json:"matcherVersion"`
	SanitizerVersion string        `json:"sanitizerVersion"`
	TestCase         string        `json:"testCase"`
	Interactions     []interaction `json:"interactions"`
}

type Result struct {
	FormatVersion      int    `json:"formatVersion"`
	MatcherVersion     string `json:"matcherVersion"`
	SanitizerVersion   string `json:"sanitizerVersion"`
	TestCase           string `json:"testCase"`
	Mode               string `json:"mode"`
	InteractionCount   int    `json:"interactionCount"`
	Sanitized          bool   `json:"sanitized"`
	UnmatchedRequests  int    `json:"unmatchedRequests"`
	UnusedInteractions int    `json:"unusedInteractions"`
	Success            bool   `json:"success"`
}

type Session struct {
	mu              sync.Mutex
	mode, dir       string
	tape            cassette
	next, unmatched int
	failure         error
	random          *mathrand.Rand
	seed            int64
	secrets         []string
	closed          bool
}

var active struct {
	sync.RWMutex
	session *Session
}

func Mode() string {
	m := os.Getenv(ModeEnv)
	if m == "" {
		return "live"
	}
	return m
}
func Enabled() bool   { return Mode() != "live" }
func Replaying() bool { return Mode() == "replay" }

// New creates an isolated session. A record directory must be new: baselines are
// immutable, including after a failed recording. Each retry uses a new directory.
func New(mode, dir, testCase string) (*Session, error) {
	if mode != "record" && mode != "replay" {
		return nil, fmt.Errorf("unsupported recording mode %q", mode)
	}
	if dir == "" || testCase == "" || !strings.HasPrefix(testCase, "Test") || strings.ContainsAny(testCase, "/\\\x00") {
		return nil, errors.New("recording requires a directory and one exact test name")
	}
	s := &Session{mode: mode, dir: dir, tape: cassette{FormatVersion: 1, MatcherVersion: "1", SanitizerVersion: "1", TestCase: testCase}}
	for _, key := range []string{"ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "ALICLOUD_SECURITY_TOKEN", "ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET", "ALIBABA_CLOUD_SECURITY_TOKEN"} {
		if value := os.Getenv(key); value != "" {
			s.secrets = append(s.secrets, value)
		}
	}
	var seed int64
	if mode == "record" {
		if err := os.Mkdir(dir, 0700); err != nil {
			return nil, fmt.Errorf("create new recording directory: %w", err)
		}
		var b [8]byte
		if _, err := rand.Read(b[:]); err != nil {
			return nil, err
		}
		seed = int64(binary.LittleEndian.Uint64(b[:]))
		if err := os.WriteFile(filepath.Join(dir, "seed"), []byte(strconv.FormatInt(seed, 10)), 0600); err != nil {
			return nil, err
		}
	} else {
		data, err := os.ReadFile(filepath.Join(dir, "cassette.json"))
		if err != nil {
			return nil, err
		}
		if len(data) > 128<<20 {
			return nil, errors.New("recording exceeds size limit")
		}
		if err = json.Unmarshal(data, &s.tape); err != nil {
			return nil, errors.New("invalid recording JSON")
		}
		if s.tape.FormatVersion != 1 || s.tape.MatcherVersion != "1" || s.tape.SanitizerVersion != "1" || s.tape.TestCase != testCase {
			return nil, errors.New("incompatible recording format, rules, or test case")
		}
		data, err = os.ReadFile(filepath.Join(dir, "seed"))
		if err != nil {
			return nil, err
		}
		seed, err = strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return nil, errors.New("invalid recording seed")
		}
	}
	s.seed = seed
	s.random = mathrand.New(mathrand.NewSource(seed))
	return s, nil
}

// WrapMain fits the testing framework's Run interface, so existing TestMain
// sweep handling remains unchanged. The caller still calls resource.TestMain.
func WrapMain(m interface{ Run() int }) interface{ Run() int } { return mainRunner{m} }

type mainRunner struct{ m interface{ Run() int } }

func (r mainRunner) Run() int {
	if !Enabled() {
		return r.m.Run()
	}
	f := flag.Lookup("test.run")
	if f == nil {
		fmt.Fprintln(os.Stderr, "recording requires -run with one exact anchored test name")
		return 1
	}
	testCase, err := selectedCase(f.Value.String(), os.Getenv(CaseEnv))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	s, err := New(Mode(), os.Getenv(PathEnv), testCase)
	if err != nil {
		fmt.Fprintln(os.Stderr, "recording initialization failed:", err)
		return 1
	}
	restoreRandom, err := seedGlobalRandom(s.seed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer restoreRandom()
	active.Lock()
	active.session = s
	active.Unlock()
	code := r.m.Run()
	active.Lock()
	active.session = nil
	active.Unlock()
	if err = s.Close(code == 0); err != nil {
		fmt.Fprintln(os.Stderr, "recording validation failed:", err)
		return 1
	}
	return code
}

// RunMain lets service-package TestMain functions use the same recording entry
// without importing the aggregate Provider. Ordinary execution is unchanged.
func RunMain(m interface{ Run() int }) int {
	flag.Parse()
	if f := flag.Lookup("sweep"); Enabled() && f != nil && f.Value.String() != "" {
		fmt.Fprintln(os.Stderr, "recording cannot be combined with resource sweeping")
		return 1
	}
	return WrapMain(m).Run()
}

// Seed the legacy math/rand source used by both HashiCorp acceptance helper
// packages and existing tests. Go 1.24+ requires randseednop=0 for Seed to act.
// Crypto randomness and ordinary live-mode tests are never changed.
func seedGlobalRandom(seed int64) (func(), error) {
	previous, existed := os.LookupEnv("GODEBUG")
	if err := os.Setenv("GODEBUG", previous+",randseednop=0"); err != nil {
		return nil, err
	}
	mathrand.Seed(seed)
	return func() {
		if existed {
			os.Setenv("GODEBUG", previous)
		} else {
			os.Unsetenv("GODEBUG")
		}
	}, nil
}

// RequireSession prevents SDK configuration/discovery before an opt-in TestMain.
func RequireSession() error {
	if !Enabled() {
		return nil
	}
	active.RLock()
	defer active.RUnlock()
	if active.session == nil {
		return errors.New("recording requires a recording-enabled TestMain")
	}
	return nil
}

// ClientToken preserves idempotency values in strict matching. Only an active
// session replaces the live generator; unique tokens replay from the saved seed.
func ClientToken(action string, live func() string) string {
	active.RLock()
	s := active.session
	active.RUnlock()
	if s == nil {
		return live()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var entropy [16]byte
	s.random.Read(entropy[:])
	if len(action) > 28 {
		action = action[:28]
	}
	return "TF-" + action + "-" + hex.EncodeToString(entropy[:])
}

// Only a literal top-level Go test identifier is accepted. Regex groups,
// alternatives, prefixes and subtest selectors could execute several cases.
func selectedCase(selection, legacyCase string) (string, error) {
	match := regexp.MustCompile(`^\^(Test[A-Za-z0-9_]+)\$$`).FindStringSubmatch(selection)
	if match == nil {
		return "", errors.New("recording requires -run with one exact anchored test name, such as ^TestExample$")
	}
	if legacyCase != "" && legacyCase != match[1] {
		return "", errors.New("legacy TF_ACC_RECORDING_CASE must match the test selected by -run")
	}
	return match[1], nil
}

// RandInt preserves existing randomness outside recording. Tests using other
// randomness helpers must be adapted before their recordings are reusable.
func RandInt(fallback func() int) int {
	active.RLock()
	s := active.session
	active.RUnlock()
	if s == nil {
		return fallback()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.random.Intn(900000000) + 100000000
}

// Transport wraps a specific SDK client; it never replaces global HTTP state.
func Transport(base http.RoundTripper) http.RoundTripper {
	if !Enabled() {
		return base
	}
	active.RLock()
	s := active.session
	active.RUnlock()
	if s == nil {
		return errorTransport{errors.New("recording requires a recording-enabled TestMain")}
	}
	return s.Transport(base)
}
func (s *Session) Transport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return transport{s, base}
}

type transport struct {
	s    *Session
	base http.RoundTripper
}
type errorTransport struct{ err error }

func (t errorTransport) RoundTrip(*http.Request) (*http.Response, error) { return nil, t.err }

// Call is the Darabonba HTTP-client adapter used by generated v2 clients.
func (t transport) Call(r *http.Request, _ *http.Transport) (*http.Response, error) {
	return t.RoundTrip(r)
}

type HTTPClient struct{ Timeout time.Duration }

func (client HTTPClient) Call(r *http.Request, base *http.Transport) (*http.Response, error) {
	return (&http.Client{Transport: Transport(base), Timeout: client.Timeout}).Do(r)
}

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	s := t.s
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, errors.New("recording session is closed")
	}
	if s.failure != nil {
		return nil, s.failure
	}
	body, err := readBody(r.Body)
	if err != nil {
		return nil, s.fail(err)
	}
	if r.Body != nil {
		r.Body = io.NopCloser(bytes.NewReader(body))
	}
	normalized, err := s.normalize(r, body)
	if err != nil {
		return nil, s.fail(err)
	}
	if s.mode == "replay" {
		if s.next >= len(s.tape.Interactions) || !reflect.DeepEqual(normalized, s.tape.Interactions[s.next].Request) {
			s.unmatched++
			return nil, s.fail(fmt.Errorf("recording request mismatch at interaction %d (%s %s); network fallback is disabled", s.next, r.Method, r.URL.Path))
		}
		i := s.tape.Interactions[s.next]
		s.next++
		return &http.Response{StatusCode: i.Status, Status: fmt.Sprintf("%d %s", i.Status, http.StatusText(i.Status)), Header: i.Headers.Clone(), Body: io.NopCloser(strings.NewReader(i.Body)), ContentLength: int64(len(i.Body)), Request: r, Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1}, nil
	}
	response, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, s.fail(errors.New("recording upstream transport failed"))
	}
	data, err := readBody(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, s.fail(err)
	}
	response.Body = io.NopCloser(bytes.NewReader(data))
	clean, err := s.payload(data, response.Header.Get("Content-Type"), false)
	if err != nil {
		return nil, s.fail(err)
	}
	headers, err := s.headers(response.Header, false)
	if err != nil {
		return nil, s.fail(err)
	}
	s.tape.Interactions = append(s.tape.Interactions, interaction{normalized, response.StatusCode, headers, clean})
	return response, nil
}

// Unsupported poisons an active recording even if a caller swallows the error.
func Unsupported(client string) error {
	err := fmt.Errorf("recording v1 does not support %s", client)
	active.RLock()
	s := active.session
	active.RUnlock()
	if s != nil {
		s.mu.Lock()
		s.fail(err)
		s.mu.Unlock()
	}
	return err
}

func (s *Session) fail(err error) error { s.failure = err; return err }
func readBody(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	data, err := io.ReadAll(io.LimitReader(r, maxBody+1))
	if len(data) > maxBody {
		return nil, errors.New("recording body exceeds limit")
	}
	return data, err
}

func (s *Session) Close(success bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errors.New("recording session already closed")
	}
	s.closed = true
	unused := 0
	if s.mode == "replay" {
		unused = len(s.tape.Interactions) - s.next
	}
	if len(s.tape.Interactions) == 0 && s.failure == nil {
		s.failure = errors.New("no interactions recorded")
	}
	if unused > 0 && s.failure == nil {
		s.failure = fmt.Errorf("%d recording interactions were not consumed", unused)
	}
	valid := success && s.failure == nil
	if s.mode == "record" && valid {
		if err := writeJSON(filepath.Join(s.dir, "cassette.json"), s.tape); err != nil {
			s.failure = err
			valid = false
		}
	}
	result := Result{1, "1", "1", s.tape.TestCase, s.mode, len(s.tape.Interactions), valid, s.unmatched, unused, valid}
	// Replay writes its result outside the immutable input directory when requested.
	resultPath := os.Getenv("TF_ACC_RECORDING_RESULT_PATH")
	if resultPath == "" {
		resultPath = filepath.Join(s.dir, "result.json")
	}
	if err := writeJSON(resultPath, result); err != nil {
		return err
	}
	if !success && s.failure == nil {
		return errors.New("test execution did not pass")
	}
	return s.failure
}
func writeJSON(name string, value interface{}) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(name, append(data, '\n'), 0600)
}

// Wait is available for explicitly adapted polling loops. It preserves cancellation
// while avoiding wall-clock delay during replay. Unadapted waits remain unchanged.
func Wait(done <-chan struct{}, delay time.Duration) error {
	if Replaying() {
		select {
		case <-done:
			return errors.New("recording wait canceled")
		default:
			return nil
		}
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-done:
		return errors.New("recording wait canceled")
	case <-timer.C:
		return nil
	}
}

func ignoredHeader(key string) bool {
	switch strings.ToLower(key) {
	case "proxy-authorization", "x-acs-accesskey-id", "x-log-date", "x-log-security-token", "x-log-credential", "x-log-signature", "authorization", "x-acs-security-token", "securitytoken", "accesskeyid", "signature", "signaturenonce", "signaturemethod", "signatureversion", "timestamp", "date", "x-acs-date", "x-acs-signature-nonce", "x-acs-signature-method", "x-acs-signature-version", "x-acs-credentials-provider", "user-agent", "x-acs-client-version", "x-sdk-client", "x-acs-content-sha256", "accept-encoding", "content-length", "connection":
		return true
	}
	return false
}
func ignoredQuery(key string) bool {
	switch strings.ToLower(key) {
	case "securitytoken", "accesskeyid", "signature", "signaturenonce", "signaturemethod", "signatureversion", "timestamp":
		return true
	}
	return false
}
func sensitive(key string) bool {
	k := strings.ToLower(strings.NewReplacer("-", "", "_", "", ".", "").Replace(key))
	for _, p := range []string{"password", "secret", "privatekey", "securitytoken", "authorization", "credential", "accesstoken", "refreshtoken", "idtoken", "clientsecret", "keymaterial"} {
		if strings.Contains(k, p) {
			return true
		}
	}
	return k == "accesskeyid" || k == "useraccesskeyid"
}
func identity(key string) bool {
	switch strings.ToLower(key) {
	case "accountid", "ownerid", "uid", "userid":
		return true
	}
	return false
}
func token(value string) string {
	if strings.HasPrefix(value, "recording-id-") {
		return value
	}
	sum := sha256.Sum256([]byte(value))
	return "recording-id-" + hex.EncodeToString(sum[:8])
}
func (s *Session) cleanString(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if (strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")) && json.Valid([]byte(trimmed)) {
		var embedded interface{}
		if err := json.Unmarshal([]byte(trimmed), &embedded); err != nil {
			return "", errors.New("invalid embedded JSON")
		}
		clean, err := s.walk(embedded, "")
		if err != nil {
			return "", err
		}
		data, err := json.Marshal(clean)
		return string(data), err
	}
	lower := strings.ToLower(value)
	for _, marker := range []string{"password=", "password:", "private_key=", "access_key_secret=", "security_token=", "securitytoken=", "accesskeyid=", "signature=", "token=", "-----begin "} {
		if strings.Contains(lower, marker) {
			return "", errors.New("potential embedded credential material in recording")
		}
	}
	for _, secret := range s.secrets {
		if strings.Contains(value, secret) {
			return "", errors.New("credential material found outside authentication fields")
		}
	}
	if strings.Contains(value, "-----BEGIN ") || strings.Contains(value, "LTAI") {
		return "", errors.New("potential credential material in recording")
	}
	return value, nil
}
func (s *Session) walk(value interface{}, key string) (interface{}, error) {
	if sensitive(key) {
		return nil, errors.New("sensitive business field is not supported by recording v1")
	}
	switch v := value.(type) {
	case map[string]interface{}:
		for k, x := range v {
			clean, err := s.walk(x, k)
			if err != nil {
				return nil, err
			}
			v[k] = clean
		}
		return v, nil
	case []interface{}:
		for i, x := range v {
			clean, err := s.walk(x, key)
			if err != nil {
				return nil, err
			}
			v[i] = clean
		}
		return v, nil
	case string:
		if identity(key) {
			return token(v), nil
		}
		return s.cleanString(v)
	default:
		if identity(key) {
			return token(fmt.Sprint(v)), nil
		}
		return v, nil
	}
}
func (s *Session) payload(body []byte, contentType string, form bool) (string, error) {
	if len(body) == 0 {
		return "", nil
	}
	if form || strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return "", errors.New("invalid form body")
		}
		if err = s.values(values); err != nil {
			return "", err
		}
		return values.Encode(), nil
	}
	if strings.Contains(contentType, "json") || json.Valid(body) {
		if !json.Valid(body) {
			return "", errors.New("invalid or trailing JSON body")
		}
		var v interface{}
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.UseNumber()
		if err := dec.Decode(&v); err != nil {
			return "", errors.New("invalid JSON body")
		}
		v, err := s.walk(v, "")
		if err != nil {
			return "", err
		}
		data, err := json.Marshal(v)
		return string(data), err
	}
	return "", errors.New("unsupported recording body type; only JSON and form bodies are supported")
}
func (s *Session) values(v url.Values) error {
	for k, vs := range v {
		if ignoredQuery(k) {
			delete(v, k)
			continue
		}
		if sensitive(k) {
			return errors.New("sensitive business query parameter is not supported by recording v1")
		}
		for i, x := range vs {
			if identity(k) {
				vs[i] = token(x)
			} else {
				clean, err := s.cleanString(x)
				if err != nil {
					return err
				}
				vs[i] = clean
			}
		}
	}
	return nil
}
func (s *Session) headers(h http.Header, request bool) (http.Header, error) {
	out := make(http.Header)
	for k, vs := range h {
		if ignoredHeader(k) || strings.EqualFold(k, "Set-Cookie") || strings.EqualFold(k, "Cookie") || strings.EqualFold(k, "Content-Length") || strings.EqualFold(k, "Server") || strings.EqualFold(k, "X-Request-Id") || strings.EqualFold(k, "X-Acs-Request-Id") {
			continue
		}
		if sensitive(k) {
			return nil, errors.New("unsupported sensitive recording header")
		}
		for _, v := range vs {
			clean, err := s.cleanString(v)
			if err != nil {
				return nil, err
			}
			out.Add(k, clean)
		}
	}
	return out, nil
}
func (s *Session) normalize(r *http.Request, body []byte) (request, error) {
	u := *r.URL
	if u.User != nil {
		return request{}, errors.New("URL credentials are not supported by recording")
	}
	u.Fragment = ""
	if _, err := s.cleanString(u.Host); err != nil {
		return request{}, err
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return request{}, errors.New("invalid recording query encoding")
	}
	if err := s.values(q); err != nil {
		return request{}, err
	}
	u.RawQuery = q.Encode()
	path, err := s.cleanString(u.EscapedPath())
	if err != nil {
		return request{}, err
	}
	_ = path
	headers, err := s.headers(r.Header, true)
	if err != nil {
		return request{}, err
	}
	clean, err := s.payload(body, headers.Get("Content-Type"), false)
	if err != nil {
		return request{}, err
	}
	return request{r.Method, u.String(), headers, clean}, nil
}
