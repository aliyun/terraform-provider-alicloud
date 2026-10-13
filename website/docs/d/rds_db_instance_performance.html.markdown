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

-> **NOTE:** Available since v2.1.1.

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

* `db_instance_id` - (Required, String, Available since v2.1.1) The ID of the RDS instance.
* `key` - (Required, String, Available since v2.1.1) The performance metrics that you want to query. Separate multiple metrics with commas (,). Whitespace around each metric is trimmed, empty entries are ignored and duplicate metrics are collapsed. You must provide at least one and at most 30 distinct metrics; otherwise the plan fails before any API call is sent. Valid values depend on the database engine of the instance, such as `MySQL_Sessions`, `MySQL_QPSTPS`, `MySQL_NetworkTraffic`, `MySQL_IOPS` and `MySQL_SpaceUsage` for MySQL instances. For more information, see [DescribeDBInstancePerformance](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-describedbinstanceperformance).

  -> **NOTE:** `MySQL_SpaceUsage` and `SQLServer_SpaceUsage` can be queried only for the last 1 day.
* `start_time` - (Required, String, Available since v2.1.1) The beginning of the time range to query. Specify the time in the `yyyy-MM-ddTHH:mmZ` format. The time must be in UTC.
* `end_time` - (Required, String, Available since v2.1.1) The end of the time range to query. Specify the time in the `yyyy-MM-ddTHH:mmZ` format. The time must be in UTC.

  -> **NOTE:** The interval between `start_time` and `end_time` must be greater than the monitoring collection frequency of the instance, otherwise an empty result may be returned.
* `node_id` - (Optional, String, Available since v2.1.1) The unique ID of the instance node. Example: `339****`.
* `output_file` - (Optional, String, Available since v2.1.1) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Read-Only, String, Available since v2.1.1) The ID of the data source.
* `engine` - (Read-Only, String, Available since v2.1.1) The database engine of the instance.
* `performance_keys` - (Read-Only, List, Available since v2.1.1) A list of performance metrics rebuilt on the provider side from the requested `key`, so that each requested metric keeps a stable index across refreshes when the requested metrics remain unchanged. Each distinct token of the request `key` appears in the comma order of the request, and a requested metric that the API did not echo is kept in place with an empty `unit`, an empty `value_format` and an empty `values` list. Matching is case-sensitive on the exact comma-separated `key` tokens; duplicate request tokens are collapsed and empty tokens are ignored. Metrics returned outside the request are appended after the requested ones, sorted by `key` in ascending lexicographical (case-sensitive) order, instead of being dropped. Changes to these extra metrics can change their indices. Because `key` must contain at least one metric, the list is never empty for a valid request. Each element contains the following attributes:
  * `key` - (Read-Only, String, Available since v2.1.1) The name of the metric.
  * `unit` - (Read-Only, String, Available since v2.1.1) The unit of the metric.
  * `value_format` - (Read-Only, String, Available since v2.1.1) The format of the metric value.
  * `values` - (Read-Only, List, Available since v2.1.1) A list of monitoring data of the metric, empty when the metric has no samples in the specified time range. Each element contains the following attributes:
    * `date` - (Read-Only, String, Available since v2.1.1) The timestamp of the monitoring data. The time is in UTC.
    * `value` - (Read-Only, String, Available since v2.1.1) The value of the metric at the specified time.
