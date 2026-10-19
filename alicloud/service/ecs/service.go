package ecs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// querySavingPlan reads a single Saving Plan by instance ID via the
// QuerySavingsPlansInstance API (BssOpenApi 2017-12-14). It distinguishes the
// three query outcomes so the resource Read/Create paths can treat them
// differently instead of collapsing them into a single (nil, nil):
//
//   - (object, nil): the plan was found; refresh state from object.
//   - (nil, errSavingPlanQueryNotFound): the query returned a not-found-class
//     signal (an HTTP 404, or a *tea.SDKError whose Code is exactly "NotFound").
//     The BssOpenApi/2017-12-14 QuerySavingsPlansInstance contract does NOT
//     guarantee that such a signal means THIS saving plan instance is gone (a
//     404 may be routing/endpoint/scoped to a different object, and the official
//     error table documents no instance-NotFound code), so this is treated as a
//     "reported not found, retain" signal, NOT a deletion authorisation. Read
//     keeps the paid ID with a warning; a paid financial commitment is never
//     dropped on a not-found signal the contract does not pin to this instance.
//   - (nil, errSavingPlanNotVisible): a legitimate BSS success whose Items list
//     does not contain the plan. This is ambiguous — the plan may be not yet
//     visible (eventual consistency after a purchase) or genuinely absent; the
//     contract cannot tell — so Read keeps the paid ID and warns rather than
//     RemoveResource'ing a financial commitment on an ambiguous empty list.
//   - (nil, err): a business-failure envelope or a malformed/transient
//     response. The caller surfaces this as an error and preserves state for the
//     next refresh; a paid Saving Plan ID is never dropped on a contract
//     violation.
//
// The previous implementation returned (nil, nil) for both the not-found error
// and the empty-Items success, which made Read RemoveResource the paid ID on
// the same empty-Items response Create treated as "keep the ID" — losing a
// financial commitment and forcing a re-purchase on the next plan.
func querySavingPlan(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"InstanceId": id,
		"PageSize":   "10",
		"PageNum":    "1",
	}
	response, err := client.RpcPost("BssOpenApi", "2017-12-14", "QuerySavingsPlansInstance", nil, request, true)
	logs.AddDebug("QuerySavingsPlansInstance", response, request)
	if err != nil {
		if isSavingPlanNotFound(err) {
			return nil, errSavingPlanQueryNotFound
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "QuerySavingsPlansInstance", errs.AlibabaCloudSdkGoERROR)
	}
	// formatError early-returns nil for several non-success envelope shapes — a
	// Code of success/ok/200/0, a Success of true, or an empty/absent Message —
	// so a contradictory envelope (e.g. Code=="Success" but Success==false) or a
	// business failure with an empty Message arrives here as err==nil, often
	// with empty Items. Validate the BSS envelope locally: only a legitimate
	// success (Success==true AND a success Code) authorises the Items list; every
	// other shape is a business failure surfaced as an error so Read keeps the
	// paid ID in state for the next refresh rather than RemoveResource'ing it.
	if isSavingPlanBusinessFailure(response) {
		return nil, fmt.Errorf("QuerySavingsPlansInstance for %s returned a non-success envelope (Success=%v Code=%q Message=%q); preserving Saving Plan state", id, response["Success"], savingPlanEnvelopeString(response, "Code"), savingPlanEnvelopeString(response, "Message"))
	}
	// A 200 response with a malformed body is a contract violation, not a
	// not-found. Only a well-formed legitimate success whose Items list does not
	// contain the plan is the ambiguous "not visible" sentinel; every other shape
	// is an error so Read keeps the existing state for the next refresh and never
	// drops a paid Saving Plan ID on a transiently malformed response.
	data, ok := response["Data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("malformed QuerySavingsPlansInstance response for %s: missing Data object: %v", id, response)
	}
	items, ok := data["Items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("malformed QuerySavingsPlansInstance response for %s: Data.Items is not a list: %v", id, data["Items"])
	}
	if len(items) == 0 {
		return nil, errSavingPlanNotVisible
	}
	object, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("malformed QuerySavingsPlansInstance response for %s: first item is not an object: %v", id, items[0])
	}
	return object, nil
}

