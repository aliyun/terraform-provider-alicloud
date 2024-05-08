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
data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"
}

resource "alicloud_polardb_account" "default" {
  db_cluster_id       = alicloud_polardb_cluster.default.id
  account_name        = "terraform_example"
  account_password    = "Example1234"
  account_description = "terraform-example"
}

resource "alicloud_polardb_database" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_name       = "terraform-example"
}

resource "alicloud_polardb_account_privilege" "default" {
  db_cluster_id     = alicloud_polardb_cluster.default.id
  account_name      = alicloud_polardb_account.default.account_name
  account_privilege = "ReadOnly"
  db_names          = [alicloud_polardb_database.default.db_name]
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

```shell
$ terraform import alicloud_polardb_account_privilege.example "pc-12345:tf_account:ReadOnly"
```
