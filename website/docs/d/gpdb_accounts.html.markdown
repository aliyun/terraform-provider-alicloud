---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_accounts"
sidebar_current: "docs-alicloud-datasource-gpdb-accounts"
description: |-
  Provides a list of Gpdb Accounts to the user.
---

# alicloud\_gpdb\_accounts

This data source provides the Gpdb Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_gpdb_accounts" "ids" {
  db_instance_id = "example_value"
  ids            = ["my-Account-1", "my-Account-2"]
}
output "gpdb_account_id_1" {
  value = data.alicloud_gpdb_accounts.ids.accounts.0.id
}

data "alicloud_gpdb_accounts" "nameRegex" {
  db_instance_id = "example_value"
  name_regex     = "^my-Account"
}
output "gpdb_account_id_2" {
  value = data.alicloud_gpdb_accounts.nameRegex.accounts.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the instance.
* `ids` - (Optional, ForceNew, Computed)  A list of Account IDs. Its element value is same as Account Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Account name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the account. Valid values: `Active`, `Creating` and `Deleting`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Account names.
* `accounts` - A list of Gpdb Accounts. Each element contains the following attributes:
	* `account_description` - The description of the account.
	* `account_name` - The name of the account.
	* `db_instance_id` - The ID of the instance.
	* `id` - The ID of the Account. Its value is same as Queue Name.
	* `status` - The status of the account. Valid values: `Active`, `Creating` and `Deleting`.
