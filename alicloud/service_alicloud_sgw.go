package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type SgwService struct {
	client *connectivity.AliyunClient
}

func (s *SgwService) DescribeCloudStorageGatewayStorageBundle(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHcsSgwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeStorageBundle"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"StorageBundleId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"StorageBundleNotExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CloudStorageGatewayStorageBundle", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SgwService) DescribeCloudStorageGatewayGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHcsSgwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeGateway"
	request := map[string]interface{}{
		"GatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway:Gateway", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SgwService) DescribeGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHcsSgwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeGateway"
	request := map[string]interface{}{
		"GatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway:Gateway", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *SgwService) DescribeCloudStorageGatewayGatewaySmbUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHcsSgwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeGatewaySMBUsers"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"GatewayId":  parts[0],
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
			if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway:Gateway", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.Users.User", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Users.User", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Username"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *SgwService) DescribeTasks(id, taskId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHcsSgwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeTasks"
	request := map[string]interface{}{
		"TargetId":   id,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	if taskId != "" {
		request["TaskId"] = taskId
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
			if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway:Gateway:Task", id+":"+taskId)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.Tasks.SimpleTask", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Tasks.SimpleTask", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway Task", id+":"+taskId)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if taskId != "" && fmt.Sprint(v.(map[string]interface{})["TaskId"]) == taskId {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway Task", id+":"+taskId)), NotFoundWithResponse, response)
	}
	return
}

func (s *SgwService) CloudStorageGatewayTaskStateRefreshFunc(id, taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTasks(id, taskId)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["StateCode"]) == failState {
				return object, fmt.Sprint(object["StateCode"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["StateCode"])))
			}
		}
		return object, fmt.Sprint(object["StateCode"]), nil
	}
}

func (s *SgwService) DescribeCloudStorageGatewayGatewayCacheDisk(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHcsSgwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeGatewayCaches"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"GatewayId": parts[0],
	}
	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway:Gateway:CacheDisk", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Caches.Cache", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Caches.Cache", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["CacheId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudStorageGateway", id)), NotFoundWithResponse, response)
	}
	return object, nil
}
