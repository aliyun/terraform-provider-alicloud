---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_pool"
description: |-
  Provides a Alicloud Vpc Ipam Ipam Pool resource.
---

# alicloud_vpc_ipam_ipam_pool

Provides a Vpc Ipam Ipam Pool resource.

IP Address Management Pool.

For information about Vpc Ipam Ipam Pool and how to use it, see [What is Ipam Pool](https://next.api.alibabacloud.com/document/VpcIpam/2023-02-28/CreateIpamPool).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipam_ipam_pool&exampleId=2a160962-cd7a-9d64-32d2-d2a0f8a3cf67cf4b594f&activeTab=example&spm=docs.r.vpc_ipam_ipam_pool.0.2a160962cd&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc_ipam_ipam_pool" "parentIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  ipam_pool_name = format("%s1", var.name)
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
}


resource "alicloud_vpc_ipam_ipam_pool" "default" {
  ipam_scope_id       = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id      = alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id
  ipam_pool_name      = var.name
  source_ipam_pool_id = alicloud_vpc_ipam_ipam_pool.parentIpamPool.id
  ip_version          = "IPv4"
}
```

## Argument Reference

The following arguments are supported:
* `allocation_default_cidr_mask` - (Optional, Int) The default network mask assigned by the IPAM address pool.
IPv4 network mask value range: **0 to 32** bits.
* `allocation_max_cidr_mask` - (Optional, Int) The maximum network mask assigned by the IPAM address pool.
IPv4 network mask value range: **0 to 32** bits.
* `allocation_min_cidr_mask` - (Optional, Int) The minimum Network mask assigned by the IPAM address pool.
IPv4 network mask value range: **0 to 32** bits.
* `auto_import` - (Optional) Whether the automatic import function is enabled for the address pool.
* `clear_allocation_default_cidr_mask` - (Optional) Whether to clear the default network mask of the IPAM address pool. Value:
  - `true`: Yes.
  - `false`: No.
* `ip_version` - (Optional, ForceNew) The IP protocol version. Currently, only `IPv4` is supported * *.
* `ipam_pool_description` - (Optional) The description of the IPAM address pool.
It must be 2 to 256 characters in length and must start with an English letter or a Chinese character, but cannot start with 'http:// 'or 'https. If it is not filled in, it is empty. The default value is empty.
* `ipam_pool_name` - (Optional) The name of the resource.
* `ipam_scope_id` - (Required, ForceNew) Ipam scope id.
* `pool_region_id` - (Optional, ForceNew) The effective region of the IPAM address pool.
* `resource_group_id` - (Optional, Computed, Available since v1.242.0) The ID of the resource group.
* `source_ipam_pool_id` - (Optional, ForceNew) The instance ID of the source IPAM address pool.

-> **NOTE:**  If this parameter is not entered, the created address pool is the parent address pool.

* `tags` - (Optional, Map) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `region_id` - The ID of the IPAM hosting region.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipam Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Ipam Pool.
* `update` - (Defaults to 5 mins) Used when update the Ipam Pool.

## Import

Vpc Ipam Ipam Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_ipam_pool.example <id>
```