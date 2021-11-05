---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_havip_attachment"
sidebar_current: "docs-alicloud-resource-havip-attachment"
description: |-
  Provides an Alicloud HaVip Attachment resource.
---

# alicloud\_havip\_attachment

Provides an Alicloud HaVip Attachment resource for associating HaVip to ECS Instance.

-> **NOTE:** Terraform will auto build havip attachment while it uses `alicloud_havip_attachment` to build a havip attachment resource.

## Example Usage

Basic Usage

```
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "test_havip_attachment"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name       = var.name
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = alicloud_vpc.foo.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_havip" "foo" {
  vswitch_id  = alicloud_vswitch.foo.id
  description = var.name
}

resource "alicloud_havip_attachment" "foo" {
  havip_id    = alicloud_havip.foo.id
  instance_id = alicloud_instance.foo.id
}

resource "alicloud_security_group" "tf_test_foo" {
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_instance" "foo" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  vswitch_id        = alicloud_vswitch.foo.id
  image_id          = data.alicloud_images.default.images[0].id

  # series III
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.tf_test_foo.id]
  instance_name              = var.name
  user_data                  = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}
```
## Argument Reference

The following arguments are supported:

* `havip_id` - (Required, ForceNew) The havip_id of the havip attachment, the field can't be changed.
* `instance_id` - (Required, ForceNew) The instance_id of the havip attachment, the field can't be changed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the havip attachment id and formates as `<havip_id>:<instance_id>`.

## Import

The havip attachment can be imported using the id, e.g.

```
$ terraform import alicloud_havip_attachment.foo havip-abc123456:i-abc123456
```
