---
subcategory: "Apsara File Storage for HDFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_file_system"
sidebar_current: "docs-alicloud-resource-dfs-file-system"
description: |-
  Provides a Alicloud DFS File System resource.
---

# alicloud\_dfs\_file\_system

Provides a DFS File System resource.

For information about DFS File System and how to use it, see [What is File System](https://www.alibabacloud.com/help/doc-detail/207144.htm).

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAccFileSystem"
}

data "alicloud_dfs_zones" "default" {}

resource "alicloud_dfs_file_system" "default" {
  storage_type     = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
  zone_id          = data.alicloud_dfs_zones.default.zones.0.zone_id
  protocol_type    = "HDFS"
  description      = var.name
  file_system_name = var.name
  throughput_mode  = "Standard"
  space_capacity   = "1024"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the File system.
* `file_system_name` - (Required) The name of the File system.
* `protocol_type` - (Required, ForceNew) The protocol type. Valid values: `HDFS`.
* `provisioned_throughput_in_mi_bps` - (Optional, ForceNew) The preset throughput of the File system. Valid values: `1` to `1024`, Unit: MB/s. **NOTE:** Only when `throughput_mode` is `Provisioned`, this param is valid.
* `space_capacity` - (Required) The capacity budget of the File system. **NOTE:** When the actual data storage reaches the file system capacity budget, the data cannot be written. The file system capacity budget does not support shrinking.
* `storage_type` - (Required, ForceNew) The storage specifications of the File system. Valid values: `PERFORMANCE`, `STANDARD`.
* `throughput_mode` - (Optional, Sensitive) The throughput mode of the File system. Valid values: `Provisioned`, `Standard`.
* `zone_id` - (Required, ForceNew) The zone ID of the File system.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of File System.

## Import

DFS File System can be imported using the id, e.g.

```
$ terraform import alicloud_dfs_file_system.example <id>
```
