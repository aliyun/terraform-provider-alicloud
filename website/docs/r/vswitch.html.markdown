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

```
resource "alicloud_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id            = "${alicloud_vpc.vpc.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "cn-beijing-b"
}
```

## Module Support

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud) 
to create a VPC and several VSwitches.

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The AZ for the switch.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `cidr_block` - (Required, ForceNew) The CIDR block for the switch.
* `name` - (Optional) The name of the switch. Defaults to null.
* `description` - (Optional) The switch description. Defaults to null.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `availability_zone` The AZ for the switch.
* `cidr_block` - The CIDR block for the switch.
* `vpc_id` - The VPC ID.
* `name` - The name of the switch.
* `description` - The description of the switch.

## Import

Vswitch can be imported using the id, e.g.

```
$ terraform import alicloud_vswitch.example vsw-abc123456
```
