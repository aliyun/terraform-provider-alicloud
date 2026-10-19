package ecs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// sgSnapshotPolicyTestClient builds an AliyunClient whose Ecs RPC calls are
// redirected at an in-process httptest server with test credentials. No cloud
// account or network egress is involved: the server owns every response,
// including error and malformed envelopes. This mirrors the established
// ecsSnapshotTestClient pattern but lives in the ecs package so the resource
// under test (and its unexported model) are reachable.
func sgSnapshotPolicyTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
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
	endpoints.Store("ecs", endpoint)
	return client
}

// newSGSnapshotPolicyResource builds the resource with a configured httptest
// client so Create/Read/Update reach the in-process server.
func newSGSnapshotPolicyResource(t *testing.T, client *connectivity.AliyunClient) (*securityGroupSnapshotPolicyResource, fwschema.Schema) {
	ctx := context.Background()
	r := NewSecurityGroupSnapshotPolicyResource().(*securityGroupSnapshotPolicyResource)
	var cfg resource.ConfigureResponse
	r.Configure(ctx, resource.ConfigureRequest{ProviderData: client}, &cfg)
	if cfg.Diagnostics.HasError() {
		t.Fatalf("Configure: %v", cfg.Diagnostics)
	}
	if r.Client() == nil {
		t.Fatal("Configure did not inject the client")
	}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema: %v", schemaResp.Diagnostics)
	}
	return r, schemaResp.Schema
}

// popJSON writes a 2xx JSON POP response.
func popJSON(t *testing.T, w http.ResponseWriter, body map[string]interface{}) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Error(err)
	}
}

// popError writes a non-2xx POP error envelope that is not a NotFound (so the
// resource treats it as a real error rather than a confirmed absence).
func popError(t *testing.T, w http.ResponseWriter, code, message string) {
	t.Helper()
	w.WriteHeader(http.StatusForbidden)
	popJSON(t, w, map[string]interface{}{"Code": code, "Message": message})
}

// stateFromModel builds a tfsdk.State pre-populated from a model so Read can be
// asserted to preserve (or RemoveResource) that prior state.
func stateFromModel(t *testing.T, ctx context.Context, sch fwschema.Schema, model securityGroupSnapshotPolicyResourceModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: sch}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Fatalf("state.Set: %v", diags)
	}
	return state
}

// F1: Create must persist the recoverable id before the post-Create refresh. A
// Describe that errors after a successful Create leaves the id in resp.State
// with a real diagnostic, so the created cloud resource is not leaked.
func TestUnitSGSnapshotPolicyCreate_PersistsIDWhenDescribeErrors(t *testing.T) {
	ctx := context.Background()
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "CreateSecurityGroupSnapshotPolicy":
			popJSON(t, w, map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"})
		case "DescribeSecurityGroupSnapshotPolicies":
			// A real error after a successful create: not a NotFound, so the
			// resource keeps the id and surfaces the failure.
			popError(t, w, "Forbidden", "describe failed")
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(ctx, &securityGroupSnapshotPolicyResourceModel{
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
	}); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("Create must surface the Describe error as a diagnostic")
	}
	var out securityGroupSnapshotPolicyResourceModel
	if diags := resp.State.Get(ctx, &out); diags.HasError() {
		t.Fatalf("resp.State.Get: %v", diags)
	}
	if out.Id.ValueString() != "sp-1" || out.SecurityGroupSnapshotPolicyId.ValueString() != "sp-1" {
		t.Fatalf("the recovered id must be persisted in state before the failing refresh; got id=%q policy_id=%q",
			out.Id.ValueString(), out.SecurityGroupSnapshotPolicyId.ValueString())
	}
}

// F1 happy path: Create with a Describe that returns the created object
// produces a clean state with no diagnostics.
func TestUnitSGSnapshotPolicyCreate_HappyPath(t *testing.T) {
	ctx := context.Background()
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "CreateSecurityGroupSnapshotPolicy":
			popJSON(t, w, map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"})
		case "DescribeSecurityGroupSnapshotPolicies":
			popJSON(t, w, map[string]interface{}{
				"SecurityGroupSnapshotPolicySets": map[string]interface{}{
					"SecurityGroupSnapshotPolicy": []interface{}{
						map[string]interface{}{
							"SecurityGroupSnapshotPolicyId":   "sp-1",
							"SecurityGroupSnapshotPolicyName": "tf-policy",
							"Status":                          "enable",
							"RegionId":                        "cn-hangzhou",
						},
					},
				},
			})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(ctx, &securityGroupSnapshotPolicyResourceModel{
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
	}); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("happy-path Create must not error: %v", resp.Diagnostics)
	}
	var out securityGroupSnapshotPolicyResourceModel
	if diags := resp.State.Get(ctx, &out); diags.HasError() {
		t.Fatalf("resp.State.Get: %v", diags)
	}
	if out.Id.ValueString() != "sp-1" {
		t.Fatalf("expected id sp-1, got %q", out.Id.ValueString())
	}
	if out.SecurityGroupSnapshotPolicyName.ValueString() != "tf-policy" {
		t.Fatalf("expected name tf-policy, got %q", out.SecurityGroupSnapshotPolicyName.ValueString())
	}
}

