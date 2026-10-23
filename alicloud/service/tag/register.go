package tag

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the TAG service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Tag",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_tag_created_by",
				Factory:  NewCreatedByResource,
			},
		},
	}
}
