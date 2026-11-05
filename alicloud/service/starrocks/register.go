package starrocks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the StarRocks service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "starrocks",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_star_rocks_instance",
				Factory:  ResourceInstance,
			},
			{
				TypeName: "alicloud_star_rocks_node_group",
				Factory:  ResourceNodeGroup,
			},
		},
	}
}
