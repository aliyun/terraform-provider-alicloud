package dcdn

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Dcdn service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dcdn",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dcdn_domain",
				Factory:  ResourceAliCloudDcdnDomain,
			},
		},
	}
}
