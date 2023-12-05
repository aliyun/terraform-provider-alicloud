---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_key_pair_attachment"
sidebar_current: "docs-alicloud-resource-ecs-key-pair-attachment"
description: |-
  Provides a Alicloud ECS Key Pair Attachment resource.
---

# alicloud_ecs_key_pair_attachment

Provides a ECS Key Pair Attachment resource.

For information about ECS Key Pair Attachment and how to use it, see [What is Key Pair Attachment](https://www.alibabacloud.com/help/en/doc-detail/51775.htm).

-> **NOTE:** Available since v1.121.0+.

## Example Usage

Basic Usage

```terraform
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
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  image_id             = data.alicloud_images.example.images.0.id
  instance_type        = data.alicloud_instance_types.example.instance_types.0.id
  availability_zone    = data.alicloud_zones.example.zones.0.id
  security_groups      = [alicloud_security_group.example.id]
  instance_name        = "terraform-example"
  internet_charge_type = "PayByBandwidth"
  vswitch_id           = alicloud_vswitch.example.id
}

resource "alicloud_ecs_key_pair" "example" {
  key_pair_name = "tf-example"
}

resource "alicloud_ecs_key_pair_attachment" "example" {
  key_pair_name = alicloud_ecs_key_pair.example.key_pair_name
  instance_ids  = [alicloud_instance.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Optional, ForceNew) The name of key pair used to bind.
* `key_name` - (Deprecated since v1.121.0+) New field 'key_pair_name' instead.
* `force` - (Optional, ForceNew) Set it to true and it will reboot instances which attached with the key pair to make key pair affect immediately.
* `instance_ids` - (Required, ForceNew) The list of ECS instance's IDs.

## Attributes Reference
 
The following attributes are exported:

* `id` - The resource ID of Key Pair Attachment. The value is formatted `<key_pair_name>:<instance_ids>`.

## Import

ECS Key Pair Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_key_pair_attachment.example <key_pair_name>:<instance_ids>
```
