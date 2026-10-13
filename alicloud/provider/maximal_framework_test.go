// The Framework counterpart of maximal_sdkv2_test.go.
package provider_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
)

const (
	maximalFrameworkName           = "alicloud_unit_test_fw_maximal"
	maximalFrameworkDataSourceName = "alicloud_unit_test_fw_maximal_data"
)

type maximalFrameworkModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type maximalFrameworkIdentityModel struct {
	ID types.String `tfsdk:"id"`
}

func setMaximalFrameworkIdentity(ctx context.Context, identity *tfsdk.ResourceIdentity, id types.String) diag.Diagnostics {
	if identity == nil {
		return nil
	}
	return identity.Set(ctx, maximalFrameworkIdentityModel{ID: id})
}

type maximalFrameworkResource struct {
	fwadapt.ResourceBase

	log *hookLog
}

var (
	_ resource.Resource                     = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithConfigure        = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithConfigValidators = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithModifyPlan       = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithValidateConfig   = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithUpgradeState     = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithMoveState        = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithUpgradeIdentity  = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithImportState      = (*maximalFrameworkResource)(nil)
	_ resource.ResourceWithIdentity         = (*maximalFrameworkResource)(nil)
	_ fwadapt.WithClient                    = (*maximalFrameworkResource)(nil)
)

func newMaximalFrameworkResource(log *hookLog) *maximalFrameworkResource {
	return &maximalFrameworkResource{log: log}
}

func (r *maximalFrameworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unit_test_fw_maximal"
}

func (r *maximalFrameworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = fwschema.Schema{
		Version: 2,
		Attributes: map[string]fwschema.Attribute{
			"id": fwschema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": fwschema.StringAttribute{Required: true},
		},
		Description:         "Feature-maximal Framework fixture for the interceptor layer.",
		MarkdownDescription: "Feature-maximal Framework fixture for the interceptor layer.",
		DeprecationMessage:  "This is a unit test fixture and is never registered in the real provider.",
	}
}

func (r *maximalFrameworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.recordCRUD("Create")
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan maximalFrameworkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue("unit-test")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(setMaximalFrameworkIdentity(ctx, resp.Identity, plan.ID)...)
}

func (r *maximalFrameworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.recordCRUD("Read")
	if req.State.Raw.IsNull() {
		return
	}
	var state maximalFrameworkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	resp.Diagnostics.Append(setMaximalFrameworkIdentity(ctx, resp.Identity, state.ID)...)
}

func (r *maximalFrameworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.recordCRUD("Update")
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan maximalFrameworkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(setMaximalFrameworkIdentity(ctx, resp.Identity, plan.ID)...)
}

func (r *maximalFrameworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.recordCRUD("Delete")
}

func (r *maximalFrameworkResource) recordCRUD(op string) {
	r.log.record(op)
	if r.Client() != nil {
		r.log.record(op + "/client")
	}
}

func (r *maximalFrameworkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	r.log.record("ConfigValidators")
	return []resource.ConfigValidator{recordingResourceValidator{log: r.log}}
}

func (r *maximalFrameworkResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.log.record("ModifyPlan")
}

func (r *maximalFrameworkResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	r.log.record("ValidateConfig")
}

func (r *maximalFrameworkResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	r.log.record("UpgradeState")
	priorSchema := fwschema.Schema{
		Attributes: map[string]fwschema.Attribute{
			"id":   fwschema.StringAttribute{Computed: true},
			"name": fwschema.StringAttribute{Required: true},
		},
	}
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &priorSchema,
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				r.log.record("UpgradeState0")
				r.carryStateForward(ctx, req, resp)
			},
		},
		1: {
			PriorSchema: &priorSchema,
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				r.log.record("UpgradeState1")
				r.carryStateForward(ctx, req, resp)
			},
		},
	}
}

func (r *maximalFrameworkResource) carryStateForward(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	if resp.State.Schema == nil || req.State == nil {
		return
	}
	var prior maximalFrameworkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &prior)...)
}

