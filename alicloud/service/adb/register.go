package adb

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the adb service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "adb",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_adb_account", Factory: ResourceAdbAccount},
			{TypeName: "alicloud_adb_backup_policy", Factory: ResourceAdbBackupPolicy},
			{TypeName: "alicloud_adb_connection", Factory: ResourceAdbConnection},
			{TypeName: "alicloud_adb_db_cluster", Factory: ResourceAdbDbCluster},
			{TypeName: "alicloud_adb_db_cluster_lake_version", Factory: ResourceAdbDbClusterLakeVersion},
			{TypeName: "alicloud_adb_lake_account", Factory: ResourceAdbLakeAccount},
			{TypeName: "alicloud_adb_resource_group", Factory: ResourceAdbResourceGroup},
			// alicloud_adb_cluster is a public alias of alicloud_adb_db_cluster.
			{TypeName: "alicloud_adb_cluster", Factory: ResourceAdbDbCluster},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_adb_db_clusters", Factory: DataSourceAdbDbClusters},
			{TypeName: "alicloud_adb_db_cluster_lake_versions", Factory: DataSourceAdbDbClusterLakeVersions},
			{TypeName: "alicloud_adb_resource_groups", Factory: DataSourceAdbResourceGroups},
			{TypeName: "alicloud_adb_zones", Factory: DataSourceAdbZones},
			// alicloud_adb_clusters is a public alias of alicloud_adb_db_clusters.
			{TypeName: "alicloud_adb_clusters", Factory: DataSourceAdbDbClusters},
		},
	}
}
