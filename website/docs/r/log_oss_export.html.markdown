---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_oss_export"
sidebar_current: "docs-alicloud-resource-log-oss-export"
description: |-
  Provides a Alicloud log oss export resource.
---

# alicloud\_log\_oss\_export
Log service data delivery management, this service provides the function of delivering data in logstore to oss product storage.
[Refer to details](https://www.alibabacloud.com/help/en/log-service/latest/ship-logs-to-oss-new-version).

-> **NOTE:** Available in 1.187.0+

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "tf-log-project"
  description = "created by terraform"
  tags        = { "test" : "test" }
}
resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "tf-log-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "alicloud_log_oss_export" "example" {
  project_name    = alicloud_log_project.example.name
  logstore_name   = alicloud_log_store.example.name
  export_name     = "oss_export_name"
  display_name    = "oss_export_display_name"
  bucket          = "test_bucket"
  prefix          = "root"
  suffix          = ""
  buffer_interval = 300
  buffer_size     = 250
  compress_type   = "none"
  path_format     = "%Y/%m/%d/%H/%M"
  content_type    = "json"
  json_enable_tag = true
  role_arn        = "role_arn_for_oss_write"
  log_read_role_arn = "role_arn_for_sls_read"
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `logstore_name` - (Required，ForceNew) The name of the log logstore.
* `export_name` - (Required，ForceNew) Delivery configuration name, it can only contain lowercase letters, numbers, dashes `-` and underscores `_`. It must start and end with lowercase letters or numbers, and the name must be 2 to 128 characters long.
* `display_name` - (Optional) The display name for oss export.
* `from_time` - (Optional) The log from when to export to oss.
* `prefix` - (Optional) The data synchronized from Log Service to OSS will be stored in this directory of Bucket.
* `suffix` - (Optional) The suffix for the objects in which the shipped data is stored.
* `bucket` - (Required) The name of the oss bucket.
* `buffer_interval` - (Required) How often is it delivered every interval.
* `buffer_size` - (Required) Automatically control the creation interval of delivery tasks and set the upper limit of an OSS object size (calculated in uncompressed), unit: `MB`.
* `role_arn` - (Optional) Used to write to oss bucket, the OSS Bucket owner creates the role mark which has the oss bucket write policy, such as `acs:ram::13234:role/logrole`.
* `log_read_role_arn` - (Optional, Available in 1.188.0+) Used for logstore reading, the role should have log read policy, such as `acs:ram::13234:role/logrole`, if `log_read_role_arn` is not set, `role_arn` is used to read logstore.
* `compress_type` - (Optional) OSS data storage compression method, support: `none`, `snappy`, `zstd`, `gzip`. Among them, none means that the original data is not compressed, and snappy means that the data is compressed using the snappy algorithm, which can reduce the storage space usage of the `OSS Bucket`.
* `path_format` - (Required) The OSS Bucket directory is dynamically generated according to the creation time of the export task, it cannot start with a forward slash `/`, the default value is `%Y/%m/%d/%H/%M`.
* `time_zone` - (Required) This time zone that is used to format the time, `+0800` e.g.
* `content_type` - (Required) Storage format, only supports three types: `json`, `parquet`, `orc`, `csv`.
  **According to the different format, please select the following parameters**
* `json_enable_tag` - (Optional) Whether to deliver the label when `content_type` = `json`.
* `csv_config_delimiter` - (Optional) Separator configuration in csv content_type.
* `csv_config_columns` - (Optional) Field configuration in csv content_type.
* `csv_config_null` - (Optional) Invalid field content in csv content_type.
* `csv_config_quote` - (Optional) Escape character in csv content_type.
* `csv_config_header` - (Optional) Indicates whether to write the field name to the CSV file, the default value is `false`.
* `csv_config_linefeed` - (Optional) lineFeed in csv content_type.
* `csv_config_escape` - (Optional) escape in csv content_type.
* `config_columns` - (Optional) Configure columns when `content_type` is `parquet` or `orc`.
  *  `name` - (Required) The name of the key.
  *  `type` - (Required) Type of configuration name.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log oss export.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when Creating LogOssexport instance. 
* `update` - (Defaults to 1 mins) Used when Updating LogOssexport instance. 
* `delete` - (Defaults to 1 mins) Used when terminating the LogOssexport instance.

## Import

Log oss export can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_oss_export.example tf-log-project:tf-log-logstore:tf-log-export
```
