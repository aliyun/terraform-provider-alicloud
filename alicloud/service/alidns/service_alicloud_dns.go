package alidns

import (
	"reflect"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DnsService struct {
	client *connectivity.AliyunClient
}

func NewDnsService(client *connectivity.AliyunClient) *DnsService {
	return &DnsService{client: client}
}

func (s *DnsService) DescribeDns(id string) (*alidnssdk.DescribeDomainInfoResponse, error) {
	response := &alidnssdk.DescribeDomainInfoResponse{}
	request := alidnssdk.CreateDescribeDomainInfoRequest()
	request.RegionId = s.client.RegionId
	request.DomainName = id

	raw, err := s.client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
		return dnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	response, _ = raw.(*alidnssdk.DescribeDomainInfoResponse)
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.DomainName != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Dns", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}

func (dns *DnsService) DescribeDnsGroup(id string) (alidnssdk.DomainGroup, error) {
	var group alidnssdk.DomainGroup
	request := alidnssdk.CreateDescribeDomainGroupsRequest()
	request.RegionId = dns.client.RegionId
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := dns.client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			return group, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidnssdk.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if domainGroup.GroupId == id {
				return domainGroup, nil
			}
		}
		if len(groups) < helper.PageSizeLarge {
			break
		}
		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			return group, errs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return group, errs.WrapErrorf(errs.NotFoundErr("DnsGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (dns *DnsService) DescribeDnsRecord(id string) (*alidnssdk.DescribeDomainRecordInfoResponse, error) {
	response := &alidnssdk.DescribeDomainRecordInfoResponse{}
	request := alidnssdk.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = id
	request.RegionId = dns.client.RegionId
	raw, err := dns.client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
		return dnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRR.NoExist"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*alidnssdk.DescribeDomainRecordInfoResponse)
	if response.RecordId != id {
		return response, errs.WrapErrorf(errs.NotFoundErr("DnsRecord", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}

func (s *DnsService) DescribeDnsInstance(id string) (object alidnssdk.DescribeDnsProductInstanceResponse, err error) {
	request := alidnssdk.CreateDescribeDnsProductInstanceRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
		return alidnsClient.DescribeDnsProductInstance(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("DnsInstance", id), errs.NotFoundMsg, errs.ProviderERROR)
			return
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidnssdk.DescribeDnsProductInstanceResponse)
	return *response, nil
}

func (s *DnsService) DescribeAlidnsDomainAttachment(id string) (object alidnssdk.DescribeInstanceDomainsResponse, err error) {
	request := alidnssdk.CreateDescribeInstanceDomainsRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
		return alidnsClient.DescribeInstanceDomains(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("DnsDomainAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
			return
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidnssdk.DescribeInstanceDomainsResponse)

	if len(response.InstanceDomains) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("DnsDomainAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
		return
	}
	return *response, nil
}

func (s *DnsService) WaitForAlidnsDomainAttachment(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlidnsDomainAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return errs.WrapError(err)
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
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, "", expected, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *DnsService) ListTagResources(id string) (object alidnssdk.ListTagResourcesResponse, err error) {
	request := alidnssdk.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId

	request.ResourceType = "DOMAIN"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
		return alidnsClient.ListTagResources(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidnssdk.ListTagResourcesResponse)
	return *response, nil
}

func (s *DnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]alidnssdk.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, alidnssdk.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := alidnssdk.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.UntagResources(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := alidnssdk.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.TagResources(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *DnsService) DescribeDnsDomain(id string) (object alidnssdk.DescribeDomainInfoResponse, err error) {
	request := alidnssdk.CreateDescribeDomainInfoRequest()
	request.RegionId = s.client.RegionId

	request.DomainName = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("DnsDomain", id), errs.NotFoundMsg, errs.ProviderERROR)
			return
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidnssdk.DescribeDomainInfoResponse)
	return *response, nil
}
