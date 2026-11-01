package rocketmq

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the RocketMQ service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "RocketMQ",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_rocketmq_acl",
				Factory:  ResourceAcl,
			},
			{
				TypeName: "alicloud_rocketmq_account",
				Factory:  ResourceAccount,
			},
			{
				TypeName: "alicloud_rocketmq_topic",
				Factory:  ResourceTopic,
			},
			{
				TypeName: "alicloud_rocketmq_consumer_group",
				Factory:  ResourceConsumerGroup,
			},
			{
				TypeName: "alicloud_rocketmq_instance",
				Factory:  ResourceInstance,
			},
		},
	}
}
