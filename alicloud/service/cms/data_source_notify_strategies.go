package cms

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const notifyStrategiesDataSourceTypeName = "alicloud_cms_notify_strategies"

// pageSizeLarge mirrors the PageSizeLarge constant the SDKv2 data source used
// for the maxResults query parameter.
const pageSizeLarge = 50

var _ datasource.DataSource = &notifyStrategiesDataSource{}

func NewNotifyStrategiesDataSource() datasource.DataSource {
	return &notifyStrategiesDataSource{}
}

type notifyStrategiesDataSource struct {
	fwadapt.DataSourceBase
}

type notifyStrategiesDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	Ids              types.List   `tfsdk:"ids"`
	NameRegex        types.String `tfsdk:"name_regex"`
	Workspace        types.String `tfsdk:"workspace"`
	Names            types.List   `tfsdk:"names"`
	NotifyStrategies types.List   `tfsdk:"notify_strategies"`
	OutputFile       types.String `tfsdk:"output_file"`
}

// notifyStrategyElementAttrTypes is the attribute-type map of a single element
// of the notify_strategies list. It mirrors the resource schema and is reused
// by the resource flatten helpers.
var notifyStrategyElementAttrTypes = map[string]attr.Type{
	"id":                           types.StringType,
	"create_time":                  types.StringType,
	"custom_template_entries":      types.ListType{ElemType: types.ObjectType{AttrTypes: customTemplateEntryAttrTypes}},
	"description":                  types.StringType,
	"enable":                       types.BoolType,
	"grouping_setting":             types.ObjectType{AttrTypes: groupingSettingAttrTypes},
	"ignore_restored_notification": types.BoolType,
	"notify_strategy_id":           types.StringType,
	"notify_strategy_name":         types.StringType,
	"routes":                       types.ListType{ElemType: types.ObjectType{AttrTypes: routeAttrTypes}},
	"update_time":                  types.StringType,
	"user_id":                      types.StringType,
	"workspace":                    types.StringType,
}

func (d *notifyStrategiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = notifyStrategiesDataSourceTypeName
}

