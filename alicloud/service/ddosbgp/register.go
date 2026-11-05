package ddosbgp

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Ddosbgp service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ddosbgp",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ddos_bgp_policy",
				Factory:  resourceAliCloudDdosBgpPolicy,
			},
			{
				TypeName: "alicloud_ddosbgp_instance",
				Factory:  resourceAliCloudDdosBgpInstance,
			},
			{
				TypeName: "alicloud_ddosbgp_ip",
				Factory:  resourceAliCloudDdosBgpIp,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ddosbgp_instances",
				Factory:  dataSourceAlicloudDdosbgpInstances,
			},
			{
				TypeName: "alicloud_ddosbgp_ips",
				Factory:  dataSourceAlicloudDdosbgpIps,
			},
		},
	}
}
