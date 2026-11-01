package oos

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the OOS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Oos",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_oos_application", Factory: ResourceApplication},
			{TypeName: "alicloud_oos_application_group", Factory: ResourceApplicationGroup},
			{TypeName: "alicloud_oos_default_patch_baseline", Factory: ResourceDefaultPatchBaseline},
			{TypeName: "alicloud_oos_execution", Factory: ResourceExecution},
			{TypeName: "alicloud_oos_parameter", Factory: ResourceParameter},
			{TypeName: "alicloud_oos_patch_baseline", Factory: ResourcePatchBaseline},
			{TypeName: "alicloud_oos_secret_parameter", Factory: ResourceSecretParameter},
			{TypeName: "alicloud_oos_service_setting", Factory: ResourceServiceSetting},
			{TypeName: "alicloud_oos_state_configuration", Factory: ResourceStateConfiguration},
			{TypeName: "alicloud_oos_template", Factory: ResourceTemplate},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_oos_application_groups", Factory: DataSourceApplicationGroups},
			{TypeName: "alicloud_oos_applications", Factory: DataSourceApplications},
			{TypeName: "alicloud_oos_executions", Factory: DataSourceExecutions},
			{TypeName: "alicloud_oos_parameters", Factory: DataSourceParameters},
			{TypeName: "alicloud_oos_patch_baselines", Factory: DataSourcePatchBaselines},
			{TypeName: "alicloud_oos_secret_parameters", Factory: DataSourceSecretParameters},
			{TypeName: "alicloud_oos_state_configurations", Factory: DataSourceStateConfigurations},
			{TypeName: "alicloud_oos_templates", Factory: DataSourceTemplates},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_oos_inventories",
				Factory:  NewInventoriesDataSource,
			},
		},
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_oos_secret_parameter",
				Factory:  NewSecretParameterEphemeralResource,
			},
		},
	}
}
