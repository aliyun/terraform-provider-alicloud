---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_metric_rule_template"
sidebar_current: "docs-alicloud-resource-cms-metric-rule-template"
description: |-
  Provides a Alicloud Cloud Monitor Service Metric Rule Template resource.
---

# alicloud_cms_metric_rule_template

Provides a Cloud Monitor Service Metric Rule Template resource.

For information about Cloud Monitor Service Metric Rule Template and how to use it, see [What is Metric Rule Template](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createmetricruletemplate).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_metric_rule_template&exampleId=9e32b0f7-c1d4-bcda-f01e-11ba1d4ce0313ea3559f&activeTab=example&spm=docs.r.cms_metric_rule_template.0.9e32b0f7c1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_cms_metric_rule_template" "example" {
  metric_rule_template_name = var.name
  alert_templates {
    category    = "ecs"
    metric_name = "cpu_total"
    namespace   = "acs_ecs_dashboard"
    rule_name   = "tf_example"
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_metric_rule_template&spm=docs.r.cms_metric_rule_template.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `metric_rule_template_name` - (Required, ForceNew) The name of the alert template.
* `description` - (Optional) The description of the alert template.
* `group_id` - (Optional) The ID of the application group.
* `apply_mode` - (Optional) The mode in which the alert template is applied. Valid values:
  - `GROUP_INSTANCE_FIRST`: The metrics in the application group take precedence.
  - `ALARM_TEMPLATE_FIRST `: The metrics specified in the alert template take precedence.
* `notify_level` - (Optional) The alert notification method. Valid values:
  - `2`: Alert notifications are sent by using Phone, SMS, Email, TradeManager and DingTalk chatbots.
  - `3 `: Alert notifications are sent by using SMS, Email, TradeManager and DingTalk chatbots.
  - `4`: Alert notifications are sent by using TradeManager and DingTalk chatbots.
* `silence_time` - (Optional, Int) The mute period during which notifications are not repeatedly sent for an alert. Unit: seconds. Default value: `86400`. Valid values: `0` to `86400`.
* `webhook` - (Optional) The callback URL to which a POST request is sent when an alert is triggered based on the alert rule.
* `enable_start_time` - (Optional) The beginning of the time period during which the alert rule is effective. Valid values: `00` to `23`. The value `00` indicates 00:00 and the value `23` indicates 23:00.
* `enable_end_time` - (Optional) The end of the time period during which the alert rule is effective. Valid values: `00` to `23`. The value `00` indicates 00:59 and the value `23` indicates 23:59.
* `alert_templates` - (Optional, Set) The details of alert rules that are generated based on the alert template. See [`alert_templates`](#alert_templates) below.

### `alert_templates`

The alert_templates supports the following: 

* `rule_name` - (Required) The name of the alert rule.
* `metric_name` - (Required) The name of the metric.
-> **NOTE:** For more information, see [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) or [Appendix 1: Metrics](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `namespace` - (Required) The namespace of the cloud service.
-> **NOTE:** For more information, see [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) or [Appendix 1: Metrics](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `category` - (Required) The abbreviation of the Alibaba Cloud service name.
-> **NOTE:** To obtain the abbreviation of an Alibaba Cloud service name, call the [DescribeProjectMeta](https://www.alibabacloud.com/help/en/cms/developer-reference/api-cms-2019-01-01-describeprojectmeta) operation. The metricCategory tag in the Labels response parameter indicates the abbreviation of the Alibaba Cloud service name.
* `webhook` - (Optional) The callback URL to which a POST request is sent when an alert is triggered based on the alert rule.
* `escalations` - (Optional, Set) The information about the trigger condition based on the alert level. See [`escalations`](#alert_templates-escalations) below. 

### `alert_templates-escalations`

The escalations supports the following: 

* `critical` - (Optional, Set) The condition for triggering critical-level alerts. See [`critical`](#alert_templates-escalations-critical) below. 
* `info` - (Optional, Set) The condition for triggering info-level alerts. See [`info`](#alert_templates-escalations-info) below. 
* `warn` - (Optional, Set) The condition for triggering warn-level alerts. See [`warn`](#alert_templates-escalations-warn) below. 

### `alert_templates-escalations-critical`

The critical supports the following:

* `comparison_operator` - (Optional) The comparison operator of the threshold for critical-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
* `statistics` - (Optional) The statistical aggregation method for critical-level alerts.
* `threshold` - (Optional) The threshold for critical-level alerts.
* `times` - (Optional) The consecutive number of times for which the metric value is measured before a critical-level alert is triggered.

### `alert_templates-escalations-info`

The info supports the following:

* `comparison_operator` - (Optional) The comparison operator of the threshold for info-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
* `statistics` - (Optional) The statistical aggregation method for info-level alerts.
* `threshold` - (Optional) The threshold for info-level alerts.
* `times` - (Optional) The consecutive number of times for which the metric value is measured before an info-level alert is triggered.

### `alert_templates-escalations-warn`

The warn supports the following: 

* `comparison_operator` - (Optional) The comparison operator of the threshold for warn-level alerts. Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
* `statistics` - (Optional) The statistical aggregation method for warn-level alerts.
* `threshold` - (Optional) The threshold for warn-level alerts.
* `times` - (Optional) The consecutive number of times for which the metric value is measured before a warn-level alert is triggered.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Metric Rule Template.
* `rest_version` - The version of the alert template.

## Import

Cloud Monitor Service Metric Rule Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_metric_rule_template.example <id>
```
