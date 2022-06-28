---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_account"
sidebar_current: "docs-alicloud-resource-kvstore-account"
description: |-
  Provides a Alicloud KVStore Account resource.
---

# alicloud\_kvstore\_account

Provides a KVStore Account resource.

For information about KVStore Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/doc-detail/95973.htm).

-> **NOTE:** Available in 1.66.0+

## Example Usage

Basic Usage

```terraform
variable "creation" {
  default = "KVStore"
}

variable "name" {
  default = "kvstoreinstancevpc"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_kvstore_instance" "default" {
  instance_class = "redis.master.small.default"
  instance_name  = var.name
  vswitch_id     = alicloud_vswitch.default.id
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "Redis"
  engine_version = "4.0"
}

resource "alicloud_kvstore_account" "example" {
  account_name     = "tftestnormal"
  account_password = "YourPassword_123"
  instance_id      = alicloud_kvstore_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) The name of the account. The name must meet the following requirements:
  * The name can contain lowercase letters, digits, and hyphens (-), and must start with a lowercase letter.
  * The name can be up to 100 characters in length.
  * The name cannot be one of the reserved words in the [Reserved words for Redis account names](https://www.alibabacloud.com/help/zh/doc-detail/92665.htm) section.
* `account_password` - (Optional, Sensitive) The password of the account. The password must be 8 to 32 characters in length. It must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `!@ # $ % ^ & * ( ) _ + - =`. You have to specify one of `account_password` and `kms_encrypted_password` fields.
* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs (The engine version of instance must be 4.0 or 4.0+).
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a KVStore account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a KVStore account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `account_type` - (Optional, ForceNew) Privilege type of account.
    - Normal: Common privilege.
    Default to Normal.
* `account_privilege` - (Optional) The privilege of account access database. Default value: `RoleReadWrite` 
    - `RoleReadOnly`: This value is only for Redis and Memcache
    - `RoleReadWrite`: This value is only for Redis and Memcache

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value is formatted `<instance_id>:<account_name>`.
* `status` - The status of KVStore Account.

### Timeouts

-> **NOTE:** Available in 1.102.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Account.
* `update` - (Defaults to 6 mins) Used when update the Account.

## Import

KVStore account can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_account.example <instance_id>:<account_name>
```
