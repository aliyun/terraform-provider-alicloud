---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_scheduled_sql"
description: |-
  Provides a Alicloud Log Service (SLS) Scheduled Sql resource.
---

# alicloud_sls_scheduled_sql

Provides a Log Service (SLS) Scheduled Sql resource.

Scheduled SQL task.

For information about Log Service (SLS) Scheduled Sql and how to use it, see [What is Scheduled Sql](https://www.alibabacloud.com/help/zh/sls/developer-reference/api-sls-2020-12-30-createscheduledsql).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "defaultKIe4KV" {
  description  = "${var.name}-${random_integer.default.result}"
  project_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_log_store" "default1LI9we" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.defaultKIe4KV.project_name
  logstore_name    = "${var.name}-${random_integer.default.result}"
}


resource "alicloud_sls_scheduled_sql" "default" {
  description = "example-tf-scheduled-sql-0006"
  schedule {
    type            = "Cron"
    time_zone       = "+0700"
    delay           = "20"
    cron_expression = "0 0/1 * * *"
  }
  display_name = "example-tf-scheduled-sql-0006"
  scheduled_sql_configuration {
    script                  = "* | select * from log"
    sql_type                = "searchQuery"
    dest_endpoint           = "ap-northeast-1.log.aliyuncs.com"
    dest_project            = "job-e2e-project-jj78kur-ap-southeast-1"
    source_logstore         = alicloud_log_store.default1LI9we.logstore_name
    dest_logstore           = "example-open-api02"
    role_arn                = "acs:ram::1395894005868720:role/aliyunlogetlrole"
    dest_role_arn           = "acs:ram::1395894005868720:role/aliyunlogetlrole"
    from_time_expr          = "@m-1m"
    to_time_expr            = "@m"
    max_run_time_in_seconds = "1800"
    resource_pool           = "enhanced"
    max_retries             = "5"
    from_time               = "1713196800"
    to_time                 = "0"
    data_format             = "log2log"
  }
  scheduled_sql_name = var.name
  project            = alicloud_log_project.defaultKIe4KV.project_name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Job description.
* `display_name` - (Required) Task display name.
* `project` - (Required, ForceNew) A short description of struct.
* `schedule` - (Required, ForceNew, Set) Schedule type. This field generally does not need to be specified. If you have strict scheduling requirements—for example, running an import job every Monday at 8:00 AM—you can use a cron expression. See [`schedule`](#schedule) below.
* `scheduled_sql_configuration` - (Required, ForceNew, Set) Task configuration. See [`scheduled_sql_configuration`](#scheduled_sql_configuration) below.
* `scheduled_sql_name` - (Required, ForceNew) The job name. The naming rules are as follows:
  - Job names must be unique within the same project.
  - The name can contain only lowercase letters, digits, hyphens (-), and underscores (_).
  - The name must start and end with a lowercase letter or digit.
  - The length must be between 2 and 64 characters.
* `status` - (Optional, Computed, Available since v1.271.0) The status of the scheduled SQL job.

### `schedule`

The schedule supports the following:
* `cron_expression` - (Optional) Cron expression with a minimum precision of minutes in 24-hour format. For example, 0 0/1 * * * means checking once every hour starting from 00:00. When type is set to Cron, cronExpression must be specified.
* `delay` - (Optional, Int) Delay duration.
* `interval` - (Optional) Time interval, such as 5m or 1h.
* `run_immediately` - (Optional) Specifies whether to run the OSS import job immediately after it is created.
* `time_zone` - (Optional) Time zone.
* `type` - (Optional) The check frequency type. Log Service checks query and analysis results based on the frequency you configure. Valid values:
FixedRate: Checks query and analysis results at fixed intervals.
Cron: Uses a cron expression to specify the interval and checks query and analysis results accordingly.
Weekly: Checks query and analysis results once at a fixed time on a specific day of the week.
Daily: Checks query and analysis results once at a fixed time each day.
Hourly: Checks query and analysis results once every hour.

### `scheduled_sql_configuration`

The scheduled_sql_configuration supports the following:
* `data_format` - (Optional, ForceNew) Write mode.  
* `dest_endpoint` - (Optional) The destination endpoint.
* `dest_logstore` - (Optional) The destination Logstore.
* `dest_project` - (Optional) The destination project.
* `dest_role_arn` - (Optional) Destination write role ARN.  
* `from_time` - (Optional, ForceNew, Int) The start time of the schedule.
* `from_time_expr` - (Optional) SQL time window - start.  
* `max_retries` - (Optional, Int) Maximum number of retries.
* `max_run_time_in_seconds` - (Optional, Int) SQL timeout.  
* `parameters` - (Optional, Map) Parameter configuration.
* `resource_pool` - (Optional) Resource pool.  
* `role_arn` - (Optional) Source read role ARN.  
* `script` - (Optional) SQL statement.
* `source_logstore` - (Optional, ForceNew) The source Logstore.
* `sql_type` - (Optional) SQL type.
* `to_time` - (Optional, ForceNew, Int) Scheduled end time.  
* `to_time_expr` - (Optional) End of the SQL time window.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<project>:<scheduled_sql_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Scheduled Sql.
* `delete` - (Defaults to 5 mins) Used when delete the Scheduled Sql.
* `update` - (Defaults to 5 mins) Used when update the Scheduled Sql.

## Import

Log Service (SLS) Scheduled Sql can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_scheduled_sql.example <project>:<scheduled_sql_name>
```