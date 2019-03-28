package alicloud

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

			"fc": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'fc' has been deprecated from provider version 1.28.0. New field 'fc' which in nested endpoints instead.",
			},
			"endpoints": endpointsSchema(),
		},
		DataSourcesMap: map[string]*schema.Resource{

			"alicloud_account":            dataSourceAlicloudAccount(),
			"alicloud_images":             dataSourceAlicloudImages(),
			"alicloud_regions":            dataSourceAlicloudRegions(),
			"alicloud_zones":              dataSourceAlicloudZones(),
			"alicloud_instance_types":     dataSourceAlicloudInstanceTypes(),
			"alicloud_instances":          dataSourceAlicloudInstances(),
			"alicloud_disks":              dataSourceAlicloudDisks(),
			"alicloud_network_interfaces": dataSourceAlicloudNetworkInterfaces(),
			"alicloud_vpcs":               dataSourceAlicloudVpcs(),
			"alicloud_vswitches":          dataSourceAlicloudVSwitches(),
			"alicloud_eips":               dataSourceAlicloudEips(),
			"alicloud_key_pairs":          dataSourceAlicloudKeyPairs(),
			"alicloud_kms_keys":           dataSourceAlicloudKmsKeys(),
			"alicloud_dns_domains":        dataSourceAlicloudDnsDomains(),
			"alicloud_dns_groups":         dataSourceAlicloudDnsGroups(),
			"alicloud_dns_records":        dataSourceAlicloudDnsRecords(),
			// alicloud_dns_domain_groups, alicloud_dns_domain_records have been deprecated.
			"alicloud_dns_domain_groups":  dataSourceAlicloudDnsGroups(),
			"alicloud_dns_domain_records": dataSourceAlicloudDnsRecords(),
			// alicloud_ram_account_alias has been deprecated
			"alicloud_ram_account_alias":              dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_aliases":            dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_groups":                     dataSourceAlicloudRamGroups(),
			"alicloud_ram_users":                      dataSourceAlicloudRamUsers(),
			"alicloud_ram_roles":                      dataSourceAlicloudRamRoles(),
			"alicloud_ram_policies":                   dataSourceAlicloudRamPolicies(),
			"alicloud_security_groups":                dataSourceAlicloudSecurityGroups(),
			"alicloud_security_group_rules":           dataSourceAlicloudSecurityGroupRules(),
			"alicloud_slbs":                           dataSourceAlicloudSlbs(),
			"alicloud_slb_attachments":                dataSourceAlicloudSlbAttachments(),
			"alicloud_slb_listeners":                  dataSourceAlicloudSlbListeners(),
			"alicloud_slb_rules":                      dataSourceAlicloudSlbRules(),
			"alicloud_slb_server_groups":              dataSourceAlicloudSlbServerGroups(),
			"alicloud_slb_acls":                       dataSourceAlicloudSlbAcls(),
			"alicloud_slb_server_certificates":        dataSourceAlicloudSlbServerCertificates(),
			"alicloud_slb_ca_certificates":            dataSourceAlicloudSlbCACertificates(),
			"alicloud_oss_bucket_objects":             dataSourceAlicloudOssBucketObjects(),
			"alicloud_oss_buckets":                    dataSourceAlicloudOssBuckets(),
			"alicloud_fc_functions":                   dataSourceAlicloudFcFunctions(),
			"alicloud_fc_services":                    dataSourceAlicloudFcServices(),
			"alicloud_fc_triggers":                    dataSourceAlicloudFcTriggers(),
			"alicloud_db_instances":                   dataSourceAlicloudDBInstances(),
			"alicloud_pvtz_zones":                     dataSourceAlicloudPvtzZones(),
			"alicloud_pvtz_zone_records":              dataSourceAlicloudPvtzZoneRecords(),
			"alicloud_router_interfaces":              dataSourceAlicloudRouterInterfaces(),
			"alicloud_vpn_gateways":                   dataSourceAlicloudVpnGateways(),
			"alicloud_vpn_customer_gateways":          dataSourceAlicloudVpnCustomerGateways(),
			"alicloud_vpn_connections":                dataSourceAlicloudVpnConnections(),
			"alicloud_mongo_instances":                dataSourceAlicloudMongoInstances(),
			"alicloud_kvstore_instances":              dataSourceAlicloudKVStoreInstances(),
			"alicloud_cen_instances":                  dataSourceAlicloudCenInstances(),
			"alicloud_cen_bandwidth_packages":         dataSourceAlicloudCenBandwidthPackages(),
			"alicloud_cen_bandwidth_limits":           dataSourceAlicloudCenBandwidthLimits(),
			"alicloud_cen_route_entries":              dataSourceAlicloudCenRouteEntries(),
			"alicloud_cen_region_route_entries":       dataSourceAlicloudCenRegionRouteEntries(),
			"alicloud_cs_kubernetes_clusters":         dataSourceAlicloudCSKubernetesClusters(),
			"alicloud_cs_managed_kubernetes_clusters": dataSourceAlicloudCSManagerKubernetesClusters(),
			"alicloud_cr_namespaces":                  dataSourceAlicloudCRNamespaces(),
			"alicloud_cr_repos":                       dataSourceAlicloudCRRepos(),
			"alicloud_mns_queues":                     dataSourceAlicloudMNSQueues(),
			"alicloud_mns_topics":                     dataSourceAlicloudMNSTopics(),
			"alicloud_mns_topic_subscriptions":        dataSourceAlicloudMNSTopicSubscriptions(),
			"alicloud_api_gateway_apis":               dataSourceAlicloudApiGatewayApis(),
			"alicloud_api_gateway_groups":             dataSourceAlicloudApiGatewayGroups(),
			"alicloud_api_gateway_apps":               dataSourceAlicloudApiGatewayApps(),
			"alicloud_elasticsearch_instances":        dataSourceAlicloudElasticsearch(),
			"alicloud_drds_instances":                 dataSourceAlicloudDRDSInstances(),
			"alicloud_nas_access_groups":              dataSourceAlicloudAccessGroups(),
			"alicloud_nas_access_rules":               dataSourceAlicloudAccessRules(),
			"alicloud_nas_mount_targets":              dataSourceAlicloudMountTargets(),
			"alicloud_nas_file_systems":               dataSourceAlicloudFileSystems(),
			"alicloud_cas_certificates":               dataSourceAlicloudCasCertificates(),
			"alicloud_actiontrails":                   dataSourceAlicloudActiontrails(),
			"alicloud_common_bandwidth_packages":      dataSourceAlicloudCommonBandwidthPackages(),
			"alicloud_route_tables":                   dataSourceAlicloudRouteTables(),
			"alicloud_route_entries":                  dataSourceAlicloudRouteEntries(),
			"alicloud_nat_gateways":                   dataSourceAlicloudNatGateways(),
			"alicloud_snat_entries":                   dataSourceAlicloudSnatEntries(),
			"alicloud_forward_entries":                dataSourceAlicloudForwardEntries(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_instance":                           resourceAliyunInstance(),
			"alicloud_ram_role_attachment":                resourceAlicloudRamRoleAttachment(),
			"alicloud_disk":                               resourceAliyunDisk(),
			"alicloud_disk_attachment":                    resourceAliyunDiskAttachment(),
			"alicloud_network_interface":                  resourceAliyunNetworkInterface(),
			"alicloud_network_interface_attachment":       resourceAliyunNetworkInterfaceAttachment(),
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
			"alicloud_db_readonly_instance":               resourceAlicloudDBReadonlyInstance(),
			"alicloud_ess_scaling_group":                  resourceAlicloudEssScalingGroup(),
			"alicloud_ess_scaling_configuration":          resourceAlicloudEssScalingConfiguration(),
			"alicloud_ess_scaling_rule":                   resourceAlicloudEssScalingRule(),
			"alicloud_ess_schedule":                       resourceAlicloudEssSchedule(),
			"alicloud_ess_attachment":                     resourceAlicloudEssAttachment(),
			"alicloud_ess_lifecycle_hook":                 resourceAlicloudEssLifecycleHook(),
			"alicloud_ess_alarm":                          resourceAlicloudEssAlarm(),
			"alicloud_vpc":                                resourceAliyunVpc(),
			"alicloud_nat_gateway":                        resourceAliyunNatGateway(),
			"alicloud_nas_file_system":                    resourceAlicloudNasFileSystem(),
			"alicloud_nas_mount_target":                   resourceAlicloudNasMountTarget(),
			"alicloud_nas_access_group":                   resourceAlicloudNasAccessGroup(),
			"alicloud_nas_access_rule":                    resourceAlicloudNasAccessRule(),
			// "alicloud_subnet" aims to match aws usage habit.
			"alicloud_subnet":                 resourceAliyunSubnet(),
			"alicloud_vswitch":                resourceAliyunSubnet(),
			"alicloud_route_entry":            resourceAliyunRouteEntry(),
			"alicloud_route_table":            resourceAliyunRouteTable(),
			"alicloud_route_table_attachment": resourceAliyunRouteTableAttachment(),
			"alicloud_snat_entry":             resourceAliyunSnatEntry(),
			"alicloud_forward_entry":          resourceAliyunForwardEntry(),
			"alicloud_eip":                    resourceAliyunEip(),
			"alicloud_eip_association":        resourceAliyunEipAssociation(),
			"alicloud_slb":                    resourceAliyunSlb(),
			"alicloud_slb_listener":           resourceAliyunSlbListener(),
			"alicloud_slb_attachment":         resourceAliyunSlbAttachment(),
			"alicloud_slb_server_group":       resourceAliyunSlbServerGroup(),
			"alicloud_slb_rule":               resourceAliyunSlbRule(),
			"alicloud_slb_acl":                resourceAlicloudSlbAcl(),
			"alicloud_slb_ca_certificate":     resourceAlicloudSlbCACertificate(),
			"alicloud_slb_server_certificate": resourceAlicloudSlbServerCertificate(),
			"alicloud_oss_bucket":             resourceAlicloudOssBucket(),
			"alicloud_oss_bucket_object":      resourceAlicloudOssBucketObject(),
			"alicloud_dns_record":             resourceAlicloudDnsRecord(),
			"alicloud_dns":                    resourceAlicloudDns(),
			"alicloud_dns_group":              resourceAlicloudDnsGroup(),
			"alicloud_key_pair":               resourceAlicloudKeyPair(),
			"alicloud_key_pair_attachment":    resourceAlicloudKeyPairAttachment(),
			"alicloud_kms_key":                resourceAlicloudKmsKey(),
			"alicloud_ram_user":               resourceAlicloudRamUser(),
			"alicloud_ram_access_key":         resourceAlicloudRamAccessKey(),
			"alicloud_ram_login_profile":      resourceAlicloudRamLoginProfile(),
			"alicloud_ram_group":              resourceAlicloudRamGroup(),
			"alicloud_ram_role":               resourceAlicloudRamRole(),
			"alicloud_ram_policy":             resourceAlicloudRamPolicy(),
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
			"alicloud_cr_namespace":                        resourceAlicloudCRNamespace(),
			"alicloud_cr_repo":                             resourceAlicloudCRRepo(),
			"alicloud_cdn_domain":                          resourceAlicloudCdnDomain(),
			"alicloud_cdn_domain_new":                      resourceAlicloudCdnDomainNew(),
			"alicloud_cdn_domain_config":                   resourceAlicloudCdnDomainConfig(),
			"alicloud_router_interface":                    resourceAlicloudRouterInterface(),
			"alicloud_router_interface_connection":         resourceAlicloudRouterInterfaceConnection(),
			"alicloud_ots_table":                           resourceAlicloudOtsTable(),
			"alicloud_ots_instance":                        resourceAlicloudOtsInstance(),
			"alicloud_ots_instance_attachment":             resourceAlicloudOtsInstanceAttachment(),
			"alicloud_cms_alarm":                           resourceAlicloudCmsAlarm(),
			"alicloud_pvtz_zone":                           resourceAlicloudPvtzZone(),
			"alicloud_pvtz_zone_attachment":                resourceAlicloudPvtzZoneAttachment(),
			"alicloud_pvtz_zone_record":                    resourceAlicloudPvtzZoneRecord(),
			"alicloud_log_project":                         resourceAlicloudLogProject(),
			"alicloud_log_store":                           resourceAlicloudLogStore(),
			"alicloud_log_store_index":                     resourceAlicloudLogStoreIndex(),
			"alicloud_log_machine_group":                   resourceAlicloudLogMachineGroup(),
			"alicloud_logtail_config":                      resourceAlicloudLogtailConfig(),
			"alicloud_logtail_attachment":                  resourceAlicloudLogtailAttachment(),
			"alicloud_fc_service":                          resourceAlicloudFCService(),
			"alicloud_fc_function":                         resourceAlicloudFCFunction(),
			"alicloud_fc_trigger":                          resourceAlicloudFCTrigger(),
			"alicloud_vpn_gateway":                         resourceAliyunVpnGateway(),
			"alicloud_vpn_customer_gateway":                resourceAliyunVpnCustomerGateway(),
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
			"alicloud_ddoscoo_instance":                    resourceAlicloudDdoscoo(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	region, ok := d.GetOk("region")
	if !ok {
		if region == "" {
			region = DEFAULT_REGION
		}
	}
	config := connectivity.Config{
		AccessKey:   strings.TrimSpace(d.Get("access_key").(string)),
		SecretKey:   strings.TrimSpace(d.Get("secret_key").(string)),
		EcsRoleName: strings.TrimSpace(d.Get("ecs_role_name").(string)),
		Region:      connectivity.Region(strings.TrimSpace(region.(string))),
		RegionId:    strings.TrimSpace(region.(string)),
	}

	if token, ok := d.GetOk("security_token"); ok && token.(string) != "" {
		config.SecurityToken = strings.TrimSpace(token.(string))
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
		config.KVStoreEndpoint = strings.TrimSpace(endpoints["kvstore"].(string))
		config.FcEndpoint = strings.TrimSpace(endpoints["fc"].(string))
		config.ApigatewayEndpoint = strings.TrimSpace(endpoints["apigateway"].(string))
		config.DatahubEndpoint = strings.TrimSpace(endpoints["datahub"].(string))
		config.MnsEndpoint = strings.TrimSpace(endpoints["mns"].(string))
		config.LocationEndpoint = strings.TrimSpace(endpoints["location"].(string))
		config.ElasticsearchEndpoint = strings.TrimSpace(endpoints["elasticsearch"].(string))
		config.NasEndpoint = strings.TrimSpace(endpoints["nas"].(string))
		config.ActionTrailEndpoint = strings.TrimSpace(endpoints["actiontrail"].(string))
		config.BssOpenApiEndpoint = strings.TrimSpace(endpoints["bssopenapi"].(string))
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

		"ecs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.",

		"rds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.",

		"slb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.",

		"vpc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.",

		"cen_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CEN endpoints.",

		"ess_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Autoscaling endpoints.",

		"oss_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom OSS endpoints.",

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
	}
}

func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
	buf.WriteString(fmt.Sprintf("%s-", m["kvstore"].(string)))
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
	return hashcode.String(buf.String())
}
