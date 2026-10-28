package cddc

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cddc service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "cddc",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cddc_dedicated_host",
				Factory:  ResourceDedicatedHost,
			},
			{
				TypeName: "alicloud_cddc_dedicated_host_account",
				Factory:  ResourceDedicatedHostAccount,
			},
			{
				TypeName: "alicloud_cddc_dedicated_host_group",
				Factory:  ResourceDedicatedHostGroup,
			},
			{
				TypeName: "alicloud_cddc_dedicated_propre_host",
				Factory:  ResourceDedicatedPropreHost,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cddc_dedicated_host_accounts",
				Factory:  DataSourceDedicatedHostAccounts,
			},
			{
				TypeName: "alicloud_cddc_dedicated_host_groups",
				Factory:  DataSourceDedicatedHostGroups,
			},
			{
				TypeName: "alicloud_cddc_dedicated_hosts",
				Factory:  DataSourceDedicatedHosts,
			},
			{
				TypeName: "alicloud_cddc_host_ecs_level_infos",
				Factory:  DataSourceHostEcsLevelInfos,
			},
			{
				TypeName: "alicloud_cddc_zones",
				Factory:  DataSourceZones,
			},
		},
	}
}
