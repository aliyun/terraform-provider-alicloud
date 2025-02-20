---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_configuration"
sidebar_current: "docs-alicloud-resource-ess-scaling-configuration"
description: |-
  Provides an ESS scaling configuration resource.
---

# alicloud_ess_scaling_configuration

Provides a ESS scaling configuration resource.

-> **NOTE:** Several instance types have outdated in some regions and availability zones, such as `ecs.t1.*`, `ecs.s2.*`, `ecs.n1.*` and so on. If you want to keep them, you should set `is_outdated` to true. For more about the upgraded instance type, refer to `alicloud_instance_types` datasource.

-> **NOTE:** Available since v1.39.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ess_scaling_configuration&exampleId=a6b0b9a7-a6cb-f13b-ed28-6769b80355a84ea58631&activeTab=example&spm=docs.r.ess_scaling_configuration.0.a6b0b9a7a6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

locals {
  name = "${var.name}-${random_integer.default.result}"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = local.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = local.name
}

resource "alicloud_security_group" "default" {
  name   = local.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = local.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id = alicloud_security_group.default.id
  force_delete      = true
  active            = true
}
```

## Module Support

You can use to the existing [autoscaling module](https://registry.terraform.io/modules/terraform-alicloud-modules/autoscaling/alicloud) 
to create a configuration, scaling group and lifecycle hook one-click.

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group of a scaling configuration.
* `image_id` - (Optional) ID of an image file, indicating the image resource selected when an instance is enabled.
* `image_name` - (Optional, Available since v1.92.0) Name of an image file, indicating the image resource selected when an instance is enabled.
* `instance_type` - (Optional) Resource type of an ECS instance.
* `security_enhancement_strategy` - (Optional, ForceNew, Available since v1.232.0) Specifies whether to enable Security Hardening. Valid values: Active, Deactive.
* `instance_description` - (Optional, Available since v1.232.0) The description of ECS instances. The description must be 2 to 256 characters in length. It can contain letters but cannot start with http:// or https://.
* `spot_duration` - (Optional, Available since v1.232.0) The protection period of preemptible instances. Unit: hours. Valid values: 1, 0.
* `instance_types` - (Optional, Available since v1.46.0) Resource types of an ECS instance.
* `instance_name` - (Optional) Name of an ECS instance. Default to "ESS-Instance". It is valid from version 1.7.1.
* `io_optimized` - (Deprecated) It has been deprecated on instance resource. All the launched alicloud instances will be I/O optimized.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.
* `security_group_id` - (Optional) ID of the security group used to create new instance. It is conflict with `security_group_ids`.
* `security_group_ids` - (Optional, Available since v1.43.0) List IDs of the security group used to create new instances. It is conflict with `security_group_id`.
* `scaling_configuration_name` - (Optional) Name shown for the scheduled task. which must contain 2-64 characters (English or Chinese), starting with numbers, English letters or Chinese characters, and can contain number, underscores `_`, hypens `-`, and decimal point `.`. If this parameter value is not specified, the default value is ScalingConfigurationId.
* `internet_charge_type` - (Optional) Network billing type, Values: PayByBandwidth or PayByTraffic. Default to `PayByBandwidth`.
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). 
* `internet_max_bandwidth_out` - (Optional, Available since v1.214.0) Maximum outgoing bandwidth from the public network, measured in Mbps (Mega bit per second). The value range for PayByBandwidth is [0,1024].
* `credit_specification` - (Optional, Available since v1.98.0) Performance mode of the t5 burstable instance. Valid values: 'Standard', 'Unlimited'.
* `system_disk_category` - (Optional) Category of the system disk. The parameter value options are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd` and `cloud`. `cloud` only is used to some no I/O optimized instance. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) Size of system disk, in GiB. Valid values: Basic disk: 20 to 500, ESSD: The valid values depend on the performance level (PL) of the system disk (PL0 ESSD: 1 to 2048, PL1 ESSD: 20 to 2048, PL2 ESSD: 461 to 2048, PL3 ESSD: 1261 to 2048) , ESSD AutoPL disk: 1 to 2048, Other disk categories: 20 to 2048. The value of this parameter must be at least 1 and greater than or equal to the image size. Default value: 40 or the size of the image, whichever is larger.
* `system_disk_name` - (Optional, Available since v1.92.0) The name of the system disk. It must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-). Default value: null.
* `system_disk_description` - (Optional, Available since v1.92.0) The description of the system disk. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `system_disk_auto_snapshot_policy_id` - (Optional, Available since v1.92.0) The id of auto snapshot policy for system disk.
* `system_disk_performance_level` - (Optional, Available since v1.124.3) The performance level of the ESSD used as the system disk.
* `system_disk_encrypted` - (Optional, Available since v1.199.0) Whether to encrypt the system disk.
* `system_disk_kms_key_id` - (Optional, Available since v1.232.0) The ID of the KMS key that you want to use to encrypt the system disk.
* `system_disk_encrypt_algorithm` - (Optional, Available since v1.232.0) The algorithm that you want to use to encrypt the system disk. Valid values: AES-256, SM4-128.
* `system_disk_provisioned_iops` - (Optional, Available since v1.232.0) IOPS measures the number of read and write operations that an EBS device can process per second. 
* `image_options_login_as_non_root` - (Optional, Available since v1.232.0) Specifies whether to use ecs-user to log on to an ECS instance. For more information, see Manage the username used to log on to an ECS instance. Valid values: true, false. Default value: false.
* `deletion_protection` - (Optional, Available since v1.232.0) Specifies whether to enable the Release Protection feature for ECS instances. This parameter is applicable to only pay-as-you-go instances. You can use this parameter to specify whether an ECS instance can be directly released by using the ECS console or calling the DeleteInstance operation. Valid values: true, false. Default value: false.
* `enable` - (Optional) Whether enable the specified scaling group(make it active) to which the current scaling configuration belongs.
* `active` - (Optional) Whether active current scaling configuration in the specified scaling group. Default to `false`.
* `substitute` - (Optional) The another scaling configuration which will be active automatically and replace current configuration when setting `active` to 'false'. It is invalid when `active` is 'true'.
* `user_data` - (Optional) User-defined data to customize the startup behaviors of the ECS instance and to pass data into the ECS instance.
* `key_name` - (Optional) The name of key pair that can login ECS instance successfully without password. If it is specified, the password would be invalid.
* `role_name` - (Optional) Instance RAM role name. The name is provided and maintained by RAM. You can use `alicloud_ram_role` to create a new one.
* `force_delete` - (Optional) The last scaling configuration will be deleted forcibly with deleting its scaling group. Default to false.
* `data_disk` - (Optional) DataDisk mappings to attach to ecs instance. See [`data_disk`](#data_disk) below for details.
* `instance_pattern_info` - (Optional, Available since v1.177.0) intelligent configuration mode. In this mode, you only need to specify the number of vCPUs, memory size, instance family, and maximum price. The system selects an instance type that is provided at the lowest price based on your configurations to create ECS instances. This mode is available only for scaling groups that reside in virtual private clouds (VPCs). This mode helps reduce the failures of scale-out activities caused by insufficient inventory of instance types.  See [`instance_pattern_info`](#instance_pattern_info) below for details.
* `network_interfaces` - (Optional, Available since v1.235.0) Specify NetworkInterfaces.N to configure primary and secondary ENIs. In this case, specify at least one primary ENI. If you set NetworkInterfaces.N.InstanceType to Primary, a primary ENI is configured. If you set NetworkInterfaces.N.InstanceType to Secondary or leave the parameter empty, a secondary ENI is configured. See [`network_interfaces`](#network_interfaces) below for details.
* `custom_priorities` - (Optional, Available since v1.238.0) You can use CustomPriorities to specify the priority of a custom ECS instance type + vSwitch combination. See [`custom_priorities`](#custom_priorities) below for details.
* `instance_type_override` - (Optional, Available since v1.216.0) specify the weight of instance type.  See [`instance_type_override`](#instance_type_override) below for details.
* `instance_ids` - (Deprecated) It has been deprecated from version 1.6.0. New resource `alicloud_ess_attachment` replaces it.
* `tags` - (Optional) A mapping of tags to assign to the resource. It will be applied for ECS instances finally.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "http://", or "https://" It can be a null string.
* `override` - (Optional, Available since v1.46.0) Indicates whether to overwrite the existing data. Default to false.
* `password_inherit` - (Optional, Available since v1.62.0) Specifies whether to use the password that is predefined in the image. If the PasswordInherit parameter is set to true, the `password` and `kms_encrypted_password` will be ignored. You must ensure that the selected image has a password configured.
* `password` - (Optional, ForceNew, Available since v1.60.0) The password of the ECS instance. The password must be 8 to 30 characters in length. It must contains at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `() ~!@#$%^&*-_+=\|{}[]:;'<>,.?/`, The password of Windows-based instances cannot start with a forward slash (/).
* `kms_encrypted_password` - (Optional, ForceNew, Available since v1.60.0) An KMS encrypts password used to a db account. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available since v1.60.0) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `resource_group_id` - (Optional, Available since v1.135.0) ID of resource group.
* `host_name` - (Optional, Available since v1.143.0) Hostname of an ECS instance.
* `spot_strategy` - (Optional, Available since v1.151.0) The spot strategy for a Pay-As-You-Go instance. Valid values: `NoSpot`, `SpotAsPriceGo`, `SpotWithPriceLimit`.
* `spot_price_limit` - (Optional, Available since v1.151.0) Sets the maximum price hourly for instance types. See [`spot_price_limit`](#spot_price_limit) below for details.

-> **NOTE:** Before enabling the scaling group, it must have a active scaling configuration.

-> **NOTE:** If the number of attached ECS instances by `instance_ids` is smaller than MinSize, the Auto Scaling Service will automatically create ECS Pay-As-You-Go instance to cater to MinSize. For example, MinSize=5 and 2 existing ECS instances has been attached to the scaling group. When the scaling group is enabled, it will create 3 instnaces automatically based on its current active scaling configuration.

-> **NOTE:** Restrictions on attaching ECS instances:

   - The attached ECS instances and the scaling group must have the same region and network type(`Classic` or `VPC`).
   - The attached ECS instances and the instance with active scaling configurations must have the same instance type.
   - The attached ECS instances must in the running state.
   - The attached ECS instances has not been attached to other scaling groups.
   - The attached ECS instances supports Subscription and Pay-As-You-Go payment methods.

-> **NOTE:** The last scaling configuration can't be set to inactive and deleted alone.


### `data_disk`

The datadisk mapping supports the following:

* `size` - (Optional) Size of data disk, in GB. The value ranges [5,2000] for a cloud disk, [5,1024] for an ephemeral disk, [5,800] for an ephemeral_ssd disk, [20,32768] for cloud_efficiency, cloud_ssd, cloud_essd disk. 
* `provisioned_iops` - (Optional, Available since v1.232.0) IOPS measures the number of read and write operations that an Elastic Block Storage (EBS) device can process per second.
* `device` - (Optional, Deprecated, Available since v1.92.0) The mount point of data disk N. Valid values of N: 1 to 16. If this parameter is not specified, the system automatically allocates a mount point to created ECS instances. The name of the mount point ranges from /dev/xvdb to /dev/xvdz in alphabetical order.
* `category` - (Optional) Category of data disk. The parameter value options are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd` , `cloud_essd` and `cloud`.
* `snapshot_id` - (Optional) Snapshot used for creating the data disk. If this parameter is specified, the size parameter is neglected, and the size of the created disk is the size of the snapshot. 
* `delete_with_instance` - (Optional) Whether to delete data disks attached on ecs when release ecs instance. Optional value: `true` or `false`, default to `true`.
* `encrypted` - (Optional, Available since v1.92.0) Specifies whether data disk N is to be encrypted. Valid values of N: 1 to 16. Valid values: `true`: encrypted, `false`: not encrypted. Default value: `false`.
* `kms_key_id` - (Optional, Available since v1.92.0) The CMK ID for data disk N. Valid values of N: 1 to 16.
* `name` - (Optional, Available since v1.92.0) The name of data disk N. Valid values of N: 1 to 16. It must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-). Default value: null.
* `description` - (Optional, Available since v1.92.0) The description of data disk N. Valid values of N: 1 to 16. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `auto_snapshot_policy_id` - (Optional, Available since v1.92.0) The id of auto snapshot policy for data disk.
* `performance_level` - (Optional, Available since v1.124.3) The performance level of the ESSD used as data disk.

### `network_interfaces`

The networkInterfaces mapping supports the following:

* `instance_type` - (Optional, Available since v1.235.0) The ENI type. If you specify NetworkInterfaces.N, specify at least one primary ENI. You cannot specify SecurityGroupId or SecurityGroupIds.N. Valid values: Primary, Secondary.
* `network_interface_traffic_mode` - (Optional, Available since v1.235.0) The communication mode of the ENI. Valid values: Standard, HighPerformance. 
* `ipv6_address_count` - (Optional, Available since v1.235.0) The number of randomly generated IPv6 addresses that you want to assign to primary ENI N.
* `security_group_ids` - (Optional, Available since v1.235.0) The ID of security group N to which ENI N belongs.

### `custom_priorities`

The customPriorities mapping supports the following:

* `instance_type` - (Optional, Available since v1.238.0) This parameter takes effect only if you set Scaling Policy to Priority Policy and the instance type specified by CustomPriorities.N.InstanceType is contained in the scaling configuration.
* `vswitch_id` - (Optional, Available since v1.238.0) This parameter takes effect only if you set Scaling Policy to Priority Policy and the vSwitch specified by CustomPriorities.N.VswitchId is included in the vSwitch list of your scaling group.

### `instance_pattern_info`

The instancePatternInfo mapping supports the following:

* `cores` - (Optional) The number of vCPUs that are specified for an instance type in instancePatternInfo.
* `instance_family_level` - (Optional) The instance family level in instancePatternInfo.
* `max_price` - (Optional) The maximum hourly price for a pay-as-you-go instance or a preemptible instance in instancePatternInfo.
* `memory` - (Optional) The memory size that is specified for an instance type in instancePatternInfo.
* `burstable_performance` - (Optional, Available since v1.220.0) Specifies whether to include burstable instance types.  Valid values: Exclude, Include, Required.
* `excluded_instance_types` - (Optional, Available since v1.220.0) Instance type N that you want to exclude. You can use wildcard characters, such as an asterisk (*), to exclude an instance type or an instance family.
* `architectures` - (Optional, Available since v1.220.0) Architecture N of instance type N. Valid values: X86, Heterogeneous, BareMetal, Arm, SuperComputeCluster.

### `instance_type_override`

The instanceTypeOverride mapping supports the following:

* `instance_type` - (Optional) The is specified for an instance type in instanceTypeOverride.
* `weighted_capacity` - (Optional) The weight of instance type in instanceTypeOverride.

### `spot_price_limit`

The spotPriceLimit mapping supports the following:

* `instance_type` - (Optional, Available since v1.151.0) Resource type of an ECS instance.
* `price_limit` - (Optional, Available since v1.151.0) Price limit hourly of instance type, 2 decimals is allowed at most.


## Attributes Reference

The following attributes are exported:

* `id` - The scaling configuration ID.

## Import

ESS scaling configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_scaling_configuration.example asg-abc123456
```
