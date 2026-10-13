package fwadapt

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type stubResource struct {
	createCalled bool
	readCalled   bool
	updateCalled bool
	deleteCalled bool
	createErr    string
	readErr      string
	updateErr    string
	deleteErr    string

	importStateCalled    bool
	identitySchemaCalled bool
}

func (s *stubResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "test_stub"
}

func (s *stubResource) Schema(_ context.Context, _ resource.SchemaRequest, _ *resource.SchemaResponse) {
}

func (s *stubResource) Create(_ context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	s.createCalled = true
	if s.createErr != "" {
		resp.Diagnostics.AddError(s.createErr, "")
	}
}

func (s *stubResource) Read(_ context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
	s.readCalled = true
	if s.readErr != "" {
		resp.Diagnostics.AddError(s.readErr, "")
	}
}

func (s *stubResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	s.updateCalled = true
	if s.updateErr != "" {
		resp.Diagnostics.AddError(s.updateErr, "")
	}
}

func (s *stubResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	s.deleteCalled = true
	if s.deleteErr != "" {
		resp.Diagnostics.AddError(s.deleteErr, "")
	}
}

// diagsResource reports diagnostics the AddError helpers cannot build, such as one
// carrying an attribute path.
type diagsResource struct {
	stubResource
	diags diag.Diagnostics
}

func (s *diagsResource) Create(_ context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	s.createCalled = true
	resp.Diagnostics.Append(s.diags...)
}

type stubDataSource struct {
	readCalled bool
	readErr    string
}

func (s *stubDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "test_stub"
}

func (s *stubDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, _ *datasource.SchemaResponse) {
}

func (s *stubDataSource) Read(_ context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	s.readCalled = true
	if s.readErr != "" {
		resp.Diagnostics.AddError(s.readErr, "")
	}
}

type recordInterceptor struct {
	name       string
	ops        []string
	beforeErr  error
	afterErr   error
	rewriteErr error
	wrapMsg    string
	swallow    bool
	metaSeen   interface{}
}

func (r *recordInterceptor) Before(ctx context.Context, call intercept.Call) error {
	r.ops = append(r.ops, "before:"+r.name)
	r.metaSeen = call.Meta
	return r.beforeErr
}

func (r *recordInterceptor) After(ctx context.Context, call intercept.Call, err error) error {
	r.ops = append(r.ops, "after:"+r.name)
	if r.swallow {
		return nil
	}
	if r.rewriteErr != nil {
		return r.rewriteErr
	}
	if r.wrapMsg != "" && err != nil {
		return fmt.Errorf("%s: %w", r.wrapMsg, err)
	}
	return err
}

type importStateStub struct{ stubResource }

func (s *importStateStub) ImportState(_ context.Context, _ resource.ImportStateRequest, _ *resource.ImportStateResponse) {
	s.importStateCalled = true
}

type identityStub struct{ stubResource }

func (s *identityStub) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, _ *resource.IdentitySchemaResponse) {
	s.identitySchemaCalled = true
}

type importStateAndIdentityStub struct{ stubResource }

func (s *importStateAndIdentityStub) ImportState(_ context.Context, _ resource.ImportStateRequest, _ *resource.ImportStateResponse) {
	s.importStateCalled = true
}

func (s *importStateAndIdentityStub) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, _ *resource.IdentitySchemaResponse) {
	s.identitySchemaCalled = true
}

type upgradeStateStub struct{ stubResource }

func (s *upgradeStateStub) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{0: {}}
}

type configureStub struct {
	stubResource
	configuredWith interface{}
}

func (s *configureStub) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	s.configuredWith = req.ProviderData
}

func TestWrapResourceNilChainIsIdentity(t *testing.T) {
	r := &stubResource{}
	got := WrapResource("test", r, nil)
	if got != resource.Resource(r) {
		t.Fatal("nil chain must return the same resource")
	}
	got = WrapResource("test", r, []intercept.Interceptor{})
	if got != resource.Resource(r) {
		t.Fatal("empty chain must return the same resource")
	}
}

func TestWrapResourceCreate(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := &stubResource{}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.CreateResponse
	wr.Create(context.Background(), resource.CreateRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diags: %v", resp.Diagnostics)
	}
	if !r.createCalled {
		t.Fatal("inner Create not called")
	}
	want := []string{"before:rec", "after:rec"}
	if fmt.Sprint(rec.ops) != fmt.Sprint(want) {
		t.Fatalf("ops = %v, want %v", rec.ops, want)
	}
}

