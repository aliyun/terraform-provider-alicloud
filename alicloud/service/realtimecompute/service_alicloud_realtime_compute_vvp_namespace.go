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
// Terraform lifecycle aborts the backoff instead of continuing to send requests. Each attempt sends
// the lifecycle context into the HTTP request itself, so cancellation also ends the socket in
// flight; this wrapper is the only retry owner, so no later retry or duplicate write is sent either.
func retryRpcPost(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPostWithContext(ctx, product, version, action, nil, request)
	}, timeout, vvpNamespaceRetryable)
}

// retryRpcGet keeps RPC APIs whose metadata only permits GET on the same shared retry path, so a
// read shares the lifecycle cancellation, the in-flight deadline and the transient classification
// with the write calls.
func retryRpcGet(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcGetWithContext(ctx, product, version, action, request, nil)
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
// returns the object when the namespace now exists. Existence alone does not prove this request
// created it: the product exposes no request correlation or idempotency token, so the caller must
// NOT adopt the object. The caller reports an uncertain result and leaves the object out of state,
// with import or ownership-coordination guidance. A lookup failure is returned so the caller can
// report both causes.
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
	return vvpNamespaceGetWithTimeout(ctx, client, id, 1*time.Minute, includeSecondary...)
}

// vvpNamespaceDescribePageSize bounds one DescribeNamespaces page. The published product contract
// (foasconsole 2021-10-28) documents Namespace, PageIndex/PageSize and TotalCount/TotalPage for
// this collection; the 2019-06-01 metadata is not published, so the read cannot assume the filter
// always returns the identity on the first page and walks the pages instead.
const vvpNamespaceDescribePageSize = 100

// vvpNamespaceDescribeMaxPages bounds the walk so a product that keeps reporting another page
// cannot turn one read into an unbounded loop.
const vvpNamespaceDescribeMaxPages = 100

