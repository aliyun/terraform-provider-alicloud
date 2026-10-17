package paiworkspace

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI Workspace
// service. The POP product code is AIWorkSpace (kept here only for catalog
// reference; the wire constants live in service.go).
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "AIWorkSpace",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_workspace_product",
				Factory:  NewProductResource,
			},
		},
	}
}
