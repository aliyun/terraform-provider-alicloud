package eci

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eci",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_eci_container_group",
				Factory:  ResourceEciContainerGroup,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_eci_container_groups",
				Factory:  DataSourceEciContainerGroups,
			},
		},
	}
}
