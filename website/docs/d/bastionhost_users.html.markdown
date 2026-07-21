---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_users"
sidebar_current: "docs-alicloud-datasource-bastionhost-users"
description: |-
  Provides a list of Bastionhost Users to the user.
---

# alicloud\_bastionhost\_users

This data source provides the Bastionhost Users of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_users" "ids" {
  instance_id = "example_value"
  ids         = ["1", "10"]
}
output "bastionhost_user_id_1" {
  value = data.alicloud_bastionhost_users.ids.users.0.id
}

data "alicloud_bastionhost_users" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-User"
}
output "bastionhost_user_id_2" {
  value = data.alicloud_bastionhost_users.nameRegex.users.0.id
}

```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional, ForceNew) Specify the New Created the User's Display Name. Supports up to 128 Characters.
* `ids` - (Optional, ForceNew, Computed)  A list of User IDs.
* `instance_id` - (Required, ForceNew) You Want to Query the User the Bastion Host ID of.
* `mobile` - (Optional, ForceNew) Specify the New of the User That Created a Different Mobile Phone Number from Your.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by User name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `source` - (Optional, ForceNew) Specify the New of the User That Created the Source. Valid Values: Local: Local User RAM: Ram User. Valid values: `Local`, `Ram`.
* `source_user_id` - (Optional, ForceNew) Specify the Newly Created User Is Uniquely Identified. Indicates That the Parameter Is a Bastion Host Corresponding to the User with the Ram User's Unique Identifier. The Newly Created User Source Grant Permission to a RAM User (That Is, Source Used as a Ram), this Parameter Is Required. You Can Call Access Control of Listusers Interface from the Return Data Userid to Obtain the Parameters.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Frozen`, `Normal`.
* `user_name` - (Optional, ForceNew) Specify the New User Name. This Parameter Is Only by Letters, Lowercase Letters, Numbers, and Underscores (_), Supports up to 128 Characters.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of User names.
* `users` - A list of Bastionhost Users. Each element contains the following attributes:
	* `comment` - Specify the New of the User That Created the Remark Information. Supports up to 500 Characters.
	* `display_name` - Specify the New Created the User's Display Name. Supports up to 128 Characters.
	* `email` - Specify the New User's Mailbox.
	* `id` - The ID of the User.
	* `instance_id` - You Want to Query the User the Bastion Host ID of.
	* `mobile` - Specify the New of the User That Created a Different Mobile Phone Number from Your.
	* `mobile_country_code` - Specify the New Create User Mobile Phone Number of the International Domain Name. The Default Value Is the CN Value: CN: Mainland China (+86) HK: hong Kong, China (+852) Mo: Macau, China (+853) TW: Taiwan, China (+886) ru: Russian (+7) SG: Singapore (+65) My: malaysia (+60) ID: Indonesia (+62) De: Germany (+49) AU: Australia (+61) US: United States (+1) AE: dubai (+971) JP: Japan (+81) Introducing the Long-Range GB: United Kingdom (+44) in: India (+91) KR: South Korea (+82) Ph: philippines (+63) Ch: Switzerland (+41) Se: Sweden (+46).
	* `source` - Specify the New of the User That Created the Source. Valid Values: Local: Local User RAM: Ram User.
	* `source_user_id` - Specify the Newly Created User Is Uniquely Identified. Indicates That the Parameter Is a Bastion Host Corresponding to the User with the Ram User's Unique Identifier. The Newly Created User Source Grant Permission to a RAM User (That Is, Source Used as a Ram), this Parameter Is Required. You Can Call Access Control of Listusers Interface from the Return Data Userid to Obtain the Parameters.
	* `status` - The status of the resource.
	* `user_id` - The User ID.
	* `user_name` - Specify the New User Name. This Parameter Is Only by Letters, Lowercase Letters, Numbers, and Underscores (_), Supports up to 128 Characters.
