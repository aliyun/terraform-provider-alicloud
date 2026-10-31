package recording

import (
	"context"
	"encoding/json"
	"fmt"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"io"
	mathrand "math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestRecordReplayLifecycle(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"RequestId":"test","State":%q}`, r.URL.Query().Get("Action"))
	}))
	defer server.Close()
	dir := filepath.Join(t.TempDir(), "TestLifecycle")
	record, err := New("record", dir, "TestLifecycle")
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Transport: record.Transport(http.DefaultTransport)}
	for _, action := range []string{"Create", "Read", "Update", "Read", "Delete", "Read"} {
		req, _ := http.NewRequest("POST", server.URL+"/?Action="+action+"&Version=2020-01-01&AccessKeyId=real-key&Signature=secret&Timestamp=first", strings.NewReader("Value=expected"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "secret-authorization")
		req.Header.Set("X-Acs-Security-Token", "secret-token")
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
	}
	if err := record.Close(true); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "cassette.json"))
	for _, secret := range []string{"real-key", "secret-authorization", "secret-token", "Timestamp", "Signature"} {
		if strings.Contains(string(data), secret) {
			t.Fatalf("persisted sensitive/dynamic field %s", secret)
		}
	}
	server.Close()
	replay, err := New("replay", dir, "TestLifecycle")
	if err != nil {
		t.Fatal(err)
	}
	client.Transport = replay.Transport(errorTransport{fmt.Errorf("network must not be called")})
	for _, action := range []string{"Create", "Read", "Update", "Read", "Delete", "Read"} {
		req, _ := http.NewRequest("POST", server.URL+"/?Version=2020-01-01&Action="+action+"&AccessKeyId=dummy&Signature=different&Timestamp=second", strings.NewReader("Value=expected"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
	}
	if err := replay.Close(true); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 6 {
		t.Fatalf("unexpected network calls: %d", calls.Load())
	}
	result, _ := os.ReadFile(filepath.Join(dir, "result.json"))
	var got Result
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatal(err)
	}
	if !got.Success || got.UnusedInteractions != 0 || got.UnmatchedRequests != 0 {
		t.Fatalf("bad replay result %s", result)
	}
}

