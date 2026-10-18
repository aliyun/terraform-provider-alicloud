package cms

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const alertMetricGroupsTypeName = "alicloud_cms_alert_metric_groups"

var (
	_ datasource.DataSource              = &alertMetricGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &alertMetricGroupsDataSource{}
)

func NewAlertMetricGroupsDataSource() datasource.DataSource {
	return &alertMetricGroupsDataSource{}
}

type alertMetricGroupsDataSource struct {
	fwadapt.DataSourceBase
}

// The muxed provider only exposes protocol v5, under which framework nested
// attributes (ListNestedAttribute / NestedAttributeObject) cannot be served.
// The nested groups structure is therefore modelled as ListAttribute with
// ObjectType elements, mirroring the SDK v2 TypeList + Resource shape that v5
// supports. Leaf attributes remain primitive attributes.

var (
	supportedOptObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"display_name_cn": types.StringType,
			"display_name_en": types.StringType,
			"value":           types.StringType,
		},
	}
	paramValueObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"label_cn": types.StringType,
			"label_en": types.StringType,
			"value":    types.StringType,
		},
	}
	filterObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"dim":             types.StringType,
			"display_name_cn": types.StringType,
			"display_name_en": types.StringType,
			"hidden":          types.BoolType,
			"opt":             types.StringType,
			"label_disabled":  types.BoolType,
			"dim_disabled":    types.BoolType,
			"supported_opts":  types.ListType{ElemType: supportedOptObjectType},
		},
	}
	paramObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"max_width":      types.Int64Type,
			"min_width":      types.Int64Type,
			"name":           types.StringType,
			"placeholder_cn": types.StringType,
			"placeholder_en": types.StringType,
			"type":           types.StringType,
			"value":          types.StringType,
			"values":         types.ListType{ElemType: paramValueObjectType},
		},
	}
	groupObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                    types.StringType,
			"alert_metric_group_id": types.StringType,
			"datasource_types":      types.StringType,
			"display_name_cn":       types.StringType,
			"display_name_en":       types.StringType,
			"description_cn":        types.StringType,
			"description_en":        types.StringType,
			"order_index":           types.Int64Type,
			"filters":               types.ListType{ElemType: filterObjectType},
			"params":                types.ListType{ElemType: paramObjectType},
		},
	}
)

// alertMetricGroupsModel is the top-level state shape of the data source.
type alertMetricGroupsModel struct {
	Ids            types.List   `tfsdk:"ids"`
	DatasourceType types.String `tfsdk:"datasource_type"`
	IncludeDetails types.Bool   `tfsdk:"include_details"`
	OutputFile     types.String `tfsdk:"output_file"`
	Groups         types.List   `tfsdk:"groups"`
	Id             types.String `tfsdk:"id"`
}

func (d *alertMetricGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = alertMetricGroupsTypeName
}

func (d *alertMetricGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Cms Alert Metric Groups of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, a hash of the listed alert metric group ids.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of alert metric group IDs to filter the results. When set, only groups whose `alert_metric_group_id` is in this list are returned.",
			},
			"datasource_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the data source. When set, the server filters groups by this value.",
			},
			"include_details": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to include detail fields (filters/params) in the response. Defaults to false.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "When set, the listed groups are also persisted to this file path.",
			},
			"groups": schema.ListAttribute{
				ElementType: groupObjectType,
				Computed:    true,
				Description: "A list of Alert Metric Group entries.",
			},
		},
	}
}

func (d *alertMetricGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config alertMetricGroupsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Client-side ids filter: when `ids` is set, only groups whose id is in
	// this set are kept.
	idsMap := map[string]string{}
	if !config.Ids.IsNull() && len(config.Ids.Elements()) > 0 {
		for _, el := range config.Ids.Elements() {
			if s, ok := el.(types.String); ok && s.ValueString() != "" {
				idsMap[s.ValueString()] = s.ValueString()
			}
		}
	}

	svc := NewService(d.Client())
	objects, err := svc.ListAlertMetricGroups(config.DatasourceType.ValueString(), config.IncludeDetails.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", alertMetricGroupsTypeName),
			err.Error(),
		)
		return
	}

	groupObjs := make([]attr.Value, 0)
	objectsFiltered := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, raw := range objects {
		groupID := strVal(raw, "alertMetricGroupId")
		if len(idsMap) > 0 {
			if _, ok := idsMap[groupID]; !ok {
				continue
			}
		}
		groupObjs = append(groupObjs, buildGroupObject(raw))
		objectsFiltered = append(objectsFiltered, raw)
		ids = append(ids, groupID)
	}

	state := config
	state.Groups = listValueFromObjects(groupObjectType, groupObjs)
	state.Ids = stringSliceToListValue(ids)
	state.Id = types.StringValue(helper.DataResourceIdHash(ids))

	if outputFile := config.OutputFile.ValueString(); outputFile != "" {
		if err := helper.WriteToFile(outputFile, objectsFiltered); err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: writing output file", alertMetricGroupsTypeName),
				err.Error(),
			)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func buildGroupObject(raw map[string]interface{}) types.Object {
	obj, diags := types.ObjectValue(groupObjectType.AttrTypes, map[string]attr.Value{
		"id":                    types.StringValue(strVal(raw, "alertMetricGroupId")),
		"alert_metric_group_id": types.StringValue(strVal(raw, "alertMetricGroupId")),
		"datasource_types":      types.StringValue(strVal(raw, "datasourceTypes")),
		"display_name_cn":       types.StringValue(strVal(raw, "displayNameCn")),
		"display_name_en":       types.StringValue(strVal(raw, "displayNameEn")),
		"description_cn":        types.StringValue(strVal(raw, "descriptionCn")),
		"description_en":        types.StringValue(strVal(raw, "descriptionEn")),
		"order_index":           types.Int64Value(int64Val(raw, "orderIndex")),
		"filters":               buildFilterList(raw["filters"]),
		"params":                buildParamList(raw["params"]),
	})
	if diags.HasError() {
		obj, _ = types.ObjectValue(groupObjectType.AttrTypes, map[string]attr.Value{})
	}
	return obj
}

