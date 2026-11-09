package polardb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitPolarDBEndpointSccModeRequest(t *testing.T) {
	for _, tc := range []struct {
		name       string
		sccMode    string
		configured bool
	}{
		{name: "configured_on", sccMode: "on", configured: true},
		{name: "configured_off", sccMode: "off", configured: true},
		{name: "unconfigured"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var created atomic.Bool
			var mu sync.Mutex
			var createForm, modifyForm map[string][]string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "DescribeDBClusterEndpoints":
					items := []map[string]string{{"DBEndpointId": "pe-primary", "EndpointType": "Primary"}}
					if created.Load() {
						items = append(items, map[string]string{
							"DBEndpointId":   "pe-new",
							"EndpointType":   "Custom",
							"ReadWriteMode":  "ReadWrite",
							"SccMode":        tc.sccMode,
							"EndpointConfig": `{"ConsistLevel":"0"}`,
						})
					}
					_ = json.NewEncoder(w).Encode(map[string]interface{}{"Items": items})
				case "CreateDBClusterEndpoint":
					mu.Lock()
					createForm = copyForm(r.Form)
					mu.Unlock()
					created.Store(true)
					_ = json.NewEncoder(w).Encode(map[string]string{"DBEndpointId": "pe-new"})
				case "ModifyDBClusterEndpoint":
					mu.Lock()
					modifyForm = copyForm(r.Form)
					mu.Unlock()
					_, _ = w.Write([]byte(`{}`))
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			client := unitTestPolarDBClient(t, server.URL)
			raw := map[string]interface{}{
				"db_cluster_id":   "pc-unit-test",
				"endpoint_type":   "Custom",
				"read_write_mode": "ReadWrite",
				"endpoint_config": map[string]interface{}{"ConsistLevel": "0"},
			}
			if tc.configured {
				raw["scc_mode"] = tc.sccMode
			}
			d := schema.TestResourceDataRaw(t, ResourceEndpoint().Schema, raw)
			d.MarkNewResource()
			if err := resourceEndpointCreate(d, client); err != nil {
				t.Fatalf("create failed: %s", err)
			}

			mu.Lock()
			defer mu.Unlock()
			if createForm == nil {
				t.Fatal("CreateDBClusterEndpoint was not called")
			}
			assertSccModeParam(t, "create", createForm, tc.sccMode, tc.configured)
			if modifyForm != nil {
				assertSccModeParam(t, "modify", modifyForm, tc.sccMode, tc.configured)
			}
		})
	}
}

func assertSccModeParam(t *testing.T, action string, form map[string][]string, want string, configured bool) {
	t.Helper()
	got, present := form["SccMode"]
	if !configured {
		if present {
			t.Fatalf("%s sent SccMode %v without configuration", action, got)
		}
		return
	}
	if !present {
		t.Fatalf("%s did not send SccMode, want %q", action, want)
	}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("%s SccMode = %v, want %q", action, got, want)
	}
}

func copyForm(form map[string][]string) map[string][]string {
	result := make(map[string][]string, len(form))
	for key, values := range form {
		result[key] = append([]string(nil), values...)
	}
	return result
}

func unitTestPolarDBClient(t *testing.T, serverURL string) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(serverURL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}
