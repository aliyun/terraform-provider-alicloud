// Package cms. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package cms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const bizTraceTypeName = "alicloud_cms_biz_trace"

var (
	_ resource.Resource                = &bizTraceResource{}
	_ resource.ResourceWithImportState = &bizTraceResource{}
)

func NewBizTraceResource() resource.Resource {
	return &bizTraceResource{}
}

type bizTraceResource struct {
	fwadapt.ResourceBase
}

type bizTraceResourceModel struct {
	Id             types.String `tfsdk:"id"`
	FwTimeouts     types.Object `tfsdk:"timeouts"`
	AdvancedConfig types.String `tfsdk:"advanced_config"`
	BizTraceCode   types.String `tfsdk:"biz_trace_code"`
	BizTraceName   types.String `tfsdk:"biz_trace_name"`
	CreateTime     types.String `tfsdk:"create_time"`
	RegionId       types.String `tfsdk:"region_id"`
	RuleConfig     types.String `tfsdk:"rule_config"`
	Workspace      types.String `tfsdk:"workspace"`
}

func (r *bizTraceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = bizTraceTypeName
}

func (r *bizTraceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"advanced_config": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Advanced Configuration",
				PlanModifiers: []planmodifier.String{
					fwJSONSubsetSuppress{},
				},
			},
			"biz_trace_code": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The service link code, which is used for the value of the service link label.",
			},
			"biz_trace_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the resource",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource",
				PlanModifiers: []planmodifier.String{
					// Server-set and stable; keep the prior state value on refresh
					// so a Computed attribute not present in config is not planned
					// as unknown ("known after apply"), which would yield a perpetual
					// non-empty refresh plan.
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_config": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Business Trace Rule Configuration",
				PlanModifiers: []planmodifier.String{
					fwJSONSubsetSuppress{},
				},
			},
			"workspace": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Work space",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *bizTraceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bizTraceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ClientToken": buildClientToken("CreateBizTrace"),
	}
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	createValues["ResourceType"] = "BizTrace"
	if !plan.AdvancedConfig.IsNull() && !plan.AdvancedConfig.IsUnknown() {
		createValues["AdvancedConfig"] = fwValueToInterface(plan.AdvancedConfig, map[string]fwWireField{}, false)
	}
	if !plan.BizTraceCode.IsNull() && !plan.BizTraceCode.IsUnknown() {
		createValues["BizTraceCode"] = fwValueToInterface(plan.BizTraceCode, map[string]fwWireField{}, false)
	}
	if !plan.BizTraceName.IsNull() && !plan.BizTraceName.IsUnknown() {
		createValues["BizTraceName"] = fwValueToInterface(plan.BizTraceName, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.RegionId.IsNull() && !plan.RegionId.IsUnknown() {
		createValues["RegionId"] = fwValueToInterface(plan.RegionId, map[string]fwWireField{}, false)
	}
	if !plan.RuleConfig.IsNull() && !plan.RuleConfig.IsUnknown() {
		createValues["RuleConfig"] = fwValueToInterface(plan.RuleConfig, map[string]fwWireField{}, false)
	}
	if !plan.Workspace.IsNull() && !plan.Workspace.IsUnknown() {
		createValues["Workspace"] = fwValueToInterface(plan.Workspace, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.RuleConfig", Target: "ruleConfig", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BizTraceCode", Target: "bizTraceCode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedConfig", Target: "advancedConfig", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BizTraceName", Target: "bizTraceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "regionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Workspace", Target: "workspace", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, true); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	// The EntityStore for the workspace must be explicitly initialised before
	// BizTrace creation; it is not auto-provisioned by workspace creation.
	if !plan.Workspace.IsNull() && !plan.Workspace.IsUnknown() {
		if err := createCmsEntityStore(r.Client(), plan.Workspace.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: initialising EntityStore", bizTraceTypeName),
				errs.WrapErrorf(err, "Creating %s: CreateEntityStore for workspace %s", bizTraceTypeName, plan.Workspace.ValueString()).Error(),
			)
			return
		}
	}

	// NotFound.EntityStore is returned right after workspace creation while the
	// backend entityStore resource is still being asynchronously provisioned;
	// retry until it becomes ready. Keep needRetry for throttling/5xx/network.
	response, err := retryRoaPostWithExtraRetry(r.Client(), "Cms", "2024-03-30", "/bizTrace", request, operationTimeout, map[string]string{"regionId": "host", "workspace": "body", "bizTraceCode": "body", "bizTraceName": "body", "ruleConfig": "body", "advancedConfig": "body", "__api": "CreateBizTrace"}, func(e error) bool {
		return strings.Contains(e.Error(), "NotFound.EntityStore")
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateBizTrace", bizTraceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateBizTrace", bizTraceTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(response, fwMappingPath("$.bizTraceId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_cms_biz_trace: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateBizTrace", bizTraceTypeName),
			fmt.Sprintf("The response carried no BizTraceId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.AdvancedConfig = fwKnownValue(trackedPlan.AdvancedConfig).(types.String)
	trackedPlan.BizTraceCode = fwKnownValue(trackedPlan.BizTraceCode).(types.String)
	trackedPlan.BizTraceName = fwKnownValue(trackedPlan.BizTraceName).(types.String)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.String)
	trackedPlan.RegionId = fwKnownValue(trackedPlan.RegionId).(types.String)
	trackedPlan.RuleConfig = fwKnownValue(trackedPlan.RuleConfig).(types.String)
	trackedPlan.Workspace = fwKnownValue(trackedPlan.Workspace).(types.String)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := bizTraceGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", bizTraceTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", bizTraceTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", bizTraceTypeName),
			fmt.Sprintf("The resource %s is not visible via GetBizTrace right after creation", id),
		)
		return
	}
	createdState := bizTraceState(object, &plan)
	// rule_config and advanced_config are JSON documents the CMS server enriches
	// on read with server-resolved keys (callType, entranceServiceId, dyingScope);
	// the configured value is the canonical subset. Preserve the configured value
	// in post-write state so the apply-phase consistency check sees the same value
	// the plan announced — the server-enriched superset is reconciled on subsequent
	// plans/refresh by fwJSONSubsetSuppress. The other fields are returned by the
	// server unchanged (biz_trace_code/biz_trace_name) or are server-only
	// (create_time/region_id), so bizTraceState already matches the plan for them.
	if !plan.RuleConfig.IsNull() && !plan.RuleConfig.IsUnknown() {
		createdState.RuleConfig = plan.RuleConfig
	}
	if !plan.AdvancedConfig.IsNull() && !plan.AdvancedConfig.IsUnknown() {
		createdState.AdvancedConfig = plan.AdvancedConfig
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, createdState)...)
}

func (r *bizTraceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bizTraceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := bizTraceGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBizTrace", bizTraceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBizTrace", bizTraceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, bizTraceState(object, &state))...)
}

