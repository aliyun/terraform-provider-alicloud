---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_account"
sidebar_current: "docs-alicloud-resource-db-account"
description: |-
  Provides an RDS account resource.
---

# alicloud\_db\_account

Provides an RDS account resource and used to manage databases.

-> **NOTE:** Currently, only MySQL、PostgreSQL、MariaDB、SQL Server（exclude SQL Server 2017 clustered edition）instance support creating a `Normal` account. Other engine instance, like PPAS and SQL Server 2017, only support creating a `Super` account, and you can log on to the database to create other accounts using this Super account.
> **NOTE:** Because the `Super` account can not be deleted, there does not suggest to manage `Super` account using this resource. Otherwise, this resource can not be deleted when account is `Super`.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbaccountmysql"
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

resource "alicloud_db_account" "account" {
  instance_id = "${alicloud_db_instance.instance.id}"
  name        = "tftestnormal"
  password    = "Test12345"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs.
* `name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `password` - (Optional, Sensitive) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters. You have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Optional, Available in 1.57.1+) An KMS encrypts password used to a db account. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `type` - (Optional, ForceNew)Privilege type of account.
    - Normal: Common privilege.
    - Super: High privilege.
    - defalut Normal to MySQL、MariaDB、SQL Server(exclude SQL Server 2017 clustered edition).
    - defalut Super to PostgreSQL, PPAS, SQL Server 2017 clustered edition.
    Currently, MySQL 5.7, SQL Server 2012/2016, PostgreSQL, and PPAS each can have only one initial account.
    Other accounts are created by the initial account that has logged on to the database. [Refer to details](https://www.alibabacloud.com/help/doc-detail/26263.htm).

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.
* `instance_id` - The Id of DB instance.
* `name` - The name of DB account.
* `description` - The account description.
* `type` - Privilege type of account.

## Import

RDS account can be imported using the id, e.g.

```
$ terraform import alicloud_db_account.example "rm-12345:tf_account"
```
