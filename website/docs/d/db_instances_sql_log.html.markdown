---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instances_sql_log"
sidebar_current: "docs-alicloud-datasource-db-instances-sql-log"
description: |-
    Provides a collection of RDS instances sql log records according to the specified filters.
---

# alicloud\_db\_instances_sql_log

The `alicloud_db_instances_sql_log` data source provides a collection of RDS instances sql log records available in Alibaba Cloud account.
Filters support start time, end time, page size, page number and so on.

## Example Usage

```
data "alicloud_db_instances_sql_log" "db_instances_sql_log" {
  db_instance_id = "rm-d324gfdhf5g"
  start_time     = "2020-01-04T15:00:00Z"
  end_time       = "2020-01-05T15:00:00Z"
}

output "first_db_instance_id" {
  value = "${data.alicloud_db_instances.db_instances_sql_log.sql_record.0.sql_text}"
}

```
E
## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required) The database instance id.
* `start_time` - (Required) Start time of the sql log records. The format is YYYY-MM-DDThh:mm:ssZ, such as 2011-05-30T12:11:4Z.
* `end_time` - (Required) End time of the the sql log records. The format is YYYY-MM-DDThh:mm:ssZ, such as 2011-05-30T12:11:4Z.
* `sql_id` - (Optional) The sql id of the sql log records.
* `query_key_words` - (Optional) The keywords for search sql log. 
* `data_base` - (Optional) The database name of the sql log records.
* `user` - (Optional) The user name of the sql log records.
* `page_size` - (Optional) Number of records per page. The number should between `30` and `100`. Default is `30`.
* `page_number` - (Optional) The page number. Default is `1`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `total_record_count` - The total record count. 
* `page_record_number` - The page record number. 
* `page_record_count` - The page record count. 
* `sql_record` - A list of RDS instances sql log record. Each element contains the following attributes:
  * `db_name` - The database name of sql log.
  * `account_name` - The account name of sql log.
  * `host_address` - The host address of sql log.
  * `sql_text` - The sql text of sql log.
  * `total_execution_times` - The total execution times of sql log.
  * `return_row_counts` - The return row counts of sql log.
  * `thread_id` - The thread id of sql log.
  * `execute_time` - The execution time of sql log.
