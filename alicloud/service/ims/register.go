package ims

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the IMS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ims",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_ims_default_domain",
				Factory:  NewDefaultDomainDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ims_user",
				Factory:  NewUserResource,
			},
		},
	}
}