// vvpNamespaceGetWithTimeout describes the resource with the read bounded by the caller's remaining
// operation budget, so a stage late in a lifecycle cannot spend more time than the operation is
// allowed to take. The context still carries the caller's cancellation. The walk covers the
// documented pagination of the collection and only accepts a row whose parent is the requested
// instance, so a same-name namespace under another workspace is never mistaken for this one.
func vvpNamespaceGetWithTimeout(ctx context.Context, client *connectivity.AliyunClient, id string, readTimeout time.Duration, includeSecondary ...bool) (map[string]interface{}, error) {
	completeRead := len(includeSecondary) == 0 || includeSecondary[0]
	getValues := map[string]interface{}{"RegionId": client.RegionId, "InstanceId": id}
	parts, identityPartsErr := fwadapt.GetIdentityParts(id, ":", 2, map[int]string{}, []int{1})
	if identityPartsErr != nil {
		return nil, fmt.Errorf("reading alicloud_realtime_compute_vvp_namespace: %w", identityPartsErr)
	}
	instanceId := id
	if len(parts) > 0 {
		instanceId = parts[0]
	}
	if len(parts) > 1 {
		getValues["Namespace"] = parts[1]
	}
	identities := []fwadapt.ReadIdentity{}
	if len(parts) > 1 {
		identities = append(identities, fwadapt.ReadIdentity{Path: "$.Namespace", Value: parts[1], AllowMissing: false})
	}
	seenRows := 0
	for page := 1; page <= vvpNamespaceDescribeMaxPages; page++ {
		getValues["InstanceId"] = instanceId
		if len(parts) > 1 {
			getValues["Namespace"] = parts[1]
		}
		request := map[string]interface{}{}
		request["RegionId"] = client.RegionId
		request["DescribeNamespacesRequest.PageIndex"] = page
		request["DescribeNamespacesRequest.PageSize"] = vvpNamespaceDescribePageSize
		if err := fwadapt.ApplyMappings(request, getValues, []fwadapt.Mapping{
			{Source: "$.InstanceId", Target: "DescribeNamespacesRequest.InstanceId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.Namespace", Target: "DescribeNamespacesRequest.Namespace", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "DescribeNamespacesRequest.Region", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		}, false); err != nil {
			return nil, fmt.Errorf("building DescribeNamespaces request: %w", err)
		}
		response, err := retryRpcGet(ctx, client, "foasconsole", "2019-06-01", "DescribeNamespaces", request, readTimeout)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ResourceNotFound.VvpNamespace", "ResourceNotFound"}) {
				return nil, nil
			}
			return nil, err
		}
		// The walk advances its own page counter, so a response that echoes a different page (a
		// repeated page 1 or a foreign page) would silently cover a page that was never received.
		// Validate the optional echoed PageIndex before the rows are used for selection or absence.
		if pageErr := vvpNamespaceResponsePage(response, page); pageErr != nil {
			return nil, fmt.Errorf("reading alicloud_realtime_compute_vvp_namespace %s: %w", id, pageErr)
		}
		// Gets RootMapping "Namespaces". Only a decoded collection is evidence of what the workspace
		// holds: the list is handed to fwadapt.SelectReadObject, which selects the row whose
		// $.Namespace equals the requested identity; multiple matching rows and list items without a
		// readable identity are errors, and the first element is never used as a fallback. A decoded
		// list without this identity is a positive not-found, while an omitted, null or mistyped
		// field is an unreadable response and fails the read instead of clearing state.
		rows, rowsErr := vvpNamespaceResponseRows(response)
		if rowsErr != nil {
			return nil, fmt.Errorf("reading alicloud_realtime_compute_vvp_namespace %s: %w", id, rowsErr)
		}
		seenRows += len(rows)
		// Scope the rows to the requested parent before selection and duplicate detection: a
		// same-name namespace under another workspace is a different resource, and leaving it in the
		// candidate set would make the requested one look ambiguous. Rows without a readable
		// InstanceId stay candidates, because the documented missing-InstanceId limitation only says
		// their parent cannot be disproved (and they are not adopted as another parent either).
		scopedRows := make([]interface{}, 0, len(rows))
		for _, row := range rows {
			if object, ok := row.(map[string]interface{}); ok {
				if parent, ok := object["InstanceId"].(string); ok && parent != "" && parent != instanceId {
					continue
				}
			}
			scopedRows = append(scopedRows, row)
		}
		inner, identityErr := fwadapt.SelectReadObject(scopedRows, identities)
		if identityErr != nil {
			return nil, fmt.Errorf("selecting DescribeNamespaces resource: %w", identityErr)
		}
		if inner != nil {
			// The collection is the one parent's list, but the identity of a row still has to match
			// the requested workspace: a same-name namespace under another parent is a different
			// resource and must not be adopted or reported as this one.
			if parent, ok := inner["InstanceId"].(string); ok && parent != "" && parent != instanceId {
				inner = nil
			}
		}
		if inner != nil {
			return vvpNamespaceCompleteRead(client, inner, id, completeRead)
		}
		complete, completeErr := vvpNamespaceSearchComplete(response, page, len(rows), seenRows)
		if completeErr != nil {
			// An unreadable or contradictory page shape means the search cannot be proven complete.
			// Reporting not-found here would remove a live namespace from state.
			return nil, fmt.Errorf("reading alicloud_realtime_compute_vvp_namespace %s: %w", id, completeErr)
		}
		if complete {
			return nil, nil
		}
	}
	return nil, fmt.Errorf("reading alicloud_realtime_compute_vvp_namespace %s: the namespace was not found within %d DescribeNamespaces pages", id, vvpNamespaceDescribeMaxPages)
}

// vvpNamespaceResponsePage validates the optional echoed PageIndex of a DescribeNamespaces response
// against the page that was requested. A present value must match: a repeated or foreign page cannot
// be used to prove absence, because the walk would otherwise move past a page it never received. A
// response without PageIndex keeps the legacy single-page shape and is handled by the completeness
// checks below.
func vvpNamespaceResponsePage(response map[string]interface{}, requested int) error {
	page, present, err := vvpNamespaceIntField(response, "PageIndex")
	if err != nil {
		return err
	}
	if !present {
		return nil
	}
	if page != requested {
		return fmt.Errorf("the response is for page %d but page %d was requested", page, requested)
	}
	return nil
}

