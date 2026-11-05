package esa

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const loadBalancerOriginStatusTypeName = "alicloud_esa_load_balancer_origin_status"

var (
	_ datasource.DataSource              = &loadBalancerOriginStatusDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerOriginStatusDataSource{}
)

func NewLoadBalancerOriginStatusDataSource() datasource.DataSource {
	return &loadBalancerOriginStatusDataSource{}
}

type loadBalancerOriginStatusDataSource struct {
	fwadapt.DataSourceBase
}

type loadBalancerOriginStatusModel struct {
	SiteId          types.String                        `tfsdk:"site_id"`
	LoadBalancerIds types.List                          `tfsdk:"load_balancer_ids"`
	PoolType        types.String                        `tfsdk:"pool_type"`
	Id              types.String                        `tfsdk:"id"`
	OriginStatus    []loadBalancerOriginStatusItemModel `tfsdk:"origin_status"`
}

type loadBalancerOriginStatusItemModel struct {
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	OriginId       types.String `tfsdk:"origin_id"`
	PoolId         types.String `tfsdk:"pool_id"`
	PoolType       types.String `tfsdk:"pool_type"`
	Reason         types.String `tfsdk:"reason"`
	Status         types.String `tfsdk:"status"`
}

var loadBalancerOriginStatusItemObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"load_balancer_id": types.StringType,
	"origin_id":        types.StringType,
	"pool_id":          types.StringType,
	"pool_type":        types.StringType,
	"reason":           types.StringType,
	"status":           types.StringType,
}}

func (d *loadBalancerOriginStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = loadBalancerOriginStatusTypeName
}

func (d *loadBalancerOriginStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the origin health status of the specified ESA load balancers. The status reflects the probe results of origins under the default address pool. Load balancer IDs can be obtained from the `alicloud_esa_load_balancer` resource or the `ListLoadBalancers` API.",
		Attributes: map[string]schema.Attribute{
			"site_id": schema.StringAttribute{
				Required:    true,
				Description: "The site ID. Load balancer origin status is queried within the scope of a site.",
			},
			"load_balancer_ids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "One or more load balancer IDs. The API accepts up to 100 IDs per request.",
			},
			"pool_type": schema.StringAttribute{
				Optional:    true,
				Description: "The origin address pool type. Only the `default_pool` origins affect the load balancer status. Pass `default_pool` to query the origins under the default address pool only.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the query conditions so that it stays stable across refreshes.",
			},
			"origin_status": schema.ListAttribute{
				Computed:    true,
				ElementType: loadBalancerOriginStatusItemObjectType,
				Description: "The origin status list of the load balancers. Each element contains the load balancer ID, origin ID, pool ID, pool type, probe reason and status.",
			},
		},
	}
}

func (d *loadBalancerOriginStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const action = "ListLoadBalancerOriginStatus"

	var model loadBalancerOriginStatusModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var loadBalancerIds []string
	resp.Diagnostics.Append(model.LoadBalancerIds.ElementsAs(ctx, &loadBalancerIds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"SiteId":          model.SiteId.ValueString(),
		"LoadBalancerIds": strings.Join(loadBalancerIds, ","),
	}
	if v := model.PoolType.ValueString(); !model.PoolType.IsNull() && !model.PoolType.IsUnknown() && v != "" {
		request["PoolType"] = v
	}

	response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return d.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
	}, 5*time.Minute, helper.NeedRetry)
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", loadBalancerOriginStatusTypeName, action),
			err.Error(),
		)
		return
	}

	items := loadBalancerOriginStatusItems(response)
	state := model
	state.Id = types.StringValue(helper.DataResourceIdHash([]string{
		model.SiteId.ValueString(),
		strings.Join(loadBalancerIds, ","),
		model.PoolType.ValueString(),
	}))
	state.OriginStatus = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// loadBalancerOriginStatusItems parses the OriginStatus array of the response.
// A missing key or a non-array value yields an empty slice; it is not an error.
func loadBalancerOriginStatusItems(response map[string]interface{}) []loadBalancerOriginStatusItemModel {
	items := make([]loadBalancerOriginStatusItemModel, 0)
	if response == nil {
		return items
	}
	raw, ok := response["OriginStatus"].([]interface{})
	if !ok {
		return items
	}
	for _, v := range raw {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, loadBalancerOriginStatusItemModel{
			LoadBalancerId: types.StringValue(esaString(item, "LoadBalancerId")),
			OriginId:       types.StringValue(esaString(item, "OriginId")),
			PoolId:         types.StringValue(esaString(item, "PoolId")),
			PoolType:       types.StringValue(esaString(item, "PoolType")),
			Reason:         types.StringValue(esaString(item, "Reason")),
			Status:         types.StringValue(esaString(item, "Status")),
		})
	}
	return items
}
