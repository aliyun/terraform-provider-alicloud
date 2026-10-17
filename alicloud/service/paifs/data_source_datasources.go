package paifs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	datasourcesTypeName = "alicloud_pai_feature_store_datasources"

	// listPageSize caps the per-page size used while paginating ListDatasources.
	listPageSize = 50
	// listPageGuard bounds the pagination loop so a misreported total cannot run
	// away forever.
	listPageGuard = 100
)

var (
	_ datasource.DataSource              = &datasourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &datasourcesDataSource{}
)

func NewDatasourcesDataSource() datasource.DataSource {
	return &datasourcesDataSource{}
}

type datasourcesDataSource struct {
	fwadapt.DataSourceBase
}

type datasourcesDataSourceModel struct {
	Id          types.String     `tfsdk:"id"`
	InstanceId  types.String     `tfsdk:"instance_id"`
	WorkspaceId types.String     `tfsdk:"workspace_id"`
	Type        types.String     `tfsdk:"type"`
	Name        types.String     `tfsdk:"name"`
	Ids         types.List       `tfsdk:"ids"`
	Datasources []datasourceItem `tfsdk:"datasources"`
	TotalCount  types.Int64      `tfsdk:"total_count"`
}

type datasourceItem struct {
	DatasourceId   types.String `tfsdk:"datasource_id"`
	InstanceId     types.String `tfsdk:"instance_id"`
	Type           types.String `tfsdk:"type"`
	DatasourceName types.String `tfsdk:"datasource_name"`
	Uri            types.String `tfsdk:"uri"`
	Config         types.String `tfsdk:"config"`
	CreateTime     types.String `tfsdk:"create_time"`
}

func (d *datasourcesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = datasourcesTypeName
}

func (d *datasourcesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of PAI Feature Store Datasources available to the user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, same as `instance_id`.",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore instance.",
			},
			"workspace_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the results by the PAI workspace ID.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the results by the data source type. Valid values: `Hologres`, `GraphCompute`, `FeatureDB`, `MaxCompute`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Hologres", "GraphCompute", "FeatureDB", "MaxCompute"),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter the results by the datasource name.",
			},
			"ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of datasource IDs.",
			},
			"datasources": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"datasource_id": schema.StringAttribute{
							Computed:    true,
							Description: "The first ID of the resource.",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the FeatureStore instance.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The data source type.",
						},
						"datasource_name": schema.StringAttribute{
							Computed:    true,
							Description: "The datasource name.",
						},
						"uri": schema.StringAttribute{
							Computed:    true,
							Description: "The resource address of the datasource.",
						},
						"config": schema.StringAttribute{
							Computed:    true,
							Description: "The resource configuration information as a JSON string.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "The creation time of the resource.",
						},
					},
				},
				Description: "A list of FeatureStore Datasources.",
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of matching datasources.",
			},
		},
	}
}

func (d *datasourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config datasourcesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := config.InstanceId.ValueString()
	if instanceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: missing filter", datasourcesTypeName),
			"The data source requires a non-empty instance_id.",
		)
		return
	}

	var items []datasourceItem
	var ids []string
	total := 0
	for page := 1; page <= listPageGuard; page++ {
		body := map[string]interface{}{
			"InstanceId": instanceId,
			"PageNumber": strconv.Itoa(page),
			"PageSize":   strconv.Itoa(listPageSize),
		}
		if !config.WorkspaceId.IsNull() {
			body["WorkspaceId"] = config.WorkspaceId.ValueString()
		}
		if !config.Type.IsNull() {
			body["Type"] = config.Type.ValueString()
		}
		if !config.Name.IsNull() {
			body["Name"] = config.Name.ValueString()
		}

		response, err := d.Client().RpcPost(paifsProduct, paifsVersion, "ListDatasources", nil, body, true)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling ListDatasources", datasourcesTypeName),
				err.Error(),
			)
			return
		}

		list, _ := response["Datasources"].([]interface{})
		if total == 0 {
			total = toInt(response["TotalCount"])
		}
		for _, raw := range list {
			object, _ := raw.(map[string]interface{})
			if object == nil {
				continue
			}
			items = append(items, datasourceItemFrom(object, instanceId))
			if v, ok := object["DatasourceId"].(string); ok && v != "" {
				ids = append(ids, v)
			}
		}
		if len(list) == 0 || len(items) >= total {
			break
		}
	}

	idList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := datasourcesDataSourceModel{
		Id:          types.StringValue(instanceId),
		InstanceId:  types.StringValue(instanceId),
		WorkspaceId: config.WorkspaceId,
		Type:        config.Type,
		Name:        config.Name,
		Ids:         idList,
		Datasources: items,
		TotalCount:  types.Int64Value(int64(total)),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// datasourceItemFrom maps a single ListDatasources entry onto a nested item.
// InstanceId is a request filter and is not echoed by the response, so it is
// carried from the data source configuration.
func datasourceItemFrom(object map[string]interface{}, instanceId string) datasourceItem {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return datasourceItem{
		DatasourceId:   str("DatasourceId"),
		InstanceId:     nonEmptyString(instanceId),
		Type:           str("Type"),
		DatasourceName: str("Name"),
		Uri:            str("Uri"),
		Config:         str("Config"),
		CreateTime:     str("GmtCreateTime"),
	}
}

// toInt coerces the loose POP RPC number encoding (float64 or string) into an int.
func toInt(v interface{}) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return 0
}
