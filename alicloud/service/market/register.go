package market

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Market service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Market",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_market_order",
				Factory:  ResourceOrder,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_market_product",
				Factory:  DataSourceProduct,
			},
			{
				TypeName: "alicloud_market_products",
				Factory:  DataSourceProducts,
			},
		},
	}
}
