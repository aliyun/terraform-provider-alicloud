---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_pool_cidrs"
sidebar_current: "docs-alicloud-datasource-vpc-ipam-ipam-pool-cidrs"
description: |-
  Provides a list of Vpc Ipam Ipam Pool Cidr owned by an Alibaba Cloud account.
---

# alicloud_vpc_ipam_ipam_pool_cidrs

This data source provides Vpc Ipam Ipam Pool Cidr available to the user.[What is Ipam Pool Cidr](https://next.api.alibabacloud.com/document/VpcIpam/2023-02-28/AddIpamPoolCidr)

-> **NOTE:** Available since v1.241.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv4"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "default" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}

data "alicloud_vpc_ipam_ipam_pool_cidrs" "default" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool_cidr.default.ipam_pool_id
}

output "alicloud_vpc_ipam_ipam_pool_cidr_example_id" {
  value = data.alicloud_vpc_ipam_ipam_pool_cidrs.default.cidrs.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cidr` - (ForceNew, Optional) The CIDR address segment to be preset.> currently, only IPv4 address segments are supported.
* `ipam_pool_id` - (Required, ForceNew) The ID of the IPAM pool instance.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `cidrs` - A list of Ipam Pool Cidr Entries. Each element contains the following attributes:
  * `cidr` - The CIDR address segment to be preset.> currently, only IPv4 address segments are supported.
  * `ipam_pool_id` - The ID of the IPAM pool instance.
  * `status` - The status of the resource
  * `id` - The ID of the resource supplied above.
