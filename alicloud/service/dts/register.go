package dts

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DTS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dts",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_dts_check_job",
				Factory:  NewCheckJobResource,
			},
			{
				TypeName: "alicloud_dts_gad_instance",
				Factory:  NewGadInstanceResource,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dts_migration_instance",
				Factory:  ResourceMigrationInstance,
			},
			{
				TypeName: "alicloud_dts_migration_job",
				Factory:  ResourceMigrationJob,
			},
			{
				TypeName: "alicloud_dts_consumer_channel",
				Factory:  ResourceConsumerChannel,
			},
			{
				TypeName: "alicloud_dts_instance",
				Factory:  ResourceInstance,
			},
			{
				TypeName: "alicloud_dts_job_monitor_rule",
				Factory:  ResourceJobMonitorRule,
			},
			{
				TypeName: "alicloud_dts_subscription_job",
				Factory:  ResourceSubscriptionJob,
			},
			{
				TypeName: "alicloud_dts_synchronization_instance",
				Factory:  ResourceSynchronizationInstance,
			},
			{
				TypeName: "alicloud_dts_synchronization_job",
				Factory:  ResourceSynchronizationJob,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_dts_consumer_channels",
				Factory:  DataSourceConsumerChannels,
			},
			{
				TypeName: "alicloud_dts_instances",
				Factory:  DataSourceInstances,
			},
			{
				TypeName: "alicloud_dts_migration_jobs",
				Factory:  DataSourceMigrationJobs,
			},
			{
				TypeName: "alicloud_dts_subscription_jobs",
				Factory:  DataSourceSubscriptionJobs,
			},
			{
				TypeName: "alicloud_dts_synchronization_jobs",
				Factory:  DataSourceSynchronizationJobs,
			},
		},
	}
}
