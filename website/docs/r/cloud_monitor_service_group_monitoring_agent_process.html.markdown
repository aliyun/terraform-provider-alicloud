---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_group_monitoring_agent_process"
description: |-
  Provides a Alicloud Cloud Monitor Service Group Monitoring Agent Process resource.
---

# alicloud_cloud_monitor_service_group_monitoring_agent_process

Provides a Cloud Monitor Service Group Monitoring Agent Process resource.

For information about Cloud Monitor Service Group Monitoring Agent Process and how to use it, see [What is Group Monitoring Agent Process](https://www.alibabacloud.com/help/en/cms/developer-reference/api-cms-2019-01-01-creategroupmonitoringagentprocess).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_monitor_service_group_monitoring_agent_process&exampleId=c306705e-2514-5f60-ef62-3aec0fe7411d33f22f24&activeTab=example&spm=docs.r.cloud_monitor_service_group_monitoring_agent_process.0.c306705e25&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
  contacts                 = ["user", "user1", "user2"]
}

resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = var.name
  contact_groups     = [alicloud_cms_alarm_contact_group.default.id]
}

resource "alicloud_cloud_monitor_service_group_monitoring_agent_process" "default" {
  group_id                      = alicloud_cms_monitor_group.default.id
  process_name                  = var.name
  match_express_filter_relation = "or"
  match_express {
    name     = var.name
    value    = "*"
    function = "all"
  }
  alert_config {
    escalations_level   = "critical"
    comparison_operator = "GreaterThanOrEqualToThreshold"
    statistics          = "Average"
    threshold           = "20"
    times               = "100"
    effective_interval  = "00:00-22:59"
    silence_time        = "85800"
    webhook             = "https://www.aliyun.com"
    target_list {
      target_list_id = "1"
      json_params    = "{}"
      level          = "WARN"
      arn            = "acs:mns:cn-hangzhou:120886317861****:/queues/test123/message"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The ID of the application group.
* `process_name` - (Required, ForceNew) The name of the process.
* `match_express_filter_relation` - (Optional, ForceNew) The logical operator used between conditional expressions that are used to match instances. Valid values: `all`, `and`, `or`.
* `match_express` - (Optional, ForceNew, Set) The expressions used to match instances. See [`match_express`](#match_express) below.
* `alert_config` - (Required, Set) The alert rule configurations. See [`alert_config`](#alert_config) below.

### `match_express`

The match_express supports the following:

* `name` - (Optional, ForceNew) The criteria based on which the instances are matched.
* `value` - (Optional, ForceNew) The keyword used to match the instance name.
* `function` - (Optional, ForceNew) The matching condition. Valid values: `all`, `startWith`, `endWith`, `contains`, `notContains`, `equals`.

### `alert_config`

The alert_config supports the following:

* `escalations_level` (Required) The alert level. Valid values: `critical`, `warn`, `info`.
* `comparison_operator` (Required) The operator that is used to compare the metric value with the threshold. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
* `statistics` (Required) The statistical method for alerts. Valid values: `Average`.
* `threshold` (Required) The alert threshold.
* `times` (Required) The number of times for which the threshold can be consecutively exceeded.
* `effective_interval` (Optional) The time period during which the alert rule is effective.
* `silence_time` (Optional, Int) The mute period during which new alert notifications are not sent even if the trigger conditions are met. Unit: seconds.
* `webhook` (Optional) The callback URL.
* `target_list` (Optional, Set) The alert triggers. See [`target_list`](#alert_config-target_list) below.

### `alert_config-target_list`

The target_list supports the following:

* `target_list_id` (Optional) The ID of the resource for which alerts are triggered.
* `json_params` (Optional) The parameters of the alert callback. Specify the parameters in the JSON format.
* `level` (Optional) The alert level. Valid values: `CRITICAL`, `WARN`, `INFO`.
* `arn` (Optional) The Alibaba Cloud Resource Name (ARN) of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Group Monitoring Agent Process. It formats as `<group_id>:<group_monitoring_agent_process_id>`.
* `group_monitoring_agent_process_id` - The ID of the Group Monitoring Agent Process.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Group Monitoring Agent Process.
* `update` - (Defaults to 5 mins) Used when update the Group Monitoring Agent Process.
* `delete` - (Defaults to 5 mins) Used when delete the Group Monitoring Agent Process.

## Import

Cloud Monitor Service Group Monitoring Agent Process can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_monitor_service_group_monitoring_agent_process.example <group_id>:<group_monitoring_agent_process_id>
```
