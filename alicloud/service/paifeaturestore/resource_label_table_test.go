package paifeaturestore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paifeaturestore"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// These tests exercise alicloud_pai_feature_store_label_table against an
// in-process HTTP mock. They intentionally do NOT call acctest.PreCheck (no real
// cloud credentials): the provider is configured with mock creds, protocol=http,
// and an ALIBABA_CLOUD_ENDPOINT_PAIFEATURESTORE pointing at the mock server, so
// the ROA client signs against localhost and the mock dispatches on
// path+method. The parent FeatureStore instance, project and datasource have no
// Terraform resources yet, so this is the only way to cover
// Create/Read/Update/Delete/Import and failure paths without a paid instance.
//
// The functions use the TestAccAliCloud prefix because the remote ACC harness
// (pre_pr_checks.py --test-case) only accepts TestAcc* names and is the sole
// execution entry. They are not real-cloud acceptance tests and must not be
// read as proof of service-side behavior; the real-cloud limitation is
// documented in the RD implementation report.

const (
	labelTableTestInstanceID   = "fsi-mock-001"
	labelTableFailInstanceID   = "fsi-fail"
	labelTableTestProjectID    = "proj-mock-001"
	labelTableTestResourceName = "alicloud_pai_feature_store_label_table.default"
	labelTableTestResourceType = "alicloud_pai_feature_store_label_table"
)

type labelTableMockStore struct {
	mu            sync.Mutex
	tables        map[string]map[string]interface{}
	counter       int
	updateCount   int
	failGet       bool
	reorderFields bool
	createBody    map[string]interface{}
	updateBody    map[string]interface{}
}

func newLabelTableMockStore(reorderFields bool) *labelTableMockStore {
	return &labelTableMockStore{
		tables:        map[string]map[string]interface{}{},
		reorderFields: reorderFields,
	}
}

func startLabelTableMock(t *testing.T) *labelTableMockStore {
	return startLabelTableMockWith(t, true)
}

// startLabelTableMockWith starts the in-process mock. When reorderFields is
// true the mock returns the Fields array in the reverse order it was sent,
// mimicking a server that canonicalizes the array, so the Read normalization
// is exercised. The import test disables it because import has no prior
// reference to restore the order from.
func startLabelTableMockWith(t *testing.T, reorderFields bool) *labelTableMockStore {
	t.Helper()
	store := newLabelTableMockStore(reorderFields)
	srv := httptest.NewServer(http.HandlerFunc(store.handle))
	t.Cleanup(srv.Close)
	host := strings.TrimPrefix(srv.URL, "http://")
	// loadEndpoint resolves ALIBABA_CLOUD_ENDPOINT_{PRODUCT} first; NO_PROXY keeps
	// the localhost traffic out of any corporate proxy.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_PAIFEATURESTORE", host)
	t.Setenv("NO_PROXY", host)
	return store
}

