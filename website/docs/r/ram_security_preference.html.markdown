---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_security_preference"
sidebar_current: "docs-alicloud-resource-ram-security-preference"
description: |-
  Provides a Alicloud RAM Security Preference resource.
---

# alicloud\_ram\_security\_preference

Provides a RAM Security Preference resource.

For information about RAM Security Preference and how to use it, see [What is Security Preference](https://www.alibabacloud.com/help/en/doc-detail/186694.htm).

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ram_security_preference" "example" {
  enable_save_mfa_ticket        = false
  allow_user_to_change_password = true
}
```

## Argument Reference

The following arguments are supported:

* `enable_save_mfa_ticket` - (Optional) Specifies whether to remember the MFA devices for seven days. Valid values: `true` and `false`
  * `true` - remembers the MFA devices for seven days.
  * `false` - does not remember the MFA devices. This is the default value.
* `allow_user_to_change_password` - (Optional) Specifies whether RAM users can change their passwords. Valid values: `true` and `false`
  * `true` - RAM users can change their passwords. This is the default value.
  * `false` - RAM users cannot change their passwords.
* `allow_user_to_manage_access_keys` - (Optional) Specifies whether RAM users can manage their AccessKey pairs. Valid values: `true` and `false`
  * `true` - RAM users can manage their AccessKey pairs.
  * `false` - RAM users cannot manage their AccessKey pairs. This is the default value.
* `allow_user_to_manage_mfa_devices` - (Optional) Specifies whether RAM users can manage their MFA devices. Valid values: `true` and `false`
  * `true` - RAM users can manage their MFA devices. This is the default value.
  * `false` - RAM users cannot manage their MFA devices.
* `login_session_duration` - (Optional) The validity period of the logon session of RAM users. Valid values: 6 to 24. Unit: hours. Default value: 6.
* `enforce_mfa_for_login` - (Optional) Specifies whether MFA is required for all RAM users when they log on to the Alibaba Cloud Management Console by using usernames and passwords. Valid values: `true` and `false`
  * `true` - MFA is required for all RAM users when they log on to the Alibaba Cloud Management Console by using usernames and passwords.
  * `false` - User-specific settings are applied. This is the default value.
* `login_network_masks` - (Optional) The subnet mask that specifies the IP addresses from which you can log on to the Alibaba Cloud Management Console. This parameter takes effect on password-based logon and single sign-on (SSO). However, this parameter does not take effect on API calls that are authenticated by using AccessKey pairs.**NOTE:** You can specify up to 25 subnet masks. The total length of the subnet masks can be a maximum of 512 characters.
  * If you specify a subnet mask, RAM users can use only the IP addresses in the subnet mask to log on to the Alibaba Cloud Management Console.  
  * If you do not specify a subnet mask, RAM users can use all IP addresses to log on to the Alibaba Cloud Management Console.
  * If you need to specify multiple subnet masks, separate the subnet masks with semicolons (;). Example: 192.168.0.0/16;10.0.0.0/8.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Security Preference. The ID is set as `RamSecurityPreference`. 

## Import

RAM Security Preference can be imported using the id, e.g.

```
$ terraform import alicloud_ram_security_preference.example <id>
```