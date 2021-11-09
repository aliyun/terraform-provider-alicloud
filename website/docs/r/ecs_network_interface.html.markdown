---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_network_interface"
sidebar_current: "docs-alicloud-resource-ecs-network-interface"
description: |-
  Provides a Alicloud ECS Network Interface resource.
---

# alicloud\_ecs\_network\_interface

Provides a ECS Network Interface resource.

For information about ECS Network Interface and how to use it, see [What is Network Interface](https://www.alibabacloud.com/help/en/doc-detail/58504.htm).

-> **NOTE:** Available in v1.123.1+.

-> **NOTE** Only one of `private_ip_addresses` or `secondary_private_ip_address_count` can be specified when assign private IPs. 

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAcc"
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
* `name` - (Optional, Computed, Deprecated in v1.123.1+) Field `name` has been deprecated from provider version 1.123.1. New field `network_interface_name` instead
* `network_interface_name` - (Optional, Computed) The name of the ENI. The name must be 2 to 128 characters in length, and can contain letters, digits, colons (:), underscores (_), and hyphens (-). It must start with a letter and cannot start with http:// or https://.
* `primary_ip_address` - (Optional, Computed, ForceNew) The primary private IP address of the ENI. The specified IP address must be available within the CIDR block of the VSwitch. If this parameter is not specified, an available IP address is assigned from the VSwitch CIDR block at random.
* `private_ip` - (Optional, Computed, ForceNew, Deprecated in v1.123.1+) Field `private_ip` has been deprecated from provider version 1.123.1. New field `primary_ip_address` instead
* `private_ip_addresses` - (Optional, Computed) Specifies secondary private IP address N of the ENI. This IP address must be an available IP address within the CIDR block of the VSwitch to which the ENI belongs.
* `private_ips` - (Optional, Computed, Deprecated in v1.123.1+) Field `private_ips` has been deprecated from provider version 1.123.1. New field `private_ip_addresses` instead
* `private_ips_count` - (Optional, Computed, Deprecated in v1.123.1+) Field `private_ips_count` has been deprecated from provider version 1.123.1. New field `secondary_private_ip_address_count` instead
* `queue_number` - (Optional, Computed) The queue number of the ENI.
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `secondary_private_ip_address_count` - (Optional, Computed) The number of private IP addresses that can be automatically created by ECS.
* `security_group_ids` - (Optional, Computed) The ID of security group N. The security groups and the ENI must belong to the same VPC. The valid values of N are based on the maximum number of security groups to which an ENI can be added.
* `security_groups` - (Optional, Computed, Deprecated in v1.123.1+) Field `security_groups` has been deprecated from provider version 1.123.1. New field `security_group_ids` instead
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch in the specified VPC. The private IP addresses assigned to the ENI must be available IP addresses within the CIDR block of the VSwitch.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Network Interface.
* `mac` - The MAC address of the ENI.
* `status` - The status of the ENI.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Network Interface.
* `delete` - (Defaults to 1 mins) Used when delete the Network Interface.
* `update` - (Defaults to 1 mins) Used when update the Network Interface.

## Import

ECS Network Interface can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_network_interface.example eni-abcd12345
```
