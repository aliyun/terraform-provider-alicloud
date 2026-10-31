package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewConfigService(client *connectivity.AliyunClient) *ConfigService {
	return &ConfigService{client: client}
}

func NewConfigServiceV2(client *connectivity.AliyunClient) *ConfigServiceV2 {
	return &ConfigServiceV2{client: client}
}

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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"ConfigRuleNotExists", "Invalid.ConfigRuleId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Config:Rule", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ConfigRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ConfigRule", response)
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
		if errs.IsExpectedErrors(err, []string{"AccountNotExisted", "DeliveryChannelNotExists"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ConfigDeliveryChannel", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.DeliveryChannels", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DeliveryChannels", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Config", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DeliveryChannelId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("Config", id), errs.NotFoundWithResponse, response)
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
		if errs.IsExpectedErrors(err, []string{"AccountNotExisted"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ConfigConfigurationRecorder", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.ConfigurationRecorder", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ConfigurationRecorder", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) ConfigConfigurationRecorderStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigConfigurationRecorder(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["ConfigurationRecorderStatus"].(string) == failState {
				return object, object["ConfigurationRecorderStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["ConfigurationRecorderStatus"].(string)))
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
		return "", errs.WrapError(err)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Config:Aggregator", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Aggregator", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Aggregator", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *ConfigService) ConfigAggregatorStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigAggregator(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["AggregatorStatus"]) == failState {
				return object, fmt.Sprint(object["AggregatorStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["AggregatorStatus"])))
			}
		}
		return object, fmt.Sprint(object["AggregatorStatus"]), nil
	}
}

func (s *ConfigService) DescribeConfigAggregateConfigRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAggregateConfigRule"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"AggregatorId": parts[0],
		"ConfigRuleId": parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"ConfigRuleNotExists", "Invalid.AggregatorId.Value", "Invalid.ConfigRuleId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Config:AggregateConfigRule", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ConfigRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ConfigRule", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigAggregateCompliancePack(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetAggregateCompliancePack"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"AggregatorId":     parts[0],
		"CompliancePackId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value", "Invalid.CompliancePackId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Config:AggregateCompliancePack", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.CompliancePack", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CompliancePack", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ConfigService) ConfigAggregateCompliancePackStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigAggregateCompliancePack(id)
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

func (s *ConfigService) DescribeConfigCompliancePack(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetCompliancePack"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"CompliancePackId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"Invalid.CompliancePackId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Config:CompliancePack", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.CompliancePack", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CompliancePack", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) ConfigCompliancePackStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigCompliancePack(id)
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

func (s *ConfigService) ConfigAggregateConfigRuleStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigAggregateConfigRule(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ConfigRuleState"]) == failState {
				return object, fmt.Sprint(object["ConfigRuleState"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["ConfigRuleState"])))
			}
		}
		return object, fmt.Sprint(object["ConfigRuleState"]), nil
	}
}

func (s *ConfigService) ConfigRuleStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigRule(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ConfigRuleState"]) == failState {
				return object, fmt.Sprint(object["ConfigRuleState"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["ConfigRuleState"])))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2019-01-08", action, request, nil)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2019-01-08", action, request, nil)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"DeliveryChannelNotExists"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ConfigDelivery", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DeliveryChannel", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DeliveryChannel", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigAggregateDelivery(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAggregateConfigDeliveryChannel"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AggregatorId":      parts[0],
		"DeliveryChannelId": parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ConfigDeliveryChannel", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DeliveryChannel", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DeliveryChannel", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

type ConfigServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeConfigRule <<< Encapsulated get interface for Config Rule.

