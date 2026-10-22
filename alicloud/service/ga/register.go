package ga

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the GA (Global
// Accelerator) service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ga",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ga_endpoint_group_ip_address_cidr_blocks",
				Factory:  DataSourceEndpointGroupIpAddressCidrBlocks,
			},
		},
	}
}
