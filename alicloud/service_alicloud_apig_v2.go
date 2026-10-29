package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

type ApigServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeApigHttpApi <<< Encapsulated get interface for Apig HttpApi.

func (s *ApigServiceV2) DescribeApigHttpApi(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	httpApiId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("HttpApi", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, WrapErrorf(NotFoundErr("HttpApi", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigHttpApiStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigHttpApiStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigHttpApi)
}

func (s *ApigServiceV2) ApigHttpApiStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigHttpApi >>> Encapsulated.

// DescribeApigDomain <<< Encapsulated get interface for Apig Domain.

func (s *ApigServiceV2) DescribeApigDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	domainId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/domains/%s", domainId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Domain", id), NotFoundMsg, response)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, WrapErrorf(NotFoundErr("Domain", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigDomainStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigDomainStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigDomain)
}

func (s *ApigServiceV2) ApigDomainStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigDomain >>> Encapsulated.

// DescribeApigGateway <<< Encapsulated get interface for Apig Gateway.

func (s *ApigServiceV2) DescribeApigGateway(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	gatewayId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/gateways/%s", gatewayId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Gateway", id), NotFoundMsg, response)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"NotFound.GatewayNotFound", "Conflict.GatewayIsDeleted"}) {
		return object, WrapErrorf(NotFoundErr("Gateway", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) DescribeGatewayQueryAvailableInstances(d *schema.ResourceData) (object map[string]interface{}, err error) {
	return s.DescribeGatewayQueryAvailableInstancesByID(d.Id(), 1*time.Minute)
}

const (
	// apigGatewayRenewalPageSize requests a whole page of BssOpenApi results; 300
	// is the documented maximum for QueryAvailableInstances.
	apigGatewayRenewalPageSize = 300
)

// DescribeGatewayQueryAvailableInstancesByID queries the billing attributes of one
// instance from BssOpenApi. Every reported page is walked because the InstanceIDs
// filter is not guaranteed to be honoured; the walk ends on an empty page or a
// matching short page, and retryTimeout bounds it as a whole so no page count is
// assumed. An empty or incomplete result is retried within retryTimeout because
// BssOpenApi may not expose a freshly created instance yet, while a request error
// is always returned as is instead of being reported as a missing instance.
func (s *ApigServiceV2) DescribeGatewayQueryAvailableInstancesByID(id string, retryTimeout time.Duration) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var instance map[string]interface{}
	request := make(map[string]interface{})
	query := make(map[string]interface{})
	// InstanceIDs is kept as the cheap server-side filter, but the response is
	// matched by InstanceID and every page is walked because the filter is not
	// guaranteed to be honoured.
	query["InstanceIDs"] = id
	query["PageSize"] = apigGatewayRenewalPageSize

	var endpoint string
	request["ProductCode"] = "apigateway"
	action := "QueryAvailableInstances"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	// The walk is bounded by the caller's overall budget instead of an arbitrary
	// page count, so a result set of any size is still searched completely.
	deadline := time.Now().Add(retryTimeout)
	err = retry.Retry(retryTimeout, func() *retry.RetryError {
		var totalCount int
		walked := 0
		for page := 1; ; page++ {
			query["PageNum"] = page
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}

			var found bool
			var entries []interface{}
			var pageErr error
			instance, entries, totalCount, found, pageErr = apigGatewayRenewalInstancePage(id, response)
			if pageErr != nil {
				wait()
				return retry.RetryableError(pageErr)
			}
			if found {
				return nil
			}
			walked += len(entries)
			// An empty page is the contract-free end of the result set.
			if len(entries) == 0 {
				break
			}
			// TotalCount is only honoured once the walk has passed more than one
			// page worth of results and the accumulated count matches it exactly,
			// so a stale or bogus total can never cut the walk short.
			if totalCount > apigGatewayRenewalPageSize && walked == totalCount && len(entries) < apigGatewayRenewalPageSize {
				break
			}
			// A total that never ends must not keep the provider paging past the
			// caller's budget.
			if !time.Now().Before(deadline) {
				wait()
				return retry.RetryableError(fmt.Errorf("BssOpenApi QueryAvailableInstances response does not contain apig gateway %s within %d result(s) after %d page(s), stopped at the %s deadline", id, totalCount, page, retryTimeout))
			}
		}
		wait()
		return retry.RetryableError(fmt.Errorf("BssOpenApi QueryAvailableInstances response does not contain apig gateway %s within %d result(s)", id, totalCount))
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return instance, nil
}

// apigGatewayRenewalTotalCount reads the TotalCount of a QueryAvailableInstances
// response. A missing or unparsable value is reported as 0 so the caller keeps
// the single-page behaviour instead of failing the whole read.
func apigGatewayRenewalTotalCount(response map[string]interface{}) int {
	raw, err := jsonpath.Get("$.Data.TotalCount", response)
	if err != nil || raw == nil {
		return 0
	}
	switch value := raw.(type) {
	case float64:
		return int(value)
	case int:
		return value
	case int64:
		return int(value)
	case json.Number:
		// The SDK decodes response numbers with UseNumber, so TotalCount is a
		// json.Number rather than a float64 when it comes from the real API.
		total, convErr := value.Int64()
		if convErr != nil {
			return 0
		}
		return int(total)
	case string:
		total, convErr := strconv.Atoi(strings.TrimSpace(value))
		if convErr != nil {
			return 0
		}
		return total
	}
	return 0
}

// apigGatewayRenewalInstancePage extracts the target entry from one
// QueryAvailableInstances page. found=false with err=nil means the page is
// well-formed but does not contain the target, so the caller walks the next page.
func apigGatewayRenewalInstancePage(id string, response map[string]interface{}) (instance map[string]interface{}, entries []interface{}, totalCount int, found bool, err error) {
	if response == nil {
		return nil, nil, 0, false, fmt.Errorf("empty response returned by the BssOpenApi QueryAvailableInstances API for apig gateway %s", id)
	}
	totalCount = apigGatewayRenewalTotalCount(response)

	raw, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return nil, nil, totalCount, false, fmt.Errorf("failed to get attribute $.Data.InstanceList[*] for apig gateway %s: %s", id, err)
	}
	instances, ok := raw.([]interface{})
	if !ok {
		return nil, nil, totalCount, false, fmt.Errorf("unexpected BssOpenApi QueryAvailableInstances response for apig gateway %s: InstanceList is %v", id, raw)
	}
	for _, rawInstance := range instances {
		entry, ok := rawInstance.(map[string]interface{})
		if !ok {
			return nil, instances, totalCount, false, fmt.Errorf("unexpected BssOpenApi QueryAvailableInstances instance entry for apig gateway %s: %v", id, rawInstance)
		}
		// The InstanceIDs filter is not guaranteed to be honoured, so the entry
		// must match the gateway being refreshed instead of trusting list order.
		if fmt.Sprint(entry["InstanceID"]) == id {
			return entry, instances, totalCount, true, nil
		}
	}
	return nil, instances, totalCount, false, nil
}

// apigGatewayRenewalInstance extracts the single instance entry of a BssOpenApi
// QueryAvailableInstances response. Request errors take precedence over an empty
// response and neither of them is reported as a missing instance.
func apigGatewayRenewalInstance(id string, response map[string]interface{}, responseErr error) (map[string]interface{}, error) {
	if responseErr != nil {
		return nil, responseErr
	}
	instance, entries, _, found, err := apigGatewayRenewalInstancePage(id, response)
	if err != nil {
		return nil, err
	}
	if found {
		return instance, nil
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("apig gateway %s is not visible in the BssOpenApi QueryAvailableInstances response yet", id)
	}
	return nil, fmt.Errorf("BssOpenApi QueryAvailableInstances response does not contain apig gateway %s: %v", id, entries)
}

// apigGatewayRenewalConverged reports whether the BssOpenApi view of the renewal
// settings already matches the values requested by SetRenewal. Renewal duration
// and its unit are only meaningful for AutoRenewal.
func apigGatewayRenewalConverged(renewalRaw map[string]interface{}, expected map[string]interface{}) bool {
	expectedStatus, ok := expected["RenewalStatus"]
	if !ok {
		return true
	}
	if fmt.Sprint(renewalRaw["RenewStatus"]) != fmt.Sprint(expectedStatus) {
		return false
	}
	if fmt.Sprint(expectedStatus) != "AutoRenewal" {
		return true
	}
	if expectedPeriod, ok := expected["RenewalPeriod"]; ok && fmt.Sprint(renewalRaw["RenewalDuration"]) != fmt.Sprint(expectedPeriod) {
		return false
	}
	if expectedUnit, ok := expected["RenewalPeriodUnit"]; ok && fmt.Sprint(renewalRaw["RenewalDurationUnit"]) != fmt.Sprint(expectedUnit) {
		return false
	}
	return true
}

func waitForApigGatewayRenewalConvergence(id string, expected map[string]interface{}, timeout time.Duration, wait func(), describe func() (map[string]interface{}, error)) error {
	deadline := time.Now().Add(timeout)
	for {
		renewalRaw, err := describe()
		if err == nil && apigGatewayRenewalConverged(renewalRaw, expected) {
			return nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return fmt.Errorf("timeout waiting for the renewal settings of apig gateway %s to converge: expected %v, but the last query failed: %w", id, expected, err)
			}
			return fmt.Errorf("timeout waiting for the renewal settings of apig gateway %s to converge: expected %v, but observed status %q, renewal duration %v, renewal duration unit %q", id, expected, renewalRaw["RenewStatus"], renewalRaw["RenewalDuration"], renewalRaw["RenewalDurationUnit"])
		}
		wait()
	}
}

