---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_account"
sidebar_current: "docs-alicloud-resource-rds-account"
description: |-
  Provides a Alicloud RDS Account resource.
---

# alicloud\_rds\_account

Provides a RDS Account resource.

For information about RDS Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/doc-detail/26263.htm).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbaccountmysql"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
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

resource "alicloud_rds_account" "account" {
  db_instance_id   = alicloud_db_instance.instance.id
  account_name     = "tftestnormal12"
  account_password = "Test12345"
}
```

## Argument Reference

The following arguments are supported:

* `account_description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `account_name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and end with letters or numbers, The length must be 2-63 characters for PostgreSQL, otherwise the length must be 2-32 characters.
* `account_password` - (Optional, Sensitive) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters. You have to specify one of `password` and `kms_encrypted_password` fields.
* `account_type` - (Optional, Computed, ForceNew) Privilege type of account. Default to `Normal`.
    `Normal`: Common privilege.
    `Super`: High privilege. 
* `db_instance_id` - (Required, ForceNew) The Id of instance in which account belongs.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.

### Deprecated Attributes

* `description` - (Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_description` instead.
* `instance_id` - (Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `db_instance_id` instead.
* `name` - (Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_name` instead.
* `password` - (Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_password` instead.
* `type` - (Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_type` instead.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value is formatted `<db_instance_id>:<account_name>`.
* `status` - The status of the resource. Valid values: `Available`, `Unavailable`.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Account.
* `update` - (Defaults to 6 mins) Used when update the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.

## Import

RDS Account can be imported using the id, e.g.

```
$ terraform import alicloud_rds_account.example <db_instance_id>:<account_name>
```
