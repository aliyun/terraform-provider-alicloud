package datahub

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataHub service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DataHub",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_datahub_project",
				Factory:  ResourceProject,
			},
			{
				TypeName: "alicloud_datahub_subscription",
				Factory:  ResourceSubscription,
			},
			{
				TypeName: "alicloud_datahub_topic",
				Factory:  ResourceTopic,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_datahub_service",
				Factory:  DataSourceService,
			},
		},
	}
}
