---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_nat_gateway_forward_entry"
description: |-
  Provides a Alicloud ENS Nat Gateway Forward Entry resource.
---

# alicloud_ens_nat_gateway_forward_entry

Provides a ENS Nat Gateway Forward Entry resource.



For information about ENS Nat Gateway Forward Entry and how to use it, see [What is Nat Gateway Forward Entry](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateForwardEntry).

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

resource "alicloud_ens_network" "default6T9qR2" {
  network_name  = "example用例_Dnat"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5BAAN2" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = "example用例-dnat"
  ens_region_id = alicloud_ens_network.default6T9qR2.ens_region_id
  network_id    = alicloud_ens_network.default6T9qR2.id
}

resource "alicloud_ens_nat_gateway" "defaultlZ7YKl" {
  vswitch_id    = alicloud_ens_vswitch.default5BAAN2.id
  ens_region_id = alicloud_ens_vswitch.default5BAAN2.ens_region_id
  network_id    = alicloud_ens_vswitch.default5BAAN2.network_id
  instance_type = "enat.default"
  nat_name      = "example用例-dnat"
}

resource "alicloud_ens_eip" "defaultLQgQB6" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "example用例-dnat"
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
  instance_name              = "example用例-dnat"
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
  eip_name             = "example用例-dnat2"
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
  instance_name              = "example用例-dnat2"
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
  forward_entry_name = "example用例-dnat"
  internal_ip        = alicloud_ens_instance.defaulth6OQ3p.private_ip_address
}
```

## Argument Reference

The following arguments are supported:
* `external_ip` - (Required) The elastic public IP address that provides public network access in the DNAT entry.
* `external_port` - (Required) The external port or port segment for port forwarding.
* `forward_entry_name` - (Optional) The name of the DNAT entry.
* `health_check_port` - (Optional, Int) The detection port of DNAT must be within the intranet port range. The default value is empty.
* `internal_ip` - (Required) The private IP address of the instance that uses the DNAT entry for public network communication.
* `internal_port` - (Required) Internal port or port segment for port forwarding.
* `ip_protocol` - (Optional) Protocol type, value:
  - `TCP`: forwards TCP packets.
  - `UDP`: forwards UDP packets.
  - `Any`: Forward messages of all protocols.
* `nat_gateway_id` - (Required, ForceNew) The ID of the NAT gateway.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `status` - DNAT entry status, value:.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Gateway Forward Entry.
* `delete` - (Defaults to 5 mins) Used when delete the Nat Gateway Forward Entry.
* `update` - (Defaults to 5 mins) Used when update the Nat Gateway Forward Entry.

## Import

ENS Nat Gateway Forward Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_nat_gateway_forward_entry.example <forward_entry_id>
```