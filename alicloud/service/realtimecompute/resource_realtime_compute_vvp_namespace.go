// Package realtimecompute. This file is generated automatically. Please do not modify it manually, thank you!
// Object-valued attributes use list of object types to stay compatible with Terraform Plugin Protocol 5.
package realtimecompute

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const vvpNamespaceTypeName = "alicloud_realtime_compute_vvp_namespace"

var (
	_ resource.Resource                = &vvpNamespaceResource{}
	_ resource.ResourceWithImportState = &vvpNamespaceResource{}
)

func NewVvpNamespaceResource() resource.Resource {
	return &vvpNamespaceResource{}
}

type vvpNamespaceResource struct {
	fwadapt.ResourceBase
}

type vvpNamespaceResourceModel struct {
	Id           types.String `tfsdk:"id"`
	FwTimeouts   types.Object `tfsdk:"timeouts"`
	GmtCreate    types.Int64  `tfsdk:"gmt_create"`
	GmtModified  types.Int64  `tfsdk:"gmt_modified"`
	Ha           types.Bool   `tfsdk:"ha"`
	InstanceId   types.String `tfsdk:"instance_id"`
	Namespace    types.String `tfsdk:"namespace"`
	ResourceSpec types.List   `tfsdk:"resource_spec"`
	ResourceUsed types.List   `tfsdk:"resource_used"`
	Status       types.String `tfsdk:"status"`
	Tags         types.Map    `tfsdk:"tags"`
}

func (r *vvpNamespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = vvpNamespaceTypeName
}

func (r *vvpNamespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_create": schema.Int64Attribute{
				Computed:    true,
				Description: "The creation time",
			},
			"gmt_modified": schema.Int64Attribute{
				Computed:    true,
				Description: "The modification time",
			},
			"ha": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "High availability, set at creation, immutable",
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The instance id of the parent vvp instance",
			},
			"namespace": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The namespace name; the workspace default namespace is named <workspace_name>-default (or <workspace_name>-default-ha for a multi-availability-zone workspace)",
			},
			"resource_spec": schema.ListAttribute{
				ElementType: vvpNamespaceResourceSpecObjType(),
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(1),
					vvpNamespaceResourceSpecValidator{},
				},
				Optional:    true,
				Computed:    true,
				Description: "Resource specification (CU) of the namespace",
			},
			"resource_used": schema.ListAttribute{
				ElementType: vvpNamespaceResourceUsedObjType(),
				Computed:    true,
				Description: "Resource usage of the namespace",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the namespace",
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Tags of the namespace",
			},
		},
	}
}

// vvpNamespaceResourceSpecValidator enforces the required cpu/memory_gb pair inside every
// resource_spec element. Protocol 5 cannot carry nested attributes, so the per-field
// required-ness declared by the CloudSpec model is checked here instead of in the schema.
type vvpNamespaceResourceSpecValidator struct{}

func (vvpNamespaceResourceSpecValidator) Description(context.Context) string {
	return "resource_spec requires cpu and memory_gb in every element"
}

func (vvpNamespaceResourceSpecValidator) MarkdownDescription(context.Context) string {
	return "resource_spec requires cpu and memory_gb in every element"
}

func (vvpNamespaceResourceSpecValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for i, element := range req.ConfigValue.Elements() {
		object, ok := element.(types.Object)
		if !ok || object.IsNull() || object.IsUnknown() {
			continue
		}
		attributes := object.Attributes()
		for _, name := range []string{"cpu", "memory_gb"} {
			if value, present := attributes[name]; !present || value.IsNull() {
				resp.Diagnostics.AddAttributeError(
					req.Path.AtListIndex(i).AtName(name),
					"Missing required value",
					fmt.Sprintf("resource_spec[%d].%s is required.", i, name),
				)
			}
		}
		cpu, okCPU := attributes["cpu"].(types.Int64)
		memory, okMemory := attributes["memory_gb"].(types.Int64)
		if !okCPU || !okMemory || cpu.IsNull() || cpu.IsUnknown() || memory.IsNull() || memory.IsUnknown() {
			continue
		}
		if reason := vvpNamespaceResourceSpecRatioError(cpu.ValueInt64(), memory.ValueInt64()); reason != "" {
			resp.Diagnostics.AddAttributeError(req.Path.AtListIndex(i).AtName("memory_gb"), "Invalid resource specification", reason)
		}
	}
}

