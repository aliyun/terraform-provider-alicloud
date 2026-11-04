package wafopenapi

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the WafOpenapi service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "WafOpenapi",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_wafv3_address_book",
				Factory:  ResourceWafv3AddressBook,
			},
			{
				TypeName: "alicloud_wafv3_defense_resource_group",
				Factory:  ResourceWafv3DefenseResourceGroup,
			},
			{
				TypeName: "alicloud_wafv3_defense_rule",
				Factory:  ResourceWafv3DefenseRule,
			},
			{
				TypeName: "alicloud_wafv3_defense_template",
				Factory:  ResourceWafv3DefenseTemplate,
			},
			{
				TypeName: "alicloud_wafv3_domain",
				Factory:  ResourceWafv3Domain,
			},
			{
				TypeName: "alicloud_wafv3_instance",
				Factory:  ResourceWafv3Instance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_wafv3_address_books",
				Factory:  DataSourceWafv3AddressBooks,
			},
			{
				TypeName: "alicloud_wafv3_defense_rules",
				Factory:  DataSourceWafv3DefenseRules,
			},
			{
				TypeName: "alicloud_wafv3_domains",
				Factory:  DataSourceWafv3Domains,
			},
			{
				TypeName: "alicloud_wafv3_instances",
				Factory:  DataSourceWafv3Instances,
			},
		},
	}
}
