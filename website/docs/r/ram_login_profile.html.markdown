---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_login_profile"
sidebar_current: "docs-alicloud-resource-ram-login-profile"
description: |-
  Provides a RAM User Login Profile resource.
---

# alicloud\_ram\_login\_profile

Provides a RAM User Login Profile resource.


## Example Usage

```terraform
# Create a RAM login profile.
resource "alicloud_ram_user" "user" {
  name         = "user_test"
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
  force        = true
}

resource "alicloud_ram_login_profile" "profile" {
  user_name = alicloud_ram_user.user.name
  password  = "Yourpassword1234"
}
```
## Argument Reference

The following arguments are supported:

* `user_name` - (Required, ForceNew) Name of the RAM user. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen.
* `password` - (Required, Sensitive) Password of the RAM user.
* `mfa_bind_required` - (Optional) This parameter indicates whether the MFA needs to be bind when the user first logs in. Default value is `false`.
* `password_reset_required` - (Optional) This parameter indicates whether the password needs to be reset when the user first logs in. Default value is `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The login profile ID.
* `user_name` - The user name.
* `mfa_bind_required` - The parameter which indicates whether the MFA needs to be bind when the user first logs in.
* `password_reset_required` - The parameter which indicates whether the password needs to be reset when the user first logs in.

## Import

RAM login profile can be imported using the id or user name, e.g.

```shell
$ terraform import alicloud_ram_login_profile.example my-login
```
