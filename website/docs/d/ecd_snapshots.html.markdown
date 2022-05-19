---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_snapshots"
sidebar_current: "docs-alicloud-datasource-ecd-snapshots"
description: |-
  Provides a list of Ecd Snapshots to the user.
---

# alicloud\_ecd\_snapshots

This data source provides the Ecd Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.169.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_snapshots" "ids" {}
output "ecd_snapshot_id_1" {
  value = data.alicloud_ecd_snapshots.ids.snapshots.0.id
}
```

## Argument Reference

The following arguments are supported:

* `desktop_id` - (Optional, ForceNew) The ID of the Desktop.
* `snapshot_id` - (Optional, ForceNew) The ID of the Snapshot.
* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Snapshot name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Snapshot names.
* `snapshots` - A list of Ecd Snapshots. Each element contains the following attributes:
  * `create_time` - The time when the snapshot was created.
  * `description` - The description of the snapshot.
  * `desktop_id` - The ID of the cloud desktop to which the snapshot belongs.
  * `id` - The ID of the Snapshot.
  * `progress` - The progress of creating the snapshot.
  * `remain_time` - The remaining time that is required to create the snapshot. Unit: seconds.
  * `snapshot_id` - The ID of the snapshot.
  * `snapshot_name` -The name of the snapshot.
  * `snapshot_type` - The type of the snapshot.
  * `source_disk_size` - The capacity of the source disk. Unit: GiB.
  * `source_disk_type` - The type of the source disk.
  * `status` - The status of the snapshot.