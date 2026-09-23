---
subcategory: "Database Backup(DBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbs_backup_plans"
sidebar_current: "docs-alicloud-datasource-dbs-backup-plans"
description: |-
  Provides a list of Dbs Backup Plans to the user.
---

# alicloud\_dbs\_backup\_plans

This data source provides the Dbs Backup Plans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.185.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dbs_backup_plans" "ids" {}
output "dbs_backup_plan_id_1" {
  value = data.alicloud_dbs_backup_plans.ids.plans.0.id
}

data "alicloud_dbs_backup_plans" "nameRegex" {
  name_regex = "^my-BackupPlan"
}
output "dbs_backup_plan_id_2" {
  value = data.alicloud_dbs_backup_plans.nameRegex.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `backup_plan_name` - (Optional, ForceNew) The name of the resource.
* `enable_details` - (Optional) Default to `true`. Set it to `false` can hide the `payment_type` to output.
* `ids` - (Optional, ForceNew, Computed) A list of Backup Plan IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Backup Plan name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `check_pass`, `init`, `locked`, `pause`, `running`, `stop`, `wait`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Backup Plan names.
* `plans` - A list of Dbs Backup Plans. Each element contains the following attributes:
	* `backup_gateway_id` - The ID of the backup gateway.
	* `backup_method` - The Backup method.
	* `backup_objects` - The backup object.
	* `backup_period` - Full backup cycle.
	* `backup_plan_id` - The first ID of the resource.
	* `backup_plan_name` - The name of the resource.
	* `backup_retention_period` - The retention time of backup data.
	* `backup_start_time` - The start time of full Backup. 
	* `backup_storage_type` - Built-in storage type.
	* `cross_aliyun_id` - The UID that is backed up across Alibaba cloud accounts.
	* `cross_role_name` - The name of the RAM role that is backed up across Alibaba cloud accounts.
	* `database_type` - The database type.
	* `duplication_archive_period` - The storage time for conversion to archive cold standby is 365 days by default.
	* `duplication_infrequent_access_period` - The storage time is converted to low-frequency access. The default time is 180 days.
	* `enable_backup_log` - Whether to enable incremental log Backup.
	* `id` - The ID of the Backup Plan.
	* `instance_class` - The Instance class.
	* `oss_bucket_name` - The OSS Bucket name.
	* `payment_type` - The payment type of the resource.
	* `resource_group_id` - The ID of the resource group.
	* `source_endpoint_database_name` - The name of the database.
	* `source_endpoint_instance_id` - The ID of the database instance.
	* `source_endpoint_instance_type` - The location of the database.
	* `source_endpoint_region` - The region of the database.
	* `source_endpoint_sid` - The Oracle SID name.
	* `source_endpoint_user_name` - The source endpoint username.
	* `status` - The status of the resource.