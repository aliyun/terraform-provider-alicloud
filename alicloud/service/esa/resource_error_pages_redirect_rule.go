package esa

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/planmodifiers/int64planmodifier"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

const (
	errorPagesRedirectRuleTypeName = "alicloud_esa_error_pages_redirect_rule"
	errorPagesRedirectRuleTimeout  = 5 * time.Minute
)

var (
	_ resource.Resource                = &errorPagesRedirectRuleResource{}
	_ resource.ResourceWithImportState = &errorPagesRedirectRuleResource{}
)

func NewErrorPagesRedirectRuleResource() resource.Resource {
	return &errorPagesRedirectRuleResource{}
}

type errorPagesRedirectRuleResource struct {
	fwadapt.ResourceBase
}

type errorPagesRedirectRuleResourceModel struct {
	SiteId             types.String `tfsdk:"site_id"`
	RuleEnable         types.String `tfsdk:"rule_enable"`
	RuleName           types.String `tfsdk:"rule_name"`
	Rule               types.String `tfsdk:"rule"`
	SiteVersion        types.Int64  `tfsdk:"site_version"`
	Sequence           types.Int64  `tfsdk:"sequence"`
	ErrorPagesRedirect types.List   `tfsdk:"error_pages_redirect"`
	Id                 types.String `tfsdk:"id"`
	ConfigId           types.Int64  `tfsdk:"config_id"`
	ConfigType         types.String `tfsdk:"config_type"`
}

var errorPagesRedirectObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"status_code": types.StringType,
	"target_url":  types.StringType,
}}

func (r *errorPagesRedirectRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = errorPagesRedirectRuleTypeName
}

func (r *errorPagesRedirectRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Esa ErrorPagesRedirectRule resource.",
		Attributes: map[string]schema.Attribute{
			"site_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the site.",
			},
			"rule_enable": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
				Description: "Whether to enable the rule. Valid values: `on` and `off`.",
			},
			"rule_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the rule.",
			},
			"rule": schema.StringAttribute{
				Optional:    true,
				Description: "The rule content, which is used to match the requests.",
			},
			"site_version": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The version of the site.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"sequence": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The execution order of the rule.",
			},
			"error_pages_redirect": schema.ListAttribute{
				Required:    true,
				ElementType: errorPagesRedirectObjectType,
				Description: "The error page redirect rules. Each element contains `status_code` and `target_url`.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The resource ID in terraform of Error Pages Redirect Rule. It formats as `<site_id>:<config_id>`",
			},
			"config_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the rule configuration.",
			},
			"config_type": schema.StringAttribute{
				Computed:    true,
				Description: "The configuration type.",
			},
		},
	}
}

