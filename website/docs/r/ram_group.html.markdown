---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group"
description: |-
  Provides a Alicloud RAM Group resource.
---

# alicloud_ram_group

Provides a RAM Group resource.

The group that users can join.

For information about RAM Group and how to use it, see [What is Group](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-creategroup).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ram_group" "group" {
  group_name = var.name
  comments   = var.name
  force      = true
}
```

## Argument Reference

The following arguments are supported:
* `comments` - (Optional) The Group comment information. The maximum length is 128 characters.
* `group_name` - (Optional, ForceNew, Available since v1.245.0) The group name. You must specify at least one of the `group_name` and `name`.
It can be 1 to 64 characters in length and can contain letters, digits, periods (.), underscores (_), and dashes (-).

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.120.0). Field 'name' has been deprecated from provider version 1.120.0. New field 'group_name' instead.
* `force` - (Optional, Bool) Specifies whether to force delete the Group. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - (Available since v1.245.0) The create time of the group.

## Timeouts

-> **NOTE:** Available since v1.245.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Group.
* `delete` - (Defaults to 5 mins) Used when delete the Group.
* `update` - (Defaults to 5 mins) Used when update the Group.

## Import

RAM Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_group.example <id>
```