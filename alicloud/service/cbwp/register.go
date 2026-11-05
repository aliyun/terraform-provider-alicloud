package cbwp

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cbwp service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cbwp",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_common_bandwidth_package",
				Factory:  ResourceCommonBandwidthPackage,
			},
			{
				TypeName: "alicloud_common_bandwidth_package_attachment",
				Factory:  ResourceCommonBandwidthPackageAttachment,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_common_bandwidth_packages",
				Factory:  DataSourceCommonBandwidthPackages,
			},
		},
	}
}
