---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group"
sidebar_current: "docs-alicloud-resource-ram-group"
description: |-
  Provides a RAM Group resource.
---

# alicloud_ram_group

Provides a RAM Group resource.

-> **NOTE:** When you want to destroy this resource forcefully(means remove all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully. 

-> **NOTE:** Available since v1.0.0+.

## Example Usage

```terraform
# Create a new RAM Group.
resource "alicloud_ram_group" "group" {
  name     = "groupName"
  comments = "this is a group comments."
  force    = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the RAM group. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `comments` - (Optional) Comment of the RAM group. This parameter can have a string of 1 to 128 characters.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The group ID.

## Import

RAM group can be imported using the id or name, e.g.

```shell
$ terraform import alicloud_ram_group.example my-group
```
