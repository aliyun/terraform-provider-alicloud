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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// savingPlanTestPlan builds a valid Framework plan with the required
// committed_amount attribute populated so Create does not reject the plan.
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
	}
	if diags := plan.Set(context.Background(), &model); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	return plan
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

// TestUnitSavingPlanReadRemovesStateOnNotFound verifies that when the API reports
// the saving plan is gone (Query returns no items), Read removes the resource
// from state rather than retaining a stale entry. This covers the failure path
// the reviewer asked about: after the cloud resource disappears (expired or
// cancelled via console), a subsequent Read must not retain the ID in state.
func TestUnitSavingPlanReadRemovesStateOnNotFound(t *testing.T) {
	client := savingPlanTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "QuerySavingsPlansInstance":
			// Official BSS success envelope: Success=true, Code=Success, Message
			// may be empty, with empty Items — a legitimate "no such plan" for
			// this instance id, which querySavingPlan treats as not-found and
			// Read RemoveResource's.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true,
				"Code":    "Success",
				"Message": "",
				"Data":    map[string]interface{}{"Items": []interface{}{}},
			})
		}
	})
	r := newSavingPlanResourceWithClient(t, client)
	resp := resource.ReadResponse{State: savingPlanTestState(t, r)}
	r.Read(context.Background(), resource.ReadRequest{
		State: savingPlanTestState(t, r),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on not-found: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatalf("Read should remove resource from state on not-found; state=%v", resp.State.Raw)
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
// genuine success.
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
						"SavingsType":    "EcsCompute",
						"PayMode":        "AllUpfront",
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
