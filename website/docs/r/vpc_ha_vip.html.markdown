---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ha_vip"
description: |-
  Provides a Alicloud VPC Ha Vip resource.
---

# alicloud_vpc_ha_vip

Provides a VPC Ha Vip resource.

Highly available virtual IP.

For information about VPC Ha Vip and how to use it, see [What is Ha Vip](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createhavip).

-> **NOTE:** Available since v1.205.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ha_vip&exampleId=02b304c1-2253-1cf3-d880-36fca9b32972866b9cfe&activeTab=example&spm=docs.r.vpc_ha_vip.0.02b304c122&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  description = var.name
  vpc_name    = var.name
  cidr_block  = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVswitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "192.168.0.0/21"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
  description  = var.name
}

resource "alicloud_resource_manager_resource_group" "defaultRg" {
  display_name        = "tf-example-defaultRg"
  resource_group_name = "${var.name}2"
}

resource "alicloud_resource_manager_resource_group" "changeRg" {
  display_name        = "tf-example-changeRg"
  resource_group_name = "${var.name}3"
}


resource "alicloud_vpc_ha_vip" "default" {
  description       = var.name
  vswitch_id        = alicloud_vswitch.defaultVswitch.id
  ha_vip_name       = var.name
  ip_address        = "192.168.1.101"
  resource_group_id = alicloud_resource_manager_resource_group.defaultRg.id
}
```

## Argument Reference

The following arguments are supported:
* `associated_instance_type` - (Optional, Computed) The type of the instance that is bound to the HaVip. Value:
  - `EcsInstance`: ECS instance.
  - `NetworkInterface`: ENI instance.
* `associated_instances` - (Optional, Computed, List) The ID of the ECS instance to be associated with the HAVIP. 
* `description` - (Optional) The description of the HAVIP. The description must be 1 to 255 characters in length and cannot start with `http://` or `https://`.
* `ha_vip_name` - (Optional, Computed) The name of the HAVIP. The name must be 1 to 128 characters in length, and cannot start with `http://` or `https://`.
* `ip_address` - (Optional, ForceNew, Computed) The IP address of the HAVIP. The specified IP address must be an idle IP address that falls within the CIDR block of the vSwitch. If this parameter is not set, an idle IP address from the CIDR block of the vSwitch is randomly assigned to the HAVIP. 
* `master_instance_id` - (Optional, Computed) The primary instance ID bound to HaVip
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the HAVIP belongs.
* `tags` - (Optional, Map) The tags of HaVip.
* `vswitch_id` - (Required, ForceNew) The switch ID to which the HaVip instance belongs

The following arguments will be discarded. Please use new fields as soon as possible:
* `havip_name` - (Deprecated since v1.259.0). Field 'havip_name' has been deprecated from provider version 1.259.0. New field 'ha_vip_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `associated_eip_addresses` - EIP bound to HaVip
* `create_time` - The creation time of the resource
* `ha_vip_id` - The ID of the HaVip instance.
* `status` - The status of this resource instance.
* `vpc_id` - The VPC ID to which the HaVip instance belongs

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ha Vip.
* `delete` - (Defaults to 5 mins) Used when delete the Ha Vip.
* `update` - (Defaults to 5 mins) Used when update the Ha Vip.

## Import

VPC Ha Vip can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ha_vip.example <id>
```