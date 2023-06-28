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
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.4.0.0/24"
  zone_id      = data.alicloud_adb_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "Cluster"
  db_node_class       = "C8"
  db_node_count       = "4"
  db_node_storage     = "400"
  mode                = "reserver"
  db_cluster_version  = "3.0"
  payment_type        = "PayAsYouGo"
  vswitch_id          = alicloud_vswitch.default.id
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  security_ips        = ["10.168.1.12", "10.168.1.11"]
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_adb_backup_policy" "default" {
  db_cluster_id           = alicloud_adb_db_cluster.default.id
  preferred_backup_period = ["Tuesday", "Thursday", "Saturday"]
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
