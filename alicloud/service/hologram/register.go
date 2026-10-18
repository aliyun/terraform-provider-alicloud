package hologram

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Hologram service.
// The Hologram instance remains an SDK v2 resource in package alicloud; only the
// Framework-native User resource and Users data source are declared here.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Hologram",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_hologram_user",
				Factory:  NewUserResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_hologram_users",
				Factory:  NewUsersDataSource,
			},
		},
	}
}
