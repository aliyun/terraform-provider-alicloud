package cen

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cen service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cen",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cen_transit_router_ecr_attachment",
				Factory:  resourceAliCloudCenTransitRouterEcrAttachment,
			},
			{
				TypeName: "alicloud_cen_instance",
				Factory:  resourceAliCloudCenCenInstance,
			},
			{
				TypeName: "alicloud_cen_instance_attachment",
				Factory:  resourceAlicloudCenInstanceAttachment,
			},
			{
				TypeName: "alicloud_cen_bandwidth_package",
				Factory:  resourceAlicloudCenBandwidthPackage,
			},
			{
				TypeName: "alicloud_cen_bandwidth_package_attachment",
				Factory:  resourceAlicloudCenBandwidthPackageAttachment,
			},
			{
				TypeName: "alicloud_cen_bandwidth_limit",
				Factory:  resourceAlicloudCenBandwidthLimit,
			},
			{
				TypeName: "alicloud_cen_route_entry",
				Factory:  resourceAlicloudCenRouteEntry,
			},
			{
				TypeName: "alicloud_cen_instance_grant",
				Factory:  resourceAlicloudCenInstanceGrant,
			},
			{
				TypeName: "alicloud_cen_transit_router",
				Factory:  resourceAliCloudCenTransitRouter,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_entry",
				Factory:  resourceAlicloudCenTransitRouterRouteEntry,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_table",
				Factory:  resourceAlicloudCenTransitRouterRouteTable,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_table_association",
				Factory:  resourceAliCloudCenTransitRouterRouteTableAssociation,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_table_propagation",
				Factory:  resourceAlicloudCenTransitRouterRouteTablePropagation,
			},
			{
				TypeName: "alicloud_cen_transit_router_vbr_attachment",
				Factory:  resourceAliCloudCenTransitRouterVbrAttachment,
			},
			{
				TypeName: "alicloud_cen_transit_router_vpc_attachment",
				Factory:  resourceAliCloudCenTransitRouterVpcAttachment,
			},
			{
				TypeName: "alicloud_cen_flowlog",
				Factory:  resourceAliCloudCenFlowLog,
			},
			{
				TypeName: "alicloud_cen_route_map",
				Factory:  resourceAlicloudCenRouteMap,
			},
			{
				TypeName: "alicloud_cen_private_zone",
				Factory:  resourceAliCloudCenPrivateZone,
			},
			{
				TypeName: "alicloud_cen_vbr_health_check",
				Factory:  resourceAlicloudCenVbrHealthCheck,
			},
			{
				TypeName: "alicloud_cen_route_service",
				Factory:  resourceAlicloudCenRouteService,
			},
			{
				TypeName: "alicloud_cen_transit_router_peer_attachment",
				Factory:  resourceAliCloudCenTransitRouterPeerAttachment,
			},
			{
				TypeName: "alicloud_cen_traffic_marking_policy",
				Factory:  resourceAliCloudCenTrafficMarkingPolicy,
			},
			{
				TypeName: "alicloud_cen_transit_router_vpn_attachment",
				Factory:  resourceAliCloudCenTransitRouterVpnAttachment,
			},
			{
				TypeName: "alicloud_cen_transit_router_grant_attachment",
				Factory:  resourceAliCloudCenTransitRouterGrantAttachment,
			},
			{
				TypeName: "alicloud_cen_transit_router_prefix_list_association",
				Factory:  resourceAlicloudCenTransitRouterPrefixListAssociation,
			},
			{
				TypeName: "alicloud_cen_transit_router_cidr",
				Factory:  resourceAliCloudCenTransitRouterCidr,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain",
				Factory:  resourceAliCloudCenTransitRouterMulticastDomain,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_source",
				Factory:  resourceAlicloudCenTransitRouterMulticastDomainSource,
			},
			{
				TypeName: "alicloud_cen_inter_region_traffic_qos_policy",
				Factory:  resourceAliCloudCenInterRegionTrafficQosPolicy,
			},
			{
				TypeName: "alicloud_cen_inter_region_traffic_qos_queue",
				Factory:  resourceAliCloudCenInterRegionTrafficQosQueue,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_peer_member",
				Factory:  resourceAliCloudCenTransitRouterMulticastDomainPeerMember,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_member",
				Factory:  resourceAliCloudCenTransitRouterMulticastDomainMember,
			},
			{
				TypeName: "alicloud_cen_child_instance_route_entry_to_attachment",
				Factory:  resourceAlicloudCenChildInstanceRouteEntryToAttachment,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_association",
				Factory:  resourceAliCloudCenTransitRouterMulticastDomainAssociation,
			},
			{
				TypeName: "alicloud_cen_transit_route_table_aggregation",
				Factory:  resourceAliCloudCenTransitRouteTableAggregation,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cen_instances",
				Factory:  dataSourceAlicloudCenInstances,
			},
			{
				TypeName: "alicloud_cen_bandwidth_packages",
				Factory:  dataSourceAlicloudCenBandwidthPackages,
			},
			{
				TypeName: "alicloud_cen_bandwidth_limits",
				Factory:  dataSourceAlicloudCenBandwidthLimits,
			},
			{
				TypeName: "alicloud_cen_route_entries",
				Factory:  dataSourceAlicloudCenRouteEntries,
			},
			{
				TypeName: "alicloud_cen_region_route_entries",
				Factory:  dataSourceAlicloudCenRegionRouteEntries,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_entries",
				Factory:  dataSourceAlicloudCenTransitRouterRouteEntries,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_table_associations",
				Factory:  dataSourceAlicloudCenTransitRouterRouteTableAssociations,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_table_propagations",
				Factory:  dataSourceAliCloudCenTransitRouterRouteTablePropagations,
			},
			{
				TypeName: "alicloud_cen_transit_router_route_tables",
				Factory:  dataSourceAliCloudCenTransitRouterRouteTables,
			},
			{
				TypeName: "alicloud_cen_transit_router_vbr_attachments",
				Factory:  dataSourceAlicloudCenTransitRouterVbrAttachments,
			},
			{
				TypeName: "alicloud_cen_transit_router_vpc_attachments",
				Factory:  dataSourceAliCloudCenTransitRouterVpcAttachments,
			},
			{
				TypeName: "alicloud_cen_transit_routers",
				Factory:  dataSourceAlicloudCenTransitRouters,
			},
			{
				TypeName: "alicloud_cen_flowlogs",
				Factory:  dataSourceAliCloudCenFlowLogs,
			},
			{
				TypeName: "alicloud_cen_route_maps",
				Factory:  dataSourceAlicloudCenRouteMaps,
			},
			{
				TypeName: "alicloud_cen_private_zones",
				Factory:  dataSourceAlicloudCenPrivateZones,
			},
			{
				TypeName: "alicloud_cen_instance_attachments",
				Factory:  dataSourceAlicloudCenInstanceAttachments,
			},
			{
				TypeName: "alicloud_cen_vbr_health_checks",
				Factory:  dataSourceAlicloudCenVbrHealthChecks,
			},
			{
				TypeName: "alicloud_cen_route_services",
				Factory:  dataSourceAlicloudCenRouteServices,
			},
			{
				TypeName: "alicloud_cen_transit_router_peer_attachments",
				Factory:  dataSourceAlicloudCenTransitRouterPeerAttachments,
			},
			{
				TypeName: "alicloud_cen_transit_router_service",
				Factory:  dataSourceAlicloudCenTransitRouterService,
			},
			{
				TypeName: "alicloud_cen_transit_router_available_resources",
				Factory:  dataSourceAliCloudCenTransitRouterAvailableResources,
			},
			{
				TypeName: "alicloud_cen_traffic_marking_policies",
				Factory:  dataSourceAlicloudCenTrafficMarkingPolicies,
			},
			{
				TypeName: "alicloud_cen_transit_router_vpn_attachments",
				Factory:  dataSourceAliCloudCenTransitRouterVpnAttachments,
			},
			{
				TypeName: "alicloud_cen_transit_router_prefix_list_associations",
				Factory:  dataSourceAlicloudCenTransitRouterPrefixListAssociations,
			},
			{
				TypeName: "alicloud_cen_transit_router_cidrs",
				Factory:  dataSourceAlicloudCenTransitRouterCidrs,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domains",
				Factory:  dataSourceAlicloudCenTransitRouterMulticastDomains,
			},
			{
				TypeName: "alicloud_cen_inter_region_traffic_qos_policies",
				Factory:  dataSourceAlicloudCenInterRegionTrafficQosPolicies,
			},
			{
				TypeName: "alicloud_cen_inter_region_traffic_qos_queues",
				Factory:  dataSourceAlicloudCenInterRegionTrafficQosQueues,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_peer_members",
				Factory:  dataSourceAlicloudCenTransitRouterMulticastDomainPeerMembers,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_members",
				Factory:  dataSourceAlicloudCenTransitRouterMulticastDomainMembers,
			},
			{
				TypeName: "alicloud_cen_child_instance_route_entry_to_attachments",
				Factory:  dataSourceAlicloudCenChildInstanceRouteEntryToAttachments,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_associations",
				Factory:  dataSourceAlicloudCenTransitRouterMulticastDomainAssociations,
			},
			{
				TypeName: "alicloud_cen_transit_router_multicast_domain_sources",
				Factory:  dataSourceAlicloudCenTransitRouterMulticastDomainSources,
			},
			{
				TypeName: "alicloud_cen_transit_route_table_aggregations",
				Factory:  dataSourceAlicloudCenTransitRouteTableAggregations,
			},
		},
	}
}
