package eci

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ECI service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eci",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_eci_commit_container_task",
				Factory:  NewCommitContainerTaskResource,
			},
		},
	}
}
