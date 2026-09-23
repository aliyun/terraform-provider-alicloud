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

type DasServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeDasSqlLogConfig <<< Encapsulated get interface for Das SqlLogConfig.

func (s *DasServiceV2) DescribeDasSqlLogConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id

	action := "DescribeSqlLogConfig"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("DAS", "2020-01-16", action, query, request, true)

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
	code, _ := jsonpath.Get("$.Data", response)
	if InArray(fmt.Sprint(code), []string{"-404"}) {
		return object, WrapErrorf(NotFoundErr("SqlLogConfig", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *DasServiceV2) DasSqlLogConfigStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.DasSqlLogConfigStateRefreshFuncWithApi(id, field, failStates, s.DescribeDasSqlLogConfig)
}

func (s *DasServiceV2) DasSqlLogConfigStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeDasSqlLogConfig >>> Encapsulated.
