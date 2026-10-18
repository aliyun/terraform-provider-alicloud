package apigateway

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Dataset API contract: CloudAPI 2016-07-14. Tags are read through
// ListTagResources (DescribeDatasetInfo does not return Tag) and written through
// TagResources/UntagResources; ModifyDataset does not accept Tag, so tags never
// ride on ModifyDataset. See task AC5.

const datasetTypeName = "alicloud_api_gateway_dataset"
const datasetAPIProduct = "CloudAPI"
const datasetAPIVersion = "2016-07-14"
const datasetTagResourceType = "dataset"

var (
	_ resource.Resource                = &datasetResource{}
	_ resource.ResourceWithImportState = &datasetResource{}
)

func NewDatasetResource() resource.Resource {
	return &datasetResource{}
}

type datasetResource struct {
	fwadapt.ResourceBase
}

type datasetResourceModel struct {
	Id           types.String `tfsdk:"id"`
	DatasetId    types.String `tfsdk:"dataset_id"`
	DatasetName  types.String `tfsdk:"dataset_name"`
	DatasetType  types.String `tfsdk:"dataset_type"`
	Description  types.String `tfsdk:"description"`
	CreateTime   types.String `tfsdk:"create_time"`
	ModifiedTime types.String `tfsdk:"modified_time"`
	Tags         types.Map    `tfsdk:"tags"`
}

func (r *datasetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = datasetTypeName
}

