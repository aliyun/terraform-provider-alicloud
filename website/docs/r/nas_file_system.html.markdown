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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_file_system&exampleId=e4c2e56e-d616-5a3d-e56d-56aa3caacf29cd5c0a18&activeTab=example&spm=docs.r.nas_file_system.0.e4c2e56ed6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

CPFS Usage

A CPFS file system is created inside a VPC, so `vpc_id` and `vswitch_id` are required. `capacity` and `zone_id` are also required for CPFS.

```terraform
resource "alicloud_nas_file_system" "cpfs" {
  protocol_type    = "cpfs"
  storage_type     = "advance_100"
  file_system_type = "cpfs"
  capacity         = 3600
  zone_id          = "cn-hangzhou-i"
  # vpc_id and vswitch_id are required when file_system_type = cpfs
  vpc_id     = "vpc-xxxxxxxxxxxxxxxxxxxxx"
  vswitch_id = "vsw-xxxxxxxxxxxxxxxxxxxxx"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nas_file_system&spm=docs.r.nas_file_system.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `capacity` - (Optional, Computed, Int, Available since v1.140.0) File system capacity.

  Unit: GiB, required and valid when FileSystemType = extreme or cpfs.

  For optional values, please refer to the actual specifications on the purchase page:
    - [Fast NAS Pay-As-You-Go Page](https://common-buy.aliyun.com/?commodityCode=nas_extreme_post#/buy)
    - [Fast NAS Package Monthly Purchase Page](https://common-buy.aliyun.com/?commodityCode=nas_extreme#/buy)
    - [Parallel File System CPFS Pay-As-You-Go Purchase Page](https://common-buy.aliyun.com/?commodityCode=nas_cpfs_post#/buy)
    - [Parallel File System CPFS Package Monthly Purchase Page](https://common-buy.aliyun.com/?commodityCode=cpfs#/buy)

* `description` - (Optional) File system description.

  Restrictions:
    - 2~128 English or Chinese characters in length.
    - Must start with upper and lower case letters or Chinese, and cannot start with 'http://' and 'https'.
    - Can contain numbers, colons (:), underscores (_), or dashes (-).

* `encrypt_type` - (Optional, ForceNew, Computed, Int, Available since v1.121.2) Whether the file system is encrypted.

  Use the KMS service hosting key to encrypt and store the file system disk data. When reading and writing encrypted data, there is no need to decrypt it.

  Value:
    - 0 (default): not encrypted.
    - 1: NAS managed key. NAS managed keys are supported when FileSystemType = standard or extreme.
    - 2: User management key. You can manage keys only when FileSystemType = extreme.

* `file_system_type` - (Optional, ForceNew, Computed, Available since v1.140.0) File system type.

  Value:
    - standard (default): Universal NAS
    - extreme: extreme NAS
    - cpfs: file storage CPFS
    - cpfsse: file storage CPFS Smart Edition

  -> **NOTE:** Whether the network fields `vpc_id` and `vswitch_id` must be configured depends on `file_system_type`. Only CPFS file systems create a resource inside a VPC; for `standard` and `extreme` these fields are reserved by the interface and have not taken effect, so they should be left unset. The configuration rule for each `file_system_type` is as follows:
    - `standard` / `extreme`: do not configure `vpc_id` or `vswitch_id` (reserved by the interface, not effective).
    - `cpfs`: both `vpc_id` and `vswitch_id` are required.
    - `cpfsse`: `vpc_id` is required; `vswitch_id` is not required.

* `keytab` - (Optional, Available since v1.248.0) String of keytab file content encrypted by base64

  -> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `keytab_md5` - (Optional, Available since v1.248.0) String of the keytab file content encrypted by MD5

  -> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `kms_key_id` - (Optional, ForceNew, Computed, Available since v1.140.0) The ID of the KMS key.

  This parameter is required only when EncryptType = 2.

* `nfs_acl` - (Optional, Computed, List, Available since v1.236.0) NFS ACL See [`nfs_acl`](#nfs_acl) below.

* `options` - (Optional, Computed, List, Available since v1.248.0) Option. See [`options`](#options) below.

* `protocol_type` - (Required, ForceNew) File transfer protocol type.
    - When FileSystemType = standard, the values are NFS and SMB.
    - When FileSystemType = extreme, the value is NFS.
    - When FileSystemType = cpfs, the value is cpfs.

* `recycle_bin` - (Optional, Computed, List) Recycle Bin See [`recycle_bin`](#recycle_bin) below.

* `redundancy_type` - (Optional, ForceNew, Computed, Available since v1.267.0) Storage redundancy type. Only effective for General CPFS. Options: Locally Redundant Storage (LRS), Zone-Redundant Storage (ZRS). Default value: LRS.

* `redundancy_vswitch_ids` - (Optional, ForceNew, List, Available since v1.267.0) Redundancy vSwitch ID list. Only set when the file system's storage redundancy type is Zone-Redundant Storage (ZRS), and must set vSwitch IDs from three different availability zones under the same VPC.

* `resource_group_id` - (Optional, Computed, Available since v1.236.0) The ID of the resource group.

* `smb_acl` - (Optional, Computed, List, Available since v1.248.0) SMB ACL See [`smb_acl`](#smb_acl) below.

* `snapshot_id` - (Optional, Available since v1.236.0) Only extreme NAS is supported.

  -> **NOTE:** A file system is created from a snapshot. The version of the created file system is the same as that of the snapshot source file system. For example, if the source file system version of the snapshot is 1 and you need to create A file system of version 2, you can first create A file system A from the snapshot, then create A file system B that meets the configuration of version 2, copy the data in file system A to file system B, and migrate the business to file system B after the copy is completed.

  -> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `storage_type` - (Required, ForceNew) The storage type.
    - When FileSystemType = standard, the values are Performance, Capacity, and Premium.
    - When FileSystemType = extreme, the value is standard or advance.
    - When FileSystemType = cpfs, the values are advance_100(100MB/s/TiB baseline) and advance_200(200MB/s/TiB baseline).

* `tags` - (Optional, Map, Available since v1.153.0) Label information collection.

* `vswitch_id` - (Optional, ForceNew, Available since v1.153.0) The ID of the switch.

  This parameter must be configured when FileSystemType = cpfs. When the FileSystemType is standard or extreme, this parameter is reserved for the interface and has not taken effect yet. You do not need to configure it.

  -> **NOTE:** This `vswitch_id` configures the network of the CPFS file system itself and is different from `vswitch_id` in [`alicloud_nas_mount_target`](nas_mount_target.html), which specifies the vSwitch of the mount target used by clients to access a NAS file system. A mount target still requires its own `vswitch_id` regardless of `file_system_type`.

  -> **NOTE:** For `standard` or `extreme` file systems, do not set `vswitch_id`. Since this field is not `Computed`, a value configured on these file system types cannot be read back from the API, which produces a permanent diff on every plan and, because the field is `ForceNew`, forces the file system to be destroyed and recreated. If you previously configured `vswitch_id` on a `standard` or `extreme` file system, remove it from the configuration before upgrading.

* `vpc_id` - (Optional, ForceNew, Available since v1.153.0) The ID of the VPC network.

  This parameter must be configured when FileSystemType = cpfs. When the FileSystemType is standard or extreme, this parameter is reserved for the interface and has not taken effect yet. You do not need to configure it.

  -> **NOTE:** For `standard` or `extreme` file systems, do not set `vpc_id`. Since this field is not `Computed`, a value configured on these file system types cannot be read back from the API, which produces a permanent diff on every plan and, because the field is `ForceNew`, forces the file system to be destroyed and recreated. If you previously configured `vpc_id` on a `standard` or `extreme` file system, remove it from the configuration before upgrading.

* `zone_id` - (Optional, ForceNew, Computed) The zone ID.

  The usable area refers to the physical area where power and network are independent of each other in the same area.

  When the FileSystemType is set to standard, this parameter is optional. By default, a zone that meets the conditions is randomly selected based on the ProtocolType and StorageType configurations. This parameter is required when FileSystemType = extreme or FileSystemType = cpfs.

  -> **NOTE:** file systems in different zones in the same region communicate with ECS cloud servers.

  -> **NOTE:** We recommend that the file system and the ECS instance belong to the same zone to avoid cross-zone latency.

### `nfs_acl`

The nfs_acl supports the following:
* `enabled` - (Optional, Computed) Whether the NFS ACL function is enabled.

### `options`

The options supports the following:
* `enable_oplock` - (Optional, Computed, Available since v1.248.0) Whether to enable the OpLock function. Value:
    - true: On.
    - false: does not turn on.

  -> **NOTE:** Description Only file systems of the SMB protocol type are supported.

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

  -> **NOTE:** Explain that user A needs to have the permission to create A directory, otherwise the/home/A directory cannot be created.

* `reject_unencrypted_access` - (Optional, Available since v1.248.0) Whether to reject non-encrypted clients.
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the File System.
* `delete` - (Defaults to 20 mins) Used when delete the File System.
* `update` - (Defaults to 10 mins) Used when update the File System.

## Import

File Storage (NAS) File System can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_file_system.example <id>
```