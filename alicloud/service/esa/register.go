package esa

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the Framework-native resources of the ESA product package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ESA",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_esa_error_pages_redirect_rule",
				Factory:  NewErrorPagesRedirectRuleResource,
			},
		},
	}
}
