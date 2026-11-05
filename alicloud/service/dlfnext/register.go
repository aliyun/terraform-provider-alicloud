package dlfnext

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DlfNext
// (Data Lake Formation, 2025-03-10) service. It registers three Framework
// resources — catalog, database and table — which model a catalog → database
// → table dependency chain.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DlfNext",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_dlf_catalog",
				Factory:  NewCatalogResource,
			},
			{
				TypeName: "alicloud_dlf_database",
				Factory:  NewDatabaseResource,
			},
			{
				TypeName: "alicloud_dlf_table",
				Factory:  NewTableResource,
			},
		},
	}
}
