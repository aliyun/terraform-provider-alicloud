---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine_route"
description: |-
  Provides a Alicloud ESA Routine Route resource.
---

# alicloud_esa_routine_route

Provides a ESA Routine Route resource.



For information about ESA Routine Route and how to use it, see [What is Routine Route](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRoutineRoute).

-> **NOTE:** Available since v1.251.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_routine" "default" {
  description = "example-routine2"
  name        = "example-routine2"
}

resource "alicloud_esa_routine_route" "default" {
  route_enable = "on"
  rule         = "(http.host eq \"video.example1.com\")"
  sequence     = "1"
  routine_name = alicloud_esa_routine.default.name
  site_id      = alicloud_esa_site.default.id
  bypass       = "off"
  route_name   = "example_routine"
}
```

## Argument Reference

The following arguments are supported:
* `bypass` - (Optional) Bypass mode. Value range:
  - on: Open
  - off: off
* `route_enable` - (Optional) Routing switch. Value range:
  - on: Open
  - off: off
* `route_name` - (Optional) The route name.
* `routine_name` - (Required, ForceNew) The edge function Routine name.
* `rule` - (Optional) The rule content.
* `sequence` - (Optional, ForceNew, Int) Rule execution order.
* `site_id` - (Required, ForceNew, Int) Site Id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<routine_name>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Routine Route.
* `delete` - (Defaults to 5 mins) Used when delete the Routine Route.
* `update` - (Defaults to 5 mins) Used when update the Routine Route.

## Import

ESA Routine Route can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_routine_route.example <site_id>:<routine_name>:<config_id>
```