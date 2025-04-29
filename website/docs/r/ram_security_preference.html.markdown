---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_security_preference"
description: |-
  Provides a Alicloud RAM Security Preference resource.
---

# alicloud_ram_security_preference

Provides a RAM Security Preference resource.



For information about RAM Security Preference and how to use it, see [What is Security Preference](https://www.alibabacloud.com/help/en/doc-detail/186694.htm).

-> **NOTE:** Available since v1.152.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_security_preference&exampleId=a8427f8a-f030-814f-bcfd-9c52d5811d5d82082bb0&activeTab=example&spm=docs.r.ram_security_preference.0.a8427f8af0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ram_security_preference" "example" {
  enable_save_mfa_ticket        = false
  allow_user_to_change_password = true
}
```

### Deleting `alicloud_ram_security_preference` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ram_security_preference`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `allow_user_to_change_password` - (Optional, Computed) Whether to allow RAM users to manage their own passwords. Value:
  - true (default): Allowed.
  - false: not allowed.
* `allow_user_to_login_with_passkey` - (Optional, Available since v1.248.0) Whether to allow RAM users to log on using a passkey. Value:
  - true (default): Allowed.
  - false: not allowed.
* `allow_user_to_manage_access_keys` - (Optional, Computed) Whether to allow RAM users to manage their own access keys. Value:
  - true: Allow.
  - false (default): Not allowed.
* `allow_user_to_manage_mfa_devices` - (Optional, Computed) Whether to allow RAM users to manage multi-factor authentication devices. Value:
  - true (default): Allowed.
  - false: not allowed.
* `allow_user_to_manage_personal_ding_talk` - (Optional, Available since v1.248.0) Whether to allow RAM users to independently manage the binding and unbinding of personal DingTalk. Value:
  - true (default): Allowed.
  - false: not allowed.
* `enable_save_mfa_ticket` - (Optional, Computed) Whether to save the verification status of a RAM user after logging in using multi-factor authentication. The validity period is 7 days. Value:
  - true: Allow.
  - false (default): Not allowed.
* `login_network_masks` - (Optional) The login mask. The logon mask determines which IP addresses are affected by the logon console, including password logon and single sign-on (SSO), but API calls made using the access key are not affected.
  - If the mask is specified, RAM users can only log on from the specified IP address.
  - If you do not specify any mask, the login console function will apply to the entire network.

When you need to configure multiple login masks, use a semicolon (;) to separate them, for example: 192.168.0.0/16;10.0.0.0/8.

Configure a maximum of 40 logon masks, with a total length of 512 characters.
* `login_session_duration` - (Optional, Computed, Int) The validity period of the logon session of RAM users.
Valid values: 1 to 24. Unit: hours.
Default value: 6.
* `mfa_operation_for_login` - (Optional, Computed, Available since v1.248.0) MFA must be used during logon (replace the original EnforceMFAForLogin parameter, the original parameter is still valid, we recommend that you update it to a new parameter). Value:
  - mandatory: mandatory for all RAM users. The original value of EnforceMFAForLogin is true.
  - independent (default): depends on the independent configuration of each RAM user. The original value of EnforceMFAForLogin is false.
  - adaptive: Used only during abnormal login.
* `operation_for_risk_login` - (Optional, Computed, Available since v1.248.0) Whether MFA is verified twice during abnormal logon. Value:
  - autonomous (default): Skip, do not force binding.
  - enforceVerify: Force binding validation.
* `verification_types` - (Optional, Set, Available since v1.248.0) Means of multi-factor authentication. Value:
  - sms: secure phone.
  - email: Secure mailbox.

The following arguments will be discarded. Please use new fields as soon as possible:
* `enforce_mfa_for_login` - (Optional, Deprecated since v1.248.0) Field `enforce_mfa_for_login` has been deprecated from provider version 1.248.0. New field `mfa_operation_for_login` instead. 
Specifies whether MFA is required for all RAM users when they log on to the Alibaba Cloud Management Console by using usernames and passwords. Valid values: `true` and `false`
  - `true` - MFA is required for all RAM users when they log on to the Alibaba Cloud Management Console by using usernames and passwords.
  - `false` - User-specific settings are applied. This is the default value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as Alibaba Account ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Security Preference.
* `update` - (Defaults to 5 mins) Used when update the Security Preference.

## Import

RAM Security Preference can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_security_preference.example 
```