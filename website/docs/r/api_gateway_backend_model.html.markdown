---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_backend_model"
sidebar_current: "docs-alicloud-resource-api-gateway-backend-model"
description: |-
  Provides a Alicloud Api Gateway Backend Model resource.
---

# alicloud_api_gateway_backend_model

Provides a Api Gateway Backend Model resource.

For information about Api Gateway Backend Model and how to use it, see [What is Backend Model](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-createbackendmodel).

-> **NOTE:** Available since v1.279.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_api_gateway_backend" "default" {
  backend_name = var.name
  description  = var.name
  backend_type = "HTTP"
}

resource "alicloud_api_gateway_backend_model" "default" {
  backend_id   = alicloud_api_gateway_backend.default.id
  backend_type = "HTTP"
  stage_name   = "RELEASE"
  description  = var.name
  backend_model_data = jsonencode({
    ServiceAddress     = "http://apigateway.alicloudapi.com:8080"
    HttpTargetHostName = "www.example.com"
  })
}
```

## Argument Reference

The following arguments are supported:

* `backend_id` - (Required, ForceNew) The ID of the Backend.
* `backend_type` - (Required, ForceNew) The type of the Backend. Valid values: `HTTP`, `VPC`, `FC_EVENT`, `FC_EVENT_V3`, `FC_HTTP`, `FC_HTTP_V3`, `OSS`, `MOCK`, `EVENTBRIDGE`.
* `stage_name` - (Required, ForceNew) The stage name of the Backend Model. Valid values: `RELEASE`, `PRE`, `TEST`.
* `description` - (Optional) The description of the Backend Model.
* `backend_model_data` - (Required) The backend model data in JSON format. The structure varies by `backend_type`. See [CreateBackendModel](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-createbackendmodel) for more details.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backend Model. The value is formatted as `<backend_id>:<stage_name>`.
* `backend_model_id` - The ID of the Backend Model.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Backend Model.
* `delete` - (Defaults to 5 mins) Used when delete the Backend Model.
* `update` - (Defaults to 5 mins) Used when update the Backend Model.

## Import

Api Gateway Backend Model can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_backend_model.example <backend_id>:<stage_name>
```
