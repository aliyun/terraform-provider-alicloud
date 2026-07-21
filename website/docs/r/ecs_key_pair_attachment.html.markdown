---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_key_pair_attachment"
description: |-
  Provides a Alicloud ECS Key Pair Attachment resource.
---

# alicloud_ecs_key_pair_attachment

Provides a ECS Key Pair Attachment resource.

For information about ECS Key Pair Attachment and how to use it, see [What is Key Pair Attachment](https://www.alibabacloud.com/help/en/doc-detail/51775.htm).

-> **NOTE:** Available since v1.121.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_key_pair_attachment&exampleId=fcd97394-cdad-9ac9-c127-6f73f7f6f5ce25e0b7e9&activeTab=example&spm=docs.r.ecs_key_pair_attachment.0.fcd97394cd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  image_id          = data.alicloud_images.default.images.0.id
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_ecs_key_pair_attachment" "default" {
  key_pair_name = alicloud_ecs_key_pair.default.id
  instance_ids  = [alicloud_instance.default.id]
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_key_pair_attachment&spm=docs.r.ecs_key_pair_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, ForceNew, List) The IDs of instances to which you want to bind the SSH key pair.
* `key_pair_name` - (Optional, ForceNew) The name of the SSH key pair.
* `force` - (Optional, ForceNew, Bool) Specifies whether to make the key pair effective immediately. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `key_name` - (Optional, ForceNew, Deprecated since v1.121.0) Field `key_name` has been deprecated from provider version 1.121.0. New field `key_pair_name` instead.

-> **WARNING:**  If `force` set to `true`, it it will reboot instances which attached with the key pair to make key pair effective immediately.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Key Pair Attachment. It formats as `<key_pair_name>:<instance_ids>`.

## Timeouts

-> **NOTE:** Available since v1.274.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Key Pair Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Key Pair Attachment.

## Import

ECS Key Pair Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_key_pair_attachment.example <key_pair_name>:<instance_ids>
```
