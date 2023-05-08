---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vswitch"
sidebar_current: "docs-alicloud-resource-vswitch"
description: |-
  Provides a Alicloud Vpc Vswitch resource.
---

# alicloud_vswitch

Provides a VPC switch resource.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id     = alicloud_vpc.vpc.id
  cidr_block = "172.16.0.0/21"
  zone_id    = "cn-beijing-b"
}

```

```terraform
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc_ipv4_cidr_block" "cidr_blocks" {
  vpc_id               = alicloud_vpc.vpc.id
  secondary_cidr_block = "192.163.0.0/16"
}

resource "alicloud_vswitch" "island-nat" {
  vpc_id       = alicloud_vpc_ipv4_cidr_block.cidr_blocks.vpc_id
  cidr_block   = "172.16.0.0/21"
  zone_id      = "cn-beijing-b"
  vswitch_name = "example_value"
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
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc_ipv4_cidr_block" "example" {
  vpc_id               = alicloud_vpc.default.id
  secondary_cidr_block = "192.163.0.0/16"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id     = alicloud_vpc_ipv4_cidr_block.example.vpc_id
  cidr_block = "192.163.0.0/24"
  zone_id    = "cn-beijing-b"
}
```

## Module Support

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud) 
to create a VPC and several VSwitches one-click.

## Argument Reference

The following arguments are supported:
* `all` - (Optional) Whether to unbind all tags of the resource. Value:-**true**: untags all resources.-**false** (default): does not remove all tags of the resource.
* `cidr_block` - (Required, ForceNew) The IPv4 CIDR block of the VSwitch.
* `description` - (Optional) The description of VSwitch.
* `enable_ipv6` - (Optional) Whether the IPv6 function is enabled in the switch. Value:-**true**: enables IPv6.-**false** (default): IPv6 is not enabled.
* `ipv6_cidr_block_mask` - (Optional) The IPv6 CIDR block of the VSwitch.
* `tags` - (Optional) The tags of VSwitch.
* `vswitch_name` - (Optional) The name of the VSwitch.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `zone_id` - (Required, ForceNew) The zone ID  of the resource

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated) Field 'name' has been deprecated from provider version 1.119.0. New field 'vswitch_name' instead.
* `availability_zone` - (Deprecated) Field 'availability_zone' has been deprecated from provider version 1.119.0. New field 'zone_id' instead.

### Timeouts

-> **NOTE:** Available in 1.79.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the vswitch (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the vswitch. 

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `available_ip_address_count` - The number of available IP addresses.
* `create_time` - The creation time of the VSwitch.
* `ipv6_cidr_block` - The IPv6 CIDR block of the VSwitch.
* `is_default` - Indicates whether the VSwitch is a default VSwitch.
* `network_acl_id` - The ID of the network ACL.
* `resource_group_id` - The resource group id of VSwitch.
* `route_table_id` - The route table id
* `status` - The status of the resource
* `vswitch_id` - The ID of the VSwitch.

## Import

Vpc Vswitch can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_vswitch.example <id>
```