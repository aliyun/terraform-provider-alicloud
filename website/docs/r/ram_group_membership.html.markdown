---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group_membership"
sidebar_current: "docs-alicloud-resource-ram-group-membership"
description: |-
  Provides a RAM Group membership resource.
---

# alicloud_ram_group_membership

Provides a RAM Group membership resource. 

-> **NOTE:** Available since v1.0.0+.

## Example Usage

```terraform
variable "name" {
  default = "tfexample"
}
resource "alicloud_ram_group" "group" {
  name     = format("%sgroup", var.name)
  comments = "this is a group comments."
}

resource "alicloud_ram_user" "user" {
  name         = format("%suser", var.name)
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
}

resource "alicloud_ram_user" "user1" {
  name         = format("%suser1", var.name)
  display_name = "user_display_name1"
  mobile       = "86-18688888889"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
}

resource "alicloud_ram_group_membership" "membership" {
  group_name = alicloud_ram_group.group.name
  user_names = [alicloud_ram_user.user.name, alicloud_ram_user.user1.name]
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Required, ForceNew) Name of the RAM group. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `user_names` - (Required) Set of user name which will be added to group. Each name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen.

## Attributes Reference

The following attributes are exported:

* `id` - The membership ID, it's set to `group_name`

## Import
RAM Group membership can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_group_membership.example my-group
```
