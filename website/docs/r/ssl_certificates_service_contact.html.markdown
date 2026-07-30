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
  name         = var.name
  mobile       = "13312345678"
  email        = "test@example.com"
  webhook_list = ["https://oapi.dingtalk.com/robot/send?access_token=xxxx"]
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The name of the contact.
* `mobile` - (Required, Sensitive) The mobile phone number of the contact.
* `email` - (Optional, Sensitive) The email address of the contact.
* `idcard` - (Optional, Sensitive) The ID card number of the contact. CFCA only consumes the idcard when the contact is used to apply for a CFCA-brand certificate; other brands ignore it. CreateContact/UpdateContact accept idcard on any account.
* `webhook_list` - (Optional) A list of DingTalk robot Webhook URLs used to receive notifications, e.g. `["https://oapi.dingtalk.com/robot/send?access_token=xxx"]`. Each element must be a DingTalk robot URL (`https://oapi.dingtalk.com/robot/send?access_token=...`); other URLs are rejected by the API.

-> **NOTE:** `email`, `mobile`, and `idcard` are sensitive. The CAS API returns these fields server-masked on read (e.g. `tes****1@example.com`, `133******78`), so they are not reconciled from the API after apply: the configured value is authoritative, out-of-band changes (e.g. editing in the console) are not detectable by Terraform, and `terraform import` will not populate them (the first plan after import reports them being set). `webhook_list` is read back in plaintext from the API.

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