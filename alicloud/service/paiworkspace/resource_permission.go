package paiworkspace

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const permissionTypeName = "alicloud_pai_workspace_permission"

// PAI Workspace (AIWorkSpace 2021-02-04) product/version used by every ROA
// call in this package. Permission is a read-only catalog object: the API
// exposes only GetPermission and ListPermissions, so the resource has no
// Create/Update/Delete lifecycle against the cloud.
const (
	permissionProduct = "AIWorkSpace"
	permissionVersion = "2021-02-04"
	permissionIdSep   = ":"
)

var _ resource.Resource = &permissionResource{}
var _ resource.ResourceWithImportState = &permissionResource{}

func NewPermissionResource() resource.Resource {
	return &permissionResource{}
}

type permissionResource struct {
	fwadapt.ResourceBase
}

// permissionResourceModel holds the resource state. The PermissionRules array
// is flattened into two parallel computed lists rather than a nested attribute:
// this provider is served over protocol v5 by a tf5muxserver that muxes the SDK
// v2 and framework halves, and the v5 schema conversion rejects nested
// attributes (a SchemaAttribute with a nested_type/Attributes representation)
// with "protocol version 5 cannot have Attributes set". Framework blocks
// (ListNestedBlock) are config-only and expose no Computed field, so they
// cannot surface a read-only computed collection either. A list of strings is
// the v5-native form the SDK v2 half emits for every []string field, so the
// two parallel lists (accessibility, entity_access_type) surface the rule
// catalog with collection-Read distinguishing power (R4.23/R7.16) while
// staying v5-clean. Index correspondence between the two lists is preserved so
// a caller can pair each accessibility with its entity_access_type.
//
// The two list fields are types.List (not []types.String) because a
// schema.ListAttribute carries an unknown value during planning (it is
// Computed), and a plain []basetypes.StringValue slice cannot hold unknowns —
// the framework's tfsdk decode would fail with "Received unknown value,
// however the target type cannot handle unknown values. Suggested Type:
// basetypes.ListValue". types.List is the canonical Go type paired with
// ListAttribute and accepts the unknown/computed values the planner emits.
type permissionResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	WorkspaceId                     types.String `tfsdk:"workspace_id"`
	PermissionCode                  types.String `tfsdk:"permission_code"`
	PermissionRuleAccessibilities   types.List   `tfsdk:"permission_rule_accessibilities"`
	PermissionRuleEntityAccessTypes types.List   `tfsdk:"permission_rule_entity_access_types"`
}

func (r *permissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = permissionTypeName
}

func (r *permissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Workspace Permission resource. " +
			"PAI Workspace permissions are read-only catalog entries defined by the `GetPermission` and " +
			"`ListPermissions` APIs; this resource imports an existing permission into Terraform state " +
			"and does not create, update or delete it on the cloud.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, formatted as `<workspace_id>:<permission_code>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PAI Workspace. How to get the workspace ID, see ListWorkspaces.",
				PlanModifiers: []planmodifier.String{
					// A different workspace scopes a different permission, so the
					// resource identity changes: replace rather than update.
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permission_code": schema.StringAttribute{
				Required: true,
				Description: "The name of the permission point, unique within the region, e.g. `PaiDLC:ListJobs`. " +
					"How to get the permission code, see ListPermissions.",
				PlanModifiers: []planmodifier.String{
					// Permission codes contain a colon (e.g. PaiDLC:ListJobs); a
					// different code is a different permission identity.
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permission_rule_accessibilities": schema.ListAttribute{
				Computed: true,
				Description: "The Accessibility value of each rule of the permission point, in catalog order. " +
					"Valid values: `PUBLIC`, `PRIVATE`, `ANY`. Index-corresponds with `permission_rule_entity_access_types`.",
				ElementType: types.StringType,
			},
			"permission_rule_entity_access_types": schema.ListAttribute{
				Computed: true,
				Description: "The EntityAccessType value of each rule of the permission point, in catalog order. " +
					"Valid values: `CREATOR`, `ANY`. Index-corresponds with `permission_rule_accessibilities`.",
				ElementType: types.StringType,
			},
		},
	}
}

