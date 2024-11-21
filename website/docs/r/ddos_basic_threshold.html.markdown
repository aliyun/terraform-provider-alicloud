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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddos_basic_threshold&exampleId=c3124a4b-1b51-c0a8-1814-d43395a0c3928862ecd8&activeTab=example&spm=docs.r.ddos_basic_threshold.0.c3124a4b1b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
}

resource "alicloud_ddos_basic_threshold" "example" {
  instance_type = "ecs"
  instance_id   = alicloud_instance.default.id
  internet_ip   = alicloud_instance.default.public_ip
  bps           = 100
  pps           = 60000
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, ForceNew) The type of the Instance. Valid values: `ecs`,`slb`,`eip`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `internet_ip` - (Required, ForceNew) The IP address of the public IP address asset.
* `bps` - (Required) Specifies the traffic scrubbing threshold. Unit: Mbit/s. The traffic scrubbing threshold cannot exceed the peak inbound or outbound Internet traffic, whichever is larger, of the asset.
* `pps` - (Required) The current message number cleaning threshold. Unit: pps.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Threshold. It formats as `<instance_type>:<instance_id>:<internet_ip>`.
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
