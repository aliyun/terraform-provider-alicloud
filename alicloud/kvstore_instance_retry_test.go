package alicloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	sdkendpoints "github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestUnitKvstoreCreateLockRetry(t *testing.T) {
	vpcEndpoint := sdkendpoints.GetEndpointFromMap("cn-hangzhou", "vpc")
	t.Cleanup(func() { sdkendpoints.AddEndpointMapping("cn-hangzhou", "vpc", vpcEndpoint) })
	for _, variable := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy", "TF_ENDPOINT_PATH"} {
		t.Setenv(variable, "")
	}
	t.Setenv("NO_PROXY", "*")
	t.Setenv("no_proxy", "*")
	for _, test := range []struct {
		name      string
		firstCode string
		calls     int
		zoneField string
		zone      string
		vswitch   bool
	}{
		{name: "lock conflict", firstCode: "CanNotAcquireLock", calls: 2},
		{name: "invalid parameter", firstCode: "InvalidParameter", calls: 1},
		{name: "VSwitch zone fallback", firstCode: "InvalidParameter", calls: 1, zone: "cn-hangzhou-j", vswitch: true},
		{name: "explicit zone", firstCode: "InvalidParameter", calls: 1, zoneField: "zone_id", zone: "cn-hangzhou-i", vswitch: true},
		{name: "legacy zone", firstCode: "InvalidParameter", calls: 1, zoneField: "availability_zone", zone: "cn-hangzhou-i", vswitch: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var tokens []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				if r.Form.Get("Action") == "DescribeVSwitchAttributes" {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, `{"VSwitchId":"vsw-test","VpcId":"vpc-test","ZoneId":"cn-hangzhou-j"}`)
					return
				}
				if r.Form.Get("Action") != "CreateInstance" {
					t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				}
				if got := r.Form.Get("ZoneId"); got != test.zone {
					t.Errorf("CreateInstance ZoneId = %q, want %q", got, test.zone)
				}
				tokens = append(tokens, r.Form.Get("Token"))
				t.Logf("received create request %d", len(tokens))
				code := "InvalidParameter"
				if len(tokens) == 1 {
					code = test.firstCode
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"Code":%q,"Message":"test failure","RequestId":"test"}`, code)
			}))
			defer server.Close()

			credential, err := credentials.NewCredential(new(credentials.Config).
				SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
			if err != nil {
				t.Fatal(err)
			}
			var endpoints, signVersion sync.Map
			endpoints.Store("r_kvstore", strings.TrimPrefix(server.URL, "http://"))
			endpoints.Store("vpc", strings.TrimPrefix(server.URL, "http://"))
			config := &connectivity.Config{
				Region: connectivity.Hangzhou, RegionId: "cn-hangzhou", Protocol: "HTTP",
				AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
				Endpoints: &endpoints, SignVersion: &signVersion,
				AccountType: "test", SkipRegionValidation: true,
				MaxRetryTimeout: 15, ClientReadTimeout: 1000, ClientConnectTimeout: 1000,
			}
			client, err := config.Client()
			if err != nil {
				t.Fatal(err)
			}
			d := schema.TestResourceDataRaw(t, resourceAliCloudKvstoreInstance().Schema, map[string]interface{}{
				"instance_class": "redis.shard.with.proxy.small.ce",
				"engine_version": "6.0",
			})
			if test.vswitch {
				d.Set("vswitch_id", "vsw-test")
			}
			if test.zoneField == "zone_id" {
				d.Set("zone_id", test.zone)
			} else if test.zoneField == "availability_zone" {
				d.Set("availability_zone", test.zone)
			}
			// A terminal response after the conflict exercises the create retry loop
			// without creating an instance or entering the instance-state waiter.
			err = resourceAliCloudKvstoreInstanceCreate(d, client)
			if err == nil || !IsExpectedErrors(err, []string{"InvalidParameter"}) {
				t.Fatalf("expected the terminal API error, got %v", err)
			}
			if len(tokens) != test.calls {
				t.Fatalf("sent %d requests, want %d", len(tokens), test.calls)
			}
			for _, token := range tokens {
				if token == "" || token != tokens[0] {
					t.Fatal("create retries must reuse the nonempty idempotency token")
				}
			}
			if d.Id() != "" {
				t.Fatal("failed creation must not set an instance ID")
			}
		})
	}
}