// vvpNamespaceSearchComplete reports whether a DescribeNamespaces page proves that the requested
// identity is not in the collection. Absence only counts once the search is complete:
//   - with a page count, this must be the last page, and the page count has to be able to cover the
//     reported total (a positive TotalPage that contradicts TotalCount/PageSize is incomplete);
//   - without a page count (the legacy single-page shape), the page must not be full and the reported
//     total must not exceed the rows read so far.
//
// Anything else keeps the walk going. An unreadable or contradictory field is an error so an
// incomplete search never turns into a not-found that removes a live namespace from state. That
// includes a reported total below the rows already read: the page metadata cannot then be trusted
// to prove how much of the collection was covered.
func vvpNamespaceSearchComplete(response map[string]interface{}, page int, pageRows int, seenRows int) (bool, error) {
	totalPage, hasTotalPage, err := vvpNamespaceIntField(response, "TotalPage")
	if err != nil {
		return false, err
	}
	totalCount, hasTotalCount, err := vvpNamespaceIntField(response, "TotalCount")
	if err != nil {
		return false, err
	}
	pageSize, hasPageSize, err := vvpNamespaceIntField(response, "PageSize")
	if err != nil {
		return false, err
	}
	// The rows read so far cannot outnumber the total the response reports. Such a contradiction
	// means the counts describe different collections (or a broken page), so neither the current
	// page nor the walk may be used to prove absence.
	if hasTotalCount && totalCount < seenRows {
		return false, fmt.Errorf("the response read %d namespace(s) but reports TotalCount %d", seenRows, totalCount)
	}
	if hasTotalPage {
		if totalPage == 0 {
			if seenRows > 0 {
				return false, fmt.Errorf("the response returned %d namespace(s) but TotalPage is 0", seenRows)
			}
			if hasTotalCount && totalCount > 0 {
				return false, fmt.Errorf("the response reports %d namespaces but TotalPage is 0", totalCount)
			}
			return true, nil
		}
		if hasTotalCount && totalCount > seenRows && page >= totalPage {
			return false, fmt.Errorf("the response reports %d namespaces but only %d namespace(s) were read across %d page(s) with TotalPage %d", totalCount, seenRows, page, totalPage)
		}
		if hasTotalCount && hasPageSize && pageSize > 0 && totalCount > totalPage*pageSize {
			return false, fmt.Errorf("the response reports %d namespaces but TotalPage %d with PageSize %d cannot cover them", totalCount, totalPage, pageSize)
		}
		// A positive page count below the page just read is contradictory: the collection cannot have
		// fewer pages than the one the service returned. Treating it as the last page would let the
		// read conclude absence and clear state, so it is an error regardless of TotalCount/PageSize.
		if page > totalPage {
			return false, fmt.Errorf("the response is for page %d but reports TotalPage %d", page, totalPage)
		}
		return page >= totalPage, nil
	}
	// No page count. A full page, or a reported total above what has been read, leaves the search
	// open; only a short page without a matching total proves absence.
	effectivePageSize := vvpNamespaceDescribePageSize
	if hasPageSize && pageSize > 0 {
		effectivePageSize = pageSize
	}
	if hasTotalCount && totalCount > seenRows {
		return false, nil
	}
	if pageRows >= effectivePageSize {
		return false, nil
	}
	return true, nil
}

// vvpNamespaceIntField reads an optional non-negative integer field from a DescribeNamespaces
// response. A missing or null field reports present=false; a mistyped, non-integral or negative value
// is an error, because guessing a value there could make an incomplete search look complete.
func vvpNamespaceIntField(response map[string]interface{}, name string) (int, bool, error) {
	raw, present := response[name]
	if !present || raw == nil {
		return 0, false, nil
	}
	var parsed int64
	switch value := raw.(type) {
	case float64:
		if float64(int64(value)) != value {
			return 0, false, fmt.Errorf("the %s field %v is not an integer", name, value)
		}
		parsed = int64(value)
	case int:
		parsed = int64(value)
	case int64:
		parsed = value
	case json.Number:
		n, err := value.Int64()
		if err != nil {
			return 0, false, fmt.Errorf("the %s field %q is not an integer", name, value.String())
		}
		parsed = n
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return 0, false, fmt.Errorf("the %s field %q is not an integer", name, value)
		}
		parsed = int64(n)
	default:
		return 0, false, fmt.Errorf("the %s field has an unreadable type %T", name, raw)
	}
	if parsed < 0 {
		return 0, false, fmt.Errorf("the %s field is negative (%d)", name, parsed)
	}
	return int(parsed), true, nil
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

// vvpNamespaceStrictIntValue converts an already-navigated response value to an integer. A
// non-integral or non-numeric value is an error rather than a silent null, because silently mapping
// it would overwrite a valid state value with a wrong one.
func vvpNamespaceStrictIntValue(raw interface{}, name string) (int64, error) {
	switch value := raw.(type) {
	case float64:
		if float64(int64(value)) != value {
			return 0, fmt.Errorf("the %s value %v is not an integer", name, value)
		}
		return int64(value), nil
	case float32:
		if float32(int64(value)) != value {
			return 0, fmt.Errorf("the %s value %v is not an integer", name, value)
		}
		return int64(value), nil
	case int:
		return int64(value), nil
	case int32:
		return int64(value), nil
	case int64:
		return value, nil
	case json.Number:
		parsed, err := value.Int64()
		if err != nil {
			return 0, fmt.Errorf("the %s value %v is not an integer", name, value)
		}
		return parsed, nil
	case string:
		// The service has returned these counters as numeric strings; parse them rather than losing
		// the value, but a non-numeric string stays an error.
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("the %s value %q is not an integer", name, value)
		}
		return parsed, nil
	}
	return 0, fmt.Errorf("the %s field has unexpected type %T", name, raw)
}