// Sentinel errors distinguishing querySavingPlan outcomes. Defined here so the
// resource Read/Create paths dispatch on errors.Is instead of collapsing the
// not-found and not-yet-visible cases: a paid Saving Plan ID is never dropped
// on a not-found signal the BSS contract does not pin to this instance, nor on
// an ambiguous empty-Items response.
var (
	errSavingPlanQueryNotFound = errors.New("QuerySavingsPlansInstance returned a not-found-class signal (404 or Code=NotFound); the BSS contract does not guarantee this means the saving plan instance is gone, so the paid ID is retained")
	errSavingPlanNotVisible    = errors.New("QuerySavingsPlansInstance returned a legitimate success with no matching plan; the plan is not yet visible or its absence cannot be authoritatively confirmed")
)

// savingPlanMaxPages bounds querySavingPlans pagination so a runaway loop
// (e.g. a service reporting a huge TotalCount but never advancing) cannot
// fetch indefinitely. 100 pages at PageSize 100 covers 10,000 plans — far
// beyond any real account — while still terminating. Hitting the cap without
// collecting every reported plan is reported as an error, never a silent
// truncation.
const savingPlanMaxPages = 100

// querySavingPlans lists all Saving Plans via QuerySavingsPlansInstance,
// optionally filtered by instance ID. Used by the saving_plans data source.
//
// It paginates by Data.TotalCount so accounts with more than one page of
// plans are not silently truncated — the previous implementation fetched only
// PageNum=1/PageSize=100 and dropped everything beyond the first 100. The
// official QuerySavingsPlansInstance response carries Data.TotalCount together
// with PageNum/PageSize/Items, which is used here as the authoritative bound.
//
// Pagination is fail-closed and bounded: a hard page cap (savingPlanMaxPages)
// prevents runaway loops; a no-progress guard errors when a page returns no
// items yet TotalCount claims more remain; a later-page transport/server error
// or a malformed later page surfaces as an error instead of collapsing the
// result to whatever was already collected. A malformed envelope — a missing
// Data object, Items that is not a list, a non-object element, or an
// unparseable/absent TotalCount — is reported as an error rather than silently
// collapsing to an empty list, mirroring the querySavingPlan resource path: a
// failed query must never fake success.
func querySavingPlans(client *connectivity.AliyunClient, instanceId string) ([]map[string]interface{}, error) {
	const pageSize = 100
	var plans []map[string]interface{}
	lastTotalCount := 0
	completed := false
	for page := 1; page <= savingPlanMaxPages; page++ {
		request := map[string]interface{}{
			"PageSize": strconv.Itoa(pageSize),
			"PageNum":  strconv.Itoa(page),
		}
		if instanceId != "" {
			request["InstanceId"] = instanceId
		}
		response, err := client.RpcPost("BssOpenApi", "2017-12-14", "QuerySavingsPlansInstance", nil, request, true)
		logs.AddDebug("QuerySavingsPlansInstance", response, request)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ecs_saving_plans", "QuerySavingsPlansInstance", errs.AlibabaCloudSdkGoERROR)
		}
		// The data source reuses the same BSS success contract as querySavingPlan:
		// a business-failure envelope (which formatError may pass through as
		// err==nil, e.g. an empty Message or a contradictory Code) must not be
		// reported as an empty plan list — that would fake success. Only a
		// legitimate success with empty Items returns an empty list.
		if isSavingPlanBusinessFailure(response) {
			return nil, fmt.Errorf("QuerySavingsPlansInstance returned a non-success envelope (Success=%v Code=%q Message=%q)", response["Success"], savingPlanEnvelopeString(response, "Code"), savingPlanEnvelopeString(response, "Message"))
		}
		data, ok := response["Data"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed QuerySavingsPlansInstance response: missing Data object: %v", response)
		}
		items, ok := data["Items"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed QuerySavingsPlansInstance response: Data.Items is not a list: %v", data["Items"])
		}
		// A non-object item is a contract violation, not a silently skippable
		// entry: skipping bad items would hide a malformed response behind a
		// partial success. The querySavingPlan resource path rejects the same
		// shape; the data source aligns with it.
		pagePlans := make([]map[string]interface{}, 0, len(items))
		for _, raw := range items {
			item, ok := raw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("malformed QuerySavingsPlansInstance response: Data.Items contains a non-object element: %v", raw)
			}
			pagePlans = append(pagePlans, item)
		}
		plans = append(plans, pagePlans...)
		totalCount, ok := savingPlanTotalCount(data)
		if !ok {
			return nil, fmt.Errorf("malformed QuerySavingsPlansInstance response: Data.TotalCount is missing or not a number: %v", data["TotalCount"])
		}
		lastTotalCount = totalCount
		// Terminal condition: every plan the service reports has been collected.
		// An empty account (TotalCount 0, no items) breaks here and returns an
		// empty list, which is a genuine absence, not a faked success.
		if len(plans) >= totalCount {
			completed = true
			break
		}
		// No-progress guard: a page that returned no items while TotalCount
		// claims more remain is a contract violation — fail closed instead of
		// looping forever or silently returning a partial list.
		if len(pagePlans) == 0 {
			return nil, fmt.Errorf("QuerySavingsPlansInstance page %d returned no items but TotalCount=%d indicates %d more remain", page, totalCount, totalCount-len(plans))
		}
	}
	// The loop exited by exhausting the safety cap rather than collecting every
	// reported plan. Returning the partial result here would silently truncate,
	// which is exactly the original defect — surface it as an error instead.
	if !completed && len(plans) < lastTotalCount {
		return nil, fmt.Errorf("QuerySavingsPlansInstance pagination exceeded the safety cap of %d pages without collecting all %d reported plans; aborting to avoid silent truncation", savingPlanMaxPages, lastTotalCount)
	}
	return plans, nil
}

