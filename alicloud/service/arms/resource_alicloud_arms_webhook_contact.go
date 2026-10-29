package arms

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	webhookContactTypeName = "alicloud_arms_webhook_contact"
	// webhookContactTimeout replaces the SDKv2 resource-level create/delete
	// default timeouts (5 minutes each); terraform-plugin-framework does not
	// support user-configurable resource timeouts yet.
	webhookContactTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &webhookContactResource{}
	_ resource.ResourceWithImportState = &webhookContactResource{}
)

func NewWebhookContactResource() resource.Resource {
	return &webhookContactResource{}
}

type webhookContactResource struct {
	fwadapt.ResourceBase
}

type webhookContactResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	WebhookContactId   types.String `tfsdk:"webhook_contact_id"`
	WebhookContactName types.String `tfsdk:"webhook_contact_name"`
	Url                types.String `tfsdk:"url"`
	Method             types.String `tfsdk:"method"`
	Body               types.String `tfsdk:"body"`
	RecoverBody        types.String `tfsdk:"recover_body"`
	BizHeaders         types.Map    `tfsdk:"biz_headers"`
	BizParams          types.Map    `tfsdk:"biz_params"`
}

func (r *webhookContactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = webhookContactTypeName
}

func (r *webhookContactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"webhook_contact_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"webhook_contact_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"method": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Get", "Post"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"body": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"recover_body": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"biz_headers": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"biz_params": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *webhookContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookContactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// CreateOrUpdateWebhookContact requires both templates when method is Post.
	// Check them before building the request so a missing companion value is
	// reported locally instead of being sent to the API.
	if plan.Method.ValueString() == "Post" {
		if plan.Body.IsNull() || plan.Body.ValueString() == "" {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: invalid body", webhookContactTypeName),
				"body must be set to a non-empty value when method is Post",
			)
		}
		if plan.RecoverBody.IsNull() || plan.RecoverBody.ValueString() == "" {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: invalid recover_body", webhookContactTypeName),
				"recover_body must be set to a non-empty value when method is Post",
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	request, err := webhookContactCreateRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building CreateOrUpdateWebhookContact request", webhookContactTypeName),
			err.Error(),
		)
		return
	}

	action := "CreateOrUpdateWebhookContact"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, webhookContactTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Client().RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
			fmt.Sprintf("Creating %s: calling CreateOrUpdateWebhookContact", webhookContactTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, webhookContactTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	rawId, err := jsonpath.Get("$.WebhookContact.WebhookId", response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading $.WebhookContact.WebhookId", webhookContactTypeName),
			errs.WrapErrorf(err, errs.FailedGetAttributeMsg, webhookContactTypeName, "$.WebhookContact.WebhookId", response).Error(),
		)
		return
	}
	id, err := WebhookContactId(rawId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading $.WebhookContact.WebhookId", webhookContactTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewService(r.Client()).DescribeWebhookContactContext(ctx, id)
	if err != nil {
		// Keep the created ID in state so the contact stays trackable even when
		// the read-back fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling DescribeWebhookContacts", webhookContactTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeWebhookContacts", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := webhookContactState(object)
	if err != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping DescribeWebhookContacts response", webhookContactTypeName),
			err.Error(),
		)
		return
	}
	webhookContactPreserveEmptyConfigured(&model, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// Update exists because the Framework resource interface requires it. Every
// configurable argument is replacement-only, so Terraform replaces the contact
// instead of calling this method.
func (r *webhookContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookContactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *webhookContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookContactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Client()).DescribeWebhookContactContext(ctx, id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeWebhookContacts", webhookContactTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeWebhookContacts", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// Build and validate the complete model before writing state: a
	// contract-drifted response must fail the refresh instead of replacing the
	// configured ForceNew values with empty ones.
	model, err := webhookContactState(object)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping DescribeWebhookContacts response", webhookContactTypeName),
			err.Error(),
		)
		return
	}
	webhookContactPreserveEmptyConfigured(&model, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *webhookContactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookContactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if err := NewService(r.Client()).DeleteWebhookContactContext(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteWebhookContact", webhookContactTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteWebhookContact", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *webhookContactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func webhookContactCreateRequest(plan webhookContactResourceModel) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"Method":      plan.Method.ValueString(),
		"Url":         plan.Url.ValueString(),
		"WebhookName": plan.WebhookContactName.ValueString(),
	}
	if !plan.Body.IsNull() && plan.Body.ValueString() != "" {
		request["Body"] = plan.Body.ValueString()
	}
	if !plan.RecoverBody.IsNull() && plan.RecoverBody.ValueString() != "" {
		request["RecoverBody"] = plan.RecoverBody.ValueString()
	}
	if !plan.BizHeaders.IsNull() && len(plan.BizHeaders.Elements()) > 0 {
		headers, err := webhookContactKvMapToJSON(plan.BizHeaders)
		if err != nil {
			return nil, fmt.Errorf("serializing biz_headers: %w", err)
		}
		request["BizHeaders"] = headers
	}
	if !plan.BizParams.IsNull() && len(plan.BizParams.Elements()) > 0 {
		params, err := webhookContactKvMapToJSON(plan.BizParams)
		if err != nil {
			return nil, fmt.Errorf("serializing biz_params: %w", err)
		}
		request["BizParams"] = params
	}
	return request, nil
}

