---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_accounts"
sidebar_current: "docs-alicloud-datasource-kvstore-accounts"
description: |-
  Provides a list of KVStore Accounts to the user.
---

# alicloud\_kvstore\_accounts

This data source provides the KVStore Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.102.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_kvstore_accounts" "example" {
  instance_id = "example_value"
}

output "first_kvstore_account_id" {
  value = data.alicloud_kvstore_accounts.example.accounts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Optional, ForceNew) The name of the account.
* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs.
* `status` - (Optional, ForceNew) The status of KVStore Account. Valid Values: `"Available` `Unavailable`
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Account names.
* `ids` - A list of Account IDs.
* `accounts` - A list of Kvstore Accounts. Each element contains the following attributes:
	* `id` - The ID of the Account.
	* `account_name` - The name of the account.
	* `account_privilege` - The privilege of account access database.
	* `account_type` - Privilege type of account.
	* `description` - The description of account.
	* `instance_id` - The Id of instance in which account belongs.
	* `status` - The status of account.
