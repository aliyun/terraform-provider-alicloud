---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_ha_vip"
description: |-
  Provides a Alicloud ENS Ha Vip resource.
---

# alicloud_ens_ha_vip

Provides a ENS Ha Vip resource.

High-Availability Virtual IP Address.

For information about ENS Ha Vip and how to use it, see [What is Ha Vip](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateHaVip).

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

resource "alicloud_ens_network" "default4wYgcV" {
  network_name  = "tf-example"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultcW3Eib" {
  cidr_block    = "10.0.9.0/24"
  vswitch_name  = "tf-example"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default4wYgcV.id
}


resource "alicloud_ens_ha_vip" "default" {
  description = "desc1"
  vswitch_id  = alicloud_ens_vswitch.defaultcW3Eib.id
  amount      = 1
  ip_address  = "10.0.9.5"
  ha_vip_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `amount` - (Optional, ForceNew, Int) The number of highly available VIPs created. Value range: 1 to 10. When the specified IP address is created, the quantity can only be 1. Defaults to 1.
* `description` - (Optional, Computed, ForceNew) The description of the HaVip instance.
* `ha_vip_name` - (Optional, Computed) The name of the HaVip instance.
* `ip_address` - (Optional, Computed, ForceNew) The IP address of the HaVip. If not specified, the system automatically assigns one from the vSwitch CIDR block.
* `vswitch_id` - (Required, ForceNew) The vSwitch ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time of the resource.
* `ens_region_id` - The node ID of ENS.
* `network_id` - The network ID.
* `status` - The HaVip status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ha Vip.
* `delete` - (Defaults to 5 mins) Used when delete the Ha Vip.
* `update` - (Defaults to 5 mins) Used when update the Ha Vip.

## Import

ENS Ha Vip can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_ha_vip.example <ha_vip_id>
```