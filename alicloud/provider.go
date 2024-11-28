package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/aliyun/credentials-go/credentials"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/go-homedir"
)

// Provider returns a schema.Provider for alicloud
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_ACCESS_KEY", "ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABACLOUD_ACCESS_KEY_ID"}, nil),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_SECRET_KEY", "ALIBABA_CLOUD_ACCESS_KEY_SECRET", "ALIBABACLOUD_ACCESS_KEY_SECRET"}, nil),
				Description: descriptions["secret_key"],
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_SECURITY_TOKEN", "ALIBABA_CLOUD_SECURITY_TOKEN", "ALIBABACLOUD_SECURITY_TOKEN"}, nil),
				Description: descriptions["security_token"],
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_ECS_ROLE_NAME", "ALIBABA_CLOUD_ECS_METADATA"}, nil),
				Description: descriptions["ecs_role_name"],
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_REGION", "ALIBABA_CLOUD_REGION"}, nil),
				Description: descriptions["region"],
			},
			"ots_instance_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'ots_instance_name' has been deprecated from provider version 1.10.0. New field 'instance_name' of resource 'alicloud_ots_table' instead.",
			},
			"log_endpoint": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'log_endpoint' has been deprecated from provider version 1.28.0. New field 'log' which in nested endpoints instead.",
			},
			"mns_endpoint": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'mns_endpoint' has been deprecated from provider version 1.28.0. New field 'mns' which in nested endpoints instead.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_ACCOUNT_ID", "ALIBABA_CLOUD_ACCOUNT_ID"}, nil),
				Description: descriptions["account_id"],
			},
			"assume_role":           assumeRoleSchema(),
			"sign_version":          signVersionSchema(),
			"assume_role_with_oidc": assumeRoleWithOidcSchema(),
			"fc": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'fc' has been deprecated from provider version 1.28.0. New field 'fc' which in nested endpoints instead.",
			},
			"endpoints": endpointsSchema(),
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_credentials_file"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_SHARED_CREDENTIALS_FILE", "ALIBABA_CLOUD_CREDENTIALS_FILE"}, nil),
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_PROFILE", "ALIBABA_CLOUD_PROFILE"}, nil),
			},
			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_region_validation"],
			},
			"configuration_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["configuration_source"],
				ValidateFunc: validation.StringLenBetween(0, 128),
				DefaultFunc:  schema.EnvDefaultFunc("TF_APPEND_USER_AGENT", ""),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTPS",
				Description:  descriptions["protocol"],
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"client_read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_READ_TIMEOUT", 60000),
				Description: descriptions["client_read_timeout"],
			},
			"client_connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_CONNECT_TIMEOUT", 60000),
				Description: descriptions["client_connect_timeout"],
			},
			"source_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_SOURCE_IP", "ALIBABA_CLOUD_SOURCE_IP"}, nil),
				Description: descriptions["source_ip"],
			},
			"security_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_SECURITY_TRANSPORT", "ALIBABA_CLOUD_SECURITY_TRANSPORT"}, nil),
				//Deprecated:  "It has been deprecated from version 1.136.0 and using new field secure_transport instead.",
			},
			"secure_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_SECURE_TRANSPORT", "ALIBABA_CLOUD_SECURE_TRANSPORT"}, nil),
				Description: descriptions["secure_transport"],
			},
			"credentials_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_CREDENTIALS_URI", "ALIBABA_CLOUD_CREDENTIALS_URI"}, nil),
				Description: descriptions["credentials_uri"],
			},
			"max_retry_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MAX_RETRY_TIMEOUT", 0),
				Description: descriptions["max_retry_timeout"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alicloud_gpdb_data_backups":      dataSourceAliCloudGpdbDataBackups(),
			"alicloud_gpdb_log_backups":       dataSourceAliCloudGpdbLogbackups(),
			"alicloud_governance_baselines":   dataSourceAliCloudGovernanceBaselines(),
			"alicloud_vpn_gateway_zones":      dataSourceAliCloudVPNGatewayZones(),
			"alicloud_account":                dataSourceAlicloudAccount(),
			"alicloud_caller_identity":        dataSourceAlicloudCallerIdentity(),
			"alicloud_images":                 dataSourceAlicloudImages(),
			"alicloud_regions":                dataSourceAlicloudRegions(),
			"alicloud_zones":                  dataSourceAlicloudZones(),
			"alicloud_db_zones":               dataSourceAlicloudDBZones(),
			"alicloud_instance_type_families": dataSourceAlicloudInstanceTypeFamilies(),
			"alicloud_instance_types":         dataSourceAlicloudInstanceTypes(),
			"alicloud_instances":              dataSourceAlicloudInstances(),
			"alicloud_disks":                  dataSourceAlicloudEcsDisks(),
			"alicloud_network_interfaces":     dataSourceAlicloudEcsNetworkInterfaces(),
			"alicloud_snapshots":              dataSourceAlicloudEcsSnapshots(),
			"alicloud_vpcs":                   dataSourceAlicloudVpcs(),
			"alicloud_vswitches":              dataSourceAlicloudVswitches(),
			"alicloud_eips":                   dataSourceAlicloudEipAddresses(),
			"alicloud_key_pairs":              dataSourceAlicloudEcsKeyPairs(),
			"alicloud_kms_keys":               dataSourceAlicloudKmsKeys(),
			"alicloud_kms_ciphertext":         dataSourceAlicloudKmsCiphertext(),
			"alicloud_kms_plaintext":          dataSourceAlicloudKmsPlaintext(),
			"alicloud_dns_resolution_lines":   dataSourceAlicloudDnsResolutionLines(),
			"alicloud_dns_domains":            dataSourceAlicloudAlidnsDomains(),
			"alicloud_dns_groups":             dataSourceAlicloudDnsGroups(),
			"alicloud_dns_records":            dataSourceAlicloudDnsRecords(),
			// alicloud_dns_domain_groups, alicloud_dns_domain_records have been deprecated.
			"alicloud_dns_domain_groups":  dataSourceAlicloudDnsGroups(),
			"alicloud_dns_domain_records": dataSourceAlicloudDnsRecords(),
			// alicloud_ram_account_alias has been deprecated
			"alicloud_ram_account_alias":                                dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_aliases":                              dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_groups":                                       dataSourceAlicloudRamGroups(),
			"alicloud_ram_users":                                        dataSourceAlicloudRamUsers(),
			"alicloud_ram_roles":                                        dataSourceAlicloudRamRoles(),
			"alicloud_ram_policies":                                     dataSourceAlicloudRamPolicies(),
			"alicloud_ram_policy_document":                              dataSourceAlicloudRamPolicyDocument(),
			"alicloud_security_groups":                                  dataSourceAlicloudSecurityGroups(),
			"alicloud_security_group_rules":                             dataSourceAlicloudSecurityGroupRules(),
			"alicloud_slbs":                                             dataSourceAlicloudSlbLoadBalancers(),
			"alicloud_slb_attachments":                                  dataSourceAlicloudSlbAttachments(),
			"alicloud_slb_backend_servers":                              dataSourceAlicloudSlbBackendServers(),
			"alicloud_slb_listeners":                                    dataSourceAlicloudSlbListeners(),
			"alicloud_slb_rules":                                        dataSourceAlicloudSlbRules(),
			"alicloud_slb_server_groups":                                dataSourceAlicloudSlbServerGroups(),
			"alicloud_slb_master_slave_server_groups":                   dataSourceAlicloudSlbMasterSlaveServerGroups(),
			"alicloud_slb_acls":                                         dataSourceAlicloudSlbAcls(),
			"alicloud_slb_server_certificates":                          dataSourceAlicloudSlbServerCertificates(),
			"alicloud_slb_ca_certificates":                              dataSourceAlicloudSlbCaCertificates(),
			"alicloud_slb_domain_extensions":                            dataSourceAlicloudSlbDomainExtensions(),
			"alicloud_slb_zones":                                        dataSourceAlicloudSlbZones(),
			"alicloud_oss_service":                                      dataSourceAlicloudOssService(),
			"alicloud_oss_bucket_objects":                               dataSourceAlicloudOssBucketObjects(),
			"alicloud_oss_buckets":                                      dataSourceAlicloudOssBuckets(),
			"alicloud_ons_instances":                                    dataSourceAlicloudOnsInstances(),
			"alicloud_ons_topics":                                       dataSourceAlicloudOnsTopics(),
			"alicloud_ons_groups":                                       dataSourceAlicloudOnsGroups(),
			"alicloud_alikafka_consumer_groups":                         dataSourceAlicloudAlikafkaConsumerGroups(),
			"alicloud_alikafka_instances":                               dataSourceAlicloudAlikafkaInstances(),
			"alicloud_alikafka_topics":                                  dataSourceAlicloudAlikafkaTopics(),
			"alicloud_alikafka_sasl_users":                              dataSourceAlicloudAlikafkaSaslUsers(),
			"alicloud_alikafka_sasl_acls":                               dataSourceAlicloudAlikafkaSaslAcls(),
			"alicloud_fc_functions":                                     dataSourceAlicloudFcFunctions(),
			"alicloud_file_crc64_checksum":                              dataSourceAlicloudFileCRC64Checksum(),
			"alicloud_fc_services":                                      dataSourceAlicloudFcServices(),
			"alicloud_fc_triggers":                                      dataSourceAlicloudFcTriggers(),
			"alicloud_fc_custom_domains":                                dataSourceAlicloudFcCustomDomains(),
			"alicloud_fc_zones":                                         dataSourceAlicloudFcZones(),
			"alicloud_db_instances":                                     dataSourceAlicloudDBInstances(),
			"alicloud_db_instance_engines":                              dataSourceAlicloudDBInstanceEngines(),
			"alicloud_db_instance_classes":                              dataSourceAlicloudDBInstanceClasses(),
			"alicloud_rds_backups":                                      dataSourceAlicloudRdsBackups(),
			"alicloud_rds_modify_parameter_logs":                        dataSourceAlicloudRdsModifyParameterLogs(),
			"alicloud_pvtz_zones":                                       dataSourceAlicloudPvtzZones(),
			"alicloud_pvtz_zone_records":                                dataSourceAlicloudPvtzZoneRecords(),
			"alicloud_router_interfaces":                                dataSourceAlicloudRouterInterfaces(),
			"alicloud_vpn_gateways":                                     dataSourceAlicloudVpnGateways(),
			"alicloud_vpn_customer_gateways":                            dataSourceAlicloudVpnCustomerGateways(),
			"alicloud_vpn_connections":                                  dataSourceAlicloudVpnConnections(),
			"alicloud_ssl_vpn_servers":                                  dataSourceAlicloudSslVpnServers(),
			"alicloud_ssl_vpn_client_certs":                             dataSourceAlicloudSslVpnClientCerts(),
			"alicloud_mongo_instances":                                  dataSourceAlicloudMongoDBInstances(),
			"alicloud_mongodb_instances":                                dataSourceAlicloudMongoDBInstances(),
			"alicloud_mongodb_zones":                                    dataSourceAlicloudMongoDBZones(),
			"alicloud_gpdb_instances":                                   dataSourceAlicloudGpdbInstances(),
			"alicloud_gpdb_zones":                                       dataSourceAlicloudGpdbZones(),
			"alicloud_kvstore_instances":                                dataSourceAlicloudKvstoreInstances(),
			"alicloud_kvstore_zones":                                    dataSourceAlicloudKVStoreZones(),
			"alicloud_kvstore_permission":                               dataSourceAlicloudKVStorePermission(),
			"alicloud_kvstore_instance_classes":                         dataSourceAlicloudKVStoreInstanceClasses(),
			"alicloud_kvstore_instance_engines":                         dataSourceAlicloudKVStoreInstanceEngines(),
			"alicloud_cen_instances":                                    dataSourceAlicloudCenInstances(),
			"alicloud_cen_bandwidth_packages":                           dataSourceAlicloudCenBandwidthPackages(),
			"alicloud_cen_bandwidth_limits":                             dataSourceAlicloudCenBandwidthLimits(),
			"alicloud_cen_route_entries":                                dataSourceAlicloudCenRouteEntries(),
			"alicloud_cen_region_route_entries":                         dataSourceAlicloudCenRegionRouteEntries(),
			"alicloud_cen_transit_router_route_entries":                 dataSourceAlicloudCenTransitRouterRouteEntries(),
			"alicloud_cen_transit_router_route_table_associations":      dataSourceAlicloudCenTransitRouterRouteTableAssociations(),
			"alicloud_cen_transit_router_route_table_propagations":      dataSourceAlicloudCenTransitRouterRouteTablePropagations(),
			"alicloud_cen_transit_router_route_tables":                  dataSourceAlicloudCenTransitRouterRouteTables(),
			"alicloud_cen_transit_router_vbr_attachments":               dataSourceAlicloudCenTransitRouterVbrAttachments(),
			"alicloud_cen_transit_router_vpc_attachments":               dataSourceAliCloudCenTransitRouterVpcAttachments(),
			"alicloud_cen_transit_routers":                              dataSourceAlicloudCenTransitRouters(),
			"alicloud_cs_kubernetes_clusters":                           dataSourceAlicloudCSKubernetesClusters(),
			"alicloud_cs_managed_kubernetes_clusters":                   dataSourceAlicloudCSManagerKubernetesClusters(),
			"alicloud_cs_edge_kubernetes_clusters":                      dataSourceAlicloudCSEdgeKubernetesClusters(),
			"alicloud_cs_serverless_kubernetes_clusters":                dataSourceAlicloudCSServerlessKubernetesClusters(),
			"alicloud_cs_kubernetes_permissions":                        dataSourceAlicloudCSKubernetesPermissions(),
			"alicloud_cs_kubernetes_addons":                             dataSourceAlicloudCSKubernetesAddons(),
			"alicloud_cs_kubernetes_version":                            dataSourceAlicloudCSKubernetesVersion(),
			"alicloud_cs_kubernetes_addon_metadata":                     dataSourceAlicloudCSKubernetesAddonMetadata(),
			"alicloud_cr_namespaces":                                    dataSourceAlicloudCRNamespaces(),
			"alicloud_cr_repos":                                         dataSourceAlicloudCRRepos(),
			"alicloud_cr_ee_instances":                                  dataSourceAlicloudCrEEInstances(),
			"alicloud_cr_ee_namespaces":                                 dataSourceAlicloudCrEENamespaces(),
			"alicloud_cr_ee_repos":                                      dataSourceAlicloudCrEERepos(),
			"alicloud_cr_ee_sync_rules":                                 dataSourceAlicloudCrEESyncRules(),
			"alicloud_mns_queues":                                       dataSourceAlicloudMNSQueues(),
			"alicloud_mns_topics":                                       dataSourceAlicloudMNSTopics(),
			"alicloud_mns_topic_subscriptions":                          dataSourceAlicloudMNSTopicSubscriptions(),
			"alicloud_api_gateway_service":                              dataSourceAlicloudApiGatewayService(),
			"alicloud_api_gateway_apis":                                 dataSourceAliCloudApiGatewayApis(),
			"alicloud_api_gateway_groups":                               dataSourceAlicloudApiGatewayGroups(),
			"alicloud_api_gateway_apps":                                 dataSourceAlicloudApiGatewayApps(),
			"alicloud_elasticsearch_instances":                          dataSourceAlicloudElasticsearch(),
			"alicloud_elasticsearch_zones":                              dataSourceAlicloudElaticsearchZones(),
			"alicloud_drds_instances":                                   dataSourceAlicloudDRDSInstances(),
			"alicloud_nas_service":                                      dataSourceAlicloudNasService(),
			"alicloud_nas_access_groups":                                dataSourceAlicloudNasAccessGroups(),
			"alicloud_nas_access_rules":                                 dataSourceAlicloudAccessRules(),
			"alicloud_nas_mount_targets":                                dataSourceAlicloudNasMountTargets(),
			"alicloud_nas_file_systems":                                 dataSourceAlicloudFileSystems(),
			"alicloud_nas_protocols":                                    dataSourceAlicloudNasProtocols(),
			"alicloud_cas_certificates":                                 dataSourceAliCloudSslCertificatesServiceCertificates(),
			"alicloud_common_bandwidth_packages":                        dataSourceAlicloudCommonBandwidthPackages(),
			"alicloud_route_tables":                                     dataSourceAlicloudRouteTables(),
			"alicloud_route_entries":                                    dataSourceAlicloudRouteEntries(),
			"alicloud_nat_gateways":                                     dataSourceAlicloudNatGateways(),
			"alicloud_snat_entries":                                     dataSourceAlicloudSnatEntries(),
			"alicloud_forward_entries":                                  dataSourceAlicloudForwardEntries(),
			"alicloud_ddoscoo_instances":                                dataSourceAlicloudDdoscooInstances(),
			"alicloud_ddosbgp_instances":                                dataSourceAlicloudDdosbgpInstances(),
			"alicloud_ess_alarms":                                       dataSourceAlicloudEssAlarms(),
			"alicloud_ess_notifications":                                dataSourceAlicloudEssNotifications(),
			"alicloud_ess_scaling_groups":                               dataSourceAlicloudEssScalingGroups(),
			"alicloud_ess_scaling_rules":                                dataSourceAlicloudEssScalingRules(),
			"alicloud_ess_scaling_configurations":                       dataSourceAlicloudEssScalingConfigurations(),
			"alicloud_ess_lifecycle_hooks":                              dataSourceAlicloudEssLifecycleHooks(),
			"alicloud_ess_scheduled_tasks":                              dataSourceAlicloudEssScheduledTasks(),
			"alicloud_ots_service":                                      dataSourceAlicloudOtsService(),
			"alicloud_ots_instances":                                    dataSourceAlicloudOtsInstances(),
			"alicloud_ots_instance_attachments":                         dataSourceAlicloudOtsInstanceAttachments(),
			"alicloud_ots_tables":                                       dataSourceAlicloudOtsTables(),
			"alicloud_ots_tunnels":                                      dataSourceAlicloudOtsTunnels(),
			"alicloud_ots_secondary_indexes":                            dataSourceAlicloudOtsSecondaryIndexes(),
			"alicloud_ots_search_indexes":                               dataSourceAlicloudOtsSearchIndexes(),
			"alicloud_cloud_connect_networks":                           dataSourceAlicloudCloudConnectNetworks(),
			"alicloud_emr_instance_types":                               dataSourceAlicloudEmrInstanceTypes(),
			"alicloud_emr_disk_types":                                   dataSourceAlicloudEmrDiskTypes(),
			"alicloud_emr_main_versions":                                dataSourceAlicloudEmrMainVersions(),
			"alicloud_sag_acls":                                         dataSourceAlicloudSagAcls(),
			"alicloud_yundun_dbaudit_instance":                          dataSourceAlicloudDbauditInstances(),
			"alicloud_yundun_bastionhost_instances":                     dataSourceAlicloudBastionhostInstances(),
			"alicloud_bastionhost_instances":                            dataSourceAlicloudBastionhostInstances(),
			"alicloud_market_product":                                   dataSourceAlicloudProduct(),
			"alicloud_market_products":                                  dataSourceAlicloudProducts(),
			"alicloud_polardb_clusters":                                 dataSourceAlicloudPolarDBClusters(),
			"alicloud_polardb_node_classes":                             dataSourceAlicloudPolarDBNodeClasses(),
			"alicloud_polardb_endpoints":                                dataSourceAlicloudPolarDBEndpoints(),
			"alicloud_polardb_accounts":                                 dataSourceAlicloudPolarDBAccounts(),
			"alicloud_polardb_databases":                                dataSourceAlicloudPolarDBDatabases(),
			"alicloud_polardb_zones":                                    dataSourceAlicloudPolarDBZones(),
			"alicloud_hbase_instances":                                  dataSourceAlicloudHBaseInstances(),
			"alicloud_hbase_zones":                                      dataSourceAlicloudHBaseZones(),
			"alicloud_hbase_instance_types":                             dataSourceAlicloudHBaseInstanceTypes(),
			"alicloud_adb_clusters":                                     dataSourceAlicloudAdbDbClusters(),
			"alicloud_adb_zones":                                        dataSourceAlicloudAdbZones(),
			"alicloud_cen_flowlogs":                                     dataSourceAlicloudCenFlowlogs(),
			"alicloud_kms_aliases":                                      dataSourceAlicloudKmsAliases(),
			"alicloud_dns_domain_txt_guid":                              dataSourceAlicloudDnsDomainTxtGuid(),
			"alicloud_edas_service":                                     dataSourceAlicloudEdasService(),
			"alicloud_fnf_service":                                      dataSourceAlicloudFnfService(),
			"alicloud_kms_service":                                      dataSourceAlicloudKmsService(),
			"alicloud_sae_service":                                      dataSourceAlicloudSaeService(),
			"alicloud_dataworks_service":                                dataSourceAlicloudDataWorksService(),
			"alicloud_data_works_service":                               dataSourceAlicloudDataWorksService(),
			"alicloud_mns_service":                                      dataSourceAlicloudMnsService(),
			"alicloud_cloud_storage_gateway_service":                    dataSourceAlicloudCloudStorageGatewayService(),
			"alicloud_vs_service":                                       dataSourceAlicloudVsService(),
			"alicloud_pvtz_service":                                     dataSourceAlicloudPvtzService(),
			"alicloud_cms_service":                                      dataSourceAlicloudCmsService(),
			"alicloud_maxcompute_service":                               dataSourceAlicloudMaxcomputeService(),
			"alicloud_brain_industrial_service":                         dataSourceAlicloudBrainIndustrialService(),
			"alicloud_iot_service":                                      dataSourceAlicloudIotService(),
			"alicloud_ack_service":                                      dataSourceAlicloudAckService(),
			"alicloud_cr_service":                                       dataSourceAlicloudCrService(),
			"alicloud_dcdn_service":                                     dataSourceAlicloudDcdnService(),
			"alicloud_datahub_service":                                  dataSourceAlicloudDatahubService(),
			"alicloud_ons_service":                                      dataSourceAlicloudOnsService(),
			"alicloud_fc_service":                                       dataSourceAlicloudFcService(),
			"alicloud_privatelink_service":                              dataSourceAlicloudPrivateLinkService(),
			"alicloud_edas_applications":                                dataSourceAlicloudEdasApplications(),
			"alicloud_edas_deploy_groups":                               dataSourceAlicloudEdasDeployGroups(),
			"alicloud_edas_clusters":                                    dataSourceAlicloudEdasClusters(),
			"alicloud_resource_manager_folders":                         dataSourceAliCloudResourceManagerFolders(),
			"alicloud_dns_instances":                                    dataSourceAlicloudAlidnsInstances(),
			"alicloud_resource_manager_policies":                        dataSourceAlicloudResourceManagerPolicies(),
			"alicloud_resource_manager_resource_groups":                 dataSourceAlicloudResourceManagerResourceGroups(),
			"alicloud_resource_manager_roles":                           dataSourceAlicloudResourceManagerRoles(),
			"alicloud_resource_manager_policy_versions":                 dataSourceAlicloudResourceManagerPolicyVersions(),
			"alicloud_alidns_domain_groups":                             dataSourceAlicloudAlidnsDomainGroups(),
			"alicloud_kms_key_versions":                                 dataSourceAlicloudKmsKeyVersions(),
			"alicloud_alidns_records":                                   dataSourceAlicloudAlidnsRecords(),
			"alicloud_resource_manager_accounts":                        dataSourceAlicloudResourceManagerAccounts(),
			"alicloud_resource_manager_resource_directories":            dataSourceAlicloudResourceManagerResourceDirectories(),
			"alicloud_resource_manager_handshakes":                      dataSourceAlicloudResourceManagerHandshakes(),
			"alicloud_waf_domains":                                      dataSourceAlicloudWafDomains(),
			"alicloud_kms_secrets":                                      dataSourceAlicloudKmsSecrets(),
			"alicloud_cen_route_maps":                                   dataSourceAlicloudCenRouteMaps(),
			"alicloud_cen_private_zones":                                dataSourceAlicloudCenPrivateZones(),
			"alicloud_dms_enterprise_instances":                         dataSourceAlicloudDmsEnterpriseInstances(),
			"alicloud_cassandra_clusters":                               dataSourceAlicloudCassandraClusters(),
			"alicloud_cassandra_data_centers":                           dataSourceAlicloudCassandraDataCenters(),
			"alicloud_cassandra_zones":                                  dataSourceAlicloudCassandraZones(),
			"alicloud_kms_secret_versions":                              dataSourceAlicloudKmsSecretVersions(),
			"alicloud_waf_instances":                                    dataSourceAlicloudWafInstances(),
			"alicloud_eci_image_caches":                                 dataSourceAlicloudEciImageCaches(),
			"alicloud_dms_enterprise_users":                             dataSourceAlicloudDmsEnterpriseUsers(),
			"alicloud_dms_user_tenants":                                 dataSourceAlicloudDmsUserTenants(),
			"alicloud_ecs_dedicated_hosts":                              dataSourceAlicloudEcsDedicatedHosts(),
			"alicloud_oos_templates":                                    dataSourceAlicloudOosTemplates(),
			"alicloud_oos_executions":                                   dataSourceAlicloudOosExecutions(),
			"alicloud_resource_manager_policy_attachments":              dataSourceAlicloudResourceManagerPolicyAttachments(),
			"alicloud_dcdn_domains":                                     dataSourceAlicloudDcdnDomains(),
			"alicloud_mse_clusters":                                     dataSourceAlicloudMseClusters(),
			"alicloud_actiontrail_trails":                               dataSourceAlicloudActiontrailTrails(),
			"alicloud_actiontrails":                                     dataSourceAlicloudActiontrailTrails(),
			"alicloud_alidns_instances":                                 dataSourceAlicloudAlidnsInstances(),
			"alicloud_alidns_domains":                                   dataSourceAlicloudAlidnsDomains(),
			"alicloud_log_alert_resource":                               dataSourceAlicloudLogAlertResource(),
			"alicloud_log_service":                                      dataSourceAlicloudLogService(),
			"alicloud_cen_instance_attachments":                         dataSourceAlicloudCenInstanceAttachments(),
			"alicloud_cdn_service":                                      dataSourceAlicloudCdnService(),
			"alicloud_cen_vbr_health_checks":                            dataSourceAlicloudCenVbrHealthChecks(),
			"alicloud_config_rules":                                     dataSourceAlicloudConfigRules(),
			"alicloud_config_configuration_recorders":                   dataSourceAlicloudConfigConfigurationRecorders(),
			"alicloud_config_delivery_channels":                         dataSourceAlicloudConfigDeliveryChannels(),
			"alicloud_cms_alarm_contacts":                               dataSourceAlicloudCmsAlarmContacts(),
			"alicloud_kvstore_connections":                              dataSourceAlicloudKvstoreConnections(),
			"alicloud_cms_alarm_contact_groups":                         dataSourceAlicloudCmsAlarmContactGroups(),
			"alicloud_enhanced_nat_available_zones":                     dataSourceAlicloudEnhancedNatAvailableZones(),
			"alicloud_cen_route_services":                               dataSourceAlicloudCenRouteServices(),
			"alicloud_kvstore_accounts":                                 dataSourceAlicloudKvstoreAccounts(),
			"alicloud_cms_group_metric_rules":                           dataSourceAlicloudCmsGroupMetricRules(),
			"alicloud_fnf_flows":                                        dataSourceAlicloudFnfFlows(),
			"alicloud_fnf_schedules":                                    dataSourceAlicloudFnfSchedules(),
			"alicloud_ros_change_sets":                                  dataSourceAlicloudRosChangeSets(),
			"alicloud_ros_stacks":                                       dataSourceAlicloudRosStacks(),
			"alicloud_ros_stack_groups":                                 dataSourceAlicloudRosStackGroups(),
			"alicloud_ros_templates":                                    dataSourceAlicloudRosTemplates(),
			"alicloud_privatelink_vpc_endpoint_services":                dataSourceAlicloudPrivatelinkVpcEndpointServices(),
			"alicloud_privatelink_vpc_endpoints":                        dataSourceAlicloudPrivatelinkVpcEndpoints(),
			"alicloud_privatelink_vpc_endpoint_connections":             dataSourceAlicloudPrivatelinkVpcEndpointConnections(),
			"alicloud_privatelink_vpc_endpoint_service_resources":       dataSourceAlicloudPrivatelinkVpcEndpointServiceResources(),
			"alicloud_privatelink_vpc_endpoint_service_users":           dataSourceAlicloudPrivatelinkVpcEndpointServiceUsers(),
			"alicloud_resource_manager_resource_shares":                 dataSourceAlicloudResourceManagerResourceShares(),
			"alicloud_privatelink_vpc_endpoint_zones":                   dataSourceAlicloudPrivatelinkVpcEndpointZones(),
			"alicloud_ga_accelerators":                                  dataSourceAlicloudGaAccelerators(),
			"alicloud_eci_container_groups":                             dataSourceAlicloudEciContainerGroups(),
			"alicloud_resource_manager_shared_resources":                dataSourceAliCloudResourceManagerSharedResources(),
			"alicloud_resource_manager_shared_targets":                  dataSourceAliCloudResourceManagerSharedTargets(),
			"alicloud_ga_listeners":                                     dataSourceAlicloudGaListeners(),
			"alicloud_tsdb_instances":                                   dataSourceAlicloudTsdbInstances(),
			"alicloud_tsdb_zones":                                       dataSourceAlicloudTsdbZones(),
			"alicloud_ga_bandwidth_packages":                            dataSourceAlicloudGaBandwidthPackages(),
			"alicloud_ga_endpoint_groups":                               dataSourceAliCloudGaEndpointGroups(),
			"alicloud_brain_industrial_pid_organizations":               dataSourceAlicloudBrainIndustrialPidOrganizations(),
			"alicloud_ga_ip_sets":                                       dataSourceAlicloudGaIpSets(),
			"alicloud_ga_forwarding_rules":                              dataSourceAlicloudGaForwardingRules(),
			"alicloud_eipanycast_anycast_eip_addresses":                 dataSourceAlicloudEipanycastAnycastEipAddresses(),
			"alicloud_brain_industrial_pid_projects":                    dataSourceAlicloudBrainIndustrialPidProjects(),
			"alicloud_cms_monitor_groups":                               dataSourceAlicloudCmsMonitorGroups(),
			"alicloud_ram_saml_providers":                               dataSourceAlicloudRamSamlProviders(),
			"alicloud_quotas_quotas":                                    dataSourceAlicloudQuotasQuotas(),
			"alicloud_quotas_application_infos":                         dataSourceAlicloudQuotasQuotaApplications(),
			"alicloud_cms_monitor_group_instanceses":                    dataSourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_cms_monitor_group_instances":                      dataSourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_quotas_quota_alarms":                              dataSourceAlicloudQuotasQuotaAlarms(),
			"alicloud_ecs_commands":                                     dataSourceAlicloudEcsCommands(),
			"alicloud_cloud_storage_gateway_storage_bundles":            dataSourceAlicloudCloudStorageGatewayStorageBundles(),
			"alicloud_ecs_hpc_clusters":                                 dataSourceAlicloudEcsHpcClusters(),
			"alicloud_brain_industrial_pid_loops":                       dataSourceAlicloudBrainIndustrialPidLoops(),
			"alicloud_quotas_quota_applications":                        dataSourceAlicloudQuotasQuotaApplications(),
			"alicloud_ecs_auto_snapshot_policies":                       dataSourceAlicloudEcsAutoSnapshotPolicies(),
			"alicloud_rds_parameter_groups":                             dataSourceAlicloudRdsParameterGroups(),
			"alicloud_rds_collation_time_zones":                         dataSourceAlicloudRdsCollationTimeZones(),
			"alicloud_ecs_launch_templates":                             dataSourceAlicloudEcsLaunchTemplates(),
			"alicloud_resource_manager_control_policies":                dataSourceAlicloudResourceManagerControlPolicies(),
			"alicloud_resource_manager_control_policy_attachments":      dataSourceAlicloudResourceManagerControlPolicyAttachments(),
			"alicloud_instance_keywords":                                dataSourceAlicloudInstanceKeywords(),
			"alicloud_rds_accounts":                                     dataSourceAlicloudRdsAccounts(),
			"alicloud_db_instance_class_infos":                          dataSourceAlicloudDBInstanceClassInfos(),
			"alicloud_rds_cross_regions":                                dataSourceAlicloudRdsCrossRegions(),
			"alicloud_rds_cross_region_backups":                         dataSourceAlicloudRdsCrossRegionBackups(),
			"alicloud_rds_character_set_names":                          dataSourceAlicloudRdsCharacterSetNames(),
			"alicloud_rds_slots":                                        dataSourceAlicloudRdsSlots(),
			"alicloud_rds_class_details":                                dataSourceAlicloudRdsClassDetails(),
			"alicloud_havips":                                           dataSourceAlicloudHavips(),
			"alicloud_ecs_snapshots":                                    dataSourceAlicloudEcsSnapshots(),
			"alicloud_ecs_key_pairs":                                    dataSourceAlicloudEcsKeyPairs(),
			"alicloud_adb_db_clusters":                                  dataSourceAlicloudAdbDbClusters(),
			"alicloud_vpc_flow_logs":                                    dataSourceAlicloudVpcFlowLogs(),
			"alicloud_network_acls":                                     dataSourceAlicloudNetworkAcls(),
			"alicloud_ecs_disks":                                        dataSourceAlicloudEcsDisks(),
			"alicloud_ddoscoo_domain_resources":                         dataSourceAlicloudDdoscooDomainResources(),
			"alicloud_ddoscoo_ports":                                    dataSourceAlicloudDdoscooPorts(),
			"alicloud_slb_load_balancers":                               dataSourceAlicloudSlbLoadBalancers(),
			"alicloud_ecs_network_interfaces":                           dataSourceAlicloudEcsNetworkInterfaces(),
			"alicloud_config_aggregators":                               dataSourceAlicloudConfigAggregators(),
			"alicloud_config_aggregate_config_rules":                    dataSourceAlicloudConfigAggregateConfigRules(),
			"alicloud_config_aggregate_compliance_packs":                dataSourceAlicloudConfigAggregateCompliancePacks(),
			"alicloud_config_compliance_packs":                          dataSourceAlicloudConfigCompliancePacks(),
			"alicloud_eip_addresses":                                    dataSourceAlicloudEipAddresses(),
			"alicloud_direct_mail_receiverses":                          dataSourceAlicloudDirectMailReceiverses(),
			"alicloud_log_projects":                                     dataSourceAlicloudLogProjects(),
			"alicloud_log_stores":                                       dataSourceAlicloudLogStores(),
			"alicloud_event_bridge_service":                             dataSourceAlicloudEventBridgeService(),
			"alicloud_event_bridge_event_buses":                         dataSourceAlicloudEventBridgeEventBuses(),
			"alicloud_amqp_virtual_hosts":                               dataSourceAlicloudAmqpVirtualHosts(),
			"alicloud_amqp_queues":                                      dataSourceAlicloudAmqpQueues(),
			"alicloud_amqp_exchanges":                                   dataSourceAlicloudAmqpExchanges(),
			"alicloud_cassandra_backup_plans":                           dataSourceAlicloudCassandraBackupPlans(),
			"alicloud_cen_transit_router_peer_attachments":              dataSourceAlicloudCenTransitRouterPeerAttachments(),
			"alicloud_amqp_instances":                                   dataSourceAlicloudAmqpInstances(),
			"alicloud_hbr_vaults":                                       dataSourceAlicloudHbrVaults(),
			"alicloud_ssl_certificates_service_certificates":            dataSourceAliCloudSslCertificatesServiceCertificates(),
			"alicloud_arms_alert_contacts":                              dataSourceAlicloudArmsAlertContacts(),
			"alicloud_event_bridge_rules":                               dataSourceAlicloudEventBridgeRules(),
			"alicloud_cloud_firewall_control_policies":                  dataSourceAliCloudCloudFirewallControlPolicies(),
			"alicloud_sae_namespaces":                                   dataSourceAlicloudSaeNamespaces(),
			"alicloud_sae_config_maps":                                  dataSourceAlicloudSaeConfigMaps(),
			"alicloud_alb_security_policies":                            dataSourceAlicloudAlbSecurityPolicies(),
			"alicloud_alb_system_security_policies":                     dataSourceAlicloudAlbSystemSecurityPolicies(),
			"alicloud_event_bridge_event_sources":                       dataSourceAlicloudEventBridgeEventSources(),
			"alicloud_ecd_policy_groups":                                dataSourceAlicloudEcdPolicyGroups(),
			"alicloud_ecp_key_pairs":                                    dataSourceAlicloudEcpKeyPairs(),
			"alicloud_hbr_ecs_backup_plans":                             dataSourceAlicloudHbrEcsBackupPlans(),
			"alicloud_hbr_nas_backup_plans":                             dataSourceAlicloudHbrNasBackupPlans(),
			"alicloud_hbr_oss_backup_plans":                             dataSourceAlicloudHbrOssBackupPlans(),
			"alicloud_scdn_domains":                                     dataSourceAlicloudScdnDomains(),
			"alicloud_alb_server_groups":                                dataSourceAlicloudAlbServerGroups(),
			"alicloud_data_works_folders":                               dataSourceAlicloudDataWorksFolders(),
			"alicloud_arms_alert_contact_groups":                        dataSourceAlicloudArmsAlertContactGroups(),
			"alicloud_express_connect_access_points":                    dataSourceAlicloudExpressConnectAccessPoints(),
			"alicloud_cloud_storage_gateway_gateways":                   dataSourceAlicloudCloudStorageGatewayGateways(),
			"alicloud_lindorm_instances":                                dataSourceAlicloudLindormInstances(),
			"alicloud_express_connect_physical_connection_service":      dataSourceAlicloudExpressConnectPhysicalConnectionService(),
			"alicloud_cddc_dedicated_host_groups":                       dataSourceAlicloudCddcDedicatedHostGroups(),
			"alicloud_hbr_ecs_backup_clients":                           dataSourceAlicloudHbrEcsBackupClients(),
			"alicloud_msc_sub_contacts":                                 dataSourceAlicloudMscSubContacts(),
			"alicloud_express_connect_physical_connections":             dataSourceAlicloudExpressConnectPhysicalConnections(),
			"alicloud_alb_load_balancers":                               dataSourceAlicloudAlbLoadBalancers(),
			"alicloud_alb_zones":                                        dataSourceAlicloudAlbZones(),
			"alicloud_sddp_rules":                                       dataSourceAlicloudSddpRules(),
			"alicloud_bastionhost_user_groups":                          dataSourceAlicloudBastionhostUserGroups(),
			"alicloud_security_center_groups":                           dataSourceAlicloudSecurityCenterGroups(),
			"alicloud_alb_acls":                                         dataSourceAlicloudAlbAcls(),
			"alicloud_hbr_snapshots":                                    dataSourceAlicloudHbrSnapshots(),
			"alicloud_bastionhost_users":                                dataSourceAlicloudBastionhostUsers(),
			"alicloud_dfs_access_groups":                                dataSourceAlicloudDfsAccessGroups(),
			"alicloud_ehpc_job_templates":                               dataSourceAlicloudEhpcJobTemplates(),
			"alicloud_sddp_configs":                                     dataSourceAlicloudSddpConfigs(),
			"alicloud_hbr_restore_jobs":                                 dataSourceAlicloudHbrRestoreJobs(),
			"alicloud_alb_listeners":                                    dataSourceAlicloudAlbListeners(),
			"alicloud_ens_key_pairs":                                    dataSourceAlicloudEnsKeyPairs(),
			"alicloud_sae_applications":                                 dataSourceAlicloudSaeApplications(),
			"alicloud_alb_rules":                                        dataSourceAliCloudAlbRules(),
			"alicloud_cms_metric_rule_templates":                        dataSourceAlicloudCmsMetricRuleTemplates(),
			"alicloud_iot_device_groups":                                dataSourceAlicloudIotDeviceGroups(),
			"alicloud_express_connect_virtual_border_routers":           dataSourceAlicloudExpressConnectVirtualBorderRouters(),
			"alicloud_imm_projects":                                     dataSourceAlicloudImmProjects(),
			"alicloud_click_house_db_clusters":                          dataSourceAlicloudClickHouseDbClusters(),
			"alicloud_direct_mail_domains":                              dataSourceAliCloudDirectMailDomains(),
			"alicloud_bastionhost_host_groups":                          dataSourceAlicloudBastionhostHostGroups(),
			"alicloud_vpc_dhcp_options_sets":                            dataSourceAlicloudVpcDhcpOptionsSets(),
			"alicloud_alb_health_check_templates":                       dataSourceAlicloudAlbHealthCheckTemplates(),
			"alicloud_cdn_real_time_log_deliveries":                     dataSourceAlicloudCdnRealTimeLogDeliveries(),
			"alicloud_click_house_accounts":                             dataSourceAlicloudClickHouseAccounts(),
			"alicloud_selectdb_db_clusters":                             dataSourceAlicloudSelectDBDbClusters(),
			"alicloud_selectdb_db_instances":                            dataSourceAlicloudSelectDBDbInstances(),
			"alicloud_direct_mail_mail_addresses":                       dataSourceAlicloudDirectMailMailAddresses(),
			"alicloud_database_gateway_gateways":                        dataSourceAlicloudDatabaseGatewayGateways(),
			"alicloud_bastionhost_hosts":                                dataSourceAlicloudBastionhostHosts(),
			"alicloud_amqp_bindings":                                    dataSourceAlicloudAmqpBindings(),
			"alicloud_slb_tls_cipher_policies":                          dataSourceAlicloudSlbTlsCipherPolicies(),
			"alicloud_cloud_sso_directories":                            dataSourceAlicloudCloudSsoDirectories(),
			"alicloud_bastionhost_host_accounts":                        dataSourceAlicloudBastionhostHostAccounts(),
			"alicloud_waf_certificates":                                 dataSourceAlicloudWafCertificates(),
			"alicloud_simple_application_server_instances":              dataSourceAlicloudSimpleApplicationServerInstances(),
			"alicloud_simple_application_server_plans":                  dataSourceAlicloudSimpleApplicationServerPlans(),
			"alicloud_simple_application_server_images":                 dataSourceAlicloudSimpleApplicationServerImages(),
			"alicloud_video_surveillance_system_groups":                 dataSourceAlicloudVideoSurveillanceSystemGroups(),
			"alicloud_msc_sub_subscriptions":                            dataSourceAlicloudMscSubSubscriptions(),
			"alicloud_sddp_instances":                                   dataSourceAlicloudSddpInstances(),
			"alicloud_vpc_nat_ip_cidrs":                                 dataSourceAlicloudVpcNatIpCidrs(),
			"alicloud_vpc_nat_ips":                                      dataSourceAlicloudVpcNatIps(),
			"alicloud_quick_bi_users":                                   dataSourceAlicloudQuickBiUsers(),
			"alicloud_vod_domains":                                      dataSourceAlicloudVodDomains(),
			"alicloud_arms_dispatch_rules":                              dataSourceAlicloudArmsDispatchRules(),
			"alicloud_open_search_app_groups":                           dataSourceAlicloudOpenSearchAppGroups(),
			"alicloud_graph_database_db_instances":                      dataSourceAlicloudGraphDatabaseDbInstances(),
			"alicloud_arms_prometheus_alert_rules":                      dataSourceAlicloudArmsPrometheusAlertRules(),
			"alicloud_dbfs_instances":                                   dataSourceAlicloudDbfsInstances(),
			"alicloud_rdc_organizations":                                dataSourceAlicloudRdcOrganizations(),
			"alicloud_eais_instances":                                   dataSourceAlicloudEaisInstances(),
			"alicloud_sae_ingresses":                                    dataSourceAlicloudSaeIngresses(),
			"alicloud_cloudauth_face_configs":                           dataSourceAlicloudCloudauthFaceConfigs(),
			"alicloud_imp_app_templates":                                dataSourceAlicloudImpAppTemplates(),
			"alicloud_mhub_products":                                    dataSourceAlicloudMhubProducts(),
			"alicloud_cloud_sso_scim_server_credentials":                dataSourceAlicloudCloudSsoScimServerCredentials(),
			"alicloud_dts_subscription_jobs":                            dataSourceAlicloudDtsSubscriptionJobs(),
			"alicloud_service_mesh_service_meshes":                      dataSourceAlicloudServiceMeshServiceMeshes(),
			"alicloud_service_mesh_versions":                            dataSourceAlicloudServiceMeshVersions(),
			"alicloud_mhub_apps":                                        dataSourceAlicloudMhubApps(),
			"alicloud_cloud_sso_groups":                                 dataSourceAlicloudCloudSsoGroups(),
			"alicloud_hbr_backup_jobs":                                  dataSourceAlicloudHbrBackupJobs(),
			"alicloud_click_house_regions":                              dataSourceAlicloudClickHouseRegions(),
			"alicloud_dts_synchronization_jobs":                         dataSourceAlicloudDtsSynchronizationJobs(),
			"alicloud_cloud_firewall_instances":                         dataSourceAlicloudCloudFirewallInstances(),
			"alicloud_cr_endpoint_acl_policies":                         dataSourceAlicloudCrEndpointAclPolicies(),
			"alicloud_cr_endpoint_acl_service":                          dataSourceAlicloudCrEndpointAclService(),
			"alicloud_actiontrail_history_delivery_jobs":                dataSourceAlicloudActiontrailHistoryDeliveryJobs(),
			"alicloud_sae_instance_specifications":                      dataSourceAlicloudSaeInstanceSpecifications(),
			"alicloud_cen_transit_router_service":                       dataSourceAlicloudCenTransitRouterService(),
			"alicloud_ecs_deployment_sets":                              dataSourceAlicloudEcsDeploymentSets(),
			"alicloud_cloud_sso_users":                                  dataSourceAlicloudCloudSsoUsers(),
			"alicloud_cloud_sso_access_configurations":                  dataSourceAlicloudCloudSsoAccessConfigurations(),
			"alicloud_dfs_file_systems":                                 dataSourceAlicloudDfsFileSystems(),
			"alicloud_dfs_zones":                                        dataSourceAlicloudDfsZones(),
			"alicloud_vpc_traffic_mirror_filters":                       dataSourceAlicloudVpcTrafficMirrorFilters(),
			"alicloud_dfs_access_rules":                                 dataSourceAlicloudDfsAccessRules(),
			"alicloud_nas_zones":                                        dataSourceAlicloudNasZones(),
			"alicloud_dfs_mount_points":                                 dataSourceAlicloudDfsMountPoints(),
			"alicloud_vpc_traffic_mirror_filter_egress_rules":           dataSourceAlicloudVpcTrafficMirrorFilterEgressRules(),
			"alicloud_ecd_simple_office_sites":                          dataSourceAlicloudEcdSimpleOfficeSites(),
			"alicloud_vpc_traffic_mirror_filter_ingress_rules":          dataSourceAlicloudVpcTrafficMirrorFilterIngressRules(),
			"alicloud_ecd_nas_file_systems":                             dataSourceAlicloudEcdNasFileSystems(),
			"alicloud_vpc_traffic_mirror_service":                       dataSourceAlicloudVpcTrafficMirrorService(),
			"alicloud_msc_sub_webhooks":                                 dataSourceAlicloudMscSubWebhooks(),
			"alicloud_ecd_users":                                        dataSourceAlicloudEcdUsers(),
			"alicloud_vpc_traffic_mirror_sessions":                      dataSourceAlicloudVpcTrafficMirrorSessions(),
			"alicloud_gpdb_accounts":                                    dataSourceAlicloudGpdbAccounts(),
			"alicloud_vpc_ipv6_gateways":                                dataSourceAlicloudVpcIpv6Gateways(),
			"alicloud_vpc_ipv6_egress_rules":                            dataSourceAlicloudVpcIpv6EgressRules(),
			"alicloud_vpc_ipv6_addresses":                               dataSourceAlicloudVpcIpv6Addresses(),
			"alicloud_hbr_server_backup_plans":                          dataSourceAlicloudHbrServerBackupPlans(),
			"alicloud_cms_dynamic_tag_groups":                           dataSourceAlicloudCmsDynamicTagGroups(),
			"alicloud_ecd_network_packages":                             dataSourceAlicloudEcdNetworkPackages(),
			"alicloud_cloud_storage_gateway_gateway_smb_users":          dataSourceAlicloudCloudStorageGatewayGatewaySmbUsers(),
			"alicloud_vpc_ipv6_internet_bandwidths":                     dataSourceAlicloudVpcIpv6InternetBandwidths(),
			"alicloud_simple_application_server_firewall_rules":         dataSourceAlicloudSimpleApplicationServerFirewallRules(),
			"alicloud_pvtz_endpoints":                                   dataSourceAlicloudPvtzEndpoints(),
			"alicloud_pvtz_resolver_zones":                              dataSourceAlicloudPvtzResolverZones(),
			"alicloud_pvtz_rules":                                       dataSourceAlicloudPvtzRules(),
			"alicloud_ecd_bundles":                                      dataSourceAlicloudEcdBundles(),
			"alicloud_simple_application_server_disks":                  dataSourceAlicloudSimpleApplicationServerDisks(),
			"alicloud_simple_application_server_snapshots":              dataSourceAlicloudSimpleApplicationServerSnapshots(),
			"alicloud_simple_application_server_custom_images":          dataSourceAlicloudSimpleApplicationServerCustomImages(),
			"alicloud_cloud_storage_gateway_stocks":                     dataSourceAlicloudCloudStorageGatewayStocks(),
			"alicloud_cloud_storage_gateway_gateway_cache_disks":        dataSourceAlicloudCloudStorageGatewayGatewayCacheDisks(),
			"alicloud_cloud_storage_gateway_gateway_block_volumes":      dataSourceAlicloudCloudStorageGatewayGatewayBlockVolumes(),
			"alicloud_direct_mail_tags":                                 dataSourceAlicloudDirectMailTags(),
			"alicloud_cloud_storage_gateway_gateway_file_shares":        dataSourceAlicloudCloudStorageGatewayGatewayFileShares(),
			"alicloud_ecd_desktops":                                     dataSourceAlicloudEcdDesktops(),
			"alicloud_cloud_storage_gateway_express_syncs":              dataSourceAlicloudCloudStorageGatewayExpressSyncs(),
			"alicloud_oos_applications":                                 dataSourceAlicloudOosApplications(),
			"alicloud_eci_virtual_nodes":                                dataSourceAlicloudEciVirtualNodes(),
			"alicloud_eci_zones":                                        dataSourceAlicloudEciZones(),
			"alicloud_ros_stack_instances":                              dataSourceAlicloudRosStackInstances(),
			"alicloud_ros_regions":                                      dataSourceAlicloudRosRegions(),
			"alicloud_ecs_dedicated_host_clusters":                      dataSourceAlicloudEcsDedicatedHostClusters(),
			"alicloud_oos_application_groups":                           dataSourceAlicloudOosApplicationGroups(),
			"alicloud_dts_consumer_channels":                            dataSourceAlicloudDtsConsumerChannels(),
			"alicloud_emr_clusters":                                     dataSourceAlicloudEmrClusters(),
			"alicloud_emrv2_clusters":                                   dataSourceAlicloudEmrV2Clusters(),
			"alicloud_ecd_images":                                       dataSourceAlicloudEcdImages(),
			"alicloud_oos_patch_baselines":                              dataSourceAlicloudOosPatchBaselines(),
			"alicloud_ecd_commands":                                     dataSourceAlicloudEcdCommands(),
			"alicloud_cddc_zones":                                       dataSourceAlicloudCddcZones(),
			"alicloud_cddc_host_ecs_level_infos":                        dataSourceAlicloudCddcHostEcsLevelInfos(),
			"alicloud_cddc_dedicated_hosts":                             dataSourceAlicloudCddcDedicatedHosts(),
			"alicloud_oos_parameters":                                   dataSourceAlicloudOosParameters(),
			"alicloud_oos_state_configurations":                         dataSourceAlicloudOosStateConfigurations(),
			"alicloud_oos_secret_parameters":                            dataSourceAliCloudOosSecretParameters(),
			"alicloud_click_house_backup_policies":                      dataSourceAlicloudClickHouseBackupPolicies(),
			"alicloud_cloud_sso_service":                                dataSourceAlicloudCloudSsoService(),
			"alicloud_mongodb_audit_policies":                           dataSourceAlicloudMongodbAuditPolicies(),
			"alicloud_mongodb_accounts":                                 dataSourceAlicloudMongodbAccounts(),
			"alicloud_mongodb_serverless_instances":                     dataSourceAlicloudMongodbServerlessInstances(),
			"alicloud_cddc_dedicated_host_accounts":                     dataSourceAlicloudCddcDedicatedHostAccounts(),
			"alicloud_cr_chart_namespaces":                              dataSourceAlicloudCrChartNamespaces(),
			"alicloud_fnf_executions":                                   dataSourceAlicloudFnFExecutions(),
			"alicloud_cr_chart_repositories":                            dataSourceAlicloudCrChartRepositories(),
			"alicloud_mongodb_sharding_network_public_addresses":        dataSourceAlicloudMongodbShardingNetworkPublicAddresses(),
			"alicloud_ga_acls":                                          dataSourceAlicloudGaAcls(),
			"alicloud_ga_additional_certificates":                       dataSourceAlicloudGaAdditionalCertificates(),
			"alicloud_alidns_custom_lines":                              dataSourceAlicloudAlidnsCustomLines(),
			"alicloud_ros_template_scratches":                           dataSourceAlicloudRosTemplateScratches(),
			"alicloud_alidns_gtm_instances":                             dataSourceAlicloudAlidnsGtmInstances(),
			"alicloud_vpc_bgp_groups":                                   dataSourceAlicloudVpcBgpGroups(),
			"alicloud_nas_snapshots":                                    dataSourceAlicloudNasSnapshots(),
			"alicloud_hbr_replication_vault_regions":                    dataSourceAlicloudHbrReplicationVaultRegions(),
			"alicloud_alidns_address_pools":                             dataSourceAlicloudAlidnsAddressPools(),
			"alicloud_ecs_prefix_lists":                                 dataSourceAlicloudEcsPrefixLists(),
			"alicloud_alidns_access_strategies":                         dataSourceAlicloudAlidnsAccessStrategies(),
			"alicloud_vpc_bgp_peers":                                    dataSourceAlicloudVpcBgpPeers(),
			"alicloud_nas_filesets":                                     dataSourceAlicloudNasFilesets(),
			"alicloud_cdn_ip_info":                                      dataSourceAlicloudCdnIpInfo(),
			"alicloud_nas_auto_snapshot_policies":                       dataSourceAlicloudNasAutoSnapshotPolicies(),
			"alicloud_nas_lifecycle_policies":                           dataSourceAlicloudNasLifecyclePolicies(),
			"alicloud_vpc_bgp_networks":                                 dataSourceAlicloudVpcBgpNetworks(),
			"alicloud_nas_data_flows":                                   dataSourceAlicloudNasDataFlows(),
			"alicloud_ecs_storage_capacity_units":                       dataSourceAlicloudEcsStorageCapacityUnits(),
			"alicloud_dbfs_snapshots":                                   dataSourceAlicloudDbfsSnapshots(),
			"alicloud_msc_sub_contact_verification_message":             dataSourceAlicloudMscSubContactVerificationMessage(),
			"alicloud_dts_migration_jobs":                               dataSourceAlicloudDtsMigrationJobs(),
			"alicloud_mse_gateways":                                     dataSourceAlicloudMseGateways(),
			"alicloud_mongodb_sharding_network_private_addresses":       dataSourceAlicloudMongodbShardingNetworkPrivateAddresses(),
			"alicloud_ecp_instances":                                    dataSourceAlicloudEcpInstances(),
			"alicloud_ecp_zones":                                        dataSourceAlicloudEcpZones(),
			"alicloud_ecp_instance_types":                               dataSourceAlicloudEcpInstanceTypes(),
			"alicloud_dcdn_ipa_domains":                                 dataSourceAlicloudDcdnIpaDomains(),
			"alicloud_sddp_data_limits":                                 dataSourceAlicloudSddpDataLimits(),
			"alicloud_ecs_image_components":                             dataSourceAlicloudEcsImageComponents(),
			"alicloud_sae_application_scaling_rules":                    dataSourceAlicloudSaeApplicationScalingRules(),
			"alicloud_sae_grey_tag_routes":                              dataSourceAlicloudSaeGreyTagRoutes(),
			"alicloud_ecs_snapshot_groups":                              dataSourceAlicloudEcsSnapshotGroups(),
			"alicloud_vpn_ipsec_servers":                                dataSourceAlicloudVpnIpsecServers(),
			"alicloud_cr_chains":                                        dataSourceAlicloudCrChains(),
			"alicloud_vpn_pbr_route_entries":                            dataSourceAlicloudVpnPbrRouteEntries(),
			"alicloud_mse_znodes":                                       dataSourceAlicloudMseZnodes(),
			"alicloud_cen_transit_router_available_resources":           dataSourceAliCloudCenTransitRouterAvailableResources(),
			"alicloud_ecs_image_pipelines":                              dataSourceAlicloudEcsImagePipelines(),
			"alicloud_hbr_ots_backup_plans":                             dataSourceAlicloudHbrOtsBackupPlans(),
			"alicloud_hbr_ots_snapshots":                                dataSourceAlicloudHbrOtsSnapshots(),
			"alicloud_bastionhost_host_share_keys":                      dataSourceAlicloudBastionhostHostShareKeys(),
			"alicloud_ecs_network_interface_permissions":                dataSourceAlicloudEcsNetworkInterfacePermissions(),
			"alicloud_mse_engine_namespaces":                            dataSourceAlicloudMseEngineNamespaces(),
			"alicloud_mse_nacos_configs":                                dataSourceAlicloudMseNacosConfigs(),
			"alicloud_ga_accelerator_spare_ip_attachments":              dataSourceAlicloudGaAcceleratorSpareIpAttachments(),
			"alicloud_smartag_flow_logs":                                dataSourceAlicloudSmartagFlowLogs(),
			"alicloud_ecs_invocations":                                  dataSourceAlicloudEcsInvocations(),
			"alicloud_ecd_snapshots":                                    dataSourceAlicloudEcdSnapshots(),
			"alicloud_tag_meta_tags":                                    dataSourceAlicloudTagMetaTags(),
			"alicloud_ecd_desktop_types":                                dataSourceAlicloudEcdDesktopTypes(),
			"alicloud_config_deliveries":                                dataSourceAlicloudConfigDeliveries(),
			"alicloud_cms_namespaces":                                   dataSourceAlicloudCmsNamespaces(),
			"alicloud_cms_sls_groups":                                   dataSourceAlicloudCmsSlsGroups(),
			"alicloud_config_aggregate_deliveries":                      dataSourceAlicloudConfigAggregateDeliveries(),
			"alicloud_edas_namespaces":                                  dataSourceAlicloudEdasNamespaces(),
			"alicloud_cdn_blocked_regions":                              dataSourceAlicloudCdnBlockedRegions(),
			"alicloud_schedulerx_namespaces":                            dataSourceAlicloudSchedulerxNamespaces(),
			"alicloud_ehpc_clusters":                                    dataSourceAlicloudEhpcClusters(),
			"alicloud_cen_traffic_marking_policies":                     dataSourceAlicloudCenTrafficMarkingPolicies(),
			"alicloud_ecd_ram_directories":                              dataSourceAlicloudEcdRamDirectories(),
			"alicloud_ecd_zones":                                        dataSourceAlicloudEcdZones(),
			"alicloud_ecd_ad_connector_directories":                     dataSourceAlicloudEcdAdConnectorDirectories(),
			"alicloud_ecd_custom_properties":                            dataSourceAlicloudEcdCustomProperties(),
			"alicloud_ecd_ad_connector_office_sites":                    dataSourceAlicloudEcdAdConnectorOfficeSites(),
			"alicloud_ecs_activations":                                  dataSourceAlicloudEcsActivations(),
			"alicloud_cms_hybrid_monitor_datas":                         dataSourceAlicloudCmsHybridMonitorDatas(),
			"alicloud_cloud_firewall_address_books":                     dataSourceAliCloudCloudFirewallAddressBooks(),
			"alicloud_hbr_hana_instances":                               dataSourceAlicloudHbrHanaInstances(),
			"alicloud_cms_hybrid_monitor_sls_tasks":                     dataSourceAlicloudCmsHybridMonitorSlsTasks(),
			"alicloud_hbr_hana_backup_plans":                            dataSourceAlicloudHbrHanaBackupPlans(),
			"alicloud_cms_hybrid_monitor_fc_tasks":                      dataSourceAlicloudCmsHybridMonitorFcTasks(),
			"alicloud_ddosbgp_ips":                                      dataSourceAlicloudDdosbgpIps(),
			"alicloud_vpn_gateway_vpn_attachments":                      dataSourceAlicloudVpnGatewayVpnAttachments(),
			"alicloud_resource_manager_delegated_administrators":        dataSourceAlicloudResourceManagerDelegatedAdministrators(),
			"alicloud_polardb_global_database_networks":                 dataSourceAlicloudPolarDBGlobalDatabaseNetworks(),
			"alicloud_vpc_ipv4_gateways":                                dataSourceAlicloudVpcIpv4Gateways(),
			"alicloud_api_gateway_backends":                             dataSourceAlicloudApiGatewayBackends(),
			"alicloud_vpc_prefix_lists":                                 dataSourceAlicloudVpcPrefixLists(),
			"alicloud_cms_event_rules":                                  dataSourceAlicloudCmsEventRules(),
			"alicloud_cen_transit_router_vpn_attachments":               dataSourceAlicloudCenTransitRouterVpnAttachments(),
			"alicloud_polardb_parameter_groups":                         dataSourceAlicloudPolarDBParameterGroups(),
			"alicloud_vpn_gateway_vco_routes":                           dataSourceAlicloudVpnGatewayVcoRoutes(),
			"alicloud_dcdn_waf_policies":                                dataSourceAlicloudDcdnWafPolicies(),
			"alicloud_hbr_service":                                      dataSourceAlicloudHbrService(),
			"alicloud_api_gateway_log_configs":                          dataSourceAlicloudApiGatewayLogConfigs(),
			"alicloud_dbs_backup_plans":                                 dataSourceAlicloudDbsBackupPlans(),
			"alicloud_dcdn_waf_domains":                                 dataSourceAlicloudDcdnWafDomains(),
			"alicloud_vpc_public_ip_address_pools":                      dataSourceAlicloudVpcPublicIpAddressPools(),
			"alicloud_nlb_server_groups":                                dataSourceAlicloudNlbServerGroups(),
			"alicloud_vpc_peer_connections":                             dataSourceAlicloudVpcPeerConnections(),
			"alicloud_ebs_regions":                                      dataSourceAlicloudEbsRegions(),
			"alicloud_ebs_disk_replica_groups":                          dataSourceAlicloudEbsDiskReplicaGroups(),
			"alicloud_nlb_security_policies":                            dataSourceAlicloudNlbSecurityPolicies(),
			"alicloud_api_gateway_models":                               dataSourceAlicloudApiGatewayModels(),
			"alicloud_resource_manager_account_deletion_check_task":     dataSourceAlicloudResourceManagerAccountDeletionCheckTask(),
			"alicloud_cs_cluster_credential":                            dataSourceAlicloudCSClusterCredential(),
			"alicloud_api_gateway_plugins":                              dataSourceAlicloudApiGatewayPlugins(),
			"alicloud_message_service_queues":                           dataSourceAlicloudMessageServiceQueues(),
			"alicloud_message_service_topics":                           dataSourceAlicloudMessageServiceTopics(),
			"alicloud_message_service_subscriptions":                    dataSourceAlicloudMessageServiceSubscriptions(),
			"alicloud_cen_transit_router_prefix_list_associations":      dataSourceAlicloudCenTransitRouterPrefixListAssociations(),
			"alicloud_dms_enterprise_proxies":                           dataSourceAlicloudDmsEnterpriseProxies(),
			"alicloud_vpc_public_ip_address_pool_cidr_blocks":           dataSourceAlicloudVpcPublicIpAddressPoolCidrBlocks(),
			"alicloud_gpdb_db_instance_plans":                           dataSourceAlicloudGpdbDbInstancePlans(),
			"alicloud_adb_db_cluster_lake_versions":                     dataSourceAlicloudAdbDbClusterLakeVersions(),
			"alicloud_nlb_load_balancers":                               dataSourceAlicloudNlbLoadBalancers(),
			"alicloud_nlb_zones":                                        dataSourceAlicloudNlbZones(),
			"alicloud_service_mesh_extension_providers":                 dataSourceAlicloudServiceMeshExtensionProviders(),
			"alicloud_nlb_listeners":                                    dataSourceAlicloudNlbListeners(),
			"alicloud_nlb_server_group_server_attachments":              dataSourceAlicloudNlbServerGroupServerAttachments(),
			"alicloud_bp_studio_applications":                           dataSourceAlicloudBpStudioApplications(),
			"alicloud_cloud_sso_access_assignments":                     dataSourceAlicloudCloudSsoAccessAssignments(),
			"alicloud_cen_transit_router_cidrs":                         dataSourceAlicloudCenTransitRouterCidrs(),
			"alicloud_ga_basic_accelerators":                            dataSourceAlicloudGaBasicAccelerators(),
			"alicloud_cms_metric_rule_black_lists":                      dataSourceAlicloudCmsMetricRuleBlackLists(),
			"alicloud_cloud_firewall_vpc_firewall_cens":                 dataSourceAlicloudCloudFirewallVpcFirewallCens(),
			"alicloud_cloud_firewall_vpc_firewalls":                     dataSourceAlicloudCloudFirewallVpcFirewalls(),
			"alicloud_cloud_firewall_instance_members":                  dataSourceAlicloudCloudFirewallInstanceMembers(),
			"alicloud_ga_basic_accelerate_ips":                          dataSourceAlicloudGaBasicAccelerateIps(),
			"alicloud_ga_basic_endpoints":                               dataSourceAlicloudGaBasicEndpoints(),
			"alicloud_cloud_firewall_vpc_firewall_control_policies":     dataSourceAlicloudCloudFirewallVpcFirewallControlPolicies(),
			"alicloud_ga_basic_accelerate_ip_endpoint_relations":        dataSourceAlicloudGaBasicAccelerateIpEndpointRelations(),
			"alicloud_threat_detection_web_lock_configs":                dataSourceAlicloudThreatDetectionWebLockConfigs(),
			"alicloud_threat_detection_backup_policies":                 dataSourceAlicloudThreatDetectionBackupPolicies(),
			"alicloud_dms_enterprise_proxy_accesses":                    dataSourceAlicloudDmsEnterpriseProxyAccesses(),
			"alicloud_threat_detection_vul_whitelists":                  dataSourceAlicloudThreatDetectionVulWhitelists(),
			"alicloud_dms_enterprise_logic_databases":                   dataSourceAlicloudDmsEnterpriseLogicDatabases(),
			"alicloud_dms_enterprise_databases":                         dataSourceAlicloudDmsEnterpriseDatabases(),
			"alicloud_amqp_static_accounts":                             dataSourceAlicloudAmqpStaticAccounts(),
			"alicloud_adb_resource_groups":                              dataSourceAlicloudAdbResourceGroups(),
			"alicloud_alb_ascripts":                                     dataSourceAlicloudAlbAscripts(),
			"alicloud_threat_detection_honeypot_nodes":                  dataSourceAlicloudThreatDetectionHoneypotNodes(),
			"alicloud_cen_transit_router_multicast_domains":             dataSourceAlicloudCenTransitRouterMulticastDomains(),
			"alicloud_cen_inter_region_traffic_qos_policies":            dataSourceAlicloudCenInterRegionTrafficQosPolicies(),
			"alicloud_threat_detection_baseline_strategies":             dataSourceAlicloudThreatDetectionBaselineStrategies(),
			"alicloud_threat_detection_assets":                          dataSourceAlicloudThreatDetectionAssets(),
			"alicloud_threat_detection_log_shipper":                     dataSourceAlicloudThreatDetectionLogShipper(),
			"alicloud_threat_detection_anti_brute_force_rules":          dataSourceAlicloudThreatDetectionAntiBruteForceRules(),
			"alicloud_threat_detection_honeypot_images":                 dataSourceAlicloudThreatDetectionHoneypotImages(),
			"alicloud_threat_detection_honey_pots":                      dataSourceAlicloudThreatDetectionHoneyPots(),
			"alicloud_threat_detection_honeypot_probes":                 dataSourceAlicloudThreatDetectionHoneypotProbes(),
			"alicloud_ecs_capacity_reservations":                        dataSourceAlicloudEcsCapacityReservations(),
			"alicloud_cen_inter_region_traffic_qos_queues":              dataSourceAlicloudCenInterRegionTrafficQosQueues(),
			"alicloud_cen_transit_router_multicast_domain_peer_members": dataSourceAlicloudCenTransitRouterMulticastDomainPeerMembers(),
			"alicloud_cen_transit_router_multicast_domain_members":      dataSourceAlicloudCenTransitRouterMulticastDomainMembers(),
			"alicloud_cen_child_instance_route_entry_to_attachments":    dataSourceAlicloudCenChildInstanceRouteEntryToAttachments(),
			"alicloud_cen_transit_router_multicast_domain_associations": dataSourceAlicloudCenTransitRouterMulticastDomainAssociations(),
			"alicloud_threat_detection_honeypot_presets":                dataSourceAlicloudThreatDetectionHoneypotPresets(),
			"alicloud_cen_transit_router_multicast_domain_sources":      dataSourceAlicloudCenTransitRouterMulticastDomainSources(),
			"alicloud_bss_open_api_products":                            dataSourceAlicloudBssOpenApiProducts(),
			"alicloud_bss_open_api_pricing_modules":                     dataSourceAlicloudBssOpenApiPricingModules(),
			"alicloud_service_catalog_provisioned_products":             dataSourceAlicloudServiceCatalogProvisionedProducts(),
			"alicloud_service_catalog_product_as_end_users":             dataSourceAlicloudServiceCatalogProductAsEndUsers(),
			"alicloud_service_catalog_product_versions":                 dataSourceAlicloudServiceCatalogProductVersions(),
			"alicloud_service_catalog_launch_options":                   dataSourceAlicloudServiceCatalogLaunchOptions(),
			"alicloud_maxcompute_projects":                              dataSourceAliCloudMaxComputeProjects(),
			"alicloud_ebs_dedicated_block_storage_clusters":             dataSourceAlicloudEbsDedicatedBlockStorageClusters(),
			"alicloud_ecs_elasticity_assurances":                        dataSourceAlicloudEcsElasticityAssurances(),
			"alicloud_express_connect_grant_rule_to_cens":               dataSourceAlicloudExpressConnectGrantRuleToCens(),
			"alicloud_express_connect_virtual_physical_connections":     dataSourceAlicloudExpressConnectVirtualPhysicalConnections(),
			"alicloud_express_connect_vbr_pconn_associations":           dataSourceAlicloudExpressConnectVbrPconnAssociations(),
			"alicloud_ebs_disk_replica_pairs":                           dataSourceAlicloudEbsDiskReplicaPairs(),
			"alicloud_ga_domains":                                       dataSourceAlicloudGaDomains(),
			"alicloud_ga_custom_routing_endpoint_groups":                dataSourceAlicloudGaCustomRoutingEndpointGroups(),
			"alicloud_ga_custom_routing_endpoint_group_destinations":    dataSourceAlicloudGaCustomRoutingEndpointGroupDestinations(),
			"alicloud_ga_custom_routing_endpoints":                      dataSourceAlicloudGaCustomRoutingEndpoints(),
			"alicloud_ga_custom_routing_endpoint_traffic_policies":      dataSourceAliCloudGaCustomRoutingEndpointTrafficPolicies(),
			"alicloud_ga_custom_routing_port_mappings":                  dataSourceAlicloudGaCustomRoutingPortMappings(),
			"alicloud_service_catalog_end_user_products":                dataSourceAlicloudServiceCatalogEndUserProducts(),
			"alicloud_dcdn_kv_account":                                  dataSourceAlicloudDcdnKvAccount(),
			"alicloud_hbr_hana_backup_clients":                          dataSourceAlicloudHbrHanaBackupClients(),
			"alicloud_dts_instances":                                    dataSourceAlicloudDtsInstances(),
			"alicloud_threat_detection_instances":                       dataSourceAlicloudThreatDetectionInstances(),
			"alicloud_cr_vpc_endpoint_linked_vpcs":                      dataSourceAlicloudCrVpcEndpointLinkedVpcs(),
			"alicloud_express_connect_router_interfaces":                dataSourceAlicloudExpressConnectRouterInterfaces(),
			"alicloud_wafv3_instances":                                  dataSourceAlicloudWafv3Instances(),
			"alicloud_wafv3_domains":                                    dataSourceAliCloudWafv3Domains(),
			"alicloud_eflo_vpds":                                        dataSourceAlicloudEfloVpds(),
			"alicloud_dcdn_waf_rules":                                   dataSourceAlicloudDcdnWafRules(),
			"alicloud_actiontrail_global_events_storage_region":         dataSourceAlicloudActiontrailGlobalEventsStorageRegion(),
			"alicloud_dbfs_auto_snap_shot_policies":                     dataSourceAlicloudDbfsAutoSnapShotPolicies(),
			"alicloud_cen_transit_route_table_aggregations":             dataSourceAlicloudCenTransitRouteTableAggregations(),
			"alicloud_arms_prometheis":                                  dataSourceAlicloudArmsPrometheis(),
			"alicloud_arms_prometheus":                                  dataSourceAlicloudArmsPrometheis(),
			"alicloud_ocean_base_instances":                             dataSourceAlicloudOceanBaseInstances(),
			"alicloud_chatbot_agents":                                   dataSourceAlicloudChatbotAgents(),
			"alicloud_arms_integration_exporters":                       dataSourceAlicloudArmsIntegrationExporters(),
			"alicloud_service_catalog_portfolios":                       dataSourceAlicloudServiceCatalogPortfolios(),
			"alicloud_arms_remote_writes":                               dataSourceAlicloudArmsRemoteWrites(),
			"alicloud_eflo_subnets":                                     dataSourceAlicloudEfloSubnets(),
			"alicloud_compute_nest_service_instances":                   dataSourceAlicloudComputeNestServiceInstances(),
			"alicloud_vpc_flow_log_service":                             dataSourceAliCloudVpcFlowLogService(),
			"alicloud_arms_prometheus_monitorings":                      dataSourceAliCloudArmsPrometheusMonitorings(),
			"alicloud_ga_endpoint_group_ip_address_cidr_blocks":         dataSourceAliCloudGaEndpointGroupIpAddressCidrBlocks(),
			"alicloud_quotas_template_applications":                     dataSourceAliCloudQuotasTemplateApplications(),
			"alicloud_cloud_monitor_service_hybrid_double_writes":       dataSourceAliCloudCloudMonitorServiceHybridDoubleWrites(),
			"alicloud_cms_site_monitors":                                dataSourceAliCloudCloudMonitorServiceSiteMonitors(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_pai_workspace_run":                                    resourceAliCloudPaiWorkspaceRun(),
			"alicloud_pai_workspace_datasetversion":                         resourceAliCloudPaiWorkspaceDatasetversion(),
			"alicloud_pai_workspace_experiment":                             resourceAliCloudPaiWorkspaceExperiment(),
			"alicloud_pai_workspace_dataset":                                resourceAliCloudPaiWorkspaceDataset(),
			"alicloud_rds_custom_deployment_set":                            resourceAliCloudRdsCustomDeploymentSet(),
			"alicloud_rds_custom":                                           resourceAliCloudRdsCustom(),
			"alicloud_esa_rate_plan_instance":                               resourceAliCloudEsaRatePlanInstance(),
			"alicloud_vpc_ipam_ipam_pool_cidr":                              resourceAliCloudVpcIpamIpamPoolCidr(),
			"alicloud_vpc_ipam_ipam_pool":                                   resourceAliCloudVpcIpamIpamPool(),
			"alicloud_vpc_ipam_ipam_scope":                                  resourceAliCloudVpcIpamIpamScope(),
			"alicloud_vpc_ipam_ipam":                                        resourceAliCloudVpcIpamIpam(),
			"alicloud_gwlb_server_group":                                    resourceAliCloudGwlbServerGroup(),
			"alicloud_gwlb_listener":                                        resourceAliCloudGwlbListener(),
			"alicloud_gwlb_load_balancer":                                   resourceAliCloudGwlbLoadBalancer(),
			"alicloud_oss_bucket_cname_token":                               resourceAliCloudOssBucketCnameToken(),
			"alicloud_oss_bucket_cname":                                     resourceAliCloudOssBucketCname(),
			"alicloud_esa_site":                                             resourceAliCloudEsaSite(),
			"alicloud_pai_workspace_workspace":                              resourceAliCloudPaiWorkspaceWorkspace(),
			"alicloud_gpdb_database":                                        resourceAliCloudGpdbDatabase(),
			"alicloud_sls_collection_policy":                                resourceAliCloudSlsCollectionPolicy(),
			"alicloud_gpdb_db_instance_ip_array":                            resourceAliCloudGpdbDBInstanceIPArray(),
			"alicloud_quotas_template_service":                              resourceAliCloudQuotasTemplateService(),
			"alicloud_fcv3_vpc_binding":                                     resourceAliCloudFcv3VpcBinding(),
			"alicloud_fcv3_layer_version":                                   resourceAliCloudFcv3LayerVersion(),
			"alicloud_service_catalog_principal_portfolio_association":      resourceAliCloudServiceCatalogPrincipalPortfolioAssociation(),
			"alicloud_service_catalog_product_version":                      resourceAliCloudServiceCatalogProductVersion(),
			"alicloud_service_catalog_product_portfolio_association":        resourceAliCloudServiceCatalogProductPortfolioAssociation(),
			"alicloud_service_catalog_product":                              resourceAliCloudServiceCatalogProduct(),
			"alicloud_gpdb_hadoop_data_source":                              resourceAliCloudGpdbHadoopDataSource(),
			"alicloud_gpdb_jdbc_data_source":                                resourceAliCloudGpdbJdbcDataSource(),
			"alicloud_fcv3_provision_config":                                resourceAliCloudFcv3ProvisionConfig(),
			"alicloud_gpdb_streaming_job":                                   resourceAliCloudGpdbStreamingJob(),
			"alicloud_data_works_project":                                   resourceAliCloudDataWorksProject(),
			"alicloud_fcv3_function_version":                                resourceAliCloudFcv3FunctionVersion(),
			"alicloud_governance_account":                                   resourceAliCloudGovernanceAccount(),
			"alicloud_fcv3_trigger":                                         resourceAliCloudFcv3Trigger(),
			"alicloud_fcv3_concurrency_config":                              resourceAliCloudFcv3ConcurrencyConfig(),
			"alicloud_fcv3_async_invoke_config":                             resourceAliCloudFcv3AsyncInvokeConfig(),
			"alicloud_fcv3_alias":                                           resourceAliCloudFcv3Alias(),
			"alicloud_fcv3_custom_domain":                                   resourceAliCloudFcv3CustomDomain(),
			"alicloud_fcv3_function":                                        resourceAliCloudFcv3Function(),
			"alicloud_aligreen_oss_stock_task":                              resourceAliCloudAligreenOssStockTask(),
			"alicloud_aligreen_keyword_lib":                                 resourceAliCloudAligreenKeywordLib(),
			"alicloud_aligreen_image_lib":                                   resourceAliCloudAligreenImageLib(),
			"alicloud_aligreen_biz_type":                                    resourceAliCloudAligreenBizType(),
			"alicloud_aligreen_callback":                                    resourceAliCloudAligreenCallback(),
			"alicloud_aligreen_audit_callback":                              resourceAliCloudAligreenAuditCallback(),
			"alicloud_cloud_firewall_vpc_cen_tr_firewall":                   resourceAliCloudCloudFirewallVpcCenTrFirewall(),
			"alicloud_governance_baseline":                                  resourceAliCloudGovernanceBaseline(),
			"alicloud_gpdb_streaming_data_source":                           resourceAliCloudGpdbStreamingDataSource(),
			"alicloud_gpdb_streaming_data_service":                          resourceAliCloudGpdbStreamingDataService(),
			"alicloud_gpdb_external_data_service":                           resourceAliCloudGpdbExternalDataService(),
			"alicloud_gpdb_remote_adb_data_source":                          resourceAliCloudGpdbRemoteADBDataSource(),
			"alicloud_ens_nat_gateway":                                      resourceAliCloudEnsNatGateway(),
			"alicloud_ens_eip_instance_attachment":                          resourceAliCloudEnsEipInstanceAttachment(),
			"alicloud_ddos_bgp_policy":                                      resourceAliCloudDdosBgpPolicy(),
			"alicloud_cen_transit_router_ecr_attachment":                    resourceAliCloudCenTransitRouterEcrAttachment(),
			"alicloud_alb_load_balancer_security_group_attachment":          resourceAliCloudAlbLoadBalancerSecurityGroupAttachment(),
			"alicloud_gpdb_db_resource_group":                               resourceAliCloudGpdbDbResourceGroup(),
			"alicloud_cloud_firewall_nat_firewall":                          resourceAliCloudCloudFirewallNatFirewall(),
			"alicloud_oss_bucket_public_access_block":                       resourceAliCloudOssBucketPublicAccessBlock(),
			"alicloud_oss_account_public_access_block":                      resourceAliCloudOssAccountPublicAccessBlock(),
			"alicloud_oss_bucket_data_redundancy_transition":                resourceAliCloudOssBucketDataRedundancyTransition(),
			"alicloud_oss_bucket_meta_query":                                resourceAliCloudOssBucketMetaQuery(),
			"alicloud_oss_bucket_access_monitor":                            resourceAliCloudOssBucketAccessMonitor(),
			"alicloud_oss_bucket_user_defined_log_fields":                   resourceAliCloudOssBucketUserDefinedLogFields(),
			"alicloud_oss_bucket_transfer_acceleration":                     resourceAliCloudOssBucketTransferAcceleration(),
			"alicloud_sls_scheduled_sql":                                    resourceAliCloudSlsScheduledSQL(),
			"alicloud_express_connect_router_express_connect_router":        resourceAliCloudExpressConnectRouterExpressConnectRouter(),
			"alicloud_express_connect_router_vpc_association":               resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation(),
			"alicloud_express_connect_router_tr_association":                resourceAliCloudExpressConnectRouterExpressConnectRouterTrAssociation(),
			"alicloud_express_connect_router_vbr_child_instance":            resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstance(),
			"alicloud_express_connect_traffic_qos_rule":                     resourceAliCloudExpressConnectTrafficQosRule(),
			"alicloud_express_connect_traffic_qos_queue":                    resourceAliCloudExpressConnectTrafficQosQueue(),
			"alicloud_express_connect_traffic_qos_association":              resourceAliCloudExpressConnectTrafficQosAssociation(),
			"alicloud_express_connect_traffic_qos":                          resourceAliCloudExpressConnectTrafficQos(),
			"alicloud_nas_access_point":                                     resourceAliCloudNasAccessPoint(),
			"alicloud_api_gateway_access_control_list":                      resourceAliCloudApiGatewayAccessControlList(),
			"alicloud_api_gateway_acl_entry_attachment":                     resourceAliCloudApiGatewayAclEntryAttachment(),
			"alicloud_api_gateway_instance_acl_attachment":                  resourceAliCloudApiGatewayInstanceAclAttachment(),
			"alicloud_cloud_firewall_nat_firewall_control_policy":           resourceAliCloudCloudFirewallNatFirewallControlPolicy(),
			"alicloud_sls_alert":                                            resourceAliCloudSlsAlert(),
			"alicloud_oss_bucket_cors":                                      resourceAliCloudOssBucketCors(),
			"alicloud_oss_bucket_server_side_encryption":                    resourceAliCloudOssBucketServerSideEncryption(),
			"alicloud_oss_bucket_logging":                                   resourceAliCloudOssBucketLogging(),
			"alicloud_oss_bucket_request_payment":                           resourceAliCloudOssBucketRequestPayment(),
			"alicloud_oss_bucket_versioning":                                resourceAliCloudOssBucketVersioning(),
			"alicloud_oss_bucket_policy":                                    resourceAliCloudOssBucketPolicy(),
			"alicloud_oss_bucket_https_config":                              resourceAliCloudOssBucketHttpsConfig(),
			"alicloud_oss_bucket_referer":                                   resourceAliCloudOssBucketReferer(),
			"alicloud_hbr_policy_binding":                                   resourceAliCloudHbrPolicyBinding(),
			"alicloud_hbr_policy":                                           resourceAliCloudHbrPolicy(),
			"alicloud_oss_bucket_acl":                                       resourceAliCloudOssBucketAcl(),
			"alicloud_wafv3_defense_template":                               resourceAliCloudWafv3DefenseTemplate(),
			"alicloud_dfs_vsc_mount_point":                                  resourceAliCloudDfsVscMountPoint(),
			"alicloud_vpc_ipv6_address":                                     resourceAliCloudVpcIpv6Address(),
			"alicloud_api_gateway_instance":                                 resourceAliCloudApiGatewayInstance(),
			"alicloud_ebs_solution_instance":                                resourceAliCloudEbsSolutionInstance(),
			"alicloud_ens_instance_security_group_attachment":               resourceAliCloudEnsInstanceSecurityGroupAttachment(),
			"alicloud_ens_disk_instance_attachment":                         resourceAliCloudEnsDiskInstanceAttachment(),
			"alicloud_ens_image":                                            resourceAliCloudEnsImage(),
			"alicloud_ebs_enterprise_snapshot_policy_attachment":            resourceAliCloudEbsEnterpriseSnapshotPolicyAttachment(),
			"alicloud_ebs_enterprise_snapshot_policy":                       resourceAliCloudEbsEnterpriseSnapshotPolicy(),
			"alicloud_ebs_replica_group_drill":                              resourceAliCloudEbsReplicaGroupDrill(),
			"alicloud_ebs_replica_pair_drill":                               resourceAliCloudEbsReplicaPairDrill(),
			"alicloud_arms_synthetic_task":                                  resourceAliCloudArmsSyntheticTask(),
			"alicloud_cloud_monitor_service_enterprise_public":              resourceAliCloudCloudMonitorServiceEnterprisePublic(),
			"alicloud_cloud_monitor_service_basic_public":                   resourceAliCloudCloudMonitorServiceBasicPublic(),
			"alicloud_express_connect_ec_failover_test_job":                 resourceAliCloudExpressConnectEcFailoverTestJob(),
			"alicloud_arms_grafana_workspace":                               resourceAliCloudArmsGrafanaWorkspace(),
			"alicloud_realtime_compute_vvp_instance":                        resourceAliCloudRealtimeComputeVvpInstance(),
			"alicloud_quotas_template_applications":                         resourceAliCloudQuotasTemplateApplications(),
			"alicloud_threat_detection_oss_scan_config":                     resourceAliCloudThreatDetectionOssScanConfig(),
			"alicloud_threat_detection_malicious_file_whitelist_config":     resourceAliCloudThreatDetectionMaliciousFileWhitelistConfig(),
			"alicloud_adb_lake_account":                                     resourceAliCloudAdbLakeAccount(),
			"alicloud_ens_security_group":                                   resourceAliCloudEnsSecurityGroup(),
			"alicloud_ens_vswitch":                                          resourceAliCloudEnsVswitch(),
			"alicloud_ens_load_balancer":                                    resourceAliCloudEnsLoadBalancer(),
			"alicloud_ens_eip":                                              resourceAliCloudEnsEip(),
			"alicloud_ens_network":                                          resourceAliCloudEnsNetwork(),
			"alicloud_ens_snapshot":                                         resourceAliCloudEnsSnapshot(),
			"alicloud_ens_disk":                                             resourceAliCloudEnsDisk(),
			"alicloud_resource_manager_saved_query":                         resourceAliCloudResourceManagerSavedQuery(),
			"alicloud_threat_detection_sas_trail":                           resourceAliCloudThreatDetectionSasTrail(),
			"alicloud_threat_detection_image_event_operation":               resourceAliCloudThreatDetectionImageEventOperation(),
			"alicloud_arms_env_custom_job":                                  resourceAliCloudArmsEnvCustomJob(),
			"alicloud_arms_env_service_monitor":                             resourceAliCloudArmsEnvServiceMonitor(),
			"alicloud_arms_env_pod_monitor":                                 resourceAliCloudArmsEnvPodMonitor(),
			"alicloud_arms_addon_release":                                   resourceAliCloudArmsAddonRelease(),
			"alicloud_arms_env_feature":                                     resourceAliCloudArmsEnvFeature(),
			"alicloud_arms_environment":                                     resourceAliCloudArmsEnvironment(),
			"alicloud_hologram_instance":                                    resourceAliCloudHologramInstance(),
			"alicloud_ack_one_cluster":                                      resourceAliCloudAckOneCluster(),
			"alicloud_drds_polardbx_instance":                               resourceAliCloudDrdsPolardbxInstance(),
			"alicloud_gpdb_backup_policy":                                   resourceAliCloudGpdbBackupPolicy(),
			"alicloud_threat_detection_file_upload_limit":                   resourceAliCloudThreatDetectionFileUploadLimit(),
			"alicloud_threat_detection_client_file_protect":                 resourceAliCloudThreatDetectionClientFileProtect(),
			"alicloud_rocketmq_topic":                                       resourceAliCloudRocketmqTopic(),
			"alicloud_rocketmq_consumer_group":                              resourceAliCloudRocketmqConsumerGroup(),
			"alicloud_rocketmq_instance":                                    resourceAliCloudRocketmqInstance(),
			"alicloud_dms_enterprise_authority_template":                    resourceAliCloudDMSEnterpriseAuthorityTemplate(),
			"alicloud_kms_application_access_point":                         resourceAliCloudKmsApplicationAccessPoint(),
			"alicloud_kms_client_key":                                       resourceAliCloudKmsClientKey(),
			"alicloud_kms_policy":                                           resourceAliCloudKmsPolicy(),
			"alicloud_kms_network_rule":                                     resourceAliCloudKmsNetworkRule(),
			"alicloud_kms_instance":                                         resourceAliCloudKmsInstance(),
			"alicloud_threat_detection_client_user_define_rule":             resourceAliCloudThreatDetectionClientUserDefineRule(),
			"alicloud_ims_oidc_provider":                                    resourceAliCloudImsOidcProvider(),
			"alicloud_cddc_dedicated_propre_host":                           resourceAliCloudCddcDedicatedPropreHost(),
			"alicloud_nlb_listener_additional_certificate_attachment":       resourceAliCloudNlbListenerAdditionalCertificateAttachment(),
			"alicloud_nlb_loadbalancer_common_bandwidth_package_attachment": resourceAliCloudNlbLoadbalancerCommonBandwidthPackageAttachment(),
			"alicloud_arms_prometheus_monitoring":                           resourceAliCloudArmsPrometheusMonitoring(),
			"alicloud_vpc_gateway_endpoint_route_table_attachment":          resourceAliCloudVpcGatewayEndpointRouteTableAttachment(),
			"alicloud_ens_instance":                                         resourceAliCloudEnsInstance(),
			"alicloud_vpc_gateway_endpoint":                                 resourceAliCloudVpcGatewayEndpoint(),
			"alicloud_eip_segment_address":                                  resourceAliCloudEipSegmentAddress(),
			"alicloud_fcv2_function":                                        resourceAliCloudFcv2Function(),
			"alicloud_quotas_template_quota":                                resourceAliCloudQuotasTemplateQuota(),
			"alicloud_redis_tair_instance":                                  resourceAliCloudRedisTairInstance(),
			"alicloud_vpc_vswitch_cidr_reservation":                         resourceAliCloudVpcVswitchCidrReservation(),
			"alicloud_vpc_ha_vip":                                           resourceAliCloudVpcHaVip(),
			"alicloud_config_remediation":                                   resourceAliCloudConfigRemediation(),
			"alicloud_instance":                                             resourceAliCloudInstance(),
			"alicloud_image":                                                resourceAliCloudEcsImage(),
			"alicloud_reserved_instance":                                    resourceAliCloudReservedInstance(),
			"alicloud_copy_image":                                           resourceAliCloudImageCopy(),
			"alicloud_image_export":                                         resourceAliCloudImageExport(),
			"alicloud_image_copy":                                           resourceAliCloudImageCopy(),
			"alicloud_image_import":                                         resourceAliCloudImageImport(),
			"alicloud_image_share_permission":                               resourceAliCloudImageSharePermission(),
			"alicloud_ram_role_attachment":                                  resourceAlicloudRamRoleAttachment(),
			"alicloud_disk":                                                 resourceAliCloudEcsDisk(),
			"alicloud_disk_attachment":                                      resourceAlicloudEcsDiskAttachment(),
			"alicloud_network_interface":                                    resourceAliCloudEcsNetworkInterface(),
			"alicloud_network_interface_attachment":                         resourceAliCloudEcsNetworkInterfaceAttachment(),
			"alicloud_snapshot":                                             resourceAliCloudEcsSnapshot(),
			"alicloud_snapshot_policy":                                      resourceAlicloudEcsAutoSnapshotPolicy(),
			"alicloud_launch_template":                                      resourceAliCloudEcsLaunchTemplate(),
			"alicloud_security_group":                                       resourceAliyunSecurityGroup(),
			"alicloud_security_group_rule":                                  resourceAliyunSecurityGroupRule(),
			"alicloud_db_database":                                          resourceAlicloudDBDatabase(),
			"alicloud_db_account":                                           resourceAlicloudRdsAccount(),
			"alicloud_db_account_privilege":                                 resourceAlicloudDBAccountPrivilege(),
			"alicloud_db_backup_policy":                                     resourceAlicloudDBBackupPolicy(),
			"alicloud_db_connection":                                        resourceAlicloudDBConnection(),
			"alicloud_db_read_write_splitting_connection":                   resourceAlicloudDBReadWriteSplittingConnection(),
			"alicloud_db_instance":                                          resourceAliCloudDBInstance(),
			"alicloud_rds_backup":                                           resourceAlicloudRdsBackup(),
			"alicloud_rds_db_proxy":                                         resourceAlicloudRdsDBProxy(),
			"alicloud_rds_clone_db_instance":                                resourceAlicloudRdsCloneDbInstance(),
			"alicloud_rds_upgrade_db_instance":                              resourceAlicloudRdsUpgradeDbInstance(),
			"alicloud_rds_instance_cross_backup_policy":                     resourceAlicloudRdsInstanceCrossBackupPolicy(),
			"alicloud_rds_ddr_instance":                                     resourceAlicloudRdsDdrInstance(),
			"alicloud_mongodb_instance":                                     resourceAliCloudMongoDBInstance(),
			"alicloud_mongodb_sharding_instance":                            resourceAliCloudMongoDBShardingInstance(),
			"alicloud_gpdb_instance":                                        resourceAliCloudGpdbInstance(),
			"alicloud_gpdb_elastic_instance":                                resourceAlicloudGpdbElasticInstance(),
			"alicloud_gpdb_connection":                                      resourceAlicloudGpdbConnection(),
			"alicloud_tag_policy":                                           resourceAlicloudTagPolicy(),
			"alicloud_tag_policy_attachment":                                resourceAlicloudTagPolicyAttachment(),
			"alicloud_db_readonly_instance":                                 resourceAlicloudDBReadonlyInstance(),
			"alicloud_auto_provisioning_group":                              resourceAlicloudAutoProvisioningGroup(),
			"alicloud_ess_scaling_group":                                    resourceAlicloudEssScalingGroup(),
			"alicloud_ess_eci_scaling_configuration":                        resourceAlicloudEssEciScalingConfiguration(),
			"alicloud_ess_scaling_configuration":                            resourceAlicloudEssScalingConfiguration(),
			"alicloud_ess_scaling_rule":                                     resourceAlicloudEssScalingRule(),
			"alicloud_ess_schedule":                                         resourceAlicloudEssScheduledTask(),
			"alicloud_ess_scheduled_task":                                   resourceAlicloudEssScheduledTask(),
			"alicloud_ess_attachment":                                       resourceAlicloudEssAttachment(),
			"alicloud_ess_suspend_process":                                  resourceAlicloudEssSuspendProcess(),
			"alicloud_ess_lifecycle_hook":                                   resourceAlicloudEssLifecycleHook(),
			"alicloud_ess_notification":                                     resourceAlicloudEssNotification(),
			"alicloud_ess_alarm":                                            resourceAlicloudEssAlarm(),
			"alicloud_ess_scalinggroup_vserver_groups":                      resourceAlicloudEssScalingGroupVserverGroups(),
			"alicloud_ess_alb_server_group_attachment":                      resourceAlicloudEssAlbServerGroupAttachment(),
			"alicloud_ess_server_group_attachment":                          resourceAliCloudEssServerGroupAttachment(),
			"alicloud_vpc":                                                  resourceAliCloudVpcVpc(),
			"alicloud_nat_gateway":                                          resourceAliCloudNatGateway(),
			"alicloud_nas_file_system":                                      resourceAliCloudNasFileSystem(),
			"alicloud_nas_mount_target":                                     resourceAlicloudNasMountTarget(),
			"alicloud_nas_access_group":                                     resourceAliCloudNasAccessGroup(),
			"alicloud_nas_access_rule":                                      resourceAliCloudNasAccessRule(),
			"alicloud_nas_smb_acl_attachment":                               resourceAlicloudNasSmbAclAttachment(),
			"alicloud_tag_meta_tag":                                         resourceAlicloudTagMetaTag(),
			// "alicloud_subnet" aims to match aws usage habit.
			"alicloud_subnet":                        resourceAliCloudVpcVswitch(),
			"alicloud_vswitch":                       resourceAliCloudVpcVswitch(),
			"alicloud_route_entry":                   resourceAliyunRouteEntry(),
			"alicloud_route_table":                   resourceAliCloudVpcRouteTable(),
			"alicloud_route_table_attachment":        resourceAliCloudVpcRouteTableAttachment(),
			"alicloud_snat_entry":                    resourceAlicloudSnatEntry(),
			"alicloud_forward_entry":                 resourceAlicloudForwardEntry(),
			"alicloud_eip":                           resourceAliCloudEipAddress(),
			"alicloud_eip_association":               resourceAliCloudEipAssociation(),
			"alicloud_slb":                           resourceAlicloudSlbLoadBalancer(),
			"alicloud_slb_listener":                  resourceAliCloudSlbListener(),
			"alicloud_slb_attachment":                resourceAliyunSlbAttachment(),
			"alicloud_slb_backend_server":            resourceAliyunSlbBackendServer(),
			"alicloud_slb_domain_extension":          resourceAlicloudSlbDomainExtension(),
			"alicloud_slb_server_group":              resourceAliyunSlbServerGroup(),
			"alicloud_slb_master_slave_server_group": resourceAliyunSlbMasterSlaveServerGroup(),
			"alicloud_slb_rule":                      resourceAliyunSlbRule(),
			"alicloud_slb_acl":                       resourceAlicloudSlbAcl(),
			"alicloud_slb_ca_certificate":            resourceAlicloudSlbCaCertificate(),
			"alicloud_slb_server_certificate":        resourceAlicloudSlbServerCertificate(),
			"alicloud_oss_bucket":                    resourceAlicloudOssBucket(),
			"alicloud_oss_bucket_object":             resourceAlicloudOssBucketObject(),
			"alicloud_oss_bucket_replication":        resourceAlicloudOssBucketReplication(),
			"alicloud_ons_instance":                  resourceAlicloudOnsInstance(),
			"alicloud_ons_topic":                     resourceAlicloudOnsTopic(),
			"alicloud_ons_group":                     resourceAlicloudOnsGroup(),
			"alicloud_alikafka_consumer_group":       resourceAlicloudAlikafkaConsumerGroup(),
			"alicloud_alikafka_instance":             resourceAliCloudAlikafkaInstance(),
			"alicloud_alikafka_topic":                resourceAlicloudAlikafkaTopic(),
			"alicloud_alikafka_sasl_user":            resourceAliCloudAlikafkaSaslUser(),
			"alicloud_alikafka_sasl_acl":             resourceAlicloudAlikafkaSaslAcl(),
			"alicloud_dns_record":                    resourceAlicloudDnsRecord(),
			"alicloud_dns":                           resourceAlicloudDns(),
			"alicloud_dns_group":                     resourceAlicloudDnsGroup(),
			"alicloud_key_pair":                      resourceAlicloudEcsKeyPair(),
			"alicloud_key_pair_attachment":           resourceAlicloudEcsKeyPairAttachment(),
			"alicloud_kms_key":                       resourceAliCloudKmsKey(),
			"alicloud_kms_ciphertext":                resourceAlicloudKmsCiphertext(),
			"alicloud_ram_user":                      resourceAlicloudRamUser(),
			"alicloud_ram_account_password_policy":   resourceAlicloudRamAccountPasswordPolicy(),
			"alicloud_ram_access_key":                resourceAlicloudRamAccessKey(),
			"alicloud_ram_login_profile":             resourceAliCloudRamLoginProfile(),
			"alicloud_ram_group":                     resourceAlicloudRamGroup(),
			"alicloud_ram_role":                      resourceAlicloudRamRole(),
			"alicloud_ram_policy":                    resourceAlicloudRamPolicy(),
			// alicloud_ram_alias has been deprecated
			"alicloud_ram_alias":                                             resourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_alias":                                     resourceAlicloudRamAccountAlias(),
			"alicloud_ram_group_membership":                                  resourceAlicloudRamGroupMembership(),
			"alicloud_ram_user_policy_attachment":                            resourceAlicloudRamUserPolicyAtatchment(),
			"alicloud_ram_role_policy_attachment":                            resourceAlicloudRamRolePolicyAttachment(),
			"alicloud_ram_group_policy_attachment":                           resourceAlicloudRamGroupPolicyAtatchment(),
			"alicloud_container_cluster":                                     resourceAlicloudCSSwarm(),
			"alicloud_cs_application":                                        resourceAlicloudCSApplication(),
			"alicloud_cs_swarm":                                              resourceAlicloudCSSwarm(),
			"alicloud_cs_kubernetes":                                         resourceAlicloudCSKubernetes(),
			"alicloud_cs_kubernetes_addon":                                   resourceAlicloudCSKubernetesAddon(),
			"alicloud_cs_managed_kubernetes":                                 resourceAlicloudCSManagedKubernetes(),
			"alicloud_cs_edge_kubernetes":                                    resourceAlicloudCSEdgeKubernetes(),
			"alicloud_cs_serverless_kubernetes":                              resourceAlicloudCSServerlessKubernetes(),
			"alicloud_cs_kubernetes_autoscaler":                              resourceAlicloudCSKubernetesAutoscaler(),
			"alicloud_cs_kubernetes_node_pool":                               resourceAliCloudAckNodepool(),
			"alicloud_cs_kubernetes_permissions":                             resourceAlicloudCSKubernetesPermissions(),
			"alicloud_cs_autoscaling_config":                                 resourceAlicloudCSAutoscalingConfig(),
			"alicloud_cr_namespace":                                          resourceAlicloudCRNamespace(),
			"alicloud_cr_repo":                                               resourceAlicloudCRRepo(),
			"alicloud_cr_ee_instance":                                        resourceAliCloudCrInstance(),
			"alicloud_cr_ee_namespace":                                       resourceAliCloudCrEENamespace(),
			"alicloud_cr_ee_repo":                                            resourceAliCloudCrEERepo(),
			"alicloud_cr_ee_sync_rule":                                       resourceAliCloudCrEESyncRule(),
			"alicloud_cdn_domain":                                            resourceAlicloudCdnDomain(),
			"alicloud_cdn_domain_new":                                        resourceAliCloudCdnDomain(),
			"alicloud_cdn_domain_config":                                     resourceAliCloudCdnDomainConfig(),
			"alicloud_router_interface":                                      resourceAlicloudRouterInterface(),
			"alicloud_router_interface_connection":                           resourceAlicloudRouterInterfaceConnection(),
			"alicloud_ots_table":                                             resourceAlicloudOtsTable(),
			"alicloud_ots_instance":                                          resourceAlicloudOtsInstance(),
			"alicloud_ots_instance_attachment":                               resourceAlicloudOtsInstanceAttachment(),
			"alicloud_ots_tunnel":                                            resourceAlicloudOtsTunnel(),
			"alicloud_ots_secondary_index":                                   resourceAlicloudOtsSecondaryIndex(),
			"alicloud_ots_search_index":                                      resourceAlicloudOtsSearchIndex(),
			"alicloud_cms_alarm":                                             resourceAliCloudCmsAlarm(),
			"alicloud_cms_site_monitor":                                      resourceAlicloudCmsSiteMonitor(),
			"alicloud_pvtz_zone":                                             resourceAlicloudPvtzZone(),
			"alicloud_pvtz_zone_attachment":                                  resourceAlicloudPvtzZoneAttachment(),
			"alicloud_pvtz_zone_record":                                      resourceAlicloudPvtzZoneRecord(),
			"alicloud_log_alert":                                             resourceAlicloudLogAlert(),
			"alicloud_log_alert_resource":                                    resourceAlicloudLogAlertResource(),
			"alicloud_log_audit":                                             resourceAlicloudLogAudit(),
			"alicloud_log_dashboard":                                         resourceAlicloudLogDashboard(),
			"alicloud_log_etl":                                               resourceAlicloudLogETL(),
			"alicloud_log_ingestion":                                         resourceAlicloudLogIngestion(),
			"alicloud_log_machine_group":                                     resourceAlicloudLogMachineGroup(),
			"alicloud_log_oss_export":                                        resourceAlicloudLogOssExport(),
			"alicloud_log_oss_shipper":                                       resourceAlicloudLogOssShipper(),
			"alicloud_log_project":                                           resourceAliCloudSlsProject(),
			"alicloud_log_resource":                                          resourceAlicloudLogResource(),
			"alicloud_log_resource_record":                                   resourceAlicloudLogResourceRecord(),
			"alicloud_log_store":                                             resourceAliCloudSlsLogStore(),
			"alicloud_log_store_index":                                       resourceAlicloudLogStoreIndex(),
			"alicloud_logtail_config":                                        resourceAlicloudLogtailConfig(),
			"alicloud_logtail_attachment":                                    resourceAlicloudLogtailAttachment(),
			"alicloud_fc_service":                                            resourceAlicloudFCService(),
			"alicloud_fc_function":                                           resourceAlicloudFCFunction(),
			"alicloud_fc_trigger":                                            resourceAlicloudFCTrigger(),
			"alicloud_fc_alias":                                              resourceAlicloudFCAlias(),
			"alicloud_fc_custom_domain":                                      resourceAlicloudFCCustomDomain(),
			"alicloud_fc_function_async_invoke_config":                       resourceAlicloudFCFunctionAsyncInvokeConfig(),
			"alicloud_vpn_gateway":                                           resourceAliCloudVPNGatewayVPNGateway(),
			"alicloud_vpn_customer_gateway":                                  resourceAliCloudVPNGatewayCustomerGateway(),
			"alicloud_vpn_route_entry":                                       resourceAliyunVpnRouteEntry(),
			"alicloud_vpn_connection":                                        resourceAliCloudVPNGatewayVpnConnection(),
			"alicloud_ssl_vpn_server":                                        resourceAliyunSslVpnServer(),
			"alicloud_ssl_vpn_client_cert":                                   resourceAliyunSslVpnClientCert(),
			"alicloud_cen_instance":                                          resourceAliCloudCenInstance(),
			"alicloud_cen_instance_attachment":                               resourceAlicloudCenInstanceAttachment(),
			"alicloud_cen_bandwidth_package":                                 resourceAlicloudCenBandwidthPackage(),
			"alicloud_cen_bandwidth_package_attachment":                      resourceAlicloudCenBandwidthPackageAttachment(),
			"alicloud_cen_bandwidth_limit":                                   resourceAlicloudCenBandwidthLimit(),
			"alicloud_cen_route_entry":                                       resourceAlicloudCenRouteEntry(),
			"alicloud_cen_instance_grant":                                    resourceAlicloudCenInstanceGrant(),
			"alicloud_cen_transit_router":                                    resourceAlicloudCenTransitRouter(),
			"alicloud_cen_transit_router_route_entry":                        resourceAlicloudCenTransitRouterRouteEntry(),
			"alicloud_cen_transit_router_route_table":                        resourceAlicloudCenTransitRouterRouteTable(),
			"alicloud_cen_transit_router_route_table_association":            resourceAlicloudCenTransitRouterRouteTableAssociation(),
			"alicloud_cen_transit_router_route_table_propagation":            resourceAlicloudCenTransitRouterRouteTablePropagation(),
			"alicloud_cen_transit_router_vbr_attachment":                     resourceAliCloudCenTransitRouterVbrAttachment(),
			"alicloud_cen_transit_router_vpc_attachment":                     resourceAliCloudCenTransitRouterVpcAttachment(),
			"alicloud_kvstore_instance":                                      resourceAliCloudKvstoreInstance(),
			"alicloud_kvstore_backup_policy":                                 resourceAlicloudKVStoreBackupPolicy(),
			"alicloud_kvstore_account":                                       resourceAlicloudKvstoreAccount(),
			"alicloud_datahub_project":                                       resourceAlicloudDatahubProject(),
			"alicloud_datahub_subscription":                                  resourceAlicloudDatahubSubscription(),
			"alicloud_datahub_topic":                                         resourceAlicloudDatahubTopic(),
			"alicloud_mns_queue":                                             resourceAlicloudMNSQueue(),
			"alicloud_mns_topic":                                             resourceAlicloudMNSTopic(),
			"alicloud_havip":                                                 resourceAliCloudVpcHaVip(),
			"alicloud_mns_topic_subscription":                                resourceAlicloudMNSSubscription(),
			"alicloud_havip_attachment":                                      resourceAliCloudVpcHaVipAttachment(),
			"alicloud_api_gateway_api":                                       resourceAliyunApigatewayApi(),
			"alicloud_api_gateway_group":                                     resourceAliyunApigatewayGroup(),
			"alicloud_api_gateway_app":                                       resourceAliyunApigatewayApp(),
			"alicloud_api_gateway_app_attachment":                            resourceAliyunApigatewayAppAttachment(),
			"alicloud_api_gateway_vpc_access":                                resourceAliyunApigatewayVpc(),
			"alicloud_common_bandwidth_package":                              resourceAliCloudCbwpCommonBandwidthPackage(),
			"alicloud_common_bandwidth_package_attachment":                   resourceAliCloudCbwpCommonBandwidthPackageAttachment(),
			"alicloud_drds_instance":                                         resourceAlicloudDRDSInstance(),
			"alicloud_elasticsearch_instance":                                resourceAlicloudElasticsearch(),
			"alicloud_cas_certificate":                                       resourceAliCloudSslCertificatesServiceCertificate(),
			"alicloud_ddoscoo_instance":                                      resourceAliCloudDdoscooInstance(),
			"alicloud_ddosbgp_instance":                                      resourceAlicloudDdosbgpInstance(),
			"alicloud_network_acl":                                           resourceAliCloudVpcNetworkAcl(),
			"alicloud_network_acl_attachment":                                resourceAliyunNetworkAclAttachment(),
			"alicloud_network_acl_entries":                                   resourceAliyunNetworkAclEntries(),
			"alicloud_emr_cluster":                                           resourceAlicloudEmrCluster(),
			"alicloud_emrv2_cluster":                                         resourceAlicloudEmrV2Cluster(),
			"alicloud_cloud_connect_network":                                 resourceAlicloudCloudConnectNetwork(),
			"alicloud_cloud_connect_network_attachment":                      resourceAlicloudCloudConnectNetworkAttachment(),
			"alicloud_cloud_connect_network_grant":                           resourceAlicloudCloudConnectNetworkGrant(),
			"alicloud_sag_acl":                                               resourceAlicloudSagAcl(),
			"alicloud_sag_acl_rule":                                          resourceAlicloudSagAclRule(),
			"alicloud_sag_qos":                                               resourceAlicloudSagQos(),
			"alicloud_sag_qos_policy":                                        resourceAlicloudSagQosPolicy(),
			"alicloud_sag_qos_car":                                           resourceAlicloudSagQosCar(),
			"alicloud_sag_snat_entry":                                        resourceAlicloudSagSnatEntry(),
			"alicloud_sag_dnat_entry":                                        resourceAlicloudSagDnatEntry(),
			"alicloud_sag_client_user":                                       resourceAlicloudSagClientUser(),
			"alicloud_yundun_dbaudit_instance":                               resourceAlicloudDbauditInstance(),
			"alicloud_yundun_bastionhost_instance":                           resourceAlicloudBastionhostInstance(),
			"alicloud_bastionhost_instance":                                  resourceAlicloudBastionhostInstance(),
			"alicloud_polardb_cluster":                                       resourceAlicloudPolarDBCluster(),
			"alicloud_polardb_cluster_endpoint":                              resourceAlicloudPolarDBClusterEndpoint(),
			"alicloud_polardb_backup_policy":                                 resourceAlicloudPolarDBBackupPolicy(),
			"alicloud_polardb_database":                                      resourceAlicloudPolarDBDatabase(),
			"alicloud_polardb_account":                                       resourceAlicloudPolarDBAccount(),
			"alicloud_polardb_account_privilege":                             resourceAlicloudPolarDBAccountPrivilege(),
			"alicloud_polardb_endpoint":                                      resourceAlicloudPolarDBEndpoint(),
			"alicloud_polardb_endpoint_address":                              resourceAlicloudPolarDBEndpointAddress(),
			"alicloud_polardb_primary_endpoint":                              resourceAlicloudPolarDBPrimaryEndpoint(),
			"alicloud_hbase_instance":                                        resourceAlicloudHBaseInstance(),
			"alicloud_market_order":                                          resourceAlicloudMarketOrder(),
			"alicloud_adb_cluster":                                           resourceAliCloudAdbDbCluster(),
			"alicloud_adb_backup_policy":                                     resourceAlicloudAdbBackupPolicy(),
			"alicloud_adb_account":                                           resourceAlicloudAdbAccount(),
			"alicloud_adb_connection":                                        resourceAlicloudAdbConnection(),
			"alicloud_cen_flowlog":                                           resourceAliCloudCenFlowLog(),
			"alicloud_kms_secret":                                            resourceAliCloudKmsSecret(),
			"alicloud_maxcompute_project":                                    resourceAliCloudMaxComputeProject(),
			"alicloud_kms_alias":                                             resourceAlicloudKmsAlias(),
			"alicloud_dns_instance":                                          resourceAlicloudAlidnsInstance(),
			"alicloud_dns_domain_attachment":                                 resourceAlicloudAlidnsDomainAttachment(),
			"alicloud_alidns_domain_attachment":                              resourceAlicloudAlidnsDomainAttachment(),
			"alicloud_edas_application":                                      resourceAlicloudEdasApplication(),
			"alicloud_edas_deploy_group":                                     resourceAlicloudEdasDeployGroup(),
			"alicloud_edas_application_scale":                                resourceAlicloudEdasApplicationScale(),
			"alicloud_edas_slb_attachment":                                   resourceAlicloudEdasSlbAttachment(),
			"alicloud_edas_cluster":                                          resourceAlicloudEdasCluster(),
			"alicloud_edas_instance_cluster_attachment":                      resourceAlicloudEdasInstanceClusterAttachment(),
			"alicloud_edas_application_deployment":                           resourceAlicloudEdasApplicationPackageAttachment(),
			"alicloud_dns_domain":                                            resourceAlicloudAlidnsDomain(),
			"alicloud_dms_enterprise_instance":                               resourceAlicloudDmsEnterpriseInstance(),
			"alicloud_waf_domain":                                            resourceAlicloudWafDomain(),
			"alicloud_cen_route_map":                                         resourceAlicloudCenRouteMap(),
			"alicloud_resource_manager_role":                                 resourceAlicloudResourceManagerRole(),
			"alicloud_resource_manager_resource_group":                       resourceAliCloudResourceManagerResourceGroup(),
			"alicloud_resource_manager_folder":                               resourceAlicloudResourceManagerFolder(),
			"alicloud_resource_manager_handshake":                            resourceAlicloudResourceManagerHandshake(),
			"alicloud_cen_private_zone":                                      resourceAlicloudCenPrivateZone(),
			"alicloud_resource_manager_policy":                               resourceAlicloudResourceManagerPolicy(),
			"alicloud_resource_manager_account":                              resourceAlicloudResourceManagerAccount(),
			"alicloud_waf_instance":                                          resourceAlicloudWafInstance(),
			"alicloud_resource_manager_resource_directory":                   resourceAliCloudResourceManagerResourceDirectory(),
			"alicloud_alidns_domain_group":                                   resourceAlicloudAlidnsDomainGroup(),
			"alicloud_resource_manager_policy_version":                       resourceAlicloudResourceManagerPolicyVersion(),
			"alicloud_kms_key_version":                                       resourceAlicloudKmsKeyVersion(),
			"alicloud_alidns_record":                                         resourceAlicloudAlidnsRecord(),
			"alicloud_ddoscoo_scheduler_rule":                                resourceAlicloudDdoscooSchedulerRule(),
			"alicloud_cassandra_cluster":                                     resourceAlicloudCassandraCluster(),
			"alicloud_cassandra_data_center":                                 resourceAlicloudCassandraDataCenter(),
			"alicloud_cen_vbr_health_check":                                  resourceAlicloudCenVbrHealthCheck(),
			"alicloud_eci_openapi_image_cache":                               resourceAlicloudEciImageCache(),
			"alicloud_eci_image_cache":                                       resourceAlicloudEciImageCache(),
			"alicloud_dms_enterprise_user":                                   resourceAlicloudDmsEnterpriseUser(),
			"alicloud_ecs_dedicated_host":                                    resourceAlicloudEcsDedicatedHost(),
			"alicloud_oos_template":                                          resourceAlicloudOosTemplate(),
			"alicloud_edas_k8s_cluster":                                      resourceAlicloudEdasK8sCluster(),
			"alicloud_oos_execution":                                         resourceAlicloudOosExecution(),
			"alicloud_resource_manager_policy_attachment":                    resourceAlicloudResourceManagerPolicyAttachment(),
			"alicloud_dcdn_domain":                                           resourceAliCloudDcdnDomain(),
			"alicloud_mse_cluster":                                           resourceAlicloudMseCluster(),
			"alicloud_actiontrail_trail":                                     resourceAlicloudActiontrailTrail(),
			"alicloud_actiontrail":                                           resourceAlicloudActiontrailTrail(),
			"alicloud_alidns_domain":                                         resourceAlicloudAlidnsDomain(),
			"alicloud_alidns_instance":                                       resourceAlicloudAlidnsInstance(),
			"alicloud_edas_k8s_application":                                  resourceAlicloudEdasK8sApplication(),
			"alicloud_edas_k8s_slb_attachment":                               resourceAlicloudEdasK8sSlbAttachment(),
			"alicloud_config_rule":                                           resourceAliCloudConfigRule(),
			"alicloud_config_configuration_recorder":                         resourceAlicloudConfigConfigurationRecorder(),
			"alicloud_config_delivery_channel":                               resourceAlicloudConfigDeliveryChannel(),
			"alicloud_cms_alarm_contact":                                     resourceAlicloudCmsAlarmContact(),
			"alicloud_cen_route_service":                                     resourceAlicloudCenRouteService(),
			"alicloud_kvstore_connection":                                    resourceAlicloudKvstoreConnection(),
			"alicloud_cms_alarm_contact_group":                               resourceAlicloudCmsAlarmContactGroup(),
			"alicloud_cms_group_metric_rule":                                 resourceAliCloudCmsGroupMetricRule(),
			"alicloud_fnf_flow":                                              resourceAlicloudFnfFlow(),
			"alicloud_fnf_schedule":                                          resourceAlicloudFnfSchedule(),
			"alicloud_ros_change_set":                                        resourceAlicloudRosChangeSet(),
			"alicloud_ros_stack":                                             resourceAlicloudRosStack(),
			"alicloud_ros_stack_group":                                       resourceAlicloudRosStackGroup(),
			"alicloud_ros_template":                                          resourceAlicloudRosTemplate(),
			"alicloud_privatelink_vpc_endpoint_service":                      resourceAliCloudPrivateLinkVpcEndpointService(),
			"alicloud_privatelink_vpc_endpoint":                              resourceAliCloudPrivateLinkVpcEndpoint(),
			"alicloud_privatelink_vpc_endpoint_connection":                   resourceAliCloudPrivateLinkVpcEndpointConnection(),
			"alicloud_privatelink_vpc_endpoint_service_resource":             resourceAliCloudPrivateLinkVpcEndpointServiceResource(),
			"alicloud_privatelink_vpc_endpoint_service_user":                 resourceAliCloudPrivateLinkVpcEndpointServiceUser(),
			"alicloud_resource_manager_resource_share":                       resourceAliCloudResourceManagerResourceShare(),
			"alicloud_privatelink_vpc_endpoint_zone":                         resourceAliCloudPrivateLinkVpcEndpointZone(),
			"alicloud_ga_accelerator":                                        resourceAliCloudGaAccelerator(),
			"alicloud_eci_container_group":                                   resourceAlicloudEciContainerGroup(),
			"alicloud_resource_manager_shared_resource":                      resourceAlicloudResourceManagerSharedResource(),
			"alicloud_resource_manager_shared_target":                        resourceAlicloudResourceManagerSharedTarget(),
			"alicloud_ga_listener":                                           resourceAliCloudGaListener(),
			"alicloud_tsdb_instance":                                         resourceAlicloudTsdbInstance(),
			"alicloud_ga_bandwidth_package":                                  resourceAliCloudGaBandwidthPackage(),
			"alicloud_ga_endpoint_group":                                     resourceAliCloudGaEndpointGroup(),
			"alicloud_brain_industrial_pid_organization":                     resourceAlicloudBrainIndustrialPidOrganization(),
			"alicloud_ga_bandwidth_package_attachment":                       resourceAliCloudGaBandwidthPackageAttachment(),
			"alicloud_ga_ip_set":                                             resourceAliCloudGaIpSet(),
			"alicloud_ga_forwarding_rule":                                    resourceAliCloudGaForwardingRule(),
			"alicloud_eipanycast_anycast_eip_address":                        resourceAliCloudEipanycastAnycastEipAddress(),
			"alicloud_brain_industrial_pid_project":                          resourceAlicloudBrainIndustrialPidProject(),
			"alicloud_cms_monitor_group":                                     resourceAlicloudCmsMonitorGroup(),
			"alicloud_eipanycast_anycast_eip_address_attachment":             resourceAliCloudEipanycastAnycastEipAddressAttachment(),
			"alicloud_ram_saml_provider":                                     resourceAliCloudRamSamlProvider(),
			"alicloud_quotas_application_info":                               resourceAliCloudQuotasQuotaApplication(),
			"alicloud_cms_monitor_group_instances":                           resourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_quotas_quota_alarm":                                    resourceAliCloudQuotasQuotaAlarm(),
			"alicloud_ecs_command":                                           resourceAlicloudEcsCommand(),
			"alicloud_cloud_storage_gateway_storage_bundle":                  resourceAlicloudCloudStorageGatewayStorageBundle(),
			"alicloud_ecs_hpc_cluster":                                       resourceAlicloudEcsHpcCluster(),
			"alicloud_vpc_flow_log":                                          resourceAliCloudVpcFlowLog(),
			"alicloud_brain_industrial_pid_loop":                             resourceAlicloudBrainIndustrialPidLoop(),
			"alicloud_quotas_quota_application":                              resourceAliCloudQuotasQuotaApplication(),
			"alicloud_ecs_auto_snapshot_policy":                              resourceAlicloudEcsAutoSnapshotPolicy(),
			"alicloud_rds_parameter_group":                                   resourceAlicloudRdsParameterGroup(),
			"alicloud_ecs_launch_template":                                   resourceAliCloudEcsLaunchTemplate(),
			"alicloud_resource_manager_control_policy":                       resourceAlicloudResourceManagerControlPolicy(),
			"alicloud_resource_manager_control_policy_attachment":            resourceAlicloudResourceManagerControlPolicyAttachment(),
			"alicloud_rds_account":                                           resourceAlicloudRdsAccount(),
			"alicloud_rds_db_node":                                           resourceAlicloudRdsDBNode(),
			"alicloud_rds_db_instance_endpoint":                              resourceAlicloudRdsDBInstanceEndpoint(),
			"alicloud_rds_db_instance_endpoint_address":                      resourceAlicloudRdsDBInstanceEndpointAddress(),
			"alicloud_ecs_snapshot":                                          resourceAliCloudEcsSnapshot(),
			"alicloud_ecs_key_pair":                                          resourceAlicloudEcsKeyPair(),
			"alicloud_ecs_key_pair_attachment":                               resourceAlicloudEcsKeyPairAttachment(),
			"alicloud_adb_db_cluster":                                        resourceAliCloudAdbDbCluster(),
			"alicloud_ecs_disk":                                              resourceAliCloudEcsDisk(),
			"alicloud_ecs_disk_attachment":                                   resourceAlicloudEcsDiskAttachment(),
			"alicloud_ecs_auto_snapshot_policy_attachment":                   resourceAlicloudEcsAutoSnapshotPolicyAttachment(),
			"alicloud_ddoscoo_domain_resource":                               resourceAliCloudDdosCooDomainResource(),
			"alicloud_ddoscoo_port":                                          resourceAliCloudDdosCooPort(),
			"alicloud_slb_load_balancer":                                     resourceAlicloudSlbLoadBalancer(),
			"alicloud_ecs_network_interface":                                 resourceAliCloudEcsNetworkInterface(),
			"alicloud_ecs_network_interface_attachment":                      resourceAliCloudEcsNetworkInterfaceAttachment(),
			"alicloud_config_aggregator":                                     resourceAlicloudConfigAggregator(),
			"alicloud_config_aggregate_config_rule":                          resourceAlicloudConfigAggregateConfigRule(),
			"alicloud_config_aggregate_compliance_pack":                      resourceAliCloudConfigAggregateCompliancePack(),
			"alicloud_config_compliance_pack":                                resourceAliCloudConfigCompliancePack(),
			"alicloud_direct_mail_receivers":                                 resourceAlicloudDirectMailReceivers(),
			"alicloud_eip_address":                                           resourceAliCloudEipAddress(),
			"alicloud_event_bridge_event_bus":                                resourceAlicloudEventBridgeEventBus(),
			"alicloud_amqp_virtual_host":                                     resourceAlicloudAmqpVirtualHost(),
			"alicloud_amqp_queue":                                            resourceAlicloudAmqpQueue(),
			"alicloud_amqp_exchange":                                         resourceAlicloudAmqpExchange(),
			"alicloud_cassandra_backup_plan":                                 resourceAlicloudCassandraBackupPlan(),
			"alicloud_cen_transit_router_peer_attachment":                    resourceAliCloudCenTransitRouterPeerAttachment(),
			"alicloud_amqp_instance":                                         resourceAliCloudAmqpInstance(),
			"alicloud_hbr_vault":                                             resourceAliCloudHbrVault(),
			"alicloud_ssl_certificates_service_certificate":                  resourceAliCloudSslCertificatesServiceCertificate(),
			"alicloud_arms_alert_contact":                                    resourceAlicloudArmsAlertContact(),
			"alicloud_event_bridge_slr":                                      resourceAlicloudEventBridgeServiceLinkedRole(),
			"alicloud_event_bridge_rule":                                     resourceAliCloudEventBridgeRule(),
			"alicloud_cloud_firewall_control_policy":                         resourceAliCloudCloudFirewallControlPolicy(),
			"alicloud_sae_namespace":                                         resourceAlicloudSaeNamespace(),
			"alicloud_sae_config_map":                                        resourceAlicloudSaeConfigMap(),
			"alicloud_alb_security_policy":                                   resourceAlicloudAlbSecurityPolicy(),
			"alicloud_kvstore_audit_log_config":                              resourceAlicloudKvstoreAuditLogConfig(),
			"alicloud_event_bridge_event_source":                             resourceAlicloudEventBridgeEventSource(),
			"alicloud_cloud_firewall_control_policy_order":                   resourceAliCloudCloudFirewallControlPolicyOrder(),
			"alicloud_ecd_policy_group":                                      resourceAlicloudEcdPolicyGroup(),
			"alicloud_ecp_key_pair":                                          resourceAlicloudEcpKeyPair(),
			"alicloud_hbr_ecs_backup_plan":                                   resourceAlicloudHbrEcsBackupPlan(),
			"alicloud_hbr_nas_backup_plan":                                   resourceAlicloudHbrNasBackupPlan(),
			"alicloud_hbr_oss_backup_plan":                                   resourceAlicloudHbrOssBackupPlan(),
			"alicloud_scdn_domain":                                           resourceAlicloudScdnDomain(),
			"alicloud_alb_server_group":                                      resourceAliCloudAlbServerGroup(),
			"alicloud_data_works_folder":                                     resourceAlicloudDataWorksFolder(),
			"alicloud_arms_alert_contact_group":                              resourceAlicloudArmsAlertContactGroup(),
			"alicloud_dcdn_domain_config":                                    resourceAliCloudDcdnDomainConfig(),
			"alicloud_scdn_domain_config":                                    resourceAlicloudScdnDomainConfig(),
			"alicloud_cloud_storage_gateway_gateway":                         resourceAliCloudCloudStorageGatewayGateway(),
			"alicloud_lindorm_instance":                                      resourceAliCloudLindormInstance(),
			"alicloud_cddc_dedicated_host_group":                             resourceAlicloudCddcDedicatedHostGroup(),
			"alicloud_hbr_ecs_backup_client":                                 resourceAlicloudHbrEcsBackupClient(),
			"alicloud_msc_sub_contact":                                       resourceAlicloudMscSubContact(),
			"alicloud_express_connect_physical_connection":                   resourceAliCloudExpressConnectPhysicalConnection(),
			"alicloud_alb_load_balancer":                                     resourceAliCloudAlbLoadBalancer(),
			"alicloud_sddp_rule":                                             resourceAliCloudSddpRule(),
			"alicloud_bastionhost_user_group":                                resourceAlicloudBastionhostUserGroup(),
			"alicloud_security_center_group":                                 resourceAlicloudSecurityCenterGroup(),
			"alicloud_alb_acl":                                               resourceAlicloudAlbAcl(),
			"alicloud_bastionhost_user":                                      resourceAlicloudBastionhostUser(),
			"alicloud_dfs_access_group":                                      resourceAliCloudDfsAccessGroup(),
			"alicloud_ehpc_job_template":                                     resourceAlicloudEhpcJobTemplate(),
			"alicloud_sddp_config":                                           resourceAlicloudSddpConfig(),
			"alicloud_hbr_restore_job":                                       resourceAlicloudHbrRestoreJob(),
			"alicloud_alb_listener":                                          resourceAliCloudAlbListener(),
			"alicloud_ens_key_pair":                                          resourceAlicloudEnsKeyPair(),
			"alicloud_sae_application":                                       resourceAliCloudSaeApplication(),
			"alicloud_alb_rule":                                              resourceAliCloudAlbRule(),
			"alicloud_cms_metric_rule_template":                              resourceAliCloudCmsMetricRuleTemplate(),
			"alicloud_iot_device_group":                                      resourceAlicloudIotDeviceGroup(),
			"alicloud_express_connect_virtual_border_router":                 resourceAlicloudExpressConnectVirtualBorderRouter(),
			"alicloud_imm_project":                                           resourceAlicloudImmProject(),
			"alicloud_click_house_db_cluster":                                resourceAlicloudClickHouseDbCluster(),
			"alicloud_direct_mail_domain":                                    resourceAlicloudDirectMailDomain(),
			"alicloud_bastionhost_host_group":                                resourceAlicloudBastionhostHostGroup(),
			"alicloud_vpc_dhcp_options_set":                                  resourceAliCloudVpcDhcpOptionsSet(),
			"alicloud_alb_health_check_template":                             resourceAliCloudAlbHealthCheckTemplate(),
			"alicloud_cdn_real_time_log_delivery":                            resourceAliCloudCdnRealTimeLogDelivery(),
			"alicloud_click_house_account":                                   resourceAlicloudClickHouseAccount(),
			"alicloud_selectdb_db_cluster":                                   resourceAlicloudSelectDBDbCluster(),
			"alicloud_selectdb_db_instance":                                  resourceAlicloudSelectDBDbInstance(),
			"alicloud_bastionhost_user_attachment":                           resourceAlicloudBastionhostUserAttachment(),
			"alicloud_direct_mail_mail_address":                              resourceAlicloudDirectMailMailAddress(),
			"alicloud_dts_job_monitor_rule":                                  resourceAlicloudDtsJobMonitorRule(),
			"alicloud_database_gateway_gateway":                              resourceAlicloudDatabaseGatewayGateway(),
			"alicloud_bastionhost_host":                                      resourceAlicloudBastionhostHost(),
			"alicloud_amqp_binding":                                          resourceAliCloudAmqpBinding(),
			"alicloud_slb_tls_cipher_policy":                                 resourceAlicloudSlbTlsCipherPolicy(),
			"alicloud_cloud_sso_directory":                                   resourceAlicloudCloudSsoDirectory(),
			"alicloud_bastionhost_host_account":                              resourceAlicloudBastionhostHostAccount(),
			"alicloud_bastionhost_host_attachment":                           resourceAlicloudBastionhostHostAttachment(),
			"alicloud_bastionhost_host_account_user_group_attachment":        resourceAlicloudBastionhostHostAccountUserGroupAttachment(),
			"alicloud_bastionhost_host_account_user_attachment":              resourceAlicloudBastionhostHostAccountUserAttachment(),
			"alicloud_bastionhost_host_group_account_user_attachment":        resourceAlicloudBastionhostHostGroupAccountUserAttachment(),
			"alicloud_bastionhost_host_group_account_user_group_attachment":  resourceAlicloudBastionhostHostGroupAccountUserGroupAttachment(),
			"alicloud_waf_certificate":                                       resourceAlicloudWafCertificate(),
			"alicloud_simple_application_server_instance":                    resourceAlicloudSimpleApplicationServerInstance(),
			"alicloud_video_surveillance_system_group":                       resourceAlicloudVideoSurveillanceSystemGroup(),
			"alicloud_msc_sub_subscription":                                  resourceAlicloudMscSubSubscription(),
			"alicloud_sddp_instance":                                         resourceAlicloudSddpInstance(),
			"alicloud_vpc_nat_ip_cidr":                                       resourceAlicloudVpcNatIpCidr(),
			"alicloud_vpc_nat_ip":                                            resourceAlicloudVpcNatIp(),
			"alicloud_quick_bi_user":                                         resourceAlicloudQuickBiUser(),
			"alicloud_vod_domain":                                            resourceAlicloudVodDomain(),
			"alicloud_arms_dispatch_rule":                                    resourceAlicloudArmsDispatchRule(),
			"alicloud_open_search_app_group":                                 resourceAlicloudOpenSearchAppGroup(),
			"alicloud_graph_database_db_instance":                            resourceAlicloudGraphDatabaseDbInstance(),
			"alicloud_arms_prometheus_alert_rule":                            resourceAlicloudArmsPrometheusAlertRule(),
			"alicloud_dbfs_instance":                                         resourceAliCloudDbfsDbfsInstance(),
			"alicloud_rdc_organization":                                      resourceAlicloudRdcOrganization(),
			"alicloud_eais_instance":                                         resourceAliCloudEaisInstance(),
			"alicloud_sae_ingress":                                           resourceAlicloudSaeIngress(),
			"alicloud_cloudauth_face_config":                                 resourceAlicloudCloudauthFaceConfig(),
			"alicloud_imp_app_template":                                      resourceAlicloudImpAppTemplate(),
			"alicloud_pvtz_user_vpc_authorization":                           resourceAlicloudPvtzUserVpcAuthorization(),
			"alicloud_mhub_product":                                          resourceAlicloudMhubProduct(),
			"alicloud_cloud_sso_scim_server_credential":                      resourceAlicloudCloudSsoScimServerCredential(),
			"alicloud_dts_subscription_job":                                  resourceAlicloudDtsSubscriptionJob(),
			"alicloud_service_mesh_service_mesh":                             resourceAliCloudServiceMeshServiceMesh(),
			"alicloud_mhub_app":                                              resourceAlicloudMhubApp(),
			"alicloud_cloud_sso_group":                                       resourceAlicloudCloudSsoGroup(),
			"alicloud_dts_synchronization_instance":                          resourceAlicloudDtsSynchronizationInstance(),
			"alicloud_dts_synchronization_job":                               resourceAlicloudDtsSynchronizationJob(),
			"alicloud_cloud_firewall_instance":                               resourceAliCloudCloudFirewallInstance(),
			"alicloud_cr_endpoint_acl_policy":                                resourceAlicloudCrEndpointAclPolicy(),
			"alicloud_actiontrail_history_delivery_job":                      resourceAlicloudActiontrailHistoryDeliveryJob(),
			"alicloud_ecs_deployment_set":                                    resourceAlicloudEcsDeploymentSet(),
			"alicloud_cloud_sso_user":                                        resourceAlicloudCloudSsoUser(),
			"alicloud_cloud_sso_access_configuration":                        resourceAliCloudCloudSsoAccessConfiguration(),
			"alicloud_dfs_file_system":                                       resourceAliCloudDfsFileSystem(),
			"alicloud_vpc_traffic_mirror_filter":                             resourceAliCloudVpcTrafficMirrorFilter(),
			"alicloud_dfs_access_rule":                                       resourceAliCloudDfsAccessRule(),
			"alicloud_vpc_traffic_mirror_filter_egress_rule":                 resourceAliCloudVpcTrafficMirrorFilterEgressRule(),
			"alicloud_dfs_mount_point":                                       resourceAliCloudDfsMountPoint(),
			"alicloud_ecd_simple_office_site":                                resourceAlicloudEcdSimpleOfficeSite(),
			"alicloud_vpc_traffic_mirror_filter_ingress_rule":                resourceAliCloudVpcTrafficMirrorFilterIngressRule(),
			"alicloud_ecd_nas_file_system":                                   resourceAlicloudEcdNasFileSystem(),
			"alicloud_cloud_sso_user_attachment":                             resourceAlicloudCloudSsoUserAttachment(),
			"alicloud_cloud_sso_access_assignment":                           resourceAlicloudCloudSsoAccessAssignment(),
			"alicloud_msc_sub_webhook":                                       resourceAlicloudMscSubWebhook(),
			"alicloud_waf_protection_module":                                 resourceAlicloudWafProtectionModule(),
			"alicloud_ecd_user":                                              resourceAlicloudEcdUser(),
			"alicloud_vpc_traffic_mirror_session":                            resourceAliCloudVpcTrafficMirrorSession(),
			"alicloud_gpdb_account":                                          resourceAliCloudGpdbAccount(),
			"alicloud_security_center_service_linked_role":                   resourceAlicloudSecurityCenterServiceLinkedRole(),
			"alicloud_event_bridge_service_linked_role":                      resourceAlicloudEventBridgeServiceLinkedRole(),
			"alicloud_vpc_ipv6_gateway":                                      resourceAliCloudVpcIpv6Gateway(),
			"alicloud_vpc_ipv6_egress_rule":                                  resourceAliCloudVpcIpv6EgressRule(),
			"alicloud_hbr_server_backup_plan":                                resourceAlicloudHbrServerBackupPlan(),
			"alicloud_cms_dynamic_tag_group":                                 resourceAliCloudCmsDynamicTagGroup(),
			"alicloud_ecd_network_package":                                   resourceAlicloudEcdNetworkPackage(),
			"alicloud_cloud_storage_gateway_gateway_smb_user":                resourceAlicloudCloudStorageGatewayGatewaySmbUser(),
			"alicloud_vpc_ipv6_internet_bandwidth":                           resourceAliCloudVpcIpv6InternetBandwidth(),
			"alicloud_simple_application_server_firewall_rule":               resourceAlicloudSimpleApplicationServerFirewallRule(),
			"alicloud_pvtz_endpoint":                                         resourceAlicloudPvtzEndpoint(),
			"alicloud_pvtz_rule":                                             resourceAlicloudPvtzRule(),
			"alicloud_pvtz_rule_attachment":                                  resourceAlicloudPvtzRuleAttachment(),
			"alicloud_simple_application_server_snapshot":                    resourceAlicloudSimpleApplicationServerSnapshot(),
			"alicloud_simple_application_server_custom_image":                resourceAlicloudSimpleApplicationServerCustomImage(),
			"alicloud_cloud_storage_gateway_gateway_cache_disk":              resourceAliCloudCloudStorageGatewayGatewayCacheDisk(),
			"alicloud_cloud_storage_gateway_gateway_logging":                 resourceAlicloudCloudStorageGatewayGatewayLogging(),
			"alicloud_cloud_storage_gateway_gateway_block_volume":            resourceAlicloudCloudStorageGatewayGatewayBlockVolume(),
			"alicloud_direct_mail_tag":                                       resourceAlicloudDirectMailTag(),
			"alicloud_cloud_storage_gateway_gateway_file_share":              resourceAlicloudCloudStorageGatewayGatewayFileShare(),
			"alicloud_ecd_desktop":                                           resourceAlicloudEcdDesktop(),
			"alicloud_cloud_storage_gateway_express_sync":                    resourceAlicloudCloudStorageGatewayExpressSync(),
			"alicloud_cloud_storage_gateway_express_sync_share_attachment":   resourceAlicloudCloudStorageGatewayExpressSyncShareAttachment(),
			"alicloud_oos_application":                                       resourceAlicloudOosApplication(),
			"alicloud_eci_virtual_node":                                      resourceAlicloudEciVirtualNode(),
			"alicloud_ros_stack_instance":                                    resourceAlicloudRosStackInstance(),
			"alicloud_ecs_dedicated_host_cluster":                            resourceAlicloudEcsDedicatedHostCluster(),
			"alicloud_oos_application_group":                                 resourceAlicloudOosApplicationGroup(),
			"alicloud_dts_consumer_channel":                                  resourceAlicloudDtsConsumerChannel(),
			"alicloud_ecd_image":                                             resourceAlicloudEcdImage(),
			"alicloud_oos_patch_baseline":                                    resourceAliCloudOosPatchBaseline(),
			"alicloud_ecd_command":                                           resourceAlicloudEcdCommand(),
			"alicloud_cddc_dedicated_host":                                   resourceAlicloudCddcDedicatedHost(),
			"alicloud_oos_service_setting":                                   resourceAlicloudOosServiceSetting(),
			"alicloud_oos_parameter":                                         resourceAlicloudOosParameter(),
			"alicloud_oos_state_configuration":                               resourceAlicloudOosStateConfiguration(),
			"alicloud_oos_secret_parameter":                                  resourceAlicloudOosSecretParameter(),
			"alicloud_click_house_backup_policy":                             resourceAlicloudClickHouseBackupPolicy(),
			"alicloud_mongodb_audit_policy":                                  resourceAlicloudMongodbAuditPolicy(),
			"alicloud_cloud_sso_access_configuration_provisioning":           resourceAlicloudCloudSsoAccessConfigurationProvisioning(),
			"alicloud_mongodb_account":                                       resourceAlicloudMongodbAccount(),
			"alicloud_mongodb_serverless_instance":                           resourceAlicloudMongodbServerlessInstance(),
			"alicloud_ecs_session_manager_status":                            resourceAlicloudEcsSessionManagerStatus(),
			"alicloud_cddc_dedicated_host_account":                           resourceAlicloudCddcDedicatedHostAccount(),
			"alicloud_cr_chart_namespace":                                    resourceAlicloudCrChartNamespace(),
			"alicloud_fnf_execution":                                         resourceAlicloudFnFExecution(),
			"alicloud_cr_chart_repository":                                   resourceAlicloudCrChartRepository(),
			"alicloud_mongodb_sharding_network_public_address":               resourceAlicloudMongodbShardingNetworkPublicAddress(),
			"alicloud_ga_acl":                                                resourceAliCloudGaAcl(),
			"alicloud_ga_acl_attachment":                                     resourceAliCloudGaAclAttachment(),
			"alicloud_ga_additional_certificate":                             resourceAliCloudGaAdditionalCertificate(),
			"alicloud_alidns_custom_line":                                    resourceAlicloudAlidnsCustomLine(),
			"alicloud_vpc_vbr_ha":                                            resourceAlicloudVpcVbrHa(),
			"alicloud_ros_template_scratch":                                  resourceAlicloudRosTemplateScratch(),
			"alicloud_alidns_gtm_instance":                                   resourceAlicloudAlidnsGtmInstance(),
			"alicloud_vpc_bgp_group":                                         resourceAlicloudVpcBgpGroup(),
			"alicloud_ram_security_preference":                               resourceAlicloudRamSecurityPreference(),
			"alicloud_nas_snapshot":                                          resourceAlicloudNasSnapshot(),
			"alicloud_hbr_replication_vault":                                 resourceAlicloudHbrReplicationVault(),
			"alicloud_alidns_address_pool":                                   resourceAlicloudAlidnsAddressPool(),
			"alicloud_ecs_prefix_list":                                       resourceAlicloudEcsPrefixList(),
			"alicloud_alidns_access_strategy":                                resourceAlicloudAlidnsAccessStrategy(),
			"alicloud_alidns_monitor_config":                                 resourceAlicloudAlidnsMonitorConfig(),
			"alicloud_vpc_dhcp_options_set_attachment":                       resourceAlicloudVpcDhcpOptionsSetAttachement(),
			"alicloud_vpc_bgp_peer":                                          resourceAliCloudExpressConnectBgpPeer(),
			"alicloud_nas_fileset":                                           resourceAlicloudNasFileset(),
			"alicloud_nas_auto_snapshot_policy":                              resourceAliCloudNasAutoSnapshotPolicy(),
			"alicloud_nas_lifecycle_policy":                                  resourceAlicloudNasLifecyclePolicy(),
			"alicloud_vpc_bgp_network":                                       resourceAlicloudVpcBgpNetwork(),
			"alicloud_nas_data_flow":                                         resourceAlicloudNasDataFlow(),
			"alicloud_ecs_storage_capacity_unit":                             resourceAlicloudEcsStorageCapacityUnit(),
			"alicloud_nas_recycle_bin":                                       resourceAlicloudNasRecycleBin(),
			"alicloud_dbfs_snapshot":                                         resourceAliCloudDbfsSnapshot(),
			"alicloud_dbfs_instance_attachment":                              resourceAliCloudDbfsInstanceAttachment(),
			"alicloud_dts_migration_job":                                     resourceAlicloudDtsMigrationJob(),
			"alicloud_dts_migration_instance":                                resourceAlicloudDtsMigrationInstance(),
			"alicloud_mse_gateway":                                           resourceAlicloudMseGateway(),
			"alicloud_dbfs_service_linked_role":                              resourceAlicloudDbfsServiceLinkedRole(),
			"alicloud_resource_manager_service_linked_role":                  resourceAliCloudResourceManagerServiceLinkedRole(),
			"alicloud_rds_service_linked_role":                               resourceAlicloudRdsServiceLinkedRole(),
			"alicloud_mongodb_sharding_network_private_address":              resourceAliCloudMongodbShardingNetworkPrivateAddress(),
			"alicloud_ecp_instance":                                          resourceAliCloudEcpInstance(),
			"alicloud_dcdn_ipa_domain":                                       resourceAlicloudDcdnIpaDomain(),
			"alicloud_sddp_data_limit":                                       resourceAlicloudSddpDataLimit(),
			"alicloud_ecs_image_component":                                   resourceAliCloudEcsImageComponent(),
			"alicloud_sae_application_scaling_rule":                          resourceAlicloudSaeApplicationScalingRule(),
			"alicloud_sae_grey_tag_route":                                    resourceAliCloudSaeGreyTagRoute(),
			"alicloud_ecs_snapshot_group":                                    resourceAlicloudEcsSnapshotGroup(),
			"alicloud_alb_listener_additional_certificate_attachment":        resourceAlicloudAlbListenerAdditionalCertificateAttachment(),
			"alicloud_vpn_ipsec_server":                                      resourceAlicloudVpnIpsecServer(),
			"alicloud_cr_chain":                                              resourceAlicloudCrChain(),
			"alicloud_vpn_pbr_route_entry":                                   resourceAlicloudVpnPbrRouteEntry(),
			"alicloud_slb_acl_entry_attachment":                              resourceAlicloudSlbAclEntryAttachment(),
			"alicloud_mse_znode":                                             resourceAlicloudMseZnode(),
			"alicloud_alikafka_instance_allowed_ip_attachment":               resourceAliCloudAliKafkaInstanceAllowedIpAttachment(),
			"alicloud_ecs_image_pipeline":                                    resourceAlicloudEcsImagePipeline(),
			"alicloud_slb_server_group_server_attachment":                    resourceAlicloudSlbServerGroupServerAttachment(),
			"alicloud_alb_listener_acl_attachment":                           resourceAliCloudAlbListenerAclAttachment(),
			"alicloud_hbr_ots_backup_plan":                                   resourceAlicloudHbrOtsBackupPlan(),
			"alicloud_sae_load_balancer_internet":                            resourceAlicloudSaeLoadBalancerInternet(),
			"alicloud_bastionhost_host_share_key":                            resourceAlicloudBastionhostHostShareKey(),
			"alicloud_cdn_fc_trigger":                                        resourceAlicloudCdnFcTrigger(),
			"alicloud_sae_load_balancer_intranet":                            resourceAlicloudSaeLoadBalancerIntranet(),
			"alicloud_bastionhost_host_account_share_key_attachment":         resourceAlicloudBastionhostHostAccountShareKeyAttachment(),
			"alicloud_alb_acl_entry_attachment":                              resourceAlicloudAlbAclEntryAttachment(),
			"alicloud_ecs_network_interface_permission":                      resourceAlicloudEcsNetworkInterfacePermission(),
			"alicloud_mse_engine_namespace":                                  resourceAlicloudMseEngineNamespace(),
			"alicloud_mse_nacos_config":                                      resourceAlicloudMseNacosConfig(),
			"alicloud_ga_accelerator_spare_ip_attachment":                    resourceAlicloudGaAcceleratorSpareIpAttachment(),
			"alicloud_smartag_flow_log":                                      resourceAlicloudSmartagFlowLog(),
			"alicloud_ecs_invocation":                                        resourceAlicloudEcsInvocation(),
			"alicloud_ddos_basic_defense_threshold":                          resourceAlicloudDdosBasicDefenseThreshold(),
			"alicloud_ecd_snapshot":                                          resourceAlicloudEcdSnapshot(),
			"alicloud_ecd_bundle":                                            resourceAlicloudEcdBundle(),
			"alicloud_config_delivery":                                       resourceAliCloudConfigDelivery(),
			"alicloud_cms_namespace":                                         resourceAliCloudCmsNamespace(),
			"alicloud_cms_sls_group":                                         resourceAlicloudCmsSlsGroup(),
			"alicloud_config_aggregate_delivery":                             resourceAliCloudConfigAggregateDelivery(),
			"alicloud_edas_namespace":                                        resourceAlicloudEdasNamespace(),
			"alicloud_schedulerx_namespace":                                  resourceAlicloudSchedulerxNamespace(),
			"alicloud_ehpc_cluster":                                          resourceAlicloudEhpcCluster(),
			"alicloud_cen_traffic_marking_policy":                            resourceAliCloudCenTrafficMarkingPolicy(),
			"alicloud_ecs_instance_set":                                      resourceAlicloudEcsInstanceSet(),
			"alicloud_ecd_ram_directory":                                     resourceAlicloudEcdRamDirectory(),
			"alicloud_service_mesh_user_permission":                          resourceAliCloudServiceMeshUserPermission(),
			"alicloud_ecd_ad_connector_directory":                            resourceAlicloudEcdAdConnectorDirectory(),
			"alicloud_ecd_custom_property":                                   resourceAlicloudEcdCustomProperty(),
			"alicloud_ecd_ad_connector_office_site":                          resourceAlicloudEcdAdConnectorOfficeSite(),
			"alicloud_ecs_activation":                                        resourceAlicloudEcsActivation(),
			"alicloud_cloud_firewall_address_book":                           resourceAliCloudCloudFirewallAddressBook(),
			"alicloud_sms_short_url":                                         resourceAlicloudSmsShortUrl(),
			"alicloud_hbr_hana_instance":                                     resourceAlicloudHbrHanaInstance(),
			"alicloud_cms_hybrid_monitor_sls_task":                           resourceAlicloudCmsHybridMonitorSlsTask(),
			"alicloud_hbr_hana_backup_plan":                                  resourceAlicloudHbrHanaBackupPlan(),
			"alicloud_cms_hybrid_monitor_fc_task":                            resourceAlicloudCmsHybridMonitorFcTask(),
			"alicloud_fc_layer_version":                                      resourceAlicloudFcLayerVersion(),
			"alicloud_ddosbgp_ip":                                            resourceAlicloudDdosbgpIp(),
			"alicloud_vpn_gateway_vpn_attachment":                            resourceAlicloudVpnGatewayVpnAttachment(),
			"alicloud_resource_manager_delegated_administrator":              resourceAlicloudResourceManagerDelegatedAdministrator(),
			"alicloud_polardb_global_database_network":                       resourceAlicloudPolarDBGlobalDatabaseNetwork(),
			"alicloud_vpc_ipv4_gateway":                                      resourceAliCloudVpcIpv4Gateway(),
			"alicloud_api_gateway_backend":                                   resourceAlicloudApiGatewayBackend(),
			"alicloud_vpc_prefix_list":                                       resourceAliCloudVpcPrefixList(),
			"alicloud_cms_event_rule":                                        resourceAliCloudCloudMonitorServiceEventRule(),
			"alicloud_ddos_basic_threshold":                                  resourceAliCloudDdosBasicThreshold(),
			"alicloud_cen_transit_router_vpn_attachment":                     resourceAlicloudCenTransitRouterVpnAttachment(),
			"alicloud_polardb_parameter_group":                               resourceAlicloudPolarDBParameterGroup(),
			"alicloud_vpn_gateway_vco_route":                                 resourceAlicloudVpnGatewayVcoRoute(),
			"alicloud_dcdn_waf_policy":                                       resourceAlicloudDcdnWafPolicy(),
			"alicloud_api_gateway_log_config":                                resourceAlicloudApiGatewayLogConfig(),
			"alicloud_dbs_backup_plan":                                       resourceAlicloudDbsBackupPlan(),
			"alicloud_dcdn_waf_domain":                                       resourceAlicloudDcdnWafDomain(),
			"alicloud_vpc_ipv4_cidr_block":                                   resourceAliCloudVpcIpv4CidrBlock(),
			"alicloud_vpc_public_ip_address_pool":                            resourceAliCloudVpcPublicIpAddressPool(),
			"alicloud_dcdn_waf_policy_domain_attachment":                     resourceAlicloudDcdnWafPolicyDomainAttachment(),
			"alicloud_nlb_server_group":                                      resourceAliCloudNlbServerGroup(),
			"alicloud_vpc_peer_connection":                                   resourceAliCloudVpcPeerPeerConnection(),
			"alicloud_ga_access_log":                                         resourceAlicloudGaAccessLog(),
			"alicloud_ebs_disk_replica_group":                                resourceAlicloudEbsDiskReplicaGroup(),
			"alicloud_nlb_security_policy":                                   resourceAliCloudNlbSecurityPolicy(),
			"alicloud_vod_editing_project":                                   resourceAlicloudVodEditingProject(),
			"alicloud_api_gateway_model":                                     resourceAlicloudApiGatewayModel(),
			"alicloud_cen_transit_router_grant_attachment":                   resourceAlicloudCenTransitRouterGrantAttachment(),
			"alicloud_api_gateway_plugin":                                    resourceAliCloudApiGatewayPlugin(),
			"alicloud_api_gateway_plugin_attachment":                         resourceAlicloudApiGatewayPluginAttachment(),
			"alicloud_message_service_queue":                                 resourceAliCloudMessageServiceQueue(),
			"alicloud_message_service_topic":                                 resourceAlicloudMessageServiceTopic(),
			"alicloud_message_service_subscription":                          resourceAlicloudMessageServiceSubscription(),
			"alicloud_cen_transit_router_prefix_list_association":            resourceAlicloudCenTransitRouterPrefixListAssociation(),
			"alicloud_dms_enterprise_proxy":                                  resourceAlicloudDmsEnterpriseProxy(),
			"alicloud_vpc_public_ip_address_pool_cidr_block":                 resourceAliCloudVpcPublicIpAddressPoolCidrBlock(),
			"alicloud_gpdb_db_instance_plan":                                 resourceAliCloudGpdbDbInstancePlan(),
			"alicloud_adb_db_cluster_lake_version":                           resourceAliCloudAdbDbClusterLakeVersion(),
			"alicloud_ga_acl_entry_attachment":                               resourceAlicloudGaAclEntryAttachment(),
			"alicloud_nlb_load_balancer":                                     resourceAliCloudNlbLoadBalancer(),
			"alicloud_service_mesh_extension_provider":                       resourceAlicloudServiceMeshExtensionProvider(),
			"alicloud_nlb_listener":                                          resourceAliCloudNlbListener(),
			"alicloud_nlb_server_group_server_attachment":                    resourceAliCloudNlbServerGroupServerAttachment(),
			"alicloud_bp_studio_application":                                 resourceAlicloudBpStudioApplication(),
			"alicloud_vpc_network_acl_attachment":                            resourceAliCloudVpcNetworkAclAttachment(),
			"alicloud_cen_transit_router_cidr":                               resourceAlicloudCenTransitRouterCidr(),
			"alicloud_das_switch_das_pro":                                    resourceAlicloudDasSwitchDasPro(),
			"alicloud_ga_basic_accelerator":                                  resourceAliCloudGaBasicAccelerator(),
			"alicloud_ga_basic_endpoint_group":                               resourceAliCloudGaBasicEndpointGroup(),
			"alicloud_cms_metric_rule_black_list":                            resourceAlicloudCmsMetricRuleBlackList(),
			"alicloud_ga_basic_ip_set":                                       resourceAliCloudGaBasicIpSet(),
			"alicloud_cloud_firewall_vpc_firewall_cen":                       resourceAliCloudCloudFirewallVpcFirewallCen(),
			"alicloud_cloud_firewall_vpc_firewall":                           resourceAliCloudCloudFirewallVpcFirewall(),
			"alicloud_cloud_firewall_instance_member":                        resourceAlicloudCloudFirewallInstanceMember(),
			"alicloud_ga_basic_accelerate_ip":                                resourceAliCloudGaBasicAccelerateIp(),
			"alicloud_ga_basic_endpoint":                                     resourceAliCloudGaBasicEndpoint(),
			"alicloud_cloud_firewall_vpc_firewall_control_policy":            resourceAliCloudCloudFirewallVpcFirewallControlPolicy(),
			"alicloud_ga_basic_accelerate_ip_endpoint_relation":              resourceAlicloudGaBasicAccelerateIpEndpointRelation(),
			"alicloud_vpc_gateway_route_table_attachment":                    resourceAliCloudVpcGatewayRouteTableAttachment(),
			"alicloud_threat_detection_web_lock_config":                      resourceAlicloudThreatDetectionWebLockConfig(),
			"alicloud_threat_detection_backup_policy":                        resourceAlicloudThreatDetectionBackupPolicy(),
			"alicloud_dms_enterprise_proxy_access":                           resourceAlicloudDmsEnterpriseProxyAccess(),
			"alicloud_threat_detection_vul_whitelist":                        resourceAlicloudThreatDetectionVulWhitelist(),
			"alicloud_dms_enterprise_logic_database":                         resourceAlicloudDmsEnterpriseLogicDatabase(),
			"alicloud_amqp_static_account":                                   resourceAliCloudAmqpStaticAccount(),
			"alicloud_adb_resource_group":                                    resourceAliCloudAdbResourceGroup(),
			"alicloud_alb_ascript":                                           resourceAlicloudAlbAscript(),
			"alicloud_threat_detection_honeypot_node":                        resourceAlicloudThreatDetectionHoneypotNode(),
			"alicloud_cen_transit_router_multicast_domain":                   resourceAliCloudCenTransitRouterMulticastDomain(),
			"alicloud_cen_transit_router_multicast_domain_source":            resourceAlicloudCenTransitRouterMulticastDomainSource(),
			"alicloud_cen_inter_region_traffic_qos_policy":                   resourceAlicloudCenInterRegionTrafficQosPolicy(),
			"alicloud_threat_detection_baseline_strategy":                    resourceAlicloudThreatDetectionBaselineStrategy(),
			"alicloud_threat_detection_anti_brute_force_rule":                resourceAlicloudThreatDetectionAntiBruteForceRule(),
			"alicloud_threat_detection_honey_pot":                            resourceAlicloudThreatDetectionHoneyPot(),
			"alicloud_threat_detection_honeypot_probe":                       resourceAlicloudThreatDetectionHoneypotProbe(),
			"alicloud_ecs_capacity_reservation":                              resourceAlicloudEcsCapacityReservation(),
			"alicloud_cen_inter_region_traffic_qos_queue":                    resourceAlicloudCenInterRegionTrafficQosQueue(),
			"alicloud_cen_transit_router_multicast_domain_peer_member":       resourceAlicloudCenTransitRouterMulticastDomainPeerMember(),
			"alicloud_cen_transit_router_multicast_domain_member":            resourceAliCloudCenTransitRouterMulticastDomainMember(),
			"alicloud_cen_child_instance_route_entry_to_attachment":          resourceAlicloudCenChildInstanceRouteEntryToAttachment(),
			"alicloud_cen_transit_router_multicast_domain_association":       resourceAliCloudCenTransitRouterMulticastDomainAssociation(),
			"alicloud_threat_detection_honeypot_preset":                      resourceAlicloudThreatDetectionHoneypotPreset(),
			"alicloud_service_catalog_provisioned_product":                   resourceAlicloudServiceCatalogProvisionedProduct(),
			"alicloud_vpc_peer_connection_accepter":                          resourceAliCloudVpcPeerPeerConnectionAccepter(),
			"alicloud_ebs_dedicated_block_storage_cluster":                   resourceAlicloudEbsDedicatedBlockStorageCluster(),
			"alicloud_ecs_elasticity_assurance":                              resourceAlicloudEcsElasticityAssurance(),
			"alicloud_express_connect_grant_rule_to_cen":                     resourceAlicloudExpressConnectGrantRuleToCen(),
			"alicloud_express_connect_virtual_physical_connection":           resourceAlicloudExpressConnectVirtualPhysicalConnection(),
			"alicloud_express_connect_vbr_pconn_association":                 resourceAlicloudExpressConnectVbrPconnAssociation(),
			"alicloud_ebs_disk_replica_pair":                                 resourceAlicloudEbsDiskReplicaPair(),
			"alicloud_ga_domain":                                             resourceAlicloudGaDomain(),
			"alicloud_ga_custom_routing_endpoint_group":                      resourceAliCloudGaCustomRoutingEndpointGroup(),
			"alicloud_ga_custom_routing_endpoint_group_destination":          resourceAliCloudGaCustomRoutingEndpointGroupDestination(),
			"alicloud_ga_custom_routing_endpoint":                            resourceAliCloudGaCustomRoutingEndpoint(),
			"alicloud_ga_custom_routing_endpoint_traffic_policy":             resourceAliCloudGaCustomRoutingEndpointTrafficPolicy(),
			"alicloud_nlb_load_balancer_security_group_attachment":           resourceAliCloudNlbLoadBalancerSecurityGroupAttachment(),
			"alicloud_dcdn_kv_namespace":                                     resourceAlicloudDcdnKvNamespace(),
			"alicloud_dcdn_kv":                                               resourceAlicloudDcdnKv(),
			"alicloud_hbr_hana_backup_client":                                resourceAlicloudHbrHanaBackupClient(),
			"alicloud_dts_instance":                                          resourceAlicloudDtsInstance(),
			"alicloud_threat_detection_instance":                             resourceAliCloudThreatDetectionInstance(),
			"alicloud_cr_vpc_endpoint_linked_vpc":                            resourceAlicloudCrVpcEndpointLinkedVpc(),
			"alicloud_express_connect_router_interface":                      resourceAlicloudExpressConnectRouterInterface(),
			"alicloud_wafv3_instance":                                        resourceAlicloudWafv3Instance(),
			"alicloud_alb_load_balancer_common_bandwidth_package_attachment": resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachment(),
			"alicloud_wafv3_domain":                                          resourceAlicloudWafv3Domain(),
			"alicloud_eflo_vpd":                                              resourceAlicloudEfloVpd(),
			"alicloud_dcdn_waf_rule":                                         resourceAlicloudDcdnWafRule(),
			"alicloud_dcdn_er":                                               resourceAlicloudDcdnEr(),
			"alicloud_actiontrail_global_events_storage_region":              resourceAlicloudActiontrailGlobalEventsStorageRegion(),
			"alicloud_dbfs_auto_snap_shot_policy":                            resourceAlicloudDbfsAutoSnapShotPolicy(),
			"alicloud_cen_transit_route_table_aggregation":                   resourceAlicloudCenTransitRouteTableAggregation(),
			"alicloud_arms_prometheus":                                       resourceAlicloudArmsPrometheus(),
			"alicloud_oos_default_patch_baseline":                            resourceAlicloudOosDefaultPatchBaseline(),
			"alicloud_ocean_base_instance":                                   resourceAliCloudOceanBaseInstance(),
			"alicloud_chatbot_publish_task":                                  resourceAlicloudChatbotPublishTask(),
			"alicloud_arms_integration_exporter":                             resourceAlicloudArmsIntegrationExporter(),
			"alicloud_service_catalog_portfolio":                             resourceAliCloudServiceCatalogPortfolio(),
			"alicloud_arms_remote_write":                                     resourceAliCloudArmsRemoteWrite(),
			"alicloud_eflo_subnet":                                           resourceAlicloudEfloSubnet(),
			"alicloud_compute_nest_service_instance":                         resourceAlicloudComputeNestServiceInstance(),
			"alicloud_cloud_monitor_service_hybrid_double_write":             resourceAliCloudCloudMonitorServiceHybridDoubleWrite(),
			"alicloud_event_bridge_connection":                               resourceAliCloudEventBridgeConnection(),
			"alicloud_event_bridge_api_destination":                          resourceAliCloudEventBridgeApiDestination(),
			"alicloud_cloud_monitor_service_monitoring_agent_process":        resourceAliCloudCloudMonitorServiceMonitoringAgentProcess(),
			"alicloud_cloud_monitor_service_group_monitoring_agent_process":  resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcess(),
		},
	}
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return providerConfigure(d, provider)
	}
	return provider
}

