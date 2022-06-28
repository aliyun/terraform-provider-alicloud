---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_network_interface"
sidebar_current: "docs-alicloud-resource-network-interface"
description: |-
  Provides an ECS Elastic Network Interface resource.
---

# alicloud\_network\_interface

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_network_interface](https://www.terraform.io/docs/providers/alicloud/r/ecs_network_interface) from version 1.123.1.

Provides an ECS Elastic Network Interface resource.

For information about Elastic Network Interface and how to use it, see [Elastic Network Interface](https://www.alibabacloud.com/help/doc-detail/58496.html).

-> **NOTE** Only one of private_ips or private_ips_count can be specified when assign private IPs. 

## Example Usage

```
variable "name" {
  default = "networkInterfaceName"
}

resource "alicloud_vpc" "vpc" {
  vpc_name       = var.name
  cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
  name              = var.name
  cidr_block        = "192.168.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vpc_id            = alicloud_vpc.vpc.id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = alicloud_vpc.vpc.id
}

resource "alicloud_network_interface" "default" {
  name              = var.name
  vswitch_id        = alicloud_vswitch.vswitch.id
  security_groups   = [alicloud_security_group.group.id]
  private_ip        = "192.168.0.2"
  private_ips_count = 3
}
```

## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The VSwitch to create the ENI in.
* `security_groups` - (Require) A list of security group ids to associate with.
* `private_ip` - (Optional, ForceNew) The primary private IP of the ENI.
* `name` - (Optional) Name of the ENI. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `description` - (Optional) Description of the ENI. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `private_ips`  - (Optional) List of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `private_ips_count` - (Optional) Number of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `resource_group_id` - (ForceNew, ForceNew, Available in 1.57.0+) The Id of resource group which the network interface belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ENI ID.
* `mac` - (Available in 1.54.0+) The MAC address of an ENI.

## Import

ENI can be imported using the id, e.g.

```
$ terraform import alicloud_network_interface.eni eni-abc1234567890000
```
