package paifeaturestore_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// These tests exercise alicloud_paifeaturestore_feature_view against an in-process
// HTTP mock. They intentionally do NOT call acctest.PreCheck (no real cloud
// credentials): the provider is configured with mock creds, protocol=http, and an
// ALIBABA_CLOUD_ENDPOINT_PAIFEATURESTORE pointing at the mock server, so the ROA
// client signs against localhost and the mock dispatches on path+method. This is
// the only honest way to cover Read/Import/no-op/failure/collection-order without
// the parent FeatureStore instance/project chain, which has no TF resources yet.

const (
	pfsTestInstanceID   = "fsi-mock-001"
	pfsFailInstanceID   = "fsi-fail"
	pfsTestProjectID    = "proj-mock-001"
	pfsTestResourceName = "alicloud_paifeaturestore_feature_view.default"
)

// pfsMockStore is the in-memory backend the mock server reads and mutates. It is
// not safe to share across concurrent *resource.Test runs; each test gets its own.
type pfsMockStore struct {
	mu      sync.Mutex
	views   map[string]map[string]interface{}
	counter int
}

func newPfsMockStore() *pfsMockStore {
	return &pfsMockStore{views: map[string]map[string]interface{}{}}
}

// startPfsMock stands up the mock ROA server and wires the provider endpoint at it.
func startPfsMock(t *testing.T) *pfsMockStore {
	t.Helper()
	store := newPfsMockStore()
	srv := httptest.NewServer(http.HandlerFunc(store.handle))
	t.Cleanup(srv.Close)
	host := strings.TrimPrefix(srv.URL, "http://")
	// loadEndpoint resolves ALIBABA_CLOUD_ENDPOINT_{PRODUCT} first; NO_PROXY keeps
	// the localhost traffic out of any corporate proxy.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_PAIFEATURESTORE", host)
	t.Setenv("NO_PROXY", host)
	return store
}

func (s *pfsMockStore) handle(w http.ResponseWriter, r *http.Request) {
	// ROA dispatch: method + path. The path is
	//   /api/v1/instances/{InstanceId}/featureviews[/{FeatureViewId}]
	// Signing headers and the RegionId query are ignored; the mock never validates
	// credentials, mirroring the rds mock precedent.
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	idx := -1
	for i, p := range parts {
		if p == "featureviews" {
			idx = i
			break
		}
	}
	if idx == -1 {
		writeJSON(w, http.StatusOK, errObj("UnknownPath", fmt.Sprintf("no featureviews in %s", r.URL.Path)))
		return
	}
	var instanceId, featureViewId string
	if idx >= 1 {
		instanceId = parts[idx-1]
	}
	if idx+1 < len(parts) {
		featureViewId = parts[idx+1]
	}

	switch {
	case r.Method == http.MethodPost && featureViewId == "":
		s.create(w, r, instanceId)
	case r.Method == http.MethodGet && featureViewId != "":
		s.get(w, featureViewId)
	case r.Method == http.MethodDelete && featureViewId != "":
		s.delete(w, featureViewId)
	default:
		writeJSON(w, http.StatusOK, errObj("UnsupportedMethod", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	}
}

func (s *pfsMockStore) create(w http.ResponseWriter, r *http.Request, instanceId string) {
	if instanceId == pfsFailInstanceID {
		// Drive the Create failure path: formatError builds a *tea.SDKError whose
		// Code the resource surfaces verbatim.
		writeJSON(w, http.StatusOK, errObj("CreateFeatureViewErrorProblem", "simulated create failure"))
		return
	}
	body := readJSONBody(r)
	s.mu.Lock()
	s.counter++
	id := fmt.Sprintf("fv-%d", s.counter)
	obj := map[string]interface{}{}
	for k, v := range body {
		obj[k] = v
	}
	obj["FeatureViewId"] = id
	obj["GmtCreateTime"] = "2026-01-01T00:00:00Z"
	obj["RegionId"] = "cn-hangzhou"
	s.views[id] = obj
	s.mu.Unlock()
	// Success: no Code field, so formatError returns nil and Create reads FeatureViewId.
	writeJSON(w, http.StatusOK, obj)
}

func (s *pfsMockStore) get(w http.ResponseWriter, featureViewId string) {
	s.mu.Lock()
	obj, ok := s.views[featureViewId]
	s.mu.Unlock()
	if !ok {
		// Not-found surfaces as the resourceNotExistErrorCodes entry; isFeatureViewNotExist
		// matches it via *tea.SDKError.Code, so Read removes the resource instead of erroring.
		writeJSON(w, http.StatusOK, errObj("DoConfigCenterRequestErrorProblem", "feature view not found"))
		return
	}
	writeJSON(w, http.StatusOK, obj)
}

func (s *pfsMockStore) delete(w http.ResponseWriter, featureViewId string) {
	s.mu.Lock()
	delete(s.views, featureViewId)
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]interface{}{"RequestId": "r-del"})
}

