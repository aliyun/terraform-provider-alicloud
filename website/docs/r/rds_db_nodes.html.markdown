---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_nodes"
sidebar_current: "docs-alicloud-resource-rds-db-nodes"
description: |-
  Provide RDS cluster instance to increase node resources.
---

# alicloud\_rds\_db\_nodes

Provide RDS cluster instance to increase node resources.
-> **NOTE:** Available in 1.199.0+.

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

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_db_zones.default.ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
  zone_id 				   = data.alicloud_db_zones.default.ids.0
  zone_id_slave_a          = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_rds_db_nodes" "node" {
  db_instance_id = alicloud_db_instance.default.id
  db_node {
    class_code = alicloud_db_instance.default.instance_type
    zone_id = alicloud_db_instance.default.zone_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `db_node` - (Required, ForceNew) List of cluster nodes.

## Block db_node

The db_node mapping supports the following:

* `class_code` - (Required) The specification information of the node.
* `zone_id` - (Required) The zone ID of the node.
* `node_id` - (Computed)The ID of the node.
* `node_role` - (Computed)The role of node.
* `node_region_id` - (Computed)The region ID of the node.

## Attributes Reference

The following attributes are exported:
* `id` - The Id of DB instance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Use when opening exclusive agent (until it reaches the initial `Running` status).
* `delete` - (Defaults to 20 mins) Use when closing exclusive agent.

## Import

RDS MySQL database cluster node agent function can be imported using id, e.g.

```shell
$ terraform import alicloud_rds_db_nodes.example abc12345678
```
