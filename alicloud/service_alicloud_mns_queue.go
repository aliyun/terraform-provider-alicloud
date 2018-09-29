package alicloud

import "github.com/dxh031/ali_mns"

func (client *AliyunClient) MnsQueueManager() (ali_mns.AliQueueManager, error) {

	mnsClient, err := client.Mnsconn()
	if err != nil {
		return nil, err
	}
	queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
	return queueManager, nil
}
