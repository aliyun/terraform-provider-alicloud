// Package dts. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package dts

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

const gadInstanceTypeName = "alicloud_dts_gad_instance"

var (
	_ resource.Resource                = &gadInstanceResource{}
	_ resource.ResourceWithImportState = &gadInstanceResource{}
)

func NewGadInstanceResource() resource.Resource {
	return &gadInstanceResource{}
}

type gadInstanceResource struct {
	fwadapt.ResourceBase
}

type gadInstanceResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	FwTimeouts             types.Object `tfsdk:"timeouts"`
	CreateTime             types.String `tfsdk:"create_time"`
	DbEngineType           types.String `tfsdk:"db_engine_type"`
	GadInstanceName        types.String `tfsdk:"gad_instance_name"`
	GadInstanceType        types.String `tfsdk:"gad_instance_type"`
	MasterDbInstanceId     types.String `tfsdk:"master_db_instance_id"`
	MasterDbInstanceRegion types.String `tfsdk:"master_db_instance_region"`
	MasterDbZoneId         types.String `tfsdk:"master_db_zone_id"`
	RegionId               types.String `tfsdk:"region_id"`
	ResourceGroupId        types.String `tfsdk:"resource_group_id"`
	Status                 types.String `tfsdk:"status"`
}

func (r *gadInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = gadInstanceTypeName
}

func (r *gadInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource",
			},
			"db_engine_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "Database engine type",
			},
			"gad_instance_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "GAD instance name",
			},
			"gad_instance_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "GAD instance type",
			},
			"master_db_instance_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "Master database instance ID",
			},
			"master_db_instance_region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "Master database instance region",
			},
			"master_db_zone_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "Master database instance zone ID",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource",
			},
			"resource_group_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "The ID of the resource group. Create-only: CreateGadInstance accepts ResourceGroupId at creation, but DTS exposes no GAD-specific resource-group relocation API and ConvertInstanceResourceGroup (the DTS instance relocation API) returns Forbidden.InstanceNotFound for GAD instance IDs (gad-*), so a resource_group_id change forces replacement.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "GAD instance status",
			},
		},
	}
}

