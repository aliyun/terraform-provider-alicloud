package alicloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/credentials-go/credentials/providers"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/go-homedir"
)

// Provider returns a schema.Provider for alicloud
func Provider() *schema.Provider {
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
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Domestic", "International"}, true),
				DefaultFunc:  schema.MultiEnvDefaultFunc([]string{"ALIBABA_CLOUD_ACCOUNT_TYPE"}, nil),
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
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SKIP_REGION_VALIDATION", false),
				Description: descriptions["skip_region_validation"],
			},
			"features": featuresSchema(),
			"configuration_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["configuration_source"],
				ValidateFunc: validation.StringLenBetween(0, 1024),
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
			// account
			"alicloud_account": dataSourceAlicloudAccount(),

			// ack

			// adb resources and data sources are registered via service/adb.ServicePackage()
			// (provider/registry); see alicloud/service/adb/register.go.

			// alikafka

			// amqp

			// api — registered via alicloud/service/apigateway.ServicePackage()

			// apig

			// bastionhost
			"alicloud_bastionhost_instances":       dataSourceAlicloudBastionhostInstances(),
			"alicloud_bastionhost_user_groups":     dataSourceAlicloudBastionhostUserGroups(),
			"alicloud_bastionhost_users":           dataSourceAlicloudBastionhostUsers(),
			"alicloud_bastionhost_host_groups":     dataSourceAlicloudBastionhostHostGroups(),
			"alicloud_bastionhost_hosts":           dataSourceAlicloudBastionhostHosts(),
			"alicloud_bastionhost_host_accounts":   dataSourceAlicloudBastionhostHostAccounts(),
			"alicloud_bastionhost_host_share_keys": dataSourceAlicloudBastionhostHostShareKeys(),

			// brain

			// bss

			// caller
			"alicloud_caller_identity": dataSourceAlicloudCallerIdentity(),

			// cassandra

			// chatbot

			// click
			// clickhouse resources/data sources migrated to alicloud/service/clickhouse

			// cloud
			"alicloud_cloud_monitor_service_metric_alarm_rules":   dataSourceAliCloudCloudMonitorServiceMetricAlarmRules(),
			"alicloud_cloud_monitor_service_hybrid_double_writes": dataSourceAliCloudCloudMonitorServiceHybridDoubleWrites(),

			// cms
			"alicloud_cms_event_notify_policies":     dataSourceAliCloudCmsEventNotifyPolicies(),
			"alicloud_cms_alert_rules_v2":            dataSourceAliCloudCmsAlertRulesV2(),
			"alicloud_cms_service":                   dataSourceAlicloudCmsService(),
			"alicloud_cms_alarm_contacts":            dataSourceAlicloudCmsAlarmContacts(),
			"alicloud_cms_alarm_contact_groups":      dataSourceAlicloudCmsAlarmContactGroups(),
			"alicloud_cms_group_metric_rules":        dataSourceAlicloudCmsGroupMetricRules(),
			"alicloud_cms_monitor_groups":            dataSourceAlicloudCmsMonitorGroups(),
			"alicloud_cms_monitor_group_instanceses": dataSourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_cms_monitor_group_instances":   dataSourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_cms_metric_rule_templates":     dataSourceAlicloudCmsMetricRuleTemplates(),
			"alicloud_cms_dynamic_tag_groups":        dataSourceAlicloudCmsDynamicTagGroups(),
			"alicloud_cms_namespaces":                dataSourceAlicloudCmsNamespaces(),
			"alicloud_cms_datasets":                  dataSourceAliCloudCmsDatasets(),
			"alicloud_cms_sls_groups":                dataSourceAlicloudCmsSlsGroups(),
			"alicloud_cms_hybrid_monitor_datas":      dataSourceAlicloudCmsHybridMonitorDatas(),
			"alicloud_cms_hybrid_monitor_sls_tasks":  dataSourceAlicloudCmsHybridMonitorSlsTasks(),
			"alicloud_cms_hybrid_monitor_fc_tasks":   dataSourceAlicloudCmsHybridMonitorFcTasks(),
			"alicloud_cms_event_rules":               dataSourceAlicloudCmsEventRules(),
			"alicloud_cms_metric_rule_black_lists":   dataSourceAlicloudCmsMetricRuleBlackLists(),
			"alicloud_cms_site_monitors":             dataSourceAliCloudCloudMonitorServiceSiteMonitors(),

			// compute

			// config
			// cr
			"alicloud_cr_artifact_subscription_rules": dataSourceAliCloudCrArtifactSubscriptionRules(),
			"alicloud_cr_artifact_lifecycle_rules":    dataSourceAliCloudCrArtifactLifecycleRules(),
			"alicloud_cr_namespaces":                  dataSourceAlicloudCRNamespaces(),
			"alicloud_cr_repos":                       dataSourceAlicloudCRRepos(),
			"alicloud_cr_ee_instances":                dataSourceAlicloudCrEEInstances(),
			"alicloud_cr_ee_namespaces":               dataSourceAlicloudCrEENamespaces(),
			"alicloud_cr_ee_repos":                    dataSourceAliCloudCrEERepos(),
			"alicloud_cr_ee_sync_rules":               dataSourceAlicloudCrEESyncRules(),
			"alicloud_cr_service":                     dataSourceAlicloudCrService(),
			"alicloud_cr_endpoint_acl_policies":       dataSourceAlicloudCrEndpointAclPolicies(),
			"alicloud_cr_endpoint_acl_service":        dataSourceAlicloudCrEndpointAclService(),
			"alicloud_cr_internet_endpoint":           dataSourceAlicloudCrInternetEndpoint(),
			"alicloud_cr_chart_namespaces":            dataSourceAlicloudCrChartNamespaces(),
			"alicloud_cr_chart_repositories":          dataSourceAlicloudCrChartRepositories(),
			"alicloud_cr_chains":                      dataSourceAlicloudCrChains(),
			"alicloud_cr_vpc_endpoint_linked_vpcs":    dataSourceAlicloudCrVpcEndpointLinkedVpcs(),

			// cs

			// das

			// data

			// database

			// dataworks

			// db (registered in alicloud/service/rds)

			// ddoscoo

			// direct

			// disks
			"alicloud_disks": dataSourceAliCloudEcsDisks(),

			// dms
			"alicloud_dms_enterprise_instances":       dataSourceAlicloudDmsEnterpriseInstances(),
			"alicloud_dms_enterprise_users":           dataSourceAlicloudDmsEnterpriseUsers(),
			"alicloud_dms_user_tenants":               dataSourceAlicloudDmsUserTenants(),
			"alicloud_dms_enterprise_proxies":         dataSourceAlicloudDmsEnterpriseProxies(),
			"alicloud_dms_enterprise_proxy_accesses":  dataSourceAlicloudDmsEnterpriseProxyAccesses(),
			"alicloud_dms_enterprise_logic_databases": dataSourceAlicloudDmsEnterpriseLogicDatabases(),
			"alicloud_dms_enterprise_databases":       dataSourceAlicloudDmsEnterpriseDatabases(),

			// drds

			// dts
			// (registered via dts.ServicePackage)

			// ecd

			// ecp

			// ecs
			"alicloud_ecs_dedicated_hosts":               dataSourceAliCloudEcsDedicatedHosts(),
			"alicloud_ecs_commands":                      dataSourceAlicloudEcsCommands(),
			"alicloud_ecs_hpc_clusters":                  dataSourceAlicloudEcsHpcClusters(),
			"alicloud_ecs_auto_snapshot_policies":        dataSourceAlicloudEcsAutoSnapshotPolicies(),
			"alicloud_ecs_launch_templates":              dataSourceAlicloudEcsLaunchTemplates(),
			"alicloud_ecs_snapshots":                     dataSourceAliCloudEcsSnapshots(),
			"alicloud_ecs_key_pairs":                     dataSourceAliCloudEcsKeyPairs(),
			"alicloud_ecs_disks":                         dataSourceAliCloudEcsDisks(),
			"alicloud_ecs_network_interfaces":            dataSourceAlicloudEcsNetworkInterfaces(),
			"alicloud_ecs_deployment_sets":               dataSourceAlicloudEcsDeploymentSets(),
			"alicloud_ecs_dedicated_host_clusters":       dataSourceAlicloudEcsDedicatedHostClusters(),
			"alicloud_ecs_prefix_lists":                  dataSourceAlicloudEcsPrefixLists(),
			"alicloud_ecs_storage_capacity_units":        dataSourceAlicloudEcsStorageCapacityUnits(),
			"alicloud_ecs_image_components":              dataSourceAlicloudEcsImageComponents(),
			"alicloud_ecs_snapshot_groups":               dataSourceAlicloudEcsSnapshotGroups(),
			"alicloud_ecs_image_pipelines":               dataSourceAlicloudEcsImagePipelines(),
			"alicloud_ecs_network_interface_permissions": dataSourceAlicloudEcsNetworkInterfacePermissions(),
			"alicloud_ecs_invocations":                   dataSourceAlicloudEcsInvocations(),
			"alicloud_ecs_activations":                   dataSourceAlicloudEcsActivations(),
			"alicloud_ecs_capacity_reservations":         dataSourceAlicloudEcsCapacityReservations(),
			"alicloud_ecs_elasticity_assurances":         dataSourceAlicloudEcsElasticityAssurances(),

			// edas

			// ehpc

			// eip

			// eips

			// elasticsearch

			// enhanced
			"alicloud_enhanced_nat_available_zones": dataSourceAlicloudEnhancedNatAvailableZones(),

			// ens

			// esa

			// event

			// express

			// fc

			// fcv3 (migrated to alicloud/service/fcv3; registered via registry)

			// file
			"alicloud_file_crc64_checksum": dataSourceAlicloudFileCRC64Checksum(),

			// forward

			// ga
			// ga types are registered via service/ga.ServicePackage().

			// governance

			// gpdb (registered via alicloud/service/gpdb.ServicePackage)
			// graph (gdb migrated to alicloud/service/gdb)

			// havips
			"alicloud_havips": dataSourceAlicloudHavips(),

			// hbase

			// hbr

			// images
			"alicloud_images": dataSourceAlicloudImages(),

			// instance
			"alicloud_instance_type_families": dataSourceAlicloudInstanceTypeFamilies(),
			"alicloud_instance_types":         dataSourceAliCloudInstanceTypes(),

			// instances
			"alicloud_instances": dataSourceAlicloudInstances(),

			// key
			"alicloud_key_pairs": dataSourceAliCloudEcsKeyPairs(),

			// kms

			// kvstore

			// market

			// maxcompute

			// message
			"alicloud_message_service_queues":        dataSourceAlicloudMessageServiceQueues(),
			"alicloud_message_service_topics":        dataSourceAlicloudMessageServiceTopics(),
			"alicloud_message_service_subscriptions": dataSourceAlicloudMessageServiceSubscriptions(),

			// mns

			// mse

			// nas

			// nat

			// network
			"alicloud_network_interfaces": dataSourceAlicloudEcsNetworkInterfaces(),
			"alicloud_network_acls":       dataSourceAlicloudNetworkAcls(),

			// ocean
			// alicloud_ocean_base_instances is declared by service/oceanbasepro.ServicePackage().

			// oos

			// oss
			// ots

			// polardb

			// privatelink

			// pvtz

			// quick
			// quickbi types are registered via service/quickbi.ServicePackage().

			// quotas

			// ram

			// rds (registered in alicloud/service/rds)

			// realtime

			// redis

			// resource
			"alicloud_resource_manager_folders":                     dataSourceAliCloudResourceManagerFolders(),
			"alicloud_resource_manager_policies":                    dataSourceAlicloudResourceManagerPolicies(),
			"alicloud_resource_manager_resource_groups":             dataSourceAlicloudResourceManagerResourceGroups(),
			"alicloud_resource_manager_roles":                       dataSourceAlicloudResourceManagerRoles(),
			"alicloud_resource_manager_policy_versions":             dataSourceAlicloudResourceManagerPolicyVersions(),
			"alicloud_resource_manager_accounts":                    dataSourceAliCloudResourceManagerAccounts(),
			"alicloud_resource_manager_resource_directories":        dataSourceAlicloudResourceManagerResourceDirectories(),
			"alicloud_resource_manager_handshakes":                  dataSourceAlicloudResourceManagerHandshakes(),
			"alicloud_resource_manager_policy_attachments":          dataSourceAlicloudResourceManagerPolicyAttachments(),
			"alicloud_resource_manager_resource_shares":             dataSourceAlicloudResourceManagerResourceShares(),
			"alicloud_resource_manager_shared_resources":            dataSourceAliCloudResourceManagerSharedResources(),
			"alicloud_resource_manager_shared_targets":              dataSourceAliCloudResourceManagerSharedTargets(),
			"alicloud_resource_manager_control_policies":            dataSourceAlicloudResourceManagerControlPolicies(),
			"alicloud_resource_manager_control_policy_attachments":  dataSourceAlicloudResourceManagerControlPolicyAttachments(),
			"alicloud_resource_manager_delegated_administrators":    dataSourceAlicloudResourceManagerDelegatedAdministrators(),
			"alicloud_resource_manager_account_deletion_check_task": dataSourceAlicloudResourceManagerAccountDeletionCheckTask(),

			// ros
			// route
			"alicloud_route_tables":  dataSourceAlicloudRouteTables(),
			"alicloud_route_entries": dataSourceAlicloudRouteEntries(),

			// router

			// sae

			// sag

			// scdn

			// security
			"alicloud_security_groups":      dataSourceAlicloudSecurityGroups(),
			"alicloud_security_group_rules": dataSourceAlicloudSecurityGroupRules(),

			// service

			// slb
			"alicloud_slb_attachments":                dataSourceAlicloudSlbAttachments(),
			"alicloud_slb_backend_servers":            dataSourceAlicloudSlbBackendServers(),
			"alicloud_slb_listeners":                  dataSourceAlicloudSlbListeners(),
			"alicloud_slb_rules":                      dataSourceAlicloudSlbRules(),
			"alicloud_slb_server_groups":              dataSourceAlicloudSlbServerGroups(),
			"alicloud_slb_master_slave_server_groups": dataSourceAlicloudSlbMasterSlaveServerGroups(),
			"alicloud_slb_acls":                       dataSourceAlicloudSlbAcls(),
			"alicloud_slb_server_certificates":        dataSourceAlicloudSlbServerCertificates(),
			"alicloud_slb_ca_certificates":            dataSourceAlicloudSlbCaCertificates(),
			"alicloud_slb_domain_extensions":          dataSourceAlicloudSlbDomainExtensions(),
			"alicloud_slb_zones":                      dataSourceAlicloudSlbZones(),
			"alicloud_slb_load_balancers":             dataSourceAlicloudSlbLoadBalancers(),
			"alicloud_slb_tls_cipher_policies":        dataSourceAlicloudSlbTlsCipherPolicies(),

			// slbs
			"alicloud_slbs": dataSourceAlicloudSlbLoadBalancers(),

			// smartag

			// snapshots
			"alicloud_snapshots": dataSourceAliCloudEcsSnapshots(),

			// snat

			// ssl
			"alicloud_ssl_vpn_servers":      dataSourceAlicloudSslVpnServers(),
			"alicloud_ssl_vpn_client_certs": dataSourceAlicloudSslVpnClientCerts(),

			// tag

			// threat

			// tsdb

			// vod

			// vpc
			"alicloud_vpc_route_target_groups":                 dataSourceAliCloudVpcRouteTargetGroups(),
			"alicloud_vpc_ipv6_cidr_blocks":                    dataSourceAliCloudVpcIpv6CidrBlocks(),
			"alicloud_vpc_flow_logs":                           dataSourceAlicloudVpcFlowLogs(),
			"alicloud_vpc_dhcp_options_sets":                   dataSourceAlicloudVpcDhcpOptionsSets(),
			"alicloud_vpc_traffic_mirror_filters":              dataSourceAlicloudVpcTrafficMirrorFilters(),
			"alicloud_vpc_traffic_mirror_filter_egress_rules":  dataSourceAlicloudVpcTrafficMirrorFilterEgressRules(),
			"alicloud_vpc_traffic_mirror_filter_ingress_rules": dataSourceAlicloudVpcTrafficMirrorFilterIngressRules(),
			"alicloud_vpc_traffic_mirror_service":              dataSourceAlicloudVpcTrafficMirrorService(),
			"alicloud_vpc_traffic_mirror_sessions":             dataSourceAlicloudVpcTrafficMirrorSessions(),
			"alicloud_vpc_ipv6_gateways":                       dataSourceAlicloudVpcIpv6Gateways(),
			"alicloud_vpc_ipv6_egress_rules":                   dataSourceAlicloudVpcIpv6EgressRules(),
			"alicloud_vpc_ipv6_addresses":                      dataSourceAlicloudVpcIpv6Addresses(),
			"alicloud_vpc_ipv6_internet_bandwidths":            dataSourceAlicloudVpcIpv6InternetBandwidths(),
			"alicloud_vpc_ipv4_gateways":                       dataSourceAlicloudVpcIpv4Gateways(),
			"alicloud_vpc_prefix_lists":                        dataSourceAlicloudVpcPrefixLists(),
			"alicloud_vpc_public_ip_address_pools":             dataSourceAlicloudVpcPublicIpAddressPools(),
			"alicloud_vpc_public_ip_address_pool_cidr_blocks":  dataSourceAlicloudVpcPublicIpAddressPoolCidrBlocks(),
			"alicloud_vpc_flow_log_service":                    dataSourceAliCloudVpcFlowLogService(),

			// vpcs
			"alicloud_vpcs": dataSourceAlicloudVpcs(),

			// vpn
			"alicloud_vpn_gateway_enhanced_vpn_gateways": dataSourceAliCloudVpnGatewayEnhancedVpnGateways(),
			"alicloud_vpn_gateway_zones":                 dataSourceAliCloudVPNGatewayZones(),
			"alicloud_vpn_gateways":                      dataSourceAlicloudVpnGateways(),
			"alicloud_vpn_customer_gateways":             dataSourceAlicloudVpnCustomerGateways(),
			"alicloud_vpn_connections":                   dataSourceAlicloudVpnConnections(),
			"alicloud_vpn_ipsec_servers":                 dataSourceAlicloudVpnIpsecServers(),
			"alicloud_vpn_pbr_route_entries":             dataSourceAlicloudVpnPbrRouteEntries(),
			"alicloud_vpn_gateway_vpn_attachments":       dataSourceAliCloudVpnGatewayVpnAttachments(),
			"alicloud_vpn_gateway_vco_routes":            dataSourceAlicloudVpnGatewayVcoRoutes(),

			// vswitches
			"alicloud_vswitches": dataSourceAlicloudVswitches(),

			// waf
			// alicloud_waf_* resources and data sources are registered via
			// service/waf.ServicePackage() in provider/registry/registry.go.
			// alicloud_wafv3_* types are registered via
			// service/wafopenapi.ServicePackage() in provider/registry/registry.go.

			// yundun
			"alicloud_yundun_bastionhost_instances": dataSourceAlicloudBastionhostInstances(),

			// zones
			"alicloud_zones": dataSourceAlicloudZones(),
		},
		ResourcesMap: map[string]*schema.Resource{
			// ack

			// adb resources and data sources are registered via service/adb.ServicePackage()
			// (provider/registry); see alicloud/service/adb/register.go.

			// aligreen

			// alikafka

			// amqp

			// api — registered via alicloud/service/apigateway.ServicePackage()

			// apig

			// auto
			"alicloud_auto_provisioning_group": resourceAlicloudAutoProvisioningGroup(),

			// bastionhost
			"alicloud_bastionhost_instance":                                 resourceAliCloudBastionhostInstance(),
			"alicloud_bastionhost_user_group":                               resourceAlicloudBastionhostUserGroup(),
			"alicloud_bastionhost_user":                                     resourceAlicloudBastionhostUser(),
			"alicloud_bastionhost_host_group":                               resourceAlicloudBastionhostHostGroup(),
			"alicloud_bastionhost_user_attachment":                          resourceAlicloudBastionhostUserAttachment(),
			"alicloud_bastionhost_host":                                     resourceAlicloudBastionhostHost(),
			"alicloud_bastionhost_host_account":                             resourceAlicloudBastionhostHostAccount(),
			"alicloud_bastionhost_host_attachment":                          resourceAlicloudBastionhostHostAttachment(),
			"alicloud_bastionhost_host_account_user_group_attachment":       resourceAlicloudBastionhostHostAccountUserGroupAttachment(),
			"alicloud_bastionhost_host_account_user_attachment":             resourceAlicloudBastionhostHostAccountUserAttachment(),
			"alicloud_bastionhost_host_group_account_user_attachment":       resourceAlicloudBastionhostHostGroupAccountUserAttachment(),
			"alicloud_bastionhost_host_group_account_user_group_attachment": resourceAlicloudBastionhostHostGroupAccountUserGroupAttachment(),
			"alicloud_bastionhost_host_share_key":                           resourceAlicloudBastionhostHostShareKey(),
			"alicloud_bastionhost_host_account_share_key_attachment":        resourceAlicloudBastionhostHostAccountShareKeyAttachment(),

			// brain

			// cassandra

			// chatbot

			// click
			// clickhouse resources migrated to alicloud/service/clickhouse

			// cloud
			"alicloud_cloud_monitor_service_metric_alarm_rule":              resourceAliCloudCloudMonitorServiceMetricAlarmRule(),
			"alicloud_cloud_monitor_service_agent_config":                   resourceAliCloudCloudMonitorServiceAgentConfig(),
			"alicloud_cloud_monitor_service_enterprise_public":              resourceAliCloudCloudMonitorServiceEnterprisePublic(),
			"alicloud_cloud_monitor_service_basic_public":                   resourceAliCloudCloudMonitorServiceBasicPublic(),
			"alicloud_cloud_monitor_service_hybrid_double_write":            resourceAliCloudCloudMonitorServiceHybridDoubleWrite(),
			"alicloud_cloud_monitor_service_monitoring_agent_process":       resourceAliCloudCloudMonitorServiceMonitoringAgentProcess(),
			"alicloud_cloud_monitor_service_group_monitoring_agent_process": resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcess(),

			// cms
			"alicloud_cms_event_notify_policy":     resourceAliCloudCmsEventNotifyPolicy(),
			"alicloud_cms_alert_rule_v2":           resourceAliCloudCmsAlertRuleV2(),
			"alicloud_cms_agg_task_group":          resourceAliCloudCmsAggTaskGroup(),
			"alicloud_cms_addon_release":           resourceAliCloudCmsAddonRelease(),
			"alicloud_cms_prometheus_view":         resourceAliCloudCmsPrometheusView(),
			"alicloud_cms_prometheus_instance":     resourceAliCloudCmsPrometheusInstance(),
			"alicloud_cms_integration_policy":      resourceAliCloudCmsIntegrationPolicy(),
			"alicloud_cms_workspace":               resourceAliCloudCmsWorkspace(),
			"alicloud_cms_dataset":                 resourceAliCloudCmsDataset(),
			"alicloud_cms_alarm":                   resourceAliCloudCmsAlarm(),
			"alicloud_cms_site_monitor":            resourceAliCloudCloudMonitorServiceSiteMonitor(),
			"alicloud_cms_alarm_contact":           resourceAlicloudCmsAlarmContact(),
			"alicloud_cms_alarm_contact_group":     resourceAlicloudCmsAlarmContactGroup(),
			"alicloud_cms_group_metric_rule":       resourceAliCloudCmsGroupMetricRule(),
			"alicloud_cms_monitor_group":           resourceAlicloudCmsMonitorGroup(),
			"alicloud_cms_monitor_group_instances": resourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_cms_metric_rule_template":    resourceAliCloudCmsMetricRuleTemplate(),
			"alicloud_cms_dynamic_tag_group":       resourceAliCloudCmsDynamicTagGroup(),
			"alicloud_cms_namespace":               resourceAliCloudCmsNamespace(),
			"alicloud_cms_sls_group":               resourceAlicloudCmsSlsGroup(),
			"alicloud_cms_hybrid_monitor_sls_task": resourceAlicloudCmsHybridMonitorSlsTask(),
			"alicloud_cms_hybrid_monitor_fc_task":  resourceAlicloudCmsHybridMonitorFcTask(),
			"alicloud_cms_event_rule":              resourceAliCloudCloudMonitorServiceEventRule(),
			"alicloud_cms_metric_rule_black_list":  resourceAlicloudCmsMetricRuleBlackList(),

			// compute

			// config
			// copy
			"alicloud_copy_image": resourceAliCloudImageCopy(),

			// cr
			"alicloud_cr_instance_customized_domain":  resourceAliCloudCrInstanceCustomizedDomain(),
			"alicloud_cr_artifact_subscription_rule":  resourceAliCloudCrArtifactSubscriptionRule(),
			"alicloud_cr_internet_endpoint":           resourceAlicloudCrInternetEndpoint(),
			"alicloud_cr_artifact_lifecycle_rule":     resourceAliCloudCrArtifactLifecycleRule(),
			"alicloud_cr_storage_domain_routing_rule": resourceAliCloudCrStorageDomainRoutingRule(),
			"alicloud_cr_scan_rule":                   resourceAliCloudCrScanRule(),
			"alicloud_cr_namespace":                   resourceAlicloudCRNamespace(),
			"alicloud_cr_repo":                        resourceAlicloudCRRepo(),
			"alicloud_cr_ee_instance":                 resourceAliCloudCrInstance(),
			"alicloud_cr_ee_namespace":                resourceAliCloudCrEENamespace(),
			"alicloud_cr_ee_repo":                     resourceAliCloudCrEERepo(),
			"alicloud_cr_ee_sync_rule":                resourceAliCloudCrRepoSyncRule(),
			"alicloud_cr_endpoint_acl_policy":         resourceAlicloudCrEndpointAclPolicy(),
			"alicloud_cr_chart_namespace":             resourceAlicloudCrChartNamespace(),
			"alicloud_cr_chart_repository":            resourceAlicloudCrChartRepository(),
			"alicloud_cr_chain":                       resourceAlicloudCrChain(),
			"alicloud_cr_vpc_endpoint_linked_vpc":     resourceAlicloudCrVpcEndpointLinkedVpc(),

			// cs

			// das

			// data

			// database

			// db (registered in alicloud/service/rds)

			// ddoscoo

			// direct

			// disk
			"alicloud_disk_attachment": resourceAlicloudEcsDiskAttachment(),

			// dms
			"alicloud_dms_enterprise_workspace":          resourceAliCloudDmsEnterpriseWorkspace(),
			"alicloud_dms_enterprise_authority_template": resourceAliCloudDMSEnterpriseAuthorityTemplate(),
			"alicloud_dms_enterprise_instance":           resourceAlicloudDmsEnterpriseInstance(),
			"alicloud_dms_enterprise_user":               resourceAlicloudDmsEnterpriseUser(),
			"alicloud_dms_enterprise_proxy":              resourceAlicloudDmsEnterpriseProxy(),
			"alicloud_dms_enterprise_proxy_access":       resourceAlicloudDmsEnterpriseProxyAccess(),
			"alicloud_dms_enterprise_logic_database":     resourceAlicloudDmsEnterpriseLogicDatabase(),

			// drds

			// dts
			// (registered via dts.ServicePackage)

			// ecd

			// ecp

			// ecs
			"alicloud_ecs_disk_default_kms_key":            resourceAliCloudEcsDiskDefaultKmsKey(),
			"alicloud_ecs_disk_encryption_by_default":      resourceAliCloudEcsDiskEncryptionByDefault(),
			"alicloud_ecs_ram_role_attachment":             resourceAliCloudEcsRamRoleAttachment(),
			"alicloud_ecs_image_pipeline_execution":        resourceAliCloudEcsImagePipelineExecution(),
			"alicloud_ecs_dedicated_host":                  resourceAlicloudEcsDedicatedHost(),
			"alicloud_ecs_command":                         resourceAlicloudEcsCommand(),
			"alicloud_ecs_hpc_cluster":                     resourceAlicloudEcsHpcCluster(),
			"alicloud_ecs_auto_snapshot_policy":            resourceAliCloudEcsAutoSnapshotPolicy(),
			"alicloud_ecs_launch_template":                 resourceAliCloudEcsLaunchTemplate(),
			"alicloud_ecs_snapshot":                        resourceAliCloudEcsSnapshot(),
			"alicloud_ecs_key_pair":                        resourceAliCloudEcsKeyPair(),
			"alicloud_ecs_key_pair_attachment":             resourceAliCloudEcsKeyPairAttachment(),
			"alicloud_ecs_disk_attachment":                 resourceAlicloudEcsDiskAttachment(),
			"alicloud_ecs_auto_snapshot_policy_attachment": resourceAliCloudEcsAutoSnapshotPolicyAttachment(),
			"alicloud_ecs_network_interface":               resourceAliCloudEcsNetworkInterface(),
			"alicloud_ecs_network_interface_attachment":    resourceAliCloudEcsNetworkInterfaceAttachment(),
			"alicloud_ecs_deployment_set":                  resourceAliCloudEcsDeploymentSet(),
			"alicloud_ecs_dedicated_host_cluster":          resourceAlicloudEcsDedicatedHostCluster(),
			"alicloud_ecs_session_manager_status":          resourceAliCloudEcsSessionManagerStatus(),
			"alicloud_ecs_prefix_list":                     resourceAlicloudEcsPrefixList(),
			"alicloud_ecs_storage_capacity_unit":           resourceAlicloudEcsStorageCapacityUnit(),
			"alicloud_ecs_image_component":                 resourceAliCloudEcsImageComponent(),
			"alicloud_ecs_snapshot_group":                  resourceAlicloudEcsSnapshotGroup(),
			"alicloud_ecs_image_pipeline":                  resourceAlicloudEcsImagePipeline(),
			"alicloud_ecs_network_interface_permission":    resourceAlicloudEcsNetworkInterfacePermission(),
			"alicloud_ecs_invocation":                      resourceAlicloudEcsInvocation(),
			"alicloud_ecs_instance_set":                    resourceAliCloudEcsInstanceSet(),
			"alicloud_ecs_activation":                      resourceAlicloudEcsActivation(),
			"alicloud_ecs_capacity_reservation":            resourceAlicloudEcsCapacityReservation(),
			"alicloud_ecs_elasticity_assurance":            resourceAliCloudEcsElasticityAssurance(),

			// edas

			// ehpc

			// eip

			// elasticsearch

			// ens

			// esa

			// event

			// express

			// fc

			// fcv3 (migrated to alicloud/service/fcv3; registered via registry)

			// forward

			// ga
			// ga types are registered via service/ga.ServicePackage().

			// governance

			// gpdb (registered via alicloud/service/gpdb.ServicePackage)
			// graph (gdb migrated to alicloud/service/gdb)

			// havip
			"alicloud_havip":            resourceAliCloudVpcHaVip(),
			"alicloud_havip_attachment": resourceAliCloudVpcHaVipAttachment(),

			// hbase

			// hbr

			// image
			"alicloud_image":                  resourceAliCloudEcsImage(),
			"alicloud_image_export":           resourceAliCloudImageExport(),
			"alicloud_image_copy":             resourceAliCloudImageCopy(),
			"alicloud_image_import":           resourceAliCloudImageImport(),
			"alicloud_image_share_permission": resourceAliCloudImageSharePermission(),

			// instance
			"alicloud_instance": resourceAliCloudInstance(),

			// key
			"alicloud_key_pair":            resourceAliCloudEcsKeyPair(),
			"alicloud_key_pair_attachment": resourceAliCloudEcsKeyPairAttachment(),

			// kms

			// kvstore

			// launch
			"alicloud_launch_template": resourceAliCloudEcsLaunchTemplate(),

			// live

			// market

			// maxcompute

			// message
			"alicloud_message_service_account_logging": resourceAliCloudMessageServiceAccountLogging(),
			"alicloud_message_service_event_rule":      resourceAliCloudMessageServiceEventRule(),
			"alicloud_message_service_service":         resourceAliCloudMessageServiceService(),
			"alicloud_message_service_endpoint_acl":    resourceAliCloudMessageServiceEndpointAcl(),
			"alicloud_message_service_endpoint":        resourceAliCloudMessageServiceEndpoint(),
			"alicloud_message_service_queue":           resourceAliCloudMessageServiceQueue(),
			"alicloud_message_service_topic":           resourceAliCloudMessageServiceTopic(),
			"alicloud_message_service_subscription":    resourceAliCloudMessageServiceSubscription(),

			// milvus

			// mns

			// mse

			// nas

			// nat

			// network
			"alicloud_network_interface":            resourceAliCloudEcsNetworkInterface(),
			"alicloud_network_interface_attachment": resourceAliCloudEcsNetworkInterfaceAttachment(),
			"alicloud_network_acl":                  resourceAliCloudVpcNetworkAcl(),
			"alicloud_network_acl_attachment":       resourceAliyunNetworkAclAttachment(),
			"alicloud_network_acl_entries":          resourceAliyunNetworkAclEntries(),

			// ocean
			// alicloud_ocean_base_instance is declared by service/oceanbasepro.ServicePackage().

			// oos

			// oss
			// ots

			// pai

			// polar

			// polardb

			// privatelink

			// pvtz

			// quick
			// quickbi types are registered via service/quickbi.ServicePackage().

			// quotas

			// ram
			"alicloud_ram_role_attachment": resourceAlicloudRamRoleAttachment(),

			// rds (registered in alicloud/service/rds)

			// realtime

			// redis

			// reserved
			"alicloud_reserved_instance": resourceAliCloudReservedInstance(),

			// resource
			"alicloud_resource_manager_resource_directory_sharing":     resourceAliCloudResourceManagerResourceDirectorySharing(),
			"alicloud_resource_manager_delivery_channel":               resourceAliCloudResourceManagerDeliveryChannel(),
			"alicloud_resource_manager_multi_account_delivery_channel": resourceAliCloudResourceManagerMultiAccountDeliveryChannel(),
			"alicloud_resource_manager_message_contact":                resourceAliCloudResourceManagerMessageContact(),
			"alicloud_resource_manager_auto_grouping_rule":             resourceAliCloudResourceManagerAutoGroupingRule(),
			"alicloud_resource_manager_saved_query":                    resourceAliCloudResourceManagerSavedQuery(),
			"alicloud_resource_manager_role":                           resourceAlicloudResourceManagerRole(),
			"alicloud_resource_manager_resource_group":                 resourceAliCloudResourceManagerResourceGroup(),
			"alicloud_resource_manager_resource_group_settings":        resourceAliCloudResourceManagerResourceGroupSettings(),
			"alicloud_resource_manager_folder":                         resourceAliCloudResourceManagerFolder(),
			"alicloud_resource_manager_handshake":                      resourceAliCloudResourceManagerHandshake(),
			"alicloud_resource_manager_handshake_acceptance":           resourceAliCloudResourceManagerHandshakeAcceptance(),
			"alicloud_resource_manager_policy":                         resourceAlicloudResourceManagerPolicy(),
			"alicloud_resource_manager_account":                        resourceAliCloudResourceManagerAccount(),
			"alicloud_resource_manager_resource_directory":             resourceAliCloudResourceManagerResourceDirectory(),
			"alicloud_resource_manager_policy_version":                 resourceAlicloudResourceManagerPolicyVersion(),
			"alicloud_resource_manager_policy_attachment":              resourceAlicloudResourceManagerPolicyAttachment(),
			"alicloud_resource_manager_resource_share":                 resourceAliCloudResourceManagerResourceShare(),
			"alicloud_resource_manager_shared_resource":                resourceAliCloudResourceManagerSharedResource(),
			"alicloud_resource_manager_shared_target":                  resourceAliCloudResourceManagerSharedTarget(),
			"alicloud_resource_manager_control_policy":                 resourceAliCloudResourceManagerControlPolicy(),
			"alicloud_resource_manager_control_policy_attachment":      resourceAliCloudResourceManagerControlPolicyAttachment(),
			"alicloud_resource_manager_service_linked_role":            resourceAliCloudResourceManagerServiceLinkedRole(),
			"alicloud_resource_manager_delegated_administrator":        resourceAliCloudResourceManagerDelegatedAdministrator(),

			// ros
			// route
			"alicloud_route_entry":            resourceAliyunRouteEntry(),
			"alicloud_route_table":            resourceAliCloudVpcRouteTable(),
			"alicloud_route_table_attachment": resourceAliCloudVpcRouteTableAttachment(),

			// router

			// sae

			// sag

			// scdn

			// security
			"alicloud_security_group":      resourceAliCloudEcsSecurityGroup(),
			"alicloud_security_group_rule": resourceAliyunSecurityGroupRule(),

			// service

			// slb
			"alicloud_slb":                                resourceAlicloudSlbLoadBalancer(),
			"alicloud_slb_listener":                       resourceAliCloudSlbListener(),
			"alicloud_slb_attachment":                     resourceAliyunSlbAttachment(),
			"alicloud_slb_backend_server":                 resourceAliyunSlbBackendServer(),
			"alicloud_slb_domain_extension":               resourceAlicloudSlbDomainExtension(),
			"alicloud_slb_server_group":                   resourceAliyunSlbServerGroup(),
			"alicloud_slb_master_slave_server_group":      resourceAliyunSlbMasterSlaveServerGroup(),
			"alicloud_slb_rule":                           resourceAliyunSlbRule(),
			"alicloud_slb_acl":                            resourceAlicloudSlbAcl(),
			"alicloud_slb_ca_certificate":                 resourceAlicloudSlbCaCertificate(),
			"alicloud_slb_server_certificate":             resourceAlicloudSlbServerCertificate(),
			"alicloud_slb_load_balancer":                  resourceAlicloudSlbLoadBalancer(),
			"alicloud_slb_tls_cipher_policy":              resourceAlicloudSlbTlsCipherPolicy(),
			"alicloud_slb_acl_entry_attachment":           resourceAlicloudSlbAclEntryAttachment(),
			"alicloud_slb_server_group_server_attachment": resourceAliCloudSlbServerGroupServerAttachment(),

			// smartag

			// sms

			// snapshot
			"alicloud_snapshot":        resourceAliCloudEcsSnapshot(),
			"alicloud_snapshot_policy": resourceAliCloudEcsAutoSnapshotPolicy(),

			// snat

			// ssl
			"alicloud_ssl_vpn_server":      resourceAliyunSslVpnServer(),
			"alicloud_ssl_vpn_client_cert": resourceAliyunSslVpnClientCert(),

			// star: alicloud_star_rocks_* migrated to alicloud/service/starrocks; see starrocks.ServicePackage().

			// subnet
			// "alicloud_subnet" aims to match aws usage habit.
			"alicloud_subnet": resourceAliCloudVpcVswitch(),

			// tag

			// threat

			// tsdb

			// vod

			// vpc
			"alicloud_vpc_route_target_group":                      resourceAliCloudVpcRouteTargetGroup(),
			"alicloud_vpc_ipv6_cidr_block":                         resourceAliCloudVpcIpv6CidrBlock(),
			"alicloud_vpc_ipv6_address":                            resourceAliCloudVpcIpv6Address(),
			"alicloud_vpc_gateway_endpoint_route_table_attachment": resourceAliCloudVpcGatewayEndpointRouteTableAttachment(),
			"alicloud_vpc_gateway_endpoint":                        resourceAliCloudVpcGatewayEndpoint(),
			"alicloud_vpc_vswitch_cidr_reservation":                resourceAliCloudVpcVswitchCidrReservation(),
			"alicloud_vpc_ha_vip":                                  resourceAliCloudVpcHaVip(),
			"alicloud_vpc":                                         resourceAliCloudVpcVpc(),
			"alicloud_vpc_route_entry":                             resourceAliCloudVpcRouteEntry(),
			"alicloud_vpc_flow_log":                                resourceAliCloudVpcFlowLog(),
			"alicloud_vpc_dhcp_options_set":                        resourceAliCloudVpcDhcpOptionsSet(),
			"alicloud_vpc_traffic_mirror_filter":                   resourceAliCloudVpcTrafficMirrorFilter(),
			"alicloud_vpc_traffic_mirror_filter_egress_rule":       resourceAliCloudVpcTrafficMirrorFilterEgressRule(),
			"alicloud_vpc_traffic_mirror_filter_ingress_rule":      resourceAliCloudVpcTrafficMirrorFilterIngressRule(),
			"alicloud_vpc_traffic_mirror_session":                  resourceAliCloudVpcTrafficMirrorSession(),
			"alicloud_vpc_ipv6_gateway":                            resourceAliCloudVpcIpv6Gateway(),
			"alicloud_vpc_ipv6_egress_rule":                        resourceAliCloudVpcIpv6EgressRule(),
			"alicloud_vpc_ipv6_internet_bandwidth":                 resourceAliCloudVpcIpv6InternetBandwidth(),
			"alicloud_vpc_dhcp_options_set_attachment":             resourceAlicloudVpcDhcpOptionsSetAttachement(),
			"alicloud_vpc_ipv4_gateway":                            resourceAliCloudVpcIpv4Gateway(),
			"alicloud_vpc_prefix_list":                             resourceAliCloudVpcPrefixList(),
			"alicloud_vpc_ipv4_cidr_block":                         resourceAliCloudVpcIpv4CidrBlock(),
			"alicloud_vpc_public_ip_address_pool":                  resourceAliCloudVpcPublicIpAddressPool(),
			"alicloud_vpc_public_ip_address_pool_cidr_block":       resourceAliCloudVpcPublicIpAddressPoolCidrBlock(),
			"alicloud_vpc_network_acl_attachment":                  resourceAliCloudVpcNetworkAclAttachment(),
			"alicloud_vpc_gateway_route_table_attachment":          resourceAliCloudVpcGatewayRouteTableAttachment(),

			// vpn
			"alicloud_vpn_gateway_enhanced_vpn_gateway": resourceAliCloudVpnGatewayEnhancedVpnGateway(),
			"alicloud_vpn_gateway":                      resourceAliCloudVPNGatewayVPNGateway(),
			"alicloud_vpn_customer_gateway":             resourceAliCloudVPNGatewayCustomerGateway(),
			"alicloud_vpn_route_entry":                  resourceAliyunVpnRouteEntry(),
			"alicloud_vpn_connection":                   resourceAliCloudVPNGatewayVpnConnection(),
			"alicloud_vpn_ipsec_server":                 resourceAlicloudVpnIpsecServer(),
			"alicloud_vpn_pbr_route_entry":              resourceAlicloudVpnPbrRouteEntry(),
			"alicloud_vpn_gateway_vpn_attachment":       resourceAliCloudVpnGatewayVpnAttachment(),
			"alicloud_vpn_gateway_vco_route":            resourceAlicloudVpnGatewayVcoRoute(),

			// vswitch
			"alicloud_vswitch": resourceAliCloudVpcVswitch(),

			// waf
			// alicloud_waf_* resources and data sources are registered via
			// service/waf.ServicePackage() in provider/registry/registry.go.

			// yundun
			"alicloud_yundun_bastionhost_instance": resourceAliCloudBastionhostInstance(),
		},
	}
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, provider)
	}

	if diags := assembleProvider(provider); diags.HasError() {
		// Registration duplicates are a startup hard error: no configuration
		// can make them work, and the provider is still being built, so there
		// is nowhere to return diagnostics to.
		panic("provider registration error: " + formatRegistrationDiags(diags))
	}

	return provider
}

