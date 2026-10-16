// Package ecs holds the Framework-native ECS resources. These are the no-cloud
// unit tests for DiagnosticReport: they drive the real Create/Read/Delete/Import
// methods against an httptest server, exercising the failure paths that the
// acceptance test cannot deterministically reach — Create persisting the ReportId
// before a Describe that fails, the strict diagnosticReportId, the legal-missing
// vs bad-response distinction in describe, the additional_options JSON encoding
// and Read closed-loop, the delete poll and the importer. They do not touch the
// cloud; helper.NeedRetry returns false for the 400 errors used here, so none of
// these paths retry-loop or hang.
package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// diagnosticReportTestClient wires a real *connectivity.AliyunClient at an httptest
// server and routes Ecs RpcPost traffic there via endpoints.Store("ecs"). It mirrors
// the ecsSnapshotTestClient recipe: an in-memory access-key credential, a Config with
// SkipRegionValidation and a fresh endpoints map, and NO_PROXY set to the server host
// so the HTTP client dials the test server directly rather than any real endpoint.
func diagnosticReportTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build credential: %v", err)
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
		t.Fatalf("build client: %v", err)
	}
	// Route Ecs RpcPost (rpcRequest lowercases the product to "ecs" and reads
	// endpoints.Load("ecs")) to the test server.
	endpoints.Store("ecs", endpoint)
	return client
}

// diagnosticReportTFType is the tftypes object the schema produces: six string
// attributes plus a map[string]string additional_options. It must name the same
// attributes as the schema for Plan/State Get/Set to round-trip.
func diagnosticReportTFType() tftypes.Type {
	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":                 tftypes.String,
		"resource_id":        tftypes.String,
		"metric_set_id":      tftypes.String,
		"start_time":         tftypes.String,
		"end_time":           tftypes.String,
		"additional_options": tftypes.Map{ElementType: tftypes.String},
	}}
}

// additionalOptionsValue builds a tftypes map value from a Go map; a nil argument
// yields a null map (Create must skip encoding it), and a non-nil empty map yields
// an empty-but-present map.
func additionalOptionsValue(kv map[string]string) tftypes.Value {
	if kv == nil {
		return tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil)
	}
	elements := map[string]tftypes.Value{}
	for k, v := range kv {
		elements[k] = tftypes.NewValue(tftypes.String, v)
	}
	return tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, elements)
}

// diagnosticReportPlanValue builds a plan object: id is null (Computed, set on Create),
// resource_id is required, the metric/time fields are the SDKv2 defaults, and
// additional_options carries the given value.
func diagnosticReportPlanValue(resourceID string, additional tftypes.Value) tftypes.Value {
	return tftypes.NewValue(diagnosticReportTFType(), map[string]tftypes.Value{
		"id":                 tftypes.NewValue(tftypes.String, nil),
		"resource_id":        tftypes.NewValue(tftypes.String, resourceID),
		"metric_set_id":      tftypes.NewValue(tftypes.String, "dms-instancedefault"),
		"start_time":         tftypes.NewValue(tftypes.String, nil),
		"end_time":           tftypes.NewValue(tftypes.String, nil),
		"additional_options": additional,
	})
}

// diagnosticReportStateValue builds a state object carrying a known id and the given
// additional_options, used to pre-populate Read's prior state.
func diagnosticReportStateValue(id string, additional tftypes.Value) tftypes.Value {
	return tftypes.NewValue(diagnosticReportTFType(), map[string]tftypes.Value{
		"id":                 tftypes.NewValue(tftypes.String, id),
		"resource_id":        tftypes.NewValue(tftypes.String, "i-test"),
		"metric_set_id":      tftypes.NewValue(tftypes.String, "dms-instancedefault"),
		"start_time":         tftypes.NewValue(tftypes.String, nil),
		"end_time":           tftypes.NewValue(tftypes.String, nil),
		"additional_options": additional,
	})
}

// nullAttributesObject builds a valid (non-null) object whose every attribute is null.
// SetAttribute works on such an object, so it is the realistic "empty" state the
// framework hands ImportState before the passthrough sets the id.
func nullAttributesObject() tftypes.Value {
	return tftypes.NewValue(diagnosticReportTFType(), map[string]tftypes.Value{
		"id":                 tftypes.NewValue(tftypes.String, nil),
		"resource_id":        tftypes.NewValue(tftypes.String, nil),
		"metric_set_id":      tftypes.NewValue(tftypes.String, nil),
		"start_time":         tftypes.NewValue(tftypes.String, nil),
		"end_time":           tftypes.NewValue(tftypes.String, nil),
		"additional_options": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
	})
}

