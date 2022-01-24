---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instance_attachment"
sidebar_current: "docs-alicloud-resource-dbfs-instance-attachment"
description: |-
  Provides a Alicloud DBFS Instance Attachment resource.
---

# alicloud\_dbfs\_instance\_attachment

Provides a DBFS Instance Attachment resource.

For information about DBFS Instance Attachment and how to use it, see [What is Instance Attachment](https://help.aliyun.com/document_detail/149726.html).

-> **NOTE:** Available in v1.156.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
locals {
  zone_id = "cn-hangzhou-i"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = local.zone_id
}
resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id      = data.alicloud_vpcs.default.ids[0]
}
data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^centos_8"
  most_recent = true
}
resource "alicloud_instance" "default" {
  image_id             = data.alicloud_images.default.images[0].id
  instance_name        = var.name
  instance_type        = "ecs.g7se.large"
  availability_zone    = local.zone_id
  vswitch_id           = data.alicloud_vswitches.default.ids[0]
  system_disk_category = "cloud_essd"
  security_groups = [
    alicloud_security_group.default.id
  ]
}
resource "alicloud_dbfs_instance" "default" {
  category          = "standard"
  zone_id           = alicloud_instance.default.availability_zone
  performance_level = "PL1"
  instance_name     = var.name
  size              = 100
}
resource "alicloud_dbfs_instance_attachment" "example" {
  ecs_id      = alicloud_instance.default.id
  instance_id = alicloud_dbfs_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `ecs_id` - (Required, ForceNew) The ID of the ECS instance.
* `instance_id` - (Required, ForceNew) The ID of the database file system.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance Attachment. The value formats as `<instance_id>:<ecs_id>`.
* `status` -The status of Database file system. Valid values: `attached`, `attaching`, `unattached`, `detaching`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.

## Import

DBFS Instance Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_dbfs_instance_attachment.example <instance_id>:<ecs_id>
```
