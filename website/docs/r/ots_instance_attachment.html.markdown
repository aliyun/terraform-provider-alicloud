---
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_instance_attachment"
sidebar_current: "docs-alicloud-resource-ots-instance-attachment"
description: |-
  Provides an OTS (Open Table Service) resource to attach VPC to instance.
---

# alicloud\_ots\_instance\_attachment

This resource will help you to bind a VPC to an OTS instance.

## Example Usage

```
# Create an OTS instance
resource "alicloud_ots_instance" "foo" {
  name = "my-ots-instance"
  description = "for table"
  accessed_by = "Vpc"
  tags {
    Created = "TF"
    For = "Building table"
  }
}

data "alicloud_zones" "foo" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/16"
  name = "for-ots-instance"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "for-ots-instance"
  cidr_block = "172.16.1.0/24"
  availability_zone = "${data.alicloud_zones.foo.zones.0.id}"
}
resource "alicloud_ots_instance_attachment" "foo" {
  instance_name = "${alicloud_ots_instance.foo.name}"
  vpc_name = "attachment1"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance.
* `vpc_name` - (Required, ForceNew) The name of attaching VPC to instance.
* `vswitch_id` - (Required, ForceNew) The ID of attaching VSwitch to instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is same as "instance_name".
* `instance_name` - The instance name.
* `vpc_name` - The name of attaching VPC to instance.
* `vswitch_id` - The ID of attaching VSwitch to instance.
* `vpc_id` - The ID of attaching VPC to instance.


