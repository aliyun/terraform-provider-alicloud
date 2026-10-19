package ecs

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	ecsAPIProduct = "Ecs"
	ecsAPIVersion = "2014-05-26"

	// Retry deadlines mirror the SDKv2 PortRangeList resource: 1 minute for
	// the read-side describe calls, 5 minutes for mutating calls.
	describeTimeout = 1 * time.Minute
	mutateTimeout   = 5 * time.Minute
)

// describePortRangeList queries DescribePortRangeLists with a server-side id
// filter (PortRangeListId.1) and returns the matching PortRangeList object. A
// nil result with no error means the list is gone: that signal is a
// structurally-correct empty PortRangeLists array, and the caller (Read/Update)
// then drops state. Anything that is not a clean empty array is not a deletion
// and is returned as an error so Read never silently removes a live resource:
// a missing or nil PortRangeLists root (a truncated/partial response), a field
// of the wrong type, a non-object element, or a non-empty list whose only
// entries are sibling ids (none match the requested id). The whole result is
// scanned for a matching id (rather than trusting lists[0]) so a server that
// returns the sibling first cannot drift Read onto it.
func describePortRangeList(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"PortRangeListId.1": id,
		"RegionId":          client.RegionId,
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.Retry(describeTimeout, func() *retry.RetryError {
		resp, callErr := client.RpcPost(ecsAPIProduct, ecsAPIVersion, "DescribePortRangeLists", nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		response = resp
		return nil
	})
	if err != nil {
		return nil, err
	}
	rawLists, hasLists := response["PortRangeLists"]
	if !hasLists || rawLists == nil {
		// A gone list surfaces as an explicit empty array, not as a missing
		// field. An absent or nil PortRangeLists root is an unexpected response
		// shape (truncation, a partial body, or an id filter that did not
		// round-trip), not a deletion: return an error so Read keeps the prior
		// id in state instead of silently dropping a live resource.
		return nil, fmt.Errorf("DescribePortRangeLists returned no PortRangeLists field for id %s (expected a list); refusing to drop the resource from state", id)
	}
	lists, ok := rawLists.([]interface{})
	if !ok {
		return nil, fmt.Errorf("DescribePortRangeLists returned a malformed PortRangeLists field of type %T (expected a list)", rawLists)
	}
	for _, raw := range lists {
		object, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("DescribePortRangeLists returned a malformed PortRangeList element of type %T (expected an object)", raw)
		}
		if fmt.Sprint(object["PortRangeListId"]) == id {
			return object, nil
		}
	}
	if len(lists) > 0 {
		// The PortRangeListId.1 filter narrows the result to the requested id.
		// A non-empty response whose only entries are sibling ids (none match)
		// is not a clean "gone" signal — a clean deletion is an empty array.
		// Return an error so Read keeps the prior id rather than treating the
		// resource as deleted.
		return nil, fmt.Errorf("DescribePortRangeLists returned %d PortRangeList(s) but none matched id %s; refusing to drop the resource from state", len(lists), id)
	}
	// Structurally-correct empty set: the list call returned no entries, which
	// is the legitimate "resource is gone" signal. The caller (Read/Update)
	// drops state.
	return nil, nil
}

// describePortRangeListEntries queries DescribePortRangeListEntries and returns
// the Entry slice of the list. A not-found id surfaces as the
// InvalidPortRangeListId.NotFound error, which the caller treats as "no entries"
// rather than a hard failure, matching the SDKv2 Read behavior. An Entries
// field that is present but not a list of objects is malformed: it is returned
// as an error so Read never silently replaces a populated entry set with an
// empty one. An absent or nil Entries field is the legitimate empty set.
func describePortRangeListEntries(client *connectivity.AliyunClient, id string) ([]interface{}, error) {
	request := map[string]interface{}{
		"PortRangeListId": id,
		"RegionId":        client.RegionId,
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.Retry(describeTimeout, func() *retry.RetryError {
		resp, callErr := client.RpcPost(ecsAPIProduct, ecsAPIVersion, "DescribePortRangeListEntries", nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		response = resp
		return nil
	})
	if err != nil {
		if isPortRangeListNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	rawEntries, hasEntries := response["Entries"]
	if !hasEntries || rawEntries == nil {
		return nil, nil
	}
	entries, ok := rawEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf("DescribePortRangeListEntries returned a malformed Entries field of type %T (expected a list)", rawEntries)
	}
	for _, raw := range entries {
		if _, ok := raw.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("DescribePortRangeListEntries returned a malformed Entry element of type %T (expected an object)", raw)
		}
	}
	return entries, nil
}

// callRpcWithRetry wraps a mutating RpcPost with the provider's incremental
// retry policy and returns the response. Create/Modify/Delete pass wait=true
// (synchronous completion); TagResources/UntagResources pass wait=false,
// matching the SDKv2 ECS flow.
func callRpcWithRetry(client *connectivity.AliyunClient, action string, request map[string]interface{}, wait bool) (map[string]interface{}, error) {
	retryWait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.Retry(mutateTimeout, func() *retry.RetryError {
		resp, callErr := client.RpcPost(ecsAPIProduct, ecsAPIVersion, action, nil, request, wait)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				retryWait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		response = resp
		return nil
	})
	return response, err
}

// isPortRangeListNotFound reports whether err carries the
// InvalidPortRangeListId.NotFound code that ECS returns when the queried list
// id does not exist. The code arrives on a *tea.SDKError from RpcPost; the
// helpers below wrap mutating and describe calls in retry.Retry, whose
// *retry.RetryError hides the SDK error behind a non-tea type, so the chain is
// unwrapped first. The substring fallback covers any other wrapped form and
// matches IsExpectedErrors' behavior on the SDKv2 line.
func isPortRangeListNotFound(err error) bool {
	if err == nil {
		return false
	}
	current := err
	for current != nil {
		if e, ok := current.(*tea.SDKError); ok && e.Code != nil && *e.Code == "InvalidPortRangeListId.NotFound" {
			return true
		}
		current = errors.Unwrap(current)
	}
	return strings.Contains(err.Error(), "InvalidPortRangeListId.NotFound")
}