func (s *pfsMockStore) drop(featureViewId string) {
	s.mu.Lock()
	delete(s.views, featureViewId)
	s.mu.Unlock()
}

func (s *pfsMockStore) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counter
}

func errObj(code, message string) map[string]interface{} {
	return map[string]interface{}{"Code": code, "Message": message, "RequestId": "r-mock"}
}

func readJSONBody(r *http.Request) map[string]interface{} {
	out := map[string]interface{}{}
	if r.Body == nil {
		return out
	}
	defer r.Body.Close()
	_ = json.NewDecoder(r.Body).Decode(&out)
	return out
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func pfsProviderConfig() string {
	return `
provider "alicloud" {
  access_key             = "test-key"
  secret_key             = "test-secret"
  region                 = "cn-hangzhou"
  protocol               = "HTTP"
  skip_region_validation = true
}
`
}

func pfsBasicConfig() string {
	return pfsProviderConfig() + `
resource "alicloud_paifeaturestore_feature_view" "default" {
  instance_id       = "` + pfsTestInstanceID + `"
  feature_view_name = "tf-test-fv"
  project_id        = "` + pfsTestProjectID + `"
  type              = "Batch"
  write_method      = "ByReadyMadeTable"
  sync_online_table = true
  ttl               = 86400

  fields {
    name = "user_id"
    type = "STRING"
    attributes = ["PrimaryKey"]
  }

  fields {
    name = "age"
    type = "INT64"
  }
}
`
}

func pfsOrderedConfig() string {
	return pfsProviderConfig() + `
resource "alicloud_paifeaturestore_feature_view" "default" {
  instance_id       = "` + pfsTestInstanceID + `"
  feature_view_name = "tf-test-fv-ordered"
  project_id        = "` + pfsTestProjectID + `"
  type              = "Batch"
  write_method      = "ByReadyMadeTable"
  sync_online_table = true
  ttl               = 3600

  fields {
    name = "z_first"
    type = "STRING"
  }

  fields {
    name = "a_second"
    type = "INT64"
  }

  fields {
    name = "m_third"
    type = "DOUBLE"
  }
}
`
}

func pfsFailConfig() string {
	return pfsProviderConfig() + `
resource "alicloud_paifeaturestore_feature_view" "default" {
  instance_id       = "` + pfsFailInstanceID + `"
  feature_view_name = "tf-test-fv-fail"
  project_id        = "` + pfsTestProjectID + `"
  type              = "Batch"
  write_method      = "ByReadyMadeTable"
  sync_online_table = true
  ttl               = 60

  fields {
    name = "user_id"
    type = "STRING"
  }
}
`
}

// TestAccPaiFeatureStoreFeatureView_basic covers Create, Read, and the no-op plan
// (re-applying the same config makes zero Create/Update/Delete calls), then Destroy.
func TestAccPaiFeatureStoreFeatureView_basic(t *testing.T) {
	store := startPfsMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkPfsDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: pfsBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pfsTestResourceName, "id", pfsTestInstanceID+":fv-1"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "instance_id", pfsTestInstanceID),
					resource.TestCheckResourceAttr(pfsTestResourceName, "feature_view_id", "fv-1"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "feature_view_name", "tf-test-fv"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "project_id", pfsTestProjectID),
					resource.TestCheckResourceAttr(pfsTestResourceName, "type", "Batch"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "write_method", "ByReadyMadeTable"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "sync_online_table", "true"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "ttl", "86400"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttrSet(pfsTestResourceName, "create_time"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.#", "2"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.0.name", "user_id"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.1.name", "age"),
				),
			},
			{
				// No-op: identical config must refresh (one GET) and plan no changes,
				// so no additional Create is issued. The counter proves it stayed at 1.
				Config: pfsBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pfsTestResourceName, "feature_view_id", "fv-1"),
					func(*terraform.State) error {
						if got := store.count(); got != 1 {
							return fmt.Errorf("expected no-op (counter=1), got %d Create calls", got)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccPaiFeatureStoreFeatureView_import covers terraform import via the
// composite "<InstanceId>:<FeatureViewId>" identity.
func TestAccPaiFeatureStoreFeatureView_import(t *testing.T) {
	store := startPfsMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkPfsDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: pfsBasicConfig(),
			},
			{
				ResourceName:      pfsTestResourceName,
				ImportState:       true,
				ImportStateId:     pfsTestInstanceID + ":fv-1",
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccPaiFeatureStoreFeatureView_collectionOrder verifies that the ordered Fields
// list preserves element order across Read, satisfying the collection-order coverage.
func TestAccPaiFeatureStoreFeatureView_collectionOrder(t *testing.T) {
	store := startPfsMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkPfsDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: pfsOrderedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.#", "3"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.0.name", "z_first"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.0.type", "STRING"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.1.name", "a_second"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.1.type", "INT64"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.2.name", "m_third"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "fields.2.type", "DOUBLE"),
				),
			},
		},
	})
}

