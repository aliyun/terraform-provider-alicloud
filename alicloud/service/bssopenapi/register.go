package bssopenapi

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage is a placeholder for the bssopenapi service; declarations move here when the product adopts the split layout.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "BssOpenApi",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_bss_open_api_products",
				Factory:  DataSourceProducts,
			},
			{
				TypeName: "alicloud_bss_open_api_pricing_modules",
				Factory:  DataSourcePricingModules,
			},
		},
	}
}
