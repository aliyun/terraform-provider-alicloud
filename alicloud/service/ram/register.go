package ram

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the RAM service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ram",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ram_access_key", Factory: ResourceAccessKey},
			{TypeName: "alicloud_ram_access_key_policy", Factory: ResourceAccessKeyPolicy},
			{TypeName: "alicloud_ram_account_alias", Factory: ResourceAccountAlias},
			{TypeName: "alicloud_ram_alias", Factory: ResourceAccountAlias},
			{TypeName: "alicloud_ram_account_password_policy", Factory: ResourceAccountPasswordPolicy},
			{TypeName: "alicloud_ram_group", Factory: ResourceGroup},
			{TypeName: "alicloud_ram_group_membership", Factory: ResourceGroupMembership},
			{TypeName: "alicloud_ram_group_policy_attachment", Factory: ResourceGroupPolicyAttachment},
			{TypeName: "alicloud_ram_login_profile", Factory: ResourceLoginProfile},
			{TypeName: "alicloud_ram_password_policy", Factory: ResourcePasswordPolicy},
			{TypeName: "alicloud_ram_policy", Factory: ResourcePolicy},
			{TypeName: "alicloud_ram_role", Factory: ResourceRole},
			{TypeName: "alicloud_ram_role_policy_attachment", Factory: ResourceRolePolicyAttachment},
			{TypeName: "alicloud_ram_saml_provider", Factory: ResourceSamlProvider},
			{TypeName: "alicloud_ram_security_preference", Factory: ResourceSecurityPreference},
			{TypeName: "alicloud_ram_user", Factory: ResourceUser},
			{TypeName: "alicloud_ram_user_group_attachment", Factory: ResourceUserGroupAttachment},
			{TypeName: "alicloud_ram_user_policy_attachment", Factory: ResourceUserPolicyAttachment},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ram_access_key_policy", Factory: DataSourceAccessKeyPolicy},
			{TypeName: "alicloud_ram_account_alias", Factory: DataSourceAccountAlias},
			{TypeName: "alicloud_ram_account_aliases", Factory: DataSourceAccountAlias},
			{TypeName: "alicloud_ram_groups", Factory: DataSourceGroups},
			{TypeName: "alicloud_ram_policies", Factory: DataSourcePolicies},
			{TypeName: "alicloud_ram_policy_document", Factory: DataSourcePolicyDocument},
			{TypeName: "alicloud_ram_role_policy_attachments", Factory: DataSourceRolePolicyAttachments},
			{TypeName: "alicloud_ram_roles", Factory: DataSourceRoles},
			{TypeName: "alicloud_ram_saml_providers", Factory: DataSourceSamlProviders},
			{TypeName: "alicloud_ram_system_policys", Factory: DataSourceSystemPolicys},
			{TypeName: "alicloud_ram_users", Factory: DataSourceUsers},
		},
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_ram_access_key",
				Factory:  NewRamAccessKeyEphemeralResource,
			},
		},
	}
}
