---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_mount_point"
sidebar_current: "docs-alicloud-resource-dfs-mount-point"
description: |-
  Provides a Alicloud DFS Mount Point resource.
---

# alicloud_dfs_mount_point

Provides a DFS Mount Point resource.

For information about DFS Mount Point and how to use it, see [What is Mount Point](https://www.alibabacloud.com/help/doc-detail/207144.htm).

-> **NOTE:** Available since v1.140.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_dfs_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_dfs_zones.default.zones.0.zone_id
}

resource "alicloud_dfs_file_system" "default" {
  storage_type                     = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
  zone_id                          = data.alicloud_dfs_zones.default.zones.0.zone_id
  protocol_type                    = "HDFS"
  description                      = var.name
  file_system_name                 = var.name
  throughput_mode                  = "Provisioned"
  space_capacity                   = "1024"
  provisioned_throughput_in_mi_bps = "512"
}

resource "alicloud_dfs_access_group" "default" {
  access_group_name = var.name
  description       = var.name
  network_type      = "VPC"
}

resource "alicloud_dfs_mount_point" "default" {
  description     = var.name
  vpc_id          = alicloud_vpc.default.id
  file_system_id  = alicloud_dfs_file_system.default.id
  access_group_id = alicloud_dfs_access_group.default.id
  network_type    = "VPC"
  vswitch_id      = alicloud_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Required) The ID of the Access Group.
* `description` - (Optional) The description of the Mount Point.
* `file_system_id` - (Required, ForceNew) The ID of the File System.
* `network_type` - (Required, ForceNew) The network type of the Mount Point. Valid values: `VPC`.
* `status` - (Optional) The status of the Mount Point. Valid values: `Active`, `Inactive`.
* `vpc_id` - (Required, ForceNew) The vpc id.
* `vswitch_id` - (Required, ForceNew) The vswitch id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the Mount Point. The value formats as `<file_system_id>:<mount_point_id>`.
* `mount_point_id` - The ID of the Mount Point.

## Import

DFS Mount Point can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_mount_point.example <file_system_id>:<mount_point_id>
```
