---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_association"
sidebar_current: "docs-alicloud-resource-eip-association"
description: |-
  Provides a Alicloud ECS EIP Association resource.
---

# alicloud_eip_association

Provides an Alicloud EIP Association resource for associating Elastic IP to ECS Instance, SLB Instance or Nat Gateway.

-> **NOTE:** `alicloud_eip_association` is useful in scenarios where EIPs are either
 pre-existing or distributed to customers or users and therefore cannot be changed.

-> **NOTE:** From version 1.7.1, the resource support to associate EIP to SLB Instance or Nat Gateway.

-> **NOTE:** One EIP can only be associated with ECS or SLB instance which in the VPC.

-> **NOTE:** Available since v1.117.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  instance_name     = var.name
  image_id          = data.alicloud_images.example.images.0.id
  instance_type     = data.alicloud_instance_types.example.instance_types.0.id
  security_groups   = [alicloud_security_group.example.id]
  vswitch_id        = alicloud_vswitch.example.id
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_eip_address" "example" {
  address_name = var.name
}

resource "alicloud_eip_association" "example" {
  allocation_id = alicloud_eip_address.example.id
  instance_id   = alicloud_instance.example.id
}
```

## Module Support

You can use the existing [eip module](https://registry.terraform.io/modules/terraform-alicloud-modules/eip/alicloud) 
to create several EIP instances and associate them with other resources one-click, like ECS instances, SLB, Nat Gateway and so on.

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required, ForceNew) The ID of the EIP that you want to associate with an instance.
* `instance_id` - (Required, ForceNew) The ID of the ECS or SLB instance or Nat Gateway or NetworkInterface or HaVip.
* `instance_type` - (Optional, ForceNew, Available since v1.46.0) The type of the instance with which you want to associate the EIP. Valid values: `Nat`, `SlbInstance`, `EcsInstance`, `NetworkInterface`, `HaVip` and `IpAddress`.
* `mode` - (Optional, ForceNew, Available since v1.217.0) The association mode. Default value: `NAT`. Valid values: `NAT`, `BINDED`, `MULTI_BINDED`. **Note:** This parameter is required only when `instance_type` is set to `NetworkInterface`.
* `vpc_id` - (Optional, ForceNew, Available since v1.203.0) The ID of the VPC that has IPv4 gateways enabled and that is deployed in the same region as the EIP. When you associate an EIP with an IP address, the system can enable the IP address to access the Internet based on VPC route configurations. **Note:** This parameter is required if `instance_type` is set to `IpAddress`.
* `private_ip_address` - (Optional, ForceNew, Available since v1.52.2) The IP address in the CIDR block of the vSwitch.
* `force` - (Optional, Bool, Available since v1.95.0) When EIP is bound to a NAT gateway, and the NAT gateway adds a DNAT or SNAT entry, set it for `true` can unassociation any way. Default value: `false`. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of EIP Association. It formats as `<allocation_id>:<instance_id>`.

## Timeouts

-> **NOTE:** Available since 1.194.1.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Elastic IP address association.
* `delete` - (Defaults to 10 mins) Used when delete the Elastic IP address association.

## Import

-> **NOTE:** Available since 1.117.0.

Elastic IP address association can be imported using the id, e.g.

```shell
$ terraform import alicloud_eip_association.example <allocation_id>:<instance_id>
```
