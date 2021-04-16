---
subcategory: "Network Attached Storage (NAS)"
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

```terraform
resource "alicloud_nas_file_system" "foo" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "tf-testAccNasConfig"
  encrypt_type = "1"
}
```
## Argument Reference

The following arguments are supported:

* `protocol_type` - (Required, ForceNew) The Protocol Type of a File System. Valid values: `NFS` and `SMB`.
* `storage_type` - (Required, ForceNew) The Storage Type of a File System. Valid values: `Capacity` and `Performance`.
* `description` - (Optional) The File System description.
* `encrypt_type` - (Optional) Whether the file system is encrypted.Using kms service escrow key to encrypt and store the file system data. When reading and writing encrypted data, there is no need to decrypt.
                              Valid values:
                                    0: The file system is not encrypted.
                                    1: The file system is encrypted with a managed secret key.
  
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the File System.

## Import

Nas File System can be imported using the id, e.g.

```
$ terraform import alicloud_nas_file_system.foo 1337849c59
```
