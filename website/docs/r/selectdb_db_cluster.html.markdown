---
subcategory: "SelectDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_selectdb_db_cluster"
sidebar_current: "docs-alicloud-resource-selectdb-db-cluster"
description: |-
  Provides a Alicloud SelectDB DBCluster resource.
---

# alicloud_selectdb_db_cluster

Provides a SelectDB DBCluster resource.

For information about SelectDB DBCluster and how to use it, see [What is DBCluster](https://www.alibabacloud.com/help/zh/selectdb/latest/api-selectdb-2023-05-22-createdbcluster).

-> **NOTE:** Available since v1.229.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_selectdb_db_cluster&exampleId=f90da21a-f1e4-6b65-bfd8-cc9d2d4c9d04e84407e0&activeTab=example&spm=docs.r.selectdb_db_cluster.0.f90da21af1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The InstanceId of DBInstance for DBCluster. Every DBCluster requires one DBInstance to rely on.
* `payment_type` - (Required, ForceNew) The payment type of the resource. Valid values: `PayAsYouGo`,`Subscription`.
* `db_cluster_class` - (Required) The DBCluster class. db_cluster_class has a range of class from `selectdb.xlarge` to `selectdb.256xlarge`.
* `cache_size` - (Required) The desired cache size on creating cluster. The number should be divided by 100.
* `db_cluster_description` - (Required) The DBCluster description.
* `desired_status` - (Optional) The desired status for the resource. Valid values: `ACTIVATION`,`STOPPED`,`STARTING`,`RESTART`.
* `desired_params` - (Optional) The modified parameter in DBCluster. See [`desired_params`](#desired_params) below.

### `desired_params`

The desired_params supports the following:

* `name` - (Optional) Parameter name.
* `value` - (Optional) The new value of Parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster. It formats as <db_instance_id>:<db_cluster_id>.
* `db_cluster_id` - The id of the cluster. 
* `engine` - The engine of DBCluster. Always `selectdb`.
* `engine_version` - The version of DBCluster. 
* `create_time` - The time when DBCluster is created.
* `status` - The current status of the resource.
* `cpu` - The cpu resource amount of DBCluster. Depends on `db_cluster_class`.
* `memory` - The memory resource amount of DBCluster. Depends on `db_cluster_class`.
* `cache_size` - The cache size of DBCluster.
* `region_id` - The ID of region for the cluster.
* `zone_id` - The ID of zone for the cluster.
* `vpc_id` - The ID of the VPC for the cluster.
* `param_change_logs` - The details about parameter changelogs in DBCluster returned.
  * `name` - Changed parameter name.
  * `old_value` - The old value of parameter.
  * `new_value` - The new value of parameter.
  * `gmt_created` - When the parameter change is created.
  * `gmt_modified` - When the parameter change is modified.
  * `config_id` - The id of parameter change.
  * `is_applied` - Whether the parameter changing is applied.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the SelectDB DBCluster (until it reaches the initial `ACTIVATION` status).
* `update` - (Defaults to 30 mins) Used when update the SelectDB DBCluster.
* `delete` - (Defaults to 30 mins) Used when delete the SelectDB DBCluster.

## Import

SelectDB DBCluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_selectdb_db_cluster.example <db_instance_id>:<db_cluster_id>
```
