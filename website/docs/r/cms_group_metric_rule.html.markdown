---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_group_metric_rule"
sidebar_current: "docs-alicloud-resource-cms-group-metric-rule"
description: |-
  Provides a Alicloud Cloud Monitor Service Group Metric Rule resource.
---

# alicloud\_cms\_group\_metric\_rule

Provides a Cloud Monitor Service Group Metric Rule resource.

For information about Cloud Monitor Service Group Metric Rule and how to use it, see [What is Group Metric Rule](https://www.alibabacloud.com/help/en/doc-detail/114943.htm).

-> **NOTE:** Available in v1.104.0+.

## Example Usage

Basic Usage

```terraform
resource "random_uuid" "this" {}

resource "alicloud_cms_group_metric_rule" "this" {
  group_id = "539****"
  rule_id  = random_uuid.this.id

  category    = "ecs"
  namespace   = "acs_ecs_dashboard"
  metric_name = "cpu_total"
  period      = "60"

  group_metric_rule_name = "tf-testacc-rule-name"
  email_subject          = "tf-testacc-rule-name-warning"
  interval               = "3600"
  silence_time           = 85800
  no_effective_interval  = "00:00-05:30"
  webhook                = "http://www.aliyun.com"
  escalations {
    warn {
      comparison_operator = "GreaterThanOrEqualToThreshold"
      statistics          = "Average"
      threshold           = "90"
      times               = 3
    }
    info {
      comparison_operator = "LessThanLastWeek"
      statistics          = "Average"
      threshold           = "90"
      times               = 5
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Required, ForceNew) The abbreviation of the service name. 
* `contact_groups` - (Optional, Computed) Alarm contact group.
* `dimensions` - (Optional, Computed) The dimensions that specify the resources to be associated with the alert rule.
* `effective_interval` - (Optional) The time period during which the alert rule is effective.
* `email_subject` - (Optional, Computed) The subject of the alert notification email.                                         .
* `escalations` - (Required) Alarm level. See the block for escalations.
* `group_id` - (Required) The ID of the application group.
* `group_metric_rule_name` - (Required) The name of the alert rule.                                      
* `interval` - (Optional, ForceNew) The interval at which Cloud Monitor checks whether the alert rule is triggered. Unit: seconds.                                    
* `metric_name` - (Required) The name of the metric.
* `namespace` - (Required, ForceNew) The namespace of the service.
* `no_effective_interval` - (Optional) The time period during which the alert rule is ineffective.                                       
* `period` - (Optional) The aggregation period of the monitoring data. Unit: seconds. The value is an integral multiple of 60. Default value: `300`.                       
* `rule_id` - (Required, ForceNew) The ID of the alert rule.
* `silence_time` - (Optional) The mute period during which new alerts are not reported even if the alert trigger conditions are met. Unit: seconds. Default value: `86400`, which is equivalent to one day.
* `webhook` - (Optional) The callback URL.                        

#### Block escalations

The escalations supports the following: 

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

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Group Metric Rule. Value as `rule_id`.
* `status` - The status of Group Metric Rule.

## Import

Cloud Monitor Service Group Metric Rule can be imported using the id, e.g.

```
$ terraform import alicloud_cms_group_metric_rule.example <rule_id>
```
