---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ots_snapshots"
sidebar_current: "docs-alicloud-datasource-hbr-ots-snapshots"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Ots Snapshots to the user.
---

# alicloud\_hbr\_ots\_snapshots

This data source provides the Hbr Ots Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.164.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_ots_snapshots" "snapshots" {
}
```

## Argument Reference

The following arguments are supported:

* `start_time` - (Optional, ForceNew)  The start time of the backup. This value must be a UNIX timestamp. Unit: milliseconds.
* `end_time` - (Optional, ForceNew)  The end time of the backup. This value must be a UNIX timestamp. Unit: milliseconds
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `snapshots` - A list of Hbr Ots Snapshots. Each element contains the following attributes:
  * `status` - The status of the backup job. Valid values: `COMPLETE`,`PARTIAL_COMPLETE`,`FAILED`.
  * `snapshot_hash` - The hash value of the backup snapshot.
  * `vault_id` - The ID of the backup vault that stores the backup snapshot.
  * `backup_type` - The backup type. Valid value: `COMPLETE`, which indicates full backup.
  * `create_time` - The time when the Table store instance was created. This value is a UNIX timestamp. Unit: seconds.
  * `actual_bytes` - The actual amount of backup snapshots after duplicates are removed. Unit: bytes.
  * `source_type` - The type of the data source. Valid values: `ECS_FILE`,`PARTIAL_COMPLETE`,`FAILED`
  * `bytes_total` - The total amount of data. Unit: bytes.
  * `complete_time` - The time when the backup snapshot was completed. This value is a UNIX timestamp. Unit: seconds.
  * `retention` - The retention period of the backup snapshot.
  * `created_time` - The time when the backup snapshot was created. This value is a UNIX timestamp. Unit: seconds.
  * `parent_snapshot_hash` - The hash value of the parent backup snapshot.
  * `start_time` - The start time of the backup snapshot. This value is a UNIX timestamp. Unit: seconds.
  * `updated_time` - The time when the backup snapshot was updated. This value is a UNIX timestamp. Unit: seconds.
  * `snapshot_id` - The ID of the backup snapshot.
  * `id` - The ID of the backup snapshot.
  * `job_id` - The ID of the backup job.
  * `instance_name` - The name of the Table store instance.
  * `table_name` - The name of the table in the Table store instance.
  * `range_start` - The time when the backup job started. This value is a UNIX timestamp. Unit: milliseconds.
  * `range_end` - The time when the backup job ended. This value is a UNIX timestamp. Unit: milliseconds.

