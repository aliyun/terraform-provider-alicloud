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
  message_contact_name = "resourceexample"
}
```

## Argument Reference

The following arguments are supported:
* `email_address` - (Required) Email address
* `message_contact_name` - (Required) The Name of Contact
* `message_types` - (Required, List) The message type, including AccountExpenses, ProductMessage, SecurityMessage, FaultMessage, ActivityMessage, and ServiceMessage.
* `phone_number` - (Optional) Phone number, such as 86-11xxxxxxxxx format
* `title` - (Required) Title

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time
* `status` - Contact Status

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