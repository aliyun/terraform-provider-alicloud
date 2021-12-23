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

-> **NOTE:** Available in 1.79.0+


## Example Usage

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
  name              = var.name
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = "ecs.n1.tiny"
  security_group_id = alicloud_security_group.default.id
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


## Argument Reference

The following arguments are supported:

* `launch_template_id` - (Required, ForceNew) The ID of the instance launch template associated with the auto provisioning group.
* `total_target_capacity` - (Required) The total target capacity of the auto provisioning group. The target capacity consists of the following three parts:PayAsYouGoTargetCapacity,SpotTargetCapacity and the supplemental capacity besides PayAsYouGoTargetCapacity and SpotTargetCapacity.
* `auto_provisioning_group_name` - (Optional) The name of the auto provisioning group to be created. It must be 2 to 128 characters in length. It must start with a letter but cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-)
* `auto_provisioning_group_type` - (Optional, ForceNew) The type of the auto provisioning group. Valid values:`request` and `maintain`,Default value: `maintain`.
* `spot_allocation_strategy` - (Optional, ForceNew) The scale-out policy for preemptible instances. Valid values:`lowest-price` and `diversified`,Default value: `lowest-price`.
* `spot_target_capacity` - (Optional) The target capacity of preemptible instances in the auto provisioning group.
* `spot_instance_interruption_behavior` - (Optional, ForceNew) The default behavior after preemptible instances are shut down. Value values: `stop` and `terminate`,Default value: `stop`.
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
* `launch_template_config` - (Required, ForceNew) DataDisk mappings to attach to ecs instance. See [Block config](#block-config) below for details.

## Block config

The config mapping supports the following:
* `instance_type` - (Optional) The instance type of the Nth extended configurations of the launch template.
* `max_price` - (Required) The maximum price of the instance type specified in the Nth extended configurations of the launch template.
* `vswitch_id` - (Required) The ID of the VSwitch in the Nth extended configurations of the launch template.
* `weighted_capacity` - (Required) The weight of the instance type specified in the Nth extended configurations of the launch template.
* `priority` - (Optional) The priority of the instance type specified in the Nth extended configurations of the launch template. A value of 0 indicates the highest priority.
                     
## Attributes Reference

The following attributes are exported:
* `id` - The ID of the auto provisioning group

## Import

ECS auto provisioning group can be imported using the id, e.g.

```
$ terraform import alicloud_auto_provisioning_group.example asg-abc123456
```
