---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk_attachment"
sidebar_current: "docs-alicloud-resource-ecs-disk-attachment"
description: |-
  Provides a ECS Disk Attachment resource.
---

# alicloud_ecs_disk_attachment

Provides an Alicloud ECS Disk Attachment as a resource, to attach and detach disks from ECS Instances.

For information about ECS Disk Attachment and how to use it, see [What is Disk Attachment](https://www.alibabacloud.com/help/en/doc-detail/25515.htm).

-> **NOTE:** Available since v1.122.0.

## Example Usage

Basic usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_disk_attachment&exampleId=530cd864-561e-a815-a96f-1ae9af756b0726383d3a&activeTab=example&spm=docs.r.ecs_disk_attachment.0.530cd86456&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a new ECS disk-attachment and use it attach one disk to a new instance.
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "10.4.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name        = "tf-example"
  description = "New security group"
  vpc_id      = alicloud_vpc.default.id
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = var.name
  host_name         = var.name
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

data "alicloud_zones" "disk" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_ecs_disk" "default" {
  zone_id              = data.alicloud_zones.disk.zones.0.id
  category             = "cloud_efficiency"
  delete_auto_snapshot = "true"
  description          = "Test For Terraform"
  disk_name            = var.name
  enable_auto_snapshot = "true"
  encrypted            = "true"
  size                 = "500"
  tags = {
    Created     = "TF"
    Environment = "Acceptance-test"
  }
}

resource "alicloud_ecs_disk_attachment" "default" {
  disk_id     = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_disk_attachment&spm=docs.r.ecs_disk_attachment.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the Instance to attach to.
* `disk_id` - (Required, ForceNew) ID of the Disk to be attached.
* `delete_with_instance` - (Optional) Indicates whether the disk is released together with the instance. Default to: `false`.
* `bootable` - (Optional, ForceNew) Whether to mount as a system disk. Default to: `false`.
* `key_pair_name` - (Optional, ForceNew) The name of key pair
* `password` - (Optional, ForceNew) When mounting the system disk, setting the user name and password of the instance is only effective for the administrator and root user names, and other user names are not effective.

## Attributes Reference

The following attributes are exported:

* `id` - The Disk Attachment ID and it formats as `<disk_id>:<instance_id>`.
* `device` - The name of the cloud disk device.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Disk.
* `delete` - (Defaults to 2 mins) Used when delete the Disk.

## Import

The disk attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_disk_attachment.example d-abc12345678:i-abc12355
```