func (s *labelTableMockStore) handle(w http.ResponseWriter, r *http.Request) {
	// ROA dispatch: method + path. The path is
	//   /api/v1/instances/{InstanceId}/labeltables[/{LabelTableId}]
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	idx := -1
	for i, p := range parts {
		if p == "labeltables" {
			idx = i
			break
		}
	}
	if idx == -1 {
		writeLabelTableJSON(w, http.StatusOK, labelTableErrObj("UnknownPath", fmt.Sprintf("no labeltables in %s", r.URL.Path)))
		return
	}
	var instanceId, labelTableId string
	if idx >= 1 {
		instanceId = parts[idx-1]
	}
	if idx+1 < len(parts) {
		labelTableId = parts[idx+1]
	}

	switch {
	case r.Method == http.MethodPost && labelTableId == "":
		s.create(w, r, instanceId)
	case r.Method == http.MethodGet && labelTableId != "":
		s.get(w, labelTableId)
	case r.Method == http.MethodPut && labelTableId != "":
		s.update(w, r, labelTableId)
	case r.Method == http.MethodDelete && labelTableId != "":
		s.delete(w, labelTableId)
	default:
		writeLabelTableJSON(w, http.StatusOK, labelTableErrObj("UnsupportedMethod", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	}
}

func (s *labelTableMockStore) create(w http.ResponseWriter, r *http.Request, instanceId string) {
	if instanceId == labelTableFailInstanceID {
		// Drive the Create failure path: formatError builds a *tea.SDKError whose
		// Code the resource surfaces verbatim.
		writeLabelTableJSON(w, http.StatusOK, labelTableErrObj("CreateLabelTableErrorProblem", "simulated create failure"))
		return
	}
	body := readLabelTableBody(r)
	s.mu.Lock()
	s.createBody = body
	s.counter++
	id := fmt.Sprintf("lt-%d", s.counter)
	obj := map[string]interface{}{
		"Name":          body["Name"],
		"ProjectId":     body["ProjectId"],
		"DatasourceId":  body["DatasourceId"],
		"Fields":        body["Fields"],
		"GmtCreateTime": "2026-01-01T00:00:00Z",
	}
	s.tables[id] = obj
	s.mu.Unlock()
	writeLabelTableJSON(w, http.StatusOK, map[string]interface{}{
		"LabelTableId": id,
	})
}

func (s *labelTableMockStore) get(w http.ResponseWriter, labelTableId string) {
	s.mu.Lock()
	if s.failGet {
		// One-shot failure: the next Read after this returns a non-NotFound error,
		// proving the resource does not clear state on arbitrary API errors.
		s.failGet = false
		s.mu.Unlock()
		writeLabelTableJSON(w, http.StatusOK, labelTableErrObj("InternalError", "simulated get failure"))
		return
	}
	obj, ok := s.tables[labelTableId]
	reorder := s.reorderFields
	s.mu.Unlock()
	if !ok {
		// Not-found surfaces as the resourceNotExistErrorCodes entry;
		// isLabelTableNotExist matches it via *tea.SDKError.Code, so Read removes
		// the resource instead of erroring.
		writeLabelTableJSON(w, http.StatusOK, labelTableErrObj("DoConfigCenterRequestErrorProblem", "label table not found"))
		return
	}
	writeLabelTableJSON(w, http.StatusOK, labelTableResponse(obj, reorder))
}

// labelTableResponse simulates server-side canonicalization: the Attributes of
// every field come back reversed, and when reorderFields is set the Fields array
// itself is returned in reverse order too. Attributes is a Set attribute, so it
// must stay converged regardless of response order; the Fields order is
// restored by Read from the plan/state reference, and the ACC steps verify that
// the following plan stays empty.
func labelTableResponse(obj map[string]interface{}, reorderFields bool) map[string]interface{} {
	out := make(map[string]interface{}, len(obj))
	for k, v := range obj {
		out[k] = v
	}
	fieldsArr, ok := out["Fields"].([]interface{})
	if !ok {
		return out
	}
	fields := make([]interface{}, len(fieldsArr))
	for i, f := range fieldsArr {
		field, ok := f.(map[string]interface{})
		if !ok {
			fields[i] = f
			continue
		}
		fieldCopy := make(map[string]interface{}, len(field))
		for k, v := range field {
			fieldCopy[k] = v
		}
		if attrs, ok := fieldCopy["Attributes"].([]interface{}); ok {
			fieldCopy["Attributes"] = reversedLabelTableStrings(attrs)
		}
		fields[i] = fieldCopy
	}
	if reorderFields {
		fields = reversedLabelTableStrings(fields)
	}
	out["Fields"] = fields
	return out
}

func reversedLabelTableStrings(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[len(in)-1-i] = v
	}
	return out
}

func (s *labelTableMockStore) update(w http.ResponseWriter, r *http.Request, labelTableId string) {
	body := readLabelTableBody(r)
	s.mu.Lock()
	obj, ok := s.tables[labelTableId]
	if !ok {
		s.mu.Unlock()
		writeLabelTableJSON(w, http.StatusOK, labelTableErrObj("DoConfigCenterRequestErrorProblem", "label table not found"))
		return
	}
	s.updateBody = body
	obj["Name"] = body["Name"]
	obj["DatasourceId"] = body["DatasourceId"]
	obj["Fields"] = body["Fields"]
	s.updateCount++
	s.mu.Unlock()
	writeLabelTableJSON(w, http.StatusOK, map[string]interface{}{})
}

func (s *labelTableMockStore) delete(w http.ResponseWriter, labelTableId string) {
	s.mu.Lock()
	delete(s.tables, labelTableId)
	s.mu.Unlock()
	writeLabelTableJSON(w, http.StatusOK, map[string]interface{}{})
}

func (s *labelTableMockStore) drop(labelTableId string) {
	s.mu.Lock()
	delete(s.tables, labelTableId)
	s.mu.Unlock()
}

func (s *labelTableMockStore) failNextGet() {
	s.mu.Lock()
	s.failGet = true
	s.mu.Unlock()
}

func (s *labelTableMockStore) createCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counter
}

