// Package tag. This file is generated automatically. Please do not modify it manually, thank you!
// Service helpers for the framework-native alicloud_tag_created_by resource.
// Shared client-token and incremental-wait helpers are imported from alicloud/helper, and the
// framework wire/state adapters from alicloud/provider/fwadapt; this file keeps the tag request,
// retry-classification and tag-specific mapping helpers.
package tag

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type fwDeleteCondition struct {
	Property, Action string
	Values           []string
}

type fwScopedReadIdentity struct{ Source, Target, Value string }

// Only attributes with an explicit JSON document contract use this string type.
type fwJSONStringType struct{ basetypes.StringType }
type fwJSONStringValue struct{ basetypes.StringValue }

func (t fwJSONStringType) String() string                       { return "fwJSONStringType" }
func (t fwJSONStringType) Equal(other attr.Type) bool           { _, ok := other.(fwJSONStringType); return ok }
func (t fwJSONStringType) ValueType(context.Context) attr.Value { return fwJSONStringValue{} }
func (t fwJSONStringType) ValueFromString(_ context.Context, value basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return fwJSONStringValue{StringValue: value}, nil
}
func (t fwJSONStringType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	converted, err := t.StringType.ValueFromTerraform(ctx, value)
	if err != nil {
		return nil, err
	}
	return fwJSONStringValue{StringValue: converted.(basetypes.StringValue)}, nil
}

func fwJSONString(value types.String) fwJSONStringValue    { return fwJSONStringValue{StringValue: value} }
func (v fwJSONStringValue) Type(context.Context) attr.Type { return fwJSONStringType{} }
func (v fwJSONStringValue) Equal(other attr.Value) bool {
	value, ok := other.(fwJSONStringValue)
	return ok && v.StringValue.Equal(value.StringValue)
}
func (v fwJSONStringValue) NullValue() attr.Value { return fwJSONString(types.StringNull()) }
func (v fwJSONStringValue) SemanticEqual(other attr.Value) bool {
	value, ok := other.(fwJSONStringValue)
	if !ok || v.IsNull() || v.IsUnknown() || value.IsNull() || value.IsUnknown() {
		return false
	}
	equal, err := fwJSONDocumentsEqual(v.ValueString(), value.ValueString())
	return err == nil && equal
}

// The shared adapters assert these exported contracts; keep the compile-time proof here so a
// rename cannot silently drop the tag type back to plain equality or a non-typed null.
var (
	_ fwadapt.SemanticValue = fwJSONStringValue{}
	_ fwadapt.NullValue     = fwJSONStringValue{}
)

func (v fwJSONStringValue) StringSemanticEquals(_ context.Context, other basetypes.StringValuable) (bool, diag.Diagnostics) {
	value, ok := other.(fwJSONStringValue)
	if !ok || v.IsNull() || v.IsUnknown() || value.IsNull() || value.IsUnknown() {
		return false, nil
	}
	equal, err := fwJSONDocumentsEqual(v.ValueString(), value.ValueString())
	if err != nil {
		return false, diag.Diagnostics{diag.NewErrorDiagnostic("Invalid JSON document", "A JSON document could not be compared. The document content is omitted.")}
	}
	return equal, nil
}
func (v fwJSONStringValue) ValidateAttribute(_ context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	if !json.Valid([]byte(v.ValueString())) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON document", "This attribute requires a valid JSON document. The document content is omitted.")
	}
}

