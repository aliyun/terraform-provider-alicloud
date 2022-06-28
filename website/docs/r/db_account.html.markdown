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

-> **DEPRECATED:**  This resource  has been deprecated from version `1.120.0`. Please use new resource [alicloud_rds_account](https://www.terraform.io/docs/providers/alicloud/r/rds_account).

## Example Usage

```
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
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_account" "account" {
  instance_id = alicloud_db_instance.instance.id
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
    
    Default to Normal.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.

## Import

RDS account can be imported using the id, e.g.

```
$ terraform import alicloud_db_account.example "rm-12345:tf_account"
```
