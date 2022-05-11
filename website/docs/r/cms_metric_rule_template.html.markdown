---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_metric_rule_template"
sidebar_current: "docs-alicloud-resource-cms-metric-rule-template"
description: |-
  Provides a Alicloud Cloud Monitor Service Metric Rule Template resource.
---

# alicloud\_cms\_metric\_rule\_template

Provides a Cloud Monitor Service Metric Rule Template resource.

For information about Cloud Monitor Service Metric Rule Template and how to use it, see [What is Metric Rule Template](https://www.alibabacloud.com/help/doc-detail/114984.html).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_metric_rule_template" "example" {
  metric_rule_template_name = "example_value"
  alert_templates {
    category    = "ecs"
    metric_name = "cpu_total"
    namespace   = "acs_ecs_dashboard"
    rule_name   = "tf_testAcc_new"
    escalations {
      critical {
        comparison_operator = "GreaterThanThreshold"
        statistics          = "Average"
        threshold           = "90"
        times               = "3"
      }
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `alert_templates` - (Optional) The details of alert rules that are generated based on the alert template. See the following `Block alert_templates`.
* `apply_mode` - (Optional) The mode in which the alert template is applied. Valid values:`GROUP_INSTANCE_FIRST`or `ALARM_TEMPLATE_FIRST`. GROUP_INSTANCE_FIRST: The metrics in the application group take precedence. If a metric specified in the alert template does not exist in the application group, the system does not generate an alert rule for the metric based on the alert template. ALARM_TEMPLATE_FIRST: The metrics specified in the alert template take precedence. If a metric specified in the alert template does not exist in the application group, the system still generates an alert rule for the metric based on the alert template.
* `description` - (Optional) The description of the alert template.
* `enable_end_time` - (Optional) The end of the time period during which the alert rule is effective. Valid values: 00 to 23. The value 00 indicates 00:59 and the value 23 indicates 23:59.
* `enable_start_time` - (Optional) The beginning of the time period during which the alert rule is effective. Valid values: 00 to 23. The value 00 indicates 00:00 and the value 23 indicates 23:00.
* `group_id` - (Optional) The ID of the application group.
* `metric_rule_template_name` - (Required) The name of the alert template.
* `notify_level` - (Optional) The alert notification method. Valid values:Set the value to 4. The value 4 indicates that alert notifications are sent by using TradeManager and DingTalk chatbots.
* `rest_version` - (Optional) The version of the alert template to be modified.

-> **NOTE:** The version changes with the number of times that the alert template is modified.
* `silence_time` - (Optional) The mute period during which notifications are not repeatedly sent for an alert.Valid values: 0 to 86400. Unit: seconds. Default value: `86400`.

-> **NOTE:** Only one alert notification is sent during each mute period even if the metric value exceeds the alert threshold several times.
* `webhook` - (Optional) The callback URL to which a POST request is sent when an alert is triggered based on the alert rule.

#### Block alert_templates

The alert_templates supports the following: 

* `category` - (Required) The abbreviation of the service name. Valid values: `ecs`, `rds`, `ads`, `slb`, `vpc`, `apigateway`, `cdn`, `cs`, `dcdn`, `ddos`, `eip`, `elasticsearch`, `emr`, `ess`, `hbase`, `iot_edge`, `kvstore_sharding`, `kvstore_splitrw`, `kvstore_standard`, `memcache`, `mns`, `mongodb`, `mongodb_cluster`, `mongodb_sharding`, `mq_topic`, `ocs`, `opensearch`, `oss`, `polardb`, `petadata`, `scdn`, `sharebandwidthpackages`, `sls`, `vpn`.
* `escalations` - (Optional) The information about the trigger condition based on the alert level. See the following `Block escalations`.
* `metric_name` - (Required) The name of the metric.

-> **NOTE:** For more information, see [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) or [Appendix 1: Metrics](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `namespace` - (Required) The namespace of the service.

-> **NOTE:** For more information, see [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) or [Appendix 1: Metrics](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `rule_name` - (Required) The name of the alert rule.
* `webhook` - (Optional) The callback URL to which a POST request is sent when an alert is triggered based on the alert rule.

#### Block escalations

The escalations supports the following: 

* `critical` - (Optional) The condition for triggering critical-level alerts. See the following `Block critical`.
* `info` - (Optional) The condition for triggering info-level alerts. See the following `Block info`.
* `warn` - (Optional) The condition for triggering warn-level alerts. See the following `Block warn`.

#### Block warn

The warn supports the following: 

* `comparison_operator` - (Optional) The comparison operator of the threshold for warn-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
* `statistics` - (Optional) The statistical aggregation method for warn-level alerts.
* `threshold` - (Optional) The threshold for warn-level alerts.
* `times` - (Optional) The consecutive number of times for which the metric value is measured before a warn-level alert is triggered.

#### Block info

The info supports the following: 

* `comparison_operator` - (Optional) The comparison operator of the threshold for info-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
* `statistics` - (Optional) The statistical aggregation method for info-level alerts.
* `threshold` - (Optional) The threshold for info-level alerts.
* `times` - (Optional) The consecutive number of times for which the metric value is measured before an info-level alert is triggered.

#### Block critical

The critical supports the following: 

* `comparison_operator` - (Optional) The comparison operator of the threshold for critical-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
* `statistics` - (Optional) The statistical aggregation method for critical-level alerts.
* `threshold` - (Optional) The threshold for critical-level alerts.
* `times` - (Optional) The consecutive number of times for which the metric value is measured before a critical-level alert is triggered.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Metric Rule Template.

## Import

Cloud Monitor Service Metric Rule Template can be imported using the id, e.g.

```
$ terraform import alicloud_cms_metric_rule_template.example <id>
```
