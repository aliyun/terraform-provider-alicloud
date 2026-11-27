package cs

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitCSKubernetesAddonsErrorContext(t *testing.T) {
	for _, read := range []bool{false, true} {
		name := "describeAvailableAddons"
		if read {
			name = "dataSourceRead"
		}
		t.Run(name, func(t *testing.T) {
			var requests atomic.Int32
			client := csKubeconfigTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"addon query denied","RequestId":"unit-test"}`))
			})
			d := schema.TestResourceDataRaw(t, DataSourceCSKubernetesAddons().Schema, map[string]interface{}{
				"cluster_id": "unit-test-cluster",
			})
			var err error
			if read {
				err = dataCSKubernetesAddonsRead(d, client)
			} else {
				_, err = describeAvailableAddons(d, client)
			}
			if err == nil {
				t.Fatal("expected the addon query error to be returned")
			}
			if requests.Load() != 1 {
				t.Fatalf("got %d requests, want one non-retryable addon query", requests.Load())
			}
			if !errs.IsExpectedErrors(err, []string{"InvalidParameter"}) {
				t.Errorf("original API error was not preserved: %s", err)
			}
			message := err.Error()
			if strings.Contains(message, "%!") {
				t.Errorf("malformed format argument in error: %s", message)
			}
			for _, want := range []string{
				"Resource resourceAlicloudCSKubernetesAddon DescribeCsKubernetesAllAvailableAddons Failed!!!",
				"addon query denied",
			} {
				if !strings.Contains(message, want) {
					t.Errorf("error missing %q: %s", want, message)
				}
			}
			if read && !strings.Contains(message, "Resource resourceAlicloudCSKubernetesAddon describeAddonsMeta Failed!!!") {
				t.Errorf("read error missing resource context: %s", message)
			}
			if d.Id() != "" {
				t.Errorf("failed query assigned an ID: %q", d.Id())
			}
		})
	}
}

// TestUnitCSKubernetesAddonInstallClusterUpdatingRetry covers the retry added
// for InvalidClusterStatus.Forbidden: while the cluster stays in the updating
// state the install request is repeated and the install succeeds once the API
// accepts it, while unrelated 400 codes stay non-retryable. The retry
// condition is shared by install/upgrade/uninstall/config update.
func TestUnitCSKubernetesAddonInstallClusterUpdatingRetry(t *testing.T) {
	for _, tc := range []struct {
		name         string
		codes        []string
		wantRequests int
		wantErr      string
	}{
		{"cluster_updating_retried", []string{"InvalidClusterStatus.Forbidden", "InvalidClusterStatus.Forbidden"}, 3, ""},
		{"unrelated_error_not_retried", []string{"InvalidParameter"}, 1, "InvalidParameter"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var installRequests atomic.Int32
			client := csKubeconfigTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/components/install"):
					if i := int(installRequests.Add(1)) - 1; i < len(tc.codes) {
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write([]byte(fmt.Sprintf(`{"code":%q,"message":"cannot operate cluster where state is updating","RequestId":"unit-test"}`, tc.codes[i])))
						return
					}
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"task_id":"T-unit-test","cluster_id":"c1","request_id":"unit-test"}`))
				case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/tasks/"):
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"task_id":"T-unit-test","state":"success","request_id":"unit-test"}`))
				default:
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"code":"InvalidParameter","message":"unexpected request"}`))
				}
			})
			roaClient, err := client.NewRoaCsClient()
			if err != nil {
				t.Fatal(err)
			}
			d := schema.TestResourceDataRaw(t, ResourceCSKubernetesAddon().Schema, map[string]interface{}{
				"cluster_id": "c1",
				"name":       "logtail-ds",
			})
			err = NewCsClient(roaClient).installAddon(d)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("installAddon returned error: %s", err)
				}
			} else if !errs.IsExpectedErrors(err, []string{tc.wantErr}) {
				t.Fatalf("installAddon error = %v, want terminal %q", err, tc.wantErr)
			}
			if got := int(installRequests.Load()); got != tc.wantRequests {
				t.Errorf("InstallClusterAddons called %d times, want %d", got, tc.wantRequests)
			}
		})
	}
}
