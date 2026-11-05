package ims

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the IMS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ims",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ims_oidc_providers",
				Factory:  DataSourceImsOidcProviders,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ims_oidc_provider",
				Factory:  ResourceImsOidcProvider,
			},
		},
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
