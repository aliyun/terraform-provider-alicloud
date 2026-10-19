package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// savingPlanTestState builds a valid Framework state by populating it from a
// model with id and saving_plan_id set and all other attributes null. The
// Schema field on tfsdk.State is of an internal Framework interface type;
// schema.Schema implements it and the field is exported, so direct
// assignment works without importing the internal package.
func savingPlanTestState(t *testing.T, r *savingPlanResource) tfsdk.State {
	t.Helper()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	state := tfsdk.State{Schema: schemaResp.Schema}
	model := savingPlanResourceModel{
		Id:           types.StringValue("spn-test-123"),
		SavingPlanId: types.StringValue("spn-test-123"),
	}
	if diags := state.Set(context.Background(), &model); diags.HasError() {
		t.Fatalf("state.Set: %v", diags)
	}
	return state
}

// savingPlanConfirmedVisibleState builds a state whose Computed-only attributes
// (create_time, region_id, payment_type) are populated, mirroring the state a
// successful Read leaves. savingPlanWasConfirmedVisible treats this as evidence
// the plan was previously confirmed to exist — used to assert that a not-found
// signal STILL retains the paid ID (the BssOpenApi/2017-12-14 contract does not
// pin a 404 / NotFound to this instance, so prior visibility is not a deletion
// authorisation). A freshly purchased-but-not-yet-visible ID
// (savingPlanTestState, all Computed null) is the contrast case.
func savingPlanConfirmedVisibleState(t *testing.T, r *savingPlanResource) tfsdk.State {
	t.Helper()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	state := tfsdk.State{Schema: schemaResp.Schema}
	model := savingPlanResourceModel{
		Id:           types.StringValue("spn-test-123"),
		SavingPlanId: types.StringValue("spn-test-123"),
		CreateTime:   types.StringValue("2026-01-01 00:00:00"),
		RegionId:     types.StringValue("cn-hangzhou"),
		PaymentType:  types.StringValue("AllUpfront"),
	}
	if diags := state.Set(context.Background(), &model); diags.HasError() {
		t.Fatalf("state.Set: %v", diags)
	}
	return state
}

// savingPlanWasConfirmedVisible reports whether prior state carries evidence
// that a prior Read refreshed real API data: a Computed-only field (create_time,
// region_id, payment_type) is non-null. Production Read no longer uses this to
// gate removal — a not-found signal never drops a paid ID, because the
// BssOpenApi/2017-12-14 contract does not pin a 404 / NotFound to this instance —
// but the marker remains a useful test assertion that a successful refresh
// populated the Computed attributes (and that a fresh-purchase persisted state
// did not).
func savingPlanWasConfirmedVisible(state savingPlanResourceModel) bool {
	return !state.CreateTime.IsNull() || !state.RegionId.IsNull() || !state.PaymentType.IsNull()
}

// savingPlanTestPlan builds a valid Framework plan with the required
// committed_amount attribute populated and a known instance_family (required
// when purchase_method is ByInstanceFamily), so Create passes its
// conditional-required validation. Every attribute is a concrete known value:
// this is the known-plan path. The real unknown-Computed plan path that
// exercises the NewState and request-construction handling is covered by
// savingPlanUnknownComputedTestPlan; a fully-known plan must not be the only
// shape the tests use or it masks unknown NewState / silent-coercion defects.
func savingPlanTestPlan(t *testing.T, r *savingPlanResource) tfsdk.Plan {
	t.Helper()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	model := savingPlanResourceModel{
		CommittedAmount: types.StringValue("0.1"),
		PlanType:        types.StringValue("EcsCompute"),
		OfferingType:    types.StringValue("AllUpfront"),
		PurchaseMethod:  types.StringValue("ByInstanceFamily"),
		ChargeType:      types.StringValue("PrePaid"),
		InstanceFamily:  types.StringValue("ecs.g6"),
	}
	if diags := plan.Set(context.Background(), &model); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	return plan
}

// savingPlanUnknownComputedTestPlan builds a real Framework Create plan whose
// Computed attributes are unknown, mirroring what the runtime produces when a
// user omits an Optional+Computed field. Default-bearing attributes
// (plan_type, offering_type, purchase_method, charge_type) are known because
// their StaticString plan modifier applies during PlanResourceChanges; the
// purely Computed attributes (region_id, saving_plan_id, create_time,
// payment_type) and the omitted Optional+Computed attributes (period,
// start_time) are unknown at Create with no prior state. instance_family is
// passed in so the caller can drive both the valid known case and the
// conditionally-required-but-unknown misconfiguration case.
func savingPlanUnknownComputedTestPlan(t *testing.T, r *savingPlanResource, instanceFamily types.String) tfsdk.Plan {
	t.Helper()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	model := savingPlanResourceModel{
		CommittedAmount: types.StringValue("0.1"),
		PlanType:        types.StringValue("EcsCompute"),
		OfferingType:    types.StringValue("AllUpfront"),
		PurchaseMethod:  types.StringValue("ByInstanceFamily"),
		ChargeType:      types.StringValue("PrePaid"),
		InstanceFamily:  instanceFamily,
		// Omitted Optional+Computed attributes (no StaticString default) arrive
		// as unknown at Create plan time, as do the purely Computed attributes
		// with no prior state to reuse. This mirrors a real Framework Create
		// plan whose NewState must be fully known after apply.
		Period:       types.Int64Unknown(),
		StartTime:    types.StringUnknown(),
		RegionId:     types.StringUnknown(),
		SavingPlanId: types.StringUnknown(),
		CreateTime:   types.StringUnknown(),
		PaymentType:  types.StringUnknown(),
	}
	if diags := plan.Set(context.Background(), &model); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	return plan
}

// assertSavingPlanModelNoUnknowns fails the test if any attribute of the model
// is unknown — the Framework ApplyResourceChange NewState for a newly-created
// resource must be fully known, so a Create that persists the purchased ID on a
// refresh miss must not leak any unknown Computed value into state.
func assertSavingPlanModelNoUnknowns(t *testing.T, m savingPlanResourceModel) {
	t.Helper()
	strs := map[string]types.String{
		"id": m.Id, "committed_amount": m.CommittedAmount, "plan_type": m.PlanType,
		"offering_type": m.OfferingType, "purchase_method": m.PurchaseMethod,
		"charge_type": m.ChargeType, "period_unit": m.PeriodUnit,
		"instance_family": m.InstanceFamily, "instance_family_set": m.InstanceFamilySet,
		"saving_plan_name": m.SavingPlanName, "description": m.Description,
		"start_time": m.StartTime, "region_id": m.RegionId,
		"saving_plan_id": m.SavingPlanId, "create_time": m.CreateTime,
		"payment_type": m.PaymentType,
	}
	for name, v := range strs {
		if v.IsUnknown() {
			t.Fatalf("state attribute %s must not be unknown after Create; NewState must be fully known", name)
		}
	}
	if m.Period.IsUnknown() {
		t.Fatal("state attribute period must not be unknown after Create; NewState must be fully known")
	}
}

// savingPlanTestCreateResponse returns a CreateResponse whose State carries
// the resource Schema, mirroring what the Framework runtime does for real
// requests. Create calls resp.State.Set, which dereferences the Schema — a
// zero-value tfsdk.State panics on Set, so tests must pre-populate it exactly
// as the Framework would when dispatching a real Create.
func savingPlanTestCreateResponse(t *testing.T, r *savingPlanResource) resource.CreateResponse {
	t.Helper()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	return resource.CreateResponse{
		State: tfsdk.State{Schema: schemaResp.Schema},
	}
}

// savingPlanTestClient returns an AliyunClient whose RpcPost calls are served
// by an in-process httptest.Server. The handler dispatches by the Action form
// field. No real cloud endpoint is contacted — the tea-rpc client is pointed
// at the local server via config.Endpoints for both "ecs" (Purchase) and
// "bssopenapi" (Query). This follows the established ecsSnapshotTestClient
// pattern. It never opens ALICLOUD_SAVING_PLAN_TEST and performs no purchase:
// the server returns synthetic responses only.
func savingPlanTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm: %v", err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("credentials.NewCredential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
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
	endpoints.Store("ecs", endpoint)
	endpoints.Store("bssopenapi", endpoint)
	return client
}

// newSavingPlanResourceWithClient injects a mock client into a Framework
// resource via Configure, mirroring how the provider wires the real client.
func newSavingPlanResourceWithClient(t *testing.T, client *connectivity.AliyunClient) *savingPlanResource {
	t.Helper()
	r := &savingPlanResource{}
	var configureResp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("Configure: %v", configureResp.Diagnostics)
	}
	return r
}

// TestUnitSavingPlanDeleteIsStateOnly verifies that Delete does not call any API.
// Saving Plans are financial commitments: the PurchaseSavingPlanOffering API
// cannot be reversed and no cancellation API exists. Delete therefore removes
// only the Terraform state entry. Using a mock server lets us assert zero API
// calls were made, proving the state-only contract without incurring cost.
func TestUnitSavingPlanDeleteIsStateOnly(t *testing.T) {
	var requestCount int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
	})
	r := newSavingPlanResourceWithClient(t, client)
	var resp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{
		State: savingPlanTestState(t, r),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete with valid state should not produce errors: %v", resp.Diagnostics)
	}
	if requestCount != 0 {
		t.Fatalf("Delete must not call any API; got %d requests", requestCount)
	}
}

