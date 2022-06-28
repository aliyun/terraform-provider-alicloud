---
subcategory: "Cassandra"
layout: "alicloud"
page_title: "Alicloud: alicloud_cassandra_backup_plan"
sidebar_current: "docs-alicloud-resource-cassandra-backup-plan"
description: |-
  Provides a Alicloud Cassandra Backup Plan resource.
---

# alicloud\_cassandra\_backup\_plan

Provides a Cassandra Backup Plan resource.

For information about Cassandra Backup Plan and how to use it, see [What is Backup Plan](https://www.alibabacloud.com/help/doc-detail/157522.htm).

-> **NOTE:** Available in v1.128.0+.

## Example Usage

Basic Usage

```terraform

variable "name" {
  default = "tf-testAccCassandrBackupPlan"
}

data "alicloud_cassandra_zones" "example" {
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  vpc_id       = alicloud_cassandra_cluster.example.id
  zone_id      = data.alicloud_cassandra_zones.example.zones[length(data.alicloud_cassandra_zones.example.ids) + (-1)].id
  cidr_block   = cidrsubnet(alicloud_vpc.example.vpcs.0.cidr_block, 8, 4)
}

resource "alicloud_cassandra_cluster" "example" {
  cluster_name        = var.name
  data_center_name    = var.name
  auto_renew          = "false"
  instance_type       = "cassandra.c.large"
  major_version       = "3.11"
  node_count          = "2"
  pay_type            = "PayAsYouGo"
  vswitch_id          = alicloud_vswitch.example[0].id
  disk_size           = "160"
  disk_type           = "cloud_ssd"
  maintain_start_time = "18:00Z"
  maintain_end_time   = "20:00Z"
  ip_white            = "127.0.0.1"
}

resource "alicloud_cassandra_backup_plan" "example" {
  cluster_id     = alicloud_cassandra_cluster.example.id
  data_center_id = alicloud_cassandra_cluster.example.zone_id
  backup_time    = "00:10Z"
  active         = false

}

```

## Argument Reference

The following arguments are supported:

* `active` - (Optional, Computed) Specifies whether to activate the backup plan. Valid values: `True`, `False`. Default value: `True`.
* `backup_period` - (Optional, Computed) The backup cycle. Valid values: `Friday`, `Monday`, `Saturday`, `Sunday`, `Thursday`, `Tuesday`, `Wednesday`.
* `backup_time` - (Required) The start time of the backup task each day. The time is displayed in UTC and denoted by Z.
* `cluster_id` - (Required, ForceNew) The ID of the cluster for the backup.
* `data_center_id` - (Required, ForceNew) The ID of the data center for the backup in the cluster.
* `retention_period` - (Optional, Computed) The duration for which you want to retain the backup. Valid values: 1 to 30. Unit: days. Default value: `30`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Backup Plan. The value is formatted `<cluster_id>:<data_center_id>`.

## Import

Cassandra Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_cassandra_backup_plan.example <cluster_id>:<data_center_id>
```
