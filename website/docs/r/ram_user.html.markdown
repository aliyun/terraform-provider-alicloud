---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_user"
sidebar_current: "docs-alicloud-resource-ram-user"
description: |-
  Provides a RAM User resource.
---

# alicloud\_ram\_user

Provides a RAM User resource.

-> **NOTE:** When you want to destroy this resource forcefully(means release all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully.

## Example Usage

```
# Create a new RAM user.
resource "alicloud_ram_user" "user" {
  name         = "user_test"
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
  force        = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the RAM user. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen.
* `display_name` - (Optional) Name of the RAM user which for display. This name can have a string of 1 to 12 characters or Chinese characters, must contain only alphanumeric characters or Chinese characters or hyphens, such as "-",".", and must not end with a hyphen.
* `mobile` - (Optional) Phone number of the RAM user. This number must contain an international area code prefix, just look like this: 86-18600008888.
* `email` - (Optional) Email of the RAM user.
* `comments` - (Optional) Comment of the RAM user. This parameter can have a string of 1 to 128 characters.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The original id is user name, but it is user id in 1.37.0+.
* `name` - The user name.
* `display_name` - The user display name.
* `mobile` - The user phone number.
* `email` - The user email.
* `comments` - The user comments.

## Import

RAM user can be imported using the id or name, e.g.

```
$ terraform import alicloud_ram_user.example user
```