func (r *gadInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan gadInstanceResourceModel
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
		"ClientToken": buildClientToken("CreateGadInstance"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.DbEngineType.IsNull() && !plan.DbEngineType.IsUnknown() {
		createValues["DbEngineType"] = fwValueToInterface(plan.DbEngineType, map[string]fwWireField{}, false)
	}
	if !plan.GadInstanceName.IsNull() && !plan.GadInstanceName.IsUnknown() {
		createValues["GadInstanceName"] = fwValueToInterface(plan.GadInstanceName, map[string]fwWireField{}, false)
	}
	if !plan.GadInstanceType.IsNull() && !plan.GadInstanceType.IsUnknown() {
		createValues["GadInstanceType"] = fwValueToInterface(plan.GadInstanceType, map[string]fwWireField{}, false)
	}
	if !plan.MasterDbInstanceId.IsNull() && !plan.MasterDbInstanceId.IsUnknown() {
		createValues["MasterDbInstanceId"] = fwValueToInterface(plan.MasterDbInstanceId, map[string]fwWireField{}, false)
	}
	if !plan.MasterDbInstanceRegion.IsNull() && !plan.MasterDbInstanceRegion.IsUnknown() {
		createValues["MasterDbInstanceRegion"] = fwValueToInterface(plan.MasterDbInstanceRegion, map[string]fwWireField{}, false)
	}
	if !plan.MasterDbZoneId.IsNull() && !plan.MasterDbZoneId.IsUnknown() {
		createValues["MasterDbZoneId"] = fwValueToInterface(plan.MasterDbZoneId, map[string]fwWireField{}, false)
	}
	if !plan.RegionId.IsNull() && !plan.RegionId.IsUnknown() {
		createValues["RegionId"] = fwValueToInterface(plan.RegionId, map[string]fwWireField{}, false)
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		createValues["ResourceGroupId"] = fwValueToInterface(plan.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		createValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.DbEngineType", Target: "MasterDbEngineType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.GadInstanceName", Target: "InstanceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.GadInstanceType", Target: "InstanceType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.MasterDbInstanceId", Target: "MasterDbInstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.MasterDbInstanceRegion", Target: "MasterDbInstanceRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.MasterDbZoneId", Target: "MasterDbZoneId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceGroupId", Target: "ResourceGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "Dts", "2020-01-01", "CreateGadInstance", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateGadInstance", gadInstanceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateGadInstance", gadInstanceTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(response, fwMappingPath("$.InstanceId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_dts_gad_instance: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateGadInstance", gadInstanceTypeName),
			fmt.Sprintf("The response carried no GadInstanceId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.String)
	trackedPlan.DbEngineType = fwKnownValue(trackedPlan.DbEngineType).(types.String)
	trackedPlan.GadInstanceName = fwKnownValue(trackedPlan.GadInstanceName).(types.String)
	trackedPlan.GadInstanceType = fwKnownValue(trackedPlan.GadInstanceType).(types.String)
	trackedPlan.MasterDbInstanceId = fwKnownValue(trackedPlan.MasterDbInstanceId).(types.String)
	trackedPlan.MasterDbInstanceRegion = fwKnownValue(trackedPlan.MasterDbInstanceRegion).(types.String)
	trackedPlan.MasterDbZoneId = fwKnownValue(trackedPlan.MasterDbZoneId).(types.String)
	trackedPlan.RegionId = fwKnownValue(trackedPlan.RegionId).(types.String)
	trackedPlan.ResourceGroupId = fwKnownValue(trackedPlan.ResourceGroupId).(types.String)
	trackedPlan.Status = fwKnownValue(trackedPlan.Status).(types.String)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := gadInstanceGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", gadInstanceTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", gadInstanceTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", gadInstanceTypeName),
			fmt.Sprintf("The resource %s is not visible via DescribeGadInstanceDetail right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, gadInstanceState(object, &plan))...)
}

func (r *gadInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state gadInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := gadInstanceGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeGadInstanceDetail", gadInstanceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeGadInstanceDetail", gadInstanceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, gadInstanceState(object, &state))...)
}

func (r *gadInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state gadInstanceResourceModel
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
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		planValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.DbEngineType.IsNull() && !plan.DbEngineType.IsUnknown() {
		planValues["DbEngineType"] = fwValueToInterface(plan.DbEngineType, map[string]fwWireField{}, false)
	}
	if !plan.GadInstanceName.IsNull() && !plan.GadInstanceName.IsUnknown() {
		planValues["GadInstanceName"] = fwValueToInterface(plan.GadInstanceName, map[string]fwWireField{}, false)
	}
	if !plan.GadInstanceType.IsNull() && !plan.GadInstanceType.IsUnknown() {
		planValues["GadInstanceType"] = fwValueToInterface(plan.GadInstanceType, map[string]fwWireField{}, false)
	}
	if !plan.MasterDbInstanceId.IsNull() && !plan.MasterDbInstanceId.IsUnknown() {
		planValues["MasterDbInstanceId"] = fwValueToInterface(plan.MasterDbInstanceId, map[string]fwWireField{}, false)
	}
	if !plan.MasterDbInstanceRegion.IsNull() && !plan.MasterDbInstanceRegion.IsUnknown() {
		planValues["MasterDbInstanceRegion"] = fwValueToInterface(plan.MasterDbInstanceRegion, map[string]fwWireField{}, false)
	}
	if !plan.MasterDbZoneId.IsNull() && !plan.MasterDbZoneId.IsUnknown() {
		planValues["MasterDbZoneId"] = fwValueToInterface(plan.MasterDbZoneId, map[string]fwWireField{}, false)
	}
	if !plan.RegionId.IsNull() && !plan.RegionId.IsUnknown() {
		planValues["RegionId"] = fwValueToInterface(plan.RegionId, map[string]fwWireField{}, false)
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		planValues["ResourceGroupId"] = fwValueToInterface(plan.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		planValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		stateValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.DbEngineType.IsNull() && !state.DbEngineType.IsUnknown() {
		stateValues["DbEngineType"] = fwValueToInterface(state.DbEngineType, map[string]fwWireField{}, false)
	}
	if !state.GadInstanceName.IsNull() && !state.GadInstanceName.IsUnknown() {
		stateValues["GadInstanceName"] = fwValueToInterface(state.GadInstanceName, map[string]fwWireField{}, false)
	}
	if !state.GadInstanceType.IsNull() && !state.GadInstanceType.IsUnknown() {
		stateValues["GadInstanceType"] = fwValueToInterface(state.GadInstanceType, map[string]fwWireField{}, false)
	}
	if !state.MasterDbInstanceId.IsNull() && !state.MasterDbInstanceId.IsUnknown() {
		stateValues["MasterDbInstanceId"] = fwValueToInterface(state.MasterDbInstanceId, map[string]fwWireField{}, false)
	}
	if !state.MasterDbInstanceRegion.IsNull() && !state.MasterDbInstanceRegion.IsUnknown() {
		stateValues["MasterDbInstanceRegion"] = fwValueToInterface(state.MasterDbInstanceRegion, map[string]fwWireField{}, false)
	}
	if !state.MasterDbZoneId.IsNull() && !state.MasterDbZoneId.IsUnknown() {
		stateValues["MasterDbZoneId"] = fwValueToInterface(state.MasterDbZoneId, map[string]fwWireField{}, false)
	}
	if !state.RegionId.IsNull() && !state.RegionId.IsUnknown() {
		stateValues["RegionId"] = fwValueToInterface(state.RegionId, map[string]fwWireField{}, false)
	}
	if !state.ResourceGroupId.IsNull() && !state.ResourceGroupId.IsUnknown() {
		stateValues["ResourceGroupId"] = fwValueToInterface(state.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		stateValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	// ModifyGadInstanceName renames the GAD instance only. Its ResourceGroupId query parameter is an
	// authorization scope, not a mutation target — the API documents it as "特殊业务专用字段,无需传参".
	// The generated Update previously bundled ResourceGroupId into ModifyGadInstanceName, so a
	// resource_group_id change was sent to the rename API, which does not move the resource, leaving
	// state drifting behind config. Rename no longer sends ResourceGroupId.
	//
	// resource_group_id is create-only (ForceNew; see Schema PlanModifiers). DTS exposes no
	// GAD-specific resource-group relocation API (Dts/2020-01-01 GAD surface: Create/Delete/
	// ModifyName/Describe/Detach only), and ConvertInstanceResourceGroup — the relocation API
	// alicloud_dts_instance uses for regular DTS instances (dts-*) — returns
	// Forbidden.InstanceNotFound for GAD instance IDs (gad-*). A resource_group_id change therefore
	// forces replacement (RequiresReplaceIfConfigured) rather than an in-place relocation call.
	if fwPlannedChanged(plan.GadInstanceName, state.GadInstanceName) {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"InstanceId":  state.Id.ValueString(),
			"ClientToken": buildClientToken("ModifyGadInstanceName"),
		}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.GadInstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.GadInstanceName", Target: "InstanceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues0, stateValues, []string{"$.GadInstanceId"}), false); err != nil {
			resp.Diagnostics.AddError("Building ModifyGadInstanceName request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Dts", "2020-01-01", "ModifyGadInstanceName", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_dts_gad_instance: ModifyGadInstanceName", err.Error())
			return
		}
	}
	// resource_group_id changes never reach Update: RequiresReplaceIfConfigured in the Schema
	// forces replacement before Update runs, so the change-detection branch that would relocate the
	// instance in place is intentionally absent. GAD has no in-place relocation API —
	// ConvertInstanceResourceGroup (used by alicloud_dts_instance for dts-* ids) returns
	// Forbidden.InstanceNotFound for gad-* ids (AccTest Step 3) — so no in-place call is made.
	// DetachGadInstanceDbMember is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	object, err := gadInstanceGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeGadInstanceDetail", gadInstanceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeGadInstanceDetail", gadInstanceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, gadInstanceState(object, &plan))...)
}

func (r *gadInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state gadInstanceResourceModel
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
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		deleteValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.DbEngineType.IsNull() && !state.DbEngineType.IsUnknown() {
		deleteValues["DbEngineType"] = fwValueToInterface(state.DbEngineType, map[string]fwWireField{}, false)
	}
	if !state.GadInstanceName.IsNull() && !state.GadInstanceName.IsUnknown() {
		deleteValues["GadInstanceName"] = fwValueToInterface(state.GadInstanceName, map[string]fwWireField{}, false)
	}
	if !state.GadInstanceType.IsNull() && !state.GadInstanceType.IsUnknown() {
		deleteValues["GadInstanceType"] = fwValueToInterface(state.GadInstanceType, map[string]fwWireField{}, false)
	}
	if !state.MasterDbInstanceId.IsNull() && !state.MasterDbInstanceId.IsUnknown() {
		deleteValues["MasterDbInstanceId"] = fwValueToInterface(state.MasterDbInstanceId, map[string]fwWireField{}, false)
	}
	if !state.MasterDbInstanceRegion.IsNull() && !state.MasterDbInstanceRegion.IsUnknown() {
		deleteValues["MasterDbInstanceRegion"] = fwValueToInterface(state.MasterDbInstanceRegion, map[string]fwWireField{}, false)
	}
	if !state.MasterDbZoneId.IsNull() && !state.MasterDbZoneId.IsUnknown() {
		deleteValues["MasterDbZoneId"] = fwValueToInterface(state.MasterDbZoneId, map[string]fwWireField{}, false)
	}
	if !state.RegionId.IsNull() && !state.RegionId.IsUnknown() {
		deleteValues["RegionId"] = fwValueToInterface(state.RegionId, map[string]fwWireField{}, false)
	}
	if !state.ResourceGroupId.IsNull() && !state.ResourceGroupId.IsUnknown() {
		deleteValues["ResourceGroupId"] = fwValueToInterface(state.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		deleteValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	deleteValues["GadInstanceId"] = state.Id.ValueString()
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
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteGadInstance")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.GadInstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ResourceGroupId", Target: "ResourceGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteGadInstance request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Dts", "2020-01-01", "DeleteGadInstance", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"DTS.Msg.GADInstanceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_dts_gad_instance: DeleteGadInstance", err.Error())
			return
		}
	}
}

func (r *gadInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// gadInstanceState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func gadInstanceState(object map[string]interface{}, prior *gadInstanceResourceModel) gadInstanceResourceModel {
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
		// F8: DescribeGadInstanceDetail returns CreateTime as a numeric millisecond
		// timestamp (e.g. 1770393523000), not a string. The StringAttribute schema needs the
		// value coerced; without this the assertion above fails and create_time reads as
		// null, failing the NotNull post-apply refresh check. strconv avoids the scientific
		// notation fmt.Sprint produces for large float64 values (1.7703935233e+12).
		if v, ok := cur.(json.Number); ok {
			return types.StringValue(v.String())
		}
		if v, ok := cur.(float64); ok {
			return types.StringValue(strconv.FormatInt(int64(v), 10))
		}
		if v, ok := cur.(int); ok {
			return types.StringValue(strconv.Itoa(v))
		}
		if v, ok := cur.(int64); ok {
			return types.StringValue(strconv.FormatInt(v, 10))
		}
		return types.StringNull()
	}
	m := gadInstanceResourceModel{
		FwTimeouts:             types.ObjectNull(fwTimeoutTypes()),
		Id:                     str("GadInstanceId"),
		CreateTime:             str("CreateTime"),
		DbEngineType:           str("DbEngineType"),
		GadInstanceName:        str("GadInstanceName"),
		GadInstanceType:        str("GadInstanceType"),
		MasterDbInstanceId:     str("MasterDbInstanceId"),
		MasterDbInstanceRegion: str("MasterDbInstanceRegion"),
		RegionId:               str("RegionId"),
		ResourceGroupId:        str("ResourceGroupId"),
		Status:                 str("Status"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.MasterDbZoneId = prior.MasterDbZoneId
	}
	return m
}
