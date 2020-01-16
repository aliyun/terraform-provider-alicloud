---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_alarms"
sidebar_current: "docs-alicloud_ess_alarms"
description: |-
    Provides a list of alarms available to the user.
---

# alicloud_ess_alarms

This data source provides available alarm resources. 

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
* `enable` - (Optional)Whether to enable specific ess alarm(`true` or `false`).
* `metric_type` - (Optional) The type for the alarm's associated metric. Supported value: system, custom. "system" means the metric data is collected by Aliyun Cloud Monitor Service(CMS), "custom" means the metric data is upload to CMS by users. Defaults to system.
* `state` - (Optional) The state of alarm task. Supported value: ALARM, OK, INSUFFICIENT_DATA.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `name` - (Optional) The name for ess alarm.
* `description` - (Optional) The description for the alarm.
* `enable` - (Optional, Available in 1.48.0+) Whether to enable specific ess alarm. Default to true.
* `alarm_actions` - (Required) The list of actions to execute when this alarm transition into an ALARM state. Each action is specified as ess scaling rule ari.
* `scaling_group_id` - (Required, ForceNew) The scaling group associated with this alarm, the 'ForceNew' attribute is available in 1.56.0+.
* `metric_type` - (Optional, ForceNew) The type for the alarm's associated metric. Supported value: system, custom. "system" means the metric data is collected by Aliyun Cloud Monitor Service(CMS), "custom" means the metric data is upload to CMS by users. Defaults to system. 
* `metric_name` - (Required) The name for the alarm's associated metric. See [Block_metricNames_and_dimensions](#block-metricnames_and_dimensions) below for details.
* `period` - (Optional, ForceNew) The period in seconds over which the specified statistic is applied. Supported value: 60, 120, 300, 900. Defaults to 300.
* `statistics` - (Optional) The statistic to apply to the alarm's associated metric. Supported value: Average, Minimum, Maximum. Defaults to Average.
* `threshold` - (Required) The value against which the specified statistics is compared.
* `comparison_operator` - (Optional) The arithmetic operation to use when comparing the specified Statistic and Threshold. The specified Statistic value is used as the first operand. Supported value: >=, <=, >, <. Defaults to >=.
* `evaluation_count` - (Optional) The number of times that needs to satisfies comparison condition before transition into ALARM state. Defaults to 3.
* `cloud_monitor_group_id` - (Optional) Defines the application group id defined by CMS which is assigned when you upload custom metric to CMS, only available for custom metirc.
* `dimensions` - (Optional) The dimension map for the alarm's associated metric (documented below). For all metrics, you can not set the dimension key as "scaling_group" or "userId", which is set by default, the second dimension for metric, such as "device" for "PackagesNetIn", need to be set by users.
* `state` - (Optional) The state of alarm task. Supported value: ALARM, OK, INSUFFICIENT_DATA.
  
