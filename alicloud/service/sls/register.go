package sls

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Sls service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sls",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_log_alert",
				Factory:  resourceAlicloudLogAlert,
			},
			{
				TypeName: "alicloud_log_alert_resource",
				Factory:  resourceAlicloudLogAlertResource,
			},
			{
				TypeName: "alicloud_log_audit",
				Factory:  resourceAlicloudLogAudit,
			},
			{
				TypeName: "alicloud_log_dashboard",
				Factory:  resourceAlicloudLogDashboard,
			},
			{
				TypeName: "alicloud_log_etl",
				Factory:  resourceAlicloudLogETL,
			},
			{
				TypeName: "alicloud_log_ingestion",
				Factory:  resourceAlicloudLogIngestion,
			},
			{
				TypeName: "alicloud_log_machine_group",
				Factory:  resourceAlicloudLogMachineGroup,
			},
			{
				TypeName: "alicloud_log_oss_export",
				Factory:  resourceAlicloudLogOssExport,
			},
			{
				TypeName: "alicloud_log_oss_shipper",
				Factory:  resourceAlicloudLogOssShipper,
			},
			{
				TypeName: "alicloud_log_project",
				Factory:  resourceAliCloudSlsProject,
			},
			{
				TypeName: "alicloud_log_resource",
				Factory:  resourceAlicloudLogResource,
			},
			{
				TypeName: "alicloud_log_resource_record",
				Factory:  resourceAlicloudLogResourceRecord,
			},
			{
				TypeName: "alicloud_log_store",
				Factory:  resourceAliCloudSlsLogStore,
			},
			{
				TypeName: "alicloud_log_store_index",
				Factory:  resourceAlicloudLogStoreIndex,
			},
			{
				TypeName: "alicloud_logtail_config",
				Factory:  resourceAlicloudLogtailConfig,
			},
			{
				TypeName: "alicloud_logtail_attachment",
				Factory:  resourceAlicloudLogtailAttachment,
			},
			{
				TypeName: "alicloud_sls_logtail_pipeline_config",
				Factory:  resourceAliCloudSlsLogtailPipelineConfig,
			},
			{
				TypeName: "alicloud_sls_index",
				Factory:  resourceAliCloudSlsIndex,
			},
			{
				TypeName: "alicloud_sls_machine_group",
				Factory:  resourceAliCloudSlsMachineGroup,
			},
			{
				TypeName: "alicloud_sls_logtail_config",
				Factory:  resourceAliCloudSlsLogtailConfig,
			},
			{
				TypeName: "alicloud_sls_etl",
				Factory:  resourceAliCloudSlsEtl,
			},
			{
				TypeName: "alicloud_sls_oss_export_sink",
				Factory:  resourceAliCloudSlsOssExportSink,
			},
			{
				TypeName: "alicloud_sls_collection_policy",
				Factory:  resourceAliCloudSlsCollectionPolicy,
			},
			{
				TypeName: "alicloud_sls_scheduled_sql",
				Factory:  resourceAliCloudSlsScheduledSql,
			},
			{
				TypeName: "alicloud_sls_alert",
				Factory:  resourceAliCloudSlsAlert,
			},
			{
				TypeName: "alicloud_sls_metric_store",
				Factory:  resourceAliCloudSlsMetricStore,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_log_alert_resource",
				Factory:  dataSourceAlicloudLogAlertResource,
			},
			{
				TypeName: "alicloud_log_service",
				Factory:  dataSourceAlicloudLogService,
			},
			{
				TypeName: "alicloud_log_projects",
				Factory:  dataSourceAlicloudLogProjects,
			},
			{
				TypeName: "alicloud_log_stores",
				Factory:  dataSourceAlicloudLogStores,
			},
			{
				TypeName: "alicloud_sls_indexs",
				Factory:  dataSourceAliCloudSlsIndexs,
			},
			{
				TypeName: "alicloud_sls_machine_groups",
				Factory:  dataSourceAliCloudSlsMachineGroups,
			},
			{
				TypeName: "alicloud_sls_logtail_configs",
				Factory:  dataSourceAliCloudSlsLogtailConfigs,
			},
			{
				TypeName: "alicloud_sls_etls",
				Factory:  dataSourceAliCloudSlsEtls,
			},
			{
				TypeName: "alicloud_sls_alerts",
				Factory:  dataSourceAliCloudSlsAlerts,
			},
			{
				TypeName: "alicloud_sls_metric_stores",
				Factory:  dataSourceAliCloudSlsMetricStores,
			},
		},
	}
}