// F2: a malformed Describe envelope (missing the list root) is an error, not a
// NotFound. Read must keep the prior state and surface a diagnostic instead of
// dropping the resource.
func TestUnitSGSnapshotPolicyRead_MalformedKeepsState(t *testing.T) {
	ctx := context.Background()
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeSecurityGroupSnapshotPolicies":
			// 2xx but no recognizable list root: malformed, not NotFound.
			popJSON(t, w, map[string]interface{}{"RequestId": "r-1"})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	prior := stateFromModel(t, ctx, sch, securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
		RegionId:                        types.StringValue("cn-hangzhou"),
	})
	resp := resource.ReadResponse{State: prior}
	r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("malformed response must surface a diagnostic, not be swallowed as NotFound")
	}
	var out securityGroupSnapshotPolicyResourceModel
	resp.State.Get(ctx, &out)
	if out.Id.ValueString() != "sp-1" {
		t.Fatalf("malformed response must not drop the resource; id still %q expected, got %q", "sp-1", out.Id.ValueString())
	}
}

// F2: a confirmed NotFound (a well-formed but empty list) drops the resource.
func TestUnitSGSnapshotPolicyRead_NotFoundRemovesResource(t *testing.T) {
	ctx := context.Background()
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeSecurityGroupSnapshotPolicies":
			popJSON(t, w, map[string]interface{}{
				"SecurityGroupSnapshotPolicySets": map[string]interface{}{
					"SecurityGroupSnapshotPolicy": []interface{}{},
				},
			})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	prior := stateFromModel(t, ctx, sch, securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
		RegionId:                        types.StringValue("cn-hangzhou"),
	})
	resp := resource.ReadResponse{State: prior}
	r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("a confirmed NotFound must not error: %v", resp.Diagnostics)
	}
	var out securityGroupSnapshotPolicyResourceModel
	resp.State.Get(ctx, &out)
	if !out.Id.IsNull() {
		t.Fatalf("a confirmed NotFound must remove the resource; id still %q", out.Id.ValueString())
	}
}

// F3: an Optional+Computed field left unset (unknown in plan) must not be sent
// as an empty string. Update with name unknown and status changed sends Status
// but not SecurityGroupSnapshotPolicyName.
func TestUnitSGSnapshotPolicyUpdate_UnknownFieldNotSent(t *testing.T) {
	ctx := context.Background()
	var modifyKeys []string
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ModifySecurityGroupSnapshotPolicy":
			for k := range r.Form {
				if k == "SecurityGroupSnapshotPolicyName" || k == "Status" {
					modifyKeys = append(modifyKeys, k)
				}
			}
			popJSON(t, w, map[string]interface{}{})
		case "DescribeSecurityGroupSnapshotPolicies":
			popJSON(t, w, map[string]interface{}{
				"SecurityGroupSnapshotPolicySets": map[string]interface{}{
					"SecurityGroupSnapshotPolicy": []interface{}{
						map[string]interface{}{
							"SecurityGroupSnapshotPolicyId":   "sp-1",
							"SecurityGroupSnapshotPolicyName": "tf-policy",
							"Status":                          "disable",
							"RegionId":                        "cn-hangzhou",
						},
					},
				},
			})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	// status is explicitly set to "disable"; name is left unknown (unset).
	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(ctx, &securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringUnknown(),
		Status:                          types.StringValue("disable"),
	}); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	state := stateFromModel(t, ctx, sch, securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
		RegionId:                        types.StringValue("cn-hangzhou"),
	})
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch}}
	r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Update must not error: %v", resp.Diagnostics)
	}
	hasStatus, hasName := false, false
	for _, k := range modifyKeys {
		if k == "Status" {
			hasStatus = true
		}
		if k == "SecurityGroupSnapshotPolicyName" {
			hasName = true
		}
	}
	if !hasStatus {
		t.Fatalf("Status (explicitly changed) must be sent to Modify; got %v", modifyKeys)
	}
	if hasName {
		t.Fatalf("an unknown (unset) name must not be sent to Modify; it would clear the field; got %v", modifyKeys)
	}
}