// Canonicalize exact decimal values without float64 rounding or expanding powers of ten.
func fwJSONCanonicalDocument(text string) (string, error) {
	if !json.Valid([]byte(text)) {
		return "", fmt.Errorf("invalid JSON document")
	}
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var value interface{}
	if err := decoder.Decode(&value); err != nil {
		return "", fmt.Errorf("invalid JSON document")
	}
	number := func(value json.Number) json.Number {
		literal, sign := value.String(), ""
		if strings.HasPrefix(literal, "-") {
			sign, literal = "-", literal[1:]
		}
		exponent := new(big.Int)
		if index := strings.IndexAny(literal, "eE"); index >= 0 {
			exponent.SetString(literal[index+1:], 10)
			literal = literal[:index]
		}
		if index := strings.IndexByte(literal, '.'); index >= 0 {
			exponent.Sub(exponent, big.NewInt(int64(len(literal)-index-1)))
			literal = literal[:index] + literal[index+1:]
		}
		literal = strings.TrimLeft(literal, "0")
		if literal == "" {
			return json.Number("0")
		}
		digits := strings.TrimRight(literal, "0")
		exponent.Add(exponent, big.NewInt(int64(len(literal)-len(digits))))
		return json.Number(sign + digits + "e" + exponent.String())
	}
	var normalize func(interface{}) interface{}
	normalize = func(current interface{}) interface{} {
		switch current := current.(type) {
		case json.Number:
			return number(current)
		case []interface{}:
			for index, child := range current {
				current[index] = normalize(child)
			}
		case map[string]interface{}:
			for key, child := range current {
				current[key] = normalize(child)
			}
		}
		return current
	}
	encoded, err := json.Marshal(normalize(value))
	if err != nil {
		return "", fmt.Errorf("invalid JSON document")
	}
	return string(encoded), nil
}
func fwJSONDocumentsEqual(left, right string) (bool, error) {
	a, err := fwJSONCanonicalDocument(left)
	if err != nil {
		return false, err
	}
	b, err := fwJSONCanonicalDocument(right)
	if err != nil {
		return false, err
	}
	return a == b, nil
}

func fwTimeoutTypes() map[string]attr.Type {
	return fwutil.TimeoutTypes("create", "update", "delete")
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	return fwutil.TimeoutBlock("create", "update", "delete")
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	return fwutil.TimeoutDuration(value, operation, fallback)
}

// retryRpcPost mirrors alicloud's retry.RetryContext+NeedRetry wrapper around RpcPost.
// An explicit endpoint pins the call to a control region that differs from the
// provider region; an empty endpoint keeps the client's regional resolution.
// ctx is the lifecycle context, so cancelling the operation stops the retries.
func retryRpcPost(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration, endpoint string) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		if endpoint == "" {
			return client.RpcPost(product, version, action, nil, request, true)
		}
		return client.RpcPostWithEndpoint(product, version, action, nil, request, true, endpoint)
	}, timeout, needRetry)
}

// needRetry mirrors alicloud.NeedRetry for the *tea.SDKError path that RpcPost produces.
func needRetry(err error) bool {
	if err == nil {
		return false
	}
	if strings.HasPrefix(err.Error(), "Post https://") {
		return true
	}
	if e, ok := err.(*tea.SDKError); ok {
		if e.Code == nil || e.Message == nil {
			return false
		}
		code := *e.Code
		msg := *e.Message
		if strings.Contains(msg, "code: 500, 您已开通过") {
			return false
		}
		if strings.Contains(msg, "Client.Timeout") {
			return true
		}
		if code == "ServiceUnavailable" || code == "Rejected.Throttling" || strings.Contains(code, "Throttling") || strings.HasPrefix(msg, "code: 5") {
			return true
		}
		// gap#5: v1 retries delete/read per-call on conflict/status/busy codes
		// (OperationConflict, LastTokenProcessing, IncorrectStatus, SystemBusy) and on
		// DependencyViolation (deps clear on retry, e.g. DeleteHaVip returning
		// 400 DependencyViolation.HaVip). The v2 generator uses a single global needRetry for
		// all RpcPost calls, so these codes never retried and Delete failed fast on
		// DependencyViolation. Including them here mirrors v1; all are transient and bounded
		// by the call's retry timeout. DependencyViolation carries a product suffix
		// (.HaVip/.GatewayRouteTable/...), so match by prefix.
		if strings.HasPrefix(code, "DependencyViolation") || code == "OperationConflict" ||
			code == "LastTokenProcessing" || code == "IncorrectStatus" || code == "SystemBusy" {
			return true
		}
	}
	return false
}

// isExpectedErrors matches structured API codes, including wrapped RPC failures.
// Error text, HTTP status alone and diagnostic Data are not resource identity evidence.
func isExpectedErrors(err error, codes []string) bool {
	for err != nil {
		if sdkError, ok := err.(*tea.SDKError); ok {
			if sdkError.Code == nil {
				return false
			}
			for _, code := range codes {
				if *sdkError.Code == code {
					return true
				}
			}
			return false
		}
		if wrapped, ok := err.(*errs.ComplexError); ok {
			err = wrapped.Cause
		} else if wrapped, ok := err.(interface{ Unwrap() error }); ok {
			err = wrapped.Unwrap()
		} else {
			return false
		}
	}
	return false
}

