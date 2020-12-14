---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_entry"
sidebar_current: "docs-alicloud-resource-cen-route-entry"
description: |-
  Provides a Alicloud CEN manage route entried resource.
---

# alicloud\_cen_route_entry

Provides a CEN route entry resource. Cloud Enterprise Network (CEN) supports publishing and withdrawing route entries of attached networks. You can publish a route entry of an attached VPC or VBR to a CEN instance, then other attached networks can learn the route if there is no route conflict. You can withdraw a published route entry when CEN does not need it any more.

For information about CEN route entries publishment and how to use it, see [Manage network routes](https://www.alibabacloud.com/help/doc-detail/86980.htm).

## Example Usage

Basic Usage

```
# Create a cen_route_entry resource and use it to publish a route entry pointing to an ECS.
provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-testAccCenRouteEntryConfig"
}

data "alicloud_zones" "default" {
  provider                    = alicloud.hz
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  provider          = alicloud.hz
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  provider    = alicloud.hz
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "vpc" {
  provider   = alicloud.hz
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  provider          = alicloud.hz
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_security_group" "default" {
  provider    = alicloud.hz
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.vpc.id
}

resource "alicloud_instance" "default" {
  provider                   = alicloud.hz
  vswitch_id                 = alicloud_vswitch.default.id
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.default.id]
  instance_name              = var.name
}

resource "alicloud_cen_instance" "cen" {
  name = var.name
}

resource "alicloud_cen_instance_attachment" "attach" {
  instance_id              = alicloud_cen_instance.cen.id
  child_instance_id        = alicloud_vpc.vpc.id
  child_instance_type      = "VPC"
  child_instance_region_id = "cn-hangzhou"
  depends_on               = [alicloud_vswitch.default]
}

resource "alicloud_route_entry" "route" {
  provider              = alicloud.hz
  route_table_id        = alicloud_vpc.vpc.route_table_id
  destination_cidrblock = "11.0.0.0/16"
  nexthop_type          = "Instance"
  nexthop_id            = alicloud_instance.default.id
}

resource "alicloud_cen_route_entry" "foo" {
  provider       = alicloud.hz
  instance_id    = alicloud_cen_instance.cen.id
  route_table_id = alicloud_vpc.vpc.route_table_id
  cidr_block     = alicloud_route_entry.route.destination_cidrblock
  depends_on     = [alicloud_cen_instance_attachment.attach]
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN.
* `route_table_id` - (Required, ForceNew) The route table of the attached VBR or VPC.
* `cidr_block` - (Required, ForceNew) The destination CIDR block of the route entry to publish.

->**NOTE:** The "alicloud_cen_instance_route_entries" resource depends on the related "alicloud_cen_instance_attachment" resource.

->**NOTE:** The "alicloud_cen_instance_attachment" resource should depend on the related "alicloud_vswitch" resource.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, formatted as `<instance_id>:<route_table_id>:<cidr_block>`.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_route_entry.example cen-abc123456:vtb-abc123:192.168.0.0/24
```

