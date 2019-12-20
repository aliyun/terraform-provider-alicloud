---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_account"
sidebar_current: "docs-alicloud-resource-kvstore-account"
description: |-
  Provides a kvstore account resource.
---

# alicloud\_kvstore\_account

Provides a kvstore account resource and used to manage databases.

-> **NOTE:** Available in 1.66.0+

## Example Usage

```
variable "creation" {
  default = "KVStore"
}
variable "name" {
  default = "kvstoreinstancevpc"
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
resource "alicloud_kvstore_instance" "default" {
  instance_class = "redis.master.small.default"
  instance_name  = "${var.name}"
  vswitch_id     = "${alicloud_vswitch.default.id}"
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "Redis"
  engine_version = "4.0"
}

resource "alicloud_kvstore_account" "account" {
  instance_id = "${alicloud_kvstore_instance.default.id}"
  account_name        = "tftestnormal"
  account_password    = "Test12345"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs. (The engine version of instance must be 4.0 or 4.0+)
* `account_name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `account_password` - (Optional, Sensitive) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters. You have to specify one of `account_password` and `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a KVStore account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a KVStore account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `account_type` - (Optional, ForceNew)Privilege type of account.
    - Normal: Common privilege.
    - Super: High privilege.
    Default to Normal.
* `account_privilege` - (Optional) The privilege of account access database. Valid values: 
    - RoleReadOnly: This value is only for Redis and Memcache
    - RoleReadWrite: This value is only for Redis and Memcache
    - RoleRepl: This value supports instance to read, write, and open SYNC / PSYNC commands.
                Only for Redis which engine version is 4.0 and architecture type is standard
     
   Default to "RoleReadWrite". 

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.

## Import

kvstore account can be imported using the id, e.g.

```
$ terraform import alicloud_KVStore_account.example "rm-12345:tf_account"
```
