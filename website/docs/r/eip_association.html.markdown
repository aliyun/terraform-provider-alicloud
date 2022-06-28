---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_association"
sidebar_current: "docs-alicloud-resource-eip-association"
description: |-
  Provides a ECS EIP Association resource.
---

# alicloud\_eip\_association

Provides an Alicloud EIP Association resource for associating Elastic IP to ECS Instance, SLB Instance or Nat Gateway.

-> **NOTE:** `alicloud_eip_association` is useful in scenarios where EIPs are either
 pre-existing or distributed to customers or users and therefore cannot be changed.

-> **NOTE:** From version 1.7.1, the resource support to associate EIP to SLB Instance or Nat Gateway.

-> **NOTE:** One EIP can only be associated with ECS or SLB instance which in the VPC.

## Example Usage

```
# Create a new EIP association and use it to associate a EIP form a instance.
data "alicloud_zones" "default" {
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "10.1.1.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id

  depends_on = [alicloud_vpc.vpc]
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "ecs_instance" {
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = data.alicloud_instance_types.default.instance_types[0].id
  availability_zone = data.alicloud_zones.default.zones[0].id
  security_groups   = [alicloud_security_group.group.id]
  vswitch_id        = alicloud_vswitch.vsw.id
  instance_name     = "hello"
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alicloud_eip_address" "eip" {
}

resource "alicloud_eip_association" "eip_asso" {
  allocation_id = alicloud_eip_address.eip.id
  instance_id   = alicloud_instance.ecs_instance.id
}

resource "alicloud_security_group" "group" {
  name        = "terraform-test-group"
  description = "New security group"
  vpc_id      = alicloud_vpc.vpc.id
}
```

## Module Support

You can use the existing [eip module](https://registry.terraform.io/modules/terraform-alicloud-modules/eip/alicloud) 
to create several EIP instances and associate them with other resources one-click, like ECS instances, SLB, Nat Gateway and so on.

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required, ForcesNew) The allocation EIP ID.
* `instance_id` - (Required, ForcesNew) The ID of the ECS or SLB instance or Nat Gateway or NetworkInterface or HaVip.
* `instance_type` - (Optional, ForceNew, Available in 1.46.0+) The type of cloud product that the eip instance to bind. Valid values: `EcsInstance`, `SlbInstance`, `Nat`, `NetworkInterface` and `HaVip`.
* `private_ip_address` - (Optional, ForceNew, Available in 1.52.2+) The private IP address in the network segment of the vswitch which has been assigned.
* `force` - (Optional, Available in 1.95.0+) When EIP is bound to a NAT gateway, and the NAT gateway adds a DNAT or SNAT entry, set it for `true` can unassociation any way. Default to `false`.


## Attributes Reference

The following attributes are exported:

* `id` - The EIP Association ID and it formats as `<allocation_id>:<instance_id>`.

## Import

-> **NOTE:** Available in 1.117.0+.

Elastic IP address association can be imported using the id, e.g.

```
$ terraform import alicloud_eip_association.example eip-abc12345678:i-abc12355
```