func (s *labelTableMockStore) updates() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateCount
}

func (s *labelTableMockStore) tableCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.tables)
}

func (s *labelTableMockStore) lastCreateBody() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.createBody
}

func (s *labelTableMockStore) lastUpdateBody() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateBody
}

// labelTableRequestBodyField returns the object at Fields[index] of a captured
// request body after checking the field name, so assertions read the values the
// resource actually serialized.
func labelTableRequestBodyField(body map[string]interface{}, index int, name string) (map[string]interface{}, error) {
	rawFields, ok := body["Fields"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("request body carries no Fields array: %v", body["Fields"])
	}
	if len(rawFields) <= index {
		return nil, fmt.Errorf("request body Fields has %d elements, want at least %d", len(rawFields), index+1)
	}
	field, ok := rawFields[index].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("request body Fields[%d] is not an object: %v", index, rawFields[index])
	}
	if got := field["Name"]; got != name {
		return nil, fmt.Errorf("request body Fields[%d].Name = %v, want %q", index, got, name)
	}
	return field, nil
}

func labelTableRequestBodyAttributes(field map[string]interface{}) ([]string, error) {
	raw, ok := field["Attributes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("request body field %v carries no Attributes array: %v", field["Name"], field["Attributes"])
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		name, _ := v.(string)
		out = append(out, name)
	}
	return out, nil
}

