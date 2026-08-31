---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_account"
description: |-
  Provides a Alicloud Resource Manager Account resource.
---

# alicloud_resource_manager_account

Provides a Resource Manager Account resource.



For information about Resource Manager Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/doc-detail/111231.htm).

-> **NOTE:** Available since v1.83.0.

## Example Usage

Basic Usage

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_resource_manager_account&spm=docs.r.resource_manager_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `abandonable_check_id` - (Optional, List, Available since v1.249.0) The ID of the check item that can choose to abandon and continue to perform member deletion.
The ID is obtained from the return parameter AbandonableChecks of [GetAccountDeletionCheckResult](~~ GetAccountDeletionCheckResult ~~).
* `account_name_prefix` - (Optional, Available since v1.114.0) Account name prefix. Empty the system randomly generated.
Format: English letters, numbers, and special characters_.-can be entered. It must start and end with an English letter or number, and continuous special characters_.-cannot be entered '_.-'.
The format of the full account name is @< ResourceDirectoryId>.aliyunid.com, for example: 'alice @ rd-3G ****.aliyunid.com'
The account name must be unique in the resource directory.

* `display_name` - (Required) Member name
* `folder_id` - (Optional, Computed) The ID of the parent folder
* `payer_account_id` - (Optional) The settlement account ID. If it is left blank, the newly created member will be used for self-settlement.
* `resell_account_type` - (Optional, Available since v1.249.0) The identity type of the member. Valid values:
  - resell: The member is an account for a reseller. This is the default value. A relationship is automatically established between the member and the reseller. The management account of the resource directory must be used as the billing account of the member.
  - non_resell: The member is not an account for a reseller. The member is an account that is not associated with a reseller. You can directly use the account to purchase Alibaba Cloud resources. The member is used as its own billing account.

-> **NOTE:**  This parameter is available only for resellers at the international site (alibabacloud.com).

* `tags` - (Optional, Map) The tag of the resource
* `force_delete` - (Optional, Available since v1.249.0) Whether to force delete the account.
* `type` - (Optional, Computed) Member type. The value of ResourceAccount indicates the resource account

The following arguments will be discarded. Please use new fields as soon as possible:
* `abandon_able_check_id` - (Deprecated since v1.249.0). Field 'abandon_able_check_id' has been deprecated from provider version 1.249.0. New field 'abandonable_check_id' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `join_method` - Ways for members to join the resource directory.  invited, created
* `join_time` - The time when the member joined the resource directory
* `modify_time` - The modification time of the invitation
* `resource_directory_id` - Resource directory ID
* `status` - Member joining status.  CreateSuccess,CreateVerifying,CreateFailed,CreateExpired,CreateCancelled,PromoteVerifying,PromoteFailed,PromoteExpired,PromoteCancelled,PromoteSuccess,InviteSuccess,Removed

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

Resource Manager Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_account.example <id>
```