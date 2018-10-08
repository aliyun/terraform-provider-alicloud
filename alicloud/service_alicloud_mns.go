package alicloud

import (
	"strings"

	"github.com/dxh031/ali_mns"
)

func (client *AliyunClient) MnsQueueManager() (ali_mns.AliQueueManager, error) {

	mnsClient, err := client.Mnsconn()
	if err != nil {
		return nil, err
	}
	queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
	return queueManager, nil
}

func (client *AliyunClient) MnsSubscriptionManager(topicName string) (ali_mns.AliMNSTopic, error) {

	mnsClient, err := client.Mnsconn()
	if err != nil {
		return nil, err
	}
	subscriptionManager := ali_mns.NewMNSTopic(topicName, *mnsClient)
	return subscriptionManager, nil
}

func (client *AliyunClient) MnsTopicManager() (ali_mns.AliTopicManager, error) {

	mnsClient, err := client.Mnsconn()
	if err != nil {
		return nil, err
	}
	topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
	return topicManager, nil
}

func GetTopicNameAndSubscriptionName(subscriptionId string) (string, string) {
	arr := strings.Split(subscriptionId, COLON_SEPARATED)
	return arr[0], arr[1]
}

func SubscriptionNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), SubscriptionNotExist)
}
func TopicNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), TopicNotExist)
}

func QueueNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), QueueNotExist)
}
