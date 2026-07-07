---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_metric_alarm_rule"
description: |-
  Provides a Alicloud Cloud Monitor Service Metric Alarm Rule resource.
---

# alicloud_cloud_monitor_service_metric_alarm_rule

Provides a Cloud Monitor Service Metric Alarm Rule resource.

Describes the time-series metric alarm rule configured by the user.

For information about Cloud Monitor Service Metric Alarm Rule and how to use it, see [What is Metric Alarm Rule](https://next.api.alibabacloud.com/document/Cms/2019-01-01/PutResourceMetricRule).

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_monitor_service_metric_alarm_rule" "default" {
  status               = true
  send_ok              = true
  contact_groups       = "云账号报警联系人"
  silence_time         = "86400"
  metric_alarm_rule_id = "SystemDefault_acs_drds_IOPSUsageOfDo"
  period               = "60"
  effective_interval   = "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7"
  no_data_policy       = "KEEP_LAST_STATE"
  namespace            = "acs_drds"
  metric_name          = "IOPSUsageOfDN"
  escalations {
    critical {
      comparison_operator = "GreaterThanThreshold"
      times               = "5"
      statistics          = "Average"
      threshold           = "80"
    }
  }
  resources = jsonencode([{ resource = "_ALL" }])
  rule_name = "SystemDefault_acs_drds_IOPSUsageOfDN"
}
```

## Argument Reference

The following arguments are supported:
* `composite_expression` - (Optional, Set) Alert condition for multiple metrics.

-> **NOTE:**  Single-metric and multi-metric conditions are mutually exclusive and cannot be configured simultaneously.
 See [`composite_expression`](#composite_expression) below.
* `contact_groups` - (Required, ForceNew) Alarm contact groups. Alarm notifications are sent to the contacts in these groups.

-> **NOTE:**  An alarm contact group is a collection of one or more alarm contacts. For information about how to create alarm contacts and alarm contact groups, see [PutContact](https://help.aliyun.com/document_detail/114923.html) and [PutContactGroup](https://help.aliyun.com/document_detail/114929.html).

* `effective_interval` - (Optional) The time range during which the alert rule is effective.
* `email_subject` - (Optional) Subject of alert emails.
* `escalations` - (Optional, Set) The trigger conditions for alert levels. See [`escalations`](#escalations) below.
* `interval` - (Optional) The trigger interval of the alarm rule. Unit: seconds.

-> **NOTE:** For information about how to query the statistical period of a metric, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).


-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `labels` - (Optional, Set) When a metric meets the alert condition and an alert is triggered, the labels are written to the metric and displayed in the alert notification.

-> **NOTE:**  This feature is equivalent to the Label in Prometheus alerts.
 See [`labels`](#labels) below.
* `metric_alarm_rule_id` - (Required, ForceNew) The ID of the alarm rule.

You can specify a new alarm rule ID or use an existing alarm rule ID from CloudMonitor. For information about how to query alarm rule IDs, see [DescribeMetricRuleList](https://help.aliyun.com/document_detail/114941.html).

-> **NOTE:**  Specifying a new alarm rule ID creates a threshold-based alarm rule.

* `metric_name` - (Required, ForceNew) The name of the metric. For information about how to query metric names, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).

-> **NOTE:**  When you create a Prometheus alert rule for Enterprise Cloud Monitoring, this parameter specifies the metric store name. For information about how to obtain the metric store name, see [DescribeHybridMonitorNamespaceList](https://help.aliyun.com/document_detail/428880.html).

* `namespace` - (Required, ForceNew) The namespace of the cloud service metric data. For information about how to query the namespace of a cloud service, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).

-> **NOTE:**  When you create a Prometheus alert rule for Enterprise Cloud Monitoring, this parameter must be set to `acs_prometheus`.

* `no_data_policy` - (Optional, Computed) The policy to apply when no monitoring data is available. Valid values:
  - KEEP_LAST_STATE (default): No action is taken.
  - INSUFFICIENT_DATA: The alert status is set to "Insufficient Data."
  - OK: The alert status is set to "OK."
* `no_effective_interval` - (Optional) The time range during which the alarm rule is inactive.
* `period` - (Optional) The statistical period of the metric. Unit: seconds. By default, this is the original reporting period of the metric.

-> **NOTE:**  For information about how to query the statistical period of a metric, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).

* `prometheus` - (Optional, Set) Prometheus alert.

-> **NOTE:**  You must specify this parameter only when you create a Prometheus alert rule for Enterprise Cloud Monitor.
 See [`prometheus`](#prometheus) below.
* `resources` - (Required) Resource information, for example: `[{"instanceId":"i-uf6j91r34rnwawoo****"}]`, `[{"userId":"100931896542****"}]`.
For dimensions supported in resource information, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).
* `rule_name` - (Required, ForceNew) Alert rule name.

You can enter a new alert rule name or use an existing alert rule name in CloudMonitor. For information about how to query alert rule names, see [DescribeMetricRuleList](https://help.aliyun.com/document_detail/114941.html).

-> **NOTE:**  Entering a new alert rule name creates a threshold-based alert rule.

* `send_ok` - (Optional, Computed) Specifies whether to send recovery notifications.
* `silence_time` - (Optional) Channel silence period. Unit: seconds. Default value: 86400.

-> **NOTE:**  The channel silence period specifies how long to wait before resending an alarm notification if the alarm condition has not returned to normal after the initial alarm.

* `status` - (Optional, Computed) The enabled status of the alarm rule. Valid values:
  - true: enabled.
  - false: disabled.
* `webhook` - (Optional) The URL address specified for callback when an alert is triggered. A POST request is sent to this URL.

### `composite_expression`

The composite_expression supports the following:
* `expression_list` - (Optional, Set) A list of alert conditions created using standard expressions. See [`expression_list`](#composite_expression-expression_list) below.
* `expression_list_join` - (Optional) The logical relationship between multiple metric-based alert conditions. Valid values:
  - `&&`: An alert is triggered only when all metrics meet their respective alert conditions. That is, an alert is triggered only when every expression in ExpressionList evaluates to `true`.
  - `||`: An alert is triggered as soon as any one metric meets its alert condition.
* `expression_raw` - (Optional) The alert condition created by an expression. This includes, but is not limited to, the following scenarios:
  - Configure an alert blacklist for specific resources. For example: `$instanceId != 'i-io8kfvcpp7x5****' && $Average > 50` means that even if the `Average` metric of instance `i-io8kfvcpp7x5****` in the alert rule exceeds 50, no alert will be triggered.
  - Set a special alert threshold for a specified instance in the rule. For example: `$Average > ($instanceId == 'i-io8kfvcpp7x5****' ? 80 : 50)` means that an alert is triggered only when the `Average` metric of instance `i-io8kfvcpp7x5****` exceeds 80, while for other instances, an alert is triggered when their `Average` exceeds 50.
  - Limit the number of instances exceeding the threshold in the rule. For example: `count($Average > 20) > 3` means that an alert is triggered only when more than three instances in the alert rule have an `Average` metric greater than 20.
* `level` - (Optional) The alert severity level. Valid values:
  - CRITICAL: Critical.
  - WARN: Warning.
  - INFO: Information.
* `times` - (Optional, Int) Number of consecutive times the alert condition must be met before an alert notification is sent.

### `composite_expression-expression_list`

The composite_expression-expression_list supports the following:
* `comparison_operator` - (Optional) Threshold comparison operator. Valid values:
  - GreaterThanOrEqualToThreshold: greater than or equal to
  - GreaterThanThreshold: greater than
  - LessThanOrEqualToThreshold: less than or equal to
  - LessThanThreshold: less than
  - NotEqualToThreshold: not equal to
  - GreaterThanYesterday: higher than the same time yesterday
  - LessThanYesterday: lower than the same time yesterday
  - GreaterThanLastWeek: higher than the same time last week
  - LessThanLastWeek: lower than the same time last week
  - GreaterThanLastPeriod: higher than the previous period
  - LessThanLastPeriod: lower than the previous period
* `metric_name` - (Optional) The name of the cloud service metric.
* `period` - (Optional, Int) Aggregation period of the metric.
Unit: seconds.
* `statistics` - (Optional) Statistical method of the metric. Valid values:
  - $Maximum: maximum value.
  - $Minimum: minimum value.
  - $Average: average value.
  - $Availability: availability rate (typically used for site monitoring).

-> **NOTE:**  `$` is the unified prefix symbol for metrics. For supported cloud services, see [Cloud Service Metrics](https://help.aliyun.com/document_detail/163515.html).

* `threshold` - (Optional) Alert threshold.

### `escalations`

The escalations supports the following:
* `critical` - (Optional, Set) The trigger condition for Critical-level alerts. See [`critical`](#escalations-critical) below.
* `info` - (Optional, Set) Trigger conditions for Info-level alerts. See [`info`](#escalations-info) below.
* `warn` - (Optional, Set) Trigger condition for Warn-level alerts.   See [`warn`](#escalations-warn) below.

### `escalations-critical`

The escalations-critical supports the following:
* `comparison_operator` - (Optional) The comparison operator for the Critical-level threshold. Valid values:
  - GreaterThanOrEqualToThreshold: greater than or equal to
  - GreaterThanThreshold: greater than
  - LessThanOrEqualToThreshold: less than or equal to
  - LessThanThreshold: less than
  - NotEqualToThreshold: not equal to
  - GreaterThanYesterday: increased compared to the same time yesterday
  - LessThanYesterday: decreased compared to the same time yesterday
  - GreaterThanLastWeek: increased compared to the same time last week
  - LessThanLastWeek: decreased compared to the same time last week
  - GreaterThanLastPeriod: increased compared to the previous period
  - LessThanLastPeriod: decreased compared to the previous period
* `statistics` - (Optional) Statistical method for Critical-level alerts.
* `threshold` - (Optional) Threshold for Critical-level alerts.
* `times` - (Optional, Int) The number of consecutive occurrences required to trigger a Critical-level alarm. An alarm is triggered only when the metric exceeds the threshold for this specified number of consecutive times.

### `escalations-info`

The escalations-info supports the following:
* `comparison_operator` - (Optional) Comparison operator for Info-level thresholds. Valid values:
  - GreaterThanOrEqualToThreshold: Greater than or equal to
  - GreaterThanThreshold: Greater than
  - LessThanOrEqualToThreshold: Less than or equal to
  - LessThanThreshold: Less than
  - NotEqualToThreshold: Not equal to
  - GreaterThanYesterday: Increased compared to the same time yesterday
  - LessThanYesterday: Decreased compared to the same time yesterday
  - GreaterThanLastWeek: Increased compared to the same time last week
  - LessThanLastWeek: Decreased compared to the same time last week
  - GreaterThanLastPeriod: Increased compared to the previous period
  - LessThanLastPeriod: Decreased compared to the previous period
* `statistics` - (Optional) Statistical method used for Info-level alerts.
* `threshold` - (Optional) Threshold value for Info-level alerts.
* `times` - (Optional, Int) Number of consecutive occurrences required to trigger an Info-level alert. The alert is triggered only when the metric exceeds the threshold for this specified number of consecutive times.

### `escalations-warn`

The escalations-warn supports the following:
* `comparison_operator` - (Optional) Comparison operator for the Warn-level threshold. Valid values:  
  - GreaterThanOrEqualToThreshold: greater than or equal to  
  - GreaterThanThreshold: greater than  
  - LessThanOrEqualToThreshold: less than or equal to  
  - LessThanThreshold: less than  
  - NotEqualToThreshold: not equal to  
  - GreaterThanYesterday: higher than the same time yesterday  
  - LessThanYesterday: lower than the same time yesterday  
  - GreaterThanLastWeek: higher than the same time last week  
  - LessThanLastWeek: lower than the same time last week  
  - GreaterThanLastPeriod: higher than the previous period  
  - LessThanLastPeriod: lower than the previous period  
* `statistics` - (Optional) Statistical method for Warn-level alerts.
* `threshold` - (Optional) Threshold for Warn-level alerts.
* `times` - (Optional, Int) Number of consecutive occurrences required to trigger a Warn-level alert. The alert is triggered only if the metric exceeds the threshold for this number of consecutive times.  

### `labels`

The labels supports the following:
* `key` - (Optional) The tag key.
* `value` - (Optional) Label value.

-> **NOTE:**  The label value supports template parameters, which are replaced with actual label values.


### `prometheus`

The prometheus supports the following:
* `annotations` - (Optional, Set) When a Prometheus alert is triggered, the key-value pairs of annotations are rendered to help you better understand the metric or alert rule.

-> **NOTE:**  This feature is equivalent to Annotations in Prometheus.
 See [`annotations`](#prometheus-annotations) below.
* `level` - (Optional) Alert severity level. Valid values:
  - CRITICAL: Critical
  - WARN: Warning
  - INFO: Information
* `prom_ql` - (Optional) The PromQL query statement.

-> **NOTE:**  The data retrieved by the PromQL query statement is used as alert data. Include the alert threshold in this statement.

* `times` - (Optional, Int) The number of times the alert condition must be met before an alert notification is sent.

### `prometheus-annotations`

The prometheus-annotations supports the following:
* `key` - (Optional) The key of the annotation.
* `value` - (Optional) The value of the annotation.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `dimensions` - The monitoring dimensions for the specified resource.
* `escalations` - The trigger conditions for alert levels.
  * `critical` - The trigger condition for Critical-level alerts.
    * `pre_condition` - The precondition for triggering a Critical-level alarm.
  * `info` - Trigger conditions for Info-level alerts.
    * `pre_condition` - Precondition for triggering an Info-level alert.
  * `warn` - Trigger condition for Warn-level alerts.
    * `pre_condition` - Precondition for triggering a Warn-level alert.
* `source_type` - The type of the alarm rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Metric Alarm Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Metric Alarm Rule.
* `update` - (Defaults to 5 mins) Used when update the Metric Alarm Rule.

## Import

Cloud Monitor Service Metric Alarm Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_monitor_service_metric_alarm_rule.example <metric_alarm_rule_id>
```