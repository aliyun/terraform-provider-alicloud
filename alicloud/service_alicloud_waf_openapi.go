package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type WafOpenapiService struct {
	client *connectivity.AliyunClient
}

func (s *WafOpenapiService) DescribeWafv3Instance(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeInstance"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	if _, ok := v.(map[string]interface{})["InstanceId"]; !ok {
		return object, WrapErrorf(NotFoundErr("Wafv3Instance", id), NotFoundMsg, ProviderERROR)
	}

	return v.(map[string]interface{}), nil
}

func (s *WafOpenapiService) Wafv3InstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeWafv3Instance(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		status84 := object["Status"]
		for _, failState := range failStates {
			if fmt.Sprint(status84) == failState {
				return object, fmt.Sprint(status84), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(status84)))
			}
		}
		return object, fmt.Sprint(status84), nil
	}
}

func (s *WafOpenapiService) DescribeWafv3Domain(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Domain":     parts[1],
		"RegionId":   s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDomainDetail"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Waf.Instance.ValidFaild"}) {
			return object, WrapErrorf(NotFoundErr("Domain", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	if _, ok := v.(map[string]interface{})["Domain"]; !ok {
		return object, WrapErrorf(NotFoundErr("Wafv3Domain", id), NotFoundWithResponse, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *WafOpenapiService) Wafv3DomainStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeWafv3Domain(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		localVar75 := object["Status"]
		for _, failState := range failStates {
			if fmt.Sprint(localVar75) == failState {
				return object, fmt.Sprint(localVar75), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(localVar75)))
			}
		}
		return object, fmt.Sprint(localVar75), nil
	}
}
