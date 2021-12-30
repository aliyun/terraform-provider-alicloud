---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_account_privilege"
sidebar_current: "docs-alicloud-resource-db-account-privilege"
description: |-
  Provides an RDS account privilege resource.
---

# alicloud\_db\_account\_privilege

Provides an RDS account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account.

-> **NOTE:** At present, a database can only have one database owner.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbaccountprivilegebasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name       = var.name
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

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  instance_id = alicloud_db_instance.instance.id
  name        = "tftestprivilege"
  password    = "Test12345"
  description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadOnly"
  db_names     = alicloud_db_database.db.*.name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs.
* `account_name` - (Required, ForceNew) A specified account name.
* `privilege` - The privilege of one account access database. Valid values: 
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

```
$ terraform import alicloud_db_account_privilege.example "rm-12345:tf_account:ReadOnly"
```
