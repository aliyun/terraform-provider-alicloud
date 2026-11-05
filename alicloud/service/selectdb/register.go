package selectdb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "selectdb",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_selectdb_db_cluster",
				Factory:  ResourceDbCluster,
			},
			{
				TypeName: "alicloud_selectdb_db_instance",
				Factory:  ResourceDbInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_selectdb_db_clusters",
				Factory:  DataSourceDbClusters,
			},
			{
				TypeName: "alicloud_selectdb_db_instances",
				Factory:  DataSourceDbInstances,
			},
		},
	}
}
