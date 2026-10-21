package eci

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "eci", Resources: []conns.Resource{{TypeName: "alicloud_eci_commit_container_task", Factory: NewCommitContainerTaskResource}}}
}
