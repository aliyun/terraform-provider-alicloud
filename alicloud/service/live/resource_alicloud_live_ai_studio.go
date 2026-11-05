package live

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	liveAIStudioTypeName = "alicloud_live_ai_studio"
	// liveAIStudioTimeout replaces the SDKv2 resource-level create/delete default
	// timeouts (5 minutes each); terraform-plugin-framework does not support
	// user-configurable resource timeouts yet.
	liveAIStudioTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &liveAIStudioResource{}
	_ resource.ResourceWithImportState = &liveAIStudioResource{}
)

func NewLiveAIStudioResource() resource.Resource {
	return &liveAIStudioResource{}
}

type liveAIStudioResource struct {
	fwadapt.ResourceBase
}

type liveAIStudioResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	StudioId              types.String `tfsdk:"studio_id"`
	StudioName            types.String `tfsdk:"studio_name"`
	Description           types.String `tfsdk:"description"`
	Width                 types.Int64  `tfsdk:"width"`
	Height                types.Int64  `tfsdk:"height"`
	BackgroundType        types.String `tfsdk:"background_type"`
	BackgroundResourceId  types.String `tfsdk:"background_resource_id"`
	BackgroundResourceUrl types.String `tfsdk:"background_resource_url"`
	MediaType             types.String `tfsdk:"media_type"`
	MediaResourceId       types.String `tfsdk:"media_resource_id"`
	MediaResourceUrl      types.String `tfsdk:"media_resource_url"`
	MediaLayout           types.String `tfsdk:"media_layout"`
	MattingType           types.String `tfsdk:"matting_type"`
	MattingLayout         types.String `tfsdk:"matting_layout"`
	RuleIds               types.List   `tfsdk:"rule_ids"`
}

func (r *liveAIStudioResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = liveAIStudioTypeName
}

