package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DnsService struct {
	client *connectivity.AliyunClient
}

func (s *DnsService) DescribeDns(id string) (*alidns.DescribeDomainInfoResponse, error) {
	request := alidns.CreateDescribeDomainInfoRequest()
	request.DomainName = id

	raw, err := s.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidDomainNameNoExist) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*alidns.DescribeDomainInfoResponse)
	addDebug(request.GetActionName(), raw)
	if response.DomainName != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Dns", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (dns *DnsService) DescribeDnsGroup(id string) (alidns.DomainGroup, error) {
	var group alidns.DomainGroup
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			return group, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if domainGroup.GroupId == id {
				return domainGroup, nil
			}
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return group, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return group, WrapErrorf(Error(GetNotFoundMessage("DnsGroup", id)), NotFoundMsg, ProviderERROR)
}

func (dns *DnsService) DescribeDnsRecord(id string) (*alidns.DescribeDomainRecordInfoResponse, error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = id

	raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{DomainRecordNotBelongToUser}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidns.DescribeDomainRecordInfoResponse)
	if response.RecordId != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DnsRecord", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}
