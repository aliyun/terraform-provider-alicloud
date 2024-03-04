---
subcategory: "SLS"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_oss_ingestion"
description: |-
  Provides a Alicloud SLS Oss Ingestion resource.
---

# alicloud_sls_oss_ingestion

Provides a SLS Oss Ingestion resource. OSS ingestion job configuration.

For information about SLS Oss Ingestion and how to use it, see [What is Oss Ingestion](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.218.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "defaultBucket" {
  bucket_name = var.name

  storage_class = "Standard"
}

resource "alicloud_log_project" "defaultProject" {
  name = var.name

}

resource "alicloud_log_store" "defaultILogStore" {
  retention_period = "30"
  project          = alicloud_log_project.defaultProject.name
  name             = var.name

}


resource "alicloud_sls_oss_ingestion" "default" {
  project     = alicloud_log_project.defaultProject.name
  description = "tf-testacc"
  configuration {
    logstore = alicloud_log_store.defaultILogStore.name
    source {
      endpoint = "oss-cn-chengdu.aliyuncs.com"
      bucket   = alicloud_oss_bucket.defaultBucket.bucket_name
      encoding = "utf-16"
      format {
      }
      interval          = "3m"
      pattern           = "1day*"
      prefix            = "1day"
      start_time        = "1706792358"
      end_time          = "1706792558"
      time_field        = "__time__"
      time_format       = "epoch"
      time_pattern      = "\\d+:\\d+:\\d+"
      time_zone         = "GMT-09:00"
      use_meta_index    = "false"
      compression_codec = "none"
      role_arn          = "acs:log:asdasd:asdas:dasdasdasd"
    }
  }
  oss_ingestion_name = var.name

  schedule {
    type            = "Resident"
    interval        = "5m"
    cron_expression = "*/3 * * * * *"
    time_zone       = "GMT+8:00"
    delay           = "5m"
  }
  display_name = "tf-testacc-442"
}
```

## Argument Reference

The following arguments are supported:
* `configuration` - (Required, ForceNew) Task Configuration. See [`configuration`](#configuration) below.
* `description` - (Optional) Task Description.
* `display_name` - (Required) Task Display Name.
* `oss_ingestion_name` - (Required, ForceNew) Job Name.
* `project` - (Required, ForceNew) Log project.
* `schedule` - (Optional, ForceNew) Check the frequency-dependent configuration. See [`schedule`](#schedule) below.

### `configuration`

The configuration supports the following:
* `logstore` - (Optional) Logstore name.
* `source` - (Optional, ForceNew) OSS import source configuration. See [`source`](#configuration-source) below.

### `configuration-source`

The configuration-source supports the following:
* `bucket` - (Optional) Oss bucket.
* `compression_codec` - (Optional) Compression type.
* `encoding` - (Optional) Encoding type.
* `end_time` - (Optional) Import files modified before a certain point in time.
* `endpoint` - (Optional) Oss endpoint.
* `format` - (Optional, Map) Data Format.
* `interval` - (Optional) Check New File Cycle.
* `pattern` - (Optional) File Path Regular Filtering.
* `prefix` - (Optional) File path prefix filtering.
* `restore_object_enabled` - (Optional) Import Archive.
* `role_arn` - (Optional) RoleArn for importing oss data.
* `start_time` - (Optional) Import files that have been modified since a point in time.
* `time_field` - (Optional) Extract Time Field.
* `time_format` - (Optional) Time Field Format.
* `time_pattern` - (Optional) Extraction time regular.
* `time_zone` - (Optional) Time Field Time Zone.
* `use_meta_index` - (Optional) Use the OSS metadata index.

### `schedule`

The schedule supports the following:
* `cron_expression` - (Optional) Cron expression, minimum precision is minutes, 24-hour clock. For example, 0 0/1 **indicates that the check is performed every one hour from 00:00. When type is set to Cron, cronExpression must be set.
* `delay` - (Optional) Delay time.
* `interval` - (Optional) Time interval, such as 5m, 1H.
* `run_immediately` - (Optional) Whether to execute the OSS import task immediately after it is created.
* `time_zone` - (Optional) Time Zone.
* `type` - (Optional) Check the frequency type. Log Service checks the query and analysis results according to the frequency you configured. The values are as follows: Fixedate: checks query and analysis results at regular intervals. Cron: specifies the time interval by using the Cron expression, and checks the query and analysis results at the specified time interval. Weekly: Check the query and analysis results at a fixed point in time on the day of the week. Daily: Check the query and analysis results at a fixed point in time every day. Hourly: Check the query and analysis results every hour.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project>:<oss_ingestion_name>`.
* `create_time` - Task creation time.
* `status` - OSS import task status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Oss Ingestion.
* `delete` - (Defaults to 5 mins) Used when delete the Oss Ingestion.
* `update` - (Defaults to 5 mins) Used when update the Oss Ingestion.

## Import

SLS Oss Ingestion can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_oss_ingestion.example <project>:<oss_ingestion_name>
```