// Only exact generic resource codes are accepted when an operation has no declared codes.
func createdByNotExist(err error) bool {
	return isExpectedErrors(err, []string{"ResourceNotFound.CreatedBy", "ResourceNotFound"})
}

// createdByGet describes the resource by id. Returns (nil, nil) when it is gone, mirroring
// v1 GwlbServiceV2.DescribeGwlbListener (Get RPC + retry + ResourceNotFound -> not found).
func createdByGet(ctx context.Context, client *connectivity.AliyunClient, id string, includeSecondary ...bool) (map[string]interface{}, error) {
	completeRead := len(includeSecondary) == 0 || includeSecondary[0]
	getValues := map[string]interface{}{"RegionId": id}
	request := map[string]interface{}{}
	if err := fwadapt.ApplyMappings(request, getValues, []fwadapt.Mapping{
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		return nil, fmt.Errorf("building CheckCreatedByEnabled request: %w", err)
	}
	if value, ok := request["RegionId"].(string); !ok || value == "" || value != id {
		return nil, fmt.Errorf("building CheckCreatedByEnabled request: required identity does not match")
	}
	response, err := retryRpcPost(ctx, client, "Tag", "2018-08-28", "CheckCreatedByEnabled", request, 1*time.Minute, createdByEndpoint())
	if err != nil {
		if createdByNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	response, err = fwSelectScopedReadObject(response, request, []fwScopedReadIdentity{
		{Source: "$.RegionId", Target: "RegionId", Value: id},
	}, []string{"$.OpenStatus"})
	if err != nil {
		return nil, fmt.Errorf("selecting CheckCreatedByEnabled resource: %w", err)
	}
	if response == nil {
		return nil, nil
	}

	object, err := createdByCompleteRead(client, response, id, completeRead)
	if err != nil {
		return nil, err
	}
	// The resource model declares a disabled switch (OpenStatus=false) as "does not exist".
	if enabled, ok := object["OpenStatus"].(bool); ok && !enabled {
		return nil, nil
	}
	return object, nil
}

func createdByCompleteRead(client *connectivity.AliyunClient, object map[string]interface{}, id string, includeSecondary bool) (map[string]interface{}, error) {
	object = fwadapt.NormalizeResponse(object, []fwadapt.Mapping{
		{Source: "$.OpenStatus", Target: "$.OpenStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "$.RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	object["RegionId"] = id
	if !includeSecondary {
		return object, nil
	}
	return object, nil
}

// createdByStateRefreshFunc polls the status field; mirrors v1 GwlbListenerStateRefreshFunc.
// A not-found resource returns a nil result; callers waiting for absence must pass an
// empty Target, because StateChangeConf only treats a nil result as success in that case.
func createdByStateRefreshFunc(ctx context.Context, client *connectivity.AliyunClient, id, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := createdByGet(ctx, client, id, false)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return nil, "", nil
		}
		current := ""
		if values := fwadapt.WalkMappedPath(object, fwadapt.MappingPath(field), nil); len(values) > 0 {
			current = fmt.Sprint(values[0].Value)
		}
		for _, fs := range failStates {
			if current == fs {
				return object, current, fmt.Errorf("status reached fail state %q", current)
			}
		}
		return object, current, nil
	}
}

// createdByWaitForState waits for the status field to reach target, mirroring v1 BuildStateConf.
// Returns the last polled object (nil once the resource disappears).
func createdByWaitForState(ctx context.Context, client *connectivity.AliyunClient, id, field string, target, failStates []string, timeout time.Duration) (map[string]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    createdByStateRefreshFunc(ctx, client, id, field, failStates),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	raw, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if m, ok := raw.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, nil
}

// convertTags mirrors alicloud.ConvertTags: map -> [{"Key":k,"Value":v}, ...] for request bodies.
func convertTags(m map[string]string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(m))
	for k, v := range m {
		out = append(out, map[string]interface{}{"Key": k, "Value": v})
	}
	return out
}

// tagsToStringMap unwraps a framework types.Map (ElementType String) to map[string]string.
func tagsToStringMap(m types.Map) map[string]string {
	out := map[string]string{}
	if m.IsNull() {
		return out
	}
	for k, v := range m.Elements() {
		if s, ok := v.(types.String); ok && !s.IsNull() {
			out[k] = s.ValueString()
		}
	}
	return out
}

// fwTagIgnored follows the SDK provider's tagIgnored prefix contract.
func fwTagIgnored(key string) bool {
	for _, pattern := range []string{"^aliyun", "^acs:", "^http://", "^https://", "^sae.do.not.delete"} {
		matched, _ := regexp.MatchString(pattern, key)
		if matched {
			return true
		}
	}
	return false
}

// tagsToTypeMap wraps an API tags payload ([]interface{} of {"Key","Value"}) into types.Map.
// Known empty API arrays/maps remain empty; omitted or null API values remain null.
func tagsToTypeMap(raw interface{}) types.Map {
	out := map[string]attr.Value{}
	switch value := raw.(type) {
	case []interface{}:
		for _, item := range value {
			if t, ok := item.(map[string]interface{}); ok {
				key, _ := t["TagKey"].(string)
				val, _ := t["TagValue"].(string)
				if key == "" {
					key, _ = t["tag_key"].(string)
					val, _ = t["tag_value"].(string)
				}
				if key == "" {
					key, _ = t["tagKey"].(string)
					val, _ = t["tagValue"].(string)
				}
				if key == "" {
					key, _ = t["key"].(string)
					val, _ = t["value"].(string)
				}
				if key == "" {
					key, _ = t["Key"].(string)
					val, _ = t["Value"].(string)
				}
				if key != "" && !fwTagIgnored(key) {
					out[key] = types.StringValue(val)
				}
			}
		}
	case map[string]string:
		for key, item := range value {
			if !fwTagIgnored(key) {
				out[key] = types.StringValue(item)
			}
		}
	case map[string]interface{}:
		for key, item := range value {
			if fwTagIgnored(key) {
				continue
			}
			text, ok := item.(string)
			if !ok {
				return types.MapNull(types.StringType)
			}
			out[key] = types.StringValue(text)
		}
	default:
		return types.MapNull(types.StringType)
	}
	return types.MapValueMust(types.StringType, out)
}

func fwTagsToAttributes(tags types.Map) []interface{} {
	out := []interface{}{}
	for key, value := range tagsToStringMap(tags) {
		out = append(out, map[string]interface{}{"TagKey": key, "TagValue": value})
	}
	return out
}

// A required exact request can identify a singleton configuration whose API omits its ID.
func fwSelectScopedReadObject(raw interface{}, request map[string]interface{}, identities []fwScopedReadIdentity, payloadPaths []string) (map[string]interface{}, error) {
	if len(identities) == 0 {
		return nil, fmt.Errorf("no exact request identity was declared")
	}
	for _, identity := range identities {
		value, exists := request[identity.Target].(string)
		if !exists || identity.Value == "" || value != identity.Value {
			return nil, fmt.Errorf("singleton request does not match the resource identity")
		}
	}
	object, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("singleton response is not an object")
	}
	for _, identity := range identities {
		values := fwadapt.WalkMappedPath(object, fwadapt.MappingPath(identity.Source), nil)
		if len(values) > 1 {
			return nil, fmt.Errorf("singleton response has ambiguous identity")
		}
		if len(values) == 1 && (values[0].Value == nil || fmt.Sprint(values[0].Value) != identity.Value) {
			return nil, fmt.Errorf("singleton response identity does not match the request")
		}
	}
	for _, path := range payloadPaths {
		values := fwadapt.WalkMappedPath(object, fwadapt.MappingPath(path), nil)
		for _, value := range values {
			if value.Value != nil {
				return object, nil
			}
		}
	}
	return nil, fmt.Errorf("singleton response has no declared configuration payload")
}

// A scoped API may expose several connection kinds; only a declared kind is owned here.
func fwSelectScopedReadList(raw interface{}, request map[string]interface{}, identities []fwScopedReadIdentity, payloadPaths []string, selectorPath, selectorValue string) (map[string]interface{}, error) {
	if len(identities) == 0 || selectorPath == "" || selectorValue == "" {
		return nil, fmt.Errorf("no bounded request-scoped list contract was declared")
	}
	for _, identity := range identities {
		value, exists := request[identity.Target].(string)
		if !exists || identity.Value == "" || value != identity.Value {
			return nil, fmt.Errorf("list request does not match the resource identity")
		}
	}
	rows, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("scoped list response is not an array")
	}
	var selected map[string]interface{}
	for _, row := range rows {
		object, ok := row.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("scoped list item is not an object")
		}
		values := fwadapt.WalkMappedPath(object, fwadapt.MappingPath(selectorPath), nil)
		if len(values) != 1 || values[0].Value == nil {
			return nil, fmt.Errorf("scoped list item has no declared kind")
		}
		if fmt.Sprint(values[0].Value) != selectorValue {
			continue
		}
		if selected != nil {
			return nil, fmt.Errorf("scoped list has multiple matching resource kinds")
		}
		var err error
		selected, err = fwSelectScopedReadObject(object, request, identities, payloadPaths)
		if err != nil {
			return nil, err
		}
	}
	return selected, nil
}

