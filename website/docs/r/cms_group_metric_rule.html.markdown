---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_group_metric_rule"
sidebar_current: "docs-alicloud-resource-cms-group-metric-rule"
description: |-
  Provides a Alicloud Cloud Monitor Service Group Metric Rule resource.
---

# alicloud_cms_group_metric_rule

Provides a Cloud Monitor Service Group Metric Rule resource.

For information about Cloud Monitor Service Group Metric Rule and how to use it, see [What is Group Metric Rule](https://www.alibabacloud.com/help/en/cloudmonitor/latest/putgroupmetricrule).

-> **NOTE:** Available since v1.104.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_group_metric_rule&exampleId=8762090b-bbee-7c03-4f59-c460b85a4837a4ebef91&activeTab=example&spm=docs.r.cms_group_metric_rule.0.8762090bbb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
  describe                 = var.name
}

resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = var.name
  contact_groups     = [alicloud_cms_alarm_contact_group.default.id]
}

resource "alicloud_cms_group_metric_rule" "this" {
  group_id               = alicloud_cms_monitor_group.default.id
  group_metric_rule_name = var.name
  category               = "ecs"
  metric_name            = "cpu_total"
  namespace              = "acs_ecs_dashboard"
  rule_id                = var.name
  period                 = "60"
  interval               = "3600"
  silence_time           = 85800
  no_effective_interval  = "00:00-05:30"
  webhook                = "http://www.aliyun.com"
  escalations {
    warn {
      comparison_operator = "GreaterThanOrEqualToThreshold"
      statistics          = "Average"
      threshold           = "90"
      times               = 3
    }
    info {
      comparison_operator = "LessThanLastWeek"
      statistics          = "Average"
      threshold           = "90"
      times               = 5
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Required, ForceNew) The ID of the alert rule.
* `group_id` - (Required) The ID of the application group.
* `group_metric_rule_name` - (Required) The name of the alert rule.
* `metric_name` - (Required) The name of the metric.
* `namespace` - (Required, ForceNew) The namespace of the service.
* `category` - (Optional) The abbreviation of the service name.
* `contact_groups` - (Optional) Alarm contact group.
* `dimensions` - (Optional) The dimensions that specify the resources to be associated with the alert rule.
* `email_subject` - (Optional) The subject of the alert notification email.
* `effective_interval` - (Optional) The time period during which the alert rule is effective.
* `no_effective_interval` - (Optional) The time period during which the alert rule is ineffective.
* `interval` - (Optional, ForceNew) The interval at which Cloud Monitor checks whether the alert rule is triggered. Unit: seconds.
* `period` - (Optional, Int) The aggregation period of the monitoring data. Unit: seconds. The value is an integral multiple of 60. Default value: `300`.
* `silence_time` - (Optional, Int) The mute period during which new alerts are not reported even if the alert trigger conditions are met. Unit: seconds. Default value: `86400`, which is equivalent to one day.
* `webhook` - (Optional) The callback URL.
* `targets` - (Optional, Set, Available since v1.189.0) The information about the resource for which alerts are triggered. See [`targets`](#targets) below.
* `escalations` - (Required, Set) Alarm level. See [`escalations`](#escalations) below.

### `targets`

The targets supports the following:

* `id` - (Optional) The ID of the resource for which alerts are triggered.
* `json_params` - (Optional) The parameters of the alert callback. The parameters are in the JSON format.
* `level` - (Optional) The level of the alert. Valid values: `Critical`, `Warn`, `Info`.
* `arn` - (Optional) The Alibaba Cloud Resource Name (ARN) of the resource.
-> **NOTE:** Currently, the Alibaba Cloud Resource Name (ARN) of the resource. To use, please [submit an application](https://www.alibabacloud.com/help/en/cloudmonitor/latest/describemetricruletargets).

### `escalations`

The escalations supports the following:

* `critical` - (Optional) The critical level. See [`critical`](#escalations-critical) below.
* `info` - (Optional) The info level. See [`info`](#escalations-info) below.
* `warn` - (Optional) The warn level. See [`warn`](#escalations-warn) below.

### `escalations-critical`

The critical supports the following:

* `comparison_operator` - (Optional) The comparison operator of the threshold for critical-level alerts.
* `statistics` - (Optional) The statistical aggregation method for critical-level alerts.
* `threshold` - (Optional) The threshold for critical-level alerts.
* `times` - (Optional, Int) The consecutive number of times for which the metric value is measured before a critical-level alert is triggered.

### `escalations-info`

The info supports the following: 

* `comparison_operator` - (Optional) The comparison operator of the threshold for info-level alerts.
* `statistics` - (Optional) The statistical aggregation method for info-level alerts.
* `threshold` - (Optional) The threshold for info-level alerts.
* `times` - (Optional, Int) The consecutive number of times for which the metric value is measured before a info-level alert is triggered.

### `escalations-warn`

The warn supports the following:

* `comparison_operator` - (Optional) The comparison operator of the threshold for warn-level alerts.
* `statistics` - (Optional) The statistical aggregation method for warn-level alerts.
* `threshold` - (Optional) The threshold for warn-level alerts.
* `times` - (Optional, Intl) The consecutive number of times for which the metric value is measured before a warn-level alert is triggered.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Group Metric Rule. Its value is same as `rule_id`.
* `status` - The status of Group Metric Rule.

## Timeouts

-> **NOTE:** Available since v1.191.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Group Metric Rule.
* `update` - (Defaults to 3 mins) Used when update the Group Metric Rule.
* `delete` - (Defaults to 3 mins) Used when delete the Group Metric Rule.

## Import

Cloud Monitor Service Group Metric Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_group_metric_rule.example <rule_id>
```
