---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_snapshots"
sidebar_current: "docs-alicloud-datasource-nas-snapshots"
description: |-
  Provides a list of Nas Snapshots to the user.
---

# alicloud\_nas\_snapshots

This data source provides the Nas Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nas_snapshots" "ids" {}
output "nas_snapshot_id_1" {
  value = data.alicloud_nas_snapshots.ids.snapshots.0.id
}

data "alicloud_nas_snapshots" "nameRegex" {
  name_regex = "^my-Snapshot"
}
output "nas_snapshot_id_2" {
  value = data.alicloud_nas_snapshots.nameRegex.snapshots.0.id
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Optional, ForceNew) The ID of the file system.
* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Snapshot name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `snapshot_name` - (Optional, ForceNew) The name of the snapshot.
* `status` - (Optional, ForceNew) Status. Valid values: `accomplished`, `failed`, `progressing`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Snapshot names.
* `snapshots` - A list of Nas Snapshots. Each element contains the following attributes:
	* `create_time` - The creation time of the resource.
	* `description` - The description of the snapshot.
	* `encrypt_type` - The type of the encryption.
	* `id` - The ID of the Snapshot.
	* `progress` - The progress of the snapshot creation. The value of this parameter is expressed as a percentage.
	* `remain_time` - The remaining time that is required to create the snapshot. Unit: seconds.
	* `retention_days` - The retention period of the automatic snapshot. Unit: days.
	* `snapshot_id` - The ID of the resource.
	* `snapshot_name` - The name of the snapshot.
	* `source_file_system_id` - The ID of the source file system.
	* `source_file_system_size` - The capacity of the source file system. Unit: GiB.
	* `source_file_system_version` - The version of the source file system.
	* `status` - The status of the snapshot.