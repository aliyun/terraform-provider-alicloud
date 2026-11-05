package realtimecompute

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Realtime Compute
// (foasconsole 2019-06-01 / ververica 2022-07-18) service. All migrated types
// are SDKv2.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "RealtimeCompute",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_realtime_compute_variables",
				Factory:  DataSourceRealtimeComputeVariables,
			},
			{
				TypeName: "alicloud_realtime_compute_members",
				Factory:  DataSourceRealtimeComputeMembers,
			},
			{
				TypeName: "alicloud_realtime_compute_sql_files",
				Factory:  DataSourceRealtimeComputeSqlFiles,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_realtime_compute_variable",
				Factory:  ResourceRealtimeComputeVariable,
			},
			{
				TypeName: "alicloud_realtime_compute_member",
				Factory:  ResourceRealtimeComputeMember,
			},
			{
				TypeName: "alicloud_realtime_compute_sql_file",
				Factory:  ResourceRealtimeComputeSqlFile,
			},
			{
				TypeName: "alicloud_realtime_compute_deployment",
				Factory:  ResourceRealtimeComputeDeployment,
			},
			{
				TypeName: "alicloud_realtime_compute_job",
				Factory:  ResourceRealtimeComputeJob,
			},
			{
				TypeName: "alicloud_realtime_compute_vvp_instance",
				Factory:  ResourceRealtimeComputeVvpInstance,
			},
		},
	}
}
