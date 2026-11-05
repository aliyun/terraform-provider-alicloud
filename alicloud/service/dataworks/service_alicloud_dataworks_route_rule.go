// Package dataworks. The alicloud_dataworks_route_rule resource is hand-written
// (Framework-native) against the dataworks-public 2018-06-01 RouteRule APIs. The
// four target APIs (CreateRouteRule/DeleteRouteRule/UpdateRouteRule/ListRouteRules)
// are visibility=Private and return an OpenApiResultModel wrapper (success /
// errorCode / errorMessage / data) instead of the standard Code/Message envelope,
// so the response success is inspected here and the business error is rebuilt
// from errorCode/errorMessage rather than left as an empty formatError.
package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// maxExactFloatInteger is the largest integer float64 represents without loss.
// The tea RPC decoder renders numbers as float64 (and sometimes json.Number),
// so a route rule id above this range can only be trusted when it arrives as a
// json.Number or string; a float64 beyond it is rejected as contract drift.
const maxExactFloatInteger = 1 << 53

// RouteRule API identifiers (dataworks-public 2018-06-01, visibility=Private).
const (
	routeRuleProduct     = "dataworks-public"
	routeRuleAPIVersion  = "2018-06-01"
	routeRuleCreateAPI   = "CreateRouteRule"
	routeRuleUpdateAPI   = "UpdateRouteRule"
	routeRuleDeleteAPI   = "DeleteRouteRule"
	routeRuleListAPI     = "ListRouteRules"
	routeRuleReadTimeout = 2 * time.Minute
)

// routeRuleWriteTimeout bounds the synchronous create/update/delete RPCs. The
// published @returnMode is async:false with a 20s backend timeout, so a few
// retries cover the transient classifier without a long fixed wait.
const routeRuleWriteTimeout = 5 * time.Minute

// routeRuleCompositeId joins the three identity parts that Read needs. ListRouteRules
// is addressed by ResourceGroupId + ResourceId (the DataWorks resource group and the
// network resource the rule is attached to) and has no single Get API, so the route
// rule id alone cannot re-read the resource. The composite keeps every Read
// self-contained, including after import. resource_group_id and route_rule_id are
// numeric, while resource_id (a network resource id such as "ns-xxxx") is the only
// segment that may contain a colon; it is therefore the middle segment and is
// reconstructed by re-joining everything between the first and last colon.
func routeRuleCompositeId(resourceGroupId, resourceId, routeRuleId string) string {
	return resourceGroupId + ":" + resourceId + ":" + routeRuleId
}

// parseRouteRuleId splits the composite id into its three parts. A non-composite id
// (for example a bare route rule id from an older state file) is rejected so the
// caller surfaces the missing identity instead of issuing a Read that cannot succeed.
func parseRouteRuleId(id string) (resourceGroupId, resourceId, routeRuleId string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid alicloud_dataworks_route_rule id %q: expected the format <resource_group_id>:<resource_id>:<route_rule_id>", id)
	}
	resourceGroupId = parts[0]
	resourceId = strings.Join(parts[1:len(parts)-1], ":")
	routeRuleId = parts[len(parts)-1]
	if resourceGroupId == "" || resourceId == "" || routeRuleId == "" {
		return "", "", "", fmt.Errorf("invalid alicloud_dataworks_route_rule id %q: every segment must be non-empty", id)
	}
	return resourceGroupId, resourceId, routeRuleId, nil
}

