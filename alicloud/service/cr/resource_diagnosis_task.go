// Package cr. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package cr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const diagnosisTaskTypeName = "alicloud_cr_diagnosis_task"

var (
	_ resource.Resource                = &diagnosisTaskResource{}
	_ resource.ResourceWithImportState = &diagnosisTaskResource{}
)

func NewDiagnosisTaskResource() resource.Resource {
	return &diagnosisTaskResource{}
}

type diagnosisTaskResource struct {
	fwadapt.ResourceBase
}

type diagnosisTaskResourceModel struct {
	Id              types.String `tfsdk:"id"`
	FwTimeouts      types.Object `tfsdk:"timeouts"`
	DiagnosisTaskId types.String `tfsdk:"diagnosis_task_id"`
	DiagnosisTime   types.String `tfsdk:"diagnosis_time"`
	DiagnosisType   types.String `tfsdk:"diagnosis_type"`
	InstanceId      types.String `tfsdk:"instance_id"`
	Status          types.String `tfsdk:"status"`
	Targets         types.List   `tfsdk:"targets"`
}

func (r *diagnosisTaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = diagnosisTaskTypeName
}

func (r *diagnosisTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": fwTimeoutBlock(),
			// targets uses a nested block rather than a nested attribute because
			// the muxed provider serves Protocol v5 (tf5muxserver +
			// NewProtocol5); nested attributes require Protocol v6.
			"targets": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"repository": schema.StringAttribute{
							Optional:    true,
							Description: "Repository name.",
						},
						"end_time": schema.StringAttribute{
							Optional:    true,
							Description: "End time of the diagnosis period.",
						},
						"start_time": schema.StringAttribute{
							Optional:    true,
							Description: "Start time of the diagnosis period.",
						},
						"tag": schema.StringAttribute{
							Optional:    true,
							Description: "Image tag.",
						},
						"related_id": schema.StringAttribute{
							Optional:    true,
							Description: "The ID of the diagnosed item.  ",
						},
						"namespace": schema.StringAttribute{
							Optional:    true,
							Description: "The namespace name.  ",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				Description: "Diagnosis targets.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"diagnosis_task_id": schema.StringAttribute{
				Computed:    true,
				Description: "A resource property field that represents the primary resource ID.  ",
			},
			"diagnosis_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the diagnosis task was triggered.  ",
			},
			"diagnosis_type": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.OneOf("BUILD_DIAGNOSIS", "SYNC_DIAGNOSIS", "INSTANCE_DIAGNOSIS", "TAG_LIFECYCLE_ANALYSIS", "PULL_PUSH_DIAGNOSIS"),
				},
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Diagnosis type.  ",
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Instance ID.  ",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "A resource property field that represents the resource status.  ",
			},
		},
	}
}

