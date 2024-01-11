---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scaling_rule"
sidebar_current: "docs-alicloud-resource-ess-scaling-rule"
description: |-
  Provides a ESS scaling rule resource.
---

# alicloud_ess_scaling_rule

Provides a ESS scaling rule resource.

For information about ess scaling rule, see [CreateScalingRule](https://www.alibabacloud.com/help/en/auto-scaling/latest/createscalingrule).

-> **NOTE:** Available since v1.39.0.

## Example Usage

```terraform
resource "random_integer" "default" {
  count = 1
  max   = 99999
  min   = 10000
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
  name       = "terraform-example-${random_integer.default[0].result}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "terraform-example-${random_integer.default[0].result}"
}

resource "alicloud_security_group" "default" {
  name   = "terraform-example-${random_integer.default[0].result}"
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
  scaling_group_name = "terraform-example-${random_integer.default[0].result}"
  vswitch_ids        = [alicloud_vswitch.default.id]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id = alicloud_security_group.default.id
  force_delete      = "true"
}

resource "alicloud_ess_scaling_rule" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  adjustment_type  = "TotalCapacity"
  adjustment_value = 1
}
```

## Module Support

You can use to the existing [autoscaling-rule module](https://registry.terraform.io/modules/terraform-alicloud-modules/autoscaling-rule/alicloud) 
to create different type rules, alarm task and scheduled task one-click.

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group of a scaling rule.
* `adjustment_type` - (Optional) Adjustment mode of a scaling rule. Optional values:
    - QuantityChangeInCapacity: It is used to increase or decrease a specified number of ECS instances.
    - PercentChangeInCapacity: It is used to increase or decrease a specified proportion of ECS instances.
    - TotalCapacity: It is used to adjust the quantity of ECS instances in the current scaling group to a specified value.
* `adjustment_value` - (Optional) The number of ECS instances to be adjusted in the scaling rule. This parameter is required and applicable only to simple scaling rules. The number of ECS instances to be adjusted in a single scaling activity cannot exceed 500. Value range:
    - QuantityChangeInCapacity：(0, 500] U (-500, 0]
    - PercentChangeInCapacity：[0, 10000] U [-100, 0]
    - TotalCapacity：[0, 1000]
* `scaling_rule_name` - (Optional) Name shown for the scaling rule, which must contain 2-64 characters (English or Chinese), starting with numbers, English letters or Chinese characters, and can contain number, underscores `_`, hypens `-`, and decimal point `.`. If this parameter value is not specified, the default value is scaling rule id. 
* `cooldown` - (Optional) The cooldown time of the scaling rule. This parameter is applicable only to simple scaling rules. Value range: [0, 86,400], in seconds. The default value is empty，if not set, the return value will be 0, which is the default value of integer.
* `scaling_rule_type` - (Optional, ForceNew, Available in 1.58.0+) The scaling rule type, either "SimpleScalingRule", "TargetTrackingScalingRule", "StepScalingRule". Default to "SimpleScalingRule".
* `estimated_instance_warmup` - (Optional, Available in 1.58.0+) The estimated time, in seconds, until a newly launched instance will contribute CloudMonitor metrics. Default to 300.
* `metric_name` - (Optional, Available in 1.58.0+) A CloudMonitor metric name.
* `target_value` - (Optional, Available in 1.58.0+) The target value for the metric.
* `disable_scale_in` - (Optional, Available in 1.58.0+) Indicates whether scale in by the target tracking policy is disabled. Default to false.
* `step_adjustment` - (Optional, Available in 1.58.0+) Steps for StepScalingRule. See [`step_adjustment`](#step_adjustment) below.
* `ari` - (Optional) The unique identifier of the scaling rule.
* `alarm_dimension` - (Optional, ForceNew, Available in 1.215.0+) AlarmDimension for StepScalingRule. See [`alarm_dimension`](#alarm_dimension) below.


### `step_adjustment`

The stepAdjustment mapping supports the following:

* `metric_interval_lower_bound` - (Optional) The lower bound of step.
* `metric_interval_upper_bound` - (Optional) The upper bound of step.
* `scaling_adjustment` - (Optional) The adjust value of step.

### `alarm_dimension`

The alarmDimension mapping supports the following:

* `dimension_key` - (Optional) The dimension key of the metric.
* `dimension_value` - (Optional, ForceNew) The dimension value of the metric.


## Attributes Reference

The following attributes are exported:

* `id` - The scaling rule ID.

## Import

ESS scaling rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_scaling_rule.example abc123456
```
