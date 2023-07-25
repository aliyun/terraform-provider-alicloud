---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instance_attachment"
sidebar_current: "docs-alicloud-resource-dbfs-instance-attachment"
description: |-
  Provides a Alicloud DBFS Instance Attachment resource.
---

# alicloud_dbfs_instance_attachment

Provides a DBFS Instance Attachment resource.

For information about DBFS Instance Attachment and how to use it.

-> **NOTE:** Available since v1.156.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-hangzhou"
}

locals {
  zone_id = "cn-hangzhou-i"
}
data "alicloud_instance_types" "example" {
  availability_zone    = local.zone_id
  instance_type_family = "ecs.g7se"
}
data "alicloud_images" "example" {
  instance_type = data.alicloud_instance_types.example.instance_types[length(data.alicloud_instance_types.example.instance_types) - 1].id
  name_regex    = "^aliyun_2"
  owners        = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.zone_id
}
resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone    = local.zone_id
  instance_name        = var.name
  image_id             = data.alicloud_images.example.images.1.id
  instance_type        = data.alicloud_instance_types.example.instance_types[length(data.alicloud_instance_types.example.instance_types) - 1].id
  security_groups      = [alicloud_security_group.example.id]
  vswitch_id           = alicloud_vswitch.example.id
  system_disk_category = "cloud_essd"
}

resource "alicloud_dbfs_instance" "example" {
  category          = "standard"
  zone_id           = local.zone_id
  performance_level = "PL1"
  instance_name     = var.name
  size              = 100
}

resource "alicloud_dbfs_instance_attachment" "example" {
  ecs_id      = alicloud_instance.example.id
  instance_id = alicloud_dbfs_instance.example.id
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.

## Import

DBFS Instance Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbfs_instance_attachment.example <instance_id>:<ecs_id>
```
