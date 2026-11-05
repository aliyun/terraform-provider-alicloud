package esa

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const loadBalancerRegionsTypeName = "alicloud_esa_load_balancer_regions"

const (
	loadBalancerRegionsAction   = "ListLoadBalancerRegions"
	loadBalancerRegionsPageSize = 100
	// loadBalancerRegionsMaxPages bounds the pagination loop. The region
	// catalog is small; hitting the bound means the endpoint ignored the page
	// parameters, and returning the accumulated rows would silently be partial.
	loadBalancerRegionsMaxPages = 1000
)

var (
	_ datasource.DataSource              = &loadBalancerRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerRegionsDataSource{}
)

func NewLoadBalancerRegionsDataSource() datasource.DataSource {
	return &loadBalancerRegionsDataSource{}
}

type loadBalancerRegionsDataSource struct {
	fwadapt.DataSourceBase
}

type loadBalancerRegionsModel struct {
	Id      types.String                  `tfsdk:"id"`
	Regions []loadBalancerRegionItemModel `tfsdk:"regions"`
}

type loadBalancerRegionItemModel struct {
	RegionCode   types.String                     `tfsdk:"region_code"`
	RegionCnName types.String                     `tfsdk:"region_cn_name"`
	RegionEnName types.String                     `tfsdk:"region_en_name"`
	SubRegions   []loadBalancerSubRegionItemModel `tfsdk:"sub_regions"`
}

type loadBalancerSubRegionItemModel struct {
	SubRegionCode   types.String `tfsdk:"sub_region_code"`
	SubRegionCnName types.String `tfsdk:"sub_region_cn_name"`
	SubRegionEnName types.String `tfsdk:"sub_region_en_name"`
}

var loadBalancerSubRegionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"sub_region_code":    types.StringType,
	"sub_region_cn_name": types.StringType,
	"sub_region_en_name": types.StringType,
}}

var loadBalancerRegionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"region_code":    types.StringType,
	"region_cn_name": types.StringType,
	"region_en_name": types.StringType,
	"sub_regions":    types.ListType{ElemType: loadBalancerSubRegionObjectType},
}}

func (d *loadBalancerRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = loadBalancerRegionsTypeName
}

func (d *loadBalancerRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the schedulable region catalog of ESA load balancers. The catalog enumerates the first-level regions and their second-level sub-regions and requires no site or load balancer context.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the returned region codes so that it stays stable across refreshes.",
			},
			"regions": schema.ListAttribute{
				Computed:    true,
				ElementType: loadBalancerRegionObjectType,
				Description: "The region list. Each element contains the first-level region code, Chinese and English names, and the second-level sub-region list.",
			},
		},
	}
}

func (d *loadBalancerRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model loadBalancerRegionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	items, err := d.listLoadBalancerRegionsAll(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", loadBalancerRegionsTypeName, loadBalancerRegionsAction),
			err.Error(),
		)
		return
	}

	idParts := make([]string, 0, len(items))
	for _, item := range items {
		idParts = append(idParts, item.RegionCode.ValueString())
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(idParts))
	model.Regions = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// listLoadBalancerRegionsAll reads every page. ListLoadBalancerRegions reports
// TotalPage; the loop also stops at the first page shorter than PageSize so an
// endpoint that omits the total cannot run unbounded.
func (d *loadBalancerRegionsDataSource) listLoadBalancerRegionsAll(ctx context.Context) ([]loadBalancerRegionItemModel, error) {
	items := make([]loadBalancerRegionItemModel, 0)
	for pageNumber := 1; pageNumber <= loadBalancerRegionsMaxPages; pageNumber++ {
		pageItems, totalPage, err := d.listLoadBalancerRegionsPage(ctx, pageNumber)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if totalPage > 0 {
			if pageNumber >= totalPage {
				return items, nil
			}
			continue
		}
		if len(pageItems) < loadBalancerRegionsPageSize {
			return items, nil
		}
	}
	return nil, fmt.Errorf("%s returned a full page for %d consecutive requests (PageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", loadBalancerRegionsAction, loadBalancerRegionsMaxPages, loadBalancerRegionsPageSize)
}

func (d *loadBalancerRegionsDataSource) listLoadBalancerRegionsPage(ctx context.Context, pageNumber int) ([]loadBalancerRegionItemModel, int, error) {
	request := map[string]interface{}{
		"PageNumber": pageNumber,
		"PageSize":   loadBalancerRegionsPageSize,
	}

	response, err := helper.RetryRpcPost(ctx, loadBalancerRegionsAction, func() (map[string]interface{}, error) {
		return d.Meta().RpcPost("ESA", "2024-09-10", loadBalancerRegionsAction, nil, request, true)
	}, 5*time.Minute, helper.NeedRetry)
	logs.AddDebug(loadBalancerRegionsAction, response, request)
	if err != nil {
		return nil, 0, err
	}

	totalPage := 0
	if v, ok := response["TotalPage"].(float64); ok {
		totalPage = int(v)
	}

	return loadBalancerRegionItems(response), totalPage, nil
}

// loadBalancerRegionItems parses the Regions array of the response.
func loadBalancerRegionItems(response map[string]interface{}) []loadBalancerRegionItemModel {
	items := make([]loadBalancerRegionItemModel, 0)
	if response == nil {
		return items
	}
	raw, ok := response["Regions"].([]interface{})
	if !ok {
		return items
	}
	for _, v := range raw {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, loadBalancerRegionItemModel{
			RegionCode:   types.StringValue(esaString(item, "RegionCode")),
			RegionCnName: types.StringValue(esaString(item, "RegionCnName")),
			RegionEnName: types.StringValue(esaString(item, "RegionEnName")),
			SubRegions:   loadBalancerSubRegionItems(item),
		})
	}
	return items
}

func loadBalancerSubRegionItems(item map[string]interface{}) []loadBalancerSubRegionItemModel {
	subRegions := make([]loadBalancerSubRegionItemModel, 0)
	raw, ok := item["SubRegions"].([]interface{})
	if !ok {
		return subRegions
	}
	for _, v := range raw {
		sub, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		subRegions = append(subRegions, loadBalancerSubRegionItemModel{
			SubRegionCode:   types.StringValue(esaString(sub, "SubRegionCode")),
			SubRegionCnName: types.StringValue(esaString(sub, "SubRegionCnName")),
			SubRegionEnName: types.StringValue(esaString(sub, "SubRegionEnName")),
		})
	}
	return subRegions
}

// esaString returns the string value of key in item, or "" when missing or not a
// string. ESA IDs are transmitted as integers but represented as strings across
// the provider, matching the existing ESA resources and data sources.
func esaString(item map[string]interface{}, key string) string {
	switch v := item[key].(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", int64(v))
	default:
		if v == nil {
			return ""
		}
		return fmt.Sprint(v)
	}
}
