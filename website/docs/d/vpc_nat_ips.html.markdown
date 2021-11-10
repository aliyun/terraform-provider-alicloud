---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_nat_ips"
sidebar_current: "docs-alicloud-datasource-vpc-nat-ips"
description: |-
  Provides a list of Vpc Nat Ips to the user.
---

# alicloud\_vpc\_nat\_ips

This data source provides the Vpc Nat Ips of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_nat_ips" "ids" {
  nat_gateway_id = "example_value"
  ids            = ["example_value-1", "example_value-2"]
}
output "vpc_nat_ip_id_1" {
  value = data.alicloud_vpc_nat_ips.ids.ips.0.id
}

data "alicloud_vpc_nat_ips" "nameRegex" {
  nat_gateway_id = "example_value"
  name_regex     = "^my-NatIp"
}
output "vpc_nat_ip_id_2" {
  value = data.alicloud_vpc_nat_ips.nameRegex.ips.0.id
}

data "alicloud_vpc_nat_ips" "natIpCidr" {
  nat_gateway_id = "example_value"
  nat_ip_cidr    = "example_value"
  name_regex     = "^my-NatIp"
}
output "vpc_nat_ip_id_3" {
  value = data.alicloud_vpc_nat_ips.natIpCidr.ips.0.id
}

data "alicloud_vpc_nat_ips" "natIpName" {
  nat_gateway_id = "example_value"
  ids            = ["example_value"]
  nat_ip_name    = ["example_value"]
}
output "vpc_nat_ip_id_4" {
  value = data.alicloud_vpc_nat_ips.natIpName.ips.0.id
}

data "alicloud_vpc_nat_ips" "natIpIds" {
  nat_gateway_id = "example_value"
  ids            = ["example_value"]
  nat_ip_ids     = ["example_value"]
}
output "vpc_nat_ip_id_5" {
  value = data.alicloud_vpc_nat_ips.natIpIds.ips.0.id
}

data "alicloud_vpc_nat_ips" "status" {
  nat_gateway_id = "example_value"
  ids            = ["example_value"]
  status         = "example_value"
}
output "vpc_nat_ip_id_6" {
  value = data.alicloud_vpc_nat_ips.status.ips.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Nat Ip IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Nat Ip name.
* `nat_gateway_id` - (Required, ForceNew) The ID of the Virtual Private Cloud (VPC) NAT gateway for which you want to create the NAT IP address.
* `nat_ip_cidr` - (Required, ForceNew) NAT IP ADDRESS of the address segment.
* `nat_ip_name` - (Optional, ForceNew) NAT IP ADDRESS the name of the root directory. Length is from `2` to `128` characters, must start with a letter or the Chinese at the beginning can contain numbers, half a period (.), underscore (_) and dash (-). But do not start with `http://` or `https://` at the beginning.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the NAT IP address. Valid values: `Available`, `Deleting` and `Creating`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Nat Ip IDs.
* `names` - A list of Nat Ip names.
* `ips` - A list of Vpc Nat Ips. Each element contains the following attributes:
	* `id` - The ID of the Nat Ip.
	* `is_default` - Indicates whether the BGP Group is the default NAT IP ADDRESS. Valid values: `true`: is the default NAT IP ADDRESS. `false`: it is not the default NAT IP ADDRESS.
	* `nat_gateway_id` - The ID of the Virtual Private Cloud (VPC) NAT gateway to which the NAT IP address belongs.
	* `nat_ip` - The NAT IP address that is queried.
	* `nat_ip_cidr` - The CIDR block to which the NAT IP address belongs.
	* `nat_ip_description` - The description of the NAT IP address.
	* `nat_ip_id` - The ID of the NAT IP address.
	* `nat_ip_name` - The name of the NAT IP address.
	* `status` - The status of the NAT IP address. Valid values: `Available`, `Deleting` and `Creating`.
