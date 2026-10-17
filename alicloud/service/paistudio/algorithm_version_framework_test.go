package paistudio

// These TestUnits drive the actual algorithmVersionResource Framework
// methods (Create/Read/Update/Delete/ImportState) — not the AliyunClient.Roa*
// helpers or the algorithmVersionState/validate helpers that
// algorithm_version_http_test.go exercises directly. They build a real
// terraform-plugin-framework provider server over an offline AliyunClient
// whose ROA requests target an in-process httptest mock, then drive the
// resource through the protocol-v5 PlanResourceChange / ApplyResourceChange /
// ReadResource / ImportResourceState RPCs the real provider uses.
//
// This covers the behaviour the helper-only wire tests cannot prove:
//   - Create persists the composite id (the recovery anchor) BEFORE the post-
//     create Get, so a Get 500 / not-yet-visible (404 NoSuchObject) / malformed
//     response cannot drop the already-created resource.
//   - Read removes the resource on a true NotFound (typed NoSuchObject) and
//     preserves prior state + emits diagnostics on a 5xx / malformed /
//     mismatched-identity response (no silent drop, no overwrite).
//   - Update issues the RoaPut with the planned spec and refreshes state.
//   - Delete preserves state on a transient error and is idempotent on a true
//     NotFound (the resource is destroyed with no error).
//   - ID-only Import sets only id + the two identifiers; a following Read fills
//     the computed fields and a no-change plan produces no RequiresReplace.
//
// The production ROA wire (roaRequest -> roa.NewClient -> DoRequest) is the
// same path the http_test exercises; it never reaches (*openapi.Client).Execute
// (the nil Spi.ModifyConfiguration site), so these tests run fully offline.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// algorithmVersionFrameworkType is the tftypes shape of the
// alicloud_pai_algorithm_version schema: every attribute is a string, so the
// plan/apply/read round-trips are uniform and the assertions only need to read
// strings back.
var algorithmVersionFrameworkType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"id":                 tftypes.String,
	"algorithm_id":       tftypes.String,
	"algorithm_version":  tftypes.String,
	"algorithm_spec":     tftypes.String,
	"tenant_id":          tftypes.String,
	"user_id":            tftypes.String,
	"algorithm_provider": tftypes.String,
	"algorithm_name":     tftypes.String,
	"gmt_modified_time":  tftypes.String,
	"region_id":          tftypes.String,
	"create_time":        tftypes.String,
}}

var paiFrameworkEmptyObjectType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}

const (
	mockSpecCreate = `{"JobType":"TfJob","Command":"python main.py","Image":"reg/echo:1.0"}`
	mockSpecUpdate = `{"JobType":"TfJob","Command":"python train.py","Image":"reg/echo:2.0"}`
)

// paiFrameworkTestProvider is a minimal Framework provider that registers the
// real algorithmVersionResource and injects the offline AliyunClient as the
// provider data; the resource's embedded fwadapt.ResourceBase.Configure stashes
// it so r.Client() resolves inside Create/Read/Update/Delete.
type paiFrameworkTestProvider struct {
	client *connectivity.AliyunClient
}

var _ fwprovider.Provider = (*paiFrameworkTestProvider)(nil)

func (p *paiFrameworkTestProvider) Metadata(_ context.Context, _ fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "alicloud"
}

func (p *paiFrameworkTestProvider) Schema(_ context.Context, _ fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *paiFrameworkTestProvider) Configure(_ context.Context, _ fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) {
	resp.ResourceData = p.client
	resp.DataSourceData = p.client
}

func (p *paiFrameworkTestProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAlgorithmVersionResource,
	}
}

func (p *paiFrameworkTestProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// newPaiFrameworkServer builds a configured protocol-v5 provider server backed
// by the offline client. ConfigureProvider runs the provider Configure above,
// so every resource instance built for a later RPC receives the client.
func newPaiFrameworkServer(t *testing.T, client *connectivity.AliyunClient) tfprotov5.ProviderServer {
	t.Helper()
	server := providerserver.NewProtocol5(&paiFrameworkTestProvider{client: client})()
	resp, err := server.ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{
		Config: dynamic(t, paiFrameworkEmptyObjectType, tftypes.NewValue(paiFrameworkEmptyObjectType, map[string]tftypes.Value{})),
	})
	if err != nil {
		t.Fatalf("ConfigureProvider: %s", err)
	}
	requireNoDiagErrors(t, "ConfigureProvider", resp.Diagnostics)
	return server
}