// Ignore cloud-managed tags in old state without sending them to UntagResources.
func fwManagedTagChanges(planned, prior types.Map) ([]interface{}, []interface{}, error) {
	added, removed := []interface{}{}, []interface{}{}
	if planned.IsUnknown() {
		return added, removed, nil
	}
	planMap, priorMap := tagsToStringMap(planned), tagsToStringMap(prior)
	for key := range planMap {
		if fwTagIgnored(key) {
			return nil, nil, fmt.Errorf("tag key %q is reserved for cloud-managed tags", key)
		}
	}
	for key, value := range planMap {
		if previous, exists := priorMap[key]; !exists || previous != value {
			added = append(added, map[string]interface{}{"TagKey": key, "TagValue": value})
		}
	}
	for key, value := range priorMap {
		if fwTagIgnored(key) {
			continue
		}
		if _, exists := planMap[key]; !exists {
			removed = append(removed, map[string]interface{}{"TagKey": key, "TagValue": value})
		}
	}
	return added, removed, nil
}

func createdBySetResourceTags(client *connectivity.AliyunClient, id, resourceType string, planned, prior types.Map, dryRun bool) error {
	added, removed, err := fwManagedTagChanges(planned, prior)
	if err != nil {
		return err
	}
	if len(removed) > 0 {
		return fmt.Errorf("the resource metadata has no remove tag operation")
	}
	if len(added) > 0 {
		return fmt.Errorf("the resource metadata has no add tag operation")
	}
	return nil
}

