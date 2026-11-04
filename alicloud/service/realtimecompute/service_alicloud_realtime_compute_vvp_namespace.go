// Package realtimecompute. This file is generated automatically. Please do not modify it manually, thank you!
// Framework helpers for the framework-native alicloud_realtime_compute_vvp_namespace resource.
// Shared behaviour comes from alicloud/helper (BuildClientToken, NeedRetry, RetryRpcPost),
// alicloud/errs (IsExpectedErrors), alicloud/tags (TagIgnored), alicloud/provider/fwutil
// (timeouts) and alicloud/provider/fwadapt (framework wire/state adapters); this file keeps only
// the vvp-namespace specific service helpers.
package realtimecompute

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// retryRpcPost runs the write call through the shared context-aware retry path, so cancelling the
// Terraform lifecycle aborts the backoff instead of continuing to send requests.
func retryRpcPost(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost(product, version, action, nil, request, true)
	}, timeout, vvpNamespaceRetryable)
}

// retryRpcGet keeps RPC APIs whose metadata only permits GET on the same shared retry path, so a
// read shares the lifecycle cancellation and the transient classification with the write calls.
func retryRpcGet(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcGet(product, version, action, request, nil)
	}, timeout, vvpNamespaceRetryable)
}

// vvpNamespaceRetryable extends the shared transient classification with the product codes this
// resource retries on top of it: 903047 means a spec change is already running, and the
// conflict/status/busy and DependencyViolation codes clear by themselves. 903099 stays fatal: the
// target is below the capacity the namespace already holds, so repeating the call cannot succeed.
func vvpNamespaceRetryable(err error) bool {
	if helper.NeedRetry(err) {
		return true
	}
	e, ok := err.(*tea.SDKError)
	if !ok || e.Code == nil || e.Message == nil {
		return false
	}
	code := *e.Code
	// DependencyViolation carries a product suffix (.HaVip/.GatewayRouteTable/...), so match by
	// prefix; the rest are the delete/read conflict, status and busy codes the API returns while a
	// previous change is still settling.
	return strings.HasPrefix(code, "DependencyViolation") || code == "OperationConflict" ||
		code == "LastTokenProcessing" || code == "IncorrectStatus" || code == "SystemBusy" ||
		code == "903047"
}

// vvpNamespaceOperationBudget bounds one lifecycle operation. Every stage of the operation (RPC,
// visibility wait, SUCCESS/gone wait) spends the same deadline instead of restarting a full
// timeout, so the configured timeouts.<operation> is a ceiling for the whole operation and a
// cancelled Terraform lifecycle still aborts the stage in flight.
func vvpNamespaceOperationBudget(ctx context.Context, timeout time.Duration) (context.Context, func() time.Duration, context.CancelFunc) {
	deadline := time.Now().Add(timeout)
	operationCtx, cancel := context.WithDeadline(ctx, deadline)
	return operationCtx, func() time.Duration { return time.Until(deadline) }, cancel
}

// vvpNamespaceCreateResultUncertain reports whether a failed CreateNamespace call may still have
// created the namespace. Transport timeouts, connection failures and server-side 5xx responses do
// not prove the request was rejected, so the caller looks the namespace up by its
// <instance_id>:<namespace> identity instead of reporting a failure and leaving an untracked
// namespace. Deterministic product answers (validation, quota, conflict codes) are not ambiguous.
func vvpNamespaceCreateResultUncertain(err error) bool {
	return vvpNamespaceCreateResultUncertainDepth(err, 0)
}

