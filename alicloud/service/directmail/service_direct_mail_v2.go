// Package directmail. This file is generated automatically. Please do not modify it manually, thank you!
package directmail

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type DirectMailServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewDirectMailServiceV2 constructs the DirectMailServiceV2 from a connectivity client.
func NewDirectMailServiceV2(client *connectivity.AliyunClient) *DirectMailServiceV2 {
	return &DirectMailServiceV2{client: client}
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.IpPools[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.IpPools[*]", response)
	}

	// Match the pool by exact ID instead of trusting the PoolId filter, so another
	// pool in the same scope is never adopted as the requested one.
	for _, item := range v.([]interface{}) {
		if pool, ok := item.(map[string]interface{}); ok && fmt.Sprint(pool["Id"]) == id {
			return pool, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("DedicatedIpPool", id), errs.NotFoundMsg, response)
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	currentStatus, err := jsonpath.Get("$.Detail.Id", response)
	if fmt.Sprint(currentStatus) == "" {
		return object, errs.WrapErrorf(errs.NotFoundErr("ConfigSet", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeDirectMailConfigSet >>> Encapsulated.
