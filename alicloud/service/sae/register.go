package sae

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Sae service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_sae_application_instances",
				Factory:  DataSourceApplicationInstances,
			},
		},
	}
}
