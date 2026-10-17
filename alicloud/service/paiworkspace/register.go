// Package paiworkspace holds the Terraform Plugin Framework resources and data
// sources for the PAI Workspace (AIWorkSpace) product. New types are
// implemented natively on the framework; the legacy SDK v2 resources for this
// product still live in the root alicloud package until a product-wide
// migration is scheduled.
package paiworkspace

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI Workspace
// service. Framework-native types are appended here as they are added.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiWorkspace",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_workspace_config",
				Factory:  NewConfigResource,
			},
		},
	}
}
