package threatdetection

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const customCheckItemsTypeName = "alicloud_threat_detection_custom_check_items"

var (
	_ datasource.DataSource              = &customCheckItemsDataSource{}
	_ datasource.DataSourceWithConfigure = &customCheckItemsDataSource{}

	// descListType is the list type of a description / assist_info / solution
	// block, matching descObjectType in resource_custom_check_item.go.
	descListType = types.ListType{ElemType: descObjectType}

	// itemsObjectType is the element type of the items output list.
	itemsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"check_id":          types.Int64Type,
		"check_show_name":   types.StringType,
		"vendor":            types.StringType,
		"instance_type":     types.StringType,
		"instance_sub_type": types.StringType,
		"risk_level":        types.StringType,
		"status":            types.StringType,
		"check_rule":        types.StringType,
		"check_type":        types.StringType,
		"remark":            types.StringType,
		"description":       descListType,
		"assist_info":       descListType,
		"solution":          descListType,
	}}
)

func NewCustomCheckItemsDataSource() datasource.DataSource {
	return &customCheckItemsDataSource{}
}

type customCheckItemsDataSource struct {
	fwadapt.DataSourceBase
}

type customCheckItemsModel struct {
	Id            types.String `tfsdk:"id"`
	Ids           types.List   `tfsdk:"ids"`
	CheckId       types.Int64  `tfsdk:"check_id"`
	CheckShowName types.String `tfsdk:"check_show_name"`
	CheckTypes    types.List   `tfsdk:"check_types"`
	Lang          types.String `tfsdk:"lang"`
	Statuses      types.List   `tfsdk:"statuses"`
	PageSize      types.Int64  `tfsdk:"page_size"`
	CurrentPage   types.Int64  `tfsdk:"current_page"`
	OutputFile    types.String `tfsdk:"output_file"`
	Items         types.List   `tfsdk:"items"`
}

func (d *customCheckItemsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = customCheckItemsTypeName
}

