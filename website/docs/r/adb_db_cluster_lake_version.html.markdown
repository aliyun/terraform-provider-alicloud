---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_db_cluster_lake_version"
sidebar_current: "docs-alicloud-resource-adb-db-cluster-lake-version"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) DB Cluster Lake Version resource.
---

# alicloud_adb_db_cluster_lake_version

Provides a AnalyticDB for MySQL (ADB) DB Cluster Lake Version resource.

For information about AnalyticDB for MySQL (ADB) DB Cluster Lake Version and how to use it, see [What is DB Cluster Lake Version](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/developer-reference/api-adb-2021-12-01-createdbcluster).

-> **NOTE:** Available since v1.190.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
locals {
  zone_id = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.4.0.0/24"
  zone_id      = local.zone_id
  vswitch_name = var.name
}

resource "alicloud_adb_db_cluster_lake_version" "default" {
  compute_resource              = "16ACU"
  db_cluster_version            = "5.0"
  payment_type                  = "PayAsYouGo"
  storage_resource              = "24ACU"
  enable_default_resource_group = false
  vswitch_id                    = alicloud_vswitch.default.id
  vpc_id                        = alicloud_vpc.default.id
  zone_id                       = local.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values are `PayAsYouGo`.
* `compute_resource` - (Required) The computing resources of the cluster.
* `db_cluster_version` - (Required, ForceNew) The version of the cluster. Value options: `5.0`.
* `storage_resource` - (Required) The storage resources of the cluster.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `vpc_id` - (Required, ForceNew) The vpc ID of the resource.
* `zone_id` - (Required, ForceNew) The zone ID of the resource.
* `enable_default_resource_group` - (Optional) Whether to enable default allocation of resources to user_default resource groups.
* `security_ips` - (Optional, Available in v1.198.0+) The IP addresses in an IP address whitelist of a cluster. Separate multiple IP addresses with commas (,). You can add a maximum of 500 different IP addresses to a whitelist. The entries in the IP address whitelist must be in one of the following formats:
  - IP addresses, such as 10.23.XX.XX.
  - CIDR blocks, such as 10.23.xx.xx/24. In this example, 24 indicates that the prefix of each IP address in the IP whitelist is 24 bits in length. You can replace 24 with a value within the range of 1 to 32.
* `db_cluster_description` - (Optional, Available in v1.198.0+) The description of the cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster Lake Version.
* `status` - The status of the resource.
* `commodity_code` - The name of the service.
* `engine` - The engine of the database.
* `engine_version` - The engine version of the database.
* `expire_time` - The time when the cluster expires.
* `expired` - Indicates whether the cluster has expired.
* `lock_mode` - The lock mode of the cluster.
* `lock_reason` - The reason why the cluster is locked.
* `port` - The port that is used to access the cluster.
* `connection_string` - The endpoint of the cluster.
* `create_time` - The createTime of the cluster.
* `resource_group_id` - The ID of the resource group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when create the DB Cluster Lake Version.
* `update` - (Defaults to 72 mins) Used when update the DB Cluster Lake Version.
* `delete` - (Defaults to 5 mins) Used when delete the DB Cluster Lake Version.

## Import

AnalyticDB for MySQL (ADB) DB Cluster Lake Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_db_cluster_lake_version.example <id>
```
