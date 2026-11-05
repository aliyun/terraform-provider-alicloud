package cloudmonitorservice

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	return []resource.ConfigValidator{subscriptionChannelValidator{}, subscriptionBlockValidator{}}
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
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"notify_strategy_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				// CreateSubscription accepts the binding, but UpdateSubscription
				// rejects a request that adds or changes it (HTTP 400
				// InvalidParameterValue), so the change is applied by creating a
				// replacement subscription with the requested binding.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								subscriptionListElementsValidator{},
							},
						},
						"restore_action_ids": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								subscriptionListElementsValidator{},
							},
						},
						"template_uuid": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"response_plan_id": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								subscriptionListElementsValidator{},
							},
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
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"receivers": schema.ListAttribute{
													Optional:    true,
													ElementType: types.StringType,
													Validators: []validator.List{
														listvalidator.SizeAtLeast(1),
														subscriptionListElementsValidator{},
													},
												},
												"enabled_sub_channels": schema.ListAttribute{
													Optional:    true,
													ElementType: types.StringType,
													Validators: []validator.List{
														listvalidator.SizeAtLeast(1),
														subscriptionListElementsValidator{},
													},
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
													Validators: []validator.String{
														stringvalidator.LengthAtLeast(1),
													},
												},
												"day_in_week": schema.ListAttribute{
													Optional:    true,
													ElementType: types.Int64Type,
													Validators: []validator.List{
														listvalidator.SizeAtLeast(1),
														listvalidator.ValueInt64sAre(int64validator.Between(0, 6)),
														subscriptionListElementsValidator{},
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
	response, err := r.roaCall(ctx, action, request, false, func() (map[string]interface{}, error) {
		// The shared ROA client cannot attach the operation context to an
		// already-issued HTTP request, so its own automatic retry is disabled
		// here: every retry goes through roaCall, which observes the context
		// before issuing the next attempt. CreateSubscription is replayed as
		// well, because the API has no idempotency token: a request that already
		// created the subscription would otherwise create a second, orphaned one.
		return r.Meta().RoaPost("Cms", "2024-03-30", action, subscriptionQuery(plan.Workspace), nil, request, false)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error()+
				subscriptionUnknownWriteSuffix(ctx, err, fmt.Sprintf(" Check the workspace for a subscription named %q before retrying: a request whose response was lost may already have created it.", plan.SubscriptionName.ValueString())),
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
		// The disable operation returns before the change is visible, and a
		// response that only carries the switch can still hold a partially
		// applied body. Wait for the requested configuration and the switch
		// together instead of accepting the first read that reports the switch.
		disabled := false
		object, err = r.waitSubscriptionUpdate(ctx, id, plan.Workspace, request, nil, &disabled)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: waiting for DisableSubscription", subscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, "/subscriptions/"+id, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	state, err := subscriptionState(object, identity, r.Meta().RegionId)
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
	next, err := subscriptionState(object, state, r.Meta().RegionId)
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
	var expectedEnable *bool

	if !plan.Enable.Equal(state.Enable) && !plan.Enable.IsNull() && !plan.Enable.IsUnknown() {
		enable := plan.Enable.ValueBool()
		expectedEnable = &enable
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
	}

	action := fmt.Sprintf("/subscriptions/%s", id)
	request := subscriptionRequestForConfig(plan, state)
	removed := subscriptionRemovedPaths(plan, state)
	if len(removed) == 0 && subscriptionRequestsEqual(expandSubscriptionRequest(state), request) {
		// Only the enable/disable state changed: the body is untouched, so no
		// body update is sent. The state is still rebuilt from a read that
		// confirms the stored configuration together with the requested switch:
		// a partial response that carries only the name and the switch must not
		// overwrite the description and the filter blocks already in state.
		object, err := r.waitSubscriptionUpdate(ctx, id, state.Workspace, expandSubscriptionRequest(state), nil, expectedEnable)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		next, err := subscriptionState(object, state, r.Meta().RegionId)
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

	_, err = r.roaCall(ctx, action, request, true, func() (map[string]interface{}, error) {
		return r.Meta().RoaPut("Cms", "2024-03-30", action, subscriptionQuery(state.Workspace), nil, request, false)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err = r.waitSubscriptionUpdate(ctx, id, state.Workspace, request, removed, expectedEnable)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSubscription", subscriptionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, subscriptionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	next, err := subscriptionState(object, state, r.Meta().RegionId)
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
	_, err := r.roaCall(ctx, action, subscriptionQuery(state.Workspace), true, func() (map[string]interface{}, error) {
		response, callErr := r.Meta().RoaDelete("Cms", "2024-03-30", action, subscriptionQuery(state.Workspace), nil, nil, false)
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

// roaCall runs one Cms write with a bounded retry loop. The loop runs in the
// caller's goroutine instead of retry.RetryContext, whose callback goroutine may
// outlive the cancelled call and race with the caller on the response value.
// The caller's context is checked before every attempt and every backoff, so no
// request is issued after the deadline; an attempt already on the wire still
// finishes and its response, including a created id, is returned to the caller.
// replay reports whether the verb is safe to resend: PUT and DELETE are applied
// at most once for the same payload, while CreateSubscription has no idempotency
// token and is therefore sent exactly once, except for a throttling rejection,
// which proves the service never processed the request.
func (r *subscriptionResource) roaCall(ctx context.Context, action string, logPayload interface{}, replay bool, call func() (map[string]interface{}, error)) (map[string]interface{}, error) {
	wait := subscriptionBackoff(ctx)
	for {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		response, err := call()
		if err == nil {
			logs.AddDebug(action, response, logPayload)
			return response, nil
		}
		if !subscriptionRetryable(err) || (!replay && !subscriptionThrottled(err)) {
			logs.AddDebug(action, response, logPayload)
			return response, err
		}
		if !wait() {
			return nil, ctx.Err()
		}
	}
}

// subscriptionRetryable reports whether one attempt may be retried. Transport
// failures carry no service status and are retried for every verb, matching the
// retry scope the SDK client applies when its automatic retry is enabled; the
// service errors retried by the provider keep their existing classification.
func subscriptionRetryable(err error) bool {
	if err == nil {
		return false
	}
	var transportErr *url.Error
	if stderrors.As(err, &transportErr) {
		return true
	}
	return helper.NeedRetry(err)
}

// subscriptionUnknownWriteSuffix marks a write whose outcome cannot be known:
// the shared connectivity client issues requests without a context, so a
// request already on the wire when the operation stopped may still complete,
// and a transport failure or a server-side error can also follow an accepted
// request whose response was lost. A definitive client-side rejection is not
// reported as unknown.
func subscriptionUnknownWriteSuffix(ctx context.Context, err error, detail string) string {
	if !subscriptionWriteOutcomeUnknown(ctx, err) {
		return ""
	}
	return " The server-side result of the write is unknown." + detail
}

func subscriptionWriteOutcomeUnknown(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return true
	}
	var sdkErr *tea.SDKError
	if stderrors.As(err, &sdkErr) {
		// A throttling rejection answers before the service processes the
		// request, so it proves the write was not applied.
		if subscriptionThrottled(err) {
			return false
		}
		// A 5xx answer does not prove the request was rejected.
		return tea.IntValue(sdkErr.StatusCode) >= http.StatusInternalServerError
	}
	var transportErr *url.Error
	return stderrors.As(err, &transportErr)
}

// subscriptionThrottled reports whether the service rejected the request for
// rate limiting, which is the one transient failure that is safe to resend for
// a non-replayable verb.
func subscriptionThrottled(err error) bool {
	var sdkErr *tea.SDKError
	if !stderrors.As(err, &sdkErr) {
		return false
	}
	if tea.IntValue(sdkErr.StatusCode) == http.StatusTooManyRequests {
		return true
	}
	code := tea.StringValue(sdkErr.Code)
	return code == "Rejected.Throttling" || strings.Contains(code, "Throttling")
}

func (r *subscriptionResource) setSubscriptionEnable(ctx context.Context, id string, workspace types.String, enable bool) error {
	suffix := "enable"
	if !enable {
		suffix = "disable"
	}
	action := fmt.Sprintf("/subscriptions/%s/%s", id, suffix)

	_, err := r.roaCall(ctx, action, action, true, func() (map[string]interface{}, error) {
		return r.Meta().RoaPut("Cms", "2024-03-30", action, subscriptionQuery(workspace), nil, nil, false)
	})
	return err
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
// to UpdateSubscription is visible on the read path, every removed key is
// gone and, when the same update also toggled the subscription, the enabled
// switch reports the requested value. The Cms event subscription reads are
// eventually consistent, so an immediate read can return the values from
// before the update.
func (r *subscriptionResource) waitSubscriptionUpdate(ctx context.Context, id string, workspace types.String, request map[string]interface{}, removed []subscriptionPath, enable *bool) (map[string]interface{}, error) {
	for {
		object, err := NewCmsService(r.Meta()).DescribeSubscription(ctx, id, workspace.ValueString())
		if err == nil {
			if cleared, clearErr := subscriptionPathsCleared(object, removed); clearErr != nil {
				return nil, clearErr
			} else if cleared && subscriptionRequestVisible(object, request) && subscriptionEnableMatches(object, enable) {
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

// subscriptionEnableMatches reports whether the read object carries the
// requested enabled state. A missing switch is not a match: the combined
// update must not be accepted on a response that does not confirm it.
func subscriptionEnableMatches(object map[string]interface{}, enable *bool) bool {
	if enable == nil {
		return true
	}
	current, ok := object["enable"].(bool)
	return ok && current == *enable
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

// subscriptionIntent resolves the plan for the request comparison: the service
// echoes a relation even when the configuration omits it, so an unknown
// computed relation falls back to the value already in state and does not look
// like a body change.
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
	expression := ""
	if !block.Expression.IsNull() && !block.Expression.IsUnknown() {
		expression = block.Expression.ValueString()
		if expression != "" {
			result["expression"] = expression
		}
	}
	// The service ignores relation when expression is set, so it is not sent in
	// that case: the value the service stores is not the configured one, and a
	// request that demanded it back could never converge. The read path still
	// reports the stored relation, which is why the argument is computed.
	if expression == "" && !block.Relation.IsNull() && !block.Relation.IsUnknown() {
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
// silent overwrite. Every block and optional value is rebuilt from the response
// only, so Refresh and a bare Import produce the same state for the same
// response.
func subscriptionState(object map[string]interface{}, identity subscriptionResourceModel, regionId string) (subscriptionResourceModel, error) {
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

	filterSetting, err := flattenSubscriptionFilterSetting(object["filterSetting"])
	if err != nil {
		return subscriptionResourceModel{}, err
	}
	pushingSetting, err := flattenSubscriptionPushingSetting(object["pushingSetting"])
	if err != nil {
		return subscriptionResourceModel{}, err
	}
	workspaceFilterSetting, err := flattenSubscriptionWorkspaceFilterSetting(object["workspaceFilterSetting"])
	if err != nil {
		return subscriptionResourceModel{}, err
	}
	agentConfig, err := flattenSubscriptionAgentConfig(object["agentConfig"])
	if err != nil {
		return subscriptionResourceModel{}, err
	}

	return subscriptionResourceModel{
		Id:                     types.StringValue(actualId),
		SubscriptionName:       subscriptionStringStateValue(object["subscriptionName"]),
		Description:            subscriptionStringStateValue(object["description"]),
		NotifyStrategyId:       subscriptionStringStateValue(object["notifyStrategyId"]),
		Workspace:              types.StringValue(actualWorkspace),
		Enable:                 subscriptionJSONBoolValue(object["enable"]),
		SubscribeLegacyEvent:   subscriptionBoolStateValue(object["subscribeLegacyEvent"]),
		FilterSetting:          filterSetting,
		PushingSetting:         pushingSetting,
		WorkspaceFilterSetting: workspaceFilterSetting,
		AgentConfig:            agentConfig,
		UserId:                 subscriptionStringStateValue(object["userId"]),
		SubscriptionType:       subscriptionStringStateValue(object["subscriptionType"]),
		RegionId:               types.StringValue(regionId),
		CreateTime:             subscriptionStringStateValue(object["createTime"]),
		UpdateTime:             subscriptionStringStateValue(object["updateTime"]),
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
func flattenSubscriptionFilterSetting(raw interface{}) ([]subscriptionFilterSettingModel, error) {
	setting, present, err := subscriptionObjectValue(raw, "filterSetting")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	conditions, err := flattenSubscriptionConditions(setting["conditions"])
	if err != nil {
		return nil, err
	}
	block := subscriptionFilterSettingModel{
		Expression: subscriptionStringStateValue(setting["expression"]),
		Relation:   subscriptionStringStateValue(setting["relation"]),
		Conditions: conditions,
	}
	// A block whose members are all empty carries no configuration: keeping it
	// would show a block in state that the configuration does not declare.
	if block.Expression.IsNull() && block.Relation.IsNull() && len(block.Conditions) == 0 {
		return nil, nil
	}
	return []subscriptionFilterSettingModel{block}, nil
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
			Field: subscriptionStringStateValue(condition["field"]),
			Op:    subscriptionStringStateValue(condition["op"]),
			Value: subscriptionStringStateValue(condition["value"]),
		})
	}
	return conditions, nil
}

func flattenSubscriptionPushingSetting(raw interface{}) ([]subscriptionPushingSettingModel, error) {
	setting, present, err := subscriptionObjectValue(raw, "pushingSetting")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	block := subscriptionPushingSettingModel{
		AlertActionIds:   subscriptionJSONStringList(setting["alertActionIds"]),
		RestoreActionIds: subscriptionJSONStringList(setting["restoreActionIds"]),
		TemplateUuid:     subscriptionStringStateValue(setting["templateUuid"]),
		ResponsePlanId:   subscriptionStringStateValue(setting["responsePlanId"]),
	}
	if block.AlertActionIds.IsNull() && block.RestoreActionIds.IsNull() &&
		block.TemplateUuid.IsNull() && block.ResponsePlanId.IsNull() {
		return nil, nil
	}
	return []subscriptionPushingSettingModel{block}, nil
}

func flattenSubscriptionWorkspaceFilterSetting(raw interface{}) ([]subscriptionWorkspaceFilterModel, error) {
	setting, present, err := subscriptionObjectValue(raw, "workspaceFilterSetting")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	block := subscriptionWorkspaceFilterModel{
		WorkspaceUuids: subscriptionJSONStringList(setting["workspaceUuids"]),
	}
	if block.WorkspaceUuids.IsNull() {
		return nil, nil
	}
	return []subscriptionWorkspaceFilterModel{block}, nil
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
	block := subscriptionAgentConfigModel{
		AgentUuid: subscriptionStringStateValue(config["agentUuid"]),
		Routes:    routes,
	}
	if block.AgentUuid.IsNull() && len(block.Routes) == 0 {
		return nil, nil
	}
	return []subscriptionAgentConfigModel{block}, nil
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
		if len(channels) == 0 && len(effectTimeRange) == 0 {
			// A route with no usable member is not a configured block.
			continue
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
		block := subscriptionAgentChannelModel{
			ChannelType:        subscriptionStringStateValue(channel["channelType"]),
			Receivers:          subscriptionJSONStringList(channel["receivers"]),
			EnabledSubChannels: subscriptionJSONStringList(channel["enabledSubChannels"]),
		}
		if block.ChannelType.IsNull() && block.Receivers.IsNull() && block.EnabledSubChannels.IsNull() {
			continue
		}
		channels = append(channels, block)
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
	block := subscriptionAgentEffectTimeRangeModel{
		TimeZone:          subscriptionStringStateValue(window["timeZone"]),
		DayInWeek:         subscriptionJSONInt64List(window["dayInWeek"]),
		StartTimeInMinute: subscriptionJSONInt64Value(window["startTimeInMinute"]),
		EndTimeInMinute:   subscriptionJSONInt64Value(window["endTimeInMinute"]),
	}
	if block.TimeZone.IsNull() && block.DayInWeek.IsNull() &&
		block.StartTimeInMinute.IsNull() && block.EndTimeInMinute.IsNull() {
		return nil, nil
	}
	return []subscriptionAgentEffectTimeRangeModel{block}, nil
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

// subscriptionStringStateValue defines one representation for every string:
// the service drops a cleared value from the response, and an explicit empty
// string is rejected by validation, so an absent or empty member always reads
// back as null. Refresh and a bare Import therefore produce the same state for
// the same GET response.
func subscriptionStringStateValue(v interface{}) types.String {
	if s := subscriptionJSONString(v); s != "" {
		return types.StringValue(s)
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

// subscriptionJSONStringList defines one representation for every optional
// list: an omitted, empty or unusable member reads back as null, including an
// empty string member, which the configuration cannot express. The response
// alone decides the value, never the previous state, so Refresh and a bare
// Import of the same response produce the same state and a configuration that
// omits the argument plans no change.
func subscriptionJSONStringList(v interface{}) types.List {
	items, ok := subscriptionJSONItems(v)
	if !ok {
		return types.ListNull(types.StringType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok && s != "" {
			values = append(values, types.StringValue(s))
		}
	}
	if len(values) == 0 {
		return types.ListNull(types.StringType)
	}
	return types.ListValueMust(types.StringType, values)
}

func subscriptionJSONInt64List(v interface{}) types.List {
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
	if len(values) == 0 {
		return types.ListNull(types.Int64Type)
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

// subscriptionBlockValidator rejects blocks that declare no member, and rejects
// a filter_setting that sets both expression and relation. A block with no
// configured member is omitted from the request, dropped from the read response
// and therefore cannot round-trip through apply. A member that is only known
// after apply counts as declared, and an unknown nested collection cannot be
// judged before apply, so neither is rejected here. The expression/relation
// conflict is only judged when both values are known: a dynamic configuration
// resolves them at apply time, so an unknown pair is left to the service.
type subscriptionBlockValidator struct{}

func (v subscriptionBlockValidator) Description(ctx context.Context) string {
	return "a declared block must configure at least one member"
}

func (v subscriptionBlockValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v subscriptionBlockValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	for i, element := range subscriptionConfigBlockElements(ctx, req.Config, path.Root("filter_setting")) {
		attributes, ok := subscriptionConfigObjectAttributes(element)
		if !ok {
			continue
		}
		if subscriptionMemberKnown(attributes, "expression") && subscriptionMemberKnown(attributes, "relation") {
			// The service ignores relation when expression is set, so a state
			// that carried both could never match the read-back value. The check
			// waits until both values are known: a configuration that resolves
			// them at apply time is not a plan-time conflict.
			resp.Diagnostics.AddAttributeError(
				path.Root("filter_setting").AtListIndex(i).AtName("relation"),
				"Conflicting filter_setting members",
				"relation must not be set when expression is set: the service ignores relation in that case.",
			)
			continue
		}
		if subscriptionMemberDeclared(attributes, "expression") || subscriptionMemberDeclared(attributes, "relation") ||
			subscriptionNestedBlockDeclared(attributes, "conditions") {
			continue
		}
		resp.Diagnostics.AddAttributeError(
			path.Root("filter_setting").AtListIndex(i),
			"Empty filter_setting block",
			"filter_setting must configure expression, relation or at least one conditions block.",
		)
	}

	for i, element := range subscriptionConfigBlockElements(ctx, req.Config, path.Root("pushing_setting")) {
		attributes, ok := subscriptionConfigObjectAttributes(element)
		if !ok {
			continue
		}
		if subscriptionMemberDeclared(attributes, "alert_action_ids") || subscriptionMemberDeclared(attributes, "restore_action_ids") ||
			subscriptionMemberDeclared(attributes, "template_uuid") || subscriptionMemberDeclared(attributes, "response_plan_id") {
			continue
		}
		resp.Diagnostics.AddAttributeError(
			path.Root("pushing_setting").AtListIndex(i),
			"Empty pushing_setting block",
			"pushing_setting must configure at least one member: alert_action_ids, restore_action_ids, template_uuid or response_plan_id.",
		)
	}

	for i, element := range subscriptionConfigBlockElements(ctx, req.Config, path.Root("workspace_filter_setting")) {
		attributes, ok := subscriptionConfigObjectAttributes(element)
		if !ok {
			continue
		}
		if subscriptionMemberDeclared(attributes, "workspace_uuids") {
			continue
		}
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace_filter_setting").AtListIndex(i),
			"Empty workspace_filter_setting block",
			"workspace_filter_setting must configure workspace_uuids.",
		)
	}

	for i, element := range subscriptionConfigBlockElements(ctx, req.Config, path.Root("agent_config")) {
		attributes, ok := subscriptionConfigObjectAttributes(element)
		if !ok {
			continue
		}
		agentPath := path.Root("agent_config").AtListIndex(i)
		if !subscriptionMemberDeclared(attributes, "agent_uuid") && !subscriptionNestedBlockDeclared(attributes, "routes") {
			resp.Diagnostics.AddAttributeError(
				agentPath,
				"Empty agent_config block",
				"agent_config must configure agent_uuid or at least one routes block.",
			)
		}
		for j, routeElement := range subscriptionConfigNestedBlockElements(attributes, "routes") {
			routeAttributes, ok := subscriptionConfigObjectAttributes(routeElement)
			if !ok {
				continue
			}
			routePath := agentPath.AtName("routes").AtListIndex(j)
			if !subscriptionNestedBlockDeclared(routeAttributes, "channels") && !subscriptionNestedBlockDeclared(routeAttributes, "effect_time_range") {
				resp.Diagnostics.AddAttributeError(routePath, "Empty routes block", "routes must configure at least one channels or effect_time_range block.")
			}
			for k, channelElement := range subscriptionConfigNestedBlockElements(routeAttributes, "channels") {
				channelAttributes, ok := subscriptionConfigObjectAttributes(channelElement)
				if !ok {
					continue
				}
				if subscriptionMemberDeclared(channelAttributes, "channel_type") || subscriptionMemberDeclared(channelAttributes, "receivers") ||
					subscriptionMemberDeclared(channelAttributes, "enabled_sub_channels") {
					continue
				}
				resp.Diagnostics.AddAttributeError(routePath.AtName("channels").AtListIndex(k), "Empty channels block", "channels must configure channel_type, receivers or enabled_sub_channels.")
			}
			for k, windowElement := range subscriptionConfigNestedBlockElements(routeAttributes, "effect_time_range") {
				windowAttributes, ok := subscriptionConfigObjectAttributes(windowElement)
				if !ok {
					continue
				}
				if subscriptionMemberDeclared(windowAttributes, "time_zone") || subscriptionMemberDeclared(windowAttributes, "day_in_week") ||
					subscriptionMemberDeclared(windowAttributes, "start_time_in_minute") || subscriptionMemberDeclared(windowAttributes, "end_time_in_minute") {
					continue
				}
				resp.Diagnostics.AddAttributeError(routePath.AtName("effect_time_range").AtListIndex(k), "Empty effect_time_range block", "effect_time_range must configure at least one member.")
			}
		}
	}
}

// subscriptionConfigBlockElements reads one nested block list from the
// configuration. A null or unknown list is not judged: the block is either
// absent or its value is only known once the plan is applied.
func subscriptionConfigBlockElements(ctx context.Context, config tfsdk.Config, blockPath path.Path) []attr.Value {
	var blocks types.List
	if diags := config.GetAttribute(ctx, blockPath, &blocks); diags.HasError() {
		return nil
	}
	return subscriptionKnownListElements(blocks)
}

func subscriptionConfigNestedBlockElements(attributes map[string]attr.Value, name string) []attr.Value {
	value, ok := attributes[name]
	if !ok {
		return nil
	}
	list, ok := value.(types.List)
	if !ok {
		return nil
	}
	return subscriptionKnownListElements(list)
}

func subscriptionKnownListElements(list types.List) []attr.Value {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	return list.Elements()
}

// subscriptionConfigObjectAttributes returns the members of one known block
// element. A null or unknown element cannot be judged before apply and is
// skipped.
func subscriptionConfigObjectAttributes(element attr.Value) (map[string]attr.Value, bool) {
	object, ok := element.(types.Object)
	if !ok || object.IsNull() || object.IsUnknown() {
		return nil, false
	}
	return object.Attributes(), true
}

// subscriptionMemberDeclared reports whether an attribute was declared.
// Unknown counts as declared: a dynamic value is resolved before apply and must
// not be mistaken for an empty block.
func subscriptionMemberDeclared(attributes map[string]attr.Value, name string) bool {
	value, ok := attributes[name]
	return ok && !value.IsNull()
}

// subscriptionMemberKnown reports whether an attribute carries a value that is
// already known at plan time. An unknown value can still resolve to null.
func subscriptionMemberKnown(attributes map[string]attr.Value, name string) bool {
	value, ok := attributes[name]
	return ok && !value.IsNull() && !value.IsUnknown()
}

// subscriptionNestedBlockDeclared reports whether a nested block list is
// unknown or carries at least one element.
func subscriptionNestedBlockDeclared(attributes map[string]attr.Value, name string) bool {
	value, ok := attributes[name]
	if !ok || value.IsNull() {
		return false
	}
	if value.IsUnknown() {
		return true
	}
	list, ok := value.(types.List)
	return ok && len(list.Elements()) > 0
}

// subscriptionListElementsValidator rejects a null or empty member. A null
// member satisfies the size validator but is dropped when the request is built,
// so the value on the wire would differ from the configuration and the read
// path could not round-trip it. An unknown member is resolved before apply and
// is accepted here.
type subscriptionListElementsValidator struct{}

func (v subscriptionListElementsValidator) Description(ctx context.Context) string {
	return "every list member must be a known, non-empty value"
}

func (v subscriptionListElementsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v subscriptionListElementsValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for index, element := range req.ConfigValue.Elements() {
		if element.IsUnknown() {
			continue
		}
		if element.IsNull() {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(index),
				"Null list member",
				"The list members must not be null: remove the member or provide a value.",
			)
			continue
		}
		if text, ok := element.(types.String); ok && text.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(index),
				"Empty list member",
				"The list members must not be empty strings: remove the member to unset it.",
			)
		}
	}
}

type subscriptionChannelValidator struct{}

func (v subscriptionChannelValidator) Description(ctx context.Context) string {
	return "enabled_sub_channels is required when channel_type is CONTACT, GROUP or DUTY"
}

func (v subscriptionChannelValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v subscriptionChannelValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	for i, agentElement := range subscriptionConfigBlockElements(ctx, req.Config, path.Root("agent_config")) {
		agentAttributes, ok := subscriptionConfigObjectAttributes(agentElement)
		if !ok {
			continue
		}
		for j, routeElement := range subscriptionConfigNestedBlockElements(agentAttributes, "routes") {
			routeAttributes, ok := subscriptionConfigObjectAttributes(routeElement)
			if !ok {
				continue
			}
			for k, channelElement := range subscriptionConfigNestedBlockElements(routeAttributes, "channels") {
				channelAttributes, ok := subscriptionConfigObjectAttributes(channelElement)
				if !ok {
					continue
				}
				channelType, ok := channelAttributes["channel_type"].(types.String)
				if !ok || channelType.IsNull() || channelType.IsUnknown() {
					continue
				}
				switch channelType.ValueString() {
				case "CONTACT", "GROUP", "DUTY":
					enabledSubChannels, declared := channelAttributes["enabled_sub_channels"].(types.List)
					if declared && enabledSubChannels.IsUnknown() {
						// The value is only known after apply; validation runs again
						// once it is known.
						continue
					}
					if !declared || enabledSubChannels.IsNull() || len(enabledSubChannels.Elements()) == 0 {
						resp.Diagnostics.AddAttributeError(
							path.Root("agent_config").AtListIndex(i).
								AtName("routes").AtListIndex(j).
								AtName("channels").AtListIndex(k).
								AtName("enabled_sub_channels"),
							"Missing enabled_sub_channels",
							fmt.Sprintf("enabled_sub_channels must contain at least one value when channel_type is %s.", channelType.ValueString()),
						)
					}
				}
			}
		}
	}
}
