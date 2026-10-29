package lindorm

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Lindorm service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Lindorm",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_lindorm_instance",
				Factory:  ResourceLindormInstance,
			},
			{
				TypeName: "alicloud_lindorm_instance_v2",
				Factory:  ResourceLindormInstanceV2,
			},
			{
				TypeName: "alicloud_lindorm_public_network",
				Factory:  ResourceLindormPublicNetwork,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_lindorm_instances",
				Factory:  DataSourceLindormInstances,
			},
		},
	}
}
