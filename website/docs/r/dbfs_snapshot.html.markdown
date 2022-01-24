---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_snapshot"
sidebar_current: "docs-alicloud-resource-dbfs-snapshot"
description: |-
  Provides a Alicloud DBFS Snapshot resource.
---

# alicloud\_dbfs\_snapshot

Provides a DBFS Snapshot resource.

For information about DBFS Snapshot and how to use it, see [What is Snapshot](https://help.aliyun.com/document_detail/149726.html).

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
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id      = data.alicloud_vpcs.default.ids.0
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

resource "alicloud_dbfs_instance_attachment" "default" {
  ecs_id      = alicloud_instance.default.id
  instance_id = alicloud_dbfs_instance.default.id
}

resource "alicloud_dbfs_snapshot" "example" {
  depends_on     = [alicloud_dbfs_instance_attachment.default]
  instance_id    = data.alicloud_dbfs_instances.default.ids.0
  snapshot_name  = "example_value"
  description    = "example_value"
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Snapshot.
* `delete` - (Defaults to 1 mins) Used when delete the Snapshot.

## Import

DBFS Snapshot can be imported using the id, e.g.

```
$ terraform import alicloud_dbfs_snapshot.example <id>
```