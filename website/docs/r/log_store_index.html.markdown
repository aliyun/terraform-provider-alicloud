---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_store_index"
sidebar_current: "docs-alicloud-resource-log-store-index"
description: |-
  Provides a Alicloud log store index resource.
---

# alicloud_log_store_index

Log Service provides the LogSearch/Analytics function to query and analyze large amounts of logs in real time.
You can use this function by enabling the index and field statistics. [Refer to details](https://www.alibabacloud.com/help/doc-detail/43772.htm)

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_store_index&exampleId=3978de35-af19-ffed-74a1-4574be38b2bbee3727b1&activeTab=example&spm=docs.r.log_store_index.0.3978de35af&intl_lang=EN_US" target="_blank">
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
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "example-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_store_index" "example" {
  project  = alicloud_log_project.example.project_name
  logstore = alicloud_log_store.example.logstore_name
  full_text {
    case_sensitive = true
    token          = " #$^*\r\n\t"
  }
  field_search {
    name             = "terraform-example"
    enable_analytics = true
    type             = "text"
    token            = " #$^*\r\n\t"
  }
}
```

## Module Support

You can use the existing [sls module](https://registry.terraform.io/modules/terraform-alicloud-modules/sls/alicloud) 
to create SLS project, store and store index one-click, like ECS instances.

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `logstore` - (Required, ForceNew) The log store name to the query index belongs.
* `log_reduce` - (Optional, Available since v1.263.0) Whether to enable log reduce. Default to false.
* `log_reduce_black_list` - (Optional, Available since v1.263.0) The black list of log reduce.
* `log_reduce_white_list` - (Optional, Available since v1.263.0) The white list of log reduce.
* `max_text_len` - (Optional, Available since v1.263.0) The max text length.
* `full_text` - (Optional) The configuration of full text index. See [`full_text`](#full_text) below.
* `field_search` - (Optional) List configurations of field search index. See [`field_search`](#field_search) below.

-> **Note:** At least one of the "full_text" and "field_search" should be specified.

### `full_text`

The full_text supports the following:

* `case_sensitive` - (Optional) Whether the case sensitive. Default to false.
* `include_chinese` - (Optional) Whether includes the chinese. Default to false.
* `token` - (Optional) The string of several split words, like "\r", "#"

### `field_search`

The field_search supports the following:

* `name` - (Required) The field name, which is unique in the same log store.
* `type` - (Optional) The type of one field. Valid values: ["long", "text", "double", "json"]. Default to "long".
* `alias` - (Optional) The alias of one field
* `case_sensitive` - (Optional) Whether the case sensitive for the field. Default to false. It is valid when "type" is "text" or "json".
* `include_chinese` - (Optional) Whether includes the chinese for the field. Default to false. It is valid when "type" is "text" or "json".
* `token` - (Optional) The string of several split words, like "\r", "#". It is valid when "type" is "text" or "json".
* `enable_analytics` - (Optional) Whether to enable field analytics. Default to true.
* `json_keys` - (Optional, Available in 1.66.0+) Use nested index when type is json. See [`json_keys`](#field_search-json_keys) below.

### `field_search-json_keys`

The field_search-json_keys supports the following:

* `name` - (Required) When using the json_keys field, this field is required.
* `type` - (Optional) The type of one field. Valid values: ["long", "text", "double"]. Default to "long"
* `alias` - (Optional) The alias of one field.
* `doc_value` - (Optional) Whether to enable statistics. default to true.

-> **Note:** At least one of the "full_text" and "field_search" should be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log store index. It formats of `<project>:<logstore>`.

## Import

Log store index can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_store_index.example tf-log:tf-log-store
```
