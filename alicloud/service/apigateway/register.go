package apigateway

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the apigateway service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ApiGateway",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_api_gateway_access_control_list", Factory: ResourceAccessControlList},
			{TypeName: "alicloud_api_gateway_acl_entry_attachment", Factory: ResourceAclEntryAttachment},
			{TypeName: "alicloud_api_gateway_api", Factory: ResourceApi},
			{TypeName: "alicloud_api_gateway_app", Factory: ResourceApp},
			{TypeName: "alicloud_api_gateway_app_attachment", Factory: ResourceAppAttachment},
			{TypeName: "alicloud_api_gateway_backend", Factory: ResourceBackend},
			{TypeName: "alicloud_api_gateway_backend_model", Factory: ResourceBackendModel},
			{TypeName: "alicloud_api_gateway_group", Factory: ResourceGroup},
			{TypeName: "alicloud_api_gateway_group_plugin_attachment", Factory: ResourceGroupPluginAttachment},
			{TypeName: "alicloud_api_gateway_instance", Factory: ResourceInstance},
			{TypeName: "alicloud_api_gateway_instance_acl_attachment", Factory: ResourceInstanceAclAttachment},
			{TypeName: "alicloud_api_gateway_log_config", Factory: ResourceLogConfig},
			{TypeName: "alicloud_api_gateway_model", Factory: ResourceModel},
			{TypeName: "alicloud_api_gateway_plugin", Factory: ResourcePlugin},
			{TypeName: "alicloud_api_gateway_plugin_attachment", Factory: ResourcePluginAttachment},
			{TypeName: "alicloud_api_gateway_stage_model", Factory: ResourceStageModel},
			{TypeName: "alicloud_api_gateway_vpc_access", Factory: ResourceVpcAccess},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_api_gateway_apis", Factory: DataSourceApis},
			{TypeName: "alicloud_api_gateway_apps", Factory: DataSourceApps},
			{TypeName: "alicloud_api_gateway_backends", Factory: DataSourceBackends},
			{TypeName: "alicloud_api_gateway_groups", Factory: DataSourceGroups},
			{TypeName: "alicloud_api_gateway_log_configs", Factory: DataSourceLogConfigs},
			{TypeName: "alicloud_api_gateway_models", Factory: DataSourceModels},
			{TypeName: "alicloud_api_gateway_plugins", Factory: DataSourcePlugins},
			{TypeName: "alicloud_api_gateway_service", Factory: DataSourceService},
		},
	}
}
