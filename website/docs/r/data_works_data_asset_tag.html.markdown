---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_asset_tag"
description: |-
  Provides a Alicloud Data Works Data Asset Tag resource.
---

# alicloud_data_works_data_asset_tag

Provides a Data Works Data Asset Tag resource.

For information about Data Works Data Asset Tag and how to use it, see [What is Data Asset Tag](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdataassettag).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_data_works_data_asset_tag" "default" {
  key         = var.name
  value_type  = "String"
  description = "tag description"
  values      = ["value1", "value2"]
  managers    = ["user1"]
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the tag.
* `key` - (Required, ForceNew) The tag key.
* `managers` - (Optional) The tag administrators.
* `value_type` - (Required, ForceNew) The type of the tag value. Valid values: `Boolean`, `Int`, `String`, `Double`.
* `values` - (Optional) The tag values.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `category` - The type of the tag. Valid values: `Normal`, `System`.
* `create_time` - The creation time of the tag.
* `modify_time` - The modification time of the tag.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Data Asset Tag.
* `delete` - (Defaults to 5 mins) Used when delete the Data Asset Tag.
* `update` - (Defaults to 5 mins) Used when update the Data Asset Tag.

## Import

Data Works Data Asset Tag can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_data_asset_tag.example <key>
```
