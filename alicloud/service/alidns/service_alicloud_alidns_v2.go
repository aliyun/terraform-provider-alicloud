// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alidns

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

type AlidnsServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewAlidnsServiceV2(client *connectivity.AliyunClient) *AlidnsServiceV2 {
	return &AlidnsServiceV2{client: client}
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
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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

	currentStatus := response["AddressPoolId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudGtmAddressPool", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *AlidnsServiceV2) AlidnsCloudGtmAddressPoolStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AlidnsCloudGtmAddressPoolStateRefreshFuncWithApi(id, field, failStates, s.DescribeAlidnsCloudGtmAddressPool)
}

func (s *AlidnsServiceV2) AlidnsCloudGtmAddressPoolStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
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
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[0]
	request["InstanceId"] = parts[1]

	action := "DescribeCloudGtmInstanceConfigFullInfo"
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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

	currentStatus := response["ConfigId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudGtmInstanceConfig", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *AlidnsServiceV2) AlidnsCloudGtmInstanceConfigStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AlidnsCloudGtmInstanceConfigStateRefreshFuncWithApi(id, field, failStates, s.DescribeAlidnsCloudGtmInstanceConfig)
}

func (s *AlidnsServiceV2) AlidnsCloudGtmInstanceConfigStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAlidnsCloudGtmInstanceConfig >>> Encapsulated.

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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)

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

	currentStatus := response["TemplateId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudGtmMonitorTemplate", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *AlidnsServiceV2) AlidnsCloudGtmMonitorTemplateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AlidnsCloudGtmMonitorTemplateStateRefreshFuncWithApi(id, field, failStates, s.DescribeAlidnsCloudGtmMonitorTemplate)
}

func (s *AlidnsServiceV2) AlidnsCloudGtmMonitorTemplateStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAlidnsCloudGtmMonitorTemplate >>> Encapsulated.

// DescribeAlidnsCloudGtmAddress <<< Encapsulated get interface for Alidns CloudGtmAddress.

func (s *AlidnsServiceV2) DescribeAlidnsCloudGtmAddress(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressId"] = id

	action := "DescribeCloudGtmAddress"
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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

	currentStatus := response["AddressId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudGtmAddress", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *AlidnsServiceV2) AlidnsCloudGtmAddressStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AlidnsCloudGtmAddressStateRefreshFuncWithApi(id, field, failStates, s.DescribeAlidnsCloudGtmAddress)
}

func (s *AlidnsServiceV2) AlidnsCloudGtmAddressStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAlidnsCloudGtmAddress >>> Encapsulated.
