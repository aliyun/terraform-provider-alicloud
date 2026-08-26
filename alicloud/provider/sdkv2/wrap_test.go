package sdkv2

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
)

type recordInterceptor struct {
	name       string
	ops        []string
	beforeErr  error
	afterErr   error
	rewriteErr error
	wrapMsg    string
	swallow    bool
}

func (r *recordInterceptor) Before(ctx context.Context, call intercept.Call) error {
	r.ops = append(r.ops, "before:"+r.name)
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

func TestWrapResourceNilChainIsIdentity(t *testing.T) {
	r := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error { return nil },
	}
	got := WrapResource("alicloud_test", r, nil)
	if got != r {
		t.Fatalf("expected same pointer, got %p != %p", got, r)
	}
	got = WrapResource("alicloud_test", r, []intercept.Interceptor{})
	if got != r {
		t.Fatalf("expected same pointer for empty chain, got %p != %p", got, r)
	}
}

func TestWrapResourceWrapsOnlyUsedFields(t *testing.T) {
	called := false
	rec := &recordInterceptor{name: "rec"}
	r := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error { called = true; return nil },
		Schema: map[string]*schema.Schema{"id": {Type: schema.TypeString}},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	if got == r {
		t.Fatal("expected a shallow copy when chain is non-empty")
	}
	if got.Create == nil || got.Read != nil || got.Update != nil || got.Delete != nil {
		t.Fatal("only the used fields may be wrapped")
	}
	if got.Schema["id"] != r.Schema["id"] {
		t.Fatal("schema map must be shared with the original")
	}
	if err := got.Create(nil, nil); err != nil || !called {
		t.Fatalf("inner create not invoked: called=%v err=%v", called, err)
	}
}

func TestWrapResourcePreservesNonCRUDBehaviour(t *testing.T) {
	upgraderCalled := false
	customizeDiffCalled := false
	importerCalled := false

	r := &schema.Resource{
		Read:          func(d *schema.ResourceData, meta interface{}) error { return nil },
		SchemaVersion: 3,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 2,
			Type:    cty.EmptyObject,
			Upgrade: func(ctx context.Context, st map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
				upgraderCalled = true
				return st, nil
			},
		}},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			customizeDiffCalled = true
			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				importerCalled = true
				return nil, nil
			},
		},
		DeprecationMessage: "use the other one",
		Timeouts:           &schema.ResourceTimeout{},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{&recordInterceptor{name: "rec"}})

	if got.SchemaVersion != r.SchemaVersion {
		t.Errorf("SchemaVersion = %d, want %d", got.SchemaVersion, r.SchemaVersion)
	}
	if got.DeprecationMessage != r.DeprecationMessage {
		t.Errorf("DeprecationMessage = %q, want %q", got.DeprecationMessage, r.DeprecationMessage)
	}
	if got.Timeouts != r.Timeouts {
		t.Error("Timeouts must be carried over")
	}
	if got.Importer != r.Importer {
		t.Error("Importer must be carried over")
	}
	if len(got.StateUpgraders) != 1 || got.StateUpgraders[0].Version != 2 {
		t.Fatalf("StateUpgraders = %v, want the declared one", got.StateUpgraders)
	}

	if _, err := got.StateUpgraders[0].Upgrade(context.Background(), nil, nil); err != nil || !upgraderCalled {
		t.Errorf("StateUpgraders[0].Upgrade not forwarded: called=%v err=%v", upgraderCalled, err)
	}
	if err := got.CustomizeDiff(context.Background(), nil, nil); err != nil || !customizeDiffCalled {
		t.Errorf("CustomizeDiff not forwarded: called=%v err=%v", customizeDiffCalled, err)
	}
	if _, err := got.Importer.StateContext(context.Background(), nil, nil); err != nil || !importerCalled {
		t.Errorf("Importer.StateContext not forwarded: called=%v err=%v", importerCalled, err)
	}
}

func TestWrapDataSourceWrapsReadOnly(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	ds := &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error { return nil },
	}
	got := WrapDataSource("alicloud_test", ds, []intercept.Interceptor{rec})
	if got.Read == nil {
		t.Fatal("Read must be wrapped")
	}
	if err := got.Read(nil, nil); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []string{"before:rec", "after:rec"}
	if fmt.Sprint(rec.ops) != fmt.Sprint(want) {
		t.Fatalf("ops = %v, want %v", rec.ops, want)
	}
}

