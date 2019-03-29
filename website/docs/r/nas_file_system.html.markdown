---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_file_system"
sidebar_current: "docs-alicloud-resource-nas-file-system"
description: |-
  Provides a Alicloud NAS File System resource.
---

# alicloud\_nas_file_system

Provides a Nas File System resource.

After activating NAS, you can create a file system and purchase a storage package for it in the NAS console. The NAS console also enables you to view the file system details and remove unnecessary file systems.

For information about NAS file system and how to use it, see [Manage file systems](https://www.alibabacloud.com/help/doc-detail/27530.htm)

-> **NOTE:** Available in v1.33.0+.

## Example Usage

Basic Usage

```
resource "alicloud_nas_file_system" "foo" {
  protocol_type = "NFS"
  storage_type = "Performance"
  description = "tf-testAccNasConfig"
  
}
```
## Argument Reference

The following arguments are supported:

* `protocol_type` - (Required) The Protocol Type of a File System. Valid values: `NFS` and `SMB`.
* `storage_type` - (Required) The Storage Type of a File System. Valid values: `Capacity` and `Performance`.
* `description` - (Optional) The File System description.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the File System.

## Import

Nas File System can be imported using the id, e.g.

```
$ terraform import alicloud_nas_file_system.default 1337849c59
```
