---
subcategory: "SelectDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_selectdb_db_clusters"
sidebar_current: "docs-alicloud-datasource-selectdb-db-clusters"
description: |-
  Provides a list of SelectDB DBCluster to the user.
---

# alicloud_selectdb_db_clusters

This data source provides the SelectDB DBCluster of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.229.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

variable "name" {
  default = "terraform_example"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_selectdb_db_instance" "default" {
  db_instance_class       = "selectdb.xlarge"
  db_instance_description = var.name
  cache_size              = 200
  payment_type            = "PayAsYouGo"
  engine_minor_version    = "3.0.12"
  vpc_id                  = data.alicloud_vswitches.default.vswitches.0.vpc_id
  zone_id                 = data.alicloud_vswitches.default.vswitches.0.zone_id
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}
resource "alicloud_selectdb_db_cluster" "default" {
  db_instance_id         = alicloud_selectdb_db_instance.default.id
  db_cluster_description = var.name
  db_cluster_class       = "selectdb.2xlarge"
  cache_size             = 400
  payment_type           = "PayAsYouGo"
}

data "alicloud_selectdb_db_clusters" "default" {
  ids = [alicloud_selectdb_db_cluster.default.id]
}
output "db_cluster" {
  value = data.alicloud_selectdb_db_clusters.default.ids.0
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of DBCluster IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `clusters` - A list of SelectDB DBClusters. Each element contains the following attributes:
  * `id` - The resource ID in terraform of DBCluster. It formats as <db_instance_id>:<db_cluster_id>.
  * `db_cluster_id` - The cluster ID.
  * `db_instance_id` - The instance ID.
  * `db_cluster_class` - The DBCluster class. db_cluster_class has a range of class from `selectdb.xlarge` to `selectdb.256xlarge`.
  * `db_cluster_description` - The DBCluster description.
  * `engine` - The Engine of the DBCluster.
  * `engine_version` - The engine version of the DBCluster.
  * `create_time` - The creation time of the resource.
  * `status` - The status of the DBCluster. Valid values: `ACTIVATION`,`CREATING`,`DELETING`,`RESTARTING`,`ORDER_PREPARING`.
  * `payment_type` - The payment type of the resource. Valid values: `PayAsYouGo`,`Subscription`.
  * `cpu` - The cpu resource amount of DBCluster. Depends on `db_cluster_class`.
  * `memory` - The memory resource amount of DBCluster. Depends on `db_cluster_class`.
  * `cache_size` - The cache size for DBCluster.
  * `region_id` - The ID of region for the cluster.
  * `zone_id` - The ID of zone for the cluster.
  * `vpc_id` - The ID of the VPC for the cluster.
  * `params` - 	The details about each parameter in DBCluster returned.
    * `name` - Parameter name.
    * `value` - The new value of Parameter.
    * `optional` - The value range of the parameter.
    * `comment` - The comments on the parameter.
    * `param_category` - The category of the parameter.
    * `default_value` - The default value of the parameter.
    * `is_dynamic` - Indicates whether the parameter immediately takes effect without requiring a restart.
    * `is_user_modifiable` - Indicates whether the parameter is modifiable.
  * `param_change_logs` - The configuration change logs of parameters.
    * `name` - Changed parameter name.
    * `old_value` - The old value of parameter.
    * `new_value` - The new value of parameter.
    * `gmt_created` - When the parameter change is created.
    * `gmt_modified` - When the parameter change is modified.
    * `config_id` - The id of parameter change.
    * `is_applied` - Whether the parameter changing is applied.
