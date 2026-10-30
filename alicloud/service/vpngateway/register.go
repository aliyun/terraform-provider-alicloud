package vpngateway

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the vpngateway service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "VpnGateway",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ssl_vpn_server", Factory: resourceAliyunSslVpnServer},
			{TypeName: "alicloud_ssl_vpn_client_cert", Factory: resourceAliyunSslVpnClientCert},
			{TypeName: "alicloud_vpn_gateway", Factory: resourceAliCloudVPNGatewayVPNGateway},
			{TypeName: "alicloud_vpn_gateway_enhanced_vpn_gateway", Factory: resourceAliCloudVpnGatewayEnhancedVpnGateway},
			{TypeName: "alicloud_vpn_customer_gateway", Factory: resourceAliCloudVPNGatewayCustomerGateway},
			{TypeName: "alicloud_vpn_connection", Factory: resourceAliCloudVPNGatewayVpnConnection},
			{TypeName: "alicloud_vpn_route_entry", Factory: resourceAliyunVpnRouteEntry},
			{TypeName: "alicloud_vpn_pbr_route_entry", Factory: resourceAlicloudVpnPbrRouteEntry},
			{TypeName: "alicloud_vpn_ipsec_server", Factory: resourceAlicloudVpnIpsecServer},
			{TypeName: "alicloud_vpn_gateway_vpn_attachment", Factory: ResourceAliCloudVpnGatewayVpnAttachment},
			{TypeName: "alicloud_vpn_gateway_vco_route", Factory: resourceAlicloudVpnGatewayVcoRoute},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ssl_vpn_servers", Factory: dataSourceAlicloudSslVpnServers},
			{TypeName: "alicloud_ssl_vpn_client_certs", Factory: dataSourceAlicloudSslVpnClientCerts},
			{TypeName: "alicloud_vpn_gateways", Factory: dataSourceAlicloudVpnGateways},
			{TypeName: "alicloud_vpn_gateway_enhanced_vpn_gateways", Factory: dataSourceAliCloudVpnGatewayEnhancedVpnGateways},
			{TypeName: "alicloud_vpn_gateway_zones", Factory: dataSourceAliCloudVPNGatewayZones},
			{TypeName: "alicloud_vpn_customer_gateways", Factory: dataSourceAlicloudVpnCustomerGateways},
			{TypeName: "alicloud_vpn_connections", Factory: dataSourceAlicloudVpnConnections},
			{TypeName: "alicloud_vpn_pbr_route_entries", Factory: dataSourceAlicloudVpnPbrRouteEntries},
			{TypeName: "alicloud_vpn_ipsec_servers", Factory: dataSourceAlicloudVpnIpsecServers},
			{TypeName: "alicloud_vpn_gateway_vpn_attachments", Factory: dataSourceAliCloudVpnGatewayVpnAttachments},
			{TypeName: "alicloud_vpn_gateway_vco_routes", Factory: dataSourceAlicloudVpnGatewayVcoRoutes},
		},
	}
}
