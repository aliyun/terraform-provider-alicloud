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

-> **NOTE:** Available since v1.287.0.

## Example Usage

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
}

data "alicloud_ens_nat_gateway_snat_entries" "default" {
  nat_gateway_id = alicloud_ens_nat_gateway.default.id
  ids            = [alicloud_ens_nat_gateway_snat_entry.default.id]
}

output "snat_entry_status" {
  value = data.alicloud_ens_nat_gateway_snat_entries.default.entries.0.status
}
```

## Argument Reference

The following arguments are supported:
* `nat_gateway_id` - (Required) The ID of the NAT gateway.
* `snat_entry_id` - (Optional) SNAT entry ID.
* `snat_entry_name` - (Optional) The name of the SNAT entry.
* `snat_ip` - (Optional) The EIPs in the SNAT entry. Separate multiple EIPs with commas (,).
* `source_cidr` - (Optional) The source CIDR block of the SNAT entry.
* `ids` - (Optional, Computed) A list of Nat Gateway Snat Entry IDs.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Nat Gateway Snat Entry IDs.
* `entries` - A list of Nat Gateway Snat Entry Entries. Each element contains the following attributes:
  * `creation_time` - **NOTE:** This field is only available when `enable_details` is `true`. The creation time, in UTC.
  * `dest_cidr` - **NOTE:** This field is only available when `enable_details` is `true`. The destination CIDR block.
  * `eip_affinity` - Whether to enable IP affinity.
  * `idle_timeout` - The idle timeout.
  * `isp_affinity` - Whether to open the operator affinity.
  * `nat_gateway_id` - The ID of the NAT gateway.
  * `snat_entry_id` - SNAT entry ID.
  * `snat_entry_name` - The name of the SNAT entry.
  * `snat_ip` - The EIPs in the SNAT entry.
  * `source_cidr` - The source CIDR block of the SNAT entry.
  * `standby_snat_ip` - The standby EIPs in the SNAT entry.
  * `standby_status` - The status of the standby EIP.
  * `status` - SNAT entry status.
  * `type` - **NOTE:** This field is only available when `enable_details` is `true`. The NAT type.
  * `id` - The ID of the resource supplied above.
