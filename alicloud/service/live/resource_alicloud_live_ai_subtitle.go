// Package live. AI subtitle template resource, Framework handwritten implementation.
//
// alicloud_live_ai_subtitle manages a Live AI subtitle template through the
// live/2016-11-01 Add/Describe/Update/DeleteLiveAISubtitle RPCs. SubtitleId is the
// primary key: Add returns it, Describe/Update/Delete address the template by it.
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
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const liveAISubtitleTypeName = "alicloud_live_ai_subtitle"

// liveAISubtitleTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); terraform-plugin-framework does not support
// user-configurable resource timeouts yet.
const liveAISubtitleTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &liveAISubtitleResource{}
	_ resource.ResourceWithImportState = &liveAISubtitleResource{}
)

func NewLiveAISubtitleResource() resource.Resource {
	return &liveAISubtitleResource{}
}

type liveAISubtitleResource struct {
	fwadapt.ResourceBase
}

type liveAISubtitleResourceModel struct {
	Id                    types.String  `tfsdk:"id"`
	SubtitleName          types.String  `tfsdk:"subtitle_name"`
	SrcLanguage           types.String  `tfsdk:"src_language"`
	DstLanguage           types.String  `tfsdk:"dst_language"`
	FontColor             types.String  `tfsdk:"font_color"`
	FontName              types.String  `tfsdk:"font_name"`
	BgColor               types.String  `tfsdk:"bg_color"`
	PositionNormalized    types.List    `tfsdk:"position_normalized"`
	ShowSourceLan         types.Bool    `tfsdk:"show_source_lan"`
	MaxLines              types.Int64   `tfsdk:"max_lines"`
	WordPerLine           types.Int64   `tfsdk:"word_per_line"`
	FontSizeNormalized    types.Float64 `tfsdk:"font_size_normalized"`
	BorderWidthNormalized types.Float64 `tfsdk:"border_width_normalized"`
	BgWidthNormalized     types.Float64 `tfsdk:"bg_width_normalized"`
	Description           types.String  `tfsdk:"description"`
	Width                 types.String  `tfsdk:"width"`
	Height                types.String  `tfsdk:"height"`
	CopyFrom              types.String  `tfsdk:"copy_from"`
	RulesRefer            types.List    `tfsdk:"rules_refer"`
}

func (r *liveAISubtitleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = liveAISubtitleTypeName
}

func (r *liveAISubtitleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subtitle_name": schema.StringAttribute{
				Required: true,
				Description: "The name of the AI subtitle template. It supports only digits, letters and hyphens (-), " +
					"and the hyphen cannot be the first character.",
			},
			"src_language": schema.StringAttribute{
				Required:    true,
				Description: "The source language. Valid values: en-US, zh-CN, ru-RU.",
			},
			"dst_language": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The translation target language. Valid values: en-US, zh-CN, es-ES, ru-RU.",
			},
			"font_color": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The font color, in RGBA format, e.g. 0xFFFFFF.",
			},
			"font_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The font. Valid values: KaiTi, AlibabaPuHuiTi-Regular, AlibabaPuHuiTi-Bold, " +
					"AlibabaPuHuiTi-Light, NotoSansHans-Regular, NotoSansHans-Bold, NotoSansHans-Light.",
			},
			"bg_color": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The subtitle background color, in RGBA format, e.g. 0xFF0000.",
			},
			"position_normalized": schema.ListAttribute{
				Required:    true,
				ElementType: types.Float64Type,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(2),
				},
				Description: "The subtitle position, a pair of normalized [x, y] coordinates measured from " +
					"the lower-left corner of the screen. The actual pixel offset is the value multiplied by " +
					"the preview height, e.g. [0.32, 0.27].",
			},
			"show_source_lan": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Whether to display the source language. Default value: false.",
			},
			"max_lines": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "The number of lines to display.",
			},
			"word_per_line": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 500),
				},
				Description: "The number of characters displayed per line. Valid range: [1, 500].",
			},
			"font_size_normalized": schema.Float64Attribute{
				Required: true,
				Validators: []validator.Float64{
					float64validator.Between(0, 1),
				},
				Description: "The font size. Valid range: [0, 1].",
			},
			"border_width_normalized": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Float64{
					float64validator.Between(0, 1),
				},
				Description: "The font weight (border width). Valid range: [0, 1].",
			},
			"bg_width_normalized": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Float64{
					float64validator.Between(0, 1),
				},
				Description: "The subtitle background size. Valid range: [0, 1].",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The description of the AI subtitle template.",
			},
			"width": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The preview screen width, in pixels (px).",
			},
			"height": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The preview screen height, in pixels (px).",
			},
			"copy_from": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of an existing subtitle template to copy settings from when creating. " +
					"Changing this forces a new resource.",
			},
			"rules_refer": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The IDs of subtitle rules that reference this template.",
			},
		},
	}
}

