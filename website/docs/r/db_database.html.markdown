---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_database"
sidebar_current: "docs-alicloud-resource-db-database"
description: |-
  Provides an RDS database resource.
---

# alicloud\_db\_database

Provides an RDS database resource. A DB database deployed in a DB instance. A DB instance can own multiple databases.

~> **NOTE:** At present, it does not support creating 'PostgreSQL' and 'PPAS' database. You have to login RDS instance to create manually.

## Example Usage

```
resource "alicloud_db_database" "default" {
	instance_id = "rm-2eps..."
	name = "tf_database"
	character_set = "utf8"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The Id of instance that can run database.
* `name` - (Required) Name of the database requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter
                      and have no more than 64 characters.
* `character_set` - (Required) Character set. The value range is limited to the following:
    - MySQL: [ utf8, gbk, latin1, utf8mb4 ] \(included in versions 5.5 and 5.6\).
    - SQLServer: [ Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ]

* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.


## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of instance ID and database name with format "<instance_id>:<name>".
* `instance_id` - The Id of DB instance.
* `name` - The name of DB database.
* `character_set` - Character set that database used.
* `description` - The database description.

## Import

RDS database can be imported using the id, e.g.

```
$ terraform import alicloud_db_database.example "rm-12345:tf_database"
```
