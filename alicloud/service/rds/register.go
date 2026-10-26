package rds

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage is a placeholder for the rds service; declarations move here when the product adopts the split layout.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Rds",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_rds_db_instance_performance",
				Factory:  NewDbInstancePerformanceDataSource,
			},
		},
	}
}
