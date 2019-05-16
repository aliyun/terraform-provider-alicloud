---
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_configuration"
sidebar_current: "docs-alicloud-resource-ess-scaling-configuration"
description: |-
  Provides a ESS scaling configuration resource.
---

# alicloud\_ess\_scaling\_configuration

Provides a ESS scaling configuration resource.

-> **NOTE:** Several instance types have outdated in some regions and availability zones, such as `ecs.t1.*`, `ecs.s2.*`, `ecs.n1.*` and so on. If you want to keep them, you should set `is_outdated` to true. For more about the upgraded instance type, refer to `alicloud_instance_types` datasource.

## Example Usage

```
resource "alicloud_security_group" "classic" {
  # Other parameters...
}
resource "alicloud_ess_scaling_group" "scaling" {
  min_size           = 1
  max_size           = 2
  removal_policies   = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "config" {
  scaling_group_id  = "${alicloud_ess_scaling_group.scaling.id}"

  image_id          = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type     = "ecs.n4.large"
  security_group_id = "${alicloud_security_group.classic.id}"
}

```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group of a scaling configuration.
* `image_id` - (Required) ID of an image file, indicating the image resource selected when an instance is enabled.
* `instance_type` - (Required) Resource type of an ECS instance.
* `instance_name` - (Optional) Name of an ECS instance. Default to "ESS-Instance". It is valid from version 1.7.1.
* `io_optimized` - (Deprecated) It has been deprecated on instance resource. All the launched alicloud instances will be I/O optimized.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `security_group_id` - (Optional) ID of the security group to which a newly created instance belongs.
* `scaling_configuration_name` - (Optional) Name shown for the scheduled task. If this parameter value is not specified, the default value is ScalingConfigurationId.
* `internet_charge_type` - (Optional) Network billing type, Values: PayByBandwidth or PayByTraffic. Default to `PayByBandwidth`.
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). The value range is [1,200].
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth from the public network, measured in Mbps (Mega bit per second). The value range for PayByBandwidth is [0,100].
* `system_disk_category` - (Optional) Category of the system disk. The parameter value options are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd` and `cloud`. `cloud` only is used to some no I/O optimized instance. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) Size of system disk, in GiB. Optional values: cloud: 40-500, cloud_efficiency: 40-500, cloud_ssd: 40-500, ephemeral_ssd: 40-500 The default value is {40, ImageSize}. If this parameter is set, the system disk size must be greater than or equal to max{40, ImageSize}.
* `enable` - (Optional) Whether enable the specified scaling group(make it active) to which the current scaling configuration belongs.
* `active` - (Optional) Whether active current scaling configuration in the specified scaling group. Default to `false`.
* `substitute` - (Optional) The another scaling configuration which will be active automatically and replace current configuration when setting `active` to 'false'. It is invalid when `active` is 'true'.
* `user_data` - (Optional) User-defined data to customize the startup behaviors of the ECS instance and to pass data into the ECS instance.
* `key_name` - (Optional) The name of key pair that can login ECS instance successfully without password. If it is specified, the password would be invalid.
* `role_name` - (Optional) Instance RAM role name. The name is provided and maintained by RAM. You can use `alicloud_ram_role` to create a new one.
* `force_delete` - (Optional) The last scaling configuration will be deleted forcibly with deleting its scaling group. Default to false.
* `data_disk` - (Optional) DataDisk mappings to attach to ecs instance. See [Block datadisk](#block-datadisk) below for details.
* `instance_ids` - (Deprecated) It has been deprecated from version 1.6.0. New resource `alicloud_ess_attachment` replaces it.
* `tags` - (Optional) A mapping of tags to assign to the resource. It will be applied for ECS instances finally.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "http://", or "https://" It can be a null string.

-> **NOTE:** Before enabling the scaling group, it must have a active scaling configuration.

-> **NOTE:** If the number of attached ECS instances by `instance_ids` is smaller than MinSize, the Auto Scaling Service will automatically create ECS Pay-As-You-Go instance to cater to MinSize. For example, MinSize=5 and 2 existing ECS instances has been attached to the scaling group. When the scaling group is enabled, it will create 3 instnaces automatically based on its current active scaling configuration.

-> **NOTE:** Restrictions on attaching ECS instances:

   - The attached ECS instances and the scaling group must have the same region and network type(`Classic` or `VPC`).
   - The attached ECS instances and the instance with active scaling configurations must have the same instance type.
   - The attached ECS instances must in the running state.
   - The attached ECS instances has not been attached to other scaling groups.
   - The attached ECS instances supports Subscription and Pay-As-You-Go payment methods.

-> **NOTE:** The last scaling configuration can't be set to inactive and deleted alone.


## Block datadisk

The datadisk mapping supports the following:

* `size` - (Optional) Size of data disk, in GB. The value ranges from 5 to 2,000 for a cloud disk and from 5 to 1,024 for an ephemeral disk. A maximum of four values can be entered. 
* `category` - (Optional) Category of data disk. The parameter value options are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd` and `cloud`.
* `snapshot_id` - (Optional) Snapshot used for creating the data disk. If this parameter is specified, the size parameter is neglected, and the size of the created disk is the size of the snapshot. 
* `delete_with_instance` - (Optional) Whether to delete data disks attached on ecs when release ecs instance. Optional value: `true` or `false`, default to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The scaling configuration ID.
* `active` - Wether the current scaling configuration is actived.
* `image_id` - The ecs instance Image id.
* `instance_type` - The ecs instance type.
* `security_group_id` - ID of the security group to which a newly created instance belongs.
* `scaling_configuration_name` - Name of scaling configuration.
* `internet_charge_type` - Internet charge type of ecs instance.
* `key_name` - The name of key pair that has been bound in ECS instance.
* `role_name` - The name of RAM role that has been bound in ECS instance.
* `user_data` - The hash value of the user data.
* `force_delete` - Whether delete the last scaling configuration forcibly with deleting its scaling group.
* `tags` - The scaling instance tags, use jsonencode(item) to display the value.
* `instance_name` - The ecs instance name.