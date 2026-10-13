package oos

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the OOS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Oos",
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_oos_secret_parameter",
				Factory:  NewSecretParameterEphemeralResource,
			},
		},
	}
}
