package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	projectRoleTypeName = "alicloud_data_works_project_role"
	// Framework does not support user-configurable resource timeouts, so the
	// create/update durations replace the SDKv2 defaults. All ProjectRole CRUD
	// APIs are synchronous and return immediately, so these bound only retries.
	projectRoleCreateTimeout = 5 * time.Minute
	projectRoleUpdateTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &projectRoleResource{}
	_ resource.ResourceWithImportState = &projectRoleResource{}
)

func NewProjectRoleResource() resource.Resource {
	return &projectRoleResource{}
}

type projectRoleResource struct {
	fwadapt.ResourceBase
}

type projectRoleResourceModel struct {
	Id                types.String            `tfsdk:"id"`
	ProjectId         types.String            `tfsdk:"project_id"`
	RoleName          types.String            `tfsdk:"role_name"`
	Code              types.String            `tfsdk:"code"`
	Type              types.String            `tfsdk:"type"`
	ModulePermissions []modulePermissionModel `tfsdk:"module_permissions"`
}

type modulePermissionModel struct {
	ModuleId       types.String `tfsdk:"module_id"`
	PermissionType types.String `tfsdk:"permission_type"`
}

func (r *projectRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = projectRoleTypeName
}

func (r *projectRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"code": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// SetNestedAttribute, not ListNestedAttribute: module_permissions is an
			// order-independent collection keyed by module_id. A List would permadiff
			// whenever a user configures modules in non-ascending order, because the
			// Read canonicalizes state by sorting on module_id. Set comparison is by
			// element value, so any configured order is stable.
			"module_permissions": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"module_id": schema.StringAttribute{
							Required: true,
						},
						"permission_type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("Write", "Read", "NotSet"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *projectRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectRoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissions, err := projectRoleModulePermissionsRequest(plan.ModulePermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building ModulePermissions request", projectRoleTypeName),
			err.Error(),
		)
		return
	}

	action := "CreateProjectRole"
	request := map[string]interface{}{
		"ProjectId":         plan.ProjectId.ValueString(),
		"Name":              plan.RoleName.ValueString(),
		"RegionId":          r.Meta().RegionId,
		"ModulePermissions": permissions,
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, projectRoleCreateTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) || errs.IsExpectedErrors(callErr, []string{"9990020002", "9990040003"}) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateProjectRole", projectRoleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, projectRoleTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	code, err := projectRoleRequiredString(response["Code"], "Code")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading Code", projectRoleTypeName),
			errs.WrapErrorf(err, errs.FailedGetAttributeMsg, projectRoleTypeName, "$.Code", response).Error(),
		)
		return
	}
	id := fmt.Sprintf("%s:%s", plan.ProjectId.ValueString(), code.ValueString())

	object, err := NewService(r.Meta()).DescribeDataWorksProjectRoleContext(ctx, id)
	if err != nil {
		// Keep the created id in state so the role stays trackable even when the
		// read-back fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling GetProjectRole", projectRoleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetProjectRole", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := projectRoleState(object)
	if err != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping GetProjectRole response", projectRoleTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *projectRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Meta()).DescribeDataWorksProjectRoleContext(ctx, id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetProjectRole", projectRoleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetProjectRole", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := projectRoleState(object)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping GetProjectRole response", projectRoleTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *projectRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectRoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// project_id and role_name are replacement-only, so the only in-place
	// change is module_permissions, which UpdateProjectRole replaces wholesale.
	permissions, err := projectRoleModulePermissionsRequest(plan.ModulePermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: building ModulePermissions request", projectRoleTypeName),
			err.Error(),
		)
		return
	}

	id := plan.Id.ValueString()
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", projectRoleTypeName),
			fmt.Errorf("invalid Resource Id %s. Expected parts' length 2, got %d", id, len(parts)).Error(),
		)
		return
	}

	action := "UpdateProjectRole"
	request := map[string]interface{}{
		"ProjectId":         parts[0],
		"Code":              parts[1],
		"RegionId":          r.Meta().RegionId,
		"ModulePermissions": permissions,
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, projectRoleUpdateTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) || errs.IsExpectedErrors(callErr, []string{"9990020002", "9990040003"}) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateProjectRole", projectRoleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDataWorksProjectRoleContext(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling GetProjectRole", projectRoleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetProjectRole", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	model, err := projectRoleState(object)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: mapping GetProjectRole response", projectRoleTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *projectRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if err := NewService(r.Meta()).DeleteDataWorksProjectRoleContext(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteProjectRole", projectRoleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteProjectRole", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *projectRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// projectRoleState maps a GetProjectRole ProjectRole object into the resource
// model. It is shared by Create, Read and Update so refresh semantics stay
// identical: every attribute is rebuilt from the server response, and a
// contract-drifted value fails the refresh instead of overwriting configured
// state with empty values.
func projectRoleState(object map[string]interface{}) (projectRoleResourceModel, error) {
	projectId, err := projectRoleCanonicalInt(object["ProjectId"], "ProjectId")
	if err != nil {
		return projectRoleResourceModel{}, err
	}
	code, err := projectRoleRequiredString(object["Code"], "Code")
	if err != nil {
		return projectRoleResourceModel{}, err
	}
	name, err := projectRoleRequiredString(object["Name"], "Name")
	if err != nil {
		return projectRoleResourceModel{}, err
	}
	roleType, err := projectRoleStringField(object["Type"], "Type")
	if err != nil {
		return projectRoleResourceModel{}, err
	}
	permissions, err := projectRoleModulePermissionsState(object["ModulePermissions"])
	if err != nil {
		return projectRoleResourceModel{}, err
	}
	return projectRoleResourceModel{
		Id:                types.StringValue(fmt.Sprintf("%s:%s", strconv.FormatInt(projectId, 10), code.ValueString())),
		ProjectId:         types.StringValue(strconv.FormatInt(projectId, 10)),
		RoleName:          name,
		Code:              code,
		Type:              roleType,
		ModulePermissions: permissions,
	}, nil
}

// projectRoleModulePermissionsRequest builds the ModulePermissions JSON array
// sent to CreateProjectRole/UpdateProjectRole. ModuleId is parsed to int64 so
// the SDK serializes it as a JSON number, matching the API contract.
func projectRoleModulePermissionsRequest(permissions []modulePermissionModel) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, 0, len(permissions))
	for i, p := range permissions {
		moduleId, err := projectRoleCanonicalInt(p.ModuleId.ValueString())
		if err != nil {
			return nil, fmt.Errorf("module_permissions[%d].module_id: %w", i, err)
		}
		out = append(out, map[string]interface{}{
			"ModuleId":       moduleId,
			"PermissionType": p.PermissionType.ValueString(),
		})
	}
	return out, nil
}

