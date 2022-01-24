---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_snapshots"
sidebar_current: "docs-alicloud-datasource-dbfs-snapshots"
description: |-
  Provides a list of Dbfs Snapshots to the user.
---

# alicloud\_dbfs\_snapshots

This data source provides the Dbfs Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.156.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dbfs_snapshots" "ids" {
  ids = ["example_id"]
}
output "dbfs_snapshot_id_1" {
  value = data.alicloud_dbfs_snapshots.ids.snapshots.0.id
}

data "alicloud_dbfs_snapshots" "nameRegex" {
  name_regex = "^my-Snapshot"
}
output "dbfs_snapshot_id_2" {
  value = data.alicloud_dbfs_snapshots.nameRegex.snapshots.0.id
}

data "alicloud_dbfs_snapshots" "status" {
  status = "accomplished"
}
output "dbfs_snapshot_id_3" {
  value = data.alicloud_dbfs_snapshots.status.snapshots.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Snapshot name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the snapshot. Valid values: `accomplished`, `failed`, `progressing`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Snapshot names.
* `snapshots` - A list of Dbfs Snapshots. Each element contains the following attributes:
	* `category` - The type of the Snapshot.
	* `create_time` - The creation time of the snapshot.
	* `description` - The description of the snapshot.
	* `id` - The ID of the Snapshot.
	* `last_modified_time` - The last modification time of the snapshot.
	* `progress` - The progress of the snapshot.
	* `remain_time` - The remaining completion time of the snapshot being created, in seconds.
	* `retention_days` - The retention days of the snapshot.
	* `snapshot_id` - The ID of the snapshot.
	* `snapshot_name` - The name of the snapshot.
	* `snapshot_type` - The creation of the snapshot.
	* `instance_id` - The ID of the database file system.
	* `source_fs_size` - Source database file system capacity.
	* `status` - The status of the snapshot. Possible values: `progressing`, `accomplished`, `failed`.