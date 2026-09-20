---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_nat_gateway_snat_entries"
sidebar_current: "docs-alicloud-datasource-ens-nat-gateway-snat-entries"
description: |-
  Provides a list of ENS Nat Gateway Snat Entry owned by an Alibaba Cloud account.
---

# alicloud_ens_nat_gateway_snat_entries

This data source provides ENS Nat Gateway Snat Entry available to the user.[What is Nat Gateway Snat Entry](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateSnatEntry)

-> **NOTE:** Available since v1.289.0.

## Example Usage

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

data "alicloud_ens_nat_gateway_snat_entries" "default" {
  ids             = ["${alicloud_ens_nat_gateway_snat_entry.default.id}"]
  snat_ip         = alicloud_ens_eip.defaultiUbwh0.ip_address
  snat_entry_name = "example-snat"
  source_cidr     = "10.0.0.0/8"
  nat_gateway_id  = alicloud_ens_nat_gateway.default2Kn0nu.id
}

output "alicloud_ens_nat_gateway_snat_entry_example_id" {
  value = data.alicloud_ens_nat_gateway_snat_entries.default.entries.0.id
}
```

## Argument Reference

The following arguments are supported:
* `snat_entry_id` - (Optional) The ID of the SNAT entry.
* `snat_entry_name` - (Optional) The name of the SNAT entry.
* `snat_ip` - (Optional) The EIPs in the SNAT entry. Separate multiple EIPs with commas (,).
* `source_cidr` - (Optional) The source CIDR block of the SNAT entry.
* `nat_gateway_id` - (Required) The ID of the NAT gateway.
* `ids` - (Optional, Computed) A list of Nat Gateway Snat Entry IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Nat Gateway Snat Entry IDs.
* `entries` - A list of Nat Gateway Snat Entry Entries. Each element contains the following attributes:
  * `eip_affinity` - Whether to enable IP affinity.
  * `idle_timeout` - The idle timeout. Valid values: 1 to 86400. Unit: seconds.
  * `isp_affinity` - Whether to open the operator affinity.
  * `nat_gateway_id` - The ID of the NAT gateway.
  * `snat_entry_id` - The ID of the SNAT entry.
  * `snat_entry_name` - The name of the SNAT entry.
  * `snat_ip` - The EIPs in the SNAT entry.
  * `source_cidr` - The source CIDR block of the SNAT entry.
  * `standby_snat_ip` - The standby EIPs in the SNAT entry.
  * `standby_status` - The status of the standby EIP.
  * `status` - SNAT entry status.
  * `id` - The ID of the resource supplied above.
