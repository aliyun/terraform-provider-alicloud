// Package quotas. This file is generated automatically. Please do not modify it manually, thank you!
package quotas

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type QuotasService struct {
	client *connectivity.AliyunClient
}

func NewQuotasService(client *connectivity.AliyunClient) *QuotasService {
	return &QuotasService{client: client}
}

type QuotasServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewQuotasServiceV2(client *connectivity.AliyunClient) *QuotasServiceV2 {
	return &QuotasServiceV2{client: client}
}

func (s *QuotasService) DescribeQuotasQuotaAlarm(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetQuotaAlarm"
	request := map[string]interface{}{
		"SourceIp": s.client.SourceIp,
		"RegionId": s.client.RegionId,
		"AlarmId":  id,
	}
	response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.QuotaAlarm", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.QuotaAlarm", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuotasService) DescribeQuotasQuotaApplication(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetQuotaApplication"
	request := map[string]interface{}{
		"SourceIp":      s.client.SourceIp,
		"RegionId":      s.client.RegionId,
		"ApplicationId": id,
	}
	response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.QuotaApplication", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.QuotaApplication", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuotasServiceV2) DescribeQuotasTemplateQuota(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "ListQuotaApplicationTemplates"
	request = make(map[string]interface{})

	request["Id"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"TEMPLATE.ITEM.NOT.FOUND"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("TemplateQuota", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.QuotaApplicationTemplates[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.QuotaApplicationTemplates[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("TemplateQuota", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasTemplateQuotaStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasTemplateQuota(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeQuotasTemplateQuota >>> Encapsulated.

// DescribeQuotasQuotaApplication <<< Encapsulated get interface for Quotas QuotaApplication.

func (s *QuotasServiceV2) DescribeQuotasQuotaApplication(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "GetQuotaApplication"
	request = make(map[string]interface{})

	request["ApplicationId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("QuotaApplication", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.QuotaApplication", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.QuotaApplication", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasQuotaApplicationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasQuotaApplication(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeQuotasQuotaApplication >>> Encapsulated.

// DescribeQuotasQuotaAlarm <<< Encapsulated get interface for Quotas QuotaAlarm.

func (s *QuotasServiceV2) DescribeQuotasQuotaAlarm(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "GetQuotaAlarm"
	request = make(map[string]interface{})
	request["AlarmId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ALARM.NOT.FOUND"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("QuotaAlarm", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.QuotaAlarm", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.QuotaAlarm", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasQuotaAlarmStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasQuotaAlarm(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeQuotasQuotaAlarm >>> Encapsulated.

// DescribeQuotasTemplateApplications <<< Encapsulated get interface for Quotas TemplateApplications.

func (s *QuotasServiceV2) DescribeQuotasTemplateApplications(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "ListQuotaApplicationsForTemplate"
	request = make(map[string]interface{})
	request["BatchQuotaApplicationId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.QuotaBatchApplications[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("TemplateApplications", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("TemplateApplications", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *QuotasServiceV2) DescribeListQuotaApplicationsDetailForTemplate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListQuotaApplicationsDetailForTemplate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BatchQuotaApplicationId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("quotas", "2020-05-10", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *QuotasServiceV2) QuotasTemplateApplicationsStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasTemplateApplications(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.Dimensions" {
			e := jsonata.MustCompile("$each($.Dimensions, function($v, $k) {{\"value\":$v, \"key\": $k}})[]")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeQuotasTemplateApplications >>> Encapsulated.
// DescribeQuotasTemplateService <<< Encapsulated get interface for Quotas TemplateService.

func (s *QuotasServiceV2) DescribeQuotasTemplateService(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 0 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 0, len(parts)))
	}
	action := "GetQuotaTemplateServiceStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("quotas", "2020-05-10", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.TemplateServiceStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TemplateServiceStatus", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasTemplateServiceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasTemplateService(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeQuotasTemplateService >>> Encapsulated.