// routeRuleIdFromResponse converts the route rule id returned by CreateRouteRule (in
// the OpenApiResultModel data field, typed int64) into the canonical decimal string
// used in state, filters and delete/update requests. The tea RPC decoder renders
// numbers as float64, while some accounts return a numeric string; both are accepted.
// Containers, fractional or non-finite numbers and non-positive ids are contract
// drift and are rejected instead of being stringified into a broken id.
func routeRuleIdFromResponse(raw interface{}) (string, error) {
	switch v := raw.(type) {
	case nil:
		return "", fmt.Errorf("CreateRouteRule response data is missing")
	case string:
		if v == "" {
			return "", fmt.Errorf("CreateRouteRule response data is empty")
		}
		return routeRuleDecimalInteger(v)
	case json.Number:
		// The tea decoder sometimes returns numbers as json.Number, which
		// preserves large ids beyond the float64 lossless range; parse it as
		// a decimal integer instead of going through float64.
		return routeRuleDecimalInteger(v.String())
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return "", fmt.Errorf("route rule id %v is not finite", v)
		}
		if v != math.Trunc(v) {
			return "", fmt.Errorf("route rule id %v is not an integer", v)
		}
		if v <= 0 {
			return "", fmt.Errorf("route rule id %v is not positive", v)
		}
		if math.Abs(v) > maxExactFloatInteger {
			return "", fmt.Errorf("route rule id %v exceeds the lossless integer range", v)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return routeRulePositiveInt(int64(v))
	case int64:
		return routeRulePositiveInt(v)
	case uint64:
		if v == 0 || v > uint64(1<<63-1) {
			return "", fmt.Errorf("route rule id %d is not a positive int64", v)
		}
		return strconv.FormatUint(v, 10), nil
	default:
		return "", fmt.Errorf("route rule id has unexpected type %T", raw)
	}
}

func routeRuleDecimalInteger(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("route rule id is empty")
	}
	if routeRuleHasNonIntegerRune(s) {
		return "", fmt.Errorf("route rule id %q is not an integer", s)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("route rule id %q is not an int64: %w", s, err)
	}
	return routeRulePositiveInt(v)
}

func routeRulePositiveInt(v int64) (string, error) {
	if v <= 0 {
		return "", fmt.Errorf("route rule id %d is not positive", v)
	}
	return strconv.FormatInt(v, 10), nil
}

func routeRuleHasNonIntegerRune(s string) bool {
	for _, c := range s {
		if c == '.' || c == 'e' || c == 'E' || c == '+' {
			return true
		}
	}
	return false
}

// routeRuleLookup reads a response field by its camelCase backend name, falling back
// to the PascalCase form. The published CloudSpec backend names are camelCase
// (success, data, routeRuleId, ...), but a few dataworks-public responses render
// PascalCase, so both spellings are tried before treating the field as absent.
func routeRuleLookup(obj map[string]interface{}, camel, pascal string) interface{} {
	if obj == nil {
		return nil
	}
	if v, ok := obj[camel]; ok {
		return v
	}
	if v, ok := obj[pascal]; ok {
		return v
	}
	return nil
}

// routeRuleSuccess reports whether the OpenApiResultModel body marks the call as
// successful. A nil response (transport error before the body was parsed) is not
// successful, so the caller still has the RPC error to report.
func routeRuleSuccess(response map[string]interface{}) bool {
	if response == nil {
		return false
	}
	switch v := routeRuleLookup(response, "success", "Success").(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true")
	default:
		return false
	}
}

// routeRuleBusinessError rebuilds the error from the OpenApiResultModel
// errorCode/errorMessage fields when success is false. The shared formatError only
// recognizes Code/Message envelopes, so without this the Private RouteRule APIs
// would surface an empty code/message on failure. Returns nil when the response is
// nil or success is true, leaving a transport error to the caller.
func routeRuleBusinessError(response map[string]interface{}, action string) error {
	if response == nil || routeRuleSuccess(response) {
		return nil
	}
	code := routeRuleString(routeRuleLookup(response, "errorCode", "ErrorCode"))
	message := routeRuleString(routeRuleLookup(response, "errorMessage", "ErrorMessage"))
	if dynamicCode := routeRuleLookup(response, "dynamicErrorCode", "DynamicErrorCode"); dynamicCode != nil && code == "" {
		code = routeRuleString(dynamicCode)
	}
	if dynamicMessage := routeRuleLookup(response, "dynamicErrorMessage", "DynamicErrorMessage"); dynamicMessage != nil && message == "" {
		message = routeRuleString(dynamicMessage)
	}
	if code == "" && message == "" {
		// The body marks success=false but carries no usable error fields; keep a
		// non-empty message so the failure is visible instead of an empty error.
		return fmt.Errorf("calling %s failed: the response marked success=false but carried no errorCode/errorMessage", action)
	}
	if message == "" {
		return fmt.Errorf("calling %s failed: %s", action, code)
	}
	if code == "" {
		return fmt.Errorf("calling %s failed: %s", action, message)
	}
	return fmt.Errorf("calling %s failed: %s: %s", action, code, message)
}

