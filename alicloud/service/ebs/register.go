package ebs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Ebs service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ebs",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ebs_dedicated_block_storage_cluster",
				Factory:  ResourceDedicatedBlockStorageCluster,
			},
			{
				TypeName: "alicloud_ebs_disk_replica_group",
				Factory:  ResourceDiskReplicaGroup,
			},
			{
				TypeName: "alicloud_ebs_disk_replica_pair",
				Factory:  ResourceDiskReplicaPair,
			},
			{
				TypeName: "alicloud_ebs_enterprise_snapshot_policy",
				Factory:  ResourceEnterpriseSnapshotPolicy,
			},
			{
				TypeName: "alicloud_ebs_enterprise_snapshot_policy_attachment",
				Factory:  ResourceEnterpriseSnapshotPolicyAttachment,
			},
			{
				TypeName: "alicloud_ebs_replica_group_drill",
				Factory:  ResourceReplicaGroupDrill,
			},
			{
				TypeName: "alicloud_ebs_replica_pair_drill",
				Factory:  ResourceReplicaPairDrill,
			},
			{
				TypeName: "alicloud_ebs_solution_instance",
				Factory:  ResourceSolutionInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ebs_dedicated_block_storage_clusters",
				Factory:  DataSourceDedicatedBlockStorageClusters,
			},
			{
				TypeName: "alicloud_ebs_disk_replica_groups",
				Factory:  DataSourceDiskReplicaGroups,
			},
			{
				TypeName: "alicloud_ebs_disk_replica_pairs",
				Factory:  DataSourceDiskReplicaPairs,
			},
			{
				TypeName: "alicloud_ebs_regions",
				Factory:  DataSourceRegions,
			},
		},
	}
}
