package threatdetection

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const componentsTypeName = "alicloud_threat_detection_components"

var (
	_ datasource.DataSource              = &componentsDataSource{}
	_ datasource.DataSourceWithConfigure = &componentsDataSource{}
)

// NewComponentsDataSource is the factory registered with the provider registry; it
// returns a pointer so the fwadapt client-injection guard sees a WithClient.
func NewComponentsDataSource() datasource.DataSource {
	return &componentsDataSource{}
}

type componentsDataSource struct {
	fwadapt.DataSourceBase
}

func (d *componentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = componentsTypeName
}

// Schema uses nested Blocks rather than ListNestedAttribute. The framework
// rewrites the prior SDK v2 *schema.Resource block form; a ListNestedAttribute
// (NestedAttributeObject) is a protocol v6 construct that tf5muxserver cannot
// convert to v5 ("protocol version 5 cannot have Attributes set"), so the
// provider schema would fail to load for every plan/refresh. Blocks convert to
// the same v5 nested-block representation the prior SDK v2 data source exposed,
// keeping the muxed schema byte-compatible while staying Framework-native.
func (d *componentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of Threat Detection (SOAR) Components available to the current Alibaba Cloud account.",
		Attributes: map[string]schema.Attribute{
			"component_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the component. The data source only returns the component that exactly matches this name.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by component name.",
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Description: "The language type of the request and received messages. Valid values: `zh`, `en`.",
			},
			"role_for": schema.Int64Attribute{
				Optional:    true,
				Description: "The account ID of the Resource Directory member.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of component names. The data source returns only the components whose names match the list.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results (after running `terraform plan`).",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, a hash of the component names listed.",
			},
		},
		Blocks: map[string]schema.Block{
			"components": schema.ListNestedBlock{
				Description: "A list of component entries.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                    schema.StringAttribute{Computed: true, Description: "The name of the component. Same as `component_name`."},
						"component_name":        schema.StringAttribute{Computed: true, Description: "The name of the component."},
						"component_alias":       schema.StringAttribute{Computed: true, Description: "The alias of the component, used for display."},
						"component_description": schema.StringAttribute{Computed: true, Description: "The description of the component."},
						"component_logo":        schema.StringAttribute{Computed: true, Description: "The logo URL of the component."},
						"component_extension":   schema.StringAttribute{Computed: true, Description: "The extension information of the component."},
						"create_time":           schema.Int64Attribute{Computed: true, Description: "The creation time of the component."},
						"update_time":           schema.Int64Attribute{Computed: true, Description: "The update time of the component."},
					},
					Blocks: map[string]schema.Block{
						"component_actions": schema.ListNestedBlock{
							Description: "The list of actions supported by the component.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"component_action_name":        schema.StringAttribute{Computed: true, Description: "The name of the component action."},
									"component_action_description": schema.StringAttribute{Computed: true, Description: "The description of the component action."},
								},
								Blocks: map[string]schema.Block{
									"input_configs": schema.ListNestedBlock{
										Description: "The input configurations of the component action.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"field_name":           schema.StringAttribute{Computed: true, Description: "The field name."},
												"field_type":           schema.StringAttribute{Computed: true, Description: "The field type."},
												"field_description":    schema.StringAttribute{Computed: true, Description: "The field description."},
												"field_display_config": schema.StringAttribute{Computed: true, Description: "The field display configuration."},
												"default_value":        schema.StringAttribute{Computed: true, Description: "The default value of the field."},
												"required":             schema.BoolAttribute{Computed: true, Description: "Whether the field is required."},
											},
										},
									},
									"output_configs": schema.ListNestedBlock{
										Description: "The output configurations of the component action.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"field_name": schema.StringAttribute{Computed: true, Description: "The field name."},
												"field_type": schema.StringAttribute{Computed: true, Description: "The field type."},
											},
										},
									},
								},
							},
						},
						"component_asset_configs": schema.ListNestedBlock{
							Description: "The asset configuration list of the component.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"field_name":        schema.StringAttribute{Computed: true, Description: "The field name."},
									"field_type":        schema.StringAttribute{Computed: true, Description: "The field type."},
									"field_description": schema.StringAttribute{Computed: true, Description: "The field description."},
									"required":          schema.BoolAttribute{Computed: true, Description: "Whether the field is required."},
									"encrypted":         schema.BoolAttribute{Computed: true, Description: "Whether the field is stored encrypted."},
									"default_value":     schema.StringAttribute{Computed: true, Description: "The default value of the field."},
								},
							},
						},
					},
				},
			},
		},
	}
}

type componentsDataSourceModel struct {
	ComponentName types.String     `tfsdk:"component_name"`
	NameRegex     types.String     `tfsdk:"name_regex"`
	Lang          types.String     `tfsdk:"lang"`
	RoleFor       types.Int64      `tfsdk:"role_for"`
	Ids           types.List       `tfsdk:"ids"`
	OutputFile    types.String     `tfsdk:"output_file"`
	Id            types.String     `tfsdk:"id"`
	Components    []componentModel `tfsdk:"components"`
}