// routeRuleString renders a response value as a canonical string. A nil field is
// distinct from an empty string: fmt.Sprint(nil) would produce "<nil>", which
// would leak into state as a destination value or mask a missing error field, so
// nil is folded to "" before the caller decides whether the field is present.
func routeRuleString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// routeRuleRetryable classifies transient RouteRule failures for the retry loop.
// The published tables are not exposed for these Private APIs, so only the shared
// transient classifier (throttling, socket/timeout, service-busy) and explicit
// internal-error tokens retry; every documented terminal answer is surfaced as a
// non-retryable error.
func routeRuleRetryable(err error) bool {
	if err == nil {
		return false
	}
	if helper.NeedRetry(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, token := range []string{"internalerror", "systembusy", "serviceunavailable", "request was rate limited"} {
		if strings.Contains(msg, token) {
			return true
		}
	}
	return false
}

// routeRuleNotFound reports whether a RouteRule error means the rule no longer
// exists. The Private APIs do not publish a stable not-found code, so an error is
// treated as not-found only when it carries an explicit not-found token, and a list
// result is treated as not-found only when the array is present and lacks the
// requested id (absent array is contract drift, not absence).
func routeRuleNotFound(err error, id string) bool {
	if err == nil {
		return false
	}
	if errs.IsTypedNotFound(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "notfound") || strings.Contains(msg, "not found") || strings.Contains(msg, "does not exist")
}

// routeRuleCallPost issues a RouteRule POST RPC, classifying the OpenApiResultModel
// business error and transient failures in one retry loop bounded by ctx.
func routeRuleCallPost(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var callErr error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(ctx, routeRuleWriteTimeout, func() *retry.RetryError {
		response, callErr = client.RpcPost(routeRuleProduct, routeRuleAPIVersion, action, nil, request, true)
		if busErr := routeRuleBusinessError(response, action); busErr != nil {
			callErr = busErr
		}
		logs.AddDebug(action, response, request)
		if callErr != nil {
			if routeRuleRetryable(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	if err != nil {
		return response, err
	}
	return response, nil
}

// routeRuleListCall issues ListRouteRules and returns the parsed data array. A nil
// response with no transport error, or a success=true body with a non-array data
// field, is contract drift and fails the read instead of being treated as empty.
func routeRuleListCall(ctx context.Context, client *connectivity.AliyunClient, resourceGroupId, resourceId string) ([]interface{}, map[string]interface{}, error) {
	action := routeRuleListAPI
	request := map[string]interface{}{
		"ResourceGroupId": resourceGroupId,
		"ResourceId":      resourceId,
		"RegionId":        client.RegionId,
	}
	var response map[string]interface{}
	var callErr error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(ctx, routeRuleReadTimeout, func() *retry.RetryError {
		response, callErr = client.RpcPost(routeRuleProduct, routeRuleAPIVersion, action, nil, request, true)
		if busErr := routeRuleBusinessError(response, action); busErr != nil {
			callErr = busErr
		}
		logs.AddDebug(action, response, request)
		if callErr != nil {
			if routeRuleRetryable(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	if err != nil {
		return nil, response, err
	}
	raw := routeRuleLookup(response, "data", "Data")
	list, ok := raw.([]interface{})
	if !ok {
		return nil, response, fmt.Errorf("reading %s: the response data is not a list (got %T); response: %v", action, raw, response)
	}
	return list, response, nil
}

// DescribeRouteRuleContext returns the route rule selected by its composite id from
// the ListRouteRules result. A present array that does not contain the requested
// routeRuleId means the rule is gone (NotFound); a missing array is contract drift.
func (s *Service) DescribeRouteRuleContext(ctx context.Context, id string) (map[string]interface{}, error) {
	resourceGroupId, resourceId, routeRuleId, err := parseRouteRuleId(id)
	if err != nil {
		return nil, err
	}
	list, response, err := routeRuleListCall(ctx, s.client, resourceGroupId, resourceId)
	if err != nil {
		if routeRuleNotFound(err, id) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("RouteRule", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, routeRuleListAPI, errs.AlibabaCloudSdkGoERROR)
	}
	for _, item := range list {
		entry, ok := item.(map[string]interface{})
		if !ok {
			return nil, errs.WrapErrorf(fmt.Errorf("unexpected ListRouteRules element type %T", item), errs.DefaultErrorMsg, id, routeRuleListAPI, errs.AlibabaCloudSdkGoERROR)
		}
		itemId, idErr := routeRuleIdFromResponse(routeRuleLookup(entry, "routeRuleId", "RouteRuleId"))
		if idErr != nil {
			return nil, errs.WrapErrorf(idErr, errs.DefaultErrorMsg, id, routeRuleListAPI, errs.AlibabaCloudSdkGoERROR)
		}
		if itemId == routeRuleId {
			return routeRuleNormalize(entry, resourceGroupId, resourceId, routeRuleId), nil
		}
	}
	return nil, errs.WrapErrorf(errs.NotFoundErr("RouteRule", id), errs.NotFoundMsg, response)
}

// DescribeRouteRule is the non-context wrapper kept for parity with the other
// DataWorks service methods; the Framework resource always passes ctx.
func (s *Service) DescribeRouteRule(id string) (map[string]interface{}, error) {
	return s.DescribeRouteRuleContext(context.Background(), id)
}

// DeleteRouteRuleContext deletes the route rule and is idempotent: a not-found
// error (the rule was already removed) is suppressed, while every other failure
// is returned. DeleteRouteRule addresses the rule by RouteRuleId only, so the
// composite id is split first.
func (s *Service) DeleteRouteRuleContext(ctx context.Context, id string) error {
	_, _, routeRuleId, err := parseRouteRuleId(id)
	if err != nil {
		return err
	}
	request := map[string]interface{}{
		"RouteRuleId": routeRuleId,
		"RegionId":    s.client.RegionId,
	}
	if _, err := routeRuleCallPost(ctx, s.client, routeRuleDeleteAPI, request); err != nil {
		if routeRuleNotFound(err, id) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, routeRuleDeleteAPI, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// routeRuleNormalize maps the ListRouteRules entry onto the response shape the
// resource state builder expects: every field is keyed by its attribute name and
// the composite id is filled in so the Read path produces a self-contained state.
func routeRuleNormalize(entry map[string]interface{}, resourceGroupId, resourceId, routeRuleId string) map[string]interface{} {
	return map[string]interface{}{
		"Id":                        routeRuleCompositeId(resourceGroupId, resourceId, routeRuleId),
		"RouteRuleId":               routeRuleId,
		"ResourceGroupId":           resourceGroupId,
		"ResourceId":                resourceId,
		"DestinationCidr":           routeRuleString(routeRuleLookup(entry, "destinationCidr", "DestinationCidr")),
		"DestinationType":           routeRuleString(routeRuleLookup(entry, "destinationType", "DestinationType")),
		"DestinationAddressingType": routeRuleString(routeRuleLookup(entry, "destinationAddressingType", "DestinationAddressingType")),
	}
}
