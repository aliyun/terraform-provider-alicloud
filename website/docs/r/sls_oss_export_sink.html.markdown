---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_oss_export_sink"
description: |-
  Provides a Alicloud Log Service (SLS) Oss Export Sink resource.
---

# alicloud_sls_oss_export_sink

Provides a Log Service (SLS) Oss Export Sink resource.

OSS export task.

For information about Log Service (SLS) Oss Export Sink and how to use it, see [What is Oss Export Sink](https://www.alibabacloud.com/help/en/sls/developer-reference/api-sls-2020-12-30-createossexport).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_oss_export_sink&exampleId=21d3deb2-4eda-f28f-57e3-941c074cfc19fa89488f&activeTab=example&spm=docs.r.sls_oss_export_sink.0.21d3deb24e&intl_lang=EN_US" target="_blank">
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

resource "alicloud_log_project" "defaulteyHJsO" {
  description  = "terraform-oss-example-910"
  project_name = format("%s1%s", var.name, random_integer.default.result)
}

resource "alicloud_log_store" "defaultxeHfXC" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.defaulteyHJsO.project_name
  logstore_name    = format("%s1%s", var.name, random_integer.default.result)
}

resource "alicloud_oss_bucket" "defaultiwj0xO" {
  bucket        = format("%s1%s", var.name, random_integer.default.result)
  storage_class = "Standard"
}


resource "alicloud_sls_oss_export_sink" "default" {
  project = alicloud_log_project.defaulteyHJsO.project_name
  configuration {
    logstore = alicloud_log_store.defaultxeHfXC.logstore_name
    role_arn = "acs:ram::12345678901234567:role/aliyunlogdefaultrole"
    sink {
      bucket           = alicloud_oss_bucket.defaultiwj0xO.bucket
      role_arn         = "acs:ram::12345678901234567:role/aliyunlogdefaultrole"
      time_zone        = "+0700"
      content_type     = "json"
      compression_type = "none"
      content_detail   = jsonencode({ "enableTag" : false })
      buffer_interval  = "300"
      buffer_size      = "256"
      endpoint         = "https://oss-cn-shanghai-internal.aliyuncs.com"
    }
    from_time = "1732165733"
    to_time   = "1732166733"
  }
  job_name     = "export-oss-1731404933-00001"
  display_name = "exampleterraform"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sls_oss_export_sink&spm=docs.r.sls_oss_export_sink.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `configuration` - (Required, List) OSSExportConfiguration See [`configuration`](#configuration) below.
* `description` - (Optional) The description of the job.
* `display_name` - (Required) The display name of the job.
* `job_name` - (Required, ForceNew) The unique identifier of the OSS data shipping job.
* `project` - (Required, ForceNew) The name of the project.

### `configuration`

The configuration supports the following:
* `from_time` - (Required, Int) The beginning of the time range to ship data. The value 1 specifies that the data shipping job ships data from the first log in the Logstore. Example value: 1718380800
* `logstore` - (Required) The name of the Logstore.
* `role_arn` - (Required) The Alibaba Cloud Resource Name (ARN) of the Resource Access Management (RAM) role that is used to read data from Simple Log Service. Example value: acs:ram::1234567890:role/aliyunlogdefaultrole

* `sink` - (Required, List) The configurations of the Object Storage Service (OSS) data shipping job. See [`sink`](#configuration-sink) below.
* `to_time` - (Required, Int) The end of the time range to ship data. The value 0 specifies that the data shipping job continuously ships data until the job is manually stopped. Example value: 1718380800

### `configuration-sink`

The configuration-sink supports the following:
* `bucket` - (Required) The OSS bucket.
* `buffer_interval` - (Required) The interval between two data shipping operations. Valid values: 300 to 900. Unit: seconds.

* `buffer_size` - (Required) The size of the OSS object to which data is shipped. Valid values: 5 to 256. Unit: MB.
* `compression_type` - (Required) Supports four compression types, such as snappy, gzip, zstd, and none.
* `content_detail` - (Required, JsonString) The OSS file content details. Note: the value of this parameter should be updated based on the value of the contentType parameter.

  If the contentType value is JSON, the parameters of the contentDetail value are as follows:

  If the tag is allowed to be posted, the value of the parameter enableTag is true. Example:{"enableTag": true}

  You are not allowed to post tags. The value of the parameter enableTag is false. Example:{"enableTag": false}

  If the contentType value is csv, the parameters of the contentDetail value are as follows:

  The parameter columns is the key of the log in the source logstore.

  The delimiter parameter, which can be ",","|","", or "\t".

  The header parameter determines whether the OSS file retains the header. The optional value is true or false.

  The lineFeed parameter. Optional values are "\t", "\n", or "".

  The invalid field content parameter is null to specify the delivery content when the field name does not exist.

  The escape character parameter "quote". Optional values are "" "," '", or" ".

  Example:{"null": "-", "header": false, "lineFeed": "\n", "quote": "", "delimiter": ",", "columns": ["a", "B", "c", "d"]}

  When the contentType value is parquet, the parameters of the contentDetail value are as follows:

  The columns parameter is the key of the log in the source Logstore and must carry the data type of the key, for example:{"columns": [{"name": "a", "type": "string"}, {"name": "B", "type": "string"}, {"name": "c", "type": "string": "string"}]}

  When the contentType value is set to orc, the parameters of the contentDetail value are as follows:

  The columns parameter is the key of the log in the source Logstore and must carry the data type of the key, for example:{"columns": [{"name": "a", "type": "string"}, {"name": "B", "type": "string"}, {"name": "c", "type": "string": "string"}]}
* `content_type` - (Required) The storage format of the OSS object. Valid values: json, parquet, csv, and orc.
* `delay_seconds` - (Optional, Int) The latency of data shipping. The value of this parameter cannot exceed the data retention period of the source Logstore.
* `endpoint` - (Required) The OSS Endpoint can only be an OSS intranet Endpoint and only supports the same region. Example value: https://oss-cn-hangzhou-internal.aliyuncs.com
* `path_format` - (Optional) The directory is dynamically generated according to the time. The default value is% Y/%m/%d/%H/%M. The corresponding generated directory is, for example, 2017/01/23/12/00. Note that the partition format cannot start and end. Example values:%Y/%m/%d
* `path_format_type` - (Optional) The partition format type. only support time
* `prefix` - (Optional) The prefix of the OSS object.
* `role_arn` - (Required) The ARN of the RAM role that is used to write data to OSS. Example value: acs:ram::xxxxxxx

* `suffix` - (Optional) The suffix of the OSS object.
* `time_zone` - (Required) The time zone. Example value: +0800

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project>:<job_name>`.
* `create_time` - Creation time. Example value: 1718787534
* `status` - The status of the post task. Example value: RUNNING

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Oss Export Sink.
* `delete` - (Defaults to 5 mins) Used when delete the Oss Export Sink.
* `update` - (Defaults to 5 mins) Used when update the Oss Export Sink.

## Import

Log Service (SLS) Oss Export Sink can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_oss_export_sink.example <project>:<job_name>
```