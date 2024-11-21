---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_db_instance_plan"
sidebar_current: "docs-alicloud-resource-gpdb-db-instance-plan"
description: |-
  Provides a Alicloud AnalyticDB for PostgreSQL (GPDB) DB Instance Plan resource.
---

# alicloud_gpdb_db_instance_plan

Provides a AnalyticDB for PostgreSQL (GPDB) DB Instance Plan resource.

For information about AnalyticDB for PostgreSQL (GPDB) DB Instance Plan and how to use it, see [What is DB Instance Plan](https://www.alibabacloud.com/help/en/analyticdb-for-postgresql/developer-reference/api-gpdb-2016-05-03-createdbinstanceplan).

-> **NOTE:** Available since v1.189.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_db_instance_plan&exampleId=183869b9-bbfd-657e-39b2-7ac473806c889eaa53d0&activeTab=example&spm=docs.r.gpdb_db_instance_plan.0.183869b9bb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_gpdb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids[0]
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
}

resource "time_static" "example" {}

resource "alicloud_gpdb_db_instance_plan" "default" {
  db_instance_plan_name = var.name
  plan_desc             = var.name
  plan_type             = "PauseResume"
  plan_schedule_type    = "Regular"
  plan_start_date       = formatdate("YYYY-MM-DD'T'hh:mm:ss'Z'", timeadd(time_static.example.rfc3339, "1h"))
  plan_end_date         = formatdate("YYYY-MM-DD'T'hh:mm:ss'Z'", timeadd(time_static.example.rfc3339, "24h"))
  plan_config {
    resume {
      plan_cron_time = "0 0 0 1/1 * ? "
    }
    pause {
      plan_cron_time = "0 0 10 1/1 * ? "
    }
  }
  db_instance_id = alicloud_gpdb_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the GPDB instance.
* `db_instance_plan_name` - (Required) The name of the Plan.
* `plan_type` - (Required, ForceNew) The type of the Plan. Valid values: `PauseResume`, `Resize`.
* `plan_schedule_type` - (Required, ForceNew) The execution mode of the plan. Valid values: `Postpone`, `Regular`.
* `plan_start_date` - (Optional) The start time of the Plan.
* `plan_end_date` - (Optional) The end time of the Plan.
* `plan_desc` - (Optional) The description of the Plan.
* `status` - (Optional) The Status of the Plan. Valid values: `active`, `cancel`.
* `plan_config` - (Required, Set) The execution information of the plan. See [`plan_config`](#plan_config) below.

### `plan_config`

The plan_config supports the following:

* `resume` - (Optional, Set) Resume instance plan config. See [`resume`](#plan_config-resume) below.
* `pause` - (Optional, Set) Pause instance plan config. See [`pause`](#plan_config-pause) below.
* `scale_in` - (Optional, Set) Scale In instance plan config. See [`scale_in`](#plan_config-scale_in) below.
* `scale_out` - (Optional, Set) Scale out instance plan config. See [`scale_out`](#plan_config-scale_out) below.
* `scale_up` - (Optional, Set, Available since v1.231.0) Scale up instance plan config. See [`scale_up`](#plan_config-scale_up) below.
* `scale_down` - (Optional, Set, Available since v1.231.0) Scale down instance plan config. See [`scale_down`](#plan_config-scale_down) below.

### `plan_config-resume`

The resume supports the following:

* `execute_time` - (Optional) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

### `plan_config-pause`

The pause supports the following:

* `execute_time` - (Optional) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

### `plan_config-scale_in`

The scale_in supports the following:

* `segment_node_num` - (Optional) The segment Node Num of the Plan.
* `execute_time` - (Optional) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

### `plan_config-scale_out`

The scale_out supports the following:

* `segment_node_num` - (Optional) The segment Node Num of the Plan.
* `execute_time` - (Optional) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

### `plan_config-scale_up`

The scale_up supports the following:

* `instance_spec` - (Optional) The specification of segment nodes of the Plan.
* `execute_time` - (Optional) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

### `plan_config-scale_down`

The scale_down supports the following:

* `instance_spec` - (Optional) The specification of segment nodes of the Plan.
* `execute_time` - (Optional) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DB Instance Plan. It formats as `<db_instance_id>:<plan_id>`.
* `plan_id` - The ID of the plan.
* `plan_config` - The execution information of the plan. Each element contains the following attributes:
  * `pause` - Pause instance plan config.
    * `plan_task_status` - (Available since v1.231.0) The status of the plan task.
  * `resume` - Resume instance plan config.
    * `plan_task_status` - (Available since v1.231.0) The status of the plan task.
  * `scale_out` - Scale out instance plan config.
    * `plan_task_status` - (Available since v1.231.0) The status of the plan task.
  * `scale_in` - Scale in instance plan config.
    * `plan_task_status` - (Available since v1.231.0) The status of the plan task.
  * `scale_up` - Scale up instance plan config.
    * `plan_task_status` - (Available since v1.231.0) The status of the plan task.
  * `scale_down` - Scale down instance plan config.
    * `plan_task_status` - (Available since v1.231.0) The status of the plan task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the DB Instance Plan.
* `update` - (Defaults to 1 mins) Used when update the DB Instance Plan.
* `delete` - (Defaults to 1 mins) Used when delete the DB Instance Plan.

## Import

GPDB DB Instance Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_db_instance_plan.example <db_instance_id>:<plan_id>
```
