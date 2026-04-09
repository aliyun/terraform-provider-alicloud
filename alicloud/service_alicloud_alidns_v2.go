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

// DescribeAlidnsCloudGtmMonitorTemplate <<< Encapsulated get interface for Alidns CloudGtmMonitorTemplate.

func (s *AlidnsServiceV2) DescribeAlidnsCloudGtmMonitorTemplate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = id

	action := "DescribeCloudGtmMonitorTemplate"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)

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

	currentStatus := response["TemplateId"]
	if currentStatus == nil {
		return object, WrapErrorf(NotFoundErr("CloudGtmMonitorTemplate", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *AlidnsServiceV2) AlidnsCloudGtmMonitorTemplateStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.AlidnsCloudGtmMonitorTemplateStateRefreshFuncWithApi(id, field, failStates, s.DescribeAlidnsCloudGtmMonitorTemplate)
}

func (s *AlidnsServiceV2) AlidnsCloudGtmMonitorTemplateStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeAlidnsCloudGtmMonitorTemplate >>> Encapsulated.
