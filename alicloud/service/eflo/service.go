// Package eflo. This file started as the generated helper layer for the framework-native
// alicloud_eflo_er_attachment resource and was then reduced by hand to the declarations that
// resource executes, so regenerating it has to re-apply the same reduction. The timeouts block
// comes from alicloud/provider/fwutil and the client token / retry loop from alicloud/helper;
// the fw mapping machinery, the Get normalisation and the state waits stay here.
package eflo

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

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

func fwTimeoutTypes() map[string]attr.Type {
	return fwutil.TimeoutTypes("create", "update", "delete")
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	return fwutil.TimeoutBlock("create", "update", "delete")
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	return fwutil.TimeoutDuration(value, operation, fallback)
}

// retryRpcPost wraps RpcPost in the shared ctx-aware retry loop with the product's
// transient-error classifier; a cancelled lifecycle stops retrying.
func retryRpcPost(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost(product, version, action, nil, request, true)
	}, timeout, retryableErAttachmentError)
}

func buildClientToken(action string) string {
	return helper.BuildClientToken(action)
}

// retryableErAttachmentError applies the shared transient classifier to the Eflo
// RPCs. The Eflo 2022-05-30 definition published to this provider carries no
// error-code table, so no product-specific code is retried on top of it.
func retryableErAttachmentError(err error) bool {
	return helper.NeedRetry(err)
}

// isErAttachmentNotFound reports whether err is the contract's not-found answer.
// The code is compared exactly: errs.IsExpectedErrors also matches prefixes and
// response bodies, which would let an unrelated failure read as "resource is gone"
// and drop it from state.
func isErAttachmentNotFound(err error) bool {
	for err != nil {
		if sdkError, ok := err.(*tea.SDKError); ok {
			return sdkError.Code != nil && *sdkError.Code == "1003"
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

// erAttachmentGet describes the resource by id. Returns (nil, nil) when it is gone, mirroring
// v1 GwlbServiceV2.DescribeGwlbListener (Get RPC + retry + ResourceNotFound -> not found).
func erAttachmentGet(ctx context.Context, client *connectivity.AliyunClient, id string, includeSecondary ...bool) (map[string]interface{}, error) {
	getValues := map[string]interface{}{"RegionId": client.RegionId, "ErId": id}
	parts, identityPartsErr := fwGetIdentityParts(id, ":", 2, map[int]string{}, []int{1})
	if identityPartsErr != nil {
		return nil, fmt.Errorf("reading alicloud_eflo_er_attachment: %w", identityPartsErr)
	}
	if len(parts) > 0 {
		getValues["ErId"] = parts[0]
	}
	if len(parts) > 1 {
		getValues["ErAttachmentId"] = parts[1]
	}
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	if err := fwApplyMappings(request, getValues, []fwMapping{
		{Source: "$.ErId", Target: "ErId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErAttachmentId", Target: "ErAttachmentId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		return nil, fmt.Errorf("building GetErAttachment request: %w", err)
	}
	identities := []fwReadIdentity{}
	if len(parts) > 0 {
		identities = append(identities, fwReadIdentity{Path: "$.ErId", Value: parts[0], AllowMissing: true})
	}
	if len(parts) > 1 {
		identities = append(identities, fwReadIdentity{Path: "$.ErAttachmentId", Value: parts[1], AllowMissing: false})
	}
	response, err := retryRpcPost(ctx, client, "eflo", "2022-05-30", "GetErAttachment", request, 1*time.Minute)
	if err != nil {
		if isErAttachmentNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	// Gets RootMapping "Content" 嵌套：沿点分路径逐段导航 Get 响应，末段值若为
	// 数组则取 [0]（列表包装响应如 $.TrailList[*]→TrailList 或 $.KeyPairs.KeyPair[*]→
	// KeyPairs.KeyPair），若为 map 则直接取（对象包装如 $.ResourceGroup→ResourceGroup）。
	// 扁平响应（readRootKey 空）原样返回。
	// Absence is signalled by error code 1003 (handled above), never by a successful
	// response, so a 2xx body without the mapped object is a malformed answer and must
	// surface as an error instead of removing the resource from state.
	cur := interface{}(response)
	for _, seg := range strings.Split("Content", ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("GetErAttachment returned an unexpected response shape for %s", id)
		}
		cur, ok = m[seg]
		if !ok || cur == nil {
			return nil, fmt.Errorf("GetErAttachment response for %s does not contain Content", id)
		}
	}
	inner, identityErr := fwSelectReadObject(cur, identities)
	if identityErr != nil {
		return nil, fmt.Errorf("selecting GetErAttachment resource: %w", identityErr)
	}
	if inner == nil {
		return nil, fmt.Errorf("GetErAttachment response for %s does not match the resource identity", id)
	}
	return erAttachmentCompleteRead(inner, id)
}

func erAttachmentCompleteRead(object map[string]interface{}, id string) (map[string]interface{}, error) {
	object = fwNormalizeResponse(object, []fwMapping{
		{Source: "$.Message", Target: "$.Message", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErId", Target: "$.ErId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceName", Target: "$.InstanceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Status", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "$.RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErAttachmentName", Target: "$.ErAttachmentName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.CreateTime", Target: "$.CreateTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AutoReceiveAllRoute", Target: "$.AutoReceiveAllRoute", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceId", Target: "$.InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErAttachmentId", Target: "$.ErAttachmentId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.GmtModified", Target: "$.GmtModified", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Across", Target: "$.Across", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceType", Target: "$.InstanceType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceTenantId", Target: "$.ResourceTenantId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TenantId", Target: "$.TenantId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErId", Target: "$.ErId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErAttachmentId", Target: "$.ErAttachmentId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	identityParts := strings.Split(id, ":")
	if len(identityParts) > 0 {
		object["ErId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		object["ErAttachmentId"] = identityParts[1]
	}
	return object, nil
}

// erAttachmentStateRefreshFunc polls the status field. A not-found resource reports
// state "" so the empty-target gone-wait used by Delete converges.
func erAttachmentStateRefreshFunc(ctx context.Context, client *connectivity.AliyunClient, id, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := erAttachmentGet(ctx, client, id, false)
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

// erAttachmentWaitForState waits for the status field to reach target, mirroring v1 BuildStateConf.
// Returns the last polled object (nil once the resource disappears).
func erAttachmentWaitForState(ctx context.Context, client *connectivity.AliyunClient, id, field string, target, failStates []string, timeout time.Duration) (map[string]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    erAttachmentStateRefreshFunc(ctx, client, id, field, failStates),
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

func fwKnownValue(value attr.Value) attr.Value {
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

func fwValuesEqual(a, b interface{}) bool {
	aa, errA := json.Marshal(a)
	bb, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aa) == string(bb)
}
