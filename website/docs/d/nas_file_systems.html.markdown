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
  protocol_type     = "NFS"
  description_regex = "${alicloud_nas_file_system.foo.description}"
}

output "alicloud_nas_file_systems_id" {
  value = "${data.alicloud_nas_file_systems.fs.systems.0.id}"
}
```
## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of FileSystemId.
* `storage_type` - (Required, ForceNew) The storage type of the file system.
  * Valid values:
    * `Performance` (Available when the `file_system_type` is `standard`)
    * `Capacity` (Available when the `file_system_type` is `standard`)
    * `standard` (Available in v1.140.0+ and when the `file_system_type` is `extreme`)
    * `advance` (Available in v1.140.0+ and when the `file_system_type` is `extreme`)
* `protocol_type` - (Required, ForceNew) The protocol type of the file system.
                                     Valid values:
                                           `NFS`,
                                           `SMB` (Available when the `file_system_type` is `standard`).
* `description_regex` - (Optional) A regex string to filter the results by the ：FileSystem description.
* `file_system_type` - (Optional, Available in v1.140.0+) The type of the file system.
                                      Valid values:
                                            `standard` (Default),
                                            `extreme`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of FileSystem Id.
* `descriptions` - A list of FileSystem descriptions.
* `systems` - A list of VPCs. Each element contains the following attributes:
  * `id` - ID of the FileSystem.
  * `region_id` - ID of the region where the FileSystem is located.
  * `description` - Description of the FileSystem.
  * `protocol_type` - ProtocolType block of the FileSystem
  * `storage_type` - StorageType block of the FileSystem.
  * `metered_size` - MeteredSize of the FileSystem.
  * `create_time` - Time of creation.
  * `capacity` - (Optional, Available in v1.140.0+) The capacity of the file system.
  * `file_system_type` - (Optional, Available in v1.140.0+) The type of the file system.
                            Valid values:
                            `standard` (Default),
                            `extreme`.
  * `encrypt_type` - (Optional, Available in v1.121.2+) Whether the file system is encrypted. 
    * Valid values:
      * `0`: The file system is not encrypted.
      * `1`: The file system is encrypted with a managed secret key.
      * `2`: User management key.
  * `kms_key_id` - (Optional, Available in v1.140.0+) The id of the KMS key.
  * `zone_id` - (Optional, Available in v1.140.0+) The id of the zone. Each region consists of multiple isolated locations known as zones. Each zone has an independent power supply and network.
 
