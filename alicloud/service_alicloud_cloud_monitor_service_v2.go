package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudMonitorServiceServiceV2 struct {
	client *connectivity.AliyunClient
}

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceHybridDoubleWrite(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeHybridDoubleWrite"

	conn, err := s.client.NewCmsClient()
	if err != nil {
		return object, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"SourceNamespace": parts[0],
		"SourceUserId":    parts[1],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-08"), StringPointer("AK"), nil, request, &runtime)
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

	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	resp, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:HybridDoubleWrite", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["SourceNamespace"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["SourceUserId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:HybridDoubleWrite", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceEventRuleTargets(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeEventRuleTargetList"

	conn, err := s.client.NewCmsClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"RuleName": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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

	conn, err := s.client.NewCmsClient()
	if err != nil {
		return object, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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

	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	resp, err := jsonpath.Get("$.NodeProcesses.NodeProcess", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NodeProcesses.NodeProcess", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:MonitoringAgentProcess", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["ProcessId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:MonitoringAgentProcess", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CloudMonitorServiceServiceV2) DescribeCloudMonitorServiceGroupMonitoringAgentProcess(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeGroupMonitoringAgentProcess"

	conn, err := s.client.NewCmsClient()
	if err != nil {
		return object, WrapError(err)
	}

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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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

		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		resp, err := jsonpath.Get("$.Processes.Process", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Processes.Process", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:GroupMonitoringAgentProcess", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:GroupMonitoringAgentProcess", id)), NotFoundWithResponse, response)
	}

	return object, nil
}
