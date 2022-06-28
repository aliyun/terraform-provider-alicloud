---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_snapshots"
sidebar_current: "docs-alicloud-datasource-simple-application-server-snapshots"
description: |-
  Provides a list of Simple Application Server Snapshots to the user.
---

# alicloud\_simple\_application\_server\_snapshots

This data source provides the Simple Application Server Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_snapshots" "ids" {
  ids = ["example_id"]
}
output "simple_application_server_snapshot_id_1" {
  value = data.alicloud_simple_application_server_snapshots.ids.snapshots.0.id
}

data "alicloud_simple_application_server_snapshots" "nameRegex" {
  name_regex = "^my-Snapshot"
}
output "simple_application_server_snapshot_id_2" {
  value = data.alicloud_simple_application_server_snapshots.nameRegex.snapshots.0.id
}

data "alicloud_simple_application_server_snapshots" "diskIdConf" {
  ids     = ["example_id"]
  disk_id = "example_value"
}
output "simple_application_server_snapshot_id_3" {
  value = data.alicloud_simple_application_server_snapshots.diskIdConf.snapshots.0.id
}

data "alicloud_simple_application_server_snapshots" "instanceIdConf" {
  ids         = ["example_id"]
  instance_id = "example_value"
}
output "simple_application_server_snapshot_id_4" {
  value = data.alicloud_simple_application_server_snapshots.instanceIdConf.snapshots.0.id
}

```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Optional, ForceNew) The ID of the disk.
* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot IDs.
* `instance_id` - (Optional, ForceNew) The ID of the simple application server.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Snapshot name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the snapshots. Valid values: `Progressing`, `Accomplished` and `Failed`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Snapshot names.
* `snapshots` - A list of Simple Application Server Snapshots. Each element contains the following attributes:
	* `create_time` - The time when the snapshot was created. The time follows the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time is displayed in UTC.
	* `disk_id` - The ID of the source disk. This parameter has a value even after the source disk is released.
	* `id` - The ID of the Snapshot.
	* `progress` - The progress of snapshot creation.
	* `remark` - The remarks of the snapshot.
	* `snapshot_id` - The ID of the snapshot.
	* `snapshot_name` - The name of the snapshot.
	* `source_disk_type` - A snapshot of the source of a disk type. Possible values: `System`, `Data`.
	* `status` - The status of the snapshots. Valid values: `Progressing`, `Accomplished` and `Failed`.