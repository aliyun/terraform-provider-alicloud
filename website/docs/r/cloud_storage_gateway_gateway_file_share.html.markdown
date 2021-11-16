---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_file_share"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-file-share"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway File Share resource.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_file\_share

Provides a Cloud Storage Gateway Gateway File Share resource.

For information about Cloud Storage Gateway Gateway File Share and how to use it, see [What is Gateway File Share](https://www.alibabacloud.com/help/zh/doc-detail/170298.htm).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "example_value"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = "example_value"
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = "example_value"
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = "example_value"
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = "example_value"
}

resource "alicloud_cloud_storage_gateway_gateway_file_share" "default" {
  gateway_file_share_name = "example_value"
  gateway_id              = alicloud_cloud_storage_gateway_gateway.default.id
  local_path              = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name         = alicloud_oss_bucket.default.bucket
  oss_endpoint            = alicloud_oss_bucket.default.extranet_endpoint
  protocol                = "NFS"
  remote_sync             = true
  polling_interval        = 4500
  fe_limit                = 0
  backend_limit           = 0
  cache_mode              = "Cache"
  squash                  = "none"
  lag_period              = 5
}

```

## Argument Reference

The following arguments are supported:

* `access_based_enumeration` - (Optional, Computed) Whether to enable Windows ABE, the prime minister, need windowsAcl parameter is set to true in the entry into force of. Default value: `false`. **NOTE:** The attribute is valid when the attribute `protocol` is `SMB`. Gateway version >= 1.0.45 above support. 
* `backend_limit` - (Optional, Computed) The Max upload speed of the gateway file share. Unit: `MB/s`, 0 means unlimited. Value range: `0` ~ `1280`. Default value: `0`. **NOTE:** at the same time if you have to limit the maximum write speed, maximum upload speed is no less than the maximum write speed. 
* `browsable` - (Optional, Computed) The whether browsable of the gateway file share (that is, in the network neighborhood of whether you can find). The attribute is valid when the attribute `protocol` is `SMB`. Default value: `true`.
* `cache_mode` - (Optional, Computed, ForceNew) The set up gateway file share cache mode. Value values: `Cache` or `Sync`. `Cache`: cached mode. `Sync`: replication mode are available. Default value: `Cache`.
* `direct_io` - (Optional, Computed, ForceNew) File sharing Whether to enable DirectIO (direct I/O mode for data transmission). Default value: `false`.
* `download_limit` - (Optional, Computed) The maximum download speed of the gateway file share. Unit: `MB/s`. `0` means unlimited. Value range: `0` ~ `1280`. **NOTE:** only in copy mode and enable download file data can be set. only when the shared opens the reverse synchronization or acceded to by the speed synchronization Group when, this parameter will not take effect. Gateway version >= 1.3.0 above support. 
* `fast_reclaim` - (Optional, Computed) The whether to enable Upload optimization of the gateway file share, which is suitable for data pure backup migration scenarios. Default value: `false`. **NOTE:** Gateway version >= 1.0.39 above support. 
* `fe_limit` - (Optional, Computed) The maximum write speed of the gateway file share. Unit: `MB/s`, `0` means unlimited. Value range: `0` ~ `1280`. Default value: `0`.
* `gateway_file_share_name` - (Required, ForceNew) The name of the file share. Length from `1` to `255` characters can contain lowercase letters, digits, (.), (_) Or (-), at the same time, must start with a lowercase letter.
* `gateway_id` - (Required, ForceNew) The ID of the gateway.
* `ignore_delete` - (Optional, Computed) The whether to ignore deleted of the gateway file share. After the opening of the Gateway side delete file or delete cloud (OSS) corresponding to the file. Default value: `false`. **NOTE:** `ignore_delete` and `remote_sync` cannot be enabled simultaneously. Gateway version >= 1.0.40 above support. 
* `in_place` - (Optional, Computed, ForceNew) The whether debris optimization of the gateway file share. Default value: `false`.
* `lag_period` - (Optional, Computed) The synchronization delay, I.e. gateway local cache sync to Alibaba Cloud Object Storage Service (oss) of the delay time. Unit: `Seconds`. Value range: `5` ~ `120`. Default value: `5`. **NOTE:** Gateway version >= 1.0.40 above support. 
* `local_path` - (Required, ForceNew) The cache disk inside the device name.
* `nfs_v4_optimization` - (Optional, Computed) The set up gateway file share NFS protocol, whether to enable NFS v4 optimization improve Mount Upload efficiency. Default value: `false`. **NOTE:** If it is enabled, NFS V3 cannot be mounted. The attribute is valid when the attribute `protocol` is `NFS`. Gateway version >= 1.2.0 above support. 
* `oss_bucket_name` - (Required, ForceNew) The name of the OSS Bucket.
* `oss_bucket_ssl` - (Optional, Computed, ForceNew) Whether they are using SSL connect to OSS Bucket.
* `oss_endpoint` - (Required, ForceNew) The gateway file share corresponds to the Object Storage SERVICE (OSS), Bucket Endpoint. **NOTE:** distinguish between intranet and internet Endpoint. We recommend that if the OSS Bucket and the gateway is in the same Region is use the RDS intranet IP Endpoint: `oss-cn-hangzhou-internal.aliyuncs.com`. 
* `partial_sync_paths` - (Optional, ForceNew) In part mode, the directory path group JSON format.
* `path_prefix` - (Optional, ForceNew) The subdirectory path under the object storage (OSS) bucket corresponding to the file share. If it is blank, it means the root directory of the bucket.
* `polling_interval` - (Optional) The reverse synchronization time intervals of the gateway file share. Value range: `15` ~ `36000`. **NOTE:** in copy mode + reverse synchronization is enabled Download file data, value range: `3600` ~ `36000`. 
* `protocol` - (Required, ForceNew) Share types. Valid values: `SMB`, `NFS`.
* `remote_sync` - (Optional, Computed) Whether to enable reverse synchronization of the gateway file share. Default value: `false`.
* `remote_sync_download` - (Optional, Computed) Copy mode, whether to download the file data. Default value: `false`. **NOTE:** only when the attribute `remote_sync` is `true` or acceded to by the speed synchronization group, this parameter will not take effect. 
* `ro_client_list` - (Optional) File sharing NFS read-only client list (IP address or IP address range). Use commas (,) to separate multiple clients.
* `ro_user_list` - (Optional) The read-only client list. When Protocol for Server Message Block (SMB) to go back to.
* `rw_client_list` - (Optional) Read and write the client list. When Protocol NFS is returned when the status is.
* `rw_user_list` - (Optional) Read-write user list. When Protocol for Server Message Block (SMB) to go back to.
* `squash` - (Optional, Computed) The NFS protocol user mapping of the gateway file share. Valid values: `none`, `root_squash`, `all_squash`, `all_anonymous`. Default value: `none`. **NOTE:** The attribute is valid when the attribute `protocol` is `NFS`.
* `support_archive` - (Optional, Computed, ForceNew) Whether to support the archive transparent read.
* `transfer_acceleration` - (Optional) The set up gateway file share whether to enable transmission acceleration needs corresponding OSS Bucket enabled transport acceleration. **NOTE:** Gateway version >= 1.3.0 above support. 
* `windows_acl` - (Optional, Computed) Whether to enable by Windows access list (requires AD domain) the permissions control. Default value: `false`. **NOTE:** The attribute is valid when the attribute `protocol` is `SMB`. Gateway version >= 1.0.45 above support. 
* `bypass_cache_read` - (Optional, Computed) Direct reading OSS of the gateway file share.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Gateway File Share. The value formats as `<gateway_id>:<index_id>`.
* `index_id` - The ID of the file share.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Gateway File Share.
* `update` - (Defaults to 5 mins) Used when update the Gateway File Share.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway File Share.

## Import

Cloud Storage Gateway Gateway File Share can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_storage_gateway_gateway_file_share.example <gateway_id>:<index_id>
```