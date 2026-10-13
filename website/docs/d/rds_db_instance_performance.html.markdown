---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_instance_performance"
description: |-
  Provides a list of RDS instance performance metrics.
---

# alicloud_rds_db_instance_performance

This data source provides the performance metrics of an ApsaraDB RDS instance.

For information on the performance metrics and how to use them, please refer to [DescribeDBInstancePerformance](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-describedbinstanceperformance).

-> **NOTE:** Available since v1.295.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_db_instance_performance" "default" {
  db_instance_id = "rm-xxxxxxxxxxxx"
  key            = "MySQL_Sessions,MySQL_QPSTPS"
  start_time     = "2026-10-12T10:00Z"
  end_time       = "2026-10-12T11:00Z"
}

output "first_performance_key" {
  value = data.alicloud_rds_db_instance_performance.default.performance_keys.0.key
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required) The ID of the RDS instance.
* `key` - (Required) The performance metrics that you want to query. Separate multiple metrics with commas (,). You can query up to 30 metrics at a time. Valid values depend on the database engine of the instance, such as `MySQL_Sessions`, `MySQL_QPSTPS`, `MySQL_NetworkTraffic`, `MySQL_IOPS` and `MySQL_SpaceUsage` for MySQL instances. For more information, see [DescribeDBInstancePerformance](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-describedbinstanceperformance).

  -> **NOTE:** `MySQL_SpaceUsage` and `SQLServer_SpaceUsage` can be queried only for the last 1 day.
* `start_time` - (Required) The beginning of the time range to query. Specify the time in the `yyyy-MM-ddTHH:mmZ` format. The time must be in UTC.
* `end_time` - (Required) The end of the time range to query. Specify the time in the `yyyy-MM-ddTHH:mmZ` format. The time must be in UTC.

  -> **NOTE:** The interval between `start_time` and `end_time` must be greater than the monitoring collection frequency of the instance, otherwise an empty result may be returned.
* `node_id` - (Optional) The unique ID of the instance node. Example: `339****`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source.
* `engine` - The database engine of the instance.
* `performance_keys` - A list of performance metrics for the requested `key` list, rebuilt in the comma order of the request `key`. Every metric echoed by the API is kept, even when the specified time range holds no monitoring data: such a metric has an empty `values` list, and only a response without any performance key yields an empty `performance_keys` list. Matching is case-sensitive on the exact comma-separated `key` tokens: a requested token the API did not echo is not synthesized, and a metric outside the request is appended after the requested ones. Because the order does not depend on the API response, each metric keeps a stable index across refreshes. Each element contains the following attributes:
  * `key` - The name of the metric.
  * `unit` - The unit of the metric.
  * `value_format` - The format of the metric value.
  * `values` - A list of monitoring data of the metric, empty when the metric has no samples in the specified time range. Each element contains the following attributes:
    * `date` - The timestamp of the monitoring data. The time is in UTC.
    * `value` - The value of the metric at the specified time.
