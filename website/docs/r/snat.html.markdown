---
layout: "alicloud"
page_title: "Alicloud: alicloud_snat_entry"
sidebar_current: "docs-alicloud-resource-vpc"
description: |-
  Provides a Alicloud snat resource.
---

# alicloud\_snat

Provides a snat resource.

## Example Usage

Basic Usage

```
variable "name" {
  default = "snat-entry-example-name"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
  vpc_id = "${alicloud_vswitch.vswitch.vpc_id}"
  specification = "Small"
  name = "${var.name}"
}

resource "alicloud_eip" "eip" {
  name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
  allocation_id = "${alicloud_eip.eip.id}"
  instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default"{
  snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
  source_vswitch_id = "${alicloud_vswitch.vswitch.id}"
  snat_ip = "${alicloud_eip.eip.ip_address}"
}
```

## Argument Reference

The following arguments are supported:

* `snat_table_id` - (Required, ForceNew) The value can get from `alicloud_nat_gateway` Attributes "snat_table_ids".
* `source_vswitch_id` - (Required, ForceNew) The vswitch ID.
* `snat_ip` - (Required) The SNAT ip address, the ip must along bandwidth package public ip which `alicloud_nat_gateway` argument `bandwidth_packages`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the snat entry. The value formats as `<snat_table_id>:<snat_entry_id>`
* `snat_entry_id` - The id of the snat entry on the server.

## Import

Snat Entry can be imported using the id, e.g.

```
$ terraform import alicloud_snat_entry.foo stb-1aece3:snat-232ce2
```
