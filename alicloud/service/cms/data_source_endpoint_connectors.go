package cms

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const endpointConnectorsTypeName = "alicloud_cms_endpoint_connectors"

var (
	_ datasource.DataSource              = &endpointConnectorsDataSource{}
	_ datasource.DataSourceWithConfigure = &endpointConnectorsDataSource{}
)

func NewEndpointConnectorsDataSource() datasource.DataSource {
	return &endpointConnectorsDataSource{}
}

type endpointConnectorsDataSource struct {
	fwadapt.DataSourceBase
}

type connectorItemModel struct {
	Id          types.String `tfsdk:"id"`
	ConnectorId types.String `tfsdk:"connector_id"`
	Workspace   types.String `tfsdk:"workspace"`
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Alias       types.String `tfsdk:"alias"`
	Endpoint    types.String `tfsdk:"endpoint"`
	Description types.String `tfsdk:"description"`
	Properties  types.Map    `tfsdk:"properties"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	UpdatedAt   types.Int64  `tfsdk:"updated_at"`
}

type endpointConnectorsDataSourceModel struct {
	Ids        []types.String       `tfsdk:"ids"`
	NameRegex  types.String         `tfsdk:"name_regex"`
	Workspace  types.String         `tfsdk:"workspace"`
	Type       types.String         `tfsdk:"type"`
	Name       types.String         `tfsdk:"name"`
	OutputFile types.String         `tfsdk:"output_file"`
	Id         types.String         `tfsdk:"id"`
	Connectors []connectorItemModel `tfsdk:"connectors"`
}

func (d *endpointConnectorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = endpointConnectorsTypeName
}

func (d *endpointConnectorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"name_regex":  schema.StringAttribute{Optional: true},
			"workspace":   schema.StringAttribute{Required: true},
			"type":        schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"output_file": schema.StringAttribute{Optional: true},
			"id":          schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			// connectors uses ListNestedBlock (not ListNestedAttribute) because the
			// provider is served via tf5muxserver (protocol v5), which cannot
			// convert nested attributes with Attributes set.
			"connectors": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.StringAttribute{Computed: true},
						"connector_id": schema.StringAttribute{Computed: true},
						"workspace":    schema.StringAttribute{Computed: true},
						"type":         schema.StringAttribute{Computed: true},
						"name":         schema.StringAttribute{Computed: true},
						"alias":        schema.StringAttribute{Computed: true},
						"endpoint":     schema.StringAttribute{Computed: true},
						"description":  schema.StringAttribute{Computed: true},
						"properties": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"created_at": schema.Int64Attribute{Computed: true},
						"updated_at": schema.Int64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *endpointConnectorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config endpointConnectorsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() && config.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", endpointConnectorsTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	for _, id := range config.Ids {
		if id.IsNull() {
			continue
		}
		idsMap[id.ValueString()] = id.ValueString()
	}

	objects, err := NewService(d.Client()).ListEndpointConnectors(
		config.Workspace.ValueString(),
		config.Type.ValueString(),
		config.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListEndpointConnectors", endpointConnectorsTypeName),
			err.Error(),
		)
		return
	}

	ids := make([]string, 0, len(objects))
	connectors := make([]connectorItemModel, 0, len(objects))
	for _, object := range objects {
		itemId := fmt.Sprintf("%v:%v", object["workspace"], object["connectorId"])
		if nameRegex != nil {
			if !nameRegex.MatchString(fmt.Sprint(object["name"])) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[itemId]; !ok {
				continue
			}
		}
		ids = append(ids, itemId)
		connectors = append(connectors, connectorItemState(object))
	}

	state := endpointConnectorsDataSourceModel{
		Workspace:  config.Workspace,
		NameRegex:  config.NameRegex,
		Type:       config.Type,
		Name:       config.Name,
		OutputFile: config.OutputFile,
		Id:         types.StringValue(helper.DataResourceIdHash(ids)),
		Ids:        make([]types.String, 0, len(ids)),
		Connectors: connectors,
	}
	for _, id := range ids {
		state.Ids = append(state.Ids, types.StringValue(id))
	}

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		out := make([]map[string]interface{}, 0, len(connectors))
		for _, c := range connectors {
			out = append(out, connectorItemToMap(c))
		}
		if err := helper.WriteToFile(config.OutputFile.ValueString(), out); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output file", endpointConnectorsTypeName),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func connectorItemState(object map[string]interface{}) connectorItemModel {
	str := func(key string) types.String {
		if v, ok := object[key]; ok && v != nil {
			return types.StringValue(fmt.Sprint(v))
		}
		return types.StringNull()
	}
	return connectorItemModel{
		Id:          types.StringValue(fmt.Sprintf("%v:%v", object["workspace"], object["connectorId"])),
		ConnectorId: str("connectorId"),
		Workspace:   str("workspace"),
		Type:        str("type"),
		Name:        str("name"),
		Alias:       str("alias"),
		Endpoint:    str("endpoint"),
		Description: str("description"),
		Properties:  mapToTypeMap(object["properties"]),
		CreatedAt:   mapInt64(object, "createdAt"),
		UpdatedAt:   mapInt64(object, "updatedAt"),
	}
}

func connectorItemToMap(c connectorItemModel) map[string]interface{} {
	return map[string]interface{}{
		"id":           c.Id.ValueString(),
		"connector_id": c.ConnectorId.ValueString(),
		"workspace":    c.Workspace.ValueString(),
		"type":         c.Type.ValueString(),
		"name":         c.Name.ValueString(),
		"alias":        c.Alias.ValueString(),
		"endpoint":     c.Endpoint.ValueString(),
		"description":  c.Description.ValueString(),
		"created_at":   c.CreatedAt.ValueInt64(),
		"updated_at":   c.UpdatedAt.ValueInt64(),
	}
}
