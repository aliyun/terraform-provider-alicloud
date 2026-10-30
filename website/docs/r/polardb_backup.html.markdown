---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_backup"
description: |-
  Provides a Alicloud Polardb Backup resource.
---

# alicloud_polardb_backup

Provides a Polardb Backup resource.

Cluster-level or database-level backup objects.

For information about Polardb Backup and how to use it, see [What is Backup](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateBackup).

-> **NOTE:** The resource tracks a backup set through the composite ID `<db_cluster_id>:<backup_id>`. If the cluster is deleted while its backups are retained (`backup_retention_policy_on_cluster_deletion` set to `ALL`, or `LATEST` for the backup set retained at deletion time), the retained backup set stays tracked: it is read from the released-cluster backup list (`DescribeDetachedBackups`) instead of being dropped from the Terraform state. While that list is still being populated, the resource keeps its state and reports a warning; the state is only removed once the released-cluster list is complete and does not contain the set, or when the retention policy is `NONE`. Because `DeleteBackup` requires a running cluster, a retained backup set cannot be deleted through the API while the cluster is released; destroying the resource then removes it from the Terraform state and leaves the backup set in the cloud according to the retention policy.

-> **NOTE:** `CreateBackup` returns the ID of the backup task, not the backup set ID. The resource waits for that task to finish and binds the created backup set by the task start time, recording the task ID as `backup_job_id` and the set ID as `backup_id`. If no backup set for that task becomes visible, or more than one manual backup set starts at the task start time, the create fails instead of tracking a backup set created by something else.

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_polardb_backup" "default" {
  db_cluster_id = "pc-bp1593bkhcc4frx36"
}
```

## Argument Reference

The following arguments are supported:
* `advanced_data_policies` - (Optional, List) Advanced data backup policies. The block is only sent when `backup_policy_level` is `Advanced`; under `Normal` it is ignored and the configured entries are kept in the state without being applied. Changing an advanced policy entry updates the backup policy in place and does not force replacement of the backup resource. See [`advanced_data_policies`](#advanced_data_policies) below.
* `backup_frequency` - (Optional, Computed, String) Backup frequency. Valid values: `Normal`, `2/24H`, `3/24H`, `4/24H`. Defaults to `Normal`.

  -> **NOTE:** After enhanced backup is enabled, all backups completed within any 24-hour period are retained. For backups older than 24 hours, the system retains only the first backup completed after 00:00 UTC each day and deletes all others. After enhanced backup is enabled, the `PreferredBackupPeriod` parameter defaults to all days of the week (Monday through Sunday).

* `backup_policy_level` - (Optional, Computed, String) The backup policy level. Valid values: `Normal`, `Advanced`. Defaults to `Normal`. `Advanced` requires the cluster to have advanced backup enabled (申请使用) beforehand; the `advanced_data_policies` block only takes effect under `Advanced` and is ignored under `Normal`.

  -> **NOTE:** Once a cluster has been switched to `Advanced`, the API does not allow switching it back to `Normal`; such a change fails at apply time.
* `backup_retention_policy_on_cluster_deletion` - (Optional, Computed, String) Whether to retain backups when a cluster is deleted. Valid values: `ALL`, `LATEST`, `NONE`. Defaults to `NONE`.

  -> **NOTE:** `ALL` retains all backups permanently, `LATEST` retains only the latest backup permanently, and `NONE` does not retain any backup set.

* `comment` - (Optional, Computed, String) The backup comment. It can only be set when the backup is created; changing it forces a new backup. An explicitly configured empty string is kept as `""` and is not sent as a comment.
* `db_cluster_id` - (Required, String) The cluster ID. Changing it forces a new backup.

  -> **NOTE:** You can call the [DescribeDBClusters](https://help.aliyun.com/document_detail/98094.html) operation to view information about all clusters in the target region, including their cluster IDs.

* `data_level1_backup_frequency` - (Optional, Computed, String) Backup frequency. Valid values: `Normal`, `2/24H`, `3/24H`, `4/24H`. Defaults to `Normal`.

  -> **NOTE:** PolarDB for PostgreSQL (compatible with Oracle) and PolarDB for PostgreSQL do not support this parameter. If cross-region backup is not supported in the region where your PolarDB for MySQL cluster is deployed, this parameter is not supported (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)). After advanced backup is enabled, we recommend that you do not use this parameter; instead, use the `AdvancedDataPolicies` parameter.

* `data_level1_backup_period` - (Optional, Computed, String) The primary backup cycle. Valid values: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.

  -> **NOTE:** You must select at least two days. Separate multiple values with commas (,). This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL. If the region where your PolarDB for MySQL cluster resides does not support cross-region backup, this parameter is unavailable (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)). After advanced backup is enabled, we recommend that you do not use this parameter; use the `AdvancedDataPolicies` parameter instead.

* `data_level1_backup_retention_period` - (Optional, Computed, String) Retention period for level-1 backups, in days. The value ranges from 3 to 14.

  -> **NOTE:** After advanced backup is enabled, this parameter no longer takes effect. Use the `AdvancedDataPolicies` parameter instead.

* `data_level1_backup_time` - (Optional, Computed, String) The time window for automatic backups, in the format `hh:mmZ-hh:mmZ` (UTC). The specified time must be on the hour and span exactly one hour, such as `14:00Z-15:00Z`.

  -> **NOTE:** This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL. If cross-region backup is not supported in the region where your PolarDB for MySQL cluster resides, this parameter is not supported (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)).

* `data_level2_backup_another_region_region` - (Optional, Computed, String) The region for cross-region Level-2 backup. For regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html).

  -> **NOTE:** After advanced backup is enabled, we recommend that you do not use this parameter. Use the `AdvancedDataPolicies` parameter instead.

* `data_level2_backup_another_region_retention_period` - (Optional, Computed, String) Retention period for cross-region Level-2 backup, in days. Valid values: `0`, `-1`, or an integer from 30 to 7300. Defaults to `0`.

  -> **NOTE:** `0` disables the Level-2 backup feature, `30` to `7300` is the retention period in days, and `-1` retains Level-2 backups permanently. After advanced backup is enabled, we recommend that you do not use this parameter; use the `AdvancedDataPolicies` parameter instead.

* `data_level2_backup_period` - (Optional, Computed, String) The secondary backup cycle. Valid values: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.

  -> **NOTE:** You must select at least two days. Separate multiple values with commas (,). This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL. If the region where your PolarDB for MySQL cluster resides does not support cross-region backup, this parameter is unavailable (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)). After advanced backup is enabled, we recommend that you do not use this parameter; use the `AdvancedDataPolicies` parameter instead.

* `data_level2_backup_retention_period` - (Optional, Computed, String) Retention period for secondary backups, in days. Valid values: `0`, `-1`, or an integer from 30 to 7300. Defaults to `0`.

  -> **NOTE:** `0` disables the secondary backup feature, `30` to `7300` is the retention period in days, and `-1` retains secondary backups permanently. After advanced backup is enabled, we recommend that you do not use this parameter; instead, use the `AdvancedDataPolicies` parameter.

* `preferred_backup_period` - (Optional, Computed, String) Data backup cycle. Valid values: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.

  -> **NOTE:** The order of weekdays is not significant. PolarDB normalizes the configured period to a canonical order; the resource treats `preferred_backup_period` as an unordered weekday set, so a reordered value does not produce a persistent plan diff. After advanced backup is enabled, we recommend that you do not use this parameter. Instead, use the `AdvancedDataPolicies` parameter.

* `preferred_backup_time` - (Optional, Computed, String) Time window for performing automatic backups, in UTC. The format is `hh:mmZ-hh:mmZ`. The specified time must be on the hour and span exactly one hour, such as `14:00Z-15:00Z`.

### `advanced_data_policies`

The `advanced_data_policies` block configures Advanced data backup policies. It only takes effect when `backup_policy_level` is set to `Advanced`, which requires the cluster to have advanced backup enabled (申请使用) beforehand; under the `Normal` level the block is ignored and is not sent to the API. Each entry is identified by `policy_id`, and `action_type` declares the operation performed on it. Updating an advanced policy entry modifies the backup policy in place and does not force replacement of the backup resource.

The advanced_data_policies supports the following:
* `action_type` - (Optional, String) Action type. Valid values: `CREATE`, `UPDATE`, `DELETE`. `CREATE` adds a new policy entry, `UPDATE` modifies the entry identified by `policy_id`, and `DELETE` removes it.

  -> **NOTE:** Removing the `advanced_data_policies` block or one of its entries does not delete the cloud policy; the update fails with a clear error instead. Set `action_type` to `DELETE` on the entry that should be removed, apply it, and then remove the entry from the configuration.

* `auto_created` - (Optional, Computed, Bool) Whether the policy is created automatically by the service. The service manages this flag; policies it creates itself are not tracked in the configuration.

  -> **NOTE:** Do not configure `auto_created`. It is reported by the API and only distinguishes service-generated policies from user-defined ones.

* `bak_type` - (Optional, Computed, String) Backup type.
* `dest_region` - (Optional, Computed, String) Destination region.
* `dest_type` - (Optional, Computed, String) Destination type.
* `dump_action` - (Optional, Computed, String) Dump action.
* `filter_key` - (Optional, Computed, String) Filter key.
* `filter_type` - (Optional, Computed, String) Filter type.
* `filter_value` - (Optional, Computed, String) Filter value.
* `only_preserve_one_each_day` - (Optional, Computed, Bool) Whether to keep only one backup each day.
* `only_preserve_one_each_hour` - (Optional, Computed, Bool) Whether to keep only one backup each hour.
* `policy_id` - (Optional, Computed, String) Policy ID. Identifies the advanced policy entry that an `UPDATE` operation targets.
* `retention_type` - (Optional, Computed, String) Retention type. For `delay`, `retention_value` is the retention period in days.
* `retention_value` - (Optional, Computed, String) Retention value. Interpreted together with `retention_type`: for `delay`, the number of days backups are kept.
* `src_region` - (Optional, Computed, String) Source region.
* `src_type` - (Optional, Computed, String) Source type.
* `storage_class` - (Optional, Computed, String) Storage class.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The format is `<db_cluster_id>:<backup_id>`.
* `backup_end_time` - The end time of the backup set.
* `backup_id` - The backup ID, which identifies the backup set in the cluster backup list.
* `backup_job_id` - The ID of the `CreateBackup` task that produced this backup set. `CreateBackup` returns this task ID, while the backup set is listed under its own `backup_id`; the resource binds the set to its task by the task start time. It is only known for backups created by this resource and is empty for imported backups.
* `backup_method` - The data backup method.
* `backup_mode` - Backup mode.
* `backup_progress_status` - The status of the backup task.
* `backup_set_size` - The size of the backup file, in bytes.
* `backup_start_time` - The start time of the backup set.
* `backup_type` - The backup type.
* `backups_level` - Backup level.
* `consistent_time` - The point in time of the consistent snapshot, formatted as a Unix timestamp in seconds.
* `create_time` - Start time of the backup job in UTC.
* `is_avail` - Indicates whether the backup set is available.
* `job_mode` - The backup mode.
* `preferred_next_backup_time` - The scheduled time for the next backup, in the format `YYYY-MM-DDThh:mmZ` (UTC).
* `process` - The progress percentage of the backup task.
* `status` - The backup status.
* `store_status` - Storage status. Only backup sets of a released cluster report this field; it is empty for backup sets of a running cluster.
* `task_action` - Backup task type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 20 mins) Used when create the Backup.
* `delete` - (Defaults to 10 mins) Used when delete the Backup.
* `update` - (Defaults to 5 mins) Used when update the Backup.

## Import

Polardb Backup can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_backup.example <db_cluster_id>:<backup_id>
```
