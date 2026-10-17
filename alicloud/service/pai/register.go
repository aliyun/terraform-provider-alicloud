package pai

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI EAS service.
//
// The SDK v2 PAI resources (alicloud_pai_service and the rest) stay registered
// through the legacy provider schema map; this package owns only the framework
// resources migrated off that line, so alicloud_pai_virtual_resource is the sole
// entry here today.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Pai",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_virtual_resource",
				Factory:  NewVirtualResource,
			},
		},
	}
}
