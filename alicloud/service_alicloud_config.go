package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
