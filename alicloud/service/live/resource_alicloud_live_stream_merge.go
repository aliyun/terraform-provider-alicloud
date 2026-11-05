package live

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const liveStreamMergeTypeName = "alicloud_live_stream_merge"

// liveStreamMergeTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 1 minute); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const liveStreamMergeTimeout = 1 * time.Minute

var (
	_ resource.Resource                = &liveStreamMergeResource{}
	_ resource.ResourceWithImportState = &liveStreamMergeResource{}
)

func NewLiveStreamMergeResource() resource.Resource {
	return &liveStreamMergeResource{}
}

type liveStreamMergeResource struct {
	fwadapt.ResourceBase
}

type liveStreamMergeResourceModel struct {
	Id                types.String `tfsdk:"id"`
	DomainName        types.String `tfsdk:"domain_name"`
	AppName           types.String `tfsdk:"app_name"`
	StreamName        types.String `tfsdk:"stream_name"`
	InAppName1        types.String `tfsdk:"in_app_name_1"`
	InStreamName1     types.String `tfsdk:"in_stream_name_1"`
	InAppName2        types.String `tfsdk:"in_app_name_2"`
	InStreamName2     types.String `tfsdk:"in_stream_name_2"`
	StartTime         types.String `tfsdk:"start_time"`
	EndTime           types.String `tfsdk:"end_time"`
	Protocol          types.String `tfsdk:"protocol"`
	LiveMerger        types.String `tfsdk:"live_merger"`
	MergeParameters   types.String `tfsdk:"merge_parameters"`
	SelectAppName     types.String `tfsdk:"select_app_name"`
	SelectStreamName  types.String `tfsdk:"select_stream_name"`
	SwitchMode        types.String `tfsdk:"switch_mode"`
	AppUsing          types.String `tfsdk:"app_using"`
	StreamUsing       types.String `tfsdk:"stream_using"`
	ExtraInAppStreams types.String `tfsdk:"extra_in_app_streams"`
}

func (r *liveStreamMergeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = liveStreamMergeTypeName
}

func (r *liveStreamMergeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// The output-stream identity (DomainName + AppName + StreamName)
			// is the resource primary key and drives Add/Delete; changing any
			// of them targets a different merge rule, so they force replacement.
			"domain_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"app_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"stream_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"in_app_name_1":    schema.StringAttribute{Required: true},
			"in_stream_name_1": schema.StringAttribute{Required: true},
			"in_app_name_2":    schema.StringAttribute{Required: true},
			"in_stream_name_2": schema.StringAttribute{Required: true},
			"start_time":       schema.StringAttribute{Required: true},
			"end_time":         schema.StringAttribute{Required: true},
			// Protocol and LiveMerger carry API defaults ("rtmp" and "off");
			// they are Optional + Computed so an unset value is sent as the
			// default and refreshed from Describe on every Read.
			"protocol": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("rtmp", "rtc"),
				},
			},
			"live_merger": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
			},
			"merge_parameters":   schema.StringAttribute{Optional: true, Computed: true},
			"select_app_name":    schema.StringAttribute{Optional: true, Computed: true},
			"select_stream_name": schema.StringAttribute{Optional: true, Computed: true},
			"switch_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("manually", "automatically"),
				},
			},
			// Server-side runtime fields only ever come back from Describe.
			"app_using":            schema.StringAttribute{Computed: true},
			"stream_using":         schema.StringAttribute{Computed: true},
			"extra_in_app_streams": schema.StringAttribute{Computed: true},
		},
	}
}

// buildAddRequest assembles the AddLiveStreamMerge request from a model.
// Protocol and LiveMerger fall back to their API defaults ("rtmp"/"off") when
// the user leaves them unset, matching the server-side behavior so that a
// refreshed plan stays drift-free.
func (m *liveStreamMergeResourceModel) addRequest(regionId string) (action string, request map[string]interface{}) {
	action = "AddLiveStreamMerge"
	protocol := m.Protocol.ValueString()
	if protocol == "" {
		protocol = "rtmp"
	}
	liveMerger := m.LiveMerger.ValueString()
	if liveMerger == "" {
		liveMerger = "off"
	}
	request = map[string]interface{}{
		"DomainName":    m.DomainName.ValueString(),
		"AppName":       m.AppName.ValueString(),
		"StreamName":    m.StreamName.ValueString(),
		"InAppName1":    m.InAppName1.ValueString(),
		"InStreamName1": m.InStreamName1.ValueString(),
		"InAppName2":    m.InAppName2.ValueString(),
		"InStreamName2": m.InStreamName2.ValueString(),
		"StartTime":     m.StartTime.ValueString(),
		"EndTime":       m.EndTime.ValueString(),
		"Protocol":      protocol,
		"LiveMerger":    liveMerger,
		"RegionId":      regionId,
		"ClientToken":   helper.BuildClientToken(action),
	}
	if v := m.MergeParameters.ValueString(); v != "" {
		request["MergeParameters"] = v
	}
	if v := m.SelectAppName.ValueString(); v != "" {
		request["SelectAppName"] = v
	}
	if v := m.SelectStreamName.ValueString(); v != "" {
		request["SelectStreamName"] = v
	}
	if v := m.SwitchMode.ValueString(); v != "" {
		request["SwitchMode"] = v
	}
	return action, request
}

