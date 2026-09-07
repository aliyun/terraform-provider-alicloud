---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_db_cluster_lake_versions"
sidebar_current: "docs-alicloud-datasource-adb-db-cluster-lake-versions"
description: |-
  Provides a list of Adb DBCluster Lake Versions to the user.
---

# alicloud\_adb\_db\_cluster\_lake\_versions

This data source provides the Adb DBCluster Lake Versions of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.190.0.

## Example Usage

Basic Usage

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
data "alicloud_adb_db_cluster_lake_versions" "ids" {
  ids = [alicloud_adb_db_cluster_lake_version.default.id]
}
output "adb_db_cluster_lake_version_id_1" {
  value = data.alicloud_adb_db_cluster_lake_versions.ids.versions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of DBCluster IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Preparing`, `Creating`, `Restoring`, `Running`, `Deleting`, `ClassChanging`, `NetAddressCreating`, `NetAddressDeleting`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `versions` - A list of Adb Db Clusters. Each element contains the following attributes:
  * `id` - The ID of the DBCluster.
  * `db_cluster_id` - The ID of the DBCluster.
  * `commodity_code` - The name of the service.
  * `compute_resource` - The specifications of computing resources in elastic mode. The increase of resources can speed up queries.
  * `connection_string` - The endpoint of the cluster.
  * `create_time` - The CreateTime of the ADB cluster.
  * `db_cluster_version` - The db cluster version.
  * `engine` - The engine of the database.
  * `engine_version` - The engine version of the database.
  * `expire_time` - The time when the cluster expires.
  * `expired` - Indicates whether the cluster has expired.
  * `lock_mode` - The lock mode of the cluster.
  * `lock_reason` - The reason why the cluster is locked.
  * `payment_type` - The payment type of the resource.
  * `port` - The port that is used to access the cluster.
  * `status` - The status of the resource.
  * `storage_resource` - The specifications of storage resources in elastic mode. The resources are used for data read and write operations.
  * `vpc_id` - The vpc id.
  * `vswitch_id` - The vswitch id.
  * `zone_id` - The zone ID  of the resource.
  * `resource_group_id` - The ID of the resource group.