// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type PolardbServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewPolardbServiceV2(client *connectivity.AliyunClient) *PolardbServiceV2 {
	return &PolardbServiceV2{client: client}
}

// DescribePolardbGlobalSecurityIpGroup <<< Encapsulated get interface for Polardb GlobalSecurityIpGroup.

func (s *PolardbServiceV2) DescribePolardbGlobalSecurityIpGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	action := "DescribeGlobalSecurityIPGroup"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.GlobalSecurityIPGroup[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("GlobalSecurityIpGroup", id), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("GlobalSecurityIpGroup", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["GlobalSecurityGroupId"]) != id {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("GlobalSecurityIpGroup", id), NotFoundMsg, response)
}

func (s *PolardbServiceV2) PolardbGlobalSecurityIpGroupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.PolardbGlobalSecurityIpGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribePolardbGlobalSecurityIpGroup)
}

func (s *PolardbServiceV2) PolardbGlobalSecurityIpGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolardbGlobalSecurityIpGroup >>> Encapsulated.
