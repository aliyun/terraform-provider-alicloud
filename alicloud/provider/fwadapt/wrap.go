// Package fwadapt adapts the interceptor chain to Framework resources and data
// sources, and injects the client into them (see client.go).
//
// The Framework detects optional capabilities by type-asserting the registered
// resource, so the wrapper's method set decides which capabilities the resource
// appears to have. Seven of the nine forward an absent implementation harmlessly.
// Two do not: a silent ImportState imports empty state, a silent IdentitySchema
// fails every plan — hence the four wrapper types.
package fwadapt

import (
	"context"
	"errors"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func WrapResource(name string, r resource.Resource, chain []intercept.Interceptor) resource.Resource {
	if r == nil || len(chain) == 0 {
		return r
	}
	core := &wrappedResource{name: name, inner: r, chain: chain}
	_, withImportState := r.(resource.ResourceWithImportState)
	_, withIdentity := r.(resource.ResourceWithIdentity)
	switch {
	case withImportState && withIdentity:
		return &resourceWithImportStateAndIdentity{wrappedResource: core}
	case withImportState:
		return &resourceWithImportState{wrappedResource: core}
	case withIdentity:
		return &resourceWithIdentity{wrappedResource: core}
	default:
		return core
	}
}

func WrapDataSource(name string, ds datasource.DataSource, chain []intercept.Interceptor) datasource.DataSource {
	if ds == nil || len(chain) == 0 {
		return ds
	}
	return &wrappedDataSource{name: name, inner: ds, chain: chain}
}

type wrappedResource struct {
	name  string
	inner resource.Resource
	chain []intercept.Interceptor

	// Not shared across operations: the Framework builds a fresh resource per RPC.
	providerData interface{}
}

var (
	_ resource.Resource                     = (*wrappedResource)(nil)
	_ resource.ResourceWithConfigure        = (*wrappedResource)(nil)
	_ resource.ResourceWithConfigValidators = (*wrappedResource)(nil)
	_ resource.ResourceWithModifyPlan       = (*wrappedResource)(nil)
	_ resource.ResourceWithValidateConfig   = (*wrappedResource)(nil)
	_ resource.ResourceWithUpgradeState     = (*wrappedResource)(nil)
	_ resource.ResourceWithMoveState        = (*wrappedResource)(nil)
	_ resource.ResourceWithUpgradeIdentity  = (*wrappedResource)(nil)
)

func (w *wrappedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	w.inner.Metadata(ctx, req, resp)
}

func (w *wrappedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	w.inner.Schema(ctx, req, resp)
}

func (w *wrappedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	runIntercepted(ctx, w.name, intercept.OpCreate, w.chain, &resp.Diagnostics, w.providerData,
		func() { w.inner.Create(ctx, req, resp) })
}

func (w *wrappedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	runIntercepted(ctx, w.name, intercept.OpRead, w.chain, &resp.Diagnostics, w.providerData,
		func() { w.inner.Read(ctx, req, resp) })
}

func (w *wrappedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	runIntercepted(ctx, w.name, intercept.OpUpdate, w.chain, &resp.Diagnostics, w.providerData,
		func() { w.inner.Update(ctx, req, resp) })
}

func (w *wrappedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	runIntercepted(ctx, w.name, intercept.OpDelete, w.chain, &resp.Diagnostics, w.providerData,
		func() { w.inner.Delete(ctx, req, resp) })
}

func (w *wrappedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	w.providerData = req.ProviderData
	if c, ok := w.inner.(resource.ResourceWithConfigure); ok {
		c.Configure(ctx, req, resp)
	}
	resp.Diagnostics.Append(clientNotInjected(w.name, w.inner, req.ProviderData, resp.Diagnostics)...)
}

func (w *wrappedResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if c, ok := w.inner.(resource.ResourceWithConfigValidators); ok {
		return c.ConfigValidators(ctx)
	}
	return nil
}

func (w *wrappedResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if c, ok := w.inner.(resource.ResourceWithModifyPlan); ok {
		c.ModifyPlan(ctx, req, resp)
	}
}

func (w *wrappedResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if c, ok := w.inner.(resource.ResourceWithValidateConfig); ok {
		c.ValidateConfig(ctx, req, resp)
	}
}

func (w *wrappedResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	if c, ok := w.inner.(resource.ResourceWithUpgradeState); ok {
		return c.UpgradeState(ctx)
	}
	return nil
}

func (w *wrappedResource) MoveState(ctx context.Context) []resource.StateMover {
	if c, ok := w.inner.(resource.ResourceWithMoveState); ok {
		return c.MoveState(ctx)
	}
	return nil
}

func (w *wrappedResource) UpgradeIdentity(ctx context.Context) map[int64]resource.IdentityUpgrader {
	if c, ok := w.inner.(resource.ResourceWithUpgradeIdentity); ok {
		return c.UpgradeIdentity(ctx)
	}
	return nil
}

func (w *wrappedResource) importState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	w.inner.(resource.ResourceWithImportState).ImportState(ctx, req, resp)
}

