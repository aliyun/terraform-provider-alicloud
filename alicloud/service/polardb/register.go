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
			{
				TypeName: "alicloud_polardb_clusters",
				Factory:  DataSourceClusters,
			},
			{
				TypeName: "alicloud_polardb_endpoints",
				Factory:  DataSourceEndpoints,
			},
			{
				TypeName: "alicloud_polardb_accounts",
				Factory:  DataSourceAccounts,
			},
			{
				TypeName: "alicloud_polardb_databases",
				Factory:  DataSourceDatabases,
			},
			{
				TypeName: "alicloud_polardb_zones",
				Factory:  DataSourceZones,
			},
			{
				TypeName: "alicloud_polardb_gateway_model_services",
				Factory:  DataSourceGatewayModelServices,
			},
			{
				TypeName: "alicloud_polardb_gateway_model_apis",
				Factory:  DataSourceGatewayModelAPIs,
			},
			{
				TypeName: "alicloud_polardb_gateway_cost_rules",
				Factory:  DataSourceGatewayCostRules,
			},
			{
				TypeName: "alicloud_polardb_gateway_consumer_groups",
				Factory:  DataSourceGatewayConsumerGroups,
			},
			{
				TypeName: "alicloud_polardb_gateway_consumers",
				Factory:  DataSourceGatewayConsumers,
			},
			{
				TypeName: "alicloud_polardb_gateway_budget_policies",
				Factory:  DataSourceGatewayBudgetPolicies,
			},
			{
				TypeName: "alicloud_polardb_gateway_rate_limit_policies",
				Factory:  DataSourceGatewayRateLimitPolicies,
			},
			{
				TypeName: "alicloud_polardb_global_database_networks",
				Factory:  DataSourceGlobalDatabaseNetworks,
			},
			{
				TypeName: "alicloud_polardb_parameter_groups",
				Factory:  DataSourceParameterGroups,
			},
			{
				TypeName: "alicloud_polardb_polar_fs_instances",
				Factory:  DataSourcePolarFsInstances,
			},
			{
				TypeName: "alicloud_polardb_lakebase_tenant_token",
				Factory:  DataSourceLakebaseTenantToken,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_polar_db_extension",
				Factory:  ResourceExtension,
			},
			{
				TypeName: "alicloud_polardb_global_security_ip_group",
				Factory:  ResourceGlobalSecurityIpGroup,
			},
			{
				TypeName: "alicloud_polardb_cluster",
				Factory:  ResourceCluster,
			},
			{
				TypeName: "alicloud_polardb_cluster_endpoint",
				Factory:  ResourceClusterEndpoint,
			},
			{
				TypeName: "alicloud_polardb_backup_policy",
				Factory:  ResourceBackupPolicy,
			},
			{
				TypeName: "alicloud_polardb_database",
				Factory:  ResourceDatabase,
			},
			{
				TypeName: "alicloud_polardb_account",
				Factory:  ResourceAccount,
			},
			{
				TypeName: "alicloud_polardb_account_privilege",
				Factory:  ResourceAccountPrivilege,
			},
			{
				TypeName: "alicloud_polardb_endpoint",
				Factory:  ResourceEndpoint,
			},
			{
				TypeName: "alicloud_polardb_dynamo_table",
				Factory:  ResourceDynamoTable,
			},
			{
				TypeName: "alicloud_polardb_dynamo_item",
				Factory:  ResourceDynamoItem,
			},
			{
				TypeName: "alicloud_polardb_endpoint_address",
				Factory:  ResourceEndpointAddress,
			},
			{
				TypeName: "alicloud_polardb_primary_endpoint",
				Factory:  ResourcePrimaryEndpoint,
			},
			{
				TypeName: "alicloud_polardb_application",
				Factory:  ResourceApplication,
			},
			{
				TypeName: "alicloud_polardb_application_endpoint",
				Factory:  ResourceApplicationEndpoint,
			},
			{
				TypeName: "alicloud_polardb_batch_task",
				Factory:  ResourceBatchTask,
			},
			{
				TypeName: "alicloud_polardb_aicluster",
				Factory:  ResourceAICluster,
			},
			{
				TypeName: "alicloud_polardb_gateway",
				Factory:  ResourceGateway,
			},
			{
				TypeName: "alicloud_polardb_polar_fs",
				Factory:  ResourcePolarFs,
			},
			{
				TypeName: "alicloud_polardb_polar_fs_endpoint",
				Factory:  ResourcePolarFsEndpoint,
			},
			{
				TypeName: "alicloud_polardb_gateway_model_service",
				Factory:  ResourceGatewayModelService,
			},
			{
				TypeName: "alicloud_polardb_gateway_model_api",
				Factory:  ResourceGatewayModelAPI,
			},
			{
				TypeName: "alicloud_polardb_gateway_cost_rule",
				Factory:  ResourceGatewayCostRule,
			},
			{
				TypeName: "alicloud_polardb_gateway_consumer_group",
				Factory:  ResourceGatewayConsumerGroup,
			},
			{
				TypeName: "alicloud_polardb_gateway_consumer",
				Factory:  ResourceGatewayConsumer,
			},
			{
				TypeName: "alicloud_polardb_gateway_budget_policy",
				Factory:  ResourceGatewayBudgetPolicy,
			},
			{
				TypeName: "alicloud_polardb_gateway_rate_limit_policy",
				Factory:  ResourceGatewayRateLimitPolicy,
			},
			{
				TypeName: "alicloud_polardb_global_database_network",
				Factory:  ResourceGlobalDatabaseNetwork,
			},
			{
				TypeName: "alicloud_polardb_parameter_group",
				Factory:  ResourceParameterGroup,
			},
			{
				TypeName: "alicloud_polardb_zonal_db_cluster",
				Factory:  ResourceZonalCluster,
			},
			{
				TypeName: "alicloud_polardb_zonal_endpoint",
				Factory:  ResourceZonalEndpoint,
			},
			{
				TypeName: "alicloud_polardb_zonal_account",
				Factory:  ResourceZonalAccount,
			},
		},
	}
}
