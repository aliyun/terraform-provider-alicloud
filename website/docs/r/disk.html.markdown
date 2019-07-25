---
layout: "alicloud"
page_title: "Alicloud: alicloud_disk"
sidebar_current: "docs-alicloud-resource-disk"
description: |-
  Provides a ECS Disk resource.
---

# alicloud\_disk

Provides a ECS disk resource.

-> **NOTE:** One of `size` or `snapshot_id` is required when specifying an ECS disk. If all of them be specified, `size` must more than the size of snapshot which `snapshot_id` represents. Currently, `alicloud_disk` doesn't resize disk.

## Example Usage

```
# Create a new ECS disk.
resource "alicloud_disk" "ecs_disk" {
  # cn-beijing
  availability_zone = "cn-beijing-b"
  name              = "New-disk"
  description       = "Hello ecs disk."
  category          = "cloud_efficiency"
  size              = "30"

  tags = {
    Name = "TerraformTest"
  }
}
```
## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The Zone to create the disk in.
* `name` - (Optional) Name of the ECS disk. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `description` - (Optional) Description of the disk. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `category` - (Optional, ForceNew) Category of the disk. Valid values are `cloud`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`. Default is `cloud_efficiency`.
* `size` - (Required) The size of the disk in GiBs. When resize the disk, the new size must be greater than the former value, or you would get an error `InvalidDiskSize.TooSmall`.
* `snapshot_id` - (Optional) A snapshot to base the disk off of. If the disk size required by snapshot is greater than `size`, the `size` will be ignored, conflict with `encrypted`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `encrypted` - (Optional) If true, the disk will be encrypted, conflict with `snapshot_id`.
* `delete_auto_snapshot` - (Optional Available in 1.53.0+) Indicates whether the automatic snapshot is deleted when the disk is released. Default value: false.
* `delete_with_instance` - (Optional Available in 1.53.0+) Indicates whether the disk is released together with the instance: Default value: false.
* `enable_auto_snapshot` - (Optional Available in 1.53.0+) Indicates whether to apply a created automatic snapshot policy to the disk. Default value: false.

-> **NOTE:** Disk category `cloud` has been outdated and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk.
* `status` - The disk status.

## Import

Cloud disk can be imported using the id, e.g.

```
$ terraform import alicloud_disk.example d-abc12345678
```
