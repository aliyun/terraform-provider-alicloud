---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_database"
sidebar_current: "docs-alicloud-resource-db-database"
description: |-
  Provides an RDS database resource.
---

# alicloud_db_database

Provides an RDS database resource. A DB database deployed in a DB instance. A DB instance can own multiple databases, see [What is DB Database](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-createdatabase).

-> **NOTE:** Available since v1.5.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_db_database&exampleId=a019d9a3-34c1-1d16-a785-63ac687f2aeafc448f0a&activeTab=example&spm=docs.r.db_database.0.a019d9a334&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_db_zones" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_db_database&spm=docs.r.db_database.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `name` - (Required, ForceNew) The name of the database. 
* -> **NOTE:** 
  The name must be 2 to 64 characters in length.
  The name must start with a lowercase letter and end with a lowercase letter or digit.
  The name can contain lowercase letters, digits, underscores (_), and hyphens (-).
  The name must be unique within the instance.
  For more information about invalid characters, see [Forbidden keywords table](https://help.aliyun.com/zh/rds/developer-reference/forbidden-keywords?spm=api-workbench.api_explorer.0.0.20e15f16d1z52p).

* `character_set` - (Optional, ForceNew) Character set. The value range is limited to the following:
    - MySQL: [ utf8, gbk, latin1, utf8mb4 ] \(`utf8mb4` only supports versions 5.5 and 5.6\).
    - SQLServer: [ Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS, Chinese_PRC_BIN ]
    - PostgreSQL: Valid values for PostgreSQL databases: a value in the `character set,<Collate>,<Ctype>` format. Example: `UTF8,C,en_US.utf8`.
    > - Valid values for the character set : [ KOI8U, UTF8, WIN866, WIN874, WIN1250, WIN1251, WIN1252, WIN1253, WIN1254, WIN1255, WIN1256, WIN1257, WIN1258, EUC_CN, EUC_KR, EUC_TW, EUC_JP, EUC_JIS_2004, KOI8R, MULE_INTERNAL, LATIN1, LATIN2, LATIN3, LATIN4, LATIN5, LATIN6, LATIN7, LATIN8, LATIN9, LATIN10, ISO_8859_5, ISO_8859_6, ISO_8859_7, ISO_8859_8, SQL_ASCII ]
    > - Valid values for the Collate field: You can execute the `SELECT DISTINCT collname FROM pg_collation;` statement to obtain the field value. The default value is `C`.
    > - Valid values for the Ctype field: You can execute the `SELECT DISTINCT collctype FROM pg_collation;` statement to obtain the field value. The default value is `en_US.utf8`.
    - MariaDB: [ utf8, gbk, latin1, utf8mb4 ]
  
   More details refer to [API Docs](https://www.alibabacloud.com/help/zh/doc-detail/26258.htm)

* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.

-> **NOTE:** The value of "name" or "character_set"  does not support modification.


## Attributes Reference

The following attributes are exported:

* `id` - The current database resource ID. Composed of instance ID and database name with format `<instance_id>:<name>`.

## Import

RDS database can be imported using the id, e.g.

```shell
$ terraform import alicloud_db_database.example "rm-12345:tf_database"
```
