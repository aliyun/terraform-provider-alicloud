---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_pool_allocation"
description: |-
  Provides a Alicloud Vpc Ipam Ipam Pool Allocation resource.
---

# alicloud_vpc_ipam_ipam_pool_allocation

Provides a Vpc Ipam Ipam Pool Allocation resource.

Allocates or reserves a CIDR from an IPAM address pool.

For information about Vpc Ipam Ipam Pool Allocation and how to use it, see [What is Ipam Pool Allocation](https://next.api.alibabacloud.com/document/VpcIpam/2023-02-28/CreateIpamPoolAllocation).

-> **NOTE:** Available since v1.238.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipam_ipam_pool_allocation&exampleId=2d463a79-b344-7e64-8af1-0600b4e9e2a48c0f69d0&activeTab=example&spm=docs.r.vpc_ipam_ipam_pool_allocation.0.2d463a79b3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
```

## Argument Reference

The following arguments are supported:
* `cidr` - (Optional, ForceNew) The allocated address segment.
* `cidr_mask` - (Optional, Int) Create a custom reserved network segment from The IPAM address pool by entering a mask.

-> **NOTE:**  Enter at least one of `Cidr` or **CidrMask.

* `ipam_pool_allocation_description` - (Optional) The description of the ipam pool alloctaion.
It must be 1 to 256 characters in length and must start with an English letter or Chinese character, but cannot start with 'http:// 'or 'https. If it is not filled in, it is empty. The default value is empty.
* `ipam_pool_allocation_name` - (Optional) The name of the ipam pool allocation.
It must be 1 to 128 characters in length and cannot start with 'http:// 'or 'https.
* `ipam_pool_id` - (Required, ForceNew) The ID of the IPAM Pool.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Instance creation time.
* `region_id` - When the IPAM Pool to which CIDR is allocated has the region attribute, this attribute is the IPAM Pool region.
When the IPAM Pool to which CIDR is allocated does not have the region attribute, this attribute is the IPAM region.
* `status` - The status of the instance. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipam Pool Allocation.
* `delete` - (Defaults to 5 mins) Used when delete the Ipam Pool Allocation.
* `update` - (Defaults to 5 mins) Used when update the Ipam Pool Allocation.

## Import

Vpc Ipam Ipam Pool Allocation can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_ipam_pool_allocation.example <id>
```