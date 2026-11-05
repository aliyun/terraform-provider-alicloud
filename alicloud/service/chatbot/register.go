package chatbot

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Chatbot service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Chatbot",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_chatbot_publish_task",
				Factory:  ResourcePublishTask,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_chatbot_agents",
				Factory:  DataSourceAgents,
			},
		},
	}
}
