---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_password_policy"
description: |-
  Provides a Alicloud RAM Password Policy resource.
---

# alicloud_ram_password_policy

Provides a RAM Password Policy resource.

Password strength information.

For information about RAM Password Policy and how to use it, see [What is Password Policy](https://next.api.alibabacloud.com/document/Ram/2015-05-01/SetPasswordPolicy).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_ram_password_policy" "default" {
  minimum_password_length              = "8"
  require_lowercase_characters         = false
  require_numbers                      = false
  max_password_age                     = "0"
  password_reuse_prevention            = "1"
  max_login_attemps                    = "1"
  hard_expiry                          = false
  require_uppercase_characters         = false
  require_symbols                      = false
  password_not_contain_user_name       = false
  minimum_password_different_character = "1"
}
```

### Deleting `alicloud_ram_password_policy` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ram_password_policy`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `hard_expiry` - (Optional) Specifies whether a password will expire. Valid values: `true` and `false`. Default value: `false`. If you leave this parameter unspecified, the default value false is used.

  - If you set this parameter to `true`, the Alibaba Cloud account to which the RAM users belong must reset the passwords before the RAM users can log on to the Alibaba Cloud Management Console.
  - If you set this parameter to `false`, the RAM users can change the passwords after the passwords expire and then log on to the Alibaba Cloud Management Console.
* `max_login_attemps` - (Optional, Int) The maximum number of password retries. If you enter the wrong passwords for the specified consecutive times, the account is locked for one hour.
Valid values: 0 to 32.
The default value is 0, which indicates that the password retries are not limited.
* `max_password_age` - (Optional, Int) The number of days for which a password is valid. If you reset a password, the password validity period restarts. Default value: 0. The default value indicates that the password never expires.
* `minimum_password_different_character` - (Optional, Int) The minimum number of unique characters in the password.
Valid values: 0 to 8.
The default value is 0, which indicates that no limits are imposed on the number of unique characters in a password.
* `minimum_password_length` - (Optional, Computed, Int) The minimum number of characters in the password.
Valid values: 8 to 32. Default value: 8.
* `password_not_contain_user_name` - (Optional) Specifies whether to exclude the username from the password. Valid values:

  - true: A password cannot contain the username.
  - false: A password can contain the username. This is the default value.
* `password_reuse_prevention` - (Optional, Int) The policy for password history check.
The default value is 0, which indicates that RAM users can reuse previous passwords.
* `require_lowercase_characters` - (Optional) Specifies whether the password must contain lowercase letters. Valid values:

  - true
  - false (default)
* `require_numbers` - (Optional) Specifies whether the password must contain digits. Valid values:

  - true
  - false (default)
* `require_symbols` - (Optional) Specifies whether the password must contain special characters. Valid values:

  - true
  - false (default)
* `require_uppercase_characters` - (Optional) Specifies whether the password must contain uppercase letters. Valid values:

  - true
  - false (default)

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Password Policy.
* `update` - (Defaults to 5 mins) Used when update the Password Policy.

## Import

RAM Password Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_password_policy.example 
```