// savingPlanTotalCount parses Data.TotalCount, returning (value, true) only
// when the field is present and a genuine number. POP responses decode numbers
// as float64 by default; json.Number (when UseNumber is set), int/int64 and
// string forms are also accepted. A missing, nil or non-numeric value returns
// (0, false) so the caller can fail closed instead of treating it as zero.
func savingPlanTotalCount(data map[string]interface{}) (int, bool) {
	raw, ok := data["TotalCount"]
	if !ok || raw == nil {
		return 0, false
	}
	switch v := raw.(type) {
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return int(n), true
		}
		if f, err := v.Float64(); err == nil {
			return int(f), true
		}
		return 0, false
	case float64:
		return int(v), true
	case float32:
		return int(v), true
	case int:
		return v, true
	case int64:
		return int(v), true
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		n, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(v)))
		if err != nil {
			return 0, false
		}
		return n, true
	}
}

// isSavingPlanNotFound reports whether err is a precise not-found-class signal
// from the QuerySavingsPlansInstance API: an HTTP 404 status, or a *tea.SDKError
// whose Code is exactly "NotFound". It deliberately does NOT match an error
// whose message merely contains the substring "NotFound" (the previous
// implementation did, via errs.NotFoundError's message regex and an explicit
// strings.Contains fallback), which misjudged unrelated errors whose text
// happened to contain "NotFound" as a Saving Plan not-found. The official
// BssOpenApi/2017-12-14 QuerySavingsPlansInstance error table does not document
// an instance NotFound code, so even a precise 404 / "NotFound" code is NOT
// treated as an authoritative "this instance is gone" contract — querySavingPlan
// surfaces it as errSavingPlanQueryNotFound and the resource Read keeps the paid
// ID with a warning, so a not-found signal the contract does not pin to this
// instance (and a propagation-window NotFound right after a successful purchase)
// never drops a paid financial commitment. The resource is removed only by
// deleting it from configuration, never by a query not-found signal.
func isSavingPlanNotFound(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok {
		if tea.IntValue(e.StatusCode) == 404 {
			return true
		}
		if e.Code != nil && *e.Code == "NotFound" {
			return true
		}
	}
	return false
}

// isSavingPlanBusinessFailure reports whether a QuerySavingsPlansInstance
// response envelope indicates a business failure rather than a legitimate
// success. RpcPost→formatError returns nil for several envelope shapes that are
// not real successes — a Code of success/ok/200/0, a Success of true, or an
// empty/absent Message each early-returns nil even when the other fields
// contradict — so a contradictory envelope (e.g. Code=="Success" but
// Success==false) or a business failure with an empty Message slips through as
// a 200 with empty Items. A real BSS success has BOTH Success==true AND a
// success Code; querySavingPlan and querySavingPlans treat anything else as a
// business failure so a paid plan ID is never dropped and a data source never
// returns an empty list for a failed query. Only a legitimate success with
// empty Items (or a real NotFound error) means the plan is genuinely absent.
func isSavingPlanBusinessFailure(response map[string]interface{}) bool {
	success, okS := savingPlanEnvelopeBool(response, "Success")
	code := savingPlanEnvelopeString(response, "Code")
	if okS && success && isSavingPlanSuccessCode(code) {
		return false
	}
	return true
}

