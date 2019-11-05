---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_database"
sidebar_current: "docs-alicloud-resource-db-database"
description: |-
  Provides an RDS database resource.
---

# alicloud\_db\_database

Provides an RDS database resource. A DB database deployed in a DB instance. A DB instance can own multiple databases.

-> **NOTE:** This resource does not support creating 'PostgreSQL' database and
you can use [Postgresql Provider](https://www.terraform.io/docs/providers/postgresql/index.html) to do it.

-> **NOTE:** This resource does not support creating 'PPAS' database. You have to login RDS instance to create manually.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbdatabasebasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = "${alicloud_vswitch.default.id}"
  instance_name    = "${var.name}"
}

resource "alicloud_db_database" "default" {
  instance_id = "${alicloud_db_instance.instance.id}"
  name        = "tftestdatabase"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `name` - (Required, ForceNew) Name of the database requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter
                      and have no more than 64 characters.
                      
-> **NOTE:** The current field `name` does not support modification.
* `character_set` - (Required) Character set. The value range is limited to the following:
    - MySQL: [ utf8, gbk, latin1, utf8mb4 ] \(`utf8mb4` only supports versions 5.5 and 5.6\).
    - SQLServer: [ Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ]

-> **NOTE:** The current field `character_set` does not support modification.
* `description` - (ForceNew) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of instance ID and database name with format `<instance_id>:<name>`.

## Import

RDS database can be imported using the id, e.g.

```
$ terraform import alicloud_db_database.example "rm-12345:tf_database"
```
