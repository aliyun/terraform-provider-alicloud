package oceanbasepro

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the OceanBasePro service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "OceanBasePro",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ocean_base_instance",
				Factory:  ResourceAliCloudOceanBaseInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ocean_base_instances",
				Factory:  DataSourceAlicloudOceanBaseInstances,
			},
		},
	}
}