func labelTableRequireAttributes(field map[string]interface{}, want ...string) error {
	got, err := labelTableRequestBodyAttributes(field)
	if err != nil {
		return err
	}
	if len(got) != len(want) {
		return fmt.Errorf("request body field %v Attributes = %v, want %v", field["Name"], got, want)
	}
	for _, w := range want {
		found := false
		for _, g := range got {
			if g == w {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("request body field %v Attributes = %v, want %v", field["Name"], got, want)
		}
	}
	return nil
}

func labelTableErrObj(code, message string) map[string]interface{} {
	return map[string]interface{}{"Code": code, "Message": message}
}

func readLabelTableBody(r *http.Request) map[string]interface{} {
	out := map[string]interface{}{}
	if r.Body == nil {
		return out
	}
	defer r.Body.Close()
	_ = json.NewDecoder(r.Body).Decode(&out)
	return out
}

func writeLabelTableJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func labelTableProviderConfig() string {
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

// labelTableBasicConfig creates a label table with two fields. Both fields
// populate every member of the flat fields object type, which requires all four
// attributes per element.
func labelTableBasicConfig() string {
	return labelTableProviderConfig() + `
resource "alicloud_pai_feature_store_label_table" "default" {
  instance_id      = "` + labelTableTestInstanceID + `"
  project_id       = "` + labelTableTestProjectID + `"
  label_table_name = "tf_testAcc_label_table"
  datasource_id    = "ds-mock-001"
  fields = [
    {
      name                = "user_id"
      type                = "STRING"
      attributes          = ["Partition"]
      aligned_entity_name = "user_id"
    },
    {
      name                = "label_value"
      type                = "DOUBLE"
      attributes          = ["LabelField"]
      aligned_entity_name = "user_id"
    },
  ]
}
`
}

// labelTableUpdatedConfig changes only in-place attributes: name,
// datasource_id and the attributes/aligned_entity_name of an existing field.
// The API supports these through UpdateLabelTable, so no replacement is
// expected.
func labelTableUpdatedConfig() string {
	return labelTableProviderConfig() + `
resource "alicloud_pai_feature_store_label_table" "default" {
  instance_id      = "` + labelTableTestInstanceID + `"
  project_id       = "` + labelTableTestProjectID + `"
  label_table_name = "tf_testAcc_label_table_updated"
  datasource_id    = "ds-mock-002"
  fields = [
    {
      name                = "user_id"
      type                = "STRING"
      attributes          = ["Partition", "EventTimeIsUnixTime"]
      aligned_entity_name = "item_id"
    },
    {
      name                = "label_value"
      type                = "DOUBLE"
      attributes          = []
      aligned_entity_name = ""
    },
  ]
}
`
}

// labelTableRenamedFieldConfig renames a field. The API cannot rename existing
// fields, so ModifyPlan must force a replacement.
func labelTableRenamedFieldConfig() string {
	return labelTableProviderConfig() + `
resource "alicloud_pai_feature_store_label_table" "default" {
  instance_id      = "` + labelTableTestInstanceID + `"
  project_id       = "` + labelTableTestProjectID + `"
  label_table_name = "tf_testAcc_label_table"
  datasource_id    = "ds-mock-001"
  fields = [
    {
      name                = "user_key"
      type                = "STRING"
      attributes          = ["Partition"]
      aligned_entity_name = "user_id"
    },
    {
      name                = "label_value"
      type                = "DOUBLE"
      attributes          = ["LabelField"]
      aligned_entity_name = "user_id"
    },
  ]
}
`
}

// labelTableReducedConfig drops the second field. Removing a list element never
// visits the removed element's plan modifiers, so this config is what forces
// the resource-level ModifyPlan replacement guard.
func labelTableReducedConfig() string {
	return labelTableProviderConfig() + `
resource "alicloud_pai_feature_store_label_table" "default" {
  instance_id      = "` + labelTableTestInstanceID + `"
  project_id       = "` + labelTableTestProjectID + `"
  label_table_name = "tf_testAcc_label_table"
  datasource_id    = "ds-mock-001"
  fields = [
    {
      name                = "user_id"
      type                = "STRING"
      attributes          = ["Partition"]
      aligned_entity_name = "user_id"
    },
  ]
}
`
}

func labelTableFailConfig() string {
	return labelTableProviderConfig() + `
resource "alicloud_pai_feature_store_label_table" "default" {
  instance_id      = "` + labelTableFailInstanceID + `"
  project_id       = "` + labelTableTestProjectID + `"
  label_table_name = "tf_testAcc_label_table_fail"
  datasource_id    = "ds-mock-001"
  fields = [
    {
      name                = "user_id"
      type                = "STRING"
      attributes          = ["Partition"]
      aligned_entity_name = "user_id"
    },
  ]
}
`
}

// TestAccAliCloudPaiFeatureStoreLabelTable_basic covers Create, Read, no-op
// refresh and Destroy.
func TestAccAliCloudPaiFeatureStoreLabelTable_basic(t *testing.T) {
	store := startLabelTableMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkLabelTableDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "id", labelTableTestInstanceID+":lt-1"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "instance_id", labelTableTestInstanceID),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_name", "tf_testAcc_label_table"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "project_id", labelTableTestProjectID),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "datasource_id", "ds-mock-001"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttrSet(labelTableTestResourceName, "create_time"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.#", "2"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.0.name", "user_id"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.0.type", "STRING"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.0.aligned_entity_name", "user_id"),
					resource.TestCheckTypeSetElemAttr(labelTableTestResourceName, "fields.0.attributes.*", "Partition"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.1.name", "label_value"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.1.type", "DOUBLE"),
					resource.TestCheckTypeSetElemAttr(labelTableTestResourceName, "fields.1.attributes.*", "LabelField"),
					func(*terraform.State) error {
						body := store.lastCreateBody()
						if got := body["Name"]; got != "tf_testAcc_label_table" {
							return fmt.Errorf("CreateLabelTable body Name = %v, want tf_testAcc_label_table", got)
						}
						if got := body["ProjectId"]; got != labelTableTestProjectID {
							return fmt.Errorf("CreateLabelTable body ProjectId = %v, want %s", got, labelTableTestProjectID)
						}
						if got := body["DatasourceId"]; got != "ds-mock-001" {
							return fmt.Errorf("CreateLabelTable body DatasourceId = %v, want ds-mock-001", got)
						}
						first, err := labelTableRequestBodyField(body, 0, "user_id")
						if err != nil {
							return err
						}
						if got := first["Type"]; got != "STRING" {
							return fmt.Errorf("CreateLabelTable body Fields[0].Type = %v, want STRING", got)
						}
						if got := first["AlignedEntityName"]; got != "user_id" {
							return fmt.Errorf("CreateLabelTable body Fields[0].AlignedEntityName = %v, want user_id", got)
						}
						if err := labelTableRequireAttributes(first, "Partition"); err != nil {
							return err
						}
						second, err := labelTableRequestBodyField(body, 1, "label_value")
						if err != nil {
							return err
						}
						if got := second["Type"]; got != "DOUBLE" {
							return fmt.Errorf("CreateLabelTable body Fields[1].Type = %v, want DOUBLE", got)
						}
						return labelTableRequireAttributes(second, "LabelField")
					},
				),
			},
			{
				// No-op: identical config must refresh and plan no changes, so no
				// additional Create is issued. The mock returns the Fields array in
				// reverse order, so this step also proves that Read restores the
				// reference order instead of producing a replacement diff.
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
					func(*terraform.State) error {
						if got := store.createCount(); got != 1 {
							return fmt.Errorf("expected no-op (createCount=1), got %d Create calls", got)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccAliCloudPaiFeatureStoreLabelTable_update covers the in-place
// UpdateLabelTable path: label_table_name, datasource_id and the attributes and
// aligned_entity_name of existing fields change without a replacement.
func TestAccAliCloudPaiFeatureStoreLabelTable_update(t *testing.T) {
	store := startLabelTableMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkLabelTableDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
				),
			},
			{
				Config: labelTableUpdatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "id", labelTableTestInstanceID+":lt-1"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_name", "tf_testAcc_label_table_updated"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "datasource_id", "ds-mock-002"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.#", "2"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.0.aligned_entity_name", "item_id"),
					resource.TestCheckTypeSetElemAttr(labelTableTestResourceName, "fields.0.attributes.*", "Partition"),
					resource.TestCheckTypeSetElemAttr(labelTableTestResourceName, "fields.0.attributes.*", "EventTimeIsUnixTime"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.1.attributes.#", "0"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.1.aligned_entity_name", ""),
					func(*terraform.State) error {
						if got := store.createCount(); got != 1 {
							return fmt.Errorf("expected in-place update (createCount=1), got %d Create calls", got)
						}
						if got := store.updates(); got != 1 {
							return fmt.Errorf("expected one UpdateLabelTable call, got %d", got)
						}
						body := store.lastUpdateBody()
						if got := body["Name"]; got != "tf_testAcc_label_table_updated" {
							return fmt.Errorf("UpdateLabelTable body Name = %v, want tf_testAcc_label_table_updated", got)
						}
						if got := body["DatasourceId"]; got != "ds-mock-002" {
							return fmt.Errorf("UpdateLabelTable body DatasourceId = %v, want ds-mock-002", got)
						}
						first, err := labelTableRequestBodyField(body, 0, "user_id")
						if err != nil {
							return err
						}
						if got := first["AlignedEntityName"]; got != "item_id" {
							return fmt.Errorf("UpdateLabelTable body Fields[0].AlignedEntityName = %v, want item_id", got)
						}
						if err := labelTableRequireAttributes(first, "Partition", "EventTimeIsUnixTime"); err != nil {
							return err
						}
						second, err := labelTableRequestBodyField(body, 1, "label_value")
						if err != nil {
							return err
						}
						if got := second["AlignedEntityName"]; got != "" {
							return fmt.Errorf("UpdateLabelTable body Fields[1].AlignedEntityName = %v, want the empty string", got)
						}
						// The API requires Attributes on every update item; an empty
						// array must still be serialized.
						return labelTableRequireAttributes(second)
					},
				),
			},
		},
	})
}

