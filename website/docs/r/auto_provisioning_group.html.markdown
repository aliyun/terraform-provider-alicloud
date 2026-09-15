---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_auto_provisioning_group"
sidebar_current: "docs-alicloud-resource-auto-provisioning-group"
description: |-
  Provides a ECS Auto Provisioning group resource.
---

# alicloud\_auto\_provisioning\_group

Provides a ECS auto provisioning group resource which is a solution that uses preemptive instances and pay_as_you_go instances to rapidly deploy clusters.

-> **NOTE:** Available since v1.79.0+


## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_auto_provisioning_group&exampleId=e055ce1b-2031-1721-49c0-10d93d461580164c92fc&activeTab=example&spm=docs.r.auto_provisioning_group.0.e055ce1b20&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "auto_provisioning_group"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_auto_provisioning_group" "default" {
  launch_template_id            = alicloud_ecs_launch_template.template.id
  total_target_capacity         = "4"
  pay_as_you_go_target_capacity = "1"
  spot_target_capacity          = "2"
  launch_template_config {
    instance_type     = "ecs.n1.small"
    vswitch_id        = alicloud_vswitch.default.id
    weighted_capacity = "2"
    max_price         = "2"
  }
}

resource "alicloud_ecs_launch_template" "template" {
  launch_template_name = var.name
  image_id             = data.alicloud_images.default.images[0].id
  instance_type        = "ecs.n1.tiny"
  security_group_id    = alicloud_security_group.default.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_auto_provisioning_group&spm=docs.r.auto_provisioning_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `launch_template_id` - (Optional, ForceNew) The ID of the instance launch template associated with the auto provisioning group. If not specified, `launch_configuration` must be provided to define instance launch parameters. When both `launch_template_id` and `launch_configuration` are specified, `launch_template_id` takes precedence.
* `total_target_capacity` - (Required) The total target capacity of the auto provisioning group. The target capacity consists of the following three parts:PayAsYouGoTargetCapacity,SpotTargetCapacity and the supplemental capacity besides PayAsYouGoTargetCapacity and SpotTargetCapacity.
* `auto_provisioning_group_name` - (Optional) The name of the auto provisioning group to be created. It must be 2 to 128 characters in length. It must start with a letter but cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-)
* `auto_provisioning_group_type` - (Optional, ForceNew) The type of the auto provisioning group. Valid values:`request` and `maintain`,Default value: `maintain`.
* `spot_allocation_strategy` - (Optional, ForceNew) The scale-out policy for preemptible instances. Valid values:`lowest-price` and `diversified`,Default value: `lowest-price`.
* `spot_target_capacity` - (Optional) The target capacity of preemptible instances in the auto provisioning group.
* `spot_instance_interruption_behavior` - (Optional, ForceNew) The default behavior after preemptible instances are shut down. Valid values: `stop` and `terminate`,Default value: `stop`.
* `spot_instance_pools_to_use_count` - (Optional, ForceNew) This parameter takes effect when the `SpotAllocationStrategy` parameter is set to `lowest-price`. The auto provisioning group selects instance types of the lowest cost to create instances.
* `pay_as_you_go_allocation_strategy` - (Optional, ForceNew) The scale-out policy for pay-as-you-go instances. Valid values: `lowest-price` and `prioritized`,Default value: `lowest-price`.
* `pay_as_you_go_target_capacity` - (Optional) The target capacity of pay-as-you-go instances in the auto provisioning group.
* `default_target_capacity_type` - (Optional) The type of supplemental instances. When the total value of `PayAsYouGoTargetCapacity` and `SpotTargetCapacity` is smaller than the value of TotalTargetCapacity, the auto provisioning group will create instances of the specified type to meet the capacity requirements. Valid values:`PayAsYouGo`: Pay-as-you-go instances; `Spot`: Preemptible instances, Default value: `Spot`.
* `launch_template_version` - (Optional, ForceNew) The version of the instance launch template associated with the auto provisioning group.
* `excess_capacity_termination_policy` - (Optional) The shutdown policy for excess preemptible instances followed when the capacity of the auto provisioning group exceeds the target capacity. Valid values: `no-termination` and `termination`,Default value: `no-termination`.
* `terminate_instances_with_expiration` - (Optional) The shutdown policy for preemptible instances when the auto provisioning group expires. Valid values: `false` and `true`, default value: `false`.
* `terminate_instances` - (Optional, ForceNew) Specifies whether to release instances of the auto provisioning group. Valid values:`false` and `true`, default value: `false`.
* `description` - (Optional, ForceNew) The description of the auto provisioning group.
* `max_spot_price` - (Optional) The global maximum price for preemptible instances in the auto provisioning group. If both the `MaxSpotPrice` and `LaunchTemplateConfig.N.MaxPrice` parameters are specified, the maximum price is the lower value of the two.
* `valid_from` - (Optional, ForceNew) The time when the auto provisioning group is started. The period of time between this point in time and the point in time specified by the `valid_until` parameter is the effective time period of the auto provisioning group.By default, an auto provisioning group is immediately started after creation.
* `valid_until` - (Optional, ForceNew) The time when the auto provisioning group expires. The period of time between this point in time and the point in time specified by the `valid_from` parameter is the effective time period of the auto provisioning group.By default, an auto provisioning group never expires.
* `launch_template_config` - (Required, ForceNew) DataDisk mappings to attach to ecs instance. See [`config`](#block-config) below for details.
* `launch_configuration` - (Optional, ForceNew) The launch configuration block that defines instance launch parameters as an alternative to `launch_template_id`. When `launch_template_id` is not specified, this block provides the instance configuration (image, disk, security, network, etc.). These fields are create-only: the Describe API does not return them, so they are not populated during refresh. See [`launch_configuration`](#launch_configuration) below for details.

### `block-config`

The config mapping supports the following:
* `instance_type` - (Optional) The instance type of the Nth extended configurations of the launch template.
* `max_price` - (Required) The maximum price of the instance type specified in the Nth extended configurations of the launch template.
* `vswitch_id` - (Required) The ID of the VSwitch in the Nth extended configurations of the launch template.
* `weighted_capacity` - (Required) The weight of the instance type specified in the Nth extended configurations of the launch template.
* `priority` - (Optional) The priority of the instance type specified in the Nth extended configurations of the launch template. A value of 0 indicates the highest priority.

### `launch_configuration`

The launch configuration supports the following:
* `image_id` - (Optional) The ID of the image used to create instances.
* `image_family` - (Optional) The name of the image family.
* `instance_name` - (Optional) The name of the instance.
* `instance_description` - (Optional) The description of the instance.
* `host_name` - (Optional) The hostname of the instance. Cannot be specified together with `host_names`.
* `host_names` - (Optional) The hostnames of the instances. Cannot be specified together with `host_name`.
* `credit_specification` - (Optional) The credit specification. Valid values: `standard`, `unlimited`, `unrestricted`.
* `deployment_set_id` - (Optional) The ID of the deployment set.
* `auto_release_time` - (Optional) The automatic release time of the instance.
* `io_optimized` - (Optional) Whether the instance is I/O optimized. Valid values: `optimized`, `none`.
* `security_group_id` - (Optional) The ID of the security group to which to assign the instance.
* `security_group_ids` - (Optional) The IDs of the security groups to which to assign the instance.
* `security_enhancement_strategy` - (Optional) The security enhancement strategy. Valid values: `DeletionProtection`, `NoDeletionProtection`.
* `internet_max_bandwidth_in` - (Optional) The maximum inbound public bandwidth.
* `internet_max_bandwidth_out` - (Optional) The maximum outbound public bandwidth.
* `internet_charge_type` - (Optional) The billing method for network usage. Valid values: `PayByTraffic`, `PayByBandwidth`.
* `password` - (Optional, Sensitive) The password of the instance.
* `password_inherit` - (Optional) Whether to inherit the password from the launch template.
* `key_pair_name` - (Optional) The name of the key pair.
* `ram_role_name` - (Optional) The name of the instance RAM role.
* `user_data` - (Optional) The user data of the instance, Base64 encoded.
* `resource_group_id` - (Optional) The ID of the resource group to which to assign the instance.
* `system_disk_category` - (Optional) The category of the system disk.
* `system_disk_size` - (Optional) The size of the system disk. Valid values: 20 to 500.
* `system_disk_performance_level` - (Optional) The performance level of the system disk.
* `system_disk_name` - (Optional) The name of the system disk.
* `system_disk_description` - (Optional) The description of the system disk.
* `system_disk` - (Optional) The system disk encryption settings. See [`system_disk`](#launch_configuration-system_disk) below.
* `data_disk` - (Optional) The data disk configurations. See [`data_disk`](#launch_configuration-data_disk) below.
* `network_interface` - (Optional) The network interface configurations. See [`network_interface`](#launch_configuration-network_interface) below.
* `tag` - (Optional) The tags of the instance. See [`tag`](#launch_configuration-tag) below.
* `arn` - (Optional) The ARN configurations for the instance. See [`arn`](#launch_configuration-arn) below.
* `period` - (Optional) The subscription period.
* `period_unit` - (Optional) The unit of the subscription period.
* `auto_renew` - (Optional) Whether to enable auto-renewal.
* `auto_renew_period` - (Optional) The auto-renewal period.
* `additional_info` - (Optional) Additional information. See [`additional_info`](#launch_configuration-additional_info) below.

### `launch_configuration-system_disk`

The system disk encryption settings:
* `encrypted` - (Optional) Whether to encrypt the system disk.
* `kms_key_id` - (Optional) The KMS key ID for encryption.
* `encrypt_algorithm` - (Optional) The encryption algorithm.
* `provisioned_iops` - (Optional) The provisioned IOPS.
* `bursting_enabled` - (Optional) Whether bursting is enabled.

### `launch_configuration-data_disk`

The data disk configurations:
* `category` - (Optional) The category of the data disk.
* `disk_name` - (Optional) The name of the data disk.
* `size` - (Optional) The size of the data disk.
* `device` - (Optional) The device name of the data disk.
* `snapshot_id` - (Optional) The snapshot ID used to create the data disk.
* `description` - (Optional) The description of the data disk.
* `delete_with_instance` - (Optional) Whether to delete the data disk with the instance.
* `encrypted` - (Optional) Whether to encrypt the data disk.
* `kms_key_id` - (Optional) The KMS key ID for encryption.
* `encrypt_algorithm` - (Optional) The encryption algorithm.
* `performance_level` - (Optional) The performance level of the data disk.
* `provisioned_iops` - (Optional) The provisioned IOPS.
* `bursting_enabled` - (Optional) Whether bursting is enabled.

### `launch_configuration-network_interface`

The network interface configurations:
* `security_group_id` - (Optional) The ID of the security group.
* `security_group_ids` - (Optional) The IDs of the security groups.
* `instance_type` - (Optional) The instance type.

### `launch_configuration-tag`

The tag configurations:
* `key` - (Optional) The key of the tag.
* `value` - (Optional) The value of the tag.

### `launch_configuration-arn`

The ARN configurations:
* `rolearn` - (Optional) The ARN of the role.
* `role_type` - (Optional) The type of the role.
* `assume_role_for` - (Optional) The assume role for.

### `launch_configuration-additional_info`

The additional information:
* `pvd_config` - (Optional) The PVD configuration.
                     
## Attributes Reference

The following attributes are exported:
* `id` - The ID of the auto provisioning group

## Import

ECS auto provisioning group can be imported using the id, e.g.

```shell
$ terraform import alicloud_auto_provisioning_group.example asg-abc123456
```
