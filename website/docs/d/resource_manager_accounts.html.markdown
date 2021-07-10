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

```terraform
data "alicloud_resource_manager_accounts" "default" {}

output "first_account_id" {
  value = "${data.alicloud_resource_manager_accounts.default.accounts.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of account IDs.
* `status` - (Optional, ForceNew, Available in v1.114.0+) The status of account, valid values: `CreateCancelled`, `CreateExpired`, `CreateFailed`, `CreateSuccess`, `CreateVerifying`, `InviteSuccess`, `PromoteCancelled`, `PromoteExpired`, `PromoteFailed`, `PromoteSuccess`, and `PromoteVerifying`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional,  Available in v1.124.3+) Default to `false`. Set it to `true` can output more details about resource attributes.

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
    * `payer_account_id` - (Available in v1.124.3+) Settlement account ID. If the value is empty, the current account will be used for settlement.
    * `account_name` - (Available in v1.125.0+) The Alibaba Cloud account name of the member account.
    
