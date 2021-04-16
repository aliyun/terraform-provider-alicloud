---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_file_systems"
sidebar_current: "docs-alicloud-datasource-nas-file-systems"
description: |-
  Provides a list of FileSystems owned by an Alibaba Cloud account.
---

# alicloud\_nas_file_systems

This data source provides FileSystems available to the user.

-> **NOTE**: Available in 1.35.0+

## Example Usage

```terraform
data "alicloud_nas_file_systems" "fs" {
  protocol_type = "NFS"
  description_regex   = "${alicloud_nas_file_system.foo.description}"
}

output "alicloud_nas_file_systems_id" {
  value = "${data.alicloud_nas_file_systems.fs.systems.0.id}"
}
```
## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of FileSystemId.
* `storage_type` - (Optional) Filter results by a specific StorageType. Valid values: `Capacity` and `Performance`.
* `protocol_type` - (Optional) Filter results by a specific ProtocolType. Valid values: `NFS` and `SMB`.
* `description_regex` - (Optional) A regex string to filter the results by the ：FileSystem description.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of FileSystem Id.
* `descriptions` - A list of FileSystem descriptions.
* `systems` - A list of VPCs. Each element contains the following attributes:
  * `id` - ID of the FileSystem.
  * `region_id` - ID of the region where the FileSystem is located.
  * `description` - Destription of the FileSystem.
  * `protocol_type` - ProtocolType block of the FileSystem
  * `storage_type` - StorageType block of the FileSystem.
  * `metered_size` - MeteredSize of the FileSystem.
  * `create_time` - Time of creation.
  * `encrypt_type` - (Optional, Available in v1.121.2+) Whether the file system is encrypted.
                      Valid values:
                      0: The file system is not encrypted.
                      1: The file system is encrypted with a managed secret key.
