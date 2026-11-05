package rkvstore

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage wires the R-kvstore product into the provider registry.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "R-kvstore",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_kvstore_instance", Factory: ResourceAliCloudKvstoreInstance},
			{TypeName: "alicloud_kvstore_backup_policy", Factory: resourceAliCloudKvStoreBackupPolicy},
			{TypeName: "alicloud_kvstore_account", Factory: resourceAliCloudRedisAccount},
			{TypeName: "alicloud_kvstore_connection", Factory: resourceAlicloudKvstoreConnection},
			{TypeName: "alicloud_kvstore_audit_log_config", Factory: ResourceAlicloudKvstoreAuditLogConfig},
			{TypeName: "alicloud_redis_global_security_ip_group", Factory: resourceAliCloudRedisGlobalSecurityIpGroup},
			{TypeName: "alicloud_redis_backup", Factory: resourceAliCloudRedisBackup},
			{TypeName: "alicloud_redis_tair_instance", Factory: ResourceAliCloudRedisTairInstance},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_kvstore_instances", Factory: dataSourceAlicloudKvstoreInstances},
			{TypeName: "alicloud_kvstore_zones", Factory: dataSourceAlicloudKVStoreZones},
			{TypeName: "alicloud_kvstore_permission", Factory: dataSourceAlicloudKVStorePermission},
			{TypeName: "alicloud_kvstore_instance_classes", Factory: dataSourceAlicloudKVStoreInstanceClasses},
			{TypeName: "alicloud_kvstore_instance_engines", Factory: dataSourceAlicloudKVStoreInstanceEngines},
			{TypeName: "alicloud_kvstore_connections", Factory: dataSourceAlicloudKvstoreConnections},
			{TypeName: "alicloud_kvstore_accounts", Factory: dataSourceAlicloudKvstoreAccounts},
			{TypeName: "alicloud_redis_global_security_ip_groups", Factory: dataSourceAliCloudRedisGlobalSecurityIpGroups},
		},
	}
}
