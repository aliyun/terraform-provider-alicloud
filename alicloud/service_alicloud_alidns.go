package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			err = WrapErrorf(Error(GetNotFoundMessage("AlidnsDomainGroup", id)), NotFoundMsg, ProviderERROR)
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
	err = WrapErrorf(Error(GetNotFoundMessage("AlidnsDomainGroup", id)), NotFoundMsg, ProviderERROR)
	return
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
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("alidnsdomain", id)), NotFoundMsg, ProviderERROR)
			return
		}
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
	for key, _ := range oldItems.(map[string]interface{}) {
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
			err = WrapErrorf(Error(GetNotFoundMessage("AlidnsDomain", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDomainInfoResponse)
	return *response, nil
}
