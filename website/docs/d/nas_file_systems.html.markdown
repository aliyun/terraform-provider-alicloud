---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_file_systems"
sidebar_current: "docs-alicloud-datasource-nas-file-systems"
description: |-
  Provides a list of FileSystems owned by an Alibaba Cloud account.
---

# alicloud\_nas_file_systems

This data source provides FileSystems available to the user.

-> NOTE: Available in 1.35.0+

## Example Usage

```
data "alicloud_nas_file_systems" "fs" {
   protocol_type = "NFS"
   description = "${alicloud_nas_file_system.foo.description}"
}

output "first_nas_filesystems_id" {
  value = "${data.alicloud_nas_filesystems.fs.filesystems.0.id}"
}
```
## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of FileSystemId.
* `storage_type` - (Optional) Filter results by a specific StorageType. 
* `protocol_type` - (Optional) Filter results by a specific ProtocolType. 
* `description_regex` - (Optional) A regex string to filter the results by the ：FileSystem description.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of FileSystem Id.
* `systems` - A list of VPCs. Each element contains the following attributes:
  * `id` - ID of the FileSystem.
  * `region_id` - ID of the region where the FileSystem is located.
  * `description` - Destription of the FileSystem.
  * `protocol_type` - ProtocolType block of the FileSystem
  * `storage_type` - StorageType block of the FileSystem.
  * `metered_size` - MeteredSize of the FileSystem.
  * `create_time` - Time of creation.