// Create is read-only: there is no CreatePermission API. The resource
// represents an existing catalog permission, so Create verifies the permission
// is readable via GetPermission and populates the computed state. No write API
// is called.
func (r *permissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId := plan.WorkspaceId.ValueString()
	permissionCode := plan.PermissionCode.ValueString()
	plan.Id = types.StringValue(permissionID(workspaceId, permissionCode))

	object, err := getPermission(r.Client(), workspaceId, permissionCode)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading permission", permissionTypeName),
			errs.WrapErrorf(err, "Creating %s: calling GetPermission for %s/%s", permissionTypeName, workspaceId, permissionCode).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: permission not found", permissionTypeName),
			fmt.Sprintf("The permission %s does not exist in workspace %s; PAI Workspace permissions are read-only catalog entries and cannot be created.", permissionCode, workspaceId),
		)
		return
	}
	plan.PermissionRuleAccessibilities, plan.PermissionRuleEntityAccessTypes = permissionRuleListsFromResponse(object)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *permissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, permissionCode, err := parsePermissionID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", permissionTypeName),
			err.Error(),
		)
		return
	}

	object, err := getPermission(r.Client(), workspaceId, permissionCode)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPermission", permissionTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPermission for %s/%s", permissionTypeName, workspaceId, permissionCode).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	state.WorkspaceId = types.StringValue(workspaceId)
	state.PermissionCode = types.StringValue(permissionCode)
	state.PermissionRuleAccessibilities, state.PermissionRuleEntityAccessTypes = permissionRuleListsFromResponse(object)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is read-only: there is no UpdatePermission API. The only user-settable
