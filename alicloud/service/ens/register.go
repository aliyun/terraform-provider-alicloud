package ens

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Ens service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ens",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ens_bucket_lifecycle", Factory: ResourceEnsBucketLifecycle},
			{TypeName: "alicloud_ens_disk", Factory: ResourceEnsDisk},
			{TypeName: "alicloud_ens_disk_instance_attachment", Factory: ResourceEnsDiskInstanceAttachment},
			{TypeName: "alicloud_ens_eip", Factory: ResourceEnsEip},
			{TypeName: "alicloud_ens_eip_instance_attachment", Factory: ResourceEnsEipInstanceAttachment},
			{TypeName: "alicloud_ens_image", Factory: ResourceEnsImage},
			{TypeName: "alicloud_ens_instance", Factory: ResourceEnsInstance},
			{TypeName: "alicloud_ens_instance_security_group_attachment", Factory: ResourceEnsInstanceSecurityGroupAttachment},
			{TypeName: "alicloud_ens_key_pair", Factory: ResourceEnsKeyPair},
			{TypeName: "alicloud_ens_load_balancer", Factory: ResourceEnsLoadBalancer},
			{TypeName: "alicloud_ens_load_balancer_http_listener", Factory: ResourceEnsLoadBalancerHttpListener},
			{TypeName: "alicloud_ens_load_balancer_udp_listener", Factory: ResourceEnsLoadBalancerUdpListener},
			{TypeName: "alicloud_ens_nat_gateway", Factory: ResourceEnsNatGateway},
			{TypeName: "alicloud_ens_network", Factory: ResourceEnsNetwork},
			{TypeName: "alicloud_ens_network_route_table", Factory: ResourceEnsNetworkRouteTable},
			{TypeName: "alicloud_ens_security_group", Factory: ResourceEnsSecurityGroup},
			{TypeName: "alicloud_ens_snapshot", Factory: ResourceEnsSnapshot},
			{TypeName: "alicloud_ens_vswitch", Factory: ResourceEnsVswitch},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ens_bucket_lifecycles", Factory: DataSourceEnsBucketLifecycles},
			{TypeName: "alicloud_ens_key_pairs", Factory: DataSourceEnsKeyPairs},
			{TypeName: "alicloud_ens_load_balancer_udp_listeners", Factory: DataSourceEnsLoadBalancerUdpListeners},
			{TypeName: "alicloud_ens_network_route_tables", Factory: DataSourceEnsNetworkRouteTables},
			{TypeName: "alicloud_ens_security_groups", Factory: DataSourceEnsSecurityGroups},
		},
	}
}
