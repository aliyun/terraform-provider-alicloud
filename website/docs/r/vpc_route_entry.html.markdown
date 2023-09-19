---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_route_entry"
description: |-
  Provides a Alicloud VPC Route Entry resource.
---

# alicloud_vpc_route_entry

Provides a VPC Route Entry resource. There are route entries in the routing table, and the next hop is judged based on the route entries.

For information about VPC Route Entry and how to use it, see [What is Route Entry](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.211.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "10.1.1.0/24"
  availability_zone = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "default"
  vpc_id      = alicloud_vpc.default.id
}

resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_instance" "default" {
  security_groups            = ["${alicloud_security_group.default.id}"]
  vswitch_id                 = alicloud_vswitch.default.id
  allocate_public_ip         = true
  instance_charge_type       = "PostPaid"
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  system_disk_category       = "cloud_efficiency"
  image_id                   = data.alicloud_images.default.images.0.id
  instance_name              = var.name
}

resource "alicloud_vpc_route_entry" "default" {
  route_table_id        = alicloud_vpc.default.route_table_id
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = alicloud_instance.default.id
  name                  = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description of the route entry.
* `destination_cidrblock` - (Optional, ForceNew) The destination network segment of the routing entry.
* `nexthop_id` - (Optional) The ID of the next hop instance of the custom route entry.
* `nexthop_type` - (Optional) The type of the next hop of the custom route entry. Valid values:
  - **Instance** (default): The ECS Instance.
  - **HaVip**: a highly available virtual IP address.
  - **RouterInterface**: indicates the router interface.
  - **Network interface**: ENI.
  - **VpnGateway**: the VPN gateway.
  - **IPv6Gateway**:IPv6 gateway.
  - **NatGateway**:NAT gateway.
  - **Attachment**: The forwarding router.
  - **VpcPeer**:VPC peer connection.
* `route_entry_name` - (Optional) The name of the route entry.
* `route_table_id` - (Required, ForceNew) Routing table ID.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.211.0). Field 'name' has been deprecated from provider version 1.211.0. New field 'route_entry_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<route_table_id>:<route_entry_id>`.
* `route_entry_id` - The ID of the route entry.
* `status` - The status of the route entry.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Route Entry.
* `delete` - (Defaults to 5 mins) Used when delete the Route Entry.
* `update` - (Defaults to 5 mins) Used when update the Route Entry.

## Import

VPC Route Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_route_entry.example <route_table_id>:<route_entry_id>
```