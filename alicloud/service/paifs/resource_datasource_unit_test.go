package paifs

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// paifsTestClient builds a credential-free *connectivity.AliyunClient whose RPC
// transport points at an in-process httptest.Server, mirroring the established
// SDKv2 unit-test construction (see ecsSnapshotTestClient / rdsAccountTestClient
// in the root alicloud package). The PaiFeatureStore endpoint is stored under
// the lowercased product code "paifeaturestore", which is how RpcPost's
// loadApiEndpoint resolves it (ConvertKebabToSnake is a no-op for PascalCase,
// then the lookup lower-cases). No real cloud account is touched.
func paifsTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("paifeaturestore", endpoint)
	return client
}

func strPtr(s string) *string { return &s }

// TestUnitPaifsDatasourceIdentifiers covers the composite-id import contract
// (AC3): attributes win when present, the id "InstanceId:DatasourceId" is the
// fallback, and a malformed id is not silently split.
func TestUnitPaifsDatasourceIdentifiers(t *testing.T) {
	for _, tc := range []struct {
		name  string
		state datasourceResourceModel
		wantI string
		wantD string
	}{
		{"both attributes present", datasourceResourceModel{InstanceId: types.StringValue("fs-1"), DatasourceId: types.StringValue("ds-1")}, "fs-1", "ds-1"},
		{"fall back to composite id", datasourceResourceModel{Id: types.StringValue("fs-2:ds-2")}, "fs-2", "ds-2"},
		{"attributes win over id", datasourceResourceModel{Id: types.StringValue("x:y"), InstanceId: types.StringValue("fs-3"), DatasourceId: types.StringValue("ds-3")}, "fs-3", "ds-3"},
		{"empty id and attributes", datasourceResourceModel{}, "", ""},
		{"id with extra colons is not split", datasourceResourceModel{Id: types.StringValue("a:b:c")}, "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotI, gotD := datasourceIdentifiers(tc.state)
			if gotI != tc.wantI || gotD != tc.wantD {
				t.Fatalf("datasourceIdentifiers got %q:%q want %q:%q", gotI, gotD, tc.wantI, tc.wantD)
			}
		})
	}
}

// TestUnitPaifsStateMapping covers Read-back state mapping: the composite id,
// the carried InstanceId/DatasourceId, field-to-attribute mapping and empty→null.
func TestUnitPaifsStateMapping(t *testing.T) {
	obj := map[string]interface{}{
		"Type": "FeatureDB", "Name": "myds", "WorkspaceId": "ws-1",
		"Uri": "featuredb://x", "Config": `{"k":"v"}`, "GmtCreateTime": "2024-01-02",
	}
	m := datasourceState(obj, "fs-1", "ds-1")
	if m.Id.ValueString() != "fs-1:ds-1" {
		t.Fatalf("Id = %q, want fs-1:ds-1", m.Id.ValueString())
	}
	if m.InstanceId.ValueString() != "fs-1" {
		t.Errorf("InstanceId = %q, want fs-1", m.InstanceId.ValueString())
	}
	if m.DatasourceId.ValueString() != "ds-1" {
		t.Errorf("DatasourceId = %q, want ds-1", m.DatasourceId.ValueString())
	}
	if m.Type.ValueString() != "FeatureDB" {
		t.Errorf("Type = %q, want FeatureDB", m.Type.ValueString())
	}
	if m.DatasourceName.ValueString() != "myds" {
		t.Errorf("DatasourceName = %q, want myds", m.DatasourceName.ValueString())
	}
	if m.WorkspaceId.ValueString() != "ws-1" {
		t.Errorf("WorkspaceId = %q, want ws-1", m.WorkspaceId.ValueString())
	}
	if m.Uri.ValueString() != "featuredb://x" {
		t.Errorf("Uri = %q, want featuredb://x", m.Uri.ValueString())
	}
	if m.Config.ValueString() != `{"k":"v"}` {
		t.Errorf("Config = %q, want {\"k\":\"v\"}", m.Config.ValueString())
	}
	if m.CreateTime.ValueString() != "2024-01-02" {
		t.Errorf("CreateTime = %q, want 2024-01-02", m.CreateTime.ValueString())
	}

	// Empty values become null so drift detection treats absent server values as
	// removed rather than "".
	empty := datasourceState(map[string]interface{}{}, "fs-1", "ds-1")
	for attr, v := range map[string]types.String{
		"Type": empty.Type, "DatasourceName": empty.DatasourceName,
		"WorkspaceId": empty.WorkspaceId, "Uri": empty.Uri,
		"Config": empty.Config, "CreateTime": empty.CreateTime,
	} {
		if !v.IsNull() {
			t.Errorf("%s must be null for an empty object, got %q", attr, v.ValueString())
		}
	}
}