// TestUnitSavingPlanReadPreservesStateWhenNotVisible verifies the P1-1 fix:
// when Query returns a legitimate success whose Items list does not contain the
// plan (the eventual-consistency "not yet visible" window, or an absence the
// contract cannot authoritatively confirm), Read must NOT RemoveResource the
// paid ID — it preserves state and warns so the next apply refreshes without
// re-purchasing. The previous implementation returned (nil, nil) for this shape
// and RemoveResource'd the ID on the same empty-Items response Create treated as
// "keep the ID", losing a financial commitment and forcing a re-purchase. No
// query signal — not even a precise NotFound — authorises removal; the resource
// is removed only by deleting it from configuration.
func TestUnitSavingPlanReadPreservesStateWhenNotVisible(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "QuerySavingsPlansInstance":
			// Legitimate BSS success with empty Items: the plan is not yet
			// visible (or its absence cannot be authoritatively confirmed).
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true,
				"Code":    "Success",
				"Message": "",
				"Data":    map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanTestState(t, r)
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read should not error when the plan is not yet visible (ID preserved): %v", resp.Diagnostics)
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("Read must preserve the paid ID when the plan is not yet visible; state was removed")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("preserved state should keep the saving plan id; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanReadRetainsUnconfirmedIDOnNotFound verifies that a precise
// NotFound signal (Code=="NotFound") NEVER removes a paid ID — here one that has
// never been confirmed visible since purchase/import. The state has only the
// id/saving_plan_id (no create_time/region_id/payment_type), mirroring a freshly
// purchased plan whose post-purchase refresh missed — the same NotFound
// response Create keeps the ID for (a propagation window). The
// BssOpenApi/2017-12-14 contract does not pin a 404 / NotFound to this instance,
// so the signal is not an authorisation to drop the financial commitment; Read
// retains the ID with a warning (not an error) so the next apply refreshes
// without re-purchasing. (The confirmed-visible counterpart — where the ID is
// STILL retained — is TestUnitSavingPlanReadRetainsIDOnNotFoundEvenWhenPreviouslyVisible.)
func TestUnitSavingPlanReadRetainsUnconfirmedIDOnNotFound(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "QuerySavingsPlansInstance":
			// 200 with a precise NotFound code envelope: formatError synthesises a
			// *tea.SDKError whose Code isSavingPlanNotFound recognises exactly
			// (no err.Error() substring match), so querySavingPlan returns
			// errSavingPlanQueryNotFound. The Message must be non-empty so
			// formatError surfaces the SDKError rather than early-returning nil.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": false,
				"Code":    "NotFound",
				"Message": "The specified instance is not found.",
				"Data":    map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanTestState(t, r) // only id — not confirmed visible
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read should warn (not error) when an unconfirmed ID is reported NotFound: %v", resp.Diagnostics)
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("Read must retain an unconfirmed paid ID on NotFound (propagation window); state was removed")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("retained state should keep the saving plan id; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanReadRetainsIDOnNotFoundEvenWhenPreviouslyVisible verifies
// that a precise not-found signal (Code=="NotFound") NEVER drops a paid Saving
// Plan from state, even when the plan was previously confirmed visible. The
// BssOpenApi/2017-12-14 QuerySavingsPlansInstance contract does not guarantee
// that a 404 / NotFound means THIS instance is gone (it may be routing/endpoint
// or scoped to a different object; the official error table documents no
// instance-NotFound code), so prior visibility only proves past existence, not
// that this signal is the instance being deleted. Read therefore retains the ID
// with a warning in every case, and the resource is removed only by deleting it
// from configuration. This is the HTTP closed loop:
// previously-visible → same generic NotFound → ID retained.
func TestUnitSavingPlanReadRetainsIDOnNotFoundEvenWhenPreviouslyVisible(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "QuerySavingsPlansInstance":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": false,
				"Code":    "NotFound",
				"Message": "The specified instance is not found.",
				"Data":    map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanConfirmedVisibleState(t, r) // create_time/region_id populated
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read should warn (not error) on a not-found signal; got error: %v", resp.Diagnostics)
	}
	if resp.State.Raw.IsNull() {
		t.Fatalf("Read must retain the paid ID on a not-found signal even when previously visible; state was removed")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("the paid ID must be retained on a not-found signal; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanReadPreservesStateOnError verifies that when the Query API
// fails with a non-not-found error (e.g. server error), Read reports the error
// and preserves the existing state — it must not silently remove the resource
// the way it does for not-found. This is the "malformed response preserves
// state" path: only a proven-absent plan drops state; a transport or server
// error keeps it for the next refresh.
func TestUnitSavingPlanReadPreservesStateOnError(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "QuerySavingsPlansInstance":
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanTestState(t, r)
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{
		State: initialState,
	}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should report an error when the Query API fails")
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("Read should preserve state on non-not-found error; state was removed")
	}
}

// TestUnitSavingPlanCreatePersistsIDWhenRefreshMisses verifies the core
// financial-commitment safeguard: when Purchase succeeds but the plan is not
// yet visible via Query (eventual-consistency delay, empty items), Create
// persists the purchased ID to state BEFORE refreshing so the next apply
// refreshes instead of re-purchasing — never lose a paid Saving Plan ID. A
// warning (not an error) signals the refresh miss.
func TestUnitSavingPlanCreatePersistsIDWhenRefreshMisses(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			// Not yet visible: a legitimate BSS success with empty Items, which
			// querySavingPlan treats as not-found so Create persists the purchased
			// ID with a warning instead of re-purchasing.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true,
				"Code":    "Success",
				"Message": "",
				"Data":    map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanTestPlan(t, r),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create should not error when refresh misses (ID persisted): %v", resp.Diagnostics)
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("purchased ID should be persisted to state; got %q", model.Id.ValueString())
	}
	if model.SavingPlanId.ValueString() != "spn-test-123" {
		t.Fatalf("saving_plan_id should be persisted to state; got %q", model.SavingPlanId.ValueString())
	}
}

// TestUnitSavingPlanCreatePersistsIDWhenRefreshErrors verifies that when Purchase
// succeeds but the subsequent Query fails with a transport/server error (not
// not-found), Create still persists the purchased ID to state and reports a
// warning rather than an error. The ID must survive a failed refresh so the
// next apply retries the refresh without re-purchasing.
func TestUnitSavingPlanCreatePersistsIDWhenRefreshErrors(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanTestPlan(t, r),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create should not error when refresh fails (ID persisted): %v", resp.Diagnostics)
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("purchased ID should be persisted even when refresh errors; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanReadPreservesStateOnMalformedMissingData verifies that a 200
// response missing the Data object is an error, not a not-found, so Read
// preserves state rather than RemoveResource'ing a paid plan ID. A
// transiently malformed response must never drop a financial commitment.
func TestUnitSavingPlanReadPreservesStateOnMalformedMissingData(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			// Legitimate success envelope but the body lacks the Data object.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "", "RequestId": "req-1",
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanTestState(t, r)
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should error when Data is missing from the response")
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("Read should preserve state on a malformed (missing Data) response; state was removed")
	}
}

// TestUnitSavingPlanReadPreservesStateOnMalformedItemsType verifies that a 200
// response whose Data.Items is not a list is an error, preserving state.
func TestUnitSavingPlanReadPreservesStateOnMalformedItemsType(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": "not-a-list"},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanTestState(t, r)
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should error when Data.Items is not a list")
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("Read should preserve state on a malformed (Items not a list) response; state was removed")
	}
}

// TestUnitSavingPlanReadPreservesStateOnMalformedItemType verifies that a 200
// response whose first item is not an object is an error, preserving state.
func TestUnitSavingPlanReadPreservesStateOnMalformedItemType(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{"not-an-object"}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanTestState(t, r)
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should error when the first item is not an object")
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("Read should preserve state on a malformed (first item not an object) response; state was removed")
	}
}

// TestUnitSavingPlanCreateRejectsMissingIDSets verifies that when the Purchase
// response carries no SavingPlanIdSets at all, Create reports an error and
// persists no ID to state — never synthesising a placeholder. Only a real,
// validated ID enters state.
func TestUnitSavingPlanCreateRejectsMissingIDSets(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "PurchaseSavingPlanOffering" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{Plan: savingPlanTestPlan(t, r)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Create must error when SavingPlanIdSets is missing")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "" {
		t.Fatalf("no synthetic ID should be persisted when SavingPlanIdSets is missing; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanCreateRejectsNonSliceIDSets verifies that a SavingPlanIdSets
// value of the wrong type (e.g. a string) is rejected with an error and no ID
// is persisted, rather than coercing or guessing.
func TestUnitSavingPlanCreateRejectsNonSliceIDSets(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "PurchaseSavingPlanOffering" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": "not-a-list",
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{Plan: savingPlanTestPlan(t, r)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Create must error when SavingPlanIdSets is not a list")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "" {
		t.Fatalf("no synthetic ID should be persisted when SavingPlanIdSets is malformed; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanCreateRejectsEmptyIDSets verifies that an empty
// SavingPlanIdSets list is rejected with an error and no ID persisted — a
// successful purchase must return the IDs purchased.
func TestUnitSavingPlanCreateRejectsEmptyIDSets(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "PurchaseSavingPlanOffering" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{Plan: savingPlanTestPlan(t, r)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Create must error when SavingPlanIdSets is empty")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "" {
		t.Fatalf("no synthetic ID should be persisted when SavingPlanIdSets is empty; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanCreateRejectsNonStringID verifies that a SavingPlanIdSets
// entry that is not a string (e.g. a number) is rejected with an error and no
// ID persisted. fmt.Sprint would stringify it and mask the malformation; the
// strict assertion rejects it instead.
func TestUnitSavingPlanCreateRejectsNonStringID(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "PurchaseSavingPlanOffering" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{12345},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{Plan: savingPlanTestPlan(t, r)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Create must error when SavingPlanIdSets[0] is not a string")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "" {
		t.Fatalf("no synthetic ID should be persisted when SavingPlanIdSets[0] is non-string; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanReadPreservesStateOnBusinessFailure verifies that a 200
// response whose BSS envelope indicates a business failure — not a legitimate
// success — is treated as an error and preserves the paid ID in state, never
// as a not-found that would RemoveResource it. RpcPost→formatError returns nil
// for several non-success shapes (a Code of success/ok/200/0, a Success of
// true, or an empty/absent Message), so a normal failure with a non-empty
// Message is caught as an SDKError while a failure with an empty Message, a
// missing Message, or a contradictory Success/Code pair slips through as
// err==nil with empty Items. querySavingPlan validates the envelope locally so
// every business-failure shape surfaces as an error and keeps the paid ID.
func TestUnitSavingPlanReadPreservesStateOnBusinessFailure(t *testing.T) {
	cases := []struct {
		name string
		body map[string]interface{}
	}{
		{
			// formatError catches this: Code not success, Success false, Message
			// non-empty -> SDKError (err != nil branch).
			"success_false_code_error_message_nonempty",
			map[string]interface{}{"Success": false, "Code": "InvalidParam.NotSupport", "Message": "not supported", "Data": map[string]interface{}{"Items": []interface{}{}}},
		},
		{
			// formatError returns nil: Message present but empty -> early nil.
			// Caught by the local envelope check (Success=false).
			"success_false_code_error_message_empty",
			map[string]interface{}{"Success": false, "Code": "InvalidParam.NotSupport", "Message": "", "Data": map[string]interface{}{"Items": []interface{}{}}},
		},
		{
			// formatError catches this: Message absent -> SDKError (err != nil).
			"success_false_code_error_message_missing",
			map[string]interface{}{"Success": false, "Code": "InvalidParam.NotSupport", "Data": map[string]interface{}{"Items": []interface{}{}}},
		},
		{
			// formatError returns nil: Code=="Success" early-returns nil even
			// though Success==false. Caught by the local envelope check.
			"success_false_code_success_contradictory",
			map[string]interface{}{"Success": false, "Code": "Success", "Message": "ok", "Data": map[string]interface{}{"Items": []interface{}{}}},
		},
		{
			// formatError returns nil: Success==true early-returns nil even
			// though Code is a business error. Caught by the local envelope check
			// (legitimate success requires BOTH Success==true AND a success Code).
			"success_true_code_error_contradictory",
			map[string]interface{}{"Success": true, "Code": "InvalidParam.NotSupport", "Message": "", "Data": map[string]interface{}{"Items": []interface{}{}}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
					_ = json.NewEncoder(w).Encode(tc.body)
				}
			})
			r := newSavingPlanResourceWithClient(t, client)
			initialState := savingPlanTestState(t, r)
			resp := resource.ReadResponse{State: initialState}
			r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
			if !resp.Diagnostics.HasError() {
				t.Fatal("Read should error on a business-failure envelope")
			}
			if resp.State.Raw.IsNull() {
				t.Fatal("Read should preserve state on a business-failure envelope; state was removed")
			}
		})
	}
}

// TestUnitSavingPlanReadRefreshsOnLegitimateSuccess verifies that a well-formed BSS
// success envelope (Success=true, Code=Success, legitimate Data.Items)
// refreshes state without error — the strict envelope check does not reject a
// genuine success. The response carries the REAL BSS enum value space
// (SavingsType="ECS", PayMode="total"), not the Purchase-form schema enums:
// savingPlanState must translate them onto plan_type/offering_type so the
// refreshed state stays inside the schema's OneOf (EcsCompute/AllUpfront)
// instead of writing the raw BSS values, which would violate OneOf and, since
// both are RequiresReplace, schedule a replacement on the next plan.
func TestUnitSavingPlanReadRefreshsOnLegitimateSuccess(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true,
				"Code":    "Success",
				"Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{
					map[string]interface{}{
						"InstanceId":     "spn-test-123",
						"PoolValue":      "0.1",
						"SavingsType":    "ECS",
						"PayMode":        "total",
						"Region":         "cn-hangzhou",
						"StartTimestamp": "2026-01-01 00:00:00",
					},
				}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	initialState := savingPlanTestState(t, r)
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on a legitimate success envelope: %v", resp.Diagnostics)
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("Read should refresh state on a legitimate success envelope; state was removed")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("refreshed state should keep the saving plan id; got %q", model.Id.ValueString())
	}
	if model.PlanType.ValueString() != "EcsCompute" {
		t.Fatalf("BSS SavingsType ECS must normalize to plan_type EcsCompute; got %q", model.PlanType.ValueString())
	}
	if model.OfferingType.ValueString() != "AllUpfront" {
		t.Fatalf("BSS PayMode total must normalize to offering_type AllUpfront; got %q", model.OfferingType.ValueString())
	}
}

// TestUnitSavingPlansDataSourceRejectsBusinessFailure verifies the data source
// reuses the same BSS success contract as querySavingPlan: a business-failure
// envelope (which formatError may pass through as err==nil, e.g. an empty
// Message) must not be returned as an empty plan list that fakes success.
func TestUnitSavingPlansDataSourceRejectsBusinessFailure(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": false, "Code": "InvalidParam.NotSupport", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans should error on a business-failure envelope; got plans=%v", plans)
	}
}

// TestUnitSavingPlansDataSourceReturnsEmptyOnLegitimateSuccess verifies that a
// legitimate success with empty Items and a zero TotalCount returns an empty
// list (no plans) without error — the data source reports genuine absence, not
// a faked success. TotalCount=0 is the authoritative signal that pagination is
// complete on the first page.
func TestUnitSavingPlansDataSourceReturnsEmptyOnLegitimateSuccess(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{}, "TotalCount": 0, "PageNum": 1, "PageSize": 100},
			})
		}
	})
	plans, err := querySavingPlans(client, "")
	if err != nil {
		t.Fatalf("querySavingPlans should not error on a legitimate success: %v", err)
	}
	if len(plans) != 0 {
		t.Fatalf("querySavingPlans should return an empty list on legitimate success with no items; got %v", plans)
	}
}

// savingPlanItem builds a single synthetic plan item used by the pagination
// regression tests, keyed by instance ID.
func savingPlanItem(id string) map[string]interface{} {
	return map[string]interface{}{
		"InstanceId": id,
		"PoolValue":  "0.1",
		"PayMode":    "AllUpfront",
		"Region":     "cn-hangzhou",
	}
}

// savingPlanPageEnvelope builds a legitimate BSS success page carrying the
// given items, TotalCount, PageNum and PageSize.
func savingPlanPageEnvelope(items []interface{}, totalCount, pageNum int) map[string]interface{} {
	return map[string]interface{}{
		"Success": true, "Code": "Success", "Message": "",
		"Data": map[string]interface{}{
			"Items":      items,
			"TotalCount": totalCount,
			"PageNum":    pageNum,
			"PageSize":   100,
		},
	}
}

// TestUnitSavingPlansDataSourcePaginatesByTotalCount verifies the core fix for
// defect #3: when TotalCount exceeds a single page, querySavingPlans follows
// PageNum 2 and collects the second page's items instead of silently
// truncating at 100. TotalCount=2 with one item per page requires a second
// request; both items must be present in the result and exactly two requests
// made.
func TestUnitSavingPlansDataSourcePaginatesByTotalCount(t *testing.T) {
	var requests int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") != "QuerySavingsPlansInstance" {
			return
		}
		requests++
		switch r.Form.Get("PageNum") {
		case "1":
			_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(
				[]interface{}{savingPlanItem("spn-1")}, 2, 1))
		default:
			_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(
				[]interface{}{savingPlanItem("spn-2")}, 2, 2))
		}
	})
	plans, err := querySavingPlans(client, "")
	if err != nil {
		t.Fatalf("querySavingPlans should paginate without error: %v", err)
	}
	if requests != 2 {
		t.Fatalf("querySavingPlans should issue 2 paginated requests; got %d", requests)
	}
	if len(plans) != 2 {
		t.Fatalf("querySavingPlans should collect both pages' items; got %d plans", len(plans))
	}
	if got := fmt.Sprint(plans[0]["InstanceId"]); got != "spn-1" {
		t.Fatalf("first plan should be spn-1; got %q", got)
	}
	if got := fmt.Sprint(plans[1]["InstanceId"]); got != "spn-2" {
		t.Fatalf("second plan should be spn-2; got %q", got)
	}
}

