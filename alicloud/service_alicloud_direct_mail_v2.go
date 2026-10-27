// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

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

	v, err := jsonpath.Get("$.IpPools[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.IpPools[*]", response)
	}

	// Match the pool by exact ID instead of trusting the PoolId filter, so another
	// pool in the same scope is never adopted as the requested one.
	for _, item := range v.([]interface{}) {
		if pool, ok := item.(map[string]interface{}); ok && fmt.Sprint(pool["Id"]) == id {
			return pool, nil
		}
	}

	return object, WrapErrorf(NotFoundErr("DedicatedIpPool", id), NotFoundMsg, response)
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
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

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

	currentStatus, err := jsonpath.Get("$.Detail.Id", response)
	if fmt.Sprint(currentStatus) == "" {
		return object, WrapErrorf(NotFoundErr("ConfigSet", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeDirectMailConfigSet >>> Encapsulated.
