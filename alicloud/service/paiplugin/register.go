package paiplugin

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAIPlugin service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PAIPlugin",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_plugin_group",
				Factory:  NewPluginGroupResource,
			},
		},
	}
}
