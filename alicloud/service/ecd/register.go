package ecd

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ECD service.
//
// Only alicloud_ecd_ad_connector_directory is split into this package; the
// remaining ECD resources and data sources stay in the flat root package until
// their own per-type migrations.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecd",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ecd_ad_connector_directory",
				Factory:  ResourceAdConnectorDirectory,
			},
		},
	}
}
