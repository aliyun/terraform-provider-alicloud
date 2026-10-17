package paiworkspace

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI Workspace
// (AIWorkSpace 2021-02-04) service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "AIWorkSpace",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_workspace_permission",
				Factory:  NewPermissionResource,
			},
		},
	}
}
