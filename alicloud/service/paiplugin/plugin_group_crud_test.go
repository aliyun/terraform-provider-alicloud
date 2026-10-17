package paiplugin

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnitPluginGroupFrameworkCRUDMatrix exercises F4 (AC-R4): it drives the
// plugin-framework pluginGroupResource.Create/Read/Delete/ImportState against
// an httptest server, asserting the state-path contracts — NOT the helper
// functions, which plugin_group_unit_test.go already covers. No cloud is
// touched: a fake access_key credential signs locally and the endpoint map is
// pre-seeded so Config.Client() (SkipRegionValidation) builds a client whose
// ROA path resolves to the test server.
//
// The matrix maps directly to AC-R4:
//   - Create → Get 500 / genuine-not-found (暂不可见) / 错形(missing Id) all
//     retain the assigned id in state; Create success refreshes to the GetGroup
//     body.
//   - Read genuine NotFound removes the resource; empty-code 404, unknown-code
//     404, InvalidApi.NotFound (routing) 404, 500, 200-with-no-Data, 200 with
//     missing/wrong-type/wrong-groupId Id all preserve the prior state + diag.
//   - Delete error preserves state; genuine NotFound is idempotent (no diag).
//   - ID-only import → Read produces a legal state (no diag, id preserved).
//
// The ROA SDK surfaces HTTP 4xx/5xx as *tea.SDKError{StatusCode, Code} parsed
// from the JSON body (Code/Message/RequestId); formatError returns it raw, so
// isPluginGroupNotExist sees the SDK error directly. CreateGroup/DeleteGroup are
// called with autoRetry=true, but a 200 is a single request and a 5xx is capped
// at maxAttempts=3 with no backoff sleep (policy "no") — never a hang.
func TestUnitPluginGroupFrameworkCRUDMatrix(t *testing.T) {
	ctx := context.Background()

	// Canned GetGroup/DeleteGroup response bodies keyed by scenario label. A 200
	// with a Data object carrying Id=123 is the "healthy group" shape.
	const (
		healthyBody    = `{"Data":{"Id":123,"Name":"g1","Source":0,"Remark":"terraform-example"}}`
		noDataBody     = `{"Message":"ok","RequestId":"r1"}`
		missingIdBody  = `{"Data":{"Name":"g1"}}`
		wrongIdBody    = `{"Data":{"Id":456,"Name":"g1"}}`
		wrongTypeBody  = `{"Data":{"Id":true,"Name":"g1"}}`
		notFoundBody   = `{"Code":"NotFound","Message":"group not found","RequestId":"r1"}`
		emptyCodeBody  = `{"Message":"err","RequestId":"r1"}`
		unknownCodeBdy = `{"Code":"SomeUnknownError","Message":"x","RequestId":"r1"}`
		invalidApiBdy  = `{"Code":"InvalidApi.NotFound","Message":"specified api is not found","RequestId":"r1"}`
		internalBody   = `{"Code":"InternalError","Message":"internal","RequestId":"r1"}`
	)

	// groupSrv returns an httptest server whose CreateGroup always succeeds
	// (Data.Id=123) and whose GetGroup/DeleteGroup behavior is driven by get/del.
	// nil get/del means "healthy 200 Data.Id=123". It also records the wire calls.
	groupSrv := func(t *testing.T, get, del func() (int, string)) (*httptest.Server, *[]string) {
		var mu sync.Mutex
		calls := []string{}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			action := r.Header.Get("x-acs-action")
			mu.Lock()
			calls = append(calls, r.Method+" "+r.URL.EscapedPath()+" action="+action)
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			switch action {
			case "CreateGroup":
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{"Data":{"Id":123,"Name":"g1"}}`)
			case "GetGroup":
				code, body := http.StatusOK, healthyBody
				if get != nil {
					code, body = get()
				}
				w.WriteHeader(code)
				_, _ = io.WriteString(w, body)
			case "DeleteGroup":
				code, body := http.StatusOK, `{}`
				if del != nil {
					code, body = del()
				}
				w.WriteHeader(code)
				_, _ = io.WriteString(w, body)
			default:
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{}`)
			}
		}))
		t.Cleanup(srv.Close)
		return srv, &calls
	}

	// newResource builds the resource, injects a client wired to srv, and returns
	// the concrete *pluginGroupResource (so the promoted Client() is callable for
	// the injection check; resource.Resource is an interface without Client()).
	// planFromModel/stateFromModel still accept resource.Resource.
	newResource := func(t *testing.T, srv *httptest.Server) *pluginGroupResource {
		client := newPluginGroupTestClient(t, srv)
		r := NewPluginGroupResource().(*pluginGroupResource)
		var cfgResp resource.ConfigureResponse
		r.Configure(ctx, resource.ConfigureRequest{ProviderData: client}, &cfgResp)
		require.False(t, cfgResp.Diagnostics.HasError(), "configure: %v", cfgResp.Diagnostics)
		require.NotNil(t, r.Client(), "client not injected")
		var sr resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &sr)
		require.False(t, sr.Diagnostics.HasError(), "schema: %v", sr.Diagnostics)
		return r
	}

	// planFromModel builds a tfsdk.Plan (Raw populated) from a model struct.
	planFromModel := func(t *testing.T, r resource.Resource, m pluginGroupResourceModel) tfsdk.Plan {
		var sr resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &sr)
		p := tfsdk.Plan{Schema: sr.Schema}
		require.False(t, p.Set(ctx, &m).HasError(), "plan.Set")
		return p
	}

	// stateFromModel builds a tfsdk.State (Raw populated) from a model struct.
	stateFromModel := func(t *testing.T, r resource.Resource, m pluginGroupResourceModel) tfsdk.State {
		var sr resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &sr)
		s := tfsdk.State{Schema: sr.Schema}
		require.False(t, s.Set(ctx, &m).HasError(), "state.Set")
		return s
	}

	// decodeState decodes a tfsdk.State back into the model for assertions.
	decodeState := func(t *testing.T, s tfsdk.State) pluginGroupResourceModel {
		var m pluginGroupResourceModel
		require.False(t, s.Get(ctx, &m).HasError(), "state.Get")
		return m
	}

	newPlanModel := func() pluginGroupResourceModel {
		// Every field is explicitly initialized (null for computed/unset) so
		// reflect.FromValue never sees a zero-value attr.Value.
		return pluginGroupResourceModel{
			Id:             types.StringNull(),
			PluginGroupId:  types.StringNull(),
			Name:           types.StringValue("g1"),
			Source:         types.Int64Value(0),
			Project:        types.StringNull(),
			Algorithm:      types.StringNull(),
			Filter:         types.StringNull(),
			Column:         types.StringNull(),
			Table:          types.StringNull(),
			Text:           types.StringValue("13900000000"),
			Uri:            types.StringNull(),
			Remark:         types.StringValue("terraform-example"),
			PhoneNumber:    types.BoolNull(),
			Amount:         types.Int64Null(),
			CampaignId:     types.StringNull(),
			CreatedTime:    types.StringNull(),
			UpdatedTime:    types.StringNull(),
			Status:         types.Int64Null(),
			InferenceJobId: types.StringNull(),
			History:        types.StringNull(),
		}
	}

	newStateModel := func(id string) pluginGroupResourceModel {
		m := newPlanModel()
		m.Id = types.StringValue(id)
		m.PluginGroupId = types.StringValue(id)
		return m
	}

	// --- Create matrix (AC-R4: Create→Get 500/暂不可见/错形 保留 ID) ---

	t.Run("Create success refreshes to GetGroup body", func(t *testing.T) {
		srv, calls := groupSrv(t, nil, nil)
		r := newResource(t, srv)
		plan := planFromModel(t, r, newPlanModel())
		resp := resource.CreateResponse{State: tfsdk.State{Schema: plan.Schema}}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
		assert.False(t, resp.Diagnostics.HasError(), "unexpected diags: %v", resp.Diagnostics)
		assert.False(t, resp.State.Raw.IsNull(), "state must be set on success")
		m := decodeState(t, resp.State)
		assert.Equal(t, "123", m.Id.ValueString(), "refreshed id from GetGroup Data.Id")
		assert.Equal(t, "123", m.PluginGroupId.ValueString())
		assert.Equal(t, "g1", m.Name.ValueString())
		require.NotEmpty(t, *calls, "expected Create + Get wire calls")
		assert.Equal(t, "CreateGroup", wireAction((*calls)[0]))
		assert.Equal(t, "GetGroup", wireAction((*calls)[1]))
	})

	t.Run("Create with GetGroup 500 retains assigned id + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusInternalServerError, internalBody }, nil)
		r := newResource(t, srv)
		plan := planFromModel(t, r, newPlanModel())
		resp := resource.CreateResponse{State: tfsdk.State{Schema: plan.Schema}}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "GetGroup 500 must surface a diag")
		assert.False(t, resp.State.Raw.IsNull(), "assigned id must be retained, not dropped")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "persisted create-state id retained")
	})

	t.Run("Create with GetGroup genuine-not-found (暂不可见) retains assigned id + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusNotFound, notFoundBody }, nil)
		r := newResource(t, srv)
		plan := planFromModel(t, r, newPlanModel())
		resp := resource.CreateResponse{State: tfsdk.State{Schema: plan.Schema}}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "not-visible-after-create must surface a diag (not silently remove)")
		assert.False(t, resp.State.Raw.IsNull(), "assigned id must be retained")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "persisted create-state id retained")
	})

	t.Run("Create with GetGroup 错形 (missing Id) retains assigned id + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusOK, missingIdBody }, nil)
		r := newResource(t, srv)
		plan := planFromModel(t, r, newPlanModel())
		resp := resource.CreateResponse{State: tfsdk.State{Schema: plan.Schema}}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "malformed GetGroup (no Id) must surface a diag (F1/F3)")
		assert.False(t, resp.State.Raw.IsNull(), "assigned id must be retained")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "persisted create-state id retained")
	})

	// --- Read matrix (AC-R4) ---

	t.Run("Read genuine NotFound removes the resource", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusNotFound, notFoundBody }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in} // pre-populated from req.State
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.False(t, resp.Diagnostics.HasError(), "genuine NotFound must not diag: %v", resp.Diagnostics)
		assert.True(t, resp.State.Raw.IsNull(), "genuine NotFound must RemoveResource")
	})

	t.Run("Read empty-code 404 preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusNotFound, emptyCodeBody }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "empty-code 404 must NOT be treated as not-found (F2)")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved, not removed")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained")
	})

	t.Run("Read unknown-code 404 preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusNotFound, unknownCodeBdy }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "unknown-code 404 must NOT be treated as not-found (F2)")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained")
	})

	t.Run("Read InvalidApi.NotFound (routing) 404 preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusNotFound, invalidApiBdy }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "routing 404 must NOT be treated as not-found (F2)")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained")
	})

	t.Run("Read 500 preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusInternalServerError, internalBody }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "500 must NOT remove the resource")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained")
	})

	t.Run("Read 200 missing Data preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusOK, noDataBody }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "200 without Data must diag (F3)")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained")
	})

	t.Run("Read 200 missing Id (错形) preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusOK, missingIdBody }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "200 with Data but no Id must diag (F1/F3)")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained")
	})

	t.Run("Read 200 wrong-type Id (错形) preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusOK, wrongTypeBody }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "200 with bool Id must diag (F1 toStr/F3)")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained")
	})

	t.Run("Read 200 wrong groupId (跨资源) preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, func() (int, string) { return http.StatusOK, wrongIdBody }, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "200 whose Data.Id != requested groupId must diag (F1 cross-resource takeover)")
		assert.False(t, resp.State.Raw.IsNull(), "state must be preserved (not adopt the foreign group)")
		assert.Equal(t, "123", decodeState(t, resp.State).Id.ValueString(), "prior state retained, foreign id NOT adopted")
	})

	t.Run("Read 200 healthy refreshes state", func(t *testing.T) {
		srv, _ := groupSrv(t, nil, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.ReadResponse{State: in}
		r.Read(ctx, resource.ReadRequest{State: in}, &resp)
		assert.False(t, resp.Diagnostics.HasError(), "healthy read: %v", resp.Diagnostics)
		m := decodeState(t, resp.State)
		assert.Equal(t, "123", m.Id.ValueString())
		assert.Equal(t, "g1", m.Name.ValueString())
	})

	// --- Delete matrix (AC-R4: Delete 错误保留,真 NotFound 幂等) ---

	t.Run("Delete 500 preserves state + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, nil, func() (int, string) { return http.StatusInternalServerError, internalBody })
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.DeleteResponse{State: in}
		r.Delete(ctx, resource.DeleteRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "Delete 500 must diag and leave the resource in place")
		assert.False(t, resp.State.Raw.IsNull(), "state preserved on delete error")
	})

	t.Run("Delete genuine NotFound is idempotent (no diag)", func(t *testing.T) {
		srv, _ := groupSrv(t, nil, func() (int, string) { return http.StatusNotFound, notFoundBody })
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.DeleteResponse{State: in}
		r.Delete(ctx, resource.DeleteRequest{State: in}, &resp)
		assert.False(t, resp.Diagnostics.HasError(), "genuine NotFound delete must be idempotent: %v", resp.Diagnostics)
	})

	t.Run("Delete routing 404 (InvalidApi.NotFound) preserves + diag", func(t *testing.T) {
		srv, _ := groupSrv(t, nil, func() (int, string) { return http.StatusNotFound, invalidApiBdy })
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.DeleteResponse{State: in}
		r.Delete(ctx, resource.DeleteRequest{State: in}, &resp)
		assert.True(t, resp.Diagnostics.HasError(), "routing 404 on delete must NOT be idempotent (F2)")
	})

	t.Run("Delete success has no diag", func(t *testing.T) {
		srv, _ := groupSrv(t, nil, nil)
		r := newResource(t, srv)
		in := stateFromModel(t, r, newStateModel("123"))
		resp := resource.DeleteResponse{State: in}
		r.Delete(ctx, resource.DeleteRequest{State: in}, &resp)
		assert.False(t, resp.Diagnostics.HasError(), "delete success: %v", resp.Diagnostics)
	})

	// --- ImportState matrix (AC-R4: ID-only import → Read 合法 state) ---

	t.Run("ID-only import then Read yields legal state", func(t *testing.T) {
		srv, _ := groupSrv(t, nil, nil)
		r := newResource(t, srv)
		// The framework pre-populates ImportStateResponse.State with a typed
		// null object of the schema type before calling ImportState (so
		// ImportStatePassthroughID's SetAttribute can add the id attribute).
		// Mirror that here: a zero tftypes.Value has no type and SetAttribute
		// would fail, so seed Raw with a typed null object.
		var sr resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &sr)
		impResp := resource.ImportStateResponse{State: tfsdk.State{
			Schema: sr.Schema,
			Raw:    tftypes.NewValue(sr.Schema.Type().TerraformType(ctx), nil),
		}}
		r.ImportState(ctx, resource.ImportStateRequest{ID: "123"}, &impResp)
		assert.False(t, impResp.Diagnostics.HasError(), "import: %v", impResp.Diagnostics)
		assert.Equal(t, "123", decodeState(t, impResp.State).Id.ValueString(), "import sets id")

		// The framework would now refresh; drive Read with the imported state.
		readResp := resource.ReadResponse{State: impResp.State}
		r.Read(ctx, resource.ReadRequest{State: impResp.State}, &readResp)
		assert.False(t, readResp.Diagnostics.HasError(), "import refresh read: %v", readResp.Diagnostics)
		assert.False(t, readResp.State.Raw.IsNull(), "imported resource must remain in state")
		m := decodeState(t, readResp.State)
		assert.Equal(t, "123", m.Id.ValueString(), "legal state id after import refresh")
		assert.Equal(t, "g1", m.Name.ValueString())
	})
}

