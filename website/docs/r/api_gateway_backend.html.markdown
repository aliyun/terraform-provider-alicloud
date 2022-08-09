---
subcategory: "API Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_backend"
sidebar_current: "docs-alicloud-resource-api-gateway-backend"
description: |-
  Provides a Alicloud Api Gateway Backend resource.
---

# alicloud\_api\_gateway\_backend

Provides a Api Gateway Backend resource.

For information about Api Gateway Backend and how to use it, see [What is Backend](https://www.alibabacloud.com/help/zh/api-gateway/latest/api-doc-cloudapi-2016-07-14-api-doc-createbackend).

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
variable "name1" {
  default = "tf-testAccBackend"
}

resource "alicloud_api_gateway_backend" "default" {
  backend_name = var.name
  description  = var.name
  backend_type = "HTTP"
}
```

## Argument Reference

The resource does not support any argument.
* `backend_type` - (Required, ForceNew) The type of the Backend. Valid values: `HTTP`, `VPC`, `FC_EVENT`, `FC_HTTP`, `OSS`, `MOCK`.
* `backend_name` - (Required) The name of the Backend.
* `create_event_bridge_service_linked_role` - (Optional, ForceNew) Whether to create an Event bus service association role.
* `description` - (Optional) The description of the Backend.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backend.

## Import

Api Gateway Backend can be imported using the id, e.g.

```
$ terraform import alicloud_api_gateway_backend.example <id>
```