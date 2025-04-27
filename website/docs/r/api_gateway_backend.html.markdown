---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_backend"
sidebar_current: "docs-alicloud-resource-api-gateway-backend"
description: |-
  Provides a Alicloud Api Gateway Backend resource.
---

# alicloud_api_gateway_backend

Provides a Api Gateway Backend resource.

For information about Api Gateway Backend and how to use it, see [What is Backend](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-createbackend).

-> **NOTE:** Available since v1.181.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_backend&exampleId=7069f394-4f0b-9a56-8915-e0d9ff969442eb4d7f9d&activeTab=example&spm=docs.r.api_gateway_backend.0.7069f3944f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_api_gateway_backend" "default" {
  backend_name = var.name
  description  = var.name
  backend_type = "HTTP"
}
```

## Argument Reference

The following arguments are supported:

* `backend_type` - (Required, ForceNew) The type of the Backend. Valid values: `HTTP`, `VPC`, `FC_EVENT`, `FC_EVENT_V3`, `FC_HTTP`, `FC_HTTP_V3`, `OSS`, `MOCK`.
* `backend_name` - (Required) The name of the Backend.
* `create_event_bridge_service_linked_role` - (Optional, ForceNew) Whether to create an Event bus service association role.
* `description` - (Optional) The description of the Backend.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backend.

## Import

Api Gateway Backend can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_backend.example <id>
```