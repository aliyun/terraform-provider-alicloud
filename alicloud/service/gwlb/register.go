package gwlb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Gwlb service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Gwlb",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_gwlb_load_balancer",
				Factory:  ResourceLoadBalancer,
			},
			{
				TypeName: "alicloud_gwlb_listener",
				Factory:  ResourceListener,
			},
			{
				TypeName: "alicloud_gwlb_server_group",
				Factory:  ResourceServerGroup,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_gwlb_zones",
				Factory:  DataSourceZones,
			},
		},
	}
}
