---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_account_privilege"
sidebar_current: "docs-alicloud-resource-db-account-privilege"
description: |-
  Provides an RDS account privilege resource.
---

# alicloud_db_account_privilege

Provides an RDS account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account, see [What is DB Account Privilege](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-grantaccountprivilege).

-> **NOTE:** At present, a database can only have one database owner.

-> **NOTE:** Available since v1.5.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_db_account_privilege&exampleId=21b01937-a938-4a74-7007-8f4af584b76926ed62d1&activeTab=example&spm=docs.r.db_account_privilege.0.21b01937a9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
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

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = "${var.name}_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.instance.id
  account_name        = "tfexample"
  account_password    = "Example12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.account_name
  privilege    = "ReadOnly"
  db_names     = alicloud_db_database.db.*.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_db_account_privilege&spm=docs.r.db_account_privilege.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs.
* `account_name` - (Required, ForceNew) A specified account name.
* `privilege` - (Optional, ForceNew) The privilege of one account access database. Valid values: 
    - ReadOnly: This value is only for MySQL, MariaDB and SQL Server
    - ReadWrite: This value is only for MySQL, MariaDB and SQL Server
    - DDLOnly: (Available in 1.64.0+) This value is only for MySQL and MariaDB
    - DMLOnly: (Available in 1.64.0+) This value is only for MySQL and MariaDB
    - DBOwner: (Available in 1.64.0+) This value is only for SQL Server and PostgreSQL.
      Default to "ReadOnly". 
* `db_names` - (Required) List of specified database name.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID, account name and privilege with format `<instance_id>:<name>:<privilege>`.

## Import

RDS account privilege can be imported using the id, e.g.

```shell
$ terraform import alicloud_db_account_privilege.example "rm-12345:tf_account:ReadOnly"
```
