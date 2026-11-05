package hbase

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the HBase service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "HBase",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_hbase_instance",
				Factory:  ResourceHBaseInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_hbase_instances",
				Factory:  DataSourceHBaseInstances,
			},
			{
				TypeName: "alicloud_hbase_instance_types",
				Factory:  DataSourceHBaseInstanceTypes,
			},
			{
				TypeName: "alicloud_hbase_zones",
				Factory:  DataSourceHBaseZones,
			},
		},
	}
}
