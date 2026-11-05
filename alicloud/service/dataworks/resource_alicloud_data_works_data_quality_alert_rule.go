package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

const dataQualityAlertRuleTypeName = "alicloud_data_works_data_quality_alert_rule"

// dataQualityAlertRuleTimeout replaces the SDKv2 resource-level default
// timeouts (Create/Update/Delete each 5 minutes); terraform-plugin-framework
// does not support user-configurable resource timeouts yet.
const dataQualityAlertRuleTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &dataQualityAlertRuleResource{}
	_ resource.ResourceWithImportState = &dataQualityAlertRuleResource{}
)

func NewDataQualityAlertRuleResource() resource.Resource {
	return &dataQualityAlertRuleResource{}
}

type dataQualityAlertRuleResource struct {
	fwadapt.ResourceBase
}

type dataQualityAlertRuleResourceModel struct {
	Id           types.String      `tfsdk:"id"`
	ProjectId    types.Int64       `tfsdk:"project_id"`
	Condition    types.String      `tfsdk:"condition"`
	Target       targetModel       `tfsdk:"target"`
	Notification notificationModel `tfsdk:"notification"`
}

type targetModel struct {
	Type types.String `tfsdk:"type"`
	Ids  types.List   `tfsdk:"ids"`
}

type notificationModel struct {
	Channels  types.List      `tfsdk:"channels"`
	Receivers []receiverModel `tfsdk:"receivers"`
}

type receiverModel struct {
	ReceiverType   types.String `tfsdk:"receiver_type"`
	ReceiverValues types.List   `tfsdk:"receiver_values"`
	Extension      types.String `tfsdk:"extension"`
}

func (r *dataQualityAlertRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = dataQualityAlertRuleTypeName
}

func (r *dataQualityAlertRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"condition": schema.StringAttribute{
				Required: true,
			},
			"target": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("DataQualityScan"),
						},
					},
					"ids": schema.ListAttribute{
						Required:    true,
						ElementType: types.Int64Type,
						// The API currently only accepts a single monitoring
						// target id, so the list is constrained to one entry.
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.SizeAtMost(1),
						},
					},
				},
			},
			"notification": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"channels": schema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(
								stringvalidator.OneOf("Dingding", "Weixin", "Mail", "Sms", "Phone", "Webhook", "Feishu"),
							),
						},
					},
					"receivers": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"receiver_type": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											"AliUid",
											"WebhookUrl",
											"DingdingUrl",
											"WeixinUrl",
											"FeishuUrl",
											"TaskOwner",
											"DataQualityScanOwner",
											"ShiftSchedule",
										),
									},
								},
								"receiver_values": schema.ListAttribute{
									Optional:    true,
									ElementType: types.StringType,
								},
								"extension": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *dataQualityAlertRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataQualityAlertRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreateDataQualityAlertRule"
	request, diags := buildDataQualityAlertRuleRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request["RegionId"] = r.Meta().RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataQualityAlertRuleTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
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
			fmt.Sprintf("Creating %s: calling %s", dataQualityAlertRuleTypeName, action),
			errs.WrapErrorf(err, "Creating %s: calling %s", dataQualityAlertRuleTypeName, action).Error(),
		)
		return
	}

	// response["Id"] is decoded by the tea-rpc client as json.Number or
	// float64; format it back to a plain decimal string for state and the
	// Get/Delete calls (which accept the long value as a numeric string).
	id := strconv.FormatInt(toInt64Value(response["Id"]), 10)
	object, err := NewService(r.Meta()).DescribeDataQualityAlertRule(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDataQualityAlertRule", dataQualityAlertRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDataQualityAlertRule", dataQualityAlertRuleTypeName).Error(),
		)
		return
	}
	state := dataQualityAlertRuleState(object, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *dataQualityAlertRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dataQualityAlertRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeDataQualityAlertRule(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDataQualityAlertRule", dataQualityAlertRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDataQualityAlertRule", dataQualityAlertRuleTypeName).Error(),
		)
		return
	}
	model := dataQualityAlertRuleState(object, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataQualityAlertRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dataQualityAlertRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state dataQualityAlertRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "UpdateDataQualityAlertRule"
	request, diags := buildDataQualityAlertRuleRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request["Id"] = state.Id.ValueString()
	request["RegionId"] = r.Meta().RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataQualityAlertRuleTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
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
			fmt.Sprintf("Updating %s: calling %s", dataQualityAlertRuleTypeName, action),
			errs.WrapErrorf(err, "Updating %s: calling %s", dataQualityAlertRuleTypeName, action).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDataQualityAlertRule(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDataQualityAlertRule", dataQualityAlertRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDataQualityAlertRule", dataQualityAlertRuleTypeName).Error(),
		)
		return
	}
	model := dataQualityAlertRuleState(object, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataQualityAlertRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dataQualityAlertRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "DeleteDataQualityAlertRule"
	request := map[string]interface{}{
		"Id":          state.Id.ValueString(),
		"RegionId":    r.Meta().RegionId,
		"ClientToken": helper.BuildClientToken(action),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataQualityAlertRuleTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
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
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling %s", dataQualityAlertRuleTypeName, action),
			errs.WrapErrorf(err, "Deleting %s: calling %s", dataQualityAlertRuleTypeName, action).Error(),
		)
		return
	}
}

