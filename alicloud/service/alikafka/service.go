package alikafka

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	alikafkasdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AlikafkaService struct {
	client *connectivity.AliyunClient
}

// NewAlikafkaService constructs an AlikafkaService bound to the shared client,
// keeping the product package self-contained (it does not reverse-import the
// root alicloud package).
func NewAlikafkaService(client *connectivity.AliyunClient) *AlikafkaService {
	return &AlikafkaService{client: client}
}

// NewAlikafkaServiceV2 constructs an AlikafkaServiceV2 (the generated V2
// service defined in service_v2.go) bound to the shared client. It lives in
// the handwritten service.go so the generated file stays as close to its
// generated form as the migration allows.
func NewAlikafkaServiceV2(client *connectivity.AliyunClient) *AlikafkaServiceV2 {
	return &AlikafkaServiceV2{client: client}
}

func (s *AlikafkaService) DescribeAlikafkaInstance(instanceId string) (*alikafkasdk.InstanceVO, error) {
	alikafkaInstance := &alikafkasdk.InstanceVO{}
	instanceListReq := alikafkasdk.CreateGetInstanceListRequest()
	instanceListReq.RegionId = s.client.RegionId

	wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
	var raw interface{}
	var err error
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(client *alikafkasdk.Client) (interface{}, error) {
			return client.GetInstanceList(instanceListReq)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)
		return nil
	})

	if err != nil {
		return alikafkaInstance, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, instanceListReq.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	instanceListResp, _ := raw.(*alikafkasdk.GetInstanceListResponse)
	logs.AddDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		// ServiceStatus equals 10 means the instance is released, do not return the instance.
		if v.InstanceId == instanceId && v.ServiceStatus != 10 {
			return &v, nil
		}
	}
	return alikafkaInstance, errs.WrapErrorf(errs.NotFoundErr("AlikafkaInstance", instanceId), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *AlikafkaService) DescribeAlikafkaInstanceByOrderId(orderId string, timeout int) (*alikafkasdk.InstanceVO, error) {
	alikafkaInstance := &alikafkasdk.InstanceVO{}
	instanceListReq := alikafkasdk.CreateGetInstanceListRequest()
	instanceListReq.RegionId = s.client.RegionId
	instanceListReq.OrderId = orderId

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {

		wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
		var raw interface{}
		var err error
		err = retry.Retry(10*time.Minute, func() *retry.RetryError {
			raw, err = s.client.WithAlikafkaClient(func(client *alikafkasdk.Client) (interface{}, error) {
				return client.GetInstanceList(instanceListReq)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)
			return nil
		})

		if err != nil {
			return alikafkaInstance, errs.WrapErrorf(err, errs.DefaultErrorMsg, orderId, instanceListReq.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		instanceListResp, _ := raw.(*alikafkasdk.GetInstanceListResponse)
		logs.AddDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)
		for _, v := range instanceListResp.InstanceList.InstanceVO {
			return &v, nil
		}
		if time.Now().After(deadline) {
			return alikafkaInstance, errs.WrapErrorf(errs.NotFoundErr("AlikafkaInstance", orderId), errs.NotFoundMsg, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) DescribeAlikafkaConsumerGroup(id string) (*alikafkasdk.ConsumerVO, error) {
	alikafkaConsumerGroup := &alikafkasdk.ConsumerVO{}

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return alikafkaConsumerGroup, errs.WrapError(err)
	}
	instanceId := parts[0]
	consumerId := parts[1]

	request := alikafkasdk.CreateGetConsumerListRequest()
	request.InstanceId = instanceId
	request.RegionId = s.client.RegionId

	wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
	var raw interface{}
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(client *alikafkasdk.Client) (interface{}, error) {
			return client.GetConsumerList(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return alikafkaConsumerGroup, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	consumerListResp, _ := raw.(*alikafkasdk.GetConsumerListResponse)
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range consumerListResp.ConsumerList.ConsumerVO {
		if v.ConsumerId == consumerId {
			return &v, nil
		}
	}
	return alikafkaConsumerGroup, errs.WrapErrorf(errs.NotFoundErr("AlikafkaConsumerGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *AlikafkaService) DescribeAlikafkaTopicStatus(id string) (*alikafkasdk.TopicStatus, error) {
	alikafkaTopicStatus := &alikafkasdk.TopicStatus{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return alikafkaTopicStatus, errs.WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafkasdk.CreateGetTopicStatusRequest()
	request.InstanceId = instanceId
	request.RegionId = s.client.RegionId
	request.Topic = topic

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.GetTopicStatus(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return alikafkaTopicStatus, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	topicStatusResp, _ := raw.(*alikafkasdk.GetTopicStatusResponse)

	if topicStatusResp.TopicStatus.OffsetTable.OffsetTableItem != nil {
		return &topicStatusResp.TopicStatus, nil
	}

	return alikafkaTopicStatus, errs.WrapErrorf(errs.NotFoundErr("AlikafkaTopicStatus "+errs.ResourceNotfound, id), errs.ResourceNotfound)
}

func (s *AlikafkaService) DescribeAlikafkaTopic(id string) (*alikafkasdk.TopicVO, error) {

	alikafkaTopic := &alikafkasdk.TopicVO{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return alikafkaTopic, errs.WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafkasdk.CreateGetTopicListRequest()
	request.InstanceId = instanceId
	request.RegionId = s.client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.GetTopicList(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return alikafkaTopic, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	topicListResp, _ := raw.(*alikafkasdk.GetTopicListResponse)

	for _, v := range topicListResp.TopicList.TopicVO {
		if v.Topic == topic {
			return &v, nil
		}
	}
	return alikafkaTopic, errs.WrapErrorf(errs.NotFoundErr("AlikafkaTopic", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *AlikafkaService) DescribeAlikafkaSaslUser(id string) (*alikafkasdk.SaslUserVO, error) {
	alikafkaSaslUser := &alikafkasdk.SaslUserVO{}

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return alikafkaSaslUser, errs.WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]

	request := alikafkasdk.CreateDescribeSaslUsersRequest()
	request.InstanceId = instanceId
	request.RegionId = s.client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return alikafkaSaslUser, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	userListResp, _ := raw.(*alikafkasdk.DescribeSaslUsersResponse)
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range userListResp.SaslUserList.SaslUserVO {
		if v.Username == username {
			return &v, nil
		}
	}
	return alikafkaSaslUser, errs.WrapErrorf(errs.NotFoundErr("AlikafkaSaslUser", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *AlikafkaService) DescribeAlikafkaSaslAcl(id string) (*alikafkasdk.KafkaAclVO, error) {
	alikafkaSaslAcl := &alikafkasdk.KafkaAclVO{}

	parts, err := helper.ParseResourceId(id, 6)
	if err != nil {
		return alikafkaSaslAcl, errs.WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	aclResourceType := parts[2]
	aclResourceName := parts[3]
	aclResourcePatternType := parts[4]
	aclOperationType := parts[5]

	request := alikafkasdk.CreateDescribeAclsRequest()
	request.InstanceId = instanceId
	request.RegionId = s.client.RegionId
	request.Username = username
	request.AclResourceType = aclResourceType
	request.AclResourceName = aclResourceName

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.DescribeAcls(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"BIZ_SUBSCRIPTION_NOT_FOUND", "BIZ_TOPIC_NOT_FOUND", "BIZ.INSTANCE.STATUS.ERROR"}) {
			return alikafkaSaslAcl, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return alikafkaSaslAcl, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	aclListResp, _ := raw.(*alikafkasdk.DescribeAclsResponse)
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range aclListResp.KafkaAclList.KafkaAclVO {
		if v.Username == username && v.AclResourceType == aclResourceType && v.AclResourceName == aclResourceName && v.AclResourcePatternType == aclResourcePatternType && v.AclOperationType == aclOperationType {
			return &v, nil
		}
	}
	return alikafkaSaslAcl, errs.WrapErrorf(errs.NotFoundErr("AlikafkaSaslAcl", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *AlikafkaService) WaitForAlikafkaInstanceUpdated(id string, topicQuota int, diskSize int, ioMax int,
	eipMax int, paidType int, specType string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaInstance(id)
		if err != nil {
			return errs.WrapError(err)
		}

		// Wait for all variables be equal.
		if object.InstanceId == id && object.TopicNumLimit == topicQuota && object.DiskSize == diskSize && object.IoMax == ioMax && object.EipMax == eipMax && object.PaidType == paidType && object.SpecType == specType {
			return nil
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.InstanceId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		// Process wait for running.
		if object.InstanceId == id && status == Running {

			// ServiceStatus equals 5, means the server is in service.
			if object.ServiceStatus == 5 {
				return nil
			}

		} else if object.InstanceId == id {

			// If target status is not deleted and found a instance, return.
			if status != Deleted {
				return nil
			} else {
				// ServiceStatus equals 10, means the server is in released.
				if object.ServiceStatus == 10 {
					return nil
				}
			}
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.InstanceId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaConsumerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaConsumerGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if object.InstanceId+":"+object.ConsumerId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.InstanceId+":"+object.ConsumerId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) KafkaTopicStatusRefreshFunc(id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlikafkaTopicStatus(id)
		if err != nil {
			if !errs.IsExpectedErrors(err, []string{errs.ResourceNotfound}) {
				return nil, "", errs.WrapError(err)
			}
		}

		if object.OffsetTable.OffsetTableItem != nil && len(object.OffsetTable.OffsetTableItem) > 0 {
			return object, "Running", errs.WrapError(err)
		}

		return object, "Creating", nil
	}
}

func (s *AlikafkaService) WaitForAlikafkaTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaTopic(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if object.InstanceId+":"+object.Topic == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.InstanceId+":"+object.Topic, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaSaslUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	instanceId := parts[0]
	for {
		object, err := s.DescribeAlikafkaSaslUser(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if instanceId+":"+object.Username == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, instanceId+":"+object.Username, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaSaslAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 6)
	if err != nil {
		return errs.WrapError(err)
	}
	instanceId := parts[0]
	for {
		object, err := s.DescribeAlikafkaSaslAcl(id)
		if err != nil {

			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if instanceId+":"+object.Username+":"+object.AclResourceType+":"+object.AclResourceName+":"+object.AclResourcePatternType+":"+object.AclOperationType == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, instanceId+":"+object.Username, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []alikafkasdk.TagResource, err error) {
	request := alikafkasdk.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []alikafkasdk.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, alikafkasdk.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.ListTagResources(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.Throttling, errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*alikafkasdk.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *AlikafkaService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := alikafkasdk.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = string(resourceType)
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				raw, err := s.client.WithAlikafkaClient(func(client *alikafkasdk.Client) (interface{}, error) {
					return client.UntagResources(request)
				})
				if err != nil {
					if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
		}

		if len(create) > 0 {
			request := alikafkasdk.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = string(resourceType)
			request.RegionId = s.client.RegionId

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				raw, err := s.client.WithAlikafkaClient(func(client *alikafkasdk.Client) (interface{}, error) {
					return client.TagResources(request)
				})
				if err != nil {
					if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
		}

	}

	return nil
}

func (s *AlikafkaService) tagsToMap(tags []alikafkasdk.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *AlikafkaService) ignoreTag(t alikafkasdk.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *AlikafkaService) tagVOTagsToMap(tags []alikafkasdk.TagVO) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.tagVOIgnoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *AlikafkaService) tagVOIgnoreTag(t alikafkasdk.TagVO) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *AlikafkaService) diffTags(oldTags, newTags []alikafkasdk.TagResourcesTag) ([]alikafkasdk.TagResourcesTag, []alikafkasdk.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []alikafkasdk.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *AlikafkaService) tagsFromMap(m map[string]interface{}) []alikafkasdk.TagResourcesTag {
	result := make([]alikafkasdk.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, alikafkasdk.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *AlikafkaService) GetAllowedIpList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAllowedIpList"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlikafkaService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("alikafka", "2019-09-16", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *AlikafkaService) DescribeAliKafkaConsumerGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetConsumerList"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ConsumerList.ConsumerVO", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ConsumerList.ConsumerVO", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["ConsumerId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *AlikafkaService) DescribeAliKafkaSaslUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeSaslUsers"
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"BIZ_INSTANCE_STATUS_ERROR", "BIZ.INSTANCE.STATUS.ERROR"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:SaslUser", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.SaslUserList.SaslUserVO", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SaslUserList.SaslUserVO", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:SaslUser", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Username"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:SaslUser", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AlikafkaService) DescribeAliKafkaInstanceAllowedIpAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetAllowedIpList"

	client := s.client

	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		err = errs.WrapError(err)
		return
	}

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"BIZ_INSTANCE_STATUS_ERROR", "BIZ.INSTANCE.STATUS.ERROR"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:InstanceAllowedIpAttachment", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	var resp interface{}
	allowedType := parts[1]

	switch allowedType {
	case "vpc":
		resp, err = jsonpath.Get("$.AllowedList.VpcList", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AllowedList.VpcList", response)
		}
	case "internet":
		resp, err = jsonpath.Get("$.AllowedList.InternetList", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AllowedList.InternetList", response)
		}
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:InstanceAllowedIpAttachment", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		ipList := v.(map[string]interface{})
		if fmt.Sprint(ipList["PortRange"]) == parts[2] {
			for _, ip := range ipList["AllowedIpList"].([]interface{}) {
				if fmt.Sprint(ip) == parts[3] {
					idExist = true
					return v.(map[string]interface{}), nil
				}
			}
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:InstanceAllowedIpAttachment", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AlikafkaService) DescribeAliKafkaInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetInstanceList"

	client := s.client

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": []string{id},
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.InstanceList.InstanceVO", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.InstanceList.InstanceVO", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:Instance", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == id && fmt.Sprint(v.(map[string]interface{})["ServiceStatus"]) != "10" {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("AliKafka:Instance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AlikafkaService) GetQuotaTip(instanceId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetQuotaTip"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": instanceId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapError(err)
	}
	v, err := jsonpath.Get("$.QuotaData", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, instanceId, "$.QuotaData", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *AlikafkaService) DescribeAliKafkaInstanceByOrderId(orderId string, timeout int) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetInstanceList"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"OrderId":  orderId,
	}

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(10*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, orderId, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.InstanceList.InstanceVO", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, orderId, "$.InstanceList.InstanceVO", response)
		}
		for _, v := range v.([]interface{}) {
			return v.(map[string]interface{}), nil
		}
		if time.Now().After(deadline) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AlikafkaInstance", orderId), errs.NotFoundMsg, errs.ProviderERROR)
		}
	}
}

func (s *AlikafkaService) AliKafkaInstanceStateRefreshFunc(id, attribute string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAliKafkaInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {

			if fmt.Sprint(object[attribute]) == failState {
				return object, fmt.Sprint(object[attribute]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object[attribute])))
			}
		}
		return object, fmt.Sprint(object[attribute]), nil
	}
}

func (s *AlikafkaService) AliKafkaConsumerStateRefreshFunc(id, attribute string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAliKafkaConsumerGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {

			if fmt.Sprint(object[attribute]) == failState {
				return object, fmt.Sprint(object[attribute]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object[attribute])))
			}
		}
		return object, fmt.Sprint(object[attribute]), nil
	}
}

func (s *AlikafkaService) AliKafkaTopicStateRefreshFunc(id, attribute string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlikafkaTopic(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		return object, "existing", nil
	}
}
