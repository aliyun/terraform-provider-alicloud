---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_cidr_blocks"
sidebar_current: "docs-alicloud-datasource-vpc-ipv6-cidr-blocks"
description: |-
  Provides a list of VPC Ipv6 Cidr Block owned by an Alibaba Cloud account.
---

# alicloud_vpc_ipv6_cidr_blocks

This data source provides VPC Ipv6 Cidr Block available to the user.[What is Ipv6 Cidr Block](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/AssociateVpcCidrBlock)

-> **NOTE:** Available since v1.280.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpv6Pool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv6"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpv6PoolCidr" {
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpv6Pool.id
  cidr         = "fd03:d00:a000::/48"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "example-ipv6-cidr-block"
}


resource "alicloud_vpc_ipv6_cidr_block" "default" {
  ipv6_ipam_pool_id = alicloud_vpc_ipam_ipam_pool_cidr.defaultIpv6PoolCidr.ipam_pool_id
  vpc_id            = alicloud_vpc.defaultVpc.id
  ipv6_cidr_block   = "fd03:d00:a000::/60"
}

data "alicloud_vpc_ipv6_cidr_blocks" "default" {
  ids    = ["${alicloud_vpc_ipv6_cidr_block.default.id}"]
  vpc_id = alicloud_vpc.defaultVpc.id
}

output "alicloud_vpc_ipv6_cidr_block_example_id" {
  value = data.alicloud_vpc_ipv6_cidr_blocks.default.blocks.0.id
}
```

## Argument Reference

The following arguments are supported:
* `vpc_id` - (Required) The ID of the VPC.
You can specify up to 20 VPC IDs, separated by commas (,).
* `ids` - (Optional, Computed) A list of Ipv6 Cidr Block IDs. The value is formulated as `<vpc_id>#<ipv6_cidr_block>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Ipv6 Cidr Block IDs.
* `blocks` - A list of Ipv6 Cidr Block Entries. Each element contains the following attributes:
    * `ipv6_cidr_block` - The additional IPv6 CIDR block to be removed.
    * `id` - The ID of the resource supplied above.
