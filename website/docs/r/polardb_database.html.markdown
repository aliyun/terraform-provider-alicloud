---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_database"
sidebar_current: "docs-alicloud-resource-polardb-database"
description: |-
  Provides an PolarDB database resource.
---

# alicloud\_db\_database

Provides an PolarDB database resource. A DB database deployed in a DB cluster. A DB cluster can own multiple databases.

-> **NOTE:** Available in v1.66.0+.

## Example Usage

```
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "testDB"
	}

	variable "clusterchargetype" {
		default = "PostPaid"
	}

	variable "engine" {
		default = "MySQL"
	}

	variable "engineversion" {
		default = "8.0"
	}

	variable "clusterclass" {
		default = "polar.mysql.x4.large"
	}

	resource "alicloud_polardb_cluster" "cluster" {
		db_type = "${var.engine}"
		db_version = "${var.engineversion}"
		pay_type = "${var.clusterchargetype}"
		db_node_class = "${var.clusterclass}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		description = "${var.name}"
	}

    resource "alicloud_polardb_database" "default" {
      cluster_id = "${alicloud_polardb_cluster.cluster.id}"
      name        = "tftestdatabase"
    }
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `name` - (Required, ForceNew) Name of the database requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letterand have no more than 64 characters.
* `character_set` - (Optional,ForceNew) Character set. The value range is limited to the following: [ utf8, gbk, latin1, utf8mb4, Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ] \(`utf8mb4` only supports versions 5.5 and 5.6\).
* `account_name` - (Required) The account_name of cluster that can run database.
* `account_privilege` - (Optional) The Id of cluster that can run database. Valid values are `ReadOnly`, `ReadWrite`,`DMLOnly`, `DDLOnly`, Default to `ReadWrite`.
* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.


## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of cluster ID and database name with format `<cluster_id>:<name>`.
* `cluster_id` - The Id of DB cluster.
* `name` - The name of DB database.
* `character_set` - Character set that database used.
* `description` - The database description.

## Import

PolarDB database can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_database.example "pc-12345:tf_database"
```
