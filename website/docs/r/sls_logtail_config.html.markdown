---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_logtail_config"
description: |-
  Provides a Alicloud Log Service (SLS) Logtail Config resource.
---

# alicloud_sls_logtail_config

Provides a Log Service (SLS) Logtail Config resource.



For information about Log Service (SLS) Logtail Config and how to use it, see [What is Logtail Config](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateConfig).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_logtail_config&exampleId=15e7d6fd-7d07-0cdf-4b27-a59bb1ce780ff4a6c0f4&activeTab=example&spm=docs.r.sls_logtail_config.0.15e7d6fd7d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

variable "name" {
  default = "tfaccsls62147"
}

variable "project_name" {
  default = "project-for-logtail-terraform"
}

resource "alicloud_log_project" "defaultuA28zS" {
  project_name = var.project_name
}

resource "alicloud_sls_logtail_config" "default" {
  project_name = alicloud_log_project.defaultuA28zS.project_name
  output_detail {
    endpoint      = "cn-hangzhou-intranet.log.aliyuncs.com"
    region        = "cn-hangzhou"
    logstore_name = "example"
  }

  output_type = "LogService"
  input_detail = jsonencode({
    "adjustTimezone" : false,
    "delayAlarmBytes" : 0,
    "delaySkipBytes" : 0,
    "discardNonUtf8" : false,
    "discardUnmatch" : true,
    "dockerFile" : false,
    "enableRawLog" : false,
    "enableTag" : false,
    "fileEncoding" : "utf8",
    "filePattern" : "access*.log",
    "filterKey" : ["key1"],
    "filterRegex" : ["regex1"],
    "key" : ["key1", "key2"],
    "localStorage" : true,
    "logBeginRegex" : ".*",
    "logPath" : "/var/log/httpd",
    "logTimezone" : "",
    "logType" : "common_reg_log",
    "maxDepth" : 1000,
    "maxSendRate" : -1,
    "mergeType" : "topic",
    "preserve" : true,
    "preserveDepth" : 0,
    "priority" : 0,
    "regex" : "(w+)(s+)",
    "sendRateExpire" : 0,
    "sensitive_keys" : [],
    "tailExisted" : false,
    "timeFormat" : "%Y/%m/%d %H:%M:%S",
    "timeKey" : "time",
    "topicFormat" : "none"
  })
  logtail_config_name = "tfaccsls62147"
  input_type          = "file"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sls_logtail_config&spm=docs.r.sls_logtail_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `create_time` - (Optional, ForceNew, Available since v1.255.0) The creation time of the resource
* `input_detail` - (Optional, ForceNew) The detailed configuration entered by logtail.
* `input_type` - (Optional, ForceNew) Method of log entry
* `last_modify_time` - (Optional, ForceNew, Computed, Int) Last modification time, unix timestamp
* `log_sample` - (Optional, ForceNew) Sample log
* `logtail_config_name` - (Optional, ForceNew, Computed, Available since v1.255.0) The name of the resource
* `output_detail` - (Optional, ForceNew, List, Available since v1.255.0) Detailed configuration of logtail output See [`output_detail`](#output_detail) below.
* `output_type` - (Optional, ForceNew) Log output mode. You can only upload data to log service.
* `project_name` - (Required, ForceNew, Available since v1.255.0) Project name

### `output_detail`

The output_detail supports the following:
* `endpoint` - (Optional, ForceNew, Available since v1.255.0) The endpoint of the log project.
* `logstore_name` - (Optional, ForceNew, Available since v1.255.0) The name of the output target logstore.
* `region` - (Optional, ForceNew, Available since v1.255.0) Region

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>:<logtail_config_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Logtail Config.
* `delete` - (Defaults to 5 mins) Used when delete the Logtail Config.

## Import

Log Service (SLS) Logtail Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_logtail_config.example <project_name>:<logtail_config_name>
```