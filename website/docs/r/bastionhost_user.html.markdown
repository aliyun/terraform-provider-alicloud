---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_user"
sidebar_current: "docs-alicloud-resource-bastionhost-user"
description: |-
  Provides a Alicloud Bastion Host User resource.
---

# alicloud_bastionhost_user

Provides a Bastion Host User resource.

For information about Bastion Host User and how to use it, see [What is User](https://www.alibabacloud.com/help/en/bastion-host/latest/api-yundun-bastionhost-2019-12-09-createuser).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bastionhost_user&exampleId=c0feb6b6-e7e7-3a99-484d-e9457651b49fa643286a&activeTab=example&spm=docs.r.bastionhost_user.0.c0feb6b6e7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "10.4.0.0/16"
}

data "alicloud_vswitches" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = data.alicloud_vpcs.default.ids.0
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_bastionhost_instance" "default" {
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  plan_code          = "cloudbastion"
  storage            = "5"
  bandwidth          = "5"
  period             = "1"
  vswitch_id         = data.alicloud_vswitches.default.ids[0]
  security_group_ids = [alicloud_security_group.default.id]
}

resource "alicloud_bastionhost_user" "local_user" {
  instance_id         = alicloud_bastionhost_instance.default.id
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Local"
  user_name           = "${var.name}_local_user"
}

resource "alicloud_ram_user" "user" {
  name         = "${var.name}_bastionhost_user"
  display_name = "${var.name}_bastionhost_user"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
  force        = true
}
data "alicloud_account" "default" {}
resource "alicloud_bastionhost_user" "ram_user" {
  instance_id    = alicloud_bastionhost_instance.default.id
  source         = "Ram"
  source_user_id = data.alicloud_account.default.id
  user_name      = alicloud_ram_user.user.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_bastionhost_user&spm=docs.r.bastionhost_user.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) Specify the New of the User That Created the Remark Information. Supports up to 500 Characters.
* `display_name` - (Optional) Specify the New Created the User's Display Name. Supports up to 128 Characters.
* `email` - (Optional) Specify the New User's Mailbox.
* `instance_id` - (Required, ForceNew) You Want to Query the User the Bastion Host ID of.
* `mobile` - (Optional) Specify the New of the User That Created a Different Mobile Phone Number from Your.
* `mobile_country_code` - (Optional) Specify the New Create User Mobile Phone Number of the International Domain Name. The Default Value Is the CN. Valid Values:
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
* `password` - (Optional, Sensitive) Specify the New User's Password. Supports up to 128 Characters. Description of the New User as the Source of the Local User That Is, Source Value for Local, this Parameter Is Required.
* `source` - (Required, ForceNew) Specify the New of the User That Created the Source. Valid Values:
  * Local: Local User
  * Ram: Ram User
  * AD: AD-authenticated User
  * LDAP: LDAP-authenticated User
-> **NOTE:** From version 1.199.0, `source` can be set to `AD`, `LDAP`.
* `source_user_id` - (Optional, ForceNew) Specify the Newly Created User Is Uniquely Identified. Indicates That the Parameter Is a Bastion Host Corresponding to the User with the Ram User's Unique Identifier. The Newly Created User Source Grant Permission to a RAM User (That Is, Source Used as a Ram), this Parameter Is Required. You Can Call Access Control of Listusers Interface from the Return Data Userid to Obtain the Parameters.
* `status` - (Optional) The status of the resource. Valid values: `Frozen`, `Normal`.
* `user_name` - (Required, ForceNew) Specify the New User Name. This Parameter Is Only by Letters, Lowercase Letters, Numbers, and Underscores (_), Supports up to 128 Characters.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User. The value formats as `<instance_id>:<user_id>`.
* `user_id` - The User ID.

## Timeouts

-> **NOTE:** Available since v1.199.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the User.
* `update` - (Defaults to 5 mins) Used when update the User.
* `delete` - (Defaults to 5 mins) Used when delete the User.

## Import

Bastion Host User can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_user.example <instance_id>:<user_id>
```
