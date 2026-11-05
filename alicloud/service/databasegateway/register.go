package databasegateway

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Database Gateway
// service. The Name is the canonical product namespace (DatabaseGateway); the
// normalized package directory is alicloud/service/databasegateway.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DatabaseGateway",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_database_gateway_gateway",
				Factory:  ResourceGateway,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_database_gateway_gateways",
				Factory:  DataSourceGateways,
			},
		},
	}
}