// TestAccPaiFeatureStoreFeatureView_readGone covers the not-found Read path: after a
// server-side deletion the next refresh's GetFeatureView returns the
// DoConfigCenterRequestErrorProblem code, Read removes the resource, and Terraform
// recreates it (a fresh Create producing a new feature_view_id).
func TestAccPaiFeatureStoreFeatureView_readGone(t *testing.T) {
	store := startPfsMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkPfsDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: pfsBasicConfig(),
				// After Check drops fv-1 from the mock store, the refresh's Read sees the
				// not-found code and removes the resource from state; the refresh plan then
				// recreates it (+ create). ExpectNonEmptyPlan asserts that recreation rather
				// than the default empty-plan idempotency check, matching the documented
				// "Terraform recreates it (a fresh Create producing a new feature_view_id)"
				// behavior verified by Step 2 (fv-2).
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pfsTestResourceName, "feature_view_id", "fv-1"),
					func(*terraform.State) error {
						store.drop("fv-1") // simulate the resource vanishing server-side
						return nil
					},
				),
			},
			{
				Config: pfsBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Recreated: fv-1 is gone, Read removed it, the plan recreated as fv-2.
					resource.TestCheckResourceAttr(pfsTestResourceName, "feature_view_id", "fv-2"),
					resource.TestCheckResourceAttr(pfsTestResourceName, "id", pfsTestInstanceID+":fv-2"),
				),
			},
		},
	})
}

// TestAccPaiFeatureStoreFeatureView_failureCreate covers the Create failure path: the
// mock returns an error code for the fsi-fail instance, and the resource surfaces it.
func TestAccPaiFeatureStoreFeatureView_failureCreate(t *testing.T) {
	startPfsMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      pfsFailConfig(),
				ExpectError: regexp.MustCompile("CreateFeatureViewErrorProblem"),
			},
		},
	})
}

// checkPfsDestroyed asserts the mock backend is empty after Destroy ran.
func checkPfsDestroyed(store *pfsMockStore) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if got := len(store.views); got != 0 {
			return fmt.Errorf("expected empty mock backend after destroy, got %d views", got)
		}
		return nil
	}
}
