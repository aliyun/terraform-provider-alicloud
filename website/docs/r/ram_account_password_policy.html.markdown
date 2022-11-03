---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_password_policy"
sidebar_current: "docs-alicloud-resource-ram-account-password-policy"
description: |-
  Provides a RAM password policy configuration for entire account.
---

# alicloud\_ram\_account\_password\_policy

Provides a RAM password policy configuration for entire account. Only one resource per account.

-> **NOTE:** This resource overwrites an existing configuration. During action `terraform destroy` it sets values the same as defaults for this resource (it does not preserve any preexisted configuration).

-> **NOTE:** Available in 1.46.0+

## Example Usage

Empty resource sets defaults values for every property.

```terraform
resource "alicloud_ram_account_password_policy" "default" {

}
```

```terraform
resource "alicloud_ram_account_password_policy" "corporate" {
  minimum_password_length      = 9
  require_lowercase_characters = false
  require_uppercase_characters = false
  require_numbers              = false
  require_symbols              = false
  hard_expiry                  = true
  max_password_age             = 12
  password_reuse_prevention    = 5
  max_login_attempts           = 3
}
```
For not specified values sets defaults.

## Argument Reference

The following arguments are supported:

* `minimum_password_length` - (Optional) Minimal required length of password for a user. Valid value range: [8-32]. Default to 12.
* `require_lowercase_characters` - (Optional) Specifies if the occurrence of a lowercase character in the password is mandatory. Default to true.
* `require_uppercase_characters` - (Optional) Specifies if the occurrence of an uppercase character in the password is mandatory. Default to true.
* `require_numbers` - (Optional) Specifies if the occurrence of a number in the password is mandatory. Default to true.
* `require_symbols` - (Optional Specifies if the occurrence of a special character in the password is mandatory. Default to true.
* `hard_expiry` - (Optional) Specifies if a password can expire in a hard way. Default to false.
* `max_password_age` - (Optional) The number of days after which password expires. A value of 0 indicates that the password never expires. Valid value range: [0-1095]. Default to 0.
* `password_reuse_prevention` - (Optional) User is not allowed to use the latest number of passwords specified in this parameter. A value of 0 indicates the password history check policy is disabled. Valid value range: [0-24]. Default to 0.
* `max_login_attempts` - (Optional, Type: int) Maximum logon attempts with an incorrect password within an hour. Valid value range: [0-32]. Default to 5.

## Import

RAM account password policy can be imported using the `id`, e.g.

```bash
$ terraform import alicloud_ram_account_password_policy.example ram-account-password-policy
```