// wireAction extracts the x-acs-action value from a recorded "METHOD path action=X" call.
func wireAction(call string) string {
	if i := strings.Index(call, "action="); i >= 0 {
		return call[i+len("action="):]
	}
	return ""
}

// newPluginGroupTestClient builds a real connectivity.AliyunClient whose ROA path
// is wired to srv, using the exported Config.Client() path. SkipRegionValidation
// avoids region validation; the pre-seeded Endpoints map short-circuits
// loadApiEndpoint so the PAIPlugin ROA calls hit srv. No cloud is touched: the
// access_key credential signs locally.
func newPluginGroupTestClient(t *testing.T, srv *httptest.Server) *connectivity.AliyunClient {
	t.Helper()
	host := strings.TrimPrefix(srv.URL, "http://")
	endpoints := &sync.Map{}
	endpoints.Store("paiplugin", host)
	cred, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("FAKEPAIPLUGINAK").
		SetAccessKeySecret("FAKEPAIPLUGINSK"))
	require.NoError(t, err)
	cfg := &connectivity.Config{
		RegionId:             "cn-hangzhou",
		AccessKey:            "FAKEPAIPLUGINAK",
		SecretKey:            "FAKEPAIPLUGINSK",
		Protocol:             "HTTP",
		Endpoints:            endpoints,
		SignVersion:          &sync.Map{},
		Credential:           cred,
		SkipRegionValidation: true,
	}
	client, err := cfg.Client()
	require.NoError(t, err)
	require.NotNil(t, client)
	return client
}