func (r *datasetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an API Gateway Dataset resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `dataset_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dataset_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the dataset.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dataset_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the dataset.",
			},
			"dataset_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the dataset. Valid values: `JWT_BLOCKING`, `IP_WHITELIST_CIDR`, `PARAMETER_ACCESS`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("JWT_BLOCKING", "IP_WHITELIST_CIDR", "PARAMETER_ACCESS"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the dataset.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the dataset was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the dataset was last modified.",
				// No UseStateForUnknown: ModifyDataset bumps the timestamp, so the
				// planned state value can never survive an apply.
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The tags of the dataset.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *datasetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan datasetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"DatasetName": plan.DatasetName.ValueString(),
		"DatasetType": plan.DatasetType.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request["Description"] = plan.Description.ValueString()
	}

	response, err := r.Client().RpcPost(datasetAPIProduct, datasetAPIVersion, "CreateDataset", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDataset", datasetTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDataset", datasetTypeName).Error(),
		)
		return
	}

	id := fmt.Sprint(response["DatasetId"])
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDataset", datasetTypeName),
			fmt.Sprintf("The response carried no DatasetId: %v", response),
		)
		return
	}

	plannedTags := tagsToMap(plan.Tags)
	if len(plannedTags) > 0 {
		if err := setDatasetTags(r.Client(), id, plannedTags, nil); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: tagging", datasetTypeName),
				err.Error(),
			)
			return
		}
	}

	model, found, err := readDatasetState(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDatasetInfo", datasetTypeName),
			err.Error(),
		)
		return
	}
	if !found {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDataset", datasetTypeName),
			fmt.Sprintf("The dataset %s is not visible via DescribeDatasetInfo right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *datasetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state datasetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, found, err := readDatasetState(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDatasetInfo", datasetTypeName),
			err.Error(),
		)
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *datasetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state datasetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.DatasetName.Equal(state.DatasetName) || !plan.Description.Equal(state.Description) {
		request := map[string]interface{}{
			"DatasetId":   state.Id.ValueString(),
			"DatasetName": plan.DatasetName.ValueString(),
		}
		if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
			request["Description"] = plan.Description.ValueString()
		}
		if _, err := r.Client().RpcPost(datasetAPIProduct, datasetAPIVersion, "ModifyDataset", nil, request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyDataset", datasetTypeName),
				errs.WrapErrorf(err, "Updating %s: calling ModifyDataset", datasetTypeName).Error(),
			)
			return
		}
	}

	if !plan.Tags.Equal(state.Tags) {
		if err := setDatasetTags(r.Client(), state.Id.ValueString(), tagsToMap(plan.Tags), tagsToMap(state.Tags)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: tagging", datasetTypeName),
				err.Error(),
			)
			return
		}
	}

	model, found, err := readDatasetState(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDatasetInfo", datasetTypeName),
			err.Error(),
		)
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *datasetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state datasetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RpcPost(datasetAPIProduct, datasetAPIVersion, "DeleteDataset", nil, map[string]interface{}{
		"DatasetId": state.Id.ValueString(),
	}, true)
	if err != nil && !errs.IsExpectedErrors(err, []string{"NotFoundDataset"}) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDataset", datasetTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteDataset", datasetTypeName).Error(),
		)
		return
	}
}

func (r *datasetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeDataset calls DescribeDatasetInfo and returns the DatasetInfo object,
// or a nil map with no error when the dataset no longer exists.
func describeDataset(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := client.RpcPost(datasetAPIProduct, datasetAPIVersion, "DescribeDatasetInfo", nil, map[string]interface{}{
		"DatasetId": id,
	}, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundDataset"}) {
			return nil, nil
		}
		return nil, errs.WrapErrorf(err, "Reading %s: calling DescribeDatasetInfo", datasetTypeName)
	}
	object, _ := response["DatasetInfo"].(map[string]interface{})
	if object == nil {
		return nil, fmt.Errorf("the DescribeDatasetInfo response carried no DatasetInfo: %v", response)
	}
	return object, nil
}

// listDatasetTags reads the dataset's tags through ListTagResources. Tags are
// never returned by DescribeDatasetInfo, so this is the only read path.
func listDatasetTags(client *connectivity.AliyunClient, id string) map[string]string {
	query := map[string]interface{}{
		"ResourceId.1": id,
		"ResourceType": datasetTagResourceType,
	}
	response, err := client.RpcPost(datasetAPIProduct, datasetAPIVersion, "ListTagResources", query, map[string]interface{}{}, true)
	if err != nil {
		return nil
	}
	tagsRaw, _ := response["TagResources"].(map[string]interface{})
	tagList, _ := tagsRaw["TagResource"].([]interface{})
	out := map[string]string{}
	for _, item := range tagList {
		t, _ := item.(map[string]interface{})
		if t == nil {
			continue
		}
		key := tagValue(t, "TagKey", "Key")
		val := tagValue(t, "TagValue", "Value")
		if key == "" {
			continue
		}
		out[key] = val
	}
	return out
}

func tagValue(t map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := t[k].(string); ok {
			return v
		}
	}
	return ""
}

// setDatasetTags reconciles the planned tag set against the existing one:
// UntagResources for removed keys, TagResources for added/changed keys.
func setDatasetTags(client *connectivity.AliyunClient, id string, planned, existing map[string]string) error {
	removed := make([]string, 0)
	for key := range existing {
		if _, ok := planned[key]; !ok {
			removed = append(removed, key)
		}
	}
	added := map[string]string{}
	for key, value := range planned {
		old, ok := existing[key]
		if !ok || old != value {
			added[key] = value
		}
	}

	if len(removed) > 0 {
		request := map[string]interface{}{
			"ResourceId.1": id,
			"ResourceType": datasetTagResourceType,
		}
		for i, key := range removed {
			request[fmt.Sprintf("TagKey.%d", i+1)] = key
		}
		if _, err := client.RpcPost(datasetAPIProduct, datasetAPIVersion, "UntagResources", nil, request, true); err != nil {
			return errs.WrapErrorf(err, "Updating %s: calling UntagResources", datasetTypeName)
		}
	}

	if len(added) > 0 {
		request := map[string]interface{}{
			"ResourceId.1": id,
			"ResourceType": datasetTagResourceType,
		}
		count := 1
		for key, value := range added {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
		if _, err := client.RpcPost(datasetAPIProduct, datasetAPIVersion, "TagResources", nil, request, true); err != nil {
			return errs.WrapErrorf(err, "Updating %s: calling TagResources", datasetTypeName)
		}
	}
	return nil
}

// readDatasetState assembles the resource model from DescribeDatasetInfo and
// ListTagResources. A false result with no error means the dataset is gone.
func readDatasetState(client *connectivity.AliyunClient, id string) (datasetResourceModel, bool, error) {
	object, err := describeDataset(client, id)
	if err != nil {
		return datasetResourceModel{}, false, err
	}
	if object == nil {
		return datasetResourceModel{}, false, nil
	}
	return datasetResourceModel{
		Id:           types.StringValue(id),
		DatasetId:    strVal(object["DatasetId"], id),
		DatasetName:  strVal(object["DatasetName"], ""),
		DatasetType:  strVal(object["DatasetType"], ""),
		Description:  strVal(object["Description"], ""),
		CreateTime:   strVal(object["CreatedTime"], ""),
		ModifiedTime: strVal(object["ModifiedTime"], ""),
		Tags:         tagsToType(listDatasetTags(client, id)),
	}, true, nil
}

// strVal renders an API field as a Framework string. An absent or empty value
// becomes null, matching the legacy behavior where dates and optional fields
// the API omits read back as absent rather than "".
func strVal(v interface{}, fallback string) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	if fallback != "" {
		return types.StringValue(fallback)
	}
	return types.StringNull()
}

// tagsToMap flattens a Framework tags map to a plain string map. A null or
// unknown map (no tags in config) yields an empty map, making setDatasetTags a
// no-op on create.
func tagsToMap(m types.Map) map[string]string {
	out := map[string]string{}
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	for key, value := range m.Elements() {
		if s, ok := value.(types.String); ok {
			out[key] = s.ValueString()
		}
	}
	return out
}

// tagsToType renders a plain string map as a Framework tags map. An empty map
// becomes null, so a dataset with no tags reads back as tags = null rather than
// tags = {}.
func tagsToType(tags map[string]string) types.Map {
	if len(tags) == 0 {
		return types.MapNull(types.StringType)
	}
	values := make(map[string]attr.Value, len(tags))
	for key, value := range tags {
		values[key] = types.StringValue(value)
	}
	m, diags := types.MapValue(types.StringType, values)
	_ = diags
	return m
}