func (r *liveAIStudioResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// studio_id mirrors id (the StudioId returned by CreateLiveAIStudio
			// and echoed as TemplateId by DescribeLiveAIStudio) so users can
			// reference the cloud identity directly.
			"studio_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// studio_name is the human-readable studio name. CreateLiveAIStudio
			// requires it and the API documents it as unique. ModifyLiveAIStudio
			// locates the studio by StudioName rather than StudioId and exposes
			// no rename parameter, so studio_name is immutable in place.
			"studio_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"width": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"height": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"background_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("VOD", "LIVE", "PIC"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// background_resource_id and background_resource_url are mutually
			// exclusive ("二选一"); the API rejects both being set.
			"background_resource_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"background_resource_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"media_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("VOD", "LIVE", "PIC"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"media_resource_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"media_resource_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// media_layout / matting_layout carry a JSON object serialized as a
			// string, e.g. {"positionX":"0","positionY":"0","heightNormalized":"0.5"}.
			// They are passed through verbatim so the API's encoding round-trips
			// exactly.
			"media_layout": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"matting_layout": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"matting_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("green", "blue", "complex"),
				},
			},
			// rule_ids is read-only; the API returns the rule list attached to a
			// studio and exposes no call to manage it through this resource.
			"rule_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *liveAIStudioResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan liveAIStudioResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := liveAIStudioCreateRequest(plan)
	request["RegionId"] = r.Meta().RegionId
	action := "CreateLiveAIStudio"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(ctx, liveAIStudioTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("live", "2016-11-01", action, map[string]interface{}{}, request, true)
		if callErr != nil {
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
			fmt.Sprintf("Creating %s: calling CreateLiveAIStudio", liveAIStudioTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, liveAIStudioTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id, _ := response["StudioId"].(string)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading StudioId from response", liveAIStudioTypeName),
			fmt.Sprintf("CreateLiveAIStudio did not return a StudioId: %v", response),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveAIStudioContext(ctx, id)
	if err != nil {
		// Keep the created id in state so the studio stays trackable even when
		// the read-back fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling DescribeLiveAIStudio", liveAIStudioTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeLiveAIStudio", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mErr := liveAIStudioState(object)
	if mErr != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping DescribeLiveAIStudio response", liveAIStudioTypeName),
			mErr.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// Update re-applies the configured studio attributes through ModifyLiveAIStudio.
// ModifyLiveAIStudio locates the studio by StudioName (not StudioId) and accepts
// the same configurable set as CreateLiveAIStudio; studio_name itself is
// replacement-only, so it is read from state as the immutable identifier.
func (r *liveAIStudioResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan liveAIStudioResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state liveAIStudioResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := liveAIStudioModifyRequest(plan, state)
	request["RegionId"] = r.Meta().RegionId
	action := "ModifyLiveAIStudio"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(ctx, liveAIStudioTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("live", "2016-11-01", action, map[string]interface{}{}, request, true)
		if callErr != nil {
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
			fmt.Sprintf("Updating %s: calling ModifyLiveAIStudio", liveAIStudioTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveAIStudioContext(ctx, state.Id.ValueString())
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling DescribeLiveAIStudio", liveAIStudioTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeLiveAIStudio", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mErr := liveAIStudioState(object)
	if mErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: mapping DescribeLiveAIStudio response", liveAIStudioTypeName),
			mErr.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *liveAIStudioResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state liveAIStudioResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveAIStudioContext(ctx, id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLiveAIStudio", liveAIStudioTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeLiveAIStudio", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mErr := liveAIStudioState(object)
	if mErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping DescribeLiveAIStudio response", liveAIStudioTypeName),
			mErr.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *liveAIStudioResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state liveAIStudioResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if err := NewService(r.Meta()).DeleteLiveAIStudioContext(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteLiveAIStudio", liveAIStudioTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteLiveAIStudio", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *liveAIStudioResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// liveAIStudioCreateRequest builds the CreateLiveAIStudio request from the
// planned state. Required attributes (StudioName, MattingType) are always sent;
// optional attributes are omitted when null or empty so the API keeps its
// defaults instead of receiving spurious empty values.
func liveAIStudioCreateRequest(plan liveAIStudioResourceModel) map[string]interface{} {
	request := map[string]interface{}{
		"StudioName":  plan.StudioName.ValueString(),
		"MattingType": plan.MattingType.ValueString(),
		"ClientToken": helper.BuildClientToken("CreateLiveAIStudio"),
	}
	if v := liveAIStudioString(plan.Description); v != "" {
		request["Description"] = v
	}
	if !plan.Width.IsNull() && !plan.Width.IsUnknown() {
		request["Width"] = plan.Width.ValueInt64()
	}
	if !plan.Height.IsNull() && !plan.Height.IsUnknown() {
		request["Height"] = plan.Height.ValueInt64()
	}
	if v := liveAIStudioString(plan.BackgroundType); v != "" {
		request["BackgroundType"] = v
	}
	if v := liveAIStudioString(plan.BackgroundResourceId); v != "" {
		request["BackgroundResourceId"] = v
	}
	if v := liveAIStudioString(plan.BackgroundResourceUrl); v != "" {
		request["BackgroundResourceUrl"] = v
	}
	if v := liveAIStudioString(plan.MediaType); v != "" {
		request["MediaType"] = v
	}
	if v := liveAIStudioString(plan.MediaResourceId); v != "" {
		request["MediaResourceId"] = v
	}
	if v := liveAIStudioString(plan.MediaResourceUrl); v != "" {
		request["MediaResourceUrl"] = v
	}
	if v := liveAIStudioString(plan.MediaLayout); v != "" {
		request["MediaLayout"] = v
	}
	if v := liveAIStudioString(plan.MattingLayout); v != "" {
		request["MattingLayout"] = v
	}
	return request
}

// liveAIStudioModifyRequest builds the ModifyLiveAIStudio request. StudioName
// is the immutable API identifier and is taken from prior state; every other
// configurable attribute is taken from the plan.
func liveAIStudioModifyRequest(plan, state liveAIStudioResourceModel) map[string]interface{} {
	request := map[string]interface{}{
		"StudioName": state.StudioName.ValueString(),
	}
	if v := liveAIStudioString(plan.Description); v != "" {
		request["Description"] = v
	}
	if !plan.Width.IsNull() && !plan.Width.IsUnknown() {
		request["Width"] = plan.Width.ValueInt64()
	}
	if !plan.Height.IsNull() && !plan.Height.IsUnknown() {
		request["Height"] = plan.Height.ValueInt64()
	}
	if v := liveAIStudioString(plan.BackgroundType); v != "" {
		request["BackgroundType"] = v
	}
	if v := liveAIStudioString(plan.BackgroundResourceId); v != "" {
		request["BackgroundResourceId"] = v
	}
	if v := liveAIStudioString(plan.BackgroundResourceUrl); v != "" {
		request["BackgroundResourceUrl"] = v
	}
	if v := liveAIStudioString(plan.MediaType); v != "" {
		request["MediaType"] = v
	}
	if v := liveAIStudioString(plan.MediaResourceId); v != "" {
		request["MediaResourceId"] = v
	}
	if v := liveAIStudioString(plan.MediaResourceUrl); v != "" {
		request["MediaResourceUrl"] = v
	}
	if v := liveAIStudioString(plan.MediaLayout); v != "" {
		request["MediaLayout"] = v
	}
	if v := liveAIStudioString(plan.MattingLayout); v != "" {
		request["MattingLayout"] = v
	}
	if v := liveAIStudioString(plan.MattingType); v != "" {
		request["MattingType"] = v
	}
	return request
}

// liveAIStudioState maps a DescribeLiveAIStudio SubtitleConfig item to the
// resource model. CreateLiveAIStudio/DeleteLiveAIStudio use StudioId/StudioName,
// while DescribeLiveAIStudio echoes them as TemplateId/TemplateName; both name
// the same identity.
func liveAIStudioState(object map[string]interface{}) (liveAIStudioResourceModel, error) {
	id, _ := object["TemplateId"].(string)
	if id == "" {
		return liveAIStudioResourceModel{}, fmt.Errorf("TemplateId is missing in the DescribeLiveAIStudio response")
	}
	studioName, _ := object["TemplateName"].(string)

	width, wErr := liveAIStudioInt64(object["Width"])
	if wErr != nil {
		return liveAIStudioResourceModel{}, fmt.Errorf("Width: %w", wErr)
	}
	height, hErr := liveAIStudioInt64(object["Height"])
	if hErr != nil {
		return liveAIStudioResourceModel{}, fmt.Errorf("Height: %w", hErr)
	}

	ruleIds, rErr := liveAIStudioRuleIds(object["RuleIds"])
	if rErr != nil {
		return liveAIStudioResourceModel{}, fmt.Errorf("RuleIds: %w", rErr)
	}

	return liveAIStudioResourceModel{
		Id:                    types.StringValue(id),
		StudioId:              types.StringValue(id),
		StudioName:            liveAIStudioStringState(studioName),
		Description:           liveAIStudioStringState(liveAIStudioStringField(object["Description"])),
		Width:                 width,
		Height:                height,
		BackgroundType:        liveAIStudioStringState(liveAIStudioStringField(object["BackgroundType"])),
		BackgroundResourceId:  liveAIStudioStringState(liveAIStudioStringField(object["BackgroundResourceId"])),
		BackgroundResourceUrl: liveAIStudioStringState(liveAIStudioStringField(object["BackgroundResourceUrl"])),
		MediaType:             liveAIStudioStringState(liveAIStudioStringField(object["MediaType"])),
		MediaResourceId:       liveAIStudioStringState(liveAIStudioStringField(object["MediaResourceId"])),
		MediaResourceUrl:      liveAIStudioStringState(liveAIStudioStringField(object["MediaResourceUrl"])),
		MediaLayout:           liveAIStudioStringState(liveAIStudioStringField(object["MediaLayout"])),
		MattingLayout:         liveAIStudioStringState(liveAIStudioStringField(object["MattingLayout"])),
		MattingType:           liveAIStudioStringState(liveAIStudioStringField(object["MattingType"])),
		RuleIds:               ruleIds,
	}, nil
}

// liveAIStudioString returns the string value of a framework String attribute,
// treating null and unknown as empty.
func liveAIStudioString(v types.String) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}
	return v.ValueString()
}

// liveAIStudioStringField coerces a response scalar to string. The tea RPC
// client may decode numbers as json.Number or float64; layout fields are
// already strings, so only a string coercion is needed here.
func liveAIStudioStringField(raw interface{}) string {
	switch v := raw.(type) {
	case nil:
		return ""
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprint(raw)
	}
}

// liveAIStudioStringState converts a string to a framework String value, mapping
// the empty string to null so unset attributes are not written as "".
func liveAIStudioStringState(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// liveAIStudioInt64 parses a response Width/Height, which the API documents as
// integer but may arrive as a string, json.Number or float64.
func liveAIStudioInt64(raw interface{}) (types.Int64, error) {
	switch v := raw.(type) {
	case nil:
		return types.Int64Null(), nil
	case string:
		if v == "" {
			return types.Int64Null(), nil
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return types.Int64Null(), fmt.Errorf("value %q is not an int64: %w", v, err)
		}
		return types.Int64Value(n), nil
	case json.Number:
		n, err := strconv.ParseInt(v.String(), 10, 64)
		if err != nil {
			return types.Int64Null(), fmt.Errorf("value %q is not an int64: %w", v.String(), err)
		}
		return types.Int64Value(n), nil
	case float64:
		return types.Int64Value(int64(v)), nil
	case int:
		return types.Int64Value(int64(v)), nil
	case int64:
		return types.Int64Value(v), nil
	default:
		return types.Int64Null(), fmt.Errorf("unexpected type %T", raw)
	}
}

// liveAIStudioRuleIds reads the RuleIds object ({ "ruleId": ["..."] }) into a
// framework list of strings.
func liveAIStudioRuleIds(raw interface{}) (types.List, error) {
	if raw == nil {
		return types.ListNull(types.StringType), nil
	}
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return types.ListNull(types.StringType), fmt.Errorf("unexpected RuleIds type %T", raw)
	}
	items, ok := obj["ruleId"].([]interface{})
	if !ok {
		return types.ListNull(types.StringType), nil
	}
	elems := make([]attr.Value, 0, len(items))
	for i, item := range items {
		s, ok := item.(string)
		if !ok {
			return types.ListNull(types.StringType), fmt.Errorf("ruleId[%d] has unexpected type %T", i, item)
		}
		elems = append(elems, types.StringValue(s))
	}
	// types.ListValue returns (types.List, diag.Diagnostics); the function
	// keeps an error signature for symmetry with liveAIStudioInt64, so any
	// diagnostic is surfaced as an error instead of being dropped silently.
	list, diags := types.ListValue(types.StringType, elems)
	if diags.HasError() {
		return types.ListNull(types.StringType), fmt.Errorf("building RuleIds list: %v", diags)
	}
	return list, nil
}
