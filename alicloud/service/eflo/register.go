package eflo

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Eflo service.
// The DataSources slice is consumed by the framework provider server; the
// SDKv2 eflo resources/data sources remain registered through assembleProvider.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eflo",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_eflo_er_route_entries",
				Factory:  NewErRouteEntriesDataSource,
			},
		},
	}
}
