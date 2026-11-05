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

const topicsTypeName = "alicloud_rocketmq_topics"

var (
	_ datasource.DataSource              = &topicsDataSource{}
	_ datasource.DataSourceWithConfigure = &topicsDataSource{}
)

func NewTopicsDataSource() datasource.DataSource {
	return &topicsDataSource{}
}

type topicsDataSource struct {
	fwadapt.DataSourceBase
}

type topicsModel struct {
	Id           types.String `tfsdk:"id"`
	InstanceId   types.String `tfsdk:"instance_id"`
	Filter       types.String `tfsdk:"filter"`
	MessageTypes types.List   `tfsdk:"message_types"`
	NameRegex    types.String `tfsdk:"name_regex"`
	Ids          types.List   `tfsdk:"ids"`
	Names        types.List   `tfsdk:"names"`
	Topics       []topicModel `tfsdk:"topics"`
}

type topicModel struct {
	TopicName           types.String `tfsdk:"topic_name"`
	MessageType         types.String `tfsdk:"message_type"`
	Status              types.String `tfsdk:"status"`
	Remark              types.String `tfsdk:"remark"`
	InstanceId          types.String `tfsdk:"instance_id"`
	RegionId            types.String `tfsdk:"region_id"`
	CreateTime          types.String `tfsdk:"create_time"`
	UpdateTime          types.String `tfsdk:"update_time"`
	MaxSendTps          types.Int64  `tfsdk:"max_send_tps"`
	LiteTopicExpiration types.Int64  `tfsdk:"lite_topic_expiration"`
}

var topicObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"topic_name":            types.StringType,
	"message_type":          types.StringType,
	"status":                types.StringType,
	"remark":                types.StringType,
	"instance_id":           types.StringType,
	"region_id":             types.StringType,
	"create_time":           types.StringType,
	"update_time":           types.StringType,
	"max_send_tps":          types.Int64Type,
	"lite_topic_expiration": types.Int64Type,
}}

func (d *topicsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = topicsTypeName
}

func (d *topicsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a collection of topics of a RocketMQ 5.0 instance.",
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
				Description: "The keyword that the server uses to filter topics. The value is sent to the server as the `filter` parameter and is not applied on the provider side.",
			},
			"message_types": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A list of message types to filter results. Valid values: `NORMAL`, `FIFO`, `DELAY`, `TRANSACTION` and `LITE`. The values are sent to the server as the `messageTypes` parameter and are also applied on the provider side, so a server that ignores the parameter cannot widen the result.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by topic name. The filter is applied on the provider side after all pages are read.",
				Validators: []validator.String{
					rocketmqNameRegexValidator{},
				},
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of topic names to filter results. When omitted, every topic that matches the other filters is returned and this attribute lists the matched topic names.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of topic names that match the configured filters.",
			},
			"topics": schema.ListAttribute{
				Computed:    true,
				ElementType: topicObjectType,
				Description: "A list of topics that match the configured filters, sorted by topic name.",
			},
		},
	}
}

func (d *topicsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model topicsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if v := model.NameRegex.ValueString(); !model.NameRegex.IsNull() && v != "" {
		compiled, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", topicsTypeName),
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

	messageTypes := make([]string, 0)
	if !model.MessageTypes.IsNull() && !model.MessageTypes.IsUnknown() {
		resp.Diagnostics.Append(model.MessageTypes.ElementsAs(ctx, &messageTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	messageTypeSet := make(map[string]struct{}, len(messageTypes))
	for _, messageType := range messageTypes {
		messageTypeSet[messageType] = struct{}{}
	}

	path := fmt.Sprintf("/instances/%s/topics", model.InstanceId.ValueString())
	query := topicsQuery(&model, messageTypes)
	items, err := rocketmqCollectPages(path, func(item map[string]interface{}) string {
		return rocketmqString(item, "topicName")
	}, func(pageNumber int) ([]map[string]interface{}, error) {
		return d.describeTopicsPage(ctx, path, query, pageNumber)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", topicsTypeName, path),
			err.Error(),
		)
		return
	}

	matched := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if len(idSet) > 0 {
			if _, ok := idSet[rocketmqString(item, "topicName")]; !ok {
				continue
			}
		}
		if len(messageTypeSet) > 0 {
			if _, ok := messageTypeSet[rocketmqString(item, "messageType")]; !ok {
				continue
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(rocketmqString(item, "topicName")) {
			continue
		}
		matched = append(matched, item)
	}
	sort.Slice(matched, func(i, j int) bool {
		return rocketmqString(matched[i], "topicName") < rocketmqString(matched[j], "topicName")
	})

	matchedIds := make([]string, 0, len(matched))
	topics := make([]topicModel, 0, len(matched))
	for _, item := range matched {
		matchedIds = append(matchedIds, rocketmqString(item, "topicName"))
		topics = append(topics, topicModel{
			TopicName:           types.StringValue(rocketmqString(item, "topicName")),
			MessageType:         types.StringValue(rocketmqString(item, "messageType")),
			Status:              types.StringValue(rocketmqString(item, "status")),
			Remark:              types.StringValue(rocketmqString(item, "remark")),
			InstanceId:          types.StringValue(rocketmqString(item, "instanceId")),
			RegionId:            types.StringValue(rocketmqString(item, "regionId")),
			CreateTime:          types.StringValue(rocketmqString(item, "createTime")),
			UpdateTime:          types.StringValue(rocketmqString(item, "updateTime")),
			MaxSendTps:          rocketmqInt64(item, "maxSendTps"),
			LiteTopicExpiration: rocketmqInt64(item, "liteTopicExpiration"),
		})
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, matchedIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(topicsId(&model, ids, messageTypes))
	model.Ids = idsList
	model.Names = idsList
	model.Topics = topics
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// topicsQuery builds the query shared by every page request. Array values use
// the indexed form that the RocketMQ ROA gateway accepts; every value is applied
// again on the provider side in Read.
func topicsQuery(model *topicsModel, messageTypes []string) map[string]*string {
	query := map[string]*string{
		"pageSize": helper.StringPointer(strconv.Itoa(rocketmqDataSourcePageSize)),
	}
	if v := model.Filter.ValueString(); !model.Filter.IsNull() && v != "" {
		query["filter"] = helper.StringPointer(v)
	}
	for i, messageType := range messageTypes {
		query[fmt.Sprintf("messageTypes.%d", i+1)] = helper.StringPointer(messageType)
	}
	return query
}

func (d *topicsDataSource) describeTopicsPage(ctx context.Context, path string, query map[string]*string, pageNumber int) ([]map[string]interface{}, error) {
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
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, topicsTypeName, path, errs.AlibabaCloudSdkGoERROR)
	}
	return rocketmqListItems(response)
}

func topicsId(model *topicsModel, ids, messageTypes []string) string {
	sortedIds := append([]string(nil), ids...)
	sort.Strings(sortedIds)
	sortedMessageTypes := append([]string(nil), messageTypes...)
	sort.Strings(sortedMessageTypes)
	return rocketmqDataSourceId(
		model.InstanceId.ValueString(),
		model.Filter.ValueString(),
		model.NameRegex.ValueString(),
		strings.Join(sortedIds, ","),
		strings.Join(sortedMessageTypes, ","),
	)
}
