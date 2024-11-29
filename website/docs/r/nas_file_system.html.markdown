---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_file_system"
sidebar_current: "docs-alicloud-resource-nas-file-system"
description: |-
  Provides a Alicloud File Storage (NAS) File System resource.
---

# alicloud_nas_file_system

Provides a File Storage (NAS) File System resource.

For information about File Storage (NAS) File System and how to use it, see [What is File System](https://www.alibabacloud.com/help/en/nas/developer-reference/api-nas-2017-06-26-createfilesystem).

-> **NOTE:** Available since v1.33.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_file_system&exampleId=f5f5e208-95b8-0c3e-c75b-1e859e5cc773ab263762&activeTab=example&spm=docs.r.nas_file_system.0.f5f5e20895&intl_lang=EN_US" target="_blank">
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_file_system&exampleId=484fb7e1-829c-f8c1-e36c-d2b66549afe95c24478d&activeTab=example&spm=docs.r.nas_file_system.1.484fb7e182&intl_lang=EN_US" target="_blank">
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
  file_system_type = "extreme"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "NFS"
  storage_type     = "standard"
  capacity         = 100
  description      = var.name
  encrypt_type     = 1
  file_system_type = "extreme"
  zone_id          = data.alicloud_nas_zones.default.zones.0.zone_id
}
```

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_file_system&exampleId=ff389e9e-bbb3-6c5d-47f6-7f9bfd68c6144140711d&activeTab=example&spm=docs.r.nas_file_system.2.ff389e9ebb&intl_lang=EN_US" target="_blank">
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
  file_system_type = "cpfs"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nas_zones.default.zones.1.zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "cpfs"
  storage_type     = "advance_100"
  capacity         = 5000
  description      = var.name
  file_system_type = "cpfs"
  vswitch_id       = alicloud_vswitch.default.id
  vpc_id           = alicloud_vpc.default.id
  zone_id          = data.alicloud_nas_zones.default.zones.1.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `capacity` - (Optional, Int, Available since v1.140.0) The capacity of the file system. Unit: GiB. **Note:** If `file_system_type` is set to `extreme` or `cpfs`, `capacity` must be set.
* `description` - (Optional) The description of the file system.
* `encrypt_type` - (Optional, ForceNew, Int, Available since v1.121.2) Specifies whether to encrypt data in the file system. Default value: `0`. Valid values:
  - `0`: The data in the file system is not encrypted.
  - `1`: A NAS-managed key is used to encrypt the data in the file system. **NOTE:** `1` is valid only when `file_system_type` is set to `standard` or `extreme`.
  - `2`: A KMS-managed key is used to encrypt the data in the file system. **NOTE:** From version 1.140.0, `encrypt_type` can be set to `2`, and `2` is valid only when `file_system_type` is set to `standard` or `extreme`.
* `file_system_type` - (Optional, ForceNew, Available since v1.140.0) The type of the file system. Default value: `standard`. Valid values: `standard`, `extreme`, `cpfs`.
* `kms_key_id` - (Optional, ForceNew, Available since v1.140.0) The ID of the KMS-managed key. **Note:** If `encrypt_type` is set to `2`, `kms_key_id` must be set.
* `protocol_type` - (Required, ForceNew) The protocol type of the file system. Valid values:
  - If `file_system_type` is set to `standard`. Valid values: `NFS`, `SMB`.
  - If `file_system_type` is set to `extreme`. Valid values: `NFS`.
  - If `file_system_type` is set to `cpfs`. Valid values: `cpfs`.
* `recycle_bin` - (Optional, Set, Available since v1.236.0) The recycle bin feature of the file system. See [`recycle_bin`](#recycle_bin) below.
-> **NOTE:** `recycle_bin` takes effect only if `file_system_type` is set to `standard`.
* `nfs_acl` - (Optional, Set, Available since v1.236.0) The NFS ACL feature of the file system. See [`nfs_acl`](#nfs_acl) below.
-> **NOTE:** `nfs_acl` takes effect only if `file_system_type` is set to `standard`.
* `resource_group_id` - (Optional, Available since v1.236.0) The ID of the resource group.
* `snapshot_id` - (Optional, Available since v1.236.0) The ID of the snapshot. **NOTE:** `snapshot_id` takes effect only if `file_system_type` is set to `extreme`.
* `storage_type` - (Required, ForceNew) The storage type of the file system. Valid values:
  - If `file_system_type` is set to `standard`. Valid values: `Performance`, `Capacity`, `Premium`.
  - If `file_system_type` is set to `extreme`. Valid values: `standard`, `advance`.
  - If `file_system_type` is set to `cpfs`. Valid values: `advance_100`, `advance_200`.
-> **NOTE:** From version 1.140.0, `storage_type` can be set to `standard`, `advance`. From version 1.153.0, `storage_type` can be set to `advance_100`, `advance_200`. From version 1.236.0, `storage_type` can be set to `Premium`.
* `tags` - (Optional, Available since v1.153.0) A mapping of tags to assign to the resource.
* `vswitch_id` - (Optional, ForceNew, Available since v1.153.0) The ID of the vSwitch. **NOTE:** `vswitch_id` takes effect only if `file_system_type` is set to `cpfs`.
* `vpc_id` - (Optional, ForceNew, Available since v1.153.0) The ID of the VPC. **NOTE:** `vpc_id` takes effect only if `file_system_type` is set to `cpfs`.
* `zone_id` - (Optional, ForceNew, Available since v1.140.0) The ID of the zone. **Note:** If `file_system_type` is set to `extreme` or `cpfs`, `zone_id` must be set.

### `recycle_bin`

The recycle_bin supports the following:

* `status` - (Optional) Specifies whether to enable the recycle bin feature. Default value: `Disable`. Valid values: `Enable`, `Disable`.
* `reserved_days` - (Optional) The retention period of the files in the recycle bin. Unit: days. Default value: `3`. Valid values: `1` to `180`. **NOTE:** `reserved_days` takes effect only if `status` is set to `Enable`.

### `nfs_acl`

The nfs_acl supports the following:

* `enabled` - (Optional, Bool) Specifies whether to enable the NFS ACL feature. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of File System.
* `create_time` - (Available since v1.236.0) The time when the file system was created.
* `status` - (Available since v1.236.0) The status of the File System.
* `recycle_bin` - (Available since v1.236.0) The recycle bin feature of the file system.
  * `size` - The size of the files that are dumped to the recycle bin.
  * `secondary_size` - The size of the Infrequent Access (IA) data that is dumped to the recycle bin.
  * `enable_time` - The time at which the recycle bin was enabled.
  
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the File System.
* `update` - (Defaults to 5 mins) Used when update the File System.
* `delete` - (Defaults to 20 mins) Used when delete the File System.

## Import

File Storage (NAS) File System can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_file_system.example <id>
```
