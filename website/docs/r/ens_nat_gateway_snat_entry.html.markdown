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

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "ens_region_id" {
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "10.0.0.0/24"
  vswitch_name  = var.name
  ens_region_id = alicloud_ens_network.default.ens_region_id
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_eip" "default" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = alicloud_ens_vswitch.default.ens_region_id
  eip_name             = var.name
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_nat_gateway" "default" {
  vswitch_id    = alicloud_ens_vswitch.default.id
  ens_region_id = alicloud_ens_vswitch.default.ens_region_id
  network_id    = alicloud_ens_vswitch.default.network_id
  instance_type = "enat.default"
  nat_name      = var.name
}

resource "alicloud_ens_eip_instance_attachment" "default" {
  instance_id   = alicloud_ens_nat_gateway.default.id
  allocation_id = alicloud_ens_eip.default.id
  instance_type = "Nat"
  standby       = false
}

resource "alicloud_ens_nat_gateway_snat_entry" "default" {
  snat_entry_name = var.name
  source_cidr     = "10.0.0.0/8"
  snat_ip         = alicloud_ens_eip.default.ip_address
  nat_gateway_id  = alicloud_ens_nat_gateway.default.id
  idle_timeout    = "50"
  isp_affinity    = false
  eip_affinity    = false
}
```

## Argument Reference

The following arguments are supported:
* `eip_affinity` - (Optional) Whether to enable IP affinity.
* `idle_timeout` - (Optional, ForceNew, Int) The idle timeout. Valid values: 1 to 86400. Unit: seconds.
* `isp_affinity` - (Optional) Whether to open the operator affinity. Value:
false: disable Operator affinity.
true: Turn on operator affinity.
* `nat_gateway_id` - (Required, ForceNew) The ID of the NAT gateway.
* `snat_entry_name` - (Optional) The name of the SNAT entry.
* `snat_ip` - (Required) The EIPs in the SNAT entry. Separate multiple EIPs with commas (,).
* `source_cidr` - (Optional, ForceNew) The source CIDR block of the SNAT entry.
* `source_network_id` - (Optional, ForceNew) The ID of the network. This parameter indicates that all ENS instances in the network can access the Internet through SNAT rules.
* `source_vswitch_id` - (Optional, ForceNew) The ID of the vSwitch that needs to access the Internet. This parameter indicates that all ENS instances in the vSwitch can access the Internet through SNAT rules.
* `standby_snat_ip` - (Optional, ForceNew) The standby EIPs in the SNAT entry. Separate multiple standby EIPs with commas (,).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `creation_time` - The creation time, in UTC.
* `dest_cidr` - The destination CIDR block.
* `standby_status` - The status of the standby EIP.
* `status` - SNAT entry status.
* `type` - The NAT type.

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