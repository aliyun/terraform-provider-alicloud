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
				TypeName: "alicloud_ecs_port_range_list",
				Factory:  NewPortRangeListResource,
			},
		},
	}
}
