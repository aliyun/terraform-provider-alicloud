---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_interface"
description: |-
  Provides a Alicloud ENS Network Interface resource.
---

# alicloud_ens_network_interface

Provides a ENS Network Interface resource.

Elastic Network Card.

For information about ENS Network Interface and how to use it, see [What is Network Interface](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateNetworkInterface).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "ens_region_id" {
  default = "cn-hangzhou-58"
}

resource "alicloud_ens_security_group" "defaultvdnzzs" {
  security_group_name = "弹性网卡example用例使用"
}

resource "alicloud_ens_network" "defaultsLhpIw" {
  network_name  = "弹性网卡example用例"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultwsMJ1N" {
  cidr_block    = "10.0.5.0/24"
  vswitch_name  = "弹性网卡example用例"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultsLhpIw.id
}


resource "alicloud_ens_network_interface" "default" {
  description            = "desc"
  network_interface_name = "弹性网卡example用例"
  security_group_ids     = ["${alicloud_ens_security_group.defaultvdnzzs.id}"]
  vswitch_id             = alicloud_ens_vswitch.defaultwsMJ1N.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the ENI.
* `network_interface_name` - (Optional) The name of the ENI.
* `security_group_ids` - (Required, ForceNew, List) The ID of the security group.
* `vswitch_id` - (Optional, ForceNew) The vSwitch ID.
* `vmnc_learn` - (Optional, Computed) Whether to enable NIC route learning. Possible values:
  - true: On
  - false: Off (default)

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time.
* `ens_region_id` - The node ID of ENS.
* `instance_id` - The ID of the instance bound to the Eni.
* `mac_address` - The MAC address of the Eni.
* `network_id` - The network ID.
* `primary_ip` - The primary private network IP address.
* `primary_ip_type` - The primary IP address type.
* `status` - Status of network card, value:.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Interface.
* `delete` - (Defaults to 5 mins) Used when delete the Network Interface.
* `update` - (Defaults to 5 mins) Used when update the Network Interface.

## Import

ENS Network Interface can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_network_interface.example <network_interface_id>
```