type componentModel struct {
	Id                    types.String                `tfsdk:"id"`
	ComponentName         types.String                `tfsdk:"component_name"`
	ComponentAlias        types.String                `tfsdk:"component_alias"`
	ComponentDescription  types.String                `tfsdk:"component_description"`
	ComponentLogo         types.String                `tfsdk:"component_logo"`
	ComponentExtension    types.String                `tfsdk:"component_extension"`
	CreateTime            types.Int64                 `tfsdk:"create_time"`
	UpdateTime            types.Int64                 `tfsdk:"update_time"`
	ComponentActions      []componentActionModel      `tfsdk:"component_actions"`
	ComponentAssetConfigs []componentAssetConfigModel `tfsdk:"component_asset_configs"`
}

type componentActionModel struct {
	ComponentActionName        types.String                      `tfsdk:"component_action_name"`
	ComponentActionDescription types.String                      `tfsdk:"component_action_description"`
	InputConfigs               []componentFieldInputConfigModel  `tfsdk:"input_configs"`
	OutputConfigs              []componentFieldOutputConfigModel `tfsdk:"output_configs"`
}

type componentFieldInputConfigModel struct {
	FieldName          types.String `tfsdk:"field_name"`
	FieldType          types.String `tfsdk:"field_type"`
	FieldDescription   types.String `tfsdk:"field_description"`
	FieldDisplayConfig types.String `tfsdk:"field_display_config"`
	DefaultValue       types.String `tfsdk:"default_value"`
	Required           types.Bool   `tfsdk:"required"`
}

type componentFieldOutputConfigModel struct {
	FieldName types.String `tfsdk:"field_name"`
	FieldType types.String `tfsdk:"field_type"`
}

type componentAssetConfigModel struct {
	FieldName        types.String `tfsdk:"field_name"`
	FieldType        types.String `tfsdk:"field_type"`
	FieldDescription types.String `tfsdk:"field_description"`
	Required         types.Bool   `tfsdk:"required"`
	Encrypted        types.Bool   `tfsdk:"encrypted"`
	DefaultValue     types.String `tfsdk:"default_value"`
}

func (d *componentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const (
		action   = "ListComponents"
		product  = "sophonsoar"
		version  = "2025-09-03"
		pageSize = 20
		notExist = "Component.External.NotExist"
	)

	var config componentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() && config.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: invalid name_regex", componentsTypeName), err.Error())
			return
		}
		nameRegex = r
	}

	var idsFilter []string
	if !config.Ids.IsNull() {
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &idsFilter, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	idsSet := make(map[string]struct{}, len(idsFilter))
	for _, id := range idsFilter {
		idsSet[id] = struct{}{}
	}

	body := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   pageSize,
	}
	if !config.ComponentName.IsNull() && config.ComponentName.ValueString() != "" {
		body["ComponentName"] = config.ComponentName.ValueString()
	}
	if !config.Lang.IsNull() && config.Lang.ValueString() != "" {
		body["Lang"] = config.Lang.ValueString()
	}
	if !config.RoleFor.IsNull() {
		body["RoleFor"] = config.RoleFor.ValueInt64()
	}

	var objects []map[string]interface{}
	for {
		response, err := d.Client().RpcPost(product, version, action, nil, body, true)
		if err != nil {
			// A filter by component_name that matches no component returns
			// Component.External.NotExist. The data source contract is that a
			// filter with no match yields an empty list rather than an error.
			if isComponentNotExist(err, notExist) {
				objects = nil
				break
			}
			resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling %s", componentsTypeName, action), err.Error())
			return
		}

		rawComponents, ok := response["Components"].([]interface{})
		if !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", componentsTypeName, action),
				fmt.Sprintf("The response carried no Components list: %v", response),
			)
			return
		}

		for _, raw := range rawComponents {
			item, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			name := strValue(item["ComponentName"])
			if len(idsSet) > 0 {
				if _, keep := idsSet[name]; !keep {
					continue
				}
			}
			if nameRegex != nil && !nameRegex.MatchString(name) {
				continue
			}
			objects = append(objects, item)
		}

		if len(rawComponents) < pageSize {
			break
		}
		body["PageNumber"] = body["PageNumber"].(int) + 1
	}

	components := make([]componentModel, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		name := strValue(object["ComponentName"])
		components = append(components, componentModel{
			Id:                    types.StringValue(name),
			ComponentName:         types.StringValue(name),
			ComponentAlias:        types.StringValue(strValue(object["ComponentAlias"])),
			ComponentDescription:  types.StringValue(strValue(object["ComponentDescription"])),
			ComponentLogo:         types.StringValue(strValue(object["ComponentLogo"])),
			ComponentExtension:    types.StringValue(strValue(object["ComponentExtension"])),
			CreateTime:            types.Int64Value(int64Value(object["CreateTime"])),
			UpdateTime:            types.Int64Value(int64Value(object["UpdateTime"])),
			ComponentActions:      expandComponentActions(object["ComponentActions"]),
			ComponentAssetConfigs: expandComponentAssetConfigs(object["ComponentAssetConfigs"]),
		})
		ids = append(ids, name)
	}

	state := componentsDataSourceModel{
		ComponentName: config.ComponentName,
		NameRegex:     config.NameRegex,
		Lang:          config.Lang,
		RoleFor:       config.RoleFor,
		OutputFile:    config.OutputFile,
		Components:    components,
	}

	idList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Ids = idList
	state.Id = types.StringValue(componentsDataResourceId(ids))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		if err := writeComponentsToFile(config.OutputFile.ValueString(), components); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: writing output_file", componentsTypeName), err.Error())
			return
		}
	}
}

