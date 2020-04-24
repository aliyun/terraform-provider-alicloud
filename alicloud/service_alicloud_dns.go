package alicloud

import (
	"reflect"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DnsService struct {
	client *connectivity.AliyunClient
}

func (s *DnsService) DescribeDns(id string) (*alidns.DescribeDomainInfoResponse, error) {
	response := &alidns.DescribeDomainInfoResponse{}
	request := alidns.CreateDescribeDomainInfoRequest()
	request.RegionId = s.client.RegionId
	request.DomainName = id

	raw, err := s.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ = raw.(*alidns.DescribeDomainInfoResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.DomainName != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Dns", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (dns *DnsService) DescribeDnsGroup(id string) (alidns.DomainGroup, error) {
	var group alidns.DomainGroup
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.RegionId = dns.client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			return group, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
	response := &alidns.DescribeDomainRecordInfoResponse{}
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = id
	request.RegionId = dns.client.RegionId
	raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRR.NoExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*alidns.DescribeDomainRecordInfoResponse)
	if response.RecordId != id {
		return response, WrapErrorf(Error(GetNotFoundMessage("DnsRecord", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (s *DnsService) DescribeDnsInstance(id string) (object alidns.DescribeDnsProductInstanceResponse, err error) {
	request := alidns.CreateDescribeDnsProductInstanceRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDnsProductInstance(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsInstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDnsProductInstanceResponse)
	return *response, nil
}

func (s *DnsService) DescribeDnsDomainAttachment(id string) (object alidns.DescribeInstanceDomainsResponse, err error) {
	request := alidns.CreateDescribeInstanceDomainsRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeInstanceDomains(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeInstanceDomainsResponse)

	if len(response.InstanceDomains) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return *response, nil
}

func (s *DnsService) WaitForAlidnsDomainAttachment(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDnsDomainAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		domainNames := make(map[string]interface{}, 0)
		for _, v := range object.InstanceDomains {
			domainNames[v.DomainName] = v.DomainName
		}

		exceptDomainNames := make(map[string]interface{}, 0)
		for _, v := range expected {
			for _, vv := range v.([]interface{}) {
				exceptDomainNames[vv.(string)] = vv.(string)
			}
		}

		if reflect.DeepEqual(domainNames, exceptDomainNames) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", expected, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}
