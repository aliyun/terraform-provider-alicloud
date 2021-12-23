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
  encrypt_type  = "1"
}
```

```terraform
resource "alicloud_nas_file_system" "foo" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = "cn-hangzhou-f"
  storage_type     = "standard"
  description      = "tf-testAccNasConfig"
  capacity         = "100"
}
```

## Argument Reference

The following arguments are supported:
* `file_system_type` - (Optional, Available in v1.140.0+) the type of the file system. 
                                    Valid values:
                                    `standard` (Default),
                                    `extreme`.
* `protocol_type` - (Required, ForceNew) The protocol type of the file system.
                               Valid values:
                                     `NFS`,
                                     `SMB` (Available when the `file_system_type` is `standard`).
* `storage_type` - (Required, ForceNew) The storage type of the file System. 
  * Valid values: 
    * `Performance` (Available when the `file_system_type` is `standard`)
    * `Capacity` (Available when the `file_system_type` is `standard`)
    * `standard` (Available in v1.140.0+ and when the `file_system_type` is `extreme`)
    * `advance` (Available in v1.140.0+ and when the `file_system_type` is `extreme`)
* `description` - (Optional) The File System description.
* `encrypt_type` - (Optional, Available in v1.121.2+) Whether the file system is encrypted. Using kms service escrow key to encrypt and store the file system data. When reading and writing encrypted data, there is no need to decrypt. 
  * Valid values:
    * `0` (Default): The file system is not encrypted. 
    * `1`: The file system is encrypted with a managed secret key.
    * `2` (Available in v1.140.0+ and when the `file_system_type` is `extreme`): User management key.
* `capacity` - (Optional, Available in v1.140.0+ and when the `file_system_type` is `extreme`) The capacity of the file system. The `capacity` is required when the `file_system_type` is `extreme`.
                            Unit: gib; **Note**: The minimum value is 100.
* `zone_id` - (Optional, Available in v1.140.0+) The available zones information that supports nas.When FileSystemType=standard, this parameter is not required. **Note:** By default, a qualified availability zone is randomly selected according to the `protocol_type` and `storage_type` configuration.
* `kms_key_id` - (Optional, Available in v1.140.0+ and when the `encrypt_type` is `2`) The id of the KMS key. The `kms_key_id` is required when the `encrypt_type` is `2`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the File System.

## Import

Nas File System can be imported using the id, e.g.

```
$ terraform import alicloud_nas_file_system.foo 1337849c59
```
