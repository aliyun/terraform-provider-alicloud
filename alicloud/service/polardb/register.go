package polardb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PolarDB service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PolarDB",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_polardb_backup",
				Factory:  NewBackupResource,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_polardb_node_classes",
				Factory:  DataSourceNodeClasses,
			},
		},
	}
}
