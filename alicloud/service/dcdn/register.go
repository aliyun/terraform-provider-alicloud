package dcdn

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Dcdn service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "dcdn",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dcdn_domain",
				Factory:  resourceAliCloudDcdnDomain,
			},
			{
				TypeName: "alicloud_dcdn_domain_config",
				Factory:  resourceAliCloudDcdnDomainConfig,
			},
			{
				TypeName: "alicloud_dcdn_ipa_domain",
				Factory:  resourceAlicloudDcdnIpaDomain,
			},
			{
				TypeName: "alicloud_dcdn_waf_policy",
				Factory:  resourceAlicloudDcdnWafPolicy,
			},
			{
				TypeName: "alicloud_dcdn_waf_domain",
				Factory:  resourceAlicloudDcdnWafDomain,
			},
			{
				TypeName: "alicloud_dcdn_waf_policy_domain_attachment",
				Factory:  resourceAlicloudDcdnWafPolicyDomainAttachment,
			},
			{
				TypeName: "alicloud_dcdn_kv_namespace",
				Factory:  resourceAlicloudDcdnKvNamespace,
			},
			{
				TypeName: "alicloud_dcdn_kv",
				Factory:  resourceAlicloudDcdnKv,
			},
			{
				TypeName: "alicloud_dcdn_waf_rule",
				Factory:  resourceAlicloudDcdnWafRule,
			},
			{
				TypeName: "alicloud_dcdn_er",
				Factory:  resourceAlicloudDcdnEr,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_dcdn_service",
				Factory:  dataSourceAlicloudDcdnService,
			},
			{
				TypeName: "alicloud_dcdn_domains",
				Factory:  dataSourceAlicloudDcdnDomains,
			},
			{
				TypeName: "alicloud_dcdn_ipa_domains",
				Factory:  dataSourceAlicloudDcdnIpaDomains,
			},
			{
				TypeName: "alicloud_dcdn_waf_policies",
				Factory:  dataSourceAlicloudDcdnWafPolicies,
			},
			{
				TypeName: "alicloud_dcdn_waf_domains",
				Factory:  dataSourceAlicloudDcdnWafDomains,
			},
			{
				TypeName: "alicloud_dcdn_kv_account",
				Factory:  dataSourceAlicloudDcdnKvAccount,
			},
			{
				TypeName: "alicloud_dcdn_waf_rules",
				Factory:  dataSourceAlicloudDcdnWafRules,
			},
		},
	}
}
