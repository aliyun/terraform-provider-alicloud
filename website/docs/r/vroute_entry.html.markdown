---
layout: "alicloud"
page_title: "Alicloud: alicloud_route_entry"
sidebar_current: "docs-alicloud-resource-route-entry"
description: |-
  Provides a Alicloud Route Entry resource.
---

# alicloud\_route\_entry

Provides a route entry resource. A route entry represents a route item of one VPC route table.

## Example Usage

Basic Usage

```
resource "alicloud_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_route_entry" "default" {
  route_table_id        = "${alicloud_vpc.default.router_table_id}"
  destination_cidrblock = "${var.entry_cidr}"
  nexthop_type          = "Instance"
  nexthop_id            = "${alicloud_instance.snat.id}"
}

resource "alicloud_instance" "snat" {
  // ...
}
```
## Argument Reference

The following arguments are supported:

* `router_id` - (Deprecated) This argument has beeb deprecated. Please use other arguments to launch a custom route entry.
* `route_table_id` - (Required, Forces new resource) The ID of the route table.
* `destination_cidrblock` - (Required, Forces new resource) The RouteEntry's target network segment.
* `nexthop_type` - (Required, Forces new resource) The next hop type. Available values:
    - `Instance` (Default): Route the traffic destined for the destination CIDR block to an ECS instance in the VPC.
    - `RouterInterface`: Route the traffic destined for the destination CIDR block to a router interface.
    - `VpnGateway`: Route the traffic destined for the destination CIDR block to a VPN Gateway.
    - `HaVip`: Route the traffic destined for the destination CIDR block to an HAVIP.

* `nexthop_id` - (Required, Forces new resource) The route entry's next hop. ECS instance ID or VPC router interface ID.

## Attributes Reference

The following attributes are exported:

* `router_id` - The ID of the virtual router attached to Vpc.
* `route_table_id` - The ID of the route table.
* `destination_cidrblock` - The RouteEntry's target network segment.
* `nexthop_type` - The next hop type.
* `nexthop_id` - The route entry's next hop.

## Import

Router entry can be imported using the id, e.g.

```
$ terraform import alicloud_route_entry.example abc123456
```

