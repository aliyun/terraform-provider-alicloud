package oss

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the OSS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Oss",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_oss_access_point", Factory: ResourceAccessPoint},
			{TypeName: "alicloud_oss_bucket", Factory: ResourceAlicloudOssBucket},
			{TypeName: "alicloud_oss_bucket_access_monitor", Factory: ResourceAliCloudOssBucketAccessMonitor},
			{TypeName: "alicloud_oss_bucket_acl", Factory: ResourceAliCloudOssBucketAcl},
			{TypeName: "alicloud_oss_bucket_archive_direct_read", Factory: ResourceAliCloudOssBucketArchiveDirectRead},
			{TypeName: "alicloud_oss_bucket_cname", Factory: ResourceAliCloudOssBucketCname},
			{TypeName: "alicloud_oss_bucket_cname_token", Factory: ResourceAliCloudOssBucketCnameToken},
			{TypeName: "alicloud_oss_bucket_cors", Factory: ResourceAliCloudOssBucketCors},
			{TypeName: "alicloud_oss_bucket_data_redundancy_transition", Factory: ResourceAliCloudOssBucketDataRedundancyTransition},
			{TypeName: "alicloud_oss_bucket_https_config", Factory: ResourceAliCloudOssBucketHttpsConfig},
			{TypeName: "alicloud_oss_bucket_inventory", Factory: ResourceAliCloudOssBucketInventory},
			{TypeName: "alicloud_oss_bucket_logging", Factory: ResourceAliCloudOssBucketLogging},
			{TypeName: "alicloud_oss_bucket_meta_query", Factory: ResourceAliCloudOssBucketMetaQuery},
			{TypeName: "alicloud_oss_bucket_object", Factory: ResourceAlicloudOssBucketObject},
			{TypeName: "alicloud_oss_bucket_object_worm_configuration", Factory: ResourceAliCloudOssBucketObjectWormConfiguration},
			{TypeName: "alicloud_oss_bucket_overwrite_config", Factory: ResourceAliCloudOssBucketOverwriteConfig},
			{TypeName: "alicloud_oss_bucket_policy", Factory: ResourceAliCloudOssBucketPolicy},
			{TypeName: "alicloud_oss_bucket_public_access_block", Factory: ResourceAliCloudOssBucketPublicAccessBlock},
			{TypeName: "alicloud_oss_bucket_referer", Factory: ResourceAliCloudOssBucketReferer},
			{TypeName: "alicloud_oss_bucket_replication", Factory: ResourceAlicloudOssBucketReplication},
			{TypeName: "alicloud_oss_bucket_request_payment", Factory: ResourceAliCloudOssBucketRequestPayment},
			{TypeName: "alicloud_oss_bucket_response_header", Factory: ResourceAliCloudOssBucketResponseHeader},
			{TypeName: "alicloud_oss_bucket_server_side_encryption", Factory: ResourceAliCloudOssBucketServerSideEncryption},
			{TypeName: "alicloud_oss_bucket_style", Factory: ResourceAliCloudOssBucketStyle},
			{TypeName: "alicloud_oss_bucket_transfer_acceleration", Factory: ResourceAliCloudOssBucketTransferAcceleration},
			{TypeName: "alicloud_oss_bucket_user_defined_log_fields", Factory: ResourceAliCloudOssBucketUserDefinedLogFields},
			{TypeName: "alicloud_oss_bucket_versioning", Factory: ResourceAliCloudOssBucketVersioning},
			{TypeName: "alicloud_oss_bucket_website", Factory: ResourceAliCloudOssBucketWebsite},
			{TypeName: "alicloud_oss_bucket_worm", Factory: ResourceAliCloudOssBucketWorm},
			{TypeName: "alicloud_oss_account_public_access_block", Factory: ResourceAliCloudOssAccountPublicAccessBlock},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_oss_bucket_inventories", Factory: DataSourceAliCloudOssBucketInventories},
			{TypeName: "alicloud_oss_bucket_objects", Factory: DataSourceAlicloudOssBucketObjects},
			{TypeName: "alicloud_oss_buckets", Factory: DataSourceAlicloudOssBuckets},
			{TypeName: "alicloud_oss_service", Factory: DataSourceAlicloudOssService},
		},
		Resources: []conns.Resource{
			{TypeName: "alicloud_oss_access_point_policy", Factory: NewAccessPointPolicyResource},
		},
	}
}
