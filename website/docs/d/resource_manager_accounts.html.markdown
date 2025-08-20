---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_accounts"
description: |-
  Provides a list of Resource Manager Accounts to the user.
---

# alicloud_resource_manager_accounts

This data source provides the Resource Manager Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.86.0.

## Example Usage

```terraform
data "alicloud_resource_manager_accounts" "default" {
}

output "resource_manager_account_id_0" {
  value = data.alicloud_resource_manager_accounts.default.accounts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Account IDs.
* `status` - (Optional, ForceNew, Available since v1.114.0) The status of account. Valid values: `CreateCancelled`, `CreateExpired`, `CreateFailed`, `CreateSuccess`, `CreateVerifying`, `InviteSuccess`, `PromoteCancelled`, `PromoteExpired`, `PromoteFailed`, `PromoteSuccess`, `PromoteVerifying`.
* `tags` - (Optional, ForceNew, Available since v1.259.0) A mapping of tags to assign to the resource.
* `enable_details` - (Optional, Bool, Available since v1.124.3) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `accounts` - A list of accounts. Each element contains the following attributes:
  * `id` - The ID of the Account.
  * `account_id`- The Alibaba Cloud account ID of the member.
  * `display_name`- The display name of the member.
  * `type` - The type of the member.
  * `folder_id` - The ID of the folder.
  * `resource_directory_id` - The ID of the resource directory.
  * `status` - The status of the member.
  * `tags` - (Available since v1.259.0) The tags that are added to the member.
  * `account_name` - (Available since v1.125.0) The Alibaba Cloud account name of the member. **Note:** `account_name` takes effect only if `enable_details` is set to `true`.
  * `payer_account_id` - (Available since v1.124.3) The ID of the settlement account. **Note:** `payer_account_id` takes effect only if `enable_details` is set to `true`.
  * `join_method` - The way in which the member joins the resource directory.
  * `join_time` - The time when the member joined the resource directory.
  * `modify_time` - The time when the member was modified.
