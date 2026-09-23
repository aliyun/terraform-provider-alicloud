---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_account_deletion_check_task"
sidebar_current: "docs-alicloud-datasource-resource-manager-account-deletion-check-task"
description: |-
  Provides a datasource to open the Resource Manager Account Deletion Check Task.
---

# alicloud\_resource\_manager\_account\_deletion\_check\_task

Using this data source can open Resource Manager Account Deletion Check Task.

For information about Resource Manager Account Deletion Check Task and how to use it, see [What is Resource Manager Account Deletion Check Task](https://www.alibabacloud.com/help/en/resource-management/latest/check-account-delete).

-> **NOTE:** Available in v1.187.0+.

-> **NOTE:** The member deletion feature is in invitational preview. You can contact the service manager of Alibaba Cloud to apply for a trial.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_account_deletion_check_task" "task" {
  account_id = "your_account_id"
}

output "abandon_able_checks_ids" {
  value = data.alicloud_resource_manager_account_deletion_check_task.task.abandon_able_checks.*.check_id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, ForceNew) The ID of the member that you want to delete.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - This ID of Resource Manager Account Deletion Check Task.
* `not_allow_reason` - The reasons why the member cannot be deleted. Each element contains the following attributes:
  * `check_id` - The ID of the check item.
  * `check_name` - The name of the cloud service to which the check item belongs.
  * `description` - The description of the check item.
* `abandon_able_checks` - The check items that you can choose to ignore for the member deletion. Each element contains the following attributes:
  * `check_id` - The ID of the check item.
  * `check_name` - The name of the cloud service to which the check item belongs.
  * `description` - The description of the check item.
* `allow_delete` - Indicates whether the member can be deleted.
* `status` - The status of the check.