func (r *maximalFrameworkResource) MoveState(ctx context.Context) []resource.StateMover {
	r.log.record("MoveState")
	return []resource.StateMover{
		{
			StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				r.log.record("MoveState0")
				if resp.TargetState.Schema == nil {
					return
				}
				resp.Diagnostics.Append(resp.TargetState.Set(ctx, maximalFrameworkModel{
					ID:   types.StringValue("unit-test"),
					Name: types.StringValue("moved"),
				})...)
			},
		},
	}
}

func (r *maximalFrameworkResource) UpgradeIdentity(ctx context.Context) map[int64]resource.IdentityUpgrader {
	r.log.record("UpgradeIdentity")
	return map[int64]resource.IdentityUpgrader{
		0: {
			IdentityUpgrader: func(ctx context.Context, req resource.UpgradeIdentityRequest, resp *resource.UpgradeIdentityResponse) {
				r.log.record("UpgradeIdentity0")
				resp.Diagnostics.Append(setMaximalFrameworkIdentity(ctx, resp.Identity, types.StringValue("unit-test"))...)
			},
		},
	}
}

func (r *maximalFrameworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.log.record("ImportState")
	if resp.State.Schema == nil {
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *maximalFrameworkResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	r.log.record("IdentitySchema")
	resp.IdentitySchema = identityschema.Schema{
		Version: 1,
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}

type recordingResourceValidator struct {
	log *hookLog
}

var _ resource.ConfigValidator = recordingResourceValidator{}

func (v recordingResourceValidator) Description(context.Context) string { return "records that it ran" }

func (v recordingResourceValidator) MarkdownDescription(context.Context) string {
	return "records that it ran"
}

func (v recordingResourceValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	v.log.record("ConfigValidator")
}

type maximalFrameworkDataSource struct {
	fwadapt.DataSourceBase

	log *hookLog
}

var (
	_ datasource.DataSource                     = (*maximalFrameworkDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*maximalFrameworkDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*maximalFrameworkDataSource)(nil)
	_ datasource.DataSourceWithValidateConfig   = (*maximalFrameworkDataSource)(nil)
	_ fwadapt.WithClient                        = (*maximalFrameworkDataSource)(nil)
)

func newMaximalFrameworkDataSource(log *hookLog) *maximalFrameworkDataSource {
	return &maximalFrameworkDataSource{log: log}
}

func (d *maximalFrameworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unit_test_fw_maximal_data"
}

func (d *maximalFrameworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"id":   dsschema.StringAttribute{Computed: true},
			"name": dsschema.StringAttribute{Required: true},
		},
		Description: "Feature-maximal Framework data source fixture for the interceptor layer.",
	}
}

