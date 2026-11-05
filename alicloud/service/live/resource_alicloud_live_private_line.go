package live

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const livePrivateLineTypeName = "alicloud_live_private_line"

// livePrivateLineTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each several minutes); terraform-plugin-framework does
// not support user-configurable resource timeouts yet.
const (
	livePrivateLineCreateTimeout = 10 * time.Minute
	livePrivateLineDeleteTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &livePrivateLineResource{}
	_ resource.ResourceWithImportState = &livePrivateLineResource{}
)

func NewLivePrivateLineResource() resource.Resource {
	return &livePrivateLineResource{}
}

type livePrivateLineResource struct {
	fwadapt.ResourceBase
}

type livePrivateLineResourceModel struct {
	Id               types.String `tfsdk:"id"`
	DomainName       types.String `tfsdk:"domain_name"`
	AppName          types.String `tfsdk:"app_name"`
	StreamName       types.String `tfsdk:"stream_name"`
	VideoCenter      types.String `tfsdk:"video_center"`
	AccelerationType types.String `tfsdk:"acceleration_type"`
	AccelerationArea types.String `tfsdk:"acceleration_area"`
	Reuse            types.String `tfsdk:"reuse"`
	MaxBandwidth     types.Int64  `tfsdk:"max_bandwidth"`
	InstanceId       types.String `tfsdk:"instance_id"`
	Status           types.String `tfsdk:"status"`
	BindingStatus    types.String `tfsdk:"binding_status"`
	IP               types.String `tfsdk:"ip"`
	Uid              types.String `tfsdk:"uid"`
}

func (r *livePrivateLineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = livePrivateLineTypeName
}

