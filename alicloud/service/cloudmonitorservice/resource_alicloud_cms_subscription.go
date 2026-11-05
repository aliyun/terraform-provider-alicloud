package cloudmonitorservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	Id                     types.String                          `tfsdk:"id"`
	SubscriptionName       types.String                          `tfsdk:"subscription_name"`
	Description            types.String                          `tfsdk:"description"`
	NotifyStrategyId       types.String                          `tfsdk:"notify_strategy_id"`
	Workspace              types.String                          `tfsdk:"workspace"`
	Enable                 types.Bool                            `tfsdk:"enable"`
	SubscribeLegacyEvent   types.Bool                            `tfsdk:"subscribe_legacy_event"`
	FilterSetting          []cmsSubscriptionFilterSettingModel   `tfsdk:"filter_setting"`
	PushingSetting         []cmsSubscriptionPushingSettingModel  `tfsdk:"pushing_setting"`
	WorkspaceFilterSetting []cmsSubscriptionWorkspaceFilterModel `tfsdk:"workspace_filter_setting"`
	AgentConfig            []cmsSubscriptionAgentConfigModel     `tfsdk:"agent_config"`
	UserId                 types.String                          `tfsdk:"user_id"`
	SubscriptionType       types.String                          `tfsdk:"subscription_type"`
	RegionId               types.String                          `tfsdk:"region_id"`
	CreateTime             types.String                          `tfsdk:"create_time"`
	UpdateTime             types.String                          `tfsdk:"update_time"`
}

type cmsSubscriptionFilterSettingModel struct {
	Conditions []cmsSubscriptionConditionModel `tfsdk:"conditions"`
	Expression types.String                    `tfsdk:"expression"`
	Relation   types.String                    `tfsdk:"relation"`
}

type cmsSubscriptionConditionModel struct {
	Field types.String `tfsdk:"field"`
	Op    types.String `tfsdk:"op"`
	Value types.String `tfsdk:"value"`
}

type cmsSubscriptionPushingSettingModel struct {
	AlertActionIds   types.List   `tfsdk:"alert_action_ids"`
	RestoreActionIds types.List   `tfsdk:"restore_action_ids"`
	TemplateUuid     types.String `tfsdk:"template_uuid"`
	ResponsePlanId   types.String `tfsdk:"response_plan_id"`
}

type cmsSubscriptionWorkspaceFilterModel struct {
	WorkspaceUuids types.List `tfsdk:"workspace_uuids"`
}

type cmsSubscriptionAgentConfigModel struct {
	AgentUuid types.String                     `tfsdk:"agent_uuid"`
	Routes    []cmsSubscriptionAgentRouteModel `tfsdk:"routes"`
}

type cmsSubscriptionAgentRouteModel struct {
	Channels        []cmsSubscriptionAgentChannelModel         `tfsdk:"channels"`
	EffectTimeRange []cmsSubscriptionAgentEffectTimeRangeModel `tfsdk:"effect_time_range"`
}

type cmsSubscriptionAgentChannelModel struct {
	ChannelType        types.String `tfsdk:"channel_type"`
	Receivers          types.List   `tfsdk:"receivers"`
	EnabledSubChannels types.List   `tfsdk:"enabled_sub_channels"`
}

type cmsSubscriptionAgentEffectTimeRangeModel struct {
	TimeZone          types.String `tfsdk:"time_zone"`
	DayInWeek         types.List   `tfsdk:"day_in_week"`
	StartTimeInMinute types.Int64  `tfsdk:"start_time_in_minute"`
	EndTimeInMinute   types.Int64  `tfsdk:"end_time_in_minute"`
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
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
		Blocks: map[string]schema.Block{
			"filter_setting": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"expression": schema.StringAttribute{
							Optional: true,
						},
						"relation": schema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]schema.Block{
						"conditions": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"field": schema.StringAttribute{
										Required: true,
									},
									"op": schema.StringAttribute{
										Required: true,
									},
									"value": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"pushing_setting": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alert_action_ids": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						"restore_action_ids": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						"template_uuid": schema.StringAttribute{
							Optional: true,
						},
						"response_plan_id": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"workspace_filter_setting": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"workspace_uuids": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"agent_config": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"agent_uuid": schema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]schema.Block{
						"routes": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"channels": schema.ListNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"channel_type": schema.StringAttribute{
													Optional: true,
												},
												"receivers": schema.ListAttribute{
													Optional:    true,
													ElementType: types.StringType,
												},
												"enabled_sub_channels": schema.ListAttribute{
													Optional:    true,
													ElementType: types.StringType,
												},
											},
										},
									},
									"effect_time_range": schema.ListNestedBlock{
										Validators: []validator.List{
											listvalidator.SizeAtMost(1),
										},
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"time_zone": schema.StringAttribute{
													Optional: true,
												},
												"day_in_week": schema.ListAttribute{
													Optional:    true,
													ElementType: types.Int64Type,
												},
												"start_time_in_minute": schema.Int64Attribute{
													Optional: true,
												},
												"end_time_in_minute": schema.Int64Attribute{
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
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

	object, err := r.waitCmsSubscriptionReadable(id, plan.Workspace)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
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
		if err := r.waitCmsSubscriptionEnable(id, plan.Workspace, false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: waiting for DisableSubscription", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		if object, err = NewCmsService(r.Meta()).DescribeCmsSubscription(id, plan.Workspace.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	previous := cmsSubscriptionResourceModel{
		Id:               types.StringValue(id),
		SubscriptionName: plan.SubscriptionName,
		Workspace:        plan.Workspace,
		Enable:           plan.Enable,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, cmsSubscriptionState(object, previous))...)
}

