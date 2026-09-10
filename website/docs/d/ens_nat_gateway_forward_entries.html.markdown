---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_nat_gateway_forward_entries"
sidebar_current: "docs-alicloud-datasource-ens-nat-gateway-forward-entries"
description: |-
  Provides a list of ENS Nat Gateway Forward Entry owned by an Alibaba Cloud account.
---

# alicloud_ens_nat_gateway_forward_entries

This data source provides ENS Nat Gateway Forward Entry available to the user.[What is Nat Gateway Forward Entry](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateForwardEntry)

-> **NOTE:** Available since v1.287.0.

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

resource "alicloud_ens_network" "default6T9qR2" {
  network_name  = "测试用例_Dnat"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5BAAN2" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = "测试用例-dnat"
  ens_region_id = alicloud_ens_network.default6T9qR2.ens_region_id
  network_id    = alicloud_ens_network.default6T9qR2.id
}

resource "alicloud_ens_nat_gateway" "defaultlZ7YKl" {
  vswitch_id    = alicloud_ens_vswitch.default5BAAN2.id
  ens_region_id = alicloud_ens_vswitch.default5BAAN2.ens_region_id
  network_id    = alicloud_ens_vswitch.default5BAAN2.network_id
  instance_type = "enat.default"
  nat_name      = "测试用例-dnat"
}

resource "alicloud_ens_eip" "defaultLQgQB6" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "测试用例-dnat"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "defaultc19VZl" {
  instance_id   = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  allocation_id = alicloud_ens_eip.defaultLQgQB6.id
  instance_type = "Nat"
}

resource "alicloud_ens_instance" "defaulth6OQ3p" {
  auto_renew = false
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.default5BAAN2.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "测试用例-dnat"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}

resource "alicloud_ens_eip" "eip2" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "测试用例-dnat2"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "default4Ph8bE" {
  instance_id   = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  allocation_id = alicloud_ens_eip.eip2.id
  instance_type = "Nat"
}

resource "alicloud_ens_instance" "instance2" {
  auto_renew = false
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.default5BAAN2.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "测试用例-dnat2"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}


resource "alicloud_ens_nat_gateway_forward_entry" "default" {
  external_port      = "100/200"
  external_ip        = alicloud_ens_eip.defaultLQgQB6.ip_address
  ip_protocol        = "TCP"
  internal_port      = "100/200"
  health_check_port  = "150"
  nat_gateway_id     = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  forward_entry_name = "测试用例-dnat"
  internal_ip        = alicloud_ens_instance.defaulth6OQ3p.private_ip_address
}

data "alicloud_ens_nat_gateway_forward_entries" "default" {
  ids                = ["${alicloud_ens_nat_gateway_forward_entry.default.id}"]
  external_ip        = alicloud_ens_eip.defaultLQgQB6.ip_address
  forward_entry_name = "测试用例-dnat"
  internal_ip        = alicloud_ens_instance.defaulth6OQ3p.private_ip_address
  ip_protocol        = "TCP"
  nat_gateway_id     = alicloud_ens_nat_gateway.defaultlZ7YKl.id
}

output "alicloud_ens_nat_gateway_forward_entry_example_id" {
  value = data.alicloud_ens_nat_gateway_forward_entries.default.entries.0.id
}
```

## Argument Reference

The following arguments are supported:
* `external_ip` - (Optional) The elastic public IP address that provides public network access in the DNAT entry.
* `forward_entry_id` - (Optional) The ID of the DNAT entry.
* `forward_entry_name` - (Optional) The name of the DNAT entry.
* `internal_ip` - (Optional) The private IP address of the instance that uses the DNAT entry for public network communication.
* `ip_protocol` - (Optional) Protocol type, value:
  - `TCP`: forwards TCP packets.
  - `UDP`: forwards UDP packets.
  - `Any`: Forward messages of all protocols.
* `nat_gateway_id` - (Required) The ID of the NAT gateway.
* `ids` - (Optional, Computed) A list of Nat Gateway Forward Entry IDs. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Nat Gateway Forward Entry IDs.
* `entries` - A list of Nat Gateway Forward Entry Entries. Each element contains the following attributes:
    * `external_ip` - The elastic public IP address that provides public network access in the DNAT entry.
    * `external_port` - The external port or port segment for port forwarding.
    * `forward_entry_id` - The ID of the DNAT entry.
    * `forward_entry_name` - The name of the DNAT entry.
    * `health_check_port` - The detection port of DNAT must be within the intranet port range.
    * `internal_ip` - The private IP address of the instance that uses the DNAT entry for public network communication.
    * `internal_port` - Internal port or port segment for port forwarding.
    * `ip_protocol` - Protocol type, value:.
    * `nat_gateway_id` - The ID of the NAT gateway.
    * `status` - DNAT entry status, value:.
    * `id` - The ID of the resource supplied above.