func (r *errorPagesRedirectRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan errorPagesRedirectRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"SiteId": plan.SiteId.ValueString(),
	}

	if !plan.RuleEnable.IsNull() && !plan.RuleEnable.IsUnknown() {
		request["RuleEnable"] = plan.RuleEnable.ValueString()
	}

	if !plan.RuleName.IsNull() && !plan.RuleName.IsUnknown() {
		request["RuleName"] = plan.RuleName.ValueString()
	}

	if !plan.Rule.IsNull() && !plan.Rule.IsUnknown() {
		request["Rule"] = plan.Rule.ValueString()
	}

	if !plan.SiteVersion.IsNull() && !plan.SiteVersion.IsUnknown() {
		request["SiteVersion"] = plan.SiteVersion.ValueInt64()
	}

	if !plan.Sequence.IsNull() && !plan.Sequence.IsUnknown() {
		request["Sequence"] = plan.Sequence.ValueInt64()
	}

	errorPagesRedirectElements := plan.ErrorPagesRedirect.Elements()
	errorPagesRedirectItems := make([]map[string]interface{}, 0, len(errorPagesRedirectElements))
	errorPagesRedirectPath := path.Root("error_pages_redirect")
	summaryPrefix := fmt.Sprintf("Invalid %s", errorPagesRedirectRuleTypeName)
	if plan.ErrorPagesRedirect.IsNull() || plan.ErrorPagesRedirect.IsUnknown() {
		resp.Diagnostics.AddAttributeError(errorPagesRedirectPath, summaryPrefix,
			"error_pages_redirect is required and must contain at least one element.")
		return
	}

	for i, element := range errorPagesRedirectElements {
		obj, ok := element.(types.Object)
		if !ok {
			continue
		}

		attributes := obj.Attributes()
		statusCode, _ := attributes["status_code"].(types.String)
		targetURL, _ := attributes["target_url"].(types.String)

		code := statusCode.ValueString()
		if statusCode.IsNull() || statusCode.IsUnknown() || code == "" {
			resp.Diagnostics.AddAttributeError(errorPagesRedirectPath.AtListIndex(i).AtName("status_code"), summaryPrefix,
				"status_code is required.")
			return
		}

		if targetURL.IsNull() || targetURL.IsUnknown() || targetURL.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(errorPagesRedirectPath.AtListIndex(i).AtName("target_url"), summaryPrefix,
				"target_url is required.")
			return
		}

		errorPagesRedirectItems = append(errorPagesRedirectItems, map[string]interface{}{
			"StatusCode": code,
			"TargetURL":  targetURL.ValueString(),
		})
	}

	encodedErrorPagesRedirect, jsonErr := json.Marshal(errorPagesRedirectItems)
	if jsonErr != nil {
		resp.Diagnostics.AddAttributeError(errorPagesRedirectPath, summaryPrefix,
			fmt.Sprintf("Serializing error_pages_redirect: %s", jsonErr))
		return
	}

	request["ErrorPagesRedirect"] = string(encodedErrorPagesRedirect)

	action := "CreateErrorPagesRedirectRule"
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(5*time.Second, 3*time.Second)
	err = retry.RetryContext(ctx, errorPagesRedirectRuleTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "LockFailed"}) || helper.NeedRetry(err) {
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
			fmt.Sprintf("Creating %s: calling %s", errorPagesRedirectRuleTypeName, action),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, errorPagesRedirectRuleTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	configId, ok := response["ConfigId"]
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading ConfigId", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("The %s response did not contain a usable ConfigId: %v", action, response),
		)
		return
	}

	id := fmt.Sprintf("%v:%v", request["SiteId"], configId)

	esaServiceV2 := EsaServiceV2{r.Meta()}
	object, err := esaServiceV2.DescribeEsaErrorPagesRedirectRule(id)
	if err != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetErrorPagesRedirectRule", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	state, stateDiags := errorPagesRedirectRuleState(object, id)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *errorPagesRedirectRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state errorPagesRedirectRuleResourceModel
	esaServiceV2 := EsaServiceV2{r.Meta()}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := esaServiceV2.DescribeEsaErrorPagesRedirectRule(id)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetErrorPagesRedirectRule", errs.AlibabaCloudSdkGoERROR).Error(),
		)

		return
	}

	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state, diags := errorPagesRedirectRuleState(object, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *errorPagesRedirectRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state errorPagesRedirectRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: invalid ID", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("Expected an ID of the form <site_id>:<config_id>, got %q.", id),
		)
		return
	}

	action := "UpdateErrorPagesRedirectRule"
	request := map[string]interface{}{
		"SiteId":   parts[0],
		"ConfigId": parts[1],
	}
	update := false

	if !plan.RuleEnable.Equal(state.RuleEnable) {
		update = true
		request["RuleEnable"] = plan.RuleEnable.ValueString()
	}

	if !plan.RuleName.Equal(state.RuleName) {
		update = true
		request["RuleName"] = plan.RuleName.ValueString()
	}

	if !plan.Rule.Equal(state.Rule) {
		update = true
		request["Rule"] = plan.Rule.ValueString()
	}

	if !plan.Sequence.Equal(state.Sequence) && !plan.Sequence.IsNull() && !plan.Sequence.IsUnknown() {
		update = true
		request["Sequence"] = plan.Sequence.ValueInt64()
	}

	if !plan.ErrorPagesRedirect.Equal(state.ErrorPagesRedirect) {
		update = true
	}
	errorPagesRedirectElements := plan.ErrorPagesRedirect.Elements()
	errorPagesRedirectItems := make([]map[string]interface{}, 0, len(errorPagesRedirectElements))
	errorPagesRedirectPath := path.Root("error_pages_redirect")
	summaryPrefix := fmt.Sprintf("Invalid %s", errorPagesRedirectRuleTypeName)
	if plan.ErrorPagesRedirect.IsNull() || plan.ErrorPagesRedirect.IsUnknown() {
		resp.Diagnostics.AddAttributeError(errorPagesRedirectPath, summaryPrefix,
			"error_pages_redirect is required and must contain at least one element.")
		return
	}

	for i, element := range errorPagesRedirectElements {
		obj, ok := element.(types.Object)
		if !ok {
			continue
		}

		attributes := obj.Attributes()
		statusCode, _ := attributes["status_code"].(types.String)
		targetURL, _ := attributes["target_url"].(types.String)

		code := statusCode.ValueString()
		if statusCode.IsNull() || statusCode.IsUnknown() || code == "" {
			resp.Diagnostics.AddAttributeError(errorPagesRedirectPath.AtListIndex(i).AtName("status_code"), summaryPrefix,
				"status_code is required.")
			return
		}

		if targetURL.IsNull() || targetURL.IsUnknown() || targetURL.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(errorPagesRedirectPath.AtListIndex(i).AtName("target_url"), summaryPrefix,
				"target_url is required.")
			return
		}

		errorPagesRedirectItems = append(errorPagesRedirectItems, map[string]interface{}{
			"StatusCode": code,
			"TargetURL":  targetURL.ValueString(),
		})
	}

	encodedErrorPagesRedirect, jsonErr := json.Marshal(errorPagesRedirectItems)
	if jsonErr != nil {
		resp.Diagnostics.AddAttributeError(errorPagesRedirectPath, summaryPrefix,
			fmt.Sprintf("Serializing error_pages_redirect: %s", jsonErr))
		return
	}

	request["ErrorPagesRedirect"] = string(encodedErrorPagesRedirect)

	if update {
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(5*time.Second, 3*time.Second)
		err = retry.RetryContext(ctx, errorPagesRedirectRuleTimeout, func() *retry.RetryError {
			response, err = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "LockFailed"}) || helper.NeedRetry(err) {
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
				fmt.Sprintf("Updating %s: calling %s", errorPagesRedirectRuleTypeName, action),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	esaServiceV2 := EsaServiceV2{r.Meta()}
	object, err := esaServiceV2.DescribeEsaErrorPagesRedirectRule(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetErrorPagesRedirectRule", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, stateDiags := errorPagesRedirectRuleState(object, id)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *errorPagesRedirectRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state errorPagesRedirectRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: invalid ID", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("Expected an ID of the form <site_id>:<config_id>, got %q.", id),
		)
		return
	}

	action := "DeleteErrorPagesRedirectRule"
	request := map[string]interface{}{
		"SiteId":   parts[0],
		"ConfigId": parts[1],
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(5*time.Second, 3*time.Second)
	err = retry.RetryContext(ctx, errorPagesRedirectRuleTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling %s", errorPagesRedirectRuleTypeName, action),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *errorPagesRedirectRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func errorPagesRedirectRuleState(object map[string]interface{}, id string) (errorPagesRedirectRuleResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := errorPagesRedirectRuleResourceModel{
		RuleEnable:  errorPagesRedirectStringValue(object["RuleEnable"]),
		RuleName:    errorPagesRedirectStringValue(object["RuleName"]),
		Rule:        errorPagesRedirectStringValue(object["Rule"]),
		SiteVersion: errorPagesRedirectInt64Value(object["SiteVersion"]),
		Sequence:    errorPagesRedirectInt64Value(object["Sequence"]),
		Id:          types.StringValue(id),
		ConfigId:    errorPagesRedirectInt64Value(object["ConfigId"]),
		ConfigType:  types.StringValue(fmt.Sprint(object["ConfigType"])),
	}

	if parts := strings.Split(id, ":"); len(parts) == 2 {
		model.SiteId = types.StringValue(parts[0])
	}

	items := make([]attr.Value, 0)
	if rawList, ok := object["ErrorPagesRedirect"].([]interface{}); ok {
		for _, raw := range rawList {
			child, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}

			items = append(items, types.ObjectValueMust(errorPagesRedirectObjectType.AttrTypes, map[string]attr.Value{
				"status_code": types.StringValue(convert.Interface2String(child["StatusCode"])),
				"target_url":  types.StringValue(convert.Interface2String(child["TargetURL"])),
			}))
		}
	}

	listValue, listDiags := types.ListValue(errorPagesRedirectObjectType, items)
	diags.Append(listDiags...)
	model.ErrorPagesRedirect = listValue

	return model, diags
}

func errorPagesRedirectInt64Value(raw interface{}) types.Int64 {
	switch v := raw.(type) {
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

func errorPagesRedirectStringValue(raw interface{}) types.String {
	if s, ok := raw.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}
