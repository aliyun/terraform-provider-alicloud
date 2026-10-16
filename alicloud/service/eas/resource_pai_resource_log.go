package eas

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const resourceLogTypeName = "alicloud_pai_resource_log"

var (
	_ resource.Resource                = &paiResourceLogResource{}
	_ resource.ResourceWithImportState = &paiResourceLogResource{}
)

// NewPaiResourceLogResource is the factory registered for alicloud_pai_resource_log.
// Returning a pointer is required: the fwadapt registration guard rejects
// value-returning factories because their embedded Configure never runs, so
// Client() would stay nil forever.
func NewPaiResourceLogResource() resource.Resource {
	return &paiResourceLogResource{}
}

type paiResourceLogResource struct {
	fwadapt.ResourceBase
}

// paiResourceLogModel is the tfsdk model. The SLS log-delivery config is a
// singleton on each EAS Resource, so the resource identity is the composite
// <cluster_id>:<resource_id> of the parent resource. There is no Update API,
// so every configurable attribute forces replacement; the resource therefore
// does not implement ResourceWithUpdate.
type paiResourceLogModel struct {
	Id          types.String `tfsdk:"id"`
	ClusterId   types.String `tfsdk:"cluster_id"`
	ResourceId  types.String `tfsdk:"resource_id"`
	ProjectName types.String `tfsdk:"project_name"`
	LogStore    types.String `tfsdk:"log_store"`
	Status      types.String `tfsdk:"status"`
	Message     types.String `tfsdk:"message"`
}

func (r *paiResourceLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = resourceLogTypeName
}

func (r *paiResourceLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI-EAS Resource SLS log delivery resource. " +
			"It enables/disables delivering the EAS resource group logs to a Log Service (SLS) " +
			"project and log store. The log delivery config is a singleton per EAS resource, " +
			"located by the parent resource's cluster_id and resource_id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, formatted as `<cluster_id>:<resource_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The cluster id of the parent EAS resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The resource id of the parent EAS resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_name": schema.StringAttribute{
				Required: true,
				Description: "The SLS project name that the resource group logs are delivered to. " +
					"The project must already exist in Log Service; this resource does not create it.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"log_store": schema.StringAttribute{
				Required:    true,
				Description: "The SLS log store that the resource group logs are delivered to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message": schema.StringAttribute{
				Computed: true,
				Description: "The SLS log delivery message returned by DescribeResourceLog. " +
					"Empty when there is no extra message.",
			},
		},
	}
}

func (r *paiResourceLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan paiResourceLogModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	action := resourceLogAction(plan.ClusterId.ValueString(), plan.ResourceId.ValueString())
	body := map[string]interface{}{
		"ProjectName": plan.ProjectName.ValueString(),
		"LogStore":    plan.LogStore.ValueString(),
	}
	request := map[string]interface{}{"body": body}

	if _, err := client.RoaPost("EAS", "2021-07-01", action, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResourceLog", resourceLogTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateResourceLog at %s", resourceLogTypeName, action).Error(),
		)
		return
	}
	logs.AddDebug(action, nil, request)

	plan.Id = types.StringValue(resourceLogID(plan.ClusterId.ValueString(), plan.ResourceId.ValueString()))

	object, err := describeResourceLog(client, plan.ClusterId.ValueString(), plan.ResourceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling DescribeResourceLog", resourceLogTypeName),
			errs.WrapErrorf(err, "Creating %s: calling DescribeResourceLog at %s", resourceLogTypeName, action).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResourceLog", resourceLogTypeName),
			fmt.Sprintf("The resource log config is not visible via DescribeResourceLog right after creation (cluster_id=%s, resource_id=%s)",
				plan.ClusterId.ValueString(), plan.ResourceId.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resourceLogState(object, &plan))...)
}

func (r *paiResourceLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state paiResourceLogModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, resourceId := splitResourceLogID(state.Id.ValueString(), state.ClusterId.ValueString(), state.ResourceId.ValueString())

	object, err := describeResourceLog(r.Client(), clusterId, resourceId)
	if err != nil {
		if isResourceLogNotExist(err) {
			resp.Diagnostics.Append(resourceLogGoneDiag(resourceLogTypeName, clusterId, resourceId, err)...)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeResourceLog", resourceLogTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeResourceLog (cluster_id=%s, resource_id=%s)", resourceLogTypeName, clusterId, resourceId).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resourceLogState(object, &state))...)
}

// Update refreshes state from DescribeResourceLog. There is no Update API, and
// every configurable attribute forces replacement, so Terraform never routes a
// change through Update; it replaces (Delete + Create) instead. The method
// exists only to satisfy the resource.Resource interface, which this provider's
// vendored terraform-plugin-framework requires on every resource. Defensively,
// it keeps the computed status/message fresh if it is ever reached.
func (r *paiResourceLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan paiResourceLogModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, resourceId := splitResourceLogID(plan.Id.ValueString(), plan.ClusterId.ValueString(), plan.ResourceId.ValueString())
	object, err := describeResourceLog(r.Client(), clusterId, resourceId)
	if err != nil {
		if isResourceLogNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling DescribeResourceLog", resourceLogTypeName),
			errs.WrapErrorf(err, "Updating %s: calling DescribeResourceLog (cluster_id=%s, resource_id=%s)", resourceLogTypeName, clusterId, resourceId).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resourceLogState(object, &plan))...)
}