func (r *diagnosisTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan diagnosisTaskResourceModel
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
		"ClientToken": buildClientToken("CreateDiagnosisTask"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Meta().RegionId
	createValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
	if !plan.DiagnosisTaskId.IsNull() && !plan.DiagnosisTaskId.IsUnknown() {
		createValues["DiagnosisTaskId"] = fwValueToInterface(plan.DiagnosisTaskId, map[string]fwWireField{}, false)
	}
	if !plan.DiagnosisTime.IsNull() && !plan.DiagnosisTime.IsUnknown() {
		createValues["DiagnosisTime"] = fwValueToInterface(plan.DiagnosisTime, map[string]fwWireField{}, false)
	}
	if !plan.DiagnosisType.IsNull() && !plan.DiagnosisType.IsUnknown() {
		createValues["DiagnosisType"] = fwValueToInterface(plan.DiagnosisType, map[string]fwWireField{}, false)
	}
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		createValues["InstanceId"] = fwValueToInterface(plan.InstanceId, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		createValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	if !plan.Targets.IsNull() && !plan.Targets.IsUnknown() {
		createValues["Targets"] = fwValueToInterface(plan.Targets, map[string]fwWireField{"repository": {Name: "Repository", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "end_time": {Name: "EndTime", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "start_time": {Name: "StartTime", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "tag": {Name: "Tag", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "related_id": {Name: "RelatedId", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "namespace": {Name: "Namespace", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.Targets.Namespace", Target: "DiagnosisTarget.Namespace", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DiagnosisType", Target: "DiagnosisType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Targets.Extra", Target: "DiagnosisTarget.Extra", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Targets.Repository", Target: "DiagnosisTarget.Repository", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Targets.EndTime", Target: "DiagnosisTarget.EndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Targets.RelatedId", Target: "DiagnosisTarget.RelatedId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Targets.StartTime", Target: "DiagnosisTarget.StartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Targets.Tag", Target: "DiagnosisTarget.Tag", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Meta(), "cr", "2018-12-01", "CreateDiagnosisTask", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDiagnosisTask", diagnosisTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDiagnosisTask", diagnosisTaskTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(request, fwMappingPath("InstanceId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_cr_diagnosis_task: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	identityValues1 := fwWalkMappedPath(response, fwMappingPath("$.DiagnosisTaskId"), nil)
	if len(identityValues1) == 0 || identityValues1[0].Value == nil || fmt.Sprint(identityValues1[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_cr_diagnosis_task: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues1[0].Value))
	// gap#7: GetDiagnosisTask requires DiagnosisType as a request parameter even though it is
	// not part of the model's uniqueKeyFields (InstanceId + DiagnosisTaskId). Encoding it as the
	// third composite-id segment lets Read/Update/Import reconstruct the full lookup without a
	// separate channel; DiagnosisType is an immutable enum (ForceNew), so it is stable in the id.
	if plan.DiagnosisType.IsNull() || plan.DiagnosisType.IsUnknown() || plan.DiagnosisType.ValueString() == "" {
		resp.Diagnostics.AddError("Creating alicloud_cr_diagnosis_task: missing identity", "The create operation did not accept a DiagnosisType identity component")
		return
	}
	identityParts = append(identityParts, plan.DiagnosisType.ValueString())
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDiagnosisTask", diagnosisTaskTypeName),
			fmt.Sprintf("The response carried no InstanceId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.DiagnosisTaskId = fwKnownValue(trackedPlan.DiagnosisTaskId).(types.String)
	trackedPlan.DiagnosisTime = fwKnownValue(trackedPlan.DiagnosisTime).(types.String)
	trackedPlan.DiagnosisType = fwKnownValue(trackedPlan.DiagnosisType).(types.String)
	trackedPlan.InstanceId = fwKnownValue(trackedPlan.InstanceId).(types.String)
	trackedPlan.Status = fwKnownValue(trackedPlan.Status).(types.String)
	trackedPlan.Targets = fwKnownValue(trackedPlan.Targets).(types.List)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := diagnosisTaskGet(r.Meta(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", diagnosisTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", diagnosisTaskTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", diagnosisTaskTypeName),
			fmt.Sprintf("The resource %s is not visible via GetDiagnosisTask right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, diagnosisTaskState(object, &plan))...)
}

func (r *diagnosisTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state diagnosisTaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := diagnosisTaskGet(r.Meta(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDiagnosisTask", diagnosisTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDiagnosisTask", diagnosisTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, diagnosisTaskState(object, &state))...)
}

func (r *diagnosisTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state diagnosisTaskResourceModel
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

	planValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
	if !plan.DiagnosisTaskId.IsNull() && !plan.DiagnosisTaskId.IsUnknown() {
		planValues["DiagnosisTaskId"] = fwValueToInterface(plan.DiagnosisTaskId, map[string]fwWireField{}, false)
	}
	if !plan.DiagnosisTime.IsNull() && !plan.DiagnosisTime.IsUnknown() {
		planValues["DiagnosisTime"] = fwValueToInterface(plan.DiagnosisTime, map[string]fwWireField{}, false)
	}
	if !plan.DiagnosisType.IsNull() && !plan.DiagnosisType.IsUnknown() {
		planValues["DiagnosisType"] = fwValueToInterface(plan.DiagnosisType, map[string]fwWireField{}, false)
	}
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		planValues["InstanceId"] = fwValueToInterface(plan.InstanceId, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		planValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	if !plan.Targets.IsNull() && !plan.Targets.IsUnknown() {
		planValues["Targets"] = fwValueToInterface(plan.Targets, map[string]fwWireField{"repository": {Name: "Repository", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "end_time": {Name: "EndTime", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "start_time": {Name: "StartTime", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "tag": {Name: "Tag", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "related_id": {Name: "RelatedId", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "namespace": {Name: "Namespace", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	stateValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
	if !state.DiagnosisTaskId.IsNull() && !state.DiagnosisTaskId.IsUnknown() {
		stateValues["DiagnosisTaskId"] = fwValueToInterface(state.DiagnosisTaskId, map[string]fwWireField{}, false)
	}
	if !state.DiagnosisTime.IsNull() && !state.DiagnosisTime.IsUnknown() {
		stateValues["DiagnosisTime"] = fwValueToInterface(state.DiagnosisTime, map[string]fwWireField{}, false)
	}
	if !state.DiagnosisType.IsNull() && !state.DiagnosisType.IsUnknown() {
		stateValues["DiagnosisType"] = fwValueToInterface(state.DiagnosisType, map[string]fwWireField{}, false)
	}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		stateValues["InstanceId"] = fwValueToInterface(state.InstanceId, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		stateValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	if !state.Targets.IsNull() && !state.Targets.IsUnknown() {
		stateValues["Targets"] = fwValueToInterface(state.Targets, map[string]fwWireField{"repository": {Name: "Repository", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "end_time": {Name: "EndTime", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "start_time": {Name: "StartTime", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "tag": {Name: "Tag", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "related_id": {Name: "RelatedId", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "namespace": {Name: "Namespace", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	updated0 := false
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"InstanceId":  state.Id.ValueString(),
			"ClientToken": buildClientToken("CancelDiagnosisTask"),
		}
		parts0 := strings.Split(state.Id.ValueString(), ":")
		if len(parts0) > 0 {
			request0["InstanceId"] = parts0[0]
		}
		if len(parts0) > 1 {
			request0["DiagnosisTaskId"] = parts0[1]
		}
		request0["RegionId"] = r.Meta().RegionId
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DiagnosisTaskId", Target: "DiagnosisTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues0, stateValues, []string{"$.InstanceId", "$.DiagnosisTaskId", "$.RegionId"}), false); err != nil {
			resp.Diagnostics.AddError("Building CancelDiagnosisTask request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Meta(), "cr", "2018-12-01", "CancelDiagnosisTask", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_cr_diagnosis_task: CancelDiagnosisTask", err.Error())
			return
		}
	}
	object, err := diagnosisTaskGet(r.Meta(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDiagnosisTask", diagnosisTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDiagnosisTask", diagnosisTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, diagnosisTaskState(object, &plan))...)
}

func (r *diagnosisTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Cannot destroy resource alicloud_cr_diagnosis_task. The API declares no
	// deletion operation (there is no DeleteDiagnosisTask). Terraform removes
	// the resource from the state file; the remote diagnosis task may remain.
	// This mirrors the v1 no-delete pattern (e.g. alicloud_aligreen_oss_stock_task).
}

func (r *diagnosisTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// diagnosisTaskState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func diagnosisTaskState(object map[string]interface{}, prior *diagnosisTaskResourceModel) diagnosisTaskResourceModel {
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
	// gap#11: composite-ID resource — build the id from all composite fields in the response,
	// joined by ":". DiagnosisType is not part of the model's uniqueKeyFields but GetDiagnosisTask
	// requires it as a request parameter, so it is carried as the third id segment; the Get
	// response (or completeRead) carries it, keeping the state id round-tripping through
	// Read/Update/Import (otherwise Read would lose DiagnosisType → GetDiagnosisTask 404).
	_cidParts := []string{
		str("InstanceId").ValueString(),
		str("DiagnosisTaskId").ValueString(),
		str("DiagnosisType").ValueString(),
	}
	m := diagnosisTaskResourceModel{
		FwTimeouts:      types.ObjectNull(fwTimeoutTypes()),
		Id:              types.StringValue(strings.Join(_cidParts, ":")),
		DiagnosisTaskId: str("DiagnosisTaskId"),
		DiagnosisTime:   str("DiagnosisTime"),
		DiagnosisType:   str("DiagnosisType"),
		InstanceId:      str("InstanceId"),
		Status:          str("Status"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.Targets = prior.Targets
	}
	return m
}
