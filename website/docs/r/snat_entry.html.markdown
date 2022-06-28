---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_snat_entry"
sidebar_current: "docs-alicloud-resource-vpc"
description: |-
  Provides a Alicloud snat resource.
---

# alicloud\_snat\_entry

Provides a snat resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "snat-entry-example-name"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.vpc.id
  cidr_block = "172.16.0.0/21"
  zone_id    = data.alicloud_zones.default.zones[0].id
  name       = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id        = alicloud_vswitch.vswitch.vpc_id
  specification = "Small"
  name          = var.name
}

resource "alicloud_eip_address" "default" {
  count        = 2
  address_name = var.name
}

resource "alicloud_eip_association" "default" {
  count         = 2
  allocation_id = element(alicloud_eip_address.default.*.id, count.index)
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_common_bandwidth_package" "default" {
  name                 = "tf_cbp"
  bandwidth            = 10
  internet_charge_type = "PayByTraffic"
  ratio                = 100
}

resource "alicloud_common_bandwidth_package_attachment" "default" {
  count                = 2
  bandwidth_package_id = alicloud_common_bandwidth_package.default.id
  instance_id          = element(alicloud_eip_address.default.*.id, count.index)
}

resource "alicloud_snat_entry" "default" {
  depends_on        = [alicloud_eip_association.default]
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vswitch.id
  snat_ip           = join(",", alicloud_eip_address.default.*.ip_address)
}
```

## Argument Reference

The following arguments are supported:

* `snat_table_id` - (Required, ForceNew) The value can get from `alicloud_nat_gateway` Attributes "snat_table_ids".
* `source_vswitch_id` - (Optional, ForceNew, Computed) The vswitch ID.
* `source_cidr` - (Optional, ForceNew, Available in 1.71.1+, Computed) The private network segment of Ecs. This parameter and the `source_vswitch_id` parameter are mutually exclusive and cannot appear at the same time.
* `snat_entry_name` - (Optional, Available in 1.71.2+) The name of snat entry.
* `snat_ip` - (Required, ForceNew) The SNAT ip address, the ip must along bandwidth package public ip which `alicloud_nat_gateway` argument `bandwidth_packages`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the snat entry. The value formats as `<snat_table_id>:<snat_entry_id>`
* `snat_entry_id` - The id of the snat entry on the server.
* `status` - (Available in 1.119.1+) The status of snat entry.

### Timeouts

-> **NOTE:** Available in 1.119.1+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the snat.
* `update` - (Defaults to 2 mins) Used when update the snat.
* `delete` - (Defaults to 2 mins) Used when delete the snat.

## Import

Snat Entry can be imported using the id, e.g.

```
$ terraform import alicloud_snat_entry.foo stb-1aece3:snat-232ce2
```
