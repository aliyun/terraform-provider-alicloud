package gdb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Gdb service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "gdb",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_graph_database_db_instance",
				Factory:  ResourceGraphDatabaseDbInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_graph_database_db_instances",
				Factory:  DataSourceGraphDatabaseDbInstances,
			},
		},
	}
}
