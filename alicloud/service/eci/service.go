// Package eci. This file is generated automatically. Please do not modify it manually, thank you!
// Self-contained service helpers for the framework-native alicloud_eci_commit_container_task resource.
// fwadapt.ResourceBase only exposes Client()/Configure(); the retry/token/tags/state-wait
// helpers live in package alicloud (circular from here) or are unexported in connectivity,
// so this file mirrors the v1 GwlbServiceV2 surface independently inside package eci.
package eci

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type fwTimeoutDurationValidator struct{}

type fwWireField struct {
	Name   string
	Single bool
	Read   bool
	Input  bool
	Fields map[string]fwWireField
}

type fwMapping struct {
	Source, Target, Encoding, Constant string
	Group, GroupKey, GroupValue        string
	Repeat                             bool
	IndexSources                       []int
	Values                             map[string]string
}

type fwMappedValue struct {
	Value   interface{}
	Indexes []int
}
type fwDeleteCondition struct {
	Property, Action string
	Values           []string
}
type fwReadIdentity struct {
	Path, Value  string
	AllowMissing bool
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
func (v fwJSONStringValue) fwNullValue() attr.Value { return fwJSONString(types.StringNull()) }
func (v fwJSONStringValue) fwSemanticEqual(other attr.Value) bool {
	value, ok := other.(fwJSONStringValue)
	if !ok || v.IsNull() || v.IsUnknown() || value.IsNull() || value.IsUnknown() {
		return false
	}
	equal, err := fwJSONDocumentsEqual(v.ValueString(), value.ValueString())
	return err == nil && equal
}
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
	return map[string]attr.Type{"create": types.StringType, "update": types.StringType, "delete": types.StringType}
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	attributes := map[string]schema.Attribute{}
	for _, operation := range []string{"create", "update", "delete"} {
		attributes[operation] = schema.StringAttribute{
			Optional:    true,
			Description: "Maximum retry and wait duration for this operation, for example 10m.",
			Validators:  []validator.String{fwTimeoutDurationValidator{}},
		}
	}
	return schema.SingleNestedBlock{Attributes: attributes}
}

