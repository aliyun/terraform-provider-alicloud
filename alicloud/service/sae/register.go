package sae

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_sae_application_instances",
				Factory:  NewApplicationInstancesDataSource,
			},
		},
	}
}