func TestWrapResourceRead(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := &stubResource{}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.ReadResponse
	wr.Read(context.Background(), resource.ReadRequest{}, &resp)

	if !r.readCalled {
		t.Fatal("inner Read not called")
	}
}

func TestWrapResourceUpdate(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := &stubResource{}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.UpdateResponse
	wr.Update(context.Background(), resource.UpdateRequest{}, &resp)

	if !r.updateCalled {
		t.Fatal("inner Update not called")
	}
}

func TestWrapResourceDelete(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := &stubResource{}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.DeleteResponse
	wr.Delete(context.Background(), resource.DeleteRequest{}, &resp)

	if !r.deleteCalled {
		t.Fatal("inner Delete not called")
	}
}

func TestWrapResourceBeforeAbort(t *testing.T) {
	rec := &recordInterceptor{name: "rec", beforeErr: errors.New("abort")}
	r := &stubResource{}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.CreateResponse
	wr.Create(context.Background(), resource.CreateRequest{}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error from Before abort")
	}
	if r.createCalled {
		t.Fatal("inner must be skipped after Before abort")
	}
}

func TestWrapResourceAfterSeesInnerError(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := &stubResource{createErr: "inner error"}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.CreateResponse
	wr.Create(context.Background(), resource.CreateRequest{}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics from inner error")
	}
	want := []string{"before:rec", "after:rec"}
	if fmt.Sprint(rec.ops) != fmt.Sprint(want) {
		t.Fatalf("ops = %v, want %v", rec.ops, want)
	}
}

func TestWrapResourceAfterRewritesError(t *testing.T) {
	rec := &recordInterceptor{name: "rec", rewriteErr: errors.New("rewritten")}
	r := &stubResource{createErr: "inner error"}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.CreateResponse
	wr.Create(context.Background(), resource.CreateRequest{}, &resp)

	if got, want := len(resp.Diagnostics), 2; got != want {
		t.Fatalf("len(diagnostics) = %d, want %d: %v", got, want, resp.Diagnostics)
	}
	if got := resp.Diagnostics[0].Summary(); got != "inner error" {
		t.Fatalf("the inner error was replaced, summary = %q", got)
	}
	if got := resp.Diagnostics[1].Summary(); got != "rewritten" {
		t.Fatalf("summary = %q, want %q", got, "rewritten")
	}
}

func TestWrapResourceAfterCannotSwallowError(t *testing.T) {
	rec := &recordInterceptor{name: "rec", swallow: true}
	r := &stubResource{readErr: "resource not found"}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.ReadResponse
	wr.Read(context.Background(), resource.ReadRequest{}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("After returned nil, but a failure must not be reported as success")
	}
	if got, want := len(resp.Diagnostics), 1; got != want {
		t.Fatalf("len(diagnostics) = %d, want %d: %v", got, want, resp.Diagnostics)
	}
	if got := resp.Diagnostics[0].Summary(); got != "resource not found" {
		t.Fatalf("summary = %q, want the inner one restored verbatim", got)
	}
	if !r.readCalled {
		t.Fatal("inner Read not called")
	}
}

func TestWrapResourceWarningsSurviveRewrite(t *testing.T) {
	rec := &recordInterceptor{name: "rec", rewriteErr: errors.New("rewritten")}
	r := &stubResource{readErr: "resource not found"}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.ReadResponse
	resp.Diagnostics.AddWarning("deprecated field", "use the other one")
	wr.Read(context.Background(), resource.ReadRequest{}, &resp)

	if got := resp.Diagnostics[0].Summary(); got != "deprecated field" {
		t.Fatalf("the warning must come first and unchanged, diagnostics = %v", resp.Diagnostics)
	}
}

// The bridge hands the chain a plain error, so a hook that only wraps it must not
// cost the attribute path and detail that error came from.
func TestWrapResourceWrapKeepsEveryInnerErrorIntact(t *testing.T) {
	rec := &recordInterceptor{name: "rec", wrapMsg: "retry gave up"}
	r := &diagsResource{diags: diag.Diagnostics{
		diag.NewAttributeErrorDiagnostic(path.Root("cidr_block"), "invalid CIDR", "must be a /16 or narrower"),
		diag.NewErrorDiagnostic("second failure", "with a detail"),
	}}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.CreateResponse
	wr.Create(context.Background(), resource.CreateRequest{}, &resp)

	if got, want := len(resp.Diagnostics), 3; got != want {
		t.Fatalf("len(diagnostics) = %d, want %d: %v", got, want, resp.Diagnostics)
	}
	first, ok := resp.Diagnostics[0].(diag.DiagnosticWithPath)
	if !ok {
		t.Fatalf("the attribute path was lost: %#v", resp.Diagnostics[0])
	}
	if got := first.Path().String(); got != "cidr_block" {
		t.Fatalf("path = %q, want %q", got, "cidr_block")
	}
	if got := resp.Diagnostics[1].Detail(); got != "with a detail" {
		t.Fatalf("detail = %q, want it preserved", got)
	}
	// The hook's own wording arrives on top, carrying both messages it saw.
	want := "retry gave up: invalid CIDR: must be a /16 or narrower; second failure: with a detail"
	if got := resp.Diagnostics[2].Summary(); got != want {
		t.Fatalf("summary = %q, want %q", got, want)
	}
}

