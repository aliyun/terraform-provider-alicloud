package arms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ARMS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ARMS",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_arms_webhook_contact",
				Factory:  NewWebhookContactResource,
			},
		},
	}
}
