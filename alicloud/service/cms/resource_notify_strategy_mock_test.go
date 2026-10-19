package cms

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// These are cloudless protocol-regression tests for the Framework-native
// alicloud_cms_notify_strategy resource. They stand up an httptest.Server that
// speaks the CMS ROA JSON contract, inject it into an AliyunClient via the
// fwadapt Configure hook, and drive Create/Read/Update/Delete directly through
// the Framework resource RPC request/response types. They cover the real CRUD
// path, state flattening, the failure path (5xx on Create, 404 tolerated on
// Read and Delete), and assert the actual HTTP method/path/body that the ROA
// client emitted — not helper internals.

// cmsNotifyStrategyTestClient mirrors the established cloud-firewall mock
// pattern: a local httptest server whose host:port is stored as the "cms"
// endpoint, so the ROA client (Protocol="http") dials the mock instead of the
// real cloud. NO_PROXY keeps localhost traffic off any proxy.
func cmsNotifyStrategyTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("credential: %v", err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("config.Client: %v", err)
	}
	endpoints.Store("cms", endpoint)
	return client
}

// notifyStrategyResourceSchema returns the resource schema and its top-level
// attribute-type map (attributes + blocks), used to build complete tftypes
// objects for Plan/State.
func notifyStrategyResourceSchema(t *testing.T) (schema.Schema, map[string]attr.Type) {
	t.Helper()
	r := &notifyStrategyResource{}
	resp := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("resource Schema diagnostics: %v", resp.Diagnostics)
	}
	objType, ok := resp.Schema.Type().(types.ObjectType)
	if !ok {
		t.Fatalf("schema Type() = %T, want types.ObjectType", resp.Schema.Type())
	}
	return resp.Schema, objType.AttrTypes
}

// nullValueFor returns the null attr.Value for the given attribute type, used
// to populate every unspecified field so ObjectValue receives a complete map.
func nullValueFor(t attr.Type) attr.Value {
	switch x := t.(type) {
	case basetypes.StringType:
		return types.StringNull()
	case basetypes.BoolType:
		return types.BoolNull()
	case basetypes.Int64Type:
		return types.Int64Null()
	case basetypes.ObjectType:
		return types.ObjectNull(x.AttrTypes)
	case basetypes.ListType:
		return types.ListNull(x.ElemType)
	default:
		return types.StringNull()
	}
}

// notifyStrategyTFTypes builds a complete tftypes.Value for the resource
// schema: every attribute/block defaults to null, overridden by the given map.
func notifyStrategyTFTypes(t *testing.T, sch schema.Schema, attrTypes map[string]attr.Type, overrides map[string]attr.Value) tftypes.Value {
	t.Helper()
	vals := make(map[string]attr.Value, len(attrTypes))
	for k, typ := range attrTypes {
		if v, ok := overrides[k]; ok {
			vals[k] = v
		} else {
			vals[k] = nullValueFor(typ)
		}
	}
	obj, diags := types.ObjectValue(attrTypes, vals)
	if diags.HasError() {
		t.Fatalf("ObjectValue diagnostics: %v", diags)
	}
	tfVal, err := obj.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("ToTerraformValue: %v", err)
	}
	return tfVal
}

// notifyStrategyClient spins up a resource bound to the mock client.
func notifyStrategyClient(t *testing.T, handler http.HandlerFunc) (*notifyStrategyResource, schema.Schema, map[string]attr.Type) {
	t.Helper()
	sch, attrTypes := notifyStrategyResourceSchema(t)
	client := cmsNotifyStrategyTestClient(t, handler)
	r := &notifyStrategyResource{}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, &resource.ConfigureResponse{})
	if r.Client() == nil {
		t.Fatal("client not injected after Configure")
	}
	return r, sch, attrTypes
}

// recordedReq captures a single HTTP request the ROA client emitted.
type recordedReq struct {
	Method string
	Path   string
	Body   string
}

// hasRequest reports whether the recorder saw a request matching method and a
// path prefix (e.g. hasRequest(recs, "POST", "/notifyStrategies")).
func hasRequest(recs []recordedReq, method, pathPrefix string) bool {
	for _, r := range recs {
		if r.Method == method && strings.HasPrefix(r.Path, pathPrefix) {
			return true
		}
	}
	return false
}

// requestBody returns the body of the first recorded request matching method and
// path prefix, or "" if none.
func requestBody(recs []recordedReq, method, pathPrefix string) string {
	for _, r := range recs {
		if r.Method == method && strings.HasPrefix(r.Path, pathPrefix) {
			return r.Body
		}
	}
	return ""
}

