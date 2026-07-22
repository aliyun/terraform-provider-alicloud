---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_contact"
description: |-
  Provides a Alicloud Certificate Management Service (Original SSL Certificate) Contact resource.
---

# alicloud_ssl_certificates_service_contact

Provides a Certificate Management Service (Original SSL Certificate) Contact resource.

Certificate Contact Person.

For information about Certificate Management Service (Original SSL Certificate) Contact and how to use it, see [What is Contact](https://next.api.alibabacloud.com/document/cas/2020-04-07/CreateContact).

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_ssl_certificates_service_contact" "default" {
  name   = var.name
  mobile = "13312345678"
  email  = "test@example.com"
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The name of the contact.
* `mobile` - (Required) The mobile phone number of the contact.
* `email` - (Optional) The email address of the contact.
* `idcard` - (Optional) The ID card number of the contact. This is only required by the CFCA certificate brand and is not accepted for other brands.

-> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `webhooks` - (Optional) The Webhook addresses of DingTalk, WeCom, or Lark bots used to receive notifications, formatted as a JSON list string, e.g. `["https://example.com/webhook"]`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Contact.
* `delete` - (Defaults to 5 mins) Used when delete the Contact.
* `update` - (Defaults to 5 mins) Used when update the Contact.

## Import

Certificate Management Service (Original SSL Certificate) Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_contact.example <contact_id>
```