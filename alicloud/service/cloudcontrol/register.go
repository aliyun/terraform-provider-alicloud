package cloudcontrol

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CloudControl service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "cloudcontrol",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cloud_control_resource",
				Factory:  ResourceCloudControlResource,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cloud_control_prices",
				Factory:  DataSourceCloudControlPrices,
			},
			{
				TypeName: "alicloud_cloud_control_products",
				Factory:  DataSourceCloudControlProducts,
			},
			{
				TypeName: "alicloud_cloud_control_resource_types",
				Factory:  DataSourceCloudControlResourceTypes,
			},
		},
	}
}
