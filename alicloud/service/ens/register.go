// Package ens holds the ENS service's Framework-native resource and data source
// declarations. New ENS managed resources are implemented on the
// terraform-plugin-framework and registered here, not in the legacy SDK v2
// ResourcesMap in package alicloud.
package ens

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ENS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ens",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ens_common_bandwidth_eip_attachment",
				Factory:  NewCommonBandwidthEipAttachmentResource,
			},
		},
	}
}
