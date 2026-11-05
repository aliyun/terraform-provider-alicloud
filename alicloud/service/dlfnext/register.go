package dlfnext

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DlfNext service.
// The factories live in this package and are wired into the provider through
// the registry.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DlfNext",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_dlf_next_catalog", Factory: ResourceCatalog},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_dlf_next_catalogs", Factory: DataSourceCatalogs},
		},
	}
}