func (s *ApigServiceV2) ApigGatewayStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigGatewayStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigGateway)
}

func (s *ApigServiceV2) ApigGatewayStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		object["chargeType"] = convertApigGatewaydatachargeTypeResponse(object["chargeType"])
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigGateway >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Apig.
func (s *ApigServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/v1/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["ResourceId"] = StringPointer(convertListToJsonString(expandSingletonToList(d.Id())))
			query["RegionId"] = StringPointer(client.RegionId)
			query["TagKey"] = StringPointer(convertListToJsonString(convertListStringToListInterface(removedTagKeys)))
			query["ResourceType"] = StringPointer(resourceType)
			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = fmt.Sprintf("/v1/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})

			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["value"] = value
				tagsMap["key"] = key
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			request["tag"] = tagsMaps

			request["resourceType"] = "gateway"
			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "resourceId.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeApigEnvironment <<< Encapsulated get interface for Apig Environment.

func (s *ApigServiceV2) DescribeApigEnvironment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	environmentId := id
	action := fmt.Sprintf("/v1/environments/%s", environmentId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["environmentId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.EnvironmentNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Environment", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigEnvironmentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigEnvironment(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigEnvironment >>> Encapsulated.

// DescribeApigService <<< Encapsulated get interface for Apig Service.

func (s *ApigServiceV2) DescribeApigService(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	serviceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/services/%s", serviceId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Service", id), NotFoundMsg, response)
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.ServiceNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Service", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"NotFound.ServiceNotFound"}) {
		return object, WrapErrorf(NotFoundErr("Service", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigServiceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigServiceStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigService)
}

func (s *ApigServiceV2) ApigServiceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigService >>> Encapsulated.

// DescribeApigPlugin <<< Encapsulated get interface for Apig Plugin.

func (s *ApigServiceV2) DescribeApigPlugin(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/plugins")

	gatewayTypes := []string{"API", "AI"}

	for _, gatewayType := range gatewayTypes {
		query["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
		query["pageNumber"] = StringPointer("1")
		query["gatewayType"] = StringPointer(gatewayType)
		for {
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(1*time.Minute, func() *retry.RetryError {
				response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)

			if err != nil {
				if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
					return object, WrapErrorf(NotFoundErr("Plugin", id), NotFoundMsg, err)
				}
				return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			}

			resp, err := jsonpath.Get("$.data.items", response)
			if err != nil {
				return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data.items", response)
			}

			for _, v := range resp.([]interface{}) {
				if fmt.Sprint(v.(map[string]interface{})["pluginId"]) == id {
					return v.(map[string]interface{}), nil
				}
			}

			if len(resp.([]interface{})) < PageSizeLarge {
				break
			}

			pageNumber, err := strconv.Atoi(*query["pageNumber"])
			if err != nil {
				return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			}

			query["pageNumber"] = StringPointer(strconv.Itoa(pageNumber + 1))
		}
	}

	return nil, WrapErrorf(NotFoundErr("Plugin", id), NotFoundMsg, response)
}

func (s *ApigServiceV2) ApigPluginStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigPluginStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigPlugin)
}

func (s *ApigServiceV2) ApigPluginStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPlugin >>> Encapsulated.
// DescribeApigPluginClass <<< Encapsulated get interface for Apig PluginClass.

func (s *ApigServiceV2) DescribeApigPluginClass(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	pluginClassId := id
	action := fmt.Sprintf("/v1/plugin-classes/%s", pluginClassId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["pluginClassId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("PluginClass", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigPluginClassStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigPluginClass(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPluginClass >>> Encapsulated.
// DescribeApigOperation <<< Encapsulated get interface for Apig Operation.

func (s *ApigServiceV2) DescribeApigOperation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	httpApiId := parts[0]
	operationId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/operations/%s", httpApiId, operationId)
	request = make(map[string]interface{})
	query = make(map[string]*string)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Operation", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(NotFoundErr("Operation", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigOperationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigOperation(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigOperation >>> Encapsulated.
// DescribeApigApiAttachment <<< Encapsulated get interface for Apig ApiAttachment.

func (s *ApigServiceV2) DescribeApigApiAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	httpApiId := parts[0]
	action := fmt.Sprintf("/v1/http-apis/%s/attachment", httpApiId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["environmentId"] = StringPointer(parts[2])
	query["routeId"] = StringPointer(parts[1])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Gateway", id), NotFoundMsg, err)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	currentStatus := v.(map[string]interface{})["deployStatus"]
	if currentStatus == "NotDeployed" {
		return object, WrapErrorf(NotFoundErr("ApiAttachment", id), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigApiAttachmentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigApiAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigApiAttachment >>> Encapsulated.

// DescribeApigRoute <<< Encapsulated get interface for Apig Route.

func (s *ApigServiceV2) DescribeApigRoute(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	httpApiId := parts[0]
	routeId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/http-apis/%s/routes/%s", httpApiId, routeId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Route", id), NotFoundMsg, response)
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.RouteNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Route", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigRouteStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigRouteStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigRoute)
}

func (s *ApigServiceV2) ApigRouteStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigRoute >>> Encapsulated.

// DescribeApigAiModelProvider <<< Encapsulated get interface for Apig AiModelProvider.

func (s *ApigServiceV2) DescribeApigAiModelProvider(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	modelProviderId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/ai-model-providers/%s", modelProviderId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("AiModelProvider", id), NotFoundMsg, response)
	}
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigAiModelProviderStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigAiModelProviderStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigAiModelProvider)
}

func (s *ApigServiceV2) ApigAiModelProviderStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigAiModelProvider >>> Encapsulated.

// DescribeApigPolicy <<< Encapsulated get interface for Apig Policy.

func (s *ApigServiceV2) DescribeApigPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	policyId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v2/policies/%s", policyId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Policy", id), NotFoundMsg, response)
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Policy", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, WrapErrorf(NotFoundErr("Policy", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigPolicyStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigPolicyStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigPolicy)
}

func (s *ApigServiceV2) ApigPolicyStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPolicy >>> Encapsulated.

// DescribeApigSource <<< Encapsulated get interface for Apig Source.

func (s *ApigServiceV2) DescribeApigSource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	sourceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/sources/%s", sourceId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Source", id), NotFoundMsg, response)
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.ServiceSourceNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Source", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"NotFound.ServiceSourceNotFound"}) {
		return object, WrapErrorf(NotFoundErr("Source", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigSourceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigSourceStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigSource)
}

func (s *ApigServiceV2) ApigSourceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigSource >>> Encapsulated.

// DescribeApigSecret <<< Encapsulated get interface for Apig Secret.

func (s *ApigServiceV2) DescribeApigSecret(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	secretId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/secrets/%s", secretId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Secret", id), NotFoundMsg, response)
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.Secret"}) {
			return object, WrapErrorf(NotFoundErr("Secret", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigSecretStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigSecretStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigSecret)
}

func (s *ApigServiceV2) ApigSecretStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigSecret >>> Encapsulated.