// TestUnitSavingPlansDataSourcePaginatesMultiplePages verifies pagination across
// more than two pages: TotalCount=250 with 100 items per page requires three
// requests and must yield all 250 plans.
func TestUnitSavingPlansDataSourcePaginatesMultiplePages(t *testing.T) {
	var requests int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") != "QuerySavingsPlansInstance" {
			return
		}
		requests++
		pageNum, _ := strconv.Atoi(r.Form.Get("PageNum"))
		switch pageNum {
		case 1, 2:
			items := make([]interface{}, 100)
			for i := range items {
				items[i] = savingPlanItem(fmt.Sprintf("spn-%d", (pageNum-1)*100+i+1))
			}
			_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(items, 250, pageNum))
		default:
			items := make([]interface{}, 50)
			for i := range items {
				items[i] = savingPlanItem(fmt.Sprintf("spn-%d", 200+i+1))
			}
			_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(items, 250, pageNum))
		}
	})
	plans, err := querySavingPlans(client, "")
	if err != nil {
		t.Fatalf("querySavingPlans should paginate across pages without error: %v", err)
	}
	if requests != 3 {
		t.Fatalf("querySavingPlans should issue 3 paginated requests; got %d", requests)
	}
	if len(plans) != 250 {
		t.Fatalf("querySavingPlans should collect all 250 plans; got %d", len(plans))
	}
}