// vvpNamespaceResourceSpecRatioError returns a non-empty reason when the pair violates the product
// contract documented on foasconsole ResourceSpec: memory must be exactly four times the CPU count.
// Unknown values are filtered by the caller, so this only judges pairs that are known at plan time.
func vvpNamespaceResourceSpecRatioError(cpu, memoryGb int64) string {
	if cpu <= 0 || memoryGb <= 0 {
		return fmt.Sprintf("resource_spec requires positive cpu and memory_gb, got cpu=%d memory_gb=%d", cpu, memoryGb)
	}
	if memoryGb != 4*cpu {
		return fmt.Sprintf("resource_spec.memory_gb must be 4 times resource_spec.cpu, got cpu=%d memory_gb=%d", cpu, memoryGb)
	}
	return ""
}

func (r *vvpNamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vvpNamespaceResourceModel
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
		"ClientToken": buildClientToken("CreateNamespace"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.Ha.IsNull() && !plan.Ha.IsUnknown() {
		createValues["Ha"] = fwValueToInterface(plan.Ha, map[string]fwWireField{}, false)
	}
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		createValues["InstanceId"] = fwValueToInterface(plan.InstanceId, map[string]fwWireField{}, false)
	}
	if !plan.Namespace.IsNull() && !plan.Namespace.IsUnknown() {
		createValues["Namespace"] = fwValueToInterface(plan.Namespace, map[string]fwWireField{}, false)
	}
	if !plan.ResourceSpec.IsNull() && !plan.ResourceSpec.IsUnknown() {
		createValues["ResourceSpec"] = fwValueToInterface(plan.ResourceSpec, map[string]fwWireField{"cpu": {Name: "Cpu", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "memory_gb": {Name: "MemoryGb", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.Ha", Target: "CreateNamespaceRequest.Ha", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.InstanceId", Target: "CreateNamespaceRequest.InstanceId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Namespace", Target: "CreateNamespaceRequest.Namespace", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "CreateNamespaceRequest.Region", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.ResourceSpec.Cpu", Target: "CreateNamespaceRequest.ResourceSpec.Cpu", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.ResourceSpec.MemoryGb", Target: "CreateNamespaceRequest.ResourceSpec.MemoryGB", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	if _, err := retryRpcPost(ctx, r.Client(), "foasconsole", "2019-06-01", "CreateNamespace", request, operationTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNamespace", vvpNamespaceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateNamespace", vvpNamespaceTypeName).Error(),
		)
		return
	}

	// The request is flattened to dotted keys ("CreateNamespaceRequest.InstanceId") before it is
	// sent, so walking it for the nested path always comes back empty. The composite identity is
	// built from the required plan fields that were mapped into that request just above.
	identityParts := []string{}
	instanceId := strings.TrimSpace(plan.InstanceId.ValueString())
	if instanceId == "" {
		resp.Diagnostics.AddError("Creating alicloud_realtime_compute_vvp_namespace: missing identity", "The create request did not carry the instance_id component")
		return
	}
	identityParts = append(identityParts, instanceId)
	namespaceName := strings.TrimSpace(plan.Namespace.ValueString())
	if namespaceName == "" {
		resp.Diagnostics.AddError("Creating alicloud_realtime_compute_vvp_namespace: missing identity", "The create request did not carry the namespace component")
		return
	}
	identityParts = append(identityParts, namespaceName)
	id := strings.Join(identityParts, ":")

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.GmtCreate = fwKnownValue(trackedPlan.GmtCreate).(types.Int64)
	trackedPlan.GmtModified = fwKnownValue(trackedPlan.GmtModified).(types.Int64)
	trackedPlan.Ha = fwKnownValue(trackedPlan.Ha).(types.Bool)
	trackedPlan.InstanceId = fwKnownValue(trackedPlan.InstanceId).(types.String)
	trackedPlan.Namespace = fwKnownValue(trackedPlan.Namespace).(types.String)
	trackedPlan.ResourceSpec = fwKnownValue(trackedPlan.ResourceSpec).(types.List)
	trackedPlan.ResourceUsed = fwKnownValue(trackedPlan.ResourceUsed).(types.List)
	trackedPlan.Status = fwKnownValue(trackedPlan.Status).(types.String)
	trackedPlan.Tags = fwKnownValue(trackedPlan.Tags).(types.Map)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The model's create async detection is assertNotEmpty on $.Namespace, so poll until
	// DescribeNamespaces returns the namespace instead of waiting for a status string.
	object, err := vvpNamespaceWaitForState(ctx, r.Client(), id, "#$.Namespace", []string{"#CHECKSET"}, nil, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for Running", vvpNamespaceTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for Running (%s)", vvpNamespaceTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", vvpNamespaceTypeName),
			fmt.Sprintf("The resource %s is not visible via DescribeNamespaces right after creation", id),
		)
		return
	}
	// The namespace is listed as soon as the create request is accepted, while its capacity
	// allocation still runs in the background; any follow-up call fails with 903047 until that
	// allocation settles, so wait for SUCCESS before reporting the create as done.
	if object, err = vvpNamespaceWaitForState(ctx, r.Client(), id, "Status", []string{"SUCCESS"}, []string{"FAILED"}, operationTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for SUCCESS", vvpNamespaceTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for SUCCESS (%s)", vvpNamespaceTypeName, id).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vvpNamespaceState(object, &plan))...)
}

