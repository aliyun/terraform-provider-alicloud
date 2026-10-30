package maxcompute

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the maxcompute service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "MaxCompute",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_maxcompute_project", Factory: resourceAliCloudMaxComputeProject},
			{TypeName: "alicloud_max_compute_quota", Factory: resourceAliCloudMaxComputeQuota},
			{TypeName: "alicloud_max_compute_quota_plan", Factory: resourceAliCloudMaxComputeQuotaPlan},
			{TypeName: "alicloud_max_compute_quota_schedule", Factory: resourceAliCloudMaxComputeQuotaSchedule},
			{TypeName: "alicloud_max_compute_role", Factory: resourceAliCloudMaxComputeRole},
			{TypeName: "alicloud_max_compute_role_user_attachment", Factory: resourceAliCloudMaxComputeRoleUserAttachment},
			{TypeName: "alicloud_max_compute_tenant_role_user_attachment", Factory: resourceAliCloudMaxComputeTenantRoleUserAttachment},
			{TypeName: "alicloud_max_compute_tunnel_quota_timer", Factory: resourceAliCloudMaxComputeTunnelQuotaTimer},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_maxcompute_projects", Factory: dataSourceAliCloudMaxComputeProjects},
			{TypeName: "alicloud_maxcompute_service", Factory: dataSourceAlicloudMaxcomputeService},
		},
	}
}
