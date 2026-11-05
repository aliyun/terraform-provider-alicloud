package cdn

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CDN service.
//
// All CDN SDKv2 resources and data sources migrated from the root flat layout
// are declared here; root provider.go no longer lists them inline. The
// factories are the same SDKv2 functions that the migrated type files define
// in this package, preserving their original (case-sensitive) names.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cdn",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cdn_domain_new",
				Factory:  resourceAliCloudCdnDomain,
			},
			{
				TypeName: "alicloud_cdn_domain_config",
				Factory:  resourceAliCloudCdnDomainConfig,
			},
			{
				TypeName: "alicloud_cdn_real_time_log_delivery",
				Factory:  resourceAliCloudCdnRealTimeLogDelivery,
			},
			{
				TypeName: "alicloud_cdn_fc_trigger",
				Factory:  resourceAlicloudCdnFcTrigger,
			},
			{
				// Deprecated alias kept on its own legacy factory and schema.
				TypeName: "alicloud_cdn_domain",
				Factory:  resourceAlicloudCdnDomain,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cdn_domain_configs",
				Factory:  dataSourceAliCloudCdnDomainConfigs,
			},
			{
				TypeName: "alicloud_cdn_service",
				Factory:  dataSourceAliCloudCdnService,
			},
			{
				TypeName: "alicloud_cdn_real_time_log_deliveries",
				Factory:  dataSourceAlicloudCdnRealTimeLogDeliveries,
			},
			{
				TypeName: "alicloud_cdn_ip_info",
				Factory:  dataSourceAlicloudCdnIpInfo,
			},
			{
				TypeName: "alicloud_cdn_blocked_regions",
				Factory:  dataSourceAlicloudCdnBlockedRegions,
			},
		},
	}
}
