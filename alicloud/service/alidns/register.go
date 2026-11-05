package alidns

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the alidns service's SDK v2 resources and data sources.
// Handwritten SDKv2 code migrated from the root alicloud package; original
// behaviour, schema and IDs are preserved. dns-prefixed entries are legacy
// aliases that share the alidns implementation factories.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Alidns",
		SDKResources: []conns.SDKResource{
			// alidns
			{TypeName: "alicloud_alidns_cloud_gtm_instance_config", Factory: resourceAliCloudAlidnsCloudGtmInstanceConfig},
			{TypeName: "alicloud_alidns_cloud_gtm_monitor_template", Factory: resourceAliCloudAlidnsCloudGtmMonitorTemplate},
			{TypeName: "alicloud_alidns_cloud_gtm_address", Factory: resourceAliCloudAlidnsCloudGtmAddress},
			{TypeName: "alicloud_alidns_cloud_gtm_address_pool", Factory: resourceAliCloudAlidnsCloudGtmAddressPool},
			{TypeName: "alicloud_alidns_domain_attachment", Factory: resourceAlicloudAlidnsDomainAttachment},
			{TypeName: "alicloud_alidns_domain_group", Factory: resourceAlicloudAlidnsDomainGroup},
			{TypeName: "alicloud_alidns_record", Factory: resourceAliCloudAlidnsRecord},
			{TypeName: "alicloud_alidns_domain", Factory: resourceAlicloudAlidnsDomain},
			{TypeName: "alicloud_alidns_instance", Factory: resourceAlicloudAlidnsInstance},
			{TypeName: "alicloud_alidns_custom_line", Factory: resourceAlicloudAlidnsCustomLine},
			{TypeName: "alicloud_alidns_gtm_instance", Factory: resourceAlicloudAlidnsGtmInstance},
			{TypeName: "alicloud_alidns_address_pool", Factory: resourceAlicloudAlidnsAddressPool},
			{TypeName: "alicloud_alidns_access_strategy", Factory: resourceAlicloudAlidnsAccessStrategy},
			{TypeName: "alicloud_alidns_monitor_config", Factory: resourceAlicloudAlidnsMonitorConfig},
			// dns (alicloud_dns_* legacy aliases share the alidns factories)
			{TypeName: "alicloud_dns_record", Factory: resourceAlicloudDnsRecord},
			{TypeName: "alicloud_dns", Factory: resourceAlicloudDns},
			{TypeName: "alicloud_dns_group", Factory: resourceAlicloudDnsGroup},
			{TypeName: "alicloud_dns_instance", Factory: resourceAlicloudAlidnsInstance},
			{TypeName: "alicloud_dns_domain_attachment", Factory: resourceAlicloudAlidnsDomainAttachment},
			{TypeName: "alicloud_dns_domain", Factory: resourceAlicloudAlidnsDomain},
		},
		SDKDataSources: []conns.SDKDataSource{
			// alidns
			{TypeName: "alicloud_alidns_domain_groups", Factory: dataSourceAlicloudAlidnsDomainGroups},
			{TypeName: "alicloud_alidns_records", Factory: dataSourceAlicloudAlidnsRecords},
			{TypeName: "alicloud_alidns_instances", Factory: dataSourceAlicloudAlidnsInstances},
			{TypeName: "alicloud_alidns_domains", Factory: dataSourceAlicloudAlidnsDomains},
			{TypeName: "alicloud_alidns_custom_lines", Factory: dataSourceAlicloudAlidnsCustomLines},
			{TypeName: "alicloud_alidns_gtm_instances", Factory: dataSourceAlicloudAlidnsGtmInstances},
			{TypeName: "alicloud_alidns_address_pools", Factory: dataSourceAlicloudAlidnsAddressPools},
			{TypeName: "alicloud_alidns_access_strategies", Factory: dataSourceAlicloudAlidnsAccessStrategies},
			// dns (alicloud_dns_* legacy aliases share the matching factories)
			{TypeName: "alicloud_dns_resolution_lines", Factory: dataSourceAlicloudDnsResolutionLines},
			{TypeName: "alicloud_dns_domains", Factory: dataSourceAlicloudAlidnsDomains},
			{TypeName: "alicloud_dns_groups", Factory: dataSourceAlicloudDnsGroups},
			{TypeName: "alicloud_dns_records", Factory: dataSourceAlicloudDnsRecords},
			{TypeName: "alicloud_dns_domain_groups", Factory: dataSourceAlicloudDnsGroups},
			{TypeName: "alicloud_dns_domain_records", Factory: dataSourceAlicloudDnsRecords},
			{TypeName: "alicloud_dns_domain_txt_guid", Factory: dataSourceAlicloudDnsDomainTxtGuid},
			{TypeName: "alicloud_dns_instances", Factory: dataSourceAlicloudAlidnsInstances},
		},
	}
}
