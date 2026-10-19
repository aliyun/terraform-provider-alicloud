package pai

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI service.
//
// Existing alicloud_pai_* resources are SDK v2 implementations registered in
// the legacy provider map. New Framework-native PAI resources live in this
// service package and are registered here; the SDK v2 set is unaffected.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Pai",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_virtual_resource",
				Factory:  NewVirtualResourceResource,
			},
		},
	}
}