func TestWrapResourceForwardsProviderDataAsMeta(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := &configureStub{}
	got := WrapResource("test", r, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	meta := struct{ client string }{client: "stub"}
	wr.Configure(context.Background(), resource.ConfigureRequest{ProviderData: meta}, &resource.ConfigureResponse{})
	if r.configuredWith != interface{}(meta) {
		t.Fatal("inner Configure did not receive the provider data")
	}

	var resp resource.CreateResponse
	wr.Create(context.Background(), resource.CreateRequest{}, &resp)
	if rec.metaSeen != interface{}(meta) {
		t.Fatalf("interceptor meta = %v, want the provider data %v", rec.metaSeen, meta)
	}
}

func TestWrapResourceMetaIsNilBeforeConfigure(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	got := WrapResource("test", &stubResource{}, []intercept.Interceptor{rec})
	wr := got.(*wrappedResource)

	var resp resource.ReadResponse
	wr.Read(context.Background(), resource.ReadRequest{}, &resp)
	if rec.metaSeen != nil {
		t.Fatalf("meta = %v, want nil before Configure", rec.metaSeen)
	}
}

func TestWrapResourceNilSafeInterfacesAlwaysPresent(t *testing.T) {
	got := WrapResource("test", &stubResource{}, []intercept.Interceptor{&recordInterceptor{name: "rec"}})

	if _, ok := got.(resource.ResourceWithConfigure); !ok {
		t.Error("wrapper must implement ResourceWithConfigure")
	}
	if _, ok := got.(resource.ResourceWithConfigValidators); !ok {
		t.Error("wrapper must implement ResourceWithConfigValidators")
	}
	if _, ok := got.(resource.ResourceWithModifyPlan); !ok {
		t.Error("wrapper must implement ResourceWithModifyPlan")
	}
	if _, ok := got.(resource.ResourceWithValidateConfig); !ok {
		t.Error("wrapper must implement ResourceWithValidateConfig")
	}
}

func TestWrapResourceDoesNotInventCapabilities(t *testing.T) {
	got := WrapResource("test", &stubResource{}, []intercept.Interceptor{&recordInterceptor{name: "rec"}})

	if _, ok := got.(resource.ResourceWithImportState); ok {
		t.Error("wrapper must not add ResourceWithImportState")
	}
	if _, ok := got.(resource.ResourceWithIdentity); ok {
		t.Error("wrapper must not add ResourceWithIdentity")
	}
}

func TestWrapResourceStateMigrationCapabilities(t *testing.T) {
	chain := []intercept.Interceptor{&recordInterceptor{name: "rec"}}

	plain := WrapResource("test", &stubResource{}, chain).(*wrappedResource)
	if got := plain.UpgradeState(context.Background()); got != nil {
		t.Errorf("UpgradeState = %v, want nil for a resource without it", got)
	}
	if got := plain.MoveState(context.Background()); got != nil {
		t.Errorf("MoveState = %v, want nil for a resource without it", got)
	}
	if got := plain.UpgradeIdentity(context.Background()); got != nil {
		t.Errorf("UpgradeIdentity = %v, want nil for a resource without it", got)
	}

	wrapped := WrapResource("test", &upgradeStateStub{}, chain)
	withUpgrade, ok := wrapped.(resource.ResourceWithUpgradeState)
	if !ok {
		t.Fatal("wrapper must implement ResourceWithUpgradeState")
	}
	if got := withUpgrade.UpgradeState(context.Background()); len(got) != 1 {
		t.Errorf("UpgradeState = %v, want the inner resource's single upgrader", got)
	}

	rec := &recordInterceptor{name: "rec"}
	inner := &upgradeStateStub{}
	got := WrapResource("test", inner, []intercept.Interceptor{rec})
	var resp resource.ReadResponse
	got.Read(context.Background(), resource.ReadRequest{}, &resp)
	if !inner.readCalled {
		t.Error("inner Read not called")
	}
	want := []string{"before:rec", "after:rec"}
	if fmt.Sprint(rec.ops) != fmt.Sprint(want) {
		t.Errorf("ops = %v, want %v", rec.ops, want)
	}
}

func TestWrapResourceCapabilitySplit(t *testing.T) {
	cases := []struct {
		name        string
		inner       resource.Resource
		importState bool
		identity    bool
	}{
		{name: "plain", inner: &stubResource{}},
		{name: "import state", inner: &importStateStub{}, importState: true},
		{name: "identity", inner: &identityStub{}, identity: true},
		{name: "both", inner: &importStateAndIdentityStub{}, importState: true, identity: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := &recordInterceptor{name: "rec"}
			got := WrapResource("test", tc.inner, []intercept.Interceptor{rec})

			withImportState, hasImportState := got.(resource.ResourceWithImportState)
			if hasImportState != tc.importState {
				t.Fatalf("ResourceWithImportState = %v, want %v", hasImportState, tc.importState)
			}
			withIdentity, hasIdentity := got.(resource.ResourceWithIdentity)
			if hasIdentity != tc.identity {
				t.Fatalf("ResourceWithIdentity = %v, want %v", hasIdentity, tc.identity)
			}

			var resp resource.ReadResponse
			got.Read(context.Background(), resource.ReadRequest{}, &resp)
			want := []string{"before:rec", "after:rec"}
			if fmt.Sprint(rec.ops) != fmt.Sprint(want) {
				t.Fatalf("ops = %v, want %v", rec.ops, want)
			}

			if tc.importState {
				withImportState.ImportState(context.Background(), resource.ImportStateRequest{}, &resource.ImportStateResponse{})
			}
			if tc.identity {
				withIdentity.IdentitySchema(context.Background(), resource.IdentitySchemaRequest{}, &resource.IdentitySchemaResponse{})
			}
			inner := innerStub(tc.inner)
			if inner.importStateCalled != tc.importState {
				t.Errorf("inner ImportState called = %v, want %v", inner.importStateCalled, tc.importState)
			}
			if inner.identitySchemaCalled != tc.identity {
				t.Errorf("inner IdentitySchema called = %v, want %v", inner.identitySchemaCalled, tc.identity)
			}
		})
	}
}

