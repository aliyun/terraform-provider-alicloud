package ims

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const applicationAssignmentTypeName = "alicloud_ims_application_assignment"

var (
	_ resource.Resource                = &applicationAssignmentResource{}
	_ resource.ResourceWithImportState = &applicationAssignmentResource{}
)

func NewApplicationAssignmentResource() resource.Resource {
	return &applicationAssignmentResource{}
}

type applicationAssignmentResource struct {
	fwadapt.ResourceBase
}

type applicationAssignmentResourceModel struct {
	Id            types.String `tfsdk:"id"`
	AppId         types.String `tfsdk:"app_id"`
	PrincipalType types.String `tfsdk:"principal_type"`
	PrincipalName types.String `tfsdk:"principal_name"`
	CreateTime    types.String `tfsdk:"create_time"`
}

func (r *applicationAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = applicationAssignmentTypeName
}

func (r *applicationAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an IMS Application Assignment resource. " +
			"For information about IMS Application Assignment and how to use it, see [What is Application Assignment](https://www.alibabacloud.com/help/en/ims/developer-reference/api-ims-2019-08-15-createapplicationassignment).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the IMS Application Assignment. It is formatted as `<app_id>:<principal_type>:<principal_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the application to which the principal is assigned.",
			},
			"principal_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The type of the principal. Valid values: `user` and `role`.",
				Validators: []validator.String{
					stringvalidator.OneOf("user", "role"),
				},
			},
			"principal_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the principal. When PrincipalType is `user`, it is the logon name of the RAM user; when PrincipalType is `role`, it is the name of the RAM role.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the assignment was created, in RFC 3339 format (UTC).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *applicationAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applicationAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"AppId":         plan.AppId.ValueString(),
		"PrincipalName": plan.PrincipalName.ValueString(),
		"PrincipalType": plan.PrincipalType.ValueString(),
	}

	if _, err := r.Meta().RpcPost("Ims", "2019-08-15", "CreateApplicationAssignment", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateApplicationAssignment", applicationAssignmentTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateApplicationAssignment", applicationAssignmentTypeName).Error(),
		)
		return
	}

	// CreateApplicationAssignment returns only RequestId and carries no resource
	// identifier, so the assignment is read back via List to confirm it exists and
	// to populate the computed CreateTime before persisting state.
	object, err := getApplicationAssignment(r.Meta(), plan.AppId.ValueString(), plan.PrincipalType.ValueString(), plan.PrincipalName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling ListApplicationAssignments", applicationAssignmentTypeName),
			errs.WrapErrorf(err, "Creating %s: calling ListApplicationAssignments", applicationAssignmentTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateApplicationAssignment", applicationAssignmentTypeName),
			fmt.Sprintf("The assignment (app_id=%s, principal_type=%s, principal_name=%s) is not visible via ListApplicationAssignments right after creation",
				plan.AppId.ValueString(), plan.PrincipalType.ValueString(), plan.PrincipalName.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, applicationAssignmentState(object))...)
}

func (r *applicationAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state applicationAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId, principalType, principalName, ok := parseApplicationAssignmentID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", applicationAssignmentTypeName),
			fmt.Sprintf("The resource id %q must be formatted as `<app_id>:<principal_type>:<principal_name>`", state.Id.ValueString()),
		)
		return
	}

	object, err := getApplicationAssignment(r.Meta(), appId, principalType, principalName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListApplicationAssignments", applicationAssignmentTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListApplicationAssignments", applicationAssignmentTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, applicationAssignmentState(object))...)
}

func (r *applicationAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No UpdateApplicationAssignment API exists and every identity attribute is
	// ForceNew, so an in-place Update never carries a real change. Re-read the
	// live assignment anyway so any server-side drift in CreateTime is picked up
	// before state is persisted, matching the package's ims_user convention.
	var state applicationAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId, principalType, principalName, ok := parseApplicationAssignmentID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", applicationAssignmentTypeName),
			fmt.Sprintf("The resource id %q must be formatted as `<app_id>:<principal_type>:<principal_name>`", state.Id.ValueString()),
		)
		return
	}

	object, err := getApplicationAssignment(r.Meta(), appId, principalType, principalName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling ListApplicationAssignments", applicationAssignmentTypeName),
			errs.WrapErrorf(err, "Updating %s: calling ListApplicationAssignments", applicationAssignmentTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, applicationAssignmentState(object))...)
}

