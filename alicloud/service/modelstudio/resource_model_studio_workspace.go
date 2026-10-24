// Package modelstudio. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package modelstudio

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

const workspaceTypeName = "alicloud_model_studio_workspace"

var (
	_ resource.Resource                = &workspaceResource{}
	_ resource.ResourceWithImportState = &workspaceResource{}
)

func NewWorkspaceResource() resource.Resource {
	return &workspaceResource{}
}

type workspaceResource struct {
	fwadapt.ResourceBase
}

type workspaceResourceModel struct {
	Id            types.String `tfsdk:"id"`
	FwTimeouts    types.Object `tfsdk:"timeouts"`
	ApiHost       types.String `tfsdk:"api_host"`
	CreateTime    types.Int64  `tfsdk:"create_time"`
	RegionId      types.String `tfsdk:"region_id"`
	ServiceSite   types.String `tfsdk:"service_site"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
}

func (r *workspaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = workspaceTypeName
}

func (r *workspaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_host": schema.StringAttribute{
				Computed:    true,
				Description: "The API host.",
			},
			"create_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The creation time.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID.",
			},
			"service_site": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "The service deployment region. For more information, [see the documentation](https://www.alibabacloud.com/help/zh/model-studio/regions/).",
			},
			"workspace_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The workspace name.",
			},
		},
	}
}

func (r *workspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workspaceResourceModel
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
		"ClientToken": buildClientToken("CreateWorkspace"),
	}
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.ApiHost.IsNull() && !plan.ApiHost.IsUnknown() {
		createValues["ApiHost"] = fwValueToInterface(plan.ApiHost, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.RegionId.IsNull() && !plan.RegionId.IsUnknown() {
		createValues["RegionId"] = fwValueToInterface(plan.RegionId, map[string]fwWireField{}, false)
	}
	if !plan.ServiceSite.IsNull() && !plan.ServiceSite.IsUnknown() {
		createValues["ServiceSite"] = fwValueToInterface(plan.ServiceSite, map[string]fwWireField{}, false)
	}
	if !plan.WorkspaceName.IsNull() && !plan.WorkspaceName.IsUnknown() {
		createValues["WorkspaceName"] = fwValueToInterface(plan.WorkspaceName, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.ServiceSite", Target: "serviceSite", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.WorkspaceName", Target: "workspaceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, true); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRoaPost(r.Client(), "ModelStudio", "2026-02-10", "/modelstudio/workspaces", request, operationTimeout, map[string]string{"workspaceName": "query", "serviceSite": "query", "__api": "CreateWorkspace"})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateWorkspace", workspaceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateWorkspace", workspaceTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(response, fwMappingPath("$.workspace.workspaceId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_model_studio_workspace: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateWorkspace", workspaceTypeName),
			fmt.Sprintf("The response carried no WorkspaceId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.ApiHost = fwKnownValue(trackedPlan.ApiHost).(types.String)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.Int64)
	trackedPlan.RegionId = fwKnownValue(trackedPlan.RegionId).(types.String)
	trackedPlan.ServiceSite = fwKnownValue(trackedPlan.ServiceSite).(types.String)
	trackedPlan.WorkspaceName = fwKnownValue(trackedPlan.WorkspaceName).(types.String)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := workspaceGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", workspaceTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", workspaceTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", workspaceTypeName),
			fmt.Sprintf("The resource %s is not visible via ListWorkspaces right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, workspaceState(object, &plan))...)
}

func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workspaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := workspaceGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListWorkspaces", workspaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListWorkspaces", workspaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, workspaceState(object, &state))...)
}

func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state workspaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *workspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workspaceResourceModel
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
	if !state.ApiHost.IsNull() && !state.ApiHost.IsUnknown() {
		deleteValues["ApiHost"] = fwValueToInterface(state.ApiHost, map[string]fwWireField{}, false)
	}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		deleteValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.RegionId.IsNull() && !state.RegionId.IsUnknown() {
		deleteValues["RegionId"] = fwValueToInterface(state.RegionId, map[string]fwWireField{}, false)
	}
	if !state.ServiceSite.IsNull() && !state.ServiceSite.IsUnknown() {
		deleteValues["ServiceSite"] = fwValueToInterface(state.ServiceSite, map[string]fwWireField{}, false)
	}
	if !state.WorkspaceName.IsNull() && !state.WorkspaceName.IsUnknown() {
		deleteValues["WorkspaceName"] = fwValueToInterface(state.WorkspaceName, map[string]fwWireField{}, false)
	}
	deleteValues["WorkspaceId"] = state.Id.ValueString()
	needsRefresh := false
	if needsRefresh {
		object, err := workspaceGet(r.Client(), state.Id.ValueString())
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
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteWorkspace")}
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.WorkspaceId", Target: "workspaceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, true); err != nil {
			resp.Diagnostics.AddError("Building DeleteWorkspace request", err.Error())
			return
		}
		if _, err := retryRoaDelete(r.Client(), "ModelStudio", "2026-02-10", "/modelstudio/workspaces/{workspaceId}", request0, operationTimeout, map[string]string{"workspaceId": "path", "__api": "DeleteWorkspace"}); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.Workspace", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_model_studio_workspace: DeleteWorkspace", err.Error())
			return
		}
	}
}

func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// workspaceState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func workspaceState(object map[string]interface{}, prior *workspaceResourceModel) workspaceResourceModel {
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
	m := workspaceResourceModel{
		FwTimeouts:    types.ObjectNull(fwTimeoutTypes()),
		Id:            str("WorkspaceId"),
		ApiHost:       str("ApiHost"),
		CreateTime:    intVal("CreateTime"),
		RegionId:      str("RegionId"),
		ServiceSite:   str("ServiceSite"),
		WorkspaceName: str("WorkspaceName"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
	}
	return m
}
