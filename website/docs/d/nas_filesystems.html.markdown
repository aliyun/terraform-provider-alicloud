---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_filesystems"
sidebar_current: "docs-alicloud-datasource-nas-filesystems"
description: |-
    Provides a list of FileSystems owned by an Alibaba Cloud account.
---

# alicloud\_nas_filesystems

This data source provides FileSystems available to the user.

## Example Usage

```
data "alicloud_nas_filesystems" "fs" {
  storage_type = "Performance"
  protocol_type = "NFS"
}

output "first_nas_filesystems_id" {
  value = "${data.alicloud_nas_filesystems.nas_filesystems_ds.filesystems.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Optional) Filter results by a specific StorageType block. 
* `protocol_type` - (Optional) Filter results by a specific ProtocolType block
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `filesystems` - A list of VPCs. Each element contains the following attributes:
  * `id`                    - ID of the FileSystem.
  * `region_id`             - ID of the region where the FileSystem is located.
  * `destription`           - Destription of the FileSystem.
  * `protocol_type`         - ProtocolType block of the FileSystem
  * `storage_type`          - StorageType block of the FileSystem.
  * `metered_size`          - MeteredSize of the FileSystem.
  * `mounttarget_domain`    - List of MountTargetDomain of the FileSystem.
  * `creation_time`         - Time of creation.