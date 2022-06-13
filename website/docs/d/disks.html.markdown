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

```
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
* `type` - (Optional) Disk type. Possible values: `system` and `data`.
* `category` - (Optional) Disk category. Possible values: `cloud` (basic cloud disk), `cloud_efficiency` (ultra cloud disk), `ephemeral_ssd` (local SSD cloud disk), `cloud_ssd` (SSD cloud disk), and `cloud_essd` (ESSD cloud disk).
* `encrypted` - (Optional) Indicate whether the disk is encrypted or not. Possible values: `on` and `off`.
* `instance_id` - (Optional) Filter the results by the specified ECS instance ID.
* `resource_group_id` - (Optional, ForceNew, Available in 1.57.0+) The Id of resource group which the disk belongs.
* `tags` - (Optional) A map of tags assigned to the disks. It must be in the format:
  ```
  data "alicloud_disks" "disks_ds" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

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
