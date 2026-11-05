package resourcemanager

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers all ResourceManager Resource/Data Source factories with the provider registry.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ResourceManager",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_resource_manager_account", Factory: ResourceAccount},
			{TypeName: "alicloud_resource_manager_auto_grouping_rule", Factory: ResourceAutoGroupingRule},
			{TypeName: "alicloud_resource_manager_control_policy", Factory: ResourceControlPolicy},
			{TypeName: "alicloud_resource_manager_control_policy_attachment", Factory: ResourceControlPolicyAttachment},
			{TypeName: "alicloud_resource_manager_delegated_administrator", Factory: ResourceDelegatedAdministrator},
			{TypeName: "alicloud_resource_manager_delivery_channel", Factory: ResourceDeliveryChannel},
			{TypeName: "alicloud_resource_manager_folder", Factory: ResourceFolder},
			{TypeName: "alicloud_resource_manager_handshake", Factory: ResourceHandshake},
			{TypeName: "alicloud_resource_manager_handshake_acceptance", Factory: ResourceHandshakeAcceptance},
			{TypeName: "alicloud_resource_manager_message_contact", Factory: ResourceMessageContact},
			{TypeName: "alicloud_resource_manager_multi_account_delivery_channel", Factory: ResourceMultiAccountDeliveryChannel},
			{TypeName: "alicloud_resource_manager_policy", Factory: ResourcePolicy},
			{TypeName: "alicloud_resource_manager_policy_attachment", Factory: ResourcePolicyAttachment},
			{TypeName: "alicloud_resource_manager_policy_version", Factory: ResourcePolicyVersion},
			{TypeName: "alicloud_resource_manager_resource_directory", Factory: ResourceResourceDirectory},
			{TypeName: "alicloud_resource_manager_resource_directory_sharing", Factory: ResourceResourceDirectorySharing},
			{TypeName: "alicloud_resource_manager_resource_group", Factory: ResourceResourceGroup},
			{TypeName: "alicloud_resource_manager_resource_group_settings", Factory: ResourceResourceGroupSettings},
			{TypeName: "alicloud_resource_manager_resource_share", Factory: ResourceResourceShare},
			{TypeName: "alicloud_resource_manager_role", Factory: ResourceRole},
			{TypeName: "alicloud_resource_manager_saved_query", Factory: ResourceSavedQuery},
			{TypeName: "alicloud_resource_manager_service_linked_role", Factory: ResourceServiceLinkedRole},
			{TypeName: "alicloud_resource_manager_shared_resource", Factory: ResourceSharedResource},
			{TypeName: "alicloud_resource_manager_shared_target", Factory: ResourceSharedTarget},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_resource_manager_account_deletion_check_task", Factory: DataSourceAccountDeletionCheckTask},
			{TypeName: "alicloud_resource_manager_accounts", Factory: DataSourceAccounts},
			{TypeName: "alicloud_resource_manager_control_policies", Factory: DataSourceControlPolicies},
			{TypeName: "alicloud_resource_manager_control_policy_attachments", Factory: DataSourceControlPolicyAttachments},
			{TypeName: "alicloud_resource_manager_delegated_administrators", Factory: DataSourceDelegatedAdministrators},
			{TypeName: "alicloud_resource_manager_folders", Factory: DataSourceFolders},
			{TypeName: "alicloud_resource_manager_handshakes", Factory: DataSourceHandshakes},
			{TypeName: "alicloud_resource_manager_policies", Factory: DataSourcePolicies},
			{TypeName: "alicloud_resource_manager_policy_attachments", Factory: DataSourcePolicyAttachments},
			{TypeName: "alicloud_resource_manager_policy_versions", Factory: DataSourcePolicyVersions},
			{TypeName: "alicloud_resource_manager_resource_directories", Factory: DataSourceResourceDirectories},
			{TypeName: "alicloud_resource_manager_resource_groups", Factory: DataSourceResourceGroups},
			{TypeName: "alicloud_resource_manager_resource_shares", Factory: DataSourceResourceShares},
			{TypeName: "alicloud_resource_manager_roles", Factory: DataSourceRoles},
			{TypeName: "alicloud_resource_manager_shared_resources", Factory: DataSourceSharedResources},
			{TypeName: "alicloud_resource_manager_shared_targets", Factory: DataSourceSharedTargets},
		},
	}
}