func (d *maximalFrameworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.log.record("Read")
	if d.Client() != nil {
		d.log.record("Read/client")
	}
	if req.Config.Raw.IsNull() {
		return
	}
	var config maximalFrameworkModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.ID = types.StringValue("unit-test")
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (d *maximalFrameworkDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	d.log.record("ConfigValidators")
	return []datasource.ConfigValidator{recordingDataSourceValidator{log: d.log}}
}

func (d *maximalFrameworkDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	d.log.record("ValidateConfig")
}

type recordingDataSourceValidator struct {
	log *hookLog
}

var _ datasource.ConfigValidator = recordingDataSourceValidator{}

func (v recordingDataSourceValidator) Description(context.Context) string {
	return "records that it ran"
}

func (v recordingDataSourceValidator) MarkdownDescription(context.Context) string {
	return "records that it ran"
}

func (v recordingDataSourceValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	v.log.record("ConfigValidator")
}

func TestUnitProviderMaximalFrameworkWrapKeepsEveryCapability(t *testing.T) {
	log := &hookLog{}
	inner := newMaximalFrameworkResource(log)
	wrapped := fwadapt.WrapResource(maximalFrameworkName, inner, []intercept.Interceptor{&chainRecorder{}})

	if wrapped == resource.Resource(inner) {
		t.Fatal("WrapResource returned the input: nothing would be intercepted")
	}

	if _, ok := wrapped.(resource.ResourceWithConfigure); !ok {
		t.Error("the wrapper lost ResourceWithConfigure")
	}
	if _, ok := wrapped.(resource.ResourceWithConfigValidators); !ok {
		t.Error("the wrapper lost ResourceWithConfigValidators")
	}
	if _, ok := wrapped.(resource.ResourceWithModifyPlan); !ok {
		t.Error("the wrapper lost ResourceWithModifyPlan")
	}
	if _, ok := wrapped.(resource.ResourceWithValidateConfig); !ok {
		t.Error("the wrapper lost ResourceWithValidateConfig")
	}
	if _, ok := wrapped.(resource.ResourceWithUpgradeState); !ok {
		t.Error("the wrapper lost ResourceWithUpgradeState")
	}
	if _, ok := wrapped.(resource.ResourceWithMoveState); !ok {
		t.Error("the wrapper lost ResourceWithMoveState")
	}
	if _, ok := wrapped.(resource.ResourceWithUpgradeIdentity); !ok {
		t.Error("the wrapper lost ResourceWithUpgradeIdentity")
	}
	if _, ok := wrapped.(resource.ResourceWithImportState); !ok {
		t.Error("the wrapper lost ResourceWithImportState")
	}
	if _, ok := wrapped.(resource.ResourceWithIdentity); !ok {
		t.Error("the wrapper lost ResourceWithIdentity")
	}
}

func TestUnitProviderMaximalFrameworkWrapForwardsEveryCapability(t *testing.T) {
	ctx := context.Background()
	log := &hookLog{}
	wrapped := fwadapt.WrapResource(maximalFrameworkName, newMaximalFrameworkResource(log), []intercept.Interceptor{&chainRecorder{}})

	t.Run("Metadata", func(t *testing.T) {
		var resp resource.MetadataResponse
		wrapped.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "alicloud"}, &resp)
		if resp.TypeName != maximalFrameworkName {
			t.Errorf("TypeName = %q, want %q", resp.TypeName, maximalFrameworkName)
		}
	})

	t.Run("Schema", func(t *testing.T) {
		var resp resource.SchemaResponse
		wrapped.Schema(ctx, resource.SchemaRequest{}, &resp)
		if resp.Schema.Version != 2 {
			t.Errorf("Schema.Version = %d, want 2", resp.Schema.Version)
		}
		if resp.Schema.DeprecationMessage == "" {
			t.Error("the schema's DeprecationMessage was lost")
		}
	})

	t.Run("ConfigValidators", func(t *testing.T) {
		validators := wrapped.(resource.ResourceWithConfigValidators).ConfigValidators(ctx)
		if len(validators) != 1 {
			t.Fatalf("got %d validators, want 1", len(validators))
		}
		validators[0].ValidateResource(ctx, resource.ValidateConfigRequest{}, &resource.ValidateConfigResponse{})
		if !log.has("ConfigValidator") {
			t.Error("the returned validator is not the fixture's")
		}
	})

	t.Run("ModifyPlan", func(t *testing.T) {
		wrapped.(resource.ResourceWithModifyPlan).ModifyPlan(ctx, resource.ModifyPlanRequest{}, &resource.ModifyPlanResponse{})
		if !log.has("ModifyPlan") {
			t.Error("ModifyPlan was not forwarded")
		}
	})

	t.Run("ValidateConfig", func(t *testing.T) {
		wrapped.(resource.ResourceWithValidateConfig).ValidateConfig(ctx, resource.ValidateConfigRequest{}, &resource.ValidateConfigResponse{})
		if !log.has("ValidateConfig") {
			t.Error("ValidateConfig was not forwarded")
		}
	})

	t.Run("UpgradeState", func(t *testing.T) {
		upgraders := wrapped.(resource.ResourceWithUpgradeState).UpgradeState(ctx)
		if len(upgraders) != 2 {
			t.Fatalf("got %d upgraders, want 2", len(upgraders))
		}
		for _, version := range []int64{0, 1} {
			u, ok := upgraders[version]
			if !ok {
				t.Fatalf("no upgrader for version %d", version)
			}
			if u.PriorSchema == nil {
				t.Errorf("upgrader %d lost its PriorSchema", version)
			}
			u.StateUpgrader(ctx, resource.UpgradeStateRequest{}, &resource.UpgradeStateResponse{})
		}
		if !log.has("UpgradeState0") || !log.has("UpgradeState1") {
			t.Errorf("an upgrader function did not run: %v", log.all())
		}
	})

	t.Run("MoveState", func(t *testing.T) {
		movers := wrapped.(resource.ResourceWithMoveState).MoveState(ctx)
		if len(movers) != 1 {
			t.Fatalf("got %d movers, want 1", len(movers))
		}
		movers[0].StateMover(ctx, resource.MoveStateRequest{}, &resource.MoveStateResponse{})
		if !log.has("MoveState0") {
			t.Error("the mover function did not run")
		}
	})

	t.Run("UpgradeIdentity", func(t *testing.T) {
		upgraders := wrapped.(resource.ResourceWithUpgradeIdentity).UpgradeIdentity(ctx)
		if len(upgraders) != 1 {
			t.Fatalf("got %d identity upgraders, want 1", len(upgraders))
		}
		u, ok := upgraders[0]
		if !ok {
			t.Fatal("no identity upgrader for version 0")
		}
		u.IdentityUpgrader(ctx, resource.UpgradeIdentityRequest{}, &resource.UpgradeIdentityResponse{})
		if !log.has("UpgradeIdentity0") {
			t.Error("the identity upgrader function did not run")
		}
	})

	t.Run("ImportState", func(t *testing.T) {
		var resp resource.ImportStateResponse
		wrapped.(resource.ResourceWithImportState).ImportState(ctx, resource.ImportStateRequest{ID: "unit-test"}, &resp)
		if !log.has("ImportState") {
			t.Error("ImportState was not forwarded")
		}
	})

	t.Run("IdentitySchema", func(t *testing.T) {
		var resp resource.IdentitySchemaResponse
		wrapped.(resource.ResourceWithIdentity).IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &resp)
		if !log.has("IdentitySchema") {
			t.Fatal("IdentitySchema was not forwarded")
		}
		if resp.IdentitySchema.Version != 1 {
			t.Errorf("IdentitySchema.Version = %d, want 1", resp.IdentitySchema.Version)
		}
		if _, ok := resp.IdentitySchema.Attributes["id"]; !ok {
			t.Error("the identity schema lost its attributes")
		}
	})
}

