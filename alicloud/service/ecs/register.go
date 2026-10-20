package ecs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ECS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecs",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ecs_snapshot_lock",
				Factory:  NewSnapshotLockResource,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ecs_disk",
				Factory:  ResourceDisk,
			},
			{
				// Deprecated alias kept on the same factory and schema.
				TypeName: "alicloud_disk",
				Factory:  ResourceDisk,
			},
		},
	}
}
