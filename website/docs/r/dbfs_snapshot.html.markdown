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
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_dbfs_snapshot&exampleId=830f5098-7bea-1db0-0d73-c8d0b342bef9e07158ee&activeTab=example&spm=docs.r.dbfs_snapshot.0.830f50987b" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

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
  name_regex    = "^aliyun_2_1903_x64_20G_alibase_20240628.vhd"
  owners        = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = local.zone_id
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids[0]
}

resource "alicloud_instance" "default" {
  availability_zone    = local.zone_id
  instance_name        = var.name
  image_id             = data.alicloud_images.example.images.0.id
  instance_type        = data.alicloud_instance_types.example.instance_types[length(data.alicloud_instance_types.example.instance_types) - 1].id
  security_groups      = [alicloud_security_group.example.id]
  vswitch_id           = data.alicloud_vswitches.default.ids.0
  system_disk_category = "cloud_essd"
}

resource "alicloud_dbfs_instance" "default" {
  category          = "enterprise"
  zone_id           = alicloud_instance.default.availability_zone
  performance_level = "PL1"
  fs_name           = var.name
  size              = 100
}

resource "alicloud_dbfs_instance_attachment" "default" {
  ecs_id      = alicloud_instance.default.id
  instance_id = alicloud_dbfs_instance.default.id
}
resource "alicloud_dbfs_snapshot" "example" {
  instance_id    = alicloud_dbfs_instance_attachment.default.instance_id
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