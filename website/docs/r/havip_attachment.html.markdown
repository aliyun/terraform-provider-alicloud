---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_havip_attachment"
sidebar_current: "docs-alicloud-resource-havip-attachment"
description: |-
  Provides an Alicloud HaVip Attachment resource.
---

# alicloud_havip_attachment

Provides an Alicloud HaVip Attachment resource for associating HaVip to ECS Instance.

-> **NOTE:** Terraform will auto build havip attachment while it uses `alicloud_havip_attachment` to build a havip attachment resource.

-> **NOTE:** Available since v1.18.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_havip" "example" {
  vswitch_id  = alicloud_vswitch.example.id
  description = var.name
}

resource "alicloud_security_group" "example" {
  name        = var.name
  description = var.name
  vpc_id      = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  vswitch_id                 = alicloud_vswitch.example.id
  image_id                   = data.alicloud_images.example.images.0.id
  instance_type              = data.alicloud_instance_types.example.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.example.id]
  instance_name              = var.name
  user_data                  = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

resource "alicloud_havip_attachment" "example" {
  havip_id    = alicloud_havip.example.id
  instance_id = alicloud_instance.example.id
}
```
## Argument Reference

The following arguments are supported:

* `havip_id` - (Required, ForceNew) The havip_id of the havip attachment, the field can't be changed.
* `instance_id` - (Required, ForceNew) The instance_id of the havip attachment, the field can't be changed.
* `force` - (Optional, Available since v1.200.0) Specifies whether to forcefully disassociate the HAVIP from the ECS instance or ENI. Default value: `False`. Valid values: `True` and `False`.
* `instance_type` - (Optional, ForceNew, Available since v1.201.0) The Type of instance to bind HaVip to. Valid values: `EcsInstance` and `NetworkInterface`. When the HaVip instance is bound to a resilient NIC, the resilient NIC instance must be filled in.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the havip attachment id and formates as `<havip_id>:<instance_id>`.
* `status` - (Available in v1.201.0+) The status of the HaVip instance.

## Timeouts

-> **NOTE:** Available since 1.194.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the HaVip Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the HaVip Attachment.

## Import

The havip attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_havip_attachment.foo havip-abc123456:i-abc123456
```
