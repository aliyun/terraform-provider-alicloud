---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_asset_tags"
sidebar_current: "docs-alicloud-datasource-data-works-data-asset-tags"
description: |-
  Provides a list of Data Works Data Asset Tags to the user.
---

# alicloud\_data\_works\_data\_asset\_tags

This data source provides the Data Works Data Asset Tags of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_data_works_data_asset_tag" "default" {
  key        = var.name
  value_type = "String"
}

data "alicloud_data_works_data_asset_tags" "default" {
  key = alicloud_data_works_data_asset_tag.default.key
  ids = [alicloud_data_works_data_asset_tag.default.key]
}

output "data_asset_tag_id" {
  value = data.alicloud_data_works_data_asset_tags.default.tags.0.id
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Optional) The type of the tag. Valid values: `Normal`, `System`.
* `ids` - (Optional, Computed) A list of Data Asset Tag IDs.
* `key` - (Optional) The tag key.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `value_type` - (Optional) The type of the tag value. Valid values: `Boolean`, `Int`, `String`, `Double`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `tags` - A list of Data Works Data Asset Tags. Each element contains the following attributes:
  * `category` - The type of the tag.
  * `create_time` - The creation time of the tag.
  * `description` - The description of the tag.
  * `id` - The ID of the Data Asset Tag.
  * `key` - The tag key.
  * `managers` - The tag administrators.
  * `modify_time` - The modification time of the tag.
  * `value_type` - The type of the tag value.
  * `values` - The tag values.
