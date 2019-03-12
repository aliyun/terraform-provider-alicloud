---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mount_target"
sidebar_current: "docs-alicloud-resource-nas-mount-target"
description: |-
  Provides a Alicloud Nas MountTarget resource.
---

# alicloud\_nas_mount_target

Provides a Nas MountTarget resource.

-> NOTE: Available in v1.34.0+.

## Example Usage

Basic Usage

```
resource "alicloud_nas_mount_target" "foo" {
  file_system_id = "192094b415"
  access_group_name = "tf-testAccNasConfigName"
  vswitch_id = "vsw-13dee3331d"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) File system ID.
* `access_group_name` - (Required) Permission group name.
* `vswitch_id` - (Optional, ForceNew) VSwitch ID.
* `status` - (Optional) Whether the MountTarget is active. An inactive MountTarget is inusable. Valid values are Active(default) and Inactive.

## Attributes Reference

The following attributes are exported:

* `id`  - This ID of this resource. The value is a mount target domain.

## Import

Nas MountTarget  can be imported using the id, e.g.

```
$ terraform import alicloud_nas_mount_target.example 192094b415-luw38.cn-beijing.nas.aliyuncs.com
```
