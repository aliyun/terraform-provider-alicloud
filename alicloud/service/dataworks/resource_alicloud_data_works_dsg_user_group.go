package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var (
	_ resource.Resource                = &dsgUserGroupResource{}
	_ resource.ResourceWithImportState = &dsgUserGroupResource{}
)

// dsgUserGroupProjectsAttrTypes is the attribute type map of the nested
// projects object, shared by the schema and the list element type.
var dsgUserGroupProjectsAttrTypes = map[string]attr.Type{
	"project_name": types.StringType,
	"accounts":     types.StringType,
}

func NewDsgUserGroupResource() resource.Resource {
	return &dsgUserGroupResource{}
}

type dsgUserGroupResource struct {
	fwadapt.ResourceBase
}

type dsgUserGroupResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Owner         types.String `tfsdk:"owner"`
	UserGroupType types.Int64  `tfsdk:"user_group_type"`
	ProjectName   types.String `tfsdk:"project_name"`
	Accounts      types.Set    `tfsdk:"accounts"`
	Projects      types.List   `tfsdk:"projects"`
	GmtCreate     types.String `tfsdk:"gmt_create"`
	GmtModified   types.String `tfsdk:"gmt_modified"`
}

func (r *dsgUserGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = dsgUserGroupTypeName
}

