---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_indexs"
sidebar_current: "docs-alicloud-datasource-sls-indexs"
description: |-
  Provides a list of Sls Index owned by an Alibaba Cloud account.
---

# alicloud_sls_indexs

This data source provides Sls Index available to the user.[What is Index](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateIndex)

-> **NOTE:** Available since v1.262.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-nanjing"
}

variable "logstore_name" {
  default = "logstore-example-1"
}

variable "project_name" {
  default = "project-for-index-terraform-example-1"
}

resource "alicloud_log_project" "default" {
  description  = "terraform example"
  project_name = var.project_name
}

resource "alicloud_log_store" "default" {
  hot_ttl          = "7"
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.default.project_name
  logstore_name    = var.logstore_name
}

resource "alicloud_sls_index" "default" {
  line {
    chn            = "true"
    case_sensitive = "true"
    token = [
      "a"
    ]
    exclude_keys = [
      "t"
    ]
  }
  keys = jsonencode(
    {
      "example" : {
        "caseSensitive" : false,
        "token" : [
          "\n",
          "\t",
          ",",
          " ",
          ";",
          "\"",
          "'",
          "(",
          ")",
          "{",
          "}",
          "[",
          "]",
          "<",
          ">",
          "?",
          "/",
          "#",
          ":"
        ],
        "type" : "text",
        "doc_value" : false,
        "alias" : "",
        "chn" : false
      }
    }
  )

  logstore_name = alicloud_log_store.default.logstore_name
  project_name  = var.project_name
}

data "alicloud_sls_indexs" "default" {
  logstore_name = alicloud_log_store.default.logstore_name
  project_name  = alicloud_log_project.default.project_name
}

output "alicloud_sls_index_example_id" {
  value = data.alicloud_sls_indexs.default.indexs.0.id
}
```

## Argument Reference

The following arguments are supported:
* `logstore_name` - (Required, ForceNew) Logstore name
* `project_name` - (Required, ForceNew) Project name
* `ids` - (Optional, ForceNew, Computed) A list of Index IDs. The value is formulated as `<project_name>:<logstore_name>`.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Index IDs.
* `indexs` - A list of Index Entries. Each element contains the following attributes:
  * `keys` - Field index
  * `line` - Full-text index
    * `case_sensitive` - Is case sensitive.
    * `chn` - Does it include Chinese.
    * `exclude_keys` - List of excluded fields.
    * `include_keys` - Include field list.
    * `token` - Delimiter.
  * `log_reduce_black_list` - The blacklist of the cluster fields of log clustering is filtered only when log clustering is enabled.
  * `log_reduce_white_list` - The whitelist of the cluster fields for log clustering. This filter is valid only when log clustering is enabled.
  * `max_text_len` - Maximum length of statistical field
  * `ttl` - Log index storage time
  * `id` - The ID of the resource supplied above.