// newDiagnosticReportTestResource instantiates the resource, fetches its schema, and
// injects the mock client through fwadapt.ResourceBase.Configure so Client() returns
// the in-memory AliyunClient for the duration of the test.
func newDiagnosticReportTestResource(t *testing.T, client *connectivity.AliyunClient) (*diagnosticReportResource, schema.Schema) {
	t.Helper()
	r := &diagnosticReportResource{}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema error: %v", schemaResp.Diagnostics)
	}

	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configure error: %v", configureResp.Diagnostics)
	}
	return r, schemaResp.Schema
}

// assertStateID reads resp.State into a model and asserts the id. Used after Create/Read
// to confirm the ReportId is (or is not) retained.
func assertStateID(t *testing.T, ctx context.Context, state tfsdk.State, want string) {
	t.Helper()
	var model diagnosticReportResourceModel
	if diags := state.Get(ctx, &model); diags.HasError() {
		t.Fatalf("state.Get: %v", diags)
	}
	if got := model.Id.ValueString(); got != want {
		t.Fatalf("state id = %q, want %q", got, want)
	}
}

// assertStateAdditionalOptions reads resp.State into a model and asserts the
// additional_options map decodes to want, proving the write-only map survives the
// Create/Read round-trip.
func assertStateAdditionalOptions(t *testing.T, ctx context.Context, state tfsdk.State, want map[string]string) {
	t.Helper()
	var model diagnosticReportResourceModel
	if diags := state.Get(ctx, &model); diags.HasError() {
		t.Fatalf("state.Get: %v", diags)
	}
	got := map[string]string{}
	if diags := model.AdditionalOptions.ElementsAs(ctx, &got, false); diags.HasError() {
		t.Fatalf("ElementsAs: %v", diags)
	}
	if len(got) != len(want) {
		t.Fatalf("additional_options = %v, want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("additional_options[%q] = %q, want %q", k, got[k], v)
		}
	}
}

// respondJSON writes a JSON body with the given status; status 200 omits WriteHeader so
// the default 200 is used.
func respondJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	fmt.Fprint(w, body)
}

// branchDiagnosticHandler routes CreateDiagnosticReport and DescribeDiagnosticReportAttributes
// to separate bodies, defaulting any other action to a 400 error so an unexpected call
// surfaces instead of silently passing.
func branchDiagnosticHandler(t *testing.T, createBody, describeBody string, describeStatus int, onCreate func(form url.Values)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch r.Form.Get("Action") {
		case "CreateDiagnosticReport":
			if onCreate != nil {
				onCreate(r.Form)
			}
			respondJSON(w, http.StatusOK, createBody)
		case "DescribeDiagnosticReportAttributes":
			respondJSON(w, describeStatus, describeBody)
		case "DeleteDiagnosticReports":
			respondJSON(w, http.StatusOK, "{}")
		default:
			respondJSON(w, http.StatusBadRequest, `{"Code":"InternalError","Message":"unexpected action","RequestId":"t"}`)
		}
	}
}

// TestUnitDiagnosticReportIdStrict exercises the strict ReportId extractor: only a
// non-empty string is accepted. A missing key, nil, or any non-string type (number,
// map, slice) is rejected so a malformed response cannot yield a pseudo-ID — the
// regression the human located against the SDKv2 fmt.Sprint fallback.
func TestUnitDiagnosticReportIdStrict(t *testing.T) {
	cases := []struct {
		name     string
		response map[string]interface{}
		wantID   string
		wantErr  bool
	}{
		{"string", map[string]interface{}{"ReportId": "dr-1"}, "dr-1", false},
		{"missing", map[string]interface{}{"RequestId": "t"}, "", true},
		{"nil", map[string]interface{}{"ReportId": nil}, "", true},
		{"number", map[string]interface{}{"ReportId": 42}, "", true},
		{"map", map[string]interface{}{"ReportId": map[string]interface{}{"k": "v"}}, "", true},
		{"slice", map[string]interface{}{"ReportId": []interface{}{"a"}}, "", true},
		{"empty_string", map[string]interface{}{"ReportId": ""}, "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := diagnosticReportId(c.response)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error, got id=%q", id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != c.wantID {
				t.Fatalf("id = %q, want %q", id, c.wantID)
			}
		})
	}
}

