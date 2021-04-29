---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disks"
sidebar_current: "docs-alicloud-datasource-ecs-disks"
description: |-
  Provides a list of Ecs Disks to the user.
---

# alicloud\_ecs\_disks

This data source provides the Ecs Disks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_disks" "example" {
  ids        = ["d-artgdsvdvxxxx"]
  name_regex = "tf-test"
}

output "first_ecs_disk_id" {
  value = data.alicloud_ecs_disks.example.disks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `additional_attributes` - (Optional, ForceNew) Other attribute values. Currently, only the incoming value of IOPS is supported, which means to query the IOPS upper limit of the current disk.
* `auto_snapshot_policy_id` - (Optional, ForceNew) Query cloud disks based on the automatic snapshot policy ID.
* `zone_id` - (Optional, ForceNew) ID of the free zone to which the disk belongs.
* `availability_zone` - (Optional, ForceNew, Deprecated in v1.122.0+) Field `availability_zone` has been deprecated from provider version 1.122.0. New field `zone_id` instead.
* `category` - (Optional, ForceNew) Disk category. Valid values: `cloud`, `cloud_efficiency`, `cloud_essd`, `cloud_ssd`, `ephemeral_ssd`.
* `delete_auto_snapshot` - (Optional, ForceNew) Indicates whether the automatic snapshot is deleted when the disk is released.
* `delete_with_instance` - (Optional, ForceNew) Indicates whether the disk is released together with the instance.
* `disk_name` - (Optional, ForceNew) The disk name.
* `disk_type` - (Optional, ForceNew) The disk type.
* `dry_run` - (Optional, ForceNew) Specifies whether to check the validity of the request without actually making the request.request Default value: false. Valid values:
    * `true`: The validity of the request is checked but the request is not made. Check items include the required parameters, request format, service limits, and available ECS resources. If the check fails, the corresponding error message is returned. If the check succeeds, the DryRunOperation error code is returned.
    * `false`: The validity of the request is checked. If the check succeeds, a 2xx HTTP status code is returned and the request is made.
* `enable_auto_snapshot` - (Optional, ForceNew) Indicates whether the automatic snapshot is deleted when the disk is released.
* `enable_automated_snapshot_policy` - (Optional, ForceNew) Whether the cloud disk has an automatic snapshot policy
* `enable_shared` - (Optional, ForceNew) Whether it is shared block storage.
* `encrypted` - (Optional, ForceNew) Indicate whether the disk is encrypted or not. Possible values: `on` and `off`.
* `ids` - (Optional, ForceNew, Computed)  A list of Disk IDs.
* `instance_id` - (Optional, ForceNew) Filter the results by the specified ECS instance ID.
* `kms_key_id` - (Optional, ForceNew) The kms key id.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Disk name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `payment_type` - (Optional, ForceNew) Payment method for disk. Valid Values: `PayAsYouGo`, `Subscription`.
* `portable` - (Optional, ForceNew) Whether the cloud disk or local disk supports uninstallation.
* `resource_group_id` - (Optional, ForceNew) The Id of resource group which the disk belongs.
* `snapshot_id` - (Optional, ForceNew) The source snapshot id.
* `status` - (Optional, ForceNew) The status of disk.
* `type` - (Optional, ForceNew, Deprecated in v1.122.0+) Field `type` has been deprecated from provider version 1.122.0. New field `disk_type` instead.
* `tags` - (Optional) A map of tags assigned to the disks.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Disk names.
* `disks` - A list of Ecs Disks. Each element contains the following attributes:
    * `id` - ID of the disk.
    * `disk_id` - ID of the disk.
    * `name` - Disk name.
    * `description` - Disk description.
    * `region_id` - Region ID the disk belongs to.
    * `availability_zone` - Availability zone of the disk.
    * `zone_id` - The zone id.
    * `status` - Current status.
    * `type` - Disk type.
    * `category` - Disk category.
    * `encrypted` - Indicate whether the disk is encrypted or not.
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
    * `auto_snapshot_policy_id` - Query cloud disks based on the automatic snapshot policy ID.
    * `delete_auto_snapshot` - Indicates whether the automatic snapshot is deleted when the disk is released.
    * `delete_with_instance` - Indicates whether the disk is released together with the instance.
    * `device` - Cloud disk or the device name of the mounted instance on the site.
    * `disk_name` - The disk name.
    * `enable_auto_snapshot` - Whether the disk implements an automatic snapshot policy.
    * `enable_automated_snapshot_policy` - Whether the disk implements an automatic snapshot policy.
    * `mount_instance_num` - Number of instances mounted on shared storage.
    * `mount_instances` - Disk mount instances.
        * `attached_time` - A mount of time.
        * `device` - The mount point of the disk.
        * `instance_id` - The instance ID of the disk mount.
    * `payment_type` - Payment method for disk.
    * `performance_level` - Performance levels of ESSD cloud disk.
    * `portable` - Whether the disk is unmountable.
    * `product_code` - The product logo of the cloud market.

