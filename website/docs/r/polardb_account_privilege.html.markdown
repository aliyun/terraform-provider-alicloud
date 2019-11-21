---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_account_privilege"
sidebar_current: "docs-alicloud-resource-polardb-account-privilege"
description: |-
  Provides an PolarDB account privilege resource.
---

# alicloud\_polardb\_account\_privilege

Provides an PolarDB account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account.

-> **NOTE:** Available in v1.66.0+.

## Example Usage

```
variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "dbaccountprivilegebasic"
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

resource "alicloud_polardb_cluster" "default" {
  db_type = "MySQL"
  db_version = "8.0"
  pay_type = "PostPaid"
  db_node_class    = "polar.mysql.x4.large"
  vswitch_id = "${alicloud_vswitch.default.id}"
  description = "${var.name}"
}

resource "alicloud_polardb_database" "db" {
  count       = 2
  instance_id = "${alicloud_polardb_instance.cluster.id}"
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_polardb_account" "account" {
  instance_id = "${alicloud_polardb_instance.cluster.id}"
  name        = "tftestprivilege"
  password    = "Test12345"
  description = "from terraform"
}

resource "alicloud_polardb_account_privilege" "privilege" {
  cluster_id    = "${alicloud_polardb_instance.cluster.id}"
  account_name  = "${alicloud_polardb_account.account.name}"
  privilege     = "ReadOnly"
  db_names      = "${alicloud_polardb_database.db.*.name}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `account_name` - (Required, ForceNew) A specified account name.
* `privilege` - (Optional, ForceNew) The privilege of one account access database. Valid values: ["ReadOnly", "ReadWrite"]. Default to "ReadOnly".
* `db_names` - (Required) List of specified database name.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID, account name and privilege with format `<cluster_id>:<name>:<privilege>`.
* `instance_id` - The Id of DB instance.
* `account_name` - The name of DB account.
* `privilege` - The specified account privilege.
* `db_names` - List of granted privilege database names.

## Import

PolarDB account privilege can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_account_privilege.example "pc-12345:tf_account:ReadOnly"
```