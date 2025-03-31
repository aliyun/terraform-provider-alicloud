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
* `hard_expiry` - (Optional) Whether to restrict logon after the password expires. Value:
  - true: After the password expires, you cannot log in to the console. You must reset the password of the RAM user through the main account or a RAM user with administrator permissions to log on normally.
  - false (default): After the password expires, the RAM user can change the password and log on normally.
* `max_login_attemps` - (Optional, Int) Password retry constraint. After entering the wrong password continuously for the set number of times, the account will be locked for one hour.
Value range: 0~32.
Default value: 0, which means that the password retry constraint is not enabled.
* `max_password_age` - (Optional, Int) Password validity period.
Value range: 0~1095. Unit: days.
Default value: 0, which means never expires.
* `minimum_password_different_character` - (Optional, Int) The minimum number of unique characters in the password.
Valid values: 0 to 8.
The default value is 0, which indicates that no limits are imposed on the number of unique characters in a password.
* `minimum_password_length` - (Optional, Computed, Int) The minimum number of characters in the password.
Valid values: 8 to 32. Default value: 8.
* `password_not_contain_user_name` - (Optional) Whether the user name is not allowed in the password. Value:
  - true: The password cannot contain the user name.
  - false (default): The user name can be included in the password.
* `password_reuse_prevention` - (Optional, Int) Historical password check policy.
Do not use the previous N Passwords. The value range of N is 0 to 24.
Default value: 0, indicating that the historical password check policy is not enabled.
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
* `id` - The ID of the resource supplied above. This field is set to your Alibaba Cloud Account ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Password Policy.
* `update` - (Defaults to 5 mins) Used when update the Password Policy.

## Import

RAM Password Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_password_policy.example <id>.
```