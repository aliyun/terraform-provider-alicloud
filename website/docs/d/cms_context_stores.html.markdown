---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_context_stores"
description: |-
  Provides a list of Cms Context Stores to the user.
---

# alicloud_cms_context_stores

This data source provides the Cms Context Stores of the current Alibaba Cloud user, filtered by workspace.

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
data "alicloud_cms_context_stores" "example" {
  workspace  = alicloud_cms_workspace.default.workspace_name
  name_regex = "example-context-store"
}

output "context_store_ids" {
  value = data.alicloud_cms_context_stores.example.ids
}

output "first_context_store_name" {
  value = data.alicloud_cms_context_stores.example.names.0
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The name of the workspace to which the context stores belong.
* `context_store_name` - (Optional) Filter the results by the name of the context store.
* `context_type` - (Optional) Filter the results by the type of the context store.
* `name_regex` - (Optional) A regex string to filter results by the context store name.
* `output_file` - (Optional) File name where to save the data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of Context Store IDs. Each ID is formatted as `<workspace>:<context_store_name>`.
* `names` - A list of Context Store names.
* `context_stores` - A list of Context Stores. Each element contains the following attributes:
  * `context_store_name` - The name of the context store.
  * `context_type` - The type of the context store.
  * `workspace` - The name of the workspace to which the context store belongs.
  * `description` - The description of the context store.
  * `status` - The status of the context store.
  * `region_id` - The region ID of the context store.
  * `create_time` - The creation time of the context store.
  * `update_time` - The last update time of the context store.
