---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_address"
description: |-
  Provides a Alicloud VPC Ipv6 Address resource.
---

# alicloud_vpc_ipv6_address

Provides a VPC Ipv6 Address resource. 

For information about VPC Ipv6 Address and how to use it, see [What is Ipv6 Address](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  ipv6_isp    = "BGP"
  cidr_block  = "172.168.0.0/16"
  enable_ipv6 = true
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vswich" {
  vpc_id               = alicloud_vpc.vpc.id
  cidr_block           = "172.168.0.0/24"
  zone_id              = data.alicloud_zones.default.zones.0.id
  vswitch_name         = var.name
  ipv6_cidr_block_mask = "1"
}

resource "alicloud_vpc_ipv6_address" "default" {
  resource_group_id        = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  vswitch_id               = alicloud_vswitch.vswich.id
  ipv6_address_description = var.name
  ipv6_address_name        = var.name
  tags = {
    Created = "TF"
    For     = "example"
  }
}
```

## Argument Reference

The following arguments are supported:
* `ipv6_address_description` - (Optional) The description of the IPv6 Address. The description must be 2 to 256 characters in length. It cannot start with http:// or https://.
* `ipv6_address_name` - (Optional) The name of the IPv6 Address. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with http:// or https://.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the instance belongs.
* `tags` - (Optional, Map) The tags for the resource.
* `vswitch_id` - (Required, ForceNew) The VSwitchId of the IPv6 address.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `ipv6_address` - IPv6 address.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.  Available, Pending and Deleting.

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