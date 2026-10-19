package cms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CMS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_cms_context_store_api_key",
				Factory:  NewContextStoreAPIKeyResource,
			},
		},
	}
}
