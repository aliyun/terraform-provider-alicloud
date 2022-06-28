---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_instance_performance"
sidebar_current: "docs-alicloud-datasource-rds-db_instance_performance"
description: |-
  Provides a list of Rds DB Instance Performance to the user.
---

# alicloud\_rds\_db\_instance\_performance

This data source provides the Rds DB Instance Performance of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.174.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_db_instance_performance" "example" {
  db_instance_id = "example_value"
  start_time     = "2022-06-23T17:00Z"
  end_time       = "2022-06-24T17:40Z"
  key            = "MemoryUsage,CpuUsage"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The db instance id.
* `end_time` - (Required, ForceNew) The end time.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `start_time` - (Required, ForceNew) The start time.
* `key` - (Required, ForceNew) The performance metric that you want to query. If you enter more than one performance metric, separate them with commas (,). For more information, see [Performance parameter table](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/performance-parameter-table).

-> **NOTE**: If you set the Key parameter to MySQL_SpaceUsage or SQLServer_SpaceUsage, you can query the performance metric only over one day.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `performance_keys` - A list of Rds DB Instance Performance. Each element contains the following attributes:
  * `key` - The name of the performance metric.
  * `unit` - The unit of the performance metric.
  * `value_format` - The format that the value of the performance metric follows. Multiple performance metric values are separated by the &amp; string. Example: com_delete&amp;com_insert&amp;com_insert_select&amp;com_replace.
  * `parameter_name` - The name of the parameter.
  * `values` - An array that consists of performance metric values.
    * `value` - The value of the performance metric.
    * `date` - The date and time when the value of the performance metric was recorded. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mmZ format. The time is displayed in UTC.