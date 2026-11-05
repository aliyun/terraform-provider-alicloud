package cloudmonitorservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const cmsSubscriptionTypeName = "alicloud_cms_subscription"

// terraform-plugin-framework has no user-configurable resource timeouts yet,
// so the SDKv2 default of five minutes per operation is kept as a constant.
const cmsSubscriptionTimeout = 5 * time.Minute

var (
	cmsSubscriptionConditionType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"field": types.StringType,
		"op":    types.StringType,
		"value": types.StringType,
	}}

	cmsSubscriptionFilterSettingType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"conditions": types.ListType{ElemType: cmsSubscriptionConditionType},
		"expression": types.StringType,
		"relation":   types.StringType,
	}}

	cmsSubscriptionPushingSettingType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"alert_action_ids":   types.ListType{ElemType: types.StringType},
		"restore_action_ids": types.ListType{ElemType: types.StringType},
		"template_uuid":      types.StringType,
		"response_plan_id":   types.StringType,
	}}

	cmsSubscriptionWorkspaceFilterSettingType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"workspace_uuids": types.ListType{ElemType: types.StringType},
	}}

	cmsSubscriptionAgentChannelType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"channel_type":         types.StringType,
		"receivers":            types.ListType{ElemType: types.StringType},
		"enabled_sub_channels": types.ListType{ElemType: types.StringType},
	}}

	cmsSubscriptionAgentEffectTimeRangeType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"time_zone":            types.StringType,
		"day_in_week":          types.ListType{ElemType: types.Int64Type},
		"start_time_in_minute": types.Int64Type,
		"end_time_in_minute":   types.Int64Type,
	}}

	cmsSubscriptionAgentRouteType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"channels":          types.ListType{ElemType: cmsSubscriptionAgentChannelType},
		"effect_time_range": types.ListType{ElemType: cmsSubscriptionAgentEffectTimeRangeType},
	}}

	cmsSubscriptionAgentConfigType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"agent_uuid": types.StringType,
		"routes":     types.ListType{ElemType: cmsSubscriptionAgentRouteType},
	}}
)

var (
	_ resource.Resource                = &cmsSubscriptionResource{}
	_ resource.ResourceWithImportState = &cmsSubscriptionResource{}
)

func NewCmsSubscriptionResource() resource.Resource {
	return &cmsSubscriptionResource{}
}

type cmsSubscriptionResource struct {
	fwadapt.ResourceBase
}

type cmsSubscriptionResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	SubscriptionName       types.String `tfsdk:"subscription_name"`
	Description            types.String `tfsdk:"description"`
	NotifyStrategyId       types.String `tfsdk:"notify_strategy_id"`
	Workspace              types.String `tfsdk:"workspace"`
	Enable                 types.Bool   `tfsdk:"enable"`
	SubscribeLegacyEvent   types.Bool   `tfsdk:"subscribe_legacy_event"`
	FilterSetting          types.List   `tfsdk:"filter_setting"`
	PushingSetting         types.List   `tfsdk:"pushing_setting"`
	WorkspaceFilterSetting types.List   `tfsdk:"workspace_filter_setting"`
	AgentConfig            types.List   `tfsdk:"agent_config"`
	UserId                 types.String `tfsdk:"user_id"`
	SubscriptionType       types.String `tfsdk:"subscription_type"`
	RegionId               types.String `tfsdk:"region_id"`
	CreateTime             types.String `tfsdk:"create_time"`
	UpdateTime             types.String `tfsdk:"update_time"`
}

func (r *cmsSubscriptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = cmsSubscriptionTypeName
}

