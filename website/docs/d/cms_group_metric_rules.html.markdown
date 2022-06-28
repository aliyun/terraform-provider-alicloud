---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_group_metric_rules"
sidebar_current: "docs-alicloud-datasource-cms-group-metric-rules"
description: |-
  Provides a list of Cms Group Metric Rules to the user.
---

# alicloud\_cms\_group\_metric\_rules

This data source provides the Cms Group Metric Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.104.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_group_metric_rules" "example" {
  ids        = ["4a9a8978-a9cc-55ca-aa7c-530ccd91ae57"]
  name_regex = "the_resource_name"
}

output "first_cms_group_metric_rule_id" {
  value = data.alicloud_cms_group_metric_rules.example.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `dimensions` - (Optional, ForceNew) The dimensions that specify the resources to be associated with the alert rule.
* `enable_state` - (Optional, ForceNew) EnableState.
* `group_id` - (Optional, ForceNew) The ID of the application group.
* `group_metric_rule_name` - (Optional, ForceNew) The name of the alert rule.
* `ids` - (Optional, ForceNew, Computed)  A list of Group Metric Rule IDs.
* `metric_name` - (Optional, ForceNew) The name of the metric.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `namespace` - (Optional, ForceNew) The namespace of the service.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of Group Metric Rule.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Group Metric Rule names.
* `rules` - A list of Cms Group Metric Rules. Each element contains the following attributes:
	* `contact_groups` - Alarm contact group.
	* `dimensions` - The dimensions that specify the resources to be associated with the alert rule.
	* `effective_interval` - The time period during which the alert rule is effective.
	* `email_subject` - The subject of the alert notification email.
	* `enable_state` - Indicates whether the alert rule is enabled.
	* `escalations` - Alarm level.
        * `critical` - (Optional) The critical level.
            * `comparison_operator` - (Optional) The comparison operator of the threshold for critical-level alerts.                                         
            * `statistics` - (Optional) The statistical aggregation method for critical-level alerts.                                
            * `threshold` - (Optional) The threshold for critical-level alerts.
            * `times` - (Optional) The consecutive number of times for which the metric value is measured before a critical-level alert is triggered.  
        * `info` - (Optional) The info level.
            * `comparison_operator` - (Optional) The comparison operator of the threshold for info-level alerts.                                         
            * `statistics` - (Optional) The statistical aggregation method for info-level alerts.                                
            * `threshold` - (Optional) The threshold for info-level alerts.
            * `times` - (Optional) The consecutive number of times for which the metric value is measured before a info-level alert is triggered.
        * `warn` - (Optional) The warn level.
            * `comparison_operator` - (Optional) The comparison operator of the threshold for warn-level alerts.                                         
            * `statistics` - (Optional) The statistical aggregation method for warn-level alerts.                                
            * `threshold` - (Optional) The threshold for warn-level alerts.
            * `times` - (Optional) The consecutive number of times for which the metric value is measured before a warn-level alert is triggered.
	* `group_id` - The ID of the application group.
	* `group_metric_rule_name` - The name of the alert rule.
	* `id` - The ID of the Group Metric Rule.
	* `metric_name` - The name of the metric.
	* `namespace` - The namespace of the service.
	* `no_effective_interval` - The time period during which the alert rule is ineffective.
	* `period` - The aggregation period of the monitoring data. Unit: seconds. The value is an integral multiple of 60. Default value: `300`.
	* `resources` - The resources that are associated with the alert rule.
	* `rule_id` - The ID of the alert rule.
	* `silence_time` - The mute period during which new alerts are not reported even if the alert trigger conditions are met. Unit: seconds. Default value: `86400`, which is equivalent to one day.
	* `source_type` - The type of the alert rule. The value is fixed to METRIC, indicating an alert rule for time series metrics.
	* `status` - The status of Group Metric Rule..
	* `webhook` -  The callback URL.