// TestUnitDiagnosticReportDescribe drives describeDiagnosticReport against a test server
// to prove the legal-missing vs bad-response distinction. The two legal "gone" paths
// (403 ReportNotExist and the 200 envelope with only RequestId) return (nil, nil); a
// wrong-type, empty, or mismatched ReportId returns an error so Read/Delete retain
// state rather than silently removing the resource.
func TestUnitDiagnosticReportDescribe(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		body       string
		wantObject bool
		wantErr    bool
	}{
		{"report_not_exist", http.StatusBadRequest, `{"Code":"InvalidParameter.ReportNotExist","Message":"gone","RequestId":"t"}`, false, false},
		{"empty_envelope", http.StatusOK, `{"RequestId":"t"}`, false, false},
		{"empty_object", http.StatusOK, `{}`, false, true},
		{"unrelated_field", http.StatusOK, `{"ResourceId":"i-other"}`, false, true},
		{"null_reportid", http.StatusOK, `{"ReportId":null,"RequestId":"t"}`, false, true},
		{"bad_type_number", http.StatusOK, `{"ReportId":42,"RequestId":"t"}`, false, true},
		{"empty_string", http.StatusOK, `{"ReportId":"","RequestId":"t"}`, false, true},
		{"id_mismatch", http.StatusOK, `{"ReportId":"other","RequestId":"t"}`, false, true},
		{"match", http.StatusOK, `{"ReportId":"dr-x","ResourceId":"i-test","RequestId":"t"}`, true, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := diagnosticReportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				_ = r.ParseForm()
				respondJSON(w, c.status, c.body)
			})
			obj, err := describeDiagnosticReport(context.Background(), client, "dr-x")
			if c.wantErr && err == nil {
				t.Fatalf("expected error, got object=%v", obj)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.wantObject && obj == nil {
				t.Fatalf("expected object, got nil")
			}
			if !c.wantObject && obj != nil {
				t.Fatalf("expected nil object, got %v", obj)
			}
		})
	}
}

// TestUnitDiagnosticReportCreateRetainsIDOnDescribeError proves the central fix:
// Create persists the ReportId into state BEFORE the post-create Describe, so a
// Describe that returns a generic error (not ReportNotExist) loses neither the id
// nor the write-only additional_options. The diagnostics carry a warning, not an
// error, so Terraform keeps the resource.
func TestUnitDiagnosticReportCreateRetainsIDOnDescribeError(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"dr-1","RequestId":"t"}`,
		`{"Code":"InternalError","Message":"boom","RequestId":"t"}`,
		http.StatusBadRequest, nil))
	r, sch := newDiagnosticReportTestResource(t, client)

	plan := tfsdk.Plan{Raw: diagnosticReportPlanValue("i-test", additionalOptionsValue(map[string]string{"key": "value"})), Schema: sch}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned an error, expected only a warning: %v", resp.Diagnostics)
	}
	if len(resp.Diagnostics) == 0 {
		t.Fatalf("Create returned no warning, expected one for the failed refresh")
	}
	assertStateID(t, ctx, resp.State, "dr-1")
	assertStateAdditionalOptions(t, ctx, resp.State, map[string]string{"key": "value"})
}

// TestUnitDiagnosticReportCreateRetainsIDOnDescribeEmpty covers the other legal
// "not yet visible" path: Describe returns a 200 envelope with no ReportId. Create
// retains the id and warns, exactly like the error case.
func TestUnitDiagnosticReportCreateRetainsIDOnDescribeEmpty(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"dr-2","RequestId":"t"}`,
		`{"RequestId":"t"}`,
		http.StatusOK, nil))
	r, sch := newDiagnosticReportTestResource(t, client)

	plan := tfsdk.Plan{Raw: diagnosticReportPlanValue("i-test", additionalOptionsValue(map[string]string{"key": "value"})), Schema: sch}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned an error, expected only a warning: %v", resp.Diagnostics)
	}
	if len(resp.Diagnostics) == 0 {
		t.Fatalf("Create returned no warning, expected one for the not-yet-visible refresh")
	}
	assertStateID(t, ctx, resp.State, "dr-2")
	assertStateAdditionalOptions(t, ctx, resp.State, map[string]string{"key": "value"})
}