// TestUnitSavingPlansDataSourceFailsOnLaterPageError verifies defect #3's
// later-page-failure path: when the first page succeeds but a subsequent page
// returns a transport/server error, querySavingPlans fails closed — it surfaces
// the error instead of returning the partial first-page result as if complete.
func TestUnitSavingPlansDataSourceFailsOnLaterPageError(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Form.Get("Action") != "QuerySavingsPlansInstance" {
			return
		}
		switch r.Form.Get("PageNum") {
		case "1":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(
				[]interface{}{savingPlanItem("spn-1")}, 2, 1))
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans must error when a later page fails; got plans=%v", plans)
	}
}

// TestUnitSavingPlansDataSourceFailsOnMalformedLaterPage verifies defect #3's
// malformed-later-page path: a second page that returns 200 but lacks the Data
// object is a contract violation surfaced as an error, not a silent partial
// result.
func TestUnitSavingPlansDataSourceFailsOnMalformedLaterPage(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") != "QuerySavingsPlansInstance" {
			return
		}
		switch r.Form.Get("PageNum") {
		case "1":
			_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(
				[]interface{}{savingPlanItem("spn-1")}, 2, 1))
		default:
			// Malformed: success envelope but no Data object.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
			})
		}
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans must error on a malformed later page; got plans=%v", plans)
	}
}

// TestUnitSavingPlansDataSourceFailsOnNoProgress verifies defect #3's no-progress
// guard: when TotalCount claims more plans remain but a page returns no items,
// querySavingPlans fails closed instead of looping forever or returning a
// partial list that fakes success.
func TestUnitSavingPlansDataSourceFailsOnNoProgress(t *testing.T) {
	var requests int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") != "QuerySavingsPlansInstance" {
			return
		}
		requests++
		// TotalCount=2 but the first page returns no items — a contract
		// violation the no-progress guard must catch.
		_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(
			[]interface{}{}, 2, 1))
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans must error on no-progress; got plans=%v", plans)
	}
	if requests > 1 {
		t.Fatalf("querySavingPlans should stop after the no-progress page; made %d requests", requests)
	}
}

// TestUnitSavingPlansDataSourceRejectsMissingData verifies defect #1 on the data
// source path: a 200 success envelope missing the Data object is an error, not
// an empty list that fakes success.
func TestUnitSavingPlansDataSourceRejectsMissingData(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
			})
		}
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans must error when Data is missing; got plans=%v", plans)
	}
}

// TestUnitSavingPlansDataSourceRejectsNonListItems verifies defect #1 on the data
// source path: Data.Items that is not a list is an error, not nil,nil.
func TestUnitSavingPlansDataSourceRejectsNonListItems(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": "not-a-list"},
			})
		}
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans must error when Items is not a list; got plans=%v", plans)
	}
}

// TestUnitSavingPlansDataSourceRejectsNonObjectItem verifies defect #1 on the data
// source path: a non-object element in Items is an error rather than silently
// skipped — the previous implementation dropped bad items and reported a
// partial list as complete.
func TestUnitSavingPlansDataSourceRejectsNonObjectItem(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(savingPlanPageEnvelope(
				[]interface{}{"not-an-object"}, 1, 1))
		}
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans must error when an item is not an object; got plans=%v", plans)
	}
}

// TestUnitSavingPlansDataSourceRejectsMissingTotalCount verifies defect #3's
// malformed bound: a success page with valid Items but no TotalCount is an
// error — without the authoritative bound pagination cannot proceed safely, and
// treating a missing TotalCount as zero would fake completion.
func TestUnitSavingPlansDataSourceRejectsMissingTotalCount(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{
					"Items":   []interface{}{savingPlanItem("spn-1")},
					"PageNum": 1, "PageSize": 100,
				},
			})
		}
	})
	plans, err := querySavingPlans(client, "")
	if err == nil {
		t.Fatalf("querySavingPlans must error when TotalCount is missing; got plans=%v", plans)
	}
}

// TestUnitSavingPlansDataSourceAcceptsStringTotalCount verifies the TotalCount
// parser accepts a numeric string form (some POP envelopes stringify counts),
// paginating correctly to collect both items.
func TestUnitSavingPlansDataSourceAcceptsStringTotalCount(t *testing.T) {
	var requests int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") != "QuerySavingsPlansInstance" {
			return
		}
		requests++
		switch r.Form.Get("PageNum") {
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{
					"Items":      []interface{}{savingPlanItem("spn-1")},
					"TotalCount": "2", "PageNum": 1, "PageSize": 100,
				},
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{
					"Items":      []interface{}{savingPlanItem("spn-2")},
					"TotalCount": "2", "PageNum": 2, "PageSize": 100,
				},
			})
		}
	})
	plans, err := querySavingPlans(client, "")
	if err != nil {
		t.Fatalf("querySavingPlans should accept a string TotalCount: %v", err)
	}
	if requests != 2 {
		t.Fatalf("querySavingPlans should issue 2 requests for string TotalCount; got %d", requests)
	}
	if len(plans) != 2 {
		t.Fatalf("querySavingPlans should collect both items; got %d", len(plans))
	}
}

// TestUnitSavingPlanCreateOmitsUnknownComputedFromRequest verifies the P1-2
// request-construction fix. A real Framework Create plan carries unknown
// Optional+Computed attributes (period, start_time — no StaticString default,
// so omitted-by-user arrives as unknown). Framework's ValueString/ValueInt64
// return "" / 0 for an unknown, and IsNull is false for unknown, so the previous
// `!IsNull()` guard would have silently sent Period=0 and StartTime="" in the
// Purchase request. The fix omits any unknown value just like a null value, so
// only known explicit values (committed_amount, instance_family, and the
// default-bearing plan_type/offering_type/purchase_method/charge_type) reach
// the API. This test drives Create with a real unknown plan and captures the
// HTTP request to prove unknown Computed parameters are not silently coerced to
// 0 / "".
func TestUnitSavingPlanCreateOmitsUnknownComputedFromRequest(t *testing.T) {
	var (
		periodSent, instanceFamilySent, startTimeSent bool
		committedAmount, planType, instanceFamily     string
	)
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_, periodSent = r.Form["Period"]
			_, instanceFamilySent = r.Form["InstanceFamily"]
			_, startTimeSent = r.Form["StartTime"]
			committedAmount = r.Form.Get("CommittedAmount")
			planType = r.Form.Get("PlanType")
			instanceFamily = r.Form.Get("InstanceFamily")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringValue("ecs.g6")),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create should not error: %v", resp.Diagnostics)
	}
	// Known explicit / default-bearing values must be sent.
	if committedAmount != "0.1" {
		t.Fatalf("CommittedAmount should be sent as the known value; got %q", committedAmount)
	}
	if planType != "EcsCompute" {
		t.Fatalf("PlanType should be sent as the default known value; got %q", planType)
	}
	if instanceFamily != "ecs.g6" {
		t.Fatalf("InstanceFamily should be sent as the known value; got %q", instanceFamily)
	}
	// Unknown Optional+Computed attributes must be omitted, not coerced to 0 / "".
	if periodSent {
		t.Fatal("Period must be omitted when unknown; the request carried Period (would silently send 0)")
	}
	if startTimeSent {
		t.Fatal("StartTime must be omitted when unknown; the request carried StartTime (would silently send empty)")
	}
	_ = instanceFamilySent // asserted via the value check above
}

// TestUnitSavingPlanCreateDiagnosesUnknownConditionalRequired verifies the P1-2
// conditional-required diagnosis. The schema documents instance_family as
// required when purchase_method is ByInstanceFamily; a user who selects
// ByInstanceFamily but omits instance_family leaves it unknown at Create plan
// time. For a financial commitment that cannot be reversed, the provider must
// diagnose this BEFORE calling Purchase (Framework would otherwise send an
// empty InstanceFamily). Create reports an error, makes no Purchase call, and
// persists no ID.
func TestUnitSavingPlanCreateDiagnosesUnknownConditionalRequired(t *testing.T) {
	var purchaseCalled bool
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			purchaseCalled = true
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	// ByInstanceFamily (the default) with instance_family left unknown — the
	// misconfiguration the provider must diagnose before attempting Purchase.
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringUnknown()),
	}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Create must error when instance_family is unknown and purchase_method is ByInstanceFamily")
	}
	if purchaseCalled {
		t.Fatal("Create must not call Purchase when a conditionally-required attribute is unknown")
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "" {
		t.Fatalf("no ID should be persisted when the conditional-required check fails; got %q", model.Id.ValueString())
	}
}

// TestUnitSavingPlanCreateUnknownPlanErrorReturnsKnownState verifies P1-2 on the
// refresh-error path: when the Create plan carries real unknown Computed
// attributes and the post-purchase refresh fails with a transient error, Create
// persists the purchased ID to a FULLY-KNOWN state (no unknowns) and warns. The
// previous persisted := plan returned the plan's unknown Computed values in
// NewState, violating the ApplyResourceChange protocol; the known-plan helper
// masked this by setting every field to a concrete value.
func TestUnitSavingPlanCreateUnknownPlanErrorReturnsKnownState(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringValue("ecs.g6")),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create should not error when refresh fails (ID persisted): %v", resp.Diagnostics)
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("purchased ID should be persisted to state; got %q", model.Id.ValueString())
	}
	if model.SavingPlanId.ValueString() != "spn-test-123" {
		t.Fatalf("saving_plan_id should be persisted to state; got %q", model.SavingPlanId.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, model)
}

