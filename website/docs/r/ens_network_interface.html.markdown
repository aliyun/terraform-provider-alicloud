---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_interface"
description: |-
  Provides a Alicloud ENS Network Interface resource.
---

# alicloud_ens_network_interface

Provides a ENS Network Interface resource.

For information about ENS Network Interface and how to use it, see [What is Network Interface](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createnetworkinterface).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "default" {
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = var.name
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_security_group" "default" {
  description         = var.name
  security_group_name = var.name
}

resource "alicloud_ens_network_interface" "default" {
  vswitch_id             = alicloud_ens_vswitch.default.id
  security_group_ids     = [alicloud_ens_security_group.default.id]
  network_interface_name = var.name
  description            = "network interface example"
}
```

## Argument Reference

The following arguments are supported:
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch to which the network interface belongs. The vSwitch must be in the same ENS region.
* `security_group_ids` - (Required, ForceNew) The list of security group IDs bound to the network interface. At most 5 security group IDs are supported.
* `network_interface_name` - (Optional) The name of the network interface.
* `description` - (Optional) The description of the network interface.
* `vmnc_learn` - (Optional) Specifies whether to enable MAC learning for the network interface. Valid values: `true`, `false`.
* `instance_id` - (Optional) The ID of the ENS instance to which the network interface is attached. Setting it attaches the network interface to the instance, and clearing it detaches the network interface.
* `private_ip_address_count` - (Optional) The number of secondary private IP addresses to maintain on the network interface. The provider reconciles the secondary private IP addresses to reach this count.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Network Interface. It is the same as `network_interface_id`.
* `network_interface_id` - The ID of the network interface.
* `status` - The status of the network interface.
* `primary_ip` - The primary private IP address of the network interface.
* `primary_ip_type` - The type of the primary private IP address.
* `mac_address` - The MAC address of the network interface.
* `network_id` - The ID of the network to which the network interface belongs.
* `ens_region_id` - The ENS region ID of the network interface.
* `create_time` - The time when the network interface was created.
* `private_ip_addresses` - The list of secondary private IP addresses of the network interface.
* `ipv6_addresses` - The list of IPv6 addresses of the network interface.
* `type` - The type of the network interface.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Interface.
* `delete` - (Defaults to 5 mins) Used when delete the Network Interface.
* `update` - (Defaults to 5 mins) Used when update the Network Interface.

## Import

ENS Network Interface can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_network_interface.example <id>
```