func (w *wrappedResource) identitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	w.inner.(resource.ResourceWithIdentity).IdentitySchema(ctx, req, resp)
}

type resourceWithImportState struct {
	*wrappedResource
}

var _ resource.ResourceWithImportState = (*resourceWithImportState)(nil)

func (w *resourceWithImportState) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	w.importState(ctx, req, resp)
}

type resourceWithIdentity struct {
	*wrappedResource
}

var _ resource.ResourceWithIdentity = (*resourceWithIdentity)(nil)

func (w *resourceWithIdentity) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	w.identitySchema(ctx, req, resp)
}

// Declares both rather than embedding the two variants above: embedding both makes
// the promoted *wrappedResource methods ambiguous and drops them.
type resourceWithImportStateAndIdentity struct {
	*wrappedResource
}

var (
	_ resource.ResourceWithImportState = (*resourceWithImportStateAndIdentity)(nil)
	_ resource.ResourceWithIdentity    = (*resourceWithImportStateAndIdentity)(nil)
)

func (w *resourceWithImportStateAndIdentity) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	w.importState(ctx, req, resp)
}

func (w *resourceWithImportStateAndIdentity) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	w.identitySchema(ctx, req, resp)
}

type wrappedDataSource struct {
	name  string
	inner datasource.DataSource
	chain []intercept.Interceptor

	providerData interface{}
}

var (
	_ datasource.DataSource                     = (*wrappedDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*wrappedDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*wrappedDataSource)(nil)
	_ datasource.DataSourceWithValidateConfig   = (*wrappedDataSource)(nil)
)

func (w *wrappedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	w.inner.Metadata(ctx, req, resp)
}

func (w *wrappedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	w.inner.Schema(ctx, req, resp)
}

func (w *wrappedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	runIntercepted(ctx, w.name, intercept.OpRead, w.chain, &resp.Diagnostics, w.providerData,
		func() { w.inner.Read(ctx, req, resp) })
}

func (w *wrappedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	w.providerData = req.ProviderData
	if c, ok := w.inner.(datasource.DataSourceWithConfigure); ok {
		c.Configure(ctx, req, resp)
	}
	resp.Diagnostics.Append(clientNotInjected(w.name, w.inner, req.ProviderData, resp.Diagnostics)...)
}

func (w *wrappedDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	if c, ok := w.inner.(datasource.DataSourceWithConfigValidators); ok {
		return c.ConfigValidators(ctx)
	}
	return nil
}

func (w *wrappedDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if c, ok := w.inner.(datasource.DataSourceWithValidateConfig); ok {
		c.ValidateConfig(ctx, req, resp)
	}
}

// The Framework reports through resp.Diagnostics, so that pointer is both the
// input and the output of the bridge.
func runIntercepted(ctx context.Context, name string, op intercept.Op, chain []intercept.Interceptor, diags *diag.Diagnostics, meta interface{}, inner func()) {
	*diags = intercept.Around(ctx, chain, intercept.Call{Name: name, Op: op, Meta: meta},
		func() diag.Diagnostics {
			inner()
			return *diags
		}, bridge)
}

var bridge = intercept.DiagBridge[diag.Diagnostics]{
	ToError:   diagToErr,
	WithError: replaceDiagErrors,
}

func diagToErr(diags diag.Diagnostics) error {
	for _, d := range diags {
		if d.Severity() == diag.SeverityError {
			if d.Detail() != "" {
				return errors.New(d.Summary() + ": " + d.Detail())
			}
			return errors.New(d.Summary())
		}
	}
	return nil
}

func replaceDiagErrors(diags diag.Diagnostics, err error) diag.Diagnostics {
	out := make(diag.Diagnostics, 0, len(diags)+1)
	for _, d := range diags {
		if d.Severity() != diag.SeverityError {
			out = append(out, d)
		}
	}
	if err != nil {
		out = append(out, diag.NewErrorDiagnostic(err.Error(), ""))
	}
	return out
}
