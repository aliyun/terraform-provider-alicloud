package threatdetection

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ThreatDetection service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ThreatDetection",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_security_center_group",
				Factory:  resourceAlicloudSecurityCenterGroup,
			},
			{
				TypeName: "alicloud_security_center_service_linked_role",
				Factory:  resourceAlicloudSecurityCenterServiceLinkedRole,
			},
			{
				TypeName: "alicloud_threat_detection_attack_path_whitelist",
				Factory:  resourceAliCloudThreatDetectionAttackPathWhitelist,
			},
			{
				TypeName: "alicloud_threat_detection_service_linked_role",
				Factory:  resourceAliCloudThreatDetectionServiceLinkedRole,
			},
			{
				TypeName: "alicloud_threat_detection_check_config",
				Factory:  resourceAliCloudThreatDetectionCheckConfig,
			},
			{
				TypeName: "alicloud_threat_detection_attack_path_sensitive_asset_config",
				Factory:  resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfig,
			},
			{
				TypeName: "alicloud_threat_detection_cycle_task",
				Factory:  resourceAliCloudThreatDetectionCycleTask,
			},
			{
				TypeName: "alicloud_threat_detection_asset_bind",
				Factory:  resourceAliCloudThreatDetectionAssetBind,
			},
			{
				TypeName: "alicloud_threat_detection_log_meta",
				Factory:  resourceAliCloudThreatDetectionLogMeta,
			},
			{
				TypeName: "alicloud_threat_detection_asset_selection_config",
				Factory:  resourceAliCloudThreatDetectionAssetSelectionConfig,
			},
			{
				TypeName: "alicloud_threat_detection_oss_scan_config",
				Factory:  resourceAliCloudThreatDetectionOssScanConfig,
			},
			{
				TypeName: "alicloud_threat_detection_malicious_file_whitelist_config",
				Factory:  resourceAliCloudThreatDetectionMaliciousFileWhitelistConfig,
			},
			{
				TypeName: "alicloud_threat_detection_monitor_account",
				Factory:  resourceAliCloudThreatDetectionMonitorAccount,
			},
			{
				TypeName: "alicloud_threat_detection_sas_trail",
				Factory:  resourceAliCloudThreatDetectionSasTrail,
			},
			{
				TypeName: "alicloud_threat_detection_image_event_operation",
				Factory:  resourceAliCloudThreatDetectionImageEventOperation,
			},
			{
				TypeName: "alicloud_threat_detection_file_upload_limit",
				Factory:  resourceAliCloudThreatDetectionFileUploadLimit,
			},
			{
				TypeName: "alicloud_threat_detection_rd_default_sync_list",
				Factory:  resourceAliCloudThreatDetectionRdDefaultSyncList,
			},
			{
				TypeName: "alicloud_threat_detection_client_file_protect",
				Factory:  resourceAliCloudThreatDetectionClientFileProtect,
			},
			{
				TypeName: "alicloud_threat_detection_client_user_define_rule",
				Factory:  resourceAliCloudThreatDetectionClientUserDefineRule,
			},
			{
				TypeName: "alicloud_threat_detection_web_lock_config",
				Factory:  resourceAlicloudThreatDetectionWebLockConfig,
			},
			{
				TypeName: "alicloud_threat_detection_backup_policy",
				Factory:  resourceAlicloudThreatDetectionBackupPolicy,
			},
			{
				TypeName: "alicloud_threat_detection_vul_whitelist",
				Factory:  resourceAlicloudThreatDetectionVulWhitelist,
			},
			{
				TypeName: "alicloud_threat_detection_honeypot_node",
				Factory:  resourceAlicloudThreatDetectionHoneypotNode,
			},
			{
				TypeName: "alicloud_threat_detection_baseline_strategy",
				Factory:  resourceAlicloudThreatDetectionBaselineStrategy,
			},
			{
				TypeName: "alicloud_threat_detection_anti_brute_force_rule",
				Factory:  resourceAliCloudThreatDetectionAntiBruteForceRule,
			},
			{
				TypeName: "alicloud_threat_detection_honey_pot",
				Factory:  resourceAlicloudThreatDetectionHoneyPot,
			},
			{
				TypeName: "alicloud_threat_detection_honeypot_probe",
				Factory:  resourceAlicloudThreatDetectionHoneypotProbe,
			},
			{
				TypeName: "alicloud_threat_detection_honeypot_preset",
				Factory:  resourceAlicloudThreatDetectionHoneypotPreset,
			},
			{
				TypeName: "alicloud_threat_detection_instance",
				Factory:  resourceAliCloudThreatDetectionInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_security_center_groups",
				Factory:  dataSourceAlicloudSecurityCenterGroups,
			},
			{
				TypeName: "alicloud_threat_detection_attack_path_whitelists",
				Factory:  dataSourceAliCloudThreatDetectionAttackPathWhitelists,
			},
			{
				TypeName: "alicloud_threat_detection_check_item_configs",
				Factory:  dataSourceAliCloudThreatDetectionCheckItemConfigs,
			},
			{
				TypeName: "alicloud_threat_detection_check_structures",
				Factory:  dataSourceAliCloudThreatDetectionCheckStructures,
			},
			{
				TypeName: "alicloud_threat_detection_web_lock_configs",
				Factory:  dataSourceAlicloudThreatDetectionWebLockConfigs,
			},
			{
				TypeName: "alicloud_threat_detection_backup_policies",
				Factory:  dataSourceAlicloudThreatDetectionBackupPolicies,
			},
			{
				TypeName: "alicloud_threat_detection_vul_whitelists",
				Factory:  dataSourceAlicloudThreatDetectionVulWhitelists,
			},
			{
				TypeName: "alicloud_threat_detection_honeypot_nodes",
				Factory:  dataSourceAlicloudThreatDetectionHoneypotNodes,
			},
			{
				TypeName: "alicloud_threat_detection_baseline_strategies",
				Factory:  dataSourceAlicloudThreatDetectionBaselineStrategies,
			},
			{
				TypeName: "alicloud_threat_detection_assets",
				Factory:  dataSourceAlicloudThreatDetectionAssets,
			},
			{
				TypeName: "alicloud_threat_detection_log_shipper",
				Factory:  dataSourceAlicloudThreatDetectionLogShipper,
			},
			{
				TypeName: "alicloud_threat_detection_rd_default_sync_list",
				Factory:  dataSourceAliCloudThreatDetectionRdDefaultSyncList,
			},
			{
				TypeName: "alicloud_threat_detection_anti_brute_force_rules",
				Factory:  dataSourceAlicloudThreatDetectionAntiBruteForceRules,
			},
			{
				TypeName: "alicloud_threat_detection_honeypot_images",
				Factory:  dataSourceAlicloudThreatDetectionHoneypotImages,
			},
			{
				TypeName: "alicloud_threat_detection_honey_pots",
				Factory:  dataSourceAlicloudThreatDetectionHoneyPots,
			},
			{
				TypeName: "alicloud_threat_detection_honeypot_probes",
				Factory:  dataSourceAlicloudThreatDetectionHoneypotProbes,
			},
			{
				TypeName: "alicloud_threat_detection_honeypot_presets",
				Factory:  dataSourceAlicloudThreatDetectionHoneypotPresets,
			},
			{
				TypeName: "alicloud_threat_detection_instances",
				Factory:  dataSourceAlicloudThreatDetectionInstances,
			},
		},
	}
}