func (r *cmsSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(state.Id.ValueString(), state.Workspace.ValueString())
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
				fmt.Sprintf("Updating %s: toggling the subscription enabled state", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id+"/enable", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		if err := r.waitCmsSubscriptionEnable(id, state.Workspace, plan.Enable.ValueBool()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: waiting for the subscription enabled state", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
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

	object, err := r.waitCmsSubscriptionUpdate(id, state.Workspace, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	updated := cmsSubscriptionState(object, state)
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

// ImportState accepts <workspace>:<subscription_id> because every subscription
// request, including the read used to refresh imported state, is scoped by its
// workspace.
func (r *cmsSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected an import ID in the form <workspace>:<subscription_id>, got %q.", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace"), parts[0])...)
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

// waitCmsSubscriptionEnable polls GetSubscription until the enabled state
// reported by the service matches the requested one. The enable and disable
// operations return before the change is visible on the read path.
func (r *cmsSubscriptionResource) waitCmsSubscriptionEnable(id string, workspace types.String, enable bool) error {
	deadline := time.Now().Add(cmsSubscriptionTimeout)
	for {
		object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(id, workspace.ValueString())
		if err != nil {
			return err
		}
		if current, ok := object["enable"].(bool); ok && current == enable {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("subscription %s did not reach enable=%t within %s", id, enable, cmsSubscriptionTimeout)
		}
		time.Sleep(3 * time.Second)
	}
}

// waitCmsSubscriptionReadable polls GetSubscription until the freshly created
// subscription is visible on the read path. Creation is asynchronous, so an
// immediate read can still report the subscription as not found.
func (r *cmsSubscriptionResource) waitCmsSubscriptionReadable(id string, workspace types.String) (map[string]interface{}, error) {
	deadline := time.Now().Add(cmsSubscriptionTimeout)
	for {
		object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(id, workspace.ValueString())
		if err == nil {
			return object, nil
		}
		if !errs.NotFoundError(err) || time.Now().After(deadline) {
			return nil, err
		}
		time.Sleep(3 * time.Second)
	}
}

// waitCmsSubscriptionUpdate polls GetSubscription until the configuration sent
// to UpdateSubscription is visible on the read path. The Cms event subscription
// reads are eventually consistent, so an immediate read can return the values
// from before the update.
func (r *cmsSubscriptionResource) waitCmsSubscriptionUpdate(id string, workspace types.String, request map[string]interface{}) (map[string]interface{}, error) {
	deadline := time.Now().Add(cmsSubscriptionTimeout)
	for {
		object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(id, workspace.ValueString())
		if err != nil {
			if errs.NotFoundError(err) && time.Now().Before(deadline) {
				time.Sleep(3 * time.Second)
				continue
			}
			return nil, err
		}
		if cmsSubscriptionUpdateApplied(object, request) {
			return object, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("subscription %s configuration did not converge within %s", id, cmsSubscriptionTimeout)
		}
		time.Sleep(3 * time.Second)
	}
}

func cmsSubscriptionUpdateApplied(object, request map[string]interface{}) bool {
	for _, key := range []string{"subscriptionName", "description", "notifyStrategyId"} {
		if want, ok := request[key].(string); ok {
			got, _ := object[key].(string)
			if got != want {
				return false
			}
		}
	}
	if want, ok := request["subscribeLegacyEvent"].(bool); ok {
		if got, ok := object["subscribeLegacyEvent"].(bool); !ok || got != want {
			return false
		}
	}
	wantFilter, ok := request["filterSetting"].(map[string]interface{})
	if !ok {
		return true
	}
	gotFilter, _ := object["filterSetting"].(map[string]interface{})
	for _, key := range []string{"expression", "relation"} {
		if want, ok := wantFilter[key].(string); ok {
			if got, _ := gotFilter[key].(string); got != want {
				return false
			}
		}
	}
	if wantConditions, ok := wantFilter["conditions"].([]map[string]interface{}); ok {
		gotConditions, _ := gotFilter["conditions"].([]interface{})
		if len(gotConditions) != len(wantConditions) {
			return false
		}
	}
	return true
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

func expandCmsSubscriptionFilterSetting(blocks []cmsSubscriptionFilterSettingModel) map[string]interface{} {
	if len(blocks) == 0 {
		return nil
	}
	block := blocks[0]
	result := map[string]interface{}{}
	if v := cmsSubscriptionValueString(block.Expression); v != "" {
		result["expression"] = v
	}
	if v := cmsSubscriptionValueString(block.Relation); v != "" {
		result["relation"] = v
	}
	if len(block.Conditions) > 0 {
		conditions := make([]map[string]interface{}, 0, len(block.Conditions))
		for _, condition := range block.Conditions {
			conditions = append(conditions, map[string]interface{}{
				"field": cmsSubscriptionValueString(condition.Field),
				"op":    cmsSubscriptionValueString(condition.Op),
				"value": cmsSubscriptionValueString(condition.Value),
			})
		}
		result["conditions"] = conditions
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandCmsSubscriptionPushingSetting(blocks []cmsSubscriptionPushingSettingModel) map[string]interface{} {
	if len(blocks) == 0 {
		return nil
	}
	block := blocks[0]
	result := map[string]interface{}{}
	if v := cmsSubscriptionValueStringList(block.AlertActionIds); len(v) > 0 {
		result["alertActionIds"] = v
	}
	if v := cmsSubscriptionValueStringList(block.RestoreActionIds); len(v) > 0 {
		result["restoreActionIds"] = v
	}
	if v := cmsSubscriptionValueString(block.TemplateUuid); v != "" {
		result["templateUuid"] = v
	}
	if v := cmsSubscriptionValueString(block.ResponsePlanId); v != "" {
		result["responsePlanId"] = v
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandCmsSubscriptionWorkspaceFilterSetting(blocks []cmsSubscriptionWorkspaceFilterModel) map[string]interface{} {
	if len(blocks) == 0 {
		return nil
	}
	result := map[string]interface{}{}
	if v := cmsSubscriptionValueStringList(blocks[0].WorkspaceUuids); len(v) > 0 {
		result["workspaceUuids"] = v
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandCmsSubscriptionAgentConfig(blocks []cmsSubscriptionAgentConfigModel) map[string]interface{} {
	if len(blocks) == 0 {
		return nil
	}
	block := blocks[0]
	result := map[string]interface{}{}
	if v := cmsSubscriptionValueString(block.AgentUuid); v != "" {
		result["agentUuid"] = v
	}
	if len(block.Routes) > 0 {
		routes := make([]map[string]interface{}, 0, len(block.Routes))
		for _, route := range block.Routes {
			routeItem := map[string]interface{}{}
			if len(route.Channels) > 0 {
				channels := make([]map[string]interface{}, 0, len(route.Channels))
				for _, channel := range route.Channels {
					channelItem := map[string]interface{}{
						"channelType": cmsSubscriptionValueString(channel.ChannelType),
					}
					if v := cmsSubscriptionValueStringList(channel.Receivers); len(v) > 0 {
						channelItem["receivers"] = v
					}
					if v := cmsSubscriptionValueStringList(channel.EnabledSubChannels); len(v) > 0 {
						channelItem["enabledSubChannels"] = v
					}
					channels = append(channels, channelItem)
				}
				routeItem["channels"] = channels
			}
			if len(route.EffectTimeRange) > 0 {
				effectTimeRange := route.EffectTimeRange[0]
				rangeItem := map[string]interface{}{}
				if v := cmsSubscriptionValueString(effectTimeRange.TimeZone); v != "" {
					rangeItem["timeZone"] = v
				}
				if v := cmsSubscriptionValueInt64List(effectTimeRange.DayInWeek); len(v) > 0 {
					rangeItem["dayInWeek"] = v
				}
				if v, ok := cmsSubscriptionValueInt64(effectTimeRange.StartTimeInMinute); ok {
					rangeItem["startTimeInMinute"] = v
				}
				if v, ok := cmsSubscriptionValueInt64(effectTimeRange.EndTimeInMinute); ok {
					rangeItem["endTimeInMinute"] = v
				}
				routeItem["effectTimeRange"] = rangeItem
			}
			routes = append(routes, routeItem)
		}
		result["routes"] = routes
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

func flattenCmsSubscriptionFilterSetting(raw interface{}) []cmsSubscriptionFilterSettingModel {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		return nil
	}
	block := cmsSubscriptionFilterSettingModel{
		Expression: cmsSubscriptionJSONStringValue(setting["expression"]),
		Relation:   cmsSubscriptionJSONStringValue(setting["relation"]),
	}
	if items, ok := setting["conditions"].([]interface{}); ok {
		for _, item := range items {
			condition, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			block.Conditions = append(block.Conditions, cmsSubscriptionConditionModel{
				Field: cmsSubscriptionJSONStringValue(condition["field"]),
				Op:    cmsSubscriptionJSONStringValue(condition["op"]),
				Value: cmsSubscriptionJSONStringValue(condition["value"]),
			})
		}
	}
	return []cmsSubscriptionFilterSettingModel{block}
}

func flattenCmsSubscriptionPushingSetting(raw interface{}) []cmsSubscriptionPushingSettingModel {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		return nil
	}
	return []cmsSubscriptionPushingSettingModel{{
		AlertActionIds:   cmsSubscriptionJSONStringList(setting["alertActionIds"]),
		RestoreActionIds: cmsSubscriptionJSONStringList(setting["restoreActionIds"]),
		TemplateUuid:     cmsSubscriptionJSONStringValue(setting["templateUuid"]),
		ResponsePlanId:   cmsSubscriptionJSONStringValue(setting["responsePlanId"]),
	}}
}

func flattenCmsSubscriptionWorkspaceFilterSetting(raw interface{}) []cmsSubscriptionWorkspaceFilterModel {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		return nil
	}
	return []cmsSubscriptionWorkspaceFilterModel{{
		WorkspaceUuids: cmsSubscriptionJSONStringList(setting["workspaceUuids"]),
	}}
}

func flattenCmsSubscriptionAgentConfig(raw interface{}) []cmsSubscriptionAgentConfigModel {
	config, ok := raw.(map[string]interface{})
	if !ok || len(config) == 0 {
		return nil
	}
	block := cmsSubscriptionAgentConfigModel{
		AgentUuid: cmsSubscriptionJSONStringValue(config["agentUuid"]),
	}
	if items, ok := config["routes"].([]interface{}); ok {
		for _, item := range items {
			route, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			block.Routes = append(block.Routes, flattenCmsSubscriptionAgentRoute(route))
		}
	}
	return []cmsSubscriptionAgentConfigModel{block}
}

func flattenCmsSubscriptionAgentRoute(route map[string]interface{}) cmsSubscriptionAgentRouteModel {
	result := cmsSubscriptionAgentRouteModel{}
	if items, ok := route["channels"].([]interface{}); ok {
		for _, item := range items {
			channel, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			result.Channels = append(result.Channels, cmsSubscriptionAgentChannelModel{
				ChannelType:        cmsSubscriptionJSONStringValue(channel["channelType"]),
				Receivers:          cmsSubscriptionJSONStringList(channel["receivers"]),
				EnabledSubChannels: cmsSubscriptionJSONStringList(channel["enabledSubChannels"]),
			})
		}
	}
	if item, ok := route["effectTimeRange"].(map[string]interface{}); ok {
		result.EffectTimeRange = append(result.EffectTimeRange, cmsSubscriptionAgentEffectTimeRangeModel{
			TimeZone:          cmsSubscriptionJSONStringValue(item["timeZone"]),
			DayInWeek:         cmsSubscriptionJSONInt64List(item["dayInWeek"]),
			StartTimeInMinute: cmsSubscriptionJSONInt64Value(item["startTimeInMinute"]),
			EndTimeInMinute:   cmsSubscriptionJSONInt64Value(item["endTimeInMinute"]),
		})
	}
	return result
}

func cmsSubscriptionValueString(v types.String) string {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueString()
	}
	return ""
}

func cmsSubscriptionValueStringList(v types.List) []string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	values := make([]string, 0, len(v.Elements()))
	for _, element := range v.Elements() {
		if s, ok := element.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
			values = append(values, s.ValueString())
		}
	}
	return values
}

func cmsSubscriptionValueInt64(v types.Int64) (int64, bool) {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueInt64(), true
	}
	return 0, false
}

func cmsSubscriptionValueInt64List(v types.List) []int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	values := make([]int64, 0, len(v.Elements()))
	for _, element := range v.Elements() {
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
