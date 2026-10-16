package eas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the EAS service.
//
// The resources here call the PAI-EAS POP product ("EAS", version "2021-07-01")
// over the ROA client, the same channel the legacy SDK v2 alicloud_pai_service
// uses. Resource type names follow the PAI governance prefix (alicloud_pai_*)
// rather than the API product code, because the governed object is PAI/Log.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eas",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_resource_log",
				Factory:  NewPaiResourceLogResource,
			},
		},
	}
}