// lastRequest returns the most recently recorded request, or the zero value.
func lastRequest(recs []recordedReq) recordedReq {
	if len(recs) == 0 {
		return recordedReq{}
	}
	return recs[len(recs)-1]
}

// cmsNotifyStrategyAPIObject is the canonical CMS GetNotifyStrategy response
// data object the mock returns; it exercises every flatten path.
func cmsNotifyStrategyAPIObject() map[string]interface{} {
	return map[string]interface{}{
		"notifyStrategyName":         "tf-test",
		"description":                "desc",
		"notifyStrategyId":           "123",
		"createTime":                 "2024-01-01 00:00:00",
		"updateTime":                 "2024-01-02 00:00:00",
		"userId":                     "u-1",
		"workspace":                  "w-1",
		"enabled":                    true,
		"ignoreRestoredNotification": false,
		"groupingSetting": map[string]interface{}{
			"groupingKeys": []interface{}{"alertName", "level"},
			"periodMin":    5,
			"silenceSec":   10,
			"times":        3,
		},
		"routes": []interface{}{
			map[string]interface{}{
				"severities": []interface{}{"critical"},
				"channels": []interface{}{
					map[string]interface{}{
						"channelType":        "dingTalk",
						"enabledSubChannels": []interface{}{"sub-1"},
						"receivers":          []interface{}{"r-1"},
					},
				},
				"effectTimeRange": map[string]interface{}{
					"timeZone":          "Asia/Shanghai",
					"dayInWeek":         []interface{}{int64(1), int64(2)},
					"startTimeInMinute": int64(0),
					"endTimeInMinute":   int64(1440),
				},
				"filterSetting": map[string]interface{}{
					"expression": "expr",
					"relation":   "and",
					"conditions": []interface{}{
						map[string]interface{}{"field": "f", "op": "eq", "value": "v"},
					},
				},
			},
		},
		"customTemplateEntries": []interface{}{
			map[string]interface{}{"targetType": "dingTalk", "templateUuid": "t-1"},
		},
	}
}

// lifecycleHandler responds per CMS ROA contract for a single strategy of id 123,
// appending every emitted request to recs so the test can assert ordering.
func lifecycleHandler(t *testing.T, recs *[]recordedReq, apiObject map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bodyBytes, _ := io.ReadAll(req.Body)
		*recs = append(*recs, recordedReq{Method: req.Method, Path: req.URL.Path, Body: string(bodyBytes)})
		w.Header().Set("Content-Type", "application/json")
		switch {
		case req.Method == http.MethodPost && req.URL.Path == "/notifyStrategies":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": "123"})
		case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/notifyStrategies/"):
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": apiObject})
		case req.Method == http.MethodPut && strings.HasPrefix(req.URL.Path, "/notifyStrategies/"):
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		case req.Method == http.MethodDelete && strings.HasPrefix(req.URL.Path, "/notifyStrategies/"):
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		default:
			http.NotFound(w, req)
		}
	}
}

