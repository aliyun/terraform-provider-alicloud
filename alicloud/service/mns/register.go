package mns

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the mns service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Mns",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_mns_queue", Factory: resourceAlicloudMNSQueue},
			{TypeName: "alicloud_mns_topic", Factory: resourceAlicloudMNSTopic},
			{TypeName: "alicloud_mns_topic_subscription", Factory: resourceAlicloudMNSSubscription},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_mns_queues", Factory: dataSourceAlicloudMNSQueues},
			{TypeName: "alicloud_mns_topics", Factory: dataSourceAlicloudMNSTopics},
			{TypeName: "alicloud_mns_topic_subscriptions", Factory: dataSourceAlicloudMNSTopicSubscriptions},
			{TypeName: "alicloud_mns_service", Factory: dataSourceAlicloudMnsService},
		},
	}
}
