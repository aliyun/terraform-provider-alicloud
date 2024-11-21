---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_address"
description: |-
  Provides a Alicloud VPC Ipv6 Address resource.
---

# alicloud_vpc_ipv6_address

Provides a VPC Ipv6 Address resource. 

For information about VPC Ipv6 Address and how to use it, see [What is Ipv6 Address](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/AllocateIpv6Address).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipv6_address&exampleId=c9582809-2e3d-7803-c6d4-a671c13f0148fcf0355d&activeTab=example&spm=docs.r.vpc_ipv6_address.0.c95828092e&intl_lang=EN_US" target="_blank">
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

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  ipv6_isp    = "BGP"
  cidr_block  = "172.168.0.0/16"
  enable_ipv6 = true
  vpc_name    = var.name

}

resource "alicloud_vswitch" "vswich" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.168.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name

  ipv6_cidr_block_mask = "1"
}


resource "alicloud_vpc_ipv6_address" "default" {
  resource_group_id        = data.alicloud_resource_manager_resource_groups.default.ids.0
  vswitch_id               = alicloud_vswitch.vswich.id
  ipv6_address_description = "create_description"
  ipv6_address_name        = var.name

}
```

## Argument Reference

The following arguments are supported:
* `ipv6_address_description` - (Optional, Computed) The description of the IPv6 Address. The description must be 2 to 256 characters in length. It cannot start with http:// or https://.
* `ipv6_address_name` - (Optional) The name of the IPv6 Address. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with http:// or https://.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the instance belongs.
* `tags` - (Optional, Map) The tags for the resource.
* `vswitch_id` - (Required, ForceNew) The VSwitchId of the IPv6 address.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.  Available, Pending and Deleting.
* `ipv6_address` - IPv6 address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipv6 Address.
* `delete` - (Defaults to 5 mins) Used when delete the Ipv6 Address.
* `update` - (Defaults to 5 mins) Used when update the Ipv6 Address.

## Import

VPC Ipv6 Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv6_address.example <id>
```