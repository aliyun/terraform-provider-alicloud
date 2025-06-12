---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_network_interface"
sidebar_current: "docs-alicloud-resource-ecs-network-interface"
description: |-
  Provides a Alicloud ECS Network Interface resource.
---

# alicloud_ecs_network_interface

Provides a ECS Network Interface resource.

For information about ECS Network Interface and how to use it, see [What is Network Interface](https://www.alibabacloud.com/help/en/doc-detail/58504.htm).

-> **NOTE:** Available since v1.123.1.

-> **NOTE** Only one of `private_ip_addresses` or `secondary_private_ip_address_count` can be specified when assign private IPs. 

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_network_interface&exampleId=424f8531-e314-8cfb-bb6e-ea1f70e0c71274f18f45&activeTab=example&spm=docs.r.ecs_network_interface.0.424f8531e3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "192.168.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vpc_id       = alicloud_vpc.default.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = "192.168.0.2"
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the ENI. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `network_interface_name` - (Optional) The name of the ENI. The name must be 2 to 128 characters in length, and can contain letters, digits, colons (:), underscores (_), and hyphens (-). It must start with a letter and cannot start with http:// or https://.
* `primary_ip_address` - (Optional, ForceNew) The primary private IP address of the ENI. The specified IP address must be available within the CIDR block of the VSwitch. If this parameter is not specified, an available IP address is assigned from the VSwitch CIDR block at random.
* `private_ip_addresses` - (Optional, List) Specifies secondary private IP address N of the ENI. This IP address must be an available IP address within the CIDR block of the VSwitch to which the ENI belongs.
* `queue_number` - (Optional, Int) The queue number of the ENI.
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `secondary_private_ip_address_count` - (Optional, Int) The number of private IP addresses that can be automatically created by ECS.
* `security_group_ids` - (Optional, List) The ID of security group N. The security groups and the ENI must belong to the same VPC. The valid values of N are based on the maximum number of security groups to which an ENI can be added. **NOTE:** Either `security_group_ids` or `security_groups` must be set with valid security group IDs.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch in the specified VPC. The private IP addresses assigned to the ENI must be available IP addresses within the CIDR block of the VSwitch.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `ipv6_address_count` - (Optional, Int, Available since v1.193.0) The number of IPv6 addresses to randomly generate for the primary ENI. Valid values: 1 to 10. **NOTE:** You cannot specify both the `ipv6_addresses` and `ipv6_address_count` parameters.
* `ipv6_addresses` - (Optional, List, Available since v1.193.0) A list of IPv6 address to be assigned to the primary ENI. Support up to 10.
* `ipv4_prefix_count` - (Optional, Int, Available since v1.213.0) The number of IPv4 prefixes that can be automatically created by ECS. Valid values: 1 to 10. **NOTE:** You cannot specify both the `ipv4_prefixes` and `ipv4_prefix_count` parameters.
* `ipv4_prefixes` - (Optional, List, Available since v1.213.0) A list of IPv4 prefixes to be assigned to the ENI. Support up to 10.
* `instance_type` - (Optional, ForceNew, Available since v1.223.0) The type of the ENI. Default value: `Secondary`. Valid values: `Secondary`, `Trunk`.
* `network_interface_traffic_mode` - (Optional, ForceNew, Available since v1.223.0) The communication mode of the ENI. Default value: `Standard`. Valid values: `Standard`, `HighPerformance`.
* `source_dest_check` - (Optional, Bool, Available since v1.252.0) Indicates whether the source and destination IP address check feature is enabled. To improve network security, enable this feature. Default value: `false`. Valid values: `true`, `false`.
* `name` - (Optional, Deprecated since v1.123.1) Field `name` has been deprecated from provider version 1.123.1. New field `network_interface_name` instead
* `private_ip` - (Optional, ForceNew, Deprecated since v1.123.1) Field `private_ip` has been deprecated from provider version 1.123.1. New field `primary_ip_address` instead
* `private_ips` - (Optional, List, Deprecated since v1.123.1) Field `private_ips` has been deprecated from provider version 1.123.1. New field `private_ip_addresses` instead
* `private_ips_count` - (Optional, Int, Deprecated since v1.123.1) Field `private_ips_count` has been deprecated from provider version 1.123.1. New field `secondary_private_ip_address_count` instead
* `security_groups` - (Optional, List, Deprecated since v1.123.1) Field `security_groups` has been deprecated from provider version 1.123.1. New field `security_group_ids` instead. **NOTE:** Either `security_group_ids` or `security_groups` must be set with valid security group IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Network Interface.
* `mac` - The MAC address of the ENI.
* `status` - The status of the ENI.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Network Interface.
* `update` - (Defaults to 1 mins) Used when update the Network Interface.
* `delete` - (Defaults to 1 mins) Used when delete the Network Interface.

## Import

ECS Network Interface can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_network_interface.example eni-abcd12345
```
