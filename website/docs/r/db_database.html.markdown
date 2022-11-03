---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_database"
sidebar_current: "docs-alicloud-resource-db-database"
description: |-
  Provides an RDS database resource.
---

# alicloud\_db\_database

Provides an RDS database resource. A DB database deployed in a DB instance. A DB instance can own multiple databases.

-> **NOTE:** This resource does not support creating 'PPAS' database. You have to login RDS instance to create manually.

## Example Usage

```terraform
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbdatabasebasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.instance.id
  name        = "tftestdatabase"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `name` - (Required, ForceNew) Name of the database requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter
                      and have no more than 64 characters.
* `character_set` - (Required) Character set. The value range is limited to the following:
    - MySQL: [ utf8, gbk, latin1, utf8mb4 ] \(`utf8mb4` only supports versions 5.5 and 5.6\).
    - SQLServer: [ Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ]
    - PostgreSQL: [ KOI8U、UTF8、WIN866、WIN874、WIN1250、WIN1251、WIN1252、WIN1253、WIN1254、WIN1255、WIN1256、WIN1257、WIN1258、EUC_CN、EUC_KR、EUC_TW、EUC_JP、EUC_JIS_2004、KOI8R、MULE_INTERNAL、LATIN1、LATIN2、LATIN3、LATIN4、LATIN5、LATIN6、LATIN7、LATIN8、LATIN9、LATIN10、ISO_8859_5、ISO_8859_6、ISO_8859_7、ISO_8859_8、SQL_ASCII ]
  
   More details refer to [API Docs](https://www.alibabacloud.com/help/zh/doc-detail/26258.htm)

* `description` - (ForceNew) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.

-> **NOTE:** The value of "name" or "character_set"  does not support modification.


## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of instance ID and database name with format `<instance_id>:<name>`.

## Import

RDS database can be imported using the id, e.g.

```shell
$ terraform import alicloud_db_database.example "rm-12345:tf_database"
```
