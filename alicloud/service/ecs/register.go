package ecs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ECS service.
// ECS is the first Framework-native ECS service package in this provider; it
// is wired into the mux provider through provider/registry.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecs",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ecs_vsc",
				Factory:  NewEcsVscResource,
			},
		},
	}
}