// --- protocol-v5 value helpers (mirrors alicloud/provider/maximal_wire_test.go) ---

func object(t *testing.T, ty tftypes.Object, attrs map[string]tftypes.Value) tftypes.Value {
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

func dynamic(t *testing.T, ty tftypes.Type, v tftypes.Value) *tfprotov5.DynamicValue {
	t.Helper()
	dv, err := tfprotov5.NewDynamicValue(ty, v)
	if err != nil {
		t.Fatalf("encoding a %s: %s", ty, err)
	}
	return &dv
}

func nullOf(t *testing.T, ty tftypes.Type) *tfprotov5.DynamicValue {
	t.Helper()
	return dynamic(t, ty, tftypes.NewValue(ty, nil))
}

func attrsOf(t *testing.T, ty tftypes.Type, dv *tfprotov5.DynamicValue) map[string]tftypes.Value {
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
		t.Fatalf("reading the attributes of a %s: %s", ty, err)
	}
	return attrs
}

func stringOf(t *testing.T, v tftypes.Value) string {
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

func assertStateIsNull(t *testing.T, ty tftypes.Type, dv *tfprotov5.DynamicValue) {
	t.Helper()
	if dv == nil {
		return
	}
	val, err := dv.Unmarshal(ty)
	if err != nil {
		t.Fatalf("unmarshalling the destroyed state: %s", err)
	}
	if !val.IsNull() {
		t.Errorf("the state after the operation is %s, want null", val)
	}
}

func requireNoDiagErrors(t *testing.T, step string, diags []*tfprotov5.Diagnostic) {
	t.Helper()
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("%s: %s: %s", step, d.Summary, d.Detail)
		}
	}
}

// requireDiagErrors fails when no error-severity diagnostic is present: the
// failure paths under test (Get 5xx, malformed, mismatched identity, delete
// 5xx) must surface a diagnostic, never complete silently.
func requireDiagErrors(t *testing.T, step string, diags []*tfprotov5.Diagnostic) {
	t.Helper()
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			return
		}
	}
	t.Fatalf("%s: expected an error diagnostic, got none (%d diagnostics)", step, len(diags))
}

// --- mock handler building blocks ---

// paiAlgorithmVersionStateObject is the full GetAlgorithmVersion response the
// mock echoes back on a successful Get, with a spec the caller can override.
func paiAlgorithmVersionStateObject(spec string) map[string]interface{} {
	obj := map[string]interface{}{
		"AlgorithmId":       mockAlgorithmID,
		"AlgorithmVersion":  mockAlgorithmVersion,
		"AlgorithmName":     "echo",
		"AlgorithmProvider": "self",
		"RegionId":          "cn-beijing",
		"CreateTime":        "2026-01-01T00:00:00Z",
	}
	if spec != "" {
		obj["AlgorithmSpec"] = mustSpec(spec)
	}
	return obj
}

// createConfigValue is the config-derived plan value for a create: the three
// Required attributes set, every Computed attribute null (the Framework fills
// them during plan).
func createConfigValue(t *testing.T, spec string) *tfprotov5.DynamicValue {
	return dynamic(t, algorithmVersionFrameworkType, object(t, algorithmVersionFrameworkType, map[string]tftypes.Value{
		"algorithm_id":      tftypes.NewValue(tftypes.String, mockAlgorithmID),
		"algorithm_version": tftypes.NewValue(tftypes.String, mockAlgorithmVersion),
		"algorithm_spec":    tftypes.NewValue(tftypes.String, spec),
	}))
}

// priorStateValue is the current/prior state for Read/Delete: the composite id
// and the two identifiers set, every computed attribute null.
func priorStateValue(t *testing.T, spec string) *tfprotov5.DynamicValue {
	attrs := map[string]tftypes.Value{
		"id":                tftypes.NewValue(tftypes.String, mockAlgorithmID+":"+mockAlgorithmVersion),
		"algorithm_id":      tftypes.NewValue(tftypes.String, mockAlgorithmID),
		"algorithm_version": tftypes.NewValue(tftypes.String, mockAlgorithmVersion),
	}
	if spec != "" {
		attrs["algorithm_spec"] = tftypes.NewValue(tftypes.String, spec)
	}
	return dynamic(t, algorithmVersionFrameworkType, object(t, algorithmVersionFrameworkType, attrs))
}

