package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DasService struct {
	client *connectivity.AliyunClient
}

func (s *DasService) DescribeInstanceDasPro(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstanceDasPro"

	conn, err := s.client.NewDasClient()
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-16"), StringPointer("AK"), nil, request, &runtime)
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

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *DasService) InstanceDasProStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstanceDasPro(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Data"]) == failState {
				return object, fmt.Sprint(object["Data"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Data"])))
			}
		}
		return object, fmt.Sprint(object["Data"]), nil
	}
}

func (s *DasService) DescribeDasSwitchDasPro(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetDasProServiceUsage"

	conn, err := s.client.NewDasClient()
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-16"), StringPointer("AK"), nil, request, &runtime)
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
		if fmt.Sprint(response["Code"]) == "-404" {
			return object, WrapErrorf(Error(GetNotFoundMessage("Das:SwitchDasPro", id)), NotFoundWithResponse, response)
		}
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	resp, err := jsonpath.Get("$.Data", response)
	if resp == nil && fmt.Sprint(response["Success"]) == "true" {
		return object, WrapErrorf(Error(GetNotFoundMessage("Das:SwitchDasPro", id)), NotFoundWithResponse, response)
	}
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = resp.(map[string]interface{})

	return object, nil
}
