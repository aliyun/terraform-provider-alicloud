---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_database"
sidebar_current: "docs-alicloud-resource-polardb-database"
description: |-
  Provides a PolarDB database resource.
---

# alicloud\_polardb\_database

Provides a PolarDB database resource. A database deployed in a PolarDB cluster. A PolarDB cluster can own multiple databases.

-> **NOTE:** Available in v1.66.0+.

## Example Usage

```terraform
variable "name" {
  default = "polardbClusterconfig"
}

variable "creation" {
  default = "PolarDB"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
  name       = var.name
}

resource "alicloud_polardb_cluster" "cluster" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "testDB"
}

resource "alicloud_polardb_database" "default" {
  db_cluster_id = alicloud_polardb_cluster.cluster.id
  db_name       = "tftestdatabase"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `db_name` - (Required, ForceNew) Name of the database requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letterand have no more than 64 characters.
* `character_set_name` - (Optional,ForceNew) Character set. The value range is limited to the following: [ utf8, gbk, latin1, utf8mb4, Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ], default is "utf8" \(`utf8mb4` only supports versions 5.5 and 5.6\).
* `db_description` - (Optional) Database description. It must start with a Chinese character or English letter, cannot start with "http://" or "https://". It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length must be 2-256 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of cluster ID and database name with format `<cluster_id>:<name>`.

## Import

PolarDB database can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_database.example "pc-12345:tf_database"
```
