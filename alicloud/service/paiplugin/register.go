package paiplugin

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAIPlugin service.
// Product code "PaiPlugin", API version "2022-01-12" (ROA style).
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiPlugin",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_plugin_schedule",
				Factory:  NewPluginScheduleResource,
			},
		},
	}
}
