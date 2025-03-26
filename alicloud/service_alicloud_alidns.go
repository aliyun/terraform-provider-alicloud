package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type AlidnsService struct {
	client *connectivity.AliyunClient
}

func (s *AlidnsService) DescribeAlidnsDomainGroup(id string) (object alidns.DomainGroup, err error) {
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.RegionId = s.client.RegionId

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)

		if len(response.DomainGroups.DomainGroup) < 1 {
			err = WrapErrorf(NotFoundErr("AlidnsDomainGroup", id), NotFoundMsg, ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.DomainGroups.DomainGroup {
			if object.GroupId == id {
				return object, nil
			}
		}
		if len(response.DomainGroups.DomainGroup) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(NotFoundErr("AlidnsDomainGroup", id), NotFoundMsg, ProviderERROR)
	return
}

func (s *AlidnsService) DescribeAlidnsRecord(id string) (object alidns.DescribeDomainRecordInfoResponse, err error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RegionId = s.client.RegionId

	request.RecordId = id

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRR.NoExist"}) {
			err = WrapErrorf(NotFoundErr("AlidnsRecord", id), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDomainRecordInfoResponse)
	return *response, nil
}

func (s *AlidnsService) ListTagResources(id string) (object alidns.ListTagResourcesResponse, err error) {
	request := alidns.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId

	request.ResourceType = "DOMAIN"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.ListTagResourcesResponse)
	return *response, nil
}

func (s *AlidnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]alidns.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, alidns.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := alidns.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := alidns.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *AlidnsService) DescribeAlidnsDomain(id string) (object alidns.DescribeDomainInfoResponse, err error) {
	request := alidns.CreateDescribeDomainInfoRequest()
	request.RegionId = s.client.RegionId

	request.DomainName = id

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = WrapErrorf(NotFoundErr("AlidnsDomain", id), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDomainInfoResponse)
	return *response, nil
}

func (s *AlidnsService) DescribeAlidnsInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDnsProductInstance"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(11*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			return object, WrapErrorf(NotFoundErr("Alidns:Instance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsCustomLine(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeCustomLine"
	request := map[string]interface{}{
		"LineId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsCustomLine.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns::CustomLine", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeCustomLine(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeCustomLine"
	request := map[string]interface{}{
		"LineId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsCustomLine.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns::CustomLine", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsGtmInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDnsGtmInstance"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsGtmInstance.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns:GtmInstance", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsAddressPool(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDnsGtmInstanceAddressPool"
	request := map[string]interface{}{
		"AddrPoolId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsGtmAddrPool.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns::AddressPool", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsAccessStrategy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDnsGtmAccessStrategy"
	request := map[string]interface{}{
		"StrategyId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsGtmAccessStrategy.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns:DnsGtmAccessStrategy", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsMonitorConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDnsGtmMonitorConfig"
	request := map[string]interface{}{
		"MonitorConfigId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
