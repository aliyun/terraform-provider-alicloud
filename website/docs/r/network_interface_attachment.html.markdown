---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_network_interface_attachment"
sidebar_current: "docs-alicloud-resource-network-interface-attachment"
description: |-
  Provides an Alicloud ECS Elastic Network Interface Attachment as a resource to attach ENI to or detach ENI from ECS Instances.
---

# alicloud\_network\_interface\_attachment

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_network_interface_attachment](https://www.terraform.io/docs/providers/alicloud/r/ecs_network_interface_attachment) from version 1.123.1.

Provides an Alicloud ECS Elastic Network Interface Attachment as a resource to attach ENI to or detach ENI from ECS Instances.

For information about Elastic Network Interface and how to use it, see [Elastic Network Interface](https://www.alibabacloud.com/help/doc-detail/58496.html).

## Example Usage

Bacis Usage

```
variable "name" {
  default = "networkInterfaceAttachment"
}

variable "number" {
  default = "2"
}

resource "alicloud_vpc" "vpc" {
  name       = var.name
  cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
  vswitch_name      = var.name
  cidr_block        = "192.168.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vpc_id            = alicloud_vpc.vpc.id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = alicloud_vpc.vpc.id
}

data "alicloud_instance_types" "instance_type" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  eni_amount        = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "instance" {
  count             = var.number
  availability_zone = data.alicloud_zones.default.zones[0].id
  security_groups   = [alicloud_security_group.group.id]

  instance_type              = data.alicloud_instance_types.instance_type.instance_types[0].id
  system_disk_category       = "cloud_efficiency"
  image_id                   = data.alicloud_images.default.images[0].id
  instance_name              = var.name
  vswitch_id                 = alicloud_vswitch.vswitch.id
  internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface" "interface" {
  count           = var.number
  name            = var.name
  vswitch_id      = alicloud_vswitch.vswitch.id
  security_groups = [alicloud_security_group.group.id]
}

resource "alicloud_network_interface_attachment" "attachment" {
  count                = var.number
  instance_id          = element(alicloud_instance.instance.*.id, count.index)
  network_interface_id = element(alicloud_network_interface.interface.*.id, count.index)
}
```

## Argument Reference

The following argument are supported:

* `instance_id` - (Required, ForceNew) The instance ID to attach.
* `network_interface_id` - (Required, ForceNew) The ENI ID to attach.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource, formatted as `<network_interface_id>:<instance_id>`.

## Import

Network Interfaces Attachment resource can be imported using the id, e.g.

```
$ terraform import alicloud_network_interface_attachment.eni eni-abc123456789000:i-abc123456789000
```
