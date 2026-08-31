---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_vswitch_cidr_reservation"
sidebar_current: "docs-alicloud-resource-vpc-vswitch-cidr-reservation"
description: |-
  Provides a Alicloud Vpc Vswitch Cidr Reservation resource.
---

# alicloud_vpc_vswitch_cidr_reservation

Provides a Vpc Vswitch Cidr Reservation resource. The reserved network segment of the vswitch. This resource type can be used only in ap-southeast region.

For information about Vpc Vswitch Cidr Reservation and how to use it, see [What is Vswitch Cidr Reservation](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/610154).

-> **NOTE:** Available since v1.205.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_vswitch_cidr_reservation&exampleId=8711b02a-db71-e3f4-a48c-8bff13d881fb796a0d85&activeTab=example&spm=docs.r.vpc_vswitch_cidr_reservation.0.8711b02adb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "ap-southeast-2"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "10.0.0.0/20"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
}


resource "alicloud_vpc_vswitch_cidr_reservation" "default" {
  ip_version                    = "IPv4"
  vswitch_id                    = alicloud_vswitch.defaultVSwitch.id
  cidr_reservation_description  = var.name
  cidr_reservation_cidr         = "10.0.10.0/24"
  vswitch_cidr_reservation_name = var.name
  cidr_reservation_type         = "Prefix"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_vswitch_cidr_reservation&spm=docs.r.vpc_vswitch_cidr_reservation.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cidr_reservation_cidr` - (Optional, ForceNew) Reserved network segment CIdrBlock.
* `cidr_reservation_description` - (Optional) The description of the reserved CIDR block.
* `cidr_reservation_mask` - (Optional, ForceNew) Reserved segment mask.
* `cidr_reservation_type` - (Optional, ForceNew) Reserved CIDR Block Type.Valid values: `Prefix`. Default value: Prefix.
* `ip_version` - (Optional, ForceNew) Reserved ip version of network segment, valid values: `IPv4`, `IPv6`, default IPv4.
* `vswitch_cidr_reservation_name` - (Optional) The name of the resource.
* `vswitch_id` - (Required, ForceNew) The Id of the switch instance.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<vswitch_id>:<vswitch_cidr_reservation_id>`.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.
* `vswitch_cidr_reservation_id` - The resource attribute field of the resource ID.
* `vpc_id` - The id of the vpc instance to which the reserved CIDR block belongs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vswitch Cidr Reservation.
* `delete` - (Defaults to 5 mins) Used when delete the Vswitch Cidr Reservation.
* `update` - (Defaults to 5 mins) Used when update the Vswitch Cidr Reservation.

## Import

Vpc Vswitch Cidr Reservation can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_vswitch_cidr_reservation.example <vswitch_id>:<vswitch_cidr_reservation_id>
```