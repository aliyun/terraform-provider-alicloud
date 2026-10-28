package ecp

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ECP (Cloudphone) service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecp",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ecp_instance",
				Factory:  ResourceEcpInstance,
			},
			{
				TypeName: "alicloud_ecp_key_pair",
				Factory:  ResourceEcpKeyPair,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ecp_instances",
				Factory:  DataSourceEcpInstances,
			},
			{
				TypeName: "alicloud_ecp_instance_types",
				Factory:  DataSourceEcpInstanceTypes,
			},
			{
				TypeName: "alicloud_ecp_key_pairs",
				Factory:  DataSourceEcpKeyPairs,
			},
			{
				TypeName: "alicloud_ecp_zones",
				Factory:  DataSourceEcpZones,
			},
		},
	}
}