func (r *bizTraceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state bizTraceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "update", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	planValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	planValues["ResourceType"] = "BizTrace"
	if !plan.AdvancedConfig.IsNull() && !plan.AdvancedConfig.IsUnknown() {
		planValues["AdvancedConfig"] = fwValueToInterface(plan.AdvancedConfig, map[string]fwWireField{}, false)
	}
	if !plan.BizTraceCode.IsNull() && !plan.BizTraceCode.IsUnknown() {
		planValues["BizTraceCode"] = fwValueToInterface(plan.BizTraceCode, map[string]fwWireField{}, false)
	}
	if !plan.BizTraceName.IsNull() && !plan.BizTraceName.IsUnknown() {
		planValues["BizTraceName"] = fwValueToInterface(plan.BizTraceName, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		planValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.RegionId.IsNull() && !plan.RegionId.IsUnknown() {
		planValues["RegionId"] = fwValueToInterface(plan.RegionId, map[string]fwWireField{}, false)
	}
	if !plan.RuleConfig.IsNull() && !plan.RuleConfig.IsUnknown() {
		planValues["RuleConfig"] = fwValueToInterface(plan.RuleConfig, map[string]fwWireField{}, false)
	}
	if !plan.Workspace.IsNull() && !plan.Workspace.IsUnknown() {
		planValues["Workspace"] = fwValueToInterface(plan.Workspace, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	stateValues["ResourceType"] = "BizTrace"
	if !state.AdvancedConfig.IsNull() && !state.AdvancedConfig.IsUnknown() {
		stateValues["AdvancedConfig"] = fwValueToInterface(state.AdvancedConfig, map[string]fwWireField{}, false)
	}
	if !state.BizTraceCode.IsNull() && !state.BizTraceCode.IsUnknown() {
		stateValues["BizTraceCode"] = fwValueToInterface(state.BizTraceCode, map[string]fwWireField{}, false)
	}
	if !state.BizTraceName.IsNull() && !state.BizTraceName.IsUnknown() {
		stateValues["BizTraceName"] = fwValueToInterface(state.BizTraceName, map[string]fwWireField{}, false)
	}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		stateValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.RegionId.IsNull() && !state.RegionId.IsUnknown() {
		stateValues["RegionId"] = fwValueToInterface(state.RegionId, map[string]fwWireField{}, false)
	}
	if !state.RuleConfig.IsNull() && !state.RuleConfig.IsUnknown() {
		stateValues["RuleConfig"] = fwValueToInterface(state.RuleConfig, map[string]fwWireField{}, false)
	}
	if !state.Workspace.IsNull() && !state.Workspace.IsUnknown() {
		stateValues["Workspace"] = fwValueToInterface(state.Workspace, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.RuleConfig, state.RuleConfig) {
		updated0 = true
	}
	if fwPlannedChanged(plan.AdvancedConfig, state.AdvancedConfig) {
		updated0 = true
	}
	if fwPlannedChanged(plan.BizTraceName, state.BizTraceName) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Workspace, state.Workspace) {
		updated0 = true
	}
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"bizTraceId":  state.Id.ValueString(),
			"ClientToken": buildClientToken("UpdateBizTrace"),
		}
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.RuleConfig", Target: "ruleConfig", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedConfig", Target: "advancedConfig", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BizTraceName", Target: "bizTraceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BizTraceId", Target: "bizTraceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "regionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Workspace", Target: "workspace", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues0, stateValues, []string{"$.BizTraceId"}), true); err != nil {
			resp.Diagnostics.AddError("Building UpdateBizTrace request", err.Error())
			return
		}
		if _, err := retryRoaPut(r.Client(), "Cms", "2024-03-30", "/bizTrace/{bizTraceId}", request0, operationTimeout, map[string]string{"regionId": "host", "bizTraceId": "path", "bizTraceName": "body", "workspace": "body", "ruleConfig": "body", "advancedConfig": "body", "__api": "UpdateBizTrace"}); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_cms_biz_trace: UpdateBizTrace", err.Error())
			return
		}
	}
	object, err := bizTraceGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBizTrace", bizTraceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBizTrace", bizTraceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	updatedState := bizTraceState(object, &plan)
	// See Create: preserve the configured rule_config/advanced_config subset in
	// post-write state so the apply-phase consistency check matches the plan;
	// the server-enriched superset is reconciled by fwJSONSubsetSuppress.
	if !plan.RuleConfig.IsNull() && !plan.RuleConfig.IsUnknown() {
		updatedState.RuleConfig = plan.RuleConfig
	}
	if !plan.AdvancedConfig.IsNull() && !plan.AdvancedConfig.IsUnknown() {
		updatedState.AdvancedConfig = plan.AdvancedConfig
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
}

