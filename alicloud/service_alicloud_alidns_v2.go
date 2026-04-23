// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AlidnsServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeAlidnsCloudGtmAddressPool <<< Encapsulated get interface for Alidns CloudGtmAddressPool.

func (s *AlidnsServiceV2) DescribeAlidnsCloudGtmAddressPool(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressPoolId"] = id

	action := "DescribeCloudGtmAddressPool"
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	currentStatus := response["AddressPoolId"]
	if currentStatus == nil {
		return object, WrapErrorf(NotFoundErr("CloudGtmAddressPool", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *AlidnsServiceV2) AlidnsCloudGtmAddressPoolStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.AlidnsCloudGtmAddressPoolStateRefreshFuncWithApi(id, field, failStates, s.DescribeAlidnsCloudGtmAddressPool)
}

func (s *AlidnsServiceV2) AlidnsCloudGtmAddressPoolStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeAlidnsCloudGtmAddressPool >>> Encapsulated.

// DescribeAlidnsCloudGtmInstanceConfig <<< Encapsulated get interface for Alidns CloudGtmInstanceConfig.

func (s *AlidnsServiceV2) DescribeAlidnsCloudGtmInstanceConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[0]
	request["InstanceId"] = parts[1]

	action := "DescribeCloudGtmInstanceConfigFullInfo"
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	currentStatus := response["ConfigId"]
	if currentStatus == nil {
		return object, WrapErrorf(NotFoundErr("CloudGtmInstanceConfig", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *AlidnsServiceV2) AlidnsCloudGtmInstanceConfigStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.AlidnsCloudGtmInstanceConfigStateRefreshFuncWithApi(id, field, failStates, s.DescribeAlidnsCloudGtmInstanceConfig)
}

func (s *AlidnsServiceV2) AlidnsCloudGtmInstanceConfigStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeAlidnsCloudGtmInstanceConfig >>> Encapsulated.