// TestUnitSavingPlanCreateUnknownPlanNotVisibleReturnsKnownState verifies P1-2
// on the empty-Items (not yet visible) path: the persisted NewState is fully
// known and the paid ID is retained — the unknown Computed values the refresh
// had not yet returned are null, not unknown.
func TestUnitSavingPlanCreateUnknownPlanNotVisibleReturnsKnownState(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringValue("ecs.g6")),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create should not error when the plan is not yet visible: %v", resp.Diagnostics)
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("purchased ID should be persisted; got %q", model.Id.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, model)
	if model.RegionId.ValueString() != "" {
		t.Fatalf("RegionId should be null (not unknown, not yet visible) after Create; got %q", model.RegionId.ValueString())
	}
}

// TestUnitSavingPlanCreateReadRetainsIDPurchaseOnce verifies the Create->Read
// purchase-count contract at the Framework protocol entry
// (resource.CreateRequest/ReadRequest): after a successful Purchase whose plan
// is not yet visible, Create persists the paid ID; a subsequent Read on that
// state (still not visible) must keep the ID. Read never calls Purchase, so the
// Purchase count stays 1 across Create+Read — a single financial commitment is
// made while the plan is not yet visible, and the persisted NewState is fully
// known at both stages. This test counts Purchase calls in a Create->Read
// sequence; it does NOT assert a next-plan diff/no-replacement, which is
// covered by the protocol-boundary closed loop in
// TestUnitSavingPlanCreateReadNextPlanNoReplacementEquivalent.
func TestUnitSavingPlanCreateReadRetainsIDPurchaseOnce(t *testing.T) {
	var purchaseCount int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			purchaseCount++
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	createResp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringValue("ecs.g6")),
	}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create should not error (ID persisted): %v", createResp.Diagnostics)
	}
	if purchaseCount != 1 {
		t.Fatalf("Create should call Purchase exactly once; got %d", purchaseCount)
	}
	var createModel savingPlanResourceModel
	createResp.Diagnostics.Append(createResp.State.Get(context.Background(), &createModel)...)
	if createModel.Id.ValueString() != "spn-test-123" {
		t.Fatalf("Create should persist the purchased ID; got %q", createModel.Id.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, createModel)

	// Simulate the next terraform plan/apply refresh: Read the state Create
	// persisted. Read must keep the ID (not RemoveResource) so the plan shows
	// no diff and does not schedule another Purchase.
	readResp := resource.ReadResponse{State: createResp.State}
	r.Read(context.Background(), resource.ReadRequest{State: createResp.State}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error while the plan is not yet visible: %v", readResp.Diagnostics)
	}
	if readResp.State.Raw.IsNull() {
		t.Fatal("Read must keep the paid ID so the next plan does not re-purchase; state was removed")
	}
	if purchaseCount != 1 {
		t.Fatalf("Read must not call Purchase again; got %d total", purchaseCount)
	}
	var model savingPlanResourceModel
	readResp.Diagnostics.Append(readResp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("ID should survive Create->Read; got %q", model.Id.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, model)
}

// TestUnitSavingPlanCreateReadLoopRetainsUnconfirmedIDOnNotFound verifies that a
// not-found signal NEVER drops a paid ID across Create->Read: Create persists the
// paid ID when the initial refresh returns NotFound (the purchase may not have
// propagated yet); a subsequent Read on that SAME NotFound response must NOT
// RemoveResource, because the BssOpenApi/2017-12-14 contract does not pin a 404 /
// NotFound to this instance, so a not-found signal is not a deletion authorisation
// (the plan has also never been confirmed visible since purchase — the persisted
// state's Computed-only attributes are null). Removing the ID on the identical
// response Create kept it for would drop the financial commitment and force a
// re-purchase on the next plan. The paid ID is retained across consecutive
// same-signal Reads; the resource is removed only by deleting it from
// configuration (see TestUnitSavingPlanCreateReadLoopRetainsIDAfterConfirmedVisibleThenNotFound
// for the confirmed-visible counterpart, where the ID is STILL retained).
func TestUnitSavingPlanCreateReadLoopRetainsUnconfirmedIDOnNotFound(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			// Same precise NotFound response for Create's refresh and the later
			// Read — the propagation-window shape Create keeps the ID for.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": false,
				"Code":    "NotFound",
				"Message": "The specified instance is not found.",
				"Data":    map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	createResp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringValue("ecs.g6")),
	}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create should not error (ID persisted): %v", createResp.Diagnostics)
	}
	var createModel savingPlanResourceModel
	createResp.Diagnostics.Append(createResp.State.Get(context.Background(), &createModel)...)
	if createModel.Id.ValueString() != "spn-test-123" {
		t.Fatalf("Create should persist the purchased ID even when the initial refresh reports NotFound; got %q", createModel.Id.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, createModel)
	if savingPlanWasConfirmedVisible(createModel) {
		t.Fatal("a freshly purchased plan whose refresh returned NotFound must not be marked confirmed visible; the convergence gate depends on this")
	}

	// Later refresh on the SAME NotFound response: the plan is still not
	// confirmed visible, so Read must retain the paid ID rather than converge.
	readResp := resource.ReadResponse{State: createResp.State}
	r.Read(context.Background(), resource.ReadRequest{State: createResp.State}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should warn (not error) retaining an unconfirmed ID on NotFound: %v", readResp.Diagnostics)
	}
	if readResp.State.Raw.IsNull() {
		t.Fatal("Read must retain the unconfirmed paid ID on the same NotFound response; state was removed")
	}
	var readModel savingPlanResourceModel
	readResp.Diagnostics.Append(readResp.State.Get(context.Background(), &readModel)...)
	if readModel.Id.ValueString() != "spn-test-123" {
		t.Fatalf("the paid ID must survive Create->Read on the same NotFound; got %q", readModel.Id.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, readModel)
}

// TestUnitSavingPlanCreateReadLoopRetainsIDAfterConfirmedVisibleThenNotFound
// verifies that a not-found signal NEVER drops a paid Saving Plan from state,
// even on the full Create->Read(visible)->Read(NotFound) loop where the plan was
// previously confirmed visible. Create persists the paid ID through the
// not-visible window; a later Read finds the plan and populates the
// Computed-only attributes (confirmed visible) and normalizes the BSS enums onto
// the schema enums; a subsequent precise NotFound (Code=="NotFound") is then
// RETAINED — not a convergence — because the BssOpenApi/2017-12-14 contract
// does not pin a 404 / NotFound to this instance, so prior visibility only
// proves past existence, not that this signal is the instance being deleted. The
// paid ID is retained in every case; the resource is removed only by deleting it
// from configuration.
func TestUnitSavingPlanCreateReadLoopRetainsIDAfterConfirmedVisibleThenNotFound(t *testing.T) {
	var queryCount int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			queryCount++
			switch queryCount {
			case 1:
				// Create's initial refresh: not yet visible (empty Items).
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Success": true, "Code": "Success", "Message": "",
					"Data": map[string]interface{}{"Items": []interface{}{}},
				})
			case 2:
				// Later Read #1: the plan becomes visible — real API data,
				// populating the Computed-only attributes that mark
				// confirmed-visible. Real BSS enum value space (ECS/total).
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Success": true, "Code": "Success", "Message": "",
					"Data": map[string]interface{}{"Items": []interface{}{
						map[string]interface{}{
							"InstanceId":     "spn-test-123",
							"PoolValue":      "0.1",
							"SavingsType":    "ECS",
							"PayMode":        "total",
							"Region":         "cn-hangzhou",
							"StartTimestamp": "2026-01-01 00:00:00",
							"InstanceFamily": "ecs.g6",
							"StartTime":      "2026-01-01 00:00:00",
							"PaymentType":    "Subscription",
						},
					}},
				})
			default:
				// Later Read #2: authoritative NotFound. The plan was confirmed
				// visible, so this is evidence of a real deletion — converge.
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Success": false,
					"Code":    "NotFound",
					"Message": "The specified instance is not found.",
					"Data":    map[string]interface{}{"Items": []interface{}{}},
				})
			}
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	createResp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringValue("ecs.g6")),
	}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create should not error (ID persisted): %v", createResp.Diagnostics)
	}
	var createModel savingPlanResourceModel
	createResp.Diagnostics.Append(createResp.State.Get(context.Background(), &createModel)...)
	if createModel.Id.ValueString() != "spn-test-123" {
		t.Fatalf("Create should persist the purchased ID through the not-visible window; got %q", createModel.Id.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, createModel)
	if savingPlanWasConfirmedVisible(createModel) {
		t.Fatal("immediately after a not-visible purchase the plan must not be confirmed visible")
	}

	// Read #1: the plan becomes visible. State is populated with real API data,
	// marking it confirmed-visible, and the BSS enums normalize to the schema
	// enums (no RequiresReplace-inducing mismatch).
	visibleResp := resource.ReadResponse{State: createResp.State}
	r.Read(context.Background(), resource.ReadRequest{State: createResp.State}, &visibleResp)
	if visibleResp.Diagnostics.HasError() {
		t.Fatalf("Read should refresh without error once visible: %v", visibleResp.Diagnostics)
	}
	if visibleResp.State.Raw.IsNull() {
		t.Fatal("Read should refresh to the full object state; state was removed")
	}
	var visibleModel savingPlanResourceModel
	visibleResp.Diagnostics.Append(visibleResp.State.Get(context.Background(), &visibleModel)...)
	if !savingPlanWasConfirmedVisible(visibleModel) {
		t.Fatal("after a successful Read the plan must be confirmed visible (Computed attributes populated) so a later NotFound can converge")
	}
	if visibleModel.PlanType.ValueString() != "EcsCompute" {
		t.Fatalf("BSS SavingsType ECS must normalize to EcsCompute; got %q", visibleModel.PlanType.ValueString())
	}
	if visibleModel.OfferingType.ValueString() != "AllUpfront" {
		t.Fatalf("BSS PayMode total must normalize to AllUpfront; got %q", visibleModel.OfferingType.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, visibleModel)

	// Read #2: a not-found signal (Code=="NotFound") on a plan that WAS
	// confirmed visible. Because the BssOpenApi/2017-12-14 contract does not pin
	// a 404 / NotFound to this instance, prior visibility only proves past
	// existence — it does not authorise dropping a paid financial commitment.
	// Read must RETAIN the ID with a warning, never converge/remove.
	notFoundResp := resource.ReadResponse{State: visibleResp.State}
	r.Read(context.Background(), resource.ReadRequest{State: visibleResp.State}, &notFoundResp)
	if notFoundResp.Diagnostics.HasError() {
		t.Fatalf("Read should warn (not error) on a not-found signal even after confirmed visible: %v", notFoundResp.Diagnostics)
	}
	if notFoundResp.State.Raw.IsNull() {
		t.Fatal("Read must retain the paid ID on a not-found signal even after the plan was confirmed visible; state was removed")
	}
	var notFoundModel savingPlanResourceModel
	notFoundResp.Diagnostics.Append(notFoundResp.State.Get(context.Background(), &notFoundModel)...)
	if notFoundModel.Id.ValueString() != "spn-test-123" {
		t.Fatalf("the paid ID must survive a not-found signal after confirmed visibility; got %q", notFoundModel.Id.ValueString())
	}
}

// TestUnitSavingPlanCreateReadLoopConvergesWhenLaterVisible verifies the P1-1
// "later visible" path: Create persists the paid ID when the plan is not yet
// visible (Computed values null, fully-known state); once the plan becomes
// visible, Read fills the Computed values from the API without re-purchasing,
// and the resulting state is fully known.
func TestUnitSavingPlanCreateReadLoopConvergesWhenLaterVisible(t *testing.T) {
	var queryCount int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			queryCount++
			if queryCount == 1 {
				// Create's initial refresh: not yet visible.
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Success": true, "Code": "Success", "Message": "",
					"Data": map[string]interface{}{"Items": []interface{}{}},
				})
				return
			}
			// Later Read: the plan is now visible. The response carries the real
			// BSS enum value space (SavingsType ECS / PayMode total), which
			// savingPlanState must normalize to the schema enums.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{
					map[string]interface{}{
						"InstanceId":     "spn-test-123",
						"PoolValue":      "0.1",
						"SavingsType":    "ECS",
						"PayMode":        "total",
						"Region":         "cn-hangzhou",
						"StartTimestamp": "2026-01-01 00:00:00",
						"PaymentType":    "Subscription",
						"InstanceFamily": "ecs.g6",
						"StartTime":      "2026-01-01 00:00:00",
					},
				}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	createResp := savingPlanTestCreateResponse(t, r)
	r.Create(context.Background(), resource.CreateRequest{
		Plan: savingPlanUnknownComputedTestPlan(t, r, types.StringValue("ecs.g6")),
	}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create should not error (ID persisted): %v", createResp.Diagnostics)
	}
	var createModel savingPlanResourceModel
	createResp.Diagnostics.Append(createResp.State.Get(context.Background(), &createModel)...)
	if createModel.Id.ValueString() != "spn-test-123" {
		t.Fatalf("Create should persist the purchased ID when not yet visible; got %q", createModel.Id.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, createModel)
	if createModel.RegionId.ValueString() != "" {
		t.Fatalf("RegionId should be null (not yet visible) after Create; got %q", createModel.RegionId.ValueString())
	}

	// Later refresh: the plan is now visible. Read fills the Computed values.
	readResp := resource.ReadResponse{State: createResp.State}
	r.Read(context.Background(), resource.ReadRequest{State: createResp.State}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should refresh without error once visible: %v", readResp.Diagnostics)
	}
	if readResp.State.Raw.IsNull() {
		t.Fatal("Read should refresh to the full object state; state was removed")
	}
	var model savingPlanResourceModel
	readResp.Diagnostics.Append(readResp.State.Get(context.Background(), &model)...)
	if model.Id.ValueString() != "spn-test-123" {
		t.Fatalf("refreshed state should keep the saving plan id; got %q", model.Id.ValueString())
	}
	if model.RegionId.ValueString() != "cn-hangzhou" {
		t.Fatalf("refreshed state should carry the API region_id; got %q", model.RegionId.ValueString())
	}
	if model.CreateTime.ValueString() != "2026-01-01 00:00:00" {
		t.Fatalf("refreshed state should carry the API create_time; got %q", model.CreateTime.ValueString())
	}
	if model.InstanceFamily.ValueString() != "ecs.g6" {
		t.Fatalf("refreshed state should carry the API instance_family; got %q", model.InstanceFamily.ValueString())
	}
	if model.PlanType.ValueString() != "EcsCompute" {
		t.Fatalf("BSS SavingsType ECS must normalize to plan_type EcsCompute once visible; got %q", model.PlanType.ValueString())
	}
	if model.OfferingType.ValueString() != "AllUpfront" {
		t.Fatalf("BSS PayMode total must normalize to offering_type AllUpfront once visible; got %q", model.OfferingType.ValueString())
	}
	assertSavingPlanModelNoUnknowns(t, model)
}