var providerConfig map[string]interface{}

func providerConfigure(d *schema.ResourceData, p *schema.Provider) (interface{}, error) {
	log.Println("using terraform version:", p.TerraformVersion)
	var getProviderConfig = func(schemaKey string, profileKey string) string {
		if schemaKey != "" {
			if v, ok := d.GetOk(schemaKey); ok && v != nil && v.(string) != "" {
				return v.(string)
			}
		}
		if v, err := getConfigFromProfile(d, profileKey); err == nil && v != nil {
			return v.(string)
		}
		return ""
	}

	accessKey := getProviderConfig("access_key", "access_key_id")
	secretKey := getProviderConfig("secret_key", "access_key_secret")
	region := getProviderConfig("region", "region_id")
	if region == "" {
		region = DEFAULT_REGION
	}
	securityToken := getProviderConfig("security_token", "sts_token")

	ecsRoleName := getProviderConfig("ecs_role_name", "ram_role_name")

	if accessKey == "" || secretKey == "" {
		if v, ok := d.GetOk("credentials_uri"); ok && v.(string) != "" {
			credentialsURIResp, err := getClientByCredentialsURI(v.(string))
			if err != nil {
				return nil, err
			}
			accessKey = credentialsURIResp.AccessKeyId
			secretKey = credentialsURIResp.AccessKeySecret
			securityToken = credentialsURIResp.SecurityToken
		}
	}

	config := &connectivity.Config{
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		EcsRoleName:          strings.TrimSpace(ecsRoleName),
		Region:               connectivity.Region(strings.TrimSpace(region)),
		RegionId:             strings.TrimSpace(region),
		SkipRegionValidation: d.Get("skip_region_validation").(bool),
		ConfigurationSource:  d.Get("configuration_source").(string),
		Protocol:             d.Get("protocol").(string),
		ClientReadTimeout:    d.Get("client_read_timeout").(int),
		ClientConnectTimeout: d.Get("client_connect_timeout").(int),
		SourceIp:             strings.TrimSpace(d.Get("source_ip").(string)),
		SecureTransport:      strings.TrimSpace(d.Get("secure_transport").(string)),
		MaxRetryTimeout:      d.Get("max_retry_timeout").(int),
		TerraformTraceId:     strings.Trim(uuid.New().String(), "-"),
		TerraformVersion:     p.TerraformVersion,
	}
	log.Println("alicloud provider trace id:", config.TerraformTraceId)
	if accessKey != "" && secretKey != "" {
		credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
		if v := strings.TrimSpace(securityToken); v != "" {
			credentialConfig.SetType("sts").SetSecurityToken(v)
		}
		credential, err := credentials.NewCredential(credentialConfig)
		if err != nil {
			return nil, err
		}
		config.Credential = credential
	}
	if v, ok := d.GetOk("security_transport"); config.SecureTransport == "" && ok && v.(string) != "" {
		config.SecureTransport = v.(string)
	}
	config.SecurityToken = strings.TrimSpace(securityToken)

	config.RamRoleArn = getProviderConfig("", "ram_role_arn")
	config.RamRoleSessionName = getProviderConfig("", "ram_session_name")
	expiredSeconds, err := getConfigFromProfile(d, "expired_seconds")
	if err == nil && expiredSeconds != nil {
		config.RamRoleSessionExpiration = (int)(expiredSeconds.(float64))
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		if assumeRole["role_arn"].(string) != "" {
			config.RamRoleArn = assumeRole["role_arn"].(string)
		}
		if assumeRole["session_name"].(string) != "" {
			config.RamRoleSessionName = assumeRole["session_name"].(string)
		}
		if config.RamRoleSessionName == "" {
			config.RamRoleSessionName = "terraform"
		}
		config.RamRolePolicy = assumeRole["policy"].(string)
		if assumeRole["session_expiration"].(int) == 0 {
			if v := os.Getenv("ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION"); v != "" {
				if expiredSeconds, err := strconv.Atoi(v); err == nil {
					config.RamRoleSessionExpiration = expiredSeconds
				}
			}
			if config.RamRoleSessionExpiration == 0 {
				config.RamRoleSessionExpiration = 3600
			}
		} else {
			config.RamRoleSessionExpiration = assumeRole["session_expiration"].(int)
		}
		if v := assumeRole["external_id"].(string); v != "" {
			config.RamRoleExternalId = v
		}

		log.Printf("[INFO] assume_role configuration set: (RamRoleArn: %q, RamRoleSessionName: %q, RamRolePolicy: %q, RamRoleSessionExpiration: %d, RamRoleExternalId: %s)",
			config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration, config.RamRoleExternalId)
	}

	if v, ok := d.GetOk("assume_role_with_oidc"); ok && len(v.([]interface{})) == 1 {
		config.AssumeRoleWithOidc, err = getAssumeRoleWithOIDCConfig(v.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		log.Printf("[INFO] assume_role_with_oidc configuration set: (RoleArn: %q, SessionName: %q, SessionExpiration: %d, OIDCProviderArn: %s)",
			config.AssumeRoleWithOidc.RoleARN, config.AssumeRoleWithOidc.RoleSessionName, config.AssumeRoleWithOidc.DurationSeconds, config.AssumeRoleWithOidc.OIDCProviderArn)
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)
	var endpointInit sync.Map
	config.Endpoints = &endpointInit

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		for key, val := range endpoints {
			endpointInit.Store(key, val)
		}
		config.EcsEndpoint = strings.TrimSpace(endpoints["ecs"].(string))
		config.RdsEndpoint = strings.TrimSpace(endpoints["rds"].(string))
		config.SlbEndpoint = strings.TrimSpace(endpoints["slb"].(string))
		config.VpcEndpoint = strings.TrimSpace(endpoints["vpc"].(string))
		config.EssEndpoint = strings.TrimSpace(endpoints["ess"].(string))
		config.OssEndpoint = strings.TrimSpace(endpoints["oss"].(string))
		config.OnsEndpoint = strings.TrimSpace(endpoints["ons"].(string))
		config.AlikafkaEndpoint = strings.TrimSpace(endpoints["alikafka"].(string))
		config.DnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
		config.RamEndpoint = strings.TrimSpace(endpoints["ram"].(string))
		config.CsEndpoint = strings.TrimSpace(endpoints["cs"].(string))
		config.CrEndpoint = strings.TrimSpace(endpoints["cr"].(string))
		config.CdnEndpoint = strings.TrimSpace(endpoints["cdn"].(string))
		config.KmsEndpoint = strings.TrimSpace(endpoints["kms"].(string))
		config.OtsEndpoint = strings.TrimSpace(endpoints["ots"].(string))
		config.CmsEndpoint = strings.TrimSpace(endpoints["cms"].(string))
		config.PvtzEndpoint = strings.TrimSpace(endpoints["pvtz"].(string))
		config.StsEndpoint = strings.TrimSpace(endpoints["sts"].(string))
		config.LogEndpoint = strings.TrimSpace(endpoints["log"].(string))
		config.DrdsEndpoint = strings.TrimSpace(endpoints["drds"].(string))
		config.DdsEndpoint = strings.TrimSpace(endpoints["dds"].(string))
		config.GpdbEnpoint = strings.TrimSpace(endpoints["gpdb"].(string))
		config.KVStoreEndpoint = strings.TrimSpace(endpoints["kvstore"].(string))
		config.PolarDBEndpoint = strings.TrimSpace(endpoints["polardb"].(string))
		config.FcEndpoint = strings.TrimSpace(endpoints["fc"].(string))
		config.ApigatewayEndpoint = strings.TrimSpace(endpoints["apigateway"].(string))
		config.DatahubEndpoint = strings.TrimSpace(endpoints["datahub"].(string))
		config.MnsEndpoint = strings.TrimSpace(endpoints["mns"].(string))
		config.LocationEndpoint = strings.TrimSpace(endpoints["location"].(string))
		config.ElasticsearchEndpoint = strings.TrimSpace(endpoints["elasticsearch"].(string))
		config.NasEndpoint = strings.TrimSpace(endpoints["nas"].(string))
		config.ActiontrailEndpoint = strings.TrimSpace(endpoints["actiontrail"].(string))
		config.BssOpenApiEndpoint = strings.TrimSpace(endpoints["bssopenapi"].(string))
		config.DdoscooEndpoint = strings.TrimSpace(endpoints["ddoscoo"].(string))
		config.DdosbgpEndpoint = strings.TrimSpace(endpoints["ddosbgp"].(string))
		config.EmrEndpoint = strings.TrimSpace(endpoints["emr"].(string))
		config.CasEndpoint = strings.TrimSpace(endpoints["cas"].(string))
		config.MarketEndpoint = strings.TrimSpace(endpoints["market"].(string))
		config.AdbEndpoint = strings.TrimSpace(endpoints["adb"].(string))
		config.CbnEndpoint = strings.TrimSpace(endpoints["cbn"].(string))
		config.MaxComputeEndpoint = strings.TrimSpace(endpoints["maxcompute"].(string))
		config.DmsEnterpriseEndpoint = strings.TrimSpace(endpoints["dms_enterprise"].(string))
		config.WafOpenapiEndpoint = strings.TrimSpace(endpoints["waf_openapi"].(string))
		config.ResourcemanagerEndpoint = strings.TrimSpace(endpoints["resourcemanager"].(string))
		config.EciEndpoint = strings.TrimSpace(endpoints["eci"].(string))
		config.OosEndpoint = strings.TrimSpace(endpoints["oos"].(string))
		config.DcdnEndpoint = strings.TrimSpace(endpoints["dcdn"].(string))
		config.Endpoints.Store("mse", strings.TrimSpace(endpoints["mse"].(string)))
		config.ConfigEndpoint = strings.TrimSpace(endpoints["config"].(string))
		config.RKvstoreEndpoint = strings.TrimSpace(endpoints["r_kvstore"].(string))
		config.FnfEndpoint = strings.TrimSpace(endpoints["fnf"].(string))
		config.RosEndpoint = strings.TrimSpace(endpoints["ros"].(string))
		config.PrivatelinkEndpoint = strings.TrimSpace(endpoints["privatelink"].(string))
		config.ResourcesharingEndpoint = strings.TrimSpace(endpoints["ressharing"].(string))
		config.GaEndpoint = strings.TrimSpace(endpoints["ga"].(string))
		config.HitsdbEndpoint = strings.TrimSpace(endpoints["hitsdb"].(string))
		config.BrainIndustrialEndpoint = strings.TrimSpace(endpoints["brain_industrial"].(string))
		config.EipanycastEndpoint = strings.TrimSpace(endpoints["eipanycast"].(string))
		config.ImsEndpoint = strings.TrimSpace(endpoints["ims"].(string))
		config.QuotasEndpoint = strings.TrimSpace(endpoints["quotas"].(string))
		config.SgwEndpoint = strings.TrimSpace(endpoints["sgw"].(string))
		config.DmEndpoint = strings.TrimSpace(endpoints["dm"].(string))
		config.EventbridgeEndpoint = strings.TrimSpace(endpoints["eventbridge"].(string))
		config.OnsproxyEndpoint = strings.TrimSpace(endpoints["onsproxy"].(string))
		config.CdsEndpoint = strings.TrimSpace(endpoints["cds"].(string))
		config.HbrEndpoint = strings.TrimSpace(endpoints["hbr"].(string))
		config.ArmsEndpoint = strings.TrimSpace(endpoints["arms"].(string))
		config.ServerlessEndpoint = strings.TrimSpace(endpoints["serverless"].(string))
		config.AlbEndpoint = strings.TrimSpace(endpoints["alb"].(string))
		config.RedisaEndpoint = strings.TrimSpace(endpoints["redisa"].(string))
		config.GwsecdEndpoint = strings.TrimSpace(endpoints["gwsecd"].(string))
		config.CloudphoneEndpoint = strings.TrimSpace(endpoints["cloudphone"].(string))
		config.ScdnEndpoint = strings.TrimSpace(endpoints["scdn"].(string))
		config.DataworkspublicEndpoint = strings.TrimSpace(endpoints["dataworkspublic"].(string))
		config.HcsSgwEndpoint = strings.TrimSpace(endpoints["hcs_sgw"].(string))
		config.CddcEndpoint = strings.TrimSpace(endpoints["cddc"].(string))
		config.MscopensubscriptionEndpoint = strings.TrimSpace(endpoints["mscopensubscription"].(string))
		config.SddpEndpoint = strings.TrimSpace(endpoints["sddp"].(string))
		config.BastionhostEndpoint = strings.TrimSpace(endpoints["bastionhost"].(string))
		config.SasEndpoint = strings.TrimSpace(endpoints["sas"].(string))
		config.AlidfsEndpoint = strings.TrimSpace(endpoints["alidfs"].(string))
		config.EhpcEndpoint = strings.TrimSpace(endpoints["ehpc"].(string))
		config.EnsEndpoint = strings.TrimSpace(endpoints["ens"].(string))
		config.IotEndpoint = strings.TrimSpace(endpoints["iot"].(string))
		config.ImmEndpoint = strings.TrimSpace(endpoints["imm"].(string))
		config.ClickhouseEndpoint = strings.TrimSpace(endpoints["clickhouse"].(string))
		config.SelectDBEndpoint = strings.TrimSpace(endpoints["selectdb"].(string))
		config.DtsEndpoint = strings.TrimSpace(endpoints["dts"].(string))
		config.DgEndpoint = strings.TrimSpace(endpoints["dg"].(string))
		config.CloudssoEndpoint = strings.TrimSpace(endpoints["cloudsso"].(string))
		config.WafEndpoint = strings.TrimSpace(endpoints["waf"].(string))
		config.SwasEndpoint = strings.TrimSpace(endpoints["swas"].(string))
		config.VsEndpoint = strings.TrimSpace(endpoints["vs"].(string))
		config.QuickbiEndpoint = strings.TrimSpace(endpoints["quickbi"].(string))
		config.VodEndpoint = strings.TrimSpace(endpoints["vod"].(string))
		config.OpensearchEndpoint = strings.TrimSpace(endpoints["opensearch"].(string))
		config.GdsEndpoint = strings.TrimSpace(endpoints["gds"].(string))
		config.DbfsEndpoint = strings.TrimSpace(endpoints["dbfs"].(string))
		config.DevopsrdcEndpoint = strings.TrimSpace(endpoints["devopsrdc"].(string))
		config.EaisEndpoint = strings.TrimSpace(endpoints["eais"].(string))
		config.CloudauthEndpoint = strings.TrimSpace(endpoints["cloudauth"].(string))
		config.ImpEndpoint = strings.TrimSpace(endpoints["imp"].(string))
		config.MhubEndpoint = strings.TrimSpace(endpoints["mhub"].(string))
		config.ServicemeshEndpoint = strings.TrimSpace(endpoints["servicemesh"].(string))
		config.AcrEndpoint = strings.TrimSpace(endpoints["acr"].(string))
		config.EdsuserEndpoint = strings.TrimSpace(endpoints["edsuser"].(string))
		config.GaplusEndpoint = strings.TrimSpace(endpoints["gaplus"].(string))
		config.DdosbasicEndpoint = strings.TrimSpace(endpoints["ddosbasic"].(string))
		config.SmartagEndpoint = strings.TrimSpace(endpoints["smartag"].(string))
		config.TagEndpoint = strings.TrimSpace(endpoints["tag"].(string))
		config.EdasEndpoint = strings.TrimSpace(endpoints["edas"].(string))
		config.EdasschedulerxEndpoint = strings.TrimSpace(endpoints["edasschedulerx"].(string))
		config.EhsEndpoint = strings.TrimSpace(endpoints["ehs"].(string))
		config.CloudfwEndpoint = strings.TrimSpace(endpoints["cloudfw"].(string))
		config.DysmsEndpoint = strings.TrimSpace(endpoints["dysms"].(string))
		config.CbsEndpoint = strings.TrimSpace(endpoints["cbs"].(string))
		config.NlbEndpoint = strings.TrimSpace(endpoints["nlb"].(string))
		config.VpcpeerEndpoint = strings.TrimSpace(endpoints["vpcpeer"].(string))
		config.EbsEndpoint = strings.TrimSpace(endpoints["ebs"].(string))
		config.DmsenterpriseEndpoint = strings.TrimSpace(endpoints["dmsenterprise"].(string))
		config.BpStudioEndpoint = strings.TrimSpace(endpoints["bpstudio"].(string))
		config.DasEndpoint = strings.TrimSpace(endpoints["das"].(string))
		config.CloudfirewallEndpoint = strings.TrimSpace(endpoints["cloudfirewall"].(string))
		config.SrvcatalogEndpoint = strings.TrimSpace(endpoints["srvcatalog"].(string))
		config.VpcPeerEndpoint = strings.TrimSpace(endpoints["vpcpeer"].(string))
		config.EfloEndpoint = strings.TrimSpace(endpoints["eflo"].(string))
		config.OceanbaseEndpoint = strings.TrimSpace(endpoints["oceanbase"].(string))
		config.BeebotEndpoint = strings.TrimSpace(endpoints["beebot"].(string))
		config.ComputeNestEndpoint = strings.TrimSpace(endpoints["computenest"].(string))
		if endpoint, ok := endpoints["alidns"]; ok {
			config.AlidnsEndpoint = strings.TrimSpace(endpoint.(string))
		} else {
			config.AlidnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
		}
		config.CassandraEndpoint = strings.TrimSpace(endpoints["cassandra"].(string))
	}

	if otsInstanceName, ok := d.GetOk("ots_instance_name"); ok && otsInstanceName.(string) != "" {
		config.OtsInstanceName = strings.TrimSpace(otsInstanceName.(string))
	}

	if logEndpoint, ok := d.GetOk("log_endpoint"); ok && logEndpoint.(string) != "" {
		config.LogEndpoint = strings.TrimSpace(logEndpoint.(string))
	}
	if mnsEndpoint, ok := d.GetOk("mns_endpoint"); ok && mnsEndpoint.(string) != "" {
		config.MnsEndpoint = strings.TrimSpace(mnsEndpoint.(string))
	}

	if account, ok := d.GetOk("account_id"); ok && account.(string) != "" {
		config.AccountId = strings.TrimSpace(account.(string))
	}

	if fcEndpoint, ok := d.GetOk("fc"); ok && fcEndpoint.(string) != "" {
		config.FcEndpoint = strings.TrimSpace(fcEndpoint.(string))
	}
	if config.StsEndpoint == "" {
		config.StsEndpoint = connectivity.LoadRegionalEndpoint(config.RegionId, "sts")
	}

	configurationSources := []string{
		fmt.Sprintf("Default/%s", config.TerraformTraceId),
	}

	// configuration source final value should also contain TF_APPEND_USER_AGENT value
	// there is need to deduplication
	config.ConfigurationSource += " " + strings.TrimSpace(os.Getenv("TF_APPEND_USER_AGENT"))
	if config.ConfigurationSource != "" {
		for _, s := range strings.Split(config.ConfigurationSource, " ") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			exist := false
			for _, con := range configurationSources {
				if s == con {
					exist = true
					break
				}
			}
			if !exist {
				configurationSources = append(configurationSources, s)
			}
		}
	}
	config.ConfigurationSource = strings.Join(configurationSources, " ") + getModuleAddr()

	var signVersion sync.Map
	config.SignVersion = &signVersion
	for _, version := range d.Get("sign_version").(*schema.Set).List() {
		for key, val := range version.(map[string]interface{}) {
			signVersion.Store(key, val)
		}
	}

	if err := config.RefreshAuthCredential(); err != nil {
		return nil, err
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("configuring Terraform Alibaba Cloud Provider: no valid credential sources for Terraform Alibaba Cloud Provider found.\n\n%s",
			"Please see https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#authentication\n"+
				"for more information about providing credentials.")
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// This is a global MutexKV for use within this plugin.
var alicloudMutexKV = mutexkv.NewMutexKV()

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key for API operations. You can retrieve this from the 'Security Management' section of the Alibaba Cloud console.",

		"secret_key": "The secret key for API operations. You can retrieve this from the 'Security Management' section of the Alibaba Cloud console.",

		"ecs_role_name": "The RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the 'Access Control' section of the Alibaba Cloud console.",

		"region": "The region where Alibaba Cloud operations will take place. Examples are cn-beijing, cn-hangzhou, eu-central-1, etc.",

		"security_token": "security token. A security token is only required if you are using Security Token Service.",

		"account_id": "The account ID for some service API operations. You can retrieve this from the 'Security Settings' section of the Alibaba Cloud console.",

		"profile": "The profile for API operations. If not set, the default profile created with `aliyun configure` will be used.",

		"shared_credentials_file": "The path to the shared credentials file. If not set this defaults to ~/.aliyun/config.json",

		"assume_role_role_arn": "The ARN of a RAM role to assume prior to making API calls.",

		"assume_role_session_name": "The session name to use when assuming the role. If omitted, `terraform` is passed to the AssumeRole call as session name.",

		"assume_role_policy": "The permissions applied when assuming a role. You cannot use, this policy to grant further permissions that are in excess to those of the, role that is being assumed.",

		"assume_role_session_expiration": "The time after which the established session for assuming role expires. Valid value range: [900-3600] seconds. Default to 0 (in this case Alicloud use own default value).",

		"skip_region_validation": "Skip static validation of region ID. Used by users of alternative AlibabaCloud-like APIs or users w/ access to regions that are not public (yet).",

		"configuration_source": "Use this to mark a terraform configuration file source.",

		"client_read_timeout":    "The maximum timeout of the client read request.",
		"client_connect_timeout": "The maximum timeout of the client connection server.",
		"source_ip":              "The source ip for the assume role invoking.",
		"secure_transport":       "The security transport for the assume role invoking.",
		"credentials_uri":        "The URI of sidecar credentials service.",
		"max_retry_timeout":      "The maximum retry timeout of the request.",

		"ecs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.",

		"rds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.",

		"slb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.",

		"vpc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.",

		"ess_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Autoscaling endpoints.",

		"oss_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom OSS endpoints.",

		"ons_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ONS endpoints.",

		"alikafka_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ALIKAFKA endpoints.",

		"dns_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DNS endpoints.",

		"ram_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RAM endpoints.",

		"cs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Container Service endpoints.",

		"cr_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Container Registry endpoints.",

		"cdn_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CDN endpoints.",

		"kms_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom KMS endpoints.",

		"ots_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Table Store endpoints.",

		"cms_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Cloud Monitor endpoints.",

		"pvtz_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Private Zone endpoints.",

		"sts_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom STS endpoints.",

		"log_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Log Service endpoints.",

		"drds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DRDS endpoints.",

		"dds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MongoDB endpoints.",

		"polardb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom PolarDB endpoints.",

		"gpdb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom GPDB endpoints.",

		"kvstore_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom R-KVStore endpoints.",

		"fc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Function Computing endpoints.",

		"apigateway_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Api Gateway endpoints.",

		"datahub_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Datahub endpoints.",

		"mns_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MNS endpoints.",

		"location_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Location Service endpoints.",

		"elasticsearch_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Elasticsearch endpoints.",

		"nas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom NAS endpoints.",

		"actiontrail_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Actiontrail endpoints.",

		"cas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CAS endpoints.",

		"bssopenapi_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BSSOPENAPI endpoints.",

		"ddoscoo_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DDOSCOO endpoints.",

		"ddosbgp_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DDOSBGP endpoints.",

		"emr_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom EMR endpoints.",

		"market_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Market Place endpoints.",

		"hbase_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom HBase endpoints.",

		"adb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom AnalyticDB endpoints.",

		"cbn_endpoint":        "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cbn endpoints.",
		"maxcompute_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MaxCompute endpoints.",

		"dms_enterprise_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dms_enterprise endpoints.",

		"waf_openapi_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom waf_openapi endpoints.",

		"resourcemanager_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom resourcemanager endpoints.",

		"alidns_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom alidns endpoints.",

		"cassandra_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cassandra endpoints.",

		"eci_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eci endpoints.",

		"oos_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom oos endpoints.",

		"dcdn_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dcdn endpoints.",

		"mse_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom mse endpoints.",

		"config_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom config endpoints.",

		"r_kvstore_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom r_kvstore endpoints.",

		"fnf_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom fnf endpoints.",

		"ros_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ros endpoints.",

		"privatelink_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom privatelink endpoints.",

		"resourcesharing_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom resourcesharing endpoints.",

		"ga_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ga endpoints.",

		"hitsdb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom hitsdb endpoints.",

		"brain_industrial_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom brain_industrial endpoints.",

		"eipanycast_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eipanycast endpoints.",

		"ims_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ims endpoints.",

		"quotas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom quotas endpoints.",

		"sgw_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom sgw endpoints.",

		"dm_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dm endpoints.",

		"eventbridge_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eventbridge_share endpoints.",

		"onsproxy_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom onsproxy endpoints.",

		"cds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cds endpoints.",

		"hbr_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom hbr endpoints.",

		"arms_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom arms endpoints.",

		"serverless_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom serverless endpoints.",

		"alb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom alb endpoints.",

		"redisa_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom redisa endpoints.",

		"gwsecd_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom gwsecd endpoints.",

		"cloudphone_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cloudphone endpoints.",

		"scdn_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom scdn endpoints.",

		"dataworkspublic_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dataworkspublic endpoints.",

		"hcs_sgw_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom hcs_sgw endpoints.",

		"cddc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cddc endpoints.",

		"mscopensubscription_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom mscopensubscription endpoints.",

		"sddp_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom sddp endpoints.",

		"bastionhost_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom bastionhost endpoints.",

		"sas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom sas endpoints.",

		"alidfs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom alidfs endpoints.",

		"ehpc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ehpc endpoints.",

		"ens_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ens endpoints.",

		"iot_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom iot endpoints.",

		"imm_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom imm endpoints.",

		"clickhouse_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom clickhouse endpoints.",

		"selectdb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom selectdb endpoints.",

		"dts_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dts endpoints.",

		"dg_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dg endpoints.",

		"cloudsso_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cloudsso endpoints.",

		"waf_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom waf endpoints.",

		"swas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom swas endpoints.",

		"vs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom vs endpoints.",

		"quickbi_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom quickbi endpoints.",

		"vod_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom vod endpoints.",

		"opensearch_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom opensearch endpoints.",

		"gds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom gds endpoints.",

		"dbfs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dbfs endpoints.",

		"devopsrdc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom devopsrdc endpoints.",

		"eais_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eais endpoints.",

		"cloudauth_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cloudauth endpoints.",

		"imp_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom imp endpoints.",

		"mhub_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom mhub endpoints.",

		"servicemesh_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom servicemesh endpoints.",

		"acr_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom acr endpoints.",

		"edsuser_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom edsuser endpoints.",

		"gaplus_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom gaplus endpoints.",

		"ddosbasic_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ddosbasic endpoints.",

		"smartag_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom smartag endpoints.",

		"tag_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom tag endpoints.",

		"edas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom edas endpoints.",

		"edasschedulerx_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom edasschedulerx endpoints.",

		"ehs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ehs endpoints.",

		"cloudfw_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cloudfw endpoints.",

		"dysmsapi_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dysmsapi endpoints.",

		"cbs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cbs endpoints.",

		"nlb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom nlb endpoints.",

		"vpcpeer_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom vpcpeer endpoints.",

		"ebs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ebs endpoints.",

		"dmsenterprise_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dmsenterprise endpoints.",

		"bpstudio_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom bpstudio endpoints.",

		"das_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom das endpoints.",

		"cloudfirewall_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cloudfirewall endpoints.",

		"srvcatalog_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom srvcatalog endpoints.",

		"eflo_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eflo endpoints.",

		"oceanbase_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom oceanbase endpoints.",

		"beebot_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom beebot endpoints.",

		"computenest_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom computenest endpoints.",
	}
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Required:    true,
					Description: descriptions["assume_role_role_arn"],
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_ASSUME_ROLE_ARN", "ALIBABA_CLOUD_ROLE_ARN"}, nil),
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ALICLOUD_ASSUME_ROLE_SESSION_NAME", "ALIBABA_CLOUD_ROLE_SESSION_NAME"}, nil),
				},
				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
				"session_expiration": {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  descriptions["assume_role_session_expiration"],
					ValidateFunc: IntBetween(900, 43200),
				},
				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["external_id"],
				},
			},
		},
	}
}

func assumeRoleWithOidcSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oidc_provider_arn": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "ARN of the OIDC IdP.",
					DefaultFunc: schema.EnvDefaultFunc("ALIBABA_CLOUD_OIDC_PROVIDER_ARN", ""),
				},
				"oidc_token_file": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The file path of OIDC token that is issued by the external IdP.",
					DefaultFunc: schema.EnvDefaultFunc("ALIBABA_CLOUD_OIDC_TOKEN_FILE", ""),
					//ExactlyOneOf: []string{"assume_role_with_oidc.0.oidc_token", "assume_role_with_oidc.0.oidc_token_file"},
				},
				"oidc_token": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(4, 20000),
					//ExactlyOneOf: []string{"assume_role_with_oidc.0.oidc_token", "assume_role_with_oidc.0.oidc_token_file"},
				},
				"role_arn": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "ARN of a RAM role to assume prior to making API calls.",
					DefaultFunc: schema.EnvDefaultFunc("ALIBABA_CLOUD_ROLE_ARN", ""),
				},
				"role_session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom name of the role session. Set this parameter based on your business requirements. In most cases, this parameter is set to the identity of the user who calls the operation, for example, the username.",
					DefaultFunc: schema.EnvDefaultFunc("ALIBABA_CLOUD_ROLE_SESSION_NAME", ""),
				},
				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The policy that specifies the permissions of the returned STS token. You can use this parameter to grant the STS token fewer permissions than the permissions granted to the RAM role.",
				},
				"session_expiration": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "The validity period of the STS token. Unit: seconds. Default value: 3600. Minimum value: 900. Maximum value: the value of the MaxSessionDuration parameter when creating a ram role.",
				},
			},
		},
	}
}

func signVersionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oss": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sls": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"computenest": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["computenest_endpoint"],
				},

				"beebot": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["beebot_endpoint"],
				},

				"eflo": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eflo_endpoint"],
				},

				"srvcatalog": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["srvcatalog_endpoint"],
				},

				"cloudfirewall": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudfirewall_endpoint"],
				},

				"das": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["das_endpoint"],
				},

				"bpstudio": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bpstudio_endpoint"],
				},

				"dmsenterprise": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dmsenterprise_endpoint"],
				},

				"ebs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ebs_endpoint"],
				},

				"nlb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["nlb_endpoint"],
				},

				"cbs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cbs_endpoint"],
				},

				"vpcpeer": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vpcpeer_endpoint"],
				},

				"dysms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dysms_endpoint"],
				},

				"edas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["edas_endpoint"],
				},

				"edasschedulerx": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["edasschedulerx_endpoint"],
				},

				"ehs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ehs_endpoint"],
				},

				"tag": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["tag_endpoint"],
				},

				"ddosbasic": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddosbasic_endpoint"],
				},

				"smartag": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["smartag_endpoint"],
				},

				"oceanbase": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["oceanbase_endpoint"],
				},

				"gaplus": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gaplus_endpoint"],
				},

				"cloudfw": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudfw_endpoint"],
				},

				"edsuser": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["edsuser_endpoint"],
				},

				"acr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["acr_endpoint"],
				},

				"imp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["imp_endpoint"],
				},
				"eais": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eais_endpoint"],
				},
				"cloudauth": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudauth_endpoint"],
				},

				"mhub": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mhub_endpoint"],
				},
				"servicemesh": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["servicemesh_endpoint"],
				},
				"quickbi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["quickbi_endpoint"],
				},
				"vod": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vod_endpoint"],
				},
				"opensearch": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["opensearch_endpoint"],
				},
				"gds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gds_endpoint"],
				},
				"dbfs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dbfs_endpoint"],
				},
				"devopsrdc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["devopsrdc_endpoint"],
				},
				"dg": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dg_endpoint"],
				},
				"waf": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["waf_endpoint"],
				},
				"vs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vs_endpoint"],
				},
				"dts": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dts_endpoint"],
				},
				"cloudsso": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudsso_endpoint"],
				},

				"iot": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["iot_endpoint"],
				},
				"swas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["swas_endpoint"],
				},

				"imm": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["imm_endpoint"],
				},
				"clickhouse": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["clickhouse_endpoint"],
				},
				"selectdb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["selectdb_endpoint"],
				},

				"alidfs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alidfs_endpoint"],
				},

				"ens": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ens_endpoint"],
				},

				"bastionhost": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bastionhost_endpoint"],
				},
				"cddc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cddc_endpoint"],
				},
				"sddp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sddp_endpoint"],
				},

				"mscopensubscription": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mscopensubscription_endpoint"],
				},

				"sas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sas_endpoint"],
				},

				"ehpc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ehpc_endpoint"],
				},

				"dataworkspublic": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dataworkspublic_endpoint"],
				},

				"hcs_sgw": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["hcs_sgw_endpoint"],
				},

				"cloudphone": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudphone_endpoint"],
				},

				"alb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alb_endpoint"],
				},
				"redisa": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["redisa_endpoint"],
				},
				"gwsecd": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gwsecd_endpoint"],
				},
				"scdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["scdn_endpoint"],
				},

				"arms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["arms_endpoint"],
				},
				"serverless": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["serverless_endpoint"],
				},

				"hbr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["hbr_endpoint"],
				},

				"onsproxy": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["onsproxy_endpoint"],
				},
				"cds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cds_endpoint"],
				},

				"dm": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dm_endpoint"],
				},

				"eventbridge": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eventbridge_endpoint"],
				},

				"sgw": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sgw_endpoint"],
				},

				"quotas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["quotas_endpoint"],
				},

				"ims": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ims_endpoint"],
				},

				"brain_industrial": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["brain_industrial_endpoint"],
				},

				"ressharing": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["resourcesharing_endpoint"],
				},
				"ga": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ga_endpoint"],
				},

				"hitsdb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["hitsdb_endpoint"],
				},

				"privatelink": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["privatelink_endpoint"],
				},

				"eipanycast": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eipanycast_endpoint"],
				},

				"fnf": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["fnf_endpoint"],
				},

				"ros": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ros_endpoint"],
				},

				"r_kvstore": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["r_kvstore_endpoint"],
				},

				"config": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["config_endpoint"],
				},

				"dcdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dcdn_endpoint"],
				},

				"mse": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mse_endpoint"],
				},

				"oos": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["oos_endpoint"],
				},

				"eci": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eci_endpoint"],
				},

				"alidns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alidns_endpoint"],
				},

				"resourcemanager": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["resourcemanager_endpoint"],
				},

				"waf_openapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["waf_openapi_endpoint"],
				},

				"dms_enterprise": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dms_enterprise_endpoint"],
				},

				"cassandra": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cassandra_endpoint"],
				},

				"cbn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cbn_endpoint"],
				},

				"ecs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ecs_endpoint"],
				},
				"rds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["rds_endpoint"],
				},
				"slb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["slb_endpoint"],
				},
				"vpc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vpc_endpoint"],
				},
				"ess": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ess_endpoint"],
				},
				"oss": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["oss_endpoint"],
				},
				"ons": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ons_endpoint"],
				},
				"alikafka": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alikafka_endpoint"],
				},
				"dns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dns_endpoint"],
				},
				"ram": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ram_endpoint"],
				},
				"cs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cs_endpoint"],
				},
				"cr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cr_endpoint"],
				},
				"cdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cdn_endpoint"],
				},

				"kms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kms_endpoint"],
				},

				"ots": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ots_endpoint"],
				},

				"cms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cms_endpoint"],
				},

				"pvtz": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["pvtz_endpoint"],
				},

				"sts": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sts_endpoint"],
				},
				// log service is sls service
				"log": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["log_endpoint"],
				},
				"drds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["drds_endpoint"],
				},
				"dds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dds_endpoint"],
				},
				"polardb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["polardb_endpoint"],
				},
				"gpdb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gpdb_endpoint"],
				},
				"kvstore": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kvstore_endpoint"],
				},
				"fc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["fc_endpoint"],
				},
				"apigateway": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["apigateway_endpoint"],
				},
				"datahub": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["datahub_endpoint"],
				},
				"mns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mns_endpoint"],
				},
				"location": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["location_endpoint"],
				},
				"elasticsearch": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["elasticsearch_endpoint"],
				},
				"nas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["nas_endpoint"],
				},
				"actiontrail": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["actiontrail_endpoint"],
				},
				"cas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cas_endpoint"],
				},
				"bssopenapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bssopenapi_endpoint"],
				},
				"ddoscoo": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddoscoo_endpoint"],
				},
				"ddosbgp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddosbgp_endpoint"],
				},
				"emr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["emr_endpoint"],
				},
				"market": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["market_endpoint"],
				},
				"adb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["adb_endpoint"],
				},
				"maxcompute": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["maxcompute_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}

func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["ecs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["rds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["slb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vpc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ess"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["oss"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ons"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alikafka"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ram"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ots"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["pvtz"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sts"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["log"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["drds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gpdb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kvstore"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["polardb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["fc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["apigateway"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["datahub"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["location"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["elasticsearch"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["nas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["actiontrail"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bssopenapi"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddoscoo"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddosbgp"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["emr"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["market"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["adb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cbn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["maxcompute"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dms_enterprise"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["waf_openapi"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["resourcemanager"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alidns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cassandra"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eci"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["oos"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dcdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mse"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["config"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["r_kvstore"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["fnf"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ros"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["privatelink"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ressharing"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ga"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["hitsdb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["brain_industrial"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eipanycast"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ims"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["quotas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sgw"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dm"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eventbridge"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["onsproxy"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["hbr"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["arms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["serverless"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["redisa"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gwsecd"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudphone"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["scdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dataworkspublic"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["hcs_sgw"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cddc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mscopensubscription"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sddp"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bastionhost"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alidfs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ehpc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ens"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["iot"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["imm"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["clickhouse"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["selectdb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dts"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dg"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudsso"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["waf"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["swas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["quickbi"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vod"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["opensearch"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dbfs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["devopsrdc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eais"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudauth"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["imp"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mhub"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["servicemesh"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["acr"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["edsuser"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gaplus"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddosbasic"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["smartag"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["tag"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["edas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["edasschedulerx"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ehs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudfw"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dysms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cbs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["nlb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vpcpeer"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ebs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dmsenterprise"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bpstudio"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["das"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudfirewall"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["srvcatalog"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vpcpeer"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eflo"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["oceanbase"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["beebot"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["computenest"].(string)))
	return hashcode.String(buf.String())
}

func getConfigFromProfile(d *schema.ResourceData, ProfileKey string) (interface{}, error) {

	if providerConfig == nil {
		if v, ok := d.GetOk("profile"); !ok && v.(string) == "" {
			return nil, nil
		}
		current := d.Get("profile").(string)
		// Set Credentials filename, expanding home directory
		profilePath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
		if err != nil {
			return nil, WrapError(err)
		}
		if profilePath == "" {
			profilePath = fmt.Sprintf("%s/.aliyun/config.json", os.Getenv("HOME"))
			if runtime.GOOS == "windows" {
				profilePath = fmt.Sprintf("%s/.aliyun/config.json", os.Getenv("USERPROFILE"))
			}
		}
		providerConfig = make(map[string]interface{})
		_, err = os.Stat(profilePath)
		if !os.IsNotExist(err) {
			data, err := ioutil.ReadFile(profilePath)
			if err != nil {
				return nil, WrapError(err)
			}
			config := map[string]interface{}{}
			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, WrapError(err)
			}
			for _, v := range config["profiles"].([]interface{}) {
				if current == v.(map[string]interface{})["name"] {
					providerConfig = v.(map[string]interface{})
				}
			}
		}
	}

	mode := ""
	if v, ok := providerConfig["mode"]; ok {
		mode = v.(string)
	} else {
		return v, nil
	}
	switch ProfileKey {
	case "access_key_id", "access_key_secret":
		if mode == "EcsRamRole" {
			return "", nil
		}
	case "ram_role_name":
		if mode != "EcsRamRole" {
			return "", nil
		}
	case "sts_token":
		if mode != "StsToken" {
			return "", nil
		}
	case "ram_role_arn", "ram_session_name":
		if mode != "RamRoleArn" {
			return "", nil
		}
	case "expired_seconds":
		if mode != "RamRoleArn" {
			return float64(0), nil
		}
	}

	return providerConfig[ProfileKey], nil
}

func getAssumeRoleWithOIDCConfig(tfMap map[string]interface{}) (*connectivity.AssumeRoleWithOidc, error) {
	if tfMap == nil {
		return nil, nil
	}

	assumeRole := connectivity.AssumeRoleWithOidc{}

	if v, ok := tfMap["session_expiration"].(int); ok && v != 0 {
		assumeRole.DurationSeconds = v
	}

	if v, ok := tfMap["policy"].(string); ok && v != "" {
		assumeRole.Policy = v
	}

	if v, ok := tfMap["role_arn"].(string); ok && v != "" {
		assumeRole.RoleARN = v
	}

	if v, ok := tfMap["role_session_name"].(string); ok && v != "" {
		assumeRole.RoleSessionName = v
	}
	if assumeRole.RoleSessionName == "" {
		assumeRole.RoleSessionName = "terraform"
	}

	if v, ok := tfMap["oidc_provider_arn"].(string); ok && v != "" {
		assumeRole.OIDCProviderArn = v
	}

	missingOidcToken := true
	if v, ok := tfMap["oidc_token"].(string); ok && v != "" {
		assumeRole.OIDCToken = v
		missingOidcToken = false
	}

	if v, ok := tfMap["oidc_token_file"].(string); ok && v != "" {
		assumeRole.OIDCTokenFile = v
		if assumeRole.OIDCToken == "" {
			token, err := os.ReadFile(v)
			if err != nil {
				return nil, fmt.Errorf("reading oidc_token_file failed. Error: %s", err)
			}
			assumeRole.OIDCToken = string(token)
		}
		missingOidcToken = false
	}
	if missingOidcToken {
		return nil, fmt.Errorf("\"assume_role_with_oidc.0.oidc_token\": one of `assume_role_with_oidc.0.oidc_token,assume_role_with_oidc.0.oidc_token_file` must be specified")
	}

	if assumeRole.OIDCToken == "" {
		return nil, fmt.Errorf("\"assume_role_with_oidc.0.oidc_token\" or \"assume_role_with_oidc.0.oidc_token_file\" content can not be empty")
	}

	return &assumeRole, nil
}

type CredentialsURIResponse struct {
	Code            string
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
}

func getClientByCredentialsURI(credentialsURI string) (*CredentialsURIResponse, error) {
	res, err := http.Get(credentialsURI)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("get Credentials from %s failed, status code %d", credentialsURI, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var response CredentialsURIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal credentials failed, the body %s", string(body))
	}

	if response.Code != "Success" {
		return nil, fmt.Errorf("fetching sts token from %s got an error and its Code is not Success", credentialsURI)
	}

	return &response, nil
}

func getModuleAddr() string {
	moduleMeta := make(map[string]interface{})
	str, err := os.ReadFile(".terraform/modules/modules.json")
	if err != nil {
		return ""
	}
	err = json.Unmarshal(str, &moduleMeta)
	if err != nil || len(moduleMeta) < 1 || moduleMeta["Modules"] == nil {
		return ""
	}
	var result string
	for _, m := range moduleMeta["Modules"].([]interface{}) {
		module := m.(map[string]interface{})
		moduleSource := fmt.Sprint(module["Source"])
		moduleVersion := fmt.Sprint(module["Version"])
		if strings.HasPrefix(moduleSource, "registry.terraform.io/") {
			parts := strings.Split(moduleSource, "/")
			if len(parts) == 4 {
				result += " " + "terraform-" + parts[3] + "-" + parts[2] + "/" + moduleVersion
			}
		}
	}
	return result
}
