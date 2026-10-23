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

-> **NOTE:** Available since v2.1.1.

## Example Usage

Basic Usage

```terraform
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
* `advanced_data_policies` - (Optional, List) Advanced data backup policies. Changing an advanced policy entry updates the backup policy in place and does not force replacement of the backup resource. See [`advanced_data_policies`](#advanced_data_policies) below.
* `backup_frequency` - (Optional, Computed) Backup frequency. Valid values are as follows:
  - `Normal` (default): Standard backup performed once per day at a scheduled time.
  - **2/24H**: Enhanced backup performed every 2 hours.
  - **3/24H**: Enhanced backup performed every 3 hours.
  - **4/24H**: Enhanced backup performed every 4 hours.

  -> **NOTE:** After enhanced backup is enabled, all backups completed within any 24-hour period are retained. For backups older than 24 hours, the system retains only the first backup completed after 00:00 UTC each day and deletes all others. After enhanced backup is enabled, the `PreferredBackupPeriod` parameter defaults to all days of the week (Monday through Sunday).

* `backup_policy_level` - (Optional, Computed) The backup policy level. Valid values: `Normal`, `Advanced`. `Advanced` requires the cluster to have advanced backup enabled (申请使用) beforehand; the `advanced_data_policies` block only takes effect under `Advanced` and is ignored under `Normal`.
* `backup_retention_policy_on_cluster_deletion` - (Optional, Computed) Whether to retain backups when a cluster is deleted:
  - `ALL`: Retains all backups permanently.
  - `LATEST`: Retains only the latest backup permanently.
  - `NONE`: Does not retain any backup sets.

  -> **NOTE:** Default value: NONE.

* `comment` - (Optional, Computed, ForceNew) The backup comment. Forces replacement only when the value is configured (changed).
* `db_cluster_id` - (Required, ForceNew) The cluster ID.

  -> **NOTE:** You can call the [DescribeDBClusters](https://help.aliyun.com/document_detail/98094.html) operation to view information about all clusters in the target region, including their cluster IDs.

* `data_level1_backup_frequency` - (Optional, Computed) Backup frequency. Valid values are as follows:
  - `Normal` (default): Standard backup performed once per day at a scheduled time.
  - **2/24H**: High-frequency backup performed every 2 hours.
  - **3/24H**: High-frequency backup performed every 3 hours.
  - **4/24H**: High-frequency backup performed every 4 hours.

  -> **NOTE:** PolarDB for PostgreSQL (compatible with Oracle) and PolarDB for PostgreSQL do not support this parameter. If cross-region backup is not supported in the region where your PolarDB for MySQL cluster is deployed, this parameter is not supported (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)). After advanced backup is enabled, we recommend that you do not use this parameter; instead, use the `AdvancedDataPolicies` parameter.

* `data_level1_backup_period` - (Optional, Computed) The primary backup cycle. Valid values are as follows:
  - `Monday`: Monday
  - `Tuesday`: Tuesday
  - `Wednesday`: Wednesday
  - `Thursday`: Thursday
  - `Friday`: Friday
  - `Saturday`: Saturday
  - `Sunday`: Sunday

  -> **NOTE:** You must select at least two days. Separate multiple values with commas (,). This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL. If the region where your PolarDB for MySQL cluster resides does not support cross-region backup, this parameter is unavailable (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)). After advanced backup is enabled, we recommend that you do not use this parameter; use the `AdvancedDataPolicies` parameter instead.

* `data_level1_backup_retention_period` - (Optional, Computed) Retention period for level-1 backups, in days. Valid values range from 3 to 14.

  -> **NOTE:** After advanced backup is enabled, this parameter no longer takes effect. Use the `AdvancedDataPolicies` parameter instead.

* `data_level1_backup_time` - (Optional, Computed) The time window for automatic backups, in the format `hh:mmZ-hh:mmZ` (UTC). The specified time must be on the hour and span exactly one hour, such as `14:00Z-15:00Z`.

  -> **NOTE:** This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL. If cross-region backup is not supported in the region where your PolarDB for MySQL cluster resides, this parameter is not supported (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)).

* `data_level2_backup_another_region_region` - (Optional, Computed) The region for cross-region Level-2 backup. For regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html).

  -> **NOTE:** After advanced backup is enabled, we recommend that you do not use this parameter. Use the `AdvancedDataPolicies` parameter instead.

* `data_level2_backup_another_region_retention_period` - (Optional, Computed) Retention period for cross-region Level-2 backup. Valid values are as follows:
  - `0`: Disables the Level-2 backup feature.
  - **30–7300**: Retention period for Level-2 backup, in days.
  - **-1**: Retains Level-2 backups permanently.

  -> **NOTE:** When you create a cluster, the default value is `0`, which disables cross-region Level-2 backup. After advanced backup is enabled, we recommend that you do not use this parameter; use the `AdvancedDataPolicies` parameter instead.

* `data_level2_backup_period` - (Optional, Computed) The secondary backup cycle. Valid values are as follows:
  - `Monday`: Monday
  - `Tuesday`: Tuesday
  - `Wednesday`: Wednesday
  - `Thursday`: Thursday
  - `Friday`: Friday
  - `Saturday`: Saturday
  - `Sunday`: Sunday

  -> **NOTE:** You must select at least two days. Separate multiple values with commas (,). This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL. If the region where your PolarDB for MySQL cluster resides does not support cross-region backup, this parameter is unavailable (for regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html)). After advanced backup is enabled, we recommend that you do not use this parameter; use the `AdvancedDataPolicies` parameter instead.

* `data_level2_backup_retention_period` - (Optional, Computed) Retention period for secondary backups. Valid values are as follows:
  - `0`: Disables the secondary backup feature.
  - **30–7300**: Retention period for secondary backups, in days.
  - **-1**: Retains secondary backups indefinitely.

  -> **NOTE:** When you create a cluster, the default value is `0`, which disables the secondary backup feature. After advanced backup is enabled, we recommend that you do not use this parameter; instead, use the `AdvancedDataPolicies` parameter.

* `preferred_backup_period` - (Optional, Computed) Data backup cycle. Valid values are as follows:
  - Monday: Monday
  - Tuesday: Tuesday
  - Wednesday: Wednesday
  - Thursday: Thursday
  - Friday: Friday
  - Saturday: Saturday
  - Sunday: Sunday

  -> **NOTE:** The order of weekdays is not significant. PolarDB normalizes the configured period to a canonical order; the resource treats `preferred_backup_period` as an unordered weekday set, so a reordered value does not produce a persistent plan diff. After advanced backup is enabled, we recommend that you do not use this parameter. Instead, use the `AdvancedDataPolicies` parameter.

* `preferred_backup_time` - (Optional, Computed) Time window for performing automatic backups, in UTC. The format is `hh:mmZ-hh:mmZ`. The specified time must be on the hour and span exactly one hour, such as `14:00Z-15:00Z`.

### `advanced_data_policies`

The `advanced_data_policies` block configures Advanced data backup policies. It only takes effect when `backup_policy_level` is set to `Advanced`, which requires the cluster to have advanced backup enabled (申请使用) beforehand; under the `Normal` level the block is ignored.

-> **NOTE:** Each entry is identified by `policy_id`. Use `action_type` to declare the operation performed on the entry (`CREATE` to add a new policy, `UPDATE` to modify an existing one, `DELETE` to remove it). `retention_value` is interpreted together with `retention_type` (for example, `delay` retains backups for the number of days given by `retention_value`). Updating an advanced policy entry modifies the backup policy in place and does not force replacement of the backup resource.

The advanced_data_policies supports the following:
* `action_type` - (Optional) Action type. Valid values: `CREATE`, `UPDATE`, `DELETE`. `CREATE` adds a new policy entry, `UPDATE` modifies the entry identified by `policy_id`, and `DELETE` removes it.
* `auto_created` - (Optional, Computed) Auto created.
* `bak_type` - (Optional, Computed) Backup type.
* `dest_region` - (Optional, Computed) Destination region.
* `dest_type` - (Optional, Computed) Destination type.
* `dump_action` - (Optional, Computed) Dump action.
* `filter_key` - (Optional, Computed) Filter key.
* `filter_type` - (Optional, Computed) Filter type.
* `filter_value` - (Optional, Computed) Filter value.
* `only_preserve_one_each_day` - (Optional, Computed, Bool) Only preserve one each day.
* `only_preserve_one_each_hour` - (Optional, Computed, Bool) Only preserve one each hour.
* `policy_id` - (Optional, Computed) Policy ID. Identifies the advanced policy entry that an `UPDATE` operation targets.
* `retention_type` - (Optional, Computed) Retention type.
* `retention_value` - (Optional, Computed) Retention value. Interpreted together with `retention_type` to determine the retention period and unit (for `delay`, the value is the retention period in days).
* `src_region` - (Optional, Computed) Source region.
* `src_type` - (Optional, Computed) Source type.
* `storage_class` - (Optional, Computed) Storage class.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The format is `<db_cluster_id>:<backup_id>`.
* `backup_end_time` - The end time of the query.
* `backup_id` - The backup ID.
* `backup_job_id` - The backup job ID.
* `backup_method` - The data backup method.
* `backup_mode` - Backup mode.
* `backup_progress_status` - The status of the backup task.
* `backup_set_size` - The size of the backup file, in bytes.
* `backup_start_time` - The start time of the query.
* `backup_type` - The backup type.
* `backups_level` - Backup level.
* `consistent_time` - The point in time of the consistent snapshot, formatted as a Unix timestamp in seconds.
* `create_time` - Start time of the backup job in UTC.
* `is_avail` - Indicates whether the backup set is available.
* `job_mode` - The backup mode.
* `preferred_next_backup_time` - The scheduled time for the next backup, in the format `YYYY-MM-DDThh:mmZ` (UTC).
* `process` - The progress percentage of the backup task.
* `status` - The backup status.
* `store_status` - Storage status.
* `task_action` - Backup task type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Backup.
* `delete` - (Defaults to 5 mins) Used when delete the Backup.
* `update` - (Defaults to 5 mins) Used when update the Backup.

## Import

Polardb Backup can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_backup.example <db_cluster_id>:<backup_id>
```
