package apig

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Apig service.
// The factories live in this package and are wired into the provider through
// the registry; the SDKResources/SDKDataSources below replace the lines
// previously maintained inline in alicloud/provider.go.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "APIG",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_apig_secret", Factory: ResourceSecret},
			{TypeName: "alicloud_apig_policy", Factory: ResourcePolicy},
			{TypeName: "alicloud_apig_source", Factory: ResourceSource},
			{TypeName: "alicloud_apig_ai_model_provider", Factory: ResourceAiModelProvider},
			{TypeName: "alicloud_apig_service", Factory: ResourceService},
			{TypeName: "alicloud_apig_plugin", Factory: ResourcePlugin},
			{TypeName: "alicloud_apig_plugin_class", Factory: ResourcePluginClass},
			{TypeName: "alicloud_apig_domain", Factory: ResourceDomain},
			{TypeName: "alicloud_apig_route", Factory: ResourceRoute},
			{TypeName: "alicloud_apig_environment", Factory: ResourceEnvironment},
			{TypeName: "alicloud_apig_gateway", Factory: ResourceGateway},
			{TypeName: "alicloud_apig_http_api", Factory: ResourceHttpApi},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_apig_secrets", Factory: DataSourceSecrets},
			{TypeName: "alicloud_apig_policies", Factory: DataSourcePolicies},
			{TypeName: "alicloud_apig_sources", Factory: DataSourceSources},
			{TypeName: "alicloud_apig_ai_model_providers", Factory: DataSourceAiModelProviders},
			{TypeName: "alicloud_apig_services", Factory: DataSourceServices},
			{TypeName: "alicloud_apig_gateways", Factory: DataSourceGateways},
			{TypeName: "alicloud_apig_plugins", Factory: DataSourcePlugins},
			{TypeName: "alicloud_apig_domains", Factory: DataSourceDomains},
			{TypeName: "alicloud_apig_routes", Factory: DataSourceRoutes},
			{TypeName: "alicloud_apig_plugin_classes", Factory: DataSourcePluginClasses},
		},
	}
}
