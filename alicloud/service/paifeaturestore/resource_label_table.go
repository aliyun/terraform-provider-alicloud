package paifeaturestore

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	labelTableTypeName = "alicloud_pai_feature_store_label_table"
	labelTableProduct  = "PaiFeatureStore"
	labelTableVersion  = "2023-06-21"
)

var (
	_ resource.Resource                = &labelTableResource{}
	_ resource.ResourceWithImportState = &labelTableResource{}
	_ resource.ResourceWithModifyPlan  = &labelTableResource{}
)

// labelTableFieldAttrTypes is the element type of the fields list. The provider
// is served over protocol 5, where the framework rejects nested attributes at
// schema load, so the field list is a flat ListAttribute of ObjectType instead
// of a ListNestedAttribute.
var labelTableFieldAttrTypes = map[string]attr.Type{
	"name":                types.StringType,
	"type":                types.StringType,
	"attributes":          types.SetType{ElemType: types.StringType},
	"aligned_entity_name": types.StringType,
}

func NewLabelTableResource() resource.Resource {
	return &labelTableResource{}
}

type labelTableResource struct {
	fwadapt.ResourceBase
}

type labelTableResourceModel struct {
	Id             types.String `tfsdk:"id"`
	InstanceId     types.String `tfsdk:"instance_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	LabelTableName types.String `tfsdk:"label_table_name"`
	DatasourceId   types.String `tfsdk:"datasource_id"`
	Fields         types.List   `tfsdk:"fields"`
	LabelTableId   types.String `tfsdk:"label_table_id"`
	CreateTime     types.String `tfsdk:"create_time"`
	RegionId       types.String `tfsdk:"region_id"`
}

type labelTableFieldModel struct {
	Name              types.String `tfsdk:"name"`
	Type              types.String `tfsdk:"type"`
	Attributes        types.Set    `tfsdk:"attributes"`
	AlignedEntityName types.String `tfsdk:"aligned_entity_name"`
}

func (r *labelTableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = labelTableTypeName
}

func (r *labelTableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Pai Feature Store Label Table resource. " +
			"A label table is the table that holds the labels used for model training, " +
			"together with the join keys of the training target and the feature entities.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Format: `<instance_id>:<label_table_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore instance, which can be obtained through the ListInstances operation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore project, which can be obtained through the ListProjects operation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"label_table_name": schema.StringAttribute{
				Required: true,
				Description: "The label table name. It can contain only letters, digits and underscores, " +
					"and the maximum length is 64 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_]+$`),
						"must contain only letters, digits and underscores"),
				},
			},
			"datasource_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the datasource where the label table is located, which can be obtained through the ListDatasources operation.",
			},
			"fields": schema.ListAttribute{
				Required: true,
				Description: "The list of fields in the label table. The server cannot add, remove or rename fields after creation, " +
					"so changing a field `name` or `type` replaces the resource; changing `attributes` or " +
					"`aligned_entity_name` of an existing field is applied in place.",
				ElementType: types.ObjectType{AttrTypes: labelTableFieldAttrTypes},
			},
			"label_table_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// ModifyPlan forces replacement when a field identity changes. The API can only
// modify attributes of existing fields, but a flat ListAttribute cannot carry
// per-element plan modifiers; comparing name/type here keeps the in-place update
// path (attributes and aligned_entity_name) while refusing unsupported edits.
func (r *labelTableResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state labelTableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// An unknown fields value cannot be compared with state; replacement is
	// decided once Terraform resolves it.
	if plan.Fields.IsUnknown() || state.Fields.IsUnknown() {
		return
	}

	planFields, diags := labelTableFieldsFromModel(ctx, plan.Fields)
	resp.Diagnostics.Append(diags...)
	stateFields, diags := labelTableFieldsFromModel(ctx, state.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if labelTableFieldIdentityChanged(planFields, stateFields) {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("fields"))
	}
}

func labelTableFieldIdentityChanged(plan, state []labelTableFieldModel) bool {
	if len(plan) != len(state) {
		return true
	}
	for i := range plan {
		if !plan[i].Name.Equal(state[i].Name) || !plan[i].Type.Equal(state[i].Type) {
			return true
		}
	}
	return false
}

