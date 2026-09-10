---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_accounts"
sidebar_current: "docs-alicloud-datasource-rds-accounts"
description: |-
  Provides a list of Rds Accounts to the user.
---

# alicloud\_rds\_accounts

This data source provides the Rds Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_accounts" "example" {
  db_instance_id = "example_value"
  name_regex     = "the_resource_name"
}

output "first_rds_account_id" {
  value = data.alicloud_rds_accounts.example.accounts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of rds account IDs.
* `db_instance_id` - (Required, ForceNew) The db instance id.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Account name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Unavailable`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Account IDs.
* `names` - A list of Account names.
* `accounts` - A list of Rds Accounts. Each element contains the following attributes:
	* `account_description` - Database description.
	* `account_name` - Name of database account.
	* `account_type` - Privilege type of account.
	* `database_privileges` - A list of database permissions the account has.
		* `account_privilege` - The type of permission for the account.
		* `account_privilege_detail` - The specific permissions corresponding to the type of account permissions.
		* `db_name` - Database name.
	* `id` - The ID of the Account.
	* `priv_exceeded` - Whether the maximum number of databases managed by the account is exceeded.
	* `status` - The status of the resource.