func fixture(t *testing.T) (string, string) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"State":"Ready"}`)
	}))
	defer server.Close()
	dir := filepath.Join(t.TempDir(), "TestOne")
	s, err := New("record", dir, "TestOne")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", server.URL+"/?Action=Read&Version=1&Value=original", nil)
	resp, err := s.Transport(nil).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if err = s.Close(true); err != nil {
		t.Fatal(err)
	}
	return dir, server.URL
}
func TestReplayRejectsChangedBusinessParameterAndExtraRequest(t *testing.T) {
	for _, change := range []string{"Action=Read&Version=1&Value=changed", "Action=Delete&Version=1&Value=original", "Action=Read&Version=2&Value=original"} {
		t.Run(change, func(t *testing.T) {
			dir, base := fixture(t)
			s, err := New("replay", dir, "TestOne")
			if err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest("GET", base+"/?"+change, nil)
			if _, err = s.Transport(errorTransport{fmt.Errorf("network")}).RoundTrip(req); err == nil {
				t.Fatal("accepted changed business request")
			}
			if err = s.Close(true); err == nil {
				t.Fatal("accepted mismatched replay")
			}
		})
	}
	dir, base := fixture(t)
	s, _ := New("replay", dir, "TestOne")
	req, _ := http.NewRequest("GET", base+"/?Action=Read&Version=1&Value=original", nil)
	resp, err := s.Transport(nil).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if _, err = s.Transport(nil).RoundTrip(req); err == nil {
		t.Fatal("extra request accepted")
	}
	if err = s.Close(true); err == nil {
		t.Fatal("extra request did not fail close")
	}
}
func TestReplayRejectsUnusedInteractions(t *testing.T) {
	dir, _ := fixture(t)
	s, err := New("replay", dir, "TestOne")
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Close(true); err == nil {
		t.Fatal("unused interactions accepted")
	}
}
func TestRecordingRejectsSensitiveFieldsBeforePersistence(t *testing.T) {
	for _, payload := range []string{`{"Password":"dont-persist"}`, `{"AccessKeySecret":"dont-persist"}`, `{"nested":{"PrivateKey":"dont-persist"}}`} {
		t.Run(payload, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, payload)
			}))
			defer server.Close()
			dir := filepath.Join(t.TempDir(), "TestSensitive")
			s, err := New("record", dir, "TestSensitive")
			if err != nil {
				t.Fatal(err)
			}
			req, _ := http.NewRequest("GET", server.URL, nil)
			if _, err = s.Transport(nil).RoundTrip(req); err == nil {
				t.Fatal("sensitive response accepted")
			}
			if err = s.Close(true); err == nil {
				t.Fatal("sensitive response finalized")
			}
			if _, err = os.Stat(filepath.Join(dir, "cassette.json")); !os.IsNotExist(err) {
				t.Fatal("failed recording persisted cassette")
			}
		})
	}
}
func TestRecordingRejectsSensitiveBusinessQueryBeforeNetwork(t *testing.T) {
	s, err := New("record", filepath.Join(t.TempDir(), "case"), "TestSensitive")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "https://example.invalid/?Password=hidden", nil)
	if _, err = s.Transport(errorTransport{fmt.Errorf("network called")}).RoundTrip(req); err == nil || strings.Contains(err.Error(), "network called") {
		t.Fatalf("sensitive query reached network: %v", err)
	}
}
func TestIncompatibleRecordingAndImmutableBaseline(t *testing.T) {
	dir, _ := fixture(t)
	if _, err := New("record", dir, "TestOne"); err == nil {
		t.Fatal("overwrote baseline")
	}
	if _, err := New("replay", dir, "TestOther"); err == nil {
		t.Fatal("accepted another test")
	}
	data, _ := os.ReadFile(filepath.Join(dir, "cassette.json"))
	data = []byte(strings.Replace(string(data), `"formatVersion": 1`, `"formatVersion": 2`, 1))
	os.WriteFile(filepath.Join(dir, "cassette.json"), data, 0600)
	if _, err := New("replay", dir, "TestOne"); err == nil {
		t.Fatal("accepted incompatible version")
	}
}
func TestSeedDeterminismAndWait(t *testing.T) {
	dir, _ := fixture(t)
	a, _ := New("replay", dir, "TestOne")
	b, _ := New("replay", dir, "TestOne")
	if a.random.Int() != b.random.Int() {
		t.Fatal("seed not deterministic")
	}
	t.Setenv(ModeEnv, "replay")
	start := time.Now()
	if err := Wait(context.Background().Done(), time.Hour); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("replay delayed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Wait(ctx.Done(), time.Hour); err == nil {
		t.Fatal("ignored cancellation")
	}
}
func TestEnabledTransportRequiresMain(t *testing.T) {
	t.Setenv(ModeEnv, "replay")
	req, _ := http.NewRequest("GET", "https://example.invalid/", nil)
	if _, err := Transport(http.DefaultTransport).RoundTrip(req); err == nil {
		t.Fatal("uninitialized session used network")
	}
}

func TestEmbeddedSecretsAreRejected(t *testing.T) {
	s := &Session{}
	for _, body := range []string{`{"Config":"{\"password\":\"not-for-storage\"}"}`, `{"UserData":"password=not-for-storage"}`, `{"Nested":[{"TokenConfig":"security_token=hidden"}]}`} {
		if _, err := s.payload([]byte(body), "application/json", false); err == nil {
			t.Fatalf("accepted embedded secret %s", body)
		}
	}
}

func TestEmbeddedAccountIdentifierIsPseudonymized(t *testing.T) {
	s := &Session{}
	body := `{"Data":"{\"OwnerId\":\"1234567890123456\"}"}`
	clean, err := s.payload([]byte(body), "application/json", false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(clean, "1234567890123456") || !strings.Contains(clean, "recording-id-") {
		t.Fatalf("embedded identity was not sanitized: %s", clean)
	}
}
func TestBusinessDateQueryIsMatched(t *testing.T) {
	s := &Session{}
	first, _ := http.NewRequest("GET", "https://example.invalid/?Action=Read&Date=2026-01-01", nil)
	second, _ := http.NewRequest("GET", "https://example.invalid/?Action=Read&Date=2026-01-02", nil)
	a, err := s.normalize(first, nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.normalize(second, nil)
	if err != nil {
		t.Fatal(err)
	}
	if a.URL == b.URL {
		t.Fatal("ignored a business date query")
	}
}

func TestMalformedPayloadsFailClosed(t *testing.T) {
	s := &Session{}
	if _, err := s.payload([]byte(`{"Value":"one"} {"Password":"hidden"}`), "application/json", false); err == nil {
		t.Fatal("accepted trailing JSON")
	}
	req, _ := http.NewRequest("GET", "https://example.invalid/?Action=Read&Value=%zz", nil)
	if _, err := s.normalize(req, nil); err == nil {
		t.Fatal("silently discarded malformed query")
	}
}

func TestSelectedCaseFromExactRun(t *testing.T) {
	for _, selection := range []string{"^TestExample$", "^TestAccExample_123$"} {
		got, err := selectedCase(selection, "")
		if err != nil || got != selection[1:len(selection)-1] {
			t.Fatalf("selection %q: case=%q err=%v", selection, got, err)
		}
	}
	for _, selection := range []string{"", "TestExample", "^Test", "^Test.*$", "^(TestOne|TestTwo)$", "^TestOne$|^TestTwo$", "^TestOne$/subtest", "^TestOne/subtest$", "^BenchmarkOne$", "^TestOne$\n"} {
		if _, err := selectedCase(selection, ""); err == nil {
			t.Errorf("accepted ambiguous selector %q", selection)
		}
	}
	if _, err := selectedCase("^TestExample$", "TestOther"); err == nil {
		t.Fatal("accepted mismatched legacy case")
	}
	if got, err := selectedCase("^TestExample$", "TestExample"); err != nil || got != "TestExample" {
		t.Fatalf("legacy compatibility: %q %v", got, err)
	}
}

func TestGlobalAcceptanceRandomSourceReplays(t *testing.T) {
	t.Setenv("GODEBUG", "randseednop=1")
	generate := func() string {
		restore, err := seedGlobalRandom(987654321)
		if err != nil {
			t.Fatal(err)
		}
		defer restore()
		return fmt.Sprintf("%d/%d/%s/%s", mathrand.Int(), sdkacctest.RandIntRange(10, 100), sdkacctest.RandString(24), sdkacctest.RandomWithPrefix("test"))
	}
	first, second := generate(), generate()
	if first != second {
		t.Fatal("existing acceptance-library random calls did not replay")
	}
	if os.Getenv("GODEBUG") != "randseednop=1" {
		t.Fatal("GODEBUG was not restored")
	}
}

func TestClientTokensReplayWithoutDroppingBusinessParameters(t *testing.T) {
	generate := func() []string {
		active.Lock()
		active.session = &Session{random: mathrand.New(mathrand.NewSource(123))}
		active.Unlock()
		defer func() { active.Lock(); active.session = nil; active.Unlock() }()
		return []string{ClientToken("CreateVpc", nil), ClientToken("CreateVpc", nil), ClientToken(strings.Repeat("LongAction", 20), nil)}
	}
	first, second := generate(), generate()
	for i, token := range first {
		if token != second[i] || len(token) > 64 {
			t.Fatal("token did not replay or exceeded length limit")
		}
	}
	if first[0] == first[1] {
		t.Fatal("idempotency tokens were reused")
	}
	if ClientToken("CreateVpc", func() string { return "live" }) != "live" {
		t.Fatal("live generator changed")
	}
}
