---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_gateways"
sidebar_current: "docs-alicloud-datasource-vpc-ipv6-gateways"
description: |-
  Provides a list of Vpc Ipv6 Gateways to the user.
---

# alicloud\_vpc\_ipv6\_gateways

This data source provides the Vpc Ipv6 Gateways of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_ipv6_gateways" "ids" {
  ids = ["example_id"]
}
output "vpc_ipv6_gateway_id_1" {
  value = data.alicloud_vpc_ipv6_gateways.ids.gateways.0.id
}

data "alicloud_vpc_ipv6_gateways" "nameRegex" {
  name_regex = "^my-Ipv6Gateway"
}
output "vpc_ipv6_gateway_id_2" {
  value = data.alicloud_vpc_ipv6_gateways.nameRegex.gateways.0.id
}

data "alicloud_vpc_ipv6_gateways" "vpcId" {
  ids    = ["example_id"]
  vpc_id = "example_value"
}
output "vpc_ipv6_gateway_id_3" {
  value = data.alicloud_vpc_ipv6_gateways.vpcId.gateways.0.id
}

data "alicloud_vpc_ipv6_gateways" "status" {
  ids    = ["example_id"]
  status = "Available"
}
output "vpc_ipv6_gateway_id_4" {
  value = data.alicloud_vpc_ipv6_gateways.status.gateways.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ipv6 Gateway IDs.
* `ipv6_gateway_name` - (Optional, ForceNew) The name of the IPv6 gateway.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ipv6 Gateway name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Deleting`, `Pending`.
* `vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC) to which the IPv6 gateway belongs.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ipv6 Gateway names.
* `gateways` - A list of Vpc Ipv6 Gateways. Each element contains the following attributes:
	* `business_status` - The status of the IPv6 gateway. Valid values:`Normal`, `FinancialLocked` and `SecurityLocked`. `Normal`: working as expected. `FinancialLocked`: locked due to overdue payments. `SecurityLocked`: locked due to security reasons.
	* `create_time` - The creation time of the resource.
	* `description` - The description of the IPv6 gateway.
	* `expired_time` - The time when the IPv6 gateway expires.
	* `id` - The ID of the Ipv6 Gateway.
	* `instance_charge_type` - The metering method of the IPv6 gateway. Valid values: `PayAsYouGo`.
	* `ipv6_gateway_id` - The first ID of the resource.
	* `ipv6_gateway_name` - The name of the IPv6 gateway.
	* `spec` - The specification of the IPv6 gateway. Valid values: `Large`, `Medium` and `Small`. `Small` (default): Free Edition. `Medium`: Enterprise Edition . `Large`: Enhanced Enterprise Edition. The throughput capacity of an IPv6 gateway varies based on the edition. For more information, see [Editions of IPv6 gateways](https://www.alibabacloud.com/help/doc-detail/98926.htm).
	* `status` - The status of the IPv6 gateway. Valid values: `Available`, `Deleting`, `Pending`.
	* `vpc_id` - The ID of the virtual private cloud (VPC) to which the IPv6 gateway belongs.
