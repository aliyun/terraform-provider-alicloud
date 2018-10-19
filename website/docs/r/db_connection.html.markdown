---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_connection"
sidebar_current: "docs-alicloud-resource-db-connection"
description: |-
  Provides an RDS instance connection resource.
---

# alicloud\_db\_connection

Provides an RDS connection resource to allocate an Internet connection string for RDS instance.

~> **NOTE:** Each RDS instance will allocate a intranet connnection string automatically and its prifix is RDS instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

```
resource "alicloud_db_connection" "default" {
	instance_id = "rm-2eps..."
	connection_prefix = "alicloud"
	port = "3306"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The Id of instance that can run database.
* `connection_prefix` - (Optional) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + 'tf'.
* `port` - (Optional) Internet connection port. Valid value: [3001-3999]. Default to 3306.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.
* `connection_prefix` - Prefix of a connection string.
* `port` - Connection instance port.
* `connection_string` - Connection instance string.
* `ip_address` - The ip address of connection string.

## Import

RDS connection can be imported using the id, e.g.

```
$ terraform import alicloud_db_connection.example abc12345678
```