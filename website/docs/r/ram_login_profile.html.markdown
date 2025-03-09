---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_login_profile"
description: |-
  Provides a Alicloud RAM Login Profile resource.
---

# alicloud_ram_login_profile

Provides a RAM Login Profile resource.



For information about RAM Login Profile and how to use it, see [What is Login Profile](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-createloginprofile).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ram_user" "user" {
  name         = "terraform_example"
  display_name = "terraform_example"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "terraform_example"
  force        = true
}

resource "alicloud_ram_login_profile" "profile" {
  user_name = alicloud_ram_user.user.name
  password  = "Example_1234"
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Required, ForceNew) The name of the RAM user. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen.
* `password` - (Required, Sensitive) The logon password of the RAM user. The password must meet the password strength requirements.
* `password_reset_required` - (Optional, Bool) Specifies whether the RAM user must change the password upon logon. Default value: `false`. Valid values: `true`, `false`.
* `mfa_bind_required` - (Optional, Bool) Specifies whether an MFA device must be attached to the RAM user upon logon. Valid values: `true`, `false`. [To enhance the security of your resources and data, the default value has been changed to `true`](https://www.alibabacloud.com/en/notice/mfa20240524?_p_lc=1) .

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Login Profile.
* `delete` - (Defaults to 5 mins) Used when delete the Login Profile.
* `update` - (Defaults to 5 mins) Used when update the Login Profile.

## Import

RAM Login Profile can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_login_profile.example <id>
```