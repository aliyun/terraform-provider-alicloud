---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_metric_list"
sidebar_current: "docs-alicloud-datasource-cms-metric-list"
description: |-
  Provides the datapoints of a specified Cloud Monitor metric to the user.
---

# alicloud_cms_metric_list

This data source provides the datapoints of a specified Cloud Monitor Service (CMS) metric within a time range. It wraps the `DescribeMetricList` API (Cms/2019-01-01) and paginates by `NextToken`.

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_metric_list" "example" {
  metric_name      = "CPUUtilization"
  metric_namespace = "acs_ecs_dashboard"
  dimensions       = "{\"instanceId\":\"i-bp1example\"}"
  start_time       = "2026-08-01 00:00:00"
  end_time         = "2026-08-01 01:00:00"
  period           = "60"
}

output "first_datapoint_average" {
  value = data.alicloud_cms_metric_list.example.datapoints.0.average
}
```

## Argument Reference

The following arguments are supported:

* `metric_name` - (Required) The name of the monitoring metric.
* `metric_namespace` - (Required) The namespace of the monitoring metric. For example, `acs_ecs_dashboard` for ECS and `acs_rds_dashboard` for RDS.
* `dimensions` - (Optional) The dimensions that classify the monitoring metric. The value is a JSON string that contains key-value pairs. For example, `{"instanceId":"i-bp1example"}`.
* `start_time` - (Optional) The beginning of the time range to query. Specify the time in the `YYYY-MM-DD HH:mm:ss` format or as a Unix timestamp in milliseconds. This value must be earlier than the end time.
* `end_time` - (Optional) The end of the time range to query. Specify the time in the `YYYY-MM-DD HH:mm:ss` format or as a Unix timestamp in milliseconds.
* `period` - (Optional) The interval at which monitoring data is collected. Unit: seconds. Valid values: `60`, `300`, `900`.
* `express` - (Optional) The aggregation method that is used to display the monitoring data. The value is a JSON string. For example, `{"downsample":"avg"}`.
* `length` - (Optional) The maximum number of entries that are returned in a single query. If you do not specify this parameter, the default value `1000` is used.
* `next_token` - (Optional) The pagination token that is used in the next request to retrieve more datapoints.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source. It is computed from the metric name and namespace.
* `datapoints` - A list of Cms Metric Datapoints. Each element contains the following attributes:
  * `uuid` - The unique identifier of the datapoint.
  * `metric_name` - The name of the monitoring metric.
  * `timestamp` - The timestamp that indicates when the metric value was collected.
  * `average` - The average of the metric values within the statistical period.
  * `sum` - The sum of the metric values within the statistical period.
  * `maximum` - The maximum metric value within the statistical period.
  * `minimum` - The minimum metric value within the statistical period.
  * `count` - The number of times that the metric value is collected within the statistical period.
  * `user_id` - The ID of the Alibaba Cloud user.
  * `instance_id` - The ID of the instance.
* `actual_period` - The actual period at which monitoring data is collected. This value is returned by the server.
* `next_token` - The pagination token that can be used in the next request to retrieve more datapoints.
