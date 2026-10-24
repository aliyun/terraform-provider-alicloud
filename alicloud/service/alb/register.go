package alb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ALB service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Alb",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_alb_acl",
				Factory:  ResourceAcl,
			},
			{
				TypeName: "alicloud_alb_acl_entry_attachment",
				Factory:  ResourceAclEntryAttachment,
			},
			{
				TypeName: "alicloud_alb_listener_acl_attachment",
				Factory:  ResourceListenerAclAttachment,
			},
		},
	}
}
