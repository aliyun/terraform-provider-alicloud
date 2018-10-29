---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_account_privilege"
sidebar_current: "docs-alicloud-resource-db-account-privilege"
description: |-
  Provides an RDS account privilege resource.
---

# alicloud\_db\_account\_privilege

Provides an RDS account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account.

## Example Usage

```
resource "alicloud_db_database" "default" {
    count = 2
	instance_id = "rm-2eps..."
	name = "tf_database"
	character_set = "utf8"
}

resource "alicloud_db_account_privilege" "default" {
	instance_id = "rm-2eps..."
	account_name = "tf_account"
	privilege = "ReadOnly"
	db_names = ["${alicloud_db_database.base.*.name}"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The Id of instance in which account belongs.
* `account_name` - (Required) A specified account name.
* `privilege` - The privilege of one account access database. Valid values: ["ReadOnly", "ReadWrite"]. Default to "ReadOnly".
* `db_names` - (Optional) List of specified database name.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID, account name and privilege with format `<instance_id>:<name>:<privilege>`.
* `instance_id` - The Id of DB instance.
* `account_name` - The name of DB account.
* `privilege` - The specified account privilege.
* `db_names` - List of granted privilege database names.

## Import

RDS account privilege can be imported using the id, e.g.

```
$ terraform import alicloud_db_account_privilege.example "rm-12345:tf_account:ReadOnly"
```