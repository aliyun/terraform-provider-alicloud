---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_stores"
sidebar_current: "docs-alicloud-datasource-log-stores"
description: |-
  Provides a list of log stores to the user.
---

# alicloud\_log\_stores

This data source provides the Log Stores of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_log_stores" "example" {
  project = "the_project_name"
  ids     = ["the_store_name"]
}

output "first_log_store_id" {
  value = data.alicloud_log_stores.example.stores.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of store IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by store name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of store names.
* `stores` - A list of Log Stores. Each element contains the following attributes:
	* `id` - The ID of the store.
	* `store_name` - The name of the store. 
