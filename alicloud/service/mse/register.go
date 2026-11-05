package mse

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the mse (2019-05-31) SDKv2 resources and data sources
// migrated into this package from the root alicloud package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "mse",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_mse_cluster", Factory: ResourceMseCluster},
			{TypeName: "alicloud_mse_gateway", Factory: ResourceMseGateway},
			{TypeName: "alicloud_mse_znode", Factory: ResourceMseZnode},
			{TypeName: "alicloud_mse_engine_namespace", Factory: ResourceMseEngineNamespace},
			{TypeName: "alicloud_mse_nacos_config", Factory: ResourceMseNacosConfig},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_mse_clusters", Factory: DataSourceMseClusters},
			{TypeName: "alicloud_mse_gateways", Factory: DataSourceMseGateways},
			{TypeName: "alicloud_mse_znodes", Factory: DataSourceMseZnodes},
			{TypeName: "alicloud_mse_engine_namespaces", Factory: DataSourceMseEngineNamespaces},
			{TypeName: "alicloud_mse_nacos_configs", Factory: DataSourceMseNacosConfigs},
		},
	}
}
