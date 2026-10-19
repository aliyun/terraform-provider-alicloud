package dataworks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataWorks service.
//
// DataWorks is onboarded to the Framework-native provider only; its read-only
// Extension data source is collected via the registry service-package loop.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DataWorks",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_data_works_extensions",
				Factory:  NewExtensionsDataSource,
			},
		},
	}
}
