---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_account"
sidebar_current: "docs-alicloud-resource-resource-manager-account"
description: |-
  Provides a Resource Manager Account resource.
---

# alicloud_resource_manager_account

Provides a Resource Manager Account resource. Member accounts are containers for resources in a resource directory. These accounts isolate resources and serve as organizational units in the resource directory. You can create member accounts in a folder and then manage them in a unified manner.
For information about Resource Manager Account and how to use it, see [What is Resource Manager Account](https://www.alibabacloud.com/help/en/doc-detail/111231.htm).

-> **NOTE:** Available since v1.83.0.

-> **NOTE:** From version 1.188.0, the resource can be destroyed. The member deletion feature is in invitational preview. You can contact the service manager of Alibaba Cloud to apply for a trial. see [how to destroy it](https://www.alibabacloud.com/help/en/resource-management/latest/delete-account).

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_account&exampleId=a3c39ae4-d702-ffd5-b41b-53e377262ab15bcbe992&activeTab=example&spm=docs.r.resource_manager_account.0.a3c39ae4d7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
variable "display_name" {
  default = "EAccount"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_resource_manager_folders" "example" {

}

resource "alicloud_resource_manager_account" "example" {
  display_name = "${var.display_name}-${random_integer.default.result}"
  folder_id    = data.alicloud_resource_manager_folders.example.ids.0
}
```

### Deleting `alicloud_resource_manager_account` or removing it from your configuration

Deleting the resource manager account or removing it from your configuration will remove it from your state file and management, 
but may not destroy the account. If there are some dependent resource in the account, 
the deleting account will enter a silence period of 45 days. After the silence period ends, 
the system automatically starts to delete the member. [See More Details](https://www.alibabacloud.com/help/en/resource-management/latest/delete-resource-account).

## Argument Reference

The following arguments are supported:

* `account_name_prefix` - (Optional, ForceNew, Available since v1.114.0) The name prefix of account.
* `display_name` - (Required) Member name. The length is 2 ~ 50 characters or Chinese characters, which can include Chinese characters, English letters, numbers, underscores (_), dots (.) And dashes (-).
* `folder_id` - (Optional) The ID of the parent folder.
* `payer_account_id` - (Optional, ForceNew) The ID of the billing account. If you leave this parameter empty, the current account is used as the billing account.
* `abandon_able_check_id` - (Optional, Available in v1.188.0+) The IDs of the check items that you can choose to ignore for the member deletion. 
  If you want to delete the account, please use datasource `alicloud_resource_manager_account_deletion_check_task` 
  to get check ids and set them.
* `tags` - (Optional, Available since v1.181.0) A mapping of tags to assign to the resource.

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

## Timeouts

-> **NOTE:** Available since v1.188.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Resource Manager Account.
* `update` - (Defaults to 3 mins) Used when update the Resource Manager Account.
* `delete` - (Defaults to 3 mins) Used when delete the Resource Manager Account.

## Import

Resource Manager Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_account.example 13148890145*****
```
