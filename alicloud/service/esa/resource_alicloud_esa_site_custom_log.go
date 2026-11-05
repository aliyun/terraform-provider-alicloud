package esa

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	siteCustomLogTypeName = "alicloud_esa_site_custom_log"
	// siteCustomLogTimeout replaces the SDKv2 resource-level create/update/delete
	// default timeouts (5 minutes each); terraform-plugin-framework does not
	// support user-configurable resource timeouts yet.
	siteCustomLogTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &siteCustomLogResource{}
	_ resource.ResourceWithImportState = &siteCustomLogResource{}
)

// NewSiteCustomLogResource returns the Framework resource for the ESA
// SiteCustomLog site-level singleton. Registered via conns.ServicePackage.Resources.
func NewSiteCustomLogResource() resource.Resource {
	return &siteCustomLogResource{}
}

type siteCustomLogResource struct {
	fwadapt.ResourceBase
}

type siteCustomLogResourceModel struct {
	Id              types.String `tfsdk:"id"`
	SiteId          types.String `tfsdk:"site_id"`
	RequestHeaders  types.List   `tfsdk:"request_headers"`
	ResponseHeaders types.List   `tfsdk:"response_headers"`
	Cookies         types.List   `tfsdk:"cookies"`
	IsExist         types.Bool   `tfsdk:"is_exist"`
	ConfigId        types.String `tfsdk:"config_id"`
}

func (r *siteCustomLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = siteCustomLogTypeName
}

