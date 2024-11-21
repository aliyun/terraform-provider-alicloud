---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_scheduled_sql"
description: |-
  Provides a Alicloud SLS Scheduled SQL resource.
---

# alicloud_sls_scheduled_sql

Provides a SLS Scheduled SQL resource. Scheduled SQL task.

For information about SLS Scheduled SQL and how to use it, see [What is Scheduled SQL](https://www.alibabacloud.com/help/zh/sls/developer-reference/api-sls-2020-12-30-createscheduledsql).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_scheduled_sql&exampleId=253dde24-19e2-e251-f325-1334d318ade82b5d4233&activeTab=example&spm=docs.r.sls_scheduled_sql.0.253dde2419&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `description` - (Optional) Task Description.
* `display_name` - (Required) Task Display Name.
* `project` - (Required, ForceNew) Log project.
* `schedule` - (Required, ForceNew) The scheduling type is generally not required by default. If there is a strong timing requirement, if it must be imported every Monday at 8 o'clock, cron can be used. See [`schedule`](#schedule) below.
* `scheduled_sql_configuration` - (Required, ForceNew) Task Configuration. See [`scheduled_sql_configuration`](#scheduled_sql_configuration) below.
* `scheduled_sql_name` - (Required, ForceNew) Timed SQL name.

### `schedule`

The schedule supports the following:
* `cron_expression` - (Optional) Cron expression, minimum precision is minutes, 24-hour clock. For example, 0 0/1 **indicates that the check is performed every one hour from 00:00. When type is set to Cron, cronExpression must be set.
* `delay` - (Optional) Delay time.
* `interval` - (Optional) Time interval, such as 5m, 1H.
* `run_immediately` - (Optional) Whether to execute the OSS import task immediately after it is created.
* `time_zone` - (Optional) Time Zone.
* `type` - (Optional) Check the frequency type. Log Service checks the query and analysis results based on the frequency you configured. The value is as follows: FixedRate: checks the query and analysis results at fixed intervals. Cron: specifies a time interval through a Cron expression, and checks the query and analysis results at the specified time interval. Weekly: Check the query and analysis results at a fixed point in time on the day of the week. Daily: checks the query and analysis results at a fixed time point every day. Hourly: Check query and analysis results every hour.

### `scheduled_sql_configuration`

The scheduled_sql_configuration supports the following:
* `data_format` - (Optional, ForceNew) Write Mode.
* `dest_endpoint` - (Optional) Target Endpoint.
* `dest_logstore` - (Optional) Target Logstore.
* `dest_project` - (Optional) Target Project.
* `dest_role_arn` - (Optional) Write target role ARN.
* `from_time` - (Optional, ForceNew) Schedule Start Time.
* `from_time_expr` - (Optional) SQL time window-start.
* `max_retries` - (Optional) Maximum retries.
* `max_run_time_in_seconds` - (Optional) SQL timeout.
* `parameters` - (Optional, Map) Parameter configuration.
* `resource_pool` - (Optional) Resource Pool.
* `role_arn` - (Optional) Read role ARN.
* `script` - (Optional) SQL statement.
* `source_logstore` - (Optional, ForceNew) Source Logstore.
* `sql_type` - (Optional) SQL type.
* `to_time` - (Optional, ForceNew) Time at end of schedule.
* `to_time_expr` - (Optional) SQL time window-end.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project>:<scheduled_sql_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Scheduled SQL.
* `delete` - (Defaults to 5 mins) Used when delete the Scheduled SQL.
* `update` - (Defaults to 5 mins) Used when update the Scheduled SQL.

## Import

SLS Scheduled SQL can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_scheduled_sql.example <project>:<scheduled_sql_name>
```