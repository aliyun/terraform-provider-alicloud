package mnsopen

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func NewMessageServiceServiceV2(client *connectivity.AliyunClient) *MessageServiceServiceV2 {
	return &MessageServiceServiceV2{client: client}
}
