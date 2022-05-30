---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_backup_policy"
sidebar_current: "docs-alicloud-resource-polardb-backup-policy"
description: |-
  Provides a PolarDB backup policy resource.
---

# alicloud\_polardb\_backup\_policy

Provides a PolarDB cluster backup policy resource and used to configure cluster backup policy.

-> **NOTE:** Available in v1.66.0+. Each PolarDB cluster has a backup policy.

## Example Usage

```terraform
variable "name" {
  default = "polardbClusterconfig"
}

variable "creation" {
  default = "PolarDB"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.large"
  pay_type      = "PostPaid"
  description   = var.name
  vswitch_id    = alicloud_vswitch.default.id
}

resource "alicloud_polardb_backup_policy" "policy" {
  db_cluster_id                               = alicloud_polardb_cluster.default.id
  preferred_backup_period                     = ["Tuesday", "Wednesday"]
  preferred_backup_time                       = "10:00Z-11:00Z"
  backup_retention_policy_on_cluster_deletion = "NONE"
}
```
### Removing alicloud_polardb_cluster from your configuration
 
The alicloud_polardb_backup_policy resource allows you to manage your polardb cluster policy, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster policy. You can resume managing the cluster via the polardb Console.
 
## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `preferred_backup_period` - (Optional) PolarDB Cluster backup period. Valid values: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]. Default to ["Tuesday", "Thursday", "Saturday"].
* `preferred_backup_time` - (Optional) PolarDB Cluster backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. Default to "02:00Z-03:00Z". China time is 8 hours behind it.
* `backup_retention_policy_on_cluster_deletion` - (Optional, Available in 1.170.0+) Specifies whether to retain backups when you delete a cluster. Valid values are `ALL`, `LATEST`, `NONE`. Default to `NONE`. Value options can refer to the latest docs [ModifyBackupPolicy](https://help.aliyun.com/document_detail/98103.html)

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'db_cluster_id'.
* `backup_retention_period` - Cluster backup retention days, Fixed for 7 days, not modified.

## Import

PolarDB backup policy can be imported using the id or cluster id, e.g.

```
$ terraform import alicloud_polardb_backup_policy.example "rm-12345678"
```