func (r *applicationAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state applicationAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"AppId":         state.AppId.ValueString(),
		"PrincipalName": state.PrincipalName.ValueString(),
		"PrincipalType": state.PrincipalType.ValueString(),
	}

	_, err := r.Meta().RpcPost("Ims", "2019-08-15", "DeleteApplicationAssignment", nil, request, true)
	if err != nil && !isApplicationAssignmentNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteApplicationAssignment", applicationAssignmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteApplicationAssignment", applicationAssignmentTypeName).Error(),
		)
		return
	}
}

func (r *applicationAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The composite id is the single source of truth: it carries the three
	// identity components. Read parses it back into AppId/PrincipalType/
	// PrincipalName before querying ListApplicationAssignments.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getApplicationAssignment queries ListApplicationAssignments for the given
// application, traversing pagination completely, and returns the assignment
// matching PrincipalType and PrincipalName, or nil when it no longer exists.
func getApplicationAssignment(client *connectivity.AliyunClient, appId, principalType, principalName string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"AppId":      appId,
		"MaxResults": 100,
	}
	for {
		response, err := client.RpcPost("Ims", "2019-08-15", "ListApplicationAssignments", nil, request, true)
		if err != nil {
			return nil, err
		}
		container, _ := response["ApplicationAssignment"].(map[string]interface{})
		items, _ := container["ApplicationAssignmentInfo"].([]interface{})
		for _, raw := range items {
			item, _ := raw.(map[string]interface{})
			if item == nil {
				continue
			}
			if fmt.Sprint(item["PrincipalType"]) == principalType && fmt.Sprint(item["PrincipalName"]) == principalName {
				return item, nil
			}
		}
		nextToken, _ := response["NextToken"].(string)
		if nextToken == "" {
			break
		}
		request["NextToken"] = nextToken
	}
	return nil, nil
}

// applicationAssignmentState maps a ListApplicationAssignments item onto the
// resource model. The composite id is rebuilt from the identity fields so
// that imported and created resources share one id representation. Empty
// strings become null, matching the IMS package convention for absent fields.
func applicationAssignmentState(item map[string]interface{}) applicationAssignmentResourceModel {
	str := func(key string) types.String {
		if v, ok := item[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	appId := str("AppId")
	principalType := str("PrincipalType")
	principalName := str("PrincipalName")
	id := types.StringNull()
	if !appId.IsNull() && !principalType.IsNull() && !principalName.IsNull() {
		id = types.StringValue(fmt.Sprintf("%s:%s:%s", appId.ValueString(), principalType.ValueString(), principalName.ValueString()))
	}
	return applicationAssignmentResourceModel{
		Id:            id,
		AppId:         appId,
		PrincipalType: principalType,
		PrincipalName: principalName,
		CreateTime:    str("CreateTime"),
	}
}

// parseApplicationAssignmentID splits the composite id back into its three
// identity components. ok is false when the id is malformed, which Read
// surfaces as an error rather than a silent NotFound.
func parseApplicationAssignmentID(id string) (appId, principalType, principalName string, ok bool) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// isApplicationAssignmentNotExist reports whether err is the IMS
// assignment-not-found error. Create returns no id and Delete is keyed by the
// three identity params, so an already-deleted assignment surfaces as a
// not-found error that Delete must treat as success. The code arrives on a
// *tea.SDKError from RpcPost; the substring fallback covers wrapped forms,
// matching isUserNotExist's behavior on the SDK line.
func isApplicationAssignmentNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		switch *e.Code {
		case "EntityNotExist.ApplicationAssignment", "NotEntityExists.ApplicationAssignment":
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "EntityNotExist.ApplicationAssignment") ||
		strings.Contains(msg, "NotEntityExists.ApplicationAssignment")
}
