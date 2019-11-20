---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_account"
sidebar_current: "docs-alicloud-resource-polardb-account"
description: |-
  Provides an RDS account resource.
---

# alicloud\_polardb\_account

Provides an PolarDB account resource and used to manage databases.

-> **NOTE:** Available in v1.66.0+. Currently, only MySQL、MariaDB、SQL Server（exclude SQL Server 2017 clustered edition）instance support creating a `Normal` account. Other engine instance, like PostgreSQL, PPAS and SQL Server 2017, only support creating a `Super` account, and you can log on to the database to create other accounts using this Super account.
> **NOTE:** Because the `Super` account can not be deleted, there does not suggest to manage `Super` account using this resource. Otherwise, this resource can not be deleted when account is `Super`.

## Example Usage

```
variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "polardbaccountmysql"
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

resource "alicloud_polardb_cluster" "cluster" {
  db_type               = "MySQL"
  db_version            = "8.0"
  db_node_class         = "polar.mysql.x4.large"
  pay_type              = "PostPaid"
  vswitch_id            = "${alicloud_vswitch.default.id}"
  description           = "${var.name}"
}

resource "alicloud_db_account" "account" {
  cluster_id    = "${alicloud_db_instance.instance.id}"
  name          = "tftestnormal"
  password      = "Test12345"
  description   = "${var.name}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `password` - (Required) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters.
* `description` - (Optional) Account description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `type` - (Optional, ForceNew) Account type, Valid values are `Normal`, `Supper`, Default to `Normal`.
## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.
* `cluster_id` - The Id of cluster.
* `name` - The name of DB account.
* `description` - The account description.
* `type` - Privilege type of account.

## Import

PolarDB account can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_account.example "pc-12345:tf_account"
```