// positionJSON serializes the configured position list into the JSON array string the
// AddLiveAISubtitle/UpdateLiveAISubtitle RPCs expect (PositionNormalized, style: json).
func positionJSON(l types.List) (string, error) {
	values := make([]float64, 0, len(l.Elements()))
	for _, e := range l.Elements() {
		f, ok := e.(types.Float64)
		if !ok {
			return "", fmt.Errorf("position_normalized element is not a float")
		}
		values = append(values, f.ValueFloat64())
	}
	b, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *liveAISubtitleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan liveAISubtitleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "AddLiveAISubtitle"
	position, err := positionJSON(plan.PositionNormalized)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: encoding position_normalized", liveAISubtitleTypeName),
			err.Error(),
		)
		return
	}
	request := map[string]interface{}{
		"RegionId":           r.Meta().RegionId,
		"ClientToken":        helper.BuildClientToken(action),
		"SubtitleName":       plan.SubtitleName.ValueString(),
		"SrcLanguage":        plan.SrcLanguage.ValueString(),
		"PositionNormalized": position,
		"WordPerLine":        plan.WordPerLine.ValueInt64(),
		"FontSizeNormalized": plan.FontSizeNormalized.ValueFloat64(),
		"ShowSourceLan":      plan.ShowSourceLan.ValueBool(),
	}
	addOptionalLiveAISubtitleRequest(request, &plan)
	// copy_from only seeds a new template from an existing one; UpdateLiveAISubtitle has no
	// such parameter, so it is sent exclusively on Create (its ForceNew semantics guarantee
	// Update is never reached on a copy_from change anyway).
	if !plan.CopyFrom.IsNull() {
		request["CopyFrom"] = plan.CopyFrom.ValueString()
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(liveAISubtitleTimeout, func() *retry.RetryError {
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
			fmt.Sprintf("Creating %s: calling AddLiveAISubtitle", liveAISubtitleTypeName),
			errs.WrapErrorf(err, "Creating %s: calling AddLiveAISubtitle", liveAISubtitleTypeName).Error(),
		)
		return
	}

	id, _ := response["SubtitleId"].(string)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: empty SubtitleId", liveAISubtitleTypeName),
			fmt.Sprintf("AddLiveAISubtitle returned no SubtitleId: %v", response),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveAISubtitle(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLiveAISubtitle", liveAISubtitleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLiveAISubtitle", liveAISubtitleTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, liveAISubtitleState(object, &plan))...)
}

func (r *liveAISubtitleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state liveAISubtitleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveAISubtitle(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLiveAISubtitle", liveAISubtitleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLiveAISubtitle", liveAISubtitleTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, liveAISubtitleState(object, &state))...)
}

