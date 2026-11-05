package cloudsso

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the cloudsso
// (2021-05-15) service: 10 SDKv2 resources and 8 SDKv2 data sources moved
// from the flat alicloud root package into this split-layout product package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "cloudsso",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cloud_sso_access_assignment", Factory: resourceAliCloudCloudSSOAccessAssignment},
			{TypeName: "alicloud_cloud_sso_access_configuration", Factory: resourceAliCloudCloudSsoAccessConfiguration},
			{TypeName: "alicloud_cloud_sso_access_configuration_provisioning", Factory: resourceAlicloudCloudSsoAccessConfigurationProvisioning},
			{TypeName: "alicloud_cloud_sso_delegate_account", Factory: resourceAliCloudCloudSSODelegateAccount},
			{TypeName: "alicloud_cloud_sso_directory", Factory: resourceAliCloudCloudSSODirectory},
			{TypeName: "alicloud_cloud_sso_group", Factory: resourceAlicloudCloudSsoGroup},
			{TypeName: "alicloud_cloud_sso_scim_server_credential", Factory: resourceAliCloudCloudSsoScimServerCredential},
			{TypeName: "alicloud_cloud_sso_user", Factory: resourceAliCloudCloudSsoUser},
			{TypeName: "alicloud_cloud_sso_user_attachment", Factory: resourceAliCloudCloudSSOUserAttachment},
			{TypeName: "alicloud_cloud_sso_user_provisioning", Factory: resourceAliCloudCloudSSOUserProvisioning},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_cloud_sso_access_assignments", Factory: dataSourceAlicloudCloudSsoAccessAssignments},
			{TypeName: "alicloud_cloud_sso_access_configurations", Factory: dataSourceAlicloudCloudSsoAccessConfigurations},
			{TypeName: "alicloud_cloud_sso_directories", Factory: dataSourceAlicloudCloudSsoDirectories},
			{TypeName: "alicloud_cloud_sso_groups", Factory: dataSourceAlicloudCloudSsoGroups},
			{TypeName: "alicloud_cloud_sso_scim_server_credentials", Factory: dataSourceAlicloudCloudSsoScimServerCredentials},
			{TypeName: "alicloud_cloud_sso_service", Factory: dataSourceAlicloudCloudSsoService},
			{TypeName: "alicloud_cloud_sso_user_provisioning_events", Factory: dataSourceAliCloudCloudSsoUserProvisioningEvents},
			{TypeName: "alicloud_cloud_sso_users", Factory: dataSourceAlicloudCloudSsoUsers},
		},
	}
}
