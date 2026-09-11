package alicloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sdkendpoints "github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func drdsInstanceReadTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	endpoint := sdkendpoints.GetEndpointFromMap("cn-hangzhou", "drds")
	t.Cleanup(func() { sdkendpoints.AddEndpointMapping("cn-hangzhou", "drds", endpoint) })
	for _, variable := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy", "TF_ENDPOINT_PATH"} {
		t.Setenv(variable, "")
	}
	t.Setenv("NO_PROXY", "*")
	t.Setenv("no_proxy", "*")
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	var endpoints, signVersion sync.Map
	endpoints.Store("drds", strings.TrimPrefix(server.URL, "http://"))
	config := &connectivity.Config{
		Region: connectivity.Hangzhou, RegionId: "cn-hangzhou", Protocol: "HTTP",
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		Endpoints: &endpoints, SignVersion: &signVersion,
		AccountType: "test", SkipRegionValidation: true,
		ClientReadTimeout: 1000, ClientConnectTimeout: 1000,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestUnitDrdsInstanceReadConsistency(t *testing.T) {
	const complete = `{"Data":{"Status":"RUN","ZoneId":"cn-hangzhou-j","Description":"test instance","InstanceSeries":"drds.sn1.4c8g","InstanceSpec":"drds.sn1.4c8g.8C16G","CommodityCode":"drdsPost","Vips":{"Vip":[{"Type":"intranet","VpcId":"vpc-test","VswitchId":"vsw-test","Dns":"db.example.com","Port":"3306"}]}}}`
	const partialSpec = `{"Data":{"Status":"RUN","ZoneId":"cn-hangzhou-j","Description":"test instance","InstanceSpec":"8C16G","Vips":{"Vip":[{"Type":"intranet","VpcId":"vpc-test","VswitchId":"vsw-test"}]}}}`
	const partialNetwork = `{"Data":{"Status":"RUN","ZoneId":"cn-hangzhou-j","Description":"test instance","InstanceSeries":"drds.sn1.4c8g","InstanceSpec":"drds.sn1.4c8g.8C16G"}}`
	for _, tc := range []struct {
		name              string
		first             string
		persistent        bool
		existing          bool
		initiallyComplete bool
	}{
		{name: "import waits for complete specification", first: partialSpec},
		{name: "specification must include series", first: strings.Replace(partialSpec, `"InstanceSpec"`, `"InstanceSeries":"drds.sn1.4c8g","InstanceSpec"`, 1)},
		{name: "import waits for network", first: partialNetwork},
		{name: "blank description is valid", first: strings.Replace(complete, `"Description":"test instance"`, `"Description":""`, 1), initiallyComplete: true},
		{name: "public VIP precedes network VIP", first: strings.Replace(complete, `"Vip":[`, `"Vip":[{"Type":"internet","Dns":"public.example.com","Port":"3306"},`, 1), initiallyComplete: true},
		{name: "lowercase dns is supported", first: strings.Replace(complete, `"Dns"`, `"dns"`, 1), initiallyComplete: true},
		{name: "numeric port is supported", first: strings.Replace(complete, `"Port":"3306"`, `"Port":3306`, 1), initiallyComplete: true},
		{name: "incomplete import fails", first: partialSpec, persistent: true},
		{name: "empty response preserves existing state", first: `{"Data":{"Status":"RUN"}}`, persistent: true, existing: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls int32
			client := drdsInstanceReadTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				if r.Form.Get("Action") != "DescribeDrdsInstance" {
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
				w.Header().Set("Content-Type", "application/json")
				response := complete
				if atomic.AddInt32(&calls, 1) == 1 || tc.persistent {
					response = tc.first
				}
				fmt.Fprint(w, response)
			})
			attributes := map[string]string{}
			if tc.existing {
				attributes = map[string]string{
					"specification": "drds.sn1.4c8g.8C16G", "instance_series": "drds.sn1.4c8g",
					"vpc_id": "vpc-test", "vswitch_id": "vsw-test", "zone_id": "cn-hangzhou-j",
					"description": "test instance", "connection_string": "db.example.com", "port": "3306",
				}
			}
			r := resourceAlicloudDRDSInstance()
			r.Timeouts.Read = schema.DefaultTimeout(time.Second)
			d := r.Data(&terraform.InstanceState{ID: "drds-test", Attributes: attributes})
			err := resourceAliCloudDRDSInstanceRead(d, client)
			if tc.persistent {
				if err == nil {
					t.Error("an incomplete response must not produce a successful Read")
				}
				for key, want := range attributes {
					if got := d.Get(key); got != want {
						t.Errorf("%s changed to %q after incomplete response, want %q", key, got, want)
					}
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				for key, want := range map[string]string{
					"specification": "drds.sn1.4c8g.8C16G", "instance_series": "drds.sn1.4c8g",
					"vswitch_id": "vsw-test", "vpc_id": "vpc-test", "status": "RUN", "instance_charge_type": "PostPaid",
					"connection_string": "db.example.com", "port": "3306",
				} {
					if got := d.Get(key); got != want {
						t.Errorf("%s = %q, want %q", key, got, want)
					}
				}
				if tc.initiallyComplete {
					if atomic.LoadInt32(&calls) != 1 {
						t.Error("a complete response must be read without retry")
					}
				} else if atomic.LoadInt32(&calls) < 2 {
					t.Error("Read returned before the complete response")
				}
			}
			if d.Id() != "drds-test" {
				t.Error("incomplete response cleared the instance ID")
			}
		})
	}
}

func TestUnitDrdsInstanceReadTransientError(t *testing.T) {
	for _, code := range []string{"InternalError", "Throttling", "InvalidDrdsInstanceId.NotFound", "Forbidden"} {
		t.Run(code, func(t *testing.T) {
			var calls int32
			client := drdsInstanceReadTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if atomic.AddInt32(&calls, 1) == 1 {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, `{"Code":%q,"Message":"test response","RequestId":"test"}`, code)
					return
				}
				fmt.Fprint(w, `{"Data":{"Status":"RUN"}}`)
			})
			service := DrdsService{client}
			object, state, err := service.DrdsInstanceStateRefreshFunc("drds-test", nil)()
			switch code {
			case "InternalError", "Throttling":
				if err != nil || object == nil || state != "RUN" || atomic.LoadInt32(&calls) < 2 {
					t.Errorf("transient error was treated as absence: object=%v state=%q err=%v calls=%d", object, state, err, atomic.LoadInt32(&calls))
				}
			case "InvalidDrdsInstanceId.NotFound":
				if err != nil || object != nil || state != "" {
					t.Errorf("NotFound must remain the delete target: object=%v state=%q err=%v", object, state, err)
				}
			case "Forbidden":
				if err == nil {
					t.Error("terminal error must not be treated as absence")
				}
			}
		})
	}
}
