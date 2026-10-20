package ram

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the RAM service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ram",
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_ram_access_key",
				Factory:  NewAccessKeyEphemeralResource,
			},
		},
	}
}