// The published ARMS API models BizHeaders/BizParams as string params carrying
// a JSON array of single-entry objects (e.g. [{"Content-Type":"application/json"}]);
// sending them as plain objects produces dotted form keys the API does not
// accept, so they are serialized explicitly here.
func webhookContactKvMapToJSON(value types.Map) (string, error) {
	elems := value.Elements()
	keys := make([]string, 0, len(elems))
	for k := range elems {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	kvList := make([]map[string]interface{}, 0, len(keys))
	for _, k := range keys {
		s, ok := elems[k].(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			return "", fmt.Errorf("unexpected value type %T for key %q", elems[k], k)
		}
		kvList = append(kvList, map[string]interface{}{k: s.ValueString()})
	}
	b, err := json.Marshal(kvList)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// webhookContactKvMapFromResponse normalizes the two documented encodings of
// BizHeaders/BizParams (a JSON string or an array of single-entry objects) into
// a flat map. Any other shape is treated as contract drift and fails the read
// instead of writing an empty map over configured ForceNew values.
func webhookContactKvMapFromResponse(field string, raw interface{}) (map[string]string, bool, error) {
	switch v := raw.(type) {
	case nil:
		return nil, false, nil
	case string:
		if v == "" {
			return nil, false, nil
		}
		var parsed interface{}
		if err := json.Unmarshal([]byte(v), &parsed); err != nil {
			return nil, false, fmt.Errorf("%s is not valid JSON: %w", field, err)
		}
		return webhookContactKvMapFromResponse(field, parsed)
	case map[string]interface{}:
		kv, err := webhookContactKvMapFromObject(field, v)
		if err != nil {
			return nil, false, err
		}
		if len(kv) == 0 {
			return nil, false, nil
		}
		return kv, true, nil
	case []interface{}:
		result := make(map[string]string, len(v))
		for i, item := range v {
			m, ok := item.(map[string]interface{})
			if !ok {
				return nil, false, fmt.Errorf("%s[%d] has unexpected type %T", field, i, item)
			}
			kv, err := webhookContactKvMapFromObject(field, m)
			if err != nil {
				return nil, false, err
			}
			for k, val := range kv {
				result[k] = val
			}
		}
		if len(result) == 0 {
			return nil, false, nil
		}
		return result, true, nil
	default:
		return nil, false, fmt.Errorf("%s has unexpected type %T", field, raw)
	}
}

func webhookContactKvMapFromObject(field string, m map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string, len(m))
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%q] has unexpected value type %T", field, k, v)
		}
		result[k] = s
	}
	return result, nil
}

