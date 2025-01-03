---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_pools"
sidebar_current: "docs-alicloud-datasource-vpc-ipam-ipam-pools"
description: |-
  Provides a list of Vpc Ipam Ipam Pool owned by an Alibaba Cloud account.
---

# alicloud_vpc_ipam_ipam_pools

This data source provides Vpc Ipam Ipam Pool available to the user.[What is Ipam Pool](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.241.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "parentIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_ipam_ipam_pool" "default" {
  ipam_scope_id         = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id        = alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id
  ipam_pool_name        = var.name
  source_ipam_pool_id   = alicloud_vpc_ipam_ipam_pool.parentIpamPool.id
  ip_version            = "IPv4"
  ipam_pool_description = var.name
}

data "alicloud_vpc_ipam_ipam_pools" "default" {
  name_regex = alicloud_vpc_ipam_ipam_pool.default.name
}

output "alicloud_vpc_ipam_ipam_pool_example_id" {
  value = data.alicloud_vpc_ipam_ipam_pools.default.pools.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ipam_pool_id` - (ForceNew, Optional) The first ID of the resource.
* `ipam_pool_name` - (ForceNew, Optional) The name of the resource.
* `ipam_scope_id` - (ForceNew, Optional) Ipam scope id.
* `pool_region_id` - (ForceNew, Optional) The effective region of the IPAM address pool.
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group.
* `source_ipam_pool_id` - (ForceNew, Optional) The instance ID of the source IPAM address pool.> If this parameter is not entered, the created address pool is the parent address pool.
* `tags` - (ForceNew, Optional) The tag of the resource.
* `ids` - (Optional, ForceNew, Computed) A list of Ipam Pool IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Ipam Pool IDs.
* `names` - A list of name of Ipam Pools.
* `pools` - A list of Ipam Pool Entries. Each element contains the following attributes:
  * `allocation_default_cidr_mask` - The default network mask assigned by the IPAM address pool.IPv4 network mask value range: **0 to 32** bits.
  * `allocation_max_cidr_mask` - The maximum network mask assigned by the IPAM address pool.IPv4 network mask value range: **0 to 32** bits.
  * `allocation_min_cidr_mask` - The minimum Network mask assigned by the IPAM address pool.IPv4 network mask value range: **0 to 32** bits.
  * `auto_import` - Whether the automatic import function is enabled for the address pool.
  * `create_time` - The creation time of the resource.
  * `has_sub_pool` - Whether it is a child address pool. Value:-**true**: Yes.-**false**: No.
  * `ip_version` - The IP protocol version. Currently, only **IPv4** is supported * *.
  * `ipam_id` - Ipam id.
  * `ipam_pool_description` - The description of the IPAM address pool.It must be 2 to 256 characters in length and must start with an English letter or a Chinese character, but cannot start with 'http:// 'or 'https. If it is not filled in, it is empty. The default value is empty.
  * `ipam_pool_id` - The first ID of the resource.
  * `ipam_pool_name` - The name of the resource.
  * `ipam_scope_id` - Ipam scope id.
  * `pool_depth` - The depth of the IPAM address pool. Value range: **0 to 10 * *.
  * `pool_region_id` - The effective region of the IPAM address pool.
  * `resource_group_id` - The ID of the resource group.
  * `source_ipam_pool_id` - The instance ID of the source IPAM address pool.> If this parameter is not entered, the created address pool is the parent address pool.
  * `status` - The status of the resource.
  * `tags` - The tag of the resource.
  * `id` - The ID of the resource supplied above.
  * `region_id` - The region ID of the resource.
