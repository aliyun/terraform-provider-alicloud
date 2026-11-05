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

const subscriptionTypeName = "alicloud_cms_subscription"

// terraform-plugin-framework has no user-configurable resource timeouts yet,
// so the SDKv2 default of five minutes per operation is kept as a constant.
const subscriptionTimeout = 5 * time.Minute

var (
	_ resource.Resource                     = &subscriptionResource{}
	_ resource.ResourceWithImportState      = &subscriptionResource{}
	_ resource.ResourceWithConfigValidators = &subscriptionResource{}
)

func NewSubscriptionResource() resource.Resource {
	return &subscriptionResource{}
}

type subscriptionResource struct {
	fwadapt.ResourceBase
}

type subscriptionResourceModel struct {
	Id                     types.String                       `tfsdk:"id"`
	SubscriptionName       types.String                       `tfsdk:"subscription_name"`
	Description            types.String                       `tfsdk:"description"`
	NotifyStrategyId       types.String                       `tfsdk:"notify_strategy_id"`
	Workspace              types.String                       `tfsdk:"workspace"`
	Enable                 types.Bool                         `tfsdk:"enable"`
	SubscribeLegacyEvent   types.Bool                         `tfsdk:"subscribe_legacy_event"`
	FilterSetting          []subscriptionFilterSettingModel   `tfsdk:"filter_setting"`
	PushingSetting         []subscriptionPushingSettingModel  `tfsdk:"pushing_setting"`
	WorkspaceFilterSetting []subscriptionWorkspaceFilterModel `tfsdk:"workspace_filter_setting"`
	AgentConfig            []subscriptionAgentConfigModel     `tfsdk:"agent_config"`
	UserId                 types.String                       `tfsdk:"user_id"`
	SubscriptionType       types.String                       `tfsdk:"subscription_type"`
	RegionId               types.String                       `tfsdk:"region_id"`
	CreateTime             types.String                       `tfsdk:"create_time"`
	UpdateTime             types.String                       `tfsdk:"update_time"`
}

type subscriptionFilterSettingModel struct {
	Conditions []subscriptionConditionModel `tfsdk:"conditions"`
	Expression types.String                 `tfsdk:"expression"`
	Relation   types.String                 `tfsdk:"relation"`
}

type subscriptionConditionModel struct {
	Field types.String `tfsdk:"field"`
	Op    types.String `tfsdk:"op"`
	Value types.String `tfsdk:"value"`
}

type subscriptionPushingSettingModel struct {
	AlertActionIds   types.List   `tfsdk:"alert_action_ids"`
	RestoreActionIds types.List   `tfsdk:"restore_action_ids"`
	TemplateUuid     types.String `tfsdk:"template_uuid"`
	ResponsePlanId   types.String `tfsdk:"response_plan_id"`
}

type subscriptionWorkspaceFilterModel struct {
	WorkspaceUuids types.List `tfsdk:"workspace_uuids"`
}

type subscriptionAgentConfigModel struct {
	AgentUuid types.String                  `tfsdk:"agent_uuid"`
	Routes    []subscriptionAgentRouteModel `tfsdk:"routes"`
}

type subscriptionAgentRouteModel struct {
	Channels        []subscriptionAgentChannelModel         `tfsdk:"channels"`
	EffectTimeRange []subscriptionAgentEffectTimeRangeModel `tfsdk:"effect_time_range"`
}

type subscriptionAgentChannelModel struct {
	ChannelType        types.String `tfsdk:"channel_type"`
	Receivers          types.List   `tfsdk:"receivers"`
	EnabledSubChannels types.List   `tfsdk:"enabled_sub_channels"`
}

type subscriptionAgentEffectTimeRangeModel struct {
	TimeZone          types.String `tfsdk:"time_zone"`
	DayInWeek         types.List   `tfsdk:"day_in_week"`
	StartTimeInMinute types.Int64  `tfsdk:"start_time_in_minute"`
	EndTimeInMinute   types.Int64  `tfsdk:"end_time_in_minute"`
}

func (r *subscriptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = subscriptionTypeName
}

func (r *subscriptionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{subscriptionChannelValidator{}}
}

