package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CdnService struct {
	client *connectivity.AliyunClient
}

func (c *CdnService) DescribeCdnDomainNew(id string) (*cdn.GetDomainDetailModel, error) {
	request := cdn.CreateDescribeCdnDomainDetailRequest()
	request.RegionId = c.client.RegionId
	request.DomainName = id

	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeCdnDomainDetail(request)
	})

	if err != nil {
		if IsExceptedError(err, InvalidDomainNotFound) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	domain, _ := raw.(*cdn.DescribeCdnDomainDetailResponse)
	if domain.GetDomainDetailModel.DomainName != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("cdn_domain", id)), NotFoundMsg, ProviderERROR)
	}
	return &domain.GetDomainDetailModel, nil
}

func (c *CdnService) DescribeCdnDomainConfig(id string) (*cdn.DomainConfig, error) {
	parts, err := ParseResourceId(id, 2)
	if err != err {
		return nil, WrapError(err)
	}
	request := cdn.CreateDescribeCdnDomainConfigsRequest()
	request.RegionId = c.client.RegionId
	request.DomainName = parts[0]

	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeCdnDomainConfigs(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidDomainNotFound) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cdn.DescribeCdnDomainConfigsResponse)
	for _, value := range response.DomainConfigs.DomainConfig {
		if value.FunctionName == parts[1] {
			return &value, nil
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("cdn_domain_config", id)), NotFoundMsg, ProviderERROR)
}

func (c *CdnService) WaitForCdnDomain(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	time.Sleep(DefaultIntervalShort * time.Second)

	for {
		domain, err := c.DescribeCdnDomainNew(id)
		if err != nil {
			if NotFoundError(err) && status == Deleted {
				break
			}
			return WrapError(err)
		}
		if domain.DomainStatus == string(status) {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, domain.DomainStatus, status, ProviderERROR)
		}
	}
	return nil
}

func (c *CdnService) DescribeDomainCertificateInfo(id string) (certInfo cdn.CertInfo, err error) {
	request := cdn.CreateDescribeDomainCertificateInfoRequest()
	request.RegionId = c.client.RegionId
	request.DomainName = id
	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeDomainCertificateInfo(request)
	})
	if err != nil {
		return certInfo, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cdn.DescribeDomainCertificateInfoResponse)
	if len(response.CertInfos.CertInfo) <= 0 {
		return certInfo, WrapErrorf(Error(GetNotFoundMessage("DomainCertificateInfo", id)), NotFoundMsg, ProviderERROR)
	}
	certInfo = response.CertInfos.CertInfo[0]
	return
}

func (c *CdnService) WaitForServerCertificateNew(id string, serverCertificate string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		certInfo, err := c.DescribeDomainCertificateInfo(id)
		if err != nil {
			return WrapError(err)
		}
		if strings.TrimSpace(certInfo.ServerCertificate) == strings.TrimSpace(serverCertificate) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strings.TrimSpace(certInfo.ServerCertificate), strings.TrimSpace(serverCertificate), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (c *CdnService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []cdn.TagItem, err error) {
	request := cdn.CreateDescribeTagResourcesRequest()
	request.RegionId = c.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cdn.DescribeTagResourcesResponse)
	if len(response.TagResources) < 1 {
		return
	}
	for _, t := range response.TagResources {
		tags = append(tags, t.Tag...)
	}
	return
}
