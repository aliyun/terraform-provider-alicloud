// Package eci. This file is generated automatically. Please do not modify it manually, thank you!
// Self-contained service helpers for the framework-native alicloud_eci_commit_container_task resource.
// fwadapt.ResourceBase only exposes Client()/Configure(); the retry, token and value-conversion
// helpers that the generated resource needs are kept local to package eci.
package eci

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
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
type fwReadIdentity struct {
	Path, Value  string
	AllowMissing bool
}

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
	return map[string]attr.Type{"create": types.StringType}
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	attributes := map[string]schema.Attribute{}
	for _, operation := range []string{"create"} {
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
		// Transient conflict/busy codes are retried by every RPC call in this package.
		// DependencyViolation carries a product suffix (.HaVip/.GatewayRouteTable/...), so
		// it is matched by prefix; all listed codes are bounded by the call retry timeout.
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

// commitContainerTaskGet describes the resource by id. Returns (nil, nil) when the task is
// no longer returned: the Get RPC is retried and a declared not-found code maps to absence.
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
	// DescribeCommitContainerTask requires ContainerGroupId or TaskId and otherwise returns the
	// first page of every task in the region (default page size 10). Without the filter the
	// identity selection below can miss the task once the account holds more than one page,
	// which would make Read drop an existing task from state. TaskId is an indexed repeatList
	// query/body parameter, so bind the resource id explicitly.
	request["TaskId.1"] = id
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
		{Source: "$.CommitContainerTaskId", Target: "$.CommitContainerTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	object["CommitContainerTaskId"] = id
	if !includeSecondary {
		return object, nil
	}
	return object, nil
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
