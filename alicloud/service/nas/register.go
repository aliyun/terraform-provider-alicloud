package nas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the NAS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "NAS",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_nas_protocol_service", Factory: resourceAliCloudNasProtocolService},
			{TypeName: "alicloud_nas_protocol_mount_target", Factory: resourceAliCloudNasProtocolMountTarget},
			{TypeName: "alicloud_nas_access_point", Factory: resourceAliCloudNasAccessPoint},
			{TypeName: "alicloud_nas_file_system", Factory: resourceAliCloudNasFileSystem},
			{TypeName: "alicloud_nas_mount_target", Factory: resourceAliCloudNasMountTarget},
			{TypeName: "alicloud_nas_access_group", Factory: resourceAliCloudNasAccessGroup},
			{TypeName: "alicloud_nas_access_rule", Factory: resourceAliCloudNasAccessRule},
			{TypeName: "alicloud_nas_log_analysis", Factory: resourceAliCloudNasLogAnalysis},
			{TypeName: "alicloud_nas_smb_acl_attachment", Factory: resourceAlicloudNasSmbAclAttachment},
			{TypeName: "alicloud_nas_snapshot", Factory: resourceAlicloudNasSnapshot},
			{TypeName: "alicloud_nas_fileset", Factory: resourceAliCloudNasFileset},
			{TypeName: "alicloud_nas_auto_snapshot_policy", Factory: resourceAliCloudNasAutoSnapshotPolicy},
			{TypeName: "alicloud_nas_lifecycle_policy", Factory: resourceAlicloudNasLifecyclePolicy},
			{TypeName: "alicloud_nas_data_flow", Factory: resourceAlicloudNasDataFlow},
			{TypeName: "alicloud_nas_recycle_bin", Factory: resourceAlicloudNasRecycleBin},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_nas_service", Factory: dataSourceAlicloudNasService},
			{TypeName: "alicloud_nas_access_groups", Factory: dataSourceAlicloudNasAccessGroups},
			{TypeName: "alicloud_nas_access_rules", Factory: dataSourceAlicloudAccessRules},
			{TypeName: "alicloud_nas_mount_targets", Factory: dataSourceAlicloudNasMountTargets},
			{TypeName: "alicloud_nas_file_systems", Factory: dataSourceAlicloudFileSystems},
			{TypeName: "alicloud_nas_protocols", Factory: dataSourceAlicloudNasProtocols},
			{TypeName: "alicloud_nas_log_analyses", Factory: dataSourceAliCloudNasLogAnalyses},
			{TypeName: "alicloud_nas_zones", Factory: dataSourceAlicloudNasZones},
			{TypeName: "alicloud_nas_snapshots", Factory: dataSourceAlicloudNasSnapshots},
			{TypeName: "alicloud_nas_filesets", Factory: dataSourceAlicloudNasFilesets},
			{TypeName: "alicloud_nas_auto_snapshot_policies", Factory: dataSourceAlicloudNasAutoSnapshotPolicies},
			{TypeName: "alicloud_nas_lifecycle_policies", Factory: dataSourceAlicloudNasLifecyclePolicies},
			{TypeName: "alicloud_nas_data_flows", Factory: dataSourceAlicloudNasDataFlows},
		},
	}
}