// vvpNamespaceCreateResultUncertainDepth bounds the wrapper unwrapping so a cyclic cause can never
// spin here.
func vvpNamespaceCreateResultUncertainDepth(err error, depth int) bool {
	if err == nil || depth > 8 {
		return false
	}
	// errs.ComplexError does not implement Unwrap, so the cause has to be unwrapped explicitly;
	// otherwise every classification below only ever sees the wrapper.
	var complexErr errs.ComplexError
	if errors.As(err, &complexErr) && complexErr.Cause != nil && complexErr.Cause != err {
		return vvpNamespaceCreateResultUncertainDepth(complexErr.Cause, depth+1)
	}
	var complexErrPtr *errs.ComplexError
	if errors.As(err, &complexErrPtr) && complexErrPtr.Cause != nil && complexErrPtr.Cause != err {
		return vvpNamespaceCreateResultUncertainDepth(complexErrPtr.Cause, depth+1)
	}
	if errors.Is(err, context.Canceled) {
		// The caller cancelled the operation; adopting a result here would fight the cancellation.
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	var sdkErr *tea.SDKError
	if errors.As(err, &sdkErr) && sdkErr.StatusCode != nil && *sdkErr.StatusCode >= 500 {
		return true
	}
	message := strings.ToLower(err.Error())
	for _, marker := range []string{
		"timeout", "timed out", "connection reset", "connection refused", "broken pipe",
		"unexpected eof", "no such host", "tls handshake", "server closed",
	} {
		if strings.Contains(message, marker) {
			return true
		}
	}
	return false
}

// vvpNamespaceRecoverCreate checks the identity after an uncertain CreateNamespace failure. It
// returns true when the namespace exists, which means the create was accepted even though its
// response was lost. A lookup failure is returned so the caller can report both causes.
func vvpNamespaceRecoverCreate(ctx context.Context, client *connectivity.AliyunClient, id string, cause error) (map[string]interface{}, error) {
	if !vvpNamespaceCreateResultUncertain(cause) {
		return nil, nil
	}
	// The create may have failed because the operation budget itself elapsed. The lookup must still
	// run, so it gets its own short deadline instead of inheriting the expired one.
	lookupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	object, err := vvpNamespaceGet(lookupCtx, client, id)
	if err != nil {
		return nil, fmt.Errorf("checking whether %s was created despite %v failed: %w", id, cause, err)
	}
	// A namespace that exists but reached the failure state is not a successful create.
	if object != nil && fmt.Sprint(object["Status"]) == "FAILED" {
		return nil, nil
	}
	return object, nil
}

// vvpNamespaceBudgetForCall returns the time left in the operation budget for one call and fails
// fast once the budget is spent, so the caller reports an exhausted budget instead of handing a
// non-positive timeout to a retry loop or a wait.
func vvpNamespaceBudgetForCall(remaining func() time.Duration) (time.Duration, error) {
	left := remaining()
	if left <= 0 {
		return 0, fmt.Errorf("the operation timeout elapsed before this stage could run")
	}
	return left, nil
}

// vvpNamespaceGet describes the resource by id. Returns (nil, nil) when it is gone, mirroring
// v1 GwlbServiceV2.DescribeGwlbListener (Get RPC + retry + ResourceNotFound -> not found).
func vvpNamespaceGet(ctx context.Context, client *connectivity.AliyunClient, id string, includeSecondary ...bool) (map[string]interface{}, error) {
	completeRead := len(includeSecondary) == 0 || includeSecondary[0]
	getValues := map[string]interface{}{"RegionId": client.RegionId, "InstanceId": id}
	parts, identityPartsErr := fwadapt.GetIdentityParts(id, ":", 2, map[int]string{}, []int{1})
	if identityPartsErr != nil {
		return nil, fmt.Errorf("reading alicloud_realtime_compute_vvp_namespace: %w", identityPartsErr)
	}
	if len(parts) > 0 {
		getValues["InstanceId"] = parts[0]
	}
	if len(parts) > 1 {
		getValues["Namespace"] = parts[1]
	}
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	if err := fwadapt.ApplyMappings(request, getValues, []fwadapt.Mapping{
		{Source: "$.InstanceId", Target: "DescribeNamespacesRequest.InstanceId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Namespace", Target: "DescribeNamespacesRequest.Namespace", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "DescribeNamespacesRequest.Region", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
	}, false); err != nil {
		return nil, fmt.Errorf("building DescribeNamespaces request: %w", err)
	}
	identities := []fwadapt.ReadIdentity{}
	if len(parts) > 1 {
		identities = append(identities, fwadapt.ReadIdentity{Path: "$.Namespace", Value: parts[1], AllowMissing: false})
	}
	response, err := retryRpcGet(ctx, client, "foasconsole", "2019-06-01", "DescribeNamespaces", request, 1*time.Minute)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.VvpNamespace", "ResourceNotFound"}) {
			return nil, nil
		}
		return nil, err
	}
	// Gets RootMapping "Namespaces". Only a decoded collection is evidence of what the workspace
	// holds: the list is handed to fwadapt.SelectReadObject, which selects the row whose $.Namespace
	// equals the requested identity; multiple matching rows and list items without a readable
	// identity are errors, and the first element is never used as a fallback. A decoded list without
	// this identity is a positive not-found, while an omitted, null or mistyped field is an
	// unreadable response and fails the read instead of clearing state.
	rows, rowsErr := vvpNamespaceResponseRows(response)
	if rowsErr != nil {
		return nil, fmt.Errorf("reading alicloud_realtime_compute_vvp_namespace %s: %w", id, rowsErr)
	}
	inner, identityErr := fwadapt.SelectReadObject(rows, identities)
	if identityErr != nil {
		return nil, fmt.Errorf("selecting DescribeNamespaces resource: %w", identityErr)
	}
	if inner == nil {
		return nil, nil
	}
	return vvpNamespaceCompleteRead(client, inner, id, completeRead)
}

