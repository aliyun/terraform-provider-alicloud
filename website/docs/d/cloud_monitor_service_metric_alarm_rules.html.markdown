---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_metric_alarm_rules"
description: |-
  Provides a list of Cloud Monitor Service Metric Alarm Rules to the user.
---

# alicloud_cloud_monitor_service_metric_alarm_rules

This data source provides the Cloud Monitor Service Metric Alarm Rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.256.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_instances" "default" {
  status = "Running"
}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_cms_alarm" "default" {
  name               = var.name
  project            = "acs_ecs_dashboard"
  metric             = "disk_writebytes"
  period             = 900
  silence_time       = 300
  webhook            = "https://www.aliyun.com"
  enabled            = true
  contact_groups     = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  effective_interval = "06:00-20:00"
  metric_dimensions  = <<EOF
  [
    {
      "instanceId": "${data.alicloud_instances.default.ids.0}",
      "device": "/dev/vda1"
    }
  ]
  EOF
  escalations_critical {
    statistics          = "Average"
    comparison_operator = "<="
    threshold           = 90
    times               = 1
  }
  escalations_info {
    statistics          = "Minimum"
    comparison_operator = "!="
    threshold           = 20
    times               = 3
  }
  escalations_warn {
    statistics          = "Average"
    comparison_operator = "=="
    threshold           = 30
    times               = 5
  }
}

data "alicloud_cloud_monitor_service_metric_alarm_rules" "ids" {
  ids = [alicloud_cms_alarm.default.id]
}

output "cloud_monitor_service_metric_alarm_rules_id_0" {
  value = data.alicloud_cloud_monitor_service_metric_alarm_rules.ids.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Metric Alarm Rule IDs.
* `dimensions` - (Optional, ForceNew) The monitoring dimensions of the specified resource.
* `metric_name` - (Optional, ForceNew) The name of the metric.
* `namespace` - (Optional, ForceNew) The namespace of the cloud service.
* `rule_name` - (Optional, ForceNew) The name of the alert rule.
* `status` - (Optional, ForceNew) Specifies whether to query enabled or disabled alert rules. Valid values: `true`, `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - A list of Hybrid Double Writes. Each element contains the following attributes:
  * `id` - The ID of the alert rule.
  * `contact_groups` - The alert contact group.
  * `dimensions` - The dimensions of the alert rule.
  * `effective_interval` - The time period during which the alert rule is effective.
  * `email_subject` - The subject of the alert notification email.
  * `metric_name` - The name of the metric.
  * `namespace` - The namespace of the cloud service.
  * `no_data_policy` - The method that is used to handle alerts when no monitoring data is found.
  * `no_effective_interval` - The time period during which the alert rule is ineffective.
  * `period` - The statistical period.
  * `resources` - The resources that are associated with the alert rule.
  * `rule_name` - The name of the alert rule.
  * `silence_time` - The mute period during which new alert notifications are not sent even if the trigger conditions are met.
  * `source_type` - The type of the alert rule.
  * `status` - Indicates whether the alert rule is enabled.
  * `webhook` - The callback URL.
  * `composite_expression` - The trigger conditions for multiple metrics.
    * `times` - The number of consecutive triggers.
    * `expression_raw` - The trigger conditions that are created by using expressions.
    * `expression_list_join` - The relationship between the trigger conditions for multiple metrics.
    * `level` - The alert level.
    * `expression_list` - The trigger conditions that are created in standard mode.
      * `metric_name` - The metric that is used to monitor the cloud service.
      * `comparison_operator` - The operator that is used to compare the metric value with the threshold.
      * `period` - The aggregation period of the metric.
      * `statistics` - The statistical method of the metric.
      * `threshold` - The alert threshold.
  * `escalations` - The conditions for triggering different levels of alerts.
    * `critical` - The conditions for triggering Critical-level alerts.
      * `comparison_operator` - The comparison operator that is used to compare the metric value with the threshold.
      * `times` - The consecutive number of times for which the metric value meets the alert condition before a Critical-level alert is triggered.
      * `pre_condition` - The additional conditions for triggering Critical-level alerts.
      * `statistics` - The statistical methods for Critical-level alerts.
      * `threshold` - The threshold for Critical-level alerts.
    * `info` - The conditions for triggering Info-level alerts.
      * `comparison_operator` - The comparison operator that is used to compare the metric value with the threshold.
      * `times` - The consecutive number of times for which the metric value meets the alert condition before a Info-level alert is triggered.
      * `pre_condition` - The additional conditions for triggering Info-level alerts.
      * `statistics` - The statistical methods for Info-level alerts.
      * `threshold` - The threshold for Info-level alerts.
    * `warn` - The conditions for triggering Warn-level alerts.
      * `comparison_operator` - The comparison operator that is used to compare the metric value with the threshold.
      * `times` - The consecutive number of times for which the metric value meets the alert condition before a Warn-level alert is triggered.
      * `pre_condition` - The additional conditions for triggering Warn-level alerts.
      * `statistics` - The statistical methods for Warn-level alerts.
      * `threshold` - The threshold for Warn-level alerts. 
  * `labels` - The tags of the alert rule.
    * `value` - The tag value of the alert rule.
    * `key` - The tag key of the alert rule.
  * `prometheus` - The Prometheus alerts.
    * `prom_ql` - The PromQL query statement.
    * `times` - The number of consecutive triggers.
    * `level` - The alert level.
    * `annotations` - The annotations of the Prometheus alert rule.
      * `value` - The value of the annotation.
      * `key` - The subject of the alert notificaThe key of the annotation.
