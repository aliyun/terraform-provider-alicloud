---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_auto_scaling_alarm_task"
description: |-
  Provides a Alicloud Auto Scaling Alarm Task resource.
---

# alicloud_auto_scaling_alarm_task

Provides a Auto Scaling Alarm Task resource.



For information about Auto Scaling Alarm Task and how to use it, see [What is Alarm Task](https://next.api.alibabacloud.com/document/Ess/2014-08-28/CreateAlarm).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

No example is available because resource test cases were not supplied.

## Argument Reference

The following arguments are supported:
* `alarm_actions` - (Optional, List) Unique identifier of the alarm task Association scaling rule
* `alarm_task_name` - (Optional) The name of the resource
* `comparison_operator` - (Optional) Comparison between the statistical value of the monitoring item and the threshold.
* `description` - (Optional) Description of the alarm task
* `dimension` - (Optional, List) The list of dimension information associated with the monitoring item. See [`dimension`](#dimension) below.
* `effective` - (Optional) Specify the effective time period of the alarm task
* `evaluation_count` - (Optional, Int) The number of times the threshold expression must be met to trigger the execution of a scaling rule.
* `expression` - (Optional, List) List of monitoring indicators for multi-indicator alarm rules See [`expression`](#expression) below.
* `expressions_logic_operator` - (Optional) Relationship between multi-indicator alarm conditions
* `metric_name` - (Optional) Monitoring Item Name
* `metric_type` - (Optional) Monitoring Item Type
* `period` - (Optional, Int) Cycle of statistical monitoring item data
* `scaling_group_id` - (Required, ForceNew) The ID of the scaling group associated with the alarm task.
* `statistics` - (Optional) Methods for statistical monitoring item data
* `status` - (Optional, Computed) The status of the resource
* `threshold` - (Optional, Float) Thresholds for monitoring metrics

### `dimension`

The dimension supports the following:
* `dimension_key` - (Optional) The Dimension Information key associated with the monitoring item.
* `dimension_value` - (Optional) The dimension information value associated with the monitoring item.

### `expression`

The expression supports the following:
* `comparison_operator` - (Optional) The comparison between the metric statistics value and the threshold value of the multi-indicator alarm rule.
* `metric_name` - (Optional) Metric name of the multi-metric alarm rule
* `period` - (Optional, Int) The period of statistical monitoring item data in multi-indicator alarm rules.
* `statistics` - (Optional) The method of statistical monitoring item data in multi-indicator alarm rules.
* `threshold` - (Optional, Float) Thresholds of monitoring metrics in multi-metric alarm rules

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `alarm_trigger_state` - The trigger status of the alarm task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Alarm Task.
* `delete` - (Defaults to 5 mins) Used when delete the Alarm Task.
* `update` - (Defaults to 5 mins) Used when update the Alarm Task.

## Import

Auto Scaling Alarm Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_auto_scaling_alarm_task.example <alarm_task_id>
```