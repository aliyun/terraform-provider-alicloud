---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_nat_ip_cidrs"
sidebar_current: "docs-alicloud-datasource-vpc-nat-ip-cidrs"
description: |-
  Provides a list of Vpc Nat Ip Cidrs to the user.
---

# alicloud\_vpc\_nat\_ip\_cidrs

This data source provides the Vpc Nat Ip Cidrs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_nat_ip_cidrs" "ids" {
  nat_gateway_id = "example_value"
  ids            = ["example_value-1", "example_value-2"]
}
output "vpc_nat_ip_cidr_id_1" {
  value = data.alicloud_vpc_nat_ip_cidrs.ids.cidrs.0.id
}

data "alicloud_vpc_nat_ip_cidrs" "nameRegex" {
  nat_gateway_id = "example_value"
  name_regex     = "^my-NatIpCidr"
}
output "vpc_nat_ip_cidr_id_2" {
  value = data.alicloud_vpc_nat_ip_cidrs.nameRegex.cidrs.0.id
}

data "alicloud_vpc_nat_ip_cidrs" "status" {
  nat_gateway_id = "example_value"
  ids            = ["example_value-1"]
  status         = "Available"
}
output "vpc_nat_ip_cidr_id_3" {
  value = data.alicloud_vpc_nat_ip_cidrs.status.cidrs.0.id
}

data "alicloud_vpc_nat_ip_cidrs" "natIpCidr" {
  nat_gateway_id = "example_value"
  nat_ip_cidrs   = ["example_value-1"]
}
output "vpc_nat_ip_cidr_id_4" {
  value = data.alicloud_vpc_nat_ip_cidrs.natIpCidr.cidrs.0.id
}

data "alicloud_vpc_nat_ip_cidrs" "atIpCidrName" {
  nat_gateway_id   = "example_value"
  nat_ip_cidr_name = ["example_value-1"]
}
output "vpc_nat_ip_cidr_id_5" {
  value = data.alicloud_vpc_nat_ip_cidrs.atIpCidrName.cidrs.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Nat Ip Cidr IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Nat Ip Cidr name.
* `nat_gateway_id` - (Required, ForceNew) NAT IP ADDRESS range to the security group of the Kafka VPC NAT gateway instance ID.
* `nat_ip_cidrs` - (Optional, ForceNew) The NAT CIDR block to be created. Support up to `20`. The CIDR block must meet the following conditions: It must be `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, or one of their subnets. The subnet mask must be `16` to `32` bits in lengths. To use a public CIDR block as the NAT CIDR block, the VPC to which the VPC NAT gateway belongs must be authorized to use public CIDR blocks. For more information, see [Create a VPC NAT gateway](https://www.alibabacloud.com/help/doc-detail/268230.htm).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `nat_ip_cidr_name` - (Optional, ForceNew) The name of the CIDR block that you want to query. Support up to `20`.
* `status` - (Optional, ForceNew) The status of the NAT IP address. Valid values:`Available`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Nat Ip Cidr names.
* `cidrs` - A list of Vpc Nat Ip Cidrs. Each element contains the following attributes:
	* `create_time` - The time when the CIDR block was created.
	* `id` - The ID of the Nat Ip Cidr.
	* `is_default` - Whether it is the default NAT IP ADDRESS. Valid values:`true` or `false`.`true`: is the default NAT IP ADDRESS. `false`: it is not the default NAT IP ADDRESS.
	* `nat_gateway_id` - The ID of the VPC NAT gateway.
	* `nat_ip_cidr` - The NAT CIDR block to be created. The CIDR block must meet the following conditions: It must be `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, or one of their subnets. The subnet mask must be `16` to `32` bits in lengths. To use a public CIDR block as the NAT CIDR block, the VPC to which the VPC NAT gateway belongs must be authorized to use public CIDR blocks. For more information, see [Create a VPC NAT gateway](https://www.alibabacloud.com/help/doc-detail/268230.htm).
	* `nat_ip_cidr_description` - NAT IP ADDRESS range to the description of. Length is from `2` to `256` characters, must start with a letter or the Chinese at the beginning, but not at the` http://` Or `https://` at the beginning.
	* `nat_ip_cidr_id` - NAT IP ADDRESS instance ID.
	* `nat_ip_cidr_name` - NAT IP ADDRESS the name of the root directory. Length is from `2` to `128` characters, must start with a letter or the Chinese at the beginning can contain numbers, half a period (.), underscore (_) and dash (-). But do not start with `http://` or `https://` at the beginning.
	* `status` - The status of the CIDR block of the NAT gateway. If the value is `Available`, the CIDR block is available.
