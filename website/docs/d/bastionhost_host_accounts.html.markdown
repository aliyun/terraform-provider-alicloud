---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_accounts"
sidebar_current: "docs-alicloud-datasource-bastionhost-host-accounts"
description: |-
  Provides a list of Bastionhost Host Accounts to the user.
---

# alicloud\_bastionhost\_host\_accounts

This data source provides the Bastionhost Host Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_host_accounts" "ids" {
  host_id     = "15"
  instance_id = "example_value"
  ids         = ["1", "2"]
}
output "bastionhost_host_account_id_1" {
  value = data.alicloud_bastionhost_host_accounts.ids.accounts.0.id
}

data "alicloud_bastionhost_host_accounts" "nameRegex" {
  host_id     = "15"
  instance_id = "example_value"
  name_regex  = "^my-HostAccount"
}
output "bastionhost_host_account_id_2" {
  value = data.alicloud_bastionhost_host_accounts.nameRegex.accounts.0.id
}

```

## Argument Reference

The following arguments are supported:

* `host_account_name` - (Optional, ForceNew) Specify the new hosting account's name, support the longest 128 characters.
* `host_id` - (Required, ForceNew) Specifies the database where you want to create your hosting account's host ID.
* `ids` - (Optional, ForceNew, Computed)  A list of Host Account IDs.
* `instance_id` - (Required, ForceNew) Specifies the database where you want to create your hosting account's host bastion host ID of.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Host Account name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `protocol_name` - (Optional, ForceNew) Specify the new hosting account of the agreement name. Valid values: USING SSH and RDP.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Host Account names.
* `accounts` - A list of Bastionhost Host Accounts. Each element contains the following attributes:
	* `has_password` - Whether to set a new password.
	* `host_account_id` - Hosting account ID.
	* `host_account_name` - Specify the new hosting account's name, support the longest 128 characters.
	* `host_id` - Specifies the database where you want to create your hosting account's host ID.
	* `id` - The ID of the Host Account.
	* `instance_id` - Specifies the database where you want to create your hosting account's host bastion host ID of.
	* `private_key_fingerprint` - The situation where the private keys of the fingerprint information.
	* `protocol_name` - Specify the new hosting account of the agreement name. Valid values: USING SSH and RDP.