func buildFilterList(raw interface{}) types.List {
	objs := make([]attr.Value, 0)
	if raw == nil {
		return listValueFromObjects(filterObjectType, objs)
	}
	for _, item := range convert.ConvertToInterfaceArray(raw) {
		child, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		obj, diags := types.ObjectValue(filterObjectType.AttrTypes, map[string]attr.Value{
			"dim":             types.StringValue(strVal(child, "dim")),
			"display_name_cn": types.StringValue(strVal(child, "displayNameCn")),
			"display_name_en": types.StringValue(strVal(child, "displayNameEn")),
			"hidden":          types.BoolValue(boolVal(child, "hidden")),
			"opt":             types.StringValue(strVal(child, "opt")),
			"label_disabled":  types.BoolValue(boolVal(child, "labelDisabled")),
			"dim_disabled":    types.BoolValue(boolVal(child, "dimDisabled")),
			"supported_opts":  buildSupportedOptsList(child["supportedOpts"]),
		})
		if diags.HasError() {
			continue
		}
		objs = append(objs, obj)
	}
	return listValueFromObjects(filterObjectType, objs)
}

func buildSupportedOptsList(raw interface{}) types.List {
	objs := make([]attr.Value, 0)
	if raw == nil {
		return listValueFromObjects(supportedOptObjectType, objs)
	}
	for _, item := range convert.ConvertToInterfaceArray(raw) {
		child, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		obj, diags := types.ObjectValue(supportedOptObjectType.AttrTypes, map[string]attr.Value{
			"display_name_cn": types.StringValue(strVal(child, "displayNameCn")),
			"display_name_en": types.StringValue(strVal(child, "displayNameEn")),
			"value":           types.StringValue(strVal(child, "value")),
		})
		if diags.HasError() {
			continue
		}
		objs = append(objs, obj)
	}
	return listValueFromObjects(supportedOptObjectType, objs)
}

func buildParamList(raw interface{}) types.List {
	objs := make([]attr.Value, 0)
	if raw == nil {
		return listValueFromObjects(paramObjectType, objs)
	}
	for _, item := range convert.ConvertToInterfaceArray(raw) {
		child, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		obj, diags := types.ObjectValue(paramObjectType.AttrTypes, map[string]attr.Value{
			"max_width":      types.Int64Value(int64Val(child, "maxWidth")),
			"min_width":      types.Int64Value(int64Val(child, "minWidth")),
			"name":           types.StringValue(strVal(child, "name")),
			"placeholder_cn": types.StringValue(strVal(child, "placeholderCn")),
			"placeholder_en": types.StringValue(strVal(child, "placeholderEn")),
			"type":           types.StringValue(strVal(child, "type")),
			"value":          types.StringValue(strVal(child, "value")),
			"values":         buildParamValuesList(child["values"]),
		})
		if diags.HasError() {
			continue
		}
		objs = append(objs, obj)
	}
	return listValueFromObjects(paramObjectType, objs)
}

func buildParamValuesList(raw interface{}) types.List {
	objs := make([]attr.Value, 0)
	if raw == nil {
		return listValueFromObjects(paramValueObjectType, objs)
	}
	for _, item := range convert.ConvertToInterfaceArray(raw) {
		child, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		obj, diags := types.ObjectValue(paramValueObjectType.AttrTypes, map[string]attr.Value{
			"label_cn": types.StringValue(strVal(child, "labelCn")),
			"label_en": types.StringValue(strVal(child, "labelEn")),
			"value":    types.StringValue(strVal(child, "value")),
		})
		if diags.HasError() {
			continue
		}
		objs = append(objs, obj)
	}
	return listValueFromObjects(paramValueObjectType, objs)
}

// listValueFromObjects builds a non-null framework list of objects, returning
// an empty list when there are no elements.
func listValueFromObjects(elemType attr.Type, objs []attr.Value) types.List {
	list, diags := types.ListValue(elemType, objs)
	if diags.HasError() {
		return types.ListNull(elemType)
	}
	return list
}

// strVal renders an interface value (string, number or nil) as a string,
// mirroring fmt.Sprint with a nil-safe empty result.
func strVal(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// boolVal parses a boolean from an interface that may arrive as a Go bool or
// as a "true"/"false" string from the API.
func boolVal(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	b, _ := strconv.ParseBool(fmt.Sprint(v))
	return b
}

// int64Val parses an integer from an interface that may arrive as a JSON number
// (float64) or as a numeric string.
func int64Val(m map[string]interface{}, key string) int64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int:
		return int64(n)
	case int64:
		return n
	}
	n, _ := strconv.ParseInt(fmt.Sprint(v), 10, 64)
	return n
}

// stringSliceToListValue converts a []string into a framework list of strings,
// returning an empty (non-null) list when there are no elements.
func stringSliceToListValue(values []string) types.List {
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	list, diags := types.ListValue(types.StringType, elems)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return list
}
