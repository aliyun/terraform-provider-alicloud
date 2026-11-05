package threatdetection

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SasService struct {
	client *connectivity.AliyunClient
}

func NewSasService(client *connectivity.AliyunClient) *SasService {
	return &SasService{client: client}
}

func (s *SasService) DescribeAllGroups(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAllGroups"
	request := map[string]interface{}{}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Groups", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Groups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SecurityCenter", id), errs.NotFoundWithResponse, response)
	}
	for _, vv := range v.([]interface{}) {
		if fmt.Sprint(vv.(map[string]interface{})["GroupId"]) == id {
			idExist = true
			return vv.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("SecurityCenter", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *SasService) DescribeSecurityCenterServiceLinkedRole(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeServiceLinkedRoleStatus"
	request := map[string]interface{}{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Sas", "2018-12-03", action, request, nil)
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
	v, err := jsonpath.Get("$.RoleStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.RoleStatus", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SasService) SecurityCenterServiceLinkedRoleStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSecurityCenterServiceLinkedRole(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ConfigList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ConfigList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("WebLockConfig", id), errs.NotFoundWithResponse, response)
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
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if helper.NeedRetry(err) {
			return retry.RetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"-100"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("ThreatDetection:BaselineStrategy", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	// This is an API behavior change. An HTTP status code of 200 will be returned when the resource deleted.
	// Customize this error message to handle this API change.
	if response == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("ThreatDetection:BaselineStrategy", id), errs.NotFoundMsg)
	}
	v, err := jsonpath.Get("$.Strategy", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Strategy", response)
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
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if helper.NeedRetry(err) {
			return retry.RetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ConsoleError"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("ThreatDetection:BaselineStrategy", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Strategies[0]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Strategies[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionBaselineStrategyStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionBaselineStrategy(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object[""])))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	count, err := jsonpath.Get("$.PageInfo.Count", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PageInfo.Count]", response)
	}
	totalCount, err := convert.FormatInt(count)
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ThreatDetection.AntiBruteForceRule", id), errs.NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.Rules[0]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Rules[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionAntiBruteForceRuleStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionAntiBruteForceRule(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object[""])))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	count, err := jsonpath.Get("$.PageInfo.Count", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PageInfo.Count", response)
	}
	totalCount, err := convert.FormatInt(count)
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ThreatDetection", id), errs.NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.List[0]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.List[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionHoneyPotStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionHoneyPot(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		state, err := jsonpath.Get("$.State[0]", object)
		if err != nil {
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(state) == failState {
				return object, fmt.Sprint(state), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(state)))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"HoneypotProbeNotReady"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	success, err := jsonpath.Get("$.Success", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Success", response)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		if success.(bool) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ThreatDetection", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *SasService) ThreatDetectionHoneypotProbeStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionHoneypotProbe(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