func (r *labelTableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan labelTableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fields, diags := labelTableFieldsToBody(ctx, plan.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"Name":         plan.LabelTableName.ValueString(),
		"ProjectId":    plan.ProjectId.ValueString(),
		"DatasourceId": plan.DatasourceId.ValueString(),
		"Fields":       fields,
	}

	pathName := fmt.Sprintf("/api/v1/instances/%s/labeltables", plan.InstanceId.ValueString())
	response, err := r.Client().RoaPost(labelTableProduct, labelTableVersion, pathName, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLabelTable", labelTableTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateLabelTable", labelTableTypeName).Error(),
		)
		return
	}

	labelTableId := labelTableStringFromResponse(response, "LabelTableId")
	if labelTableId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLabelTable", labelTableTypeName),
			fmt.Sprintf("The CreateLabelTable response carried no LabelTableId: %v", response),
		)
		return
	}

	object, err := getLabelTable(r.Client(), plan.InstanceId.ValueString(), labelTableId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetLabelTable", labelTableTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetLabelTable", labelTableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLabelTable", labelTableTypeName),
			fmt.Sprintf("The label table %s is not visible via GetLabelTable right after creation", labelTableId),
		)
		return
	}

	state, diags := labelTableState(ctx, object, plan.InstanceId.ValueString(), labelTableId, r.Client().RegionId, plan.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *labelTableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state labelTableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, labelTableId, ok := parseLabelTableID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", labelTableTypeName),
			fmt.Sprintf("The resource id %q is not in the expected <instance_id>:<label_table_id> format", state.Id.ValueString()),
		)
		return
	}

	object, err := getLabelTable(r.Client(), instanceId, labelTableId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetLabelTable", labelTableTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetLabelTable", labelTableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState, diags := labelTableState(ctx, object, instanceId, labelTableId, r.Client().RegionId, state.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *labelTableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state labelTableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, labelTableId, ok := parseLabelTableID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", labelTableTypeName),
			fmt.Sprintf("The resource id %q is not in the expected <instance_id>:<label_table_id> format", state.Id.ValueString()),
		)
		return
	}

	fields, diags := labelTableFieldsToBody(ctx, plan.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"Name":         plan.LabelTableName.ValueString(),
		"DatasourceId": plan.DatasourceId.ValueString(),
		"Fields":       fields,
	}

	pathName := fmt.Sprintf("/api/v1/instances/%s/labeltables/%s", instanceId, labelTableId)
	if _, err := r.Client().RoaPut(labelTableProduct, labelTableVersion, pathName, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateLabelTable", labelTableTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateLabelTable", labelTableTypeName).Error(),
		)
		return
	}

	object, err := getLabelTable(r.Client(), instanceId, labelTableId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetLabelTable", labelTableTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetLabelTable", labelTableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState, diags := labelTableState(ctx, object, instanceId, labelTableId, r.Client().RegionId, plan.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *labelTableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state labelTableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, labelTableId, ok := parseLabelTableID(state.Id.ValueString())
	if !ok {
		// An unparsable id cannot be deleted server-side; dropping local state is
		// the only sane action, mirroring Read's treatment of a malformed id.
		return
	}

	pathName := fmt.Sprintf("/api/v1/instances/%s/labeltables/%s", instanceId, labelTableId)
	_, err := r.Client().RoaDelete(labelTableProduct, labelTableVersion, pathName, nil, nil, nil, true)
	if err != nil && !isLabelTableNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteLabelTable", labelTableTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteLabelTable", labelTableTypeName).Error(),
		)
		return
	}
}

func (r *labelTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getLabelTable calls GetLabelTable and returns the label table object, or nil
// when it no longer exists. The API reports a missing resource with the
// DoConfigCenterRequestErrorProblem code declared by the CloudSpec
// resourceNotExistCondition; the code arrives on the *tea.SDKError returned by
// RoaGet, so only that typed code is treated as not-found.
func getLabelTable(client *connectivity.AliyunClient, instanceId, labelTableId string) (map[string]interface{}, error) {
	pathName := fmt.Sprintf("/api/v1/instances/%s/labeltables/%s", instanceId, labelTableId)
	response, err := client.RoaGet(labelTableProduct, labelTableVersion, pathName, nil, nil, nil)
	if err != nil {
		if isLabelTableNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("the GetLabelTable response was empty")
	}
	return response, nil
}

func isLabelTableNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "DoConfigCenterRequestErrorProblem" {
		return true
	}
	return false
}

// labelTableState maps a GetLabelTable response onto the resource model. The
// response carries the resource fields at its root (rootMapping "$"); InstanceId,
// LabelTableId and RegionId are not part of it and are carried in from the
// request. Empty strings become empty values for the required field members so
// the flat ObjectType stays fully populated.
func labelTableState(ctx context.Context, object map[string]interface{}, instanceId, labelTableId, regionId string, referenceFields types.List) (labelTableResourceModel, diag.Diagnostics) {
	fields, diags := labelTableFieldsFromResponse(ctx, object["Fields"], referenceFields)
	return labelTableResourceModel{
		Id:             types.StringValue(fmt.Sprintf("%s:%s", instanceId, labelTableId)),
		InstanceId:     types.StringValue(instanceId),
		ProjectId:      labelTableStringAttr(object, "ProjectId"),
		LabelTableName: labelTableStringAttr(object, "Name"),
		DatasourceId:   labelTableStringAttr(object, "DatasourceId"),
		Fields:         fields,
		LabelTableId:   types.StringValue(labelTableId),
		CreateTime:     labelTableStringAttr(object, "GmtCreateTime"),
		RegionId:       types.StringValue(regionId),
	}, diags
}

