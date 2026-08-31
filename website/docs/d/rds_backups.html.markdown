---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_backups"
sidebar_current: "docs-alicloud-datasource-rds-backups"
description: |-
  Provides a list of Rds Backups to the user.
---

# alicloud\_rds\_backups

This data source provides the Rds Backups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_backups" "example" {
  db_instance_id = "example_value"
}

output "first_rds_backup_id" {
  value = data.alicloud_rds_backups.example.backups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `backup_mode` - (Optional, ForceNew) BackupMode. Valid values: `Automated` and `Manual`.
* `db_instance_id` - (Required, ForceNew) The db instance id.
* `end_time` - (Optional, ForceNew) The end time.
* `ids` - (Optional, ForceNew, Computed)  A list of Backup IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `start_time` - (Optional, ForceNew) The start time.
* `backup_status` - (Optional, ForceNew) Backup task status. Valid values: `Automated` and `Manual`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `backups` - A list of Rds Backups. Each element contains the following attributes:
  * `backup_download_url` - The backup download url.
  * `backup_end_time` - BackupEndTime.
  * `backup_id` - BackupId.
  * `backup_initiator` - The initiator of the backup task. Value:
    * **System**: automatically initiated by the System
    * **User**: manually initiated by the User.
  * `backup_intranet_download_url` - The backup intranet download url.
  * `backup_method` - BackupMethod.
  * `backup_mode` - BackupMode.
  * `backup_size` - BackupSize.
  * `backup_start_time` - BackupStartTime.
  * `backup_type` - BackupType.
  * `consistent_time` - The consistency point of the backup set. The return value is a timestamp. **NOTE:** only MySQL 5.6 returns this parameter, and other versions return 0.
  * `copy_only_backup` - The backup mode is divided into the normal backup mode (full and incremental recovery is supported) and the replication-only mode (full recovery is supported only). **NOTE:** Only SQL Server returns this parameter. Valid values:
    * **0**: General Backup Mode
    * **1**: Copy only mode
  * `db_instance_id` - The db instance id.
  * `encryption` - The encrypted information of the backup set.
  * `host_instance_id` - HostInstanceID.
  * `id` - The ID of the Backup.
  * `is_avail` - Whether the backup set is available, the value is:
    * **0**: Not available
    * **1**: Available.
  * `meta_status` - The backup set status of the database table. **NOTE:** an empty string indicates that the backup set for database table recovery is not enabled. Valid values:
    * **OK**: normal.
    * **LARGE**: There are too many tables that cannot be used for database and table recovery.
    * **EMPTY**: The backup set that failed to be backed up.
  * `backup_status` - Backup task status. **NOTE:** This parameter will only be returned when a task is executed. Value:
    * **NoStart**: Not started
    * **Checking**: check the backup
    * **Preparing**: Prepare a backup
    * **Waiting**: Waiting for backup
    * **Uploading**: Upload backup
    * **Finished**: Complete backup
    * **Failed**: backup Failed
  * `storage_class` - The storage medium for the backup set. Valid values:
    * **0**: Regular storage
    * **1**: Archive storage.
  * `store_status` - StoreStatus.