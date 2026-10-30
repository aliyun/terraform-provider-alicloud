package oos

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const inventoriesTypeName = "alicloud_oos_inventories"

var (
	_ datasource.DataSource              = &inventoriesDataSource{}
	_ datasource.DataSourceWithConfigure = &inventoriesDataSource{}
)

func NewInventoriesDataSource() datasource.DataSource {
	return &inventoriesDataSource{}
}

type inventoriesDataSource struct {
	fwadapt.DataSourceBase
}

type inventoriesModel struct {
	NextToken       types.String `tfsdk:"next_token"`
	MaxResults      types.Int64  `tfsdk:"max_results"`
	Aggregator      types.List   `tfsdk:"aggregator"`
	Filter          types.List   `tfsdk:"filter"`
	Id              types.String `tfsdk:"id"`
	NextTokenResult types.String `tfsdk:"next_token_result"`
	Entities        types.List   `tfsdk:"inventories"`
}

var inventoriesFilterAttrTypes = map[string]attr.Type{
	"name":     types.StringType,
	"operator": types.StringType,
	"value":    types.ListType{ElemType: types.StringType},
}

func (d *inventoriesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = inventoriesTypeName
}

func (d *inventoriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the configuration inventory of the managed instances, or the aggregated inventory information, in Operation Orchestration Service.",
		Attributes: map[string]schema.Attribute{
			"aggregator": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The aggregators of the inventory.",
			},
			"filter": schema.ListAttribute{
				Optional:    true,
				ElementType: types.ObjectType{AttrTypes: inventoriesFilterAttrTypes},
				Description: "The filter rule of the components. Each element contains `name`, `operator` and `value`. A maximum of 5 rules are supported, and each rule accepts a maximum of 20 values. Valid values of `operator`: `Equal`, `NotEqual`, `BeginWith`, `LessThan` and `GreaterThan`.",
			},
			"next_token": schema.StringAttribute{
				Optional:    true,
				Description: "The token that starts the query. Leave it unset to start from the first page.",
			},
			"max_results": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: `1` to `100`. Defaults to `50`.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is derived from the query condition.",
			},
			"next_token_result": schema.StringAttribute{
				Computed:    true,
				Description: "The token that is used to retrieve the next page of results.",
			},
			"inventories": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The inventory entities. Each element is the JSON representation of one entity. The keys of the entity are the component names, which vary with the queried component.",
			},
		},
	}
}

func (d *inventoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const action = "SearchInventory"

	var config inventoriesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	request := map[string]interface{}{}
	request["RegionId"] = d.Meta().RegionId

	if !config.NextToken.IsNull() && !config.NextToken.IsUnknown() && config.NextToken.ValueString() != "" {
		request["NextToken"] = config.NextToken.ValueString()
	}

	if !config.MaxResults.IsNull() && !config.MaxResults.IsUnknown() {
		request["MaxResults"] = config.MaxResults.ValueInt64()
	}

	if !config.Aggregator.IsNull() {
		var aggregators []string
		resp.Diagnostics.Append(config.Aggregator.ElementsAs(ctx, &aggregators, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(aggregators) > 0 {
			request["Aggregator"] = aggregators
		}
	}

	if !config.Filter.IsNull() {
		elements := config.Filter.Elements()
		filters := make([]map[string]interface{}, 0, len(elements))
		for _, elementValue := range elements {
			item := map[string]interface{}{}
			element, ok := elementValue.(types.Object)
			if !ok {
				continue
			}

			attributes := element.Attributes()
			name, _ := attributes["name"].(types.String)
			operator, _ := attributes["operator"].(types.String)
			value, _ := attributes["value"].(types.List)

			if !name.IsNull() && !name.IsUnknown() && name.ValueString() != "" {
				item["Name"] = name.ValueString()
			}

			if !operator.IsNull() && !operator.IsUnknown() && operator.ValueString() != "" {
				item["Operator"] = operator.ValueString()
			}

			var values []string
			if !value.IsNull() && !value.IsUnknown() {
				resp.Diagnostics.Append(value.ElementsAs(ctx, &values, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
			}

			if len(values) > 0 {
				item["Value"] = values
			}

			filters = append(filters, item)
		}

		if len(filters) > 0 {
			request["Filter"] = filters
		}
	}

	var objects basetypes.ListValue
	var diags diag.Diagnostics

	response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return d.Meta().RpcPost("oos", "2019-06-01", action, nil, request, true)
	}, 5*time.Minute, helper.NeedRetry)
	logs.AddDebug(action, response, request)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", inventoriesTypeName, action),
			err.Error(),
		)
		return
	}

	rawEntities, _ := response["Entities"].([]interface{})
	entities := make([]attr.Value, 0, len(rawEntities))
	for _, raw := range rawEntities {
		encoded, err := json.Marshal(raw)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: encoding %s", inventoriesTypeName, action),
				err.Error(),
			)
			return
		}

		entities = append(entities, types.StringValue(string(encoded)))
	}

	objects, diags = types.ListValue(types.StringType, entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nextTokenResult := types.StringNull()
	if reported, ok := response["NextToken"].(string); ok && reported != "" {
		nextTokenResult = types.StringValue(reported)
	}

	state := config
	state.Id = types.StringValue(helper.DataResourceIdHash([]string{fmt.Sprint(objects)}))
	state.NextTokenResult = nextTokenResult
	state.Entities = objects

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
