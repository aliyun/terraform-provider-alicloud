package cloudstoragegateway

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the cloudstoragegateway service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "CloudStorageGateway",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cloud_storage_gateway_storage_bundle", Factory: ResourceAlicloudCloudStorageGatewayStorageBundle},
			{TypeName: "alicloud_cloud_storage_gateway_gateway", Factory: ResourceAliCloudCloudStorageGatewayGateway},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_smb_user", Factory: ResourceAlicloudCloudStorageGatewayGatewaySmbUser},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_cache_disk", Factory: ResourceAliCloudCloudStorageGatewayGatewayCacheDisk},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_logging", Factory: ResourceAlicloudCloudStorageGatewayGatewayLogging},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_block_volume", Factory: ResourceAlicloudCloudStorageGatewayGatewayBlockVolume},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_file_share", Factory: ResourceAlicloudCloudStorageGatewayGatewayFileShare},
			{TypeName: "alicloud_cloud_storage_gateway_express_sync", Factory: ResourceAlicloudCloudStorageGatewayExpressSync},
			{TypeName: "alicloud_cloud_storage_gateway_express_sync_share_attachment", Factory: ResourceAlicloudCloudStorageGatewayExpressSyncShareAttachment},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_cloud_storage_gateway_service", Factory: DataSourceAlicloudCloudStorageGatewayService},
			{TypeName: "alicloud_cloud_storage_gateway_storage_bundles", Factory: DataSourceAlicloudCloudStorageGatewayStorageBundles},
			{TypeName: "alicloud_cloud_storage_gateway_gateways", Factory: DataSourceAlicloudCloudStorageGatewayGateways},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_smb_users", Factory: DataSourceAlicloudCloudStorageGatewayGatewaySmbUsers},
			{TypeName: "alicloud_cloud_storage_gateway_stocks", Factory: DataSourceAlicloudCloudStorageGatewayStocks},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_cache_disks", Factory: DataSourceAlicloudCloudStorageGatewayGatewayCacheDisks},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_block_volumes", Factory: DataSourceAlicloudCloudStorageGatewayGatewayBlockVolumes},
			{TypeName: "alicloud_cloud_storage_gateway_gateway_file_shares", Factory: DataSourceAlicloudCloudStorageGatewayGatewayFileShares},
			{TypeName: "alicloud_cloud_storage_gateway_express_syncs", Factory: DataSourceAlicloudCloudStorageGatewayExpressSyncs},
		},
	}
}
