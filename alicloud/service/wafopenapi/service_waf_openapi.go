package wafopenapi

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

type WafOpenapiService struct {
	client *connectivity.AliyunClient
}

func NewWafOpenapiService(client *connectivity.AliyunClient) *WafOpenapiService {
	return &WafOpenapiService{client: client}
}

func (s *WafOpenapiService) DescribeWafv3Instance(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeInstance"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, true)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	if _, ok := v.(map[string]interface{})["InstanceId"]; !ok {
		return object, errs.WrapErrorf(errs.NotFoundErr("Wafv3Instance", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return v.(map[string]interface{}), nil
}

func (s *WafOpenapiService) Wafv3InstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeWafv3Instance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		status84 := object["Status"]
		for _, failState := range failStates {
			if fmt.Sprint(status84) == failState {
				return object, fmt.Sprint(status84), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(status84)))
			}
		}
		return object, fmt.Sprint(status84), nil
	}
}

func (s *WafOpenapiService) DescribeWafv3Domain(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Domain":     parts[1],
		"RegionId":   s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDomainDetail"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"Waf.Instance.ValidFaild"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	if _, ok := v.(map[string]interface{})["Domain"]; !ok {
		return object, errs.WrapErrorf(errs.NotFoundErr("Wafv3Domain", id), errs.NotFoundWithResponse, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *WafOpenapiService) Wafv3DomainStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeWafv3Domain(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		localVar75 := object["Status"]
		for _, failState := range failStates {
			if fmt.Sprint(localVar75) == failState {
				return object, fmt.Sprint(localVar75), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(localVar75)))
			}
		}
		return object, fmt.Sprint(localVar75), nil
	}
}

func (s *WafOpenapiService) DescribeWafDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDomain"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"Domain":     parts[1],
		"InstanceId": parts[0],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ComboError", "DomainNotExist"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("WafDomain", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Domain", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Domain", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
