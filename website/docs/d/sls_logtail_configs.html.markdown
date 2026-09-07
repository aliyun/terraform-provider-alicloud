---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_logtail_configs"
sidebar_current: "docs-alicloud-datasource-sls-logtail-configs"
description: |-
  Provides a list of Sls Logtail Config owned by an Alibaba Cloud account.
---

# alicloud_sls_logtail_configs

This data source provides Sls Logtail Config available to the user.[What is Logtail Config](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateConfig)

-> **NOTE:** Available since v1.259.0.

## Example Usage

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

data "alicloud_sls_logtail_configs" "default" {
  logtail_config_name = alicloud_sls_logtail_config.default.logtail_config_name
  logstore_name       = "example"
  project_name        = var.project_name
  offset              = 0
  size                = 100
}

output "alicloud_sls_logtail_config_example_id" {
  value = "${data.alicloud_sls_logtail_configs.default.configs.0.id}"
}
```

## Argument Reference

The following arguments are supported:
* `logstore_name` - (Required, ForceNew) Logstore name.
* `logtail_config_name` - (ForceNew, Optional) The name of the resource
* `offset` - (Required, ForceNew) Query start row. The default value is 0.
* `project_name` - (Required, ForceNew) Project name
* `size` - (Required, ForceNew) The number of rows per page set for a pagination query. The maximum value is 500.
* `ids` - (Optional, ForceNew, Computed) A list of Logtail Config IDs. The value is formulated as `<project_name>:<logtail_config_name>`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Logtail Config IDs.
* `names` - A list of name of Logtail Configs.
* `configs` - A list of Logtail Config Entries. Each element contains the following attributes:
  * `logtail_config_name` - The name of the resource
  * `id` - The ID of the resource supplied above.