// vvpNamespaceResponseRows extracts the "Namespaces" collection from a DescribeNamespaces response.
// The 2019-06-01 response schema is not published in the API metadata, so an omitted, null or
// mistyped Namespaces field is treated as an unreadable response instead of an empty one: guessing
// "empty" would let Read drop a live namespace from state on a decoding or product-shape change.
// Only an explicitly decoded array is accepted, and an empty array is a legal empty collection.
func vvpNamespaceResponseRows(response map[string]interface{}) ([]interface{}, error) {
	if response == nil {
		return nil, fmt.Errorf("DescribeNamespaces returned an empty response")
	}
	raw, ok := response["Namespaces"]
	if !ok {
		return nil, fmt.Errorf("DescribeNamespaces response does not contain Namespaces")
	}
	if raw == nil {
		return nil, fmt.Errorf("DescribeNamespaces response Namespaces is null")
	}
	rows, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("DescribeNamespaces response Namespaces has unexpected type %T", raw)
	}
	return rows, nil
}

func vvpNamespaceCompleteRead(client *connectivity.AliyunClient, object map[string]interface{}, id string, includeSecondary bool) (map[string]interface{}, error) {
	object = fwadapt.NormalizeResponse(object, []fwadapt.Mapping{
		{Source: "$.GmtCreate", Target: "$.GmtCreate", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.GmtModified", Target: "$.GmtModified", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Ha", Target: "$.Ha", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Namespace", Target: "$.Namespace", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceSpec.Cpu", Target: "$.ResourceSpec.Cpu", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceSpec.MemoryGB", Target: "$.ResourceSpec.MemoryGb", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceUsed.Cpu", Target: "$.ResourceUsed.Cpu", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceUsed.Cu", Target: "$.ResourceUsed.Cu", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceUsed.MemoryGB", Target: "$.ResourceUsed.MemoryGb", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Status", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Tags[*]", Target: "$.Tags[*]", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Tags[*].Key", Target: "$.Tags[*].Key", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Tags[*].Value", Target: "$.Tags[*].Value", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceId", Target: "$.InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Namespace", Target: "$.Namespace", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	identityParts := strings.Split(id, ":")
	if len(identityParts) > 0 {
		object["InstanceId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		object["Namespace"] = identityParts[1]
	}
	if !includeSecondary {
		return object, nil
	}
	return object, nil
}

// vvpNamespaceStateRefreshFunc polls the status field; mirrors v1 GwlbListenerStateRefreshFunc.
// A not-found resource reports state "" so an empty-target delete wait (no target statuses)
// converges; a non-empty target such as [""] would instead fail with a NotFoundError.
func vvpNamespaceStateRefreshFunc(ctx context.Context, client *connectivity.AliyunClient, id, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := vvpNamespaceGet(ctx, client, id, false)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return nil, "", nil
		}
		current := ""
		if strings.HasPrefix(field, "#") {
			if values := fwadapt.WalkMappedPath(object, fwadapt.MappingPath(strings.TrimPrefix(field, "#")), nil); len(values) > 0 && values[0].Value != nil {
				current = "#CHECKSET"
			}
		} else {
			if values := fwadapt.WalkMappedPath(object, fwadapt.MappingPath(field), nil); len(values) > 0 {
				current = fmt.Sprint(values[0].Value)
			}
		}
		for _, fs := range failStates {
			if current == fs {
				return object, current, fmt.Errorf("status reached fail state %q", current)
			}
		}
		return object, current, nil
	}
}

// vvpNamespaceWaitForState waits for the status field to reach target, mirroring v1 BuildStateConf.
// Returns the last polled object (nil once the resource disappears).
// DescribeVvpNamespace returns the namespace object for the <instance_id>:<namespace> identity,
// or nil when DescribeNamespaces no longer returns it.
func DescribeVvpNamespace(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	return vvpNamespaceGet(context.Background(), client, id)
}

func vvpNamespaceWaitForState(ctx context.Context, client *connectivity.AliyunClient, id, field string, target, failStates []string, timeout time.Duration) (map[string]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    vvpNamespaceStateRefreshFunc(ctx, client, id, field, failStates),
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
				if key != "" && !tags.TagIgnored(key, nil) {
					out[key] = types.StringValue(val)
				}
			}
		}
	case map[string]string:
		for key, item := range value {
			if !tags.TagIgnored(key, nil) {
				out[key] = types.StringValue(item)
			}
		}
	case map[string]interface{}:
		for key, item := range value {
			if tags.TagIgnored(key, nil) {
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

// fwFloat navigates a dotted key in m and returns a types.Float64. The CloudSpec model types
// ResourceUsed.Cpu/Cu/MemoryGb as number/float, so usage values must not be truncated.
func fwFloat(m map[string]interface{}, key string) types.Float64 {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.Float64Null()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.Float64Null()
		}
	}
	switch v := cur.(type) {
	case float64:
		return types.Float64Value(v)
	case float32:
		return types.Float64Value(float64(v))
	case int:
		return types.Float64Value(float64(v))
	case int64:
		return types.Float64Value(float64(v))
	case json.Number:
		if n, err := v.Float64(); err == nil {
			return types.Float64Value(n)
		}
	case string:
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return types.Float64Value(n)
		}
	}
	return types.Float64Null()
}

