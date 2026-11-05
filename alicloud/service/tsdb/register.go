package tsdb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the tsdb service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "tsdb",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_tsdb_instance", Factory: ResourceTsdbInstance},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_tsdb_instances", Factory: DataSourceTsdbInstances},
			{TypeName: "alicloud_tsdb_zones", Factory: DataSourceTsdbZones},
		},
	}
}