func (r *siteCustomLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"request_headers": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"response_headers": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"cookies": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"is_exist": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"config_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *siteCustomLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan siteCustomLogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, err := siteCustomLogRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building CreateSiteCustomLog request", siteCustomLogTypeName),
			err.Error(),
		)
		return
	}

	action := "CreateSiteCustomLog"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, siteCustomLogTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
		if callErr != nil {
			if errs.IsExpectedErrors(callErr, []string{"Site.ServiceBusy", "TooManyRequests", "InternalError"}) || helper.NeedRetry(callErr) {
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
			fmt.Sprintf("Creating %s: calling CreateSiteCustomLog", siteCustomLogTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, siteCustomLogTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// SiteCustomLog is a site-level singleton: the Terraform id is the SiteId.
	plan.Id = plan.SiteId

	if found, err := siteCustomLogReadAndRefresh(ctx, r.Meta(), &plan); err != nil {
		// Keep the id in state so the resource stays trackable when read-back fails.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.SiteId)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling GetSiteCustomLog", siteCustomLogTypeName),
			err.Error(),
		)
		return
	} else if !found {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.SiteId)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: GetSiteCustomLog reports the configuration does not exist", siteCustomLogTypeName),
			fmt.Sprintf("CreateSiteCustomLog for SiteId %s succeeded but GetSiteCustomLog returned IsExist=false", plan.SiteId.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *siteCustomLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state siteCustomLogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.ValueString() == "" {
		return
	}

	// site_id is the id; keep it for read-back so a refresh never blanks it.
	state.SiteId = state.Id

	found, err := siteCustomLogReadAndRefresh(ctx, r.Meta(), &state)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSiteCustomLog", siteCustomLogTypeName),
			err.Error(),
		)
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *siteCustomLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan siteCustomLogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the established id across update: site_id is RequiresReplace, so
	// it cannot change without recreation, but the id attribute is Computed.
	var state siteCustomLogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Id.IsNull() || plan.Id.IsUnknown() {
		plan.Id = state.Id
	}

	request, err := siteCustomLogRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: building UpdateSiteCustomLog request", siteCustomLogTypeName),
			err.Error(),
		)
		return
	}

	action := "UpdateSiteCustomLog"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, siteCustomLogTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
		if callErr != nil {
			if errs.IsExpectedErrors(callErr, []string{"Site.ServiceBusy", "TooManyRequests", "InternalError"}) || helper.NeedRetry(callErr) {
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
			fmt.Sprintf("Updating %s: calling UpdateSiteCustomLog", siteCustomLogTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	if found, err := siteCustomLogReadAndRefresh(ctx, r.Meta(), &plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling GetSiteCustomLog", siteCustomLogTypeName),
			err.Error(),
		)
		return
	} else if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete: ESA has no DeleteSiteCustomLog API. Following the established ESA
// no-op delete pattern (resource_alicloud_esa_kv_account.go,
// resource_alicloud_esa_routine_code_deployment.go), Terraform removes the
// resource from state while the server-side configuration is left in place.
func (r *siteCustomLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	log.Printf("[WARN] Cannot destroy resource %s. Terraform will remove this resource from the state file, however the server-side SiteCustomLog configuration may remain.", siteCustomLogTypeName)
}

func (r *siteCustomLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// siteCustomLogRequest builds the Create/Update request body. RequestHeaders,
// ResponseHeaders and Cookies are formData params with style=json carrying a
// JSON array of strings; they are serialized explicitly, matching the ESA
// convention used by resource_alicloud_esa_http_request_header_modification_rule.
func siteCustomLogRequest(plan siteCustomLogResourceModel) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"SiteId": plan.SiteId.ValueString(),
	}
	reqHeaders, err := siteCustomLogJSONList(plan.RequestHeaders)
	if err != nil {
		return nil, fmt.Errorf("serializing request_headers: %w", err)
	}
	request["RequestHeaders"] = reqHeaders
	respHeaders, err := siteCustomLogJSONList(plan.ResponseHeaders)
	if err != nil {
		return nil, fmt.Errorf("serializing response_headers: %w", err)
	}
	request["ResponseHeaders"] = respHeaders
	cookies, err := siteCustomLogJSONList(plan.Cookies)
	if err != nil {
		return nil, fmt.Errorf("serializing cookies: %w", err)
	}
	request["Cookies"] = cookies
	return request, nil
}

// siteCustomLogJSONList serializes a types.List of strings into a JSON array
// string. A null or unknown list is encoded as an empty array so the API
// receives a deterministic, full-configuration request on every apply.
func siteCustomLogJSONList(value types.List) (string, error) {
	if value.IsNull() || value.IsUnknown() {
		return "[]", nil
	}
	values := make([]string, 0, len(value.Elements()))
	for _, el := range value.Elements() {
		s, ok := el.(types.String)
		if !ok {
			return "", fmt.Errorf("unexpected element type %T in list", el)
		}
		values = append(values, s.ValueString())
	}
	b, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// siteCustomLogReadAndRefresh calls GetSiteCustomLog for the model's SiteId and
// repopulates the model from the response. It returns found=false (without
// error) when the custom log configuration does not exist, so the caller can
// remove the resource from state. A genuine RPC or contract error is returned.
func siteCustomLogReadAndRefresh(ctx context.Context, client *connectivity.AliyunClient, m *siteCustomLogResourceModel) (bool, error) {
	siteId := m.SiteId.ValueString()

	action := "GetSiteCustomLog"
	query := map[string]interface{}{"SiteId": siteId}
	request := make(map[string]interface{})

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	var err error
	err = retry.RetryContext(ctx, siteCustomLogTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcGet("ESA", "2024-09-10", action, query, request)
		if callErr != nil {
			if errs.IsExpectedErrors(callErr, []string{"Site.ServiceBusy", "TooManyRequests", "InternalError"}) || helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SiteNotFound"}) || errs.NotFoundError(err) {
			return false, nil
		}
		return false, errs.WrapErrorf(err, errs.DefaultErrorMsg, siteId, action, errs.AlibabaCloudSdkGoERROR)
	}

	isExist, _ := response["IsExist"].(bool)
	if !isExist {
		return false, nil
	}

	logCustomField, ok := response["LogCustomField"].(map[string]interface{})
	if !ok {
		// LogCustomField may be returned as a JSON string; parse it once here so
		// the field extractors below receive a normalized object.
		if raw, strOk := response["LogCustomField"].(string); strOk && raw != "" {
			var parsed map[string]interface{}
			if jsonErr := json.Unmarshal([]byte(raw), &parsed); jsonErr == nil {
				logCustomField = parsed
				ok = true
			}
		}
	}

	var parseErr error
	if m.RequestHeaders, parseErr = siteCustomLogListState(logCustomField, ok, "RequestHeaders"); parseErr != nil {
		return false, parseErr
	}
	if m.ResponseHeaders, parseErr = siteCustomLogListState(logCustomField, ok, "ResponseHeaders"); parseErr != nil {
		return false, parseErr
	}
	if m.Cookies, parseErr = siteCustomLogListState(logCustomField, ok, "Cookies"); parseErr != nil {
		return false, parseErr
	}
	m.IsExist = types.BoolValue(true)
	m.ConfigId = siteCustomLogStringState(response["ConfigId"])
	m.Id = types.StringValue(siteId)
	return true, nil
}

// siteCustomLogListState extracts a string array from a LogCustomField sub-key
// and converts it to a types.List. A missing or empty array maps to a null list
// so an unconfigured attribute stays null and produces no perpetual diff.
func siteCustomLogListState(logCustomField map[string]interface{}, ok bool, key string) (types.List, error) {
	if !ok {
		return types.ListNull(types.StringType), nil
	}
	values, found, err := siteCustomLogStringSlice(logCustomField[key])
	if err != nil {
		return types.List{}, fmt.Errorf("LogCustomField.%s: %w", key, err)
	}
	if !found || len(values) == 0 {
		return types.ListNull(types.StringType), nil
	}
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	return types.ListValueMust(types.StringType, elems), nil
}

// siteCustomLogStringSlice normalizes the documented encodings of a
// RequestHeaders/ResponseHeaders/Cookies field (a JSON string or a native
// []interface{}) into a string slice. Any other shape is treated as contract
// drift and fails the read instead of writing an empty list over state.
func siteCustomLogStringSlice(raw interface{}) ([]string, bool, error) {
	switch v := raw.(type) {
	case nil:
		return nil, false, nil
	case string:
		if v == "" {
			return nil, false, nil
		}
		var parsed interface{}
		if err := json.Unmarshal([]byte(v), &parsed); err != nil {
			return nil, false, fmt.Errorf("not valid JSON: %w", err)
		}
		return siteCustomLogStringSlice(parsed)
	case []interface{}:
		if len(v) == 0 {
			return nil, false, nil
		}
		out := make([]string, 0, len(v))
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, false, fmt.Errorf("[%d] has unexpected type %T", i, item)
			}
			out = append(out, s)
		}
		return out, true, nil
	default:
		return nil, false, fmt.Errorf("has unexpected type %T", raw)
	}
}

// siteCustomLogStringState converts a ConfigId-like scalar to a String state
// value; nil or empty yields null so a missing computed value stays unknown
// until a real value arrives.
func siteCustomLogStringState(raw interface{}) types.String {
	if raw == nil {
		return types.StringNull()
	}
	s := fmt.Sprintf("%v", raw)
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
