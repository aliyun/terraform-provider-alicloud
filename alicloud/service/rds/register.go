package rds

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the RDS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Rds",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_rds_db_instance_performance",
				Factory:  NewDbInstancePerformanceDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_rds_rc_cluster_node_pool",
				Factory:  NewRCClusterNodePoolResource,
			},
			{
				TypeName: "alicloud_rds_rc_cluster",
				Factory:  NewRCClusterResource,
			},
		},
	}
}
