package config

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the config service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Config",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_config_aggregate_compliance_pack", Factory: ResourceAggregateCompliancePack},
			{TypeName: "alicloud_config_aggregate_config_rule", Factory: ResourceAggregateConfigRule},
			{TypeName: "alicloud_config_aggregate_delivery", Factory: ResourceAggregateDelivery},
			{TypeName: "alicloud_config_aggregate_remediation", Factory: ResourceAggregateRemediation},
			{TypeName: "alicloud_config_aggregator", Factory: ResourceAggregator},
			{TypeName: "alicloud_config_compliance_pack", Factory: ResourceCompliancePack},
			{TypeName: "alicloud_config_configuration_recorder", Factory: ResourceConfigurationRecorder},
			{TypeName: "alicloud_config_delivery_channel", Factory: ResourceDeliveryChannel},
			{TypeName: "alicloud_config_delivery", Factory: ResourceDelivery},
			{TypeName: "alicloud_config_remediation", Factory: ResourceRemediation},
			{TypeName: "alicloud_config_report_template", Factory: ResourceReportTemplate},
			{TypeName: "alicloud_config_rule", Factory: ResourceRule},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_config_aggregate_compliance_packs", Factory: DataSourceAggregateCompliancePacks},
			{TypeName: "alicloud_config_aggregate_config_rules", Factory: DataSourceAggregateConfigRules},
			{TypeName: "alicloud_config_aggregate_deliveries", Factory: DataSourceAggregateDeliveries},
			{TypeName: "alicloud_config_aggregators", Factory: DataSourceAggregators},
			{TypeName: "alicloud_config_compliance_packs", Factory: DataSourceCompliancePacks},
			{TypeName: "alicloud_config_configuration_recorders", Factory: DataSourceConfigurationRecorders},
			{TypeName: "alicloud_config_deliveries", Factory: DataSourceDeliveries},
			{TypeName: "alicloud_config_delivery_channels", Factory: DataSourceDeliveryChannels},
			{TypeName: "alicloud_config_rules", Factory: DataSourceRules},
		},
	}
}