func (r *bizTraceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bizTraceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	deleteValues["ResourceType"] = "BizTrace"
	if !state.AdvancedConfig.IsNull() && !state.AdvancedConfig.IsUnknown() {
		deleteValues["AdvancedConfig"] = fwValueToInterface(state.AdvancedConfig, map[string]fwWireField{}, false)
	}
	if !state.BizTraceCode.IsNull() && !state.BizTraceCode.IsUnknown() {
		deleteValues["BizTraceCode"] = fwValueToInterface(state.BizTraceCode, map[string]fwWireField{}, false)
	}
	if !state.BizTraceName.IsNull() && !state.BizTraceName.IsUnknown() {
		deleteValues["BizTraceName"] = fwValueToInterface(state.BizTraceName, map[string]fwWireField{}, false)
	}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		deleteValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.RegionId.IsNull() && !state.RegionId.IsUnknown() {
		deleteValues["RegionId"] = fwValueToInterface(state.RegionId, map[string]fwWireField{}, false)
	}
	if !state.RuleConfig.IsNull() && !state.RuleConfig.IsUnknown() {
		deleteValues["RuleConfig"] = fwValueToInterface(state.RuleConfig, map[string]fwWireField{}, false)
	}
	if !state.Workspace.IsNull() && !state.Workspace.IsUnknown() {
		deleteValues["Workspace"] = fwValueToInterface(state.Workspace, map[string]fwWireField{}, false)
	}
	deleteValues["BizTraceId"] = state.Id.ValueString()
	needsRefresh := false
	if needsRefresh {
		object, err := bizTraceGet(r.Client(), state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Reading deletion conditions", err.Error())
			return
		}
		if object == nil {
			return
		}
		for key, value := range object {
			if value != nil {
				deleteValues[key] = value
			}
		}
	}
	selected := false
	deleteEligible0, deleteConditionErr0 := fwDeleteEligible(deleteValues, []fwDeleteCondition{})
	if deleteConditionErr0 != nil {
		resp.Diagnostics.AddError("Selecting deletion operation", deleteConditionErr0.Error())
		return
	}
	selected = selected || deleteEligible0
	if !selected {
		resp.Diagnostics.AddError("Selecting deletion operation", "No deletion operation matches the resource state; the resource remains tracked")
		return
	}
	if deleteEligible0 {
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteBizTrace")}
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.BizTraceId", Target: "bizTraceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "regionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, true); err != nil {
			resp.Diagnostics.AddError("Building DeleteBizTrace request", err.Error())
			return
		}
		if _, err := retryRoaDelete(r.Client(), "Cms", "2024-03-30", "/bizTrace/{bizTraceId}", request0, operationTimeout, map[string]string{"bizTraceId": "path", "regionId": "host", "__api": "DeleteBizTrace"}); err != nil && !isExpectedErrors(err, []string{"404"}) {
			resp.Diagnostics.AddError("Deleting alicloud_cms_biz_trace: DeleteBizTrace", err.Error())
			return
		}
	}
}

func (r *bizTraceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// bizTraceState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func bizTraceState(object map[string]interface{}, prior *bizTraceResourceModel) bizTraceResourceModel {
	// Get has already merged and normalized every response mapping.
	str := func(key string) types.String {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.StringNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.StringNull()
			}
		}
		// gap#7c: present-but-empty strings must return ValueString(""), not Null.
		// v1 d.Set(field, "") treats empty as set (not absence); TestCheckResourceAttrSet
		// expects a non-null value for API-returned "" (e.g. DhcpOptionsSetId when VPC has
		// no DHCP set). The earlier && v != "" collapsed present-empty to Null, failing Set
		// checks. Absent fields still return Null via the !ok branches above.
		if v, ok := cur.(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m := bizTraceResourceModel{
		FwTimeouts:     types.ObjectNull(fwTimeoutTypes()),
		Id:             str("BizTraceId"),
		AdvancedConfig: str("AdvancedConfig"),
		BizTraceCode:   str("BizTraceCode"),
		BizTraceName:   str("BizTraceName"),
		CreateTime:     str("CreateTime"),
		RegionId:       str("RegionId"),
		RuleConfig:     str("RuleConfig"),
		Workspace:      str("Workspace"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
	}
	return m
}
