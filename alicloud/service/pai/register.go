package pai

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI Workspace
// service. It is the only entry point the provider assembly loop reads; the
// Framework-native resources and data sources of this package are declared
// here so they are wired through conns.ServicePackage and surface in both
// the SDKv2 and the Plugin Framework provider servers.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "AIWorkSpace",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_workspace_default_workspace",
				Factory:  NewDefaultWorkspaceResource,
			},
		},
	}
}
