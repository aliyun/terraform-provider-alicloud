package ga

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the GA service. The
// Framework-native service extension resource and service extensions data
// source are declared here so the muxed provider serves them through fwadapt,
// alongside the SDKv2 GA resources still living in the flat alicloud package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ga",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_ga_service_extensions",
				Factory:  NewServiceExtensionsDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ga_service_extension",
				Factory:  NewServiceExtensionResource,
			},
		},
	}
}
