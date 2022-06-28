---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_account_privilege"
sidebar_current: "docs-alicloud-resource-polardb-account-privilege"
description: |-
  Provides a PolarDB account privilege resource.
---

# alicloud\_polardb\_account\_privilege

Provides a PolarDB account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account.

-> **NOTE:** Available in v1.67.0+.

## Example Usage

```terraform
variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "dbaccountprivilegebasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
  name       = var.name
}

resource "alicloud_polardb_cluster" "cluster" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_database" "db" {
  db_cluster_id = alicloud_polardb_cluster.cluster.id
  db_name       = "tftestdatabase"
}

resource "alicloud_polardb_account" "account" {
  db_cluster_id       = alicloud_polardb_cluster.cluster.id
  account_name        = "tftestnormal"
  account_password    = "Test12345"
  account_description = var.name
}

resource "alicloud_polardb_account_privilege" "privilege" {
  db_cluster_id     = alicloud_polardb_cluster.cluster.id
  account_name      = alicloud_polardb_account.account.account_name
  account_privilege = "ReadOnly"
  db_names          = [alicloud_polardb_database.db.db_name]
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `account_name` - (Required, ForceNew) A specified account name.
* `account_privilege` - (Optional, ForceNew) The privilege of one account access database. Valid values: ["ReadOnly", "ReadWrite"], ["DMLOnly", "DDLOnly"] added since version v1.101.0. Default to "ReadOnly".
* `db_names` - (Required) List of specified database name.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID, account name and privilege with format `<db_cluster_id>:<account_name>:<account_privilege>`.

## Import

PolarDB account privilege can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_account_privilege.example "pc-12345:tf_account:ReadOnly"
```
