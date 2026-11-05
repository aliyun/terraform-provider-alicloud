package ga

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Ga (Global
// Accelerator) service. The 24 resource types and 20 data source types mirror
// the registrations previously declared inline in the root provider package;
// their implementations and factories now live in this product subpackage.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ga",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ga_accelerator",
				Factory:  ResourceAliCloudGaAccelerator,
			},
			{
				TypeName: "alicloud_ga_listener",
				Factory:  ResourceAliCloudGaListener,
			},
			{
				TypeName: "alicloud_ga_bandwidth_package",
				Factory:  ResourceAliCloudGaBandwidthPackage,
			},
			{
				TypeName: "alicloud_ga_endpoint_group",
				Factory:  ResourceAliCloudGaEndpointGroup,
			},
			{
				TypeName: "alicloud_ga_bandwidth_package_attachment",
				Factory:  ResourceAliCloudGaBandwidthPackageAttachment,
			},
			{
				TypeName: "alicloud_ga_ip_set",
				Factory:  ResourceAliCloudGaIpSet,
			},
			{
				TypeName: "alicloud_ga_forwarding_rule",
				Factory:  ResourceAliCloudGaForwardingRule,
			},
			{
				TypeName: "alicloud_ga_acl",
				Factory:  ResourceAliCloudGaAcl,
			},
			{
				TypeName: "alicloud_ga_acl_attachment",
				Factory:  ResourceAliCloudGaAclAttachment,
			},
			{
				TypeName: "alicloud_ga_additional_certificate",
				Factory:  ResourceAliCloudGaAdditionalCertificate,
			},
			{
				TypeName: "alicloud_ga_accelerator_spare_ip_attachment",
				Factory:  ResourceAlicloudGaAcceleratorSpareIpAttachment,
			},
			{
				TypeName: "alicloud_ga_access_log",
				Factory:  ResourceAlicloudGaAccessLog,
			},
			{
				TypeName: "alicloud_ga_acl_entry_attachment",
				Factory:  ResourceAliCloudGaAclEntryAttachment,
			},
			{
				TypeName: "alicloud_ga_basic_accelerator",
				Factory:  ResourceAliCloudGaBasicAccelerator,
			},
			{
				TypeName: "alicloud_ga_basic_endpoint_group",
				Factory:  ResourceAliCloudGaBasicEndpointGroup,
			},
			{
				TypeName: "alicloud_ga_basic_ip_set",
				Factory:  ResourceAliCloudGaBasicIpSet,
			},
			{
				TypeName: "alicloud_ga_basic_accelerate_ip",
				Factory:  ResourceAliCloudGaBasicAccelerateIp,
			},
			{
				TypeName: "alicloud_ga_basic_endpoint",
				Factory:  ResourceAliCloudGaBasicEndpoint,
			},
			{
				TypeName: "alicloud_ga_basic_accelerate_ip_endpoint_relation",
				Factory:  ResourceAlicloudGaBasicAccelerateIpEndpointRelation,
			},
			{
				TypeName: "alicloud_ga_domain",
				Factory:  ResourceAlicloudGaDomain,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoint_group",
				Factory:  ResourceAliCloudGaCustomRoutingEndpointGroup,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoint_group_destination",
				Factory:  ResourceAliCloudGaCustomRoutingEndpointGroupDestination,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoint",
				Factory:  ResourceAliCloudGaCustomRoutingEndpoint,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoint_traffic_policy",
				Factory:  ResourceAliCloudGaCustomRoutingEndpointTrafficPolicy,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ga_accelerators",
				Factory:  DataSourceAlicloudGaAccelerators,
			},
			{
				TypeName: "alicloud_ga_listeners",
				Factory:  DataSourceAlicloudGaListeners,
			},
			{
				TypeName: "alicloud_ga_bandwidth_packages",
				Factory:  DataSourceAlicloudGaBandwidthPackages,
			},
			{
				TypeName: "alicloud_ga_endpoint_groups",
				Factory:  DataSourceAliCloudGaEndpointGroups,
			},
			{
				TypeName: "alicloud_ga_ip_sets",
				Factory:  DataSourceAlicloudGaIpSets,
			},
			{
				TypeName: "alicloud_ga_forwarding_rules",
				Factory:  DataSourceAlicloudGaForwardingRules,
			},
			{
				TypeName: "alicloud_ga_acls",
				Factory:  DataSourceAlicloudGaAcls,
			},
			{
				TypeName: "alicloud_ga_additional_certificates",
				Factory:  DataSourceAlicloudGaAdditionalCertificates,
			},
			{
				TypeName: "alicloud_ga_accelerator_spare_ip_attachments",
				Factory:  DataSourceAlicloudGaAcceleratorSpareIpAttachments,
			},
			{
				TypeName: "alicloud_ga_basic_accelerators",
				Factory:  DataSourceAliCloudGaBasicAccelerators,
			},
			{
				TypeName: "alicloud_ga_basic_accelerate_ips",
				Factory:  DataSourceAlicloudGaBasicAccelerateIps,
			},
			{
				TypeName: "alicloud_ga_basic_endpoints",
				Factory:  DataSourceAlicloudGaBasicEndpoints,
			},
			{
				TypeName: "alicloud_ga_basic_accelerate_ip_endpoint_relations",
				Factory:  DataSourceAlicloudGaBasicAccelerateIpEndpointRelations,
			},
			{
				TypeName: "alicloud_ga_domains",
				Factory:  DataSourceAlicloudGaDomains,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoint_groups",
				Factory:  DataSourceAlicloudGaCustomRoutingEndpointGroups,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoint_group_destinations",
				Factory:  DataSourceAlicloudGaCustomRoutingEndpointGroupDestinations,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoints",
				Factory:  DataSourceAlicloudGaCustomRoutingEndpoints,
			},
			{
				TypeName: "alicloud_ga_custom_routing_endpoint_traffic_policies",
				Factory:  DataSourceAliCloudGaCustomRoutingEndpointTrafficPolicies,
			},
			{
				TypeName: "alicloud_ga_custom_routing_port_mappings",
				Factory:  DataSourceAlicloudGaCustomRoutingPortMappings,
			},
			{
				TypeName: "alicloud_ga_endpoint_group_ip_address_cidr_blocks",
				Factory:  DataSourceAliCloudGaEndpointGroupIpAddressCidrBlocks,
			},
		},
	}
}