func TestUnitCmsNotifyStrategyLifecycle(t *testing.T) {
	ctx := context.Background()
	apiObject := cmsNotifyStrategyAPIObject()
	var recs []recordedReq
	r, sch, attrTypes := notifyStrategyClient(t, lifecycleHandler(t, &recs, apiObject))

	// ---- Create: POST /notifyStrategies, then GET to flatten state ----
	createOverrides := map[string]attr.Value{
		"notify_strategy_name":    types.StringValue("tf-test"),
		"description":             types.StringValue("desc"),
		"workspace":               types.StringValue("w-1"),
		"grouping_setting":        flattenGroupingSetting(apiObject["groupingSetting"]),
		"routes":                  flattenRoutes(apiObject["routes"]),
		"custom_template_entries": flattenCustomTemplateEntries(apiObject["customTemplateEntries"]),
	}
	planVal := notifyStrategyTFTypes(t, sch, attrTypes, createOverrides)
	createResp := resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{
		Plan:   tfsdk.Plan{Raw: planVal, Schema: sch},
		Config: tfsdk.Config{Raw: planVal, Schema: sch},
	}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create diagnostics: %v", createResp.Diagnostics)
	}
	if !hasRequest(recs, http.MethodPost, "/notifyStrategies") {
		t.Fatalf("Create request not emitted; recorded=%+v", recs)
	}
	createBody := requestBody(recs, http.MethodPost, "/notifyStrategies")
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(createBody), &bodyMap); err != nil {
		t.Fatalf("Create body not JSON: %v (body=%q)", err, createBody)
	}
	if bodyMap["name"] != "tf-test" {
		t.Fatalf("Create body name=%v, want tf-test", bodyMap["name"])
	}
	if _, ok := bodyMap["groupingSetting"]; !ok {
		t.Fatalf("Create body missing groupingSetting: %v", bodyMap)
	}
	if _, ok := bodyMap["routes"]; !ok {
		t.Fatalf("Create body missing routes: %v", bodyMap)
	}
	var created notifyStrategyResourceModel
	if diags := createResp.State.Get(ctx, &created); diags.HasError() {
		t.Fatalf("Create State.Get: %v", diags)
	}
	if created.Id.ValueString() != "123" {
		t.Fatalf("created.Id=%q want 123", created.Id)
	}
	if created.NotifyStrategyId.ValueString() != "123" {
		t.Fatalf("created.NotifyStrategyId=%q want 123", created.NotifyStrategyId)
	}
	if created.NotifyStrategyName.ValueString() != "tf-test" {
		t.Fatalf("created.NotifyStrategyName=%q want tf-test", created.NotifyStrategyName)
	}
	if created.Workspace.ValueString() != "w-1" {
		t.Fatalf("created.Workspace=%q want w-1", created.Workspace)
	}
	if created.UserId.ValueString() != "u-1" {
		t.Fatalf("created.UserId=%q want u-1", created.UserId)
	}
	if !created.Enable.ValueBool() {
		t.Fatalf("created.Enable=%v want true", created.Enable)
	}
	if created.GroupingSetting.IsNull() {
		t.Fatal("created.GroupingSetting is null")
	}
	if created.Routes.IsNull() || len(created.Routes.Elements()) != 1 {
		t.Fatalf("created.Routes not populated: %+v", created.Routes)
	}
	if created.CustomTemplateEntries.IsNull() || len(created.CustomTemplateEntries.Elements()) != 1 {
		t.Fatalf("created.CustomTemplateEntries not populated: %+v", created.CustomTemplateEntries)
	}

	// ---- Read: GET /notifyStrategies/123 ----
	stateOverrides := map[string]attr.Value{
		"id":                 types.StringValue("123"),
		"notify_strategy_id": types.StringValue("123"),
		"workspace":          types.StringValue("w-1"),
	}
	stateVal := notifyStrategyTFTypes(t, sch, attrTypes, stateOverrides)
	readResp := resource.ReadResponse{State: tfsdk.State{Schema: sch}}
	r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Raw: stateVal, Schema: sch}}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read diagnostics: %v", readResp.Diagnostics)
	}
	last := lastRequest(recs)
	if last.Method != http.MethodGet || last.Path != "/notifyStrategies/123" {
		t.Fatalf("Read request: got %s %s, want GET /notifyStrategies/123", last.Method, last.Path)
	}
	var readModel notifyStrategyResourceModel
	if diags := readResp.State.Get(ctx, &readModel); diags.HasError() {
		t.Fatalf("Read State.Get: %v", diags)
	}
	if readModel.Id.ValueString() != "123" {
		t.Fatalf("read Id=%q want 123 (must be preserved)", readModel.Id)
	}
	if readModel.UserId.ValueString() != "u-1" {
		t.Fatalf("read UserId=%q want u-1", readModel.UserId)
	}
	if !readModel.Enable.ValueBool() {
		t.Fatalf("read Enable=false want true")
	}
	if readModel.Routes.IsNull() || len(readModel.Routes.Elements()) != 1 {
		t.Fatalf("read Routes not populated: %+v", readModel.Routes)
	}
	// grouping_setting round-trips with grouping_keys of length 2
	gsAttrs := readModel.GroupingSetting.Attributes()
	if gk, ok := gsAttrs["grouping_keys"].(types.List); !ok || gk.IsNull() || len(gk.Elements()) != 2 {
		t.Fatalf("read grouping_keys not round-tripped: %+v", gsAttrs["grouping_keys"])
	}

	// ---- Update: PUT /notifyStrategies/123 ----
	updateOverrides := map[string]attr.Value{
		"notify_strategy_name": types.StringValue("tf-test-updated"),
		"description":          types.StringValue("desc2"),
		"workspace":            types.StringValue("w-1"),
		"grouping_setting":     flattenGroupingSetting(apiObject["groupingSetting"]),
	}
	updatePlanVal := notifyStrategyTFTypes(t, sch, attrTypes, updateOverrides)
	updateResp := resource.UpdateResponse{State: tfsdk.State{Schema: sch}}
	r.Update(ctx, resource.UpdateRequest{
		Plan:   tfsdk.Plan{Raw: updatePlanVal, Schema: sch},
		State:  tfsdk.State{Raw: stateVal, Schema: sch},
		Config: tfsdk.Config{Raw: updatePlanVal, Schema: sch},
	}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update diagnostics: %v", updateResp.Diagnostics)
	}
	if !hasRequest(recs, http.MethodPut, "/notifyStrategies/123") {
		t.Fatalf("Update request not emitted; recorded=%+v", recs)
	}
	updateBody := requestBody(recs, http.MethodPut, "/notifyStrategies/123")
	var updateBodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(updateBody), &updateBodyMap); err != nil {
		t.Fatalf("Update body not JSON: %v (body=%q)", err, updateBody)
	}
	if updateBodyMap["name"] != "tf-test-updated" {
		t.Fatalf("Update body name=%v want tf-test-updated", updateBodyMap["name"])
	}

	// ---- Delete: DELETE /notifyStrategies/123 ----
	deleteResp := resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: tfsdk.State{Raw: stateVal, Schema: sch}}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete diagnostics: %v", deleteResp.Diagnostics)
	}
	last = lastRequest(recs)
	if last.Method != http.MethodDelete || last.Path != "/notifyStrategies/123" {
		t.Fatalf("Delete request: got %s %s, want DELETE /notifyStrategies/123", last.Method, last.Path)
	}
}