// applyCreate runs the plan + apply for a create against server and returns the
// created NewState along with the apply diagnostics. PlanResourceChange never
// reaches PaiStudio (no cloud call during plan), so its diagnostics must be
// clean; the caller decides whether the apply diagnostics must be clean
// (success path) or must carry an error (the recovery-anchor path, where the
// post-create Get fails but the persisted id survives).
func applyCreate(t *testing.T, ctx context.Context, server tfprotov5.ProviderServer, configValue *tfprotov5.DynamicValue) (*tfprotov5.DynamicValue, []*tfprotov5.Diagnostic) {
	t.Helper()
	plan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         algorithmVersionTypeName,
		PriorState:       nullOf(t, algorithmVersionFrameworkType),
		ProposedNewState: configValue,
		Config:           configValue,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (create): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (create)", plan.Diagnostics)

	created, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       algorithmVersionTypeName,
		PriorState:     nullOf(t, algorithmVersionFrameworkType),
		PlannedState:   plan.PlannedState,
		Config:         configValue,
		PlannedPrivate: plan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (create): %s", err)
	}
	return created.NewState, created.Diagnostics
}

// mustSpec parses an algorithm_spec JSON string through the production
// parseAlgorithmSpec; it fails the test on invalid JSON rather than building a
// corrupt mock response.
func mustSpec(raw string) interface{} {
	spec, err := parseAlgorithmSpec(raw)
	if err != nil {
		panic(fmt.Sprintf("invalid mock algorithm_spec JSON: %s", err))
	}
	return spec
}

// specJSONEqual compares two algorithm_spec JSON strings semantically: the
// production Read canonicalizes the spec via json.Marshal, so the persisted
// string may differ in key order/whitespace from the input. Parsing both and
// deep-equal avoids a spurious order mismatch.
func specJSONEqual(a, b string) bool {
	var av, bv interface{}
	if err := json.Unmarshal([]byte(a), &av); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &bv); err != nil {
		return false
	}
	return reflect.DeepEqual(av, bv)
}

// --- TestUnits ---

