package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type BssOpenApiService struct {
	client *connectivity.AliyunClient
}

func (s *BssOpenApiService) QueryAvailableInstances(id, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBssopenapiClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"InstanceIDs": id,
		"ProductCode": productCode,
		"ProductType": productType,
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable", "SignatureDoesNotMatch"}) {
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				request["ProductCode"] = productCodeIntl
				request["ProductType"] = productTypeIntl
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp, _ := jsonpath.Get("$.Data.InstanceList", response)
		if len(resp.([]interface{})) < 1 {
			request["ProductCode"] = productCodeIntl
			request["ProductType"] = productTypeIntl
			conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage(productCode+"Instance", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage(productCode+"Instance", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *BssOpenApiService) QueryAvailableInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBssopenapiClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"InstanceIDs": id,
		"ProductCode": "vipcloudfw",
		"ProductType": "vipcloudfw",
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "cfw"
				request["ProductType"] = "cfw_pre_intl"
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		resp, _ := jsonpath.Get("$.Data.InstanceList", response)
		if len(resp.([]interface{})) < 1 {
			request["ProductCode"] = "cfw"
			request["ProductType"] = "cfw_pre_intl"
			conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)

			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
		}

		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "Success" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}

	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundWithResponse, response)
		}
	}

	object = v.([]interface{})[0].(map[string]interface{})

	return object, nil
}