// hcl2shimUnknownValue mirrors terraform-plugin-sdk's internal
// hcl2shim.UnknownVariableValue. Under SDK v2 the shim hands this sentinel string
// to d.Get/d.GetOk for any provider-config attribute whose value is unknown until
// apply (for example a role_arn that interpolates an attribute of a resource
// created in the same apply). The package is internal and cannot be imported, so
// the constant is mirrored here. Under SDK v1 such attributes read back as empty
// and eager credential resolution was deferred to apply; forwarding the sentinel
// into the AssumeRole call instead fails with "Invalid role arn format:74D93920-...".
// Treating it as unset restores the SDK v1 deferred behavior.
const hcl2shimUnknownValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

// cfgString reads a string provider-config attribute, treating the SDK v2
// plan-unknown sentinel (hcl2shimUnknownValue) as "not set". Route any config
// value that drives an eager configure-time side effect (HTTP/file/credential
// call) through this so an unknown-until-apply value is deferred to apply rather
// than fed as a garbage literal. Returns ("", false) when unset or unknown.
func cfgString(d *schema.ResourceData, key string) (string, bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return "", false
	}
	s, _ := v.(string)
	if s == "" || s == hcl2shimUnknownValue {
		return "", false
	}
	return s, true
}

// mapString is the map-read counterpart of cfgString for values read out of a
// schema block/set element (e.g. an assume_role or assume_role_with_oidc map),
// applying the identical unknown-sentinel handling. Returns ("", false) when the
// key is absent, empty, or the plan-unknown sentinel. Route every credential
// string read from such a map through this so no field is left unguarded.
func mapString(m map[string]interface{}, key string) (string, bool) {
	s, ok := m[key].(string)
	if !ok || s == "" || s == hcl2shimUnknownValue {
		return "", false
	}
	return s, true
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	log.Println("using terraform version:", p.TerraformVersion)
	// profileConfig is resolved per Configure call from this instance's own
	// profile / shared_credentials_file, so multiple provider aliases in one
	// Terraform process never share (or contaminate) each other's profile.
	profileConfig, err := loadProfileConfig(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	var getProviderConfig = func(schemaKey string, profileKey string) (string, error) {
		if schemaKey != "" {
			if v, ok := cfgString(d, schemaKey); ok {
				return v, nil
			}
		}
		v, err := getConfigFromProfile(profileConfig, profileKey)
		if err != nil {
			return "", err
		}
		if v != nil {
			s, ok := v.(string)
			if !ok {
				return "", WrapError(fmt.Errorf("provider attribute %q in profile %q must be a string, got %T", profileKey, d.Get("profile").(string), v))
			}
			return s, nil
		}
		return "", nil
	}

	accessKey, err := getProviderConfig("access_key", "access_key_id")
	if err != nil {
		return nil, diag.FromErr(err)
	}
	secretKey, err := getProviderConfig("secret_key", "access_key_secret")
	if err != nil {
		return nil, diag.FromErr(err)
	}
	region, err := getProviderConfig("region", "region_id")
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if region == "" {
		region = DEFAULT_REGION
	}
	securityToken, err := getProviderConfig("security_token", "sts_token")
	if err != nil {
		return nil, diag.FromErr(err)
	}

	ecsRoleName, err := getProviderConfig("ecs_role_name", "ram_role_name")
	if err != nil {
		return nil, diag.FromErr(err)
	}
	var profileName string
	var credential credentials.Credential
	if v, ok := cfgString(d, "profile"); ok {
		profileName = v
	}

	// TODO: supports all of profile modes after credentials supporting setting timeout
	profileMode := fmt.Sprint(profileConfig["mode"])
	// Support advanced profile modes: CloudSSO, ChainableRamRoleArn, External, OAuth
	if (accessKey == "" || secretKey == "") && profileName != "" &&
		(profileMode == "ChainableRamRoleArn" || profileMode == "CloudSSO" ||
			profileMode == "External" || profileMode == "OAuth") {
		var profileFile string
		if v, ok := cfgString(d, "shared_credentials_file"); ok {
			profileFile = absPath(v)
		}
		provider, err := providers.NewCLIProfileCredentialsProviderBuilder().WithProfileName(profileName).WithProfileFile(profileFile).Build()
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("failed to create profile credentials provider: %v", err))
		}
		credential = credentials.FromCredentialsProvider("cli_profile", provider)
		creds, err := credential.GetCredential()
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("failed to get credential from profile: %v", err))
		}
		accessKey, secretKey, securityToken = *creds.AccessKeyId, *creds.AccessKeySecret, *creds.SecurityToken
	}

	if accessKey == "" || secretKey == "" {
		if v, ok := cfgString(d, "credentials_uri"); ok {
			credentialsURIResp, err := getClientByCredentialsURI(v)
			if err != nil {
				return nil, diag.FromErr(err)
			}
			accessKey = credentialsURIResp.AccessKeyId
			secretKey = credentialsURIResp.AccessKeySecret
			securityToken = credentialsURIResp.SecurityToken
		}
	}

	config := &connectivity.Config{
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		SecurityToken:        strings.TrimSpace(securityToken),
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
	config.Features = expandFeatures(d.Get("features").([]interface{}))
	if credential != nil {
		config.Credential = credential
	}
	log.Println("alicloud provider trace id:", config.TerraformTraceId)
	if accessKey != "" && secretKey != "" && credential == nil {
		credentialConfig := new(credentials.Config).SetType("access_key").
			SetAccessKeyId(accessKey).
			SetAccessKeySecret(secretKey).
			SetTimeout(config.ClientReadTimeout).
			SetConnectTimeout(config.ClientConnectTimeout)
		if v := strings.TrimSpace(securityToken); v != "" {
			credentialConfig.SetType("sts").SetSecurityToken(v)
		}
		if config.ClientConnectTimeout != 0 {
			credentialConfig.SetConnectTimeout(config.ClientConnectTimeout)
		}
		if config.ClientReadTimeout != 0 {
			credentialConfig.SetTimeout(config.ClientReadTimeout)
		}
		credential, err := credentials.NewCredential(credentialConfig)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		config.Credential = credential
	}
	if account, ok := d.GetOk("account_id"); ok && account.(string) != "" {
		config.AccountId = strings.TrimSpace(account.(string))
	}
	if v, ok := d.GetOk("account_type"); ok && v.(string) != "" {
		config.AccountType = v.(string)
	}
	if v, ok := d.GetOk("security_transport"); config.SecureTransport == "" && ok && v.(string) != "" {
		config.SecureTransport = v.(string)
	}

	config.RamRoleArn, err = getProviderConfig("", "ram_role_arn")
	if err != nil {
		return nil, diag.FromErr(err)
	}
	config.RamRoleSessionName, err = getProviderConfig("", "ram_session_name")
	if err != nil {
		return nil, diag.FromErr(err)
	}
	expiredSeconds, err := getConfigFromProfile(profileConfig, "expired_seconds")
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if expiredSeconds != nil {
		expired, ok := expiredSeconds.(float64)
		if !ok {
			return nil, diag.FromErr(WrapError(fmt.Errorf("provider attribute \"expired_seconds\" in profile %q must be a number, got %T", d.Get("profile").(string), expiredSeconds)))
		}
		config.RamRoleSessionExpiration = int(expired)
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		if v, ok := mapString(assumeRole, "role_arn"); ok {
			config.RamRoleArn = v
		}
		if v, ok := mapString(assumeRole, "session_name"); ok {
			config.RamRoleSessionName = v
		}
		if config.RamRoleSessionName == "" {
			config.RamRoleSessionName = "terraform"
		}
		if v, ok := mapString(assumeRole, "policy"); ok {
			config.RamRolePolicy = v
		}
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
		if v, ok := mapString(assumeRole, "external_id"); ok {
			config.RamRoleExternalId = v
		}

		log.Printf("[INFO] assume_role configuration set: (RamRoleArn: %q, RamRoleSessionName: %q, RamRolePolicy: %q, RamRoleSessionExpiration: %d, RamRoleExternalId: %s)",
			config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration, config.RamRoleExternalId)
	}

	if v, ok := d.GetOk("assume_role_with_oidc"); ok && len(v.([]interface{})) == 1 {
		assumeRoleWithOidc, err := getAssumeRoleWithOIDCConfig(v.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		config.AssumeRoleWithOidc = assumeRoleWithOidc
		log.Printf("[INFO] assume_role_with_oidc configuration set: (RoleArn: %q, SessionName: %q, SessionExpiration: %d, OIDCProviderArn: %s)",
			config.AssumeRoleWithOidc.RoleARN, config.AssumeRoleWithOidc.RoleSessionName, config.AssumeRoleWithOidc.DurationSeconds, config.AssumeRoleWithOidc.OIDCProviderArn)
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)
	var endpointInit sync.Map
	config.Endpoints = &endpointInit

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		for key, val := range endpoints {
			// Compatible with the deprecated endpoint setting
			if val == nil || val.(string) == "" {
				if v, ok := deprecatedEndpointMap[key]; ok {
					val = endpoints[v]
				}
			}
			endpointInit.Store(key, val.(string))
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
		return nil, diag.FromErr(err)
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		return nil, diag.FromErr(fmt.Errorf("configuring Terraform Alibaba Cloud Provider: no valid credential sources for Terraform Alibaba Cloud Provider found.\n\n%s",
			"Please see https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#authentication\n"+
				"for more information about providing credentials."))
	}

	client, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}

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

		"skip_region_validation": "Skip static validation of region ID. Used by users of alternative AlibabaCloud-like APIs or users w/ access to regions that are not public (yet). It can also be sourced from the `ALICLOUD_SKIP_REGION_VALIDATION` environment variable.",

		"features": "Customize the behaviour of certain resources. Every toggle it holds is optional, and leaving the block out keeps the provider's default behaviour.",

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

		"eflo_controller_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom efloctrl endpoints.",

		"eflo_cnp": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eflocnp endpoints.",

		"oceanbase_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom oceanbase endpoints.",

		"beebot_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom beebot endpoints.",

		"computenest_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom computenest endpoints.",
	}
}

