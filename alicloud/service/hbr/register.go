package hbr

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Hbr service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "hbr",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_hbr_cross_account",
				Factory:  ResourceCrossAccount,
			},
			{
				TypeName: "alicloud_hbr_ecs_backup_client",
				Factory:  ResourceEcsBackupClient,
			},
			{
				TypeName: "alicloud_hbr_ecs_backup_plan",
				Factory:  ResourceEcsBackupPlan,
			},
			{
				TypeName: "alicloud_hbr_hana_backup_client",
				Factory:  ResourceHanaBackupClient,
			},
			{
				TypeName: "alicloud_hbr_hana_backup_plan",
				Factory:  ResourceHanaBackupPlan,
			},
			{
				TypeName: "alicloud_hbr_hana_instance",
				Factory:  ResourceHanaInstance,
			},
			{
				TypeName: "alicloud_hbr_nas_backup_plan",
				Factory:  ResourceNasBackupPlan,
			},
			{
				TypeName: "alicloud_hbr_oss_backup_plan",
				Factory:  ResourceOssBackupPlan,
			},
			{
				TypeName: "alicloud_hbr_ots_backup_plan",
				Factory:  ResourceOtsBackupPlan,
			},
			{
				TypeName: "alicloud_hbr_policy",
				Factory:  ResourcePolicy,
			},
			{
				TypeName: "alicloud_hbr_policy_binding",
				Factory:  ResourcePolicyBinding,
			},
			{
				TypeName: "alicloud_hbr_replication_vault",
				Factory:  ResourceReplicationVault,
			},
			{
				TypeName: "alicloud_hbr_restore_job",
				Factory:  ResourceRestoreJob,
			},
			{
				TypeName: "alicloud_hbr_server_backup_plan",
				Factory:  ResourceServerBackupPlan,
			},
			{
				TypeName: "alicloud_hbr_vault",
				Factory:  ResourceVault,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_hbr_backup_jobs",
				Factory:  DataSourceBackupJobs,
			},
			{
				TypeName: "alicloud_hbr_ecs_backup_clients",
				Factory:  DataSourceEcsBackupClients,
			},
			{
				TypeName: "alicloud_hbr_ecs_backup_plans",
				Factory:  DataSourceEcsBackupPlans,
			},
			{
				TypeName: "alicloud_hbr_hana_backup_clients",
				Factory:  DataSourceHanaBackupClients,
			},
			{
				TypeName: "alicloud_hbr_hana_backup_plans",
				Factory:  DataSourceHanaBackupPlans,
			},
			{
				TypeName: "alicloud_hbr_hana_instances",
				Factory:  DataSourceHanaInstances,
			},
			{
				TypeName: "alicloud_hbr_nas_backup_plans",
				Factory:  DataSourceNasBackupPlans,
			},
			{
				TypeName: "alicloud_hbr_oss_backup_plans",
				Factory:  DataSourceOssBackupPlans,
			},
			{
				TypeName: "alicloud_hbr_ots_backup_plans",
				Factory:  DataSourceOtsBackupPlans,
			},
			{
				TypeName: "alicloud_hbr_ots_snapshots",
				Factory:  DataSourceOtsSnapshots,
			},
			{
				TypeName: "alicloud_hbr_replication_vault_regions",
				Factory:  DataSourceReplicationVaultRegions,
			},
			{
				TypeName: "alicloud_hbr_restore_jobs",
				Factory:  DataSourceRestoreJobs,
			},
			{
				TypeName: "alicloud_hbr_server_backup_plans",
				Factory:  DataSourceServerBackupPlans,
			},
			{
				TypeName: "alicloud_hbr_service",
				Factory:  DataSourceService,
			},
			{
				TypeName: "alicloud_hbr_snapshots",
				Factory:  DataSourceSnapshots,
			},
			{
				TypeName: "alicloud_hbr_udm_snapshots",
				Factory:  DataSourceUdmSnapshots,
			},
			{
				TypeName: "alicloud_hbr_vaults",
				Factory:  DataSourceVaults,
			},
		},
	}
}
