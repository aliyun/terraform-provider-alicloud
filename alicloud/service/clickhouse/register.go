package clickhouse

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the clickhouse service.
// SDKv2 resources and data sources migrated from the root alicloud package; the
// V2 (2023-05-22) and V1 (2019-11-11) API helpers live in service.go.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "clickhouse",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_click_house_account", Factory: ResourceClickHouseAccount},
			{TypeName: "alicloud_click_house_backup_policy", Factory: ResourceClickHouseBackupPolicy},
			{TypeName: "alicloud_click_house_db_cluster", Factory: ResourceClickHouseDbCluster},
			{TypeName: "alicloud_click_house_enterprise_db_cluster", Factory: ResourceClickHouseEnterpriseDbCluster},
			{TypeName: "alicloud_click_house_enterprise_db_cluster_account", Factory: ResourceClickHouseEnterpriseDBClusterAccount},
			{TypeName: "alicloud_click_house_enterprise_db_cluster_backup_policy", Factory: ResourceClickHouseEnterpriseDbClusterBackupPolicy},
			{TypeName: "alicloud_click_house_enterprise_db_cluster_computing_group", Factory: ResourceClickHouseEnterpriseDbClusterComputingGroup},
			{TypeName: "alicloud_click_house_enterprise_db_cluster_public_endpoint", Factory: ResourceClickHouseEnterpriseDbClusterPublicEndpoint},
			{TypeName: "alicloud_click_house_enterprise_db_cluster_security_ip", Factory: ResourceClickHouseEnterpriseDbClusterSecurityIP},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_click_house_accounts", Factory: DataSourceClickHouseAccounts},
			{TypeName: "alicloud_click_house_backup_policies", Factory: DataSourceClickHouseBackupPolicies},
			{TypeName: "alicloud_click_house_db_clusters", Factory: DataSourceClickHouseDbClusters},
			{TypeName: "alicloud_click_house_regions", Factory: DataSourceClickHouseRegions},
		},
	}
}
