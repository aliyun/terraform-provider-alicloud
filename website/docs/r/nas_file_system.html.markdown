---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_file_system"
description: |-
  Provides a Alicloud File Storage (NAS) File System resource.
---

# alicloud_nas_file_system

Provides a File Storage (NAS) File System resource.

File System Instance.

For information about File Storage (NAS) File System and how to use it, see [What is File System](https://www.alibabacloud.com/help/en/nas/developer-reference/api-nas-2017-06-26-createfilesystem).

-> **NOTE:** Available since v1.33.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "standard"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "NFS"
  storage_type     = "Capacity"
  description      = var.name
  encrypt_type     = 1
  file_system_type = "standard"
  recycle_bin {
    status        = "Enable"
    reserved_days = "10"
  }
  nfs_acl {
    enabled = true
  }
  zone_id = data.alicloud_nas_zones.default.zones.0.zone_id
}
```

## Argument Reference

The following arguments are supported:
* `capacity` - (Optional, Computed, Int, Available since v1.140.0) File system capacity.

Unit: GiB, required and valid when FileSystemType = extreme or cpfs.

For optional values, please refer to the actual specifications on the purchase page:
  -[Fast NAS Pay-As-You-Go Page](https://common-buy.aliyun.com/? commodityCode=nas_extreme_post#/buy)
  -[Fast NAS Package Monthly Purchase Page](https://common-buy.aliyun.com/? commodityCode=nas_extreme#/buy)
  -[Parallel File System CPFS Pay-As-You-Go Purchase Page](https://common-buy.aliyun.com/? commodityCode=nas_cpfs_post#/buy)
  -[Parallel File System CPFS Package Monthly Purchase Page](https://common-buy.aliyun.com/? commodityCode=cpfs#/buy)
* `description` - (Optional) File system description.

Restrictions:
  - 2~128 English or Chinese characters in length.
  - Must start with upper and lower case letters or Chinese, and cannot start with'http: // 'and'https.
  - Can contain numbers, colons (:), underscores (_), or dashes (-).
* `encrypt_type` - (Optional, ForceNew, Computed, Int, Available since v1.121.2) Whether the file system is encrypted.

Use the KMS service hosting key to encrypt and store the file system disk data. When reading and writing encrypted data, there is no need to decrypt it.

Value:
  - 0 (default): not encrypted.
  - 1:NAS managed key. NAS managed keys are supported when FileSystemType = standard or extreme.
  - 2: User management key. You can manage keys only when FileSystemType = extreme.
* `file_system_type` - (Optional, ForceNew, Computed, Available since v1.140.0) File system type.

Value:
  - standard (default): Universal NAS
  - extreme: extreme NAS
  - cpfs: file storage CPFS
* `keytab` - (Optional, Available since v1.248.0) String of keytab file content encrypted by base64
* `keytab_md5` - (Optional, Available since v1.248.0) String of the keytab file content encrypted by MD5
* `kms_key_id` - (Optional, ForceNew, Computed, Available since v1.140.0) The ID of the KMS key.
This parameter is required only when EncryptType = 2.
* `nfs_acl` - (Optional, Computed, List, Available since v1.236.0) NFS ACL See [`nfs_acl`](#nfs_acl) below.
* `options` - (Optional, Computed, List, Available since v1.248.0) Option. See [`options`](#options) below.
* `protocol_type` - (Required, ForceNew) File transfer protocol type.
  - When FileSystemType = standard, the values are NFS and SMB.
  - When FileSystemType = extreme, the value is NFS.
  - When FileSystemType = cpfs, the value is cpfs.
* `recycle_bin` - (Optional, Computed, List) Recycle Bin See [`recycle_bin`](#recycle_bin) below.
* `resource_group_id` - (Optional, Computed, Available since v1.236.0) The ID of the resource group.
* `smb_acl` - (Optional, Computed, List, Available since v1.248.0) SMB ACL See [`smb_acl`](#smb_acl) below.
* `snapshot_id` - (Optional, Available since v1.236.0) Only extreme NAS is supported.

-> **NOTE:** A file system is created from a snapshot. The version of the created file system is the same as that of the snapshot source file system. For example, if the source file system version of the snapshot is 1 and you need to create A file system of version 2, you can first create A file system A from the snapshot, then create A file system B that meets the configuration of version 2, copy the data in file system A to file system B, and migrate the business to file system B after the copy is completed.

* `storage_type` - (Required, ForceNew) The storage type.
  - When FileSystemType = standard, the values are Performance, Capacity, and Premium.
  - When FileSystemType = extreme, the value is standard or advance.
  - When FileSystemType = cpfs, the values are advance_100(100MB/s/TiB baseline) and advance_200(200MB/s/TiB baseline).
* `tags` - (Optional, Map, Available since v1.153.0) Label information collection.
* `vswitch_id` - (Optional, ForceNew, Available since v1.153.0) The ID of the switch.
This parameter must be configured when FileSystemType = cpfs.
When the FileSystemType is standard or extreme, this parameter is reserved for the interface and has not taken effect yet. You do not need to configure it.
* `vpc_id` - (Optional, ForceNew, Available since v1.153.0) The ID of the VPC network.
This parameter must be configured when FileSystemType = cpfs.
When the FileSystemType is standard or extreme, this parameter is reserved for the interface and has not taken effect yet. You do not need to configure it.
* `zone_id` - (Optional, ForceNew, Computed) The zone ID.

The usable area refers to the physical area where power and network are independent of each other in the same area.

When the FileSystemType is set to standard, this parameter is optional. By default, a zone that meets the conditions is randomly selected based on the ProtocolType and StorageType configurations. This parameter is required when FileSystemType = extreme or FileSystemType = cpfs.

-> **NOTE:** - file systems in different zones in the same region communicate with ECS cloud servers.

-> **NOTE:** - We recommend that the file system and the ECS instance belong to the same zone to avoid cross-zone latency.


### `nfs_acl`

The nfs_acl supports the following:
* `enabled` - (Optional, Computed) Whether the NFS ACL function is enabled.

### `options`

The options supports the following:
* `enable_oplock` - (Optional, Computed, Available since v1.248.0) Whether to enable the OpLock function. Value:
  - true: On.
  - false: does not turn on.

-> **NOTE:**  Description Only file systems of the SMB protocol type are supported.


### `recycle_bin`

The recycle_bin supports the following:
* `reserved_days` - (Optional, Computed, Int) Retention time of files in the Recycle Bin. Unit: days.
* `status` - (Optional, Computed) Recycle Bin Status

### `smb_acl`

The smb_acl supports the following:
* `enable_anonymous_access` - (Optional, Computed, Available since v1.248.0) Whether to allow anonymous access.
  - true: Allow anonymous access.
  - false (default): Anonymous access is not allowed.
* `enabled` - (Optional, Computed, Available since v1.248.0) Whether SMB ACL is enabled
* `encrypt_data` - (Optional, Available since v1.248.0) Whether transmission encryption is enabled.
  - true: Enables encryption in transit.
  - false (default): Transport encryption is not enabled.
* `home_dir_path` - (Optional, Available since v1.248.0) The user directory home path for each user. The file path format is as follows:
  - A forward slash (/) or backslash (\) as a separator.
  - Each paragraph cannot contain ":|? *.
  - The length of each segment ranges from 0 to 255.
  - The total length range is 0~32767.

For example, if the user directory is/home, the file system will automatically create A directory of/home/A when user A logs in. Skip if/home/A already exists.

-> **NOTE:**  Explain that user A needs to have the permission to create A directory, otherwise the/home/A directory cannot be created.

* `reject_unencrypted_access` - (Optional, Available since v1.248.0) 
Whether to reject non-encrypted clients.
  - true: Deny non-encrypted clients.
  - false (default): Non-encrypted clients are not rejected.
* `super_admin_sid` - (Optional, Available since v1.248.0) The ID of the Super User. The ID rules are as follows:
  - Must start with S and no other letters can appear after the S at the beginning.
  - At least three dashes (-) apart.

Such as S-1-5-22 or S-1-5-22-23.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - CreateTime
* `recycle_bin` - Recycle Bin
  * `enable_time` - Recycle Bin open time
  * `secondary_size` - Amount of low-frequency data stored in the recycle bin. Unit: Byte.
  * `size` - The amount of files stored in the Recycle Bin. Unit: Byte.
* `region_id` - RegionId
* `status` - File system status. Includes:(such as creating a mount point) can only be performed when the file system is in the Running state.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the File System.
* `delete` - (Defaults to 20 mins) Used when delete the File System.
* `update` - (Defaults to 10 mins) Used when update the File System.

## Import

File Storage (NAS) File System can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_file_system.example <id>
```