// TestUnitPaiAlgorithmVersionFrameworkCreateSuccess drives the real Create
// method through a protocol-v5 server: RoaPost creates the version, the
// post-create Get returns it, and the persisted state carries the composite id
// and the echoed computed fields.
func TestUnitPaiAlgorithmVersionFrameworkCreateSuccess(t *testing.T) {
	store := newMockVersion()
	server0 := httptest.NewServer(store.handler(t))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	created, createDiags := applyCreate(t, ctx, server, createConfigValue(t, mockSpecCreate))
	requireNoDiagErrors(t, "ApplyResourceChange (create)", createDiags)
	createdAttrs := attrsOf(t, algorithmVersionFrameworkType, created)

	if got := stringOf(t, createdAttrs["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("created state id: got %q, want %q", got, mockAlgorithmID+":"+mockAlgorithmVersion)
	}
	if got := stringOf(t, createdAttrs["algorithm_id"]); got != mockAlgorithmID {
		t.Fatalf("created state algorithm_id: got %q", got)
	}
	if got := stringOf(t, createdAttrs["algorithm_version"]); got != mockAlgorithmVersion {
		t.Fatalf("created state algorithm_version: got %q", got)
	}
	if spec := stringOf(t, createdAttrs["algorithm_spec"]); spec == "" {
		t.Fatal("created state algorithm_spec must be populated from the Get response")
	}
	// The post-create Get echoed computed fields the plan left null.
	if got := stringOf(t, createdAttrs["region_id"]); got != "cn-beijing" {
		t.Fatalf("created state region_id: got %q, want cn-beijing", got)
	}
	// Create issues RoaPost then a read-back RoaGet, so inspect the method
	// history (not the trailing lastMethod, which is the Get) to prove the
	// write reached RoaPost.
	if !store.hasMethod(http.MethodPost) {
		t.Fatalf("create did not reach RoaPost: method history %v", store.methods)
	}
}

// TestUnitPaiAlgorithmVersionFrameworkCreateGetFailurePreservesID proves the
// recovery anchor: Create persists the composite id (and the original plan)
// BEFORE the post-create Get, so a Get that fails (5xx / not-yet-visible 404
// NoSuchObject / malformed 200) cannot drop the already-created resource. Each
// sub-case surfaces an error diagnostic but keeps the persisted id.
func TestUnitPaiAlgorithmVersionFrameworkCreateGetFailurePreservesID(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name string
		get  func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "get_500",
			get: func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(t, w, http.StatusInternalServerError, map[string]interface{}{
					"Code": "InternalError", "Message": "boom", "RequestId": "r-500",
				})
			},
		},
		{
			name: "get_not_yet_visible_404",
			get: func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(t, w, http.StatusNotFound, map[string]interface{}{
					"Code": "NoSuchObject", "Message": "not visible yet", "RequestId": "r-nf",
				})
			},
		},
		{
			name: "get_malformed_200",
			get: func(w http.ResponseWriter, _ *http.Request) {
				// 200 but missing AlgorithmVersion: would compose "algo-mock:".
				writeJSON(t, w, http.StatusOK, map[string]interface{}{
					"AlgorithmId": mockAlgorithmID,
				})
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Create always RoaPosts first; return a valid 200 so the write
				// succeeds and the recovery-anchor Set runs before the Get.
				if r.Method == http.MethodPost {
					writeJSON(t, w, http.StatusOK, paiAlgorithmVersionStateObject(mockSpecCreate))
					return
				}
				if r.Method == http.MethodGet {
					tc.get(w, r)
					return
				}
				writeJSON(t, w, http.StatusMethodNotAllowed, map[string]interface{}{"Code": "MethodNotAllowed"})
			})
			server0 := httptest.NewServer(handler)
			defer server0.Close()
			client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
			server := newPaiFrameworkServer(t, client)

			created, createDiags := applyCreate(t, ctx, server, createConfigValue(t, mockSpecCreate))
			requireDiagErrors(t, "ApplyResourceChange (create with failed Get)", createDiags)
			createdAttrs := attrsOf(t, algorithmVersionFrameworkType, created)

			// The persisted id survives the failed Get: the resource is not lost.
			if got := stringOf(t, createdAttrs["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
				t.Fatalf("%s: persisted id after failed post-create Get: got %q, want %q (recovery anchor lost)",
					tc.name, got, mockAlgorithmID+":"+mockAlgorithmVersion)
			}
			// The original plan survives too: the configured spec is intact.
			if got := stringOf(t, createdAttrs["algorithm_spec"]); got != mockSpecCreate {
				t.Fatalf("%s: persisted algorithm_spec after failed Get: got %q, want the original plan",
					tc.name, got)
			}
		})
	}
}

// TestUnitPaiAlgorithmVersionFrameworkReadNotFoundRemoves proves Read removes
// the resource from state on a true NotFound (typed NoSuchObject on the
// production RoaGet path): no error diagnostic, and the new state is null.
func TestUnitPaiAlgorithmVersionFrameworkReadNotFoundRemoves(t *testing.T) {
	server0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusNotFound, map[string]interface{}{
			"Code": "NoSuchObject", "Message": "algorithm version not found", "RequestId": "r-nf",
		})
	}))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	refreshed, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:     algorithmVersionTypeName,
		CurrentState: priorStateValue(t, mockSpecCreate),
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	requireNoDiagErrors(t, "ReadResource (not found)", refreshed.Diagnostics)
	assertStateIsNull(t, algorithmVersionFrameworkType, refreshed.NewState)
}

// TestUnitPaiAlgorithmVersionFrameworkReadTransientErrorPreservesState proves a
// 5xx on Read is not treated as not-found: the prior state is preserved and an
// error diagnostic is surfaced, so a transient server error cannot silently
// drop the resource.
func TestUnitPaiAlgorithmVersionFrameworkReadTransientErrorPreservesState(t *testing.T) {
	server0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusInternalServerError, map[string]interface{}{
			"Code": "InternalError", "Message": "boom", "RequestId": "r-500",
		})
	}))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	refreshed, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:     algorithmVersionTypeName,
		CurrentState: priorStateValue(t, mockSpecCreate),
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	requireDiagErrors(t, "ReadResource (5xx)", refreshed.Diagnostics)
	if refreshed.NewState == nil {
		t.Fatal("ReadResource (5xx) dropped state; the production Read relies on the Framework preserving prior state on a transient error")
	}
	if got := stringOf(t, attrsOf(t, algorithmVersionFrameworkType, refreshed.NewState)["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("ReadResource (5xx) must preserve the prior id: got %q", got)
	}
}

