package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type OnsService struct {
	client *connectivity.AliyunClient
}

func (s *OnsService) DescribeOnsInstance(id string) (object ons.InstanceBaseInfo, err error) {
	request := ons.CreateOnsInstanceBaseInfoRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceBaseInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OnsInstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ons.OnsInstanceBaseInfoResponse)
	return response.InstanceBaseInfo, nil
}

func (s *OnsService) DescribeOnsTopic(id string) (object ons.PublishInfoDo, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := ons.CreateOnsTopicListRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = parts[0]
	request.Topic = parts[1]

	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsTopicList(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OnsTopic", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ons.OnsTopicListResponse)

	if len(response.Data.PublishInfoDo) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("OnsTopic", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Data.PublishInfoDo[0], nil
}

func (s *OnsService) DescribeOnsGroup(id string) (object ons.SubscribeInfoDo, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := ons.CreateOnsGroupListRequest()
	request.RegionId = s.client.RegionId
	request.GroupId = parts[1]
	request.InstanceId = parts[0]

	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsGroupList(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OnsGroup", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ons.OnsGroupListResponse)

	if len(response.Data.SubscribeInfoDo) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("OnsGroup", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Data.SubscribeInfoDo[0], nil
}

func (s *OnsService) WaitForOnsTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]
	for {
		response, err := s.DescribeOnsTopic(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if response.InstanceId+":"+response.Topic == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, instanceId+":"+topic, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *OnsService) WaitForOnsGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		response, err := s.DescribeOnsGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if response.InstanceId+":"+response.GroupId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, response.InstanceId+":"+response.GroupId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *OnsService) ListTagResources(id string) (object ons.ListTagResourcesResponse, err error) {
	request := ons.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id
	request.ResourceType = "INSTANCE"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.ListTagResources(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OnsInstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ons.ListTagResourcesResponse)
	return *response, nil
}

func (s *OnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]ons.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, ons.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := ons.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := ons.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *OnsService) OnsTopicStatus(id string) (object ons.Data, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := ons.CreateOnsTopicStatusRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = parts[0]
	request.Topic = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(11*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicStatus(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("OnsTopic", id)), NotFoundMsg, ProviderERROR)
				return resource.NonRetryableError(err)
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ons.OnsTopicStatusResponse)
		object = response.Data
		return nil
	})
	return object, WrapError(err)
}

func (s *OnsService) SetResourceTagsForTopic(d *schema.ResourceData, resourceType string) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	oldItems, newItems := d.GetChange("tags")
	added := make([]ons.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, ons.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := ons.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.InstanceId = parts[0]
		request.ResourceId = &[]string{parts[1]}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := ons.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.InstanceId = parts[0]
		request.ResourceId = &[]string{parts[1]}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *OnsService) SetResourceTagsForGroup(d *schema.ResourceData, resourceType string) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	oldItems, newItems := d.GetChange("tags")
	added := make([]ons.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, ons.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := ons.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.InstanceId = parts[0]
		request.ResourceId = &[]string{parts[1]}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := ons.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.InstanceId = parts[0]
		request.ResourceId = &[]string{parts[1]}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}
