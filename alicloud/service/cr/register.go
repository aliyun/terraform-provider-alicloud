package cr

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CR service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cr",
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_cr_ee_authorization_token",
				Factory:  NewAuthorizationTokenEphemeralResource,
			},
		}, Resources: []conns.Resource{{TypeName: "alicloud_cr_repository", Factory: NewRepositoryResource}},
	}
}
