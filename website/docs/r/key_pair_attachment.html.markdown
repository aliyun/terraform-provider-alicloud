---
layout: "alicloud"
page_title: "Alicloud: alicloud_key_pair_attachment"
sidebar_current: "docs-alicloud-resource-key-pair-attachment"
description: |-
  Provides a Alicloud key pair attachment resource to bind key pair for several ECS instances.
---

# alicloud\_key\_pair\_attachment

Provides a key pair attachment resource to bind key pair for several ECS instances.

## Example Usage

Basic Usage

```
resource "alicloud_key_pair" "key" {
	key_name = "terraform-test-key-pair"
}

resource "alicloud_instance" "instance" {
  instance_name = "test-keypair-${format(var.count_format, count.index+1)}"
  image_id = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type = "ecs.n4.small"
  count = 2
  availability_zone = "${var.availability_zones}"
  ...
}

resource "alicloud_key_pair_attachment" "attach" {
  key_name = "${alicloud_key_pair.key.id}"
  instance_ids = ["${alicloud_instance.instance.*.id}"]
}
```
## Argument Reference

The following arguments are supported:

* `key_name` - (Required, Force new resource) The name of key pair used to bind.
* `instance_ids` - (Required, Force new resource) The list of ECS instance's IDs.

## Attributes Reference

* `key_name` - The name of the key pair.
* `instance_ids` The list of ECS instance's IDs.
