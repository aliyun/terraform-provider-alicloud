package cloudmonitorservice

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CloudMonitorService
// (Cms) service. It ports the 53 previously root-registered types verbatim,
// preserving the original SDK v2 factories and the single Framework data source
// (alicloud_cms_metric_list). The Name is the API product code "Cms".
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cloud_monitor_service_metric_alarm_rule",
				Factory:  resourceAliCloudCloudMonitorServiceMetricAlarmRule,
			},
			{
				TypeName: "alicloud_cloud_monitor_service_agent_config",
				Factory:  resourceAliCloudCloudMonitorServiceAgentConfig,
			},
			{
				TypeName: "alicloud_cloud_monitor_service_enterprise_public",
				Factory:  resourceAliCloudCloudMonitorServiceEnterprisePublic,
			},
			{
				TypeName: "alicloud_cloud_monitor_service_basic_public",
				Factory:  resourceAliCloudCloudMonitorServiceBasicPublic,
			},
			{
				TypeName: "alicloud_cloud_monitor_service_hybrid_double_write",
				Factory:  resourceAliCloudCloudMonitorServiceHybridDoubleWrite,
			},
			{
				TypeName: "alicloud_cloud_monitor_service_monitoring_agent_process",
				Factory:  resourceAliCloudCloudMonitorServiceMonitoringAgentProcess,
			},
			{
				TypeName: "alicloud_cloud_monitor_service_group_monitoring_agent_process",
				Factory:  resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcess,
			},
			{
				TypeName: "alicloud_cms_event_notify_policy",
				Factory:  resourceAliCloudCmsEventNotifyPolicy,
			},
			{
				TypeName: "alicloud_cms_alert_rule_v2",
				Factory:  resourceAliCloudCmsAlertRuleV2,
			},
			{
				TypeName: "alicloud_cms_agg_task_group",
				Factory:  resourceAliCloudCmsAggTaskGroup,
			},
			{
				TypeName: "alicloud_cms_addon_release",
				Factory:  resourceAliCloudCmsAddonRelease,
			},
			{
				TypeName: "alicloud_cms_prometheus_view",
				Factory:  resourceAliCloudCmsPrometheusView,
			},
			{
				TypeName: "alicloud_cms_prometheus_instance",
				Factory:  resourceAliCloudCmsPrometheusInstance,
			},
			{
				TypeName: "alicloud_cms_integration_policy",
				Factory:  resourceAliCloudCmsIntegrationPolicy,
			},
			{
				TypeName: "alicloud_cms_workspace",
				Factory:  resourceAliCloudCmsWorkspace,
			},
			{
				TypeName: "alicloud_cms_dataset",
				Factory:  resourceAliCloudCmsDataset,
			},
			{
				TypeName: "alicloud_cms_alarm",
				Factory:  resourceAliCloudCmsAlarm,
			},
			{
				TypeName: "alicloud_cms_site_monitor",
				Factory:  resourceAliCloudCloudMonitorServiceSiteMonitor,
			},
			{
				TypeName: "alicloud_cms_alarm_contact",
				Factory:  resourceAlicloudCmsAlarmContact,
			},
			{
				TypeName: "alicloud_cms_alarm_contact_group",
				Factory:  resourceAlicloudCmsAlarmContactGroup,
			},
			{
				TypeName: "alicloud_cms_group_metric_rule",
				Factory:  resourceAliCloudCmsGroupMetricRule,
			},
			{
				TypeName: "alicloud_cms_monitor_group",
				Factory:  resourceAlicloudCmsMonitorGroup,
			},
			{
				TypeName: "alicloud_cms_monitor_group_instances",
				Factory:  resourceAlicloudCmsMonitorGroupInstances,
			},
			{
				TypeName: "alicloud_cms_metric_rule_template",
				Factory:  resourceAliCloudCmsMetricRuleTemplate,
			},
			{
				TypeName: "alicloud_cms_dynamic_tag_group",
				Factory:  resourceAliCloudCmsDynamicTagGroup,
			},
			{
				TypeName: "alicloud_cms_namespace",
				Factory:  resourceAliCloudCmsNamespace,
			},
			{
				TypeName: "alicloud_cms_sls_group",
				Factory:  resourceAlicloudCmsSlsGroup,
			},
			{
				TypeName: "alicloud_cms_hybrid_monitor_sls_task",
				Factory:  resourceAlicloudCmsHybridMonitorSlsTask,
			},
			{
				TypeName: "alicloud_cms_hybrid_monitor_fc_task",
				Factory:  resourceAlicloudCmsHybridMonitorFcTask,
			},
			{
				TypeName: "alicloud_cms_event_rule",
				Factory:  resourceAliCloudCloudMonitorServiceEventRule,
			},
			{
				TypeName: "alicloud_cms_metric_rule_black_list",
				Factory:  resourceAlicloudCmsMetricRuleBlackList,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cloud_monitor_service_metric_alarm_rules",
				Factory:  dataSourceAliCloudCloudMonitorServiceMetricAlarmRules,
			},
			{
				TypeName: "alicloud_cloud_monitor_service_hybrid_double_writes",
				Factory:  dataSourceAliCloudCloudMonitorServiceHybridDoubleWrites,
			},
			{
				TypeName: "alicloud_cms_event_notify_policies",
				Factory:  dataSourceAliCloudCmsEventNotifyPolicies,
			},
			{
				TypeName: "alicloud_cms_alert_rules_v2",
				Factory:  dataSourceAliCloudCmsAlertRulesV2,
			},
			{
				TypeName: "alicloud_cms_service",
				Factory:  dataSourceAlicloudCmsService,
			},
			{
				TypeName: "alicloud_cms_alarm_contacts",
				Factory:  dataSourceAlicloudCmsAlarmContacts,
			},
			{
				TypeName: "alicloud_cms_alarm_contact_groups",
				Factory:  dataSourceAlicloudCmsAlarmContactGroups,
			},
			{
				TypeName: "alicloud_cms_group_metric_rules",
				Factory:  dataSourceAlicloudCmsGroupMetricRules,
			},
			{
				TypeName: "alicloud_cms_monitor_groups",
				Factory:  dataSourceAlicloudCmsMonitorGroups,
			},
			{
				TypeName: "alicloud_cms_monitor_group_instanceses",
				Factory:  dataSourceAlicloudCmsMonitorGroupInstances,
			},
			{
				TypeName: "alicloud_cms_monitor_group_instances",
				Factory:  dataSourceAlicloudCmsMonitorGroupInstances,
			},
			{
				TypeName: "alicloud_cms_metric_rule_templates",
				Factory:  dataSourceAlicloudCmsMetricRuleTemplates,
			},
			{
				TypeName: "alicloud_cms_dynamic_tag_groups",
				Factory:  dataSourceAlicloudCmsDynamicTagGroups,
			},
			{
				TypeName: "alicloud_cms_namespaces",
				Factory:  dataSourceAlicloudCmsNamespaces,
			},
			{
				TypeName: "alicloud_cms_datasets",
				Factory:  dataSourceAliCloudCmsDatasets,
			},
			{
				TypeName: "alicloud_cms_sls_groups",
				Factory:  dataSourceAlicloudCmsSlsGroups,
			},
			{
				TypeName: "alicloud_cms_hybrid_monitor_datas",
				Factory:  dataSourceAlicloudCmsHybridMonitorDatas,
			},
			{
				TypeName: "alicloud_cms_hybrid_monitor_sls_tasks",
				Factory:  dataSourceAlicloudCmsHybridMonitorSlsTasks,
			},
			{
				TypeName: "alicloud_cms_hybrid_monitor_fc_tasks",
				Factory:  dataSourceAlicloudCmsHybridMonitorFcTasks,
			},
			{
				TypeName: "alicloud_cms_event_rules",
				Factory:  dataSourceAlicloudCmsEventRules,
			},
			{
				TypeName: "alicloud_cms_metric_rule_black_lists",
				Factory:  dataSourceAlicloudCmsMetricRuleBlackLists,
			},
			{
				TypeName: "alicloud_cms_site_monitors",
				Factory:  dataSourceAliCloudCloudMonitorServiceSiteMonitors,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_cms_metric_list",
				Factory:  NewMetricListDataSource,
			},
		},
	}
}
