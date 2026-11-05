package drds

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Drds service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Drds",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_drds_instance",
				Factory:  resourceAlicloudDRDSInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_drds_instances",
				Factory:  dataSourceAlicloudDRDSInstances,
			},
		},
	}
}
