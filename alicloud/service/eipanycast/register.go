package eipanycast

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Eipanycast service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eipanycast",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_eipanycast_anycast_eip_address",
				Factory:  ResourceAnycastEipAddress,
			},
			{
				TypeName: "alicloud_eipanycast_anycast_eip_address_attachment",
				Factory:  ResourceAnycastEipAddressAttachment,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_eipanycast_anycast_eip_addresses",
				Factory:  DataSourceAnycastEipAddresses,
			},
		},
	}
}
