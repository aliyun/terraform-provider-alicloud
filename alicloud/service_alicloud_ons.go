package alicloud

import (
	"time"

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

func (s *OnsService) DescribeOnsTopic(id string) (*ons.PublishInfoDo, error) {
	onsTopic := &ons.PublishInfoDo{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return onsTopic, WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := ons.CreateOnsTopicListRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = instanceId

	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsTopicList(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			return onsTopic, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return onsTopic, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	topicListResp, _ := raw.(*ons.OnsTopicListResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range topicListResp.Data.PublishInfoDo {
		if v.Topic == topic {
			return &v, nil
		}
	}
	return onsTopic, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
}

func (s *OnsService) DescribeOnsGroup(id string) (*ons.SubscribeInfoDo, error) {
	onsGroup := &ons.SubscribeInfoDo{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return onsGroup, WrapError(err)
	}
	instanceId := parts[0]
	groupId := parts[1]

	request := ons.CreateOnsGroupListRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = instanceId

	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsGroupList(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			return onsGroup, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return onsGroup, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	groupListResp, _ := raw.(*ons.OnsGroupListResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range groupListResp.Data.SubscribeInfoDo {
		if v.GroupId == groupId {
			return &v, nil
		}
	}

	return onsGroup, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
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