// TestAccAliCloudPaiFeatureStoreLabelTable_replace covers a field rename, which
// the API cannot apply to an existing table: ModifyPlan must force a new
// resource and never call UpdateLabelTable.
func TestAccAliCloudPaiFeatureStoreLabelTable_replace(t *testing.T) {
	store := startLabelTableMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkLabelTableDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
				),
			},
			{
				Config: labelTableRenamedFieldConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-2"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "id", labelTableTestInstanceID+":lt-2"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.0.name", "user_key"),
					func(*terraform.State) error {
						if got := store.createCount(); got != 2 {
							return fmt.Errorf("expected replacement (createCount=2), got %d Create calls", got)
						}
						if got := store.updates(); got != 0 {
							return fmt.Errorf("expected no UpdateLabelTable call during replacement, got %d", got)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccAliCloudPaiFeatureStoreLabelTable_fieldRemoval covers the list-element
// removal gap: deleting a fields entry must also recreate the label table
// instead of planning an in-place update the API would reject.
func TestAccAliCloudPaiFeatureStoreLabelTable_fieldRemoval(t *testing.T) {
	store := startLabelTableMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkLabelTableDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.#", "2"),
				),
			},
			{
				Config: labelTableReducedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-2"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "id", labelTableTestInstanceID+":lt-2"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.#", "1"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "fields.0.name", "user_id"),
					func(*terraform.State) error {
						if got := store.createCount(); got != 2 {
							return fmt.Errorf("expected replacement (createCount=2), got %d Create calls", got)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccAliCloudPaiFeatureStoreLabelTable_import covers terraform import via
// the composite "<InstanceId>:<LabelTableId>" identity.
func TestAccAliCloudPaiFeatureStoreLabelTable_import(t *testing.T) {
	// Import has no plan or prior state to restore the field order from, so the
	// mock returns the stored order here; order normalization is covered by the
	// basic and update cases.
	store := startLabelTableMockWith(t, false)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkLabelTableDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
				),
			},
			{
				ResourceName:      labelTableTestResourceName,
				ImportState:       true,
				ImportStateId:     labelTableTestInstanceID + ":lt-1",
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudPaiFeatureStoreLabelTable_readGone covers the not-found Read
// path: after a server-side deletion the next refresh's GetLabelTable returns
// the DoConfigCenterRequestErrorProblem code, Read removes the resource, and
// Terraform recreates it.
func TestAccAliCloudPaiFeatureStoreLabelTable_readGone(t *testing.T) {
	store := startLabelTableMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkLabelTableDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config:             labelTableBasicConfig(),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
					func(*terraform.State) error {
						store.drop("lt-1") // simulate the label table vanishing server-side
						return nil
					},
				),
			},
			{
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-2"),
					resource.TestCheckResourceAttr(labelTableTestResourceName, "id", labelTableTestInstanceID+":lt-2"),
				),
			},
		},
	})
}