func (r *dataQualityAlertRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildDataQualityAlertRuleRequest assembles the Create/Update request body.
// Target and Notification are JSON-style form parameters: the structured
// attributes are marshaled to JSON strings before being placed on the request.
func buildDataQualityAlertRuleRequest(ctx context.Context, plan dataQualityAlertRuleResourceModel) (map[string]interface{}, diag.Diagnostics) {
	request := map[string]interface{}{
		"ProjectId": plan.ProjectId.ValueInt64(),
		"Condition": plan.Condition.ValueString(),
	}
	diags := diag.Diagnostics{}

	// Target -> {"Type": "...", "Ids": [...]}
	ids := make([]int64, 0)
	diags.Append(plan.Target.Ids.ElementsAs(ctx, &ids, false)...)
	if diags.HasError() {
		return request, diags
	}
	targetMap := map[string]interface{}{
		"Type": plan.Target.Type.ValueString(),
		"Ids":  ids,
	}
	targetJson, err := json.Marshal(targetMap)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Building %s request: encoding Target", dataQualityAlertRuleTypeName),
			errs.WrapError(err).Error(),
		)
		return request, diags
	}
	request["Target"] = string(targetJson)

	// Notification -> {"Channels": [...], "Receivers": [...]}
	channels := make([]string, 0)
	diags.Append(plan.Notification.Channels.ElementsAs(ctx, &channels, false)...)
	if diags.HasError() {
		return request, diags
	}
	receiversArr := make([]interface{}, 0, len(plan.Notification.Receivers))
	for _, rcv := range plan.Notification.Receivers {
		receiverMap := map[string]interface{}{
			"ReceiverType": rcv.ReceiverType.ValueString(),
		}
		values := make([]string, 0)
		diags.Append(rcv.ReceiverValues.ElementsAs(ctx, &values, false)...)
		if diags.HasError() {
			return request, diags
		}
		receiverMap["ReceiverValues"] = values
		if ext := rcv.Extension.ValueString(); ext != "" {
			receiverMap["Extension"] = ext
		}
		receiversArr = append(receiversArr, receiverMap)
	}
	notificationMap := map[string]interface{}{
		"Channels":  channels,
		"Receivers": receiversArr,
	}
	notificationJson, err := json.Marshal(notificationMap)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Building %s request: encoding Notification", dataQualityAlertRuleTypeName),
			errs.WrapError(err).Error(),
		)
		return request, diags
	}
	request["Notification"] = string(notificationJson)

	return request, diags
}

