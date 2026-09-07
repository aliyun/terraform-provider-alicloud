---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_addresses"
sidebar_current: "docs-alicloud-datasource-vpc-ipv6-addresses"
description: |-
  Provides a list of Vpc Ipv6 Addresses to the user.
---

# alicloud\_vpc\_ipv6\_addresses

This data source provides the Vpc Ipv6 Addresses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_ipv6_addresses" "associatedInstanceId" {
  associated_instance_id = "example_value"
}
output "vpc_ipv6_address_id_1" {
  value = data.alicloud_vpc_ipv6_addresses.associatedInstanceId.addresses.0.id
}

data "alicloud_vpc_ipv6_addresses" "vswitchId" {
  vswitch_id = "example_value"
}
output "vpc_ipv6_address_id_2" {
  value = data.alicloud_vpc_ipv6_addresses.vswitchId.addresses.0.id
}

data "alicloud_vpc_ipv6_addresses" "vpcId" {
  vpc_id = "example_value"
}
output "vpc_ipv6_address_id_3" {
  value = data.alicloud_vpc_ipv6_addresses.vpcId.addresses.0.id
}

data "alicloud_vpc_ipv6_addresses" "status" {
  status = "Available"
}
output "vpc_ipv6_address_id_4" {
  value = data.alicloud_vpc_ipv6_addresses.status.addresses.0.id
}

```

## Argument Reference

The following arguments are supported:

* `associated_instance_id` - (Optional, ForceNew) The ID of the instance that is assigned the IPv6 address.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the IPv6 address. Valid values:`Pending` or `Available`. 
  - `Pending`: The IPv6 address is being configured. 
  - `Available`: The IPv6 address is available.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC to which the IPv6 address belongs.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch to which the IPv6 address belongs.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ipv6 Address names.
* `addresses` - A list of Vpc Ipv6 Addresses. Each element contains the following attributes:
	* `associated_instance_id` - The ID of the instance that is assigned the IPv6 address.
	* `associated_instance_type` - The type of the instance that is assigned the IPv6 address.
	* `create_time` - The time when the IPv6 address was created.
	* `id` - The ID of the Ipv6 Address.
	* `ipv6_address` - The address of the Ipv6 Address.
	* `ipv6_address_id` - The ID of the IPv6 address.
	* `ipv6_address_name` - The name of the IPv6 address.
	* `ipv6_gateway_id` - The ID of the IPv6 gateway to which the IPv6 address belongs.
	* `network_type` - The type of communication supported by the IPv6 address. Valid values:`Private` or `Public`. `Private`: communication within the private network. `Public`: communication over the public network
	* `status` - The status of the IPv6 address. Valid values:`Pending` or `Available`.
	* `vpc_id` - The ID of the VPC to which the IPv6 address belongs.
	* `vswitch_id` - The ID of the vSwitch to which the IPv6 address belongs.
