---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_database"
description: |-
  Provides a Alicloud RDS Database resource.
---

# alicloud_db_database

Provides a RDS Database resource.

Supports creating a database under an instance.

For information about RDS Database and how to use it, see [What is Database](https://next.api.alibabacloud.com/document/Rds/2014-08-15/CreateDatabase).

-> **NOTE:** Available since v1.5.0.

-> **NOTE:** The `collation_name`, `account_name` and `account_privilege` arguments take effect only when the database is created. Changing any of them forces a new database to be created.

-> **NOTE:** The account permission granted through `account_name` and `account_privilege` is not read back from the API. If the permission of the same account on this database is also managed by `alicloud_db_account_privilege`, both configurations write the same permission and the last applied one wins.

## Example Usage


<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_db_database&exampleId=7331e217-bfa5-6f58-36cf-e9568a53de6e55163ae0&activeTab=example&spm=docs.r.db_database.0.7331e217bf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
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
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = var.name
  instance_charge_type     = "Postpaid"
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = "terraform_example"
  account_password = "Test12345"
}

resource "alicloud_db_database" "default" {
  instance_id       = alicloud_db_instance.default.id
  data_base_name    = var.name
  character_set     = "gbk"
  collation_name    = "gbk_chinese_ci"
  account_name      = alicloud_rds_account.default.account_name
  account_privilege = "ReadWrite"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_db_database&spm=docs.r.db_database.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The instance ID. You can call the DescribeDBInstances operation to query the instance ID.
* `data_base_name` - (Optional, ForceNew, Available since v1.267.0) The name of the database.

  -> **NOTE:**
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
* `account_name` - (Optional, ForceNew, String, Available since v2.2.0) The name of the account to be granted permissions on the database when you create the database. You can call the DescribeAccounts operation to query the account name.

  -> **NOTE:** The account must belong to the instance specified by `instance_id` and must exist before the database is created.
* `account_privilege` - (Optional, ForceNew, String, Available since v2.2.0) The privilege granted to the account specified by `account_name` when you create the database. Valid values: `ReadOnly`, `ReadWrite`, `DDLOnly`, `DMLOnly`, `DBOwner`.

  -> **NOTE:** The valid values depend on the database engine. `DDLOnly` and `DMLOnly` are supported only by MySQL and MariaDB; `DBOwner` is supported only by SQL Server and PostgreSQL. This argument requires `account_name` to be set.
* `collation_name` - (Optional, ForceNew, String, Available since v2.2.0) The collation of the database. This parameter is supported only for MySQL instances. Specify a collation that matches `character_set`, for example `gbk_chinese_ci` for `gbk` or `utf8mb4_general_ci` for `utf8mb4`.
* `name` - (Optional, ForceNew, Deprecated from v1.267.0) The attribute has been deprecated from 1.267.0 and using `data_base_name` instead.

  -> **NOTE:** The value of "data_base_name" or "character_set"  does not support modification.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<data_base_name>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Database.
* `delete` - (Defaults to 5 mins) Used when delete the Database.
* `update` - (Defaults to 5 mins) Used when update the Database.

## Import

RDS Database can be imported using the id, e.g.

```shell
$ terraform import alicloud_db_database.example <instance_id>:<data_base_name>
```
