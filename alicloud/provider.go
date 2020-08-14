package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

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
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY")),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY")),
				Description: descriptions["secret_key"],
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECURITY_TOKEN", os.Getenv("SECURITY_TOKEN")),
				Description: descriptions["security_token"],
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ECS_ROLE_NAME", os.Getenv("ALICLOUD_ECS_ROLE_NAME")),
				Description: descriptions["ecs_role_name"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_REGION", os.Getenv("ALICLOUD_REGION")),
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
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCOUNT_ID", os.Getenv("ALICLOUD_ACCOUNT_ID")),
				Description: descriptions["account_id"],
			},
			"assume_role": assumeRoleSchema(),
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
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SHARED_CREDENTIALS_FILE", ""),
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_PROFILE", ""),
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
				Default:      "",
				Description:  descriptions["configuration_source"],
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTPS",
				Description:  descriptions["protocol"],
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{

			"alicloud_account":                dataSourceAlicloudAccount(),
			"alicloud_caller_identity":        dataSourceAlicloudCallerIdentity(),
			"alicloud_images":                 dataSourceAlicloudImages(),
			"alicloud_regions":                dataSourceAlicloudRegions(),
			"alicloud_zones":                  dataSourceAlicloudZones(),
			"alicloud_db_zones":               dataSourceAlicloudDBZones(),
			"alicloud_instance_type_families": dataSourceAlicloudInstanceTypeFamilies(),
			"alicloud_instance_types":         dataSourceAlicloudInstanceTypes(),
			"alicloud_instances":              dataSourceAlicloudInstances(),
			"alicloud_disks":                  dataSourceAlicloudDisks(),
			"alicloud_network_interfaces":     dataSourceAlicloudNetworkInterfaces(),
			"alicloud_snapshots":              dataSourceAlicloudSnapshots(),
			"alicloud_vpcs":                   dataSourceAlicloudVpcs(),
			"alicloud_vswitches":              dataSourceAlicloudVSwitches(),
			"alicloud_eips":                   dataSourceAlicloudEips(),
			"alicloud_key_pairs":              dataSourceAlicloudKeyPairs(),
			"alicloud_kms_keys":               dataSourceAlicloudKmsKeys(),
			"alicloud_kms_ciphertext":         dataSourceAlicloudKmsCiphertext(),
			"alicloud_kms_plaintext":          dataSourceAlicloudKmsPlaintext(),
			"alicloud_dns_resolution_lines":   dataSourceAlicloudDnsResolutionLines(),
			"alicloud_dns_domains":            dataSourceAlicloudDnsDomains(),
			"alicloud_dns_groups":             dataSourceAlicloudDnsGroups(),
			"alicloud_dns_records":            dataSourceAlicloudDnsRecords(),
			// alicloud_dns_domain_groups, alicloud_dns_domain_records have been deprecated.
			"alicloud_dns_domain_groups":  dataSourceAlicloudDnsGroups(),
			"alicloud_dns_domain_records": dataSourceAlicloudDnsRecords(),
			// alicloud_ram_account_alias has been deprecated
			"alicloud_ram_account_alias":                     dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_aliases":                   dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_groups":                            dataSourceAlicloudRamGroups(),
			"alicloud_ram_users":                             dataSourceAlicloudRamUsers(),
			"alicloud_ram_roles":                             dataSourceAlicloudRamRoles(),
			"alicloud_ram_policies":                          dataSourceAlicloudRamPolicies(),
			"alicloud_security_groups":                       dataSourceAlicloudSecurityGroups(),
			"alicloud_security_group_rules":                  dataSourceAlicloudSecurityGroupRules(),
			"alicloud_slbs":                                  dataSourceAlicloudSlbs(),
			"alicloud_slb_attachments":                       dataSourceAlicloudSlbAttachments(),
			"alicloud_slb_backend_servers":                   dataSourceAlicloudSlbBackendServers(),
			"alicloud_slb_listeners":                         dataSourceAlicloudSlbListeners(),
			"alicloud_slb_rules":                             dataSourceAlicloudSlbRules(),
			"alicloud_slb_server_groups":                     dataSourceAlicloudSlbServerGroups(),
			"alicloud_slb_master_slave_server_groups":        dataSourceAlicloudSlbMasterSlaveServerGroups(),
			"alicloud_slb_acls":                              dataSourceAlicloudSlbAcls(),
			"alicloud_slb_server_certificates":               dataSourceAlicloudSlbServerCertificates(),
			"alicloud_slb_ca_certificates":                   dataSourceAlicloudSlbCACertificates(),
			"alicloud_slb_domain_extensions":                 dataSourceAlicloudSlbDomainExtensions(),
			"alicloud_slb_zones":                             dataSourceAlicloudSlbZones(),
			"alicloud_oss_bucket_objects":                    dataSourceAlicloudOssBucketObjects(),
			"alicloud_oss_buckets":                           dataSourceAlicloudOssBuckets(),
			"alicloud_ons_instances":                         dataSourceAlicloudOnsInstances(),
			"alicloud_ons_topics":                            dataSourceAlicloudOnsTopics(),
			"alicloud_ons_groups":                            dataSourceAlicloudOnsGroups(),
			"alicloud_alikafka_consumer_groups":              dataSourceAlicloudAlikafkaConsumerGroups(),
			"alicloud_alikafka_instances":                    dataSourceAlicloudAlikafkaInstances(),
			"alicloud_alikafka_topics":                       dataSourceAlicloudAlikafkaTopics(),
			"alicloud_alikafka_sasl_users":                   dataSourceAlicloudAlikafkaSaslUsers(),
			"alicloud_alikafka_sasl_acls":                    dataSourceAlicloudAlikafkaSaslAcls(),
			"alicloud_fc_functions":                          dataSourceAlicloudFcFunctions(),
			"alicloud_file_crc64_checksum":                   dataSourceAlicloudFileCRC64Checksum(),
			"alicloud_fc_services":                           dataSourceAlicloudFcServices(),
			"alicloud_fc_triggers":                           dataSourceAlicloudFcTriggers(),
			"alicloud_fc_zones":                              dataSourceAlicloudFcZones(),
			"alicloud_db_instances":                          dataSourceAlicloudDBInstances(),
			"alicloud_db_instance_engines":                   dataSourceAlicloudDBInstanceEngines(),
			"alicloud_db_instance_classes":                   dataSourceAlicloudDBInstanceClasses(),
			"alicloud_pvtz_zones":                            dataSourceAlicloudPvtzZones(),
			"alicloud_pvtz_zone_records":                     dataSourceAlicloudPvtzZoneRecords(),
			"alicloud_router_interfaces":                     dataSourceAlicloudRouterInterfaces(),
			"alicloud_vpn_gateways":                          dataSourceAlicloudVpnGateways(),
			"alicloud_vpn_customer_gateways":                 dataSourceAlicloudVpnCustomerGateways(),
			"alicloud_vpn_connections":                       dataSourceAlicloudVpnConnections(),
			"alicloud_ssl_vpn_servers":                       dataSourceAlicloudSslVpnServers(),
			"alicloud_ssl_vpn_client_certs":                  dataSourceAlicloudSslVpnClientCerts(),
			"alicloud_mongo_instances":                       dataSourceAlicloudMongoDBInstances(),
			"alicloud_mongodb_instances":                     dataSourceAlicloudMongoDBInstances(),
			"alicloud_mongodb_zones":                         dataSourceAlicloudMongoDBZones(),
			"alicloud_gpdb_instances":                        dataSourceAlicloudGpdbInstances(),
			"alicloud_gpdb_zones":                            dataSourceAlicloudGpdbZones(),
			"alicloud_kvstore_instances":                     dataSourceAlicloudKVStoreInstances(),
			"alicloud_kvstore_zones":                         dataSourceAlicloudKVStoreZones(),
			"alicloud_kvstore_instance_classes":              dataSourceAlicloudKVStoreInstanceClasses(),
			"alicloud_kvstore_instance_engines":              dataSourceAlicloudKVStoreInstanceEngines(),
			"alicloud_cen_instances":                         dataSourceAlicloudCenInstances(),
			"alicloud_cen_bandwidth_packages":                dataSourceAlicloudCenBandwidthPackages(),
			"alicloud_cen_bandwidth_limits":                  dataSourceAlicloudCenBandwidthLimits(),
			"alicloud_cen_route_entries":                     dataSourceAlicloudCenRouteEntries(),
			"alicloud_cen_region_route_entries":              dataSourceAlicloudCenRegionRouteEntries(),
			"alicloud_cs_kubernetes_clusters":                dataSourceAlicloudCSKubernetesClusters(),
			"alicloud_cs_managed_kubernetes_clusters":        dataSourceAlicloudCSManagerKubernetesClusters(),
			"alicloud_cs_serverless_kubernetes_clusters":     dataSourceAlicloudCSServerlessKubernetesClusters(),
			"alicloud_cr_namespaces":                         dataSourceAlicloudCRNamespaces(),
			"alicloud_cr_repos":                              dataSourceAlicloudCRRepos(),
			"alicloud_cr_ee_instances":                       dataSourceAlicloudCrEEInstances(),
			"alicloud_cr_ee_namespaces":                      dataSourceAlicloudCrEENamespaces(),
			"alicloud_cr_ee_repos":                           dataSourceAlicloudCrEERepos(),
			"alicloud_cr_ee_sync_rules":                      dataSourceAlicloudCrEESyncRules(),
			"alicloud_mns_queues":                            dataSourceAlicloudMNSQueues(),
			"alicloud_mns_topics":                            dataSourceAlicloudMNSTopics(),
			"alicloud_mns_topic_subscriptions":               dataSourceAlicloudMNSTopicSubscriptions(),
			"alicloud_api_gateway_apis":                      dataSourceAlicloudApiGatewayApis(),
			"alicloud_api_gateway_groups":                    dataSourceAlicloudApiGatewayGroups(),
			"alicloud_api_gateway_apps":                      dataSourceAlicloudApiGatewayApps(),
			"alicloud_elasticsearch_instances":               dataSourceAlicloudElasticsearch(),
			"alicloud_elasticsearch_zones":                   dataSourceAlicloudElaticsearchZones(),
			"alicloud_drds_instances":                        dataSourceAlicloudDRDSInstances(),
			"alicloud_nas_access_groups":                     dataSourceAlicloudAccessGroups(),
			"alicloud_nas_access_rules":                      dataSourceAlicloudAccessRules(),
			"alicloud_nas_mount_targets":                     dataSourceAlicloudMountTargets(),
			"alicloud_nas_file_systems":                      dataSourceAlicloudFileSystems(),
			"alicloud_nas_protocols":                         dataSourceAlicloudNasProtocols(),
			"alicloud_cas_certificates":                      dataSourceAlicloudCasCertificates(),
			"alicloud_actiontrails":                          dataSourceAlicloudActiontrails(),
			"alicloud_common_bandwidth_packages":             dataSourceAlicloudCommonBandwidthPackages(),
			"alicloud_route_tables":                          dataSourceAlicloudRouteTables(),
			"alicloud_route_entries":                         dataSourceAlicloudRouteEntries(),
			"alicloud_nat_gateways":                          dataSourceAlicloudNatGateways(),
			"alicloud_snat_entries":                          dataSourceAlicloudSnatEntries(),
			"alicloud_forward_entries":                       dataSourceAlicloudForwardEntries(),
			"alicloud_ddoscoo_instances":                     dataSourceAlicloudDdoscooInstances(),
			"alicloud_ddosbgp_instances":                     dataSourceAlicloudDdosbgpInstances(),
			"alicloud_ess_alarms":                            dataSourceAlicloudEssAlarms(),
			"alicloud_ess_notifications":                     dataSourceAlicloudEssNotifications(),
			"alicloud_ess_scaling_groups":                    dataSourceAlicloudEssScalingGroups(),
			"alicloud_ess_scaling_rules":                     dataSourceAlicloudEssScalingRules(),
			"alicloud_ess_scaling_configurations":            dataSourceAlicloudEssScalingConfigurations(),
			"alicloud_ess_lifecycle_hooks":                   dataSourceAlicloudEssLifecycleHooks(),
			"alicloud_ess_scheduled_tasks":                   dataSourceAlicloudEssScheduledTasks(),
			"alicloud_ots_instances":                         dataSourceAlicloudOtsInstances(),
			"alicloud_ots_instance_attachments":              dataSourceAlicloudOtsInstanceAttachments(),
			"alicloud_ots_tables":                            dataSourceAlicloudOtsTables(),
			"alicloud_cloud_connect_networks":                dataSourceAlicloudCloudConnectNetworks(),
			"alicloud_emr_instance_types":                    dataSourceAlicloudEmrInstanceTypes(),
			"alicloud_emr_disk_types":                        dataSourceAlicloudEmrDiskTypes(),
			"alicloud_emr_main_versions":                     dataSourceAlicloudEmrMainVersions(),
			"alicloud_sag_acls":                              dataSourceAlicloudSagAcls(),
			"alicloud_yundun_dbaudit_instance":               dataSourceAlicloudDbauditInstances(),
			"alicloud_yundun_bastionhost_instances":          dataSourceAlicloudBastionhostInstances(),
			"alicloud_market_product":                        dataSourceAlicloudProduct(),
			"alicloud_market_products":                       dataSourceAlicloudProducts(),
			"alicloud_polardb_clusters":                      dataSourceAlicloudPolarDBClusters(),
			"alicloud_polardb_node_classes":                  dataSourceAlicloudPolarDBNodeClasses(),
			"alicloud_polardb_endpoints":                     dataSourceAlicloudPolarDBEndpoints(),
			"alicloud_polardb_accounts":                      dataSourceAlicloudPolarDBAccounts(),
			"alicloud_polardb_databases":                     dataSourceAlicloudPolarDBDatabases(),
			"alicloud_polardb_zones":                         dataSourceAlicloudPolarDBZones(),
			"alicloud_hbase_instances":                       dataSourceAlicloudHBaseInstances(),
			"alicloud_hbase_zones":                           dataSourceAlicloudHBaseZones(),
			"alicloud_adb_clusters":                          dataSourceAlicloudAdbClusters(),
			"alicloud_adb_zones":                             dataSourceAlicloudAdbZones(),
			"alicloud_cen_flowlogs":                          dataSourceAlicloudCenFlowlogs(),
			"alicloud_kms_aliases":                           dataSourceAlicloudKmsAliases(),
			"alicloud_dns_domain_txt_guid":                   dataSourceAlicloudDnsDomainTxtGuid(),
			"alicloud_edas_applications":                     dataSourceAlicloudEdasApplications(),
			"alicloud_edas_deploy_groups":                    dataSourceAlicloudEdasDeployGroups(),
			"alicloud_edas_clusters":                         dataSourceAlicloudEdasClusters(),
			"alicloud_resource_manager_folders":              dataSourceAlicloudResourceManagerFolders(),
			"alicloud_dns_instances":                         dataSourceAlicloudDnsInstances(),
			"alicloud_resource_manager_policies":             dataSourceAlicloudResourceManagerPolicies(),
			"alicloud_resource_manager_resource_groups":      dataSourceAlicloudResourceManagerResourceGroups(),
			"alicloud_resource_manager_roles":                dataSourceAlicloudResourceManagerRoles(),
			"alicloud_resource_manager_policy_versions":      dataSourceAlicloudResourceManagerPolicyVersions(),
			"alicloud_alidns_domain_groups":                  dataSourceAlicloudAlidnsDomainGroups(),
			"alicloud_kms_key_versions":                      dataSourceAlicloudKmsKeyVersions(),
			"alicloud_alidns_records":                        dataSourceAlicloudAlidnsRecords(),
			"alicloud_resource_manager_accounts":             dataSourceAlicloudResourceManagerAccounts(),
			"alicloud_resource_manager_resource_directories": dataSourceAlicloudResourceManagerResourceDirectories(),
			"alicloud_resource_manager_handshakes":           dataSourceAlicloudResourceManagerHandshakes(),
			"alicloud_waf_domains":                           dataSourceAlicloudWafDomains(),
			"alicloud_kms_secrets":                           dataSourceAlicloudKmsSecrets(),
			"alicloud_cen_route_maps":                        dataSourceAlicloudCenRouteMaps(),
			"alicloud_cen_private_zones":                     dataSourceAlicloudCenPrivateZones(),
			"alicloud_dms_enterprise_instances":              dataSourceAlicloudDmsEnterpriseInstances(),
			"alicloud_cassandra_clusters":                    dataSourceAlicloudCassandraClusters(),
			"alicloud_cassandra_data_centers":                dataSourceAlicloudCassandraDataCenters(),
			"alicloud_cassandra_zones":                       dataSourceAlicloudCassandraZones(),
			"alicloud_kms_secret_versions":                   dataSourceAlicloudKmsSecretVersions(),
			"alicloud_waf_instances":                         dataSourceAlicloudWafInstances(),
			"alicloud_eci_image_caches":                      dataSourceAlicloudEciImageCaches(),
			"alicloud_dms_enterprise_users":                  dataSourceAlicloudDmsEnterpriseUsers(),
			"alicloud_ecs_dedicated_hosts":                   dataSourceAlicloudEcsDedicatedHosts(),
			"alicloud_oos_templates":                         dataSourceAlicloudOosTemplates(),
			"alicloud_oos_executions":                        dataSourceAlicloudOosExecutions(),
			"alicloud_resource_manager_policy_attachments":   dataSourceAlicloudResourceManagerPolicyAttachments(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_instance":                           resourceAliyunInstance(),
			"alicloud_image":                              resourceAliCloudImage(),
			"alicloud_reserved_instance":                  resourceAliCloudReservedInstance(),
			"alicloud_copy_image":                         resourceAliCloudImageCopy(),
			"alicloud_image_export":                       resourceAliCloudImageExport(),
			"alicloud_image_copy":                         resourceAliCloudImageCopy(),
			"alicloud_image_import":                       resourceAliCloudImageImport(),
			"alicloud_image_share_permission":             resourceAliCloudImageSharePermission(),
			"alicloud_ram_role_attachment":                resourceAlicloudRamRoleAttachment(),
			"alicloud_disk":                               resourceAliyunDisk(),
			"alicloud_disk_attachment":                    resourceAliyunDiskAttachment(),
			"alicloud_network_interface":                  resourceAliyunNetworkInterface(),
			"alicloud_network_interface_attachment":       resourceAliyunNetworkInterfaceAttachment(),
			"alicloud_snapshot":                           resourceAliyunSnapshot(),
			"alicloud_snapshot_policy":                    resourceAliyunSnapshotPolicy(),
			"alicloud_launch_template":                    resourceAliyunLaunchTemplate(),
			"alicloud_security_group":                     resourceAliyunSecurityGroup(),
			"alicloud_security_group_rule":                resourceAliyunSecurityGroupRule(),
			"alicloud_db_database":                        resourceAlicloudDBDatabase(),
			"alicloud_db_account":                         resourceAlicloudDBAccount(),
			"alicloud_db_account_privilege":               resourceAlicloudDBAccountPrivilege(),
			"alicloud_db_backup_policy":                   resourceAlicloudDBBackupPolicy(),
			"alicloud_db_connection":                      resourceAlicloudDBConnection(),
			"alicloud_db_read_write_splitting_connection": resourceAlicloudDBReadWriteSplittingConnection(),
			"alicloud_db_instance":                        resourceAlicloudDBInstance(),
			"alicloud_mongodb_instance":                   resourceAlicloudMongoDBInstance(),
			"alicloud_mongodb_sharding_instance":          resourceAlicloudMongoDBShardingInstance(),
			"alicloud_gpdb_instance":                      resourceAlicloudGpdbInstance(),
			"alicloud_gpdb_connection":                    resourceAlicloudGpdbConnection(),
			"alicloud_db_readonly_instance":               resourceAlicloudDBReadonlyInstance(),
			"alicloud_auto_provisioning_group":            resourceAlicloudAutoProvisioningGroup(),
			"alicloud_ess_scaling_group":                  resourceAlicloudEssScalingGroup(),
			"alicloud_ess_scaling_configuration":          resourceAlicloudEssScalingConfiguration(),
			"alicloud_ess_scaling_rule":                   resourceAlicloudEssScalingRule(),
			"alicloud_ess_schedule":                       resourceAlicloudEssScheduledTask(),
			"alicloud_ess_scheduled_task":                 resourceAlicloudEssScheduledTask(),
			"alicloud_ess_attachment":                     resourceAlicloudEssAttachment(),
			"alicloud_ess_lifecycle_hook":                 resourceAlicloudEssLifecycleHook(),
			"alicloud_ess_notification":                   resourceAlicloudEssNotification(),
			"alicloud_ess_alarm":                          resourceAlicloudEssAlarm(),
			"alicloud_ess_scalinggroup_vserver_groups":    resourceAlicloudEssScalingGroupVserverGroups(),
			"alicloud_vpc":                                resourceAliyunVpc(),
			"alicloud_nat_gateway":                        resourceAliyunNatGateway(),
			"alicloud_nas_file_system":                    resourceAlicloudNasFileSystem(),
			"alicloud_nas_mount_target":                   resourceAlicloudNasMountTarget(),
			"alicloud_nas_access_group":                   resourceAlicloudNasAccessGroup(),
			"alicloud_nas_access_rule":                    resourceAlicloudNasAccessRule(),
			// "alicloud_subnet" aims to match aws usage habit.
			"alicloud_subnet":                        resourceAliyunSubnet(),
			"alicloud_vswitch":                       resourceAliyunSubnet(),
			"alicloud_route_entry":                   resourceAliyunRouteEntry(),
			"alicloud_route_table":                   resourceAliyunRouteTable(),
			"alicloud_route_table_attachment":        resourceAliyunRouteTableAttachment(),
			"alicloud_snat_entry":                    resourceAliyunSnatEntry(),
			"alicloud_forward_entry":                 resourceAliyunForwardEntry(),
			"alicloud_eip":                           resourceAliyunEip(),
			"alicloud_eip_association":               resourceAliyunEipAssociation(),
			"alicloud_slb":                           resourceAliyunSlb(),
			"alicloud_slb_listener":                  resourceAliyunSlbListener(),
			"alicloud_slb_attachment":                resourceAliyunSlbAttachment(),
			"alicloud_slb_backend_server":            resourceAliyunSlbBackendServer(),
			"alicloud_slb_domain_extension":          resourceAlicloudSlbDomainExtension(),
			"alicloud_slb_server_group":              resourceAliyunSlbServerGroup(),
			"alicloud_slb_master_slave_server_group": resourceAliyunSlbMasterSlaveServerGroup(),
			"alicloud_slb_rule":                      resourceAliyunSlbRule(),
			"alicloud_slb_acl":                       resourceAlicloudSlbAcl(),
			"alicloud_slb_ca_certificate":            resourceAlicloudSlbCACertificate(),
			"alicloud_slb_server_certificate":        resourceAlicloudSlbServerCertificate(),
			"alicloud_oss_bucket":                    resourceAlicloudOssBucket(),
			"alicloud_oss_bucket_object":             resourceAlicloudOssBucketObject(),
			"alicloud_ons_instance":                  resourceAlicloudOnsInstance(),
			"alicloud_ons_topic":                     resourceAlicloudOnsTopic(),
			"alicloud_ons_group":                     resourceAlicloudOnsGroup(),
			"alicloud_alikafka_consumer_group":       resourceAlicloudAlikafkaConsumerGroup(),
			"alicloud_alikafka_instance":             resourceAlicloudAlikafkaInstance(),
			"alicloud_alikafka_topic":                resourceAlicloudAlikafkaTopic(),
			"alicloud_alikafka_sasl_user":            resourceAlicloudAlikafkaSaslUser(),
			"alicloud_alikafka_sasl_acl":             resourceAlicloudAlikafkaSaslAcl(),
			"alicloud_dns_record":                    resourceAlicloudDnsRecord(),
			"alicloud_dns":                           resourceAlicloudDns(),
			"alicloud_dns_group":                     resourceAlicloudDnsGroup(),
			"alicloud_key_pair":                      resourceAlicloudKeyPair(),
			"alicloud_key_pair_attachment":           resourceAlicloudKeyPairAttachment(),
			"alicloud_kms_key":                       resourceAlicloudKmsKey(),
			"alicloud_kms_ciphertext":                resourceAlicloudKmsCiphertext(),
			"alicloud_ram_user":                      resourceAlicloudRamUser(),
			"alicloud_ram_account_password_policy":   resourceAlicloudRamAccountPasswordPolicy(),
			"alicloud_ram_access_key":                resourceAlicloudRamAccessKey(),
			"alicloud_ram_login_profile":             resourceAlicloudRamLoginProfile(),
			"alicloud_ram_group":                     resourceAlicloudRamGroup(),
			"alicloud_ram_role":                      resourceAlicloudRamRole(),
			"alicloud_ram_policy":                    resourceAlicloudRamPolicy(),
			// alicloud_ram_alias has been deprecated
			"alicloud_ram_alias":                           resourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_alias":                   resourceAlicloudRamAccountAlias(),
			"alicloud_ram_group_membership":                resourceAlicloudRamGroupMembership(),
			"alicloud_ram_user_policy_attachment":          resourceAlicloudRamUserPolicyAtatchment(),
			"alicloud_ram_role_policy_attachment":          resourceAlicloudRamRolePolicyAttachment(),
			"alicloud_ram_group_policy_attachment":         resourceAlicloudRamGroupPolicyAtatchment(),
			"alicloud_container_cluster":                   resourceAlicloudCSSwarm(),
			"alicloud_cs_application":                      resourceAlicloudCSApplication(),
			"alicloud_cs_swarm":                            resourceAlicloudCSSwarm(),
			"alicloud_cs_kubernetes":                       resourceAlicloudCSKubernetes(),
			"alicloud_cs_managed_kubernetes":               resourceAlicloudCSManagedKubernetes(),
			"alicloud_cs_serverless_kubernetes":            resourceAlicloudCSServerlessKubernetes(),
			"alicloud_cs_kubernetes_autoscaler":            resourceAlicloudCSKubernetesAutoscaler(),
			"alicloud_cr_namespace":                        resourceAlicloudCRNamespace(),
			"alicloud_cr_repo":                             resourceAlicloudCRRepo(),
			"alicloud_cr_ee_namespace":                     resourceAlicloudCrEENamespace(),
			"alicloud_cr_ee_repo":                          resourceAlicloudCrEERepo(),
			"alicloud_cr_ee_sync_rule":                     resourceAlicloudCrEESyncRule(),
			"alicloud_cdn_domain":                          resourceAlicloudCdnDomain(),
			"alicloud_cdn_domain_new":                      resourceAlicloudCdnDomainNew(),
			"alicloud_cdn_domain_config":                   resourceAlicloudCdnDomainConfig(),
			"alicloud_router_interface":                    resourceAlicloudRouterInterface(),
			"alicloud_router_interface_connection":         resourceAlicloudRouterInterfaceConnection(),
			"alicloud_ots_table":                           resourceAlicloudOtsTable(),
			"alicloud_ots_instance":                        resourceAlicloudOtsInstance(),
			"alicloud_ots_instance_attachment":             resourceAlicloudOtsInstanceAttachment(),
			"alicloud_cms_alarm":                           resourceAlicloudCmsAlarm(),
			"alicloud_cms_site_monitor":                    resourceAlicloudCmsSiteMonitor(),
			"alicloud_pvtz_zone":                           resourceAlicloudPvtzZone(),
			"alicloud_pvtz_zone_attachment":                resourceAlicloudPvtzZoneAttachment(),
			"alicloud_pvtz_zone_record":                    resourceAlicloudPvtzZoneRecord(),
			"alicloud_log_project":                         resourceAlicloudLogProject(),
			"alicloud_log_store":                           resourceAlicloudLogStore(),
			"alicloud_log_store_index":                     resourceAlicloudLogStoreIndex(),
			"alicloud_log_machine_group":                   resourceAlicloudLogMachineGroup(),
			"alicloud_logtail_config":                      resourceAlicloudLogtailConfig(),
			"alicloud_logtail_attachment":                  resourceAlicloudLogtailAttachment(),
			"alicloud_log_dashboard":                       resourceAlicloudLogDashboard(),
			"alicloud_log_alert":                           resourceAlicloudLogAlert(),
			"alicloud_log_audit":                           resourceAlicloudLogAudit(),
			"alicloud_fc_service":                          resourceAlicloudFCService(),
			"alicloud_fc_function":                         resourceAlicloudFCFunction(),
			"alicloud_fc_trigger":                          resourceAlicloudFCTrigger(),
			"alicloud_vpn_gateway":                         resourceAliyunVpnGateway(),
			"alicloud_vpn_customer_gateway":                resourceAliyunVpnCustomerGateway(),
			"alicloud_vpn_route_entry":                     resourceAliyunVpnRouteEntry(),
			"alicloud_vpn_connection":                      resourceAliyunVpnConnection(),
			"alicloud_ssl_vpn_server":                      resourceAliyunSslVpnServer(),
			"alicloud_ssl_vpn_client_cert":                 resourceAliyunSslVpnClientCert(),
			"alicloud_cen_instance":                        resourceAlicloudCenInstance(),
			"alicloud_cen_instance_attachment":             resourceAlicloudCenInstanceAttachment(),
			"alicloud_cen_bandwidth_package":               resourceAlicloudCenBandwidthPackage(),
			"alicloud_cen_bandwidth_package_attachment":    resourceAlicloudCenBandwidthPackageAttachment(),
			"alicloud_cen_bandwidth_limit":                 resourceAlicloudCenBandwidthLimit(),
			"alicloud_cen_route_entry":                     resourceAlicloudCenRouteEntry(),
			"alicloud_cen_instance_grant":                  resourceAlicloudCenInstanceGrant(),
			"alicloud_kvstore_instance":                    resourceAlicloudKVStoreInstance(),
			"alicloud_kvstore_backup_policy":               resourceAlicloudKVStoreBackupPolicy(),
			"alicloud_kvstore_account":                     resourceAlicloudKVstoreAccount(),
			"alicloud_datahub_project":                     resourceAlicloudDatahubProject(),
			"alicloud_datahub_subscription":                resourceAlicloudDatahubSubscription(),
			"alicloud_datahub_topic":                       resourceAlicloudDatahubTopic(),
			"alicloud_mns_queue":                           resourceAlicloudMNSQueue(),
			"alicloud_mns_topic":                           resourceAlicloudMNSTopic(),
			"alicloud_havip":                               resourceAliyunHaVip(),
			"alicloud_mns_topic_subscription":              resourceAlicloudMNSSubscription(),
			"alicloud_havip_attachment":                    resourceAliyunHaVipAttachment(),
			"alicloud_api_gateway_api":                     resourceAliyunApigatewayApi(),
			"alicloud_api_gateway_group":                   resourceAliyunApigatewayGroup(),
			"alicloud_api_gateway_app":                     resourceAliyunApigatewayApp(),
			"alicloud_api_gateway_app_attachment":          resourceAliyunApigatewayAppAttachment(),
			"alicloud_api_gateway_vpc_access":              resourceAliyunApigatewayVpc(),
			"alicloud_common_bandwidth_package":            resourceAliyunCommonBandwidthPackage(),
			"alicloud_common_bandwidth_package_attachment": resourceAliyunCommonBandwidthPackageAttachment(),
			"alicloud_drds_instance":                       resourceAlicloudDRDSInstance(),
			"alicloud_elasticsearch_instance":              resourceAlicloudElasticsearch(),
			"alicloud_actiontrail":                         resourceAlicloudActiontrail(),
			"alicloud_cas_certificate":                     resourceAlicloudCasCertificate(),
			"alicloud_ddoscoo_instance":                    resourceAlicloudDdoscooInstance(),
			"alicloud_ddosbgp_instance":                    resourceAlicloudDdosbgpInstance(),
			"alicloud_network_acl":                         resourceAliyunNetworkAcl(),
			"alicloud_network_acl_attachment":              resourceAliyunNetworkAclAttachment(),
			"alicloud_network_acl_entries":                 resourceAliyunNetworkAclEntries(),
			"alicloud_emr_cluster":                         resourceAlicloudEmrCluster(),
			"alicloud_cloud_connect_network":               resourceAlicloudCloudConnectNetwork(),
			"alicloud_cloud_connect_network_attachment":    resourceAlicloudCloudConnectNetworkAttachment(),
			"alicloud_cloud_connect_network_grant":         resourceAlicloudCloudConnectNetworkGrant(),
			"alicloud_sag_acl":                             resourceAlicloudSagAcl(),
			"alicloud_sag_acl_rule":                        resourceAlicloudSagAclRule(),
			"alicloud_sag_qos":                             resourceAlicloudSagQos(),
			"alicloud_sag_qos_policy":                      resourceAlicloudSagQosPolicy(),
			"alicloud_sag_qos_car":                         resourceAlicloudSagQosCar(),
			"alicloud_sag_snat_entry":                      resourceAlicloudSagSnatEntry(),
			"alicloud_sag_dnat_entry":                      resourceAlicloudSagDnatEntry(),
			"alicloud_sag_client_user":                     resourceAlicloudSagClientUser(),
			"alicloud_yundun_dbaudit_instance":             resourceAlicloudDbauditInstance(),
			"alicloud_yundun_bastionhost_instance":         resourceAlicloudBastionhostInstance(),
			"alicloud_polardb_cluster":                     resourceAlicloudPolarDBCluster(),
			"alicloud_polardb_backup_policy":               resourceAlicloudPolarDBBackupPolicy(),
			"alicloud_polardb_database":                    resourceAlicloudPolarDBDatabase(),
			"alicloud_polardb_account":                     resourceAlicloudPolarDBAccount(),
			"alicloud_polardb_account_privilege":           resourceAlicloudPolarDBAccountPrivilege(),
			"alicloud_polardb_endpoint":                    resourceAlicloudPolarDBEndpoint(),
			"alicloud_polardb_endpoint_address":            resourceAlicloudPolarDBEndpointAddress(),
			"alicloud_hbase_instance":                      resourceAlicloudHBaseInstance(),
			"alicloud_market_order":                        resourceAlicloudMarketOrder(),
			"alicloud_adb_cluster":                         resourceAlicloudAdbCluster(),
			"alicloud_adb_backup_policy":                   resourceAlicloudAdbBackupPolicy(),
			"alicloud_adb_account":                         resourceAlicloudAdbAccount(),
			"alicloud_adb_connection":                      resourceAlicloudAdbConnection(),
			"alicloud_cen_flowlog":                         resourceAlicloudCenFlowlog(),
			"alicloud_kms_secret":                          resourceAlicloudKmsSecret(),
			"alicloud_maxcompute_project":                  resourceAlicloudMaxComputeProject(),
			"alicloud_kms_alias":                           resourceAlicloudKmsAlias(),
			"alicloud_dns_instance":                        resourceAlicloudDnsInstance(),
			"alicloud_dns_domain_attachment":               resourceAlicloudDnsDomainAttachment(),
			"alicloud_edas_application":                    resourceAlicloudEdasApplication(),
			"alicloud_edas_deploy_group":                   resourceAlicloudEdasDeployGroup(),
			"alicloud_edas_application_scale":              resourceAlicloudEdasInstanceApplicationAttachment(),
			"alicloud_edas_slb_attachment":                 resourceAlicloudEdasSlbAttachment(),
			"alicloud_edas_cluster":                        resourceAlicloudEdasCluster(),
			"alicloud_edas_instance_cluster_attachment":    resourceAlicloudEdasInstanceClusterAttachment(),
			"alicloud_edas_application_deployment":         resourceAlicloudEdasApplicationPackageAttachment(),
			"alicloud_dns_domain":                          resourceAlicloudDnsDomain(),
			"alicloud_dms_enterprise_instance":             resourceAlicloudDmsEnterpriseInstance(),
			"alicloud_waf_domain":                          resourceAlicloudWafDomain(),
			"alicloud_cen_route_map":                       resourceAlicloudCenRouteMap(),
			"alicloud_resource_manager_role":               resourceAlicloudResourceManagerRole(),
			"alicloud_resource_manager_resource_group":     resourceAlicloudResourceManagerResourceGroup(),
			"alicloud_resource_manager_folder":             resourceAlicloudResourceManagerFolder(),
			"alicloud_resource_manager_handshake":          resourceAlicloudResourceManagerHandshake(),
			"alicloud_cen_private_zone":                    resourceAlicloudCenPrivateZone(),
			"alicloud_resource_manager_policy":             resourceAlicloudResourceManagerPolicy(),
			"alicloud_resource_manager_account":            resourceAlicloudResourceManagerAccount(),
			"alicloud_waf_instance":                        resourceAlicloudWafInstance(),
			"alicloud_resource_manager_resource_directory": resourceAlicloudResourceManagerResourceDirectory(),
			"alicloud_alidns_domain_group":                 resourceAlicloudAlidnsDomainGroup(),
			"alicloud_resource_manager_policy_version":     resourceAlicloudResourceManagerPolicyVersion(),
			"alicloud_kms_key_version":                     resourceAlicloudKmsKeyVersion(),
			"alicloud_alidns_record":                       resourceAlicloudAlidnsRecord(),
			"alicloud_ddoscoo_scheduler_rule":              resourceAlicloudDdoscooSchedulerRule(),
			"alicloud_cassandra_cluster":                   resourceAlicloudCassandraCluster(),
			"alicloud_cassandra_data_center":               resourceAlicloudCassandraDataCenter(),
			"alicloud_cen_vbr_health_check":                resourceAlicloudCenVbrHealthCheck(),
			"alicloud_eci_openapi_image_cache":             resourceAlicloudEciImageCache(),
			"alicloud_eci_image_cache":                     resourceAlicloudEciImageCache(),
			"alicloud_dms_enterprise_user":                 resourceAlicloudDmsEnterpriseUser(),
			"alicloud_ecs_dedicated_host":                  resourceAlicloudEcsDedicatedHost(),
			"alicloud_oos_template":                        resourceAlicloudOosTemplate(),
			"alicloud_edas_k8s_cluster":                    resourceAlicloudEdasK8sCluster(),
			"alicloud_oos_execution":                       resourceAlicloudOosExecution(),
			"alicloud_resource_manager_policy_attachment":  resourceAlicloudResourceManagerPolicyAttachment(),
			"alicloud_edas_k8s_application":                resourceAlicloudEdasK8sApplication(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var providerConfig map[string]interface{}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var getProviderConfig = func(str string, key string) string {
		if str == "" {
			value, err := getConfigFromProfile(d, key)
			if err == nil && value != nil {
				str = value.(string)
			}
		}
		return str
	}

	accessKey := getProviderConfig(d.Get("access_key").(string), "access_key_id")
	secretKey := getProviderConfig(d.Get("secret_key").(string), "access_key_secret")
	region := getProviderConfig(d.Get("region").(string), "region_id")
	if region == "" {
		region = DEFAULT_REGION
	}

	ecsRoleName := getProviderConfig(d.Get("ecs_role_name").(string), "ram_role_name")

	config := &connectivity.Config{
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		EcsRoleName:          strings.TrimSpace(ecsRoleName),
		Region:               connectivity.Region(strings.TrimSpace(region)),
		RegionId:             strings.TrimSpace(region),
		SkipRegionValidation: d.Get("skip_region_validation").(bool),
		ConfigurationSource:  d.Get("configuration_source").(string),
		Protocol:             d.Get("protocol").(string),
	}
	token := getProviderConfig(d.Get("security_token").(string), "sts_token")
	config.SecurityToken = strings.TrimSpace(token)

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

		log.Printf("[INFO] assume_role configuration set: (RamRoleArn: %q, RamRoleSessionName: %q, RamRolePolicy: %q, RamRoleSessionExpiration: %d)",
			config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration)
	}

	if err := config.MakeConfigByEcsRoleName(); err != nil {
		return nil, err
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		config.EcsEndpoint = strings.TrimSpace(endpoints["ecs"].(string))
		config.RdsEndpoint = strings.TrimSpace(endpoints["rds"].(string))
		config.SlbEndpoint = strings.TrimSpace(endpoints["slb"].(string))
		config.VpcEndpoint = strings.TrimSpace(endpoints["vpc"].(string))
		config.CenEndpoint = strings.TrimSpace(endpoints["cen"].(string))
		config.EssEndpoint = strings.TrimSpace(endpoints["ess"].(string))
		config.OssEndpoint = strings.TrimSpace(endpoints["oss"].(string))
		config.OssEndpoint = strings.TrimSpace(endpoints["ons"].(string))
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
		config.ActionTrailEndpoint = strings.TrimSpace(endpoints["actiontrail"].(string))
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
		if endpoint, ok := endpoints["alidns"]; ok {
			config.AlidnsEndpoint = strings.TrimSpace(endpoint.(string))
		} else {
			config.AlidnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
		}
		config.CassandraEndpoint = strings.TrimSpace(endpoints["cassandra"].(string))
	}

	if config.RamRoleArn != "" {
		config.AccessKey, config.SecretKey, config.SecurityToken, err = getAssumeRoleAK(config.AccessKey, config.SecretKey, config.SecurityToken, region, config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration, config.StsEndpoint)
		if err != nil {
			return nil, err
		}
	}

	if ots_instance_name, ok := d.GetOk("ots_instance_name"); ok && ots_instance_name.(string) != "" {
		config.OtsInstanceName = strings.TrimSpace(ots_instance_name.(string))
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

	if config.ConfigurationSource == "" {
		sourceName := fmt.Sprintf("Default/%s:%s", config.AccessKey, strings.Trim(uuid.New().String(), "-"))
		if len(sourceName) > 64 {
			sourceName = sourceName[:64]
		}
		config.ConfigurationSource = sourceName
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

		"ecs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.",

		"rds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.",

		"slb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.",

		"vpc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.",

		"cen_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CEN endpoints.",

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
					DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ASSUME_ROLE_ARN", ""),
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
					DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ASSUME_ROLE_SESSION_NAME", ""),
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
					ValidateFunc: intBetween(900, 3600),
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
				"cen": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cen_endpoint"],
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
	buf.WriteString(fmt.Sprintf("%s-", m["cen"].(string)))
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
	return hashcode.String(buf.String())
}

func getConfigFromProfile(d *schema.ResourceData, ProfileKey string) (interface{}, error) {

	if providerConfig == nil {
		if v, ok := d.GetOk("profile"); !ok && v.(string) == "" {
			return nil, nil
		}
		current := d.Get("profile").(string)
		// Set CredsFilename, expanding home directory
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

func getAssumeRoleAK(accessKey, secretKey, stsToken, region, roleArn, sessionName, policy string, sessionExpiration int, stsEndpoint string) (string, string, string, error) {
	request := sts.CreateAssumeRoleRequest()
	request.RoleArn = roleArn
	request.RoleSessionName = sessionName
	request.DurationSeconds = requests.NewInteger(sessionExpiration)
	request.Policy = policy
	request.Scheme = "https"
	request.Domain = stsEndpoint

	var client *sts.Client
	var err error
	if stsToken == "" {
		client, err = sts.NewClientWithAccessKey(region, accessKey, secretKey)
	} else {
		client, err = sts.NewClientWithStsToken(region, accessKey, secretKey, stsToken)
	}

	if err != nil {
		return "", "", "", err
	}

	response, err := client.AssumeRole(request)
	if err != nil {
		return "", "", "", err
	}

	return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
}
