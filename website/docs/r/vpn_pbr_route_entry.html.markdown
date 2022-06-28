---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_pbr_route_entry"
sidebar_current: "docs-alicloud-resource-vpn-pbr-route-entry"
description: |-
  Provides a Alicloud VPN Pbr Route Entry resource.
---

# alicloud\_vpn\_pbr\_route\_entry

Provides a VPN Pbr Route Entry resource.

-> **NOTE:** Available in 1.162.0+.

For information about VPN Pbr Route Entry and how to use it, see [What is VPN Pbr Route Entry](https://www.alibabacloud.com/help/en/doc-detail/127248.html).


## Example Usage

Basic Usage

```
variable "name" {
  default = "tfacc"
}

data "alicloud_vpn_gateways" "default" {
}

resource "alicloud_vpn_customer_gateway" "default" {
  name       = var.name
  ip_address = "192.168.1.1"
}

resource "alicloud_vpn_connection" "default" {
  name                = var.name
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  vpn_gateway_id      = data.alicloud_vpn_gateways.default.ids.0
  local_subnet        = ["192.168.2.0/24"]
  remote_subnet       = ["192.168.3.0/24"]
}

resource alicloud_vpn_pbr_route_entry default {
  vpn_gateway_id = data.alicloud_vpn_gateways.default.ids.0
  route_source   = "192.168.1.0/24"
  route_dest     = "10.0.0.0/24"
  next_hop       = alicloud_vpn_connection.default.id
  weight         = 0
  publish_vpc    = false
}
```
## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Required, ForceNew) The ID of the vpn gateway.
* `next_hop` - (Required, ForceNew) The next hop of the policy-based route.
* `publish_vpc` - (Required) Whether to issue the destination route to the VPC.
* `route_source` - (Required, ForceNew) The source CIDR block of the policy-based route.
* `route_dest` - (Required, ForceNew) The destination CIDR block of the policy-based route.
* `weight` - (Required) The weight of the policy-based route. Valid values: 0 and 100.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the vpn pbr route entry. The value formats as `<vpn_gateway_id>:<next_hop>:<route_source>:<route_dest>`.
* `status` - The status of the vpn pbr route entry.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the vpn pbr route entry.
* `update` - (Defaults to 5 mins) Used when update the vpn pbr route entry.
* `delete` - (Defaults to 5 mins) Used when delete the vpn pbr route entry.

## Import

VPN Pbr route entry can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_pbr_route_entry.example <vpn_gateway_id>:<next_hop>:<route_source>:<route_dest>
```
