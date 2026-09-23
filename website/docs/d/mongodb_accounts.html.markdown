---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_accounts"
sidebar_current: "docs-alicloud-datasource-mongodb-accounts"
description: |-
  Provides a list of Mongodb Accounts to the user.
---

# alicloud\_mongodb\_accounts

This data source provides the Mongodb Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_accounts" "example" {
  instance_id  = "example_value"
  account_name = "root"
}
output "mongodb_account_id_1" {
  value = data.alicloud_mongodb_accounts.example.accounts.0.id
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `account_name` - (Optional) The name of the account. Valid values: `root`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `accounts` - A list of Mongodb Accounts. Each element contains the following attributes:
	* `account_description` - The description of the account.
	* `account_name` - The name of the account.
	* `character_type` - The role of the account. Valid values: `db`, `cs`, `mongos`, `logic`, `normal`.
	* `id` - The ID of the Account. The value formats as `<instance_id>:<account_name>`.
	* `instance_id` - The id of the instance to which the account belongs.
	* `status` - The status of the account. Valid values: `Unavailable`, `Available`.