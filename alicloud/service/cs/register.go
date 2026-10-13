package cs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cs",
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_cs_cluster_credential",
				Factory:  NewClusterCredentialEphemeralResource,
			},
		},
	}
}
