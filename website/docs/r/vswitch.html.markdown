---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vswitch"
sidebar_current: "docs-alicloud-resource-vswitch"
description: |-
  Provides a Alicloud VPC switch resource.
---

# alicloud\_vswitch

Provides a VPC switch resource.

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

resource "alicloud_vswitch" "foo" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.16.0.0/21"
  vpc_id       = alicloud_vpc.foo.id
  zone_id      = data.alicloud_zones.foo.zones.0.id
}

```

```terraform
data "alicloud_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc_ipv4_cidr_block" "cidr_blocks" {
  vpc_id               = alicloud_vpc.vpc.id
  secondary_cidr_block = "192.163.0.0/16"
}

resource "alicloud_vswitch" "island-nat" {
  vpc_id       = alicloud_vpc_ipv4_cidr_block.cidr_blocks.vpc_id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.foo.zones.0.id
  vswitch_name = "terraform-example"
  tags = {
    BuiltBy     = "example_value"
    cnm_version = "example_value"
    Environment = "example_value"
    ManagedBy   = "example_value"
  }
}

```

Create a switch associated with the additional network segment

```terraform
data "alicloud_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc_ipv4_cidr_block" "foo" {
  vpc_id               = alicloud_vpc.foo.id
  secondary_cidr_block = "192.163.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  vpc_id     = alicloud_vpc_ipv4_cidr_block.foo.vpc_id
  cidr_block = "192.163.0.0/24"
  zone_id    = data.alicloud_zones.foo.zones.0.id
}
```

## Module Support

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud) 
to create a VPC and several VSwitches one-click.

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, ForceNew, Deprecated from v1.119.0+) Field `availability_zone` has been deprecated from provider version 1.119.0. New field `zone_id` instead.
* `zone_id` - (Optional, ForceNew, Available in 1.119.0+) The AZ for the switch. **Note:** Required for a VPC switch.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `cidr_block` - (Required, ForceNew) The CIDR block for the switch.
* `name` - (Optional, Deprecated from v1.119.0+) Field `name` has been deprecated from provider version 1.119.0. New field `vswitch_name` instead.
* `vswitch_name` - (Optional, Available in 1.119.0+) The name of the switch. Defaults to null.
* `description` - (Optional) The switch description. Defaults to null.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `enable_ipv6` - (Optional, Available in 1.201.1+) Specifies whether to enable the IPv6 CIDR block. Valid values: `false` (Default): disables IPv6 CIDR blocks. `true`: enables IPv6 CIDR blocks.
* `ipv6_cidr_block_mask` - (Optional, Available in 1.201.1+) The last 8 bits of the switch's IPv6 segment, taking values: 0~255. This parameter is only supported to be configured when the VPC to which the switch belongs is IPv6 enabled.

### Timeouts

-> **NOTE:** Available in 1.79.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the vswitch (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the vswitch. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `availability_zone` The AZ for the switch.
* `cidr_block` - The CIDR block for the switch.
* `vpc_id` - The VPC ID.
* `name` - The name of the switch.
* `description` - The description of the switch.
* `status` - (Available in 1.119.0+) The status of the switch.
* `ipv6_cidr_block` - (Available in 1.201.1+) The IPv6 CIDR block for the switch.

## Import

Vswitch can be imported using the id, e.g.

```shell
$ terraform import alicloud_vswitch.example vsw-abc123456
```
