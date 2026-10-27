package alidns_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	sdkalidns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestRecordingDnsTerraformCLI uses Terraform's acceptance harness, not direct
// handler calls. Only Configure is adapted to point real SDK clients at a local
// fixture with explicit dummy credentials. No cloud account is used.
func TestRecordingDnsTerraformCLI(t *testing.T) {
	if recording.Enabled() {
		t.Fatal("this offline test manages its recording sessions")
	}
	terraform, err := exec.LookPath("terraform")
	if err != nil {
		t.Skip("offline Terraform CLI test requires an installed terraform executable")
	}
	terraform, err = filepath.Abs(terraform)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("TF_ACC_TERRAFORM_PATH", terraform)
	t.Setenv("CHECKPOINT_DISABLE", "1")
	for _, key := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy", "TF_LOG", "TF_LOG_PATH"} {
		t.Setenv(key, "")
	}
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	t.Setenv("no_proxy", "127.0.0.1,localhost")
	t.Setenv("ALICLOUD_ACCESS_KEY", "recording-access-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "recording-secret-key")
	t.Setenv("ALICLOUD_SECURITY_TOKEN", "")
	t.Setenv("ALIBABA_CLOUD_ACCOUNT_TYPE", "Domestic")

	var mu sync.Mutex
	var calls atomic.Int32
	deleted, value := true, "192.0.2.1"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		calls.Add(1)
		if err := r.ParseForm(); err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		result := map[string]interface{}{"RequestId": "test", "RecordId": "fixture-record"}
		switch r.Form.Get("Action") {
		case "AddDomainRecord":
			deleted = false
			value = r.Form.Get("Value")
		case "UpdateDomainRecord":
			value = r.Form.Get("Value")
		case "DescribeDomainRecordInfo":
			if deleted {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]string{"Code": "InvalidRR.NoExist", "Message": "fixture record removed", "RequestId": "test"})
				return
			}
			result = map[string]interface{}{"RecordId": "fixture-record", "DomainName": "example.invalid", "RR": "www", "Type": "A", "Value": value, "TTL": 600, "Priority": 0, "Line": "default", "Status": "ENABLE", "Locked": false}
		case "DeleteDomainRecord":
			deleted = true
		default:
			t.Errorf("unexpected fixture API action %s", r.Form.Get("Action"))
			w.WriteHeader(400)
		}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()
	endpoint, _ := url.Parse(server.URL)
	directory := filepath.Join(t.TempDir(), t.Name())
	if output := os.Getenv("TF_RECORDING_TEST_OUTPUT"); output != "" {
		if err := os.MkdirAll(output, 0700); err != nil {
			t.Fatal(err)
		}
		directory = filepath.Join(output, t.Name())
	}
	t.Setenv(recording.PathEnv, directory)
	t.Setenv(recording.CaseEnv, t.Name())

	factory := func() (tfprotov5.ProviderServer, error) {
		p := alicloud.Provider()
		// Limit schema validation to the actual registered pilot resource. Resource
		// implementation, SDK calls, state operations and Terraform RPC remain real.
		p.ResourcesMap = map[string]*schema.Resource{"alicloud_dns_record": p.ResourcesMap["alicloud_dns_record"]}
		p.DataSourcesMap = nil
		p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			credential, err := credentials.NewCredential(&credentials.Config{Type: ptr("access_key"), AccessKeyId: ptr("recording-access-key"), AccessKeySecret: ptr("recording-secret-key")})
			if err != nil {
				return nil, diag.FromErr(err)
			}
			endpoints := new(sync.Map)
			endpoints.Store("alidns", endpoint.Host)
			client, err := (&connectivity.Config{AccessKey: "recording-access-key", SecretKey: "recording-secret-key", Region: connectivity.Hangzhou, RegionId: string(connectivity.Hangzhou), AccountType: "Domestic", Protocol: "HTTP", Credential: credential, Endpoints: endpoints, SignVersion: new(sync.Map)}).Client()
			if err != nil {
				return nil, diag.FromErr(err)
			}
			_, err = client.WithDnsClient(func(c *sdkalidns.Client) (interface{}, error) { c.Domain = endpoint.Host; return nil, nil })
			if err != nil {
				return nil, diag.FromErr(err)
			}
			return client, nil
		}
		return schema.NewGRPCProviderServer(p), nil
	}
	config := func(value string) string {
		return fmt.Sprintf(`
provider "alicloud" { region = "cn-hangzhou" }
resource "alicloud_dns_record" "test" {
 name = "example.invalid"
 host_record = "www"
 type = "A"
 value = %q
 ttl = 600
 routing = "default"
}
`, value)
	}
	run := func(phase *testing.T) int {
		passed := phase.Run("terraform", func(lifecycle *testing.T) {
			resource.Test(lifecycle, resource.TestCase{
				IsUnitTest:               true,
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){"alicloud": factory},
				Steps: []resource.TestStep{
					{Config: config("192.0.2.1"), Check: resource.TestCheckResourceAttr("alicloud_dns_record.test", "value", "192.0.2.1")},
					{ResourceName: "alicloud_dns_record.test", ImportState: true, ImportStateVerify: true},
					{Config: config("192.0.2.2"), Check: resource.TestCheckResourceAttr("alicloud_dns_record.test", "value", "192.0.2.2")},
				},
			})
		})
		if !passed {
			return 1
		}
		return 0
	}
	baselineFiles := make(map[string][]byte)
	for _, mode := range []string{"record", "replay"} {
		if mode == "replay" {
			server.Close()
		}
		t.Setenv(recording.ModeEnv, mode)
		t.Setenv("TF_ACC_RECORDING_RESULT_PATH", directory+"-"+mode+".json")
		t.Run(mode, func(phase *testing.T) {
			if runRecordingFixture(phase, func() int { return run(phase) }) != 0 {
				phase.Fatal("Terraform recording session failed")
			}
		})
		if t.Failed() {
			return
		}
		for _, name := range []string{"cassette.json", "seed"} {
			content, err := os.ReadFile(filepath.Join(directory, name))
			if err != nil {
				t.Fatal(err)
			}
			if mode == "record" {
				baselineFiles[name] = content
			} else if !bytes.Equal(baselineFiles[name], content) {
				t.Fatalf("replay modified immutable baseline file %s", name)
			}
		}
	}
	var result recording.Result
	data, err := os.ReadFile(directory + "-replay.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.UnmatchedRequests != 0 || result.UnusedInteractions != 0 || result.InteractionCount == 0 {
		t.Fatalf("incomplete replay: %+v", result)
	}
	if int(calls.Load()) != result.InteractionCount {
		t.Fatalf("HTTP fixture was called during replay: calls=%d interactions=%d", calls.Load(), result.InteractionCount)
	}
	t.Logf("Terraform CLI record/replay: interactions=%d, replay network calls=0, artifacts=%s", result.InteractionCount, directory)
}
