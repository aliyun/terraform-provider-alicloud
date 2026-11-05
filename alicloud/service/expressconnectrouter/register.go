package expressconnectrouter

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the ExpressConnectRouter product's SDKv2 resources
// and data sources that have migrated into this product subpackage.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ExpressConnectRouter",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_express_connect_router_express_connect_router", Factory: ResourceExpressConnectRouter},
			{TypeName: "alicloud_express_connect_router_grant_association", Factory: ResourceGrantAssociation},
			{TypeName: "alicloud_express_connect_router_tr_association", Factory: ResourceTRAssociation},
			{TypeName: "alicloud_express_connect_router_vbr_child_instance", Factory: ResourceVbrChildInstance},
			{TypeName: "alicloud_express_connect_router_vpc_association", Factory: ResourceVpcAssociation},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_express_connect_router_express_connect_routers", Factory: DataSourceExpressConnectRouters},
			{TypeName: "alicloud_express_connect_router_grant_associations", Factory: DataSourceGrantAssociations},
			{TypeName: "alicloud_express_connect_router_tr_associations", Factory: DataSourceTRAssociations},
			{TypeName: "alicloud_express_connect_router_vbr_child_instances", Factory: DataSourceVbrChildInstances},
			{TypeName: "alicloud_express_connect_router_vpc_associations", Factory: DataSourceVpcAssociations},
		},
	}
}
