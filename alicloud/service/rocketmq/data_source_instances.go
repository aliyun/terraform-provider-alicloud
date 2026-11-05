package rocketmq

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	instancesTypeName = "alicloud_rocketmq_instances"
	instancesPath     = "/instances"
)

var (
	_ datasource.DataSource              = &instancesDataSource{}
	_ datasource.DataSourceWithConfigure = &instancesDataSource{}
)

func NewInstancesDataSource() datasource.DataSource {
	return &instancesDataSource{}
}

type instancesDataSource struct {
	fwadapt.DataSourceBase
}

type instancesModel struct {
	Id              types.String    `tfsdk:"id"`
	Ids             types.List      `tfsdk:"ids"`
	NameRegex       types.String    `tfsdk:"name_regex"`
	Names           types.List      `tfsdk:"names"`
	Filter          types.String    `tfsdk:"filter"`
	ResourceGroupId types.String    `tfsdk:"resource_group_id"`
	SeriesCodes     types.List      `tfsdk:"series_codes"`
	Tags            types.Map       `tfsdk:"tags"`
	Instances       []instanceModel `tfsdk:"instances"`
}

type instanceModel struct {
	InstanceId      types.String               `tfsdk:"instance_id"`
	InstanceName    types.String               `tfsdk:"instance_name"`
	Status          types.String               `tfsdk:"status"`
	RegionId        types.String               `tfsdk:"region_id"`
	SeriesCode      types.String               `tfsdk:"series_code"`
	SubSeriesCode   types.String               `tfsdk:"sub_series_code"`
	PaymentType     types.String               `tfsdk:"payment_type"`
	CommodityCode   types.String               `tfsdk:"commodity_code"`
	ServiceCode     types.String               `tfsdk:"service_code"`
	ResourceGroupId types.String               `tfsdk:"resource_group_id"`
	Remark          types.String               `tfsdk:"remark"`
	TopicCount      types.Int64                `tfsdk:"topic_count"`
	GroupCount      types.Int64                `tfsdk:"group_count"`
	CreateTime      types.String               `tfsdk:"create_time"`
	UpdateTime      types.String               `tfsdk:"update_time"`
	StartTime       types.String               `tfsdk:"start_time"`
	ExpireTime      types.String               `tfsdk:"expire_time"`
	ReleaseTime     types.String               `tfsdk:"release_time"`
	Tags            []rocketmqTagModel         `tfsdk:"tags"`
	ProductInfo     []instanceProductInfoModel `tfsdk:"product_info"`
}

type instanceProductInfoModel struct {
	TraceOn      types.Bool   `tfsdk:"trace_on"`
	CapacityType types.String `tfsdk:"capacity_type"`
}

var instanceProductInfoObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"trace_on":      types.BoolType,
	"capacity_type": types.StringType,
}}

var instanceObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"instance_id":       types.StringType,
	"instance_name":     types.StringType,
	"status":            types.StringType,
	"region_id":         types.StringType,
	"series_code":       types.StringType,
	"sub_series_code":   types.StringType,
	"payment_type":      types.StringType,
	"commodity_code":    types.StringType,
	"service_code":      types.StringType,
	"resource_group_id": types.StringType,
	"remark":            types.StringType,
	"topic_count":       types.Int64Type,
	"group_count":       types.Int64Type,
	"create_time":       types.StringType,
	"update_time":       types.StringType,
	"start_time":        types.StringType,
	"expire_time":       types.StringType,
	"release_time":      types.StringType,
	"tags":              types.ListType{ElemType: rocketmqTagObjectType},
	"product_info":      types.ListType{ElemType: instanceProductInfoObjectType},
}}

func (d *instancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = instancesTypeName
}

