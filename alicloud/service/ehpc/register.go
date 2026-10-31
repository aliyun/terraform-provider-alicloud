package ehpc

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the EHPC service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "EHPC",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ehpc_clusters",
				Factory:  dataSourceAlicloudEhpcClusters,
			},
			{
				TypeName: "alicloud_ehpc_job_templates",
				Factory:  dataSourceAlicloudEhpcJobTemplates,
			},
			{
				TypeName: "alicloud_ehpc_users",
				Factory:  dataSourceAliCloudEhpcUsers,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ehpc_cluster",
				Factory:  resourceAlicloudEhpcCluster,
			},
			{
				TypeName: "alicloud_ehpc_cluster_v2",
				Factory:  resourceAliCloudEhpcClusterV2,
			},
			{
				TypeName: "alicloud_ehpc_job_template",
				Factory:  resourceAlicloudEhpcJobTemplate,
			},
			{
				TypeName: "alicloud_ehpc_queue",
				Factory:  resourceAliCloudEhpcQueue,
			},
			{
				TypeName: "alicloud_ehpc_user",
				Factory:  resourceAliCloudEhpcUser,
			},
		},
	}
}
