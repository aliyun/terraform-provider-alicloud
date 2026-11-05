package privatelink

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the privatelink service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Privatelink",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_privatelink_vpc_endpoint_service", Factory: ResourceVpcEndpointService},
			{TypeName: "alicloud_privatelink_vpc_endpoint", Factory: ResourceVpcEndpoint},
			{TypeName: "alicloud_privatelink_vpc_endpoint_connection", Factory: ResourceVpcEndpointConnection},
			{TypeName: "alicloud_privatelink_vpc_endpoint_service_resource", Factory: ResourceVpcEndpointServiceResource},
			{TypeName: "alicloud_privatelink_vpc_endpoint_service_user", Factory: ResourceVpcEndpointServiceUser},
			{TypeName: "alicloud_privatelink_vpc_endpoint_zone", Factory: ResourceVpcEndpointZone},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_privatelink_service", Factory: DataSourceService},
			{TypeName: "alicloud_privatelink_vpc_endpoint_services", Factory: DataSourceVpcEndpointServices},
			{TypeName: "alicloud_privatelink_vpc_endpoints", Factory: DataSourceVpcEndpoints},
			{TypeName: "alicloud_privatelink_vpc_endpoint_connections", Factory: DataSourceVpcEndpointConnections},
			{TypeName: "alicloud_privatelink_vpc_endpoint_service_resources", Factory: DataSourceVpcEndpointServiceResources},
			{TypeName: "alicloud_privatelink_vpc_endpoint_service_users", Factory: DataSourceVpcEndpointServiceUsers},
			{TypeName: "alicloud_privatelink_vpc_endpoint_zones", Factory: DataSourceVpcEndpointZones},
		},
	}
}
