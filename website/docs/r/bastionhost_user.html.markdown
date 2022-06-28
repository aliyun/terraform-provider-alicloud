---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_user"
sidebar_current: "docs-alicloud-resource-bastionhost-user"
description: |-
  Provides a Alicloud Bastion Host User resource.
---

# alicloud\_bastionhost\_user

Provides a Bastion Host User resource.

For information about Bastion Host User and how to use it, see [What is User](https://www.alibabacloud.com/help/doc-detail/204503.htm).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_user" "Local" {
  instance_id         = "example_value"
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Local"
  user_name           = "my-local-user"
}

resource "alicloud_bastionhost_user" "Ram" {
  instance_id         = "example_value"
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Ram"
  source_user_id      = "1234567890"
  user_name           = "my-ram-user"
}
```

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) Specify the New of the User That Created the Remark Information. Supports up to 500 Characters.
* `display_name` - (Optional, Computed) Specify the New Created the User's Display Name. Supports up to 128 Characters.
* `email` - (Optional) Specify the New User's Mailbox.
* `instance_id` - (Required, ForceNew) You Want to Query the User the Bastion Host ID of.
* `mobile` - (Optional) Specify the New of the User That Created a Different Mobile Phone Number from Your.
* `mobile_country_code` - (Optional, Computed) Specify the New Create User Mobile Phone Number of the International Domain Name. The Default Value Is the CN. Valid Values:
  * CN: Mainland China (+86) 
  * HK: hong Kong, China (+852) 
  * MO: Macau, China (+853) 
  * TW: Taiwan, China (+886) 
  * RU: Russian (+7)
  * SG: Singapore (+65) 
  * MY: malaysia (+60) 
  * ID: Indonesia (+62) 
  * DE: Germany (+49) 
  * AU: Australia (+61) 
  * US: United States (+1) 
  * AE: dubai (+971) 
  * JP: Japan (+81) Introducing the Long-Range 
  * GB: United Kingdom (+44) 
  * IN: India (+91) 
  * KR: South Korea (+82) 
  * PH: philippines (+63) 
  * CH: Switzerland (+41) 
  * SE: Sweden (+46)
  
* `password` - (Optional, Sensitive) Specify the New User's Password. Supports up to 128 Characters. Description of the New User as the Source of the Local User (That Is, Source Value for Local, this Parameter Is Required.
* `source` - (Required, ForceNew) Specify the New of the User That Created the Source. Valid Values: 
  * Local: Local User
  * RAM: Ram User
  
* `source_user_id` - (Optional, ForceNew) Specify the Newly Created User Is Uniquely Identified. Indicates That the Parameter Is a Bastion Host Corresponding to the User with the Ram User's Unique Identifier. The Newly Created User Source Grant Permission to a RAM User (That Is, Source Used as a Ram), this Parameter Is Required. You Can Call Access Control of Listusers Interface from the Return Data Userid to Obtain the Parameters.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Frozen`, `Normal`.
* `user_name` - (Required, ForceNew) Specify the New User Name. This Parameter Is Only by Letters, Lowercase Letters, Numbers, and Underscores (_), Supports up to 128 Characters.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User. The value formats as `<instance_id>:<user_id>`.
* `user_id` - The User ID.

## Import

Bastion Host User can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_user.example <instance_id>:<user_id>
```