func TestBeforeAbortSkipsInnerAndRunsAllAfters(t *testing.T) {
	innerCalled := false
	a := &recordInterceptor{name: "a"}
	b := &recordInterceptor{name: "b", beforeErr: errors.New("abort")}
	c := &recordInterceptor{name: "c"}
	chain := []intercept.Interceptor{a, b, c}
	r := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error { innerCalled = true; return nil },
	}
	got := WrapResource("alicloud_test", r, chain)
	err := got.Create(nil, nil)
	if err == nil {
		t.Fatal("expected the abort error")
	}
	if innerCalled {
		t.Fatal("inner must be skipped after a Before abort")
	}
	want := []string{"before:a", "before:b", "after:c", "after:b", "after:a"}
	if fmt.Sprint(a.ops) != fmt.Sprint([]string{"before:a", "after:a"}) ||
		fmt.Sprint(b.ops) != fmt.Sprint([]string{"before:b", "after:b"}) ||
		fmt.Sprint(c.ops) != fmt.Sprint([]string{"after:c"}) {
		t.Fatalf("ops out of contract: a=%v b=%v c=%v, want %v", a.ops, b.ops, c.ops, want)
	}
}

func TestAfterRewritesError(t *testing.T) {
	rec := &recordInterceptor{name: "rec", rewriteErr: errors.New("rewritten")}
	r := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error { return errors.New("inner") },
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	if err := got.Create(nil, nil); err == nil || err.Error() != "rewritten" {
		t.Fatalf("expected rewritten error, got %v", err)
	}
}

func TestAfterSeesInnerError(t *testing.T) {
	inner := errors.New("inner")
	rec := &recordInterceptor{name: "rec"}
	r := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error { return inner },
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	if err := got.Create(nil, nil); err != inner {
		t.Fatalf("expected the inner error unchanged, got %v", err)
	}
}

