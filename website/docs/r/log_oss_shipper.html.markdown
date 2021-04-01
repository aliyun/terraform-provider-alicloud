---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_oss_shipper"
sidebar_current: "docs-alicloud-resource-log-oss-shipper"
description: |-
  Provides a Alicloud log oss shipper resource.
---

# alicloud\_log\_oss_shipper
Log service data delivery management, this service provides the function of delivering data in logstore to oss product storage.
[Refer to details](https://www.alibabacloud.com/help/zh/doc-detail/43724.htm).


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

resource "alicloud_log_oss_shipper" "example" {
  project_name    = alicloud_log_project.example.name
  logstore_name   = alicloud_log_logstore.example.name
  shipper_name    = "oss_shipper_name"
  oss_bucket      = "test_bucket"
  oss_prefix      = "/root"
  buffer_interval = 300
  buffer_size     = 250
  compress_type   = "none"
  path_format     = "%Y/%m/%d/%H/%M"
  format          = "json"
  json_enable_tag = true
}

```


## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `logstore_name` - (Required，ForceNew) The name of the log logstore.
* `shipper_name` - (Required，ForceNew) Delivery configuration name.
* `oss_prefix` - (Optional) The data synchronized from Log Service to OSS will be stored in this directory of Bucket.
* `oss_bucket` - (Required) The name of the oss bucket.
* `buffer_interval` - (Required) How often is it delivered every interval.
* `buffer_size` - (Required) Automatically control the creation interval of delivery tasks and set the upper limit of an OSS object size (calculated in uncompressed), unit: `MB`.
* `role_arn` - (Optional) Used for access control, the OSS Bucket owner creates the role mark, such as "acs:ram::13234:role/logrole"
* `compress_type` - (Optional) OSS data storage compression method, support: none, snappy. Among them, none means that the original data is not compressed, and snappy means that the data is compressed using the snappy algorithm, which can reduce the storage space usage of the `OSS Bucket`.
* `path_format` - (Required) The storage format only supports three types: json, parquet, csv.
* `format` - (Required) Storage format.
* `json_enable_tag` - (Optional) Whether to deliver the label.
* `csv_config_delimiter` - (Optional) Separator configuration in csv configuration format.
* `csv_config_columns` - (Optional) Field configuration in csv configuration format.
* `csv_config_nullidentifier` - (Optional) Invalid field content.
* `csv_config_quote` - (Optional) Escape character under csv configuration.
* `csv_config_header` - (Optional) Indicates whether to write the field name to the CSV file,The default value is false.
* `csv_config_linefeed` - (Optional) Separator in csv configuration.
* `parquet_config` - (Optional) Configure to use parquet storage format.
    * `name` - (Required) The name of the key.
    * `type` - (Required) Type of configuration name.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log oss shipper. It sames as its name.


## Import

Log oss shipper can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_oss_shipper.example tf-log-project:tf-log-logstore:tf-log-shipper
```