// dataQualityAlertRuleState maps a GetDataQualityAlertRule object onto the
// resource model. The tea-rpc client decodes API numbers as json.Number or
// float64 rather than int, so every numeric encoding the transport produces
// must map back to Int64/String or planned values would collapse to null
// ("inconsistent result after apply").
//
// The reference model (plan on Create/Update, prior state on Read) supplies
// the user's intended ordering for notification.channels and the
// notification.receivers list. The API does not document a stable return
// order for these collections, so the response is reordered to match the
// reference before being stored; otherwise a Framework List value comparison
// would surface a spurious reordering diff (permadiff) on every plan when the
// server returns the same elements in a different sequence (R4.19). Elements
// present in the response but absent from the reference are appended, so no
// server-side value is silently dropped.
func dataQualityAlertRuleState(object map[string]interface{}, ref dataQualityAlertRuleResourceModel) dataQualityAlertRuleResourceModel {
	str := func(key string) types.String {
		switch v := object[key].(type) {
		case string:
			if v != "" {
				return types.StringValue(v)
			}
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
	// idStr handles the rule identifier which the API returns as a number;
	// it is stored in Terraform state as a string.
	idStr := func(key string) types.String {
		switch v := object[key].(type) {
		case string:
			if v != "" {
				return types.StringValue(v)
			}
		case json.Number:
			return types.StringValue(v.String())
		case float64:
			return types.StringValue(strconv.FormatInt(int64(v), 10))
		case int:
			return types.StringValue(strconv.Itoa(v))
		case int64:
			return types.StringValue(strconv.FormatInt(v, 10))
		}
		return types.StringNull()
	}
	int64List := func(raw []interface{}) types.List {
		elems := make([]attr.Value, 0, len(raw))
		for _, e := range raw {
			elems = append(elems, types.Int64Value(toInt64Value(e)))
		}
		list, _ := types.ListValue(types.Int64Type, elems)
		return list
	}
	stringList := func(raw []interface{}) types.List {
		elems := make([]attr.Value, 0, len(raw))
		for _, e := range raw {
			if s, ok := e.(string); ok {
				elems = append(elems, types.StringValue(s))
			}
		}
		list, _ := types.ListValue(types.StringType, elems)
		return list
	}

	model := dataQualityAlertRuleResourceModel{
		Id:        idStr("Id"),
		ProjectId: int64Val("ProjectId"),
		Condition: str("Condition"),
		Target: targetModel{
			Type: types.StringNull(),
			Ids:  types.ListNull(types.Int64Type),
		},
		Notification: notificationModel{
			Channels:  types.ListNull(types.StringType),
			Receivers: nil,
		},
	}

	if t, ok := object["Target"].(map[string]interface{}); ok {
		if v, ok := t["Type"].(string); ok && v != "" {
			model.Target.Type = types.StringValue(v)
		}
		if raw, ok := t["Ids"].([]interface{}); ok && len(raw) > 0 {
			model.Target.Ids = int64List(raw)
		}
	}

	if n, ok := object["Notification"].(map[string]interface{}); ok {
		if raw, ok := n["Channels"].([]interface{}); ok && len(raw) > 0 {
			apiChannels := toStringSlice(raw)
			refChannels := listStringSlice(ref.Notification.Channels)
			model.Notification.Channels = stringList(toInterfaceSlice(normalizeStringOrder(apiChannels, refChannels)))
		}
		if rawReceivers, ok := n["Receivers"].([]interface{}); ok {
			ordered := normalizeReceiverOrder(rawReceivers, ref.Notification.Receivers)
			for _, rRaw := range ordered {
				r, ok := rRaw.(map[string]interface{})
				if !ok {
					continue
				}
				rm := receiverModel{
					ReceiverType:   nestedStr(r, "ReceiverType"),
					ReceiverValues: types.ListNull(types.StringType),
					Extension:      nestedStr(r, "Extension"),
				}
				if rawVals, ok := r["ReceiverValues"].([]interface{}); ok && len(rawVals) > 0 {
					rm.ReceiverValues = stringList(rawVals)
				}
				model.Notification.Receivers = append(model.Notification.Receivers, rm)
			}
		}
	}

	return model
}

// nestedStr reads a string field from a nested object map.
func nestedStr(object map[string]interface{}, key string) types.String {
	switch v := object[key].(type) {
	case string:
		if v != "" {
			return types.StringValue(v)
		}
	}
	return types.StringNull()
}

// toInt64Value converts a tea-rpc decoded scalar to an int64 value.
func toInt64Value(v interface{}) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	case string:
		i, _ := strconv.ParseInt(n, 10, 64)
		return i
	}
	return 0
}

