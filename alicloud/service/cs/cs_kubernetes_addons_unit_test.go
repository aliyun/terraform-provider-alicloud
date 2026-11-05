package cs

import (
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
