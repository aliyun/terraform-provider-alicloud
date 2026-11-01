---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_disks"
sidebar_current: "docs-alicloud-datasource-disks"
description: |-
    Provides a list of disks to the user.
---

# alicloud\_disks

-> **DEPRECATED:** This datasource has been renamed to [alicloud_ecs_disks](https://www.terraform.io/docs/providers/alicloud/d/ecs_disks) from version 1.122.0.

This data source provides the disks of the current Alibaba Cloud user.

## Example Usage

```terraform
data "alicloud_disks" "disks_ds" {
  name_regex = "sample_disk"
}

output "first_disk_id" {
  value = "${data.alicloud_disks.disks_ds.disks.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of disks IDs.
* `name_regex` - (Optional) A regex string to filter results by disk name.
* `type` - (Optional, Deprecated) Disk type. Possible values: `system` and `data`.
* `category` - (Optional) Disk category. Possible values: `cloud` (basic cloud disk), `cloud_efficiency` (ultra cloud disk), `ephemeral_ssd` (local SSD cloud disk), `cloud_ssd` (SSD cloud disk), and `cloud_essd` (ESSD cloud disk), `cloud_essd_entry`.
* `encrypted` - (Optional) Indicate whether the disk is encrypted or not. Possible values: `on` and `off`.
* `instance_id` - (Optional) Filter the results by the specified ECS instance ID.
* `resource_group_id` - (Optional, Available since v1.57.0) The Id of resource group which the disk belongs.
* `tags` - (Optional) A map of tags assigned to the disks. It must be in the format:

  ```terraform
  data "alicloud_disks" "disks_ds" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `additional_attributes` - (Optional) Additional attributes.
* `auto_snapshot_policy_id` - (Optional) The ID of the automatic snapshot policy.
* `delete_auto_snapshot` - (Optional) Whether to delete the automatic snapshot.
* `delete_with_instance` - (Optional) Whether to delete the disk with the instance.
* `disk_name` - (Optional) The name of the disk.
* `disk_type` - (Optional) Disk type. Valid values: `system`, `data`, `all`.
* `dry_run` - (Optional) Whether to check validity without making the request.
* `enable_auto_snapshot` - (Optional) Whether to enable the automatic snapshot.
* `enable_automated_snapshot_policy` - (Optional) Whether to enable the automated snapshot policy.
* `enable_shared` - (Optional) Whether to enable shared.
* `kms_key_id` - (Optional) The KMS key ID.
* `operation_locks` - (Optional) Operation locks. See [`operation_locks`](#operation_locks) below.
* `payment_type` - (Optional) The payment type. Valid values: `PayAsYouGo`, `Subscription`.
* `portable` - (Optional) Whether the disk is portable.
* `snapshot_id` - (Optional) The snapshot ID.
* `status` - (Optional) The status of the disk.
* `zone_id` - (Optional) The zone ID.
* `availability_zone` - (Optional, Deprecated) The availability zone. Deprecated; use `zone_id` instead.

### `operation_locks`

The `operation_locks` supports the following:

* `lock_reason` - (Optional) The reason of the operation lock.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of disk names.
* `disks` - A list of disks. Each element contains the following attributes:
  * `id` - ID of the disk.
  * `name` - Disk name.
  * `description` - Disk description.
  * `region_id` - Region ID the disk belongs to.
  * `availability_zone` - Availability zone of the disk.
  * `status` - Current status. Possible values: `In_use`, `Available`, `Attaching`, `Detaching`, `Creating` and `ReIniting`.
  * `type` - Disk type. Possible values: `system` and `data`.
  * `category` - Disk category. Possible values: `cloud` (basic cloud disk), `cloud_efficiency` (ultra cloud disk), `ephemeral_ssd` (local SSD cloud disk), `cloud_ssd` (SSD cloud disk), and `cloud_essd` (ESSD cloud disk).
  * `encrypted` - Indicate whether the disk is encrypted or not. Possible values: `on` and `off`.
  * `size` - Disk size in GiB.
  * `image_id` - ID of the image from which the disk is created. It is null unless the disk is created using an image.
  * `snapshot_id` - Snapshot used to create the disk. It is null if no snapshot is used to create the disk.
  * `instance_id` - ID of the related instance. It is `null` unless the `status` is `In_use`.
  * `creation_time` - Disk creation time.
  * `attached_time` - Disk attachment time.
  * `detached_time` - Disk detachment time.
  * `expiration_time` - Disk expiration time.
  * `tags` - A map of tags assigned to the disk.
  * `resource_group_id` - The Id of resource group.
  * `auto_snapshot_policy_id` - The ID of the automatic snapshot policy.
  * `delete_auto_snapshot` - Whether to delete the automatic snapshot.
  * `delete_with_instance` - Whether to delete the disk with the instance.
  * `device` - The device name.
  * `disk_id` - The ID of the disk.
  * `disk_name` - The name of the disk.
  * `disk_type` - The disk type.
  * `enable_auto_snapshot` - Whether to enable the automatic snapshot.
  * `enable_automated_snapshot_policy` - Whether to enable the automated snapshot policy.
  * `expired_time` - The expired time of the disk.
  * `iops` - The IOPS of the disk.
  * `iops_read` - The read IOPS.
  * `iops_write` - The write IOPS.
  * `kms_key_id` - The KMS key ID.
  * `mount_instance_num` - The number of mounted instances.
  * `mount_instances` - A list of mounted instances.
    * `attached_time` - The attached time.
    * `device` - The device name.
    * `instance_id` - The ID of the instance.
  * `operation_locks` - Operation locks.
    * `lock_reason` - The reason of the operation lock.
  * `payment_type` - The payment type.
  * `performance_level` - The performance level of the ESSD.
  * `portable` - Whether the disk is portable.
  * `product_code` - The product code.
  * `zone_id` - The zone ID.
