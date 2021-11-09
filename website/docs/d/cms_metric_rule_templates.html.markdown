---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_metric_rule_templates"
sidebar_current: "docs-alicloud-datasource-cms-metric-rule-templates"
description: |- 
    Provides a list of Cms Metric Rule Templates to the user.
---

# alicloud\_cms\_metric\_rule\_templates

This data source provides the Cms Metric Rule Templates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_metric_rule_templates" "ids" {
  ids = ["example_value"]
}
output "cms_metric_rule_template_id_1" {
  value = data.alicloud_cms_metric_rule_templates.ids.templates.0.id
}

data "alicloud_cms_metric_rule_templates" "nameRegex" {
  name_regex = "^my-MetricRuleTemplate"
}
output "cms_metric_rule_template_id_2" {
  value = data.alicloud_cms_metric_rule_templates.nameRegex.templates.0.id
}

data "alicloud_cms_metric_rule_templates" "keyword" {
  keyword = "^my-MetricRuleTemplate"
}
output "cms_metric_rule_template_id_3" {
  value = data.alicloud_cms_metric_rule_templates.nameRegex.templates.0.id
}

data "alicloud_cms_metric_rule_templates" "template_id" {
  template_id = "example_value"
}
output "cms_metric_rule_template_id_4" {
  value = data.alicloud_cms_metric_rule_templates.nameRegex.templates.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Valid values: `true` or `false`. Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Metric Rule Template IDs.
* `keyword` - (Optional, ForceNew) The name of the alert template. You can perform fuzzy search based on the template name.
* `metric_rule_template_name` - (Optional, ForceNew) The name of the alert template.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Metric Rule Template name.
* `template_id` - (Optional, ForceNew) The ID of the alert template.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Metric Rule Template names.
* `templates` - A list of Cms Metric Rule Templates. Each element contains the following attributes:
    * `alert_templates` - The details of alert rules that are generated based on the alert template.
        * `category` - The abbreviation of the service name. Valid values: `ecs`, `rds`, `ads`, `slb`, `vpc`, `apigateway`, `cdn`, `cs`, `dcdn`, `ddos`, `eip`, `elasticsearch`, `emr`, `ess`, `hbase`, `iot_edge`, `kvstore_sharding`, `kvstore_splitrw`, `kvstore_standard`, `memcache`, `mns`, `mongodb`, `mongodb_cluster`, `mongodb_sharding`, `mq_topic`, `ocs`, `opensearch`, `oss`, `polardb`, `petadata`, `scdn`, `sharebandwidthpackages`, `sls`, `vpn`.
        * `escalations` - The information about the trigger condition based on the alert level.
            * `critical` - The condition for triggering critical-level alerts.
                * `threshold` - The threshold for critical-level alerts.
                * `times` - The consecutive number of times for which the metric value is measured before a
                  critical-level alert is triggered.
                * `comparison_operator` - The comparison operator of the threshold for critical-level alerts.Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
                * `statistics` - The statistical aggregation method for critical-level alerts.
            * `info` - The condition for triggering info-level alerts.
                * `statistics` - The statistical aggregation method for info-level alerts.
                * `threshold` - The threshold for info-level alerts.
                * `times` - The consecutive number of times for which the metric value is measured before an info-level
                  alert is triggered.
                * `comparison_operator` - The comparison operator of the threshold for info-level alerts.Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
            * `warn` - The condition for triggering warn-level alerts.
                * `comparison_operator` - The comparison operator of the threshold for warn-level alerts.Valid values: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`.
                * `statistics` - The statistical aggregation method for warn-level alerts.
                * `threshold` - The threshold for warn-level alerts.
                * `times` - The consecutive number of times for which the metric value is measured before a warn-level
                  alert is triggered.
        * `metric_name` - The name of the metric.
        * `namespace` - The namespace of the service.
        * `rule_name` - The name of the alert rule.
        * `webhook` - The callback URL to which a POST request is sent when an alert is triggered based on the alert rule.
    * `description` - The description of the alert template.
    * `group_id` - GroupId.
    * `id` - The ID of the Metric Rule Template.
    * `metric_rule_template_name` - The name of the alert template.
    * `rest_version` - The version of the alert template.

	-> **NOTE:** The version changes with the number of times that the alert template is modified.
    * `template_id` - The ID of the alert template.