// TestUnitDiagnosticReportCreateSendsAdditionalOptions asserts the write-only map is
// actually serialized onto the HTTP request as a JSON-encoded string (the real API
// contract: style json, a single AdditionalOptions form field holding a JSON object),
// and that it survives into state via the post-create refresh.
func TestUnitDiagnosticReportCreateSendsAdditionalOptions(t *testing.T) {
	ctx := context.Background()
	var captured string
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"dr-3","RequestId":"t"}`,
		`{"ReportId":"dr-3","ResourceId":"i-test","RequestId":"t"}`,
		http.StatusOK, func(form url.Values) {
			captured = form.Get("AdditionalOptions")
		}))
	r, sch := newDiagnosticReportTestResource(t, client)

	plan := tfsdk.Plan{Raw: diagnosticReportPlanValue("i-test", additionalOptionsValue(map[string]string{"key": "value"})), Schema: sch}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create error: %v", resp.Diagnostics)
	}
	// The map must be marshaled to a JSON object string, not flattened or dropped.
	want, _ := json.Marshal(map[string]string{"key": "value"})
	if captured != string(want) {
		t.Fatalf("AdditionalOptions form value = %q, want %q", captured, string(want))
	}
	assertStateID(t, ctx, resp.State, "dr-3")
	assertStateAdditionalOptions(t, ctx, resp.State, map[string]string{"key": "value"})
}

// TestUnitDiagnosticReportCreateSkipsNullAdditionalOptions asserts a null map is not
// sent and the resource still creates cleanly.
func TestUnitDiagnosticReportCreateSkipsNullAdditionalOptions(t *testing.T) {
	ctx := context.Background()
	var captured string
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"dr-4","RequestId":"t"}`,
		`{"ReportId":"dr-4","ResourceId":"i-test","RequestId":"t"}`,
		http.StatusOK, func(form url.Values) {
			captured = form.Get("AdditionalOptions")
		}))
	r, sch := newDiagnosticReportTestResource(t, client)

	plan := tfsdk.Plan{Raw: diagnosticReportPlanValue("i-test", additionalOptionsValue(nil)), Schema: sch}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create error: %v", resp.Diagnostics)
	}
	if captured != "" {
		t.Fatalf("AdditionalOptions should not be sent for a null map, got %q", captured)
	}
	assertStateID(t, ctx, resp.State, "dr-4")
}

// TestUnitDiagnosticReportCreateSkipsUnknownAdditionalOptions asserts an unknown map
// (the Computed/unknown plan Terraform sends before the value is known) is not
// encoded and does not error.
func TestUnitDiagnosticReportCreateSkipsUnknownAdditionalOptions(t *testing.T) {
	ctx := context.Background()
	var captured string
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"dr-4b","RequestId":"t"}`,
		`{"ReportId":"dr-4b","ResourceId":"i-test","RequestId":"t"}`,
		http.StatusOK, func(form url.Values) {
			captured = form.Get("AdditionalOptions")
		}))
	r, sch := newDiagnosticReportTestResource(t, client)

	unknown := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, tftypes.UnknownValue)
	plan := tfsdk.Plan{Raw: diagnosticReportPlanValue("i-test", unknown), Schema: sch}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create error: %v", resp.Diagnostics)
	}
	if captured != "" {
		t.Fatalf("AdditionalOptions should not be sent for an unknown map, got %q", captured)
	}
	assertStateID(t, ctx, resp.State, "dr-4b")
}

// TestUnitDiagnosticReportReadBadResponseRetainsID proves Read does not lose the id
// when Describe returns a malformed response (wrong-type ReportId). The framework
// pre-populates resp.State with the prior state; Read's AddError+return leaves it
// intact, so the resource is retained rather than written to a null id — the bug the
// human located.
func TestUnitDiagnosticReportReadBadResponseRetainsID(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"unused","RequestId":"t"}`,
		`{"ReportId":42,"RequestId":"t"}`,
		http.StatusOK, nil))
	r, sch := newDiagnosticReportTestResource(t, client)

	prior := diagnosticReportStateValue("dr-5", additionalOptionsValue(map[string]string{"key": "value"}))
	// Mirror the framework server: resp.State is pre-populated with the prior state.
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: prior, Schema: sch}}
	r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Raw: prior, Schema: sch}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("Read returned no error, expected one for a wrong-type ReportId")
	}
	assertStateID(t, ctx, resp.State, "dr-5")
	assertStateAdditionalOptions(t, ctx, resp.State, map[string]string{"key": "value"})
}

