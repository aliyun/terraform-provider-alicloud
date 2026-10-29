package dts_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dtsSynchronizationJobImportTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{AccessKey: "test-key", SecretKey: "test-secret", Credential: credential, RegionId: "eu-west-1", AccountType: "test", Protocol: "http", Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("dts", endpoint)
	return client
}

func TestUnitDtsJobImportIdentity(t *testing.T) {
	for kind, factory := range map[string]func() *schema.Resource{"sync": dts.ResourceSynchronizationJob, "subscription": dts.ResourceSubscriptionJob} {
		t.Run(kind, func(t *testing.T) {
			testDtsJobImportIdentity(t, factory)
		})
	}
}

func testDtsJobImportIdentity(t *testing.T, factory func() *schema.Resource) {
	for _, tc := range []struct {
		name, input, returned, status, apiCode, wantError string
	}{
		{name: "canonical_job", input: "njl106g62735gmd", returned: "njl106g62735gmd"},
		{name: "rebuilt_job", input: "is210es227j5gr6", returned: "is210es227j5gr6"},
		{name: "not_prefix_based", input: "dts-job-id", returned: "dts-job-id"},
		{name: "instance_id", input: "dtsnjl106g62735gmd", returned: "njl106g62735gmd", wantError: "does not match"},
		{name: "instance_with_deleted_history", input: "dtsnjl106g62735gmd", returned: "is210es227j5gr6", wantError: "does not match"},
		{name: "different_job", input: "njl106g62735gmd", returned: "is210es227j5gr6", wantError: "does not match"},
		{name: "missing_job_id", input: "njl106g62735gmd", wantError: "DtsJobId"},
		{name: "deleted_job", input: "njl106g62735gmd", returned: "njl106g62735gmd", status: "Finished", wantError: "cannot import"},
		{name: "not_found", input: "njl106g62735gmd", apiCode: "Forbidden.InstanceNotFound", wantError: "cannot import"},
		{name: "permission_denied", input: "njl106g62735gmd", apiCode: "Forbidden.NoPermission", wantError: "Forbidden.NoPermission"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			client := dtsSynchronizationJobImportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				requests++
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				if r.Form.Get("Action") != "DescribeDtsJobDetail" || r.Form.Get("DtsJobId") != tc.input {
					t.Errorf("unexpected request: %v", r.Form)
				}
				w.Header().Set("Content-Type", "application/json")
				if tc.apiCode != "" {
					w.WriteHeader(http.StatusForbidden)
					json.NewEncoder(w).Encode(map[string]interface{}{"Code": tc.apiCode, "Message": "test error", "RequestId": "test-request"})
					return
				}
				json.NewEncoder(w).Encode(map[string]interface{}{"Success": true, "DtsJobId": tc.returned, "DtsInstanceID": "dtsnjl106g62735gmd", "Status": tc.status})
			})
			resource := factory()
			d := schema.TestResourceDataRaw(t, resource.Schema, nil)
			d.SetId(tc.input)
			states, err := resource.Importer.State(d, client)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("expected %q, got states=%v err=%v", tc.wantError, states, err)
				}
				if len(states) != 0 {
					t.Fatalf("failed import returned state: %v", states)
				}
			} else if err != nil || len(states) != 1 || states[0].Id() != tc.input {
				t.Fatalf("unexpected result states=%v err=%v", states, err)
			}
			if d.Id() != tc.input {
				t.Fatalf("identity switched from %q to %q", tc.input, d.Id())
			}
			if requests != 1 {
				t.Fatalf("expected one exact-ID query, got %d requests", requests)
			}
		})
	}
}

func TestUnitDtsJobReadIdentity(t *testing.T) {
	for kind, factory := range map[string]func() *schema.Resource{"sync": dts.ResourceSynchronizationJob, "subscription": dts.ResourceSubscriptionJob} {
		t.Run(kind, func(t *testing.T) {
			testDtsJobReadIdentity(t, factory)
		})
	}
}