func (d *instancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a collection of RocketMQ 5.0 instances in the current region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the filters so that it stays stable across refreshes.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of instance IDs to filter results. When omitted, every instance that matches the other filters is returned and this attribute lists the matched IDs.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by instance name. The filter is applied on the provider side after all pages are read.",
				Validators: []validator.String{
					rocketmqNameRegexValidator{},
				},
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of instance names that match the configured filters.",
			},
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: "The keyword that the server uses to filter instances, for example an instance ID. The value is sent to the server as the `filter` parameter and is not applied on the provider side.",
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the resource group that the instances belong to. The value is sent to the server as the `resourceGroupId` parameter and is also applied on the provider side, so a server that ignores the parameter cannot widen the result.",
			},
			"series_codes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A list of instance series codes to filter results. Valid values: `standard`, `ultimate`, `professional`, `dedicated` and `shared`. The values are sent to the server as the `seriesCodes` parameter and are also applied on the provider side, so a server that ignores the parameter cannot widen the result.",
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A map of tags assigned to the instances. An instance must carry every configured tag to be returned. The values are sent to the server as the `tags` parameter and are also applied on the provider side, so a server that ignores the parameter cannot widen the result.",
			},
			"instances": schema.ListAttribute{
				Computed:    true,
				ElementType: instanceObjectType,
				Description: "A list of instances that match the configured filters, sorted by instance ID.",
			},
		},
	}
}

