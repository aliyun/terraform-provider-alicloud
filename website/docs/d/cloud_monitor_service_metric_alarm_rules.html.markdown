---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_metric_alarm_rules"
sidebar_current: "docs-alicloud-datasource-cloud-monitor-service-metric-alarm-rules"
description: |-
  Provides a list of Cloud Monitor Service Metric Alarm Rule owned by an Alibaba Cloud account.
---

# alicloud_cloud_monitor_service_metric_alarm_rules

This data source provides Cloud Monitor Service Metric Alarm Rule available to the user.[What is Metric Alarm Rule](https://next.api.alibabacloud.com/document/Cms/2019-01-01/PutResourceMetricRule)

-> **NOTE:** Available since v1.283.0.

## Example Usage

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
* `dimensions` - (ForceNew, Optional) The monitoring dimensions for the specified resource.
Format: a set of key:value pairs, for example: `{"userId":"120886317861****"}` and `{"instanceId":"i-2ze2d6j5uhg20x47****"}`.
* `metric_alarm_rule_id` - (ForceNew, Optional) The ID of the alarm rule.

You can specify a new alarm rule ID or use an existing alarm rule ID from CloudMonitor. For information about how to query alarm rule IDs, see [DescribeMetricRuleList](https://help.aliyun.com/document_detail/114941.html).

-> **NOTE:**  Specifying a new alarm rule ID creates a threshold-based alarm rule.

* `metric_name` - (ForceNew, Optional) The name of the metric. For information about how to query metric names, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).

-> **NOTE:**  When you create a Prometheus alert rule for Enterprise Cloud Monitoring, this parameter specifies the metric store name. For information about how to obtain the metric store name, see [DescribeHybridMonitorNamespaceList](https://help.aliyun.com/document_detail/428880.html).

* `namespace` - (ForceNew, Optional) The namespace of the cloud service metric data. For information about how to query the namespace of a cloud service, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).

-> **NOTE:**  When you create a Prometheus alert rule for Enterprise Cloud Monitoring, this parameter must be set to `acs_prometheus`.

* `rule_name` - (ForceNew, Optional) Alert rule name.

You can enter a new alert rule name or use an existing alert rule name in CloudMonitor. For information about how to query alert rule names, see [DescribeMetricRuleList](https://help.aliyun.com/document_detail/114941.html).

-> **NOTE:**  Entering a new alert rule name creates a threshold-based alert rule.

* `status` - (ForceNew, Optional) The enabled status of the alarm rule. Valid values:
  - true: enabled.
  - false: disabled.
* `ids` - (Optional, ForceNew, Computed) A list of Metric Alarm Rule IDs. 
* `enable_details` - (Optional, ForceNew) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Metric Alarm Rule IDs.
* `rules` - A list of Metric Alarm Rule Entries. Each element contains the following attributes:
    * `composite_expression` - Alert condition for multiple metrics.
        * `expression_list` - A list of alert conditions created using standard expressions.
            * `comparison_operator` - Threshold comparison operator.
            * `metric_name` - The name of the cloud service metric.
            * `period` - Aggregation period of the metric.
            * `statistics` - Statistical method of the metric.
            * `threshold` - Alert threshold.
        * `expression_list_join` - The logical relationship between multiple metric-based alert conditions.
        * `expression_raw` - The alert condition created by an expression.
        * `level` - The alert severity level.
        * `times` - Number of consecutive times the alert condition must be met before an alert notification is sent.
    * `contact_groups` - Alarm contact groups.
    * `dimensions` - The monitoring dimensions for the specified resource.
    * `effective_interval` - The time range during which the alert rule is effective.
    * `email_subject` - Subject of alert emails.
    * `escalations` - The trigger conditions for alert levels.
        * `critical` - The trigger condition for Critical-level alerts.
            * `comparison_operator` - The comparison operator for the Critical-level threshold.
            * `pre_condition` - The precondition for triggering a Critical-level alarm.
            * `statistics` - Statistical method for Critical-level alerts.
            * `threshold` - Threshold for Critical-level alerts.
            * `times` - The number of consecutive occurrences required to trigger a Critical-level alarm.
        * `info` - Trigger conditions for Info-level alerts.
            * `comparison_operator` - Comparison operator for Info-level thresholds.
            * `pre_condition` - Precondition for triggering an Info-level alert.
            * `statistics` - Statistical method used for Info-level alerts.
            * `threshold` - Threshold value for Info-level alerts.
            * `times` - Number of consecutive occurrences required to trigger an Info-level alert.
        * `warn` - Trigger condition for Warn-level alerts.
            * `comparison_operator` - Comparison operator for the Warn-level threshold.
            * `pre_condition` - Precondition for triggering a Warn-level alert.
            * `statistics` - Statistical method for Warn-level alerts.
            * `threshold` - Threshold for Warn-level alerts.
            * `times` - Number of consecutive occurrences required to trigger a Warn-level alert.
    * `labels` - When a metric meets the alert condition and an alert is triggered, the labels are written to the metric and displayed in the alert notification.
        * `key` - The tag key.
        * `value` - Label value.
    * `metric_alarm_rule_id` - The ID of the alarm rule.
    * `metric_name` - The name of the metric.
    * `namespace` - The namespace of the cloud service metric data.
    * `no_data_policy` - The policy to apply when no monitoring data is available.
    * `no_effective_interval` - The time range during which the alarm rule is inactive.
    * `period` - The statistical period of the metric.
    * `prometheus` - Prometheus alert.
        * `annotations` - When a Prometheus alert is triggered, the key-value pairs of annotations are rendered to help you better understand the metric or alert rule.
            * `key` - The key of the annotation.
            * `value` - The value of the annotation.
        * `level` - Alert severity level.
        * `prom_ql` - The PromQL query statement.
        * `times` - The number of times the alert condition must be met before an alert notification is sent.
    * `resources` - Resource information, for example: `[{"instanceId":"i-uf6j91r34rnwawoo****"}]`, `[{"userId":"100931896542****"}]`.
    * `rule_name` - Alert rule name.
    * `send_ok` - Specifies whether to send recovery notifications.
    * `silence_time` - Channel silence period.
    * `source_type` - The type of the alarm rule.
    * `status` - The enabled status of the alarm rule.
    * `webhook` - The URL address specified for callback when an alert is triggered.
    * `targets` - The alert callback targets returned only when `enable_details` is set to `true`.
      * `arn` - Resource ARN of the alert callback target (typically an MNS queue or topic).
      * `json_params` - JSON-formatted parameters carried in the alert callback.
      * `level` - Alert severity level. Valid values: `INFO`, `WARN`, `CRITICAL`.
      * `target_id` - Alert trigger target ID.
    * `id` - The ID of the resource supplied above.
