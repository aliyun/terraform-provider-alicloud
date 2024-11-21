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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_db_cluster_lake_version&exampleId=44a0d0ac-5708-e0a7-dcbd-889815af140cbc29256a&activeTab=example&spm=docs.r.adb_db_cluster_lake_version.0.44a0d0ac57&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "ap-southeast-1"
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

resource "alicloud_adb_db_cluster_lake_version" "default" {
  db_cluster_version            = "5.0"
  vpc_id                        = data.alicloud_vpcs.default.ids.0
  vswitch_id                    = data.alicloud_vswitches.default.ids.0
  zone_id                       = data.alicloud_adb_zones.default.ids.0
  compute_resource              = "16ACU"
  storage_resource              = "0ACU"
  payment_type                  = "PayAsYouGo"
  enable_default_resource_group = false
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_version` - (Required, ForceNew) The version of the cluster. Valid values: `5.0`.
* `vpc_id` - (Required, ForceNew) The vpc ID of the resource.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `zone_id` - (Required, ForceNew) The zone ID of the resource.
* `compute_resource` - (Required) The computing resources of the cluster.
* `storage_resource` - (Required) The storage resources of the cluster.
* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `PayAsYouGo`.
* `security_ips` - (Optional, Available since v1.198.0) The IP addresses in an IP address whitelist of a cluster. Separate multiple IP addresses with commas (,). You can add a maximum of 500 different IP addresses to a whitelist. The entries in the IP address whitelist must be in one of the following formats:
  - IP addresses, such as 10.23.XX.XX.
  - CIDR blocks, such as 10.23.xx.xx/24. In this example, 24 indicates that the prefix of each IP address in the IP whitelist is 24 bits in length. You can replace 24 with a value within the range of 1 to 32.
* `db_cluster_description` - (Optional, Available since v1.198.0) The description of the cluster.
* `resource_group_id` - (Optional, Available since v1.211.1) The ID of the resource group.
* `enable_default_resource_group` - (Optional, Bool) Whether to enable default allocation of resources to user_default resource groups.
* `source_db_cluster_id` - (Optional, Available since v1.211.1) The ID of the source AnalyticDB for MySQL Data Warehouse Edition cluster.
* `backup_set_id` - (Optional, Available since v1.211.1) The ID of the backup set that you want to use to restore data.
* `restore_type` - (Optional, Available since v1.211.1) The method that you want to use to restore data. Valid values:
  - `backup`: Restores data from a backup set. **NOTE:** You must set `source_db_cluster_id` and `backup_set_id`.
  - `timepoint `: Restores data to a point in time. **NOTE:** You must set `source_db_cluster_id` and `restore_to_time`.
* `restore_to_time` - (Optional, Available since v1.211.1) The point in time to which you want to restore data from the backup set.

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
