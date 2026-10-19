// Package ens holds the Plugin Framework native implementations of ENS
// resources and data sources. New ENS types are added here rather than the
// legacy SDKv2 root package; existing SDKv2 ENS resources in package alicloud
// keep working through the SDKv2 half of the muxed provider.
package ens

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ENS service.
// It is consumed by both provider halves through the registry; only the
// Framework Resources/DataSources declared here are served natively.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ens",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_ens_network_interfaces",
				Factory:  NewNetworkInterfacesDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ens_network_interface",
				Factory:  NewNetworkInterfaceResource,
			},
		},
	}
}
