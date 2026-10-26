package oos

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const inventoryItemsTypeName = "alicloud_oos_inventory_items"

const (
	inventoryItemsAction            = "SearchInventory"
	inventoryItemsMaxFilterGroups   = 5
	inventoryItemsMaxFilterValues   = 20
	inventoryItemsDefaultMaxResults = 50
	inventoryItemsQueryTimeout      = 5 * time.Minute
)

var inventoryItemsOperators = []string{"Equal", "NotEqual", "BeginWith", "LessThan", "GreaterThan"}

var (
	_ datasource.DataSource              = &inventoryItemsDataSource{}
	_ datasource.DataSourceWithConfigure = &inventoryItemsDataSource{}
)

func NewInventoryItemsDataSource() datasource.DataSource {
	return &inventoryItemsDataSource{}
}

type inventoryItemsDataSource struct {
	fwadapt.DataSourceBase
}

type inventoryItemsFilterModel struct {
	Name     types.String   `tfsdk:"name"`
	Operator types.String   `tfsdk:"operator"`
	Value    []types.String `tfsdk:"value"`
}

type inventoryItemsModel struct {
	Id         types.String                `tfsdk:"id"`
	RegionId   types.String                `tfsdk:"region_id"`
	Filter     []inventoryItemsFilterModel `tfsdk:"filter"`
	Aggregator []types.String              `tfsdk:"aggregator"`
	MaxResults types.Int64                 `tfsdk:"max_results"`
	Entities   types.String                `tfsdk:"entities"`
}

// inventoryItemsFilterObjectType describes the element shape of the filter list attribute.
var inventoryItemsFilterObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name":     types.StringType,
	"operator": types.StringType,
	"value":    types.ListType{ElemType: types.StringType},
}}

// inventoryItemsFilterValidator rejects filter entries the API would reject, so the
// mismatch surfaces at plan time and no request is sent.
type inventoryItemsFilterValidator struct{}

func (v inventoryItemsFilterValidator) Description(context.Context) string {
	return fmt.Sprintf("Each filter entry must set `name`, `value` and a supported `operator` (%s); `value` accepts 1 to %d entries.",
		inventoryItemsOperatorList(), inventoryItemsMaxFilterValues)
}

func (v inventoryItemsFilterValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v inventoryItemsFilterValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	for i, element := range req.ConfigValue.Elements() {
		entryPath := req.Path.AtListIndex(i)
		entry, ok := element.(types.Object)
		// An unknown entry is checked once its value is known; a null entry
		// carries none of the required attributes.
		if !ok || entry.IsUnknown() {
			continue
		}
		if entry.IsNull() {
			resp.Diagnostics.AddAttributeError(entryPath.AtName("name"), "Invalid Filter Name",
				"`name` must be set for every filter entry, for example \"ACS:InstanceInformation.InstanceId\".")
			continue
		}

		attributes := entry.Attributes()
		// Unknown values can come from a resource attribute or another data
		// source and are only known after apply, so they are deferred instead
		// of being rejected as missing.
		if name, ok := attributes["name"].(types.String); ok && !name.IsUnknown() {
			if name.IsNull() || name.ValueString() == "" {
				resp.Diagnostics.AddAttributeError(entryPath.AtName("name"), "Invalid Filter Name",
					"`name` must be set for every filter entry, for example \"ACS:InstanceInformation.InstanceId\".")
			}
		}

		if operator, ok := attributes["operator"].(types.String); ok && !operator.IsUnknown() {
			switch {
			case operator.IsNull() || operator.ValueString() == "":
				resp.Diagnostics.AddAttributeError(entryPath.AtName("operator"), "Invalid Filter Operator",
					fmt.Sprintf("`operator` must be set for every filter entry. Valid values: %s.", inventoryItemsOperatorList()))
			case !inventoryItemsOperatorKnown(operator.ValueString()):
				resp.Diagnostics.AddAttributeError(entryPath.AtName("operator"), "Invalid Filter Operator",
					fmt.Sprintf("`operator` must be one of %s, got %q.", inventoryItemsOperatorList(), operator.ValueString()))
			}
		}

		if value, ok := attributes["value"].(types.List); ok && !value.IsUnknown() {
			switch {
			case value.IsNull() || len(value.Elements()) == 0:
				resp.Diagnostics.AddAttributeError(entryPath.AtName("value"), "Invalid Filter Value",
					"`value` must contain at least one entry.")
			case len(value.Elements()) > inventoryItemsMaxFilterValues:
				resp.Diagnostics.AddAttributeError(entryPath.AtName("value"), "Invalid Filter Value",
					fmt.Sprintf("`value` accepts at most %d entries, got %d.", inventoryItemsMaxFilterValues, len(value.Elements())))
			}
		}
	}
}