// fwStr navigates a dotted key in m and returns a types.String (missing -> StringNull).
// Used by nested-list ObjectValues helpers to read scalar sub-fields from each list item;
// the top-level State closures in resource.go capture the Get response, so they cannot read
// per-item sub-maps.
func fwStr(m map[string]interface{}, key string) types.String {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.StringNull()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.StringNull()
		}
	}
	if v, ok := cur.(string); ok {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// fwBool navigates a dotted key in m and returns a types.Bool.
func fwMap(m map[string]interface{}, key string) types.Map {
	raw, exists := m[key]
	if !exists || raw == nil {
		return types.MapNull(types.StringType)
	}
	values := map[string]attr.Value{}
	switch entries := raw.(type) {
	case map[string]interface{}:
		for name, value := range entries {
			if value != nil {
				values[name] = types.StringValue(fmt.Sprint(value))
			}
		}
	case map[string]string:
		for name, value := range entries {
			values[name] = types.StringValue(value)
		}
	default:
		return types.MapNull(types.StringType)
	}
	return types.MapValueMust(types.StringType, values)
}

func fwBool(m map[string]interface{}, key string) types.Bool {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.BoolNull()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.BoolNull()
		}
	}
	if v, ok := cur.(bool); ok {
		return types.BoolValue(v)
	}
	return types.BoolNull()
}

