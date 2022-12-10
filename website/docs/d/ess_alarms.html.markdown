---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_alarms"
sidebar_current: "docs-alicloud_ess_alarms"
description: |-
    Provides a list of alarms available to the user.
---

# alicloud_ess_alarms

This data source provides available alarm resources. 

-> **NOTE** Available in 1.72.0+

## Example Usage

```
data "alicloud_ess_alarm" "alarm_ds" {
  scaling_group_id = "scaling_group_id"
  ids              = ["alarm_id1", "alarm_id2"]
  name_regex       = "alarm_name"
}

output "first_scaling_rule" {
  value = "${data.alicloud_alarms.alarm_ds.configurations.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional) Scaling group id the alarms belong to.
* `name_regex` - (Optional) A regex string to filter resulting alarms by name.
* `ids` - (Optional) A list of alarm IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `metric_type` - (Optional) The type for the alarm's associated metric. Supported value: system, custom. "system" means the metric data is collected by Aliyun Cloud Monitor Service(CMS), "custom" means the metric data is upload to CMS by users. Defaults to system.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of alarm ids.
* `names` - A list of alarm names.
* `alarms` - A list of alarms. Each element contains the following attributes:
  * `id` - The id of alarm.
  * `name` -  The name for ess alarm.
  * `description` -  The description for the alarm.
  * `enable` - Whether to enable specific ess alarm.
  * `alarm_actions` - The list of actions to execute when this alarm transition into an ALARM state. Each action is specified as ess scaling rule ari.
  * `scaling_group_id` -  The scaling group associated with this alarm.
  * `metric_type` -  The type for the alarm's associated metric. 
  * `metric_name` -  The name for the alarm's associated metric. See [Block_metricNames_and_dimensions](#block-metricnames_and_dimensions) below for details.
  * `period` -  The period in seconds over which the specified statistic is applied.
  * `statistics` -  The statistic to apply to the alarm's associated metric. 
  * `threshold` -  The value against which the specified statistics is compared.
  * `comparison_operator` -  The arithmetic operation to use when comparing the specified Statistic and Threshold. The specified Statistic value is used as the first operand. 
  * `evaluation_count` -  The number of times that needs to satisfies comparison condition before transition into ALARM state. 
  * `cloud_monitor_group_id` -  Defines the application group id defined by CMS which is assigned when you upload custom metric to CMS, only available for custom metirc.
  * `dimensions` -  The dimension map for the alarm's associated metric. 
  * `state` -  The state of alarm task. 
  