// attributes force replacement, so this is reached to refresh the computed
// rule lists from the API rather than to mutate the cloud object.
func (r *permissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan permissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId := plan.WorkspaceId.ValueString()
	permissionCode := plan.PermissionCode.ValueString()
	plan.Id = types.StringValue(permissionID(workspaceId, permissionCode))

	object, err := getPermission(r.Client(), workspaceId, permissionCode)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling GetPermission", permissionTypeName),
			errs.WrapErrorf(err, "Updating %s: calling GetPermission for %s/%s", permissionTypeName, workspaceId, permissionCode).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	plan.PermissionRuleAccessibilities, plan.PermissionRuleEntityAccessTypes = permissionRuleListsFromResponse(object)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is read-only: there is no DeletePermission API. Terraform removes the
// resource from state; the cloud permission (a catalog entry) is untouched.
func (r *permissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state permissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[WARN] Cannot destroy %s %s. Terraform will remove this resource from the state file; the PAI Workspace permission catalog entry remains untouched.", permissionTypeName, state.Id.ValueString())
}

// ImportState discovers the permission by enumerating the workspace's
// permissions through ListPermissions and matching the requested code, so
// that importing one permission among many resolves the correct identity
// (rather than a blind id passthrough). The subsequent Read refreshes it via
// GetPermission.
func (r *permissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	workspaceId, permissionCode, err := parsePermissionID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: parsing resource id", permissionTypeName),
			err.Error(),
		)
		return
	}

	permissions, err := listPermissions(r.Client(), workspaceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: calling ListPermissions", permissionTypeName),
			errs.WrapErrorf(err, "Importing %s: calling ListPermissions for workspace %s", permissionTypeName, workspaceId).Error(),
		)
		return
	}
	matched := findPermission(permissions, permissionCode)
	if matched == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: permission not found", permissionTypeName),
			fmt.Sprintf("The permission %s was not found among the %d permissions listed in workspace %s.", permissionCode, len(permissions), workspaceId),
		)
		return
	}

	accessibilities, entityAccessTypes := permissionRuleListsFromObject(matched)
	model := permissionResourceModel{
		Id:                              types.StringValue(permissionID(workspaceId, permissionCode)),
		WorkspaceId:                     types.StringValue(workspaceId),
		PermissionCode:                  types.StringValue(permissionCode),
		PermissionRuleAccessibilities:   accessibilities,
		PermissionRuleEntityAccessTypes: entityAccessTypes,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// permissionID formats the composite resource id. Permission codes contain a
// colon (e.g. PaiDLC:ListJobs), so the id is split on the FIRST separator
// only; parsePermissionID mirrors this.
func permissionID(workspaceId, permissionCode string) string {
	return fmt.Sprintf("%s%s%s", workspaceId, permissionIdSep, permissionCode)
}

func parsePermissionID(id string) (string, string, error) {
	if id == "" {
		return "", "", fmt.Errorf("empty resource id")
	}
	parts := strings.SplitN(id, permissionIdSep, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid resource id %q: expected <workspace_id>%s<permission_code>", id, permissionIdSep)
	}
	return parts[0], parts[1], nil
}

// getPermission queries GetPermission and returns the response object, or nil
// when the permission no longer exists.
func getPermission(client *connectivity.AliyunClient, workspaceId, permissionCode string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v1/workspaces/%s/permissions/%s", workspaceId, permissionCode)
	response, err := client.RoaGet(permissionProduct, permissionVersion, action, nil, nil, nil)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// listPermissions queries ListPermissions and returns the Permissions array
// from the response.
func listPermissions(client *connectivity.AliyunClient, workspaceId string) ([]interface{}, error) {
	action := fmt.Sprintf("/api/v1/workspaces/%s/permissions", workspaceId)
	response, err := client.RoaGet(permissionProduct, permissionVersion, action, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	permissions, _ := response["Permissions"].([]interface{})
	return permissions, nil
}

func findPermission(permissions []interface{}, permissionCode string) map[string]interface{} {
	for _, raw := range permissions {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if code, _ := item["PermissionCode"].(string); code == permissionCode {
			return item
		}
	}
	return nil
}

// permissionRuleListsFromResponse extracts the Accessibility and EntityAccessType
// values from a GetPermission response's PermissionRules array as two parallel
// computed lists.
func permissionRuleListsFromResponse(response map[string]interface{}) (types.List, types.List) {
	rules, _ := response["PermissionRules"].([]interface{})
	return permissionRuleLists(rules)
}

// permissionRuleListsFromObject extracts the same two lists from a single
// permission object in a ListPermissions response.
func permissionRuleListsFromObject(object map[string]interface{}) (types.List, types.List) {
	rules, _ := object["PermissionRules"].([]interface{})
	return permissionRuleLists(rules)
}

// permissionRuleLists builds the Accessibility and EntityAccessType lists in
// catalog order, preserving index correspondence between the two. Each list
// is returned as a types.List (the Go type paired with schema.ListAttribute),
// built from []attr.Value so the framework can serialize the computed
// collection and accept the unknown values the planner produces before Read.
func permissionRuleLists(rules []interface{}) (types.List, types.List) {
	accessVals := make([]attr.Value, 0, len(rules))
	entityVals := make([]attr.Value, 0, len(rules))
	for _, raw := range rules {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		accessVals = append(accessVals, stringAttrValue(item, "Accessibility"))
		entityVals = append(entityVals, stringAttrValue(item, "EntityAccessType"))
	}
	return stringListValue(accessVals), stringListValue(entityVals)
}

// stringListValue wraps a slice of string attr.Values into a types.List of
// StringType. The elements are validated StringValue/Null values, so the
// conversion cannot fail in practice; on a diagnostic error (a provider bug)
// it falls back to an empty list rather than surfacing an invalid value.
func stringListValue(vals []attr.Value) types.List {
	list, diags := types.ListValue(types.StringType, vals)
	if diags.HasError() {
		list, _ = types.ListValue(types.StringType, nil)
	}
	return list
}

func stringAttrValue(object map[string]interface{}, key string) attr.Value {
	if v, ok := object[key].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}