func (s *ConfigServiceV2) DescribeConfigRule(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetConfigRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	query["ConfigRuleId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"ConfigRuleNotExists", "Invalid.ConfigRuleId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Rule", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ConfigRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ConfigRule", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ConfigServiceV2) ConfigRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigRule(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
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

// DescribeConfigRule >>> Encapsulated.

// DescribeConfigRemediation <<< Encapsulated get interface for Config Remediation.

func (s *ConfigServiceV2) DescribeConfigRemediation(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeRemediation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	query["RemediationId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Config", "2020-09-07", action, query, request)

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
		if errs.IsExpectedErrors(err, []string{"RemediationConfigNotExist", "RemediationNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Remediation", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Remediation", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Remediation", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ConfigServiceV2) ConfigRemediationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigRemediation(id)
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

// DescribeConfigRemediation >>> Encapsulated.

// DescribeConfigDelivery <<< Encapsulated get interface for Config Delivery.

func (s *ConfigServiceV2) DescribeConfigDelivery(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetConfigDeliveryChannel"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DeliveryChannelId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DeliveryChannelNotExists"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Delivery", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DeliveryChannel", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DeliveryChannel", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ConfigServiceV2) ConfigDeliveryStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigDelivery(id)
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

// DescribeConfigDelivery >>> Encapsulated.

// DescribeConfigAggregateDelivery <<< Encapsulated get interface for Config AggregateDelivery.

func (s *ConfigServiceV2) DescribeConfigAggregateDelivery(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetAggregateConfigDeliveryChannel"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AggregatorId"] = parts[0]
	query["DeliveryChannelId"] = parts[1]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DeliveryChannelNotExists", "Invalid.AggregatorId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AggregateDelivery", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DeliveryChannel", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DeliveryChannel", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ConfigServiceV2) ConfigAggregateDeliveryStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeConfigAggregateDelivery(id)
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

// DescribeConfigAggregateDelivery >>> Encapsulated.

// DescribeConfigAggregator <<< Encapsulated get interface for Config Aggregator.

func (s *ConfigServiceV2) DescribeConfigAggregator(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AggregatorId"] = id

	action := "GetAggregator"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Aggregator", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Aggregator", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Aggregator", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ConfigServiceV2) ConfigAggregatorStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ConfigAggregatorStateRefreshFuncWithApi(id, field, failStates, s.DescribeConfigAggregator)
}

func (s *ConfigServiceV2) ConfigAggregatorStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		object["AggregatorStatus"] = convertConfigAggregatorAggregatorAggregatorStatusResponse(object["AggregatorStatus"])
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

func (s *ConfigServiceV2) DescribeAsyncConfigAggregatorStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncGetAggregator(d, res)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
		}
		object["AggregatorStatus"] = convertConfigAggregatorAggregatorAggregatorStatusResponse(object["AggregatorStatus"])
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeConfigAggregator >>> Encapsulated.

// DescribeAsyncGetAggregator <<< Encapsulated for Config.
func (s *ConfigServiceV2) DescribeAsyncGetAggregator(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AggregatorId"], err = jsonpath.Get("$.AggregatorId", res)

	action := "GetAggregator"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)

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
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncGetAggregator >>> Encapsulated.
// DescribeConfigReportTemplate <<< Encapsulated get interface for Config ReportTemplate.

func (s *ConfigServiceV2) DescribeConfigReportTemplate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ReportTemplateId"] = id

	action := "GetReportTemplate"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"ReportTemplateNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ReportTemplate", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ReportTemplate", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ReportTemplate", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ConfigServiceV2) ConfigReportTemplateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ConfigReportTemplateStateRefreshFuncWithApi(id, field, failStates, s.DescribeConfigReportTemplate)
}

func (s *ConfigServiceV2) ConfigReportTemplateStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeConfigReportTemplate >>> Encapsulated.
// DescribeConfigAggregateRemediation <<< Encapsulated get interface for Config AggregateRemediation.

func (s *ConfigServiceV2) DescribeConfigAggregateRemediation(id string) (object map[string]interface{}, err error) {
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
	request["AggregatorId"] = parts[0]
	request["RemediationId"] = parts[1]

	action := "DescribeAggregateRemediation"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"RemediationNotExist", "Invalid.AggregatorId.Value"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AggregateRemediation", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Remediation", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Remediation", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ConfigServiceV2) ConfigAggregateRemediationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ConfigAggregateRemediationStateRefreshFuncWithApi(id, field, failStates, s.DescribeConfigAggregateRemediation)
}

func (s *ConfigServiceV2) ConfigAggregateRemediationStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeConfigAggregateRemediation >>> Encapsulated.
