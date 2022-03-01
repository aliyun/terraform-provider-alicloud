---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_ingestion"
sidebar_current: "docs-alicloud-resource-log-ingestion"
description: |-
  Provides a Alicloud log ingestion resource.
---

# alicloud\_log\_ingestion
Log service ingestion, this service provides the function of importing logs of various data sources(OSS, MaxCompute) into logstore.
[Refer to details](https://www.alibabacloud.com/help/en/doc-detail/147819.html).

-> **NOTE:** Available in 1.161.0+

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
resource "alicloud_log_ingestion" "example" {
  project         =  alicloud_log_project.example.name
  logstore        =  alicloud_log_store.example.name
  ingestion_name  =  "ingestion_name"
  display_name    =  "display_name"
  description     =  "oss2sls"
  interval        =  "30m"
  run_immediately =  true
  time_zone       =  "+0800"
  source          =  <<DEFINITION
        {
          "bucket": "bucket_name",
          "compressionCodec": "none",
          "encoding": "UTF-8",
          "endpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
          "format": {
            "escapeChar": "\\",
            "fieldDelimiter": ",",
            "fieldNames": [],
            "firstRowAsHeader": true,
            "maxLines": 1,
            "quoteChar": "\"",
            "skipLeadingRows": 0,
            "timeField": "",
            "type": "DelimitedText"
          },
          "pattern": "",
          "prefix": "test-prefix/",
          "restoreObjectEnabled": false,
          "roleARN": "acs:ram::1049446484210612:role/aliyunlogimportossrole",
          "type": "AliyunOSS"
        }
  DEFINITION
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `logstore` - (Required，ForceNew) The name of the target logstore.
* `ingestion_name` - (Required，ForceNew) Ingestion job name, it can only contain lowercase letters, numbers, dashes `-` and underscores `_`. It must start and end with lowercase letters or numbers, and the name must be 2 to 128 characters long.
* `display_name` - (Required) The name displayed on the web page.
* `description` - (Optional) Ingestion job description.
* `interval` - (Required) Task execution interval, support minute `m`, hour `h`, day `d`, for example 30 minutes `30m`.
* `run_immediately` - (Required) Whether to run the ingestion job immediately, if false, wait for an interval before starting the ingestion.
* `time_zone` - (Optional) Which time zone is the log time imported in, e.g. `+0800`.
* `source` - (Required) Data source and data format details. [Refer to details](https://www.alibabacloud.com/help/en/doc-detail/147819.html).


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log ingestion. It formats of `<project>:<logstore>:<ingetion_name>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when Creating LogIngestion instance.
* `update` - (Defaults to 1 mins) Used when Updating LogIngestion instance.
* `delete` - (Defaults to 1 mins) Used when terminating the LogIngestion instance.

## Import

Log ingestion can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_ingestion.example tf-log-project:tf-log-logstore:ingestion_name
```