// fwListToJSONString converts a framework types.List of strings to a JSON array
// string (e.g. ["a","b"]) for request params whose API expects a JSON-encoded list.
// Returns "[]" for null/empty. Used by Create/Update/Delete request builders for
// simple string-list fields (e.g. Lindorm ip_white_list) whose fwValueMethod is empty.
// Emitted under fwReqList (any assignable array, incl. request-only responseKey=”)
// rather than fwStateList so request-only arrays get the helper without dragging in
// the response-side state helpers (fwStr/fwBool/fwadapt.Int/ToTypeList) that Go would flag unused.
func fwListToJSONString(l interface {
	IsNull() bool
	Elements() []attr.Value
}) string {
	if l.IsNull() {
		return "[]"
	}
	arr := make([]interface{}, 0, len(l.Elements()))
	for _, e := range l.Elements() {
		if raw := fwadapt.ValueToInterface(e, nil, false); raw != nil {
			arr = append(arr, raw)
		}
	}
	b, _ := json.Marshal(arr)
	return string(b)
}

// fwadapt.WireField preserves API/source names independently of Terraform schema names.

// fwConfiguredChanged distinguishes user unknown references from computed unknown outputs.
// Terraform can defer immutable user input, so an unknown configuration is conservative.
func fwConfiguredChanged(config, plan, state attr.Value) bool {
	if config != nil && config.IsUnknown() {
		return true
	}
	if plan == nil || plan.IsUnknown() {
		return false
	}
	if config == nil || config.IsNull() || state == nil || state.IsNull() || state.IsUnknown() || plan.IsNull() {
		return fwadapt.PlannedChanged(plan, state)
	}
	switch desired := plan.(type) {
	case types.Object:
		configured, configOK := config.(types.Object)
		previous, stateOK := state.(types.Object)
		if !configOK || !stateOK {
			return true
		}
		configItems, oldItems := configured.Attributes(), previous.Attributes()
		for key, value := range desired.Attributes() {
			if fwConfiguredChanged(configItems[key], value, oldItems[key]) {
				return true
			}
		}
		return false
	case types.List:
		configured, configOK := config.(types.List)
		previous, stateOK := state.(types.List)
		if !configOK || !stateOK {
			return true
		}
		items, configItems, oldItems := desired.Elements(), configured.Elements(), previous.Elements()
		if len(items) != len(oldItems) || len(items) != len(configItems) {
			return true
		}
		for index, value := range items {
			if fwConfiguredChanged(configItems[index], value, oldItems[index]) {
				return true
			}
		}
		return false
	case types.Set:
		configured, configOK := config.(types.Set)
		previous, stateOK := state.(types.Set)
		if !configOK || !stateOK || fwContainsUnknown(configured) {
			return true
		}
		return !fwadapt.SetMatches(configured.Elements(), previous.Elements(), fwadapt.ConfiguredMatch)
	case types.Map:
		configured, configOK := config.(types.Map)
		previous, stateOK := state.(types.Map)
		if !configOK || !stateOK {
			return true
		}
		items, configItems, oldItems := desired.Elements(), configured.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for key, value := range items {
			if fwConfiguredChanged(configItems[key], value, oldItems[key]) {
				return true
			}
		}
		return false
	}
	if semantic, ok := plan.(fwadapt.SemanticValue); ok {
		return !semantic.SemanticEqual(state)
	}
	return !plan.Equal(state)
}

func fwContainsUnknown(value attr.Value) bool {
	if value == nil || value.IsNull() {
		return false
	}
	if value.IsUnknown() {
		return true
	}
	switch v := value.(type) {
	case types.Object:
		for _, child := range v.Attributes() {
			if fwContainsUnknown(child) {
				return true
			}
		}
	case types.List:
		for _, child := range v.Elements() {
			if fwContainsUnknown(child) {
				return true
			}
		}
	case types.Set:
		for _, child := range v.Elements() {
			if fwContainsUnknown(child) {
				return true
			}
		}
	case types.Map:
		for _, child := range v.Elements() {
			if fwContainsUnknown(child) {
				return true
			}
		}
	}
	return false
}

