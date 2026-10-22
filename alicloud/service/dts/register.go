package dts

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DTS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dts",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dts_migration_instance",
				Factory:  ResourceMigrationInstance,
			},
			{
				TypeName: "alicloud_dts_migration_job",
				Factory:  ResourceMigrationJob,
			},
		},
	}
}