// TestAccAliCloudPaiFeatureStoreLabelTable_readError covers a non-NotFound Read
// failure: the refresh must surface the API error instead of removing the
// resource from state, so a transient internal error cannot make Terraform
// believe the label table vanished.
func TestAccAliCloudPaiFeatureStoreLabelTable_readError(t *testing.T) {
	store := startLabelTableMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkLabelTableDestroyed(store),
		Steps: []resource.TestStep{
			{
				Config: labelTableBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(labelTableTestResourceName, "label_table_id", "lt-1"),
				),
			},
			{
				PreConfig:   func() { store.failNextGet() },
				Config:      labelTableBasicConfig(),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("InternalError"),
			},
		},
	})
}

// TestAccAliCloudPaiFeatureStoreLabelTable_failureCreate covers the Create
// failure path: the mock returns an error code for the fsi-fail instance, and
// the resource surfaces it.
func TestAccAliCloudPaiFeatureStoreLabelTable_failureCreate(t *testing.T) {
	startLabelTableMock(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      labelTableFailConfig(),
				ExpectError: regexp.MustCompile("CreateLabelTableErrorProblem"),
			},
		},
	})
}

func checkLabelTableDestroyed(store *labelTableMockStore) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if got := store.tableCount(); got != 0 {
			return fmt.Errorf("expected empty mock backend after destroy, got %d label tables", got)
		}
		return nil
	}
}

