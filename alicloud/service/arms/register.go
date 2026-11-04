package arms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ARMS service.
// The public SDKv2 types (17 resources + 16 data sources, with
// alicloud_arms_prometheus as a public alias reusing DataSourcePrometheis)
// are declared here; the matching entries were removed from provider.go so
// assembleProvider does not flag them as duplicate declarations.
// alicloud_arms_webhook_contact is a Framework resource (NewWebhookContactResource)
// and alicloud_arms_notification_policy is a hand-written SDKv2 resource, both
// migrated into this package from the root package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ARMS",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_arms_webhook_contact",
				Factory:  NewWebhookContactResource,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_arms_synthetic_task",
				Factory:  ResourceSyntheticTask,
			},
			{
				TypeName: "alicloud_arms_grafana_workspace",
				Factory:  ResourceGrafanaWorkspace,
			},
			{
				TypeName: "alicloud_arms_env_custom_job",
				Factory:  ResourceEnvCustomJob,
			},
			{
				TypeName: "alicloud_arms_env_service_monitor",
				Factory:  ResourceEnvServiceMonitor,
			},
			{
				TypeName: "alicloud_arms_env_pod_monitor",
				Factory:  ResourceEnvPodMonitor,
			},
			{
				TypeName: "alicloud_arms_addon_release",
				Factory:  ResourceAddonRelease,
			},
			{
				TypeName: "alicloud_arms_env_feature",
				Factory:  ResourceEnvFeature,
			},
			{
				TypeName: "alicloud_arms_environment",
				Factory:  ResourceEnvironment,
			},
			{
				TypeName: "alicloud_arms_prometheus_monitoring",
				Factory:  ResourcePrometheusMonitoring,
			},
			{
				TypeName: "alicloud_arms_alert_contact",
				Factory:  ResourceAlertContact,
			},
			{
				TypeName: "alicloud_arms_alert_robot",
				Factory:  ResourceAlertRobot,
			},
			{
				TypeName: "alicloud_arms_alert_contact_group",
				Factory:  ResourceAlertContactGroup,
			},
			{
				TypeName: "alicloud_arms_dispatch_rule",
				Factory:  ResourceDispatchRule,
			},
			{
				TypeName: "alicloud_arms_prometheus_alert_rule",
				Factory:  ResourcePrometheusAlertRule,
			},
			{
				TypeName: "alicloud_arms_prometheus",
				Factory:  ResourcePrometheus,
			},
			{
				TypeName: "alicloud_arms_integration_exporter",
				Factory:  ResourceIntegrationExporter,
			},
			{
				TypeName: "alicloud_arms_remote_write",
				Factory:  ResourceRemoteWrite,
			},
			{
				TypeName: "alicloud_arms_notification_policy",
				Factory:  ResourceNotificationPolicy,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_arms_addon_releases",
				Factory:  DataSourceAddonReleases,
			},
			{
				TypeName: "alicloud_arms_env_features",
				Factory:  DataSourceEnvFeatures,
			},
			{
				TypeName: "alicloud_arms_env_custom_jobs",
				Factory:  DataSourceEnvCustomJobs,
			},
			{
				TypeName: "alicloud_arms_env_service_monitors",
				Factory:  DataSourceEnvServiceMonitors,
			},
			{
				TypeName: "alicloud_arms_env_pod_monitors",
				Factory:  DataSourceEnvPodMonitors,
			},
			{
				TypeName: "alicloud_arms_environments",
				Factory:  DataSourceEnvironments,
			},
			{
				TypeName: "alicloud_arms_alert_contacts",
				Factory:  DataSourceAlertContacts,
			},
			{
				TypeName: "alicloud_arms_alert_robots",
				Factory:  DataSourceAlertRobots,
			},
			{
				TypeName: "alicloud_arms_alert_contact_groups",
				Factory:  DataSourceAlertContactGroups,
			},
			{
				TypeName: "alicloud_arms_dispatch_rules",
				Factory:  DataSourceDispatchRules,
			},
			{
				TypeName: "alicloud_arms_prometheus_alert_rules",
				Factory:  DataSourcePrometheusAlertRules,
			},
			{
				TypeName: "alicloud_arms_prometheis",
				Factory:  DataSourcePrometheis,
			},
			{
				TypeName: "alicloud_arms_prometheus",
				Factory:  DataSourcePrometheis,
			},
			{
				TypeName: "alicloud_arms_integration_exporters",
				Factory:  DataSourceIntegrationExporters,
			},
			{
				TypeName: "alicloud_arms_remote_writes",
				Factory:  DataSourceRemoteWrites,
			},
			{
				TypeName: "alicloud_arms_prometheus_monitorings",
				Factory:  DataSourcePrometheusMonitorings,
			},
		},
	}
}
