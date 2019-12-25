---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_route_entries"
sidebar_current: "docs-alicloud-datasource-route-entries"
description: |-
    Provides a list of Route Entries owned by an Alibaba Cloud account.
---

# alicloud\_route\_entries

This data source provides a list of Route Entries owned by an Alibaba Cloud account.

-> **NOTE:** Available in 1.37.0+.

## Example Usage

```
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*_64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccRouteEntryConfig"
}
resource "alicloud_vpc" "foo" {
  name       = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = alicloud_vpc.foo.id
  cidr_block        = "10.1.1.0/24"
  availability_zone = data.alicloud_zones.default.zones.0.id
  name              = var.name
}

resource "alicloud_route_entry" "foo" {
  route_table_id        = alicloud_vpc.foo.route_table_id
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = alicloud_instance.foo.id
}

resource "alicloud_security_group" "tf_test_foo" {
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_security_group_rule" "ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.tf_test_foo.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
  # cn-beijing
  security_groups = [alicloud_security_group.tf_test_foo.id]

  vswitch_id         = alicloud_vswitch.foo.id
  allocate_public_ip = true

  # series III
  instance_charge_type       = "PostPaid"
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5

  system_disk_category = "cloud_efficiency"
  image_id             = data.alicloud_images.default.images.0.id
  instance_name        = var.name
}

data "alicloud_route_entries" "foo" {
  route_table_id = alicloud_route_entry.foo.route_table_id
}

```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, ForceNew) The ID of the router table to which the route entry belongs.
* `instance_id` - (Optional) The instance ID of the next hop.
* `type` - (Optional) The type of the route entry.
* `cidr_block` - (Optional) The destination CIDR block of the route entry.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `entries` - A list of Route Entries. Each element contains the following attributes:
  * `type` - The type of the route entry.
  * `next_hop_type` - The type of the next hop.
  * `status` - The status of the route entry.
  * `instance_id` - The instance ID of the next hop.
  * `route_table_id` - The ID of the router table to which the route entry belongs.
  * `cidr_block` - The destination CIDR block of the route entry.