// toStringSlice converts a tea-rpc decoded []interface{} of strings.
func toStringSlice(raw []interface{}) []string {
	out := make([]string, 0, len(raw))
	for _, e := range raw {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// listStringSlice reads the string elements of a Framework types.List.
func listStringSlice(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	out := make([]string, 0, len(list.Elements()))
	for _, e := range list.Elements() {
		if v, ok := e.(types.String); ok && !v.IsNull() && !v.IsUnknown() {
			out = append(out, v.ValueString())
		}
	}
	return out
}

// toInterfaceSlice converts a string slice back to []interface{} for the
// shared stringList builder.
func toInterfaceSlice(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}

// normalizeStringOrder reorders api so that elements matching the reference
// order come first (in reference order), followed by any extra api elements
// not present in the reference. Elements are matched by value; duplicates in
// the reference are ignored beyond their first occurrence. This keeps state
// stable across reads when the server returns the same set in a different
// order (R4.19), without dropping server-side values the user did not
// configure.
func normalizeStringOrder(api, ref []string) []string {
	if len(api) <= 1 || len(ref) == 0 {
		return api
	}
	seen := make(map[string]bool, len(api))
	for _, s := range api {
		seen[s] = true
	}
	ordered := make([]string, 0, len(api))
	consumed := make(map[string]bool, len(ref))
	for _, s := range ref {
		if seen[s] && !consumed[s] {
			ordered = append(ordered, s)
			consumed[s] = true
		}
	}
	for _, s := range api {
		if !consumed[s] {
			ordered = append(ordered, s)
			consumed[s] = true
		}
	}
	return ordered
}

// receiverKey is the stable composite key used to pair a server-returned
// receiver with the user-configured reference receiver. receiver_type alone
// is not unique (a rule may carry several receivers of the same type), so the
// receiver_values and extension are folded in.
func receiverKey(receiverType string, values []string, extension string) string {
	return receiverType + "|" + strings.Join(values, ",") + "|" + extension
}

// serverReceiverKey computes the composite key for a server-returned receiver
// object (map[string]interface{}).
func serverReceiverKey(r map[string]interface{}) string {
	rt := ""
	if v, ok := r["ReceiverType"].(string); ok {
		rt = v
	}
	values := make([]string, 0)
	if raw, ok := r["ReceiverValues"].([]interface{}); ok {
		for _, e := range raw {
			if s, ok := e.(string); ok {
				values = append(values, s)
			}
		}
	}
	ext := ""
	if v, ok := r["Extension"].(string); ok {
		ext = v
	}
	return receiverKey(rt, values, ext)
}

// refReceiverKey computes the composite key for a user-configured reference
// receiver model.
func refReceiverKey(r receiverModel) string {
	rt := ""
	if !r.ReceiverType.IsNull() && !r.ReceiverType.IsUnknown() {
		rt = r.ReceiverType.ValueString()
	}
	values := listStringSlice(r.ReceiverValues)
	ext := ""
	if !r.Extension.IsNull() && !r.Extension.IsUnknown() {
		ext = r.Extension.ValueString()
	}
	return receiverKey(rt, values, ext)
}

// normalizeReceiverOrder reorders the server-returned receivers so that those
// matching the reference order come first (in reference order), followed by
// any extra receivers not present in the reference. Pairs are matched by the
// receiver composite key; an unmatched server receiver is appended so no
// server-side value is silently dropped (R4.19).
func normalizeReceiverOrder(api []interface{}, ref []receiverModel) []interface{} {
	if len(api) <= 1 || len(ref) == 0 {
		return api
	}
	byKey := make(map[string]interface{}, len(api))
	for _, rRaw := range api {
		r, ok := rRaw.(map[string]interface{})
		if !ok {
			continue
		}
		byKey[serverReceiverKey(r)] = rRaw
	}
	ordered := make([]interface{}, 0, len(api))
	consumed := make(map[string]bool, len(api))
	for _, refR := range ref {
		key := refReceiverKey(refR)
		if rRaw, ok := byKey[key]; ok && !consumed[key] {
			ordered = append(ordered, rRaw)
			consumed[key] = true
		}
	}
	for _, rRaw := range api {
		r, ok := rRaw.(map[string]interface{})
		if !ok {
			if !containsValue(ordered, rRaw) {
				ordered = append(ordered, rRaw)
			}
			continue
		}
		key := serverReceiverKey(r)
		if !consumed[key] {
			ordered = append(ordered, rRaw)
			consumed[key] = true
		}
	}
	return ordered
}

// containsValue reports whether v is present in the slice.
func containsValue(slice []interface{}, v interface{}) bool {
	for _, e := range slice {
		if e == v {
			return true
		}
	}
	return false
}