// TestUnitDiagnosticReportReadMalformedRetainsID proves Read does not lose the id when
// Describe returns one of the malformed 200 shapes that are NOT the legal only-RequestId
// envelope: an empty object ({}), an object carrying an unrelated business field
// ({"ResourceId":"i-other"}), or an explicitly null ReportId
// ({"ReportId":null,"RequestId":"t"}). describeDiagnosticReport must surface each as an
// error so the framework retains the prior state; treating any of them as "gone" would
// silently remove a real report. The prior id and the write-only additional_options survive.
func TestUnitDiagnosticReportReadMalformedRetainsID(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name string
		body string
	}{
		{"empty_object", `{}`},
		{"unrelated_field", `{"ResourceId":"i-other"}`},
		{"null_reportid", `{"ReportId":null,"RequestId":"t"}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
				`{"ReportId":"unused","RequestId":"t"}`,
				c.body,
				http.StatusOK, nil))
			r, sch := newDiagnosticReportTestResource(t, client)

			prior := diagnosticReportStateValue("dr-root", additionalOptionsValue(map[string]string{"key": "value"}))
			resp := &resource.ReadResponse{State: tfsdk.State{Raw: prior, Schema: sch}}
			r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Raw: prior, Schema: sch}}, resp)

			if !resp.Diagnostics.HasError() {
				t.Fatalf("Read returned no error, expected one for the malformed %s response", c.name)
			}
			assertStateID(t, ctx, resp.State, "dr-root")
			assertStateAdditionalOptions(t, ctx, resp.State, map[string]string{"key": "value"})
		})
	}
}

// TestUnitDiagnosticReportReadLegalGone proves Read removes the resource from state
// only on a legal "gone" signal (403 ReportNotExist), nilling the state without an
// error. This is the legitimate NotFound path, distinct from a bad response.
func TestUnitDiagnosticReportReadLegalGone(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"unused","RequestId":"t"}`,
		`{"Code":"InvalidParameter.ReportNotExist","Message":"gone","RequestId":"t"}`,
		http.StatusBadRequest, nil))
	r, sch := newDiagnosticReportTestResource(t, client)

	prior := diagnosticReportStateValue("dr-6", additionalOptionsValue(map[string]string{"key": "value"}))
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: prior, Schema: sch}}
	r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Raw: prior, Schema: sch}}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned an error on a legal-gone report: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatalf("expected removed (null) state, got %s", resp.State.Raw.String())
	}
}