// TestUnitSavingPlanReadNormalizesBSSEnumsToSchema verifies defect #2: the BSS
// QuerySavingsPlansInstance response carries SavingsType in {ECS, universal} and
// PayMode in {total, half, zero} — its OWN value space — while the resource
// schema validates plan_type against {EcsCompute, General} and offering_type
// against {AllUpfront, HalfUpfront, NoUpfront} (both OneOf + RequiresReplace).
// savingPlanState must translate the BSS enums onto the schema enums; writing the
// raw BSS value would violate OneOf and, since both are RequiresReplace, schedule
// a replacement on the next plan. The cases cover both plan types and all three
// payment modes with the real BSS strings (including the Chinese 通用型 form).
func TestUnitSavingPlanReadNormalizesBSSEnumsToSchema(t *testing.T) {
	cases := []struct {
		name             string
		savingsType      string
		payMode          string
		wantPlanType     string
		wantOfferingType string
	}{
		{"ecs_total", "ECS", "total", "EcsCompute", "AllUpfront"},
		{"ecs_half", "ECS", "half", "EcsCompute", "HalfUpfront"},
		{"ecs_zero", "ECS", "zero", "EcsCompute", "NoUpfront"},
		{"universal_total", "universal", "total", "General", "AllUpfront"},
		{"chinese_universal_half", "通用型", "half", "General", "HalfUpfront"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"Success": true, "Code": "Success", "Message": "",
						"Data": map[string]interface{}{"Items": []interface{}{
							map[string]interface{}{
								"InstanceId":     "spn-test-123",
								"SavingsType":    tc.savingsType,
								"PayMode":        tc.payMode,
								"Region":         "cn-hangzhou",
								"StartTimestamp": "2026-01-01 00:00:00",
							},
						}},
					})
				}
			})
			r := newSavingPlanResourceWithClient(t, client)
			initialState := savingPlanTestState(t, r)
			resp := resource.ReadResponse{State: initialState}
			r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("Read should not error: %v", resp.Diagnostics)
			}
			var model savingPlanResourceModel
			resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
			if model.PlanType.ValueString() != tc.wantPlanType {
				t.Fatalf("SavingsType %q should normalize to plan_type %q; got %q", tc.savingsType, tc.wantPlanType, model.PlanType.ValueString())
			}
			if model.OfferingType.ValueString() != tc.wantOfferingType {
				t.Fatalf("PayMode %q should normalize to offering_type %q; got %q", tc.payMode, tc.wantOfferingType, model.OfferingType.ValueString())
			}
			assertSavingPlanModelNoUnknowns(t, model)
		})
	}
}

// TestUnitSavingPlanReadRetainsLegalConfigOnUnrecognizedBSSEnum verifies defect #2's
// "unknown value not mis-changed" rule: when the BSS response carries a
// SavingsType/PayMode value the translation does not recognize, savingPlanState
// must NOT write that raw value into plan_type/offering_type (it would violate
// the schema's OneOf and, with RequiresReplace, schedule a replacement). The
// configured legal value is retained instead. The prior state carries the
// configured EcsCompute/AllUpfront; an unrecognized BSS "FutureType"/"quarterly"
// leaves them unchanged.
func TestUnitSavingPlanReadRetainsLegalConfigOnUnrecognizedBSSEnum(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") == "QuerySavingsPlansInstance" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{
					map[string]interface{}{
						"InstanceId":     "spn-test-123",
						"SavingsType":    "FutureType",
						"PayMode":        "quarterly",
						"Region":         "cn-hangzhou",
						"StartTimestamp": "2026-01-01 00:00:00",
					},
				}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	// Prior state carries the configured legal enums (as a successful Create /
	// prior Read would), so an unrecognized BSS enum must retain them rather
	// than mis-change state into a schema-violating value.
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	initialState := tfsdk.State{Schema: schemaResp.Schema}
	priorModel := savingPlanResourceModel{
		Id:           types.StringValue("spn-test-123"),
		SavingPlanId: types.StringValue("spn-test-123"),
		PlanType:     types.StringValue("EcsCompute"),
		OfferingType: types.StringValue("AllUpfront"),
		CreateTime:   types.StringValue("2026-01-01 00:00:00"),
		RegionId:     types.StringValue("cn-hangzhou"),
	}
	if diags := initialState.Set(context.Background(), &priorModel); diags.HasError() {
		t.Fatalf("state.Set: %v", diags)
	}
	resp := resource.ReadResponse{State: initialState}
	r.Read(context.Background(), resource.ReadRequest{State: initialState}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on an unrecognized BSS enum: %v", resp.Diagnostics)
	}
	var model savingPlanResourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	if model.PlanType.ValueString() != "EcsCompute" {
		t.Fatalf("an unrecognized SavingsType must retain the configured plan_type EcsCompute, not mis-change to %q", model.PlanType.ValueString())
	}
	if model.OfferingType.ValueString() != "AllUpfront" {
		t.Fatalf("an unrecognized PayMode must retain the configured offering_type AllUpfront, not mis-change to %q", model.OfferingType.ValueString())
	}
}

