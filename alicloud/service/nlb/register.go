package nlb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Nlb service.
// The factories live in this package and are wired into the provider through
// the registry; the SDKResources/SDKDataSources below replace the 16 lines
// previously maintained inline in alicloud/provider.go.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Nlb",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_nlb_load_balancer",
				Factory:  ResourceLoadBalancer,
			},
			{
				TypeName: "alicloud_nlb_listener",
				Factory:  ResourceListener,
			},
			{
				TypeName: "alicloud_nlb_server_group",
				Factory:  ResourceServerGroup,
			},
			{
				TypeName: "alicloud_nlb_server_group_server_attachment",
				Factory:  ResourceServerGroupServerAttachment,
			},
			{
				TypeName: "alicloud_nlb_security_policy",
				Factory:  ResourceSecurityPolicy,
			},
			{
				TypeName: "alicloud_nlb_load_balancer_security_group_attachment",
				Factory:  ResourceLoadBalancerSecurityGroupAttachment,
			},
			{
				TypeName: "alicloud_nlb_load_balancer_zone_shifted_attachment",
				Factory:  ResourceLoadBalancerZoneShiftedAttachment,
			},
			{
				TypeName: "alicloud_nlb_loadbalancer_common_bandwidth_package_attachment",
				Factory:  ResourceLoadbalancerCommonBandwidthPackageAttachment,
			},
			{
				TypeName: "alicloud_nlb_listener_additional_certificate_attachment",
				Factory:  ResourceListenerAdditionalCertificateAttachment,
			},
			{
				TypeName: "alicloud_nlb_hd_monitor_region_config",
				Factory:  ResourceHdMonitorRegionConfig,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_nlb_zones",
				Factory:  DataSourceZones,
			},
			{
				TypeName: "alicloud_nlb_load_balancers",
				Factory:  DataSourceLoadBalancers,
			},
			{
				TypeName: "alicloud_nlb_listeners",
				Factory:  DataSourceListeners,
			},
			{
				TypeName: "alicloud_nlb_server_groups",
				Factory:  DataSourceServerGroups,
			},
			{
				TypeName: "alicloud_nlb_server_group_server_attachments",
				Factory:  DataSourceServerGroupServerAttachments,
			},
			{
				TypeName: "alicloud_nlb_security_policies",
				Factory:  DataSourceSecurityPolicies,
			},
		},
	}
}
