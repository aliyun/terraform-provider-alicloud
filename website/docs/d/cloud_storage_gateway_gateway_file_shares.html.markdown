---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_file_shares"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-gateway-file-shares"
description: |-
  Provides a list of Cloud Storage Gateway Gateway File Shares to the user.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_file\_shares

This data source provides the Cloud Storage Gateway Gateway File Shares of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_gateway_file_shares" "ids" {
  gateway_id = "example_value"
  ids        = ["example_value-1", "example_value-2"]
}
output "cloud_storage_gateway_gateway_file_share_id_1" {
  value = data.alicloud_cloud_storage_gateway_gateway_file_shares.ids.shares.0.id
}

data "alicloud_cloud_storage_gateway_gateway_file_shares" "nameRegex" {
  gateway_id = "example_value"
  name_regex = "^my-GatewayFileShare"
}
output "cloud_storage_gateway_gateway_file_share_id_2" {
  value = data.alicloud_cloud_storage_gateway_gateway_file_shares.nameRegex.shares.0.id
}

```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The ID of the gateway.
* `ids` - (Optional, ForceNew, Computed)  A list of Gateway File Share IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Gateway File Share name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Gateway File Share names.
* `shares` - A list of Cloud Storage Gateway Gateway File Shares. Each element contains the following attributes:
  * `access_based_enumeration` - The set up gateway file share Server Message Block (SMB) protocol, whether to enable Windows ABE, the prime minister, need windowsAcl parameter is set to true in the entry into force of. Default value: `false`. **NOTE:** Gateway version >= 1.0.45 above support. 
	* `address` - Share the private IP address of the RDS instance.
	* `backend_limit` - The set up gateway file share Max upload speed. Unit: `MB/s`, `0` means unlimited. Value range: `0` ~ `1280`. Default value: `0`. **NOTE:** at the same time if you have to limit the maximum write speed, maximum upload speed is no less than the maximum write speed. 
	* `browsable` - The set up gateway file share Server Message Block (SMB) protocol whether browsable (that is, in the network neighborhood of whether you can find). The parameters in the NFS protocol not valid under. Default value: `true`.
	* `bucket_infos` - Multi-Bucket information.
	* `buckets_stub` - Whether there are multiple buckets.
	* `cache_mode` - The cache mode of the gateway file share. Value range: Cache: cached mode. Sync: replication mode are available.
	* `client_side_cmk` - File share is enabled to client-side encryption, the encryption by the use of the KMS key. **NOTE:** note: This KMS key must be the gateway and is in the same Region. 
	* `client_side_encryption` - Whether to enabled to client-side encryption of the gateway file share. Default value: `false`. **NOTE:** need to contact us open whitelist before you can the settings, and only supports enhanced more than online gateway, at the same time, server-side encryption and to client-side encryption can not simultaneously configuration. 
	* `direct_io` - Whether directio (direct I/O data transfer) is enabled for file share. Default: `false`.
	* `disk_id` - The ID of the disk.
	* `disk_type` - The cache disk type. Valid values: `cloud_efficiency`: Ultra cloud disk. `cloud_ssd`:SSD cloud disk.
	* `download_limit` - The set up gateway file share maximum download speed. Unit: `MB/s`. `0` means unlimited. Value range: `0` ~ `1280`. **NOTE:** only in copy mode and enable download file data can be set. only when the shared opens the reverse synchronization or acceded to by the speed synchronization Group when, this parameter will not take effect. Gateway version >= 1.3.0 above support. 
	* `enabled` - Shared whether the changes take effect.
	* `express_sync_id` - Speed synchronization group ID.
	* `fast_reclaim` - The set up gateway file share whether to enable Upload optimization, which is suitable for data pure backup migration scenarios. Default value: `false`. **NOTE:** Gateway version >= 1.0.39 above support. 
	* `fe_limit` - The set up gateway file share and the maximum write speed. Unit: `MB/s`, `0` means unlimited. Value range: `0` ~ `1280`. Default value: `0`.
	* `file_num_limit` - Supported by the file system file number.
	* `fs_size_limit` - File system capacity. Unit: `B`.
	* `gateway_file_share_name` - The name of the file share. Length from `1` to `255` characters can contain lowercase letters, digits, (.), (_) Or (-), at the same time, must start with a lowercase letter.
	* `gateway_id` - The ID of the gateway.
	* `id` - The ID of the Gateway File Share.
	* `ignore_delete` - Whether to ignore deleted of the gateway file share. After the opening of the Gateway side delete file or delete cloud (OSS) corresponding to the file. Default value: `false`. **NOTE:** Gateway version >= 1.0.40 above support. 
	* `in_place` - Whether debris optimization of the gateway file share. Default value: `false`.
	* `in_rate` - Cache growth. Unit: `B/s`.
	* `index_id` - The ID of the file share.
	* `kms_rotate_period` - File share is enabled to client-side encryption, key rotation period of time. Seconds. 0 represents no rotation. Rotation of the value range: `3600` ~ `86400`. Default value: `0`.
	* `lag_period` - The synchronization delay, I.e. gateway local cache sync to Alibaba Cloud Object Storage Service (oss) of the delay time. Unit: `Seconds`. Value range: `5` ~ `120`. Default value: `5`. **NOTE:** Gateway version >= 1.0.40 above support. 
	* `local_path` - The cache disk inside the device name.
	* `mns_health` - The messages from the queue health types. Valid values: `TopicAndQueueFailure`: A Message Queuing message theme can be accessed during the black hole period. `TopicFailure`: a message theme can be accessed during the black hole period. `MNSFullSyncInit`: full synchronization wait. `MNSFullSyncing`: full synchronization in progress. `QueueFailure`: a message queue can be accessed during the black hole period. `MNSNotEnabled`: Top speed synchronization is not enabled. `MNSHealthy`: sync fine.
	* `nfs_v4_optimization` - The set up gateway file share NFS protocol, whether to enable NFS v4 optimization improve Mount Upload efficiency. Default value: `false`. **NOTE:** turns on after I will not support NFS v3 mount the filesystem on a. Gateway version >= 1.2.0 above support. 
	* `obsolete_buckets` - Multi-Bucket, removing the Bucket.
	* `oss_bucket_name` - The name of the Bucket.
	* `oss_bucket_ssl` - Whether they are using SSL connect to OSS Bucket.
	* `oss_endpoint` - The set up gateway file share corresponds to the Object Storage SERVICE (OSS), Bucket Endpoint. **NOTE:** distinguish between intranet and internet Endpoint. We recommend that if the OSS Bucket and the gateway is in the same Region is use the RDS intranet IP Endpoint:oss-cn-hangzhou-internal.aliyuncs.com. 
	* `oss_health` - The OSS Bucket of type. Valid values: `BucketHealthy`: OSS connectivity. `BucketAccessDenied`: OBJECT STORAGE Service (OSS) access to an exception. `BucketMiscFailure`: OBJECT STORAGE Service (OSS) access to additional exception. `BucketNetworkFailure`: OBJECT STORAGE Service (OSS) access network an exception. `BucketNotExist`: OSS Bucket does not exist. `Nothing returns`: We may not have ever known existed.
	* `oss_used` - For a cloud-based data is. Unit: `B`.
	* `out_rate` - Upload speed. Unit: `B/s`.
	* `partial_sync_paths` - In part mode, the directory path group JSON format.
	* `path_prefix` - The prefix of the OSS.
	* `polling_interval` - The reverse synchronization time intervals of the gateway file share. Value range: `15` ~ `36000`. **NOTE:** in copy mode + reverse synchronization is enabled Download file data, value range: `3600` ~ `36000`. 
	* `protocol` - Share types. Valid values: `SMB`, `NFS`.
	* `remaining_meta_space` - You can use the metadata space. Unit: `B`.
	* `remote_sync` - Whether to enable reverse synchronization of the gateway file share. Default value: `false`.
	* `remote_sync_download` - Copy mode, whether to download the file data. Default value: `false`. **NOTE:** only when the shared opens the reverse synchronization or acceded to by the speed synchronization group, this parameter will not take effect. 
	* `ro_client_list` - The read-only client list. When Protocol NFS is returned when the status is.
	* `ro_user_list` - The read-only client list. When Protocol for Server Message Block (SMB) to go back to.
	* `rw_client_list` - Read and write the client list. When Protocol NFS is returned when the status is.
	* `rw_user_list` - Read-write user list. When Protocol for Server Message Block (SMB) to go back to.
	* `server_side_cmk` - File share is enabled server-side encryption, encryption used by the KMS key.
	* `server_side_encryption` - If the OSS Bucket side encryption.
	* `size` - The caching capacity. Unit: `B`.
	* `squash` - The set up gateway file share NFS protocol user mapping. Valid values: `none`, `root_squash`, `all_squash`, `all_anonymous`. Default value: `none`.
	* `state` - File synchronization types. Valid values: `clean`, `dirty`. `clean`: synchronization is complete. `dirty`: synchronization has not been completed.
	* `support_archive` - Whether to support the archive transparent read.
	* `sync_progress` - Full synchronization progress. When the share has been added for a synchronization group, the return parameters are valid, that shared full synchronization progress (0~100). `-2`: indicates that share the Gateway version does not support this feature. `-1`: the share does not occur full synchronization.
	* `total_download` - The OSS Bucket to the Gateway total downloads. Unit: `B`.
	* `total_upload` - The OSS Bucket to the Gateway total Upload amount. Unit: `B`.
	* `transfer_acceleration` - The set up gateway file share whether to enable transmission acceleration needs corresponding OSS Bucket enabled transport acceleration. **NOTE:** Gateway version >= 1.3.0 above support. 
	* `used` - Used cache. Unit: `B`.
	* `windows_acl` - The set up gateway file share Server Message Block (SMB) protocol, whether to enable by Windows access list (requires AD domain) the permissions control. Default value: `false`. **NOTE:** Gateway version >= 1.0.45 above support. 
	* `bypass_cache_read` - Direct reading OSS of the gateway file share.