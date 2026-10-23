---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_backup_policy"
sidebar_current: "docs-alicloud-resource-polardb-backup-policy"
description: |-
  Provides a PolarDB backup policy resource.
---

# alicloud_polardb_backup_policy

Provides a PolarDB cluster backup policy resource and used to configure cluster backup policy.

-> **NOTE:** Available since v1.66.0+. Each PolarDB cluster has a backup policy.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_backup_policy&exampleId=f5565a5b-5785-c4dd-9d94-fc3aba351fc27569cf8d&activeTab=example&spm=docs.r.polardb_backup_policy.0.f5565a5b57&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"
}

resource "alicloud_polardb_backup_policy" "default" {
  db_cluster_id                               = alicloud_polardb_cluster.default.id
  preferred_backup_period                     = ["Tuesday", "Wednesday"]
  preferred_backup_time                       = "10:00Z-11:00Z"
  backup_retention_policy_on_cluster_deletion = "NONE"
}
```

### Advanced backup policies

The following example manages a custom advanced data backup rule. Use a separate cluster from a Normal backup policy example: changing `backup_policy_level` from `Normal` to `Advanced` is irreversible. The cluster must support advanced backup policies (`advanced_policy_option` must not be `notSupport`).

For a same-region rule such as this example, the provider fills omitted `src_region` and `dest_region` with its configured region because the API requires both values. For a cross-region rule (`src_type` or `dest_type` is `level2Cross`), set both regions explicitly.

```terraform
resource "alicloud_polardb_backup_policy" "advanced" {
  db_cluster_id                               = alicloud_polardb_cluster.default.id
  backup_policy_level                         = "Advanced"
  backup_retention_policy_on_cluster_deletion = "NONE"

  advanced_data_policies {
    retention_type  = "delay"
    retention_value = "7"
    filter_type     = "crontab"
    filter_key      = "dayOfWeek"
    filter_value    = "1,2,3,4,5,6,7"
    src_type        = "db"
    dest_type       = "level1"
  }
}
```

`policy_id` is assigned by the service after creation. To manage an existing custom rule after importing the cluster policy, specify that rule's `policy_id` in the `advanced_data_policies` block. Rules generated automatically by the service are exposed through `auto_created_advanced_data_policies` and are never modified or deleted by this resource.

Optional advanced-rule fields returned by the service retain their current values if removed from Terraform configuration. Removing an attribute is not a way to clear it; the API does not document clearing semantics for omitted fields. To remove a custom rule entirely, remove its `advanced_data_policies` block.

In Advanced mode, the legacy data-backup and log-backup arguments cannot be changed through this resource; configure data backup behavior with `advanced_data_policies` instead. Removing this Terraform resource deletes only its managed custom advanced rules. It does not delete the cluster backup policy or restore the `Normal` level.

### Removing alicloud_polardb_cluster from your configuration
 
The `alicloud_polardb_backup_policy` resource manages a cluster policy, but Terraform cannot destroy the policy itself. Removing a Normal policy resource only removes it from Terraform state. Removing an Advanced policy resource deletes its managed custom advanced rules but leaves the cluster policy, its Advanced level, and service-generated rules in place.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_backup_policy&spm=docs.r.polardb_backup_policy.example&intl_lang=EN_US)
 
## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `backup_policy_level` - (Optional, Available since v2.1.1) Backup policy level, `Normal` or `Advanced`. Transitioning from `Normal` to `Advanced` cannot be reversed.
* `advanced_data_policies` - (Optional, Available since v2.1.1) Custom advanced data backup policy blocks. Available only when `backup_policy_level` is `Advanced`. The provider creates, updates, or deletes these rules to match the configuration; service-generated rules are not managed. See [Advanced Data Policies](#advanced-data-policies).
* `preferred_backup_period` - (Optional) PolarDB Cluster backup period. Valid values: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]. Default to ["Tuesday", "Thursday", "Saturday"].
* `preferred_backup_time` - (Optional) PolarDB Cluster backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.
* `backup_retention_policy_on_cluster_deletion` - (Optional, Available in 1.170.0+) Specifies whether to retain backups when you delete a cluster. Valid values are `ALL`, `LATEST`, `NONE`. Default to `NONE`. Value options can refer to the latest docs [ModifyBackupPolicy](https://www.alibabacloud.com/help/en/polardb/latest/modifybackuppolicy)
* `data_level1_backup_retention_period` - (Optional, Available in 1.207.0+) The retention period of level-1 backups. Valid values: 3 to 14. Unit: days.
* `data_level2_backup_retention_period` - (Optional, Available in 1.207.0+) The retention period of level-2 backups. Valid values are `0`, `30 to 7300`, `-1`. Default to `0`.
* `backup_frequency` - (Optional, Available in 1.207.0+) The backup frequency. Valid values are `Normal`, `2/24H`, `3/24H`, `4/24H`.Default to `Normal`.
* `data_level1_backup_frequency` - (Optional, Available in 1.207.0+) The Id of cluster that can run database.The backup frequency. Valid values are `Normal`, `2/24H`, `3/24H`, `4/24H`.Default to `Normal`.
* `data_level1_backup_time` - (Optional, Available in 1.207.0+) The time period during which automatic backup is performed. The format is HH: MMZ HH: MMZ (UTC time), and the entered value must be an hour apart, such as 14:00z-15:00z.
* `data_level1_backup_period` - (Optional, Available in 1.207.0+) PolarDB Cluster of level-1 backup period. Valid values: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"].

  -> **NOTE:** Note Select at least two values. Separate multiple values with commas (,).
* `data_level2_backup_period` - (Optional, Available in 1.207.0+) PolarDB Cluster of level-2 backup period. Valid values: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"].

  -> **NOTE:** Note Select at least two values. Separate multiple values with commas (,).
* `data_level2_backup_another_region_region` - (Optional, Available in 1.207.0+) PolarDB Cluster of level-2 backup is a cross regional backup area.
* `data_level2_backup_another_region_retention_period` - (Optional, Available in 1.207.0+) PolarDB Cluster of level-2 backup cross region backup retention period. Valid values are `0`, `30 to 7300`, `-1`. Default to `0`.
* `log_backup_retention_period` - (Optional, Available in 1.207.0+) The retention period of the log backups. Valid values are `3 to 7300`, `-1`.
* `log_backup_another_region_region` - (Optional, Available in 1.207.0+) The region in which you want to store cross-region log backups. For information about regions that support the cross-region backup feature, see [Overview.](https://www.alibabacloud.com/help/en/polardb/latest/backup-and-restoration-overview)
* `log_backup_another_region_retention_period` - (Optional, Available in 1.207.0+) The retention period of cross-region log backups. Default value: OFF. Valid values are `0`, `30 to 7300`, `-1`.

  -> **NOTE:** Note When you create a cluster, the default value of this parameter is 0.
* `backup_retention_period` - (Optional) Cluster backup retention days, Fixed for 7 days, not modified.

### Advanced Data Policies

The `advanced_data_policies` block supports:

* `policy_id` - (Optional, Computed) Service-assigned rule ID. Specify it to take over and update an existing custom rule after importing the cluster policy.
* `retention_type` - (Required) `never` or `delay`.
* `retention_value` - (Optional, Computed) Retention duration in days; required when `retention_type` is `delay`.
* `filter_type` - (Required) `crontab` or `event`.
* `filter_key` - (Optional, Computed) `dayOfWeek`, `dayOfMonth`, `dayOfYear`, or `backupInterval`; required for `crontab`.
* `filter_value` - (Optional, Computed) Filter expression corresponding to `filter_key`.
* `src_type` - (Required) `db`, `level1`, `level2`, or `level2Cross`.
* `src_region` - (Optional, Computed) Source region. Defaults to the provider region for same-region rules; required explicitly for cross-region rules.
* `dest_type` - (Required) `level1`, `level2`, or `level2Cross`.
* `dest_region` - (Optional, Computed) Destination region. Defaults to the provider region for same-region rules; required explicitly for cross-region rules.
* `dump_action` - (Optional, Computed) `copy`.
* `only_preserve_one_each_day` - (Optional, Computed) Whether to retain only one backup per day. If omitted, the service default is used.
* `storage_class` - (Optional, Computed) Destination storage class.
* `auto_created` - (Computed) Whether the service generated the rule.
* `only_preserve_one_each_hour` - (Computed) Hourly preservation setting supplied by the service.
* `bak_type` - (Computed) Backup type supplied by the service.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'db_cluster_id'.
* `enable_backup_log` - Indicates whether the log backup feature was enabled. Valid values are `0`, `1`. `1` By default, the log backup feature is enabled and cannot be disabled.
* `advanced_policy_option` - Whether advanced policies are supported or enabled by the cluster.
* `auto_created_advanced_data_policies` - Service-generated advanced backup rules. Fields have the same names as `advanced_data_policies`, but are read-only.

## Import

PolarDB backup policy can be imported using the id or cluster id, e.g.

```shell
$ terraform import alicloud_polardb_backup_policy.example "rm-12345678"
```
