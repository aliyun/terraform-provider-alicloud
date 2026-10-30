package alidns_test

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
	"strings"
	"sync"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sdkalidns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This is also the TestMain integration used after a product's tests move out of
// the root package. The recording identity is the exact test name, not its path.
func TestMain(m *testing.M) { acctest.RecordingTestMain(m) }

type recordingTestRun struct{ run func() int }

func (r recordingTestRun) Run() int { return r.run() }

// A fixture runs one explicit recording session inside an ordinary Go test. Its
// TestMain has already selected the Go tests; temporarily scope the recorder's
// single-case guard to this fixture and restore the caller's selection afterward.
// Live acceptance-test processes still use the unmodified TestMain guard.
func runRecordingFixture(t *testing.T, run func() int) int {
	t.Helper()
	selection := flag.Lookup("test.run").Value
	previous := selection.String()
	if err := selection.Set("^" + regexp.QuoteMeta(strings.SplitN(t.Name(), "/", 2)[0]) + "$"); err != nil {
		t.Fatal(err)
	}
	defer selection.Set(previous)
	return recording.WrapMain(recordingTestRun{run}).Run()
}

// TestRecordingDnsLifecycleOffline drives the registered provider resource and
// the real SDK against an HTTP fixture, then closes the server and repeats using
// the recorded responses. It does not provision cloud infrastructure.
func TestRecordingDnsLifecycleOffline(t *testing.T) {
	if recording.Enabled() {
		t.Fatal("this offline integration test manages its own recording sessions")
	}
	var mu sync.Mutex
	deleted := false
	firstCreateRetry := true
	value := "192.0.2.1"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		result := map[string]interface{}{"RequestId": "test"}
		action := r.Form.Get("Action")
		if action == "" {
			action = r.Header.Get("X-Acs-Action")
		}
		if action == "AddDomainRecord" && firstCreateRetry {
			firstCreateRetry = false
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"Code": "ServiceUnavailable", "Message": "synthetic transient error"})
			return
		}
		switch action {
		case "AddDomainRecord":
			deleted = false
			value = r.Form.Get("Value")
			result["RecordId"] = "synthetic-record"
		case "DescribeDomainRecordInfo":
			if deleted {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]string{"Code": "InvalidRR.NoExist", "Message": "synthetic record removed", "RequestId": "test"})
				return
			}
			result = map[string]interface{}{"RecordId": "synthetic-record", "DomainName": "example.invalid", "RR": "www", "Type": "A", "Value": value, "TTL": 600, "Line": "default", "Status": "ENABLE", "Locked": false}
		case "UpdateDomainRecord":
			value = r.Form.Get("Value")
			result["RecordId"] = "synthetic-record"
		case "DescribeRecordingProbe":
			result["State"] = "Ready"
		case "DeleteDomainRecord":
			deleted = true
			result["RecordId"] = "synthetic-record"
		default:
			t.Errorf("unexpected API action %s", r.Form.Get("Action"))
			w.WriteHeader(400)
		}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()
	endpoint, _ := url.Parse(server.URL)
	dir := filepath.Join(t.TempDir(), "TestRecordingDnsLifecycleOffline")
	t.Setenv(recording.PathEnv, dir)
	t.Setenv("ALICLOUD_ACCESS_KEY", "recording-access-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "recording-secret-key")
	t.Setenv("ALICLOUD_SECURITY_TOKEN", "")
	t.Setenv("HTTPS_PROXY", "")
	t.Setenv("HTTP_PROXY", "")
	t.Setenv("http_proxy", "")
	t.Setenv("https_proxy", "")
	t.Setenv("ALL_PROXY", "")
	t.Setenv("all_proxy", "")
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	t.Setenv("no_proxy", "127.0.0.1,localhost")
	run := func() int {
		credential, err := credentials.NewCredential(&credentials.Config{Type: ptr("access_key"), AccessKeyId: ptr("recording-access-key"), AccessKeySecret: ptr("recording-secret-key")})
		if err != nil {
			t.Error(err)
			return 1
		}
		endpoints := new(sync.Map)
		endpoints.Store("alidns", endpoint.Host)
		client, err := (&connectivity.Config{AccessKey: "recording-access-key", SecretKey: "recording-secret-key", Region: connectivity.Hangzhou, RegionId: string(connectivity.Hangzhou), AccountType: "Domestic", Protocol: "HTTP", Credential: credential, Endpoints: endpoints, SignVersion: new(sync.Map)}).Client()
		if err != nil {
			t.Error(err)
			return 1
		}
		// The legacy DNS client has its own built-in global endpoint map. Set
		// its explicit domain so this offline fixture cannot use that map.
		if _, err := client.WithDnsClient(func(c *sdkalidns.Client) (interface{}, error) { c.Domain = endpoint.Host; return nil, nil }); err != nil {
			t.Error(err)
			return 1
		}
		res := alicloud.Provider().ResourcesMap["alicloud_dns_record"]
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"name": "example.invalid", "host_record": "www", "type": "A", "value": "192.0.2.1", "ttl": 600, "routing": "default"})
		if err = res.Create(d, client); err != nil {
			t.Error(err)
			return 1
		}
		if d.Id() != "synthetic-record" {
			t.Errorf("unexpected resource identity %q", d.Id())
			return 1
		}
		if err = res.Read(d, client); err != nil {
			t.Error(err)
			return 1
		}
		d.Set("value", "192.0.2.2")
		if err = res.Update(d, client); err != nil {
			t.Error(err)
			return 1
		}
		if got := d.Get("value"); got != "192.0.2.2" {
			t.Errorf("read after update returned %v", got)
			return 1
		}
		imported, err := res.Importer.State(d, client)
		if err != nil || len(imported) != 1 || imported[0].Id() != d.Id() {
			t.Errorf("import failed: %v", err)
			return 1
		}
		if err = res.Delete(d, client); err != nil {
			t.Error(err)
			return 1
		}
		if err = res.Read(d, client); err != nil {
			t.Error(err)
			return 1
		}
		if d.Id() != "" {
			t.Error("read after delete did not remove state")
			return 1
		}
		// The shared Darabonba v2 adapter must use the same session as legacy SDK
		// calls, including Configure-created clients used by custom checks.
		params := &openapi.Params{Action: ptr("DescribeRecordingProbe"), Version: ptr("2020-01-01"), Protocol: ptr("HTTP"), Pathname: ptr("/"), Method: ptr("POST"), AuthType: ptr("AK"), Style: ptr("RPC"), ReqBodyType: ptr("formData"), BodyType: ptr("json")}
		probe, err := client.Do("alidns", params, nil, nil, nil, nil, false)
		if err != nil {
			t.Error(err)
			return 1
		}
		if probe["State"] != "Ready" {
			t.Errorf("unexpected v2 probe response %#v", probe)
			return 1
		}
		return 0
	}
	t.Setenv(recording.ModeEnv, "record")
	if code := runRecordingFixture(t, run); code != 0 {
		t.Fatal("record lifecycle failed")
	}
	server.Close()
	t.Setenv(recording.ModeEnv, "replay")
	if code := runRecordingFixture(t, run); code != 0 {
		t.Fatal("offline replay lifecycle failed")
	}
	data, err := os.ReadFile(filepath.Join(dir, "result.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"success": true`) {
		t.Fatalf("invalid replay report %s", data)
	}
	fmt.Printf("offline DNS provider lifecycle passed: %s\n", t.Name())
}
func ptr(s string) *string { return &s }