func (d *customCheckItemsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the custom check items of the current Alibaba Cloud account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source. It is a hash of the returned check item IDs.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of Custom Check Item IDs.",
			},
			"check_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the custom check item.",
			},
			"check_show_name": schema.StringAttribute{
				Optional:    true,
				Description: "The display name of the check item.",
			},
			"check_types": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The types of the check item to query.",
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Description: "The language of the check item.",
			},
			"statuses": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The statuses of the check item to query.",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The page size of the list query. Defaults to 50.",
			},
			"current_page": schema.Int64Attribute{
				Optional:    true,
				Description: "The current page number of the list query. Defaults to 1.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "Save the result to a file.",
			},
		},
		// items is a top-level ListNestedBlock rather than a ListNestedAttribute
		// because the provider rides the tf5 mux server (protocol v5), and
		// protocol v5 cannot represent nested attributes (Attributes set on a
		// nested object). Blocks are the v5-native nested structure, which also
		// matches the SDK v2 Elem-based shape of the original data source.
		Blocks: map[string]schema.Block{
			"items": schema.ListNestedBlock{
				Description: "A list of custom check items.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"check_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The ID of the custom check item.",
						},
						"check_show_name": schema.StringAttribute{
							Computed:    true,
							Description: "The display name of the check item.",
						},
						"vendor": schema.StringAttribute{
							Computed:    true,
							Description: "The vendor to which the check item belongs.",
						},
						"instance_type": schema.StringAttribute{
							Computed:    true,
							Description: "The asset type to which the check item belongs.",
						},
						"instance_sub_type": schema.StringAttribute{
							Computed:    true,
							Description: "The asset subtype to which the check item belongs.",
						},
						"risk_level": schema.StringAttribute{
							Computed:    true,
							Description: "The risk level of the check item.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The status of the check item.",
						},
						"check_rule": schema.StringAttribute{
							Computed:    true,
							Description: "The check rule of the check item.",
						},
						"check_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the check item.",
						},
						"remark": schema.StringAttribute{
							Computed:    true,
							Description: "The remarks of the check item.",
						},
					},
					Blocks: map[string]schema.Block{
						"description": schema.ListNestedBlock{
							Description: "The description of the check item.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"type":  schema.StringAttribute{Computed: true},
									"value": schema.StringAttribute{Computed: true},
								},
							},
						},
						"assist_info": schema.ListNestedBlock{
							Description: "The assist info of the check item.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"type":  schema.StringAttribute{Computed: true},
									"value": schema.StringAttribute{Computed: true},
								},
							},
						},
						"solution": schema.ListNestedBlock{
							Description: "The solution of the check item.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"type":  schema.StringAttribute{Computed: true},
									"value": schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *customCheckItemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config customCheckItemsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pageSize := int64(50)
	if !config.PageSize.IsNull() {
		pageSize = config.PageSize.ValueInt64()
	}
	currentPage := int64(1)
	singlePage := false
	if !config.CurrentPage.IsNull() {
		currentPage = config.CurrentPage.ValueInt64()
		singlePage = true
	}

	query := map[string]interface{}{
		"PageSize":    pageSize,
		"CurrentPage": currentPage,
	}
	if !config.CheckId.IsNull() {
		query["CheckId"] = config.CheckId.ValueInt64()
	}
	if !config.CheckShowName.IsNull() {
		query["CheckShowName"] = config.CheckShowName.ValueString()
	}
	if checkTypes := stringSlice(config.CheckTypes); len(checkTypes) > 0 {
		query["CheckTypes"] = toInterfaceSlice(checkTypes)
	}
	if !config.Lang.IsNull() {
		query["Lang"] = config.Lang.ValueString()
	}
	if statuses := stringSlice(config.Statuses); len(statuses) > 0 {
		query["Statuses"] = toInterfaceSlice(statuses)
	}

	var objects []map[string]interface{}
	ids := make([]string, 0)
	for {
		// ListCheckItems is an RPC GET, so the filter parameters travel in the
		// signed query map; the request body stays empty.
		response, err := d.Client().RpcGet("Sas", "2018-12-03", "ListCheckItems", query, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling ListCheckItems", customCheckItemsTypeName),
				errs.WrapErrorf(err, "Reading %s: calling ListCheckItems", customCheckItemsTypeName).Error(),
			)
			return
		}

		items, _ := response["CheckItems"].([]interface{})
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			if item == nil {
				continue
			}
			objects = append(objects, item)
			ids = append(ids, fmt.Sprint(item["CheckId"]))
		}

		if singlePage || int64(len(items)) < pageSize {
			break
		}
		currentPage++
		query["CurrentPage"] = currentPage
	}

	idsList := stringList(ids)
	state := customCheckItemsModel{
		Id:            types.StringValue(dataResourceIdHash(ids)),
		Ids:           idsList,
		CheckId:       config.CheckId,
		CheckShowName: config.CheckShowName,
		CheckTypes:    config.CheckTypes,
		Lang:          config.Lang,
		Statuses:      config.Statuses,
		PageSize:      config.PageSize,
		CurrentPage:   config.CurrentPage,
		OutputFile:    config.OutputFile,
		Items:         buildItemsList(objects),
	}

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		if err := writeToFile(config.OutputFile.ValueString(), objects); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output_file", customCheckItemsTypeName),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// buildItemsList maps the ListCheckItems items onto the items output list.
func buildItemsList(objects []map[string]interface{}) types.List {
	if len(objects) == 0 {
		return types.ListNull(itemsObjectType)
	}
	values := make([]attr.Value, 0, len(objects))
	for _, o := range objects {
		obj, _ := types.ObjectValue(itemsObjectType.AttrTypes, map[string]attr.Value{
			"check_id":          int64Value(o["CheckId"]),
			"check_show_name":   stringValue(o["CheckShowName"]),
			"vendor":            stringValue(o["Vendor"]),
			"instance_type":     stringValue(o["InstanceType"]),
			"instance_sub_type": stringValue(o["InstanceSubType"]),
			"risk_level":        stringValue(o["RiskLevel"]),
			"status":            stringValue(o["Status"]),
			"check_rule":        stringValue(o["CheckRule"]),
			"check_type":        stringValue(o["CheckType"]),
			"remark":            stringValue(o["Remark"]),
			"description":       structList(o["Description"]),
			"assist_info":       structList(o["AssistInfo"]),
			"solution":          structList(o["Solution"]),
		})
		values = append(values, obj)
	}
	l, _ := types.ListValue(itemsObjectType, values)
	return l
}

// stringSlice extracts the string elements of a List of strings.
func stringSlice(list types.List) []string {
	out := make([]string, 0, len(list.Elements()))
	for _, e := range list.Elements() {
		if v, ok := e.(types.String); ok && !v.IsNull() {
			out = append(out, v.ValueString())
		}
	}
	return out
}

// stringList builds a List of strings from a []string.
func stringList(ids []string) types.List {
	if len(ids) == 0 {
		return types.ListNull(types.StringType)
	}
	values := make([]attr.Value, 0, len(ids))
	for _, id := range ids {
		values = append(values, types.StringValue(id))
	}
	l, _ := types.ListValue(types.StringType, values)
	return l
}

// toInterfaceSlice converts []string to []interface{} for RPC parameters.
func toInterfaceSlice(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}

// dataResourceIdHash produces a stable hash of the returned IDs, mirroring the
// SDK v2 helper used by the original data source.
func dataResourceIdHash(ids []string) string {
	var buf []byte
	for _, id := range ids {
		buf = append(buf, id...)
		buf = append(buf, '-')
	}
	return strconv.FormatInt(int64(fnv1aHash(buf)), 10)
}

// fnv1aHash is a small FNV-1a 64-bit hash used for the data source ID.
func fnv1aHash(b []byte) uint64 {
	h := uint64(1469598103934665603)
	for _, c := range b {
		h ^= uint64(c)
		h *= 1099511628211
	}
	return h
}

// writeToFile writes the data source items to a file as indented JSON, mirroring
// the SDK v2 helper without importing package alicloud.
func writeToFile(filePath string, data interface{}) error {
	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
	}
	return os.WriteFile(filePath, bs, 0o644)
}
