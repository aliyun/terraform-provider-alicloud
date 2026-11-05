package gpdb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Gpdb service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "gpdb",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_gpdb_account", Factory: ResourceAccount},
			{TypeName: "alicloud_gpdb_api_key", Factory: ResourceApiKey},
			{TypeName: "alicloud_gpdb_backup_policy", Factory: ResourceBackupPolicy},
			{TypeName: "alicloud_gpdb_connection", Factory: ResourceConnection},
			{TypeName: "alicloud_gpdb_database", Factory: ResourceDatabase},
			{TypeName: "alicloud_gpdb_db_extension", Factory: ResourceDbExtension},
			{TypeName: "alicloud_gpdb_db_instance_ip_array", Factory: ResourceDBInstanceIPArray},
			{TypeName: "alicloud_gpdb_db_instance_plan", Factory: ResourceDbInstancePlan},
			{TypeName: "alicloud_gpdb_db_resource_group", Factory: ResourceDbResourceGroup},
			{TypeName: "alicloud_gpdb_elastic_instance", Factory: ResourceElasticInstance},
			{TypeName: "alicloud_gpdb_external_data_service", Factory: ResourceExternalDataService},
			{TypeName: "alicloud_gpdb_hadoop_data_source", Factory: ResourceHadoopDataSource},
			{TypeName: "alicloud_gpdb_instance", Factory: ResourceInstance},
			{TypeName: "alicloud_gpdb_jdbc_data_source", Factory: ResourceJdbcDataSource},
			{TypeName: "alicloud_gpdb_remote_adb_data_source", Factory: ResourceRemoteADBDataSource},
			{TypeName: "alicloud_gpdb_streaming_data_service", Factory: ResourceStreamingDataService},
			{TypeName: "alicloud_gpdb_streaming_data_source", Factory: ResourceStreamingDataSource},
			{TypeName: "alicloud_gpdb_streaming_job", Factory: ResourceStreamingJob},
			{TypeName: "alicloud_gpdb_supabase_project", Factory: ResourceSupabaseProject},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_gpdb_accounts", Factory: DataSourceAccounts},
			{TypeName: "alicloud_gpdb_api_keys", Factory: DataSourceApiKeys},
			{TypeName: "alicloud_gpdb_data_backups", Factory: DataSourceDataBackups},
			{TypeName: "alicloud_gpdb_db_extensions", Factory: DataSourceDbExtensions},
			{TypeName: "alicloud_gpdb_db_instance_plans", Factory: DataSourceDbInstancePlans},
			{TypeName: "alicloud_gpdb_instances", Factory: DataSourceInstances},
			{TypeName: "alicloud_gpdb_log_backups", Factory: DataSourceLogbackups},
			{TypeName: "alicloud_gpdb_zones", Factory: DataSourceZones},
		},
	}
}