// F3: both Optional+Computed fields explicitly changed are both sent.
func TestUnitSGSnapshotPolicyUpdate_BothChangedSent(t *testing.T) {
	ctx := context.Background()
	var modifyKeys []string
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ModifySecurityGroupSnapshotPolicy":
			for k := range r.Form {
				if k == "SecurityGroupSnapshotPolicyName" || k == "Status" {
					modifyKeys = append(modifyKeys, k)
				}
			}
			popJSON(t, w, map[string]interface{}{})
		case "DescribeSecurityGroupSnapshotPolicies":
			popJSON(t, w, map[string]interface{}{
				"SecurityGroupSnapshotPolicySets": map[string]interface{}{
					"SecurityGroupSnapshotPolicy": []interface{}{
						map[string]interface{}{
							"SecurityGroupSnapshotPolicyId":   "sp-1",
							"SecurityGroupSnapshotPolicyName": "tf-policy-2",
							"Status":                          "disable",
							"RegionId":                        "cn-hangzhou",
						},
					},
				},
			})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(ctx, &securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy-2"),
		Status:                          types.StringValue("disable"),
	}); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	state := stateFromModel(t, ctx, sch, securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
		RegionId:                        types.StringValue("cn-hangzhou"),
	})
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch}}
	r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Update must not error: %v", resp.Diagnostics)
	}
	hasStatus, hasName := false, false
	for _, k := range modifyKeys {
		if k == "Status" {
			hasStatus = true
		}
		if k == "SecurityGroupSnapshotPolicyName" {
			hasName = true
		}
	}
	if !hasStatus || !hasName {
		t.Fatalf("both changed fields must be sent to Modify; got status=%v name=%v (%v)", hasStatus, hasName, modifyKeys)
	}
}

// F2 (#8/#9): a Describe returning HTTP 200 with a well-formed list root but a
// malformed element — a non-map item, a map missing its id, or a map whose id
// is the wrong type — must not be treated as a confirmed NotFound. The describe
// helper reports it as an error so Read keeps the prior state and surfaces a
// diagnostic, instead of returning (nil, nil) and dropping the resource. Each
// case drives the real Resource.Read through the configured client and an
// injected prior state, not a copied helper.
func TestUnitSGSnapshotPolicyRead_MalformedElementsKeepState(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name    string
		element interface{}
	}{
		{"non-map item", "not-an-object"},
		{"map missing id", map[string]interface{}{"Status": "enable"}},
		{"map wrong-type id", map[string]interface{}{"SecurityGroupSnapshotPolicyId": map[string]interface{}{"x": "sp-1"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DescribeSecurityGroupSnapshotPolicies":
					popJSON(t, w, map[string]interface{}{
						"SecurityGroupSnapshotPolicySets": []interface{}{c.element},
					})
				default:
					t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				}
			})
			r, sch := newSGSnapshotPolicyResource(t, client)

			prior := stateFromModel(t, ctx, sch, securityGroupSnapshotPolicyResourceModel{
				Id:                              types.StringValue("sp-1"),
				SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
				SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
				Status:                          types.StringValue("enable"),
				RegionId:                        types.StringValue("cn-hangzhou"),
			})
			resp := resource.ReadResponse{State: prior}
			r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

			if !resp.Diagnostics.HasError() {
				t.Fatalf("a malformed element must surface a diagnostic, not be swallowed as NotFound")
			}
			var out securityGroupSnapshotPolicyResourceModel
			resp.State.Get(ctx, &out)
			if out.Id.ValueString() != "sp-1" || out.SecurityGroupSnapshotPolicyId.ValueString() != "sp-1" {
				t.Fatalf("a malformed element must not drop the resource; id/policy_id still %q expected, got id=%q policy_id=%q",
					"sp-1", out.Id.ValueString(), out.SecurityGroupSnapshotPolicyId.ValueString())
			}
		})
	}
}

