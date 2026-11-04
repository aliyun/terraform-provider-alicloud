package das

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the DAS product's SDKv2 resources and data sources
// with the provider registry. Factories are the package-level Resource* /
// DataSource* constructors; assembleProvider wires each entry under its public
// Terraform type name. The corresponding map-literal entries in provider.go
// are removed so assembleProvider's duplicateDeclaration guard does not fire at
// startup.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DAS",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_das_sql_log_config",
				Factory:  ResourceDasSqlLogConfig,
			},
			{
				TypeName: "alicloud_das_switch_das_pro",
				Factory:  ResourceDasSwitchDasPro,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_das_sql_log_configs",
				Factory:  DataSourceDasSqlLogConfigs,
			},
		},
	}
}
