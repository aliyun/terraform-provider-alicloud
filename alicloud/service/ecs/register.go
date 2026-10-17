package ecs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ECS service's
// framework-native resources. The legacy ECS resources remain in the flat
// alicloud package's ResourcesMap (SDKv2); this package registers only the
// framework-native resources, starting with the security group snapshot
// policy. A TypeName colliding with the flat map is a startup error caught by
// the registry validator, so the two dialects cannot silently shadow one
// another.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecs",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_security_group_snapshot_policy",
				Factory:  NewSecurityGroupSnapshotPolicyResource,
			},
		},
	}
}
