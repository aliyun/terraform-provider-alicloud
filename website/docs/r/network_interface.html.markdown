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

```terraform
variable "name" {
  default = "networkInterfaceName"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
  name       = var.name
  cidr_block = "192.168.0.0/24"
  zone_id    = data.alicloud_zones.default.zones[0].id
  vpc_id     = alicloud_vpc.vpc.id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = alicloud_vpc.vpc.id
}

resource "alicloud_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.vswitch.id
  security_group_ids     = [alicloud_security_group.group.id]
  private_ip             = "192.168.0.2"
  private_ips_count      = 3
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_network_interface&spm=docs.r.network_interface.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The VSwitch to create the ENI in.
* `network_interface_name` - (Optional) Name of the ENI.
* `description` - (Optional) Description of the ENI. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `primary_ip_address` - (Optional, ForceNew) The primary private IP of the ENI.
* `private_ip` - (Optional, ForceNew, Deprecated) Deprecated, use `primary_ip_address` instead.
* `private_ip_addresses` - (Optional) List of secondary private IPs to assign to the ENI.
* `private_ips` - (Optional, Deprecated) Deprecated, use `private_ip_addresses` instead.
* `queue_number` - (Optional) The queue number of the ENI.
* `resource_group_id` - (Optional, ForceNew) The Id of resource group which the network interface belongs.
* `secondary_private_ip_address_count` - (Optional) Number of secondary private IPs to assign to the ENI.
* `private_ips_count` - (Optional, Deprecated) Deprecated, use `secondary_private_ip_address_count` instead.
* `security_group_ids` - (Optional) A list of security group ids to associate with.
* `security_groups` - (Optional, Deprecated) Deprecated, use `security_group_ids` instead.
* `name` - (Optional, Deprecated) Deprecated, use `network_interface_name` instead.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `ipv6_address_count` - (Optional) The number of IPv6 addresses to assign to the ENI.
* `ipv6_addresses` - (Optional) A list of IPv6 addresses to assign to the ENI.
* `ipv4_prefix_count` - (Optional) The number of IPv4 prefixes to assign to the ENI.
* `ipv4_prefixes` - (Optional) A list of IPv4 prefixes to assign to the ENI.
* `instance_type` - (Optional, ForceNew) The instance type of the ENI.
* `network_interface_traffic_mode` - (Optional, ForceNew) The traffic mode of the ENI.
* `source_dest_check` - (Optional) Whether to enable SourceDestCheck of the ENI.
* `delete_on_release` - (Optional) Whether to delete the ENI when it is released.

## Attributes Reference

The following attributes are exported:

* `id` - The ENI ID.
* `mac` - (Available since v1.54.0) The MAC address of an ENI.
* `status` - The status of the ENI.

## Import

ENI can be imported using the id, e.g.

```shell
$ terraform import alicloud_network_interface.eni eni-abc1234567890000
```
