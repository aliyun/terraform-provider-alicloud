---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instance"
description: |-
  Provides a Alicloud DBFS Dbfs Instance resource.
---

# alicloud_dbfs_instance

Provides a DBFS Dbfs Instance resource. An instance of a database file system is equivalent to a file system and can store data of file types.

For information about DBFS Dbfs Instance and how to use it, see [What is Dbfs Instance](https://next.api.alibabacloud.com/document/DBFS/2020-04-18/CreateDbfs).

-> **NOTE:** Need to contact us open whitelist before you can use the resource.

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dbfs_instance&exampleId=b2f245ec-6bb0-7762-d74f-982d057ac847b9f79dca&activeTab=example&spm=docs.r.dbfs_instance.0.b2f245ec6b&intl_lang=EN_US" target="_blank">
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

resource "alicloud_dbfs_instance" "example" {
  category          = "standard"
  zone_id           = "cn-hangzhou-i"
  performance_level = "PL1"
  fs_name           = var.name
  size              = 100
}
```

## Argument Reference

The following arguments are supported:
* `advanced_features` - (Optional, Computed, Available since v1.212.0) The number of CPU cores and the upper limit of memory used by the database file storage instance.
* `category` - (Required, ForceNew) Category of database file system.
* `delete_snapshot` - (Optional) Whether to delete the original snapshot after creating DBFS using the snapshot.
* `enable_raid` - (Optional, ForceNew) Whether to create DBFS in RAID mode. If created in RAID mode, the capacity is at least 66GB.Valid values: true or false. Default value: false.
* `encryption` - (Optional, ForceNew) Whether to encrypt DBFS.Valid values: true or false. Default value: false.
* `fs_name` - (Optional, Available since v1.212.0) Database file system name.
* `instance_type` - (Optional, Available since v1.212.0) Instance type. Value range:
  - dbfs.small
  - dbfs.medium
  - dbfs.large (default)
* `kms_key_id` - (Optional, ForceNew) The ID of the KMS key used by DBFS.
* `performance_level` - (Optional, Computed) When you create a DBFS instance, set the performance level of the DBFS instance. Value range:
  - PL0: single disk maximum random read-write IOPS 10000
  - PL1: highest random read-write IOPS 50000 per disk (default)
  - PL2: single disk maximum random read-write IOPS 100000
  - PL3: single disk maximum random read-write IOPS 1 million.
* `raid_stripe_unit_number` - (Optional, ForceNew) Number of strips. Required when the EnableRaid parameter is true.Value range: Currently, only 8 stripes are supported.
* `size` - (Required) Size of database file system, unit GiB.
* `snapshot_id` - (Optional, ForceNew, Computed) The ID of the snapshot used to create the DBFS instance.
* `used_scene` - (Optional, Available since v1.212.0) The usage scenario of DBFS. Value range:
  - MySQL 5.7
  - PostgreSQL
  - MongoDB.
* `zone_id` - (Required, ForceNew) The ID of the zone to which the database file system belongs.
* `ecs_list` - (Optional, Deprecated from v1.156.0) The collection of ECS instances mounted to the Database file system. See [`ecs_list`](#ecs_list) below.  **NOTE:** Field 'ecs_list' has been deprecated from provider version 1.156.0 and it will be removed in the future version. Please use the new resource 'alicloud_dbfs_instance_attachment' to attach ECS and DBFS. See [`ecs_list`](#ecs_list) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

The following arguments will be discarded. Please use new fields as soon as possible:
* `instance_name` - (Deprecated since v1.212.0). Field 'instance_name' has been deprecated from provider version 1.212.0. New field 'fs_name' instead.

### `ecs_list`

The ecs_list supports the following:
* `ecs_id` - (Optional) The ID of the ECS instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dbfs Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Dbfs Instance.
* `update` - (Defaults to 5 mins) Used when update the Dbfs Instance.

## Import

DBFS Dbfs Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbfs_instance.example <id>
```