func (r *subscriptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Default:  booldefault.StaticBool(false),
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
							// The service defaults the relation to AND and echoes it even
							// when the configuration uses expression instead.
							Optional: true,
							Computed: true,
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

func (r *subscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, subscriptionTimeout)
	defer cancel()

	action := "/subscriptions"
	request := expandSubscriptionRequest(plan)
	response, err := r.roaCall(ctx, action, request, func() (map[string]interface{}, error) {
		return r.Meta().RoaPost("Cms", "2024-03-30", action, subscriptionQuery(plan.Workspace), nil, request, true)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id, _ := response["subscriptionId"].(string)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSubscription", subscriptionTypeName),
			"CreateSubscription returned an empty subscriptionId; the subscription was created but cannot be tracked in state.",
		)
		return
	}

	// Persist the identity as soon as it is known: every later step can fail,
	// and without it the created subscription could neither be refreshed nor destroyed.
	identity := subscriptionIdentityState(plan, id)
	resp.Diagnostics.Append(resp.State.Set(ctx, identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.waitSubscriptionReadable(ctx, id, plan.Workspace, request, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	if !plan.Enable.IsNull() && !plan.Enable.IsUnknown() && !plan.Enable.ValueBool() {
		if err := r.setSubscriptionEnable(ctx, id, plan.Workspace, false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling DisableSubscription", subscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, "/subscriptions/"+id+"/disable", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		// The converged read already reports the requested state, so it is used
		// directly instead of issuing a second, unverified read.
		object, err = r.waitSubscriptionEnable(ctx, id, plan.Workspace, false)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: waiting for DisableSubscription", subscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	state, err := subscriptionState(object, identity, identity, r.Meta().RegionId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *subscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, subscriptionTimeout)
	defer cancel()

	object, err := NewCmsService(r.Meta()).DescribeSubscription(ctx, state.Id.ValueString(), state.Workspace.ValueString())
	if err != nil {
		if subscriptionNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, "/subscriptions/"+state.Id.ValueString(), errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	next, err := subscriptionState(object, state, state, r.Meta().RegionId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, next)...)
}

func (r *subscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state subscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, subscriptionTimeout)
	defer cancel()

	id := state.Id.ValueString()
	var object map[string]interface{}
	var err error

	if !plan.Enable.Equal(state.Enable) && !plan.Enable.IsNull() && !plan.Enable.IsUnknown() {
		enable := plan.Enable.ValueBool()
		suffix := "enable"
		if !enable {
			suffix = "disable"
		}
		if err := r.setSubscriptionEnable(ctx, id, state.Workspace, enable); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling /%s", subscriptionTypeName, suffix),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, "/subscriptions/"+id+"/"+suffix, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		object, err = r.waitSubscriptionEnable(ctx, id, state.Workspace, enable)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: waiting for the subscription enabled state", subscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	action := fmt.Sprintf("/subscriptions/%s", id)
	request := subscriptionRequestForConfig(plan, state)
	removed := subscriptionRemovedPaths(plan, state)
	if len(removed) == 0 && subscriptionRequestsEqual(expandSubscriptionRequest(state), request) {
		// Only the enable/disable state changed: the body is untouched, so do not
		// rewrite the configuration with a full-body update.
		if object == nil {
			object, err = NewCmsService(r.Meta()).DescribeSubscription(ctx, id, state.Workspace.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
					errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
				)
				return
			}
		}
		next, err := subscriptionState(object, subscriptionIntent(plan, state), state, r.Meta().RegionId)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
				err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, next)...)
		return
	}

	_, err = r.roaCall(ctx, action, request, func() (map[string]interface{}, error) {
		return r.Meta().RoaPut("Cms", "2024-03-30", action, subscriptionQuery(state.Workspace), nil, request, true)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err = r.waitSubscriptionUpdate(ctx, id, state.Workspace, request, removed)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	next, err := subscriptionState(object, subscriptionIntent(plan, state), state, r.Meta().RegionId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, next)...)
}

func (r *subscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, subscriptionTimeout)
	defer cancel()

	action := fmt.Sprintf("/subscriptions/%s", state.Id.ValueString())
	var rawErr error
	_, err := r.roaCall(ctx, action, subscriptionQuery(state.Workspace), func() (map[string]interface{}, error) {
		response, callErr := r.Meta().RoaDelete("Cms", "2024-03-30", action, subscriptionQuery(state.Workspace), nil, nil, true)
		rawErr = callErr
		return response, callErr
	})
	if err != nil {
		if subscriptionIsResourceNotFound(rawErr) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
	}
}

// ImportState accepts <workspace>:<subscription_id> because every subscription
// request, including the read used to refresh imported state, is scoped by its
// workspace.
func (r *subscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// roaCall runs one Cms write with a bounded retry loop. The caller owns the
// deadline: every attempt, every backoff and every follow-up read share the
// same context, so a cancelled operation stops issuing requests.
func (r *subscriptionResource) roaCall(ctx context.Context, action string, logPayload interface{}, call func() (map[string]interface{}, error)) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := subscriptionBackoff(ctx)
	err = retry.RetryContext(ctx, subscriptionTimeout, func() *retry.RetryError {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return retry.NonRetryableError(ctxErr)
		}
		response, err = call()
		if err != nil {
			if helper.NeedRetry(err) {
				if !wait() {
					return retry.NonRetryableError(ctx.Err())
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, logPayload)
	return response, err
}

func (r *subscriptionResource) setSubscriptionEnable(ctx context.Context, id string, workspace types.String, enable bool) error {
	suffix := "enable"
	if !enable {
		suffix = "disable"
	}
	action := fmt.Sprintf("/subscriptions/%s/%s", id, suffix)

	_, err := r.roaCall(ctx, action, action, func() (map[string]interface{}, error) {
		return r.Meta().RoaPut("Cms", "2024-03-30", action, subscriptionQuery(workspace), nil, nil, true)
	})
	return err
}

// waitSubscriptionEnable polls GetSubscription until the enabled state reported
// by the service matches the requested one and returns that converged object.
// The enable and disable operations return before the change is visible on the
// read path.
func (r *subscriptionResource) waitSubscriptionEnable(ctx context.Context, id string, workspace types.String, enable bool) (map[string]interface{}, error) {
	for {
		object, err := NewCmsService(r.Meta()).DescribeSubscription(ctx, id, workspace.ValueString())
		if err != nil {
			return nil, err
		}
		if current, ok := object["enable"].(bool); ok && current == enable {
			return object, nil
		}
		if err := subscriptionSleep(ctx, 3*time.Second); err != nil {
			return nil, err
		}
	}
}

// waitSubscriptionReadable polls GetSubscription until the freshly created
// subscription is visible on the read path and reflects the requested
// configuration. Creation is asynchronous, so an immediate read can still
// report the subscription as not found or return a partially applied object.
func (r *subscriptionResource) waitSubscriptionReadable(ctx context.Context, id string, workspace types.String, request map[string]interface{}, removed []subscriptionPath) (map[string]interface{}, error) {
	for {
		object, err := NewCmsService(r.Meta()).DescribeSubscription(ctx, id, workspace.ValueString())
		if err == nil {
			if cleared, clearErr := subscriptionPathsCleared(object, removed); clearErr != nil {
				return nil, clearErr
			} else if cleared && subscriptionRequestVisible(object, request) {
				return object, nil
			}
		} else if !subscriptionNotFound(err) {
			return nil, err
		}
		if err := subscriptionSleep(ctx, 3*time.Second); err != nil {
			return nil, err
		}
	}
}

// waitSubscriptionUpdate polls GetSubscription until the configuration sent
// to UpdateSubscription is visible on the read path and every removed key is
// gone. The Cms event subscription reads are eventually consistent, so an
// immediate read can return the values from before the update.
func (r *subscriptionResource) waitSubscriptionUpdate(ctx context.Context, id string, workspace types.String, request map[string]interface{}, removed []subscriptionPath) (map[string]interface{}, error) {
	for {
		object, err := NewCmsService(r.Meta()).DescribeSubscription(ctx, id, workspace.ValueString())
		if err == nil {
			if cleared, clearErr := subscriptionPathsCleared(object, removed); clearErr != nil {
				return nil, clearErr
			} else if cleared && subscriptionRequestVisible(object, request) {
				return object, nil
			}
		} else if !subscriptionNotFound(err) {
			return nil, err
		}
		if err := subscriptionSleep(ctx, 3*time.Second); err != nil {
			return nil, err
		}
	}
}

func subscriptionSleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func subscriptionBackoff(ctx context.Context) func() bool {
	attempt := 0
	return func() bool {
		delay := time.Duration(attempt)*2*time.Second + 3*time.Second
		attempt++
		return subscriptionSleep(ctx, delay) == nil
	}
}

// subscriptionRequestVisible reports whether every value sent to the API is
// already reflected by the read response. It is a subset comparison so fields
// the service adds on its own (ids, timestamps, defaults) do not block
// convergence, while every managed member, list element and nested block does.
// Lists are compared position by position: conditions are referenced by their
// one-based position in `expression`, and every list in the schema is ordered.
func subscriptionRequestVisible(object, request map[string]interface{}) bool {
	for key, want := range request {
		if !subscriptionValueVisible(object[key], want) {
			return false
		}
	}
	return true
}

func subscriptionValueVisible(actual, want interface{}) bool {
	switch wantValue := want.(type) {
	case map[string]interface{}:
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			return false
		}
		for key, value := range wantValue {
			if !subscriptionValueVisible(actualMap[key], value) {
				return false
			}
		}
		return true
	case []map[string]interface{}:
		actualItems, ok := subscriptionJSONItems(actual)
		if !ok || len(actualItems) != len(wantValue) {
			return false
		}
		for index, wantItem := range wantValue {
			actualItem, ok := actualItems[index].(map[string]interface{})
			if !ok || !subscriptionValueVisible(actualItem, wantItem) {
				return false
			}
		}
		return true
	case []string:
		actualItems, ok := subscriptionJSONItems(actual)
		if !ok || len(actualItems) != len(wantValue) {
			return false
		}
		for index, wantItem := range wantValue {
			got, ok := actualItems[index].(string)
			if !ok || got != wantItem {
				return false
			}
		}
		return true
	case []int64:
		actualItems, ok := subscriptionJSONItems(actual)
		if !ok || len(actualItems) != len(wantValue) {
			return false
		}
		for index, wantItem := range wantValue {
			got, ok := subscriptionJSONInt64(actualItems[index])
			if !ok || got != wantItem {
				return false
			}
		}
		return true
	case string:
		if wantValue == "" {
			// The service may drop empty strings from the response.
			return subscriptionValueEmpty(actual)
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
		got, ok := subscriptionJSONInt64(actual)
		return ok && got == wantValue
	case nil:
		return subscriptionValueEmpty(actual)
	}
	return false
}

func subscriptionJSONItems(actual interface{}) ([]interface{}, bool) {
	switch items := actual.(type) {
	case nil:
		return nil, true
	case []interface{}:
		return items, true
	case []map[string]interface{}:
		result := make([]interface{}, 0, len(items))
		for _, item := range items {
			result = append(result, item)
		}
		return result, true
	}
	return nil, false
}

func subscriptionJSONInt64(v interface{}) (int64, bool) {
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

func subscriptionValueEmpty(v interface{}) bool {
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

// subscriptionPath is one location inside the API request and response. String
// elements address object members, index elements address ordered list members.
type subscriptionPath []subscriptionPathElement

type subscriptionPathElement struct {
	key     string
	index   int
	isIndex bool
}

func (p subscriptionPath) describe(length int) string {
	var builder strings.Builder
	for _, element := range p[:length] {
		if element.isIndex {
			fmt.Fprintf(&builder, "[%d]", element.index)
			continue
		}
		if builder.Len() > 0 {
			builder.WriteString(".")
		}
		builder.WriteString(element.key)
	}
	return builder.String()
}

// subscriptionPathsCleared reports whether every location removed from the
// configuration no longer carries a value on the read path. A wrong container
// shape is an error, not an empty value: a malformed response must not be
// mistaken for a cleared member.
func subscriptionPathsCleared(object map[string]interface{}, paths []subscriptionPath) (bool, error) {
	for _, fieldPath := range paths {
		cleared, err := subscriptionPathCleared(object, fieldPath)
		if err != nil {
			return false, err
		}
		if !cleared {
			return false, nil
		}
	}
	return true, nil
}

func subscriptionPathCleared(object map[string]interface{}, fieldPath subscriptionPath) (bool, error) {
	var current interface{} = object
	for position, element := range fieldPath {
		if element.isIndex {
			items, ok := subscriptionJSONItems(current)
			if !ok {
				return false, fmt.Errorf("the response value at %s is not a list", fieldPath.describe(position))
			}
			if element.index >= len(items) {
				return true, nil
			}
			current = items[element.index]
			continue
		}
		container, ok := current.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("the response value at %s is not an object", fieldPath.describe(position))
		}
		value, exists := container[element.key]
		if !exists {
			return true, nil
		}
		current = value
	}
	return subscriptionValueEmpty(current), nil
}

// subscriptionIntent resolves the plan for comparison and read-back: the
// Service echoes a relation even when the configuration omits it, so an
// unknown computed relation falls back to the value already in state and does
// not look like a body change.
func subscriptionIntent(plan, state subscriptionResourceModel) subscriptionResourceModel {
	resolved := plan
	if len(plan.FilterSetting) > 0 && len(state.FilterSetting) > 0 {
		if plan.FilterSetting[0].Relation.IsUnknown() {
			resolved.FilterSetting = append([]subscriptionFilterSettingModel(nil), plan.FilterSetting...)
			resolved.FilterSetting[0].Relation = state.FilterSetting[0].Relation
		}
	}
	return resolved
}

func subscriptionRequestForConfig(plan, state subscriptionResourceModel) map[string]interface{} {
	return expandSubscriptionRequest(subscriptionIntent(plan, state))
}

// subscriptionRemovedPaths lists the request locations that were present in
// state and are absent from the new plan, including members of ordered lists,
// so the update waiter can require them to be cleared instead of accepting a
// stale response.
func subscriptionRemovedPaths(plan, state subscriptionResourceModel) []subscriptionPath {
	before := expandSubscriptionRequest(state)
	after := subscriptionRequestForConfig(plan, state)
	var paths []subscriptionPath
	subscriptionCollectRemoved(before, after, nil, &paths)
	return paths
}

func subscriptionCollectRemoved(before, after interface{}, prefix subscriptionPath, out *[]subscriptionPath) {
	switch beforeValue := before.(type) {
	case map[string]interface{}:
		afterMap, ok := after.(map[string]interface{})
		if !ok {
			return
		}
		for key, value := range beforeValue {
			keyPath := append(append(subscriptionPath{}, prefix...), subscriptionPathElement{key: key})
			next, exists := afterMap[key]
			if !exists {
				*out = append(*out, keyPath)
				continue
			}
			subscriptionCollectRemoved(value, next, keyPath, out)
		}
	case []map[string]interface{}:
		afterItems, ok := after.([]map[string]interface{})
		if !ok {
			return
		}
		for index := range beforeValue {
			memberPath := append(append(subscriptionPath{}, prefix...), subscriptionPathElement{index: index, isIndex: true})
			if index >= len(afterItems) {
				*out = append(*out, memberPath)
				continue
			}
			subscriptionCollectRemoved(beforeValue[index], afterItems[index], memberPath, out)
		}
	case []string:
		afterItems, ok := after.([]string)
		if !ok {
			return
		}
		subscriptionCollectRemovedTail(len(afterItems), len(beforeValue), prefix, out)
	case []int64:
		afterItems, ok := after.([]int64)
		if !ok {
			return
		}
		subscriptionCollectRemovedTail(len(afterItems), len(beforeValue), prefix, out)
	}
}

func subscriptionCollectRemovedTail(afterLength, beforeLength int, prefix subscriptionPath, out *[]subscriptionPath) {
	for index := afterLength; index < beforeLength; index++ {
		*out = append(*out, append(append(subscriptionPath{}, prefix...), subscriptionPathElement{index: index, isIndex: true}))
	}
}

func subscriptionRequestsEqual(left, right map[string]interface{}) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && string(leftJSON) == string(rightJSON)
}

func subscriptionQuery(workspace types.String) map[string]*string {
	query := make(map[string]*string)
	if !workspace.IsNull() && !workspace.IsUnknown() && workspace.ValueString() != "" {
		query["workspace"] = helper.StringPointer(workspace.ValueString())
	}
	return query
}

// expandSubscriptionRequest always returns the complete body: UpdateSubscription
// replaces the stored configuration, so unchanged attributes must be re-sent.
// Values that are configured but empty ("", [] or false) are sent explicitly so
// the service can clear them instead of keeping the previous value.
func expandSubscriptionRequest(plan subscriptionResourceModel) map[string]interface{} {
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
	if v := expandSubscriptionFilterSetting(plan.FilterSetting); v != nil {
		request["filterSetting"] = v
	}
	if v := expandSubscriptionPushingSetting(plan.PushingSetting); v != nil {
		request["pushingSetting"] = v
	}
	if v := expandSubscriptionWorkspaceFilterSetting(plan.WorkspaceFilterSetting); v != nil {
		request["workspaceFilterSetting"] = v
	}
	if v := expandSubscriptionAgentConfig(plan.AgentConfig); v != nil {
		request["agentConfig"] = v
	}
	return request
}

func expandSubscriptionFilterSetting(blocks []subscriptionFilterSettingModel) map[string]interface{} {
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
				"field": subscriptionValueString(condition.Field),
				"op":    subscriptionValueString(condition.Op),
				"value": subscriptionValueString(condition.Value),
			})
		}
		result["conditions"] = conditions
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandSubscriptionPushingSetting(blocks []subscriptionPushingSettingModel) map[string]interface{} {
	if len(blocks) == 0 {
		return nil
	}
	block := blocks[0]
	result := map[string]interface{}{}
	if !block.AlertActionIds.IsNull() && !block.AlertActionIds.IsUnknown() {
		result["alertActionIds"] = subscriptionValueStringList(block.AlertActionIds)
	}
	if !block.RestoreActionIds.IsNull() && !block.RestoreActionIds.IsUnknown() {
		result["restoreActionIds"] = subscriptionValueStringList(block.RestoreActionIds)
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

func expandSubscriptionWorkspaceFilterSetting(blocks []subscriptionWorkspaceFilterModel) map[string]interface{} {
	if len(blocks) == 0 {
		return nil
	}
	result := map[string]interface{}{}
	if !blocks[0].WorkspaceUuids.IsNull() && !blocks[0].WorkspaceUuids.IsUnknown() {
		result["workspaceUuids"] = subscriptionValueStringList(blocks[0].WorkspaceUuids)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func expandSubscriptionAgentConfig(blocks []subscriptionAgentConfigModel) map[string]interface{} {
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
						channelItem["receivers"] = subscriptionValueStringList(channel.Receivers)
					}
					if !channel.EnabledSubChannels.IsNull() && !channel.EnabledSubChannels.IsUnknown() {
						channelItem["enabledSubChannels"] = subscriptionValueStringList(channel.EnabledSubChannels)
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
					rangeItem["dayInWeek"] = subscriptionValueInt64List(effectTimeRange.DayInWeek)
				}
				if v, ok := subscriptionValueInt64(effectTimeRange.StartTimeInMinute); ok {
					rangeItem["startTimeInMinute"] = v
				}
				if v, ok := subscriptionValueInt64(effectTimeRange.EndTimeInMinute); ok {
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

// subscriptionIdentityState resolves the planned value into a persistable state
// that carries the created id. Unknown values, including nested ones, become
// null so a failed apply never returns a state with unknown members.
func subscriptionIdentityState(plan subscriptionResourceModel, id string) subscriptionResourceModel {
	state := plan
	state.Id = types.StringValue(id)
	state.Description = subscriptionKnownString(plan.Description)
	state.NotifyStrategyId = subscriptionKnownString(plan.NotifyStrategyId)
	state.Enable = subscriptionKnownBool(plan.Enable)
	state.SubscribeLegacyEvent = subscriptionKnownBool(plan.SubscribeLegacyEvent)
	state.FilterSetting = subscriptionKnownFilterSetting(plan.FilterSetting)
	state.PushingSetting = subscriptionKnownPushingSetting(plan.PushingSetting)
	state.WorkspaceFilterSetting = subscriptionKnownWorkspaceFilterSetting(plan.WorkspaceFilterSetting)
	state.AgentConfig = subscriptionKnownAgentConfig(plan.AgentConfig)
	state.UserId = types.StringNull()
	state.SubscriptionType = types.StringNull()
	state.RegionId = types.StringNull()
	state.CreateTime = types.StringNull()
	state.UpdateTime = types.StringNull()
	return state
}

func subscriptionKnownString(v types.String) types.String {
	if v.IsUnknown() {
		return types.StringNull()
	}
	return v
}

func subscriptionKnownBool(v types.Bool) types.Bool {
	if v.IsUnknown() {
		return types.BoolNull()
	}
	return v
}

func subscriptionKnownStringList(v types.List) types.List {
	if v.IsUnknown() {
		return types.ListNull(types.StringType)
	}
	return v
}

func subscriptionKnownInt64List(v types.List) types.List {
	if v.IsUnknown() {
		return types.ListNull(types.Int64Type)
	}
	return v
}

func subscriptionKnownInt64(v types.Int64) types.Int64 {
	if v.IsUnknown() {
		return types.Int64Null()
	}
	return v
}

func subscriptionKnownFilterSetting(blocks []subscriptionFilterSettingModel) []subscriptionFilterSettingModel {
	if len(blocks) == 0 {
		return nil
	}
	result := make([]subscriptionFilterSettingModel, 0, len(blocks))
	for _, block := range blocks {
		conditions := make([]subscriptionConditionModel, 0, len(block.Conditions))
		for _, condition := range block.Conditions {
			conditions = append(conditions, subscriptionConditionModel{
				Field: subscriptionKnownString(condition.Field),
				Op:    subscriptionKnownString(condition.Op),
				Value: subscriptionKnownString(condition.Value),
			})
		}
		result = append(result, subscriptionFilterSettingModel{
			Expression: subscriptionKnownString(block.Expression),
			Relation:   subscriptionKnownString(block.Relation),
			Conditions: conditions,
		})
	}
	return result
}

func subscriptionKnownPushingSetting(blocks []subscriptionPushingSettingModel) []subscriptionPushingSettingModel {
	if len(blocks) == 0 {
		return nil
	}
	result := make([]subscriptionPushingSettingModel, 0, len(blocks))
	for _, block := range blocks {
		result = append(result, subscriptionPushingSettingModel{
			AlertActionIds:   subscriptionKnownStringList(block.AlertActionIds),
			RestoreActionIds: subscriptionKnownStringList(block.RestoreActionIds),
			TemplateUuid:     subscriptionKnownString(block.TemplateUuid),
			ResponsePlanId:   subscriptionKnownString(block.ResponsePlanId),
		})
	}
	return result
}

func subscriptionKnownWorkspaceFilterSetting(blocks []subscriptionWorkspaceFilterModel) []subscriptionWorkspaceFilterModel {
	if len(blocks) == 0 {
		return nil
	}
	result := make([]subscriptionWorkspaceFilterModel, 0, len(blocks))
	for _, block := range blocks {
		result = append(result, subscriptionWorkspaceFilterModel{
			WorkspaceUuids: subscriptionKnownStringList(block.WorkspaceUuids),
		})
	}
	return result
}

func subscriptionKnownAgentConfig(blocks []subscriptionAgentConfigModel) []subscriptionAgentConfigModel {
	if len(blocks) == 0 {
		return nil
	}
	result := make([]subscriptionAgentConfigModel, 0, len(blocks))
	for _, block := range blocks {
		routes := make([]subscriptionAgentRouteModel, 0, len(block.Routes))
		for _, route := range block.Routes {
			channels := make([]subscriptionAgentChannelModel, 0, len(route.Channels))
			for _, channel := range route.Channels {
				channels = append(channels, subscriptionAgentChannelModel{
					ChannelType:        subscriptionKnownString(channel.ChannelType),
					Receivers:          subscriptionKnownStringList(channel.Receivers),
					EnabledSubChannels: subscriptionKnownStringList(channel.EnabledSubChannels),
				})
			}
			effectTimeRange := make([]subscriptionAgentEffectTimeRangeModel, 0, len(route.EffectTimeRange))
			for _, window := range route.EffectTimeRange {
				effectTimeRange = append(effectTimeRange, subscriptionAgentEffectTimeRangeModel{
					TimeZone:          subscriptionKnownString(window.TimeZone),
					DayInWeek:         subscriptionKnownInt64List(window.DayInWeek),
					StartTimeInMinute: subscriptionKnownInt64(window.StartTimeInMinute),
					EndTimeInMinute:   subscriptionKnownInt64(window.EndTimeInMinute),
				})
			}
			routes = append(routes, subscriptionAgentRouteModel{
				Channels:        channels,
				EffectTimeRange: effectTimeRange,
			})
		}
		result = append(result, subscriptionAgentConfigModel{
			AgentUuid: subscriptionKnownString(block.AgentUuid),
			Routes:    routes,
		})
	}
	return result
}

// subscriptionState converts one GetSubscription response into state. The
// confirmed identity is authoritative: a response that omits it keeps the
// confirmed value, and a response that disagrees is an error rather than a
// silent overwrite. intent carries the configuration that produced this
// request, so an argument that was removed stays null instead of being turned
// into an empty value.
func subscriptionState(object map[string]interface{}, intent, identity subscriptionResourceModel, regionId string) (subscriptionResourceModel, error) {
	expectedId := subscriptionValueString(identity.Id)
	actualId := subscriptionJSONString(object["subscriptionId"])
	if actualId == "" {
		actualId = expectedId
	}
	if actualId == "" {
		return subscriptionResourceModel{}, fmt.Errorf("the response does not contain the subscription id and none was confirmed")
	}
	if expectedId != "" && actualId != expectedId {
		return subscriptionResourceModel{}, fmt.Errorf("the response reports subscription id %q but %q was requested", actualId, expectedId)
	}

	expectedWorkspace := subscriptionValueString(identity.Workspace)
	actualWorkspace := subscriptionJSONString(object["workspace"])
	if actualWorkspace == "" {
		actualWorkspace = expectedWorkspace
	}
	if actualWorkspace == "" {
		return subscriptionResourceModel{}, fmt.Errorf("the response does not contain the workspace and none was configured")
	}
	if expectedWorkspace != "" && actualWorkspace != expectedWorkspace {
		return subscriptionResourceModel{}, fmt.Errorf("the response reports workspace %q but %q was requested", actualWorkspace, expectedWorkspace)
	}

	filterSetting, err := flattenSubscriptionFilterSetting(object["filterSetting"], intent.FilterSetting)
	if err != nil {
		return subscriptionResourceModel{}, err
	}
	pushingSetting, err := flattenSubscriptionPushingSetting(object["pushingSetting"], intent.PushingSetting)
	if err != nil {
		return subscriptionResourceModel{}, err
	}
	workspaceFilterSetting, err := flattenSubscriptionWorkspaceFilterSetting(object["workspaceFilterSetting"], intent.WorkspaceFilterSetting)
	if err != nil {
		return subscriptionResourceModel{}, err
	}
	agentConfig, err := flattenSubscriptionAgentConfig(object["agentConfig"])
	if err != nil {
		return subscriptionResourceModel{}, err
	}

	return subscriptionResourceModel{
		Id:                     types.StringValue(actualId),
		SubscriptionName:       subscriptionJSONStringValue(object["subscriptionName"]),
		Description:            subscriptionStringStateValue(object["description"], intent.Description),
		NotifyStrategyId:       subscriptionStringStateValue(object["notifyStrategyId"], intent.NotifyStrategyId),
		Workspace:              types.StringValue(actualWorkspace),
		Enable:                 subscriptionJSONBoolValue(object["enable"]),
		SubscribeLegacyEvent:   subscriptionBoolStateValue(object["subscribeLegacyEvent"]),
		FilterSetting:          filterSetting,
		PushingSetting:         pushingSetting,
		WorkspaceFilterSetting: workspaceFilterSetting,
		AgentConfig:            agentConfig,
		UserId:                 subscriptionJSONStringValue(object["userId"]),
		SubscriptionType:       subscriptionJSONStringValue(object["subscriptionType"]),
		RegionId:               types.StringValue(regionId),
		CreateTime:             subscriptionJSONStringValue(object["createTime"]),
		UpdateTime:             subscriptionJSONStringValue(object["updateTime"]),
	}, nil
}

func subscriptionObjectValue(raw interface{}, attribute string) (map[string]interface{}, bool, error) {
	if raw == nil {
		return nil, false, nil
	}
	value, ok := raw.(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("the response value of $.%s has unexpected type %T", attribute, raw)
	}
	if len(value) == 0 {
		return nil, false, nil
	}
	return value, true, nil
}

// flattenSubscriptionFilterSetting rebuilds the block from the response only.
// An absent block means the subscription has no filter; the prior state is not
// reused, so Refresh and a bare Import see the same remote configuration.
func flattenSubscriptionFilterSetting(raw interface{}, previous []subscriptionFilterSettingModel) ([]subscriptionFilterSettingModel, error) {
	setting, present, err := subscriptionObjectValue(raw, "filterSetting")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	previousExpression := types.StringNull()
	previousRelation := types.StringNull()
	if len(previous) > 0 {
		previousExpression = previous[0].Expression
		previousRelation = previous[0].Relation
	}
	conditions, err := flattenSubscriptionConditions(setting["conditions"])
	if err != nil {
		return nil, err
	}
	return []subscriptionFilterSettingModel{{
		Expression: subscriptionStringStateValue(setting["expression"], previousExpression),
		Relation:   subscriptionStringStateValue(setting["relation"], previousRelation),
		Conditions: conditions,
	}}, nil
}

func flattenSubscriptionConditions(raw interface{}) ([]subscriptionConditionModel, error) {
	items, ok := subscriptionJSONItems(raw)
	if !ok {
		return nil, fmt.Errorf("the response value of $.filterSetting.conditions has unexpected type %T", raw)
	}
	conditions := make([]subscriptionConditionModel, 0, len(items))
	for index, item := range items {
		condition, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("the response value of $.filterSetting.conditions[%d] has unexpected type %T", index, item)
		}
		conditions = append(conditions, subscriptionConditionModel{
			Field: subscriptionJSONStringValue(condition["field"]),
			Op:    subscriptionJSONStringValue(condition["op"]),
			Value: subscriptionJSONStringValue(condition["value"]),
		})
	}
	return conditions, nil
}

func flattenSubscriptionPushingSetting(raw interface{}, previous []subscriptionPushingSettingModel) ([]subscriptionPushingSettingModel, error) {
	setting, present, err := subscriptionObjectValue(raw, "pushingSetting")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	previousAlert := types.ListNull(types.StringType)
	previousRestore := types.ListNull(types.StringType)
	previousTemplate := types.StringNull()
	previousResponsePlan := types.StringNull()
	if len(previous) > 0 {
		previousAlert = previous[0].AlertActionIds
		previousRestore = previous[0].RestoreActionIds
		previousTemplate = previous[0].TemplateUuid
		previousResponsePlan = previous[0].ResponsePlanId
	}
	return []subscriptionPushingSettingModel{{
		AlertActionIds:   subscriptionJSONStringListState(setting["alertActionIds"], previousAlert),
		RestoreActionIds: subscriptionJSONStringListState(setting["restoreActionIds"], previousRestore),
		TemplateUuid:     subscriptionStringStateValue(setting["templateUuid"], previousTemplate),
		ResponsePlanId:   subscriptionStringStateValue(setting["responsePlanId"], previousResponsePlan),
	}}, nil
}

func flattenSubscriptionWorkspaceFilterSetting(raw interface{}, previous []subscriptionWorkspaceFilterModel) ([]subscriptionWorkspaceFilterModel, error) {
	setting, present, err := subscriptionObjectValue(raw, "workspaceFilterSetting")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	previousUuids := types.ListNull(types.StringType)
	if len(previous) > 0 {
		previousUuids = previous[0].WorkspaceUuids
	}
	return []subscriptionWorkspaceFilterModel{{
		WorkspaceUuids: subscriptionJSONStringListState(setting["workspaceUuids"], previousUuids),
	}}, nil
}

func flattenSubscriptionAgentConfig(raw interface{}) ([]subscriptionAgentConfigModel, error) {
	config, present, err := subscriptionObjectValue(raw, "agentConfig")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	routes, err := flattenSubscriptionAgentRoutes(config["routes"])
	if err != nil {
		return nil, err
	}
	return []subscriptionAgentConfigModel{{
		AgentUuid: subscriptionJSONStringValue(config["agentUuid"]),
		Routes:    routes,
	}}, nil
}

func flattenSubscriptionAgentRoutes(raw interface{}) ([]subscriptionAgentRouteModel, error) {
	items, ok := subscriptionJSONItems(raw)
	if !ok {
		return nil, fmt.Errorf("the response value of $.agentConfig.routes has unexpected type %T", raw)
	}
	routes := make([]subscriptionAgentRouteModel, 0, len(items))
	for index, item := range items {
		route, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("the response value of $.agentConfig.routes[%d] has unexpected type %T", index, item)
		}
		channels, err := flattenSubscriptionAgentChannels(route["channels"])
		if err != nil {
			return nil, err
		}
		effectTimeRange, err := flattenSubscriptionAgentEffectTimeRange(route["effectTimeRange"])
		if err != nil {
			return nil, err
		}
		routes = append(routes, subscriptionAgentRouteModel{
			Channels:        channels,
			EffectTimeRange: effectTimeRange,
		})
	}
	return routes, nil
}

func flattenSubscriptionAgentChannels(raw interface{}) ([]subscriptionAgentChannelModel, error) {
	items, ok := subscriptionJSONItems(raw)
	if !ok {
		return nil, fmt.Errorf("the response value of $.agentConfig.routes.channels has unexpected type %T", raw)
	}
	channels := make([]subscriptionAgentChannelModel, 0, len(items))
	for index, item := range items {
		channel, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("the response value of $.agentConfig.routes.channels[%d] has unexpected type %T", index, item)
		}
		channels = append(channels, subscriptionAgentChannelModel{
			ChannelType:        subscriptionJSONStringValue(channel["channelType"]),
			Receivers:          subscriptionJSONStringList(channel["receivers"]),
			EnabledSubChannels: subscriptionJSONStringList(channel["enabledSubChannels"]),
		})
	}
	return channels, nil
}

func flattenSubscriptionAgentEffectTimeRange(raw interface{}) ([]subscriptionAgentEffectTimeRangeModel, error) {
	window, present, err := subscriptionObjectValue(raw, "effectTimeRange")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	return []subscriptionAgentEffectTimeRangeModel{{
		TimeZone:          subscriptionJSONStringValue(window["timeZone"]),
		DayInWeek:         subscriptionJSONInt64List(window["dayInWeek"]),
		StartTimeInMinute: subscriptionJSONInt64Value(window["startTimeInMinute"]),
		EndTimeInMinute:   subscriptionJSONInt64Value(window["endTimeInMinute"]),
	}}, nil
}

func subscriptionValueString(v types.String) string {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueString()
	}
	return ""
}

func subscriptionValueStringList(v types.List) []string {
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

func subscriptionValueInt64(v types.Int64) (int64, bool) {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueInt64(), true
	}
	return 0, false
}

func subscriptionValueInt64List(v types.List) []int64 {
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

func subscriptionJSONString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func subscriptionJSONStringValue(v interface{}) types.String {
	if s := subscriptionJSONString(v); s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// subscriptionStringStateValue mirrors the configuration: the service drops a
// cleared string from the response, so a value that was configured (including
// an explicit empty string) reads back as an empty string, while an argument
// that was never configured stays null.
func subscriptionStringStateValue(v interface{}, previous types.String) types.String {
	if s := subscriptionJSONString(v); s != "" {
		return types.StringValue(s)
	}
	if !previous.IsNull() && !previous.IsUnknown() {
		return types.StringValue("")
	}
	return types.StringNull()
}

func subscriptionJSONBoolValue(v interface{}) types.Bool {
	if b, ok := v.(bool); ok {
		return types.BoolValue(b)
	}
	return types.BoolNull()
}

func subscriptionBoolStateValue(v interface{}) types.Bool {
	if b, ok := v.(bool); ok {
		return types.BoolValue(b)
	}
	// The contract states that null and false both mean "not subscribed", so an
	// omitted response materializes the service default instead of a null that
	// would differ between apply and import.
	return types.BoolValue(false)
}

func subscriptionJSONStringList(v interface{}) types.List {
	if v == nil {
		return types.ListNull(types.StringType)
	}
	items, ok := subscriptionJSONItems(v)
	if !ok {
		return types.ListNull(types.StringType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			values = append(values, types.StringValue(s))
		}
	}
	return types.ListValueMust(types.StringType, values)
}

func subscriptionJSONStringListState(v interface{}, previous types.List) types.List {
	if v == nil {
		if !previous.IsNull() && !previous.IsUnknown() {
			return types.ListValueMust(types.StringType, []attr.Value{})
		}
		return types.ListNull(types.StringType)
	}
	return subscriptionJSONStringList(v)
}

func subscriptionJSONInt64List(v interface{}) types.List {
	if v == nil {
		return types.ListNull(types.Int64Type)
	}
	items, ok := subscriptionJSONItems(v)
	if !ok {
		return types.ListNull(types.Int64Type)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		value := subscriptionJSONInt64Value(item)
		if !value.IsNull() {
			values = append(values, value)
		}
	}
	return types.ListValueMust(types.Int64Type, values)
}

func subscriptionJSONInt64Value(v interface{}) types.Int64 {
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

type subscriptionChannelValidator struct{}

func (v subscriptionChannelValidator) Description(ctx context.Context) string {
	return "enabled_sub_channels is required when channel_type is CONTACT, GROUP or DUTY"
}

func (v subscriptionChannelValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v subscriptionChannelValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config subscriptionResourceModel
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
					if channel.EnabledSubChannels.IsUnknown() {
						// The value is only known after apply; validation runs again
						// once it is known.
						continue
					}
					if channel.EnabledSubChannels.IsNull() || len(channel.EnabledSubChannels.Elements()) == 0 {
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
