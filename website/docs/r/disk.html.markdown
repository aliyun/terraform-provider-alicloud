---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_disk"
sidebar_current: "docs-alicloud-resource-disk"
description: |-
  Provides a ECS Disk resource.
---

# alicloud\_disk

Provides a ECS disk resource.

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_disk](https://www.terraform.io/docs/providers/alicloud/r/ecs_disk) from version 1.122.0.

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
  encrypted         = true
  kms_key_id        = "2a6767f0-a16c-4679-a60f-13bf*****"
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
* `snapshot_id` - (Optional, ForceNew) A snapshot to base the disk off of. If the disk size required by snapshot is greater than `size`, the `size` will be ignored, conflict with `encrypted`.
* `kms_key_id` - (Optional, Available in 1.89.0+, ForceNew) The ID of the KMS key corresponding to the data disk, The specified parameter `Encrypted` must be `true` when KmsKeyId is not empty.
* `performance_level` - (Optional, Available in 1.95.0+) Specifies the performance level of an ESSD when you create the ESSD. Default value: `PL1`. Valid values:                                                       
    * `PL1`: A single ESSD delivers up to 50,000 random read/write IOPS.
    * `PL2`: A single ESSD delivers up to 100,000 random read/write IOPS.
    * `PL3`: A single ESSD delivers up to 1,000,000 random read/write IOPS.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `encrypted` - (Optional, ForceNew) If true, the disk will be encrypted, conflict with `snapshot_id`.
* `delete_auto_snapshot` - (Optional Available in 1.53.0+) Indicates whether the automatic snapshot is deleted when the disk is released. Default value: false.
* `delete_with_instance` - (Optional Available in 1.53.0+) Indicates whether the disk is released together with the instance: Default value: false.
* `enable_auto_snapshot` - (Optional Available in 1.53.0+) Indicates whether to apply a created automatic snapshot policy to the disk. Default value: false.
* `resource_group_id` - (Optional, Available in 1.57.0+, Modifiable in 1.115.0+) The Id of resource group which the disk belongs.
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
