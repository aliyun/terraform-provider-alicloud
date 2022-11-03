---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk_attachment"
sidebar_current: "docs-alicloud-resource-ecs-disk-attachment"
description: |-
  Provides a ECS Disk Attachment resource.
---

# alicloud\_ecs\_disk\_attachment

Provides an Alicloud ECS Disk Attachment as a resource, to attach and detach disks from ECS Instances.

For information about ECS Disk Attachment and how to use it, see [What is Disk Attachment](https://www.alibabacloud.com/help/en/doc-detail/25515.htm).

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic usage

```terraform
# Create a new ECS disk-attachment and use it attach one disk to a new instance.
resource "alicloud_security_group" "ecs_sg" {
  name        = "terraform-test-group"
  description = "New security group"
}

resource "alicloud_ecs_disk" "ecs_disk" {
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

resource "alicloud_ecs_disk_attachment" "ecs_disk_att" {
  disk_id     = alicloud_ecs_disk.ecs_disk.id
  instance_id = alicloud_instance.ecs_instance.id
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the Instance to attach to.
* `disk_id` - (Required, ForceNew) ID of the Disk to be attached.
* `delete_with_instance` - (Optional, ForceNew) Indicates whether the disk is released together with the instance. Default to: `false`.
* `bootable` - (Optional, ForceNew) Whether to mount as a system disk. Default to: `false`.
* `key_pair_name` - (Optional, ForceNew) The name of key pair
* `password` - (Optional, ForceNew) When mounting the system disk, setting the user name and password of the instance is only effective for the administrator and root user names, and other user names are not effective.
* `device_name` - (Deprecated) The device name has been deprecated, and when attaching disk, it will be allocated automatically by system according to default order from /dev/xvdb to /dev/xvdz.

## Attributes Reference

The following attributes are exported:

* `id` - The Disk Attachment ID and it formats as `<disk_id>:<instance_id>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Disk.
* `delete` - (Defaults to 2 mins) Used when delete the Disk.

## Import

The disk attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_disk_attachment.example d-abc12345678:i-abc12355
```
