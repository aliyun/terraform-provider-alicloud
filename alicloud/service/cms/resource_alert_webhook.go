package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const alertWebhookTypeName = "alicloud_cms_alert_webhook"

var (
	_ resource.Resource                = &alertWebhookResource{}
	_ resource.ResourceWithImportState = &alertWebhookResource{}
)

func NewAlertWebhookResource() resource.Resource {
	return &alertWebhookResource{}
}

// alertWebhookResource is a native terraform-plugin-framework resource. The
// CRUD actions, request body field names (camelCase) and the batch-delete
// single-id wrap are inherited verbatim from PR #10446's SDKv2 implementation;
// only the schema expression moved from helper/schema to the framework.
type alertWebhookResource struct {
	fwadapt.ResourceBase
}

type alertWebhookModel struct {
	Id               types.String `tfsdk:"id"`
	AlertWebhookName types.String `tfsdk:"alert_webhook_name"`
	Url              types.String `tfsdk:"url"`
	ContentType      types.String `tfsdk:"content_type"`
	Method           types.String `tfsdk:"method"`
	Lang             types.String `tfsdk:"lang"`
	Headers          types.Map    `tfsdk:"headers"`
	Workspace        types.String `tfsdk:"workspace"`
}

func (r *alertWebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = alertWebhookTypeName
}

func (r *alertWebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a CloudMonitor (Cms) Alert Webhook resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the alert webhook.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_webhook_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the alert webhook.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL of the alert webhook.",
			},
			"content_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The content type of the alert webhook request body. Valid values: `JSON`, `FORM`.",
				Validators: []validator.String{
					stringvalidator.OneOf("JSON", "FORM"),
				},
			},
			"method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The HTTP method of the alert webhook request. Valid values: `GET`, `POST`.",
				Validators: []validator.String{
					stringvalidator.OneOf("GET", "POST"),
				},
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The language. Valid values: `zh_CN`, `en_US`.",
				Validators: []validator.String{
					stringvalidator.OneOf("zh_CN", "en_US"),
				},
			},
			"headers": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The custom request headers of the alert webhook.",
			},
			"workspace": schema.StringAttribute{
				Computed:    true,
				Description: "The workspace to which the alert webhook belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *alertWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alertWebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	action := "/webhook"

	request := map[string]interface{}{
		"name": plan.AlertWebhookName.ValueString(),
		"url":  plan.Url.ValueString(),
	}
	if v := plan.ContentType.ValueString(); v != "" {
		request["contentType"] = v
	}
	if v := plan.Method.ValueString(); v != "" {
		request["method"] = v
	}
	if v := plan.Lang.ValueString(); v != "" {
		request["lang"] = v
	}
	if headers, hdiags := headersFromPlan(ctx, plan.Headers); headers != nil {
		request["headers"] = headers
	} else {
		resp.Diagnostics.Append(hdiags...)
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 0)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var e error
		response, e = client.RoaPost(cmsProduct, cmsAPIVersion, action, nil, nil, request, true)
		if e != nil {
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAlertWebhook", alertWebhookTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, alertWebhookTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id, _ := jsonpath.Get("$.alertWebhookId", response)
	if id == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAlertWebhook", alertWebhookTypeName),
			fmt.Sprintf("CreateAlertWebhook succeeded but the response carried no alertWebhookId: %v", response),
		)
		return
	}
	plan.Id = types.StringValue(fmt.Sprint(id))

	object, err := NewService(client).DescribeCmsAlertWebhook(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListAlertWebhooks", alertWebhookTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), "ListAlertWebhooks", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAlertWebhook", alertWebhookTypeName),
			fmt.Sprintf("The alert webhook %s is not visible via ListAlertWebhooks right after creation", plan.Id.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertWebhookState(ctx, object))...)
}

func (r *alertWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertWebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Client()).DescribeCmsAlertWebhook(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListAlertWebhooks", alertWebhookTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "ListAlertWebhooks", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertWebhookState(ctx, object))...)
}

