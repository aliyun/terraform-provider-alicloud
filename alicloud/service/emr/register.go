package emr

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Emr service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Emr",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_emr_cluster",
				Factory:  ResourceEmrCluster,
			},
			{
				TypeName: "alicloud_emrv2_cluster",
				Factory:  ResourceEmrV2Cluster,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_emr_clusters",
				Factory:  DataSourceEmrClusters,
			},
			{
				TypeName: "alicloud_emr_disk_types",
				Factory:  DataSourceEmrDiskTypes,
			},
			{
				TypeName: "alicloud_emr_instance_types",
				Factory:  DataSourceEmrInstanceTypes,
			},
			{
				TypeName: "alicloud_emr_main_versions",
				Factory:  DataSourceEmrMainVersions,
			},
			{
				TypeName: "alicloud_emrv2_cluster_instances",
				Factory:  DataSourceEmrV2ClusterInstances,
			},
			{
				TypeName: "alicloud_emrv2_clusters",
				Factory:  DataSourceEmrV2Clusters,
			},
		},
	}
}
