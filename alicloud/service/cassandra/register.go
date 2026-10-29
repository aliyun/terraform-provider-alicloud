package cassandra

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cassandra service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cassandra",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cassandra_cluster", Factory: ResourceCassandraCluster},
			{TypeName: "alicloud_cassandra_data_center", Factory: ResourceCassandraDataCenter},
			{TypeName: "alicloud_cassandra_backup_plan", Factory: ResourceCassandraBackupPlan},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_cassandra_clusters", Factory: DataSourceCassandraClusters},
			{TypeName: "alicloud_cassandra_data_centers", Factory: DataSourceCassandraDataCenters},
			{TypeName: "alicloud_cassandra_zones", Factory: DataSourceCassandraZones},
			{TypeName: "alicloud_cassandra_backup_plans", Factory: DataSourceCassandraBackupPlans},
		},
	}
}
