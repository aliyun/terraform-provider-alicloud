package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"time"
)

type CdnService struct {
	client *connectivity.AliyunClient
}

func (c *CdnService) DescribeCdnDomain(id string) (*cdn.GetDomainDetailModel, error) {
	request := cdn.CreateDescribeCdnDomainDetailRequest()
	request.DomainName = id

	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeCdnDomainDetail(request)
	})

	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	domain, _ := raw.(*cdn.DescribeCdnDomainDetailResponse)
	addDebug(request.GetActionName(), raw)
	if domain.GetDomainDetailModel.DomainName != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("cdn_domain", id)), NotFoundMsg, ProviderERROR)
	}
	return &domain.GetDomainDetailModel, nil
}

func (c *CdnService) DescribeDomainConfig(id string) (*cdn.DomainConfig, error) {
	strs := strings.Split(id, ":")
	request := cdn.CreateDescribeCdnDomainConfigsRequest()
	request.DomainName = strs[0]

	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeCdnDomainConfigs(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	configs, _ := raw.(*cdn.DescribeCdnDomainConfigsResponse)
	addDebug(request.GetActionName(), raw)
	for _, value := range configs.DomainConfigs.DomainConfig {
		if value.FunctionName == strs[1] {
			return &value, nil
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("cdn_domain_config", id)), NotFoundMsg, ProviderERROR)
}

func (c *CdnService) WaitForCdnDomain(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		domain, err := c.DescribeCdnDomain(id)
		if err != nil {
			if !IsExceptedError(err, InvalidDomainNotFound) {
				return WrapError(err)
			}
			continue
		}
		if domain.DomainStatus == string(status) {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return WrapError(GetTimeErrorFromString(GetTimeoutMessage("Domain", string(status))))
		}
	}
	return nil
}
