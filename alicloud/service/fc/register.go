package fc

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the FC service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "FC",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_fc_service",
				Factory:  resourceAlicloudFCService,
			},
			{
				TypeName: "alicloud_fc_function",
				Factory:  resourceAlicloudFCFunction,
			},
			{
				TypeName: "alicloud_fc_trigger",
				Factory:  resourceAlicloudFCTrigger,
			},
			{
				TypeName: "alicloud_fc_alias",
				Factory:  resourceAlicloudFCAlias,
			},
			{
				TypeName: "alicloud_fc_custom_domain",
				Factory:  resourceAlicloudFCCustomDomain,
			},
			{
				TypeName: "alicloud_fc_function_async_invoke_config",
				Factory:  resourceAlicloudFCFunctionAsyncInvokeConfig,
			},
			{
				TypeName: "alicloud_fc_layer_version",
				Factory:  resourceAlicloudFcLayerVersion,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_fc_functions",
				Factory:  dataSourceAlicloudFcFunctions,
			},
			{
				TypeName: "alicloud_fc_services",
				Factory:  dataSourceAlicloudFcServices,
			},
			{
				TypeName: "alicloud_fc_triggers",
				Factory:  dataSourceAlicloudFcTriggers,
			},
			{
				TypeName: "alicloud_fc_custom_domains",
				Factory:  dataSourceAlicloudFcCustomDomains,
			},
			{
				TypeName: "alicloud_fc_zones",
				Factory:  dataSourceAlicloudFcZones,
			},
			{
				TypeName: "alicloud_fc_service",
				Factory:  dataSourceAlicloudFcService,
			},
		},
	}
}