func (fwTimeoutDurationValidator) Description(context.Context) string {
	return "A positive Go duration, for example 10m or 1h."
}
func (v fwTimeoutDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (fwTimeoutDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	duration, err := time.ParseDuration(req.ConfigValue.ValueString())
	if err != nil || duration <= 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
	}
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if value.IsNull() {
		return fallback, diagnostics
	}
	if value.IsUnknown() {
		diagnostics.AddError("Unknown operation timeout", "Resolve the timeouts block before applying.")
		return 0, diagnostics
	}
	raw, ok := value.Attributes()[operation]
	if !ok || raw.IsNull() {
		return fallback, diagnostics
	}
	configured, ok := raw.(types.String)
	if !ok || configured.IsUnknown() {
		diagnostics.AddError("Invalid operation timeout", "Resolve the operation timeout to a duration string before applying.")
		return 0, diagnostics
	}
	duration, err := time.ParseDuration(configured.ValueString())
	if err != nil || duration <= 0 {
		diagnostics.AddError("Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
		return 0, diagnostics
	}
	return duration, diagnostics
}

// retryRpcPost mirrors alicloud's retry.Retry+NeedRetry wrapper around RpcPost.
// RpcPost returns *tea.SDKError, so the classifier only handles that path.
func retryRpcPost(client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.Retry(timeout, func() *retry.RetryError {
		resp, e := client.RpcPost(product, version, action, nil, request, true)
		response = resp
		if e != nil {
			if needRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		return response, errs.WrapErrorf(err, "calling %s failed", action)
	}
	return response, nil
}

// retryRpcGet preserves query parameters for RPC APIs whose metadata only permits GET.
func retryRpcGet(client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.Retry(timeout, func() *retry.RetryError {
		resp, e := client.RpcGet(product, version, action, request, nil)
		response = resp
		if e != nil {
			if needRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		return response, errs.WrapErrorf(err, "calling %s failed", action)
	}
	return response, nil
}

// buildClientToken mirrors alicloud.buildClientToken (TF-<action>-<unix>-<uuid>, <=64 chars).
func buildClientToken(action string) string {
	token := strings.TrimSpace(fmt.Sprintf("TF-%s-%d-%s", action, time.Now().Unix(), strings.Trim(uuid.New().String(), "-")))
	if len(token) > 64 {
		token = token[0:64]
	}
	return token
}

// incrementalWait mirrors alicloud.incrementalWait: exponential backoff capped at 20s with jitter.
func incrementalWait(first, incr time.Duration) func() {
	retryCount := 0
	return func() {
		upper := float64(first) + (math.Pow(2, float64(retryCount))-1)*float64(incr)
		if upper > float64(20*time.Second) {
			upper = float64(20 * time.Second)
		}
		waitTime := time.Duration(upper * (0.8 + 0.4*rand.Float64()))
		time.Sleep(waitTime)
		retryCount++
	}
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
func commitContainerTaskNotExist(err error) bool {
	return isExpectedErrors(err, []string{"ResourceNotFound.CommitContainerTask", "ResourceNotFound"})
}

// commitContainerTaskGet describes the resource by id. Returns (nil, nil) when it is gone, mirroring
// v1 GwlbServiceV2.DescribeGwlbListener (Get RPC + retry + ResourceNotFound -> not found).
func commitContainerTaskGet(client *connectivity.AliyunClient, id string, includeSecondary ...bool) (map[string]interface{}, error) {
	completeRead := len(includeSecondary) == 0 || includeSecondary[0]
	getValues := map[string]interface{}{"RegionId": client.RegionId, "CommitContainerTaskId": id}
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	if err := fwApplyMappings(request, getValues, []fwMapping{
		{Source: "$.ContainerGroupId", Target: "ContainerGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		return nil, fmt.Errorf("building DescribeCommitContainerTask request: %w", err)
	}
	identities := []fwReadIdentity{}
	identities = append(identities, fwReadIdentity{Path: "$.TaskId", Value: id, AllowMissing: false})
	response, err := retryRpcPost(client, "Eci", "2018-08-08", "DescribeCommitContainerTask", request, 1*time.Minute)
	if err != nil {
		if isExpectedErrors(err, []string{"ResourceNotFound.CommitContainerTask", "ResourceNotFound"}) || commitContainerTaskNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	// Gets RootMapping "CommitTasks" 嵌套：沿点分路径逐段导航 Get 响应，末段值若为
	// 数组则取 [0]（列表包装响应如 $.TrailList[*]→TrailList 或 $.KeyPairs.KeyPair[*]→
	// KeyPairs.KeyPair），若为 map 则直接取（对象包装如 $.ResourceGroup→ResourceGroup）。
	// 扁平响应（readRootKey 空）原样返回。
	cur := interface{}(response)
	for _, seg := range strings.Split("CommitTasks", ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, nil
		}
		cur, ok = m[seg]
		if !ok {
			return nil, nil
		}
	}
	inner, identityErr := fwSelectReadObject(cur, identities)
	if identityErr != nil {
		return nil, fmt.Errorf("selecting DescribeCommitContainerTask resource: %w", identityErr)
	}
	if inner != nil {

		return commitContainerTaskCompleteRead(client, inner, id, completeRead)
	}
	return nil, nil
}

func commitContainerTaskCompleteRead(client *connectivity.AliyunClient, object map[string]interface{}, id string, includeSecondary bool) (map[string]interface{}, error) {
	object = fwNormalizeResponse(object, []fwMapping{
		{Source: "$.TaskId", Target: "$.CommitContainerTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.CreationTime", Target: "$.CreateTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.NextToken", Target: "$.NextToken", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.MaxResults", Target: "$.PageTruncated", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TotalCount", Target: "$.TotalCount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.CommitContainerTaskId", Target: "$.CommitContainerTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	object["CommitContainerTaskId"] = id
	if !includeSecondary {
		return object, nil
	}
	return object, nil
}

// commitContainerTaskStateRefreshFunc polls the status field; mirrors v1 GwlbListenerStateRefreshFunc.
// A not-found resource reports state "" so delete state-wait (target [""]) converges.
func commitContainerTaskStateRefreshFunc(client *connectivity.AliyunClient, id, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := commitContainerTaskGet(client, id, false)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return nil, "", nil
		}
		current := ""
		if values := fwWalkMappedPath(object, fwMappingPath(field), nil); len(values) > 0 {
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

// commitContainerTaskWaitForState waits for the status field to reach target, mirroring v1 BuildStateConf.
// Returns the last polled object (nil once the resource disappears).
func commitContainerTaskWaitForState(client *connectivity.AliyunClient, id, field string, target, failStates []string, timeout time.Duration) (map[string]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    commitContainerTaskStateRefreshFunc(client, id, field, failStates),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	raw, err := stateConf.WaitForState()
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

// Missing request-scope components may use only explicitly declared API defaults.
func fwGetIdentityParts(id, separator string, count int, defaults map[int]string, essential []int) ([]string, error) {
	if separator == "" || count < 1 {
		return nil, fmt.Errorf("invalid composite identity definition")
	}
	parts := strings.Split(id, separator)
	if len(parts) > count {
		return nil, fmt.Errorf("composite identity has too many components; expected %d", count)
	}
	for len(parts) < count {
		parts = append(parts, "")
	}
	for index, part := range parts {
		if strings.TrimSpace(part) != "" {
			continue
		}
		required := false
		for _, identityIndex := range essential {
			if index == identityIndex {
				required = true
				break
			}
		}
		if value, declared := defaults[index]; !required && declared && strings.TrimSpace(value) != "" {
			parts[index] = value
			continue
		}
		return nil, fmt.Errorf("composite identity component %d is missing", index+1)
	}
	return parts, nil
}

// Select by real response identity before CompleteRead can restore request-scoped fields.
func fwSelectReadObject(raw interface{}, identities []fwReadIdentity) (map[string]interface{}, error) {
	if raw == nil {
		return nil, nil
	}
	if len(identities) == 0 {
		return nil, fmt.Errorf("no readable resource identity was declared")
	}
	rows, isList := raw.([]interface{})
	if !isList {
		if object, ok := raw.(map[string]interface{}); ok {
			rows = []interface{}{object}
		} else {
			return nil, fmt.Errorf("unexpected resource response shape")
		}
	}
	var selected map[string]interface{}
	for _, row := range rows {
		object, ok := row.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("resource response item is not an object")
		}
		matches, readable := true, false
		for _, identity := range identities {
			values := fwWalkMappedPath(object, fwMappingPath(identity.Path), nil)
			if len(values) > 1 {
				return nil, fmt.Errorf("resource identity %s has multiple values", identity.Path)
			}
			if len(values) == 0 || values[0].Value == nil || fmt.Sprint(values[0].Value) == "" {
				if !identity.AllowMissing {
					matches = false
				}
				continue
			}
			readable = true
			if identity.Value == "" || fmt.Sprint(values[0].Value) != identity.Value {
				matches = false
			}
		}
		if !readable {
			if isList {
				return nil, fmt.Errorf("resource response item has no readable identity")
			}
			continue
		}
		if matches {
			if selected != nil {
				return nil, fmt.Errorf("resource response contains multiple matching identities")
			}
			selected = object
		}
	}
	return selected, nil
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
		values := fwWalkMappedPath(object, fwMappingPath(identity.Source), nil)
		if len(values) > 1 {
			return nil, fmt.Errorf("singleton response has ambiguous identity")
		}
		if len(values) == 1 && (values[0].Value == nil || fmt.Sprint(values[0].Value) != identity.Value) {
			return nil, fmt.Errorf("singleton response identity does not match the request")
		}
	}
	for _, path := range payloadPaths {
		values := fwWalkMappedPath(object, fwMappingPath(path), nil)
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
		values := fwWalkMappedPath(object, fwMappingPath(selectorPath), nil)
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

// Optional update parameters only participate when their mapped value changed.
func fwChangedUpdateMappings(bindings []fwMapping, desired, prior map[string]interface{}, requiredSources []string) []fwMapping {
	selected := make([]fwMapping, 0, len(bindings))
	for _, binding := range bindings {
		required := false
		for _, source := range requiredSources {
			if binding.Source == source || strings.HasPrefix(binding.Source, source+".") || strings.HasPrefix(binding.Source, source+"[*]") {
				required = true
				break
			}
		}
		if required || !fwValuesEqual(fwWalkMappedPath(desired, fwMappingPath(binding.Source), nil), fwWalkMappedPath(prior, fwMappingPath(binding.Source), nil)) {
			selected = append(selected, binding)
		}
	}
	return selected
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

func commitContainerTaskSetResourceTags(client *connectivity.AliyunClient, id, resourceType string, planned, prior types.Map, dryRun bool) error {
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

// fwInt navigates a dotted key in m and returns a types.Int64.
func fwInt(m map[string]interface{}, key string) types.Int64 {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.Int64Null()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.Int64Null()
		}
	}
	switch v := cur.(type) {
	case float64:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return types.Int64Value(n)
		}
	case string:
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return types.Int64Value(n)
		}
	}
	return types.Int64Null()
}

// fwListToJSONString converts a framework types.List of strings to a JSON array
// string (e.g. ["a","b"]) for request params whose API expects a JSON-encoded list.
// Returns "[]" for null/empty. Used by Create/Update/Delete request builders for
// simple string-list fields (e.g. Lindorm ip_white_list) whose fwValueMethod is empty.
// Emitted under fwReqList (any assignable array, incl. request-only responseKey=”)
// rather than fwStateList so request-only arrays get the helper without dragging in
// the response-side state helpers (fwStr/fwBool/fwInt/ToTypeList) that Go would flag unused.
func fwListToJSONString(l interface {
	IsNull() bool
	Elements() []attr.Value
}) string {
	if l.IsNull() {
		return "[]"
	}
	arr := make([]interface{}, 0, len(l.Elements()))
	for _, e := range l.Elements() {
		if raw := fwValueToInterface(e, nil, false); raw != nil {
			arr = append(arr, raw)
		}
	}
	b, _ := json.Marshal(arr)
	return string(b)
}

// fwWireField preserves API/source names independently of Terraform schema names.

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
		return fwPlannedChanged(plan, state)
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
		return !fwSetMatches(configured.Elements(), previous.Elements(), fwConfiguredMatch)
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
	if semantic, ok := plan.(interface{ fwSemanticEqual(attr.Value) bool }); ok {
		return !semantic.fwSemanticEqual(state)
	}
	return !plan.Equal(state)
}

// fwPlannedChanged ignores computed unknown leaves while comparing configured nested values.
func fwPlannedChanged(plan, state attr.Value) bool {
	if plan == nil || plan.IsUnknown() {
		return false
	}
	if state == nil || state.IsNull() || state.IsUnknown() || plan.IsNull() {
		return state == nil || !plan.Equal(state)
	}
	switch desired := plan.(type) {
	case types.Object:
		previous, ok := state.(types.Object)
		if !ok {
			return true
		}
		oldItems := previous.Attributes()
		for key, value := range desired.Attributes() {
			if fwPlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	case types.List:
		previous, ok := state.(types.List)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for index, value := range items {
			if fwPlannedChanged(value, oldItems[index]) {
				return true
			}
		}
		return false
	case types.Set:
		previous, ok := state.(types.Set)
		if !ok {
			return true
		}
		return !fwSetMatches(desired.Elements(), previous.Elements(), func(a, b attr.Value) bool { return !fwPlannedChanged(a, b) })
	case types.Map:
		previous, ok := state.(types.Map)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for key, value := range items {
			if fwPlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	}
	if semantic, ok := plan.(interface{ fwSemanticEqual(attr.Value) bool }); ok {
		return !semantic.fwSemanticEqual(state)
	}
	return !plan.Equal(state)
}

// fwSetMatches pairs unordered members without depending on response or plan order.
func fwSetMatches(left, right []attr.Value, compatible func(attr.Value, attr.Value) bool) bool {
	if len(left) != len(right) {
		return false
	}
	owners := make([]int, len(right))
	for index := range owners {
		owners[index] = -1
	}
	var assign func(int, []bool) bool
	assign = func(index int, seen []bool) bool {
		for candidate := range right {
			if seen[candidate] || !compatible(left[index], right[candidate]) {
				continue
			}
			seen[candidate] = true
			if owners[candidate] == -1 || assign(owners[candidate], seen) {
				owners[candidate] = index
				return true
			}
		}
		return false
	}
	for index := range left {
		if !assign(index, make([]bool, len(right))) {
			return false
		}
	}
	return true
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

// Null configuration leaves are computed output, not member identity.
func fwConfiguredMatch(config, state attr.Value) bool {
	if config == nil || config.IsNull() {
		return true
	}
	if config.IsUnknown() || state == nil || state.IsNull() || state.IsUnknown() {
		return false
	}
	switch value := config.(type) {
	case types.Object:
		previous, ok := state.(types.Object)
		if !ok {
			return false
		}
		fields := previous.Attributes()
		for key, child := range value.Attributes() {
			if !fwConfiguredMatch(child, fields[key]) {
				return false
			}
		}
		return true
	case types.List:
		previous, ok := state.(types.List)
		if !ok || len(value.Elements()) != len(previous.Elements()) {
			return false
		}
		for index, child := range value.Elements() {
			if !fwConfiguredMatch(child, previous.Elements()[index]) {
				return false
			}
		}
		return true
	case types.Set:
		previous, ok := state.(types.Set)
		return ok && fwSetMatches(value.Elements(), previous.Elements(), fwConfiguredMatch)
	case types.Map:
		previous, ok := state.(types.Map)
		if !ok || len(value.Elements()) != len(previous.Elements()) {
			return false
		}
		for key, child := range value.Elements() {
			if !fwConfiguredMatch(child, previous.Elements()[key]) {
				return false
			}
		}
		return true
	}
	if semantic, ok := config.(interface{ fwSemanticEqual(attr.Value) bool }); ok {
		return semantic.fwSemanticEqual(state)
	}
	return config.Equal(state)
}

// Only readable configured inputs correlate set members; computed outputs may change.
func fwSetMemberMatches(observed, prior attr.Value, fields map[string]fwWireField) bool {
	if prior == nil || prior.IsNull() || prior.IsUnknown() {
		return true
	}
	if observed == nil || observed.IsNull() || observed.IsUnknown() {
		return false
	}
	switch current := observed.(type) {
	case types.Object:
		previous, ok := prior.(types.Object)
		if !ok {
			return false
		}
		now, old := current.Attributes(), previous.Attributes()
		for key, field := range fields {
			if !field.Read || !field.Input {
				continue
			}
			if len(field.Fields) > 0 {
				if !fwSetMemberMatches(now[key], old[key], field.Fields) {
					return false
				}
			} else if !fwConfiguredMatch(old[key], now[key]) {
				return false
			}
		}
		return true
	case types.List:
		previous, ok := prior.(types.List)
		if !ok || len(current.Elements()) != len(previous.Elements()) {
			return false
		}
		for index, value := range current.Elements() {
			if !fwSetMemberMatches(value, previous.Elements()[index], fields) {
				return false
			}
		}
		return true
	case types.Set:
		previous, ok := prior.(types.Set)
		return ok && fwSetMatches(current.Elements(), previous.Elements(), func(a, b attr.Value) bool { return fwSetMemberMatches(a, b, fields) })
	}
	return fwConfiguredMatch(prior, observed)
}

// Ambiguous set members with different request-only values must retain prior state and diagnose.
func fwValidateSetCorrelation(observed, prior attr.Value, fields map[string]fwWireField) error {
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
				if !fwSetMemberMatches(item, old, fields) {
					continue
				}
				candidate := fwPreserveRequestValues(item, old, fields)
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

// fwPreserveRequestValues retains configuration-only nested values absent from Get mappings.
// Readable attributes, including API defaults, always come from the observed response.
func fwPreserveRequestValues(observed, prior attr.Value, fields map[string]fwWireField) attr.Value {
	if observed == nil || prior == nil || prior.IsNull() || prior.IsUnknown() {
		return observed
	}
	switch current := observed.(type) {
	case types.List:
		previous, ok := prior.(types.List)
		if !ok || current.IsNull() || current.IsUnknown() {
			return observed
		}
		items, oldItems := current.Elements(), previous.Elements()
		for index := range items {
			if index < len(oldItems) {
				items[index] = fwPreserveRequestValues(items[index], oldItems[index], fields)
			}
		}
		return types.ListValueMust(current.ElementType(context.Background()), items)
	case types.Set:
		previous, ok := prior.(types.Set)
		if !ok || current.IsNull() || current.IsUnknown() {
			return observed
		}
		items, oldItems := current.Elements(), previous.Elements()
		for index, item := range items {
			var retained attr.Value
			ambiguous := false
			for _, old := range oldItems {
				if !fwSetMemberMatches(item, old, fields) {
					continue
				}
				candidate := fwPreserveRequestValues(item, old, fields)
				if retained != nil && !retained.Equal(candidate) {
					ambiguous = true
					break
				}
				retained = candidate
			}
			// Never attach a different member's request-only configuration on ambiguity.
			if retained != nil && !ambiguous {
				items[index] = retained
			}
		}
		return types.SetValueMust(current.ElementType(context.Background()), items)
	case types.Object:
		previous, ok := prior.(types.Object)
		if !ok || current.IsNull() || current.IsUnknown() {
			return observed
		}
		items, oldItems := current.Attributes(), previous.Attributes()
		for key, field := range fields {
			old, exists := oldItems[key]
			if !exists {
				continue
			}
			if !field.Read {
				items[key] = fwKnownValue(old)
			} else if len(field.Fields) > 0 {
				items[key] = fwPreserveRequestValues(items[key], old, field.Fields)
			}
		}
		return types.ObjectValueMust(current.AttributeTypes(context.Background()), items)
	}
	return observed
}

func fwDeleteEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	return fwOperationEligible(values, conditions)
}

func fwOperationEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	eligible := true
	for _, condition := range conditions {
		found := fwWalkMappedPath(values, fwMappingPath(condition.Property), nil)
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

func fwKnownValue(value attr.Value) attr.Value {
	if value != nil && value.IsUnknown() {
		if custom, ok := value.(interface{ fwNullValue() attr.Value }); ok {
			return custom.fwNullValue()
		}
	}
	switch v := value.(type) {
	case types.String:
		if v.IsUnknown() {
			return types.StringNull()
		}
	case types.Bool:
		if v.IsUnknown() {
			return types.BoolNull()
		}
	case types.Int64:
		if v.IsUnknown() {
			return types.Int64Null()
		}
	case types.Float64:
		if v.IsUnknown() {
			return types.Float64Null()
		}
	case types.List:
		if v.IsUnknown() {
			return types.ListNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, fwKnownValue(item))
			}
			return types.ListValueMust(v.ElementType(context.Background()), items)
		}
	case types.Set:
		if v.IsUnknown() {
			return types.SetNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, fwKnownValue(item))
			}
			return types.SetValueMust(v.ElementType(context.Background()), items)
		}
	case types.Map:
		if v.IsUnknown() {
			return types.MapNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Elements() {
				items[key] = fwKnownValue(item)
			}
			return types.MapValueMust(v.ElementType(context.Background()), items)
		}
	case types.Object:
		if v.IsUnknown() {
			return types.ObjectNull(v.AttributeTypes(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Attributes() {
				items[key] = fwKnownValue(item)
			}
			return types.ObjectValueMust(v.AttributeTypes(context.Background()), items)
		}
	}
	return value
}

func fwValueToInterface(value attr.Value, fields map[string]fwWireField, single bool) interface{} {
	if value == nil || value.IsNull() || value.IsUnknown() {
		return nil
	}
	switch v := value.(type) {
	case interface{ ValueString() string }:
		return v.ValueString()
	case types.Bool:
		return v.ValueBool()
	case types.Int64:
		return v.ValueInt64()
	case types.Float64:
		return v.ValueFloat64()
	case types.Object:
		out := map[string]interface{}{}
		for key, child := range v.Attributes() {
			field, mapped := fields[key]
			if !mapped {
				field.Name = key
			}
			if raw := fwValueToInterface(child, field.Fields, field.Single); raw != nil {
				out[field.Name] = raw
			}
		}
		return out
	case types.Map:
		out := map[string]interface{}{}
		for key, child := range v.Elements() {
			if raw := fwValueToInterface(child, nil, false); raw != nil {
				out[key] = raw
			}
		}
		return out
	case types.List:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := fwValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		if single {
			if len(out) == 0 {
				return nil
			}
			return out[0]
		}
		return out
	case types.Set:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := fwValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		return out
	}
	return nil
}

func fwNestedListToSlice(l attr.Value, fields map[string]fwWireField, single bool) interface{} {
	return fwValueToInterface(l, fields, single)
}

func fwMappingPath(path string) []string {
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	path = strings.ReplaceAll(path, "[*]", ".*")
	if path == "" {
		return nil
	}
	return strings.Split(strings.TrimPrefix(path, "."), ".")
}

func fwWalkMappedPath(value interface{}, path []string, indexes []int) []fwMappedValue {
	if len(path) == 0 {
		return []fwMappedValue{{Value: value, Indexes: indexes}}
	}
	if path[0] == "*" {
		var out []fwMappedValue
		if items, ok := value.([]interface{}); ok {
			for i, item := range items {
				next := append(append([]int{}, indexes...), i)
				out = append(out, fwWalkMappedPath(item, path[1:], next)...)
			}
		}
		return out
	}
	if m, ok := value.(map[string]interface{}); ok {
		if child, exists := m[path[0]]; exists {
			return fwWalkMappedPath(child, path[1:], indexes)
		}
	}
	return nil
}

func fwSetMappedPath(value interface{}, path []string, indexes []int, raw interface{}) interface{} {
	if len(path) == 0 {
		return raw
	}
	if path[0] == "*" {
		if len(indexes) == 0 {
			return value
		}
		items, _ := value.([]interface{})
		for len(items) <= indexes[0] {
			items = append(items, nil)
		}
		items[indexes[0]] = fwSetMappedPath(items[indexes[0]], path[1:], indexes[1:], raw)
		return items
	}
	m, _ := value.(map[string]interface{})
	if m == nil {
		m = map[string]interface{}{}
	}
	m[path[0]] = fwSetMappedPath(m[path[0]], path[1:], indexes, raw)
	return m
}

func fwApplyMappings(request, values map[string]interface{}, bindings []fwMapping, roa bool) error {
	mapped := map[string]interface{}{}
	repeated := map[string]bool{}
	for _, binding := range bindings {
		source, target := binding.Source, binding.Target
		if binding.Encoding == "comma" || binding.Encoding == "first" || binding.Encoding == "wrap-json" || binding.Encoding == "wrap-array" {
			source = strings.TrimSuffix(source, "[*]")
			target = strings.TrimSuffix(target, "[*]")
		}
		entries := fwWalkMappedPath(values, fwMappingPath(source), nil)
		if binding.Encoding == "scoped-const" {
			parent := source[:strings.LastIndex(source, ".")]
			entries = fwWalkMappedPath(values, fwMappingPath(parent), nil)
		}
		if strings.HasPrefix(binding.Encoding, "concat|") {
			parts := strings.Split(binding.Encoding, "|")
			fields := []string{}
			complete := true
			for _, propertyPath := range parts[2:] {
				found := fwWalkMappedPath(values, fwMappingPath(propertyPath), nil)
				if len(found) == 0 || found[0].Value == nil {
					complete = false
					break
				}
				fields = append(fields, fmt.Sprint(found[0].Value))
			}
			entries = nil
			if complete {
				entries = []fwMappedValue{{Value: strings.Join(fields, parts[1])}}
			}
		}
		if strings.Contains(source, "[*]") && !strings.Contains(target, "[*]") && binding.Encoding == "" {
			items := []interface{}{}
			for _, entry := range entries {
				if entry.Value != nil {
					items = append(items, entry.Value)
				}
			}
			if len(entries) > 0 {
				entries = []fwMappedValue{{Value: items}}
			}
		}
		concreteEntries := entries[:0]
		for _, entry := range entries {
			if entry.Value != nil {
				concreteEntries = append(concreteEntries, entry)
			}
		}
		entries = concreteEntries
		if binding.Constant != "" && !(binding.Encoding == "default" && len(entries) > 0) {
			if (binding.Encoding == "key" || binding.Encoding == "scoped-const") && len(entries) == 0 {
				continue
			}
			var raw interface{}
			if err := json.Unmarshal([]byte(binding.Constant), &raw); err != nil {
				return err
			}
			if len(entries) == 0 {
				entries = []fwMappedValue{{Value: raw}}
			} else {
				for index := range entries {
					entries[index].Value = raw
				}
			}
		}
		if strings.HasPrefix(binding.Encoding, "uniform-scalar|") && len(entries) > 0 {
			scalarType := strings.TrimPrefix(binding.Encoding, "uniform-scalar|")
			for _, entry := range entries {
				valid := false
				switch scalarType {
				case "boolean":
					_, valid = entry.Value.(bool)
				case "integer":
					switch value := entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
						valid = true
					case float64:
						valid = value == float64(int64(value))
					case json.Number:
						_, err := value.Int64()
						valid = err == nil
					}
				case "number":
					switch entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, json.Number:
						valid = true
					}
				}
				if !valid {
					return fmt.Errorf("configured value does not match scalar API parameter %s type %s", target, scalarType)
				}
			}
			first, err := json.Marshal(entries[0].Value)
			if err != nil {
				return err
			}
			for _, entry := range entries[1:] {
				next, err := json.Marshal(entry.Value)
				if err != nil {
					return err
				}
				if string(next) != string(first) {
					return fmt.Errorf("conflicting configured values for scalar API parameter %s", target)
				}
			}
			entries = entries[:1]
		}
		for _, entry := range entries {
			raw := entry.Value
			if raw == nil {
				continue
			}
			if mappedValue, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mappedValue
			}
			if strings.HasPrefix(binding.Encoding, "suffix|") {
				raw = fmt.Sprint(raw) + strings.TrimPrefix(binding.Encoding, "suffix|")
			}
			switch binding.Encoding {
			case "first":
				if list, ok := raw.([]interface{}); ok {
					if len(list) == 0 {
						continue
					}
					raw = list[0]
				}
			case "comma":
				if list, ok := raw.([]interface{}); ok {
					parts := make([]string, len(list))
					for i, item := range list {
						parts[i] = fmt.Sprint(item)
					}
					raw = strings.Join(parts, ",")
				}
			case "wrap-array":
				raw = []interface{}{raw}
			case "wrap-json":
				encoded, err := json.Marshal([]interface{}{raw})
				if err != nil {
					return err
				}
				raw = string(encoded)
			}
			indexes := append([]int(nil), entry.Indexes...)
			if len(binding.IndexSources) > 0 {
				if len(binding.IndexSources) != strings.Count(target, "[*]") {
					return fmt.Errorf("request wildcard alignment has an unexpected target dimension count")
				}
				indexes = make([]int, len(binding.IndexSources))
				for position, sourceIndex := range binding.IndexSources {
					if sourceIndex == -1 {
						continue
					}
					if sourceIndex < 0 || sourceIndex >= len(entry.Indexes) {
						return fmt.Errorf("request wildcard alignment is missing a source dimension")
					}
					indexes[position] = entry.Indexes[sourceIndex]
				}
			} else {
				for len(indexes) < strings.Count(target, "[*]") {
					indexes = append(indexes, 0)
				}
			}
			mapped = fwSetMappedPath(mapped, fwMappingPath(target), indexes, raw).(map[string]interface{})
		}
		if binding.Repeat {
			if parts := fwMappingPath(target); len(parts) > 0 {
				repeated[parts[0]] = true
			}
		}
	}
	if err := fwCompactKvPairs(mapped, bindings); err != nil {
		return err
	}
	for key, raw := range mapped {
		if roa {
			request[key] = raw
			continue
		}
		switch raw.(type) {
		case map[string]interface{}, []interface{}:
			if repeated[key] {
				fwFlattenRequest(request, key, raw)
			} else {
				encoded, err := json.Marshal(raw)
				if err != nil {
					return err
				}
				request[key] = string(encoded)
			}
		default:
			request[key] = raw
		}
	}
	return nil
}

// Only explicitly declared kvPairs groups are compacted; ordinary indexed lists retain shape.
func fwCompactKvPairs(mapped map[string]interface{}, bindings []fwMapping) error {
	groups := map[string][2]string{}
	for _, binding := range bindings {
		if binding.Group == "" {
			continue
		}
		fields := [2]string{binding.GroupKey, binding.GroupValue}
		if fields[0] == "" || fields[1] == "" || fields[0] == fields[1] {
			return fmt.Errorf("Invalid kvPairs member declaration")
		}
		if previous, exists := groups[binding.Group]; exists && previous != fields {
			return fmt.Errorf("Conflicting kvPairs member declaration")
		}
		groups[binding.Group] = fields
	}
	for group, fields := range groups {
		raw, exists := mapped[group]
		if !exists {
			continue
		}
		members, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("kvPairs group is not an indexed object")
		}
		indexes := []int{}
		for key := range members {
			index, err := strconv.Atoi(key)
			if err != nil || index < 1 || strconv.Itoa(index) != key {
				return fmt.Errorf("Invalid kvPairs member index")
			}
			indexes = append(indexes, index)
		}
		sort.Ints(indexes)
		compact := map[string]interface{}{}
		codes := map[string]bool{}
		for position, index := range indexes {
			pair, valid := members[strconv.Itoa(index)].(map[string]interface{})
			if !valid || len(pair) != 2 || pair[fields[0]] == nil || pair[fields[1]] == nil {
				return fmt.Errorf("Incomplete kvPairs member")
			}
			code, valid := pair[fields[0]].(string)
			if !valid || code == "" || codes[code] {
				return fmt.Errorf("Invalid or duplicate kvPairs code")
			}
			codes[code] = true
			compact[strconv.Itoa(position+1)] = pair
		}
		mapped[group] = compact
	}
	return nil
}

func fwFlattenRequest(request map[string]interface{}, key string, raw interface{}) {
	switch v := raw.(type) {
	case []interface{}:
		for i, item := range v {
			fwFlattenRequest(request, fmt.Sprintf("%s.%d", key, i+1), item)
		}
	case map[string]interface{}:
		for name, item := range v {
			fwFlattenRequest(request, key+"."+name, item)
		}
	default:
		if v != nil {
			request[key] = v
		}
	}
}

func fwNormalizeResponse(object map[string]interface{}, bindings []fwMapping) map[string]interface{} {
	// Whole-object mappings must precede their leaf aliases, regardless of AM order.
	// Stable insertion order preserves same-depth precedence and source array indexes.
	bindings = append([]fwMapping(nil), bindings...)
	for index := 1; index < len(bindings); index++ {
		for current := index; current > 0 && len(fwMappingPath(bindings[current].Target)) < len(fwMappingPath(bindings[current-1].Target)); current-- {
			bindings[current], bindings[current-1] = bindings[current-1], bindings[current]
		}
	}
	out := map[string]interface{}{}
	for key, value := range object {
		out[key] = value
	}
	for _, binding := range bindings {
		if strings.HasSuffix(binding.Source, "[*]") && strings.HasSuffix(binding.Target, "[*]") {
			source := strings.TrimSuffix(binding.Source, "[*]")
			target := strings.TrimSuffix(binding.Target, "[*]")
			for _, container := range fwWalkMappedPath(object, fwMappingPath(source), nil) {
				if values, ok := container.Value.([]interface{}); ok && len(values) == 0 {
					out = fwSetMappedPath(out, fwMappingPath(target), container.Indexes, []interface{}{}).(map[string]interface{})
				}
			}
		}
		for _, entry := range fwWalkMappedPath(object, fwMappingPath(binding.Source), nil) {
			raw := entry.Value
			target := binding.Target
			if binding.Encoding == "split-comma" {
				target = strings.TrimSuffix(target, "[*]")
				items := []interface{}{}
				if text, ok := raw.(string); ok && text != "" {
					for _, part := range strings.Split(text, ",") {
						items = append(items, strings.TrimSpace(part))
					}
				}
				raw = items
			}
			if binding.Encoding == "scalar-string" && raw != nil {
				switch raw.(type) {
				case string:
				case int, int32, int64, float32, float64, json.Number:
					raw = fmt.Sprint(raw)
				}
			}
			if mapped, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mapped
			}
			out = fwSetMappedPath(out, fwMappingPath(target), entry.Indexes, raw).(map[string]interface{})
		}
	}
	return out
}

func fwValuesEqual(a, b interface{}) bool {
	aa, errA := json.Marshal(a)
	bb, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aa) == string(bb)
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
				if fwValuesEqual(item, previous) {
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
	if fwValuesEqual(oldValue, newValue) {
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