// savingPlanEnvelopeBool reads a boolean field from the BSS envelope, checking
// both the Title and lower-case key (POP responses use either form), mirroring
// formatError's dual lookup.
func savingPlanEnvelopeBool(response map[string]interface{}, key string) (bool, bool) {
	for _, k := range []string{key, strings.ToLower(key)} {
		if v, ok := response[k]; ok {
			return fmt.Sprint(v) == "true", true
		}
	}
	return false, false
}

// savingPlanEnvelopeString reads a string field from the BSS envelope, checking
// both the Title and lower-case key.
func savingPlanEnvelopeString(response map[string]interface{}, key string) string {
	for _, k := range []string{key, strings.ToLower(key)} {
		if v, ok := response[k]; ok {
			return fmt.Sprint(v)
		}
	}
	return ""
}

// isSavingPlanSuccessCode mirrors the success-code set formatError treats as
// success: success, ok, 200, 0 (case-insensitive).
func isSavingPlanSuccessCode(code string) bool {
	switch strings.ToLower(code) {
	case "success", "ok", "200", "0":
		return true
	}
	return false
}

// parseSavingPlanCycle parses the Cycle field returned by
// QuerySavingsPlansInstance, which follows the format "N:Unit" (e.g.
// "1:Year"), and returns the leading integer.
func parseSavingPlanCycle(cycle string) (int, error) {
	parts := strings.SplitN(cycle, ":", 2)
	return strconv.Atoi(parts[0])
}

// apiStr converts an API value to types.String, returning null for missing
// or empty values. Non-string values (e.g. float64 for PoolValue) are
// stringified via fmt.Sprint.
func apiStr(m map[string]interface{}, key string) types.String {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok && s != "" {
			return types.StringValue(s)
		}
		return types.StringValue(fmt.Sprint(v))
	}
	return types.StringNull()
}

// periodFromCycle converts the Cycle field to types.Int64, returning null
// when the cycle is missing or unparseable.
func periodFromCycle(v interface{}) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	cycle, ok := v.(string)
	if !ok || cycle == "" {
		return types.Int64Null()
	}
	period, err := parseSavingPlanCycle(cycle)
	if err != nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(period))
}

// savingPlanOutputMap builds a plain map for output_file serialization,
// matching the SDKv2 data-source output field names.
func savingPlanOutputMap(item map[string]interface{}, savingPlanId string) map[string]interface{} {
	m := map[string]interface{}{
		"id":             savingPlanId,
		"saving_plan_id": savingPlanId,
	}
	if v, ok := item["PoolValue"]; ok && v != nil {
		m["committed_amount"] = fmt.Sprint(v)
	}
	if v, ok := item["StartTimestamp"]; ok && v != nil {
		m["create_time"] = fmt.Sprint(v)
	}
	if v, ok := item["InstanceFamily"]; ok && v != nil {
		m["instance_family"] = fmt.Sprint(v)
	}
	if v, ok := item["PayMode"]; ok && v != nil {
		m["offering_type"] = fmt.Sprint(v)
	}
	if v, ok := item["PaymentType"]; ok && v != nil {
		m["payment_type"] = fmt.Sprint(v)
	}
	if v, ok := item["Cycle"]; ok {
		if cycle, ok := v.(string); ok && cycle != "" {
			if period, err := parseSavingPlanCycle(cycle); err == nil {
				m["period"] = period
			}
		}
	}
	if v, ok := item["SavingsType"]; ok && v != nil {
		m["plan_type"] = fmt.Sprint(v)
	}
	if v, ok := item["Region"]; ok && v != nil {
		m["region_id"] = fmt.Sprint(v)
	}
	if v, ok := item["StartTime"]; ok && v != nil {
		m["start_time"] = fmt.Sprint(v)
	}
	return m
}

// writeToFile serializes data as indented JSON and writes it to filePath,
// matching the SDK data-source output_file helper.
func writeToFile(filePath string, data interface{}) error {
	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
	}
	return os.WriteFile(filePath, bs, 0644)
}
