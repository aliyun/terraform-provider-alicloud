---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_data_backups"
sidebar_current: "docs-alicloud-datasource-gpdb-data-backups"
description: |-
  Provides a list of Gpdb Data Backup owned by an Alibaba Cloud account.
---

# alicloud_gpdb_data_backups

This data source provides Gpdb Data Backup available to the user.[What is Data Backup](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.231.0.

## Example Usage

```terraform
data "alicloud_gpdb_instances" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_gpdb_data_backups" "default" {
  db_instance_id = data.alicloud_gpdb_instances.default.ids.0
}
output "alicloud_gpdb_data_backup_example_id" {
  value = data.alicloud_gpdb_data_backups.default.backups.0.db_instance_id
}
```

## Argument Reference

The following arguments are supported:
* `backup_mode` - (ForceNew, Optional) Backup mode.Full Backup Value Description:-**Automated**: The system is automatically backed up.-**Manual**: Manual backup.Recovery point value description:-**Automated**: The recovery point after a full backup.-**Manual**: The recovery point triggered manually by the user.-**Period**: The recovery point triggered periodically because of the backup policy.
* `db_instance_id` - (Required, ForceNew) The instance ID.
* `data_backup_id` - (ForceNew, Optional) The first ID of the resource
* `data_type` - (ForceNew, Optional) The backup type. Value Description:-**DATA**: Full backup.-**RESTOREPOI**: Recoverable point.
* `end_time` - (Optional, ForceNew) The query end time, which must be greater than the query start time. Format: yyyy-MM-ddTHH:mmZ(UTC time).
* `page_number` - (ForceNew, Optional) Current page number.
* `page_size` - (ForceNew, Optional) Number of records per page.
* `start_time` - (Optional, ForceNew) The query start time. Format: yyyy-MM-ddTHH:mmZ(UTC time).
* `status` - (ForceNew, Optional) Backup set status. Value Description:-Success: The backup has been completed.-Failed: Backup Failed.If not, return all.
* `ids` - (Optional, ForceNew, Computed) A list of Databackup IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Data Backup IDs.
* `backups` - A list of Data Backup Entries. Each element contains the following attributes:
  * `backup_end_time` - The backup end time. Format: yyyy-MM-ddTHH:mm:ssZ(UTC time).
  * `backup_end_time_local` - The end time of the backup (local time).
  * `backup_method` - Backup method. Value Description:-**Physical**: Physical backup.-**Snapshot**: the Snapshot backup.
  * `backup_mode` - Backup mode.Full Backup Value Description:-**Automated**: The system is automatically backed up.-**Manual**: Manual backup.Recovery point value description:-**Automated**: The recovery point after a full backup.-**Manual**: The recovery point triggered manually by the user.-**Period**: The recovery point triggered periodically because of the backup policy.
  * `backup_set_id` - The ID of the backup set.
  * `backup_size` - The size of the backup file. Unit: Byte.
  * `backup_start_time` - The backup start time. Format: yyyy-MM-ddTHH:mm:ssZ(UTC time).
  * `backup_start_time_local` - The start time of the backup (local time).
  * `bakset_name` - The name of the recovery point or full backup set.
  * `consistent_time` - -Full backup: Returns the timestamp of the consistent point in time.-Recoverable point: Returns the timestamp of the recoverable point in time.
  * `db_instance_id` - The instance ID.
  * `data_type` - The backup type. Value Description:-**DATA**: Full backup.-**RESTOREPOI**: Recoverable point.
  * `status` - Backup set status. Value Description:-Success: The backup has been completed.-Failed: Backup Failed.If not, return all.