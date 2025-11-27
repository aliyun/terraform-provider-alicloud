---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polar_db_extension"
description: |-
  Provides a Alicloud Polar Db Extension resource.
---

# alicloud_polar_db_extension

Provides a Polar Db Extension resource.



For information about Polar Db Extension and how to use it, see [What is Extension](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateExtensions).

-> **NOTE:** Available since v1.264.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

data "alicloud_polardb_node_classes" "default" {
  db_type  = "PostgreSQL"
  pay_type = "PostPaid"
  category = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "dbcluster" {
  default_time_zone = "SYSTEM"
  creation_category = "Normal"
  zone_id           = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  creation_option   = "Normal"
  db_version        = "14"
  pay_type          = "PostPaid"
  db_type           = "PostgreSQL"
  db_node_class     = "polar.pg.x4.medium.c"
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_polardb_account" "account" {
  account_type     = "Normal"
  account_name     = "nzh"
  account_password = "Ali123456"
  db_cluster_id    = alicloud_polardb_cluster.dbcluster.id
}

resource "alicloud_polardb_database" "database" {
  character_set_name = "UTF8"
  db_description     = var.name
  db_cluster_id      = alicloud_polardb_cluster.dbcluster.id
  db_name            = "nzh"
  account_name       = alicloud_polardb_account.account.db_cluster_id
}


resource "alicloud_polar_db_extension" "default" {
  extension_name = "postgres_fdw"
  db_cluster_id  = alicloud_polardb_cluster.dbcluster.id
  account_name   = alicloud_polardb_account.account.account_name
  db_name        = alicloud_polardb_database.database.db_name
}
```

## Argument Reference

The following arguments are supported:
* `account_name` - (Required, ForceNew) The database account name of the associated PolarDB cluster. Only support `Super` account.
* `db_cluster_id` - (Required, ForceNew) The ID of the cluster.
* `db_name` - (Required, ForceNew) PolarDB cluster database name.
* `extension_name` - (Required, ForceNew) Information about the installed plug-ins under the specified database.
* `installed_version` - (Optional, Computed) Installed version, only supports upgrading to the default version.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_cluster_id>:<account_name>:<db_name>:<extension_name>`.
* `default_version` - Default version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Extension.
* `delete` - (Defaults to 5 mins) Used when delete the Extension.
* `update` - (Defaults to 5 mins) Used when update the Extension.

## Import

Polar Db Extension can be imported using the id, e.g.

```shell
$ terraform import alicloud_polar_db_extension.example <db_cluster_id>:<account_name>:<db_name>:<extension_name>
```