// TestUnitPaiAlgorithmVersionFrameworkReadMalformedPreservesState proves a 200
// response missing a required identifier (the path that would compose a ":"
// id) is rejected by the response validation: prior state is kept and an error
// diagnostic is surfaced, the malformed object is not flattened onto state.
func TestUnitPaiAlgorithmVersionFrameworkReadMalformedPreservesState(t *testing.T) {
	server0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// 200 but AlgorithmVersion is missing: would compose "algo-mock:".
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"AlgorithmId": mockAlgorithmID,
		})
	}))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	refreshed, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:     algorithmVersionTypeName,
		CurrentState: priorStateValue(t, mockSpecCreate),
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	requireDiagErrors(t, "ReadResource (malformed)", refreshed.Diagnostics)
	if refreshed.NewState == nil {
		t.Fatal("ReadResource (malformed) dropped state; the production Read relies on the Framework preserving prior state on a malformed response")
	}
	if got := stringOf(t, attrsOf(t, algorithmVersionFrameworkType, refreshed.NewState)["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("ReadResource (malformed) must preserve the prior id: got %q", got)
	}
}

// TestUnitPaiAlgorithmVersionFrameworkReadMismatchedIdentityPreservesState
// proves a GetAlgorithmVersion response whose AlgorithmId does not match the
// requested resource is rejected: the prior identity is kept and an error
// diagnostic is surfaced, so a mismatched response cannot overwrite the
// resource's identity.
func TestUnitPaiAlgorithmVersionFrameworkReadMismatchedIdentityPreservesState(t *testing.T) {
	server0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"AlgorithmId":      "other-algorithm",
			"AlgorithmVersion": mockAlgorithmVersion,
		})
	}))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	refreshed, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:     algorithmVersionTypeName,
		CurrentState: priorStateValue(t, mockSpecCreate),
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	requireDiagErrors(t, "ReadResource (mismatched identity)", refreshed.Diagnostics)
	if refreshed.NewState == nil {
		t.Fatal("ReadResource (mismatched identity) dropped state; the production Read relies on the Framework preserving prior state")
	}
	attrs := attrsOf(t, algorithmVersionFrameworkType, refreshed.NewState)
	if got := stringOf(t, attrs["algorithm_id"]); got != mockAlgorithmID {
		t.Fatalf("ReadResource (mismatched identity) must keep the prior algorithm_id: got %q, want %q", got, mockAlgorithmID)
	}
	if got := stringOf(t, attrs["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("ReadResource (mismatched identity) must keep the prior composite id: got %q", got)
	}
}