func (r *liveStreamMergeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan liveStreamMergeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, request := plan.addRequest(r.Meta().RegionId)
	if _, err := r.callLive(action, request); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: calling AddLiveStreamMerge", liveStreamMergeTypeName), err.Error())
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveStreamMerge(plan.DomainName.ValueString(), plan.AppName.ValueString(), plan.StreamName.ValueString(), plan.Protocol.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling DescribeLiveStreamMerge", liveStreamMergeTypeName), err.Error())
		return
	}
	state := liveStreamMergeState(object)
	// select_app_name and select_stream_name are write-only (Describe omits
	// them), so the configured values must be carried over from the plan.
	state.SelectAppName = plan.SelectAppName
	state.SelectStreamName = plan.SelectStreamName
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *liveStreamMergeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state liveStreamMergeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainName, appName, streamName, ok := parseLiveStreamMergeId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", liveStreamMergeTypeName), fmt.Sprintf("invalid resource id %q: expected DomainName:AppName:StreamName", state.Id.ValueString()))
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveStreamMerge(domainName, appName, streamName, state.Protocol.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling DescribeLiveStreamMerge", liveStreamMergeTypeName), err.Error())
		return
	}
	model := liveStreamMergeState(object)
	// Preserve write-only fields the Describe response never carries back.
	model.SelectAppName = state.SelectAppName
	model.SelectStreamName = state.SelectStreamName
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *liveStreamMergeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan liveStreamMergeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state liveStreamMergeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// There is no dedicated Update API; AddLiveStreamMerge is idempotent and
	// re-applies the whole merge rule for the same output stream, so Update
	// re-issues Add with the new configuration.
	action, request := plan.addRequest(r.Meta().RegionId)
	if _, err := r.callLive(action, request); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Updating %s: calling AddLiveStreamMerge", liveStreamMergeTypeName), err.Error())
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveStreamMerge(plan.DomainName.ValueString(), plan.AppName.ValueString(), plan.StreamName.ValueString(), plan.Protocol.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling DescribeLiveStreamMerge", liveStreamMergeTypeName), err.Error())
		return
	}
	model := liveStreamMergeState(object)
	model.SelectAppName = plan.SelectAppName
	model.SelectStreamName = plan.SelectStreamName
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *liveStreamMergeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state liveStreamMergeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainName, appName, streamName, ok := parseLiveStreamMergeId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting %s", liveStreamMergeTypeName), fmt.Sprintf("invalid resource id %q: expected DomainName:AppName:StreamName", state.Id.ValueString()))
		return
	}

	action := "DeleteLiveStreamMerge"
	request := map[string]interface{}{
		"DomainName":  domainName,
		"AppName":     appName,
		"StreamName":  streamName,
		"RegionId":    r.Meta().RegionId,
		"ClientToken": helper.BuildClientToken(action),
	}
	_, err := r.callLive(action, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidStream.NotFound", "InvalidDomain.NotFound", "NotFound"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting %s: calling DeleteLiveStreamMerge", liveStreamMergeTypeName), err.Error())
		return
	}
}

func (r *liveStreamMergeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// callLive posts a Live RPC action with the standard incremental retry; it
// returns the response and a wrapped error suitable for Framework diagnostics.
func (r *liveStreamMergeResource) callLive(action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(liveStreamMergeTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("live", "2016-11-01", action, query, request, true)
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
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, request["DomainName"], action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// parseLiveStreamMergeId splits the composite id DomainName:AppName:StreamName.
func parseLiveStreamMergeId(id string) (domainName, appName, streamName string, ok bool) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// liveStreamMergeState maps a DescribeLiveStreamMerge item onto the resource
// model. The DescribeLiveStreamMerge response carries no "Id" field, so the
// composite id "DomainName:AppName:StreamName" is built from the response's
// own identity fields. Write-only fields (select_app_name, select_stream_name)
// are left null and overlaid by each caller from the prior plan/state.
func liveStreamMergeState(object map[string]interface{}) liveStreamMergeResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			// TrimSpace normalizes values the API returns with surrounding
			// whitespace (e.g. MergeParameters comes back with a trailing
			// newline); without it the refreshed state would never match the
			// configured value and the plan would not converge.
			if t := strings.TrimSpace(v); t != "" {
				return types.StringValue(t)
			}
			return types.StringNull()
		}
		return types.StringNull()
	}
	id := types.StringNull()
	if d, a, s, ok := parseLiveStreamMergeId(fmt.Sprintf("%v:%v:%v", object["DomainName"], object["AppName"], object["StreamName"])); ok {
		id = types.StringValue(fmt.Sprintf("%s:%s:%s", d, a, s))
	}
	return liveStreamMergeResourceModel{
		Id:                id,
		DomainName:        str("DomainName"),
		AppName:           str("AppName"),
		StreamName:        str("StreamName"),
		InAppName1:        str("InAppName1"),
		InStreamName1:     str("InStreamName1"),
		InAppName2:        str("InAppName2"),
		InStreamName2:     str("InStreamName2"),
		StartTime:         str("StartTime"),
		EndTime:           str("EndTime"),
		Protocol:          str("Protocol"),
		LiveMerger:        str("LiveMerger"),
		MergeParameters:   str("MergeParameters"),
		SwitchMode:        str("SwitchMode"),
		AppUsing:          str("AppUsing"),
		StreamUsing:       str("StreamUsing"),
		ExtraInAppStreams: str("ExtraInAppStreams"),
	}
}