func testDtsJobReadIdentity(t *testing.T, factory func() *schema.Resource) {
	for _, tc := range []struct {
		name, input, returned string
		missing, mismatch     bool
	}{
		{name: "canonical", input: "is210es227j5gr6", returned: "is210es227j5gr6"},
		{name: "legacy_instance_state", input: "dtsnjl106g62735gmd", returned: "is210es227j5gr6", mismatch: true},
		{name: "different_job", input: "njl106g62735gmd", returned: "is210es227j5gr6", mismatch: true},
		{name: "missing_identity", input: "njl106g62735gmd", mismatch: true},
		{name: "deleted_old_job", input: "njl106g62735gmd", missing: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			client := dtsSynchronizationJobImportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				requests++
				r.ParseForm()
				if r.Form.Get("Action") != "ListTagResources" && r.Form.Get("DtsJobId") != tc.input {
					t.Errorf("query switched identity: %v", r.Form)
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "DescribeDtsJobDetail":
					if tc.missing {
						w.WriteHeader(http.StatusNotFound)
						json.NewEncoder(w).Encode(map[string]interface{}{"Code": "Forbidden.InstanceNotFound", "Message": "deleted", "RequestId": "test-request"})
						return
					}
					json.NewEncoder(w).Encode(map[string]interface{}{"Success": true, "DtsJobId": tc.returned, "DtsInstanceID": "dtsnjl106g62735gmd", "Status": "Synchronizing", "MigrationMode": map[string]interface{}{}, "DestinationEndpoint": map[string]interface{}{}, "SourceEndpoint": map[string]interface{}{}, "SubscriptionDataType": map[string]interface{}{}, "Reserved": "{}"})
				case "ListTagResources":
					if r.Form.Get("ResourceId.1") != "dtsnjl106g62735gmd" {
						t.Errorf("unexpected tag resource: %v", r.Form)
					}
					json.NewEncoder(w).Encode(map[string]interface{}{"TagResources": map[string]interface{}{"TagResource": []interface{}{}}})
				case "DescribeDtsJobConfig":
					json.NewEncoder(w).Encode(map[string]interface{}{"Parameters": []interface{}{}})
				default:
					t.Errorf("unexpected mutation or lookup: %v", r.Form)
					w.WriteHeader(http.StatusBadRequest)
				}
			})
			resource := factory()
			d := schema.TestResourceDataRaw(t, resource.Schema, nil)
			d.SetId(tc.input)
			err := resource.Read(d, client)
			if tc.mismatch {
				if err == nil {
					t.Fatal("expected identity diagnostic")
				}
				if d.Id() != tc.input {
					t.Fatalf("state identity changed to %q", d.Id())
				}
				if requests != 1 {
					t.Fatalf("unexpected follow-up requests: %d", requests)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				want := tc.input
				if tc.missing {
					want = ""
				}
				if d.Id() != want {
					t.Fatalf("want ID %q got %q", want, d.Id())
				}
			}
		})
	}
}

// A failed mutation must not discover or take over a replacement job via its parent instance.
func TestUnitDtsJobUpdateDoesNotRebind(t *testing.T) {
	for kind, factory := range map[string]func() *schema.Resource{"sync": dts.ResourceSynchronizationJob, "subscription": dts.ResourceSubscriptionJob} {
		t.Run(kind, func(t *testing.T) {
			testDtsJobUpdateDoesNotRebind(t, factory)
		})
	}
}

