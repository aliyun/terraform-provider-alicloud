---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_accounts"
sidebar_current: "docs-alicloud-datasource-resource-manager-accounts"
description: |-
    Provides a list of Resource Manager Accounts to the user.
---

# alicloud\_resource\_manager\_accounts

This data source provides the Resource Manager Accounts of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.86.0+.

## Example Usage

```
data "alicloud_resource_manager_accounts" "default" {}

output "first_account_id" {
  value = "${data.alicloud_resource_manager_accounts.default.accounts.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of account IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of account IDs.
* `accounts` - A list of accounts. Each element contains the following attributes:
    * `id` - The ID of the resource.
    * `account_id`- The ID of the account.
    * `display_name`- The name of the member account.
    * `folder_id` - The ID of the folder.
    * `join_method` - The way in which the member account joined the resource directory. 
    * `join_time` - The time when the member account joined the resource directory.
    * `modify_time` - The time when the member account was modified.
    * `resource_directory_id` - The ID of the resource directory.
    * `status` - The status of the member account. 
    * `type` - The type of the member account. 
    
