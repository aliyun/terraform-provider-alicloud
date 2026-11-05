package alb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ALB service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Alb",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_alb_acl",
				Factory:  ResourceAcl,
			},
			{
				TypeName: "alicloud_alb_acl_entry_attachment",
				Factory:  ResourceAclEntryAttachment,
			},
			{
				TypeName: "alicloud_alb_listener_acl_attachment",
				Factory:  ResourceListenerAclAttachment,
			},
			{
				TypeName: "alicloud_alb_ascript",
				Factory:  ResourceAscript,
			},
			{
				TypeName: "alicloud_alb_health_check_template",
				Factory:  ResourceHealthCheckTemplate,
			},
			{
				TypeName: "alicloud_alb_listener",
				Factory:  ResourceListener,
			},
			{
				TypeName: "alicloud_alb_listener_additional_certificate_attachment",
				Factory:  ResourceListenerAdditionalCertificateAttachment,
			},
			{
				TypeName: "alicloud_alb_load_balancer",
				Factory:  ResourceLoadBalancer,
			},
			{
				TypeName: "alicloud_alb_load_balancer_access_log_config_attachment",
				Factory:  ResourceLoadBalancerAccessLogConfigAttachment,
			},
			{
				TypeName: "alicloud_alb_load_balancer_common_bandwidth_package_attachment",
				Factory:  ResourceLoadBalancerCommonBandwidthPackageAttachment,
			},
			{
				TypeName: "alicloud_alb_load_balancer_security_group_attachment",
				Factory:  ResourceLoadBalancerSecurityGroupAttachment,
			},
			{
				TypeName: "alicloud_alb_load_balancer_zone_shifted_attachment",
				Factory:  ResourceLoadBalancerZoneShiftedAttachment,
			},
			{
				TypeName: "alicloud_alb_rule",
				Factory:  ResourceRule,
			},
			{
				TypeName: "alicloud_alb_security_policy",
				Factory:  ResourceSecurityPolicy,
			},
			{
				TypeName: "alicloud_alb_server_group",
				Factory:  ResourceServerGroup,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_alb_acls",
				Factory:  DataSourceAcls,
			},
			{
				TypeName: "alicloud_alb_ascripts",
				Factory:  DataSourceAscripts,
			},
			{
				TypeName: "alicloud_alb_health_check_templates",
				Factory:  DataSourceHealthCheckTemplates,
			},
			{
				TypeName: "alicloud_alb_listeners",
				Factory:  DataSourceListeners,
			},
			{
				TypeName: "alicloud_alb_load_balancers",
				Factory:  DataSourceLoadBalancers,
			},
			{
				TypeName: "alicloud_alb_rules",
				Factory:  DataSourceRules,
			},
			{
				TypeName: "alicloud_alb_security_policies",
				Factory:  DataSourceSecurityPolicies,
			},
			{
				TypeName: "alicloud_alb_server_groups",
				Factory:  DataSourceServerGroups,
			},
			{
				TypeName: "alicloud_alb_system_security_policies",
				Factory:  DataSourceSystemSecurityPolicies,
			},
			{
				TypeName: "alicloud_alb_zones",
				Factory:  DataSourceZones,
			},
		},
	}
}
