---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_account"
sidebar_current: "docs-alicloud-resource-db-account"
description: |-
  Provides an RDS account resource.
---

# alicloud\_db\_account

Provides an RDS account resource and used to manage databases. A RDS instance supports multiple database account.

## Example Usage

```
resource "alicloud_db_account" "default" {
	instance_id = "rm-2eps..."
	name = "tf_account"
	password = "..."
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs.
* `name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `password` - (Required) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters.
* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `type` - (Optional, ForceNew)Privilege type of account.
    - Normal: Common privilege.
    - Super: High privilege.

    Default to Normal. It is is valid for MySQL 5.5/5.6 only.
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