// TestUnitPaifsIsNotExist covers the failure/no-op detection used by Read and
// Delete: a matching tea error and a wrapped form are recognized, while a nil
// error and an unrelated error are not.
func TestUnitPaifsIsNotExist(t *testing.T) {
	if isDatasourceNotExist(nil) {
		t.Fatal("nil error must not be treated as not-exist")
	}
	if !isDatasourceNotExist(&tea.SDKError{Code: strPtr(datasourceNotExist)}) {
		t.Fatal("matching tea error code must be recognized")
	}
	if !isDatasourceNotExist(errors.New("wrapped: " + datasourceNotExist)) {
		t.Fatal("substring fallback must recognize a wrapped error")
	}
	if isDatasourceNotExist(errors.New("SomeUnrelatedError")) {
		t.Fatal("unrelated error must not be treated as not-exist")
	}
}

// TestUnitPaifsToInt covers the loose POP number coercion used for TotalCount
// across float64, int and string encodings.
func TestUnitPaifsToInt(t *testing.T) {
	for _, tc := range []struct {
		name string
		v    interface{}
		want int
	}{
		{"float64", float64(5), 5},
		{"int", int(7), 7},
		{"numeric string", "9", 9},
		{"non-numeric string", "x", 0},
		{"nil", nil, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := toInt(tc.v); got != tc.want {
				t.Fatalf("toInt(%v) = %d, want %d", tc.v, got, tc.want)
			}
		})
	}
}

// TestUnitPaifsGetDatasourceReadFound covers the Read success path: a wrapped
// Datasource object is unwrapped and returned, with GetDatasource called once.
func TestUnitPaifsGetDatasourceReadFound(t *testing.T) {
	var getDatasourceCalls int32
	client := paifsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Form.Get("Action") != "GetDatasource" {
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Form.Get("InstanceId") != "fs-1" || r.Form.Get("DatasourceId") != "ds-1" {
			t.Errorf("unexpected identifiers: instance=%q datasource=%q", r.Form.Get("InstanceId"), r.Form.Get("DatasourceId"))
		}
		atomic.AddInt32(&getDatasourceCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Datasource": map[string]interface{}{
				"DatasourceId": "ds-1", "Type": "FeatureDB", "Name": "myds",
				"WorkspaceId": "ws-1", "Uri": "featuredb://x", "GmtCreateTime": "2024-01-02",
			},
		})
	})
	obj, err := getDatasource(client, "fs-1", "ds-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if obj == nil {
		t.Fatal("expected the datasource object")
	}
	if obj["DatasourceId"] != "ds-1" {
		t.Fatalf("DatasourceId = %v, want ds-1", obj["DatasourceId"])
	}
	if obj["Type"] != "FeatureDB" {
		t.Errorf("Type = %v, want FeatureDB", obj["Type"])
	}
	if got := atomic.LoadInt32(&getDatasourceCalls); got != 1 {
		t.Fatalf("GetDatasource called %d times, want 1", got)
	}
}

// TestUnitPaifsGetDatasourceReadRoot covers the rootMapping "$" fallback: when
// the API returns the fields at the root instead of under "Datasource", the
// object is still resolved.
func TestUnitPaifsGetDatasourceReadRoot(t *testing.T) {
	client := paifsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"DatasourceId": "ds-1", "Type": "FeatureDB", "Name": "myds",
		})
	})
	obj, err := getDatasource(client, "fs-1", "ds-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if obj == nil || obj["DatasourceId"] != "ds-1" {
		t.Fatalf("expected root-mapped object with DatasourceId=ds-1, got %v", obj)
	}
}

// TestUnitPaifsGetDatasourceNotFound covers the no-op path: a not-found error
// is mapped to (nil, nil) so Read removes the resource and Delete stays
// idempotent, never surfacing a hard failure for a missing datasource.
func TestUnitPaifsGetDatasourceNotFound(t *testing.T) {
	client := paifsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"Code": datasourceNotExist, "Message": "datasource not found"})
	})
	obj, err := getDatasource(client, "fs-1", "missing")
	if err != nil {
		t.Fatalf("not-found must return a nil error, got %v", err)
	}
	if obj != nil {
		t.Fatalf("not-found must return a nil object, got %v", obj)
	}
}

// TestUnitPaifsGetDatasourceOtherError covers the failure path: an unrelated API
// error is surfaced, not swallowed into a silent no-op.
func TestUnitPaifsGetDatasourceOtherError(t *testing.T) {
	client := paifsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"Code": "Throttling.OtherCase", "Message": "unrelated failure"})
	})
	obj, err := getDatasource(client, "fs-1", "ds-1")
	if err == nil {
		t.Fatal("expected a non-nil error for an unrelated API failure")
	}
	if obj != nil {
		t.Fatalf("expected a nil object on error, got %v", obj)
	}
}
