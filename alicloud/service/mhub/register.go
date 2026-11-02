package mhub

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the mhub service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Mhub",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_mhub_app",
				Factory:  ResourceApp,
			},
			{
				TypeName: "alicloud_mhub_product",
				Factory:  ResourceProduct,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_mhub_apps",
				Factory:  DataSourceApps,
			},
			{
				TypeName: "alicloud_mhub_products",
				Factory:  DataSourceProducts,
			},
		},
	}
}
