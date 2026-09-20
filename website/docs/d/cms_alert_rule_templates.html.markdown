---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_rule_templates"
sidebar_current: "docs-alicloud-datasource-cms-alert-rule-templates"
description: |-
  Provides a list of Cms Alert Rule Templates to the user.
---

# alicloud\_cms\_alert\_rule\_templates

This data source provides the Cms Alert Rule Templates of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alert_rule_templates" "ids" {
  biz_source = "CI"
  ids        = ["example_value"]
}
output "cms_alert_rule_template_id_1" {
  value = data.alicloud_cms_alert_rule_templates.ids.templates.0.id
}

data "alicloud_cms_alert_rule_templates" "nameRegex" {
  biz_source = "CI"
  name_regex = "^my-AlertRuleTemplate"
}
output "cms_alert_rule_template_id_2" {
  value = data.alicloud_cms_alert_rule_templates.nameRegex.templates.0.id
}

data "alicloud_cms_alert_rule_templates" "detailed" {
  biz_source     = "CI"
  enable_details = true
}
output "cms_alert_rule_template_id_3" {
  value = data.alicloud_cms_alert_rule_templates.detailed.templates.0.id
}
```

## Argument Reference

The following arguments are supported:

* `biz_source` - (Required) The business source (monitoring type) of the alert rule template. Valid values: `CI` (CloudInsight), `CMS_ENT` (Enterprise Cloud Monitoring). The backing `ListAlertRuleTemplates` API requires this parameter for every request.
* `enable_details` - (Optional) Valid values: `true` or `false`. Default to `false`. Set it to `true` can output more details about resource attributes. When set to `false`, only `alert_rule_template_id`, `display_name_cn`, `display_name_en`, `description_cn` and `description_en` are returned for each template.
* `ids` - (Optional, Computed) A list of Alert Rule Template IDs.
* `name_regex` - (Optional) A regex string to filter results by the English display name of the alert rule template.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Alert Rule Template English display names.
* `templates` - A list of Cms Alert Rule Templates. Each element contains the following attributes:
  * `alert_rule_template_id` - The ID of the alert rule template.
  * `biz_source` - The business source of the alert rule template.
  * `description_cn` - The Chinese description of the alert rule template.
  * `description_en` - The English description of the alert rule template.
  * `display_name_cn` - The Chinese display name of the alert rule template.
  * `display_name_en` - The English display name of the alert rule template.
  * `id` - The ID of the Alert Rule Template. It is the same as `alert_rule_template_id`.
  * `interval` - The interval of the alert rule template.
  * `labels` - The labels of the alert rule template.
  * `level` - The alert level of the alert rule template.
  * `message_cn` - The Chinese alert message of the alert rule template.
  * `message_en` - The English alert message of the alert rule template.
  * `condition` - The condition configuration of the alert rule template.
    * `alert_count` - The number of times the condition is met before an alert is triggered. Valid only when `type` is `SLS_CONDITION`.
    * `case_list` - The list of alert condition branches. Valid only when `type` is `SLS_CONDITION`.
      * `condition` - The matching expression of the condition branch.
      * `count_condition` - The count matching expression of the condition branch.
      * `level` - The alert level when the condition branch is met. Valid values: `CRITICAL`, `ERROR`, `WARNING`, `INFO`.
      * `type` - The matching type. Valid values: `HasData`, `HasDataCount`, `HasDataMatch`, `HasDataMatchCount`.
    * `compare_list` - The threshold comparison condition list. Valid only for `APM_CONDITION`.
      * `aggregate` - The post-aggregation function. Valid values: `count`, `sum`, `avg`, `min`, `max`, `p90`, `p95`, `p99`.
      * `oper` - The comparison operator. Valid values: `GT`, `GTE`, `LT`, `LTE`, `EQ`, `NE`, `YOY_UP`, `YOY_DOWN`.
      * `threshold` - The threshold value.
      * `value_level_list` - The threshold and level list.
        * `level` - The alert level.
        * `value` - The threshold value.
      * `yoy_time_unit` - The year-on-year time unit. Valid only when `oper` is `YOY_UP` or `YOY_DOWN`. Valid values: `minute`, `hour`, `day`, `week`, `month`.
      * `yoy_time_value` - The year-on-year time value.
    * `no_data_append_value` - The compensation value when there is no data. Valid only for `APM_CONDITION`.
    * `nodata_alert_level` - The alert level when there is no data.
    * `type` - The condition type. Valid values: `SLS_CONDITION`, `APM_CONDITION`.
  * `datasource` - The datasource configuration of the alert rule template.
    * `ds_list` - The sub datasource list. Required when `type` is `SLS_MULTI_DS`.
      * `project` - The SLS project name.
      * `region_id` - The region ID of the SLS project.
      * `store` - The name of the log store or metric store.
    * `instance_id` - The instance ID. Required when `type` is `PROMETHEUS_DS` or `ENTERPRISE_DS`.
    * `namespace` - The name of the Enterprise Cloud Monitoring metric repository. Valid only when `type` is `ENTERPRISE_DS`.
    * `type` - The datasource type. Valid values: `ENTERPRISE_DS`, `PROMETHEUS_DS`, `SLS_MULTI_DS`.
  * `query` - The query configuration of the alert rule template.
    * `duration` - The alert data duration in seconds. Valid when `type` is `PROMQL_QUERY`.
    * `expr` - The query expression. Required when `type` is `PROMQL_QUERY`.
    * `first_join` - The first collection operation of query results. Valid only when `type` is `SLS_MULTI_QUERY`.
      * `conditions` - The collection join condition list.
        * `first_field` - The left parameter of the join condition.
        * `oper` - The join operator.
        * `second_field` - The right parameter of the join condition.
      * `type` - The collection operation type.
    * `group_field_list` - The grouping field list.
    * `group_type` - The result grouping type. Valid values: `none`, `label`, `custom`.
    * `queries` - The sub query list. Required when `type` is `SLS_MULTI_QUERY`.
      * `duration` - The alert data duration in seconds.
      * `end` - The relative time offset end.
      * `expr` - The query expression.
      * `start` - The relative time offset start.
      * `time_unit` - The time unit of start, end and window. Valid values: `second`, `minute`, `hour`, `day`.
      * `window` - The query time window.
    * `second_join` - The second collection operation of query results. Valid only when `type` is `SLS_MULTI_QUERY`.
      * `conditions` - The collection join condition list.
        * `first_field` - The left parameter of the join condition.
        * `oper` - The join operator.
        * `second_field` - The right parameter of the join condition.
      * `type` - The collection operation type.
    * `type` - The query type. Valid values: `PROMQL_QUERY`, `SLS_MULTI_QUERY`, `APM_MULTI_QUERY`.