func (r *cmsSubscriptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subscription_name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"notify_strategy_id": schema.StringAttribute{
				Optional: true,
			},
			"workspace": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"subscribe_legacy_event": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"filter_setting": schema.ListAttribute{
				Optional:    true,
				ElementType: cmsSubscriptionFilterSettingType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"pushing_setting": schema.ListAttribute{
				Optional:    true,
				ElementType: cmsSubscriptionPushingSettingType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"workspace_filter_setting": schema.ListAttribute{
				Optional:    true,
				ElementType: cmsSubscriptionWorkspaceFilterSettingType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"agent_config": schema.ListAttribute{
				Optional:    true,
				ElementType: cmsSubscriptionAgentConfigType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"user_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subscription_type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *cmsSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "/subscriptions"
	request := expandCmsSubscriptionRequest(plan)
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(cmsSubscriptionTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPost("Cms", "2024-03-30", action, cmsSubscriptionQuery(plan.Workspace), nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id, _ := response["subscriptionId"].(string)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSubscription", cmsSubscriptionTypeName),
			"CreateSubscription returned an empty subscriptionId; the subscription was created but cannot be tracked in state.",
		)
		return
	}

	if !plan.Enable.IsNull() && !plan.Enable.IsUnknown() && !plan.Enable.ValueBool() {
		if err := r.setCmsSubscriptionEnable(id, plan.Workspace, false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling DisableSubscription", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id+"/disable", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	state := cmsSubscriptionResourceModel{
		Id:               types.StringValue(id),
		SubscriptionName: plan.SubscriptionName,
		Workspace:        plan.Workspace,
		Enable:           plan.Enable,
	}
	object, err := NewService(r.Meta()).DescribeCmsSubscription(id, plan.Workspace.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, cmsSubscriptionState(object, state))...)
}

func (r *cmsSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeCmsSubscription(state.Id.ValueString(), state.Workspace.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+state.Id.ValueString(), errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, cmsSubscriptionState(object, state))...)
}

func (r *cmsSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if !plan.Enable.Equal(state.Enable) {
		if err := r.setCmsSubscriptionEnable(id, state.Workspace, plan.Enable.ValueBool()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: toggling subscription enable state", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id+"/enable", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	action := fmt.Sprintf("/subscriptions/%s", id)
	request := expandCmsSubscriptionRequest(plan)
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(cmsSubscriptionTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPut("Cms", "2024-03-30", action, cmsSubscriptionQuery(state.Workspace), nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeCmsSubscription(id, state.Workspace.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	updated := cmsSubscriptionState(object, state)
	updated.Enable = plan.Enable
	resp.Diagnostics.Append(resp.State.Set(ctx, updated)...)
}

func (r *cmsSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := fmt.Sprintf("/subscriptions/%s", state.Id.ValueString())
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(cmsSubscriptionTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaDelete("Cms", "2024-03-30", action, cmsSubscriptionQuery(state.Workspace), nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, cmsSubscriptionQuery(state.Workspace))
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
	}
}

func (r *cmsSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *cmsSubscriptionResource) setCmsSubscriptionEnable(id string, workspace types.String, enable bool) error {
	suffix := "enable"
	if !enable {
		suffix = "disable"
	}
	action := fmt.Sprintf("/subscriptions/%s/%s", id, suffix)

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(cmsSubscriptionTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPut("Cms", "2024-03-30", action, cmsSubscriptionQuery(workspace), nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, action)
	return err
}

func cmsSubscriptionQuery(workspace types.String) map[string]*string {
	query := make(map[string]*string)
	if !workspace.IsNull() && !workspace.IsUnknown() && workspace.ValueString() != "" {
		query["workspace"] = helper.StringPointer(workspace.ValueString())
	}
	return query
}

// expandCmsSubscriptionRequest always returns the complete body: UpdateSubscription
// replaces the stored configuration, so unchanged attributes must be re-sent.
func expandCmsSubscriptionRequest(plan cmsSubscriptionResourceModel) map[string]interface{} {
	request := map[string]interface{}{
		"subscriptionName": plan.SubscriptionName.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request["description"] = plan.Description.ValueString()
	}
	if !plan.NotifyStrategyId.IsNull() && !plan.NotifyStrategyId.IsUnknown() {
		request["notifyStrategyId"] = plan.NotifyStrategyId.ValueString()
	}
	if !plan.SubscribeLegacyEvent.IsNull() && !plan.SubscribeLegacyEvent.IsUnknown() {
		request["subscribeLegacyEvent"] = plan.SubscribeLegacyEvent.ValueBool()
	}
	if v := expandCmsSubscriptionFilterSetting(plan.FilterSetting); v != nil {
		request["filterSetting"] = v
	}
	if v := expandCmsSubscriptionPushingSetting(plan.PushingSetting); v != nil {
		request["pushingSetting"] = v
	}
	if v := expandCmsSubscriptionWorkspaceFilterSetting(plan.WorkspaceFilterSetting); v != nil {
		request["workspaceFilterSetting"] = v
	}
	if v := expandCmsSubscriptionAgentConfig(plan.AgentConfig); v != nil {
		request["agentConfig"] = v
	}
	return request
}

func expandCmsSubscriptionFilterSetting(l types.List) map[string]interface{} {
	object, ok := cmsSubscriptionFirstObject(l)
	if !ok {
		return nil
	}
	attrs := object.Attributes()
	result := map[string]interface{}{}
	if v := cmsSubscriptionAttrString(attrs["expression"]); v != "" {
		result["expression"] = v
	}
	if v := cmsSubscriptionAttrString(attrs["relation"]); v != "" {
		result["relation"] = v
	}
	if conditions, ok := attrs["conditions"].(types.List); ok && !conditions.IsNull() && !conditions.IsUnknown() {
		items := make([]map[string]interface{}, 0, len(conditions.Elements()))
		for _, element := range conditions.Elements() {
			condition, ok := element.(types.Object)
			if !ok || condition.IsNull() || condition.IsUnknown() {
				continue
			}
			conditionAttrs := condition.Attributes()
			items = append(items, map[string]interface{}{
				"field": cmsSubscriptionAttrString(conditionAttrs["field"]),
				"op":    cmsSubscriptionAttrString(conditionAttrs["op"]),
				"value": cmsSubscriptionAttrString(conditionAttrs["value"]),
			})
		}
		if len(items) > 0 {
			result["conditions"] = items
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandCmsSubscriptionPushingSetting(l types.List) map[string]interface{} {
	object, ok := cmsSubscriptionFirstObject(l)
	if !ok {
		return nil
	}
	attrs := object.Attributes()
	result := map[string]interface{}{}
	if v := cmsSubscriptionAttrStringList(attrs["alert_action_ids"]); len(v) > 0 {
		result["alertActionIds"] = v
	}
	if v := cmsSubscriptionAttrStringList(attrs["restore_action_ids"]); len(v) > 0 {
		result["restoreActionIds"] = v
	}
	if v := cmsSubscriptionAttrString(attrs["template_uuid"]); v != "" {
		result["templateUuid"] = v
	}
	if v := cmsSubscriptionAttrString(attrs["response_plan_id"]); v != "" {
		result["responsePlanId"] = v
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandCmsSubscriptionWorkspaceFilterSetting(l types.List) map[string]interface{} {
	object, ok := cmsSubscriptionFirstObject(l)
	if !ok {
		return nil
	}
	attrs := object.Attributes()
	result := map[string]interface{}{}
	if v := cmsSubscriptionAttrStringList(attrs["workspace_uuids"]); len(v) > 0 {
		result["workspaceUuids"] = v
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandCmsSubscriptionAgentConfig(l types.List) map[string]interface{} {
	object, ok := cmsSubscriptionFirstObject(l)
	if !ok {
		return nil
	}
	attrs := object.Attributes()
	result := map[string]interface{}{}
	if v := cmsSubscriptionAttrString(attrs["agent_uuid"]); v != "" {
		result["agentUuid"] = v
	}
	if routes, ok := attrs["routes"].(types.List); ok && !routes.IsNull() && !routes.IsUnknown() {
		routeItems := make([]map[string]interface{}, 0, len(routes.Elements()))
		for _, element := range routes.Elements() {
			route, ok := element.(types.Object)
			if !ok || route.IsNull() || route.IsUnknown() {
				continue
			}
			routeAttrs := route.Attributes()
			routeItem := map[string]interface{}{}
			if channels, ok := routeAttrs["channels"].(types.List); ok && !channels.IsNull() && !channels.IsUnknown() {
				channelItems := make([]map[string]interface{}, 0, len(channels.Elements()))
				for _, channelElement := range channels.Elements() {
					channel, ok := channelElement.(types.Object)
					if !ok || channel.IsNull() || channel.IsUnknown() {
						continue
					}
					channelAttrs := channel.Attributes()
					channelItem := map[string]interface{}{
						"channelType": cmsSubscriptionAttrString(channelAttrs["channel_type"]),
					}
					if v := cmsSubscriptionAttrStringList(channelAttrs["receivers"]); len(v) > 0 {
						channelItem["receivers"] = v
					}
					if v := cmsSubscriptionAttrStringList(channelAttrs["enabled_sub_channels"]); len(v) > 0 {
						channelItem["enabledSubChannels"] = v
					}
					channelItems = append(channelItems, channelItem)
				}
				if len(channelItems) > 0 {
					routeItem["channels"] = channelItems
				}
			}
			if effectTimeRanges, ok := routeAttrs["effect_time_range"].(types.List); ok && !effectTimeRanges.IsNull() && !effectTimeRanges.IsUnknown() {
				rangeItems := make([]map[string]interface{}, 0, len(effectTimeRanges.Elements()))
				for _, rangeElement := range effectTimeRanges.Elements() {
					effectTimeRange, ok := rangeElement.(types.Object)
					if !ok || effectTimeRange.IsNull() || effectTimeRange.IsUnknown() {
						continue
					}
					rangeAttrs := effectTimeRange.Attributes()
					rangeItem := map[string]interface{}{}
					if v := cmsSubscriptionAttrString(rangeAttrs["time_zone"]); v != "" {
						rangeItem["timeZone"] = v
					}
					if v := cmsSubscriptionAttrInt64List(rangeAttrs["day_in_week"]); len(v) > 0 {
						rangeItem["dayInWeek"] = v
					}
					if v, ok := cmsSubscriptionAttrInt64(rangeAttrs["start_time_in_minute"]); ok {
						rangeItem["startTimeInMinute"] = v
					}
					if v, ok := cmsSubscriptionAttrInt64(rangeAttrs["end_time_in_minute"]); ok {
						rangeItem["endTimeInMinute"] = v
					}
					rangeItems = append(rangeItems, rangeItem)
				}
				if len(rangeItems) > 0 {
					routeItem["effectTimeRange"] = rangeItems[0]
				}
			}
			routeItems = append(routeItems, routeItem)
		}
		if len(routeItems) > 0 {
			result["routes"] = routeItems
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func cmsSubscriptionState(object map[string]interface{}, previous cmsSubscriptionResourceModel) cmsSubscriptionResourceModel {
	state := cmsSubscriptionResourceModel{
		Id:                     types.StringValue(cmsSubscriptionJSONString(object["subscriptionId"])),
		SubscriptionName:       cmsSubscriptionJSONStringValue(object["subscriptionName"]),
		Description:            cmsSubscriptionJSONStringValue(object["description"]),
		NotifyStrategyId:       cmsSubscriptionJSONStringValue(object["notifyStrategyId"]),
		Workspace:              cmsSubscriptionJSONStringValue(object["workspace"]),
		Enable:                 cmsSubscriptionJSONBoolValue(object["enable"]),
		SubscribeLegacyEvent:   cmsSubscriptionJSONBoolValue(object["subscribeLegacyEvent"]),
		FilterSetting:          flattenCmsSubscriptionFilterSetting(object["filterSetting"]),
		PushingSetting:         flattenCmsSubscriptionPushingSetting(object["pushingSetting"]),
		WorkspaceFilterSetting: flattenCmsSubscriptionWorkspaceFilterSetting(object["workspaceFilterSetting"]),
		AgentConfig:            flattenCmsSubscriptionAgentConfig(object["agentConfig"]),
		UserId:                 cmsSubscriptionJSONStringValue(object["userId"]),
		SubscriptionType:       cmsSubscriptionJSONStringValue(object["subscriptionType"]),
		RegionId:               cmsSubscriptionJSONStringValue(object["regionId"]),
		CreateTime:             cmsSubscriptionJSONStringValue(object["createTime"]),
		UpdateTime:             cmsSubscriptionJSONStringValue(object["updateTime"]),
	}
	if state.Id.ValueString() == "" && !previous.Id.IsNull() {
		state.Id = previous.Id
	}
	if state.Workspace.IsNull() && !previous.Workspace.IsNull() && !previous.Workspace.IsUnknown() {
		// The API echoes the workspace only for workspace-scoped subscriptions;
		// keep the configured value so workspace-scoped resources stay addressable.
		state.Workspace = previous.Workspace
	}
	return state
}

func flattenCmsSubscriptionFilterSetting(raw interface{}) types.List {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		return types.ListNull(cmsSubscriptionFilterSettingType)
	}
	conditions := make([]attr.Value, 0)
	if items, ok := setting["conditions"].([]interface{}); ok {
		for _, item := range items {
			condition, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			conditions = append(conditions, types.ObjectValueMust(cmsSubscriptionConditionType.AttrTypes, map[string]attr.Value{
				"field": cmsSubscriptionJSONStringValue(condition["field"]),
				"op":    cmsSubscriptionJSONStringValue(condition["op"]),
				"value": cmsSubscriptionJSONStringValue(condition["value"]),
			}))
		}
	}
	conditionList := types.ListValueMust(cmsSubscriptionConditionType, conditions)
	if len(conditions) == 0 {
		conditionList = types.ListNull(cmsSubscriptionConditionType)
	}
	object := types.ObjectValueMust(cmsSubscriptionFilterSettingType.AttrTypes, map[string]attr.Value{
		"conditions": conditionList,
		"expression": cmsSubscriptionJSONStringValue(setting["expression"]),
		"relation":   cmsSubscriptionJSONStringValue(setting["relation"]),
	})
	return types.ListValueMust(cmsSubscriptionFilterSettingType, []attr.Value{object})
}

func flattenCmsSubscriptionPushingSetting(raw interface{}) types.List {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		return types.ListNull(cmsSubscriptionPushingSettingType)
	}
	object := types.ObjectValueMust(cmsSubscriptionPushingSettingType.AttrTypes, map[string]attr.Value{
		"alert_action_ids":   cmsSubscriptionJSONStringList(setting["alertActionIds"]),
		"restore_action_ids": cmsSubscriptionJSONStringList(setting["restoreActionIds"]),
		"template_uuid":      cmsSubscriptionJSONStringValue(setting["templateUuid"]),
		"response_plan_id":   cmsSubscriptionJSONStringValue(setting["responsePlanId"]),
	})
	return types.ListValueMust(cmsSubscriptionPushingSettingType, []attr.Value{object})
}

func flattenCmsSubscriptionWorkspaceFilterSetting(raw interface{}) types.List {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		return types.ListNull(cmsSubscriptionWorkspaceFilterSettingType)
	}
	object := types.ObjectValueMust(cmsSubscriptionWorkspaceFilterSettingType.AttrTypes, map[string]attr.Value{
		"workspace_uuids": cmsSubscriptionJSONStringList(setting["workspaceUuids"]),
	})
	return types.ListValueMust(cmsSubscriptionWorkspaceFilterSettingType, []attr.Value{object})
}

func flattenCmsSubscriptionAgentConfig(raw interface{}) types.List {
	config, ok := raw.(map[string]interface{})
	if !ok || len(config) == 0 {
		return types.ListNull(cmsSubscriptionAgentConfigType)
	}
	routes := make([]attr.Value, 0)
	if items, ok := config["routes"].([]interface{}); ok {
		for _, item := range items {
			route, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			routes = append(routes, flattenCmsSubscriptionAgentRoute(route))
		}
	}
	routeList := types.ListValueMust(cmsSubscriptionAgentRouteType, routes)
	if len(routes) == 0 {
		routeList = types.ListNull(cmsSubscriptionAgentRouteType)
	}
	object := types.ObjectValueMust(cmsSubscriptionAgentConfigType.AttrTypes, map[string]attr.Value{
		"agent_uuid": cmsSubscriptionJSONStringValue(config["agentUuid"]),
		"routes":     routeList,
	})
	return types.ListValueMust(cmsSubscriptionAgentConfigType, []attr.Value{object})
}

func flattenCmsSubscriptionAgentRoute(route map[string]interface{}) attr.Value {
	channels := make([]attr.Value, 0)
	if items, ok := route["channels"].([]interface{}); ok {
		for _, item := range items {
			channel, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			channels = append(channels, types.ObjectValueMust(cmsSubscriptionAgentChannelType.AttrTypes, map[string]attr.Value{
				"channel_type":         cmsSubscriptionJSONStringValue(channel["channelType"]),
				"receivers":            cmsSubscriptionJSONStringList(channel["receivers"]),
				"enabled_sub_channels": cmsSubscriptionJSONStringList(channel["enabledSubChannels"]),
			}))
		}
	}
	channelList := types.ListValueMust(cmsSubscriptionAgentChannelType, channels)
	if len(channels) == 0 {
		channelList = types.ListNull(cmsSubscriptionAgentChannelType)
	}

	effectTimeRanges := make([]attr.Value, 0)
	if item, ok := route["effectTimeRange"].(map[string]interface{}); ok {
		effectTimeRanges = append(effectTimeRanges, types.ObjectValueMust(cmsSubscriptionAgentEffectTimeRangeType.AttrTypes, map[string]attr.Value{
			"time_zone":            cmsSubscriptionJSONStringValue(item["timeZone"]),
			"day_in_week":          cmsSubscriptionJSONInt64List(item["dayInWeek"]),
			"start_time_in_minute": cmsSubscriptionJSONInt64Value(item["startTimeInMinute"]),
			"end_time_in_minute":   cmsSubscriptionJSONInt64Value(item["endTimeInMinute"]),
		}))
	}
	effectTimeRangeList := types.ListValueMust(cmsSubscriptionAgentEffectTimeRangeType, effectTimeRanges)
	if len(effectTimeRanges) == 0 {
		effectTimeRangeList = types.ListNull(cmsSubscriptionAgentEffectTimeRangeType)
	}

	return types.ObjectValueMust(cmsSubscriptionAgentRouteType.AttrTypes, map[string]attr.Value{
		"channels":          channelList,
		"effect_time_range": effectTimeRangeList,
	})
}

func cmsSubscriptionFirstObject(l types.List) (types.Object, bool) {
	if l.IsNull() || l.IsUnknown() || len(l.Elements()) == 0 {
		return types.Object{}, false
	}
	object, ok := l.Elements()[0].(types.Object)
	if !ok || object.IsNull() || object.IsUnknown() {
		return types.Object{}, false
	}
	return object, true
}

func cmsSubscriptionAttrString(v attr.Value) string {
	if s, ok := v.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
		return s.ValueString()
	}
	return ""
}

func cmsSubscriptionAttrStringList(v attr.Value) []string {
	list, ok := v.(types.List)
	if !ok || list.IsNull() || list.IsUnknown() {
		return nil
	}
	values := make([]string, 0, len(list.Elements()))
	for _, element := range list.Elements() {
		if s, ok := element.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
			values = append(values, s.ValueString())
		}
	}
	return values
}

func cmsSubscriptionAttrInt64(v attr.Value) (int64, bool) {
	if i, ok := v.(types.Int64); ok && !i.IsNull() && !i.IsUnknown() {
		return i.ValueInt64(), true
	}
	return 0, false
}

func cmsSubscriptionAttrInt64List(v attr.Value) []int64 {
	list, ok := v.(types.List)
	if !ok || list.IsNull() || list.IsUnknown() {
		return nil
	}
	values := make([]int64, 0, len(list.Elements()))
	for _, element := range list.Elements() {
		if i, ok := element.(types.Int64); ok && !i.IsNull() && !i.IsUnknown() {
			values = append(values, i.ValueInt64())
		}
	}
	return values
}

func cmsSubscriptionJSONString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func cmsSubscriptionJSONStringValue(v interface{}) types.String {
	if s := cmsSubscriptionJSONString(v); s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

func cmsSubscriptionJSONBoolValue(v interface{}) types.Bool {
	if b, ok := v.(bool); ok {
		return types.BoolValue(b)
	}
	return types.BoolNull()
}

func cmsSubscriptionJSONStringList(v interface{}) types.List {
	elementType := types.ListType{ElemType: types.StringType}
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return types.ListNull(elementType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			values = append(values, types.StringValue(s))
		}
	}
	if len(values) == 0 {
		return types.ListNull(elementType)
	}
	return types.ListValueMust(elementType, values)
}

func cmsSubscriptionJSONInt64List(v interface{}) types.List {
	elementType := types.ListType{ElemType: types.Int64Type}
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return types.ListNull(elementType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		value := cmsSubscriptionJSONInt64Value(item)
		if !value.IsNull() {
			values = append(values, value)
		}
	}
	if len(values) == 0 {
		return types.ListNull(elementType)
	}
	return types.ListValueMust(elementType, values)
}

func cmsSubscriptionJSONInt64Value(v interface{}) types.Int64 {
	switch n := v.(type) {
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return types.Int64Value(i)
		}
	case float64:
		return types.Int64Value(int64(n))
	case int:
		return types.Int64Value(int64(n))
	case int64:
		return types.Int64Value(n)
	}
	return types.Int64Null()
}
