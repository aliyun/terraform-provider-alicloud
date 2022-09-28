---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk"
sidebar_current: "docs-alicloud-resource-ecs-disk"
description: |-
  Provides a Alicloud ECS Disk resource.
---

# alicloud\_ecs\_disk

Provides a ECS Disk resource.

For information about ECS Disk and how to use it, see [What is Disk](https://www.alibabacloud.com/help/en/doc-detail/25513.htm).

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_disk" "example" {
  zone_id     = "cn-beijing-b"
  disk_name   = "tf-test"
  description = "Hello ecs disk."
  category    = "cloud_efficiency"
  size        = "30"
  encrypted   = true
  kms_key_id  = "2a6767f0-a16c-4679-a60f-13bf*****"
  tags = {
    Name = "TerraformTest"
  }
}

```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional, ForceNew) ID of the free zone to which the disk belongs. One of the `zone_id` and `instance_id` must be set but can not be set at the same time.
* `availability_zone` - (Optional, ForceNew) Field `availability_zone` has been deprecated from provider version 1.122.0. New field `zone_id` instead.
* `category` - (Optional) Category of the disk. Valid values are `cloud`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`. Default is `cloud_efficiency`.
* `delete_auto_snapshot` - (Optional) Indicates whether the automatic snapshot is deleted when the disk is released. Default value: `false`.
* `delete_with_instance` - (Optional) Indicates whether the disk is released together with the instance. Default value: `false`.
* `description` - (Optional) Description of the disk. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `disk_name` - (Optional, Computed) Name of the ECS disk. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with `http://` or `https://`. Default value is `null`.
* `dry_run` - (Optional) Specifies whether to check the validity of the request without actually making the request.request Default value: false. Valid values:
    * `true`: The validity of the request is checked but the request is not made. Check items include the required parameters, request format, service limits, and available ECS resources. If the check fails, the corresponding error message is returned. If the check succeeds, the DryRunOperation error code is returned.
    * `false`: The validity of the request is checked. If the check succeeds, a 2xx HTTP status code is returned and the request is made.
* `enable_auto_snapshot` - (Optional) Indicates whether to enable creating snapshot automatically.
* `encrypted` - (Optional, ForceNew) If true, the disk will be encrypted, conflict with `snapshot_id`.
* `instance_id` - (Optional, ForceNew) The ID of the instance to which the created subscription disk is automatically attached.
    * After you specify the instance ID, the specified `resource_group_id`, `tags`, and `kms_key_id` parameters are ignored.
    * One of the `zone_id` and `instance_id` must be set but can not be set at the same time.
* `kms_key_id` - (Optional, ForceNew) The ID of the KMS key corresponding to the data disk, The specified parameter `Encrypted` must be `true` when KmsKeyId is not empty.
* `name` - (Optional, Computed, Deprecated in v1.122.0+) Field `name` has been deprecated from provider version 1.122.0. New field `disk_name` instead.
* `payment_type` - (Optional) Payment method for disk. Valid values: `PayAsYouGo`, `Subscription`. Default to `PayAsYouGo`. If you want to change the disk payment type, the `instance_id` is required.
* `performance_level` - (Optional) Specifies the performance level of an ESSD when you create the ESSD. Valid values:                                                       
    * `PL1`: A single ESSD delivers up to 50,000 random read/write IOPS.
    * `PL2`: A single ESSD delivers up to 100,000 random read/write IOPS.
    * `PL3`: A single ESSD delivers up to 1,000,000 random read/write IOPS.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional) The Id of resource group which the disk belongs.
* `size` - (Optional) The size of the disk in GiBs. When resize the disk, the new size must be greater than the former value, or you would get an error `InvalidDiskSize.TooSmall`.
* `snapshot_id` - (Optional, ForceNew) A snapshot to base the disk off of. If the disk size required by snapshot is greater than `size`, the `size` will be ignored, conflict with `encrypted`.
* `storage_set_id` - (Optional, ForceNew) The ID of the storage set.
* `storage_set_partition_number` - (Optional, ForceNew) The number of partitions in the storage set.
* `type` - (Optional, Available in v1.122.0+) The type to expand cloud disks. Valid Values: `online`, `offline`. Default to `offline`.
    * `offline`: After you resize a disk offline, you must restart the instance by using the console or by calling the RebootInstance operation for the resizing operation to take effect. For more information, see Restart the instance and RebootInstance.
    * `online`: After you resize a disk online, the resizing operation takes effect immediately and you do not need to restart the instance. You can resize ultra disks, standard SSDs, and ESSDs online.

-> **NOTE:** Disk category `cloud` has been outdated and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Disk.
* `status` - The disk status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Disk.
* `update` - (Defaults to 10 mins) Used when update the Disk.
* `delete` - (Defaults to 2 mins) Used when delete the Disk.

## Import

ECS Disk can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_disk.example d-abcd12345
```
