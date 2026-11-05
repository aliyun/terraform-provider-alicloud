package sddp

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Sddp service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sddp",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_sddp_config",
				Factory:  ResourceSddpConfig,
			},
			{
				TypeName: "alicloud_sddp_data_limit",
				Factory:  ResourceSddpDataLimit,
			},
			{
				TypeName: "alicloud_sddp_instance",
				Factory:  ResourceSddpInstance,
			},
			{
				TypeName: "alicloud_sddp_rule",
				Factory:  ResourceSddpRule,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_sddp_configs",
				Factory:  DataSourceSddpConfigs,
			},
			{
				TypeName: "alicloud_sddp_data_limits",
				Factory:  DataSourceSddpDataLimits,
			},
			{
				TypeName: "alicloud_sddp_instances",
				Factory:  DataSourceSddpInstances,
			},
			{
				TypeName: "alicloud_sddp_rules",
				Factory:  DataSourceSddpRules,
			},
		},
	}
}