func inventoryItemsOperatorKnown(operator string) bool {
	for _, known := range inventoryItemsOperators {
		if operator == known {
			return true
		}
	}
	return false
}

func inventoryItemsOperatorList() string {
	list := ""
	for i, operator := range inventoryItemsOperators {
		if i > 0 {
			list += ", "
		}
		list += operator
	}
	return list
}

func (d *inventoryItemsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = inventoryItemsTypeName
}

func (d *inventoryItemsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of configuration inventory entries collected by CloudOps Orchestration Service (OOS).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source.",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the region where the configuration inventory is queried. It must be the region configured on the provider, because the request is sent to the OOS endpoint of the provider region. Defaults to the region of the provider.",
			},
			"filter": schema.ListAttribute{
				Optional:    true,
				ElementType: inventoryItemsFilterObjectType,
				Description: fmt.Sprintf("Component filters. Up to %d filter groups are supported. Every entry must set `name`, `operator` and `value`.", inventoryItemsMaxFilterGroups),
				Validators: []validator.List{
					listvalidator.SizeAtMost(inventoryItemsMaxFilterGroups),
					inventoryItemsFilterValidator{},
				},
			},
			"aggregator": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: fmt.Sprintf("Aggregators used to query inventory aggregation information. Up to %d aggregators are supported. Each value has the form `ACS:<component>.<property>`, for example `ACS:InstanceInformation.PlatformName`; the aggregatable properties depend on the component. An unsupported combination is rejected by the API.", inventoryItemsMaxFilterGroups),
				Validators: []validator.List{
					listvalidator.SizeAtMost(inventoryItemsMaxFilterGroups),
				},
			},
			"max_results": schema.Int64Attribute{
				Optional:    true,
				Description: fmt.Sprintf("The number of entries to return on each page. Valid values: 1 to 100. Default value: %d. All pages are collected, so this only controls the request page size.", inventoryItemsDefaultMaxResults),
				Validators:  []validator.Int64{int64validator.Between(1, 100)},
			},
			"entities": schema.StringAttribute{
				Computed:    true,
				Description: "The inventory entries returned by `SearchInventory`, encoded as a JSON array. Each array element is one component record, keyed by the component name (for example `ACS:InstanceInformation` or `ACS:Application`) and carrying the component content. The component attribute names inside the content can be discovered with the `GetInventorySchema` API. The value mirrors the `Entities` array of the response exactly and concatenates the entries of all pages in page order; a query without any entry yields `[]`, and an aggregation query can return a component envelope with an empty content list.",
			},
		},
	}
}

