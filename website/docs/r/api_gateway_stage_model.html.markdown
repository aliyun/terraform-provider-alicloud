---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_stage_model"
sidebar_current: "docs-alicloud-resource-api-gateway-stage-model"
description: |-
  Provides a Alicloud Api Gateway Stage Model resource.
---

# alicloud_api_gateway_stage_model

Provides a Api Gateway Stage Model resource.

For information about Api Gateway Stage Model and how to use it, see [What is Stage Model](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-createstagemodel).

-> **NOTE:** Available since v1.280.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_stage_model" "default" {
  stage_model_name  = "DEVELOP"
  stage_model_alias = "Develop Environment"
  description       = "Develop stage for testing"
}
```

## Argument Reference

The following arguments are supported:

* `stage_model_name` - (Required, ForceNew) The name of the Stage Model. It must be 2-10 uppercase letters or digits. Valid values: cannot be `RELEASE`, `PRE`, or `TEST`.
* `stage_model_alias` - (Required) The alias of the Stage Model.
* `description` - (Optional) The description of the Stage Model.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Stage Model. The value is the `stage_model_id`.
* `stage_model_id` - The ID of the Stage Model.
* `type` - The type of the Stage Model. Valid values: `SYSTEM`, `CUSTOM`.
* `created_time` - The creation time of the Stage Model.
* `modified_time` - The last modified time of the Stage Model.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Stage Model.
* `delete` - (Defaults to 5 mins) Used when delete the Stage Model.
* `update` - (Defaults to 5 mins) Used when update the Stage Model.

## Import

Api Gateway Stage Model can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_stage_model.example <stage_model_id>
```