// vvpNamespaceStrictFloatValue converts an already-navigated response value to a float. A
// non-numeric value is an error for the same reason as vvpNamespaceStrictIntValue.
func vvpNamespaceStrictFloatValue(raw interface{}, name string) (float64, error) {
	switch value := raw.(type) {
	case float64:
		return value, nil
	case float32:
		return float64(value), nil
	case int:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case json.Number:
		parsed, err := value.Float64()
		if err != nil {
			return 0, fmt.Errorf("the %s value %v is not a number", name, value)
		}
		return parsed, nil
	case string:
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Errorf("the %s value %q is not a number", name, value)
		}
		return parsed, nil
	}
	return 0, fmt.Errorf("the %s field has unexpected type %T", name, raw)
}

// tagsToTypeMap wraps an API tags payload ([]interface{} of {"Key","Value"}) into types.Map.
// Known empty API arrays/maps remain empty and omitted or null API values remain null; an element or
// value whose type cannot be read is an error instead of being silently skipped or nulled.
func tagsToTypeMap(raw interface{}) (types.Map, error) {
	keyNames := []string{"TagKey", "tag_key", "tagKey", "key", "Key"}
	valueNames := []string{"TagValue", "tag_value", "tagValue", "value", "Value"}
	out := map[string]attr.Value{}
	switch value := raw.(type) {
	case nil:
		return types.MapNull(types.StringType), nil
	case []interface{}:
		for index, item := range value {
			t, ok := item.(map[string]interface{})
			if !ok {
				return types.MapNull(types.StringType), fmt.Errorf("the Tags[%d] element has unexpected type %T", index, item)
			}
			key, readable, err := vvpNamespaceTagField(t, keyNames)
			if err != nil {
				return types.MapNull(types.StringType), fmt.Errorf("Tags[%d]: %w", index, err)
			}
			if !readable {
				return types.MapNull(types.StringType), fmt.Errorf("the Tags[%d] element has no readable tag key", index)
			}
			tagValue, _, err := vvpNamespaceTagField(t, valueNames)
			if err != nil {
				return types.MapNull(types.StringType), fmt.Errorf("Tags[%d]: %w", index, err)
			}
			if key != "" && !tags.TagIgnored(key, nil) {
				out[key] = types.StringValue(tagValue)
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
				return types.MapNull(types.StringType), fmt.Errorf("the Tags[%q] value has unexpected type %T", key, item)
			}
			out[key] = types.StringValue(text)
		}
	default:
		return types.MapNull(types.StringType), fmt.Errorf("the Tags field has unexpected type %T", raw)
	}
	mapped, diags := types.MapValue(types.StringType, out)
	if diags.HasError() {
		return types.MapNull(types.StringType), fmt.Errorf("building the Tags map: %v", diags)
	}
	return mapped, nil
}

