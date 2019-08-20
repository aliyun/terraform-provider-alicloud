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
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_14.*_64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "RouteEntryConfig"
}
resource "alicloud_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = "${alicloud_vpc.foo.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
  name        = "${var.name}"
  description = "foo"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

  vswitch_id = "${alicloud_vswitch.foo.id}"

  instance_charge_type       = "PostPaid"
  instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5

  system_disk_category = "cloud_efficiency"
  image_id             = "${data.alicloud_images.default.images.0.id}"
  instance_name        = "${var.name}"
}
resource "alicloud_route_entry" "foo" {
  route_table_id        = "${alicloud_vpc.foo.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = "${alicloud_instance.foo.id}"
}
```
## Argument Reference

The following arguments are supported:

* `router_id` - (Deprecated) This argument has beeb deprecated. Please use other arguments to launch a custom route entry.
* `route_table_id` - (Required, ForceNew) The ID of the route table.
* `destination_cidrblock` - (ForceNew) The RouteEntry's target network segment.
* `nexthop_type` - (ForceNew) The next hop type. Available values:
    - `Instance` (Default): Route the traffic destined for the destination CIDR block to an ECS instance in the VPC.
    - `RouterInterface`: Route the traffic destined for the destination CIDR block to a router interface.
    - `VpnGateway`: Route the traffic destined for the destination CIDR block to a VPN Gateway.
    - `HaVip`: Route the traffic destined for the destination CIDR block to an HAVIP.
    - `NetworkInterface`: Route the traffic destined for the destination CIDR block to an NetworkInterface.
    - `NatGateway`: Route the traffic destined for the destination CIDR block to an Nat Gateway.

* `nexthop_id` - (ForceNew) The route entry's next hop. ECS instance ID or VPC router interface ID.
* `name` - (Optional, ForceNew, Available in 1.54.1+) The name of the route entry. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://.

## Attributes Reference

The following attributes are exported:

* `id` - The route entry id,it formats of `<route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id>`.
* `router_id` - The ID of the virtual router attached to Vpc.
* `route_table_id` - The ID of the route table.
* `destination_cidrblock` - The RouteEntry's target network segment.
* `nexthop_type` - The next hop type.
* `nexthop_id` - The route entry's next hop.

## Import

Router entry can be imported using the id, e.g (formatted as<route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id>).

```
$ terraform import alicloud_route_entry.example vtb-123456:vrt-123456:0.0.0.0/0:NatGateway:ngw-123456
```

