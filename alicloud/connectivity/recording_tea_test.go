package connectivity

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
)

type recordingSDKRunner struct{ run func() int }

func (r recordingSDKRunner) Run() int { return r.run() }

// Real SDK signing, flattening, retries and SDK error decoding remain active in
// replay. Closing the HTTP server proves both legacy paths consume the tape.
func TestRecordingLegacyTeaHTTPPaths(t *testing.T) {
	var calls atomic.Int32
	var rpcAttempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		t.Logf("fixture request: method=%s path=%s action=%s", r.Method, r.URL.Path, r.URL.Query().Get("Action"))
		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		action := r.Form.Get("Action")
		if action == "ProbeRPC" && rpcAttempts.Add(1) == 1 {
			w.WriteHeader(503)
			fmt.Fprint(w, `{"Code":"ServiceUnavailable","Message":"fixture retry","RequestId":"test"}`)
			return
		}
		if r.URL.Path == "/empty" {
			w.WriteHeader(204)
			return
		}
		if action == "FailRPC" || r.URL.Path == "/missing" {
			w.WriteHeader(404)
			fmt.Fprint(w, `{"Code":"InvalidResource.NotFound","Message":"missing fixture","RequestId":"test"}`)
			return
		}
		if action == "ProbeRPC" && r.Form.Get("Value") != "business-value" {
			t.Error("RPC business body missing")
		}
		if r.URL.Path == "/action" && r.Header.Get("x-acs-action") != "ProbeAction" {
			t.Error("ROA action header missing")
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"Value": "ok", "RequestId": "test"})
	}))
	defer server.Close()
	endpoint, _ := url.Parse(server.URL)
	for _, key := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy"} {
		t.Setenv(key, "")
	}
	t.Setenv(recording.PathEnv, filepath.Join(t.TempDir(), "case"))
	t.Setenv(recording.CaseEnv, "")
	selection := flag.Lookup("test.run").Value
	previous := selection.String()
	selection.Set("^" + regexp.QuoteMeta(t.Name()) + "$")
	defer selection.Set(previous)
	run := func() int {
		cred, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("recording-access-key").SetAccessKeySecret("recording-secret-key"))
		if err != nil {
			t.Error(err)
			return 1
		}
		endpoints := new(sync.Map)
		endpoints.Store("fixture", endpoint.Host)
		endpoints.Store("bssopenapi", endpoint.Host)
		c, err := (&Config{AccessKey: "recording-access-key", SecretKey: "recording-secret-key", Region: Hangzhou, RegionId: string(Hangzhou), Protocol: "HTTP", Credential: cred, Endpoints: endpoints, SignVersion: new(sync.Map)}).Client()
		if err != nil {
			t.Error(err)
			return 1
		}
		if c.config.AccountType != "Domestic" {
			t.Error("account discovery not preserved")
			return 1
		}
		if _, err = c.RpcPost("fixture", "2020-01-01", "ProbeRPC", nil, map[string]interface{}{"Value": "business-value", "ClientToken": "stable-token"}, true); err != nil {
			t.Error(err)
			return 1
		}
		if _, err = c.RoaPostWithApiName("fixture", "2020-01-01", "ProbeAction", "/action", nil, nil, map[string]interface{}{"Value": "business-value"}, true); err != nil {
			t.Error(err)
			return 1
		}
		if _, err = c.RoaGet("fixture", "2020-01-01", "/probe", nil, nil, nil); err != nil {
			t.Error(err)
			return 1
		}
		if response, err := c.RoaDelete("fixture", "2020-01-01", "/empty", nil, nil, nil, false); err != nil || response["headers"] == nil {
			t.Errorf("ROA 204 semantics changed: %v", err)
			return 1
		}
		for _, invoke := range []func() error{
			func() error { _, err := c.RpcPost("fixture", "2020-01-01", "FailRPC", nil, nil, false); return err },
			func() error { _, err := c.RoaGet("fixture", "2020-01-01", "/missing", nil, nil, nil); return err },
		} {
			sdkError, ok := invoke().(*tea.SDKError)
			if !ok || tea.StringValue(sdkError.Code) != "InvalidResource.NotFound" || tea.IntValue(sdkError.StatusCode) != 404 {
				t.Errorf("SDK error semantics changed: %v", sdkError)
				return 1
			}
		}
		return 0
	}
	for _, mode := range []string{"live", "record", "replay"} {
		if mode == "replay" {
			server.Close()
		}
		rpcAttempts.Store(0)
		t.Setenv(recording.ModeEnv, mode)
		if code := recording.WrapMain(recordingSDKRunner{run}).Run(); code != 0 {
			t.Fatalf("%s legacy SDK paths failed", mode)
		}
	}
	// Each live phase includes discovery, a 503 retry, ROA action/204 and SDK errors.
	if calls.Load() != 16 {
		t.Fatalf("unexpected real HTTP calls: %d", calls.Load())
	}
}

func TestRecordingLegacyTeaBodyTimeoutFailsClosed(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"Value":`)
		w.(http.Flusher).Flush()
		close(started)
		<-r.Context().Done()
	}))
	defer server.Close()
	endpoint, _ := url.Parse(server.URL)
	for _, key := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy"} {
		t.Setenv(key, "")
	}
	directory := filepath.Join(t.TempDir(), "case")
	t.Setenv(recording.ModeEnv, "record")
	t.Setenv(recording.PathEnv, directory)
	t.Setenv(recording.CaseEnv, "")
	selection := flag.Lookup("test.run").Value
	previous := selection.String()
	selection.Set("^" + regexp.QuoteMeta(t.Name()) + "$")
	defer selection.Set(previous)
	run := func() int {
		cred, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("recording-access-key").SetAccessKeySecret("recording-secret-key"))
		if err != nil {
			t.Error(err)
			return 1
		}
		endpoints := new(sync.Map)
		endpoints.Store("fixture", endpoint.Host)
		c, err := (&Config{AccessKey: "recording-access-key", SecretKey: "recording-secret-key", AccountType: "Domestic", Region: Hangzhou, RegionId: string(Hangzhou), Protocol: "HTTP", Credential: cred, Endpoints: endpoints, SignVersion: new(sync.Map), ClientReadTimeout: 100, ClientConnectTimeout: 100}).Client()
		if err != nil {
			t.Error(err)
			return 1
		}
		if _, err = c.RpcPost("fixture", "2020-01-01", "Timeout", nil, nil, false); err == nil {
			t.Error("incomplete response body did not time out")
			return 1
		}
		select {
		case <-started:
		default:
			t.Error("server never started the response body")
			return 1
		}
		return 0 // A swallowed SDK timeout must still invalidate the session.
	}
	if recording.WrapMain(recordingSDKRunner{run}).Run() == 0 {
		t.Fatal("failed transport produced a successful recording")
	}
	if _, err := os.Stat(filepath.Join(directory, "cassette.json")); !os.IsNotExist(err) {
		t.Fatal("failed timeout recording wrote a cassette")
	}
}
