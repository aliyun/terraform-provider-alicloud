---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_key_pair_attachment"
sidebar_current: "docs-alicloud-resource-key-pair-attachment"
description: |-
  Provides a Alicloud key pair attachment resource to bind key pair for several ECS instances.
---

# alicloud\_key\_pair\_attachment

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_key_pair_attachment](https://www.terraform.io/docs/providers/alicloud/r/ecs_key_pair_attachment) from version 1.121.0.

Provides a key pair attachment resource to bind key pair for several ECS instances.

-> **NOTE:** After the key pair is attached with sone instances, there instances must be rebooted to make the key pair affect.

## Example Usage

Basic Usage

```
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_ssd"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "type" {
  avaiability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "images" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "keyPairAttachmentName"
}

resource "alicloud_vpc" "vpc" {
  vpc_name       = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "10.1.1.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_security_group" "group" {
  name        = var.name
  description = "New security group"
  vpc_id      = alicloud_vpc.vpc.id
}

resource "alicloud_instance" "instance" {
  instance_name   = "${var.name}-${count.index + 1}"
  image_id        = data.alicloud_images.images.images[0].id
  instance_type   = data.alicloud_instance_types.type.instance_types[0].id
  count           = 2
  security_groups = [alicloud_security_group.group.id]
  vswitch_id      = alicloud_vswitch.vswitch.id

  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  password                   = "Test12345"

  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_ssd"
}

resource "alicloud_key_pair" "pair" {
  key_name = var.name
}

resource "alicloud_key_pair_attachment" "attachment" {
  key_name     = alicloud_key_pair.pair.id
  instance_ids = alicloud_instance.instance.*.id
}
```
## Argument Reference

The following arguments are supported:

* `key_name` - (Required, ForceNew) The name of key pair used to bind.
* `instance_ids` - (Required, ForceNew) The list of ECS instance's IDs.
* `force` - (ForceNew) Set it to true and it will reboot instances which attached with the key pair to make key pair affect immediately.

## Attributes Reference

* `key_name` - The name of the key pair.
* `instance_ids` The list of ECS instance's IDs.
