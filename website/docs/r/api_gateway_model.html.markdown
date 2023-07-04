---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_model"
sidebar_current: "docs-alicloud-resource-api-gateway-model"
description: |-
  Provides a Alicloud Api Gateway Model resource.
---

# alicloud_api_gateway_model

Provides a Api Gateway Model resource.

For information about Api Gateway Model and how to use it, see [What is Model](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-createmodel).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_group" "default" {
  name        = "example_value"
  description = "example_value"
}

resource "alicloud_api_gateway_model" "default" {
  group_id    = alicloud_api_gateway_group.default.id
  model_name  = "example_value"
  schema      = "{\"type\":\"object\",\"properties\":{\"id\":{\"format\":\"int64\",\"maximum\":100,\"exclusiveMaximum\":true,\"type\":\"integer\"},\"name\":{\"maxLength\":10,\"type\":\"string\"}}}"
  description = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The group of the model belongs to.
* `model_name` - (Required, ForceNew) The name of the model.
* `schema` - (Required) The schema of the model.
* `description` - (Optional) The description of the model.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Model. The value formats as `<group_id>:<model_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Api Gateway Model.
* `update` - (Defaults to 3 mins) Used when update the Api Gateway Model.
* `delete` - (Defaults to 3 mins) Used when delete the Api Gateway Model.

## Import

Api Gateway Model can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_model.example <group_id>:<model_name>
```