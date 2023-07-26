---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mount_target"
sidebar_current: "docs-alicloud-resource-nas-mount-target"
description: |-
  Provides a Alicloud NAS MountTarget resource.
---

# alicloud_nas_mount_target

Provides a NAS Mount Target resource.
For information about NAS Mount Target and how to use it, see [Manage NAS Mount Targets](https://www.alibabacloud.com/help/en/doc-detail/27531.htm).

-> **NOTE:** Available since v1.34.0.

## Example Usage

Basic Usage

```terraform

data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

resource "alicloud_vpc" "main" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "main" {
  vswitch_name = alicloud_vpc.main.vpc_name
  cidr_block   = alicloud_vpc.main.cidr_block
  vpc_id       = alicloud_vpc.main.id
  zone_id      = local.zone_id
}

resource "alicloud_nas_file_system" "example" {
  protocol_type    = "NFS"
  storage_type     = "advance"
  file_system_type = "extreme"
  capacity         = "100"
  zone_id          = local.zone_id
}

resource "alicloud_nas_access_group" "example" {
  access_group_name = "access_group_xxx"
  access_group_type = "Vpc"
  description       = "test_access_group"
  file_system_type  = "extreme"
}

resource "alicloud_nas_mount_target" "example" {
  file_system_id    = alicloud_nas_file_system.example.id
  access_group_name = alicloud_nas_access_group.example.access_group_name
  vswitch_id        = alicloud_vswitch.example.id
  vpc_id            = alicloud_vpc.example.id
  network_type      = alicloud_nas_access_group.example.access_group_type
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `access_group_name` - (Optional) The name of the permission group that applies to the mount target.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch in the VPC where the mount target resides.
* `status` - (Optional) Whether the MountTarget is active. The status of the mount target. Valid values: `Active` and `Inactive`, Default value is `Active`. Before you mount a file system, make sure that the mount target is in the Active state.
* `vpc_id` - (Optional, ForceNew, Available since v1.209.0.) The ID of VPC.
* `network_type` - (Optional, ForceNew, Available since v1.209.0.) mount target network type. Valid values: `VPC`. The classic network's mount targets are not supported.
* `security_group_id` - (Optional, ForceNew, Available in v1.95.0+.) The ID of security group.

## Attributes Reference

The following attributes are exported:

*`id` - This ID of this resource. It is formatted to `<file_system_id>:<mount_target_domain>`. Before version 1.95.0, the value is `<mount_target_domain>`.
* `mount_target_domain` - The IPv4 domain name of the mount target. **NOTE:** Available in v1.161.0+.

## Timeouts

-> **NOTE:** Available in v1.153.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 40 mins) Used when create the mount target (until it reaches the initial `Active` status).
* `update` - (Defaults to 40 mins) Used when update the mount target.
* `delete` - (Defaults to 40 mins) Used when delete the mount target.

## Import

NAS MountTarget can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_mount_target.foo 192094b415:192094b415-luw38.cn-beijing.nas.aliyuncs.com
```