func (r *livePrivateLineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Binds a Global Accelerator (GA) instance to a live streaming link, providing private-line acceleration for push/pull streams.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "The main pull-stream domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"app_name": schema.StringAttribute{
				Required:    true,
				Description: "The application name, matching the AppName in the pull-stream URL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"stream_name": schema.StringAttribute{
				Required:    true,
				Description: "The live stream name, matching the StreamName in the pull-stream URL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"video_center": schema.StringAttribute{
				Required:    true,
				Description: "The live center, e.g. cn-beijing, cn-shanghai, cn-shenzhen, cn-qingdao, ap-northeast-1, ap-southeast-5, eu-central-1, ap-southeast-1, ap-south-1.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"acceleration_type": schema.StringAttribute{
				Required:    true,
				Description: "The acceleration type. Valid values: `play` (pull-stream acceleration), `publish` (push-stream acceleration).",
				Validators: []validator.String{
					stringvalidator.OneOf("play", "publish"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"acceleration_area": schema.StringAttribute{
				Required:    true,
				Description: "The acceleration channel (GA endpoint region), e.g. ap-southeast-1.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reuse": schema.StringAttribute{
				Required:    true,
				Description: "Whether to reuse an existing acceleration channel. Valid values: `yes` (reuse an existing GA instance), `no` (create a new acceleration channel).",
				Validators: []validator.String{
					stringvalidator.OneOf("yes", "no"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// MaxBandwidth is required when reuse=no and write-only: DescribeLivePrivateLineAvailGA
			// never returns it, so refresh keeps the configured value (import leaves it null).
			"max_bandwidth": schema.Int64Attribute{
				Optional:    true,
				Description: "The acceleration bandwidth in Mbps. Required when `reuse` is `no`.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			// InstanceId is required when reuse=yes (the GA instance to bind) and computed when
			// reuse=no (the API creates a GA instance and Describe returns its id).
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The GA instance id. Required when `reuse` is `yes`; computed when `reuse` is `no`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The instance status. `active`: available; `inactive`: unavailable.",
			},
			"binding_status": schema.StringAttribute{
				Computed:    true,
				Description: "Whether the GA instance is bound to the live link. `yes`: bound; `no`: not bound.",
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Description: "The acceleration IP address.",
			},
			"uid": schema.StringAttribute{
				Computed:    true,
				Description: "The UID.",
			},
		},
	}
}

func (r *livePrivateLineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan livePrivateLineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reuse := plan.Reuse.ValueString()
	if reuse == "no" {
		if plan.MaxBandwidth.IsNull() || plan.MaxBandwidth.IsUnknown() || plan.MaxBandwidth.ValueInt64() <= 0 {
			resp.Diagnostics.AddError(
				"Creating alicloud_live_private_line: missing max_bandwidth",
				"max_bandwidth is required when reuse is \"no\".",
			)
			return
		}
	}
	if reuse == "yes" {
		if plan.InstanceId.IsNull() || plan.InstanceId.IsUnknown() || plan.InstanceId.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Creating alicloud_live_private_line: missing instance_id",
				"instance_id is required when reuse is \"yes\".",
			)
			return
		}
	}

	id := buildLivePrivateLineID(
		plan.DomainName.ValueString(),
		plan.AppName.ValueString(),
		plan.StreamName.ValueString(),
		plan.AccelerationType.ValueString(),
	)

	action := "CreateLivePrivateLine"
	request := map[string]interface{}{
		"DomainName":       plan.DomainName.ValueString(),
		"AppName":          plan.AppName.ValueString(),
		"StreamName":       plan.StreamName.ValueString(),
		"VideoCenter":      plan.VideoCenter.ValueString(),
		"AccelerationType": plan.AccelerationType.ValueString(),
		"AccelerationArea": plan.AccelerationArea.ValueString(),
		"Reuse":            reuse,
		"RegionId":         r.Meta().RegionId,
	}
	if reuse == "no" {
		request["MaxBandwidth"] = plan.MaxBandwidth.ValueInt64()
	}
	if reuse == "yes" {
		request["InstanceId"] = plan.InstanceId.ValueString()
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(livePrivateLineCreateTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("live", "2016-11-01", action, nil, request, true)
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
			fmt.Sprintf("Creating %s: calling CreateLivePrivateLine", livePrivateLineTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, livePrivateLineTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	plan.Id = types.StringValue(id)

	// The binding (and, for reuse=no, the GA instance it provisions) is not
	// immediately visible; wait until Describe surfaces an active entry.
	liveService := NewService(r.Meta())
	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, livePrivateLineCreateTimeout, 5*time.Second, liveService.LivePrivateLineStateRefreshFunc(id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for active binding", livePrivateLineTypeName),
			errs.WrapErrorf(err, errs.IdMsg, id).Error(),
		)
		return
	}

	object, err := liveService.DescribeLivePrivateLine(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLivePrivateLineAvailGA", livePrivateLineTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeLivePrivateLineAvailGA", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	state := livePrivateLineState(object)
	// Preserve write-only create-time configuration the Describe entry never echoes back.
	state.MaxBandwidth = plan.MaxBandwidth
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *livePrivateLineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state livePrivateLineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeLivePrivateLine(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLivePrivateLineAvailGA", livePrivateLineTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeLivePrivateLineAvailGA", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	model := livePrivateLineState(object)
	// MaxBandwidth is write-only; keep the previously configured value on refresh.
	model.MaxBandwidth = state.MaxBandwidth
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *livePrivateLineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No Update API exists for LivePrivateLine; every configurable attribute is
	// RequiresReplace, so an in-place Update only needs to reflect the refreshed
	// state back to the framework.
	var state livePrivateLineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *livePrivateLineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state livePrivateLineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainName, appName, streamName, accelerationType, err := parseLivePrivateLineID(state.Id.ValueString())
	if err != nil {
		// Malformed id; treat as already gone rather than looping on an unparseable id.
		return
	}

	action := "DeleteLivePrivateLine"
	request := map[string]interface{}{
		"DomainName":       domainName,
		"AppName":          appName,
		"StreamName":       streamName,
		"AccelerationType": accelerationType,
		"RegionId":         r.Meta().RegionId,
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(livePrivateLineDeleteTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("live", "2016-11-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteLivePrivateLine", livePrivateLineTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// DeleteLivePrivateLine unbinds the GA instance from the live link; wait for
	// the binding entry to disappear. (The GA instance itself is not released by
	// this call and must be destroyed in the GA console — documented separately.)
	liveService := NewService(r.Meta())
	stateConf := helper.BuildStateConf([]string{}, []string{""}, livePrivateLineDeleteTimeout, 2*time.Second, liveService.LivePrivateLineStateRefreshFunc(state.Id.ValueString(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for binding to disappear", livePrivateLineTypeName),
			errs.WrapErrorf(err, errs.IdMsg, state.Id.ValueString()).Error(),
		)
		return
	}
}

func (r *livePrivateLineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The id is the composite identity DomainName:AppName:StreamName:AccelerationType.
	// Read parses it and fills the readable attributes; write-only fields
	// (max_bandwidth) are not reconstructable and stay null after import.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// livePrivateLineState maps a DescribeLivePrivateLineAvailGA item onto the
// resource model. MaxBandwidth is deliberately left null: it is write-only and
// every caller overlays the plan/state value before writing state.
func livePrivateLineState(object map[string]interface{}) livePrivateLineResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return livePrivateLineResourceModel{
		Id:               compositeLivePrivateLineID(object),
		DomainName:       str("DomainName"),
		AppName:          str("AppName"),
		StreamName:       str("StreamName"),
		VideoCenter:      str("VideoCenter"),
		AccelerationType: str("AccelerationType"),
		AccelerationArea: str("AccelerationArea"),
		InstanceId:       str("InstanceId"),
		Status:           str("Status"),
		BindingStatus:    str("BindingStatus"),
		IP:               str("IP"),
		Uid:              str("Uid"),
	}
}

// compositeLivePrivateLineID rebuilds the composite id from a Describe item.
func compositeLivePrivateLineID(object map[string]interface{}) types.String {
	domainName, _ := object["DomainName"].(string)
	appName, _ := object["AppName"].(string)
	streamName, _ := object["StreamName"].(string)
	accelerationType, _ := object["AccelerationType"].(string)
	if domainName == "" || appName == "" || streamName == "" || accelerationType == "" {
		return types.StringNull()
	}
	return types.StringValue(buildLivePrivateLineID(domainName, appName, streamName, accelerationType))
}