func TestUnitProviderMaximalFrameworkWrapRunsChainOnlyOnCRUD(t *testing.T) {
	ctx := context.Background()
	log := &hookLog{}
	rec := &chainRecorder{}
	client := &connectivity.AliyunClient{}
	wrapped := fwadapt.WrapResource(maximalFrameworkName, newMaximalFrameworkResource(log), []intercept.Interceptor{rec})

	var configureResp resource.ConfigureResponse
	wrapped.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("Configure through the wrapper failed: %v", configureResp.Diagnostics)
	}

	var schemaResp resource.SchemaResponse
	wrapped.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	wrapped.(resource.ResourceWithModifyPlan).ModifyPlan(ctx, resource.ModifyPlanRequest{}, &resource.ModifyPlanResponse{})
	wrapped.(resource.ResourceWithValidateConfig).ValidateConfig(ctx, resource.ValidateConfigRequest{}, &resource.ValidateConfigResponse{})
	wrapped.(resource.ResourceWithUpgradeState).UpgradeState(ctx)
	wrapped.(resource.ResourceWithMoveState).MoveState(ctx)
	wrapped.(resource.ResourceWithUpgradeIdentity).UpgradeIdentity(ctx)
	wrapped.(resource.ResourceWithImportState).ImportState(ctx, resource.ImportStateRequest{ID: "unit-test"}, &resource.ImportStateResponse{})
	wrapped.(resource.ResourceWithIdentity).IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &resource.IdentitySchemaResponse{})

	if got := rec.snapshotBefore(); len(got) != 0 {
		t.Fatalf("the chain ran outside CRUD: %v", got)
	}

	log.reset()
	var (
		createResp resource.CreateResponse
		readResp   resource.ReadResponse
		updateResp resource.UpdateResponse
		deleteResp resource.DeleteResponse
	)
	wrapped.Create(ctx, resource.CreateRequest{}, &createResp)
	wrapped.Read(ctx, resource.ReadRequest{}, &readResp)
	wrapped.Update(ctx, resource.UpdateRequest{}, &updateResp)
	wrapped.Delete(ctx, resource.DeleteRequest{}, &deleteResp)

	for _, diags := range []interface{ HasError() bool }{
		createResp.Diagnostics, readResp.Diagnostics, updateResp.Diagnostics, deleteResp.Diagnostics,
	} {
		if diags.HasError() {
			t.Errorf("a CRUD method reported an error: %v", diags)
		}
	}

	wantOps := []string{"Create", "Read", "Update", "Delete"}
	if got := rec.beforeOps(); !reflect.DeepEqual(got, wantOps) {
		t.Errorf("Before ops = %v, want %v", got, wantOps)
	}
	if got := rec.afterOps(); !reflect.DeepEqual(got, wantOps) {
		t.Errorf("After ops = %v, want %v", got, wantOps)
	}
	for _, name := range namesOf(rec.snapshotBefore()) {
		if name != maximalFrameworkName {
			t.Errorf("the chain saw name %q, want %q", name, maximalFrameworkName)
		}
	}

	for _, call := range rec.snapshotBefore() {
		if call.Meta != interface{}(client) {
			t.Errorf("Call.Meta = %v, want the injected client", call.Meta)
		}
	}

	want := []string{
		"Create", "Create/client",
		"Read", "Read/client",
		"Update", "Update/client",
		"Delete", "Delete/client",
	}
	if got := log.all(); !reflect.DeepEqual(got, want) {
		t.Errorf("inner CRUD ran as %v, want %v", got, want)
	}
}

