package vpc

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the VPC service.
// The root alicloud package cannot be imported from a service package, so the
// encapsulated VPC API calls live here with the resources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Vpc",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cen_instance_grant", Factory: ResourceAlicloudCenInstanceGrant},
			{TypeName: "alicloud_havip", Factory: ResourceAliCloudVpcHaVip},
			{TypeName: "alicloud_havip_attachment", Factory: ResourceAliCloudVpcHaVipAttachment},
			{TypeName: "alicloud_network_acl", Factory: ResourceAliCloudVpcNetworkAcl},
			{TypeName: "alicloud_network_acl_attachment", Factory: ResourceAliyunNetworkAclAttachment},
			{TypeName: "alicloud_network_acl_entries", Factory: ResourceAliyunNetworkAclEntries},
			{TypeName: "alicloud_route_entry", Factory: ResourceAliyunRouteEntry},
			{TypeName: "alicloud_route_table", Factory: ResourceAliCloudVpcRouteTable},
			{TypeName: "alicloud_route_table_attachment", Factory: ResourceAliCloudVpcRouteTableAttachment},
			// Deprecated alias kept on the same factory and schema.
			{TypeName: "alicloud_subnet", Factory: ResourceAliCloudVpcVswitch},
			{TypeName: "alicloud_vpc", Factory: ResourceAliCloudVpcVpc},
			{TypeName: "alicloud_vpc_dhcp_options_set", Factory: ResourceAliCloudVpcDhcpOptionsSet},
			{TypeName: "alicloud_vpc_dhcp_options_set_attachment", Factory: ResourceAlicloudVpcDhcpOptionsSetAttachement},
			{TypeName: "alicloud_vpc_flow_log", Factory: ResourceAliCloudVpcFlowLog},
			{TypeName: "alicloud_vpc_gateway_endpoint", Factory: ResourceAliCloudVpcGatewayEndpoint},
			{TypeName: "alicloud_vpc_gateway_endpoint_route_table_attachment", Factory: ResourceAliCloudVpcGatewayEndpointRouteTableAttachment},
			{TypeName: "alicloud_vpc_gateway_route_table_attachment", Factory: ResourceAliCloudVpcGatewayRouteTableAttachment},
			{TypeName: "alicloud_vpc_ha_vip", Factory: ResourceAliCloudVpcHaVip},
			{TypeName: "alicloud_vpc_ipv4_cidr_block", Factory: ResourceAliCloudVpcIpv4CidrBlock},
			{TypeName: "alicloud_vpc_ipv4_gateway", Factory: ResourceAliCloudVpcIpv4Gateway},
			{TypeName: "alicloud_vpc_ipv6_address", Factory: ResourceAliCloudVpcIpv6Address},
			{TypeName: "alicloud_vpc_ipv6_cidr_block", Factory: ResourceAliCloudVpcIpv6CidrBlock},
			{TypeName: "alicloud_vpc_ipv6_egress_rule", Factory: ResourceAliCloudVpcIpv6EgressRule},
			{TypeName: "alicloud_vpc_ipv6_gateway", Factory: ResourceAliCloudVpcIpv6Gateway},
			{TypeName: "alicloud_vpc_ipv6_internet_bandwidth", Factory: ResourceAliCloudVpcIpv6InternetBandwidth},
			{TypeName: "alicloud_vpc_network_acl_attachment", Factory: ResourceAliCloudVpcNetworkAclAttachment},
			{TypeName: "alicloud_vpc_prefix_list", Factory: ResourceAliCloudVpcPrefixList},
			{TypeName: "alicloud_vpc_public_ip_address_pool", Factory: ResourceAliCloudVpcPublicIpAddressPool},
			{TypeName: "alicloud_vpc_public_ip_address_pool_cidr_block", Factory: ResourceAliCloudVpcPublicIpAddressPoolCidrBlock},
			{TypeName: "alicloud_vpc_route_entry", Factory: ResourceAliCloudVpcRouteEntry},
			{TypeName: "alicloud_vpc_route_target_group", Factory: ResourceAliCloudVpcRouteTargetGroup},
			{TypeName: "alicloud_vpc_traffic_mirror_filter", Factory: ResourceAliCloudVpcTrafficMirrorFilter},
			{TypeName: "alicloud_vpc_traffic_mirror_filter_egress_rule", Factory: ResourceAliCloudVpcTrafficMirrorFilterEgressRule},
			{TypeName: "alicloud_vpc_traffic_mirror_filter_ingress_rule", Factory: ResourceAliCloudVpcTrafficMirrorFilterIngressRule},
			{TypeName: "alicloud_vpc_traffic_mirror_session", Factory: ResourceAliCloudVpcTrafficMirrorSession},
			{TypeName: "alicloud_vpc_vswitch_cidr_reservation", Factory: ResourceAliCloudVpcVswitchCidrReservation},
			{TypeName: "alicloud_vswitch", Factory: ResourceAliCloudVpcVswitch},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_enhanced_nat_available_zones", Factory: DataSourceAlicloudEnhancedNatAvailableZones},
			{TypeName: "alicloud_havips", Factory: DataSourceAlicloudHavips},
			{TypeName: "alicloud_network_acls", Factory: DataSourceAlicloudNetworkAcls},
			{TypeName: "alicloud_route_entries", Factory: DataSourceAlicloudRouteEntries},
			{TypeName: "alicloud_route_tables", Factory: DataSourceAlicloudRouteTables},
			{TypeName: "alicloud_vpcs", Factory: DataSourceAlicloudVpcs},
			{TypeName: "alicloud_vswitches", Factory: DataSourceAlicloudVswitches},
			{TypeName: "alicloud_vpc_dhcp_options_sets", Factory: DataSourceAlicloudVpcDhcpOptionsSets},
			{TypeName: "alicloud_vpc_flow_log_service", Factory: DataSourceAliCloudVpcFlowLogService},
			{TypeName: "alicloud_vpc_flow_logs", Factory: DataSourceAlicloudVpcFlowLogs},
			{TypeName: "alicloud_vpc_ipv4_gateways", Factory: DataSourceAlicloudVpcIpv4Gateways},
			{TypeName: "alicloud_vpc_ipv6_addresses", Factory: DataSourceAlicloudVpcIpv6Addresses},
			{TypeName: "alicloud_vpc_ipv6_cidr_blocks", Factory: DataSourceAliCloudVpcIpv6CidrBlocks},
			{TypeName: "alicloud_vpc_ipv6_egress_rules", Factory: DataSourceAlicloudVpcIpv6EgressRules},
			{TypeName: "alicloud_vpc_ipv6_gateways", Factory: DataSourceAlicloudVpcIpv6Gateways},
			{TypeName: "alicloud_vpc_ipv6_internet_bandwidths", Factory: DataSourceAlicloudVpcIpv6InternetBandwidths},
			{TypeName: "alicloud_vpc_prefix_lists", Factory: DataSourceAlicloudVpcPrefixLists},
			{TypeName: "alicloud_vpc_public_ip_address_pool_cidr_blocks", Factory: DataSourceAlicloudVpcPublicIpAddressPoolCidrBlocks},
			{TypeName: "alicloud_vpc_public_ip_address_pools", Factory: DataSourceAlicloudVpcPublicIpAddressPools},
			{TypeName: "alicloud_vpc_route_target_groups", Factory: DataSourceAliCloudVpcRouteTargetGroups},
			{TypeName: "alicloud_vpc_traffic_mirror_filter_egress_rules", Factory: DataSourceAlicloudVpcTrafficMirrorFilterEgressRules},
			{TypeName: "alicloud_vpc_traffic_mirror_filter_ingress_rules", Factory: DataSourceAlicloudVpcTrafficMirrorFilterIngressRules},
			{TypeName: "alicloud_vpc_traffic_mirror_filters", Factory: DataSourceAlicloudVpcTrafficMirrorFilters},
			{TypeName: "alicloud_vpc_traffic_mirror_service", Factory: DataSourceAlicloudVpcTrafficMirrorService},
			{TypeName: "alicloud_vpc_traffic_mirror_sessions", Factory: DataSourceAlicloudVpcTrafficMirrorSessions},
		},
	}
}
