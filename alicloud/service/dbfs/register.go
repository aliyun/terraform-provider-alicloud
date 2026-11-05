package dbfs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DBFS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DBFS",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dbfs_instance",
				Factory:  ResourceInstance,
			},
			{
				TypeName: "alicloud_dbfs_snapshot",
				Factory:  ResourceSnapshot,
			},
			{
				TypeName: "alicloud_dbfs_instance_attachment",
				Factory:  ResourceInstanceAttachment,
			},
			{
				TypeName: "alicloud_dbfs_service_linked_role",
				Factory:  ResourceServiceLinkedRole,
			},
			{
				TypeName: "alicloud_dbfs_auto_snap_shot_policy",
				Factory:  ResourceAutoSnapShotPolicy,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_dbfs_instances",
				Factory:  DataSourceInstances,
			},
			{
				TypeName: "alicloud_dbfs_snapshots",
				Factory:  DataSourceSnapshots,
			},
			{
				TypeName: "alicloud_dbfs_auto_snap_shot_policies",
				Factory:  DataSourceAutoSnapShotPolicies,
			},
		},
	}
}
