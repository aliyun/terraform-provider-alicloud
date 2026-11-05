package ess

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the public ess (autoscaling) Resource and Data Source
// types that were migrated from the flat alicloud root package. All factories are
// SDKv2 handwritten returning *schema.Resource; alicloud_ess_schedule is a
// deprecated alias that shares the same factory as alicloud_ess_scheduled_task.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ess",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ess_scaling_group", Factory: resourceAlicloudEssScalingGroup},
			{TypeName: "alicloud_ess_eci_scaling_configuration", Factory: resourceAlicloudEssEciScalingConfiguration},
			{TypeName: "alicloud_ess_scaling_configuration", Factory: resourceAlicloudEssScalingConfiguration},
			{TypeName: "alicloud_ess_instance_refresh", Factory: resourceAlicloudEssInstanceRefresh},
			{TypeName: "alicloud_ess_scaling_rule", Factory: resourceAlicloudEssScalingRule},
			{TypeName: "alicloud_ess_scheduled_task", Factory: resourceAlicloudEssScheduledTask},
			// alicloud_ess_schedule is a deprecated alias that points to the same
			// factory as alicloud_ess_scheduled_task (resourceAlicloudEssScheduledTask).
			{TypeName: "alicloud_ess_schedule", Factory: resourceAlicloudEssScheduledTask},
			{TypeName: "alicloud_ess_attachment", Factory: resourceAlicloudEssAttachment},
			{TypeName: "alicloud_ess_suspend_process", Factory: resourceAlicloudEssSuspendProcess},
			{TypeName: "alicloud_ess_lifecycle_hook", Factory: resourceAlicloudEssLifecycleHook},
			{TypeName: "alicloud_ess_notification", Factory: resourceAlicloudEssNotification},
			{TypeName: "alicloud_ess_alarm", Factory: resourceAlicloudEssAlarm},
			{TypeName: "alicloud_ess_scalinggroup_vserver_groups", Factory: resourceAlicloudEssScalingGroupVserverGroups},
			{TypeName: "alicloud_ess_alb_server_group_attachment", Factory: resourceAlicloudEssAlbServerGroupAttachment},
			{TypeName: "alicloud_ess_server_group_attachment", Factory: resourceAliCloudEssServerGroupAttachment},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ess_alarms", Factory: dataSourceAlicloudEssAlarms},
			{TypeName: "alicloud_ess_eci_scaling_configurations", Factory: dataSourceAlicloudEssEciScalingConfigurations},
			{TypeName: "alicloud_ess_notifications", Factory: dataSourceAlicloudEssNotifications},
			{TypeName: "alicloud_ess_scaling_groups", Factory: dataSourceAliCloudEssScalingGroups},
			{TypeName: "alicloud_ess_scaling_rules", Factory: dataSourceAliCloudEssScalingRules},
			{TypeName: "alicloud_ess_scaling_configurations", Factory: dataSourceAlicloudEssScalingConfigurations},
			{TypeName: "alicloud_ess_lifecycle_hooks", Factory: dataSourceAlicloudEssLifecycleHooks},
			{TypeName: "alicloud_ess_scheduled_tasks", Factory: dataSourceAlicloudEssScheduledTasks},
		},
	}
}
