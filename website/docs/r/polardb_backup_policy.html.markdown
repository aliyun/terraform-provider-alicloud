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
### Removing alicloud_polardb_cluster from your configuration
 
The alicloud_polardb_backup_policy resource allows you to manage your polardb cluster policy, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster policy. You can resume managing the cluster via the polardb Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_backup_policy&spm=docs.r.polardb_backup_policy.example&intl_lang=EN_US)
 
## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
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

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'db_cluster_id'.
* `enable_backup_log` - Indicates whether the log backup feature was enabled. Valid values are `0`, `1`. `1` By default, the log backup feature is enabled and cannot be disabled.

## Import

PolarDB backup policy can be imported using the id or cluster id, e.g.

```shell
$ terraform import alicloud_polardb_backup_policy.example "rm-12345678"
```
