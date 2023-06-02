---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vswitch"
sidebar_current: "docs-alicloud-resource-vswitch"
description: |-
  Provides a Alicloud VPC Vswitch resource.
---

# alicloud_vswitch

Provides a VPC Vswitch resource. ## Module Support

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud)  to create a VPC and several VSwitches one-click.

For information about VPC Vswitch and how to use it, see [What is Vswitch](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/work-with-vswitches).

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

## Argument Reference

The following arguments are supported:
* `cidr_block` - (Required, ForceNew) The IPv4 CIDR block of the VSwitch.
* `description` - (Optional) The description of VSwitch.
* `zone_id` - (Optional, ForceNew, Available in 1.119.0+) The AZ for the VSwitch. **Note:** Required for a VPC VSwitch.
* `enable_ipv6` - (Optional, Available in v1.201.1+) Whether the IPv6 function is enabled in the switch. Value:
  - **true**: enables IPv6.
  - **false** (default): IPv6 is not enabled.
* `ipv6_cidr_block_mask` - (Optional, Available in v1.115+) The IPv6 CIDR block of the VSwitch.
* `tags` - (Optional, Map, Available in v1.55.3+) The tags of VSwitch.
* `vswitch_name` - (Optional, Available in v1.119.0+) The name of the VSwitch.
* `vpc_id` - (Required, ForceNew) The VPC ID.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated from v1.119.0+) Field 'name' has been deprecated from provider version 1.119.0. New field 'vswitch_name' instead.
* `availability_zone` - (Deprecated from v1.119.0+) Field 'availability_zone' has been deprecated from provider version 1.119.0. New field 'zone_id' instead.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VSwitch.
* `ipv6_cidr_block` - The IPv6 CIDR block of the VSwitch.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vswitch.
* `delete` - (Defaults to 5 mins) Used when delete the Vswitch.
* `update` - (Defaults to 5 mins) Used when update the Vswitch.

## Import

VPC Vswitch can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_vswitch.example <id>
```