---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_pool_allocations"
sidebar_current: "docs-alicloud-datasource-vpc-ipam-ipam-pool-allocations"
description: |-
  Provides a list of Vpc Ipam Ipam Pool Allocation owned by an Alibaba Cloud account.
---

# alicloud_vpc_ipam_ipam_pool_allocations

This data source provides Vpc Ipam Ipam Pool Allocation available to the user.[What is Ipam Pool Allocation](https://www.alibabacloud.com/help/en/)

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

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpamPoolCidr" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}

resource "alicloud_vpc_ipam_ipam_pool_allocation" "default" {
  ipam_pool_allocation_description = "init alloc desc"
  ipam_pool_allocation_name        = var.name
  cidr                             = "10.0.0.0/20"
  ipam_pool_id                     = alicloud_vpc_ipam_ipam_pool_cidr.defaultIpamPoolCidr.ipam_pool_id
}

data "alicloud_vpc_ipam_ipam_pool_allocations" "default" {
  ids = ["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}"]
}

output "alicloud_vpc_ipam_ipam_pool_allocation_example_id" {
  value = data.alicloud_vpc_ipam_ipam_pool_allocations.default.allocations.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cidr` - (ForceNew, Optional) The allocated address segment.
* `ipam_pool_allocation_id` - (ForceNew, Optional) The instance ID of the ipam pool allocation.
* `ipam_pool_allocation_name` - (ForceNew, Optional) The name of the ipam pool allocation.It must be 1 to 128 characters in length and cannot start with 'http:// 'or 'https.
* `ipam_pool_id` - (Required, ForceNew) The ID of the IPAM Pool.
* `ids` - (Optional, ForceNew, Computed) A list of Ipam Pool Allocation IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Ipam Pool Allocation IDs.
* `names` - A list of name of Ipam Pool Allocations.
* `allocations` - A list of Ipam Pool Allocation Entries. Each element contains the following attributes:
  * `cidr` - The allocated address segment.
  * `create_time` - Instance creation time.
  * `ipam_pool_allocation_description` - The description of the ipam pool alloctaion.It must be 1 to 256 characters in length and must start with an English letter or Chinese character, but cannot start with 'http:// 'or 'https. If it is not filled in, it is empty. The default value is empty.
  * `ipam_pool_allocation_id` - The instance ID of the ipam pool allocation.
  * `ipam_pool_allocation_name` - The name of the ipam pool allocation.It must be 1 to 128 characters in length and cannot start with 'http:// 'or 'https.
  * `ipam_pool_id` - The ID of the IPAM Pool.
  * `region_id` - When the IPAM Pool to which CIDR is allocated has the region attribute, this attribute is the IPAM Pool region.When the IPAM Pool to which CIDR is allocated does not have the region attribute, this attribute is the IPAM region.
  * `resource_id` - The ID of the resource.
  * `resource_owner_id` - The ID of the Alibaba Cloud account (primary account) to which the resource belongs.
  * `resource_region_id` - The region of the resource.
  * `resource_type` - The type of resource. Value:-**VPC**: indicates that the resource type is VPC.-**IpamPool**: indicates that the resource type is a child address pool.-**Custom**: indicates that the resource type is a Custom reserved CIDR block.
  * `source_cidr` - The source address segment.
  * `status` - The status of the instance. Value:-**Created**: indicates that the creation is complete.
  * `total_count` - Total number of records.
  * `id` - The ID of the resource supplied above.