func (r *vvpNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vvpNamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := vvpNamespaceGet(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNamespaces", vvpNamespaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNamespaces", vvpNamespaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vvpNamespaceState(object, &state))...)
}

func (r *vvpNamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state vvpNamespaceResourceModel
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
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		planValues["InstanceId"] = fwValueToInterface(plan.InstanceId, map[string]fwWireField{}, false)
	}
	if !plan.Namespace.IsNull() && !plan.Namespace.IsUnknown() {
		planValues["Namespace"] = fwValueToInterface(plan.Namespace, map[string]fwWireField{}, false)
	}
	if !plan.ResourceSpec.IsNull() && !plan.ResourceSpec.IsUnknown() {
		planValues["ResourceSpec"] = fwValueToInterface(plan.ResourceSpec, map[string]fwWireField{"cpu": {Name: "Cpu", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "memory_gb": {Name: "MemoryGb", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		stateValues["InstanceId"] = fwValueToInterface(state.InstanceId, map[string]fwWireField{}, false)
	}
	if !state.Namespace.IsNull() && !state.Namespace.IsUnknown() {
		stateValues["Namespace"] = fwValueToInterface(state.Namespace, map[string]fwWireField{}, false)
	}
	if !state.ResourceSpec.IsNull() && !state.ResourceSpec.IsUnknown() {
		stateValues["ResourceSpec"] = fwValueToInterface(state.ResourceSpec, map[string]fwWireField{"cpu": {Name: "Cpu", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "memory_gb": {Name: "MemoryGb", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	updated0 := false
	if fwPlannedChanged(plan.ResourceSpec, state.ResourceSpec) {
		updated0 = true
	}
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{"ClientToken": buildClientToken("ModifyNamespaceSpec")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.InstanceId", Target: "ModifyNamespaceSpecRequest.InstanceId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.Namespace", Target: "ModifyNamespaceSpecRequest.Namespace", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "ModifyNamespaceSpecRequest.Region", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.ResourceSpec.Cpu", Target: "ModifyNamespaceSpecRequest.ResourceSpec.Cpu", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.ResourceSpec.MemoryGb", Target: "ModifyNamespaceSpecRequest.ResourceSpec.MemoryGB", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		}, operationValues0, stateValues, []string{"$.InstanceId", "$.Namespace", "$.RegionId", "$.ResourceSpec.Cpu", "$.ResourceSpec.MemoryGb"}), false); err != nil {
			resp.Diagnostics.AddError("Building ModifyNamespaceSpec request", err.Error())
			return
		}
		if _, err := retryRpcPost(ctx, r.Client(), "foasconsole", "2019-06-01", "ModifyNamespaceSpec", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_realtime_compute_vvp_namespace: ModifyNamespaceSpec", err.Error())
			return
		}
		// ModifyNamespaceSpec returns before the new capacity is in place; wait for SUCCESS so the
		// state written below and the next operation see a settled namespace.
		if _, err := vvpNamespaceWaitForState(ctx, r.Client(), state.Id.ValueString(), "Status", []string{"SUCCESS"}, []string{"FAILED"}, operationTimeout); err != nil {
			resp.Diagnostics.AddError(
				"Updating alicloud_realtime_compute_vvp_namespace: waiting for SUCCESS",
				errs.WrapErrorf(err, "Updating %s: waiting for SUCCESS (%s)", vvpNamespaceTypeName, state.Id.ValueString()).Error(),
			)
			return
		}
	}
	object, err := vvpNamespaceGet(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNamespaces", vvpNamespaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNamespaces", vvpNamespaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vvpNamespaceState(object, &plan))...)
}

func (r *vvpNamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vvpNamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A namespace that is already gone — for example because its workspace was deleted — needs no
	// delete call; returning lets the framework drop the state entry.
	if state.Id.ValueString() != "" {
		existing, err := vvpNamespaceGet(ctx, r.Client(), state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Deleting alicloud_realtime_compute_vvp_namespace: DescribeNamespaces", err.Error())
			return
		}
		if existing == nil {
			return
		}
	}

	deleteValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		deleteValues["InstanceId"] = fwValueToInterface(state.InstanceId, map[string]fwWireField{}, false)
	}
	if !state.Namespace.IsNull() && !state.Namespace.IsUnknown() {
		deleteValues["Namespace"] = fwValueToInterface(state.Namespace, map[string]fwWireField{}, false)
	}
	identityParts := strings.Split(state.Id.ValueString(), ":")
	if len(identityParts) > 0 {
		deleteValues["InstanceId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		deleteValues["Namespace"] = identityParts[1]
	}
	if state.InstanceId.ValueString() == "" || state.Namespace.ValueString() == "" {
		resp.Diagnostics.AddError("Deleting alicloud_realtime_compute_vvp_namespace", "The namespace identity is incomplete; the resource remains tracked")
		return
	}
	request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteNamespace")}
	request0["RegionId"] = r.Client().RegionId
	if err := fwApplyMappings(request0, deleteValues, []fwMapping{
		{Source: "$.InstanceId", Target: "DeleteNamespaceRequest.InstanceId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Namespace", Target: "DeleteNamespaceRequest.Namespace", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "DeleteNamespaceRequest.Region", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building DeleteNamespace request", err.Error())
		return
	}
	if _, err := retryRpcPost(ctx, r.Client(), "foasconsole", "2019-06-01", "DeleteNamespace", request0, operationTimeout); err != nil {
		if isExpectedErrors(err, []string{"ResourceNotFound.VvpNamespace", "ResourceNotFound"}) {
			// Already gone; fall through to the gone-wait so the state entry is dropped.
		} else if isExpectedErrors(err, []string{"903054"}) {
			// 903054: the namespace is the last valid one in its workspace, so the product refuses
			// to delete it on its own and the object still exists. Dropping it from state here would
			// make Terraform forget a live resource and hide the drift, so keep the state and hand
			// the decision back to the user.
			resp.Diagnostics.AddError(
				fmt.Sprintf("Deleting %s: the workspace keeps its last namespace", vvpNamespaceTypeName),
				fmt.Sprintf("%s is the last valid namespace of workspace %s, and Realtime Compute refuses to delete it on its own (903054). The namespace still exists and stays in Terraform state. Delete the parent workspace instead, or remove this resource from state with `terraform state rm` once you have confirmed the namespace should disappear with the workspace.",
					state.Namespace.ValueString(), state.InstanceId.ValueString()),
			)
			return
		} else {
			resp.Diagnostics.AddError("Deleting alicloud_realtime_compute_vvp_namespace: DeleteNamespace", err.Error())
			return
		}
	}
	// gap: v1 Delete 用 BuildStateConf(target=[]) 轮询 StateRefreshFunc 直到 Get 返回 nil（gone），
	// 解决 gateway-route-table 异步删除后 DeleteVpc 竞态 DependencyViolation.GatewayRouteTable。
	// v2 WaitForState：空 deleteTargetStatus（gone-wait）发 Target=[]string{}（空 slice, len 0），
	// 镜像 v1 BuildStateConf(target=[])。hashicorp retry 当 Refresh 返回 (nil,"",nil)（gone）且
	// Target 空 → success（done）；若误用 Target=[""]（len 1, 非空），gone 时 helper 报
	// NotFoundError("couldn't find resource (N retries)") 而非 success——此前 route_table destroy
	// 即栽此。非空 deleteTargetStatus 则发 Target=["<status>"] 等待该终态。Get 在 not-found 返回
	// (nil,nil)（非 error），故 wait 不挂；同步删除资源首检即 gone（nil,"",nil→Target 空→done），无 hang。
	if _, err := vvpNamespaceWaitForState(ctx, r.Client(), state.Id.ValueString(), "Namespace", []string{}, nil, operationTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for gone", vvpNamespaceTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for gone (%s)", vvpNamespaceTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
}

func (r *vvpNamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// vvpNamespaceState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func vvpNamespaceState(object map[string]interface{}, prior *vvpNamespaceResourceModel) vvpNamespaceResourceModel {
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
	boolVal := func(key string) types.Bool {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.BoolNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.BoolNull()
			}
		}
		if v, ok := cur.(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolNull()
	}
	intVal := func(key string) types.Int64 {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.Int64Null()
			}
			cur, ok = m[seg]
			if !ok {
				return types.Int64Null()
			}
		}
		switch v := cur.(type) {
		case float64:
			return types.Int64Value(int64(v))
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	// lookup traverses a dotted responseKey (e.g. "RegionStatuses.RegionStatus") segment by segment
	// from the Get response object, returning the raw value (array/map/scalar) or nil. gap#5: array
	// fields' responseKey is the array's true path from the read-tree [*] child, which may be a
	// two-level unwrap (object["RegionStatuses"]["RegionStatus"]) that a literal object[key] lookup
	// cannot reach. Single-segment keys (e.g. "IpWhiteList") traverse one step == object[key].
	lookup := func(key string) interface{} {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return nil
			}
			cur, ok = m[seg]
			if !ok {
				return nil
			}
		}
		return cur
	}
	// gap#11: composite-ID resource — build the id from all composite fields in the response,
	// joined by ":" (e.g. AccessGroupName:FileSystemType). The Get response carries every
	// part, so the state id round-trips through Read/Delete correctly (otherwise Read would
	// lose FileSystemType and DescribeAccessGroups would query the wrong namespace → 404).
	_cidParts := []string{
		str("InstanceId").ValueString(),
		str("Namespace").ValueString(),
	}
	m := vvpNamespaceResourceModel{
		FwTimeouts:   types.ObjectNull(fwTimeoutTypes()),
		Id:           types.StringValue(strings.Join(_cidParts, ":")),
		GmtCreate:    intVal("GmtCreate"),
		GmtModified:  intVal("GmtModified"),
		Ha:           boolVal("Ha"),
		InstanceId:   str("InstanceId"),
		Namespace:    str("Namespace"),
		ResourceSpec: vvpNamespaceResourceSpecToTypeList(lookup("ResourceSpec")),
		ResourceUsed: vvpNamespaceResourceUsedToTypeList(lookup("ResourceUsed")),
		Status:       str("Status"),
		Tags:         tagsToTypeMap(lookup("Tags")),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.ResourceSpec = fwPreserveRequestValues(m.ResourceSpec, prior.ResourceSpec, map[string]fwWireField{"cpu": {Name: "Cpu", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "memory_gb": {Name: "MemoryGb", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}).(types.List)
		m.ResourceUsed = fwPreserveRequestValues(m.ResourceUsed, prior.ResourceUsed, map[string]fwWireField{}).(types.List)
	}
	return m
}
