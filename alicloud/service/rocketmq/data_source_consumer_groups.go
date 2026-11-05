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

const consumerGroupsTypeName = "alicloud_rocketmq_consumer_groups"

var (
	_ datasource.DataSource              = &consumerGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &consumerGroupsDataSource{}
)

func NewConsumerGroupsDataSource() datasource.DataSource {
	return &consumerGroupsDataSource{}
}

type consumerGroupsDataSource struct {
	fwadapt.DataSourceBase
}

type consumerGroupsModel struct {
	Id         types.String         `tfsdk:"id"`
	InstanceId types.String         `tfsdk:"instance_id"`
	Filter     types.String         `tfsdk:"filter"`
	NameRegex  types.String         `tfsdk:"name_regex"`
	Ids        types.List           `tfsdk:"ids"`
	Names      types.List           `tfsdk:"names"`
	Groups     []consumerGroupModel `tfsdk:"groups"`
}

type consumerGroupModel struct {
	ConsumerGroupId types.String `tfsdk:"consumer_group_id"`
	Status          types.String `tfsdk:"status"`
	Remark          types.String `tfsdk:"remark"`
	InstanceId      types.String `tfsdk:"instance_id"`
	RegionId        types.String `tfsdk:"region_id"`
	MessageModel    types.String `tfsdk:"message_model"`
	TopicName       types.String `tfsdk:"topic_name"`
	CreateTime      types.String `tfsdk:"create_time"`
	UpdateTime      types.String `tfsdk:"update_time"`
	MaxReceiveTps   types.Int64  `tfsdk:"max_receive_tps"`
}

var consumerGroupObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"consumer_group_id": types.StringType,
	"status":            types.StringType,
	"remark":            types.StringType,
	"instance_id":       types.StringType,
	"region_id":         types.StringType,
	"message_model":     types.StringType,
	"topic_name":        types.StringType,
	"create_time":       types.StringType,
	"update_time":       types.StringType,
	"max_receive_tps":   types.Int64Type,
}}

func (d *consumerGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = consumerGroupsTypeName
}

func (d *consumerGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a collection of consumer groups of a RocketMQ 5.0 instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the instance and filters so that it stays stable across refreshes.",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the RocketMQ 5.0 instance.",
			},
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: "The keyword that the server uses to filter consumer groups. The value is sent to the server as the `filter` parameter and is not applied on the provider side.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by consumer group ID. The filter is applied on the provider side after all pages are read.",
				Validators: []validator.String{
					rocketmqNameRegexValidator{},
				},
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of consumer group IDs to filter results. When omitted, every consumer group that matches the other filters is returned and this attribute lists the matched consumer group IDs.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of consumer group IDs that match the configured filters.",
			},
			"groups": schema.ListAttribute{
				Computed:    true,
				ElementType: consumerGroupObjectType,
				Description: "A list of consumer groups that match the configured filters, sorted by consumer group ID.",
			},
		},
	}
}

func (d *consumerGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model consumerGroupsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if v := model.NameRegex.ValueString(); !model.NameRegex.IsNull() && v != "" {
		compiled, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", consumerGroupsTypeName),
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

	path := fmt.Sprintf("/instances/%s/consumerGroups", model.InstanceId.ValueString())
	query := consumerGroupsQuery(&model)
	items, err := rocketmqCollectPages(path, func(item map[string]interface{}) string {
		return rocketmqString(item, "consumerGroupId")
	}, func(pageNumber int) ([]map[string]interface{}, error) {
		return d.describeConsumerGroupsPage(ctx, path, query, pageNumber)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", consumerGroupsTypeName, path),
			err.Error(),
		)
		return
	}

	matched := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if len(idSet) > 0 {
			if _, ok := idSet[rocketmqString(item, "consumerGroupId")]; !ok {
				continue
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(rocketmqString(item, "consumerGroupId")) {
			continue
		}
		matched = append(matched, item)
	}
	sort.Slice(matched, func(i, j int) bool {
		return rocketmqString(matched[i], "consumerGroupId") < rocketmqString(matched[j], "consumerGroupId")
	})

	matchedIds := make([]string, 0, len(matched))
	groups := make([]consumerGroupModel, 0, len(matched))
	for _, item := range matched {
		matchedIds = append(matchedIds, rocketmqString(item, "consumerGroupId"))
		groups = append(groups, consumerGroupModel{
			ConsumerGroupId: types.StringValue(rocketmqString(item, "consumerGroupId")),
			Status:          types.StringValue(rocketmqString(item, "status")),
			Remark:          types.StringValue(rocketmqString(item, "remark")),
			InstanceId:      types.StringValue(rocketmqString(item, "instanceId")),
			RegionId:        types.StringValue(rocketmqString(item, "regionId")),
			MessageModel:    types.StringValue(rocketmqString(item, "messageModel")),
			TopicName:       types.StringValue(rocketmqString(item, "topicName")),
			CreateTime:      types.StringValue(rocketmqString(item, "createTime")),
			UpdateTime:      types.StringValue(rocketmqString(item, "updateTime")),
			MaxReceiveTps:   rocketmqInt64(item, "maxReceiveTps"),
		})
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, matchedIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(consumerGroupsId(&model, ids))
	model.Ids = idsList
	model.Names = idsList
	model.Groups = groups
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func consumerGroupsQuery(model *consumerGroupsModel) map[string]*string {
	query := map[string]*string{
		"pageSize": helper.StringPointer(strconv.Itoa(rocketmqDataSourcePageSize)),
	}
	if v := model.Filter.ValueString(); !model.Filter.IsNull() && v != "" {
		query["filter"] = helper.StringPointer(v)
	}
	return query
}

func (d *consumerGroupsDataSource) describeConsumerGroupsPage(ctx context.Context, path string, query map[string]*string, pageNumber int) ([]map[string]interface{}, error) {
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
		result, rpcErr := client.RoaGet("RocketMQ", "2022-08-01", path, requestQuery, nil, nil)
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
	logs.AddDebug(path, response, requestQuery)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, consumerGroupsTypeName, path, errs.AlibabaCloudSdkGoERROR)
	}
	return rocketmqListItems(response)
}

func consumerGroupsId(model *consumerGroupsModel, ids []string) string {
	sortedIds := append([]string(nil), ids...)
	sort.Strings(sortedIds)
	return rocketmqDataSourceId(
		model.InstanceId.ValueString(),
		model.Filter.ValueString(),
		model.NameRegex.ValueString(),
		strings.Join(sortedIds, ","),
	)
}
