---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_pool_cidr"
description: |-
  Provides a Alicloud Vpc Ipam Ipam Pool Cidr resource.
---

# alicloud_vpc_ipam_ipam_pool_cidr

Provides a Vpc Ipam Ipam Pool Cidr resource.

Ipam address pool preset CIDR.

For information about Vpc Ipam Ipam Pool Cidr and how to use it, see [What is Ipam Pool Cidr](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipam_ipam_pool_cidr&exampleId=2bdf0a0e-120f-0618-a92c-ba0a26dc8d1d4cbc241d&activeTab=example&spm=docs.r.vpc_ipam_ipam_pool_cidr.0.2bdf0a0e12&intl_lang=EN_US" target="_blank">
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
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv4"
}


resource "alicloud_vpc_ipam_ipam_pool_cidr" "default" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}
```

## Argument Reference

The following arguments are supported:
* `cidr` - (Required, ForceNew) The CIDR address segment to be preset.

-> **NOTE:**  currently, only IPv4 address segments are supported.

* `ipam_pool_id` - (Required, ForceNew) The ID of the IPAM pool instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ipam_pool_id>:<cidr>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipam Pool Cidr.
* `delete` - (Defaults to 5 mins) Used when delete the Ipam Pool Cidr.

## Import

Vpc Ipam Ipam Pool Cidr can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_ipam_pool_cidr.example <ipam_pool_id>:<cidr>
```