// TestUnitDiagnosticReportReadLegalGoneEmpty covers the second legal-gone path: the
// 200 envelope with only RequestId. Read removes the resource without error.
func TestUnitDiagnosticReportReadLegalGoneEmpty(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"unused","RequestId":"t"}`,
		`{"RequestId":"t"}`,
		http.StatusOK, nil))
	r, sch := newDiagnosticReportTestResource(t, client)

	prior := diagnosticReportStateValue("dr-6b", additionalOptionsValue(map[string]string{"key": "value"}))
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: prior, Schema: sch}}
	r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Raw: prior, Schema: sch}}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned an error on a legal-gone empty envelope: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatalf("expected removed (null) state, got %s", resp.State.Raw.String())
	}
}

// TestUnitDiagnosticReportReadSuccess proves a matching Describe refreshes state and
// preserves the write-only additional_options from the prior state (the API never
// returns it).
func TestUnitDiagnosticReportReadSuccess(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, branchDiagnosticHandler(t,
		`{"ReportId":"unused","RequestId":"t"}`,
		`{"ReportId":"dr-7","ResourceId":"i-test","MetricSetId":"dms-instancedefault","RequestId":"t"}`,
		http.StatusOK, nil))
	r, sch := newDiagnosticReportTestResource(t, client)

	prior := diagnosticReportStateValue("dr-7", additionalOptionsValue(map[string]string{"key": "value"}))
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: prior, Schema: sch}}
	r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Raw: prior, Schema: sch}}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read error: %v", resp.Diagnostics)
	}
	assertStateID(t, ctx, resp.State, "dr-7")
	assertStateAdditionalOptions(t, ctx, resp.State, map[string]string{"key": "value"})
}

// TestUnitDiagnosticReportDeleteGone proves Delete succeeds when the report is already
// gone after the Delete RPC: the post-delete Describe returns ReportNotExist and the
// poll exits without error.
func TestUnitDiagnosticReportDeleteGone(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		switch r.Form.Get("Action") {
		case "DeleteDiagnosticReports":
			respondJSON(w, http.StatusOK, "{}")
		case "DescribeDiagnosticReportAttributes":
			respondJSON(w, http.StatusBadRequest, `{"Code":"InvalidParameter.ReportNotExist","Message":"gone","RequestId":"t"}`)
		default:
			respondJSON(w, http.StatusBadRequest, `{"Code":"InternalError","Message":"unexpected","RequestId":"t"}`)
		}
	})
	r, sch := newDiagnosticReportTestResource(t, client)

	state := tfsdk.State{Raw: diagnosticReportStateValue("dr-8", additionalOptionsValue(nil)), Schema: sch}
	resp := &resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete error: %v", resp.Diagnostics)
	}
}

// TestUnitDiagnosticReportDeleteAlreadyNotFound proves Delete treats a ReportNotExist
// from the Delete RPC itself as a no-op success (the resource is already gone), then
// the poll confirms absence.
func TestUnitDiagnosticReportDeleteAlreadyNotFound(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		// Both Delete and Describe report not-exist: the resource was already gone.
		respondJSON(w, http.StatusBadRequest, `{"Code":"InvalidParameter.ReportNotExist","Message":"gone","RequestId":"t"}`)
	})
	r, sch := newDiagnosticReportTestResource(t, client)

	state := tfsdk.State{Raw: diagnosticReportStateValue("dr-9", additionalOptionsValue(nil)), Schema: sch}
	resp := &resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete error on already-not-found: %v", resp.Diagnostics)
	}
}

// TestUnitDiagnosticReportWaitForDeletionTimeout proves the delete poll returns a
// timeout error (not an infinite loop) when the report never disappears. The timeout
// is a single nanosecond, smaller than one localhost RPC round-trip, so the deadline
// fires after the first Describe without ever sleeping.
func TestUnitDiagnosticReportWaitForDeletionTimeout(t *testing.T) {
	client := diagnosticReportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		// Report still exists.
		respondJSON(w, http.StatusOK, `{"ReportId":"dr-t","RequestId":"t"}`)
	})
	err := waitForDiagnosticReportDeletion(context.Background(), client, "dr-t", time.Nanosecond)
	if err == nil {
		t.Fatal("expected a timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "waiting for deletion") {
		t.Fatalf("expected a timeout error, got: %v", err)
	}
}

// TestUnitDiagnosticReportWaitForDeletionGone proves the poll exits cleanly and fast
// when the report is already gone (ReportNotExist on the first Describe).
func TestUnitDiagnosticReportWaitForDeletionGone(t *testing.T) {
	client := diagnosticReportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		respondJSON(w, http.StatusBadRequest, `{"Code":"InvalidParameter.ReportNotExist","Message":"gone","RequestId":"t"}`)
	})
	if err := waitForDiagnosticReportDeletion(context.Background(), client, "dr-g", 5*time.Second); err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

// TestUnitDiagnosticReportImport proves ImportState passes the import id through to the
// id attribute, leaving the rest null.
func TestUnitDiagnosticReportImport(t *testing.T) {
	ctx := context.Background()
	client := diagnosticReportTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Import does not call RPC, but the server must exist.
	})
	r, sch := newDiagnosticReportTestResource(t, client)

	resp := &resource.ImportStateResponse{State: tfsdk.State{Raw: nullAttributesObject(), Schema: sch}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "dr-10"}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState error: %v", resp.Diagnostics)
	}
	assertStateID(t, ctx, resp.State, "dr-10")
}
