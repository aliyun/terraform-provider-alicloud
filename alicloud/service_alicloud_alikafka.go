package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type AlikafkaService struct {
	client *connectivity.AliyunClient
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaInstance(instanceId string) (alikafkaInstance *alikafka.InstanceVO, err error) {

	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = alikafkaService.client.RegionId

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})

	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, instanceId, instanceListReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)
	addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)

	for _, v := range instanceListResp.InstanceList.InstanceVO {
		if v.InstanceId == instanceId {
			alikafkaInstance = &v
			return
		}
	}
	return alikafkaInstance, WrapErrorf(Error(GetNotFoundMessage("AlikafkaInstance", instanceId)), NotFoundMsg, ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaTopic(id string) (alikafkaTopic *alikafka.TopicVO, err error) {

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateGetTopicListRequest()
	request.InstanceId = instanceId
	request.RegionId = alikafkaService.client.RegionId

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetTopicList(request)
	})

	if err != nil {
		return alikafkaTopic, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	topicListResp, _ := raw.(*alikafka.GetTopicListResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range topicListResp.TopicList.TopicVO {
		if v.Topic == topic {
			alikafkaTopic = &v
			return
		}
	}
	return alikafkaTopic, WrapErrorf(Error(GetNotFoundMessage("AlikafkaTopic", id)), NotFoundMsg, ProviderERROR)
}

func (s *AlikafkaService) WaitForAlikafkaTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaTopic(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.InstanceId+":"+object.Topic == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId+":"+object.Topic, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}