// TestUnitPaiAlgorithmVersionFrameworkUpdateIssuesRoaPutAndRefreshes proves the
// Update method issues RoaPut with the planned spec and refreshes state from
// the post-update Get: the mock store receives a PUT, and the new state carries
// the updated spec.
func TestUnitPaiAlgorithmVersionFrameworkUpdateIssuesRoaPutAndRefreshes(t *testing.T) {
	store := newMockVersion()
	server0 := httptest.NewServer(store.handler(t))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	// Seed the version via a real Create so the store and the persisted state are
	// consistent before the update.
	created, createDiags := applyCreate(t, ctx, server, createConfigValue(t, mockSpecCreate))
	requireNoDiagErrors(t, "ApplyResourceChange (create seed)", createDiags)
	createdAttrs := attrsOf(t, algorithmVersionFrameworkType, created)

	// The update changes only algorithm_spec (RequiresReplace fields untouched).
	updateConfig := createConfigValue(t, mockSpecUpdate)
	updateProposed := dynamic(t, algorithmVersionFrameworkType, object(t, algorithmVersionFrameworkType, map[string]tftypes.Value{
		"id":                createdAttrs["id"],
		"algorithm_id":      createdAttrs["algorithm_id"],
		"algorithm_version": createdAttrs["algorithm_version"],
		"algorithm_spec":    tftypes.NewValue(tftypes.String, mockSpecUpdate),
	}))
	plan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         algorithmVersionTypeName,
		PriorState:       created,
		ProposedNewState: updateProposed,
		Config:           updateConfig,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (update): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (update)", plan.Diagnostics)
	if len(plan.RequiresReplace) != 0 {
		t.Fatalf("planning a spec change requires replacing %v; algorithm_spec must not be ForceNew", plan.RequiresReplace)
	}

	updated, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       algorithmVersionTypeName,
		PriorState:     created,
		PlannedState:   plan.PlannedState,
		Config:         updateConfig,
		PlannedPrivate: plan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (update): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (update)", updated.Diagnostics)

	// Update issues RoaPut then a read-back RoaGet, so inspect the method
	// history (not the trailing lastMethod, which is the Get) to prove the
	// write reached RoaPut.
	if !store.hasMethod(http.MethodPut) {
		t.Fatalf("update did not reach RoaPut: method history %v", store.methods)
	}
	updatedAttrs := attrsOf(t, algorithmVersionFrameworkType, updated.NewState)
	if !specJSONEqual(stringOf(t, updatedAttrs["algorithm_spec"]), mockSpecUpdate) {
		t.Fatalf("updated state algorithm_spec: got %q, want %q", stringOf(t, updatedAttrs["algorithm_spec"]), mockSpecUpdate)
	}
	if got := stringOf(t, updatedAttrs["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("updated state id: got %q", got)
	}
}

// TestUnitPaiAlgorithmVersionFrameworkDeleteErrorPreservesState proves a
// transient delete error keeps the resource in state with an error diagnostic
// (the resource is not silently destroyed).
func TestUnitPaiAlgorithmVersionFrameworkDeleteErrorPreservesState(t *testing.T) {
	server0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusInternalServerError, map[string]interface{}{
			"Code": "InternalError", "Message": "boom", "RequestId": "r-500",
		})
	}))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	prior := priorStateValue(t, mockSpecCreate)
	destroyed, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:     algorithmVersionTypeName,
		PriorState:   prior,
		PlannedState: nullOf(t, algorithmVersionFrameworkType),
		Config:       nullOf(t, algorithmVersionFrameworkType),
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (delete): %s", err)
	}
	requireDiagErrors(t, "ApplyResourceChange (delete 5xx)", destroyed.Diagnostics)
	if destroyed.NewState == nil {
		t.Fatal("ApplyResourceChange (delete 5xx) dropped state; the production Delete relies on the Framework preserving prior state on a transient error")
	}
	if got := stringOf(t, attrsOf(t, algorithmVersionFrameworkType, destroyed.NewState)["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("delete (5xx) must preserve the prior id: got %q", got)
	}
}

// TestUnitPaiAlgorithmVersionFrameworkDeleteNotFoundIsIdempotent proves a
// delete of an already-absent version (typed NoSuchObject) completes with no
// error and destroys the resource: the production Delete treats NoSuchObject as
// a success, so the destroy is idempotent.
func TestUnitPaiAlgorithmVersionFrameworkDeleteNotFoundIsIdempotent(t *testing.T) {
	server0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusNotFound, map[string]interface{}{
			"Code": "NoSuchObject", "Message": "algorithm version not found", "RequestId": "r-del",
		})
	}))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	prior := priorStateValue(t, mockSpecCreate)
	destroyed, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:     algorithmVersionTypeName,
		PriorState:   prior,
		PlannedState: nullOf(t, algorithmVersionFrameworkType),
		Config:       nullOf(t, algorithmVersionFrameworkType),
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (delete): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (delete not-found)", destroyed.Diagnostics)
	assertStateIsNull(t, algorithmVersionFrameworkType, destroyed.NewState)
}