func (d *instancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model instancesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if v := model.NameRegex.ValueString(); !model.NameRegex.IsNull() && v != "" {
		compiled, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", instancesTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = compiled
	}

	ids := make([]string, 0)
	if !model.Ids.IsNull() && !model.Ids.IsUnknown() {
		resp.Diagnostics.Append(model.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	seriesCodes := make([]string, 0)
	if !model.SeriesCodes.IsNull() && !model.SeriesCodes.IsUnknown() {
		resp.Diagnostics.Append(model.SeriesCodes.ElementsAs(ctx, &seriesCodes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	seriesCodeSet := make(map[string]struct{}, len(seriesCodes))
	for _, code := range seriesCodes {
		seriesCodeSet[code] = struct{}{}
	}

	tags := make(map[string]string)
	if !model.Tags.IsNull() && !model.Tags.IsUnknown() {
		resp.Diagnostics.Append(model.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	query := instancesQuery(&model, tags, seriesCodes)
	items, err := rocketmqCollectPages(instancesPath, func(item map[string]interface{}) string {
		return rocketmqString(item, "instanceId")
	}, func(pageNumber int) ([]map[string]interface{}, error) {
		return d.describeInstancesPage(ctx, query, pageNumber)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", instancesTypeName, instancesPath),
			err.Error(),
		)
		return
	}

	matched := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if len(idSet) > 0 {
			if _, ok := idSet[rocketmqString(item, "instanceId")]; !ok {
				continue
			}
		}
		if len(seriesCodeSet) > 0 {
			if _, ok := seriesCodeSet[rocketmqString(item, "seriesCode")]; !ok {
				continue
			}
		}
		if resourceGroupId := model.ResourceGroupId.ValueString(); !model.ResourceGroupId.IsNull() && resourceGroupId != "" {
			if rocketmqString(item, "resourceGroupId") != resourceGroupId {
				continue
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(rocketmqString(item, "instanceName")) {
			continue
		}
		if !rocketmqTagsMatch(item, tags) {
			continue
		}
		matched = append(matched, item)
	}
	sort.Slice(matched, func(i, j int) bool {
		return rocketmqString(matched[i], "instanceId") < rocketmqString(matched[j], "instanceId")
	})

	matchedIds := make([]string, 0, len(matched))
	names := make([]string, 0, len(matched))
	instances := make([]instanceModel, 0, len(matched))
	for _, item := range matched {
		matchedIds = append(matchedIds, rocketmqString(item, "instanceId"))
		names = append(names, rocketmqString(item, "instanceName"))
		productInfo := make([]instanceProductInfoModel, 0)
		if raw, err := jsonpath.Get("$.productInfo", item); err == nil {
			if info, ok := raw.(map[string]interface{}); ok {
				productInfo = append(productInfo, instanceProductInfoModel{
					TraceOn:      rocketmqBool(info, "traceOn"),
					CapacityType: types.StringValue(rocketmqString(info, "capacityType")),
				})
			}
		}
		instances = append(instances, instanceModel{
			InstanceId:      types.StringValue(rocketmqString(item, "instanceId")),
			InstanceName:    types.StringValue(rocketmqString(item, "instanceName")),
			Status:          types.StringValue(rocketmqString(item, "status")),
			RegionId:        types.StringValue(rocketmqString(item, "regionId")),
			SeriesCode:      types.StringValue(rocketmqString(item, "seriesCode")),
			SubSeriesCode:   types.StringValue(rocketmqString(item, "subSeriesCode")),
			PaymentType:     types.StringValue(rocketmqString(item, "paymentType")),
			CommodityCode:   types.StringValue(rocketmqString(item, "commodityCode")),
			ServiceCode:     types.StringValue(rocketmqString(item, "serviceCode")),
			ResourceGroupId: types.StringValue(rocketmqString(item, "resourceGroupId")),
			Remark:          types.StringValue(rocketmqString(item, "remark")),
			TopicCount:      rocketmqInt64(item, "topicCount"),
			GroupCount:      rocketmqInt64(item, "groupCount"),
			CreateTime:      types.StringValue(rocketmqString(item, "createTime")),
			UpdateTime:      types.StringValue(rocketmqString(item, "updateTime")),
			StartTime:       types.StringValue(rocketmqString(item, "startTime")),
			ExpireTime:      types.StringValue(rocketmqString(item, "expireTime")),
			ReleaseTime:     types.StringValue(rocketmqString(item, "releaseTime")),
			Tags:            rocketmqTagsFromAPI(item["tags"]),
			ProductInfo:     productInfo,
		})
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, matchedIds)
	resp.Diagnostics.Append(diags...)
	namesList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(instancesId(&model, tags, ids, seriesCodes))
	model.Ids = idsList
	model.Names = namesList
	model.Instances = instances
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// instancesQuery builds the query shared by every page request. Array values use
// the indexed form that the RocketMQ ROA gateway accepts; every value is applied
// again on the provider side in Read.
func instancesQuery(model *instancesModel, tags map[string]string, seriesCodes []string) map[string]*string {
	query := map[string]*string{
		"pageSize": helper.StringPointer(strconv.Itoa(rocketmqDataSourcePageSize)),
	}
	if v := model.Filter.ValueString(); !model.Filter.IsNull() && v != "" {
		query["filter"] = helper.StringPointer(v)
	}
	if v := model.ResourceGroupId.ValueString(); !model.ResourceGroupId.IsNull() && v != "" {
		query["resourceGroupId"] = helper.StringPointer(v)
	}
	if len(tags) > 0 {
		query["tags"] = helper.StringPointer(rocketmqTagsRequestValue(tags))
	}
	for i, code := range seriesCodes {
		query[fmt.Sprintf("seriesCodes.%d", i+1)] = helper.StringPointer(code)
	}
	return query
}

func (d *instancesDataSource) describeInstancesPage(ctx context.Context, query map[string]*string, pageNumber int) ([]map[string]interface{}, error) {
	requestQuery := make(map[string]*string, len(query)+1)
	for k, v := range query {
		requestQuery[k] = v
	}
	requestQuery["pageNumber"] = helper.StringPointer(strconv.Itoa(pageNumber))

	client := d.Meta()
	var response map[string]interface{}
	// Cancellation can return before the RPC callback finishes.
	var responseMu sync.Mutex
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := sdkretry.RetryContext(ctx, 5*time.Minute, func() *sdkretry.RetryError {
		result, rpcErr := client.RoaGet("RocketMQ", "2022-08-01", instancesPath, requestQuery, nil, nil)
		responseMu.Lock()
		response = result
		responseMu.Unlock()
		if rpcErr != nil {
			if helper.NeedRetry(rpcErr) {
				wait()
				return sdkretry.RetryableError(rpcErr)
			}
			return sdkretry.NonRetryableError(rpcErr)
		}
		return nil
	})
	responseMu.Lock()
	defer responseMu.Unlock()
	logs.AddDebug(instancesPath, response, requestQuery)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, instancesTypeName, instancesPath, errs.AlibabaCloudSdkGoERROR)
	}
	return rocketmqListItems(response)
}

func instancesId(model *instancesModel, tags map[string]string, ids, seriesCodes []string) string {
	sortedIds := append([]string(nil), ids...)
	sort.Strings(sortedIds)
	sortedSeriesCodes := append([]string(nil), seriesCodes...)
	sort.Strings(sortedSeriesCodes)
	return rocketmqDataSourceId(
		model.Filter.ValueString(),
		model.NameRegex.ValueString(),
		model.ResourceGroupId.ValueString(),
		rocketmqTagsRequestValue(tags),
		strings.Join(sortedIds, ","),
		strings.Join(sortedSeriesCodes, ","),
	)
}
