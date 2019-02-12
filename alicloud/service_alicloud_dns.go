package alicloud

import (
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
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	domain, _ := raw.(*alidns.DescribeDomainInfoResponse)
	addDebug(request.GetActionName(), raw)
	if domain.DomainName != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Dns", id)), NotFoundMsg, ProviderERROR)
	}
	return domain, nil
}

func (dns *DnsService) DescribeDnsGroup(id string) (group alidns.DomainGroup, err error) {
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.KeyWord = id

	raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainGroups(request)
	})
	if err != nil {
		return group, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	groups, _ := raw.(*alidns.DescribeDomainGroupsResponse)
	domainGroup := groups.DomainGroups.DomainGroup
	for _, v := range domainGroup {
		if v.GroupName == id {
			group = v
			return group, nil
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
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	recordInfo, _ := raw.(*alidns.DescribeDomainRecordInfoResponse)
	if recordInfo.RecordId != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DnsRecord", id)), NotFoundMsg, ProviderERROR)
	}
	return recordInfo, nil
}
