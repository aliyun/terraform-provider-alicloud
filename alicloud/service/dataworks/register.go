package dataworks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataWorks service.
//
// alicloud_data_works_table_lineage is a Framework-native data source: it is
// served by the framework half of the muxed provider, not the SDKv2 map literal
// in provider.go.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DataWorks",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_data_works_table_lineage",
				Factory:  NewTableLineageDataSource,
			},
		},
	}
}