func (r *paiResourceLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state paiResourceLogModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, resourceId := splitResourceLogID(state.Id.ValueString(), state.ClusterId.ValueString(), state.ResourceId.ValueString())
	client := r.Client()
	action := resourceLogAction(clusterId, resourceId)

	_, err := client.RoaDelete("EAS", "2021-07-01", action, nil, nil, nil, true)
	if err != nil && !isResourceLogNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteResourceLog", resourceLogTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteResourceLog at %s", resourceLogTypeName, action).Error(),
		)
		return
	}
	logs.AddDebug(action, nil, nil)
}

func (r *paiResourceLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	clusterId, resourceId, ok := parseResourceLogImportID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s", resourceLogTypeName),
			fmt.Sprintf("Unexpected import ID format %q: expected `<cluster_id>:<resource_id>`", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, paiResourceLogModel{
		Id:         types.StringValue(resourceLogID(clusterId, resourceId)),
		ClusterId:  types.StringValue(clusterId),
		ResourceId: types.StringValue(resourceId),
	})...)
}

// describeResourceLog calls DescribeResourceLog and returns the response body,
// or nil when the log delivery config is not enabled (the resource is gone).
func describeResourceLog(client *connectivity.AliyunClient, clusterId, resourceId string) (map[string]interface{}, error) {
	action := resourceLogAction(clusterId, resourceId)
	response, err := client.RoaGet("EAS", "2021-07-01", action, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// resourceLogState maps a DescribeResourceLog response onto the model,
// preserving the configured cluster_id/resource_id from prior state and
// refreshing project_name, log_store, status and message from the API.
func resourceLogState(object map[string]interface{}, prior *paiResourceLogModel) paiResourceLogModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m := paiResourceLogModel{
		Status:  str("Status"),
		Message: str("Message"),
	}
	if prior != nil {
		m.Id = prior.Id
		m.ClusterId = prior.ClusterId
		m.ResourceId = prior.ResourceId
		m.ProjectName = prior.ProjectName
		m.LogStore = prior.LogStore
	}
	// Echoed fields win over the configured value so state reflects the
	// service-side truth; DescribeResourceLog returns ProjectName and LogStore.
	if v := str("ProjectName"); !v.IsNull() {
		m.ProjectName = v
	}
	if v := str("LogStore"); !v.IsNull() {
		m.LogStore = v
	}
	if m.Id.IsNull() && !m.ClusterId.IsNull() && !m.ResourceId.IsNull() {
		m.Id = types.StringValue(resourceLogID(m.ClusterId.ValueString(), m.ResourceId.ValueString()))
	}
	return m
}

// isResourceLogNotExist reports whether err means the log delivery config is
// absent. EAS ROA returns a 404 / NotFound-shaped error when the parent
// resource has no log config; the exact code is not published, so the standard
// errs.NotFoundError detection is widened with an EAS-specific substring net.
func isResourceLogNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	if e, ok := err.(*tea.SDKError); ok {
		if tea.IntValue(e.StatusCode) == 404 {
			return true
		}
		if code := tea.StringValue(e.Code); code != "" &&
			(strings.Contains(code, "ResourceLog") || strings.Contains(code, "NotFound") || strings.Contains(code, "NotExist")) {
			return true
		}
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "resourcelog") && (strings.Contains(msg, "not exist") || strings.Contains(msg, "not found"))
}

func resourceLogGoneDiag(typeName, clusterId, resourceId string, err error) diag.Diagnostics {
	return diag.Diagnostics{diag.NewWarningDiagnostic(
		fmt.Sprintf("Reading %s: resource log not found", typeName),
		fmt.Sprintf("DescribeResourceLog reported the log delivery config as absent for cluster_id=%s, resource_id=%s; removing from state. Original error: %s", clusterId, resourceId, err),
	)}
}

func resourceLogAction(clusterId, resourceId string) string {
	return fmt.Sprintf("/api/v2/resources/%s/%s/log", clusterId, resourceId)
}

func resourceLogID(clusterId, resourceId string) string {
	return fmt.Sprintf("%s:%s", clusterId, resourceId)
}

// splitResourceLogID recovers cluster_id/resource_id from the composite id,
// falling back to the stored attributes when the id is empty or malformed.
func splitResourceLogID(id, fallbackClusterId, fallbackResourceId string) (string, string) {
	if clusterId, resourceId, ok := parseResourceLogImportID(id); ok {
		return clusterId, resourceId
	}
	return fallbackClusterId, fallbackResourceId
}

// parseResourceLogImportID splits "<cluster_id>:<resource_id>" on the first
// colon. Resource ids are opaque and may contain colons, so only the leading
// segment up to the first ":" is treated as cluster_id and the remainder as
// resource_id, matching the import contract.
func parseResourceLogImportID(id string) (string, string, bool) {
	if id == "" {
		return "", "", false
	}
	idx := strings.Index(id, ":")
	if idx <= 0 || idx == len(id)-1 {
		return "", "", false
	}
	return id[:idx], id[idx+1:], true
}
