---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_backup_policy"
sidebar_current: "docs-alicloud-resource-adb-backup-policy"
description: |-
  Provides a ADB backup policy resource.
---

# alicloud_adb_backup_policy

Provides a [ADB](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/api-doc-adb-2019-03-15-api-doc-modifybackuppolicy) cluster backup policy resource and used to configure cluster backup policy.

-> **NOTE:** Available since v1.71.0.

-> Each DB cluster has a backup policy and it will be set default values when destroying the resource.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_adb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}


locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_adb_db_cluster" "cluster" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "8Core32GB"
  vswitch_id          = local.vswitch_id
  description         = var.name
}

resource "alicloud_adb_backup_policy" "default" {
  db_cluster_id           = alicloud_adb_db_cluster.cluster.id
  preferred_backup_period = ["Tuesday", "Wednesday"]
  preferred_backup_time   = "10:00Z-11:00Z"
}
```
### Removing alicloud_adb_cluster from your configuration
 
The alicloud_adb_backup_policy resource allows you to manage your adb cluster policy, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster policy. You can resume managing the cluster via the adb Console.
 
## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `preferred_backup_period` - (Required) ADB Cluster backup period. Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday].
* `preferred_backup_time` - (Required) ADB Cluster backup time, in the format of HH:mmZ- HH:mmZ. Time setting interval is one hour. China time is 8 hours behind it.

## Attributes Reference

The following attributes are exported:

* `id` - The current backup policy resource ID. It is same as 'db_cluster_id'.
* `backup_retention_period` - Cluster backup retention days, Fixed for 7 days, not modified.

## Import

ADB backup policy can be imported using the id or cluster id, e.g.

```shell
$ terraform import alicloud_adb_backup_policy.example "am-12345678"
```
