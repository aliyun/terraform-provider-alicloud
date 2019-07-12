---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mount_target"
sidebar_current: "docs-alicloud-resource-nas-mount-target"
description: |-
  Provides a Alicloud Nas MountTarget resource.
---

# alicloud\_nas_mount_target

Provides a Nas Mount Target resource.

-> NOTE: Available in v1.34.0+.

-> NOTE: Currently this resource support create a mount point in a classic network only when current region is China mainland regions.

-> NOTE: You must grant NAS with specific RAM permissions when creating a classic mount targets,
and it only can be achieved by creating a classic mount target mannually.
See [Add a mount point](https://www.alibabacloud.com/help/doc-detail/60431.htm) and [Why do I need RAM permissions to create a mount point in a classic network](https://www.alibabacloud.com/help/faq-detail/42176.htm).

## Example Usage

Basic Usage

```
resource "alicloud_nas_file_system" "foo" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "tf-testAccNasConfigFs"
}
resource "alicloud_nas_access_group" "foo" {
  name        = "tf-NasConfig-%d"
  type        = "Classic"
  description = "tf-testAccNasConfig"
}
resource "alicloud_nas_access_group" "bar" {
  name        = "tf-cNasConfig-2-%d"
  type        = "Classic"
  description = "tf-testAccNasConfig-2"
}
resource "alicloud_nas_mount_target" "foo" {
  file_system_id    = "${alicloud_nas_file_system.foo.id}"
  access_group_name = "${alicloud_nas_access_group.foo.id}"
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
$ terraform import alicloud_nas_mount_target.foo 192094b415-luw38.cn-beijing.nas.aliyuncs.com
```
