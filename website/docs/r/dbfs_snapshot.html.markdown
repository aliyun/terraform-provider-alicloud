---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_snapshot"
sidebar_current: "docs-alicloud-resource-dbfs-snapshot"
description: |-
  Provides a Alicloud DBFS Snapshot resource.
---

# alicloud_dbfs_snapshot

Provides a DBFS Snapshot resource.

For information about DBFS Snapshot and how to use it.

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

resource "alicloud_dbfs_snapshot" "example" {
  instance_id    = alicloud_dbfs_instance_attachment.example.instance_id
  snapshot_name  = var.name
  description    = var.name
  retention_days = 30
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) Description of the snapshot. The description must be `2` to `256` characters in length. It must start with a letter, and cannot start with `http://` or `https://`.
* `force` - (Optional) Whether to force deletion of snapshots.
* `instance_id` - (Required, ForceNew) The ID of the database file system.
* `retention_days` - (Optional, ForceNew) The retention time of the snapshot. Unit: days. Snapshots are automatically released after the retention time expires. Valid values: `1` to `65536`.
* `snapshot_name` - (Optional, ForceNew) The display name of the snapshot. The length is `2` to `128` characters. It must start with a large or small letter or Chinese, and cannot start with `http://` and `https://`. It can contain numbers, colons (:), underscores (_), or hyphens (-). To prevent name conflicts with automatic snapshots, you cannot start with `auto`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot.
* `status` - The status of the Snapshot.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Snapshot.
* `delete` - (Defaults to 1 mins) Used when delete the Snapshot.

## Import

DBFS Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbfs_snapshot.example <id>
```