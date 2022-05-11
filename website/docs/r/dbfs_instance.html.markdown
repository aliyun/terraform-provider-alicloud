---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instance"
sidebar_current: "docs-alicloud-resource-dbfs-instance"
description: |-
  Provides a Alicloud DBFS Instance resource.
---

# alicloud\_dbfs\_instance

Provides a DBFS Instance resource.

For information about DBFS Instance and how to use it, see [What is Instance](https://help.aliyun.com/document_detail/149726.html).

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dbfs_instance" "example" {
  category      = "standard"
  instance_name = "example_value"
  size          = 1
  zone_id       = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, Computed, ForceNew) The type of the Database file system. Valid values: `standard`.
* `delete_snapshot` - (Optional) Whether to delete the original snapshot after the DBFS is created using the snapshot. Valid values : `true` anf `false`.
* `ecs_list` - (Optional) The collection of ECS instances mounted to the Database file system. See the following `Block ecs_list`. **NOTE:** Field 'ecs_list' has been deprecated from provider version 1.156.0 and it will be removed in the future version. Please use the new resource 'alicloud_dbfs_instance_attachment' to attach ECS and DBFS.
* `enable_raid` - (Optional, ForceNew) Whether to create the Database file system in RAID way. Valid values : `true` anf `false`.
* `encryption` - (Optional, ForceNew) Whether to encrypt the database file system. Valid values: `true` and `false`.
* `instance_name` - (Required) The name of the Database file system.
* `kms_key_id` - (Optional, ForceNew) The KMS key ID of the Database file system used. This parameter is valid When `encryption` parameter is set to `true`.
* `performance_level` - (Optional, Computed, ForceNew) The performance level of the Database file system. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
* `raid_stripe_unit_number` - (Optional, ForceNew) The number of strip. This parameter is valid When `enable_raid` parameter is set to `true`.
* `size` - (Required) The size Of the Database file system. Unit: GiB.
* `snapshot_id` - (Optional) The snapshot id of the Database file system.
* `zone_id` - (Required, ForceNew) The Zone ID of the Database file system.
* `tags` - (Optional) A mapping of tags to assign to the resource.


### Block ecs_list

The ecs_list supports the following:

* `ecs_id` - (Optional) The ID of the ECS instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` -The status of Database file system. Valid values: `attached`, `attaching`, `creating`, `deleted`, `deleting`, `detaching`, `resizing`, `snapshotting`, `unattached`, `upgrading`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.
* `update` - (Defaults to 15 mins) Used when update the Instance.
* `delete` - (Defaults to 10 mins) Used when delete the Instance.

## Import

DBFS Instance can be imported using the id, e.g.

```
$ terraform import alicloud_dbfs_instance.example <id>
```
