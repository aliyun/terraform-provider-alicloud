package fcv3

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Fcv3 (FC 2023-03-30) service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "FC",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_fcv3_alias",
				Factory:  ResourceAlias,
			},
			{
				TypeName: "alicloud_fcv3_async_invoke_config",
				Factory:  ResourceAsyncInvokeConfig,
			},
			{
				TypeName: "alicloud_fcv3_concurrency_config",
				Factory:  ResourceConcurrencyConfig,
			},
			{
				TypeName: "alicloud_fcv3_custom_domain",
				Factory:  ResourceCustomDomain,
			},
			{
				TypeName: "alicloud_fcv3_function",
				Factory:  ResourceFunction,
			},
			{
				TypeName: "alicloud_fcv3_function_version",
				Factory:  ResourceFunctionVersion,
			},
			{
				TypeName: "alicloud_fcv3_layer_version",
				Factory:  ResourceLayerVersion,
			},
			{
				TypeName: "alicloud_fcv3_provision_config",
				Factory:  ResourceProvisionConfig,
			},
			{
				TypeName: "alicloud_fcv3_trigger",
				Factory:  ResourceTrigger,
			},
			{
				TypeName: "alicloud_fcv3_vpc_binding",
				Factory:  ResourceVpcBinding,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_fcv3_functions",
				Factory:  DataSourceFunctions,
			},
			{
				TypeName: "alicloud_fcv3_triggers",
				Factory:  DataSourceTriggers,
			},
		},
	}
}
