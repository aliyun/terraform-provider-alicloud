---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_account"
description: |-
  Provides a Alicloud Tair (Redis OSS-Compatible) And Memcache (KVStore) Account resource.
---

# alicloud_kvstore_account

Provides a Tair (Redis OSS-Compatible) And Memcache (KVStore) Account resource.



For information about Tair (Redis OSS-Compatible) And Memcache (KVStore) Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/redis/developer-reference/api-r-kvstore-2015-01-01-createaccount-redis).

-> **NOTE:** Available since v1.66.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kvstore_account&exampleId=2debea4c-9da4-77c8-78e6-ad3dd92814f2a421dbc5&activeTab=example&spm=docs.r.kvstore_account.0.2debea4c9d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_kvstore_zones" "default" {

}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_kvstore_instance" "default" {
  db_instance_name  = var.name
  vswitch_id        = alicloud_vswitch.default.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  zone_id           = data.alicloud_kvstore_zones.default.zones.0.id
  instance_class    = "redis.master.large.default"
  instance_type     = "Redis"
  engine_version    = "5.0"
  security_ips      = ["10.23.12.24"]
  config = {
    appendonly             = "yes"
    lazyfree-lazy-eviction = "yes"
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_kvstore_account" "default" {
  account_name     = "tfexamplename"
  account_password = "YourPassword_123"
  instance_id      = alicloud_kvstore_instance.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_kvstore_account&spm=docs.r.kvstore_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) The name of the account. The name must meet the following requirements:
  * The name must start with a lowercase letter and can contain lowercase letters, digits, and underscores (_).
  * The name can be up to 100 characters in length.
  * The name cannot be one of the reserved words listed in the [Reserved words for Redis account names](https://www.alibabacloud.com/help/en/redis/user-guide/create-and-manage-database-accounts) section.
* `account_password` - (Optional, Sensitive) The password of the account. The password must be 8 to 32 characters in length. It must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `!@ # $ % ^ & * ( ) _ + - =`. You have to specify one of `account_password` and `kms_encrypted_password` fields.
* `description` - (Optional) Database description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs (The engine version of instance must be 4.0 or 4.0+).
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a Tair (Redis OSS-Compatible) And Memcache (KVStore) account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a Tair (Redis OSS-Compatible) And Memcache (KVStore) account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `account_type` - (Optional, ForceNew) Privilege type of account.
    - Normal: Common privilege.
    Default to Normal.
* `account_privilege` - (Optional) The privilege of account access database. Default value: `RoleReadWrite` 
    - `RoleReadOnly`: This value is only for Redis and Memcache
    - `RoleReadWrite`: This value is only for Redis and Memcache

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value is formatted `<instance_id>:<account_name>`.
* `status` - The status of Tair (Redis OSS-Compatible) And Memcache (KVStore) Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 6 mins) Used when create the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.
* `update` - (Defaults to 7 mins) Used when update the Account.

## Import

Tair (Redis OSS-Compatible) And Memcache (KVStore) Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_kvstore_account.example <instance_id>:<account_name>
```