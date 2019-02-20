---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_filesystem"
sidebar_current: "docs-alicloud-resource-nas_filesystem"
description: |-
  Provides a Alicloud NAS_FileSystem resource.
---

# alicloud\_nas_filesystem

Provides a NAS_FileSystem resource.

~> **NOTE:** Terraform will auto build a filesystem while it uses `alicloud_nas_filesystem` to build a nas_filesystem resource.

## Example Usage

Basic Usage

```
resource "alicloud_nas_filesystem" "foo" {
  protocol_type = "NFS"
  storage_type = "Performance"
  description = "test_wang"
  
}
```
## Argument Reference

The following arguments are supported:

* `protocol_type` - (Required, Forces new resource) The ProtocolType block for the FileSystem.
* `storage_type` - (Required, Forces new resource) The StorageType block for the FileSystem
* `description` - (Optional) The FileSystem description. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `FileSystemId`    - The ID of the FileSystem.
* `storage_type`    - The StorageType block for the FileSystem.
* `protocol_type`   - The ProtocolType block for the FileSystem.
* `description`     - The description of the FileSystem.
* `metered_size`    - The MeteredSize of the FileSystem.
* `create_time`     - The CreateTime of the FileSystem.


