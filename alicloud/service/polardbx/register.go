package polardbx

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the polardbx service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "polardbx",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_drds_polardbx_instance",
				Factory:  resourceAliCloudDrdsPolardbxInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_drds_polardbx_instances",
				Factory:  dataSourceAliCloudDrdsPolardbxInstances,
			},
		},
	}
}
