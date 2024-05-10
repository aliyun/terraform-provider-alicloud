---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_mount_point"
description: |-
  Provides a Alicloud DFS Mount Point resource.
---

# alicloud_dfs_mount_point

Provides a DFS Mount Point resource. 

For information about DFS Mount Point and how to use it, see [What is Mount Point](https://www.alibabacloud.com/help/en/aibaba-cloud-storage-services/latest/apsara-file-storage-for-hdfs).

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
* `access_group_id` - (Required) The id of the permission group associated with the Mount point, which is used to set the access permissions of the Mount point.
* `alias_prefix` - (Optional, Available since v1.218.0) The mount point alias prefix, which specifies the mount point alias prefix.
* `description` - (Optional) The description of the Mount point.  No more than 32 characters in length.
* `file_system_id` - (Required, ForceNew) Unique file system identifier, used to retrieve specified file system resources.
* `network_type` - (Required, ForceNew) The network type of the Mount point.  Only VPC (VPC) is supported.
* `status` - (Optional, Computed) Mount point status. Value: Inactive: Disable mount points Active: Activate the mount point.
* `vswitch_id` - (Required, ForceNew) VSwitch ID, which specifies the VSwitch resource used to create the mount point.
* `vpc_id` - (Required, ForceNew) The ID of the VPC. Specifies the VPC environment to which the mount point belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<file_system_id>:<mount_point_id>`.
* `create_time` - The creation time of the Mount point resource.
* `mount_point_id` - The unique identifier of the Mount point, which is used to retrieve the specified mount point resources.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Mount Point.
* `delete` - (Defaults to 5 mins) Used when delete the Mount Point.
* `update` - (Defaults to 5 mins) Used when update the Mount Point.

## Import

DFS Mount Point can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_mount_point.example <file_system_id>:<mount_point_id>
```