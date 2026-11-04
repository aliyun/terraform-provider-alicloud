package polardb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitPolarDBEndpointCreateRetryIsolation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		codes []string
	}{
		{"endpoint_status", []string{"ClusterEndpoint.StatusNotValid", "InvalidParameter"}},
		{"shared_status", []string{"OperationDenied.DBStatus", "InvalidParameter"}},
		{"repeated_endpoint_status", []string{"ClusterEndpoint.StatusNotValid", "ClusterEndpoint.StatusNotValid", "InvalidParameter"}},
		{"non_retryable", []string{"InvalidParameter"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := append([]string(nil), errs.OperationDeniedDBStatus...)
			// Restore the original slice even when running this regression against
			// an implementation that mutates the global list.
			original := errs.OperationDeniedDBStatus
			t.Cleanup(func() { errs.OperationDeniedDBStatus = original })
			var createCalls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "DescribeDBClusterEndpoints":
					_, _ = w.Write([]byte(`{"Items":[{"DBEndpointId":"pe-primary","EndpointType":"Primary"}]}`))
				case "CreateDBClusterEndpoint":
					if r.Form.Get("DBClusterId") != "pc-unit-test" || r.Form.Get("EndpointType") != "Custom" {
						t.Errorf("unexpected create parameters: %v", r.Form)
					}
					i := int(createCalls.Add(1)) - 1
					code := "InvalidParameter"
					if i < len(tc.codes) {
						code = tc.codes[i]
					} else {
						t.Errorf("unexpected extra CreateDBClusterEndpoint request")
					}
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(map[string]string{
						"Code": code, "Message": "endpoint create fixture", "RequestId": "unit-test",
					})
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)
			credential, err := credentials.NewCredential(new(credentials.Config).
				SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
			if err != nil {
				t.Fatal(err)
			}
			endpoint := strings.TrimPrefix(server.URL, "http://")
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
			d := schema.TestResourceDataRaw(t, ResourceEndpoint().Schema, map[string]interface{}{
				"db_cluster_id": "pc-unit-test",
			})
			d.MarkNewResource()
			err = resourceEndpointCreate(d, client)
			if !errs.IsExpectedErrors(err, []string{"InvalidParameter"}) {
				t.Errorf("expected terminal API error, got %v", err)
			}
			if got := int(createCalls.Load()); got != len(tc.codes) {
				t.Errorf("CreateDBClusterEndpoint called %d times, want %d", got, len(tc.codes))
			}
			if !reflect.DeepEqual(errs.OperationDeniedDBStatus, before) {
				t.Errorf("create changed shared retry codes: got %v, want %v", errs.OperationDeniedDBStatus, before)
			}
			if d.Id() != "" {
				t.Errorf("failed create assigned an ID: %q", d.Id())
			}
		})
	}
}
