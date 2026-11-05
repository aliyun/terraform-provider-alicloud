package dbs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Dbs service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dbs",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dbs_backup_plan",
				Factory:  ResourceDbsBackupPlan,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_dbs_backup_plans",
				Factory:  DataSourceDbsBackupPlans,
			},
		},
	}
}
