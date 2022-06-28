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

```terraform
data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nas_zones.default.zones.0.zone_id
}

resource "alicloud_nas_file_system" "foo" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "tf-testacc"
  zone_id          = data.alicloud_nas_zones.default.zones.0.zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `file_system_type` - (Optional, Available in v1.140.0+) the type of the file system. 
                                    Valid values:
                                    `standard` (Default),
                                    `extreme`,
                                    `cpfs`.
* `protocol_type` - (Required, ForceNew) The protocol type of the file system.
                               Valid values:
                                     `NFS`,
                                     `SMB` (Available when the `file_system_type` is `standard`),
                                     `cpfs` (Available when the `file_system_type` is `cpfs`).
* `storage_type` - (Required, ForceNew) The storage type of the file System. 
  * Valid values: 
    * `Performance` (Available when the `file_system_type` is `standard`)
    * `Capacity` (Available when the `file_system_type` is `standard`)
    * `standard` (Available in v1.140.0+ and when the `file_system_type` is `extreme`)
    * `advance` (Available in v1.140.0+ and when the `file_system_type` is `extreme`)
    * `advance_100` (Available in v1.153.0+ and when the `file_system_type` is `cpfs`)
    * `advance_200` (Available in v1.153.0+ and when the `file_system_type` is `cpfs`)
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
* `vpc_id` - (Optional, Available in v1.153.0+) The id of the VPC. The `vpc_id` is required when the `file_system_type` is `cpfs`.
* `vswitch_id` - (Optional, Available in v1.153.0+) The id of the vSwitch. The `vswitch_id` is required when the `file_system_type` is `cpfs`.
* `tags` - (Optional, Available in v1.153.0+) A mapping of tags to assign to the resource.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the File System.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the File System.
* `update` - (Defaults to 10 mins) Used when update the File System.
* `delete` - (Defaults to 10 mins) Used when delete the File System.

## Import

Nas File System can be imported using the id, e.g.

```
$ terraform import alicloud_nas_file_system.foo 1337849c59
```
