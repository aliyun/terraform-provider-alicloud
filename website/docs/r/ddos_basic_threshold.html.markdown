---
subcategory: "Ddos Basic"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddos_basic_threshold"
sidebar_current: "docs-alicloud-resource-ddos-basic-threshold"
description: |-
  Provides a Alicloud Ddos Basic Threshold resource.
---

# alicloud_ddos_basic_threshold

Provides a Ddos Basic Threshold resource.

For information about Ddos Basic Threshold and how to use it, see [What is Threshold](https://www.alibabacloud.com/help/en/ddos-protection/latest/describe-ip-ddosthreshold).

-> **NOTE:** Available since v1.183.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  owners     = "system"
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "New security group"
  vpc_id      = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_name              = var.name
  host_name                  = var.name
  internet_max_bandwidth_out = 10
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = [alicloud_security_group.default.id]
  vswitch_id                 = alicloud_vswitch.default.id
}

resource "alicloud_ddos_basic_threshold" "example" {
  pps           = 60000
  bps           = 100
  internet_ip   = alicloud_instance.default.public_ip
  instance_id   = alicloud_instance.default.id
  instance_type = "ecs"
}
```

## Argument Reference

The following arguments are supported:

* `bps` - (Required) Specifies the traffic scrubbing threshold. Unit: Mbit/s. The traffic scrubbing threshold cannot exceed the peak inbound or outbound Internet traffic, whichever is larger, of the asset.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `instance_type` - (Required, ForceNew) The type of the Instance. Valid values: `ecs`,`slb`,`eip`.
* `internet_ip` - (Required, ForceNew) The IP address of the public IP address asset.
* `pps` - (Required) The current message number cleaning threshold. Unit: pps.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Threshold. The value formats as `<instance_type>:<instance_id>:<internet_ip>`.
* `max_bps` - Maximum flow cleaning threshold. Unit: Mbps.
* `max_pps` - The maximum number of messages cleaning threshold. Unit: pps.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ddos Threshold.
* `update` - (Defaults to 1 mins) Used when updating the Ddos Threshold.


## Import

Ddos Basic Threshold can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddos_basic_threshold.example <instance_type>:<instance_id>:<internet_ip>
```