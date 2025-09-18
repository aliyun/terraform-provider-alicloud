---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_index"
description: |-
  Provides a Alicloud Log Service (SLS) Index resource.
---

# alicloud_sls_index

Provides a Log Service (SLS) Index resource.



For information about Log Service (SLS) Index and how to use it, see [What is Index](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateIndex).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-nanjing"
}

variable "logstore_name" {
  default = "logstore-example"
}

variable "project_name" {
  default = "project-for-index-terraform-example"
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
```

## Argument Reference

The following arguments are supported:
* `keys` - (Optional, Map) Field index
* `line` - (Optional, List) Full-text index See [`line`](#line) below.
* `log_reduce` - (Optional) Whether log clustering is enabled
* `log_reduce_black_list` - (Optional, List) The blacklist of the cluster fields of log clustering is filtered only when log clustering is enabled.
* `log_reduce_white_list` - (Optional, List) The whitelist of the cluster fields for log clustering. This filter is valid only when log clustering is enabled.
* `logstore_name` - (Required, ForceNew) Logstore name
* `max_text_len` - (Optional, Int) Maximum length of statistical field
* `project_name` - (Required, ForceNew) Project name

### `line`

The line supports the following:
* `case_sensitive` - (Required) Is case sensitive
* `chn` - (Required) Does it include Chinese
* `exclude_keys` - (Optional, List) List of excluded fields
* `include_keys` - (Optional, List) Include field list
* `token` - (Required, List) Delimiter

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>:<logstore_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Index.
* `delete` - (Defaults to 5 mins) Used when delete the Index.
* `update` - (Defaults to 5 mins) Used when update the Index.

## Import

Log Service (SLS) Index can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_index.example <project_name>:<logstore_name>
```