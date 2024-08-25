---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_disk_attachment"
sidebar_current: "docs-alicloud-resource-disk-attachment"
description: |-
  Provides a ECS Disk Attachment resource.
---

# alicloud\_disk\_attachment

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_disk_attachment](https://www.terraform.io/docs/providers/alicloud/r/ecs_disk_attachment) from version 1.122.0.

Provides an Alicloud ECS Disk Attachment as a resource, to attach and detach disks from ECS Instances.

## Example Usage

Basic usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_disk_attachment&exampleId=4c20ad79-88f7-275a-07b1-1ba5ccde1666703baca6&activeTab=example&spm=docs.r.disk_attachment.0.4c20ad7988&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a new ECS disk-attachment and use it attach one disk to a new instance.
resource "alicloud_security_group" "ecs_sg" {
  name        = "terraform-test-group"
  description = "New security group"
}

resource "alicloud_disk" "ecs_disk" {
  availability_zone = "cn-beijing-a"
  size              = "50"

  tags = {
    Name = "TerraformTest-disk"
  }
}

resource "alicloud_instance" "ecs_instance" {
  image_id             = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_type        = "ecs.n4.small"
  availability_zone    = "cn-beijing-a"
  security_groups      = [alicloud_security_group.ecs_sg.id]
  instance_name        = "Hello"
  internet_charge_type = "PayByBandwidth"

  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alicloud_disk_attachment" "ecs_disk_att" {
  disk_id     = alicloud_disk.ecs_disk.id
  instance_id = alicloud_instance.ecs_instance.id
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, Forces new resource) ID of the Instance to attach to.
* `disk_id` - (Required, Forces new resource) ID of the Disk to be attached.
* `device_name` - (Deprecated) The device name has been deprecated, and when attaching disk, it will be allocated automatically by system according to default order from /dev/xvdb to /dev/xvdz.

## Attributes Reference

The following attributes are exported:

* `id` - The Disk Attachment ID and it formats as `<disk_id>:<instance_id>`.

## Import

The disk attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_disk_attachment.example d-abc12345678:i-abc12355
```
