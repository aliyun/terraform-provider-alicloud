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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	_ resource.Resource                     = &cmsSubscriptionResource{}
	_ resource.ResourceWithImportState      = &cmsSubscriptionResource{}
	_ resource.ResourceWithConfigValidators = &cmsSubscriptionResource{}
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

func (r *cmsSubscriptionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{cmsSubscriptionChannelValidator{}}
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
							Validators: []validator.String{
								stringvalidator.OneOf("AND", "OR"),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"conditions": schema.ListNestedBlock{
							Validators: []validator.List{
								listvalidator.SizeAtMost(10),
							},
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"field": schema.StringAttribute{
										Required: true,
									},
									"op": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("IN", "CONTAIN", "EQ"),
										},
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
													Validators: []validator.List{
														listvalidator.ValueInt64sAre(int64validator.Between(0, 6)),
													},
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
	err = retry.RetryContext(ctx, cmsSubscriptionTimeout, func() *retry.RetryError {
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

	// Persist the identity as soon as it is known: every later step can fail,
	// and without it the created subscription could neither be refreshed nor destroyed.
	resp.Diagnostics.Append(resp.State.Set(ctx, cmsSubscriptionIdentityState(plan, id))...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.waitCmsSubscriptionReadable(ctx, id, plan.Workspace, request, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if !plan.Enable.IsNull() && !plan.Enable.IsUnknown() && !plan.Enable.ValueBool() {
		if err := r.setCmsSubscriptionEnable(ctx, id, plan.Workspace, false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling DisableSubscription", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id+"/disable", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		if err := r.waitCmsSubscriptionEnable(ctx, id, plan.Workspace, false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: waiting for DisableSubscription", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		if object, err = NewCmsService(r.Meta()).DescribeCmsSubscription(ctx, id, plan.Workspace.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, cmsSubscriptionState(object, plan, r.Meta().RegionId))...)
}

func (r *cmsSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(ctx, state.Id.ValueString(), state.Workspace.ValueString())
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
	resp.Diagnostics.Append(resp.State.Set(ctx, cmsSubscriptionState(object, state, r.Meta().RegionId))...)
}

func (r *cmsSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if !plan.Enable.Equal(state.Enable) && !plan.Enable.IsNull() && !plan.Enable.IsUnknown() {
		enable := plan.Enable.ValueBool()
		suffix := "enable"
		if !enable {
			suffix = "disable"
		}
		if err := r.setCmsSubscriptionEnable(ctx, id, state.Workspace, enable); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling /%s", cmsSubscriptionTypeName, suffix),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, "/subscriptions/"+id+"/"+suffix, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		if err := r.waitCmsSubscriptionEnable(ctx, id, state.Workspace, enable); err != nil {
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
	err = retry.RetryContext(ctx, cmsSubscriptionTimeout, func() *retry.RetryError {
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

	removed := cmsSubscriptionRemovedPaths(plan, state)
	object, err := r.waitCmsSubscriptionUpdate(ctx, id, state.Workspace, request, removed)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", cmsSubscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, cmsSubscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, cmsSubscriptionState(object, plan, r.Meta().RegionId))...)
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
	err = retry.RetryContext(ctx, cmsSubscriptionTimeout, func() *retry.RetryError {
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
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"ResourceNotFound"}) {
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

func (r *cmsSubscriptionResource) setCmsSubscriptionEnable(ctx context.Context, id string, workspace types.String, enable bool) error {
	suffix := "enable"
	if !enable {
		suffix = "disable"
	}
	action := fmt.Sprintf("/subscriptions/%s/%s", id, suffix)

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, cmsSubscriptionTimeout, func() *retry.RetryError {
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
func (r *cmsSubscriptionResource) waitCmsSubscriptionEnable(ctx context.Context, id string, workspace types.String, enable bool) error {
	deadline := time.Now().Add(cmsSubscriptionTimeout)
	for {
		object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(ctx, id, workspace.ValueString())
		if err != nil {
			return err
		}
		if current, ok := object["enable"].(bool); ok && current == enable {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("subscription %s did not reach enable=%t within %s", id, enable, cmsSubscriptionTimeout)
		}
		if err := cmsSubscriptionSleep(ctx, 3*time.Second); err != nil {
			return err
		}
	}
}

// waitCmsSubscriptionReadable polls GetSubscription until the freshly created
// subscription is visible on the read path and reflects the requested
// configuration. Creation is asynchronous, so an immediate read can still
// report the subscription as not found or return a partially applied object.
func (r *cmsSubscriptionResource) waitCmsSubscriptionReadable(ctx context.Context, id string, workspace types.String, request map[string]interface{}, removed [][]string) (map[string]interface{}, error) {
	deadline := time.Now().Add(cmsSubscriptionTimeout)
	for {
		object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(ctx, id, workspace.ValueString())
		if err == nil && cmsSubscriptionRequestVisible(object, request) && cmsSubscriptionPathsEmpty(object, removed) {
			return object, nil
		}
		if err != nil && !errs.NotFoundError(err) {
			return nil, err
		}
		if time.Now().After(deadline) {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("subscription %s configuration did not become visible within %s", id, cmsSubscriptionTimeout)
		}
		if err := cmsSubscriptionSleep(ctx, 3*time.Second); err != nil {
			return nil, err
		}
	}
}

// waitCmsSubscriptionUpdate polls GetSubscription until the configuration sent
// to UpdateSubscription is visible on the read path and every removed key is
// gone. The Cms event subscription reads are eventually consistent, so an
// immediate read can return the values from before the update.
func (r *cmsSubscriptionResource) waitCmsSubscriptionUpdate(ctx context.Context, id string, workspace types.String, request map[string]interface{}, removed [][]string) (map[string]interface{}, error) {
	deadline := time.Now().Add(cmsSubscriptionTimeout)
	for {
		object, err := NewCmsService(r.Meta()).DescribeCmsSubscription(ctx, id, workspace.ValueString())
		if err == nil && cmsSubscriptionRequestVisible(object, request) && cmsSubscriptionPathsEmpty(object, removed) {
			return object, nil
		}
		if err != nil && !errs.NotFoundError(err) {
			return nil, err
		}
		if time.Now().After(deadline) {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("subscription %s configuration did not converge within %s", id, cmsSubscriptionTimeout)
		}
		if err := cmsSubscriptionSleep(ctx, 3*time.Second); err != nil {
			return nil, err
		}
	}
}

func cmsSubscriptionSleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// cmsSubscriptionRequestVisible reports whether every value sent to the API is
// already reflected by the read response. It is a subset comparison so fields
// the service adds on its own (ids, timestamps, defaults) do not block
// convergence, while every managed member, list element and nested block does.
func cmsSubscriptionRequestVisible(object, request map[string]interface{}) bool {
	for key, want := range request {
		if !cmsSubscriptionValueVisible(object[key], want) {
			return false
		}
	}
	return true
}

func cmsSubscriptionValueVisible(actual, want interface{}) bool {
	switch wantValue := want.(type) {
	case map[string]interface{}:
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			return false
		}
		return cmsSubscriptionRequestVisible(actualMap, wantValue)
	case []map[string]interface{}:
		actualItems := cmsSubscriptionJSONItems(actual)
		if len(actualItems) != len(wantValue) {
			return false
		}
		used := make([]bool, len(actualItems))
		for _, wantItem := range wantValue {
			matched := false
			for i, actualItem := range actualItems {
				if used[i] {
					continue
				}
				if cmsSubscriptionValueVisible(actualItem, wantItem) {
					used[i] = true
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
		return true
	case []string:
		actualItems := cmsSubscriptionJSONItems(actual)
		if len(actualItems) != len(wantValue) {
			return false
		}
		used := make([]bool, len(actualItems))
		for _, wantItem := range wantValue {
			matched := false
			for i, actualItem := range actualItems {
				if used[i] {
					continue
				}
				if got, ok := actualItem.(string); ok && got == wantItem {
					used[i] = true
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
		return true
	case []int64:
		actualItems := cmsSubscriptionJSONItems(actual)
		if len(actualItems) != len(wantValue) {
			return false
		}
		used := make([]bool, len(actualItems))
		for _, wantItem := range wantValue {
			matched := false
			for i, actualItem := range actualItems {
				if used[i] {
					continue
				}
				if got, ok := cmsSubscriptionJSONInt64(actualItem); ok && got == wantItem {
					used[i] = true
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
		return true
	case string:
		if wantValue == "" {
			// The service may drop empty strings from the response.
			return cmsSubscriptionValueEmpty(actual)
		}
		got, ok := actual.(string)
		return ok && got == wantValue
	case bool:
		got, ok := actual.(bool)
		if ok {
			return got == wantValue
		}
		// subscribeLegacyEvent documents false and null as equivalent.
		return !wantValue
	case int64:
		got, ok := cmsSubscriptionJSONInt64(actual)
		return ok && got == wantValue
	case nil:
		return cmsSubscriptionValueEmpty(actual)
	}
	return false
}

func cmsSubscriptionJSONItems(actual interface{}) []interface{} {
	switch items := actual.(type) {
	case []interface{}:
		return items
	case []map[string]interface{}:
		result := make([]interface{}, 0, len(items))
		for _, item := range items {
			result = append(result, item)
		}
		return result
	}
	return nil
}

func cmsSubscriptionJSONInt64(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return i, true
		}
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	}
	return 0, false
}

func cmsSubscriptionValueEmpty(v interface{}) bool {
	switch value := v.(type) {
	case nil:
		return true
	case string:
		return value == ""
	case bool:
		return !value
	case []interface{}:
		return len(value) == 0
	case []map[string]interface{}:
		return len(value) == 0
	case []string:
		return len(value) == 0
	case map[string]interface{}:
		return len(value) == 0
	}
	return false
}

// cmsSubscriptionPathsEmpty reports whether every path removed from the
// configuration no longer carries a value on the read path.
func cmsSubscriptionPathsEmpty(object map[string]interface{}, paths [][]string) bool {
	for _, fieldPath := range paths {
		if !cmsSubscriptionPathEmpty(object, fieldPath) {
			return false
		}
	}
	return true
}

func cmsSubscriptionPathEmpty(object map[string]interface{}, fieldPath []string) bool {
	var current interface{} = object
	for i, key := range fieldPath {
		container, ok := current.(map[string]interface{})
		if !ok {
			return true
		}
		current = container[key]
		if i == len(fieldPath)-1 {
			return cmsSubscriptionValueEmpty(current)
		}
	}
	return cmsSubscriptionValueEmpty(current)
}

// cmsSubscriptionRemovedPaths lists the configuration paths that were present
// in state and are absent from the new plan, so the update waiter can require
// them to be cleared rather than accepting a stale response.
func cmsSubscriptionRemovedPaths(plan, state cmsSubscriptionResourceModel) [][]string {
	paths := make([][]string, 0)
	if !state.Description.IsNull() && plan.Description.IsNull() {
		paths = append(paths, []string{"description"})
	}
	if !state.NotifyStrategyId.IsNull() && plan.NotifyStrategyId.IsNull() {
		paths = append(paths, []string{"notifyStrategyId"})
	}
	if !state.SubscribeLegacyEvent.IsNull() && plan.SubscribeLegacyEvent.IsNull() {
		paths = append(paths, []string{"subscribeLegacyEvent"})
	}
	if len(state.FilterSetting) > 0 {
		if len(plan.FilterSetting) == 0 {
			paths = append(paths, []string{"filterSetting"})
		} else {
			before := state.FilterSetting[0]
			after := plan.FilterSetting[0]
			if !before.Expression.IsNull() && after.Expression.IsNull() {
				paths = append(paths, []string{"filterSetting", "expression"})
			}
			if !before.Relation.IsNull() && after.Relation.IsNull() {
				paths = append(paths, []string{"filterSetting", "relation"})
			}
			if len(before.Conditions) > 0 && len(after.Conditions) == 0 {
				paths = append(paths, []string{"filterSetting", "conditions"})
			}
		}
	}
	if len(state.PushingSetting) > 0 {
		if len(plan.PushingSetting) == 0 {
			paths = append(paths, []string{"pushingSetting"})
		} else {
			before := state.PushingSetting[0]
			after := plan.PushingSetting[0]
			if !before.AlertActionIds.IsNull() && after.AlertActionIds.IsNull() {
				paths = append(paths, []string{"pushingSetting", "alertActionIds"})
			}
			if !before.RestoreActionIds.IsNull() && after.RestoreActionIds.IsNull() {
				paths = append(paths, []string{"pushingSetting", "restoreActionIds"})
			}
			if !before.TemplateUuid.IsNull() && after.TemplateUuid.IsNull() {
				paths = append(paths, []string{"pushingSetting", "templateUuid"})
			}
			if !before.ResponsePlanId.IsNull() && after.ResponsePlanId.IsNull() {
				paths = append(paths, []string{"pushingSetting", "responsePlanId"})
			}
		}
	}
	if len(state.WorkspaceFilterSetting) > 0 {
		if len(plan.WorkspaceFilterSetting) == 0 {
			paths = append(paths, []string{"workspaceFilterSetting"})
		} else if !state.WorkspaceFilterSetting[0].WorkspaceUuids.IsNull() && plan.WorkspaceFilterSetting[0].WorkspaceUuids.IsNull() {
			paths = append(paths, []string{"workspaceFilterSetting", "workspaceUuids"})
		}
	}
	if len(state.AgentConfig) > 0 {
		if len(plan.AgentConfig) == 0 {
			paths = append(paths, []string{"agentConfig"})
		} else {
			before := state.AgentConfig[0]
			after := plan.AgentConfig[0]
			if !before.AgentUuid.IsNull() && after.AgentUuid.IsNull() {
				paths = append(paths, []string{"agentConfig", "agentUuid"})
			}
			if len(before.Routes) > 0 && len(after.Routes) == 0 {
				paths = append(paths, []string{"agentConfig", "routes"})
			}
		}
	}
	return paths
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
// Values that are configured but empty ("", [] or false) are sent explicitly so
// the service can clear them instead of keeping the previous value.
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
	if !block.Expression.IsNull() && !block.Expression.IsUnknown() {
		result["expression"] = block.Expression.ValueString()
	}
	if !block.Relation.IsNull() && !block.Relation.IsUnknown() {
		result["relation"] = block.Relation.ValueString()
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
	if !block.AlertActionIds.IsNull() && !block.AlertActionIds.IsUnknown() {
		result["alertActionIds"] = cmsSubscriptionValueStringList(block.AlertActionIds)
	}
	if !block.RestoreActionIds.IsNull() && !block.RestoreActionIds.IsUnknown() {
		result["restoreActionIds"] = cmsSubscriptionValueStringList(block.RestoreActionIds)
	}
	if !block.TemplateUuid.IsNull() && !block.TemplateUuid.IsUnknown() {
		result["templateUuid"] = block.TemplateUuid.ValueString()
	}
	if !block.ResponsePlanId.IsNull() && !block.ResponsePlanId.IsUnknown() {
		result["responsePlanId"] = block.ResponsePlanId.ValueString()
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
	if !blocks[0].WorkspaceUuids.IsNull() && !blocks[0].WorkspaceUuids.IsUnknown() {
		result["workspaceUuids"] = cmsSubscriptionValueStringList(blocks[0].WorkspaceUuids)
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
	if !block.AgentUuid.IsNull() && !block.AgentUuid.IsUnknown() {
		result["agentUuid"] = block.AgentUuid.ValueString()
	}
	if len(block.Routes) > 0 {
		routes := make([]map[string]interface{}, 0, len(block.Routes))
		for _, route := range block.Routes {
			routeItem := map[string]interface{}{}
			if route.Channels != nil {
				channels := make([]map[string]interface{}, 0, len(route.Channels))
				for _, channel := range route.Channels {
					channelItem := map[string]interface{}{}
					if !channel.ChannelType.IsNull() && !channel.ChannelType.IsUnknown() {
						channelItem["channelType"] = channel.ChannelType.ValueString()
					}
					if !channel.Receivers.IsNull() && !channel.Receivers.IsUnknown() {
						channelItem["receivers"] = cmsSubscriptionValueStringList(channel.Receivers)
					}
					if !channel.EnabledSubChannels.IsNull() && !channel.EnabledSubChannels.IsUnknown() {
						channelItem["enabledSubChannels"] = cmsSubscriptionValueStringList(channel.EnabledSubChannels)
					}
					channels = append(channels, channelItem)
				}
				routeItem["channels"] = channels
			}
			if len(route.EffectTimeRange) > 0 {
				effectTimeRange := route.EffectTimeRange[0]
				rangeItem := map[string]interface{}{}
				if !effectTimeRange.TimeZone.IsNull() && !effectTimeRange.TimeZone.IsUnknown() {
					rangeItem["timeZone"] = effectTimeRange.TimeZone.ValueString()
				}
				if !effectTimeRange.DayInWeek.IsNull() && !effectTimeRange.DayInWeek.IsUnknown() {
					rangeItem["dayInWeek"] = cmsSubscriptionValueInt64List(effectTimeRange.DayInWeek)
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

// cmsSubscriptionIdentityState resolves the planned value into a persistable
// state that carries the created id. Computed values stay null until the
// follow-up read succeeds; configured values are kept so a failed apply still
// leaves a state Terraform can refresh and destroy.
func cmsSubscriptionIdentityState(plan cmsSubscriptionResourceModel, id string) cmsSubscriptionResourceModel {
	state := plan
	state.Id = types.StringValue(id)
	if state.Description.IsUnknown() {
		state.Description = types.StringNull()
	}
	if state.NotifyStrategyId.IsUnknown() {
		state.NotifyStrategyId = types.StringNull()
	}
	if state.Enable.IsUnknown() {
		state.Enable = types.BoolNull()
	}
	if state.SubscribeLegacyEvent.IsUnknown() {
		state.SubscribeLegacyEvent = types.BoolNull()
	}
	state.UserId = types.StringNull()
	state.SubscriptionType = types.StringNull()
	state.RegionId = types.StringNull()
	state.CreateTime = types.StringNull()
	state.UpdateTime = types.StringNull()
	return state
}

func cmsSubscriptionState(object map[string]interface{}, previous cmsSubscriptionResourceModel, regionId string) cmsSubscriptionResourceModel {
	state := cmsSubscriptionResourceModel{
		Id:                     types.StringValue(cmsSubscriptionJSONString(object["subscriptionId"])),
		SubscriptionName:       cmsSubscriptionJSONStringValue(object["subscriptionName"]),
		Description:            cmsSubscriptionStringStateValue(object["description"], previous.Description),
		NotifyStrategyId:       cmsSubscriptionStringStateValue(object["notifyStrategyId"], previous.NotifyStrategyId),
		Workspace:              cmsSubscriptionJSONStringValue(object["workspace"]),
		Enable:                 cmsSubscriptionJSONBoolValue(object["enable"]),
		SubscribeLegacyEvent:   cmsSubscriptionBoolStateValue(object["subscribeLegacyEvent"]),
		FilterSetting:          flattenCmsSubscriptionFilterSetting(object["filterSetting"], previous.FilterSetting),
		PushingSetting:         flattenCmsSubscriptionPushingSetting(object["pushingSetting"], previous.PushingSetting),
		WorkspaceFilterSetting: flattenCmsSubscriptionWorkspaceFilterSetting(object["workspaceFilterSetting"], previous.WorkspaceFilterSetting),
		AgentConfig:            flattenCmsSubscriptionAgentConfig(object["agentConfig"], previous.AgentConfig),
		UserId:                 cmsSubscriptionJSONStringValue(object["userId"]),
		SubscriptionType:       cmsSubscriptionJSONStringValue(object["subscriptionType"]),
		RegionId:               types.StringValue(regionId),
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

func flattenCmsSubscriptionFilterSetting(raw interface{}, previous []cmsSubscriptionFilterSettingModel) []cmsSubscriptionFilterSettingModel {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		if len(previous) > 0 {
			return previous
		}
		return nil
	}
	block := cmsSubscriptionFilterSettingModel{
		Expression: cmsSubscriptionStringStateValue(setting["expression"], cmsSubscriptionPreviousFilterString(previous, "expression")),
		Relation:   cmsSubscriptionStringStateValue(setting["relation"], cmsSubscriptionPreviousFilterString(previous, "relation")),
	}
	previousConditions := cmsSubscriptionPreviousFilterConditions(previous)
	if items, ok := setting["conditions"].([]interface{}); ok {
		conditions := make([]cmsSubscriptionConditionModel, 0, len(items))
		for _, item := range items {
			condition, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			conditions = append(conditions, cmsSubscriptionConditionModel{
				Field: cmsSubscriptionJSONStringValue(condition["field"]),
				Op:    cmsSubscriptionJSONStringValue(condition["op"]),
				Value: cmsSubscriptionJSONStringValue(condition["value"]),
			})
		}
		block.Conditions = conditions
		if len(previousConditions) == len(conditions) && cmsSubscriptionConditionsSameSet(previousConditions, conditions) {
			// Keep the configured order when the service returns the same set in a
			// different order, so the ordered list does not produce a permadiff.
			block.Conditions = previousConditions
		}
	} else {
		block.Conditions = previousConditions
	}
	return []cmsSubscriptionFilterSettingModel{block}
}

func cmsSubscriptionConditionsSameSet(left, right []cmsSubscriptionConditionModel) bool {
	if len(left) != len(right) {
		return false
	}
	used := make([]bool, len(right))
	for _, want := range left {
		matched := false
		for i, got := range right {
			if used[i] {
				continue
			}
			if want.Field.Equal(got.Field) && want.Op.Equal(got.Op) && want.Value.Equal(got.Value) {
				used[i] = true
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

func cmsSubscriptionPreviousFilterString(previous []cmsSubscriptionFilterSettingModel, attribute string) types.String {
	if len(previous) == 0 {
		return types.StringNull()
	}
	if attribute == "expression" {
		return previous[0].Expression
	}
	return previous[0].Relation
}

func cmsSubscriptionPreviousFilterConditions(previous []cmsSubscriptionFilterSettingModel) []cmsSubscriptionConditionModel {
	if len(previous) == 0 {
		return nil
	}
	return previous[0].Conditions
}

func flattenCmsSubscriptionPushingSetting(raw interface{}, previous []cmsSubscriptionPushingSettingModel) []cmsSubscriptionPushingSettingModel {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		if len(previous) > 0 {
			return previous
		}
		return nil
	}
	block := cmsSubscriptionPushingSettingModel{
		AlertActionIds:   cmsSubscriptionJSONStringListState(setting["alertActionIds"], cmsSubscriptionPreviousPushingList(previous, "alert")),
		RestoreActionIds: cmsSubscriptionJSONStringListState(setting["restoreActionIds"], cmsSubscriptionPreviousPushingList(previous, "restore")),
		TemplateUuid:     cmsSubscriptionStringStateValue(setting["templateUuid"], cmsSubscriptionPreviousPushingString(previous, "template")),
		ResponsePlanId:   cmsSubscriptionStringStateValue(setting["responsePlanId"], cmsSubscriptionPreviousPushingString(previous, "response")),
	}
	return []cmsSubscriptionPushingSettingModel{block}
}

func cmsSubscriptionPreviousPushingList(previous []cmsSubscriptionPushingSettingModel, attribute string) types.List {
	if len(previous) == 0 {
		return types.ListNull(types.ListType{ElemType: types.StringType})
	}
	if attribute == "alert" {
		return previous[0].AlertActionIds
	}
	return previous[0].RestoreActionIds
}

func cmsSubscriptionPreviousPushingString(previous []cmsSubscriptionPushingSettingModel, attribute string) types.String {
	if len(previous) == 0 {
		return types.StringNull()
	}
	if attribute == "template" {
		return previous[0].TemplateUuid
	}
	return previous[0].ResponsePlanId
}

func flattenCmsSubscriptionWorkspaceFilterSetting(raw interface{}, previous []cmsSubscriptionWorkspaceFilterModel) []cmsSubscriptionWorkspaceFilterModel {
	setting, ok := raw.(map[string]interface{})
	if !ok || len(setting) == 0 {
		if len(previous) > 0 {
			return previous
		}
		return nil
	}
	previousList := types.ListNull(types.ListType{ElemType: types.StringType})
	if len(previous) > 0 {
		previousList = previous[0].WorkspaceUuids
	}
	return []cmsSubscriptionWorkspaceFilterModel{{
		WorkspaceUuids: cmsSubscriptionJSONStringListState(setting["workspaceUuids"], previousList),
	}}
}

func flattenCmsSubscriptionAgentConfig(raw interface{}, previous []cmsSubscriptionAgentConfigModel) []cmsSubscriptionAgentConfigModel {
	config, ok := raw.(map[string]interface{})
	if !ok || len(config) == 0 {
		if len(previous) > 0 {
			return previous
		}
		return nil
	}
	previousUuid := types.StringNull()
	if len(previous) > 0 {
		previousUuid = previous[0].AgentUuid
	}
	block := cmsSubscriptionAgentConfigModel{
		AgentUuid: cmsSubscriptionStringStateValue(config["agentUuid"], previousUuid),
	}
	if items, ok := config["routes"].([]interface{}); ok {
		routes := make([]cmsSubscriptionAgentRouteModel, 0, len(items))
		for _, item := range items {
			route, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			routes = append(routes, flattenCmsSubscriptionAgentRoute(route))
		}
		block.Routes = routes
	} else if len(previous) > 0 {
		block.Routes = previous[0].Routes
	}
	return []cmsSubscriptionAgentConfigModel{block}
}

func flattenCmsSubscriptionAgentRoute(route map[string]interface{}) cmsSubscriptionAgentRouteModel {
	result := cmsSubscriptionAgentRouteModel{}
	if items, ok := route["channels"].([]interface{}); ok {
		channels := make([]cmsSubscriptionAgentChannelModel, 0, len(items))
		for _, item := range items {
			channel, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			channels = append(channels, cmsSubscriptionAgentChannelModel{
				ChannelType:        cmsSubscriptionJSONStringValue(channel["channelType"]),
				Receivers:          cmsSubscriptionJSONStringList(channel["receivers"]),
				EnabledSubChannels: cmsSubscriptionJSONStringList(channel["enabledSubChannels"]),
			})
		}
		result.Channels = channels
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

func cmsSubscriptionStringStateValue(v interface{}, previous types.String) types.String {
	if s := cmsSubscriptionJSONString(v); s != "" {
		return types.StringValue(s)
	}
	if !previous.IsNull() && !previous.IsUnknown() {
		return types.StringValue("")
	}
	return types.StringNull()
}

func cmsSubscriptionJSONBoolValue(v interface{}) types.Bool {
	if b, ok := v.(bool); ok {
		return types.BoolValue(b)
	}
	return types.BoolNull()
}

func cmsSubscriptionBoolStateValue(v interface{}) types.Bool {
	if b, ok := v.(bool); ok {
		return types.BoolValue(b)
	}
	// The contract states that null and false both mean "not subscribed", so an
	// omitted response materializes the service default instead of a null that
	// would differ between apply and import.
	return types.BoolValue(false)
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

func cmsSubscriptionJSONStringListState(v interface{}, previous types.List) types.List {
	list := cmsSubscriptionJSONStringList(v)
	if !list.IsNull() {
		return list
	}
	if !previous.IsNull() && !previous.IsUnknown() {
		return types.ListValueMust(types.ListType{ElemType: types.StringType}, []attr.Value{})
	}
	return list
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

type cmsSubscriptionChannelValidator struct{}

func (v cmsSubscriptionChannelValidator) Description(ctx context.Context) string {
	return "enabled_sub_channels is required when channel_type is CONTACT, GROUP or DUTY"
}

func (v cmsSubscriptionChannelValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v cmsSubscriptionChannelValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config cmsSubscriptionResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	for i, agent := range config.AgentConfig {
		for j, route := range agent.Routes {
			for k, channel := range route.Channels {
				if channel.ChannelType.IsNull() || channel.ChannelType.IsUnknown() {
					continue
				}
				switch channel.ChannelType.ValueString() {
				case "CONTACT", "GROUP", "DUTY":
					if channel.EnabledSubChannels.IsNull() || channel.EnabledSubChannels.IsUnknown() || len(channel.EnabledSubChannels.Elements()) == 0 {
						resp.Diagnostics.AddAttributeError(
							path.Root("agent_config").AtListIndex(i).
								AtName("routes").AtListIndex(j).
								AtName("channels").AtListIndex(k).
								AtName("enabled_sub_channels"),
							"Missing enabled_sub_channels",
							fmt.Sprintf("enabled_sub_channels must contain at least one value when channel_type is %s.", channel.ChannelType.ValueString()),
						)
					}
				}
			}
		}
	}
}