func (r *dsgUserGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owner": schema.StringAttribute{
				Required: true,
			},
			"user_group_type": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2, 3),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			// project_name is a write-only field: it is required when
			// user_group_type is 3 (MaxCompute role) but DsgUserGroupQueryList
			// never returns it, so the configured value is preserved in state.
			"project_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"accounts": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"projects": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"project_name": schema.StringAttribute{
							Required: true,
						},
						"accounts": schema.StringAttribute{
							Optional: true,
						},
					},
				},
				Optional: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_create": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dsgUserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dsgUserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userGroup, err := dsgUserGroupBuildUserGroup(0, false, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building DsgUserGroupAddOrUpdate request", dsgUserGroupTypeName),
			err.Error(),
		)
		return
	}
	if err := dsgUserGroupValidate(plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: invalid configuration", dsgUserGroupTypeName),
			err.Error(),
		)
		return
	}

	userGroupsJSON, err := json.Marshal([]interface{}{userGroup})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: serializing UserGroups", dsgUserGroupTypeName),
			err.Error(),
		)
		return
	}
	request := map[string]interface{}{
		"RegionId":   r.Meta().RegionId,
		"UserGroups": string(userGroupsJSON),
	}

	if err := NewService(r.Meta()).DsgUserGroupAddOrUpdate(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling DsgUserGroupAddOrUpdate", dsgUserGroupTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, dsgUserGroupTypeName, "DsgUserGroupAddOrUpdate", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// DsgUserGroupAddOrUpdate returns a boolean only; the created group's Id is
	// resolved by querying DsgUserGroupQueryList by name and retrying until the
	// new group becomes visible.
	name := plan.Name.ValueString()
	var object map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, dsgUserGroupTimeout, func() *retry.RetryError {
		item, callErr := NewService(r.Meta()).DsgUserGroupDescribe(ctx, "", name)
		if callErr != nil {
			if errs.NotFoundError(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		object = item
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: locating the created user group via DsgUserGroupQueryList", dsgUserGroupTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, name, "DsgUserGroupQueryList", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mErr := dsgUserGroupState(object, plan)
	if mErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: mapping DsgUserGroupQueryList response", dsgUserGroupTypeName),
			mErr.Error(),
		)
		return
	}
	dsgUserGroupPreserveEmptyConfigured(&model, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dsgUserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dsgUserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Meta()).DsgUserGroupDescribe(ctx, id, state.Name.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DsgUserGroupQueryList", dsgUserGroupTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DsgUserGroupQueryList", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mErr := dsgUserGroupState(object, state)
	if mErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping DsgUserGroupQueryList response", dsgUserGroupTypeName),
			mErr.Error(),
		)
		return
	}
	dsgUserGroupPreserveEmptyConfigured(&model, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dsgUserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dsgUserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state dsgUserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := dsgUserGroupValidate(plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: invalid configuration", dsgUserGroupTypeName),
			err.Error(),
		)
		return
	}

	idInt, perr := dsgUserGroupParseId(state.Id.ValueString())
	if perr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing id", dsgUserGroupTypeName),
			perr.Error(),
		)
		return
	}
	userGroup, err := dsgUserGroupBuildUserGroup(idInt, true, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: building DsgUserGroupAddOrUpdate request", dsgUserGroupTypeName),
			err.Error(),
		)
		return
	}
	userGroupsJSON, err := json.Marshal([]interface{}{userGroup})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: serializing UserGroups", dsgUserGroupTypeName),
			err.Error(),
		)
		return
	}
	request := map[string]interface{}{
		"RegionId":   r.Meta().RegionId,
		"UserGroups": string(userGroupsJSON),
	}

	if err := NewService(r.Meta()).DsgUserGroupAddOrUpdate(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling DsgUserGroupAddOrUpdate", dsgUserGroupTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DsgUserGroupAddOrUpdate", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// Refresh from the source of truth so computed and mutated fields reflect
	// the real state after the edit.
	object, err := NewService(r.Meta()).DsgUserGroupDescribe(ctx, state.Id.ValueString(), plan.Name.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: refreshing via DsgUserGroupQueryList", dsgUserGroupTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DsgUserGroupQueryList", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	model, mErr := dsgUserGroupState(object, plan)
	if mErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: mapping DsgUserGroupQueryList response", dsgUserGroupTypeName),
			mErr.Error(),
		)
		return
	}
	dsgUserGroupPreserveEmptyConfigured(&model, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dsgUserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dsgUserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.Id.ValueString()
	if id == "" {
		return
	}
	if err := NewService(r.Meta()).DsgUserGroupDelete(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DsgUserGroupDelete", dsgUserGroupTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DsgUserGroupDelete", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *dsgUserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// dsgUserGroupValidate enforces conditional requirements: a MaxCompute user
// group (user_group_type = 3) must set project_name.
func dsgUserGroupValidate(plan dsgUserGroupResourceModel) error {
	if plan.UserGroupType.ValueInt64() == 3 {
		if plan.ProjectName.IsNull() || plan.ProjectName.IsUnknown() || plan.ProjectName.ValueString() == "" {
			return fmt.Errorf("project_name must be set when user_group_type is 3 (MaxCompute role)")
		}
	}
	return nil
}

// dsgUserGroupBuildUserGroup assembles a single UserGroups element. When
// hasId is false the group is created; when true it is edited by Id.
func dsgUserGroupBuildUserGroup(id int64, hasId bool, plan dsgUserGroupResourceModel) (map[string]interface{}, error) {
	userGroup := map[string]interface{}{
		"Name":          plan.Name.ValueString(),
		"Owner":         plan.Owner.ValueString(),
		"UserGroupType": plan.UserGroupType.ValueInt64(),
	}
	if hasId {
		userGroup["Id"] = id
	}
	if !plan.ProjectName.IsNull() && !plan.ProjectName.IsUnknown() && plan.ProjectName.ValueString() != "" {
		userGroup["ProjectName"] = plan.ProjectName.ValueString()
	}
	accounts := dsgUserGroupAccountsFromSet(plan.Accounts)
	if len(accounts) > 0 {
		userGroup["Accounts"] = accounts
	}
	projectsJSON, err := dsgUserGroupProjectsToJSON(plan.Projects)
	if err != nil {
		return nil, fmt.Errorf("serializing projects: %w", err)
	}
	if projectsJSON != "" {
		userGroup["projects"] = projectsJSON
	}
	return userGroup, nil
}

// dsgUserGroupState maps a DsgUserGroupQueryList item to the resource model.
// user_group_type and project_name are write-only (not returned by QueryList),
// so they are preserved from the source (plan/state).
func dsgUserGroupState(object map[string]interface{}, source dsgUserGroupResourceModel) (dsgUserGroupResourceModel, error) {
	id, err := dsgUserGroupRequiredString(object["Id"], "Id")
	if err != nil {
		return dsgUserGroupResourceModel{}, err
	}
	name, err := dsgUserGroupRequiredString(object["Name"], "Name")
	if err != nil {
		return dsgUserGroupResourceModel{}, err
	}
	owner, err := dsgUserGroupRequiredString(object["Owner"], "Owner")
	if err != nil {
		return dsgUserGroupResourceModel{}, err
	}
	projects, err := dsgUserGroupProjectsFromResponse(object["mcAggregationInfo"])
	if err != nil {
		return dsgUserGroupResourceModel{}, err
	}
	return dsgUserGroupResourceModel{
		Id:            id,
		Name:          name,
		Owner:         owner,
		UserGroupType: source.UserGroupType,
		ProjectName:   source.ProjectName,
		Accounts:      dsgUserGroupAccountsSet(object["Accounts"]),
		Projects:      projects,
		GmtCreate:     dsgUserGroupString(object["GmtCreate"]),
		GmtModified:   dsgUserGroupString(object["GmtModified"]),
	}, nil
}

// dsgUserGroupPreserveEmptyConfigured keeps an explicitly empty configured
// collection in state. DsgUserGroupAddOrUpdate omits Accounts and projects when
// they are empty, and DsgUserGroupQueryList returns nothing for an empty
// collection, so normalizing them to null would fail the post-apply consistency
// check and plan a perpetual in-place update whenever a user sets an empty
// collection explicitly.
func dsgUserGroupPreserveEmptyConfigured(model *dsgUserGroupResourceModel, source dsgUserGroupResourceModel) {
	if model.Accounts.IsNull() && !source.Accounts.IsNull() && !source.Accounts.IsUnknown() && len(source.Accounts.Elements()) == 0 {
		model.Accounts = source.Accounts
	}
	if model.Projects.IsNull() && !source.Projects.IsNull() && !source.Projects.IsUnknown() && len(source.Projects.Elements()) == 0 {
		model.Projects = source.Projects
	}
}

// dsgUserGroupRequiredString rejects a missing or non-string value for a
// required attribute instead of silently writing null over configured state.
func dsgUserGroupRequiredString(raw interface{}, field string) (types.String, error) {
	if v, ok := raw.(string); ok && v != "" {
		return types.StringValue(v), nil
	}
	return types.StringNull(), fmt.Errorf("%s is required in the DsgUserGroupQueryList response, got %T %v", field, raw, raw)
}

func dsgUserGroupString(raw interface{}) types.String {
	if v, ok := raw.(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// dsgUserGroupAccountsFromSet flattens a set of strings into a slice for the
// Accounts request field.
func dsgUserGroupAccountsFromSet(value types.Set) []string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	out := make([]string, 0, len(value.Elements()))
	for _, v := range value.Elements() {
		if s, ok := v.(types.String); ok && !s.IsNull() && !s.IsUnknown() && s.ValueString() != "" {
			out = append(out, s.ValueString())
		}
	}
	return out
}

// dsgUserGroupAccountsSet builds a set of strings from the Accounts response
// array.
func dsgUserGroupAccountsSet(raw interface{}) types.Set {
	arr, _ := raw.([]interface{})
	elems := make([]attr.Value, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			elems = append(elems, types.StringValue(s))
		}
	}
	if len(elems) == 0 {
		return types.SetNull(types.StringType)
	}
	return types.SetValueMust(types.StringType, elems)
}

// dsgUserGroupProjectsToJSON serializes the projects list into the JSON string
// the API expects ([{"projectName":...,"accounts":...}]).
func dsgUserGroupProjectsToJSON(value types.List) (string, error) {
	if value.IsNull() || value.IsUnknown() || len(value.Elements()) == 0 {
		return "", nil
	}
	items := make([]map[string]interface{}, 0, len(value.Elements()))
	for _, v := range value.Elements() {
		obj, ok := v.(types.Object)
		if !ok {
			return "", fmt.Errorf("projects element has unexpected type %T", v)
		}
		item := map[string]interface{}{}
		if pn, ok := obj.Attributes()["project_name"]; ok {
			if s, ok := pn.(types.String); ok && !s.IsNull() && !s.IsUnknown() && s.ValueString() != "" {
				item["projectName"] = s.ValueString()
			}
		}
		if ac, ok := obj.Attributes()["accounts"]; ok {
			if s, ok := ac.(types.String); ok && !s.IsNull() && !s.IsUnknown() && s.ValueString() != "" {
				item["accounts"] = s.ValueString()
			}
		}
		if len(item) == 0 {
			continue
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return "", nil
	}
	b, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// dsgUserGroupProjectsFromResponse parses the mcAggregationInfo JSON string
// returned by DsgUserGroupQueryList into the projects list.
func dsgUserGroupProjectsFromResponse(raw interface{}) (types.List, error) {
	objType := types.ObjectType{AttrTypes: dsgUserGroupProjectsAttrTypes}
	s, ok := raw.(string)
	if !ok || s == "" {
		return types.ListNull(objType), nil
	}
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return types.List{}, fmt.Errorf("mcAggregationInfo is not valid JSON: %w", err)
	}
	if len(items) == 0 {
		return types.ListNull(objType), nil
	}
	elems := make([]attr.Value, 0, len(items))
	for _, it := range items {
		attrs := map[string]attr.Value{
			"project_name": types.StringNull(),
			"accounts":     types.StringNull(),
		}
		if pn, ok := it["projectName"].(string); ok && pn != "" {
			attrs["project_name"] = types.StringValue(pn)
		}
		if ac, ok := it["accounts"].(string); ok && ac != "" {
			attrs["accounts"] = types.StringValue(ac)
		}
		objVal := types.ObjectValueMust(dsgUserGroupProjectsAttrTypes, attrs)
		elems = append(elems, objVal)
	}
	return types.ListValueMust(objType, elems), nil
}

// dsgUserGroupParseId converts the resource id string to an int64 for the API.
func dsgUserGroupParseId(id string) (int64, error) {
	var i int64
	if _, err := fmt.Sscanf(id, "%d", &i); err != nil {
		return 0, fmt.Errorf("invalid user group id %q: %w", id, err)
	}
	return i, nil
}
