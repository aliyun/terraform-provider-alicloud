---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_user"
sidebar_current: "docs-alicloud-resource-ram-user"
description: |-
  Provides a RAM User resource.
---

# alicloud_ram_user

Provides a RAM User resource.

For information about RAM User and how to use it, see [What is User](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-createuser).

-> **NOTE:** When you want to destroy this resource forcefully(means release all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully.

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

```terraform
# Create a new RAM user.
resource "alicloud_ram_user" "user" {
  name         = "terraform-example"
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
* `display_name` - (Optional) Name of the RAM user which for display. This name can have a string of 1 to 128 characters or Chinese characters, must contain only alphanumeric characters or Chinese characters or hyphens, such as "-",".", and must not end with a hyphen.
* `mobile` - (Optional) Phone number of the RAM user. This number must contain an international area code prefix, just look like this: 86-18600008888.
* `email` - (Optional) Email of the RAM user.
* `comments` - (Optional) Comment of the RAM user. This parameter can have a string of 1 to 128 characters.
* `force` - (Optional, Bool) This parameter is used for resource destroy. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of User.

## Timeouts

-> **NOTE:** Available since v1.209.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the User.
* `update` - (Defaults to 3 mins) Used when update the User.
* `delete` - (Defaults to 3 mins) Used when delete the User.

## Import

RAM User can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_user.example 123456789xxx
```
