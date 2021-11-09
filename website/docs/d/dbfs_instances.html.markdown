---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instances"
sidebar_current: "docs-alicloud-datasource-dbfs-instances"
description: |-
  Provides a list of DBFS Instances to the user.
---

# alicloud\_dbfs\_instances

This data source provides the DBFS Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dbfs_instances" "ids" {
  ids = ["example_id"]
}
output "dbfs_instance_id_1" {
  value = data.alicloud_dbfs_instances.ids.instances.0.id
}

data "alicloud_dbfs_instances" "nameRegex" {
  name_regex = "^my-Instance"
}
output "dbfs_instance_id_2" {
  value = data.alicloud_dbfs_instances.nameRegex.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) Database file system state. Valid values: `attached`, `attaching`, `creating`, `deleted`, `deleting`, `detaching`, `resizing`, `snapshotting`, `unattached`, `upgrading`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of DBFS Instances. Each element contains the following attributes:
    * `attach_node_number` - the number of nodes of the Database file system.
    * `category` -  The type of the Database file system. Valid values: `standard`.
    * `create_time` - The create time of the Database file system.
    * `dbfs_cluster_id` - The cluster ID of the Database file system.
    * `ecs_list` - The collection of ECS instances mounted to the Database file system.
        * `ecs_id` - The ID of the ECS instance.
    * `enable_raid` - Whether to create the Database file system in RAID way. Valid values : `true` anf `false`. 
    * `encryption` - Whether to encrypt the Database file system. Valid values: `true` and `false`.
    * `id` - The ID of the Instance.
    * `instance_id` -  The ID of the Database File System
    * `instance_name` - The name of the Database file system.
    * `kms_key_id` - The KMS key ID of the Database file system used. This parameter is valid When `encryption` parameter is set to `true`.
    * `payment_type` - Thr payment type of the Database file system. Valid value: `PayAsYouGo`.
    * `performance_level` - The performance level of the Database file system. Valid values: `PL0`, `PL1`, `PL2`, `PL3`.
    * `raid_stripe_unit_number` - The number of strip . When `enable_raid` parameter is set to `true` will transfer. This parameter is valid When `enable_raid` parameter is set to `true`.
    * `size` - The size Of the Database file system. Unit: GiB.
    * `status` - The status of the Database file system.
    * `zone_id` - The Zone ID of the Database file system.
