package rds

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the RDS product's SDKv2 resources and data sources
// plus the Framework data source that already lived in this package. The
// factories are the unexported SDKv2 builders moved here from the flat layout;
// they stay unexported because registration happens within the package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Rds",
		SDKResources: []conns.SDKResource{
			// db (legacy RDS type names)
			{TypeName: "alicloud_db_database", Factory: resourceAliCloudRdsDatabase},
			{TypeName: "alicloud_db_account", Factory: resourceAliCloudRdsAccount},
			{TypeName: "alicloud_db_account_privilege", Factory: resourceAlicloudDBAccountPrivilege},
			{TypeName: "alicloud_db_backup_policy", Factory: resourceAlicloudDBBackupPolicy},
			{TypeName: "alicloud_db_connection", Factory: resourceAlicloudDBConnection},
			{TypeName: "alicloud_db_read_write_splitting_connection", Factory: resourceAlicloudDBReadWriteSplittingConnection},
			{TypeName: "alicloud_db_instance", Factory: resourceAliCloudDBInstance},
			{TypeName: "alicloud_db_readonly_instance", Factory: resourceAlicloudDBReadonlyInstance},
			// rds
			{TypeName: "alicloud_rds_custom_disk_attachment", Factory: resourceAliCloudRdsCustomDiskAttachment},
			{TypeName: "alicloud_rds_custom_disk", Factory: resourceAliCloudRdsCustomDisk},
			{TypeName: "alicloud_rds_custom_deployment_set", Factory: resourceAliCloudRdsCustomDeploymentSet},
			{TypeName: "alicloud_rds_custom", Factory: resourceAliCloudRdsCustom},
			{TypeName: "alicloud_rds_whitelist_template", Factory: resourceAliCloudRdsWhitelistTemplate},
			{TypeName: "alicloud_rds_backup", Factory: resourceAliCloudRdsBackup},
			{TypeName: "alicloud_rds_db_proxy", Factory: resourceAlicloudRdsDBProxy},
			{TypeName: "alicloud_rds_db_proxy_public", Factory: resourceAliCloudRdsDBProxyPublic},
			{TypeName: "alicloud_rds_clone_db_instance", Factory: resourceAlicloudRdsCloneDbInstance},
			{TypeName: "alicloud_rds_upgrade_db_instance", Factory: resourceAlicloudRdsUpgradeDbInstance},
			{TypeName: "alicloud_rds_instance_cross_backup_policy", Factory: resourceAlicloudRdsInstanceCrossBackupPolicy},
			{TypeName: "alicloud_rds_ddr_instance", Factory: resourceAlicloudRdsDdrInstance},
			{TypeName: "alicloud_rds_parameter_group", Factory: resourceAlicloudRdsParameterGroup},
			{TypeName: "alicloud_rds_account", Factory: resourceAliCloudRdsAccount},
			{TypeName: "alicloud_rds_db_node", Factory: resourceAlicloudRdsDBNode},
			{TypeName: "alicloud_rds_db_instance_endpoint", Factory: resourceAlicloudRdsDBInstanceEndpoint},
			{TypeName: "alicloud_rds_db_instance_endpoint_address", Factory: resourceAlicloudRdsDBInstanceEndpointAddress},
			{TypeName: "alicloud_rds_service_linked_role", Factory: resourceAlicloudRdsServiceLinkedRole},
		},
		SDKDataSources: []conns.SDKDataSource{
			// db (legacy RDS data source type names)
			{TypeName: "alicloud_db_zones", Factory: dataSourceAlicloudDBZones},
			{TypeName: "alicloud_db_instances", Factory: dataSourceAlicloudDBInstances},
			{TypeName: "alicloud_db_instance_engines", Factory: dataSourceAlicloudDBInstanceEngines},
			{TypeName: "alicloud_db_instance_classes", Factory: dataSourceAliCloudDBInstanceClasses},
			{TypeName: "alicloud_db_instance_class_infos", Factory: dataSourceAlicloudDBInstanceClassInfos},
			{TypeName: "alicloud_instance_keywords", Factory: dataSourceAlicloudInstanceKeywords},
			// rds
			{TypeName: "alicloud_rds_backups", Factory: dataSourceAlicloudRdsBackups},
			{TypeName: "alicloud_rds_modify_parameter_logs", Factory: dataSourceAlicloudRdsModifyParameterLogs},
			{TypeName: "alicloud_rds_parameter_groups", Factory: dataSourceAlicloudRdsParameterGroups},
			{TypeName: "alicloud_rds_collation_time_zones", Factory: dataSourceAlicloudRdsCollationTimeZones},
			{TypeName: "alicloud_rds_accounts", Factory: dataSourceAlicloudRdsAccounts},
			{TypeName: "alicloud_rds_cross_regions", Factory: dataSourceAlicloudRdsCrossRegions},
			{TypeName: "alicloud_rds_cross_region_backups", Factory: dataSourceAlicloudRdsCrossRegionBackups},
			{TypeName: "alicloud_rds_character_set_names", Factory: dataSourceAlicloudRdsCharacterSetNames},
			{TypeName: "alicloud_rds_slots", Factory: dataSourceAlicloudRdsSlots},
			{TypeName: "alicloud_rds_class_details", Factory: dataSourceAlicloudRdsClassDetails},
		},
		DataSources: []conns.DataSource{
			{TypeName: "alicloud_rds_db_instance_performance", Factory: NewDbInstancePerformanceDataSource},
		},
	}
}
