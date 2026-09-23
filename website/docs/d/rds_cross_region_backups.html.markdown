---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_cross_region_backups"
sidebar_current: "docs-alicloud-datasource-rds-cross-region-backups"
description: |-
  Provides a list of Rds Cross Region Backups to the user.
---

# alicloud\_rds\_cross\_region\_backups

This data source provides the Rds Parameter Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.196.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_cross_region_backups" "backups" {
  db_instance_id = "example_value"
  start_time     = "2022-12-01T00:00:00Z"
  end_time       = "2022-12-16T00:00:00Z"
}

output "first_rds_cross_region_backups" {
  value = data.alicloud_rds_cross_region_backups.backups.backups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The db instance id.
* `cross_backup_region` - (Optional, ForceNew) The ID of the destination region where the cross-region data backup file is stored.

-> **NOTE:** Note You must specify the `cross_backup_id` parameter. Alternatively, you must specify the `start_time` and `end_time` parameters.

* `start_time` - (Optional, ForceNew) The beginning of the time range to query. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.
* `end_time` - (Optional, ForceNew) The end of the time range to query. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.
* `cross_backup_id` - (Optional, ForceNew) The ID of the cross-region data backup file.
* `ids` - (Optional, ForceNew, Computed)  A list of Cross Region Backup IDs.
* `backup_id` - (Optional, ForceNew) The ID of the cross-region data backup file.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Cross Region Backup IDs.
* `backups` - An array that consists of details of the cross-region data backup files:
	* `consistent_time` - The point in time that is indicated by the data in the cross-region data backup file.
	* `db_instance_storage_type` - The storage type.
	* `cross_backup_id` - The ID of the cross-region data backup file.
	* `id` - The ID of the cross-region data backup file.
	* `backup_type` - The type of the cross-region data backup. Valid values:F: full data backup
		`F` - full data backup.
		`I` - incremental data backup.
	* `backup_start_time` - The time when the cross-region data backup started.
	* `backup_end_time` - The time when the cross-region data backup file was generated.
	* `cross_backup_set_location` - The location where the cross-region data backup file is stored.
	* `instance_id` - The ID of the instance. This parameter is used to determine whether the instance that generates the cross-region data backup file is a primary or secondary instance.
	* `cross_backup_download_link` - The external URL from which you can download the cross-region data backup file.
	* `engine_version` - The version of the database engine.
	* `cross_backup_set_file` - The name of the compressed package that contains the cross-region data backup file.
	* `backup_set_scale` - The level at which the cross-region data backup file is generated.
		`0` - instance-level backup.
		`1` - database-level backup.	
	* `cross_backup_set_size` - The size of the cross-region data backup file. Unit: bytes.
	* `backup_set_status` - TThe status of the cross-region data backup. Valid values:
	  	`0` - The cross-region data backup is successful.
	  	`1` - The cross-region data backup failed.	
	* `cross_backup_region` - The ID of the destination region where the cross-region data backup file of the instance is stored.
	* `category` - The RDS edition of the instance. Valid values:
	  	`Basic` - Basic Edition.
	  	`HighAvailability` - High-availability Edition.
	  	`Finance` - Enterprise Edition. This edition is supported only by the China site (aliyun.com).
	* `engine` - The engine of the database.
	* `backup_method` - The method that is used to generate the cross-region data backup file. Valid values:
	  	`L` - logical backup.
	  	`P` - physical backup.	
	* `restore_regions` - An array that consists of the regions to which the cross-region data backup file can be restored.
	* `recovery_begin_time` - The start time to which data can be restored. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
	* `recovery_end_time` - The end time to which data can be restored. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
