---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pools"
sidebar_current: "docs-alicloud-datasource-vpc-public-ip-address-pools"
description: |-
  Provides a list of Vpc Public Ip Address Pools to the user.
---

# alicloud\_vpc\_public\_ip\_address\_pools

This data source provides the Vpc Public Ip Address Pools of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.186.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_public_ip_address_pools" "ids" {
  ids = ["example_id"]
}

output "vpc_public_ip_address_pool_id_1" {
  value = data.alicloud_vpc_public_ip_address_pools.ids.pools.0.id
}

data "alicloud_vpc_public_ip_address_pools" "nameRegex" {
  name_regex = "example_name"
}

output "vpc_public_ip_address_pool_id_2" {
  value = data.alicloud_vpc_public_ip_address_pools.nameRegex.pools.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Vpc Public Ip Address Pool IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Vpc Public Ip Address Pool name.
* `public_ip_address_pool_ids` - (Optional, ForceNew) The IDs of the Vpc Public IP address pools.
* `public_ip_address_pool_name` - (Optional, ForceNew) The name of the VPC Public IP address pool.
* `isp` - (Optional, ForceNew) The Internet service provider. Valid values: `BGP`, `BGP_PRO`, `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2`, `BGP_FinanceCloud`.
* `status` - (Optional, ForceNew) The status of the Vpc Public Ip Address Pool. Valid values: `Created`, `Deleting`, `Modifying`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Public Ip Address Pool names.
* `pools` - A list of Vpc Public Ip Address Pools. Each element contains the following attributes:
	* `id` - The ID of the Vpc Public Ip Address Pool.
	* `public_ip_address_pool_id` - The ID of the Vpc Public Ip Address Pool.
	* `public_ip_address_pool_name` - The name of the Vpc Public Ip Address Pool.
	* `isp` - The Internet service provider.
	* `description` - The description of the Vpc Public Ip Address Pool.
	* `status` - The status of the Vpc Public Ip Address Pool.
	* `region_id` - The region ID of the Vpc Public Ip Address Pool.
	* `user_type` - The user type.
	* `total_ip_num` - The total number of IP addresses in the Vpc Public Ip Address Pool.
	* `used_ip_num` - The number of occupied IP addresses in the Vpc Public Ip Address Pool.
	* `create_time` - The time when the Vpc Public Ip Address Pool was created. The time is displayed in YYYY-MM-DDThh:mm:ssZ format.
	* `ip_address_remaining` - Indicates whether the Vpc Public Ip Address Pool has idle IP addresses.