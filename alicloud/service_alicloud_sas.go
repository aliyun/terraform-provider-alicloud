package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type SasService struct {
	client *connectivity.AliyunClient
}

func (s *SasService) DescribeAllGroups(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAllGroups"
	request := map[string]interface{}{}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Groups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Groups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("SecurityCenter", id)), NotFoundWithResponse, response)
	}
	for _, vv := range v.([]interface{}) {
		if fmt.Sprint(vv.(map[string]interface{})["GroupId"]) == id {
			idExist = true
			return vv.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("SecurityCenter", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *SasService) DescribeSecurityCenterServiceLinkedRole(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeServiceLinkedRoleStatus"
	request := map[string]interface{}{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Sas", "2018-12-03", action, request, nil)
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
	v, err := jsonpath.Get("$.RoleStatus", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RoleStatus", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SasService) SecurityCenterServiceLinkedRoleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSecurityCenterServiceLinkedRole(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *SasService) DescribeThreatDetectionWebLockConfig(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"Uuid": id,
	}

	var response map[string]interface{}
	action := "DescribeWebLockConfigList"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ConfigList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ConfigList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("WebLockConfig", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})

	return object, nil
}

func (s *SasService) DescribeThreatDetectionBaselineStrategy(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"Id": id,
	}

	var response map[string]interface{}
	action := "DescribeStrategyDetail"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"-100"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("ThreatDetection:BaselineStrategy", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Strategy", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Strategy", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) DescribeStrategy(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"StrategyIds": id,
	}

	var response map[string]interface{}
	action := "DescribeStrategy"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ConsoleError"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("ThreatDetection:BaselineStrategy", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Strategies[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Strategies[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionBaselineStrategyStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionBaselineStrategy(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}

func (s *SasService) DescribeThreatDetectionAntiBruteForceRule(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"Id": id,
	}

	var response map[string]interface{}
	action := "DescribeAntiBruteForceRules"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	count, err := jsonpath.Get("$.PageInfo.Count", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PageInfo.Count]", response)
	}
	if formatInt(count) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ThreatDetection.AntiBruteForceRule", id)), NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.Rules[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Rules[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionAntiBruteForceRuleStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionAntiBruteForceRule(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}

func (s *SasService) DescribeThreatDetectionHoneyPot(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"HoneypotIds.1": id,
	}

	var response map[string]interface{}
	action := "ListHoneypot"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	count, err := jsonpath.Get("$.PageInfo.Count", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PageInfo.Count", response)
	}
	if formatInt(count) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ThreatDetection", id)), NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.List[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.List[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionHoneyPotStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionHoneyPot(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		state, err := jsonpath.Get("$.State[0]", object)
		if err != nil {
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(state) == failState {
				return object, fmt.Sprint(state), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(state)))
			}
		}
		return object, fmt.Sprint(state), nil
	}
}

func (s *SasService) DescribeThreatDetectionHoneypotProbe(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"ProbeId": id,
	}

	var response map[string]interface{}
	action := "GetHoneypotProbe"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"HoneypotProbeNotReady"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	success, err := jsonpath.Get("$.Success", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Success", response)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		if success.(bool) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ThreatDetection", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionHoneypotProbeStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionHoneypotProbe(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
