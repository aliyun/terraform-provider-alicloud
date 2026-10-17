package pai_workspace

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI Workspace
// service. The connection resource is implemented natively on the Plugin
// Framework; the existing SDKv2 pai_workspace resources remain registered on
// the SDKv2 line in provider.go.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiWorkspace",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_workspace_connection",
				Factory:  NewConnectionResource,
			},
		},
	}
}
