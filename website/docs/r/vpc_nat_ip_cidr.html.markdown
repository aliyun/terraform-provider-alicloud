---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_nat_ip_cidr"
sidebar_current: "docs-alicloud-resource-vpc-nat-ip-cidr"
description: |-
  Provides a Alicloud VPC Nat Ip Cidr resource.
---

# alicloud\_vpc\_nat\_ip\_cidr

Provides a VPC Nat Ip Cidr resource.

For information about VPC Nat Ip Cidr and how to use it, see [What is Nat Ip Cidr](https://www.alibabacloud.com/help/doc-detail/281972.htm).

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "example_value"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.example.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_nat_gateway" "example" {
  vpc_id               = alicloud_vpc.default.id
  internet_charge_type = "PayByLcu"
  nat_gateway_name     = "example_value"
  description          = "example_value"
  nat_type             = "Enhanced"
  vswitch_id           = alicloud_vswitch.example.id
  network_type         = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "example" {
  nat_gateway_id   = alicloud_nat_gateway.example.id
  nat_ip_cidr_name = "example_value"
  nat_ip_cidr      = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional, Computed) Specifies whether to precheck this request only. Valid values: `true` and `false`.
* `nat_gateway_id` - (Required, ForceNew) The ID of the Virtual Private Cloud (VPC) NAT gateway where you want to create the NAT CIDR block.
* `nat_ip_cidr_description` - (Optional) The description of the NAT CIDR block. The description must be `2` to `256` characters in length. It must start with a letter but cannot start with `http://` or `https://`.
* `nat_ip_cidr_name` - (Optional) The name of the NAT CIDR block. The name must be `2` to `128` characters in length and can contain digits, periods (.), underscores (_), and hyphens (-). It must start with a letter. It must start with a letter but cannot start with `http://` or `https://`.
* `nat_ip_cidr` - The NAT CIDR block to be created. The CIDR block must meet the following conditions: It must be `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, or one of their subnets. The subnet mask must be `16` to `32` bits in lengths. To use a public CIDR block as the NAT CIDR block, the VPC to which the VPC NAT gateway belongs must be authorized to use public CIDR blocks. For more information, see [Create a VPC NAT gateway](https://www.alibabacloud.com/help/doc-detail/268230.htm).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Nat Ip Cidr. The value formats as `<nat_gateway_id>:<nat_ip_cidr>`.
* `status` - The status of the CIDR block of the NAT gateway. Valid values: `Available`.

## Import

VPC Nat Ip Cidr can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_nat_ip_cidr.example <nat_gateway_id>:<nat_ip_cidr>
```