func expandComponentActions(raw interface{}) []componentActionModel {
	items, _ := raw.([]interface{})
	out := make([]componentActionModel, 0, len(items))
	for _, rawItem := range items {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, componentActionModel{
			ComponentActionName:        types.StringValue(strValue(item["ComponentActionName"])),
			ComponentActionDescription: types.StringValue(strValue(item["ComponentActionDescription"])),
			InputConfigs:               expandInputConfigs(item["InputConfigs"]),
			OutputConfigs:              expandOutputConfigs(item["OutputConfigs"]),
		})
	}
	return out
}

func expandInputConfigs(raw interface{}) []componentFieldInputConfigModel {
	items, _ := raw.([]interface{})
	out := make([]componentFieldInputConfigModel, 0, len(items))
	for _, rawItem := range items {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, componentFieldInputConfigModel{
			FieldName:          types.StringValue(strValue(item["FieldName"])),
			FieldType:          types.StringValue(strValue(item["FieldType"])),
			FieldDescription:   types.StringValue(strValue(item["FieldDescription"])),
			FieldDisplayConfig: types.StringValue(strValue(item["FieldDisplayConfig"])),
			DefaultValue:       types.StringValue(strValue(item["DefaultValue"])),
			Required:           types.BoolValue(boolValue(item["Required"])),
		})
	}
	return out
}

func expandOutputConfigs(raw interface{}) []componentFieldOutputConfigModel {
	items, _ := raw.([]interface{})
	out := make([]componentFieldOutputConfigModel, 0, len(items))
	for _, rawItem := range items {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, componentFieldOutputConfigModel{
			FieldName: types.StringValue(strValue(item["FieldName"])),
			FieldType: types.StringValue(strValue(item["FieldType"])),
		})
	}
	return out
}

func expandComponentAssetConfigs(raw interface{}) []componentAssetConfigModel {
	items, _ := raw.([]interface{})
	out := make([]componentAssetConfigModel, 0, len(items))
	for _, rawItem := range items {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, componentAssetConfigModel{
			FieldName:        types.StringValue(strValue(item["FieldName"])),
			FieldType:        types.StringValue(strValue(item["FieldType"])),
			FieldDescription: types.StringValue(strValue(item["FieldDescription"])),
			Required:         types.BoolValue(boolValue(item["Required"])),
			Encrypted:        types.BoolValue(boolValue(item["Encrypted"])),
			DefaultValue:     types.StringValue(strValue(item["DefaultValue"])),
		})
	}
	return out
}

// isComponentNotExist reports whether err carries the
// Component.External.NotExist code that the ListComponents API returns when a
// component_name filter matches nothing.
func isComponentNotExist(err error, code string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), code)
}

// strValue renders a PascalCase response field as a string, empty for nil.
func strValue(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// int64Value coerces numeric response fields to int64, tolerating JSON's
// default float64 and numeric strings.
func int64Value(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	case float32:
		return int64(n)
	case string:
		if n == "" {
			return 0
		}
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return 0
		}
		return i
	default:
		if s := fmt.Sprintf("%v", v); s != "" {
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				return i
			}
		}
		return 0
	}
}

// boolValue coerces response fields to bool, tolerating string and numeric
// truthiness.
func boolValue(v interface{}) bool {
	switch b := v.(type) {
	case bool:
		return b
	case string:
		bb, _ := strconv.ParseBool(b)
		return bb
	case float64:
		return b != 0
	case int:
		return b != 0
	case int64:
		return b != 0
	default:
		return false
	}
}

// componentsDataResourceId mirrors the SDK data source id: a stable hash of the
// listed component names so the data source address is deterministic for a given
// result set.
func componentsDataResourceId(ids []string) string {
	var sb strings.Builder
	for _, id := range ids {
		sb.WriteString(id)
		sb.WriteString("-")
	}
	return fmt.Sprintf("%d", helper.Hashcode(sb.String()))
}

// writeComponentsToFile writes the component list as indented JSON, expanding a
// leading ~ to the user's home directory and replacing any existing file.
func writeComponentsToFile(path string, components []componentModel) error {
	if path == "" {
		return nil
	}
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil && home != "" {
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	bs, err := json.MarshalIndent(components, "", "\t")
	if err != nil {
		return fmt.Errorf("marshaling components: %w", err)
	}
	return os.WriteFile(path, bs, 0o644)
}
