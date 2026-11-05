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

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_metric_list" "example" {
  metric_name      = "CPUUtilization"
  metric_namespace = "acs_ecs_dashboard"
  dimensions       = "{\"instanceId\":\"i-bp1example\"}"
  start_time       = "1785542400000"
  end_time         = "1785546000000"
  period           = "60"
}

output "first_datapoint_average" {
  value = try(data.alicloud_cms_metric_list.example.datapoints[0].average, null)
}
```

## Argument Reference

The following arguments are supported:

* `metric_name` - (Required, String) The name of the monitoring metric.
* `metric_namespace` - (Required, String) The namespace of the monitoring metric. For example, `acs_ecs_dashboard` for ECS and `acs_rds_dashboard` for RDS.
* `dimensions` - (Optional, String) The dimensions that classify the monitoring metric. The value is a JSON string that contains key-value pairs. For example, `{"instanceId":"i-bp1example"}`.
* `start_time` - (Optional, String) The beginning of the time range to query. Specify the time in the `YYYY-MM-DD HH:mm:ss` format or as a Unix timestamp in milliseconds. For an unambiguous UTC-based query window, use a Unix timestamp in milliseconds, measured from 00:00:00 UTC on January 1, 1970. This value must be earlier than the end time.
* `end_time` - (Optional, String) The end of the time range to query. Specify the time in the `YYYY-MM-DD HH:mm:ss` format or as a Unix timestamp in milliseconds. For an unambiguous UTC-based query window, use a Unix timestamp in milliseconds, measured from 00:00:00 UTC on January 1, 1970.
* `period` - (Optional, String) The interval at which monitoring data is collected. Unit: seconds. Valid values: `15`, `60`, `900`, `3600`.
* `express` - (Optional, String) The aggregation method that is used to display the monitoring data. The value is a JSON string. Only the `groupby` field is supported. For example, `{"groupby":["userId","instanceId"]}`.
* `length` - (Optional, Int) The maximum number of entries that are returned in a single query. The maximum value is `1440`. Defaults to `1000`.
* `next_token` - (Optional, String) The pagination token that is used in the next request to retrieve more datapoints.
* `output_file` - (Optional, String) File name where to save data source results (after running `terraform plan`). The file is overwritten after every successful read, including when the metric does not exist (an empty result is written as `[]`); the read fails if the file cannot be written.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (String) The ID of the data source. It is computed from the metric name and namespace.
* `actual_period` - (String) The actual period at which monitoring data is collected. This value is returned by the server.
* `next_token` - (String) The pagination token that can be used in the next request to retrieve more datapoints.
* `datapoints` - (List) A list of Cms Metric Datapoints. Each element contains the following attributes:
  * `uuid` - (String) The unique identifier of the datapoint.
  * `metric_name` - (String) The name of the monitoring metric.
  * `timestamp` - (String) The timestamp that indicates when the metric value was collected.
  * `average` - (Float) The average of the metric values within the statistical period.
  * `sum` - (Float) The sum of the metric values within the statistical period.
  * `maximum` - (Float) The maximum metric value within the statistical period.
  * `minimum` - (Float) The minimum metric value within the statistical period.
  * `count` - (Float) The number of times that the metric value is collected within the statistical period.
  * `user_id` - (String) The ID of the Alibaba Cloud user.
  * `instance_id` - (String) The ID of the instance.
