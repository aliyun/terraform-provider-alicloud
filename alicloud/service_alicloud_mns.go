package alicloud

import (
	"strings"
)

type MnsService struct {
}

func (s *MnsService) GetTopicNameAndSubscriptionName(subscriptionId string) (string, string) {
	arr := strings.Split(subscriptionId, COLON_SEPARATED)
	return arr[0], arr[1]
}

func (s *MnsService) SubscriptionNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), SubscriptionNotExist)
}
func (s *MnsService) TopicNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), TopicNotExist)
}

func (s *MnsService) QueueNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), QueueNotExist)
}