// TestUnitPaiAlgorithmVersionFrameworkImportReadPlanNoOp proves the ID-only
// import path: ImportState parses the composite id and seeds only id + the two
// identifiers; a following Read against the mock fills the computed fields, and
// a no-change plan over the refreshed state produces no RequiresReplace (the
// imported state is consistent with config, no diff churn).
func TestUnitPaiAlgorithmVersionFrameworkImportReadPlanNoOp(t *testing.T) {
	store := newMockVersion()
	server0 := httptest.NewServer(store.handler(t))
	defer server0.Close()
	client := newOfflinePaiStudioClient(t, paiStudioHost(server0.URL))
	server := newPaiFrameworkServer(t, client)
	ctx := context.Background()

	// Seed the version so the Read after import finds it.
	if _, err := client.RoaPost(PaiStudioProduct, PaiStudioVersion,
		algorithmVersionPath(mockAlgorithmID, mockAlgorithmVersion), nil, nil,
		map[string]interface{}{"AlgorithmSpec": mustSpec(mockSpecCreate)}, true); err != nil {
		t.Fatalf("seeding version: %s", err)
	}

	importResp, err := server.ImportResourceState(ctx, &tfprotov5.ImportResourceStateRequest{
		TypeName: algorithmVersionTypeName,
		ID:       mockAlgorithmID + ":" + mockAlgorithmVersion,
	})
	if err != nil {
		t.Fatalf("ImportResourceState: %s", err)
	}
	requireNoDiagErrors(t, "ImportResourceState", importResp.Diagnostics)
	if len(importResp.ImportedResources) != 1 {
		t.Fatalf("import produced %d resources, want 1", len(importResp.ImportedResources))
	}
	importedState := importResp.ImportedResources[0].State
	importedAttrs := attrsOf(t, algorithmVersionFrameworkType, importedState)
	if got := stringOf(t, importedAttrs["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("imported state id: got %q", got)
	}
	if got := stringOf(t, importedAttrs["algorithm_id"]); got != mockAlgorithmID {
		t.Fatalf("imported state algorithm_id: got %q", got)
	}
	if got := stringOf(t, importedAttrs["algorithm_version"]); got != mockAlgorithmVersion {
		t.Fatalf("imported state algorithm_version: got %q", got)
	}
	// Computed fields are null after a bare import: they are filled by Read.
	if !importedAttrs["region_id"].IsNull() || !importedAttrs["create_time"].IsNull() {
		t.Fatal("import must seed only id + the two identifiers; computed fields must be null until Read")
	}

	// Read fills the computed fields from the Get response.
	refreshed, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:     algorithmVersionTypeName,
		CurrentState: importedState,
	})
	if err != nil {
		t.Fatalf("ReadResource after import: %s", err)
	}
	requireNoDiagErrors(t, "ReadResource after import", refreshed.Diagnostics)
	readAttrs := attrsOf(t, algorithmVersionFrameworkType, refreshed.NewState)
	if got := stringOf(t, readAttrs["region_id"]); got != "cn-beijing" {
		t.Fatalf("read after import must fill region_id: got %q", got)
	}
	if got := stringOf(t, readAttrs["create_time"]); got != "2026-01-01T00:00:00Z" {
		t.Fatalf("read after import must fill create_time: got %q", got)
	}

	// A no-change plan over the refreshed state is a no-op: no RequiresReplace
	// and the planned identity/spec match the refreshed state.
	noopPlan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         algorithmVersionTypeName,
		PriorState:       refreshed.NewState,
		ProposedNewState: refreshed.NewState,
		Config:           refreshed.NewState,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (no-op): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (no-op)", noopPlan.Diagnostics)
	if len(noopPlan.RequiresReplace) != 0 {
		t.Fatalf("a no-change plan requires replacing %v; the imported+read state must be stable", noopPlan.RequiresReplace)
	}
	planAttrs := attrsOf(t, algorithmVersionFrameworkType, noopPlan.PlannedState)
	if got := stringOf(t, planAttrs["id"]); got != mockAlgorithmID+":"+mockAlgorithmVersion {
		t.Fatalf("no-op plan id: got %q", got)
	}
	if got := stringOf(t, planAttrs["algorithm_id"]); got != mockAlgorithmID {
		t.Fatalf("no-op plan algorithm_id: got %q", got)
	}
	if got := stringOf(t, planAttrs["algorithm_version"]); got != mockAlgorithmVersion {
		t.Fatalf("no-op plan algorithm_version: got %q", got)
	}
	if got := stringOf(t, planAttrs["algorithm_spec"]); !specJSONEqual(got, mockSpecCreate) {
		t.Fatalf("no-op plan algorithm_spec: got %q, want %q", got, mockSpecCreate)
	}
}
