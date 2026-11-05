package eci

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Eci service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eci",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_eci_container_groups",
				Factory:  DataSourceEciContainerGroups,
			},
			{
				TypeName: "alicloud_eci_image_caches",
				Factory:  DataSourceEciImageCaches,
			},
			{
				TypeName: "alicloud_eci_virtual_nodes",
				Factory:  DataSourceEciVirtualNodes,
			},
			{
				TypeName: "alicloud_eci_zones",
				Factory:  DataSourceEciZones,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_eci_container_group",
				Factory:  ResourceEciContainerGroup,
			},
			{
				TypeName: "alicloud_eci_image_cache",
				Factory:  ResourceEciImageCache,
			},
			{
				TypeName: "alicloud_eci_openapi_image_cache",
				Factory:  ResourceEciImageCache,
			},
			{
				TypeName: "alicloud_eci_virtual_node",
				Factory:  ResourceEciVirtualNode,
			},
		},
	}
}