// Ambiguous set members with different request-only values must retain prior state and diagnose.
func fwValidateSetCorrelation(observed, prior attr.Value, fields map[string]fwadapt.WireField) error {
	if observed == nil || prior == nil || observed.IsNull() || observed.IsUnknown() || prior.IsNull() || prior.IsUnknown() {
		return nil
	}
	switch current := observed.(type) {
	case types.Set:
		previous, ok := prior.(types.Set)
		if !ok {
			return fmt.Errorf("unordered collection has an incompatible prior type")
		}
		for _, item := range current.Elements() {
			var retained attr.Value
			for _, old := range previous.Elements() {
				if !fwadapt.SetMemberMatches(item, old, fields) {
					continue
				}
				candidate := fwadapt.PreserveRequestValues(item, old, fields)
				if retained != nil && !retained.Equal(candidate) {
					return fmt.Errorf("cannot uniquely correlate unordered members containing request-only configuration")
				}
				retained = candidate
				if err := fwValidateSetCorrelation(item, old, fields); err != nil {
					return err
				}
			}
		}
	case types.List:
		previous, ok := prior.(types.List)
		if !ok {
			return nil
		}
		for index, item := range current.Elements() {
			if index < len(previous.Elements()) {
				if err := fwValidateSetCorrelation(item, previous.Elements()[index], fields); err != nil {
					return err
				}
			}
		}
	case types.Object:
		previous, ok := prior.(types.Object)
		if !ok {
			return nil
		}
		for key, field := range fields {
			if len(field.Fields) > 0 {
				if err := fwValidateSetCorrelation(current.Attributes()[key], previous.Attributes()[key], field.Fields); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func fwDeleteEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	return fwOperationEligible(values, conditions)
}

func fwOperationEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	eligible := true
	for _, condition := range conditions {
		found := fwadapt.WalkMappedPath(values, fwadapt.MappingPath(condition.Property), nil)
		if len(found) == 0 || found[0].Value == nil {
			return false, fmt.Errorf("cannot select an operation: %s is unknown", condition.Property)
		}
		matches := false
		for _, expected := range condition.Values {
			if fmt.Sprint(found[0].Value) == expected {
				matches = true
				break
			}
		}
		switch condition.Action {
		case "valid", "required":
			eligible = eligible && matches
		case "invalid":
			eligible = eligible && !matches
		default:
			return false, fmt.Errorf("unsupported operation dependency action: %s", condition.Action)
		}
	}
	return eligible, nil
}

func fwNestedListToSlice(l attr.Value, fields map[string]fwadapt.WireField, single bool) interface{} {
	return fwadapt.ValueToInterface(l, fields, single)
}

// fwDeltaBatches keeps removals and additions distinct, including scalar attach/detach.
func fwDeltaBatches(oldValue, newValue interface{}, direction string, batchSize int) []interface{} {
	toNumber := func(value interface{}) (float64, bool) {
		switch v := value.(type) {
		case int64:
			return float64(v), true
		case int:
			return float64(v), true
		case float64:
			return v, true
		}
		return 0, false
	}
	if previous, ok := toNumber(oldValue); ok {
		if desired, ok := toNumber(newValue); ok {
			delta := desired - previous
			if direction == "REMOVE" {
				delta = previous - desired
			}
			if delta <= 0 {
				return nil
			}
			return []interface{}{delta}
		}
	}
	oldList, oldIsList := oldValue.([]interface{})
	newList, newIsList := newValue.([]interface{})
	if oldIsList || newIsList {
		left, right := newList, oldList
		if direction == "REMOVE" {
			left, right = oldList, newList
		}
		var delta []interface{}
		for _, item := range left {
			found := false
			for _, previous := range right {
				if fwadapt.ValuesEqual(item, previous) {
					found = true
					break
				}
			}
			if !found {
				delta = append(delta, item)
			}
		}
		if batchSize <= 0 {
			batchSize = len(delta)
		}
		var out []interface{}
		for i := 0; i < len(delta); i += batchSize {
			end := i + batchSize
			if end > len(delta) {
				end = len(delta)
			}
			out = append(out, delta[i:end])
		}
		return out
	}
	if fwadapt.ValuesEqual(oldValue, newValue) {
		return nil
	}
	value := newValue
	if direction == "REMOVE" {
		value = oldValue
	}
	if value == nil || value == "" {
		return nil
	}
	return []interface{}{value}
}
