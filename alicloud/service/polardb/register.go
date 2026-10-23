package polardb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PolarDB service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PolarDB",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_polardb_node_classes",
				Factory:  DataSourceNodeClasses,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_polardb_db_node",
				Factory:  NewDbNodeResource,
			},
		},
	}
}