// vvpNamespaceResourceSpecToTypeList wraps the API object array into a framework
// types.List of nested objects (types.Object). nil/empty/non-array -> empty list
// (.#=0) to match v1 d.Set(..., []), so TestCheckResourceAttrSet("x.#") passes and
// UseStateForUnknown unknown->0 stays consistent (ListNull leaves .# unset, failing
// the .# check).
func vvpNamespaceResourceSpecToTypeList(raw interface{}) types.List {
	objType := vvpNamespaceResourceSpecObjType()
	elems := make([]attr.Value, 0)
	// Handle both array (nested list) and single object (MaxItems:1 object type, gap#7)
	if raw != nil {
		var arr []interface{}
		switch v := raw.(type) {
		case []interface{}:
			arr = v
		case map[string]interface{}:
			arr = []interface{}{v}
		}
		for _, item := range arr {
			m, _ := item.(map[string]interface{})
			obj, _ := types.ObjectValue(vvpNamespaceResourceSpecAttrTypes(), vvpNamespaceResourceSpecObjectValues(m))
			elems = append(elems, obj)
		}
	}
	l, _ := types.ListValue(objType, elems)
	return l
}

// vvpNamespaceResourceSpecAttrTypes returns the framework attr types for each nested field.
func vvpNamespaceResourceSpecAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cpu":       types.Int64Type,
		"memory_gb": types.Int64Type,
	}
}

// vvpNamespaceResourceSpecObjType is the element ObjectType for the nested list.
func vvpNamespaceResourceSpecObjType() types.ObjectType {
	return types.ObjectType{AttrTypes: vvpNamespaceResourceSpecAttrTypes()}
}

// vvpNamespaceResourceSpecObjectValues maps a single API object (map) to framework attr
// values. Nested sub-lists recurse via their own ToTypeList helpers; scalars use fwStr/fwBool/fwadapt.Int.
func vvpNamespaceResourceSpecObjectValues(m map[string]interface{}) map[string]attr.Value {
	out := map[string]attr.Value{}
	out["cpu"] = fwadapt.Int(m, "Cpu")
	out["memory_gb"] = fwadapt.Int(m, "MemoryGb")
	return out
}

// vvpNamespaceResourceUsedToTypeList wraps the API object array into a framework
// types.List of nested objects (types.Object). nil/empty/non-array -> empty list
// (.#=0) to match v1 d.Set(..., []), so TestCheckResourceAttrSet("x.#") passes and
// UseStateForUnknown unknown->0 stays consistent (ListNull leaves .# unset, failing
// the .# check).
func vvpNamespaceResourceUsedToTypeList(raw interface{}) types.List {
	objType := vvpNamespaceResourceUsedObjType()
	elems := make([]attr.Value, 0)
	// Handle both array (nested list) and single object (MaxItems:1 object type, gap#7)
	if raw != nil {
		var arr []interface{}
		switch v := raw.(type) {
		case []interface{}:
			arr = v
		case map[string]interface{}:
			arr = []interface{}{v}
		}
		for _, item := range arr {
			m, _ := item.(map[string]interface{})
			obj, _ := types.ObjectValue(vvpNamespaceResourceUsedAttrTypes(), vvpNamespaceResourceUsedObjectValues(m))
			elems = append(elems, obj)
		}
	}
	l, _ := types.ListValue(objType, elems)
	return l
}

// vvpNamespaceResourceUsedAttrTypes returns the framework attr types for each nested field.
func vvpNamespaceResourceUsedAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cpu":       types.Float64Type,
		"cu":        types.Float64Type,
		"memory_gb": types.Float64Type,
	}
}

// vvpNamespaceResourceUsedObjType is the element ObjectType for the nested list.
func vvpNamespaceResourceUsedObjType() types.ObjectType {
	return types.ObjectType{AttrTypes: vvpNamespaceResourceUsedAttrTypes()}
}

// vvpNamespaceResourceUsedObjectValues maps a single API object (map) to framework attr
// values. Nested sub-lists recurse via their own ToTypeList helpers; scalars use fwStr/fwBool/fwadapt.Int.
func vvpNamespaceResourceUsedObjectValues(m map[string]interface{}) map[string]attr.Value {
	out := map[string]attr.Value{}
	out["cpu"] = fwFloat(m, "Cpu")
	out["cu"] = fwFloat(m, "Cu")
	out["memory_gb"] = fwFloat(m, "MemoryGb")
	return out
}

// fwadapt.WireField preserves API/source names independently of Terraform schema names.