func labelTableFieldsFromModel(ctx context.Context, list types.List) ([]labelTableFieldModel, diag.Diagnostics) {
	var fields []labelTableFieldModel
	if list.IsNull() || list.IsUnknown() {
		return fields, nil
	}
	diags := list.ElementsAs(ctx, &fields, false)
	return fields, diags
}

func labelTableFieldsToBody(ctx context.Context, list types.List) ([]map[string]interface{}, diag.Diagnostics) {
	fields, diags := labelTableFieldsFromModel(ctx, list)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]map[string]interface{}, 0, len(fields))
	for _, field := range fields {
		item := map[string]interface{}{
			"Name": field.Name.ValueString(),
			"Type": field.Type.ValueString(),
		}
		attributes := make([]string, 0)
		if !field.Attributes.IsNull() && !field.Attributes.IsUnknown() {
			diags.Append(field.Attributes.ElementsAs(ctx, &attributes, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}
		// The API marks Attributes as required in both CreateLabelTable and
		// UpdateLabelTable and accepts an empty array, so it is always sent; an
		// explicit empty string for AlignedEntityName is what clears an existing
		// alignment.
		item["Attributes"] = attributes
		item["AlignedEntityName"] = field.AlignedEntityName.ValueString()
		out = append(out, item)
	}
	return out, diags
}

func labelTableFieldsFromResponse(ctx context.Context, raw interface{}, reference types.List) (types.List, diag.Diagnostics) {
	objectType := types.ObjectType{AttrTypes: labelTableFieldAttrTypes}
	arr, ok := raw.([]interface{})
	if !ok {
		return types.ListValueMust(objectType, []attr.Value{}), nil
	}

	fields := make([]labelTableFieldModel, 0, len(arr))
	for _, el := range arr {
		field, ok := el.(map[string]interface{})
		if !ok {
			continue
		}
		attributeValues := make([]attr.Value, 0)
		if rawAttributes, ok := field["Attributes"].([]interface{}); ok {
			for _, attribute := range rawAttributes {
				if name, ok := attribute.(string); ok && name != "" {
					attributeValues = append(attributeValues, types.StringValue(name))
				}
			}
		}
		fields = append(fields, labelTableFieldModel{
			Name:              types.StringValue(labelTableStringFromResponse(field, "Name")),
			Type:              types.StringValue(labelTableStringFromResponse(field, "Type")),
			Attributes:        types.SetValueMust(types.StringType, attributeValues),
			AlignedEntityName: types.StringValue(labelTableStringFromResponse(field, "AlignedEntityName")),
		})
	}

	// The reference is the plan (Create/Update) or the prior state (Read), both
	// already decoded by the framework; a decode failure here simply falls back
	// to the response order.
	referenceFields, _ := labelTableFieldsFromModel(ctx, reference)
	fields = labelTableOrderByReference(fields, referenceFields)

	values := make([]attr.Value, 0, len(fields))
	for _, field := range fields {
		values = append(values, types.ObjectValueMust(labelTableFieldAttrTypes, map[string]attr.Value{
			"name":                field.Name,
			"type":                field.Type,
			"attributes":          field.Attributes,
			"aligned_entity_name": field.AlignedEntityName,
		}))
	}
	return types.ListValueMust(objectType, values), nil
}

// labelTableOrderByReference restores the reference order on the server
// response. The field list is positional: ModifyPlan compares plan and state by
// index and replaces the resource when a position holds a different field
// identity, and the API cannot reorder existing fields. Without this
// normalization any server-side reordering would turn into a replacement on
// every plan. Fields the reference does not know about keep their response
// order after the matched ones.
func labelTableOrderByReference(fields, reference []labelTableFieldModel) []labelTableFieldModel {
	if len(reference) == 0 || len(fields) == 0 {
		return fields
	}
	used := make([]bool, len(fields))
	ordered := make([]labelTableFieldModel, 0, len(fields))
	for _, ref := range reference {
		for i, field := range fields {
			if used[i] || !field.Name.Equal(ref.Name) {
				continue
			}
			used[i] = true
			ordered = append(ordered, field)
			break
		}
	}
	for i, field := range fields {
		if !used[i] {
			ordered = append(ordered, field)
		}
	}
	return ordered
}

func labelTableStringAttr(object map[string]interface{}, key string) types.String {
	if v, ok := object[key].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func labelTableStringFromResponse(object map[string]interface{}, key string) string {
	if v, ok := object[key].(string); ok {
		return v
	}
	return ""
}

func parseLabelTableID(id string) (string, string, bool) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}
