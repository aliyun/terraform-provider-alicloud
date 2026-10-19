package cms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CMS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_cms_notify_strategies",
				Factory:  NewNotifyStrategiesDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_cms_notify_strategy",
				Factory:  NewNotifyStrategyResource,
			},
		},
	}
}
