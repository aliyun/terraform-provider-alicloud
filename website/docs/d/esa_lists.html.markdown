---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_lists"
description: |-
  Provides a list of ESA Lists to the user.
---

# alicloud_esa_lists

This data source provides the ESA Lists of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

resource "alicloud_esa_list" "default" {
  kind        = "ip"
  name        = var.name
  description = var.name
  items       = ["10.1.1.1", "10.1.1.2"]
}

data "alicloud_esa_lists" "ids" {
  ids = [alicloud_esa_list.default.id]
}

output "esa_lists_id_0" {
  value = data.alicloud_esa_lists.ids.lists.0.id
}
```

## Argument Reference

This data source automatically retrieves all pages of results; you do not need to specify `page_number` or `page_size`.

The following arguments are supported:

* `ids` - (Optional, List) A list of List IDs used to filter the returned results locally.
* `name_regex` - (Optional) A regex string used to filter the returned results locally by List name.
* `query_args` - (Optional, Set) The query parameters sent to the server. At most one block is supported. See [`query_args`](#query_args) below.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `query_args`

The query_args supports the following:

* `id_like` - (Optional) The fuzzy search for list ID.
* `name_like` - (Optional) The fuzzy search for list name.
* `description_like` - (Optional) The fuzzy search for list description.
* `name_item_like` - (Optional) The fuzzy search matching a list name or its items.
* `item_like` - (Optional) The value passed to the ListLists API to filter list contents.
* `kind` - (Optional) The type of the custom list, e.g. `ip`.
* `order_by` - (Optional) Specify the column to sort by.
* `desc` - (Optional, Bool) Whether to sort in descending order. Valid values: `true`, `false`.

QueryArgs values are passed unchanged for server-side interpretation; filtering behavior depends on the ESA API.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of List names.
* `lists` - A list of Lists. Each element contains the following attributes:
  * `id` - The ID of the List.
  * `list_id` - The ID of the List.
  * `name` - The name of the List.
  * `kind` - The type of the List.
  * `description` - The description of the List.
  * `length` - The number of items contained in the List.
  * `update_time` - The last modification time of the List.