func testDtsJobUpdateDoesNotRebind(t *testing.T, factory func() *schema.Resource) {
	for _, tc := range []struct {
		id            string
		businessError bool
	}{
		{id: "njl106g62735gmd"}, {id: "dtsnjl106g62735gmd"},
		{id: "njl106g62735gmd", businessError: true}, {id: "dtsnjl106g62735gmd", businessError: true},
	} {
		name := tc.id
		if tc.businessError {
			name += "_success_false"
		}
		t.Run(name, func(t *testing.T) {
			id := tc.id
			requests := 0
			client := dtsSynchronizationJobImportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				requests++
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				if r.Form.Get("Action") != "ModifyDtsJobName" || r.Form.Get("DtsJobId") != id {
					t.Errorf("unexpected recovery or mutation: %v", r.Form)
				}
				w.Header().Set("Content-Type", "application/json")
				if tc.businessError {
					json.NewEncoder(w).Encode(map[string]interface{}{"Success": false, "ErrCode": "Forbidden.InstanceNotFound", "ErrMessage": "old job missing", "RequestId": "test-request"})
					return
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]interface{}{"Code": "Forbidden.InstanceNotFound", "Message": "old job missing", "RequestId": "test-request"})
			})
			resource := factory()
			d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{"dts_instance_id": "dtsnjl106g62735gmd", "dts_job_name": "renamed"})
			d.SetId(id)
			err := resource.Update(d, client)
			if err == nil || !strings.Contains(err.Error(), "Forbidden.InstanceNotFound") {
				t.Fatalf("expected original error, got %v", err)
			}
			if d.Id() != id {
				t.Fatalf("state identity changed from %q to %q", id, d.Id())
			}
			if requests != 1 {
				t.Fatalf("expected failed mutation only, got %d requests", requests)
			}
		})
	}
}

// TestUnitDtsJobImportRefreshSequence exercises the Import→Refresh state-persistence path
// that Terraform drives end-to-end, instead of calling Importer/Read in isolation with fixed
// mocks. The canonical JobId passes import, the same identity drives the subsequent Read
// refresh, and the resource ID never switches across both phases; no candidate-discovery
// request (DescribeDtsJobs) is issued at any point. It complements the per-phase identity
// cases above and the failed-rename coverage in TestUnitDtsJobUpdateDoesNotRebind.
func TestUnitDtsJobImportRefreshSequence(t *testing.T) {
	for kind, factory := range map[string]func() *schema.Resource{"sync": dts.ResourceSynchronizationJob, "subscription": dts.ResourceSubscriptionJob} {
		t.Run(kind, func(t *testing.T) {
			id := "njl106g62735gmd"
			describeRequests := 0
			client := dtsSynchronizationJobImportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "DescribeDtsJobDetail":
					describeRequests++
					if r.Form.Get("DtsJobId") != id {
						t.Errorf("DescribeDtsJobDetail switched identity: %v", r.Form)
					}
					json.NewEncoder(w).Encode(map[string]interface{}{"Success": true, "DtsJobId": id, "DtsInstanceID": "dtsnjl106g62735gmd", "Status": "Synchronizing", "MigrationMode": map[string]interface{}{}, "DestinationEndpoint": map[string]interface{}{}, "SourceEndpoint": map[string]interface{}{}, "SubscriptionDataType": map[string]interface{}{}, "Reserved": "{}"})
				case "ListTagResources":
					json.NewEncoder(w).Encode(map[string]interface{}{"TagResources": map[string]interface{}{"TagResource": []interface{}{}}})
				case "DescribeDtsJobConfig":
					json.NewEncoder(w).Encode(map[string]interface{}{"Parameters": []interface{}{}})
				default:
					t.Errorf("unexpected action during sequence (candidate discovery must not run): %v", r.Form)
					w.WriteHeader(http.StatusBadRequest)
				}
			})
			resource := factory()
			// Phase 1: Import with the canonical JobId. importDtsJob must accept it and leave
			// the identity untouched, returning exactly this resource as the import result.
			d := schema.TestResourceDataRaw(t, resource.Schema, nil)
			d.SetId(id)
			states, err := resource.Importer.State(d, client)
			if err != nil || len(states) != 1 || states[0].Id() != id {
				t.Fatalf("import phase failed: states=%v err=%v", states, err)
			}
			// Phase 2: Refresh with the same identity. Read must keep it, surface no identity
			// diagnostic, and issue no follow-up lookup against a different job.
			if err := resource.Read(states[0], client); err != nil {
				t.Fatalf("refresh phase failed: %v", err)
			}
			if states[0].Id() != id {
				t.Fatalf("identity switched from %q to %q after refresh", id, states[0].Id())
			}
			if describeRequests != 2 {
				t.Fatalf("expected Import+Refresh to issue two exact-ID DescribeDtsJobDetail queries, got %d", describeRequests)
			}
		})
	}
}
