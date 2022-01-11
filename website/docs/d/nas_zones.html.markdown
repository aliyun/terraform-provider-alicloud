---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_zones"
sidebar_current: "docs-alicloud-datasource-nas-zones"
description: |-
    Provides a list of FileType owned by an Alibaba Cloud account.
---

# alicloud\_nas_zones

Provide  a data source to retrieve the type of zone used to create NAS file system.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

```terraform
data "alicloud_nas_zones" "default" {}

output "alicloud_nas_zones_id" {
  value = "${data.alicloud_nas_zones.default.zones.0.zone_id}"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_type` - (Optional, ForceNew, Available in v1.152.0+) The type of the file system.  Valid values: `standard`, `extreme`, `cpfs`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of availability zone information collection.
    * `zone_id` - String to filter results by zone id.
    * `instance_types` - A list of instance type information collection
        * `storage_type` - The storage type of the nas zones. Valid values:
          * `standard` - When FileSystemType is standard. Valid values: `Performance` and `Capacity`.
          * `extreme` - When FileSystemType is extreme. Valid values: `Standard` and `Advance`.
          * `cpfs` - When FileSystemType is cpfs. Valid values: `advance_100` and `advance_200` .
        * `protocol_type` - File transfer protocol type. Valid values:
          * `standard` - When FileSystemType is standard. Valid values: `NFS` and `SMB`.
          * `extreme` - When FileSystemType is extreme. Valid values: `NFS`.
          * `cpfs` - When FileSystemType is cpfs. Valid values: `cpfs`.
          
