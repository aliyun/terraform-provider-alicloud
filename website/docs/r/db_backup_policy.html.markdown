---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_backup_policy"
sidebar_current: "docs-alicloud-resource-db-backup-policy"
description: |-
  Provides an RDS backup policy resource.
---

# alicloud_db_backup_policy

Provides an RDS instance backup policy resource and used to configure instance backup policy, see [What is DB Backup Policy](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-modifybackuppolicy).


-> **NOTE:** Each DB instance has a backup policy and it will be set default values when destroying the resource.

-> **NOTE:** Available since v1.5.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_db_backup_policy&exampleId=18e9f7ca-02e6-8e35-5036-ddc188df6a4dd42066fb&activeTab=example&spm=docs.r.db_backup_policy.0.18e9f7ca02&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_db_zones" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_backup_policy" "policy" {
  instance_id = alicloud_db_instance.instance.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `backup_period` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'preferred_backup_period' instead.
* `preferred_backup_period` - (Optional, Available since v1.69.0) DB Instance backup period. Please set at least two days to ensure backing up at least twice a week. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday].
* `backup_time` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'preferred_backup_time' instead.
* `preferred_backup_time` - (Optional, Available since v1.69.0) DB instance backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.
* `retention_period` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'backup_retention_period' instead.
* `backup_retention_period` - (Optional, Available since v1.69.0) Instance backup retention days. Valid values: [7-730]. Default to 7. But mysql local disk is unlimited.
* `log_backup` - (Deprecated) It has been deprecated from version 1.68.0, and use field 'enable_backup_log' instead.
* `enable_backup_log` - (Optional, Available since v1.68.0) Whether to backup instance log. Valid values are `true`, `false`, Default to `true`. Note: The 'Basic Edition' category Rds instance does not support setting log backup. [What is Basic Edition](https://www.alibabacloud.com/help/doc-detail/48980.htm).
* `log_retention_period` - (Deprecated) It has been deprecated from version 1.69.0, and use field 'log_backup_retention_period' instead.
* `log_backup_retention_period` - (Optional, Available since v1.69.0) Instance log backup retention days. Valid when the `enable_backup_log` is `1`. Valid values: [7-730]. Default to 7. It cannot be larger than `backup_retention_period`.
* `local_log_retention_hours` - (Optional, Available since v1.69.0) Instance log backup local retention hours. Valid when the `enable_backup_log` is `true`. Valid values: [0-7*24].
* `local_log_retention_space` - (Optional, Available since v1.69.0) Instance log backup local retention space. Valid when the `enable_backup_log` is `true`. Valid values: [0-50].
* `high_space_usage_protection` - (Optional, Available since v1.69.0) Instance high space usage protection policy. Valid when the `enable_backup_log` is `true`. Valid values are `Enable`, `Disable`.
* `log_backup_frequency` - (Optional, Available since v1.69.0) Instance log backup frequency. Valid when the instance engine is `SQLServer`. Valid values are `LogInterval`.
* `compress_type` - (Optional, Available since v1.69.0) The compress type of instance policy. Valid values are `1`, `4`, `8`.
* `archive_backup_retention_period` - (Optional, Available since v1.69.0) Instance archive backup retention days. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. Valid values: [30-1095], and `archive_backup_retention_period` must larger than `backup_retention_period` 730.
* `archive_backup_keep_count` - (Optional, Available since v1.69.0) Instance archive backup keep count. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. When `archive_backup_keep_policy` is `ByMonth` Valid values: [1-31]. When `archive_backup_keep_policy` is `ByWeek` Valid values: [1-7].
* `archive_backup_keep_policy` - (Optional, Available since v1.69.0) Instance archive backup keep policy. Valid when the `enable_backup_log` is `true` and instance is mysql local disk. Valid values are `ByMonth`, `ByWeek`, `KeepAll`.
* `released_keep_policy` - (Optional, Available since v1.147.0) The policy based on which ApsaraDB RDS retains archived backup files if the instance is released. Default value: None. Valid values:
  * **None**: No archived backup files are retained.
  * **Lastest**: Only the most recent archived backup file is retained.
  * **All**: All archived backup files are retained.
* `category` - (Optional, Available since v1.190.0) Whether to enable second level backup.Valid values are `Flash`, `Standard`, Note:It only takes effect when the BackupPolicyMode parameter is DataBackupPolicy.
  -> **NOTE:** You can configure a backup policy by using this parameter and the PreferredBackupPeriod parameter. For example, if you set the PreferredBackupPeriod parameter to Saturday,Sunday and the BackupInterval parameter to -1, a snapshot backup is performed on every Saturday and Sunday.If the instance runs PostgreSQL, the BackupInterval parameter is supported only when the instance is equipped with standard SSDs or enhanced SSDs (ESSDs).This parameter takes effect only when you set the BackupPolicyMode parameter to DataBackupPolicy.
* `backup_interval` - (Optional, Available since v1.194.0) The frequency at which you want to perform a snapshot backup on the instance. Valid values:
  - -1: No backup frequencies are specified.
  - 30: A snapshot backup is performed once every 30 minutes.
  - 60: A snapshot backup is performed once every 60 minutes.
  - 120: A snapshot backup is performed once every 120 minutes.
  - 240: A snapshot backup is performed once every 240 minutes.
  - 360: A snapshot backup is performed once every 360 minutes.
  - 480: A snapshot backup is performed once every 480 minutes.
  - 720: A snapshot backup is performed once every 720 minutes.
* `backup_priority` - (Optional, Int, Available since v1.229.1) Specifies whether the backup settings of a secondary instance are configured. Valid values:
  - 1: secondary instance preferred
  - 2: primary instance preferred
    ->**NOTE:** This parameter is suitable only for instances that run SQL Server on RDS Cluster Edition. This parameter takes effect only when BackupMethod is set to Physical. If BackupMethod is set to Snapshot, backups are forcefully performed on the primary instance that runs SQL Server on RDS Cluster Edition.
* `enable_increment_data_backup` - (Optional, Bool, Available since v1.229.1) Specifies whether to enable incremental backup. Valid values:
  - false (default): disables the feature.
  - true: enables the feature.
    ->**NOTE:** This parameter takes effect only on instances that run SQL Server with cloud disks. This parameter takes effect only when BackupPolicyMode is set to DataBackupPolicy.
* `log_backup_local_retention_number` - (Optional, Int, Available since v1.229.1)  The number of binary log files that you want to retain on the instance. Default value: 60. Valid values: 6 to 100.
  ->**NOTE:** This parameter takes effect only when you set the BackupPolicyMode parameter to LogBackupPolicy. If the instance runs MySQL, you can set this parameter to -1. The value -1 specifies that an unlimited number of binary log files can be retained on the instance.
* `backup_method` - (Optional, Available since v1.229.1)  The backup method of the instance. Valid values:
  - Physical: physical backup
  - Snapshot: snapshot backup
    ->**NOTE:** This parameter takes effect only on instances that run SQL Server with cloud disks. This parameter takes effect only when BackupPolicyMode is set to DataBackupPolicy.

-> **NOTE:** Currently, the SQLServer instance does not support to modify `log_backup_retention_period`.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'instance_id'.

## Import

RDS backup policy can be imported using the id or instance id, e.g.

```shell
$ terraform import alicloud_db_backup_policy.example "rm-12345678"
```