// lintignore: S018
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
					DefaultFunc:  schema.EnvDefaultFunc("ALIBABA_CLOUD_OIDC_TOKEN", nil),
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

// lintignore: S018
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

// featuresSchema returns the provider's features block. Grouping the behaviour toggles keeps the
// next one from becoming another top level provider argument, at the cost of an environment
// variable: a DefaultFunc on a nested field only runs when its block is present.
func featuresSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ecs_instance": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"replace_on_image_update": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Whether a change to `image_id` on an `alicloud_instance` is planned as a replacement of the instance instead of an in-place replacement of its system disk.",
							},
						},
					},
					Description: "The behaviour toggles of the `alicloud_instance` resource.",
				},
			},
		},
		Description: descriptions["features"],
	}
}

// expandFeatures reads the features block into the struct the client carries, starting from the
// defaults because every level of the block is optional.
func expandFeatures(featuresList []interface{}) features.Features {
	expanded := features.Default()
	if len(featuresList) == 0 || featuresList[0] == nil {
		return expanded
	}
	featuresMap := featuresList[0].(map[string]interface{})
	if ecsInstanceList, ok := featuresMap["ecs_instance"].([]interface{}); ok && len(ecsInstanceList) > 0 && ecsInstanceList[0] != nil {
		ecsInstance := ecsInstanceList[0].(map[string]interface{})
		if v, ok := ecsInstance["replace_on_image_update"].(bool); ok {
			expanded.EcsInstance.ReplaceOnImageUpdate = v
		}
	}
	return expanded
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
				"chatbot": {
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

				"eflo_controller": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eflo_controller_endpoint"],
				},

				"eflo_cnp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eflo_cnp_endpoint"],
				},

				"srvcatalog": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["srvcatalog_endpoint"],
				},
				"servicecatalog": {
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
				"dbs": {
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
				"dysmsapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dysmsapi_endpoint"],
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
				"schedulerx2": {
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

				"antiddos_public": {
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
				"oceanbasepro": {
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
				"eds_user": {
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
				"quickbi_public": {
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
				"gdb": {
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
				"swas_open": {
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
				"dfs": {
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
				"dataworks_public": {
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
				"ecd": {
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
				"sae": {
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

				"amqp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["onsproxy_endpoint"],
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
				"resourcesharing": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
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
				"polardbx": {
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
				"fc_open": {
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
				"cloudapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["apigateway_endpoint"],
				},
				"apig": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"datahub": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["datahub_endpoint"],
				},
				"devops_rdc": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"mns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mns_endpoint"],
				},
				"mns_open": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mns_endpoint"],
				},
				"rocketmq": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
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
				"aiworkspace": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"vpcipam": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"gwlb": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"esa": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
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
	return helper.Hashcode(buf.String())
}

// deprecatedEndpointMap is used to map old service name to new service name,
// and its value equals to the gateway code of the API after converting it to lowercase and using underscores
// key: new endpoint key
// value: deprecated endpoint key
var deprecatedEndpointMap = map[string]string{
	"resourcesharing":  "ressharing",
	"ga":               "gaplus",
	"dms_enterprise":   "dmsenterprise",
	"sgw":              "hcs_sgw",
	"amqp":             "onsproxy",
	"cassandra":        "cds",
	"cloudfw":          "cloudfirewall",
	"sae":              "serverless",
	"r_kvstore":        "redisa",
	"ecd":              "gwsecd",
	"dataworks_public": "dataworkspublic",
	"dfs":              "alidfs",
	"swas_open":        "swas",
	"quickbi_public":   "quickbi",
	"gdb":              "gds",
	"cr":               "acr",
	"eds_user":         "edsuser",
	"antiddos_public":  "ddosbasic",
	"schedulerx2":      "edasschedulerx",
	"ehpc":             "ehs",
	"dysmsapi":         "dysms",
	"dbs":              "cbs",
	"mns_open":         "mns",
	"servicecatalog":   "srvcatalog",
	"oceanbasepro":     "oceanbase",
	"chatbot":          "beebot",
	"cloudapi":         "apigateway",
}

// defaultConfigPath returns the default location of the aliyun CLI shared
// config file under the user's home directory. It is used by
// getConfigFromProfile when the caller has not configured
// shared_credentials_file.
//
// The home directory is resolved through homedir.Dir(), which falls back to
// the OS user database when HOME / USERPROFILE are unset. This avoids the
// previous behaviour of building the path with
// fmt.Sprintf("%s/.aliyun/config.json", os.Getenv("USERPROFILE")) which, when
// USERPROFILE is empty (Windows service accounts, minimal containers), yields
// "/.aliyun/config.json" that os.Stat resolves to the current drive root and
// silently leaves the profile unloaded, surfacing as an opaque authentication
// failure. Using filepath.Join also ensures OS-correct path separators.
func defaultConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", WrapError(err)
	}
	return filepath.Join(home, ".aliyun", "config.json"), nil
}

// loadProfileConfig resolves the aliyun CLI profile selected by this provider
// instance's own "profile" / "shared_credentials_file" attributes. It reads
// the profile file on every call and returns a fresh map, so the result is
// scoped to the calling provider instance and can never leak into (or be
// replaced by) another provider alias's configuration. Returns (nil, nil)
// when no profile is configured or the file does not exist; a malformed or
// unreadable file surfaces as an error carrying the file path, profile name
// and reason instead of a silent empty config.
func loadProfileConfig(d *schema.ResourceData) (map[string]interface{}, error) {
	if v, ok := d.GetOk("profile"); !ok || v.(string) == "" {
		return nil, nil
	}
	current := d.Get("profile").(string)
	// Set Credentials filename, expanding home directory
	profilePath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
	if err != nil {
		return nil, WrapError(fmt.Errorf("failed to expand shared_credentials_file for profile %q: %v", current, err))
	}
	if profilePath == "" {
		configPath, err := defaultConfigPath()
		if err != nil {
			return nil, WrapError(fmt.Errorf("failed to resolve default credentials file for profile %q: %v", current, err))
		}
		profilePath = configPath
	}
	profileConfig := make(map[string]interface{})
	_, err = os.Stat(profilePath)
	if !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(profilePath)
		if err != nil {
			return nil, WrapError(fmt.Errorf("failed to read credentials file %s for profile %q: %v", profilePath, current, err))
		}
		config := map[string]interface{}{}
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, WrapError(fmt.Errorf("failed to parse credentials file %s for profile %q: invalid JSON: %v", profilePath, current, err))
		}
		profiles, ok := config["profiles"].([]interface{})
		if !ok {
			return nil, WrapError(fmt.Errorf("invalid credentials file %s for profile %q: field \"profiles\" must be an array of profile objects", profilePath, current))
		}
		for _, entry := range profiles {
			profileMap, ok := entry.(map[string]interface{})
			if !ok {
				return nil, WrapError(fmt.Errorf("invalid credentials file %s for profile %q: profile entries must be JSON objects", profilePath, current))
			}
			if current == profileMap["name"] {
				profileConfig = profileMap
			}
		}
		if modeValue, ok := profileConfig["mode"]; ok {
			if _, ok := modeValue.(string); !ok {
				return nil, WrapError(fmt.Errorf("invalid credentials file %s for profile %q: field \"mode\" must be a string, got %T", profilePath, current, modeValue))
			}
		}
	}
	return profileConfig, nil
}

// getConfigFromProfile extracts ProfileKey from an already-loaded profile
// config map (see loadProfileConfig), applying the per-mode gating rules.
// A nil map or an absent "mode" yields (nil, nil); loadProfileConfig already
// guards the mode type at load time, so the defensive check here only fires
// for a hand-assembled config map.
func getConfigFromProfile(profileConfig map[string]interface{}, ProfileKey string) (interface{}, error) {
	modeValue, ok := profileConfig["mode"]
	if !ok {
		return nil, nil
	}
	mode, ok := modeValue.(string)
	if !ok {
		return nil, WrapError(fmt.Errorf("invalid profile entry: field \"mode\" must be a string, got %T", modeValue))
	}
	if ProfileKey == "region_id" {
		return profileConfig["region_id"], nil
	}
	// For advanced profile modes, credentials are obtained dynamically via CLIProfileCredentialsProvider
	if mode == "ChainableRamRoleArn" || mode == "CloudSSO" ||
		mode == "External" || mode == "OAuth" {
		return nil, nil
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

	return profileConfig[ProfileKey], nil
}

func getAssumeRoleWithOIDCConfig(tfMap map[string]interface{}) (*connectivity.AssumeRoleWithOidc, error) {
	if tfMap == nil {
		return nil, nil
	}

	assumeRole := connectivity.AssumeRoleWithOidc{}

	if v, ok := tfMap["session_expiration"].(int); ok && v != 0 {
		assumeRole.DurationSeconds = v
	}

	if v, ok := mapString(tfMap, "policy"); ok {
		assumeRole.Policy = v
	}

	if v, ok := mapString(tfMap, "role_arn"); ok {
		assumeRole.RoleARN = v
	}

	if v, ok := mapString(tfMap, "role_session_name"); ok {
		assumeRole.RoleSessionName = v
	}
	if assumeRole.RoleSessionName == "" {
		assumeRole.RoleSessionName = "terraform"
	}

	if v, ok := mapString(tfMap, "oidc_provider_arn"); ok {
		assumeRole.OIDCProviderArn = v
	}

	missingOidcToken := true
	if v, ok := mapString(tfMap, "oidc_token"); ok {
		assumeRole.OIDCToken = v
		missingOidcToken = false
	}

	if v, ok := mapString(tfMap, "oidc_token_file"); ok {
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

func absPath(filePath string) string {
	if v, err := homedir.Expand(filePath); err != nil {
		log.Printf("[WARN] failed to expand profile file path: %v", err)
	} else {
		filePath = v
	}

	if v, err := filepath.Abs(filePath); err != nil {
		log.Printf("[WARN] failed to get absolute path of profile file: %v", err)
	} else {
		filePath = v
	}

	return filePath
}
