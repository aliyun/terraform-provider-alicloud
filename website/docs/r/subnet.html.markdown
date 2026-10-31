---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_subnet"
sidebar_current: "docs-alicloud-resource-subnet"
description: |-
  Provides a Alicloud VPC Subnet resource.
---

# alicloud_subnet

Provides a VPC Subnet resource. `alicloud_subnet` is an alias of [`alicloud_vswitch`](vswitch.html.markdown); the two resource types share the same schema and behavior and may be used interchangeably.

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud)  to create a VPC and several VSwitches one-click.

For information about VPC Vswitch and how to use it, see [What is Vswitch](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/work-with-vswitches).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

```terraform

data "alicloud_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_subnet" "foo" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.16.0.0/21"
  vpc_id       = alicloud_vpc.foo.id
  zone_id      = data.alicloud_zones.foo.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional, ForceNew) The IPv4 CIDR block of the VSwitch. The subnet mask must be `16` to `29` bits in length. **NOTE:** From version 1.233.0, if you do not set `is_default`, or set `is_default` to `false`, `cidr_block` is required.
* `description` - (Optional) The description of VSwitch. The description must be `1` to `256` characters in length, and cannot start with `http://` or `https://`.
* `zone_id` - (Optional, ForceNew, Available since v1.119.0) The AZ for the VSwitch. **Note:** Required for a VPC VSwitch.
* `enable_ipv6` - (Optional, Computed, Available since v1.201.0) Whether the IPv6 function is enabled in the switch. Value:
  - `true`: enables IPv6.
  - `false` (default): IPv6 is not enabled.

  -> **NOTE:** Since v1.281.0, `ipv6_cidr_block_mask` only takes effect when `enable_ipv6` is `true`, and it is required when `enable_ipv6` is `true`.

  -> **NOTE:** When `enable_ipv6` is `true`, modifying `ipv6_cidr_block_mask` will disable IPv6 and then re-enable it to assign the new subnet. This re-assignment may interrupt the IPv6 service of the VSwitch, so plan such changes with care.

* `ipv6_cidr_block_mask` - (Optional, Available since v1.201.0) The last 4, 8, or 12 bits of the IPv6 CIDR block of the VSwitch, corresponding to a VPC IPv6 address mask of `60`, `56`, or `52` respectively. It only takes effect and is required when `enable_ipv6` is `true`, and is used only for create and update operations. The valid values are determined by the IPv6 address mask of the VPC:
  - When the VPC IPv6 address mask is `52`: `0` to `4095`.
  - When the VPC IPv6 address mask is `56`: `0` to `255`.
  - When the VPC IPv6 address mask is `60`: `0` to `15`.
* `tags` - (Optional, Map, Available since v1.55.3) The tags of VSwitch.
* `vswitch_name` - (Optional, Available since v1.119.0) The name of the VSwitch. The name must be `1` to `128` characters in length, and cannot start with `http://` or `https://`.
* `vpc_id` - (Optional, ForceNew) The VPC ID. **NOTE:** From version 1.233.0, if you do not set `is_default`, or set `is_default` to `false`, `vpc_id` is required.
* `is_default` - (Optional, Computed, ForceNew, Bool, Available since v1.233.0) Specifies whether to create the default VSwitch. Default value: `false`. Valid values:
  - `true`: Creates a default vSwitch.
  - `false`: Creates a vSwitch.

  -> **NOTE:** When `is_default` is `true`, a default VPC must already exist in the region. A zone can contain only one default vSwitch, and creating another one in the same zone fails. After a default vSwitch is created, its CIDR block cannot be modified. The first and last three IP addresses in each default vSwitch CIDR block are reserved by the system, default vSwitches do not support multicast or broadcast, and the number of cloud product instances in a default vSwitch cannot exceed the remaining available instance quota of the VPC.
* `vpc_ipv6_cidr_block` - (Optional, Available since v1.280.0) The IPv6 CIDR block of the VPC. If the VPC has multiple IPv6 CIDR blocks, you can use this parameter to specify the IPv6 CIDR block range to which the VSwitch belongs. This parameter is used only for create and update operations.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.119.0) Field `name` has been deprecated from provider version 1.119.0. New field `vswitch_name` instead.
* `availability_zone` - (Deprecated since v1.119.0) Field `availability_zone` has been deprecated from provider version 1.119.0. New field `zone_id` instead.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VSwitch.
* `ipv6_cidr_block` - The IPv6 CIDR block of the VSwitch.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vswitch.
* `delete` - (Defaults to 10 mins) Used when delete the Vswitch.
* `update` - (Defaults to 5 mins) Used when update the Vswitch.

## Import

VPC Vswitch can be imported using the id, e.g.

```shell
$ terraform import alicloud_subnet.example <id>
```
