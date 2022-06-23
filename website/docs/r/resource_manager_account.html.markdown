---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_account"
sidebar_current: "docs-alicloud-resource-resource-manager-account"
description: |-
  Provides a Resource Manager Account resource.
---

# alicloud\_resource\_manager\_account

Provides a Resource Manager Account resource. Member accounts are containers for resources in a resource directory. These accounts isolate resources and serve as organizational units in the resource directory. You can create member accounts in a folder and then manage them in a unified manner.
For information about Resource Manager Account and how to use it, see [What is Resource Manager Account](https://www.alibabacloud.com/help/en/doc-detail/111231.htm).

-> **NOTE:** Available in v1.83.0+.

## Example Usage

```terraform
# Add a Resource Manager Account.
resource "alicloud_resource_manager_folder" "f1" {
  folder_name = "test1"
}

resource "alicloud_resource_manager_account" "example" {
  display_name = "RDAccount"
  folder_id    = alicloud_resource_manager_folder.f1.id
}
```
## Argument Reference

The following arguments are supported:

* `account_name_prefix` - (Optional, ForceNew, Available in v1.114.0) The name prefix of account.
* `display_name` - (Required) Member name. The length is 2 ~ 50 characters or Chinese characters, which can include Chinese characters, English letters, numbers, underscores (_), dots (.) And dashes (-).
* `folder_id` - (Optional) The ID of the parent folder.
* `payer_account_id` - (Optional, ForceNew) The ID of the billing account. If you leave this parameter empty, the current account is used as the billing account.

-> **NOTE:** The member name must be unique within the resource directory.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of Resource Manager Account.  
* `join_method` - Ways for members to join the resource directory. Valid values: `invited`, `created`.
* `join_time` - The time when the member joined the resource directory.
* `modify_time` - The modification time of the invitation.
* `resource_directory_id` - Resource directory ID.
* `status` - Member joining status. Valid values: `CreateSuccess`,`CreateVerifying`,`CreateFailed`,`CreateExpired`,`CreateCancelled`,`PromoteVerifying`,`PromoteFailed`,`PromoteExpired`,`PromoteCancelled`,`PromoteSuccess`,`InviteSuccess`,`Removed`. 
* `type` - Member type. The value of `ResourceAccount` indicates the resource account. 

## Import

Resource Manager Account can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_account.example 13148890145*****
```
