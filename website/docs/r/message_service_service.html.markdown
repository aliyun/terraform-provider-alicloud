---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_service"
description: |-
  Provides a Alicloud Message Service Service resource.
---

# alicloud_message_service_service

Provides a Message Service Service resource.

MNS Service Open Status.

For information about Message Service Service and how to use it, see [What is Service](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance).

-> **NOTE:** Available since v1.252.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_message_service_service" "default" {
}
```

### Deleting `alicloud_message_service_service` or removing it from your configuration

Terraform cannot destroy resource `alicloud_message_service_service`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service.
