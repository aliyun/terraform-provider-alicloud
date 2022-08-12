---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv4_gateways"
sidebar_current: "docs-alicloud-datasource-vpc-ipv4-gateways"
description: |-
  Provides a list of Vpc Ipv4 Gateways to the user.
---

# alicloud\_vpc\_ipv4\_gateways

This data source provides the Vpc Ipv4 Gateways of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_ipv4_gateways" "ids" {}
output "vpc_ipv4_gateway_id_1" {
  value = data.alicloud_vpc_ipv4_gateways.ids.gateways.0.id
}

data "alicloud_vpc_ipv4_gateways" "nameRegex" {
  name_regex = "^my-Ipv4Gateway"
}
output "vpc_ipv4_gateway_id_2" {
  value = data.alicloud_vpc_ipv4_gateways.nameRegex.gateways.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ipv4 Gateway IDs.
* `ipv4_gateway_name` - (Optional, ForceNew) The name of the IPv4 gateway.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ipv4 Gateway name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Creating`, `Created`, `Deleting`, `Pending`, `Deleted`.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC associated with the IPv4 Gateway.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ipv4 Gateway names.
* `gateways` - A list of Vpc Ipv4 Gateways. Each element contains the following attributes:
	* `create_time` - The creation time of the resource.
	* `enabled` - Indicates whether the IPv4 gateway is activated.
	* `id` - The ID of the Ipv4 Gateway.
	* `ipv4_gateway_description` - The description of the IPv4 gateway.
	* `ipv4_gateway_id` - The resource attribute field that represents the resource level 1 ID.
	* `ipv4_gateway_name` - The name of the IPv4 gateway.
	* `ipv4_gateway_route_table_id` - ID of the route table associated with IPv4 Gateway.
	* `status` - The status of the resource.
	* `vpc_id` - The ID of the VPC associated with the IPv4 Gateway.