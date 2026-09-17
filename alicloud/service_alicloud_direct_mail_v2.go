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

type DirectMailServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeDirectMailDedicatedIpPool <<< Encapsulated get interface for DirectMail DedicatedIpPool.

func (s *DirectMailServiceV2) DescribeDirectMailDedicatedIpPool(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PoolId"] = id
	request["RegionId"] = client.RegionId
	action := "DedicatedIpPoolList"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

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

	v, err := jsonpath.Get("$.IpPools[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.IpPools[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("DedicatedIpPool", id), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["Id"]
	if fmt.Sprint(currentStatus) == "" {
		return object, WrapErrorf(NotFoundErr("DedicatedIpPool", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *DirectMailServiceV2) DirectMailDedicatedIpPoolStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.DirectMailDedicatedIpPoolStateRefreshFuncWithApi(id, field, failStates, s.DescribeDirectMailDedicatedIpPool)
}

func (s *DirectMailServiceV2) DirectMailDedicatedIpPoolStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeDirectMailDedicatedIpPool >>> Encapsulated.

// DescribeDirectMailConfigSet <<< Encapsulated get interface for DirectMail ConfigSet.

func (s *DirectMailServiceV2) DescribeDirectMailConfigSet(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = id
	request["RegionId"] = client.RegionId
	action := "ConfigSetDetail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

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

	currentStatus, err := jsonpath.Get("$.Detail.Id", response)
	if fmt.Sprint(currentStatus) == "" {
		return object, WrapErrorf(NotFoundErr("ConfigSet", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *DirectMailServiceV2) DirectMailConfigSetStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.DirectMailConfigSetStateRefreshFuncWithApi(id, field, failStates, s.DescribeDirectMailConfigSet)
}

func (s *DirectMailServiceV2) DirectMailConfigSetStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeDirectMailConfigSet >>> Encapsulated.
