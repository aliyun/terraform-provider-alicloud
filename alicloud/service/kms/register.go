package kms

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the KMS service's SDK v2 resources, data sources
// and framework ephemeral resources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Kms",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_kms_alias", Factory: resourceAlicloudKmsAlias},
			{TypeName: "alicloud_kms_application_access_point", Factory: resourceAliCloudKmsApplicationAccessPoint},
			{TypeName: "alicloud_kms_ciphertext", Factory: resourceAlicloudKmsCiphertext},
			{TypeName: "alicloud_kms_client_key", Factory: resourceAliCloudKmsClientKey},
			{TypeName: "alicloud_kms_instance", Factory: resourceAliCloudKmsInstance},
			{TypeName: "alicloud_kms_key", Factory: resourceAliCloudKmsKey},
			{TypeName: "alicloud_kms_key_version", Factory: resourceAlicloudKmsKeyVersion},
			{TypeName: "alicloud_kms_network_rule", Factory: resourceAliCloudKmsNetworkRule},
			{TypeName: "alicloud_kms_policy", Factory: resourceAliCloudKmsPolicy},
			{TypeName: "alicloud_kms_secret", Factory: resourceAliCloudKmsSecret},
			{TypeName: "alicloud_kms_value_added_service", Factory: resourceAliCloudKmsValueAddedService},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_kms_aliases", Factory: dataSourceAlicloudKmsAliases},
			{TypeName: "alicloud_kms_ciphertext", Factory: dataSourceAlicloudKmsCiphertext},
			{TypeName: "alicloud_kms_instances", Factory: dataSourceAliCloudKmsInstances},
			{TypeName: "alicloud_kms_keys", Factory: dataSourceAlicloudKmsKeys},
			{TypeName: "alicloud_kms_key_versions", Factory: dataSourceAlicloudKmsKeyVersions},
			{TypeName: "alicloud_kms_plaintext", Factory: dataSourceAlicloudKmsPlaintext},
			{TypeName: "alicloud_kms_secrets", Factory: dataSourceAlicloudKmsSecrets},
			{TypeName: "alicloud_kms_secret_versions", Factory: dataSourceAlicloudKmsSecretVersions},
			{TypeName: "alicloud_kms_service", Factory: dataSourceAlicloudKmsService},
		},
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_kms_secret",
				Factory:  NewSecretEphemeralResource,
			},
			{
				TypeName: "alicloud_kms_plaintext",
				Factory:  NewPlaintextEphemeralResource,
			},
		},
	}
}
