package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudMonitorServiceServiceV2 struct {
	client *connectivity.AliyunClient
}

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceHybridDoubleWrite(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeHybridDoubleWrite"
	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"SourceNamespace": parts[0],
		"SourceUserId":    parts[1],
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2018-03-08", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("CloudMonitorService:HybridDoubleWrite", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["SourceNamespace"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["SourceUserId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("CloudMonitorService:HybridDoubleWrite", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceEventRuleTargets(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeEventRuleTargetList"
	client := s.client

	request := map[string]interface{}{
		"RuleName": id,
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceMonitoringAgentProcess(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeMonitoringAgentProcesses"
	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.NodeProcesses.NodeProcess", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NodeProcesses.NodeProcess", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("CloudMonitorService:MonitoringAgentProcess", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["ProcessId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("CloudMonitorService:MonitoringAgentProcess", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceGroupMonitoringAgentProcess(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeGroupMonitoringAgentProcess"
	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"GroupId":    parts[0],
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Processes.Process", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Processes.Process", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("CloudMonitorService:GroupMonitoringAgentProcess", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["GroupId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["Id"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("CloudMonitorService:GroupMonitoringAgentProcess", id), NotFoundWithResponse, response)
	}

	return object, nil
}

// DescribeCloudMonitorServiceBasicPublic <<< Encapsulated get interface for CloudMonitorService BasicPublic.

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceBasicPublic(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var endpoint string
	var query map[string]interface{}
	action := "QueryAvailableInstances"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = id

	request["SubscriptionType"] = "PayAsYouGo"
	request["ProductCode"] = "cms"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("BasicPublic", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["InstanceID"] != id {
			continue
		}
		if item["SubscriptionType"] != "PayAsYouGo" {
			continue
		}
		if item["ProductCode"] != "cms" {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("BasicPublic", id), NotFoundMsg, response)
}

func (s *CloudMonitorServiceServiceV2) CloudMonitorServiceBasicPublicStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudMonitorServiceBasicPublic(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeCloudMonitorServiceBasicPublic >>> Encapsulated.

// DescribeCloudMonitorServiceEnterprisePublic <<< Encapsulated get interface for CloudMonitorService EnterprisePublic.

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceEnterprisePublic(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = id

	request["SubscriptionType"] = "PayAsYouGo"
	request["ProductCode"] = "cms"
	request["ProductType"] = "cms_enterprise_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "cms_enterprise_public_intl"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "cms_enterprise_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("EnterprisePublic", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["InstanceID"] != id {
			continue
		}
		if item["SubscriptionType"] != "PayAsYouGo" {
			continue
		}
		if item["ProductCode"] != "cms" {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("EnterprisePublic", id), NotFoundMsg, response)
}

func (s *CloudMonitorServiceServiceV2) CloudMonitorServiceEnterprisePublicStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudMonitorServiceEnterprisePublic(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeCloudMonitorServiceEnterprisePublic >>> Encapsulated.

// DescribeCloudMonitorServiceNaamPublic <<< Encapsulated get interface for CloudMonitorService NaamPublic.

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceNaamPublic(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var endpoint string
	var query map[string]interface{}
	action := "QueryAvailableInstances"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = id

	request["SubscriptionType"] = "PayAsYouGo"
	request["ProductCode"] = "cms"
	request["ProductType"] = "cms_naam_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "cms_naam_public_intl"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "cms_naam_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("NaamPublic", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["InstanceID"] != id {
			continue
		}
		if item["SubscriptionType"] != "PayAsYouGo" {
			continue
		}
		if item["ProductCode"] != "cms" {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("NaamPublic", id), NotFoundMsg, response)
}

// DescribeCloudMonitorServiceNaamPublic >>> Encapsulated.

// DescribeCloudMonitorServiceSiteMonitor <<< Encapsulated get interface for CloudMonitorService SiteMonitor.

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceSiteMonitor(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TaskId"] = id

	action := "DescribeSiteMonitorAttribute"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden"}) {
			return object, WrapErrorf(NotFoundErr("SiteMonitor", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SiteMonitors", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SiteMonitors", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *CloudMonitorServiceServiceV2) CloudMonitorServiceSiteMonitorStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CloudMonitorServiceSiteMonitorStateRefreshFuncWithApi(id, field, failStates, s.DescribeCloudMonitorServiceSiteMonitor)
}

func (s *CloudMonitorServiceServiceV2) CloudMonitorServiceSiteMonitorStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeCloudMonitorServiceSiteMonitor >>> Encapsulated.
