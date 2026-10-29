package mns

import (
	"strings"
	"time"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// Status mirrors the root-package lifecycle status string used by the legacy
// WaitFor helpers; it is retained locally so the migrated polling semantics stay
// identical without importing the root alicloud package.
type Status string

const (
	Deleted         = Status("Deleted")
	DefaultTimeout  = 120
	COLON_SEPARATED = ":"
)

type MnsService struct {
	client *connectivity.AliyunClient
}

func NewMnsService(client *connectivity.AliyunClient) *MnsService {
	return &MnsService{client: client}
}

func (s *MnsService) SubscriptionNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), "SubscriptionNotExist")
}
func (s *MnsService) TopicNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), "TopicNotExist")
}

func (s *MnsService) QueueNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), "QueueNotExist")
}

func (s *MnsService) DescribeMnsQueue(id string) (response ali_mns.QueueAttribute, err error) {
	raw, err := s.client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
		return queueManager.GetQueueAttributes(id)
	})
	if err != nil {
		if s.QueueNotExistFunc(err) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliMnsERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetQueueAttributes", errs.AliMnsERROR)
	}
	logs.AddDebug("GetQueueAttributes", raw)
	response, _ = raw.(ali_mns.QueueAttribute)
	if response.QueueName == "" {
		return response, errs.WrapErrorf(errs.NotFoundErr("MnsQueue", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return
}

func (s *MnsService) WaitForMnsQueue(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeMnsQueue(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.QueueName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.QueueName, id, errs.ProviderERROR)
		}
	}
}

func (s *MnsService) DescribeMnsTopic(id string) (*ali_mns.TopicAttribute, error) {
	response := &ali_mns.TopicAttribute{}
	raw, err := s.client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
		return topicManager.GetTopicAttributes(id)
	})
	if err != nil {
		if s.TopicNotExistFunc(err) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliMnsERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetTopicAttributes", errs.AliMnsERROR)
	}
	logs.AddDebug("GetTopicAttributes", raw)
	resp, _ := raw.(ali_mns.TopicAttribute)
	if resp.TopicName == "" {
		return response, errs.WrapErrorf(errs.NotFoundErr("MnsTopic", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return &resp, nil
}

func (s *MnsService) WaitForMnsTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeMnsTopic(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.TopicName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.TopicName, id, errs.ProviderERROR)
		}
	}
}

func (s *MnsService) DescribeMnsTopicSubscription(id string) (*ali_mns.SubscriptionAttribute, error) {
	response := &ali_mns.SubscriptionAttribute{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return response, errs.WrapError(err)
	}
	topicName, name := parts[0], parts[1]

	raw, err := s.client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
		return subscriptionManager.GetSubscriptionAttributes(name)
	})
	if err != nil {
		if s.TopicNotExistFunc(err) || s.SubscriptionNotExistFunc(err) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliMnsERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetSubscriptionAttributes", errs.AliMnsERROR)
	}
	logs.AddDebug("GetSubscriptionAttributes", raw)
	resp, _ := raw.(ali_mns.SubscriptionAttribute)
	response = &resp
	if response.SubscriptionName == "" {
		return response, errs.WrapErrorf(errs.NotFoundErr("MnsTopicSubscription", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}

func (s *MnsService) WaitForMnsTopicSubscription(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeMnsTopicSubscription(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.TopicName+":"+object.SubscriptionName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.SubscriptionName, id, errs.ProviderERROR)
		}
	}
}