// F2 (#8/#9): a non-empty well-formed list of parseable policies, none of which
// matches the requested id, IS a confirmed absence (the API returned real
// policies, none is the requested one). Read removes the resource with no error,
// distinguishing a genuine NotFound from a malformed response.
func TestUnitSGSnapshotPolicyRead_OtherValidIdRemovesResource(t *testing.T) {
	ctx := context.Background()
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeSecurityGroupSnapshotPolicies":
			popJSON(t, w, map[string]interface{}{
				"SecurityGroupSnapshotPolicySets": []interface{}{
					map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-other", "Status": "enable"},
				},
			})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	prior := stateFromModel(t, ctx, sch, securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
		RegionId:                        types.StringValue("cn-hangzhou"),
	})
	resp := resource.ReadResponse{State: prior}
	r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("a well-formed list of other policies must be a confirmed NotFound, not an error: %v", resp.Diagnostics)
	}
	var out securityGroupSnapshotPolicyResourceModel
	resp.State.Get(ctx, &out)
	if !out.Id.IsNull() {
		t.Fatalf("a confirmed NotFound must remove the resource; id still %q", out.Id.ValueString())
	}
}

// F3/#5: characterize the real Framework plan/state→HTTP behavior of an
// explicitly-null Optional+Computed field, distinct from an unknown (unset) one.
// An unknown field is not sent (TestUnitSGSnapshotPolicyUpdate_UnknownFieldNotSent);
// an explicitly-null field is known (not unknown) and differs from the prior
// state, so the current guard sends it as an empty string. This verifies the
// actual semantics instead of assuming them from ValueString. It is a
// characterization, not a contract assertion: whether the API should clear the
// field on null (vs ignore it) is cloud-contract-unverified and left for QA/cloud.
func TestUnitSGSnapshotPolicyUpdate_NullFieldCharacterization(t *testing.T) {
	ctx := context.Background()
	var namePresent, statusPresent bool
	var nameVal, statusVal string
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "ModifySecurityGroupSnapshotPolicy":
			if _, ok := r.Form["SecurityGroupSnapshotPolicyName"]; ok {
				namePresent = true
				nameVal = r.Form.Get("SecurityGroupSnapshotPolicyName")
			}
			if _, ok := r.Form["Status"]; ok {
				statusPresent = true
				statusVal = r.Form.Get("Status")
			}
			popJSON(t, w, map[string]interface{}{})
		case "DescribeSecurityGroupSnapshotPolicies":
			popJSON(t, w, map[string]interface{}{
				"SecurityGroupSnapshotPolicySets": map[string]interface{}{
					"SecurityGroupSnapshotPolicy": []interface{}{
						map[string]interface{}{
							"SecurityGroupSnapshotPolicyId":   "sp-1",
							"SecurityGroupSnapshotPolicyName": "tf-policy",
							"Status":                          "disable",
							"RegionId":                        "cn-hangzhou",
						},
					},
				},
			})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	// name is explicitly null in plan; status is explicitly changed.
	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(ctx, &securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringNull(),
		Status:                          types.StringValue("disable"),
	}); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	state := stateFromModel(t, ctx, sch, securityGroupSnapshotPolicyResourceModel{
		Id:                              types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyId:   types.StringValue("sp-1"),
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
		RegionId:                        types.StringValue("cn-hangzhou"),
	})
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch}}
	r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update must not error: %v", resp.Diagnostics)
	}
	// Characterization: null is known and differs from the prior state, so the
	// current guard sends SecurityGroupSnapshotPolicyName as an empty string —
	// distinct from an unknown (unset) field, which is not sent at all.
	if !namePresent {
		t.Fatalf("null name characterization: expected SecurityGroupSnapshotPolicyName to be sent (null is known, not unknown)")
	}
	if nameVal != "" {
		t.Fatalf("null name characterization: expected empty string value, got %q", nameVal)
	}
	if !statusPresent || statusVal != "disable" {
		t.Fatalf("explicitly-changed status must be sent as disable, got present=%v val=%q", statusPresent, statusVal)
	}
}

