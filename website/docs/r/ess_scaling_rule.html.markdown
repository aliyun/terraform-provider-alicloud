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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ess_scaling_rule&exampleId=a4276f56-b4bf-7f65-e800-a402fcec1c1a4c58fcaf&activeTab=example&spm=docs.r.ess_scaling_rule.0.a4276f56b4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "ap-southeast-5"
}

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
  security_group_name = local.name
  vpc_id              = alicloud_vpc.default.id
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
* `scaling_rule_type` - (Optional, ForceNew, Available since v1.58.0) The scaling rule type, either "SimpleScalingRule", "TargetTrackingScalingRule", "StepScalingRule", "PredictiveScalingRule". Default to "SimpleScalingRule".
* `estimated_instance_warmup` - (Optional, Available since v1.58.0) The estimated time, in seconds, until a newly launched instance will contribute CloudMonitor metrics. Default to 300.
* `min_adjustment_magnitude` - (Optional, Available since v1.221.0) The minimum number of instances that must be scaled. This parameter takes effect if you set ScalingRuleType to SimpleScalingRule or StepScalingRule, and AdjustmentType to PercentChangeInCapacity.
* `scale_in_evaluation_count` - (Optional, Available since v1.221.0) The number of consecutive times that the event-triggered task created for scale-ins must meet the threshold conditions before an alert is triggered. After a target tracking scaling rule is created, an event-triggered task is automatically created and associated with the target tracking scaling rule.
* `scale_out_evaluation_count` - (Optional, Available since v1.221.0) The number of consecutive times that the event-triggered task created for scale-outs must meet the threshold conditions before an alert is triggered. After a target tracking scaling rule is created, an event-triggered task is automatically created and associated with the target tracking scaling rule.
* `metric_name` - (Optional, Available since v1.58.0) A CloudMonitor metric name.
* `target_value` - (Optional, Available since v1.58.0) The target value for the metric.
* `disable_scale_in` - (Optional, Available since v1.58.0) Indicates whether scale in by the target tracking policy is disabled. Default to false.
* `step_adjustment` - (Optional, Available since v1.58.0) Steps for StepScalingRule. See [`step_adjustment`](#step_adjustment) below.
* `ari` - (Optional) The unique identifier of the scaling rule.
* `alarm_dimension` - (Optional, ForceNew, Available since v1.216.0) AlarmDimension for StepScalingRule. See [`alarm_dimension`](#alarm_dimension) below.
* `predictive_scaling_mode` - (Optional, Available since v1.222.0) The mode of the predictive scaling rule. Valid values: PredictAndScale, PredictOnly.
* `initial_max_size` - (Optional, Available since v1.222.0) The maximum number of ECS instances that can be added to the scaling group. If you specify InitialMaxSize, you must also specify PredictiveValueBehavior.
* `predictive_value_behavior` - (Optional, Available since v1.222.0) The action on the predicted maximum value. Valid values: MaxOverridePredictiveValue, PredictiveValueOverrideMax, PredictiveValueOverrideMaxWithBuffer.
* `predictive_value_buffer` - (Optional, Available since v1.222.0) The ratio based on which the predicted value is increased if you set PredictiveValueBehavior to PredictiveValueOverrideMaxWithBuffer. If the predicted value increased by this ratio is greater than the initial maximum capacity, the increased value is used as the maximum value for prediction tasks. Valid values: 0 to 100.
* `predictive_task_buffer_time` - (Optional, Available since v1.222.0) The amount of buffer time before the prediction task runs. By default, all prediction tasks that are automatically created by a predictive scaling rule run on the hour. You can specify a buffer time to run prediction tasks and prepare resources in advance. Valid values: 0 to 60. Unit: minutes.
* `hybrid_monitor_namespace` - (Optional, Available since v1.245.0) The ID of the Hybrid Cloud Monitoring metric repository.
* `metric_type` - (Optional, Available since v1.245.0) The type of the metric. Valid values: system, custom, hybrid.
* `hybrid_metrics` - (Optional, Available since v1.245.0) The Hybrid Cloud Monitoring metrics. See [`hybrid_metrics`](#hybrid_metrics) below.


### `step_adjustment`

The stepAdjustment mapping supports the following:

* `metric_interval_lower_bound` - (Optional) The lower bound of step.
* `metric_interval_upper_bound` - (Optional) The upper bound of step.
* `scaling_adjustment` - (Optional) The adjust value of step.

### `hybrid_metrics`

The HybridMetrics mapping supports the following:

* `id` - (Optional) The reference ID of the metric in the metric expression.
* `metric_name` - (Optional) The name of the Hybrid Cloud Monitoring metric.
* `statistic` - (Optional) The statistical method of the metric value. Valid values: Average, Minimum, Maximum.
* `expression` - (Optional) The metric expression that consists of multiple Hybrid Cloud Monitoring metrics. It calculates a result used to trigger scaling events. The expression must comply with the Reverse Polish Notation (RPN) specification, and the operators can only be + - × /.
* `metric_name` - (Optional) The name of the Hybrid Cloud Monitoring metric.
* `dimensions` - (Optional) The structure of volumeMounts.
  See [`dimensions`](#hybrid_metrics-dimensions) below for details.

### `hybrid_metrics-dimensions`

The dimensions supports the following:

* `dimension_key` - (Optional) The key of the metric dimension.
* `dimension_value` - (Optional) The value of the metric dimension.

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