// projectRoleModulePermissionsState maps the ModulePermissions array returned by
// GetProjectRole into the model, sorted by module_id so a later plan with the
// same configured set produces no diff regardless of the server-side ordering.
func projectRoleModulePermissionsState(raw interface{}) ([]modulePermissionModel, error) {
	if raw == nil {
		return []modulePermissionModel{}, nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("ModulePermissions has unexpected type %T", raw)
	}
	type parsed struct {
		id   int64
		mid  string
		perm string
	}
	entries := make([]parsed, 0, len(list))
	for i, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ModulePermissions[%d] has unexpected type %T", i, item)
		}
		moduleId, err := projectRoleCanonicalInt(m["ModuleId"], "ModuleId")
		if err != nil {
			return nil, fmt.Errorf("ModulePermissions[%d].ModuleId: %w", i, err)
		}
		permissionType, ok := m["PermissionType"].(string)
		if !ok {
			return nil, fmt.Errorf("ModulePermissions[%d].PermissionType has unexpected type %T", i, m["PermissionType"])
		}
		entries = append(entries, parsed{id: moduleId, mid: strconv.FormatInt(moduleId, 10), perm: permissionType})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].id < entries[j].id })
	out := make([]modulePermissionModel, 0, len(entries))
	for _, e := range entries {
		out = append(out, modulePermissionModel{
			ModuleId:       types.StringValue(e.mid),
			PermissionType: types.StringValue(e.perm),
		})
	}
	return out, nil
}

// projectRoleRequiredString rejects a missing or non-string value for a
// required attribute (Code, Name) instead of silently writing null over
// configured state.
func projectRoleRequiredString(raw interface{}, field string) (types.String, error) {
	v, ok := raw.(string)
	if !ok || v == "" {
		return types.StringNull(), fmt.Errorf("%s is required in the GetProjectRole response, got %T %v", field, raw, raw)
	}
	return types.StringValue(v), nil
}

// projectRoleStringField accepts a string or nil for a computed display field
// (Type), returning null when the value is absent and failing only on a
// malformed non-string shape.
func projectRoleStringField(raw interface{}, field string) (types.String, error) {
	switch v := raw.(type) {
	case nil:
		return types.StringNull(), nil
	case string:
		if v == "" {
			return types.StringNull(), nil
		}
		return types.StringValue(v), nil
	default:
		return types.StringNull(), fmt.Errorf("%s has unexpected type %T", field, raw)
	}
}

// projectRoleCanonicalInt converts a numeric identity (ProjectId, ModuleId)
// from either the API response or the configured string into a canonical
// int64. The SDK decodes JSON numbers as float64, but json.Number is also
// accepted for forward safety. Fractional, missing or non-numeric values are
// contract drift and are rejected instead of being stringified into state.
func projectRoleCanonicalInt(raw interface{}, field ...string) (int64, error) {
	name := "value"
	if len(field) > 0 {
		name = field[0]
	}
	switch v := raw.(type) {
	case nil:
		return 0, fmt.Errorf("%s is missing", name)
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("%s %q is not an int64: %w", name, v.String(), err)
		}
		return n, nil
	case float64:
		if v != float64(int64(v)) {
			return 0, fmt.Errorf("%s %v is not an integer", name, v)
		}
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s %q is not an int64: %w", name, v, err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("%s has unexpected type %T", name, raw)
	}
}