// D1/T1: a Create response whose SecurityGroupSnapshotPolicyId is the wrong
// type (map/slice/bool) must be diagnosed and rejected — not coerced into a
// fmt.Sprintf pseudo-id such as "map[x:sp-1]" / "[sp-1]" / "true". Create fails
// immediately at extract and never issues the follow-up Describe, so no
// pseudo-id is persisted to state. This drives the real Resource.Create through
// the configured client and the live HTTP chain (Framework request/response +
// RpcPost + httptest server), not a copied helper; the previous default branch
// would have turned a map response into "map[x:sp-1]" and persisted it.
func TestUnitSGSnapshotPolicyCreate_MalformedIdDiagnoses(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name string
		id   interface{}
	}{
		{"map id", map[string]interface{}{"x": "sp-1"}},
		{"slice id", []interface{}{"sp-1"}},
		{"bool id", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var describeCalled bool
			client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "CreateSecurityGroupSnapshotPolicy":
					popJSON(t, w, map[string]interface{}{"SecurityGroupSnapshotPolicyId": c.id})
				case "DescribeSecurityGroupSnapshotPolicies":
					// Create must not reach Describe when the id is malformed; if it
					// does, a pseudo-id is being chased.
					describeCalled = true
					popJSON(t, w, map[string]interface{}{
						"SecurityGroupSnapshotPolicySets": []interface{}{},
					})
				default:
					t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				}
			})
			r, sch := newSGSnapshotPolicyResource(t, client)

			plan := tfsdk.Plan{Schema: sch}
			if diags := plan.Set(ctx, &securityGroupSnapshotPolicyResourceModel{
				SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
				Status:                          types.StringValue("enable"),
			}); diags.HasError() {
				t.Fatalf("plan.Set: %v", diags)
			}
			resp := resource.CreateResponse{State: tfsdk.State{Schema: sch}}
			r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

			if !resp.Diagnostics.HasError() {
				t.Fatalf("a wrong-type id must be diagnosed, not coerced into a pseudo-id")
			}
			if describeCalled {
				t.Fatalf("Create must not issue a follow-up Describe when the id is malformed; a pseudo-id must not be chased")
			}
			var out securityGroupSnapshotPolicyResourceModel
			resp.State.Get(ctx, &out)
			if out.Id.ValueString() != "" || out.SecurityGroupSnapshotPolicyId.ValueString() != "" {
				t.Fatalf("a malformed id must not be persisted to state; got id=%q policy_id=%q",
					out.Id.ValueString(), out.SecurityGroupSnapshotPolicyId.ValueString())
			}
		})
	}
}

// D2/T2: a numeric (float64) id returned by both Create and Describe must
// persist to state as the same normalized string ("42") across the whole
// Create->Describe->state chain, not collapse to null and overwrite the id
// persisted by Create. This drives the real Resource.Create through the
// configured client and the live HTTP chain (Framework request/response +
// RpcPost + httptest server), exercising the shared securityGroupSnapshotPolicyIdOf
// contract (strconv.FormatFloat) in extract, policyIdOfOk and the state writer.
// Before D2, the state writer only accepted a string id, so a numeric id from
// Describe became StringNull() and overwrote the "42" persisted by Create.
func TestUnitSGSnapshotPolicyCreate_NumericIdPersistsToState(t *testing.T) {
	ctx := context.Background()
	client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "CreateSecurityGroupSnapshotPolicy":
			popJSON(t, w, map[string]interface{}{"SecurityGroupSnapshotPolicyId": float64(42)})
		case "DescribeSecurityGroupSnapshotPolicies":
			// Describe echoes the same id back as a JSON number; the state writer
			// must normalize it the same way extract did, not drop it to null.
			popJSON(t, w, map[string]interface{}{
				"SecurityGroupSnapshotPolicySets": []interface{}{
					map[string]interface{}{
						"SecurityGroupSnapshotPolicyId":   float64(42),
						"SecurityGroupSnapshotPolicyName": "tf-policy",
						"Status":                          "enable",
						"RegionId":                        "cn-hangzhou",
					},
				},
			})
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
		}
	})
	r, sch := newSGSnapshotPolicyResource(t, client)

	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(ctx, &securityGroupSnapshotPolicyResourceModel{
		SecurityGroupSnapshotPolicyName: types.StringValue("tf-policy"),
		Status:                          types.StringValue("enable"),
	}); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("a numeric id must round-trip through Create/Describe without error: %v", resp.Diagnostics)
	}
	var out securityGroupSnapshotPolicyResourceModel
	if diags := resp.State.Get(ctx, &out); diags.HasError() {
		t.Fatalf("resp.State.Get: %v", diags)
	}
	if out.Id.ValueString() != "42" {
		t.Fatalf("numeric id must persist as the normalized string %q, got id=%q", "42", out.Id.ValueString())
	}
	if out.SecurityGroupSnapshotPolicyId.ValueString() != "42" {
		t.Fatalf("SecurityGroupSnapshotPolicyId must equal the normalized id %q, got %q", "42", out.SecurityGroupSnapshotPolicyId.ValueString())
	}
	if out.Id.ValueString() != out.SecurityGroupSnapshotPolicyId.ValueString() {
		t.Fatalf("id (%q) must equal security_group_snapshot_policy_id (%q)", out.Id.ValueString(), out.SecurityGroupSnapshotPolicyId.ValueString())
	}
}