// --- Protocol-boundary (tfprotov6) closed loop -----------------------------
// The tests above drive the resource's Create/Read through the Framework
// tfsdk.Plan/State entry points, one layer below the real
// PlanResourceChange/ApplyResourceChange/ReadResource RPCs that `terraform
// plan/apply` dispatches. The two tests below run a real tfprotov6 provider
// server — the exact code path terraform uses — and round-trip
// Plan -> Apply -> Read -> next Plan through it, asserting that the NewState is
// fully known, the purchased ID is retained, and the next plan schedules NO
// RequiresReplace (the BSS enums normalized onto the schema enums make
// config-vs-state equal on every RequiresReplace attribute). This is the v6
// mirror of alicloud/provider/maximal_wire_test.go's lifecycle harness.

// savingPlanProtocolType is the tftypes Object for the saving plan resource,
// hand-written to match the Framework schema's implied type (17 attributes: all
// String except period which is Int64; no timeouts block). Mirrors the
// hand-written maximalSDKv2Type in maximal_wire_test.go rather than deriving
// from GetProviderSchema, to avoid configschema API churn.
var savingPlanProtocolType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"id":                  tftypes.String,
	"committed_amount":    tftypes.String,
	"plan_type":           tftypes.String,
	"offering_type":       tftypes.String,
	"purchase_method":     tftypes.String,
	"charge_type":         tftypes.String,
	"period":              tftypes.Number,
	"period_unit":         tftypes.String,
	"instance_family":     tftypes.String,
	"instance_family_set": tftypes.String,
	"saving_plan_name":    tftypes.String,
	"description":         tftypes.String,
	"start_time":          tftypes.String,
	"region_id":           tftypes.String,
	"saving_plan_id":      tftypes.String,
	"create_time":         tftypes.String,
	"payment_type":        tftypes.String,
}}

var savingPlanEmptyObjectType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}

// savingPlanProtocolTestProvider is a minimal Framework provider that serves
// only the saving plan resource, injecting the mock AliyunClient through
// Configure so fwadapt.ResourceBase.Configure receives it exactly as the real
// muxed provider wires it. Mirrors maximalFrameworkTestProvider.
type savingPlanProtocolTestProvider struct {
	client *connectivity.AliyunClient
}

var _ fwprovider.Provider = (*savingPlanProtocolTestProvider)(nil)

func (p *savingPlanProtocolTestProvider) Metadata(ctx context.Context, req fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "alicloud"
}

func (p *savingPlanProtocolTestProvider) Schema(ctx context.Context, req fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *savingPlanProtocolTestProvider) Configure(ctx context.Context, req fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) {
	resp.ResourceData = p.client
	resp.DataSourceData = p.client
}

func (p *savingPlanProtocolTestProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &savingPlanResource{} },
	}
}

func (p *savingPlanProtocolTestProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

// savingPlanProtocolServer builds a real tfprotov6 ProviderServer, configures it
// with an empty provider config, and returns the ready server. Its
// PlanResourceChange/ApplyResourceChange/ReadResource are the real Framework
// protocol entry points terraform uses, so a round trip through them proves the
// NewState/next-plan contract without hand-asserting on individual fields.
func savingPlanProtocolServer(t *testing.T, client *connectivity.AliyunClient) tfprotov6.ProviderServer {
	t.Helper()
	server, err := providerserver.NewProtocol6WithError(&savingPlanProtocolTestProvider{client: client})()
	if err != nil {
		t.Fatalf("NewProtocol6WithError: %s", err)
	}
	cfgResp, err := server.ConfigureProvider(context.Background(), &tfprotov6.ConfigureProviderRequest{
		Config: savingPlanProtoDynamic(t, savingPlanEmptyObjectType, tftypes.NewValue(savingPlanEmptyObjectType, map[string]tftypes.Value{})),
	})
	if err != nil {
		t.Fatalf("ConfigureProvider: %s", err)
	}
	if savingPlanProtoDiagsHasError(cfgResp.Diagnostics) {
		t.Fatalf("ConfigureProvider diagnostics: %v", cfgResp.Diagnostics)
	}
	return server
}

func savingPlanProtoObject(t *testing.T, ty tftypes.Object, attrs map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	for name := range attrs {
		if _, ok := ty.AttributeTypes[name]; !ok {
			t.Fatalf("%q is not an attribute of %s", name, ty)
		}
	}
	full := make(map[string]tftypes.Value, len(ty.AttributeTypes))
	for name, attrType := range ty.AttributeTypes {
		if v, ok := attrs[name]; ok {
			full[name] = v
			continue
		}
		full[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(ty, full)
}

func savingPlanProtoDynamic(t *testing.T, ty tftypes.Type, v tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	dv, err := tfprotov6.NewDynamicValue(ty, v)
	if err != nil {
		t.Fatalf("encoding a %s: %s", ty, err)
	}
	return &dv
}

func savingPlanProtoNull(t *testing.T, ty tftypes.Type) *tfprotov6.DynamicValue {
	t.Helper()
	return savingPlanProtoDynamic(t, ty, tftypes.NewValue(ty, nil))
}

func savingPlanProtoAttrs(t *testing.T, ty tftypes.Object, dv *tfprotov6.DynamicValue) map[string]tftypes.Value {
	t.Helper()
	if dv == nil {
		t.Fatal("the response carried no value")
	}
	val, err := dv.Unmarshal(ty)
	if err != nil {
		t.Fatalf("decoding a %s: %s", ty, err)
	}
	attrs := map[string]tftypes.Value{}
	if err := val.As(&attrs); err != nil {
		t.Fatalf("reading attributes of a %s: %s", ty, err)
	}
	return attrs
}

func savingPlanProtoString(t *testing.T, v tftypes.Value) string {
	t.Helper()
	if v.IsNull() {
		return ""
	}
	var s string
	if err := v.As(&s); err != nil {
		t.Fatalf("reading a string: %s", err)
	}
	return s
}

// savingPlanProtoDiagsHasError reports whether any tfprotov6 diagnostic is an
// error. tfprotov6 response Diagnostics is a []*tfprotov6.Diagnostic slice
// (not a diag.Diagnostics), so it has no HasError method; the real Framework
// protocol surfaces Plan/Apply/Read failures as error-severity diagnostics
// alongside a nil RPC error, so the tests gate on this rather than only on the
// returned error.
func savingPlanProtoDiagsHasError(diags []*tfprotov6.Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			return true
		}
	}
	return false
}

// savingPlanProtocolCreateConfig builds the Config / ProposedNewState
// DynamicValue for a create: the configured RequiresReplace-bearing attributes
// (plan_type, offering_type, committed_amount) plus the known
// conditional-required instance_family and the default-bearing
// purchase_method/charge_type; every other attribute is null (the Framework
// turns omitted Optional+Computed attributes into unknown during planning).
func savingPlanProtocolCreateConfig(t *testing.T) *tfprotov6.DynamicValue {
	t.Helper()
	return savingPlanProtoDynamic(t, savingPlanProtocolType, savingPlanProtoObject(t, savingPlanProtocolType, map[string]tftypes.Value{
		"committed_amount": tftypes.NewValue(tftypes.String, "0.1"),
		"plan_type":        tftypes.NewValue(tftypes.String, "EcsCompute"),
		"offering_type":    tftypes.NewValue(tftypes.String, "AllUpfront"),
		"purchase_method":  tftypes.NewValue(tftypes.String, "ByInstanceFamily"),
		"charge_type":      tftypes.NewValue(tftypes.String, "PrePaid"),
		"instance_family":  tftypes.NewValue(tftypes.String, "ecs.g6"),
	}))
}

// assertSavingPlanProtoStateKnown fails if any attribute of the given NewState
// DynamicValue is unknown — the Framework ApplyResourceChange NewState for a
// newly-created resource must be fully known, so a Create/Read that persists
// the paid ID on a refresh miss must not leak unknown Computed values.
func assertSavingPlanProtoStateKnown(t *testing.T, dv *tfprotov6.DynamicValue) {
	t.Helper()
	for name, v := range savingPlanProtoAttrs(t, savingPlanProtocolType, dv) {
		if !v.IsKnown() {
			t.Fatalf("state attribute %s must not be unknown; NewState must be fully known", name)
		}
	}
}

// savingPlanProtocolID reads the id attribute from a NewState DynamicValue.
func savingPlanProtocolID(t *testing.T, dv *tfprotov6.DynamicValue) string {
	t.Helper()
	return savingPlanProtoString(t, savingPlanProtoAttrs(t, savingPlanProtocolType, dv)["id"])
}

// savingPlanProtocolRoundTrip runs the full protocol closed loop on one config:
// Plan(create) -> Apply(create) -> Read -> next Plan. At every stage it asserts
// the NewState is fully known and the purchased ID is "spn-test-123", and at the
// end it asserts the next plan carries NO RequiresReplace (the normalized state
// equals the configured schema enums, so config-vs-state yields no diff on any
// RequiresReplace attribute). The server's Query handler must already be set up
// so the single Read returns the visible (normalized) shape.
func savingPlanProtocolRoundTrip(t *testing.T, server tfprotov6.ProviderServer, createConfig *tfprotov6.DynamicValue) {
	ctx := context.Background()

	createPlan, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         savingPlanTypeName,
		PriorState:       savingPlanProtoNull(t, savingPlanProtocolType),
		ProposedNewState: createConfig,
		Config:           createConfig,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (create): %s", err)
	}
	if savingPlanProtoDiagsHasError(createPlan.Diagnostics) {
		t.Fatalf("PlanResourceChange (create) diagnostics: %v", createPlan.Diagnostics)
	}

	created, err := server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
		TypeName:       savingPlanTypeName,
		PriorState:     savingPlanProtoNull(t, savingPlanProtocolType),
		PlannedState:   createPlan.PlannedState,
		Config:         createConfig,
		PlannedPrivate: createPlan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (create): %s", err)
	}
	if savingPlanProtoDiagsHasError(created.Diagnostics) {
		t.Fatalf("ApplyResourceChange (create) diagnostics: %v", created.Diagnostics)
	}
	assertSavingPlanProtoStateKnown(t, created.NewState)
	if got := savingPlanProtocolID(t, created.NewState); got != "spn-test-123" {
		t.Fatalf("create NewState id = %q, want spn-test-123", got)
	}

	// Read refreshes the persisted state through the real ReadResource RPC.
	refreshed, err := server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     savingPlanTypeName,
		CurrentState: created.NewState,
		Private:      created.Private,
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	if savingPlanProtoDiagsHasError(refreshed.Diagnostics) {
		t.Fatalf("ReadResource diagnostics: %v", refreshed.Diagnostics)
	}
	assertSavingPlanProtoStateKnown(t, refreshed.NewState)
	if got := savingPlanProtocolID(t, refreshed.NewState); got != "spn-test-123" {
		t.Fatalf("read NewState id = %q, want spn-test-123 (paid ID retained)", got)
	}

	// The decisive next plan: with the refreshed (normalized) state as
	// PriorState and the same configured values as ProposedNewState, the
	// Framework must schedule NO RequiresReplace. A non-empty RequiresReplace
	// here means a replacement would be scheduled on the next apply,
	// re-purchasing — the exact defect this guards against.
	nextPlan, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         savingPlanTypeName,
		PriorState:       refreshed.NewState,
		ProposedNewState: createConfig,
		Config:           createConfig,
		PriorPrivate:     refreshed.Private,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (next): %s", err)
	}
	if savingPlanProtoDiagsHasError(nextPlan.Diagnostics) {
		t.Fatalf("PlanResourceChange (next) diagnostics: %v", nextPlan.Diagnostics)
	}
	if len(nextPlan.RequiresReplace) != 0 {
		t.Fatalf("next plan must not require replacement (normalized state == configured enums); got RequiresReplace=%v", nextPlan.RequiresReplace)
	}
	if got := savingPlanProtocolID(t, nextPlan.PlannedState); got != "spn-test-123" {
		t.Fatalf("next plan PlannedState id = %q, want spn-test-123 (ID retained)", got)
	}
}