// vvpNamespaceTagField reads the first recognised tag field spelling. readable=false means none of
// the spellings is present; a present field of another type is an error.
func vvpNamespaceTagField(t map[string]interface{}, names []string) (string, bool, error) {
	for _, name := range names {
		raw, present := t[name]
		if !present {
			continue
		}
		text, ok := raw.(string)
		if !ok {
			return "", false, fmt.Errorf("the tag field %q has unexpected type %T", name, raw)
		}
		return text, true, nil
	}
	return "", false, nil
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

// vvpNamespaceOptionalStrictInt reads an optional integer field of a response element. A missing or
// null value reports present=false; a present value of another type is an error.
func vvpNamespaceOptionalStrictInt(object map[string]interface{}, name string) (int64, bool, error) {
	raw, present := object[name]
	if !present || raw == nil {
		return 0, false, nil
	}
	value, err := vvpNamespaceStrictIntValue(raw, name)
	if err != nil {
		return 0, false, err
	}
	return value, true, nil
}

// vvpNamespaceOptionalStrictFloat is the float counterpart of vvpNamespaceOptionalStrictInt.
func vvpNamespaceOptionalStrictFloat(object map[string]interface{}, name string) (float64, bool, error) {
	raw, present := object[name]
	if !present || raw == nil {
		return 0, false, nil
	}
	value, err := vvpNamespaceStrictFloatValue(raw, name)
	if err != nil {
		return 0, false, err
	}
	return value, true, nil
}

// vvpNamespaceResourceSpecToTypeList wraps the API object array into a framework types.List of
// nested objects. A nil/absent value stays an empty list (.#=0) to match the earlier behaviour; a
// present value of another type, or an element/field that cannot be read, is an error so a
// malformed row cannot replace a valid resource_spec with an empty one.
func vvpNamespaceResourceSpecToTypeList(raw interface{}) (types.List, error) {
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
		default:
			return types.ListNull(objType), fmt.Errorf("the ResourceSpec field has unexpected type %T", raw)
		}
		for index, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return types.ListNull(objType), fmt.Errorf("the ResourceSpec[%d] element has unexpected type %T", index, item)
			}
			values, err := vvpNamespaceResourceSpecObjectValues(m)
			if err != nil {
				return types.ListNull(objType), fmt.Errorf("ResourceSpec[%d]: %w", index, err)
			}
			obj, diags := types.ObjectValue(vvpNamespaceResourceSpecAttrTypes(), values)
			if diags.HasError() {
				return types.ListNull(objType), fmt.Errorf("building ResourceSpec[%d]: %v", index, diags)
			}
			elems = append(elems, obj)
		}
	}
	list, diags := types.ListValue(objType, elems)
	if diags.HasError() {
		return types.ListNull(objType), fmt.Errorf("building the ResourceSpec list: %v", diags)
	}
	return list, nil
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
func vvpNamespaceResourceSpecObjectValues(m map[string]interface{}) (map[string]attr.Value, error) {
	out := map[string]attr.Value{}
	cpu, present, err := vvpNamespaceOptionalStrictInt(m, "Cpu")
	if err != nil {
		return nil, err
	}
	if present {
		out["cpu"] = types.Int64Value(cpu)
	} else {
		out["cpu"] = types.Int64Null()
	}
	memoryGb, present, err := vvpNamespaceOptionalStrictInt(m, "MemoryGb")
	if err != nil {
		return nil, err
	}
	if present {
		out["memory_gb"] = types.Int64Value(memoryGb)
	} else {
		out["memory_gb"] = types.Int64Null()
	}
	return out, nil
}

// vvpNamespaceResourceUsedToTypeList wraps the API object array into a framework types.List of
// nested objects. A nil/absent value stays an empty list (.#=0) to match the earlier behaviour; a
// present value of another type, or an element/field that cannot be read, is an error.
func vvpNamespaceResourceUsedToTypeList(raw interface{}) (types.List, error) {
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
		default:
			return types.ListNull(objType), fmt.Errorf("the ResourceUsed field has unexpected type %T", raw)
		}
		for index, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return types.ListNull(objType), fmt.Errorf("the ResourceUsed[%d] element has unexpected type %T", index, item)
			}
			values, err := vvpNamespaceResourceUsedObjectValues(m)
			if err != nil {
				return types.ListNull(objType), fmt.Errorf("ResourceUsed[%d]: %w", index, err)
			}
			obj, diags := types.ObjectValue(vvpNamespaceResourceUsedAttrTypes(), values)
			if diags.HasError() {
				return types.ListNull(objType), fmt.Errorf("building ResourceUsed[%d]: %v", index, diags)
			}
			elems = append(elems, obj)
		}
	}
	list, diags := types.ListValue(objType, elems)
	if diags.HasError() {
		return types.ListNull(objType), fmt.Errorf("building the ResourceUsed list: %v", diags)
	}
	return list, nil
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
func vvpNamespaceResourceUsedObjectValues(m map[string]interface{}) (map[string]attr.Value, error) {
	out := map[string]attr.Value{}
	cpu, present, err := vvpNamespaceOptionalStrictFloat(m, "Cpu")
	if err != nil {
		return nil, err
	}
	if present {
		out["cpu"] = types.Float64Value(cpu)
	} else {
		out["cpu"] = types.Float64Null()
	}
	cu, present, err := vvpNamespaceOptionalStrictFloat(m, "Cu")
	if err != nil {
		return nil, err
	}
	if present {
		out["cu"] = types.Float64Value(cu)
	} else {
		out["cu"] = types.Float64Null()
	}
	memoryGb, present, err := vvpNamespaceOptionalStrictFloat(m, "MemoryGb")
	if err != nil {
		return nil, err
	}
	if present {
		out["memory_gb"] = types.Float64Value(memoryGb)
	} else {
		out["memory_gb"] = types.Float64Null()
	}
	return out, nil
}

// fwadapt.WireField preserves API/source names independently of Terraform schema names.
