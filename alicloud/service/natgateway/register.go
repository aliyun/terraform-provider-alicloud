package natgateway

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the natgateway service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "NatGateway",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_forward_entry", Factory: resourceAliCloudNatGatewayForwardEntry},
			{TypeName: "alicloud_nat_gateway", Factory: resourceAliCloudNatGateway},
			{TypeName: "alicloud_snat_entry", Factory: resourceAliCloudNATGatewaySnatEntry},
			{TypeName: "alicloud_vpc_nat_ip_cidr", Factory: resourceAliCloudNatGatewayNatIpCidr},
			{TypeName: "alicloud_vpc_nat_ip", Factory: resourceAliCloudNatGatewayNatIp},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_forward_entries", Factory: dataSourceAlicloudForwardEntries},
			{TypeName: "alicloud_nat_gateways", Factory: dataSourceAlicloudNatGateways},
			{TypeName: "alicloud_snat_entries", Factory: dataSourceAlicloudSnatEntries},
			{TypeName: "alicloud_vpc_nat_ip_cidrs", Factory: dataSourceAlicloudVpcNatIpCidrs},
			{TypeName: "alicloud_vpc_nat_ips", Factory: dataSourceAlicloudVpcNatIps},
		},
	}
}