func TestWrapResourceCreateContext(t *testing.T) {
	innerCalled := false
	rec := &recordInterceptor{name: "rec"}
	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			innerCalled = true
			return nil
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	if got.CreateContext == nil {
		t.Fatal("CreateContext must be wrapped")
	}
	if got.Create != nil {
		t.Fatal("plain Create must be left untouched when CreateContext is wrapped")
	}
	diags := got.CreateContext(context.Background(), nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if !innerCalled {
		t.Fatal("inner not called")
	}
	want := []string{"before:rec", "after:rec"}
	if fmt.Sprint(rec.ops) != fmt.Sprint(want) {
		t.Fatalf("ops = %v, want %v", rec.ops, want)
	}
}

func TestWrapResourceCreateWithoutTimeout(t *testing.T) {
	withoutCalled := false
	contextCalled := false
	rec := &recordInterceptor{name: "rec"}
	r := &schema.Resource{
		CreateWithoutTimeout: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			withoutCalled = true
			return nil
		},
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			contextCalled = true
			return nil
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	if got.CreateWithoutTimeout == nil {
		t.Fatal("CreateWithoutTimeout must be wrapped")
	}
	diags := got.CreateWithoutTimeout(context.Background(), nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	if !withoutCalled {
		t.Fatal("WithoutTimeout inner not called")
	}
	if contextCalled {
		t.Fatal("CreateContext must not be called — SDK dispatches to WithoutTimeout")
	}
}

func TestWrapResourceMixedForms(t *testing.T) {
	createCalled, readCalled, updateCalled := false, false, false
	rec := &recordInterceptor{name: "rec"}
	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			createCalled = true
			return nil
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			readCalled = true
			return nil
		},
		UpdateWithoutTimeout: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			updateCalled = true
			return nil
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	if got.CreateContext == nil {
		t.Fatal("CreateContext must be wrapped")
	}
	if got.Read == nil {
		t.Fatal("Read must be wrapped")
	}
	if got.UpdateWithoutTimeout == nil {
		t.Fatal("UpdateWithoutTimeout must be wrapped")
	}
	if d := got.CreateContext(context.Background(), nil, nil); d.HasError() || !createCalled {
		t.Fatal("CreateContext failed")
	}
	if err := got.Read(nil, nil); err != nil || !readCalled {
		t.Fatal("Read failed")
	}
	if d := got.UpdateWithoutTimeout(context.Background(), nil, nil); d.HasError() || !updateCalled {
		t.Fatal("UpdateWithoutTimeout failed")
	}
}

func TestWrapResourceContextBeforeAbort(t *testing.T) {
	innerCalled := false
	rec := &recordInterceptor{name: "rec", beforeErr: errors.New("abort")}
	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			innerCalled = true
			return nil
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	diags := got.CreateContext(context.Background(), nil, nil)
	if !diags.HasError() {
		t.Fatal("expected error from Before abort")
	}
	if innerCalled {
		t.Fatal("inner must be skipped after Before abort")
	}
}

func TestContextFormDiagnosticsPassThroughVerbatim(t *testing.T) {
	rec := &recordInterceptor{name: "rec"}
	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Diagnostics{{Severity: diag.Error, Summary: "the API said no", Detail: "code: Throttling"}}
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	diags := got.CreateContext(context.Background(), nil, nil)

	if len(diags) != 1 {
		t.Fatalf("len(diags) = %d, want 1: %v", len(diags), diags)
	}
	if diags[0].Summary != "the API said no" || diags[0].Detail != "code: Throttling" {
		t.Fatalf("diagnostic was rewritten: %+v", diags[0])
	}
}

func TestContextFormAfterCannotSwallowError(t *testing.T) {
	rec := &recordInterceptor{name: "rec", swallow: true}
	r := &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Errorf("resource not found")
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	diags := got.ReadContext(context.Background(), nil, nil)

	if !diags.HasError() {
		t.Fatal("After returned nil, but a failure must not be reported as success")
	}
	if len(diags) != 1 || diags[0].Summary != "resource not found" {
		t.Fatalf("want the inner error restored verbatim, got %v", diags)
	}
}

func TestContextFormAfterRewritesError(t *testing.T) {
	rec := &recordInterceptor{name: "rec", rewriteErr: errors.New("rewritten")}
	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Errorf("inner")
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	diags := got.CreateContext(context.Background(), nil, nil)

	if len(diags) != 2 {
		t.Fatalf("len(diags) = %d, want 2 (inner kept, rewritten appended): %v", len(diags), diags)
	}
	if diags[0].Summary != "inner" {
		t.Fatalf("the inner error was replaced, summary = %q", diags[0].Summary)
	}
	if diags[1].Summary != "rewritten" {
		t.Fatalf("summary = %q, want %q", diags[1].Summary, "rewritten")
	}
}

// A hook that only wraps the chain's error must not cost the attribute path and
// detail that error came from.
func TestContextFormWrapKeepsInnerErrorsIntact(t *testing.T) {
	rec := &recordInterceptor{name: "rec", wrapMsg: "retry gave up"}
	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Diagnostics{
				{
					Severity:      diag.Error,
					Summary:       "the API said no",
					Detail:        "code: Throttling",
					AttributePath: cty.GetAttrPath("cidr_block"),
				},
				{Severity: diag.Error, Summary: "second failure", Detail: "with a detail"},
			}
		},
	}
	got := WrapResource("alicloud_test", r, []intercept.Interceptor{rec})
	diags := got.CreateContext(context.Background(), nil, nil)

	if len(diags) != 3 {
		t.Fatalf("len(diags) = %d, want 3: %v", len(diags), diags)
	}
	if got := diags[0].AttributePath; !got.Equals(cty.GetAttrPath("cidr_block")) {
		t.Fatalf("the attribute path was lost: %v", got)
	}
	if got := diags[1].Detail; got != "with a detail" {
		t.Fatalf("detail = %q, want it preserved", got)
	}
	// The hook's own wording arrives on top, carrying both messages it saw.
	want := "retry gave up: the API said no: code: Throttling; second failure: with a detail"
	if got := diags[2].Summary; got != want {
		t.Fatalf("summary = %q, want %q", got, want)
	}
}

func TestContextFormWarningsSurvive(t *testing.T) {
	warning := diag.Diagnostic{Severity: diag.Warning, Summary: "deprecated field"}
	for _, tc := range []struct {
		name string
		rec  *recordInterceptor
	}{
		{name: "swallow", rec: &recordInterceptor{name: "rec", swallow: true}},
		{name: "rewrite", rec: &recordInterceptor{name: "rec", rewriteErr: errors.New("rewritten")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &schema.Resource{
				CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
					return diag.Diagnostics{warning, {Severity: diag.Error, Summary: "inner"}}
				},
			}
			got := WrapResource("alicloud_test", r, []intercept.Interceptor{tc.rec})
			diags := got.CreateContext(context.Background(), nil, nil)

			if len(diags) == 0 || diags[0].Severity != diag.Warning || diags[0].Summary != warning.Summary {
				t.Fatalf("the warning must come first and unchanged, got %v", diags)
			}
		})
	}
}

func TestWrapDataSourceReadContext(t *testing.T) {
	innerCalled := false
	rec := &recordInterceptor{name: "rec"}
	ds := &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			innerCalled = true
			return nil
		},
	}
	got := WrapDataSource("alicloud_test", ds, []intercept.Interceptor{rec})
	if got.ReadContext == nil {
		t.Fatal("ReadContext must be wrapped")
	}
	if d := got.ReadContext(context.Background(), nil, nil); d.HasError() || !innerCalled {
		t.Fatal("ReadContext failed")
	}
}