func (r *alertWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state alertWebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	id := state.Id.ValueString()
	action := fmt.Sprintf("/webhook/%s", id)

	// Inherited UpdateAlertWebhook contract: name/url are always carried;
	// contentType/method/lang only when changed and non-empty; headers only
	// when changed. The PATCH is skipped entirely when no field changed.
	request := map[string]interface{}{
		"name": plan.AlertWebhookName.ValueString(),
		"url":  plan.Url.ValueString(),
	}
	update := false
	if !plan.ContentType.Equal(state.ContentType) {
		update = true
		if v := plan.ContentType.ValueString(); v != "" {
			request["contentType"] = v
		}
	}
	if !plan.Method.Equal(state.Method) {
		update = true
		if v := plan.Method.ValueString(); v != "" {
			request["method"] = v
		}
	}
	if !plan.Lang.Equal(state.Lang) {
		update = true
		if v := plan.Lang.ValueString(); v != "" {
			request["lang"] = v
		}
	}
	if !plan.Headers.Equal(state.Headers) {
		update = true
		if headers, hdiags := headersFromPlan(ctx, plan.Headers); headers != nil {
			request["headers"] = headers
		} else {
			resp.Diagnostics.Append(hdiags...)
		}
	}

	if update {
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 0)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			var e error
			response, e = client.RoaPatch(cmsProduct, cmsAPIVersion, action, nil, nil, request, true)
			if e != nil {
				if helper.NeedRetry(e) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateAlertWebhook", alertWebhookTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	object, err := NewService(client).DescribeCmsAlertWebhook(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListAlertWebhooks", alertWebhookTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "ListAlertWebhooks", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertWebhookState(ctx, object))...)
}

func (r *alertWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertWebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	id := state.Id.ValueString()
	// DeleteAlertWebhooks is a batch API: the single identifier is wrapped
	// into the webhookIds query parameter as a JSON array string.
	webhookIdsJson, err := json.Marshal([]string{id})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: building webhookIds", alertWebhookTypeName),
			err.Error(),
		)
		return
	}

	action := "/webhooks"
	query := map[string]*string{
		"webhookIds": ptr(string(webhookIdsJson)),
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 0)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		var e error
		response, e = client.RoaDelete(cmsProduct, cmsAPIVersion, action, query, nil, nil, true)
		if e != nil {
			if isAlertWebhookNotFound(e) {
				return nil
			}
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteAlertWebhooks", alertWebhookTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *alertWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// alertWebhookState maps a ListAlertWebhooks item onto the resource model.
// Empty API strings become null, matching the inherited PR Read behavior so
// optional/computed fields the API omits read back as absent rather than "".
func alertWebhookState(ctx context.Context, object map[string]interface{}) alertWebhookModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return alertWebhookModel{
		Id:               str("webhookId"),
		AlertWebhookName: str("name"),
		Url:              str("url"),
		ContentType:      str("contentType"),
		Method:           str("method"),
		Lang:             str("lang"),
		Workspace:        str("workspace"),
		Headers:          headersToState(ctx, object["headers"]),
	}
}

// headersFromPlan converts a plan/state headers map to the map carried in the
// Create/Update request body. A null or unknown map yields nil so the field is
// omitted, matching the inherited GetOk contract.
func headersFromPlan(ctx context.Context, m types.Map) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	if m.IsNull() || m.IsUnknown() {
		return nil, diags
	}
	var sm map[string]string
	diags.Append(m.ElementsAs(ctx, &sm, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make(map[string]interface{}, len(sm))
	for k, v := range sm {
		out[k] = v
	}
	return out, diags
}

// headersToState converts the headers field of a ListAlertWebhooks item to a
// framework map. A missing/empty block yields a null map.
func headersToState(ctx context.Context, raw interface{}) types.Map {
	m, _ := raw.(map[string]interface{})
	if len(m) == 0 {
		return types.MapNull(types.StringType)
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = fmt.Sprint(v)
	}
	result, diags := types.MapValueFrom(ctx, types.StringType, out)
	_ = diags
	return result
}
