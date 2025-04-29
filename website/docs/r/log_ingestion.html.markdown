---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_ingestion"
sidebar_current: "docs-alicloud-resource-log-ingestion"
description: |-
  Provides a Alicloud log ingestion resource.
---

# alicloud_log_ingestion

Log service ingestion, this service provides the function of importing logs of various data sources(OSS, MaxCompute) into logstore. [Refer to details](https://www.alibabacloud.com/help/en/doc-detail/147819.html).

-> **NOTE:** Available since v1.161.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_ingestion&exampleId=0a8f2c75-2eb9-48cf-cc0a-927585587c699e5d6cbe&activeTab=example&spm=docs.r.log_ingestion.0.0a8f2c752e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "example-store"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_ingestion" "example" {
  project         = alicloud_log_project.example.project_name
  logstore        = alicloud_log_store.example.logstore_name
  ingestion_name  = "terraform-example"
  display_name    = "terraform-example"
  description     = "terraform-example"
  interval        = "30m"
  run_immediately = true
  time_zone       = "+0800"
  source          = <<DEFINITION
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when Creating LogIngestion instance.
* `update` - (Defaults to 1 mins) Used when Updating LogIngestion instance.
* `delete` - (Defaults to 1 mins) Used when terminating the LogIngestion instance.

## Import

Log ingestion can be imported using the id or name, e.g.

```shell
$ terraform import alicloud_log_ingestion.example tf-log-project:tf-log-logstore:ingestion_name
```
