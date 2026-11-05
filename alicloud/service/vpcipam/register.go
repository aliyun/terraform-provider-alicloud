package vpcipam

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the VpcIpam service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "VpcIpam",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_vpc_ipam_ipam", Factory: ResourceVpcIpamIpam},
			{TypeName: "alicloud_vpc_ipam_ipam_pool", Factory: ResourceVpcIpamIpamPool},
			{TypeName: "alicloud_vpc_ipam_ipam_pool_allocation", Factory: ResourceVpcIpamIpamPoolAllocation},
			{TypeName: "alicloud_vpc_ipam_ipam_pool_cidr", Factory: ResourceVpcIpamIpamPoolCidr},
			{TypeName: "alicloud_vpc_ipam_ipam_resource_discovery", Factory: ResourceVpcIpamIpamResourceDiscovery},
			{TypeName: "alicloud_vpc_ipam_ipam_scope", Factory: ResourceVpcIpamIpamScope},
			{TypeName: "alicloud_vpc_ipam_service", Factory: ResourceVpcIpamService},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_vpc_ipam_ipams", Factory: DataSourceVpcIpamIpams},
			{TypeName: "alicloud_vpc_ipam_ipam_pools", Factory: DataSourceVpcIpamIpamPools},
			{TypeName: "alicloud_vpc_ipam_ipam_pool_allocations", Factory: DataSourceVpcIpamIpamPoolAllocations},
			{TypeName: "alicloud_vpc_ipam_ipam_pool_cidrs", Factory: DataSourceVpcIpamIpamPoolCidrs},
			{TypeName: "alicloud_vpc_ipam_ipam_scopes", Factory: DataSourceVpcIpamIpamScopes},
		},
	}
}
