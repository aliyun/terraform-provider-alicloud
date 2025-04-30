---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_havip_attachment"
description: |-
  Provides a Alicloud VPC Ha Vip Attachment resource.
---

# alicloud_havip_attachment

Provides a VPC Ha Vip Attachment resource. Attaching ECS instance to Havip.

-> **NOTE:** Terraform will auto build havip attachment while it uses `alicloud_havip_attachment` to build a havip attachment resource.

-> **NOTE:** Available since v1.18.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_havip_attachment&exampleId=cecc4d2a-58ce-bba9-6f2f-e03b832212d0d331f5e3&activeTab=example&spm=docs.r.havip_attachment.0.cecc4d2a58&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
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
  name_regex = "^ubuntu_18.*64"
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
  ha_vip_id   = alicloud_havip.example.id
  instance_id = alicloud_instance.example.id
}
```

## Argument Reference

The following arguments are supported:
* `force` - (Optional, Available since v1.18.0) Whether to force the ECS instance or Eni instance bound to AVIP to be unbound. The value is:
  - **True**: Force unbinding.
  - **False** (default): unbinding is not forced.
-> **NOTE:**  If the value of this parameter is **False**, the Master instance bound to HaVip cannot be unbound.
* `ha_vip_id` - (Optional, ForceNew, Available since v1.211.0) The ID of the HaVip instance.
* `instance_id` - (Required, ForceNew, Available since v1.18.0) The ID of the ECS instance bound to the HaVip instance.
* `instance_type` - (Optional, ForceNew, Available since v1.18.0) The type of the instance associated with the VIIP.

The following arguments will be discarded. Please use new fields as soon as possible:
* `havip_id` - (Deprecated since v1.211.0). Field 'havip_id' has been deprecated from provider version 1.211.0. New field 'ha_vip_id' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ha_vip_id>:<instance_id>`.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ha Vip Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Ha Vip Attachment.

## Import

VPC Ha Vip Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_havip_attachment.example <ha_vip_id>:<instance_id>
```