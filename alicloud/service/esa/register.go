// Package esa registers the ESA service's Terraform Plugin Framework resources
// and data sources. The framework provider collects these declarations through
// the conns.ServicePackage registry; none of them are wired into the SDK v2
// provider.go ResourcesMap.
package esa

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

const (
	errorPagesRedirectRuleTypeName  = "alicloud_esa_error_pages_redirect_rule"
	errorPagesRedirectRulesTypeName = "alicloud_esa_error_pages_redirect_rules"
)

// ServicePackage returns the registration declaration of the ESA service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Esa",
		Resources: []conns.Resource{
			{
				TypeName: errorPagesRedirectRuleTypeName,
				Factory:  NewErrorPagesRedirectRuleResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: errorPagesRedirectRulesTypeName,
				Factory:  NewErrorPagesRedirectRulesDataSource,
			},
		},
	}
}
