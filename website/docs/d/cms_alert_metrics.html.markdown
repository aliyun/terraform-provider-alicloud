---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_metrics"
sidebar_current: "docs-alicloud-datasource-cms-alert-metrics"
description: |-
  Provides a list of Cms Alert Metric owned by an Alibaba Cloud account.
---

# alicloud_cms_alert_metrics

This data source provides the predefined alert metrics of CloudMonitor 2.0 for the current Alibaba Cloud user.

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
data "alicloud_cms_alert_metrics" "default" {
  include_details = true
}

output "first_alert_metric_id" {
  value = data.alicloud_cms_alert_metrics.default.metrics.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, List) A list of Alert Metric IDs.
* `group` - (Optional) The ID of the metric group, used to filter the metrics under the group.
* `include_details` - (Optional) Whether the returned objects contain the detailed attribute fields. Default value: `false`. If it is set to `true`, the returned results contain the `filters`, `params` and `expr_template` fields.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `metrics` - A list of Alert Metrics. Each element contains the following attributes:
  * `id` - The ID of the Alert Metric. It is the same as `alert_metric_id`.
  * `alert_metric_id` - The unique key of the alert metric.
  * `group` - The key of the group to which the metric belongs.
  * `alert_level` - The alert level that is pre-filled into the alert rule.
  * `alert_message_cn` - The Chinese message template of the alert, which can be pre-filled into the message field of the alert rule.
  * `alert_message_en` - The English message template of the alert, which can be pre-filled into the message field of the alert rule.
  * `display_name_cn` - The Chinese description of the metric.
  * `display_name_en` - The English description of the metric.
  * `display_statement_cn` - The Chinese condition description of the alert rule. The params referenced in it are rendered as user-fillable fields.
  * `display_statement_en` - The English condition description of the alert rule. The params referenced in it are rendered as user-fillable fields.
  * `duration` - The data duration that is pre-filled into the alert rule. Unit: seconds.
  * `unit_cn` - The Chinese unit of the metric.
  * `unit_en` - The English unit of the metric.
  * `ext_info` - The extended properties related to business scenarios.
  * `expr_template` - The query expression template. It is available when `include_details` is set to `true`.
    * `tpl` - The expression template string.
    * `type` - The type of the expression template. Valid values: `SIMPLE_EXPR_TPL`, `LOOP_JOIN_EXPR_TPL`.
  * `filters` - The list of filter conditions, including the filter conditions automatically inherited from the group to which the metric belongs. It is available when `include_details` is set to `true`.
    * `dim` - The dimension of the filter condition.
    * `display_name_cn` - The Chinese display name of the filter condition.
    * `display_name_en` - The English display name of the filter condition.
    * `hidden` - Indicates whether the filter condition is hidden. Hidden filter conditions do not need to be displayed.
    * `opt` - The default operator.
    * `supported_opts` - The list of optional operators.
      * `display_name_cn` - The Chinese name of the optional operator.
      * `display_name_en` - The English name of the optional operator.
      * `value` - The value of the operator.
  * `params` - The list of parameter configurations of the metric. The metric automatically inherits the parameter definitions of the group. It is available when `include_details` is set to `true`.
    * `max_width` - The maximum width of the input box that the front end displays for parameters of the input type.
    * `min_width` - The minimum width of the input box that the front end displays for parameters of the input type.
    * `name` - The name of the parameter.
    * `type` - The type of the parameter. Valid values: `TEXT_PARAM` (read-only text parameter defined by the backend, no user input control is displayed), `INPUT_PARAM` (input box parameter), `SELECT_PARAM` (select box parameter).
    * `value` - The default value of the parameter.
    * `values` - The list of optional values of the parameter. It is valid only for `SELECT_PARAM`.
      * `label_cn` - The Chinese display name of the option.
      * `label_en` - The English display name of the option.
      * `value` - The value of the option.
