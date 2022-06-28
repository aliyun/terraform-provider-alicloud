---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_account"
sidebar_current: "docs-alicloud-resource-polardb-account"
description: |-
  Provides a RDS account resource.
---

# alicloud\_polardb\_account

Provides a PolarDB account resource and used to manage databases.

-> **NOTE:** Available in v1.67.0+. 

## Example Usage

```terraform
variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "polardbaccountmysql"
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
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "cluster" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.large"
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_account" "account" {
  db_cluster_id       = alicloud_polardb_cluster.cluster.id
  account_name        = "tftestnormal"
  account_password    = "Test12345"
  account_description = var.name
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `account_name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `account_password` - (Required) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `account_description` - (Optional) Account description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `account_type` - (Optional, ForceNew) Account type, Valid values are `Normal`, `Super`, Default to `Normal`.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.

## Import

PolarDB account can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_account.example "pc-12345:tf_account"
```