// TestUnitCmsNotifyStrategyCreateError drives the failure path: a 5xx from
// CreateNotifyStrategy must surface as a Create diagnostic, not a silent pass.
func TestUnitCmsNotifyStrategyCreateError(t *testing.T) {
	ctx := context.Background()
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if req.Method == http.MethodPost && req.URL.Path == "/notifyStrategies" {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "InternalError", "Message": "boom"})
			return
		}
		http.NotFound(w, req)
	}
	r, sch, attrTypes := notifyStrategyClient(t, handler)
	planVal := notifyStrategyTFTypes(t, sch, attrTypes, map[string]attr.Value{
		"notify_strategy_name": types.StringValue("tf-test"),
		"workspace":            types.StringValue("w-1"),
	})
	createResp := resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{
		Plan:   tfsdk.Plan{Raw: planVal, Schema: sch},
		Config: tfsdk.Config{Raw: planVal, Schema: sch},
	}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create on HTTP 500 must report diagnostics, got none")
	}
}

// TestUnitCmsNotifyStrategyReadNotFound verifies a 404 on GetNotifyStrategy is
// tolerated and removes the resource from state rather than erroring.
func TestUnitCmsNotifyStrategyReadNotFound(t *testing.T) {
	ctx := context.Background()
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/notifyStrategies/") {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "ResourceNotFound", "Message": "not found"})
			return
		}
		http.NotFound(w, req)
	}
	r, sch, attrTypes := notifyStrategyClient(t, handler)
	stateVal := notifyStrategyTFTypes(t, sch, attrTypes, map[string]attr.Value{
		"id":                 types.StringValue("999"),
		"notify_strategy_id": types.StringValue("999"),
		"workspace":          types.StringValue("w-1"),
	})
	readResp := resource.ReadResponse{State: tfsdk.State{Schema: sch}}
	r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Raw: stateVal, Schema: sch}}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read on 404 must be tolerated, got: %v", readResp.Diagnostics)
	}
}

// TestUnitCmsNotifyStrategyDeleteNotFound verifies a 404 on DeleteNotifyStrategy is
// tolerated (idempotent destroy) rather than surfacing an error.
func TestUnitCmsNotifyStrategyDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if req.Method == http.MethodDelete && strings.HasPrefix(req.URL.Path, "/notifyStrategies/") {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "ResourceNotFound", "Message": "not found"})
			return
		}
		http.NotFound(w, req)
	}
	r, sch, attrTypes := notifyStrategyClient(t, handler)
	stateVal := notifyStrategyTFTypes(t, sch, attrTypes, map[string]attr.Value{
		"id":                 types.StringValue("999"),
		"notify_strategy_id": types.StringValue("999"),
		"workspace":          types.StringValue("w-1"),
	})
	deleteResp := resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: tfsdk.State{Raw: stateVal, Schema: sch}}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete on 404 must be tolerated, got: %v", deleteResp.Diagnostics)
	}
}

// Compile-time guard: the resource carries the injected client, so the mock
// test cannot drift away from the fwadapt contract.
var _ fwadapt.WithClient = (*notifyStrategyResource)(nil)