func (d *inventoryItemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model inventoryItemsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()
	providerRegion := ""
	if client != nil {
		providerRegion = client.RegionId
	}
	regionId := model.RegionId.ValueString()
	if model.RegionId.IsNull() || model.RegionId.IsUnknown() || regionId == "" {
		regionId = providerRegion
	}
	// The request is sent to the OOS endpoint of the provider region, so a
	// different region_id would silently return the inventory of the endpoint
	// region. Reject it before any request is sent.
	if providerRegion != "" && regionId != providerRegion {
		resp.Diagnostics.AddAttributeError(path.Root("region_id"), "Region Mismatch",
			fmt.Sprintf("`region_id` is %q but the provider is configured for region %q. This data source calls the OOS endpoint of the provider region, so it would return the inventory of %q instead of %q. Configure a provider alias with `region = %q` and use that provider in this data source, or set `region_id` to %q.",
				regionId, providerRegion, providerRegion, regionId, regionId, providerRegion))
		return
	}
	maxResults := int64(inventoryItemsDefaultMaxResults)
	if !model.MaxResults.IsNull() && !model.MaxResults.IsUnknown() {
		maxResults = model.MaxResults.ValueInt64()
	}
	filters := expandInventoryFilters(model.Filter)
	aggregators := inventoryAggregators(model.Aggregator)

	entities := make([]interface{}, 0)
	nextToken := ""
	for {
		request := map[string]interface{}{
			"RegionId":   regionId,
			"MaxResults": maxResults,
		}
		if nextToken != "" {
			request["NextToken"] = nextToken
		}
		if len(filters) > 0 {
			request["Filter"] = filters
		}
		if len(aggregators) > 0 {
			request["Aggregator"] = aggregators
		}

		response, err := d.searchInventory(ctx, request)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", inventoryItemsTypeName, inventoryItemsAction),
				err.Error(),
			)
			return
		}

		entities = append(entities, inventoryResponseEntities(response)...)
		nextToken = inventoryResponseNextToken(response)
		if nextToken == "" {
			break
		}
	}

	encoded, err := json.Marshal(entities)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: encoding %s response", inventoryItemsTypeName, inventoryItemsAction),
			err.Error(),
		)
		return
	}

	model.Id = types.StringValue(inventoryItemsId(regionId, maxResults, model.Filter, model.Aggregator))
	model.RegionId = types.StringValue(regionId)
	model.Entities = types.StringValue(string(encoded))

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// searchInventory reads one page and retries the throttling and transient errors
// that a multi-page collection can hit.
func (d *inventoryItemsDataSource) searchInventory(ctx context.Context, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := sdkretry.RetryContext(ctx, inventoryItemsQueryTimeout, func() *sdkretry.RetryError {
		result, rpcErr := d.Client().RpcPost("Oos", "2019-06-01", inventoryItemsAction, nil, request, true)
		if rpcErr != nil {
			if helper.NeedRetry(rpcErr) {
				wait()
				return sdkretry.RetryableError(rpcErr)
			}
			return sdkretry.NonRetryableError(rpcErr)
		}
		response = result
		return nil
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, inventoryItemsTypeName, inventoryItemsAction, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(inventoryItemsAction, response, request)
	return response, nil
}

func expandInventoryFilters(filters []inventoryItemsFilterModel) []interface{} {
	expanded := make([]interface{}, 0, len(filters))
	for _, filter := range filters {
		entry := map[string]interface{}{}
		if !filter.Name.IsNull() && !filter.Name.IsUnknown() {
			entry["Name"] = filter.Name.ValueString()
		}
		if !filter.Operator.IsNull() && !filter.Operator.IsUnknown() && filter.Operator.ValueString() != "" {
			entry["Operator"] = filter.Operator.ValueString()
		}
		values := make([]string, 0, len(filter.Value))
		for _, value := range filter.Value {
			if value.IsNull() || value.IsUnknown() {
				continue
			}
			values = append(values, value.ValueString())
		}
		if len(values) > 0 {
			entry["Value"] = values
		}
		expanded = append(expanded, entry)
	}
	return expanded
}

func inventoryAggregators(aggregators []types.String) []string {
	expanded := make([]string, 0, len(aggregators))
	for _, aggregator := range aggregators {
		if aggregator.IsNull() || aggregator.IsUnknown() {
			continue
		}
		expanded = append(expanded, aggregator.ValueString())
	}
	return expanded
}

func inventoryResponseEntities(response map[string]interface{}) []interface{} {
	raw, ok := response["Entities"]
	if !ok || raw == nil {
		return nil
	}
	entities, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	return entities
}

func inventoryResponseNextToken(response map[string]interface{}) string {
	raw, ok := response["NextToken"]
	if !ok || raw == nil {
		return ""
	}
	nextToken, ok := raw.(string)
	if !ok {
		return ""
	}
	return nextToken
}

func inventoryItemsId(regionId string, maxResults int64, filters []inventoryItemsFilterModel, aggregators []types.String) string {
	parts := []string{regionId, fmt.Sprintf("%d", maxResults)}
	for _, filter := range filters {
		values := make([]string, 0, len(filter.Value))
		for _, value := range filter.Value {
			values = append(values, value.ValueString())
		}
		parts = append(parts, fmt.Sprintf("%s|%s|%s", filter.Name.ValueString(), filter.Operator.ValueString(), strings.Join(values, ",")))
	}
	for _, aggregator := range aggregators {
		parts = append(parts, aggregator.ValueString())
	}
	return fmt.Sprintf("%d", helper.Hashcode(strings.Join(parts, ";")))
}
