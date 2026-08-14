---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_nat_gateway_snat_entry"
description: |-
  Provides a Alicloud ENS Nat Gateway Snat Entry resource.
---

# alicloud_ens_nat_gateway_snat_entry

Provides a ENS Nat Gateway Snat Entry resource.



For information about ENS Nat Gateway Snat Entry and how to use it, see [What is Nat Gateway Snat Entry](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateSnatEntry).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

variable "ens_region_id" {
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "defaultXqhlfk" {
  network_name  = "example-snat"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultzkXvut" {
  cidr_block    = "10.0.0.0/24"
  vswitch_name  = "example-snat"
  ens_region_id = alicloud_ens_network.defaultXqhlfk.ens_region_id
  network_id    = alicloud_ens_network.defaultXqhlfk.id
}

resource "alicloud_ens_eip" "defaultiUbwh0" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = alicloud_ens_vswitch.defaultzkXvut.ens_region_id
  eip_name             = "example-snat"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_nat_gateway" "default2Kn0nu" {
  vswitch_id    = alicloud_ens_vswitch.defaultzkXvut.id
  ens_region_id = alicloud_ens_vswitch.defaultzkXvut.ens_region_id
  network_id    = alicloud_ens_vswitch.defaultzkXvut.network_id
  instance_type = "enat.default"
  nat_name      = "example-snat"
}

resource "alicloud_ens_eip_instance_attachment" "defaultlI0M0t" {
  instance_id   = alicloud_ens_nat_gateway.default2Kn0nu.id
  allocation_id = alicloud_ens_eip.defaultiUbwh0.id
  instance_type = "Nat"
  standby       = false
}

resource "alicloud_ens_eip" "eip2" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "example-snat2"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "defaultbMMEpj" {
  instance_id   = alicloud_ens_nat_gateway.default2Kn0nu.id
  allocation_id = alicloud_ens_eip.eip2.id
  instance_type = "Nat"
  standby       = false
}


resource "alicloud_ens_nat_gateway_snat_entry" "default" {
  snat_entry_name = "example-snat"
  source_cidr     = "10.0.0.0/8"
  snat_ip         = alicloud_ens_eip.defaultiUbwh0.ip_address
  nat_gateway_id  = alicloud_ens_nat_gateway.default2Kn0nu.id
  idle_timeout    = "50"
  isp_affinity    = false
  eip_affinity    = false
}
```

## Argument Reference

The following arguments are supported:
* `eip_affinity` - (Optional) Whether to enable IP affinity.
* `idle_timeout` - (Optional, ForceNew, Int) Timeout (seconds), value range: 0~86400
* `isp_affinity` - (Optional) Whether to enable the operator affinity. Valid values:
false: disable operator affinity.
true: enable operator affinity.
* `nat_gateway_id` - (Required, ForceNew) The ID of the NAT gateway.
* `snat_entry_name` - (Optional) The name of the SNAT entry.
* `snat_ip` - (Required) The EIPs in the SNAT entry. Separate multiple EIPs with commas (,).
* `source_cidr` - (Optional, ForceNew) The source CIDR block of the SNAT entry.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - SNAT entry status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Gateway Snat Entry.
* `delete` - (Defaults to 5 mins) Used when delete the Nat Gateway Snat Entry.
* `update` - (Defaults to 5 mins) Used when update the Nat Gateway Snat Entry.

## Import

ENS Nat Gateway Snat Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_nat_gateway_snat_entry.example <snat_entry_id>
```