// TestUnitAliCloudPaiFeatureStoreLabelTableSchema pins the parts of the schema
// that the protocol-5 provider depends on: the resource type name, the
// required/computed split of the top-level attributes, and the fact that
// fields is a flat ListAttribute with an ObjectType element rather than a
// nested attribute (which protocol 5 rejects at schema load).
func TestUnitAliCloudPaiFeatureStoreLabelTableSchema(t *testing.T) {
	implementation := paifeaturestore.NewLabelTableResource()

	metadataResponse := fwresource.MetadataResponse{}
	implementation.Metadata(context.Background(), fwresource.MetadataRequest{ProviderTypeName: "alicloud"}, &metadataResponse)
	if metadataResponse.TypeName != labelTableTestResourceType {
		t.Fatalf("expected type name %q, got %q", labelTableTestResourceType, metadataResponse.TypeName)
	}

	schemaResponse := fwresource.SchemaResponse{}
	implementation.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResponse)

	required := map[string]bool{
		"instance_id":      true,
		"project_id":       true,
		"label_table_name": true,
		"datasource_id":    true,
		"fields":           true,
	}
	computed := map[string]bool{
		"id":             true,
		"label_table_id": true,
		"create_time":    true,
		"region_id":      true,
	}
	for name, attribute := range schemaResponse.Schema.Attributes {
		if _, ok := attribute.(schema.NestedAttribute); ok {
			t.Fatalf("attribute %q must not be a nested attribute under protocol 5", name)
		}
		if want := required[name]; want != attribute.IsRequired() {
			t.Fatalf("attribute %q: expected required=%v, got %v", name, want, attribute.IsRequired())
		}
		if want := computed[name]; want != attribute.IsComputed() {
			t.Fatalf("attribute %q: expected computed=%v, got %v", name, want, attribute.IsComputed())
		}
	}

	fieldsAttribute, ok := schemaResponse.Schema.Attributes["fields"].(schema.ListAttribute)
	if !ok {
		t.Fatalf("expected fields to be a schema.ListAttribute, got %T", schemaResponse.Schema.Attributes["fields"])
	}
	objectType, ok := fieldsAttribute.ElementType.(types.ObjectType)
	if !ok {
		t.Fatalf("expected fields element type to be types.ObjectType, got %T", fieldsAttribute.ElementType)
	}
	for _, name := range []string{"name", "type", "attributes", "aligned_entity_name"} {
		if _, ok := objectType.AttrTypes[name]; !ok {
			t.Fatalf("expected fields object type to carry attribute %q, got %v", name, objectType.AttrTypes)
		}
	}
}
