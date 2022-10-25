---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_db_instance_plan"
sidebar_current: "docs-alicloud-resource-gpdb-db-instance-plan"
description: |-
  Provides a Alicloud GPDB DB Instance Plan resource.
---

# alicloud\_gpdb\_db\_instance\_plan

Provides a GPDB DB Instance Plan resource.

For information about GPDB DB Instance Plan and how to use it, see [What is DB Instance Plan](https://www.alibabacloud.com/help/zh/analyticdb-for-postgresql/latest/createdbinstanceplan).

-> **NOTE:** Available in v1.189.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_gpdb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.ids.0
  vswitch_name = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "Serverless"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "4C16G"
  master_node_num       = 1
  payment_type          = "PayAsYouGo"
  private_ip_address    = "1.1.1.1"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = local.vswitch_id
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
resource "alicloud_gpdb_db_instance_plan" "default" {
  db_instance_plan_name = "${var.name}"
  plan_desc             = "${var.name}"
  plan_type             = "PauseResume"
  plan_schedule_type    = "Regular"
  plan_start_date       = "2022-10-16T00:00:00Z"
  plan_end_date         = "2023-10-16T00:00:00Z"
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

* `db_instance_id` - (Required, ForceNew) The ID of the Database instance.
* `db_instance_plan_name` - (Required) The name of the Plan.
* `plan_desc` - (Optional, Computed) The description of the Plan.
* `plan_end_date` - (Optional, Computed) The end time of the Plan.
* `plan_schedule_type` - (Required, ForceNew) Plan scheduling type. Valid values: `Postpone`, `Regular`.
* `plan_start_date` - (Optional, Computed) The start time of the Plan.
* `plan_type` - (Required, ForceNew) The type of the Plan. Valid values: `PauseResume`, `Resize`.
* `status` - (Optional, Computed) The Status of the Plan. Valid values: `active`, `cancel`.
* `plan_config` - (Required) The plan config. See the following `Block plan_config`.

#### Block plan_config

The plan_config supports the following:

* `pause` - (Optional) Pause instance plan config. See the following `Block pause`.
* `resume` - (Optional) Resume instance plan config. See the following `Block resume`.
* `scale_in` - (Optional) Scale In instance plan config. See the following `Block scale_in`.
* `scale_out` - (Optional) Scale out instance plan config. See the following `Block scale_out`.

#### Block pause

The pause supports the following:

* `execute_time` - (Optional, Computed) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

#### Block resume

The resume supports the following:

* `execute_time` - (Optional, Computed) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.

#### Block scale_in

The scale_in supports the following:

* `execute_time` - (Optional, Computed) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.
* `segment_node_num` - (Optional) The segment Node Num of the Plan.

#### Block scale_out

The scale_out supports the following:

* `execute_time` - (Optional, Computed) The executed time of the Plan.
* `plan_cron_time` - (Optional) The Cron Time of the plan.
* `segment_node_num` - (Optional) The segment Node Num of the Plan.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of DB Instance Plan. The value formats as `<db_instance_id>:<plan_id>`.
* `plan_id` - The ID of DB Instance Plan.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the DB Instance Plan.
* `delete` - (Defaults to 1 mins) Used when delete the DB Instance Plan.
* `update` - (Defaults to 1 mins) Used when update the DB Instance Plan.


## Import

GPDB DB Instance Plan can be imported using the id, e.g.

```
$ terraform import alicloud_gpdb_db_instance_plan.example <db_instance_id>:<plan_id>
```