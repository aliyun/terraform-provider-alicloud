package mnsopen

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage declares the resources and data sources owned by this package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Mns-open",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_message_service_topic", Factory: resourceAliCloudMessageServiceTopic},
			{TypeName: "alicloud_message_service_account_logging", Factory: resourceAliCloudMessageServiceAccountLogging},
			{TypeName: "alicloud_message_service_endpoint", Factory: resourceAliCloudMessageServiceEndpoint},
			{TypeName: "alicloud_message_service_event_rule", Factory: resourceAliCloudMessageServiceEventRule},
			{TypeName: "alicloud_message_service_service", Factory: resourceAliCloudMessageServiceService},
			{TypeName: "alicloud_message_service_queue", Factory: resourceAliCloudMessageServiceQueue},
			{TypeName: "alicloud_message_service_endpoint_acl", Factory: resourceAliCloudMessageServiceEndpointAcl},
			{TypeName: "alicloud_message_service_subscription", Factory: resourceAliCloudMessageServiceSubscription},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_message_service_subscriptions", Factory: dataSourceAlicloudMessageServiceSubscriptions},
			{TypeName: "alicloud_message_service_queues", Factory: dataSourceAlicloudMessageServiceQueues},
			{TypeName: "alicloud_message_service_topics", Factory: dataSourceAlicloudMessageServiceTopics},
		},
	}
}
