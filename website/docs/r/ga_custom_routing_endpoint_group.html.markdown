---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_group"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint-group"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint Group resource.
---

# alicloud_ga_custom_routing_endpoint_group

Provides a Global Accelerator (GA) Custom Routing Endpoint Group resource.

For information about Global Accelerator (GA) Custom Routing Endpoint Group and how to use it, see [What is Custom Routing Endpoint Group](https://www.alibabacloud.com/help/en/global-accelerator/latest/createcustomroutingendpointgroups).

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region = var.region
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  listener_type  = "CustomRouting"
  port_ranges {
    from_port = 10000
    to_port   = 16000
  }
}

resource "alicloud_ga_custom_routing_endpoint_group" "default" {
  accelerator_id                     = alicloud_ga_listener.default.accelerator_id
  listener_id                        = alicloud_ga_listener.default.id
  endpoint_group_region              = var.region
  custom_routing_endpoint_group_name = "terraform-example"
  description                        = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `listener_id` - (Required, ForceNew) The ID of the custom routing listener.
* `endpoint_group_region` - (Required, ForceNew) The ID of the region in which to create the endpoint group.
* `custom_routing_endpoint_group_name` - (Optional) The name of the endpoint group.
* `description` - (Optional) The description of the endpoint group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint Group.
* `status` - The status of the Custom Routing Endpoint Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint Group.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint Group.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint Group.

## Import

Global Accelerator (GA) Custom Routing Endpoint Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint_group.example <id>
```
