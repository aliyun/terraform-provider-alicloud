package ens

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	ensLoadBalancerHttpsListenerTypeName = "alicloud_ens_load_balancer_https_listener"
	// terraform-plugin-framework does not support user-configurable resource
	// timeouts; these mirror the SDKv2 sibling listener defaults.
	ensLoadBalancerHttpsListenerCreateTimeout = 15 * time.Minute
	ensLoadBalancerHttpsListenerUpdateTimeout = 5 * time.Minute
	ensLoadBalancerHttpsListenerDeleteTimeout = 15 * time.Minute
)

var (
	_ resource.Resource                = &ensLoadBalancerHttpsListenerResource{}
	_ resource.ResourceWithImportState = &ensLoadBalancerHttpsListenerResource{}
)

func NewEnsLoadBalancerHttpsListenerResource() resource.Resource {
	return &ensLoadBalancerHttpsListenerResource{}
}

type ensLoadBalancerHttpsListenerResource struct {
	fwadapt.ResourceBase
}

type ensLoadBalancerHttpsListenerResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	LoadBalancerId         types.String `tfsdk:"load_balancer_id"`
	ListenerPort           types.Int64  `tfsdk:"listener_port"`
	ServerCertificateId    types.String `tfsdk:"server_certificate_id"`
	BackendServerPort      types.Int64  `tfsdk:"backend_server_port"`
	Description            types.String `tfsdk:"description"`
	Scheduler              types.String `tfsdk:"scheduler"`
	HealthCheck            types.String `tfsdk:"health_check"`
	HealthCheckDomain      types.String `tfsdk:"health_check_domain"`
	HealthCheckURI         types.String `tfsdk:"health_check_uri"`
	HealthyThreshold       types.Int64  `tfsdk:"healthy_threshold"`
	UnhealthyThreshold     types.Int64  `tfsdk:"unhealthy_threshold"`
	HealthCheckTimeout     types.Int64  `tfsdk:"health_check_timeout"`
	HealthCheckConnectPort types.Int64  `tfsdk:"health_check_connect_port"`
	HealthCheckInterval    types.Int64  `tfsdk:"health_check_interval"`
	HealthCheckHttpCode    types.String `tfsdk:"health_check_http_code"`
	HealthCheckMethod      types.String `tfsdk:"health_check_method"`
	IdleTimeout            types.Int64  `tfsdk:"idle_timeout"`
	RequestTimeout         types.Int64  `tfsdk:"request_timeout"`
	ListenerForward        types.String `tfsdk:"listener_forward"`
	ForwardPort            types.Int64  `tfsdk:"forward_port"`
	StickySessionType      types.String `tfsdk:"sticky_session_type"`
	CookieTimeout          types.Int64  `tfsdk:"cookie_timeout"`
	Cookie                 types.String `tfsdk:"cookie"`
	Status                 types.String `tfsdk:"status"`
}

func (r *ensLoadBalancerHttpsListenerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ensLoadBalancerHttpsListenerTypeName
}

func (r *ensLoadBalancerHttpsListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_port": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"server_certificate_id": schema.StringAttribute{
				Required: true,
			},
			"backend_server_port": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"scheduler": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("wrr"),
				Validators: []validator.String{
					stringvalidator.OneOf("wrr", "wlc", "rr", "sch", "qch", "iqch"),
				},
			},
			"health_check": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("off"),
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
			},
			// When health_check is off, the ENS Describe API omits the health
			// check sub-attributes, so they must neither be planned into the diff
			// nor set into state. Each carries a plan modifier that suppresses the
			// field while health_check is off, mirroring the SDKv2 sibling's
			// DiffSuppressFunc.
			"health_check_domain": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier{},
				},
			},
			"health_check_uri": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier{},
				},
			},
			"healthy_threshold": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(3),
				Validators: []validator.Int64{
					int64validator.Between(2, 10),
				},
				PlanModifiers: []planmodifier.Int64{
					ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier{},
				},
			},
			"unhealthy_threshold": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(3),
				Validators: []validator.Int64{
					int64validator.Between(2, 10),
				},
				PlanModifiers: []planmodifier.Int64{
					ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier{},
				},
			},
			"health_check_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(5),
				Validators: []validator.Int64{
					int64validator.Between(1, 300),
				},
				PlanModifiers: []planmodifier.Int64{
					ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier{},
				},
			},
			"health_check_connect_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{
					ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier{},
				},
			},
			"health_check_interval": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(2),
				Validators: []validator.Int64{
					int64validator.Between(1, 50),
				},
				PlanModifiers: []planmodifier.Int64{
					ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier{},
				},
			},
			"health_check_http_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("http_2xx"),
				Validators: []validator.String{
					stringvalidator.OneOf("http_2xx", "http_3xx", "http_4xx", "http_5xx"),
				},
				PlanModifiers: []planmodifier.String{
					ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier{},
				},
			},
			"health_check_method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("head"),
				Validators: []validator.String{
					stringvalidator.OneOf("head", "get"),
				},
				PlanModifiers: []planmodifier.String{
					ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier{},
				},
			},
			"idle_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(15),
				Validators: []validator.Int64{
					int64validator.Between(1, 60),
				},
			},
			"request_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(60),
				Validators: []validator.Int64{
					int64validator.Between(1, 180),
				},
			},
			"listener_forward": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("off"),
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"forward_port": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// The sticky session attributes are write-only: the ENS Describe API
			// does not return them, so Read preserves them from prior state to
			// avoid spurious replacements. They are replacement-only.
			"sticky_session_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("insert", "server"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cookie_timeout": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 86400),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"cookie": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("running", "stopped"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ensLoadBalancerHttpsListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ensLoadBalancerHttpsListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ensLoadBalancerHttpsListenerCreateRequest(plan)
	action := "CreateLoadBalancerHTTPSListener"
	var response map[string]interface{}
	wait := helper.IncrementalWait(10*time.Second, 15*time.Second)
	err := retry.RetryContext(ctx, ensLoadBalancerHttpsListenerCreateTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("Ens", "2017-11-10", action, nil, request, true)
		if callErr != nil {
			// A same-port ForceNew recreate may leave the listener port still
			// releasing on the ENS backend; absorb the consistency window.
			if errs.IsExpectedErrors(callErr, []string{"ListenerAlreadyExists"}) {
				wait()
				return retry.RetryableError(callErr)
			}
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLoadBalancerHTTPSListener", ensLoadBalancerHttpsListenerTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ensLoadBalancerHttpsListenerTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s:%d", plan.LoadBalancerId.ValueString(), plan.ListenerPort.ValueInt64()))

	ensService := NewService(r.Meta())
	stateConf := helper.BuildStateConf([]string{}, []string{"stopped"}, ensLoadBalancerHttpsListenerCreateTimeout, 5*time.Second, ensService.EnsLoadBalancerHttpsListenerStateRefreshFunc(plan.Id.ValueString(), "Status", []string{}))
	if _, waitErr := stateConf.WaitForState(); waitErr != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.Id)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for stopped", ensLoadBalancerHttpsListenerTypeName),
			errs.WrapErrorf(waitErr, errs.IdMsg, plan.Id.ValueString()).Error(),
		)
		return
	}

	// Start the listener when the desired status is running.
	if plan.Status.ValueString() == "running" {
		if convergeErr := ensLoadBalancerHttpsListenerConvergeStatus(ctx, ensService, plan.Id.ValueString(), "running", ensLoadBalancerHttpsListenerCreateTimeout); convergeErr != nil {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.Id)...)
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: starting the listener", ensLoadBalancerHttpsListenerTypeName),
				convergeErr.Error(),
			)
			return
		}
	}

	object, describeErr := ensService.DescribeEnsLoadBalancerHttpsListener(plan.Id.ValueString())
	if describeErr != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.Id)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling DescribeLoadBalancerHTTPSListenerAttribute", ensLoadBalancerHttpsListenerTypeName),
			errs.WrapErrorf(describeErr, errs.DefaultErrorMsg, plan.Id.ValueString(), "DescribeLoadBalancerHTTPSListenerAttribute", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mapErr := ensLoadBalancerHttpsListenerState(object)
	if mapErr != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.Id)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping Describe response", ensLoadBalancerHttpsListenerTypeName),
			mapErr.Error(),
		)
		return
	}
	ensLoadBalancerHttpsListenerPreserveWriteOnly(&model, plan)
	model.Id = plan.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *ensLoadBalancerHttpsListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ensLoadBalancerHttpsListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.Id.ValueString() == "" {
		return
	}

	ensService := NewService(r.Meta())
	object, err := ensService.DescribeEnsLoadBalancerHttpsListener(state.Id.ValueString())
	if err != nil {
		if errs.IsProviderNotFound(err) || errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLoadBalancerHTTPSListenerAttribute", ensLoadBalancerHttpsListenerTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeLoadBalancerHTTPSListenerAttribute", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mapErr := ensLoadBalancerHttpsListenerState(object)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping Describe response", ensLoadBalancerHttpsListenerTypeName),
			mapErr.Error(),
		)
		return
	}
	ensLoadBalancerHttpsListenerPreserveWriteOnly(&model, state)
	model.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *ensLoadBalancerHttpsListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ensLoadBalancerHttpsListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ensLoadBalancerHttpsListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ensService := NewService(r.Meta())

	attributeChanged := !plan.Description.Equal(state.Description) ||
		!plan.Scheduler.Equal(state.Scheduler) ||
		!plan.HealthCheck.Equal(state.HealthCheck) ||
		!plan.HealthCheckDomain.Equal(state.HealthCheckDomain) ||
		!plan.HealthCheckURI.Equal(state.HealthCheckURI) ||
		!plan.HealthyThreshold.Equal(state.HealthyThreshold) ||
		!plan.UnhealthyThreshold.Equal(state.UnhealthyThreshold) ||
		!plan.HealthCheckTimeout.Equal(state.HealthCheckTimeout) ||
		!plan.HealthCheckConnectPort.Equal(state.HealthCheckConnectPort) ||
		!plan.HealthCheckInterval.Equal(state.HealthCheckInterval) ||
		!plan.HealthCheckHttpCode.Equal(state.HealthCheckHttpCode) ||
		!plan.HealthCheckMethod.Equal(state.HealthCheckMethod) ||
		!plan.IdleTimeout.Equal(state.IdleTimeout) ||
		!plan.RequestTimeout.Equal(state.RequestTimeout) ||
		!plan.ServerCertificateId.Equal(state.ServerCertificateId)

	if attributeChanged {
		request := ensLoadBalancerHttpsListenerSetRequest(plan)
		action := "SetLoadBalancerHTTPSListenerAttribute"
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.RetryContext(ctx, ensLoadBalancerHttpsListenerUpdateTimeout, func() *retry.RetryError {
			var callErr error
			response, callErr = r.Meta().RpcPost("Ens", "2017-11-10", action, nil, request, true)
			if callErr != nil {
				// The Set API requires the listener to be stopped. When it is
				// running, stop it and retry the modification.
				if errs.IsExpectedErrors(callErr, []string{"IncorrectListenerStatus", "IncorrectInstanceStatus"}) {
					_ = ensService.StopEnsLoadBalancerHttpsListener(plan.Id.ValueString())
					wait()
					return retry.RetryableError(callErr)
				}
				if helper.NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling SetLoadBalancerHTTPSListenerAttribute", ensLoadBalancerHttpsListenerTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	// Converge the listener status to the desired value after an explicit
	// status change or an attribute change (the Set API leaves it stopped).
	if !plan.Status.Equal(state.Status) || attributeChanged {
		target := "stopped"
		if plan.Status.ValueString() == "running" {
			target = "running"
		}
		if convergeErr := ensLoadBalancerHttpsListenerConvergeStatus(ctx, ensService, plan.Id.ValueString(), target, ensLoadBalancerHttpsListenerUpdateTimeout); convergeErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: converging status to %s", ensLoadBalancerHttpsListenerTypeName, target),
				convergeErr.Error(),
			)
			return
		}
	}

	object, err := ensService.DescribeEnsLoadBalancerHttpsListener(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling DescribeLoadBalancerHTTPSListenerAttribute", ensLoadBalancerHttpsListenerTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), "DescribeLoadBalancerHTTPSListenerAttribute", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	model, mapErr := ensLoadBalancerHttpsListenerState(object)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: mapping Describe response", ensLoadBalancerHttpsListenerTypeName),
			mapErr.Error(),
		)
		return
	}
	ensLoadBalancerHttpsListenerPreserveWriteOnly(&model, plan)
	model.Id = plan.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *ensLoadBalancerHttpsListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ensLoadBalancerHttpsListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.Id.ValueString() == "" {
		return
	}

	ensService := NewService(r.Meta())
	parts := splitEnsLoadBalancerHttpsListenerId(state.Id.ValueString())
	if parts == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing id", ensLoadBalancerHttpsListenerTypeName),
			fmt.Sprintf("invalid Resource Id %s, expected <load_balancer_id>:<listener_port>", state.Id.ValueString()),
		)
		return
	}
	action := "DeleteLoadBalancerListener"
	request := map[string]interface{}{
		"LoadBalancerId":   parts[0],
		"ListenerPort":     parts[1],
		"ListenerProtocol": "https",
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(ctx, ensLoadBalancerHttpsListenerDeleteTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("Ens", "2017-11-10", action, nil, request, true)
		if callErr != nil {
			// A running listener must be stopped before it can be deleted.
			if errs.IsExpectedErrors(callErr, []string{"IncorrectListenerStatus", "IncorrectInstanceStatus"}) {
				_ = ensService.StopEnsLoadBalancerHttpsListener(state.Id.ValueString())
				wait()
				return retry.RetryableError(callErr)
			}
			if errs.IsExpectedErrors(callErr, []string{"ListenerNotFound", "LoadBalancerNotFound"}) {
				return nil
			}
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsProviderNotFound(err) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteLoadBalancerListener", ensLoadBalancerHttpsListenerTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// Wait until the listener is confirmed gone so a same-port ForceNew recreate
	// does not collide with a not-yet-released listener on the ENS backend.
	goneWait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	waitErr := retry.RetryContext(ctx, ensLoadBalancerHttpsListenerDeleteTimeout, func() *retry.RetryError {
		_, e := ensService.DescribeEnsLoadBalancerHttpsListener(state.Id.ValueString())
		if e != nil {
			if errs.IsProviderNotFound(e) || errs.NotFoundError(e) {
				return nil
			}
			if helper.NeedRetry(e) {
				goneWait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		goneWait()
		return retry.RetryableError(fmt.Errorf("LoadBalancerHttpsListener %s still exists after delete", state.Id.ValueString()))
	})
	if waitErr != nil && !errs.IsProviderNotFound(waitErr) && !errs.NotFoundError(waitErr) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for removal", ensLoadBalancerHttpsListenerTypeName),
			errs.WrapErrorf(waitErr, errs.IdMsg, state.Id.ValueString()).Error(),
		)
	}
}

func (r *ensLoadBalancerHttpsListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ensLoadBalancerHttpsListenerConvergeStatus starts or stops the listener and
// blocks on the state refresh until the target status is reached.
func ensLoadBalancerHttpsListenerConvergeStatus(ctx context.Context, ensService *Service, id, target string, timeout time.Duration) error {
	if target == "running" {
		if err := ensService.StartEnsLoadBalancerHttpsListener(id); err != nil {
			return errs.WrapError(err)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"running"}, timeout, 5*time.Second, ensService.EnsLoadBalancerHttpsListenerStateRefreshFunc(id, "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, id)
		}
		return nil
	}
	if err := ensService.StopEnsLoadBalancerHttpsListener(id); err != nil {
		return errs.WrapError(err)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"stopped"}, timeout, 5*time.Second, ensService.EnsLoadBalancerHttpsListenerStateRefreshFunc(id, "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, id)
	}
	return nil
}

// splitEnsLoadBalancerHttpsListenerId splits "<load_balancer_id>:<listener_port>".
func splitEnsLoadBalancerHttpsListenerId(id string) []string {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil
	}
	return parts
}

// ensLoadBalancerHttpsListenerCreateRequest builds the CreateLoadBalancerHTTPSListener
// request from the planned state.
func ensLoadBalancerHttpsListenerCreateRequest(plan ensLoadBalancerHttpsListenerResourceModel) map[string]interface{} {
	request := map[string]interface{}{
		"LoadBalancerId":      plan.LoadBalancerId.ValueString(),
		"ListenerPort":        plan.ListenerPort.ValueInt64(),
		"ServerCertificateId": plan.ServerCertificateId.ValueString(),
		"Scheduler":           plan.Scheduler.ValueString(),
		"HealthCheck":         plan.HealthCheck.ValueString(),
		"IdleTimeout":         plan.IdleTimeout.ValueInt64(),
		"RequestTimeout":      plan.RequestTimeout.ValueInt64(),
		"ListenerForward":     plan.ListenerForward.ValueString(),
	}
	if !plan.BackendServerPort.IsNull() && !plan.BackendServerPort.IsUnknown() {
		request["BackendServerPort"] = plan.BackendServerPort.ValueInt64()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() && plan.Description.ValueString() != "" {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.ForwardPort.IsNull() && !plan.ForwardPort.IsUnknown() {
		request["ForwardPort"] = plan.ForwardPort.ValueInt64()
	}
	if !plan.StickySessionType.IsNull() && !plan.StickySessionType.IsUnknown() && plan.StickySessionType.ValueString() != "" {
		request["StickySessionType"] = plan.StickySessionType.ValueString()
	}
	if !plan.CookieTimeout.IsNull() && !plan.CookieTimeout.IsUnknown() && plan.CookieTimeout.ValueInt64() > 0 {
		request["CookieTimeout"] = plan.CookieTimeout.ValueInt64()
	}
	if !plan.Cookie.IsNull() && !plan.Cookie.IsUnknown() && plan.Cookie.ValueString() != "" {
		request["Cookie"] = plan.Cookie.ValueString()
	}
	if plan.HealthCheck.ValueString() == "on" {
		if !plan.HealthCheckDomain.IsNull() && !plan.HealthCheckDomain.IsUnknown() && plan.HealthCheckDomain.ValueString() != "" {
			request["HealthCheckDomain"] = plan.HealthCheckDomain.ValueString()
		}
		if !plan.HealthCheckURI.IsNull() && !plan.HealthCheckURI.IsUnknown() && plan.HealthCheckURI.ValueString() != "" {
			request["HealthCheckURI"] = plan.HealthCheckURI.ValueString()
		}
		request["HealthyThreshold"] = plan.HealthyThreshold.ValueInt64()
		request["UnhealthyThreshold"] = plan.UnhealthyThreshold.ValueInt64()
		request["HealthCheckTimeout"] = plan.HealthCheckTimeout.ValueInt64()
		if !plan.HealthCheckConnectPort.IsNull() && !plan.HealthCheckConnectPort.IsUnknown() && plan.HealthCheckConnectPort.ValueInt64() > 0 {
			request["HealthCheckConnectPort"] = plan.HealthCheckConnectPort.ValueInt64()
		}
		request["HealthCheckInterval"] = plan.HealthCheckInterval.ValueInt64()
		request["HealthCheckHttpCode"] = plan.HealthCheckHttpCode.ValueString()
		request["HealthCheckMethod"] = plan.HealthCheckMethod.ValueString()
	}
	return request
}

// ensLoadBalancerHttpsListenerSetRequest builds the SetLoadBalancerHTTPSListenerAttribute
// request from the planned state. Only updatable attributes are sent.
func ensLoadBalancerHttpsListenerSetRequest(plan ensLoadBalancerHttpsListenerResourceModel) map[string]interface{} {
	parts := splitEnsLoadBalancerHttpsListenerId(plan.Id.ValueString())
	request := map[string]interface{}{
		"LoadBalancerId":      parts[0],
		"ListenerPort":        parts[1],
		"ServerCertificateId": plan.ServerCertificateId.ValueString(),
		"Description":         plan.Description.ValueString(),
		"Scheduler":           plan.Scheduler.ValueString(),
		"HealthCheck":         plan.HealthCheck.ValueString(),
		"IdleTimeout":         plan.IdleTimeout.ValueInt64(),
		"RequestTimeout":      plan.RequestTimeout.ValueInt64(),
	}
	if plan.HealthCheck.ValueString() == "on" {
		if !plan.HealthCheckDomain.IsNull() && !plan.HealthCheckDomain.IsUnknown() && plan.HealthCheckDomain.ValueString() != "" {
			request["HealthCheckDomain"] = plan.HealthCheckDomain.ValueString()
		}
		if !plan.HealthCheckURI.IsNull() && !plan.HealthCheckURI.IsUnknown() && plan.HealthCheckURI.ValueString() != "" {
			request["HealthCheckURI"] = plan.HealthCheckURI.ValueString()
		}
		request["HealthyThreshold"] = plan.HealthyThreshold.ValueInt64()
		request["UnhealthyThreshold"] = plan.UnhealthyThreshold.ValueInt64()
		request["HealthCheckTimeout"] = plan.HealthCheckTimeout.ValueInt64()
		if !plan.HealthCheckConnectPort.IsNull() && !plan.HealthCheckConnectPort.IsUnknown() && plan.HealthCheckConnectPort.ValueInt64() > 0 {
			request["HealthCheckConnectPort"] = plan.HealthCheckConnectPort.ValueInt64()
		}
		request["HealthCheckInterval"] = plan.HealthCheckInterval.ValueInt64()
		request["HealthCheckHttpCode"] = plan.HealthCheckHttpCode.ValueString()
		request["HealthCheckMethod"] = plan.HealthCheckMethod.ValueString()
	}
	return request
}

// ensLoadBalancerHttpsListenerState maps the DescribeLoadBalancerHTTPSListenerAttribute
// response onto the resource model. The status is already normalised to
// lowercase by the service Describe function. Health-check sub-attributes are
// only mapped while health_check is on: the ENS Describe API omits them when
// health check is off (HealthCheckConnectPort is still returned but is
// meaningless), so mapping them in that state would store spurious values.
func ensLoadBalancerHttpsListenerState(object map[string]interface{}) (ensLoadBalancerHttpsListenerResourceModel, error) {
	model := ensLoadBalancerHttpsListenerResourceModel{
		LoadBalancerId:      ensLoadBalancerHttpsListenerString(object["LoadBalancerId"]),
		ServerCertificateId: ensLoadBalancerHttpsListenerString(object["ServerCertificateId"]),
		Description:         ensLoadBalancerHttpsListenerString(object["Description"]),
		Scheduler:           ensLoadBalancerHttpsListenerString(object["Scheduler"]),
		HealthCheck:         ensLoadBalancerHttpsListenerString(object["HealthCheck"]),
		ListenerForward:     ensLoadBalancerHttpsListenerString(object["ListenerForward"]),
		Status:              ensLoadBalancerHttpsListenerString(object["Status"]),
	}

	listenerPort, err := convert.FormatInt(object["ListenerPort"])
	if err != nil {
		return model, err
	}
	model.ListenerPort = ensLoadBalancerHttpsListenerInt64(listenerPort)
	backendServerPort, err := convert.FormatInt(object["BackendServerPort"])
	if err != nil {
		return model, err
	}
	model.BackendServerPort = ensLoadBalancerHttpsListenerInt64(backendServerPort)
	idleTimeout, err := convert.FormatInt(object["IdleTimeout"])
	if err != nil {
		return model, err
	}
	model.IdleTimeout = ensLoadBalancerHttpsListenerInt64(idleTimeout)
	requestTimeout, err := convert.FormatInt(object["RequestTimeout"])
	if err != nil {
		return model, err
	}
	model.RequestTimeout = ensLoadBalancerHttpsListenerInt64(requestTimeout)
	forwardPort, err := convert.FormatInt(object["ForwardPort"])
	if err != nil {
		return model, err
	}
	model.ForwardPort = ensLoadBalancerHttpsListenerInt64(forwardPort)

	if ensLoadBalancerHttpsListenerString(object["HealthCheck"]).ValueString() == "on" {
		model.HealthCheckDomain = ensLoadBalancerHttpsListenerString(object["HealthCheckDomain"])
		model.HealthCheckURI = ensLoadBalancerHttpsListenerString(object["HealthCheckURI"])
		healthyThreshold, err := convert.FormatInt(object["HealthyThreshold"])
		if err != nil {
			return model, err
		}
		model.HealthyThreshold = ensLoadBalancerHttpsListenerInt64(healthyThreshold)
		unhealthyThreshold, err := convert.FormatInt(object["UnhealthyThreshold"])
		if err != nil {
			return model, err
		}
		model.UnhealthyThreshold = ensLoadBalancerHttpsListenerInt64(unhealthyThreshold)
		healthCheckTimeout, err := convert.FormatInt(object["HealthCheckTimeout"])
		if err != nil {
			return model, err
		}
		model.HealthCheckTimeout = ensLoadBalancerHttpsListenerInt64(healthCheckTimeout)
		healthCheckConnectPort, err := convert.FormatInt(object["HealthCheckConnectPort"])
		if err != nil {
			return model, err
		}
		model.HealthCheckConnectPort = ensLoadBalancerHttpsListenerInt64(healthCheckConnectPort)
		healthCheckInterval, err := convert.FormatInt(object["HealthCheckInterval"])
		if err != nil {
			return model, err
		}
		model.HealthCheckInterval = ensLoadBalancerHttpsListenerInt64(healthCheckInterval)
		model.HealthCheckHttpCode = ensLoadBalancerHttpsListenerString(object["HealthCheckHttpCode"])
		model.HealthCheckMethod = ensLoadBalancerHttpsListenerString(object["HealthCheckMethod"])
	}
	return model, nil
}

// ensLoadBalancerHttpsListenerPreserveWriteOnly keeps the write-only attributes
// (sticky_session_type, cookie_timeout, cookie) and, when the Describe API omits
// it, server_certificate_id from the configured/prior value so they do not
// surface as empty and trigger spurious replacements or updates.
func ensLoadBalancerHttpsListenerPreserveWriteOnly(model *ensLoadBalancerHttpsListenerResourceModel, source ensLoadBalancerHttpsListenerResourceModel) {
	if model.StickySessionType.IsNull() && !source.StickySessionType.IsNull() && !source.StickySessionType.IsUnknown() {
		model.StickySessionType = source.StickySessionType
	}
	if model.CookieTimeout.IsNull() && !source.CookieTimeout.IsNull() && !source.CookieTimeout.IsUnknown() {
		model.CookieTimeout = source.CookieTimeout
	}
	if model.Cookie.IsNull() && !source.Cookie.IsNull() && !source.Cookie.IsUnknown() {
		model.Cookie = source.Cookie
	}
	if model.ServerCertificateId.IsNull() && !source.ServerCertificateId.IsNull() && !source.ServerCertificateId.IsUnknown() {
		model.ServerCertificateId = source.ServerCertificateId
	}
}

func ensLoadBalancerHttpsListenerString(raw interface{}) types.String {
	if v, ok := raw.(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// ensLoadBalancerHttpsListenerInt64 returns a known Int64 for a positive value
// and null otherwise; the ENS Describe API omits many numeric fields when they
// are not applicable (e.g. health check fields while health_check is off).
func ensLoadBalancerHttpsListenerInt64(v int) types.Int64 {
	if v > 0 {
		return types.Int64Value(int64(v))
	}
	return types.Int64Null()
}

// ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier suppresses a
// string health-check sub-field while health_check is off.
type ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier struct{}

func (m ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var healthCheck types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("health_check"), &healthCheck)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if healthCheck.ValueString() == "off" {
		// Keep the prior state value (or null) so the configured-but-ignored
		// value does not surface as a diff while health check is off.
		resp.PlanValue = req.StateValue
	}
}

func (m ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier) Description(ctx context.Context) string {
	return "Suppresses the attribute while health_check is off."
}

func (m ensLoadBalancerHttpsListenerHealthCheckOffStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Suppresses the attribute while health_check is off."
}

// ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier suppresses an
// int64 health-check sub-field while health_check is off.
type ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier struct{}

func (m ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	var healthCheck types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("health_check"), &healthCheck)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if healthCheck.ValueString() == "off" {
		resp.PlanValue = req.StateValue
	}
}

func (m ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier) Description(ctx context.Context) string {
	return "Suppresses the attribute while health_check is off."
}

func (m ensLoadBalancerHttpsListenerHealthCheckOffInt64PlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Suppresses the attribute while health_check is off."
}
