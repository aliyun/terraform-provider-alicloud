package polardb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PolarDB service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PolarDB",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_polardb_database",
				Factory:  ResourceDatabase,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_polardb_node_classes",
				Factory:  DataSourceNodeClasses,
			},
		},
	}
}
