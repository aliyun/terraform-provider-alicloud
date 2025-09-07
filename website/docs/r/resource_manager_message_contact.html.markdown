---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_message_contact"
description: |-
  Provides a Alicloud Resource Manager Message Contact resource.
---

# alicloud_resource_manager_message_contact

Provides a Resource Manager Message Contact resource.

Message contact for Resource Directory account.

For information about Resource Manager Message Contact and how to use it, see [What is Message Contact](https://next.api.alibabacloud.com/document/ResourceDirectoryMaster/2022-04-19/AddMessageContact).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_message_contact&exampleId=803bade7-4789-3569-46c5-d80d7014dfe24c3f16a3&activeTab=example&spm=docs.r.resource_manager_message_contact.0.803bade747&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_resource_manager_message_contact" "default" {
  message_types        = ["AccountExpenses"]
  phone_number         = "86-18626811111"
  title                = "TechnicalDirector"
  email_address        = "resourceexample@126.com"
  message_contact_name = "example"
}
```

## Argument Reference

The following arguments are supported:
* `email_address` - (Required) The email address of the contact.
After you specify an email address, you need to call [SendEmailVerificationForMessageContact](~~SendEmailVerificationForMessageContact~~) to send verification information to the email address. After the verification is passed, the email address takes effect.
* `message_contact_name` - (Required) The name of the contact.
The name must be unique in your resource directory.
The name must be 2 to 12 characters in length and can contain only letters.
* `message_types` - (Required, List) The types of messages received by the contact.
* `phone_number` - (Optional) The mobile phone number of the contact.

Specify the mobile phone number in the `-` format.

-> **NOTE:**  Only mobile phone numbers in the `86-` format in the Chinese mainland are supported.

After you specify a mobile phone number, you need to call [SendPhoneVerificationForMessageContact](~~SendPhoneVerificationForMessageContact~~) to send verification information to the mobile phone number. After the verification is passed, the mobile phone number takes effect.
* `title` - (Required) The job title of the contact.Valid values:
  - FinanceDirector
  - TechnicalDirector
  - MaintenanceDirector
  - CEO
  - ProjectDirector
  - Other

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the contact was created.
* `status` - The status of the contact. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Message Contact.
* `delete` - (Defaults to 5 mins) Used when delete the Message Contact.
* `update` - (Defaults to 5 mins) Used when update the Message Contact.

## Import

Resource Manager Message Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_message_contact.example <id>
```