func webhookContactState(object map[string]interface{}) (webhookContactResourceModel, error) {
	webhookRaw, ok := object["Webhook"].(map[string]interface{})
	if !ok {
		return webhookContactResourceModel{}, fmt.Errorf("unexpected Webhook response type %T", object["Webhook"])
	}
	id, err := WebhookContactId(object["WebhookId"])
	if err != nil {
		return webhookContactResourceModel{}, err
	}
	name, err := webhookContactRequiredString(object["WebhookName"], "WebhookName")
	if err != nil {
		return webhookContactResourceModel{}, err
	}
	url, err := webhookContactRequiredString(webhookRaw["Url"], "Url")
	if err != nil {
		return webhookContactResourceModel{}, err
	}
	method, err := webhookContactRequiredString(webhookRaw["Method"], "Method")
	if err != nil {
		return webhookContactResourceModel{}, err
	}
	headers, headersFound, err := webhookContactKvMapFromResponse("BizHeaders", webhookRaw["BizHeaders"])
	if err != nil {
		return webhookContactResourceModel{}, err
	}
	params, paramsFound, err := webhookContactKvMapFromResponse("BizParams", webhookRaw["BizParams"])
	if err != nil {
		return webhookContactResourceModel{}, err
	}

	return webhookContactResourceModel{
		Id:                 types.StringValue(id),
		WebhookContactId:   types.StringValue(id),
		WebhookContactName: name,
		Url:                url,
		Method:             method,
		Body:               webhookContactString(webhookRaw["Body"]),
		RecoverBody:        webhookContactString(webhookRaw["RecoverBody"]),
		BizHeaders:         webhookContactStringMap(headers, headersFound),
		BizParams:          webhookContactStringMap(params, paramsFound),
	}, nil
}

// webhookContactRequiredString rejects a missing or non-string value for a
// required attribute instead of silently writing null over configured state.
func webhookContactRequiredString(raw interface{}, field string) (types.String, error) {
	v, ok := raw.(string)
	if !ok || v == "" {
		return types.StringNull(), fmt.Errorf("%s is required in the DescribeWebhookContacts response, got %T %v", field, raw, raw)
	}
	return types.StringValue(v), nil
}

// webhookContactPreserveEmptyConfigured keeps an explicitly empty configured
// value in state. CreateOrUpdateWebhookContact omits body, recover_body,
// biz_headers and biz_params when they are empty, so normalizing them to null
// would fail the post-apply consistency check and plan a replacement on every
// refresh for these replacement-only attributes.
func webhookContactPreserveEmptyConfigured(model *webhookContactResourceModel, source webhookContactResourceModel) {
	if model.Body.IsNull() && !source.Body.IsNull() && !source.Body.IsUnknown() && source.Body.ValueString() == "" {
		model.Body = source.Body
	}
	if model.RecoverBody.IsNull() && !source.RecoverBody.IsNull() && !source.RecoverBody.IsUnknown() && source.RecoverBody.ValueString() == "" {
		model.RecoverBody = source.RecoverBody
	}
	if model.BizHeaders.IsNull() && !source.BizHeaders.IsNull() && !source.BizHeaders.IsUnknown() && len(source.BizHeaders.Elements()) == 0 {
		model.BizHeaders = source.BizHeaders
	}
	if model.BizParams.IsNull() && !source.BizParams.IsNull() && !source.BizParams.IsUnknown() && len(source.BizParams.Elements()) == 0 {
		model.BizParams = source.BizParams
	}
}

func webhookContactString(raw interface{}) types.String {
	if v, ok := raw.(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func webhookContactStringMap(values map[string]string, found bool) types.Map {
	if !found {
		return types.MapNull(types.StringType)
	}
	elems := make(map[string]attr.Value, len(values))
	for k, v := range values {
		elems[k] = types.StringValue(v)
	}
	return types.MapValueMust(types.StringType, elems)
}
