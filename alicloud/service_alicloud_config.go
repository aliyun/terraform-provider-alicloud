package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ConfigService struct {
	client *connectivity.AliyunClient
}

func (s *ConfigService) DescribeConfigRule(id string) (object config.ConfigRule, err error) {
	request := config.CreateDescribeConfigRuleRequest()
	request.RegionId = s.client.RegionId

	request.ConfigRuleId = id

	raw, err := s.client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.DescribeConfigRule(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted", "ConfigRuleNotExists"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ConfigRule", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*config.DescribeConfigRuleResponse)
	return response.ConfigRule, nil
}
func (s *ConfigService) DescribeConfigDeliveryChannel(id string) (object config.DeliveryChannel, err error) {
	request := config.CreateDescribeDeliveryChannelsRequest()
	request.RegionId = s.client.RegionId

	request.DeliveryChannelIds = id

	raw, err := s.client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.DescribeDeliveryChannels(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted", "DeliveryChannelNotExists"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ConfigDeliveryChannel", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*config.DescribeDeliveryChannelsResponse)

	if len(response.DeliveryChannels) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("ConfigDeliveryChannel", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.DeliveryChannels[0], nil
}

func (s *ConfigService) DescribeConfigConfigurationRecorder(id string) (object config.ConfigurationRecorder, err error) {
	request := config.CreateDescribeConfigurationRecorderRequest()
	request.RegionId = s.client.RegionId

	raw, err := s.client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.DescribeConfigurationRecorder(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ConfigConfigurationRecorder", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*config.DescribeConfigurationRecorderResponse)
	return response.ConfigurationRecorder, nil
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
			if object.ConfigurationRecorderStatus == failState {
				return object, object.ConfigurationRecorderStatus, WrapError(Error(FailedToReachTargetStatus, object.ConfigurationRecorderStatus))
			}
		}
		return object, object.ConfigurationRecorderStatus, nil
	}
}