func innerStub(r resource.Resource) *stubResource {
	switch v := r.(type) {
	case *stubResource:
		return v
	case *importStateStub:
		return &v.stubResource
	case *identityStub:
		return &v.stubResource
	case *importStateAndIdentityStub:
		return &v.stubResource
	}
	panic("unknown stub type")
}

func TestWrapDataSourceRead(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	ds := &stubDataSource{}
	got := WrapDataSource("test", ds, []intercept.Interceptor{rec})
	wd := got.(*wrappedDataSource)

	var resp datasource.ReadResponse
	wd.Read(context.Background(), datasource.ReadRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diags: %v", resp.Diagnostics)
	}
	if !ds.readCalled {
		t.Fatal("inner Read not called")
	}
	want := []string{"before:rec", "after:rec"}
	if fmt.Sprint(rec.ops) != fmt.Sprint(want) {
		t.Fatalf("ops = %v, want %v", rec.ops, want)
	}
}

func TestWrapDataSourceNilChainIsIdentity(t *testing.T) {
	ds := &stubDataSource{}
	got := WrapDataSource("test", ds, nil)
	if got != datasource.DataSource(ds) {
		t.Fatal("nil chain must return the same data source")
	}
}

func TestWrapDataSourceBeforeAbort(t *testing.T) {
	rec := &recordInterceptor{name: "rec", beforeErr: errors.New("abort")}
	ds := &stubDataSource{}
	got := WrapDataSource("test", ds, []intercept.Interceptor{rec})
	wd := got.(*wrappedDataSource)

	var resp datasource.ReadResponse
	wd.Read(context.Background(), datasource.ReadRequest{}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error from Before abort")
	}
	if ds.readCalled {
		t.Fatal("inner must be skipped after Before abort")
	}
}

func TestDiagToErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("summary", "detail")
	err := diagToErr(diags)
	if err == nil || err.Error() != "summary: detail" {
		t.Fatalf("expected 'summary: detail', got %v", err)
	}

	diags = diag.Diagnostics{}
	diags.AddError("summary only", "")
	err = diagToErr(diags)
	if err == nil || err.Error() != "summary only" {
		t.Fatalf("expected 'summary only', got %v", err)
	}

	err = diagToErr(diag.Diagnostics{})
	if err != nil {
		t.Fatalf("expected nil for empty diagnostics, got %v", err)
	}
}
