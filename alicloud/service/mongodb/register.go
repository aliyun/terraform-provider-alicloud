package mongodb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the MongoDB (Dds) service package registration for the
// unified provider. Factories remain unexported and are referenced directly
// within this package, matching the handwritten stock split-layout precedent.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dds",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_mongodb_instance", Factory: resourceAliCloudMongoDBInstance},
			{TypeName: "alicloud_mongodb_account", Factory: resourceAliCloudMongodbAccount},
			{TypeName: "alicloud_mongodb_audit_policy", Factory: resourceAliCloudMongodbAuditPolicy},
			{TypeName: "alicloud_mongodb_backup", Factory: resourceAliCloudMongodbBackup},
			{TypeName: "alicloud_mongodb_global_security_ip_group", Factory: resourceAliCloudMongodbGlobalSecurityIPGroup},
			{TypeName: "alicloud_mongodb_node", Factory: resourceAliCloudMongodbNode},
			{TypeName: "alicloud_mongodb_private_srv_network_address", Factory: resourceAliCloudMongodbPrivateSrvNetworkAddress},
			{TypeName: "alicloud_mongodb_public_network_address", Factory: resourceAlicloudMongoDBPublicNetworkAddress},
			{TypeName: "alicloud_mongodb_replica_set_role", Factory: resourceAlicloudMongoDBReplicaSetRole},
			{TypeName: "alicloud_mongodb_serverless_instance", Factory: resourceAlicloudMongodbServerlessInstance},
			{TypeName: "alicloud_mongodb_sharding_audit_filter", Factory: resourceAliCloudMongodbShardingAuditFilter},
			{TypeName: "alicloud_mongodb_sharding_instance", Factory: resourceAliCloudMongoDBShardingInstance},
			{TypeName: "alicloud_mongodb_sharding_network_private_address", Factory: resourceAliCloudMongodbShardingNetworkPrivateAddress},
			{TypeName: "alicloud_mongodb_sharding_network_public_address", Factory: resourceAlicloudMongodbShardingNetworkPublicAddress},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_mongodb_instances", Factory: dataSourceAliCloudMongoDBInstances},
			{TypeName: "alicloud_mongo_instances", Factory: dataSourceAliCloudMongoDBInstances},
			{TypeName: "alicloud_mongodb_zones", Factory: dataSourceAlicloudMongoDBZones},
			{TypeName: "alicloud_mongodb_audit_policies", Factory: dataSourceAlicloudMongodbAuditPolicies},
			{TypeName: "alicloud_mongodb_accounts", Factory: dataSourceAlicloudMongodbAccounts},
			{TypeName: "alicloud_mongodb_serverless_instances", Factory: dataSourceAlicloudMongodbServerlessInstances},
			{TypeName: "alicloud_mongodb_sharding_network_public_addresses", Factory: dataSourceAlicloudMongodbShardingNetworkPublicAddresses},
			{TypeName: "alicloud_mongodb_sharding_network_private_addresses", Factory: dataSourceAlicloudMongodbShardingNetworkPrivateAddresses},
			{TypeName: "alicloud_mongodb_sharding_audit_filters", Factory: dataSourceAlicloudMongodbShardingAuditFilters},
			{TypeName: "alicloud_mongodb_backups", Factory: dataSourceAliCloudMongodbBackups},
		},
	}
}
