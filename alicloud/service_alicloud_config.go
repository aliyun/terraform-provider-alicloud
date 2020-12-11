package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ConfigService struct {
	client *connectivity.AliyunClient
}

func (s *ConfigService) DescribeConfigRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewConfigClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeConfigRule"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ConfigRuleId": id,
	}
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-01-08"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted", "ConfigRuleNotExists"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ConfigRule", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ConfigRule", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ConfigRule", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigDeliveryChannel(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewConfigClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDeliveryChannels"
	request := map[string]interface{}{
		"RegionId":           s.client.RegionId,
		"DeliveryChannelIds": id,
	}
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-01-08"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted", "DeliveryChannelNotExists"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ConfigDeliveryChannel", id)), NotFoundMsg, ProviderERROR)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Config", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DeliveryChannelId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("Config", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ConfigService) DescribeConfigConfigurationRecorder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewConfigClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeConfigurationRecorder"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-01-08"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ConfigConfigurationRecorder", id)), NotFoundMsg, ProviderERROR)
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
