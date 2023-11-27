---
subcategory: "DBFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instance"
description: |-
  Provides a Alicloud DBFS Dbfs Instance resource.
---

# alicloud_dbfs_instance

Provides a DBFS Dbfs Instance resource. An instance of a database file system is equivalent to a file system and can store data of file types.

For information about DBFS Dbfs Instance and how to use it, see [What is Dbfs Instance](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_dbfs_instance" "default" {
  category                = "standard"
  zone_id                 = "cn-hangzhou-i"
  size                    = "20"
  performance_level       = "PL1"
  fs_name                 = "rmc-test"
  used_scene              = "MongoDB"
  instance_type           = "dbfs.small"
  raid_stripe_unit_number = "2"
  advanced_features       = "{\"memorySize\":1024,\"pageCacheSize\":128,\"cpuCoreCount\":0.5}"
  kms_key_id              = "00000000-0000-0000-0000-000000000000"
  snapshot_id             = "none"
}
```

## Argument Reference

The following arguments are supported:
* `advanced_features` - (Optional, Computed) The number of CPU cores and the upper limit of memory used by the database file storage instance.
* `category` - (Required, ForceNew) Category of database file system.
* `delete_snapshot` - (Optional) Whether to delete the original snapshot after creating DBFS using the snapshot.
* `enable_raid` - (Optional, ForceNew) Whether to create DBFS in RAID mode. If created in RAID mode, the capacity is at least 66GB.Valid values: true or false. Default value: false.
* `encryption` - (Optional, ForceNew) Whether to encrypt DBFS.Valid values: true or false. Default value: false.
* `fs_name` - (Required) Database file system name.
* `instance_type` - (Optional) Instance type. Value range:
  - dbfs.small
  - dbfs.medium
  - dbfs.large (default)
.
* `kms_key_id` - (Optional, ForceNew) The ID of the KMS key used by DBFS.
* `performance_level` - (Optional, Computed) When you create a DBFS instance, set the performance level of the DBFS instance. Value range:
  - PL0: single disk maximum random read-write IOPS 10000
  - PL1: highest random read-write IOPS 50000 per disk (default)
  - PL2: single disk maximum random read-write IOPS 100000
  - PL3: single disk maximum random read-write IOPS 1 million.
* `raid_stripe_unit_number` - (Optional, ForceNew) Number of strips. Required when the EnableRaid parameter is true.Value range: Currently, only 8 stripes are supported.
* `size` - (Required) Size of database file system, unit GiB.
* `snapshot_id` - (Optional, ForceNew, Computed) The ID of the snapshot used to create the DBFS instance.
* `used_scene` - (Optional) The usage scenario of DBFS. Value range:
  - MySQL 5.7
  - PostgreSQL
  - MongoDB.
* `zone_id` - (Required, ForceNew) The ID of the zone to which the database file system belongs.

The following arguments will be discarded. Please use new fields as soon as possible:
* `instance_name` - (Deprecated since v1.214.0). Field 'instance_name' has been deprecated from provider version 1.214.0. New field 'fs_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dbfs Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Dbfs Instance.
* `update` - (Defaults to 5 mins) Used when update the Dbfs Instance.

## Import

DBFS Dbfs Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbfs_instance.example <id>
```