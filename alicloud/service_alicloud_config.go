package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ConfigService struct {
	client *connectivity.AliyunClient
}

func (s *ConfigService) DescribeConfigRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetConfigRule"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ConfigRuleId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if IsExpectedErrors(err, []string{"ConfigRuleNotExists", "Invalid.ConfigRuleId.Value"}) {
			return object, WrapErrorf(NotFoundErr("Config:Rule", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ConfigRule", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ConfigRule", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigDeliveryChannel(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDeliveryChannels"
	request := map[string]interface{}{
		"RegionId":           s.client.RegionId,
		"DeliveryChannelIds": id,
	}
	response, err = client.RpcGet("Config", "2019-01-08", action, request, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted", "DeliveryChannelNotExists"}) {
			err = WrapErrorf(NotFoundErr("ConfigDeliveryChannel", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.DeliveryChannels", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DeliveryChannels", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Config", id), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DeliveryChannelId"].(string) != id {
			return object, WrapErrorf(NotFoundErr("Config", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigConfigurationRecorder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeConfigurationRecorder"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	response, err = client.RpcGet("Config", "2019-01-08", action, request, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted"}) {
			err = WrapErrorf(NotFoundErr("ConfigConfigurationRecorder", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ConfigurationRecorder", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ConfigurationRecorder", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) ConfigConfigurationRecorderStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigConfigurationRecorder(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ConfigurationRecorderStatus"].(string) == failState {
				return object, object["ConfigurationRecorderStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ConfigurationRecorderStatus"].(string)))
			}
		}
		return object, object["ConfigurationRecorderStatus"].(string), nil
	}
}

func (s *ConfigService) convertAggregatorAccountsToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = map[string]interface{}{
			"AccountId":   item["account_id"],
			"AccountName": item["account_name"],
			"AccountType": item["account_type"],
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *ConfigService) DescribeConfigAggregator(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"AggregatorId": id,
	}

	var response map[string]interface{}
	action := "GetAggregator"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value"}) {
			return object, WrapErrorf(NotFoundErr("Config:Aggregator", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Aggregator", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Aggregator", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *ConfigService) ConfigAggregatorStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigAggregator(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["AggregatorStatus"]) == failState {
				return object, fmt.Sprint(object["AggregatorStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["AggregatorStatus"])))
			}
		}
		return object, fmt.Sprint(object["AggregatorStatus"]), nil
	}
}

func (s *ConfigService) DescribeConfigAggregateConfigRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAggregateConfigRule"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"AggregatorId": parts[0],
		"ConfigRuleId": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if IsExpectedErrors(err, []string{"ConfigRuleNotExists", "Invalid.AggregatorId.Value", "Invalid.ConfigRuleId.Value"}) {
			return object, WrapErrorf(NotFoundErr("Config:AggregateConfigRule", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ConfigRule", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ConfigRule", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigAggregateCompliancePack(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetAggregateCompliancePack"

	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"AggregatorId":     parts[0],
		"CompliancePackId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value", "Invalid.CompliancePackId.Value"}) {
			return object, WrapErrorf(NotFoundErr("Config:AggregateCompliancePack", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.CompliancePack", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CompliancePack", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ConfigService) ConfigAggregateCompliancePackStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigAggregateCompliancePack(id)
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

func (s *ConfigService) DescribeConfigCompliancePack(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetCompliancePack"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"CompliancePackId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if IsExpectedErrors(err, []string{"Invalid.CompliancePackId.Value"}) {
			return object, WrapErrorf(NotFoundErr("Config:CompliancePack", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.CompliancePack", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CompliancePack", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) ConfigCompliancePackStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigCompliancePack(id)
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

func (s *ConfigService) ConfigAggregateConfigRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigAggregateConfigRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ConfigRuleState"]) == failState {
				return object, fmt.Sprint(object["ConfigRuleState"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["ConfigRuleState"])))
			}
		}
		return object, fmt.Sprint(object["ConfigRuleState"]), nil
	}
}

func (s *ConfigService) ConfigRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ConfigRuleState"]) == failState {
				return object, fmt.Sprint(object["ConfigRuleState"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["ConfigRuleState"])))
			}
		}
		return object, fmt.Sprint(object["ConfigRuleState"]), nil
	}
}

func (s *ConfigService) StopConfigRule(id string) (err error) {
	var response map[string]interface{}
	client := s.client
	action := "StopConfigRules"
	request := map[string]interface{}{
		"ConfigRuleIds": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2019-01-08", action, request, nil)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func (s *ConfigService) ActiveConfigRule(id string) (err error) {
	var response map[string]interface{}
	client := s.client
	action := "ActiveConfigRules"
	request := map[string]interface{}{
		"ConfigRuleIds": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2019-01-08", action, request, nil)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *ConfigService) ActiveAggregateConfigRules(id, aggregatorId string) (err error) {
	var response map[string]interface{}
	client := s.client
	action := "ActiveAggregateConfigRules"
	request := map[string]interface{}{
		"ConfigRuleIds": id,
		"AggregatorId":  aggregatorId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *ConfigService) DeactiveAggregateConfigRules(id, aggregatorId string) (err error) {
	var response map[string]interface{}
	client := s.client
	action := "DeactiveAggregateConfigRules"
	request := map[string]interface{}{
		"ConfigRuleIds": id,
		"AggregatorId":  aggregatorId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *ConfigService) DescribeConfigDelivery(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetConfigDeliveryChannel"
	request := map[string]interface{}{
		"DeliveryChannelId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"DeliveryChannelNotExists"}) {
			err = WrapErrorf(NotFoundErr("ConfigDelivery", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DeliveryChannel", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DeliveryChannel", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigAggregateDelivery(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAggregateConfigDeliveryChannel"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AggregatorId":      parts[0],
		"DeliveryChannelId": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value"}) {
			err = WrapErrorf(NotFoundErr("ConfigDeliveryChannel", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DeliveryChannel", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DeliveryChannel", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
