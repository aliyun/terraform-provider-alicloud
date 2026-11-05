package expressconnect

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ExpressConnect service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ExpressConnect",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_express_connect_traffic_qos_rule", Factory: resourceAliCloudExpressConnectTrafficQosRule},
			{TypeName: "alicloud_express_connect_traffic_qos_queue", Factory: resourceAliCloudExpressConnectTrafficQosQueue},
			{TypeName: "alicloud_express_connect_traffic_qos_association", Factory: resourceAliCloudExpressConnectTrafficQosAssociation},
			{TypeName: "alicloud_express_connect_traffic_qos", Factory: resourceAliCloudExpressConnectTrafficQos},
			{TypeName: "alicloud_express_connect_ec_failover_test_job", Factory: resourceAliCloudExpressConnectEcFailoverTestJob},
			{TypeName: "alicloud_express_connect_physical_connection", Factory: resourceAliCloudExpressConnectPhysicalConnection},
			{TypeName: "alicloud_express_connect_virtual_border_router", Factory: resourceAliCloudExpressConnectVirtualBorderRouter},
			{TypeName: "alicloud_express_connect_grant_rule_to_cen", Factory: resourceAliCloudExpressConnectGrantRuleToCen},
			{TypeName: "alicloud_express_connect_virtual_physical_connection", Factory: resourceAlicloudExpressConnectVirtualPhysicalConnection},
			{TypeName: "alicloud_express_connect_vbr_pconn_association", Factory: resourceAliCloudExpressConnectVbrPconnAssociation},
			{TypeName: "alicloud_express_connect_router_interface", Factory: resourceAliCloudExpressConnectRouterInterface},
			{TypeName: "alicloud_router_interface", Factory: resourceAlicloudRouterInterface},
			{TypeName: "alicloud_router_interface_connection", Factory: resourceAlicloudRouterInterfaceConnection},
			{TypeName: "alicloud_vpc_vbr_ha", Factory: resourceAliCloudExpressConnectVbrHa},
			{TypeName: "alicloud_vpc_bgp_group", Factory: resourceAliCloudExpressConnectBgpGroup},
			{TypeName: "alicloud_vpc_bgp_peer", Factory: resourceAliCloudExpressConnectBgpPeer},
			{TypeName: "alicloud_vpc_bgp_network", Factory: resourceAliCloudExpressConnectBgpNetwork},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_express_connect_access_points", Factory: dataSourceAlicloudExpressConnectAccessPoints},
			{TypeName: "alicloud_express_connect_physical_connection_service", Factory: dataSourceAlicloudExpressConnectPhysicalConnectionService},
			{TypeName: "alicloud_express_connect_physical_connections", Factory: dataSourceAlicloudExpressConnectPhysicalConnections},
			{TypeName: "alicloud_express_connect_virtual_border_routers", Factory: dataSourceAlicloudExpressConnectVirtualBorderRouters},
			{TypeName: "alicloud_express_connect_grant_rule_to_cens", Factory: dataSourceAlicloudExpressConnectGrantRuleToCens},
			{TypeName: "alicloud_express_connect_virtual_physical_connections", Factory: dataSourceAlicloudExpressConnectVirtualPhysicalConnections},
			{TypeName: "alicloud_express_connect_vbr_pconn_associations", Factory: dataSourceAlicloudExpressConnectVbrPconnAssociations},
			{TypeName: "alicloud_express_connect_router_interfaces", Factory: dataSourceAlicloudExpressConnectRouterInterfaces},
			{TypeName: "alicloud_router_interfaces", Factory: dataSourceAlicloudRouterInterfaces},
			{TypeName: "alicloud_vpc_bgp_groups", Factory: dataSourceAlicloudVpcBgpGroups},
			{TypeName: "alicloud_vpc_bgp_peers", Factory: dataSourceAlicloudVpcBgpPeers},
			{TypeName: "alicloud_vpc_bgp_networks", Factory: dataSourceAlicloudVpcBgpNetworks},
		},
	}
}
