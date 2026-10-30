package directmail

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DirectMail (Dm,
// 2015-11-23) service. All resources and data sources are SDKv2 factories.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dm",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_direct_mail_domains",
				Factory:  DataSourceDomains,
			},
			{
				TypeName: "alicloud_direct_mail_mail_addresses",
				Factory:  DataSourceMailAddresses,
			},
			{
				TypeName: "alicloud_direct_mail_receiverses",
				Factory:  DataSourceReceiverses,
			},
			{
				TypeName: "alicloud_direct_mail_tags",
				Factory:  DataSourceTags,
			},
			{
				TypeName: "alicloud_direct_mail_config_sets",
				Factory:  DataSourceConfigSets,
			},
			{
				TypeName: "alicloud_direct_mail_dedicated_ip_pools",
				Factory:  DataSourceDedicatedIpPools,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_direct_mail_domain",
				Factory:  ResourceDomain,
			},
			{
				TypeName: "alicloud_direct_mail_mail_address",
				Factory:  ResourceMailAddress,
			},
			{
				TypeName: "alicloud_direct_mail_receivers",
				Factory:  ResourceReceivers,
			},
			{
				TypeName: "alicloud_direct_mail_tag",
				Factory:  ResourceTag,
			},
			{
				TypeName: "alicloud_direct_mail_config_set",
				Factory:  ResourceConfigSet,
			},
			{
				TypeName: "alicloud_direct_mail_dedicated_ip_pool",
				Factory:  ResourceDedicatedIpPool,
			},
		},
	}
}
