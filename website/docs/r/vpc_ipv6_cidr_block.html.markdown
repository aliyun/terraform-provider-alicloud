---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_cidr_block"
description: |-
  Provides a Alicloud VPC Ipv6 Cidr Block resource.
---

# alicloud_vpc_ipv6_cidr_block

Provides a VPC Ipv6 Cidr Block resource.

VPC IPv6 supplementary CIDR block.

For information about VPC Ipv6 Cidr Block and how to use it, see [What is Ipv6 Cidr Block](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/AssociateVpcCidrBlock).

-> **NOTE:** Available since v1.280.0.

## Example Usage

Basic Usage

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
  ipv6_ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpv6Pool.id
  vpc_id            = alicloud_vpc.defaultVpc.id
  ipv6_cidr_block   = "fd03:d00:a000::/60"
}
```

## Argument Reference

The following arguments are supported:
* `ipv6_cidr_block` - (Optional, ForceNew, Computed) The additional IPv6 CIDR block to be removed.

-> **NOTE:**  You must specify either the `Ipv6CidrBlock` parameter or the `SecondaryCidrBlock` parameter, but not both.

* `ipv6_cidr_mask` - (Optional, Int) Add an IPv6 CIDR block to the VPC from an IPAM pool by specifying a subnet mask.

-> **NOTE:**  When adding an additional IPv6 CIDR block to a VPC from an IPAM pool, you must specify at least one of the `Ipv6CidrBlock` or `Ipv6CidrMask` parameters.


-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `ipv6_ipam_pool_id` - (Optional) The ID of the IPAM pool instance.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `vpc_id` - (Required, ForceNew) The ID of the VPC.
You can specify up to 20 VPC IDs, separated by commas (,).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<vpc_id>#<ipv6_cidr_block>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipv6 Cidr Block.
* `delete` - (Defaults to 5 mins) Used when delete the Ipv6 Cidr Block.

## Import

VPC Ipv6 Cidr Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv6_cidr_block.example <vpc_id>#<ipv6_cidr_block>
```