func (r *liveAISubtitleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state liveAISubtitleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "UpdateLiveAISubtitle"
	position, err := positionJSON(plan.PositionNormalized)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: encoding position_normalized", liveAISubtitleTypeName),
			err.Error(),
		)
		return
	}
	request := map[string]interface{}{
		"RegionId":           r.Meta().RegionId,
		"ClientToken":        helper.BuildClientToken(action),
		"SubtitleId":         state.Id.ValueString(),
		"SubtitleName":       plan.SubtitleName.ValueString(),
		"SrcLanguage":        plan.SrcLanguage.ValueString(),
		"PositionNormalized": position,
		"WordPerLine":        plan.WordPerLine.ValueInt64(),
		"FontSizeNormalized": plan.FontSizeNormalized.ValueFloat64(),
		"ShowSourceLan":      plan.ShowSourceLan.ValueBool(),
	}
	addOptionalLiveAISubtitleRequest(request, &plan)

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(liveAISubtitleTimeout, func() *retry.RetryError {
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
			fmt.Sprintf("Updating %s: calling UpdateLiveAISubtitle", liveAISubtitleTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateLiveAISubtitle", liveAISubtitleTypeName).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeLiveAISubtitle(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLiveAISubtitle", liveAISubtitleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLiveAISubtitle", liveAISubtitleTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, liveAISubtitleState(object, &plan))...)
}

func (r *liveAISubtitleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state liveAISubtitleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "DeleteLiveAISubtitle"
	request := map[string]interface{}{
		"RegionId":     r.Meta().RegionId,
		"ClientToken":  helper.BuildClientToken(action),
		"SubtitleId":   state.Id.ValueString(),
		"SubtitleName": state.SubtitleName.ValueString(),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(liveAISubtitleTimeout, func() *retry.RetryError {
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
		// DeleteLiveAISubtitle reports no dedicated NotFound code; if the delete fails because
		// the template is already gone, a follow-up Describe confirms absence so the destroy is
		// still considered successful.
		if _, getErr := NewService(r.Meta()).DescribeLiveAISubtitle(state.Id.ValueString()); getErr != nil && errs.NotFoundError(getErr) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteLiveAISubtitle", liveAISubtitleTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteLiveAISubtitle", liveAISubtitleTypeName).Error(),
		)
		return
	}
}

func (r *liveAISubtitleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// addOptionalLiveAISubtitleRequest fills optional fields only when configured, so unset
// options are not sent and the server keeps its own defaults; this mirrors the request
// definition where these parameters are not required.
func addOptionalLiveAISubtitleRequest(request map[string]interface{}, plan *liveAISubtitleResourceModel) {
	if !plan.DstLanguage.IsNull() {
		request["DstLanguage"] = plan.DstLanguage.ValueString()
	}
	if !plan.FontColor.IsNull() {
		request["FontColor"] = plan.FontColor.ValueString()
	}
	if !plan.FontName.IsNull() {
		request["FontName"] = plan.FontName.ValueString()
	}
	if !plan.BgColor.IsNull() {
		request["BgColor"] = plan.BgColor.ValueString()
	}
	if !plan.MaxLines.IsNull() {
		request["MaxLines"] = plan.MaxLines.ValueInt64()
	}
	if !plan.BorderWidthNormalized.IsNull() {
		request["BorderWidthNormalized"] = plan.BorderWidthNormalized.ValueFloat64()
	}
	if !plan.BgWidthNormalized.IsNull() {
		request["BgWidthNormalized"] = plan.BgWidthNormalized.ValueFloat64()
	}
	if !plan.Description.IsNull() {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.Width.IsNull() {
		request["Width"] = plan.Width.ValueString()
	}
	if !plan.Height.IsNull() {
		request["Height"] = plan.Height.ValueString()
	}
}

// liveAISubtitleState maps a DescribeLiveAISubtitle item onto the resource model. The
// response key names differ from the request in two places: WordPerLine is returned as
// WordPerline and FontSizeNormalized is returned as a string, while PositionNormalized
// is wrapped in a {Position: [...]} object. rules_refer is read-only and never sent back.
func liveAISubtitleState(object map[string]interface{}, prior *liveAISubtitleResourceModel) liveAISubtitleResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	int64Val := func(key string) types.Int64 {
		switch v := object[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	float64Val := func(key string) types.Float64 {
		switch v := object[key].(type) {
		case float64:
			return types.Float64Value(v)
		case json.Number:
			if f, err := v.Float64(); err == nil {
				return types.Float64Value(f)
			}
		case int:
			return types.Float64Value(float64(v))
		case int64:
			return types.Float64Value(float64(v))
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return types.Float64Value(f)
			}
		}
		return types.Float64Null()
	}

	positionList := func() types.List {
		values := make([]float64, 0, 2)
		if pn, ok := object["PositionNormalized"].(map[string]interface{}); ok {
			if arr, ok := pn["Position"].([]interface{}); ok {
				for _, p := range arr {
					values = append(values, toFloat64Value(p))
				}
			}
		}
		l, diags := types.ListValueFrom(context.Background(), types.Float64Type, values)
		if diags.HasError() {
			return types.ListNull(types.Float64Type)
		}
		return l
	}

	rulesList := func() types.List {
		values := make([]string, 0)
		if rr, ok := object["RulesRefer"].(map[string]interface{}); ok {
			if ids, ok := rr["RulesId"].([]interface{}); ok {
				for _, id := range ids {
					if s, ok := id.(string); ok {
						values = append(values, s)
					}
				}
			}
		}
		l, diags := types.ListValueFrom(context.Background(), types.StringType, values)
		if diags.HasError() {
			return types.ListNull(types.StringType)
		}
		return l
	}

	showSourceLan := func() types.Bool {
		switch v := object["ShowSourceLan"].(type) {
		case bool:
			return types.BoolValue(v)
		case float64:
			return types.BoolValue(v != 0)
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.BoolValue(n != 0)
			}
		case int:
			return types.BoolValue(v != 0)
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return types.BoolValue(b)
			}
		}
		return types.BoolNull()
	}

	m := liveAISubtitleResourceModel{
		Id:                    str("SubtitleId"),
		SubtitleName:          str("SubtitleName"),
		SrcLanguage:           str("SrcLanguage"),
		DstLanguage:           str("DstLanguage"),
		FontColor:             str("FontColor"),
		FontName:              str("FontName"),
		BgColor:               str("BgColor"),
		PositionNormalized:    positionList(),
		ShowSourceLan:         showSourceLan(),
		MaxLines:              int64Val("MaxLines"),
		WordPerLine:           int64Val("WordPerline"),
		FontSizeNormalized:    float64Val("FontSizeNormalized"),
		BorderWidthNormalized: float64Val("BorderWidthNormalized"),
		BgWidthNormalized:     float64Val("BgWidthNormalized"),
		Description:           str("Description"),
		Width:                 str("Width"),
		Height:                str("Height"),
		RulesRefer:            rulesList(),
	}
	// copy_from is write-only and never returned by DescribeLiveAISubtitle, so the prior
	// configured value is preserved across refresh to avoid a permanent diff.
	if prior != nil {
		m.CopyFrom = prior.CopyFrom
	}
	return m
}

func toFloat64Value(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case json.Number:
		if f, err := n.Float64(); err == nil {
			return f
		}
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case string:
		if f, err := strconv.ParseFloat(n, 64); err == nil {
			return f
		}
	}
	return 0
}