func TestUnitProviderMaximalFrameworkDataSourceWrap(t *testing.T) {
	ctx := context.Background()
	log := &hookLog{}
	rec := &chainRecorder{}
	client := &connectivity.AliyunClient{}
	wrapped := fwadapt.WrapDataSource(maximalFrameworkDataSourceName, newMaximalFrameworkDataSource(log), []intercept.Interceptor{rec})

	var metadataResp datasource.MetadataResponse
	wrapped.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "alicloud"}, &metadataResp)
	if metadataResp.TypeName != maximalFrameworkDataSourceName {
		t.Errorf("TypeName = %q, want %q", metadataResp.TypeName, maximalFrameworkDataSourceName)
	}

	var configureResp datasource.ConfigureResponse
	wrapped.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("Configure through the wrapper failed: %v", configureResp.Diagnostics)
	}

	validators := wrapped.(datasource.DataSourceWithConfigValidators).ConfigValidators(ctx)
	if len(validators) != 1 {
		t.Fatalf("got %d validators, want 1", len(validators))
	}
	validators[0].ValidateDataSource(ctx, datasource.ValidateConfigRequest{}, &datasource.ValidateConfigResponse{})
	if !log.has("ConfigValidator") {
		t.Error("the returned validator is not the fixture's")
	}

	wrapped.(datasource.DataSourceWithValidateConfig).ValidateConfig(ctx, datasource.ValidateConfigRequest{}, &datasource.ValidateConfigResponse{})
	if !log.has("ValidateConfig") {
		t.Error("ValidateConfig was not forwarded")
	}

	if got := rec.snapshotBefore(); len(got) != 0 {
		t.Fatalf("the chain ran outside Read: %v", got)
	}

	log.reset()
	var readResp datasource.ReadResponse
	wrapped.Read(ctx, datasource.ReadRequest{}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Errorf("Read reported an error: %v", readResp.Diagnostics)
	}

	wantOps := []string{"Read"}
	if got := rec.beforeOps(); !reflect.DeepEqual(got, wantOps) {
		t.Errorf("Before ops = %v, want %v", got, wantOps)
	}
	if got := rec.afterOps(); !reflect.DeepEqual(got, wantOps) {
		t.Errorf("After ops = %v, want %v", got, wantOps)
	}
	if got := log.all(); !reflect.DeepEqual(got, []string{"Read", "Read/client"}) {
		t.Errorf("inner Read ran as %v, want [Read Read/client]", got)
	}
	for _, call := range rec.snapshotBefore() {
		if call.Meta != interface{}(client) {
			t.Errorf("Call.Meta = %v, want the injected client", call.Meta)
		}
	}
}
