package dataworks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataWorks service.
//
// DataWorks scheduling resources and data sources remain in the SDK v2 provider.go
// map literal; this package only declares Framework-native additions, starting with
// the alicloud_data_works_node_io data source.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DataWorks",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_data_works_node_io",
				Factory:  NewNodeIoDataSource,
			},
		},
	}
}
