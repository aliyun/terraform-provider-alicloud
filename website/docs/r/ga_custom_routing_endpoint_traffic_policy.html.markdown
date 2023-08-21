---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_traffic_policy"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint-traffic-policy"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint Traffic Policy resource.
---

# alicloud_ga_custom_routing_endpoint_traffic_policy

Provides a Global Accelerator (GA) Custom Routing Endpoint Traffic Policy resource.

For information about Global Accelerator (GA) Custom Routing Endpoint Traffic Policy and how to use it, see [What is Custom Routing Endpoint Traffic Policy](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createcustomroutingendpointtrafficpolicies).

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

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
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

resource "alicloud_ga_custom_routing_endpoint" "default" {
  endpoint_group_id          = alicloud_ga_custom_routing_endpoint_group.default.id
  endpoint                   = alicloud_vswitch.default.id
  type                       = "PrivateSubNet"
  traffic_to_endpoint_policy = "AllowCustom"
}

resource "alicloud_ga_custom_routing_endpoint_group_destination" "default" {
  endpoint_group_id = alicloud_ga_custom_routing_endpoint_group.default.id
  protocols         = ["TCP"]
  from_port         = 1
  to_port           = 10
}

resource "alicloud_ga_custom_routing_endpoint_traffic_policy" "default" {
  endpoint_id = alicloud_ga_custom_routing_endpoint.default.custom_routing_endpoint_id
  address     = "172.17.3.0"
  port_ranges {
    from_port = 1
    to_port   = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Required, ForceNew) The ID of the Custom Routing Endpoint.
* `address` - (Required) The IP address of the destination to which traffic is allowed.
* `port_ranges` - (Optional, Set) Port rangeSee the following. See [`port_ranges`](#port_ranges) below.

### `port_ranges`

The port_ranges supports the following:

* `from_port` - (Optional, Int) The start port of the port range of the traffic destination. The specified port must fall within the port range of the specified endpoint group.
* `to_port` - (Optional, Int) The end port of the port range of the traffic destination. The specified port must fall within the port range of the specified endpoint group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint Traffic Policy. It formats as `<endpoint_id>:<custom_routing_endpoint_traffic_policy_id>`.
* `accelerator_id` - The ID of the GA instance.
* `listener_id` - The ID of the listener.
* `endpoint_group_id` - The ID of the endpoint group.
* `custom_routing_endpoint_traffic_policy_id` - The ID of the Custom Routing Endpoint Traffic Policy.
* `status` - The status of the Custom Routing Endpoint Traffic Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint Traffic Policy.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint Traffic Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint Traffic Policy.

## Import

Global Accelerator (GA) Custom Routing Endpoint Traffic Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint_traffic_policy.example <endpoint_id>:<custom_routing_endpoint_traffic_policy_id>
```
