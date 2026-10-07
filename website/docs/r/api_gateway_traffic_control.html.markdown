---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_traffic_control"
description: |-
  Provides a Alicloud Api Gateway Traffic Control resource.
---

# alicloud_api_gateway_traffic_control

Provides a Api Gateway Traffic Control resource. Flow control.

For information about Api Gateway Traffic Control and how to use it, see [What is Traffic Control](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-createtrafficcontrol).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_api_gateway_traffic_control" "default" {
  traffic_control_name = var.name
  traffic_control_unit = "MINUTE"
  api_default          = 100
  user_default         = 50
  app_default          = 30
  description          = "traffic control for demo"
}
```

## Argument Reference

The following arguments are supported:
* `traffic_control_name` - (Required) The name of the Traffic Control. It must be unique in the current region. The name must be 1 to 64 characters in length, and can contain letters, digits, underscores (_), and hyphens (-).
* `traffic_control_unit` - (Required) The unit of the Traffic Control. Valid values: `MINUTE`, `HOUR`, `DAY`.
* `api_default` - (Required) The default flow control value for each API.
* `user_default` - (Optional) The default flow control value for each user.
* `app_default` - (Optional) The default flow control value for each app.
* `description` - (Optional) The description of the Traffic Control.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `special_policies` - The special flow control policies. See [`special_policies`](#special_policies) below.
* `create_time` - The creation time of the Traffic Control.
* `modified_time` - The last modification time of the Traffic Control.
* `region_id` - The region ID of the Traffic Control.

### `special_policies`

The special_policies supports the following:
* `special_type` - The type of the special flow control. Valid values: `APP`, `USER`.
* `specials` - The list of special flow control items. See [`specials`](#special_policies-specials) below.

### `special_policies-specials`

The specials supports the following:
* `traffic_value` - The flow control value of the special.
* `special_key` - The AppId or user account corresponding to the SpecialType.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Traffic Control.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Control.
* `update` - (Defaults to 5 mins) Used when update the Traffic Control.

## Import

Api Gateway Traffic Control can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_traffic_control.example <id>
```