// TestUnitSavingPlanProtocolCreateReadNextPlanNoReplacement drives the saving
// plan resource through a real tfprotov6 provider server (the same
// PlanResourceChange/ApplyResourceChange/ReadResource RPCs terraform runs) on
// the initial-visible shape: Purchase returns the paid ID and the immediate
// Query refresh is visible, carrying the real BSS enum value space (SavingsType
// ECS / PayMode total). savingPlanState must normalize those onto the schema
// enums (EcsCompute/AllUpfront) so the next plan's config-vs-state comparison
// yields no diff on any RequiresReplace attribute. Asserts NewState fully known
// at every stage, the paid ID retained, and the next plan RequiresReplace empty.
func TestUnitSavingPlanProtocolCreateReadNextPlanNoReplacement(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			// Visible immediately with the real BSS enum value space, which
			// savingPlanState must normalize onto the schema enums.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{
					map[string]interface{}{
						"InstanceId":     "spn-test-123",
						"PoolValue":      "0.1",
						"SavingsType":    "ECS",
						"PayMode":        "total",
						"Region":         "cn-hangzhou",
						"StartTimestamp": "2026-01-01 00:00:00",
						"InstanceFamily": "ecs.g6",
						"StartTime":      "2026-01-01 00:00:00",
						"PaymentType":    "Subscription",
					},
				}},
			})
		}
	})
	server := savingPlanProtocolServer(t, client)
	savingPlanProtocolRoundTrip(t, server, savingPlanProtocolCreateConfig(t))
}

// TestUnitSavingPlanProtocolPropagationDelayNextPlanNoReplacement drives the
// real tfprotov6 server on the propagation-delay shape: Create's refresh is
// not-yet-visible (empty Items), so Create persists a fully-known state with the
// configured enums and null Computed-only attributes; a later Read is still
// not-visible (ID retained); once visible, Read fills the Computed attributes and
// normalizes the BSS enums. At every stage the NewState is fully known, the paid
// ID is retained, and the final next plan (against the visible normalized state)
// schedules NO RequiresReplace.
func TestUnitSavingPlanProtocolPropagationDelayNextPlanNoReplacement(t *testing.T) {
	var queryCount int
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "PurchaseSavingPlanOffering":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"SavingPlanIdSets": []interface{}{"spn-test-123"},
			})
		case "QuerySavingsPlansInstance":
			queryCount++
			if queryCount <= 2 {
				// Create's refresh + Read #1: not yet visible (empty Items).
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Success": true, "Code": "Success", "Message": "",
					"Data": map[string]interface{}{"Items": []interface{}{}},
				})
				return
			}
			// Read #2: visible, real BSS enum value space.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true, "Code": "Success", "Message": "",
				"Data": map[string]interface{}{"Items": []interface{}{
					map[string]interface{}{
						"InstanceId":     "spn-test-123",
						"PoolValue":      "0.1",
						"SavingsType":    "ECS",
						"PayMode":        "total",
						"Region":         "cn-hangzhou",
						"StartTimestamp": "2026-01-01 00:00:00",
						"InstanceFamily": "ecs.g6",
						"StartTime":      "2026-01-01 00:00:00",
						"PaymentType":    "Subscription",
					},
				}},
			})
		}
	})
	ctx := context.Background()
	server := savingPlanProtocolServer(t, client)
	createConfig := savingPlanProtocolCreateConfig(t)

	createPlan, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         savingPlanTypeName,
		PriorState:       savingPlanProtoNull(t, savingPlanProtocolType),
		ProposedNewState: createConfig,
		Config:           createConfig,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (create): %s", err)
	}
	if savingPlanProtoDiagsHasError(createPlan.Diagnostics) {
		t.Fatalf("PlanResourceChange (create) diagnostics: %v", createPlan.Diagnostics)
	}
	created, err := server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
		TypeName:       savingPlanTypeName,
		PriorState:     savingPlanProtoNull(t, savingPlanProtocolType),
		PlannedState:   createPlan.PlannedState,
		Config:         createConfig,
		PlannedPrivate: createPlan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (create): %s", err)
	}
	if savingPlanProtoDiagsHasError(created.Diagnostics) {
		t.Fatalf("ApplyResourceChange (create) diagnostics: %v", created.Diagnostics)
	}
	assertSavingPlanProtoStateKnown(t, created.NewState)
	if got := savingPlanProtocolID(t, created.NewState); got != "spn-test-123" {
		t.Fatalf("create NewState id = %q, want spn-test-123", got)
	}

	// Read #1: still not visible. ID retained, NewState fully known.
	read1, err := server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     savingPlanTypeName,
		CurrentState: created.NewState,
		Private:      created.Private,
	})
	if err != nil {
		t.Fatalf("ReadResource #1: %s", err)
	}
	if savingPlanProtoDiagsHasError(read1.Diagnostics) {
		t.Fatalf("ReadResource #1 diagnostics: %v", read1.Diagnostics)
	}
	assertSavingPlanProtoStateKnown(t, read1.NewState)
	if got := savingPlanProtocolID(t, read1.NewState); got != "spn-test-123" {
		t.Fatalf("read #1 id = %q, want spn-test-123 (paid ID retained while not visible)", got)
	}

	// Read #2: now visible. Computed filled, BSS enums normalized, ID stable.
	read2, err := server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     savingPlanTypeName,
		CurrentState: read1.NewState,
		Private:      read1.Private,
	})
	if err != nil {
		t.Fatalf("ReadResource #2: %s", err)
	}
	if savingPlanProtoDiagsHasError(read2.Diagnostics) {
		t.Fatalf("ReadResource #2 diagnostics: %v", read2.Diagnostics)
	}
	assertSavingPlanProtoStateKnown(t, read2.NewState)
	if got := savingPlanProtocolID(t, read2.NewState); got != "spn-test-123" {
		t.Fatalf("read #2 id = %q, want spn-test-123", got)
	}

	// Next plan with the visible normalized state as PriorState: no
	// RequiresReplace (config == normalized state on every RequiresReplace
	// attribute).
	nextPlan, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         savingPlanTypeName,
		PriorState:       read2.NewState,
		ProposedNewState: createConfig,
		Config:           createConfig,
		PriorPrivate:     read2.Private,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (next): %s", err)
	}
	if savingPlanProtoDiagsHasError(nextPlan.Diagnostics) {
		t.Fatalf("PlanResourceChange (next) diagnostics: %v", nextPlan.Diagnostics)
	}
	if len(nextPlan.RequiresReplace) != 0 {
		t.Fatalf("next plan must not require replacement (normalized state == configured enums); got RequiresReplace=%v", nextPlan.RequiresReplace)
	}
	if got := savingPlanProtocolID(t, nextPlan.PlannedState); got != "spn-test-123" {
		t.Fatalf("next plan PlannedState id = %q, want spn-test-123 (ID retained)", got)
	}
}
