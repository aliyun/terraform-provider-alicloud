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

-> **NOTE:** From version 1.188.0, the resource can be destroyed. The member deletion feature is in invitational preview. You can contact the service manager of Alibaba Cloud to apply for a trial. see [how to destroy it](https://www.alibabacloud.com/help/en/resource-management/latest/delete-account).

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
* `abandon_able_check_id` - (Optional, ForceNew, Available in v1.188.0+) The IDs of the check items that you can choose to ignore for the member deletion. You can obtain the IDs from the datasource `alicloud_resource_manager_account_deletion_check_task` operation.
* `tags` - (Optional, Available in v1.181.0+) A mapping of tags to assign to the resource.

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

### Timeouts

-> **NOTE:** Available in 1.188.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Resource Manager Account.
* `update` - (Defaults to 3 mins) Used when update the Resource Manager Account.
* `delete` - (Defaults to 3 mins) Used when delete the Resource Manager Account.

## Import

Resource Manager Account can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_account.example 13148890145*****
```
