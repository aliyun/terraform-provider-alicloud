---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_contacts"
sidebar_current: "docs-alicloud-datasource-ssl-certificates-service-contacts"
description: |-
  Provides a list of Ssl Certificates Service Contact owned by an Alibaba Cloud account.
---

# alicloud_ssl_certificates_service_contacts

This data source provides Ssl Certificates Service Contact available to the user.[What is Contact](https://next.api.alibabacloud.com/document/cas/2020-04-07/CreateContact)

-> **NOTE:** Available since v1.285.0.

## Example Usage

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

data "alicloud_ssl_certificates_service_contacts" "default" {
  ids  = ["${alicloud_ssl_certificates_service_contact.default.id}"]
  name = var.name
}

output "alicloud_ssl_certificates_service_contact_example_id" {
  value = data.alicloud_ssl_certificates_service_contacts.default.contacts.0.id
}
```

## Argument Reference

The following arguments are supported:
* `name` - (ForceNew, Optional) The name of the resource
* `ids` - (Optional, Computed) A list of Contact IDs. 
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Contact IDs.
* `contacts` - A list of Contact Entries. Each element contains the following attributes:
    * `contact_id` - The first ID of the resource.
    * `email` - The email address of the contact. **NOTE:** This field is only available when `enable_details` is `true`.
    * `mobile` - The mobile phone number of the contact. **NOTE:** This field is only available when `enable_details` is `true`.
    * `name` - The name of the contact.
    * `webhooks` - The Webhook address used to receive notifications. **NOTE:** This field is only available when `enable_details` is `true`.
    * `id` - The ID of the resource supplied above.
