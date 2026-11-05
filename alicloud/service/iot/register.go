package iot

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Iot service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Iot",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_iot_device_group",
				Factory:  ResourceDeviceGroup,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_iot_service",
				Factory:  DataSourceService,
			},
			{
				TypeName: "alicloud_iot_device_groups",
				Factory:  DataSourceDeviceGroups,
			},
		},
	}
}
