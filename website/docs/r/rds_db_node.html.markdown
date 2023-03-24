---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_node"
sidebar_current: "docs-alicloud-resource-rds-db-node"
description: |-
  Provide RDS cluster instance to increase node resources.
---

# alicloud\_rds\_db\_node

Provide RDS cluster instance to increase node resources.
-> **NOTE:** Available in 1.202.0+.

## Example Usage

```
variable "name" {
  default = "tf-testaccrdsdbnodes"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id 		 = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
  zone_id 				   = data.alicloud_db_zones.default.ids.0
  zone_id_slave_a          = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_rds_db_node" "node" {
  db_instance_id = alicloud_db_instance.default.id
  class_code     = alicloud_db_instance.default.instance_type
  zone_id        = alicloud_db_instance.default.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `class_code` - (Required, ForceNew) The specification information of the node.
* `zone_id` - (Required, ForceNew) The zone ID of the node.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of node.The value formats as `<db_instance_id>:<node_id>`.
* `node_id` - The ID of the node.
* `node_role` - The role of node.
* `node_region_id` - The region ID of the node.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Use when opening exclusive agent (until it reaches the initial `Running` status).
* `delete` - (Defaults to 20 mins) Use when closing exclusive agent.

## Import

RDS MySQL database cluster node agent function can be imported using id, e.g.

```shell
$ terraform import alicloud_rds_db_node.example <db_instance_id>:<node_id>
```
