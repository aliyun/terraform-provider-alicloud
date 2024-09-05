package alicloud

import (
	"fmt"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	utilv2 "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type QuotasServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeQuotasTemplateQuota <<< Encapsulated get interface for Quotas TemplateQuota.

func (s *QuotasServiceV2) DescribeQuotasTemplateQuota(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := "ListQuotaApplicationTemplates"
	conn, err := client.NewQuotasClientV2()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})

	request["Id"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.CallApi(rpcParam(action, "POST", "2020-05-10"), &openapi.OpenApiRequest{Query: nil, Body: request, HostMap: nil}, &utilv2.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"TEMPLATE.ITEM.NOT.FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("TemplateQuota", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	response = response["body"].(map[string]interface{})

	v, err := jsonpath.Get("$.QuotaApplicationTemplates[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QuotaApplicationTemplates[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("TemplateQuota", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasTemplateQuotaStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasTemplateQuota(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
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
	conn, err := client.NewQuotasClientV2()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})

	request["ApplicationId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.CallApi(rpcParam(action, "POST", "2020-05-10"), &openapi.OpenApiRequest{Query: nil, Body: request, HostMap: nil}, &utilv2.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("QuotaApplication", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	response = response["body"].(map[string]interface{})
	v, err := jsonpath.Get("$.QuotaApplication", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QuotaApplication", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasQuotaApplicationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasQuotaApplication(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
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
	conn, err := client.NewQuotasClientV2()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	request["AlarmId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.CallApi(rpcParam(action, "POST", "2020-05-10"), &openapi.OpenApiRequest{Query: nil, Body: request, HostMap: nil}, &utilv2.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ALARM.NOT.FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("QuotaAlarm", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.QuotaAlarm", response["body"].(map[string]interface{}))
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QuotaAlarm", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasQuotaAlarmStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasQuotaAlarm(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
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
	conn, err := client.NewQuotasClientV2()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	request["BatchQuotaApplicationId"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.CallApi(rpcParam(action, "POST", "2020-05-10"), &openapi.OpenApiRequest{Query: nil, Body: request, HostMap: nil}, &utilv2.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	response = response["body"].(map[string]interface{})

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.QuotaBatchApplications[*]", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("TemplateApplications", id)), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("TemplateApplications", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *QuotasServiceV2) DescribeListQuotaApplicationsDetailForTemplate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListQuotaApplicationsDetailForTemplate"
	conn, err := client.NewQuotasClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BatchQuotaApplicationId"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *QuotasServiceV2) QuotasTemplateApplicationsStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasTemplateApplications(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
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
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 0, len(parts)))
	}
	action := "GetQuotaTemplateServiceStatus"
	conn, err := client.NewQuotasClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.TemplateServiceStatus", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TemplateServiceStatus", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *QuotasServiceV2) QuotasTemplateServiceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeQuotasTemplateService(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeQuotasTemplateService >>> Encapsulated.
