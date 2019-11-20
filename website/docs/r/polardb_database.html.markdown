---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_database"
sidebar_current: "docs-alicloud-resource-polardb-database"
description: |-
  Provides a PolarDB database resource.
---

# alicloud\_db\_database

Provides a PolarDB database resource. A DB database deployed in a DB cluster. A DB cluster can own multiple databases.

-> **NOTE:** Available in v1.66.0+.

## Example Usage

```
resource "alicloud_polardb_cluster" "cluster" {
    db_type = "MySQL"
    db_version = "8.0"
    pay_type = "PostPaid"
    db_node_class = "${var.clusterclass}"
    vswitch_id = "polar.mysql.x4.large"
    description = "testDB"
}

resource "alicloud_polardb_database" "default" {
  db_cluster_id = "${alicloud_polardb_cluster.cluster.id}"
  db_name        = "tftestdatabase"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `db_name` - (Required, ForceNew) Name of the database requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letterand have no more than 64 characters.
* `character_set_name` - (Optional,ForceNew) Character set. The value range is limited to the following: [ utf8, gbk, latin1, utf8mb4, Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ], default is "utf8" \(`utf8mb4` only supports versions 5.5 and 5.6\).
* `db_description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of cluster ID and database name with format `<cluster_id>:<name>`.

## Import

PolarDB database can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_database.example "pc-12345:tf_database"
```
