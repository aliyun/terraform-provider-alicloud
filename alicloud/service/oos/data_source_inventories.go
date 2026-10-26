package oos

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	inventoriesTypeName   = "alicloud_oos_inventories"
	inventoriesAction     = "SearchInventory"
	inventoriesAPIVersion = "2019-06-01"
	inventoriesReadWindow = 3 * time.Minute
)

// inventoriesFilterAttrTypes must stay flat: the provider is served over
// protocol version 5, which rejects nested attributes.
var inventoriesFilterAttrTypes = map[string]attr.Type{
	"name":     types.StringType,
	"operator": types.StringType,
	"value":    types.ListType{ElemType: types.StringType},
}

var inventoriesFilterOperators = []string{"Equal", "NotEqual", "BeginWith", "LessThan", "GreaterThan"}

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
	Id              types.String `tfsdk:"id"`
	RegionId        types.String `tfsdk:"region_id"`
	NextToken       types.String `tfsdk:"next_token"`
	MaxResults      types.Int64  `tfsdk:"max_results"`
	Filter          types.List   `tfsdk:"filter"`
	Aggregator      types.List   `tfsdk:"aggregator"`
	Entities        types.List   `tfsdk:"entities"`
	NextTokenResult types.String `tfsdk:"next_token_result"`
}

func (d *inventoriesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = inventoriesTypeName
}

func (d *inventoriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the configuration inventory of the managed instances, or the aggregated inventory information, in Operation Orchestration Service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is derived from the query condition.",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the region. Defaults to the region of the provider.",
			},
			"next_token": schema.StringAttribute{
				Optional:    true,
				Description: "The token that starts the query. Leave it unset to start from the first page.",
			},
			"next_token_result": schema.StringAttribute{
				Computed:    true,
				Description: "The token that is used to retrieve the next page of results.",
			},
			"max_results": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: `1` to `100`. Defaults to `50`.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"filter": schema.ListAttribute{
				Optional:    true,
				Description: "The filter rule of the components. Each element contains `name`, `operator` and `value`. A maximum of 5 rules are supported, and each rule accepts a maximum of 20 values. Valid values of `operator`: `Equal`, `NotEqual`, `BeginWith`, `LessThan` and `GreaterThan`.",
				ElementType: types.ObjectType{AttrTypes: inventoriesFilterAttrTypes},
				Validators: []validator.List{
					listvalidator.SizeAtMost(5),
				},
			},
			"aggregator": schema.ListAttribute{
				Optional:    true,
				Description: "The aggregators of the inventory. Valid values: `ACS:Application.Name` and `ACS:Application.Version`. A maximum of 5 aggregators are supported.",
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(5),
					listvalidator.ValueStringsAre(stringvalidator.OneOf("ACS:Application.Name", "ACS:Application.Version")),
				},
			},
			"entities": schema.ListAttribute{
				Computed:    true,
				Description: "The inventory entities. Each element is the JSON representation of one entity. The keys of the entity are the component names, which vary with the queried component.",
				ElementType: types.StringType,
			},
		},
	}
}

func (d *inventoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data inventoriesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}
	regionId := data.RegionId.ValueString()
	if regionId == "" {
		regionId = d.Client().RegionId
	}
	if regionId != "" {
		request["RegionId"] = regionId
	}
	if !data.NextToken.IsNull() && !data.NextToken.IsUnknown() && data.NextToken.ValueString() != "" {
		request["NextToken"] = data.NextToken.ValueString()
	}
	if !data.MaxResults.IsNull() && !data.MaxResults.IsUnknown() {
		request["MaxResults"] = data.MaxResults.ValueInt64()
	}

	if !data.Filter.IsNull() && !data.Filter.IsUnknown() {
		elements := data.Filter.Elements()
		filters := make([]map[string]interface{}, 0, len(elements))
		for i, elementValue := range elements {
			element, ok := elementValue.(types.Object)
			if !ok {
				continue
			}
			attributes := element.Attributes()
			name, _ := attributes["name"].(types.String)
			operator, _ := attributes["operator"].(types.String)
			value, _ := attributes["value"].(types.List)

			if operator.IsNull() || operator.IsUnknown() || !slices.Contains(inventoriesFilterOperators, operator.ValueString()) {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Reading %s: invalid filter rule", inventoriesTypeName),
					fmt.Sprintf("filter.%d.operator must be one of %v.", i, inventoriesFilterOperators),
				)
				return
			}

			var values []string
			if !value.IsNull() && !value.IsUnknown() {
				resp.Diagnostics.Append(value.ElementsAs(ctx, &values, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
			}
			if len(values) > 20 {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Reading %s: invalid filter rule", inventoriesTypeName),
					fmt.Sprintf("filter.%d.value accepts at most 20 values.", i),
				)
				return
			}

			item := map[string]interface{}{
				"Operator": operator.ValueString(),
			}
			if !name.IsNull() && !name.IsUnknown() && name.ValueString() != "" {
				item["Name"] = name.ValueString()
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

	if !data.Aggregator.IsNull() && !data.Aggregator.IsUnknown() {
		var aggregators []string
		resp.Diagnostics.Append(data.Aggregator.ElementsAs(ctx, &aggregators, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(aggregators) > 0 {
			request["Aggregator"] = aggregators
		}
	}

	response, err := helper.RetryRpcPost(ctx, inventoriesAction, func() (map[string]interface{}, error) {
		return d.Client().RpcPost("Oos", inventoriesAPIVersion, inventoriesAction, nil, request, true)
	}, inventoriesReadWindow, helper.NeedRetry)
	logs.AddDebug(inventoriesAction, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", inventoriesTypeName, inventoriesAction),
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
				fmt.Sprintf("Reading %s: encoding %s", inventoriesTypeName, inventoriesAction),
				err.Error(),
			)
			return
		}
		entities = append(entities, types.StringValue(string(encoded)))
	}
	entityList, diags := types.ListValue(types.StringType, entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nextTokenResult := types.StringNull()
	if reported, ok := response["NextToken"].(string); ok && reported != "" {
		nextTokenResult = types.StringValue(reported)
	}

	state := inventoriesModel{
		Id:              types.StringValue(helper.DataResourceIdHashFromMap(request)),
		RegionId:        data.RegionId,
		NextToken:       data.NextToken,
		NextTokenResult: nextTokenResult,
		MaxResults:      data.MaxResults,
		Filter:          data.Filter,
		Aggregator:      data.Aggregator,
		Entities:        entityList,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
