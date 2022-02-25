---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application_scaling_rules"
sidebar_current: "docs-alicloud-datasource-sae-application-scaling-rules"
description: |-
  Provides a list of Sae Application Scaling Rules to the user.
---

# alicloud\_sae\_application\_scaling\_rules

This data source provides the Sae Application Scaling Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.159.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_sae_application_scaling_rules" "ids" {
  app_id = "example_value"
  ids    = ["example_value-1", "example_value-2"]
}
output "sae_application_scaling_rule_id_1" {
  value = data.alicloud_sae_application_scaling_rules.ids.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Application Scaling Rule IDs.
* `app_id` - (Required) The ID of the Application.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - A list of Sae Application Scaling Rules. Each element contains the following attributes:
  * `app_id` - The ID of the Application.
  * `create_time` - The CreateTime of the Application Scaling Rule.
  * `id` - The ID of the Application Scaling Rule.
  * `scaling_rule_enable` - Whether to enable the auto scaling policy.
  * `scaling_rule_metric` - Monitoring indicators for elastic scaling.
    * `max_replicas` - The maximum number of instances.
    * `metrics` - The auto scaling list of monitoring indicators.
      * `metric_target_average_utilization` - The target value of the monitoring indicator.
      * `metric_type` - The metric type of the Application Scaling Rule.
    * `min_replicas` - The minimum number of instances.
    * `metrics_status` - Monitor indicator elasticity status.
      * `desired_replicas` - The number of target instances.
      * `next_scale_time_period` - The next cycle of monitoring indicator elasticity.
      * `current_replicas` - The number of current instances.
      * `last_scale_time` - The time of the last elastic expansion.
      * `max_replicas` - The maximum number of instances.
      * `min_replicas` - The minimum number of instances.
      * `current_metrics` - The current monitoring indicator elasticity list.
        * `type` - The metric type. Associated with monitoring indicators.
        * `current_value` - The current value.
        * `name` - The name of the trigger condition.
      * `next_scale_metrics` - Next monitoring indicator elasticity list
        * `next_scale_out_average_utilization` - The percentage value of the monitoring indicator elasticity that triggers the expansion condition next time.
        * `next_scale_in_average_utilization` - The percentage value of the monitoring indicator elasticity that triggers the shrinkage condition next time.
        * `name` - The name of the trigger condition.
    * `scale_down_rules` - The shrink rule.
      * `disabled` - Whether shrinkage is prohibited.
      * `stabilization_window_seconds` - Shrinkage cooling time.
      * `step` - Elastic shrinkage step. The maximum number of instances per unit time.
    * `scale_up_rules` - The expansion rules.
      * `step` - Flexible expansion step. The maximum number of instances per unit time.
      * `disabled` - Whether shrinkage is prohibited. The values are described as follows:
      * `stabilization_window_seconds` - Expansion cooling time.
  * `scaling_rule_name` - The name of the scaling rule.
  * `scaling_rule_timer` - Timing elastic expansion.
    * `begin_date` - The short-term start date of the timed elastic scaling strategy.
    * `end_date` - The short-term end date of the timed elastic scaling strategy.
    * `period` - The period in which a timed elastic scaling strategy is executed. 
    * `schedules` - Trigger point in time within a single day.
      * `at_time` - Time point. Format: `hours:minutes`.
      * `target_replicas` - The number of target instances.
      * `max_replicas` - The maximum number of instances.
      * `min_replicas` - The minimum number of instances.
  * `scaling_rule_type` - Flexible strategy type.