func (d *notifyStrategiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Cms Notify Strategies of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The provider-assigned unique ID of the data source.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of notify strategy IDs.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by the notify strategy name.",
			},
			"workspace": schema.StringAttribute{
				Optional:    true,
				Description: "The workspace id of the notify strategies.",
			},
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of notify strategy names.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File path where to write the result data in JSON format.",
			},
			"notify_strategies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the notify strategy.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "The time when the notify strategy was created.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the notify strategy.",
						},
						"enable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the notify strategy is enabled.",
						},
						"ignore_restored_notification": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to ignore restored notifications.",
						},
						"notify_strategy_id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the notify strategy.",
						},
						"notify_strategy_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the notify strategy.",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "The time when the notify strategy was last updated.",
						},
						"user_id": schema.StringAttribute{
							Computed:    true,
							Description: "The user id of the notify strategy.",
						},
						"workspace": schema.StringAttribute{
							Computed:    true,
							Description: "The workspace id of the notify strategy.",
						},
						"grouping_setting": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"grouping_keys": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
								"period_min": schema.Int64Attribute{
									Computed: true,
								},
								"silence_sec": schema.Int64Attribute{
									Computed: true,
								},
								"times": schema.Int64Attribute{
									Computed: true,
								},
							},
						},
						"routes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"severities": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
									"channels": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"channel_type": schema.StringAttribute{
													Computed: true,
												},
												"enabled_sub_channels": schema.ListAttribute{
													ElementType: types.StringType,
													Computed:    true,
												},
												"receivers": schema.ListAttribute{
													ElementType: types.StringType,
													Computed:    true,
												},
											},
										},
									},
									"effect_time_range": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"day_in_week": schema.ListAttribute{
												ElementType: types.Int64Type,
												Computed:    true,
											},
											"end_time_in_minute": schema.Int64Attribute{
												Computed: true,
											},
											"start_time_in_minute": schema.Int64Attribute{
												Computed: true,
											},
											"time_zone": schema.StringAttribute{
												Computed: true,
											},
										},
									},
									"filter_setting": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"expression": schema.StringAttribute{
												Computed: true,
											},
											"relation": schema.StringAttribute{
												Computed: true,
											},
											"conditions": schema.ListNestedAttribute{
												Computed: true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"field": schema.StringAttribute{
															Computed: true,
														},
														"op": schema.StringAttribute{
															Computed: true,
														},
														"value": schema.StringAttribute{
															Computed: true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"custom_template_entries": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"target_type": schema.StringAttribute{
										Computed: true,
									},
									"template_uuid": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *notifyStrategiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config notifyStrategiesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	var nameRegex *regexp.Regexp
	if v := config.NameRegex.ValueString(); v != "" {
		r, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", notifyStrategiesDataSourceTypeName),
				errs.WrapErrorf(err, "Reading %s: compiling name_regex %q", notifyStrategiesDataSourceTypeName, v).Error(),
			)
			return
		}
		nameRegex = r
	}

	idsSet := make(map[string]struct{})
	for _, el := range config.Ids.Elements() {
		if s, ok := el.(types.String); ok && s.ValueString() != "" {
			idsSet[s.ValueString()] = struct{}{}
		}
	}

	var objects []map[string]interface{}
	nextToken := ""
	for {
		query := map[string]*string{
			"RegionId":   stringPtr(client.RegionId),
			"maxResults": stringPtr(strconv.Itoa(pageSizeLarge)),
		}
		if v := config.Workspace.ValueString(); v != "" {
			query["workspace"] = stringPtr(v)
		}
		if nextToken != "" {
			query["nextToken"] = stringPtr(nextToken)
		}
		response, err := client.RoaGet("Cms", "2024-03-30", "/notifyStrategies", query, nil, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling ListNotifyStrategies", notifyStrategiesDataSourceTypeName),
				errs.WrapErrorf(err, "Reading %s: calling ListNotifyStrategies", notifyStrategiesDataSourceTypeName).Error(),
			)
			return
		}

		dataList, _ := response["dataList"].([]interface{})
		for _, item := range dataList {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(m["notifyStrategyName"])) {
				continue
			}
			if len(idsSet) > 0 {
				if _, ok := idsSet[fmt.Sprint(m["notifyStrategyId"])]; !ok {
					continue
				}
			}
			objects = append(objects, m)
		}

		if nt, ok := response["nextToken"].(string); ok && nt != "" {
			nextToken = nt
		} else {
			break
		}
	}

	ids := make([]string, 0, len(objects))
	names := make([]attr.Value, 0, len(objects))
	elements := make([]attr.Value, 0, len(objects))
	for _, m := range objects {
		elements = append(elements, flattenNotifyStrategyObject(m))
		idStr := fmt.Sprint(m["notifyStrategyId"])
		ids = append(ids, idStr)
		names = append(names, flattenString(m["notifyStrategyName"]))
	}

	state := notifyStrategiesDataSourceModel{
		Id:               types.StringValue(notifyStrategiesIDHash(ids)),
		Ids:              stringListValue(ids),
		Names:            stringListFromValues(names),
		NotifyStrategies: listValueFromObjects(notifyStrategyElementAttrTypes, elements),
		OutputFile:       config.OutputFile,
	}

	if v := config.OutputFile.ValueString(); v != "" {
		if err := writeNotifyStrategiesToFile(v, objects); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output_file", notifyStrategiesDataSourceTypeName),
				errs.WrapErrorf(err, "Reading %s: writing output_file %q", notifyStrategiesDataSourceTypeName, v).Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenNotifyStrategyObject(m map[string]interface{}) types.Object {
	values := map[string]attr.Value{
		"id":                           flattenString(m["notifyStrategyId"]),
		"create_time":                  flattenString(m["createTime"]),
		"description":                  flattenString(m["description"]),
		"enable":                       flattenBool(m["enabled"]),
		"ignore_restored_notification": flattenBool(m["ignoreRestoredNotification"]),
		"notify_strategy_id":           flattenString(m["notifyStrategyId"]),
		"notify_strategy_name":         flattenString(m["notifyStrategyName"]),
		"update_time":                  flattenString(m["updateTime"]),
		"user_id":                      flattenString(m["userId"]),
		"workspace":                    flattenString(m["workspace"]),
		"grouping_setting":             flattenGroupingSetting(m["groupingSetting"]),
		"routes":                       flattenRoutes(m["routes"]),
		"custom_template_entries":      flattenCustomTemplateEntries(m["customTemplateEntries"]),
	}
	obj, _ := types.ObjectValue(notifyStrategyElementAttrTypes, values)
	return obj
}

func flattenBool(v interface{}) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	if b, ok := v.(bool); ok {
		return types.BoolValue(b)
	}
	return types.BoolNull()
}

func stringListValue(ss []string) types.List {
	vals := make([]attr.Value, 0, len(ss))
	for _, s := range ss {
		if s != "" {
			vals = append(vals, types.StringValue(s))
		}
	}
	list, _ := types.ListValue(types.StringType, vals)
	return list
}

func stringListFromValues(vals []attr.Value) types.List {
	list, _ := types.ListValue(types.StringType, vals)
	return list
}

func listValueFromObjects(attrTypes map[string]attr.Type, elements []attr.Value) types.List {
	if len(elements) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, elements)
	return list
}

func notifyStrategiesIDHash(ids []string) string {
	h := sha1.New()
	for _, id := range ids {
		h.Write([]byte(id))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func writeNotifyStrategiesToFile(path string, objects []map[string]interface{}) error {
	data, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
