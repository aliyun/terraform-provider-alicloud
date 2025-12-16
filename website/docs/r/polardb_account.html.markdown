---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_account"
description: |-
  Provides a Alicloud Polar Db Account resource.
---

# alicloud_polardb_account

Provides a Polar Db Account resource.

Database account information.

For information about Polar Db Account and how to use it, see [What is Account](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateAccount).

-> **NOTE:** Available since v1.67.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_account&exampleId=d304ccab-d255-da40-7c04-6026efedb0938bbe617a&activeTab=example&spm=docs.r.polardb_account.0.d304ccabd2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"
}

resource "alicloud_polardb_account" "default" {
  db_cluster_id       = alicloud_polardb_cluster.default.id
  account_name        = "terraform_example"
  account_password    = "Example1234"
  account_description = "terraform-example"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_account&spm=docs.r.polardb_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `account_description` - (Optional) The description of the database account.
* `account_lock_state` - (Optional, Available since v1.265.0) The lock status of the account. Valid values:
  - `UnLock`: The account is not locked.
  - `Lock`: The account is locked.
* `account_name` - (Required, ForceNew) The account name. Must meet the following requirements:
  - Start with a lowercase letter and end with a letter or number.
  - Consists of lowercase letters, numbers, or underscores.
  - The length is 2 to 16 characters.
  - You cannot use some reserved usernames, such as root and admin.
* `account_password` - (Optional) The account password. You have to specify one of `account_password` and `kms_encrypted_password` fields. Must  meet the following requirements:
  - Contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters.
  - Be 8 to 32 characters in length.
  - Special characters include !@#$%^&*()_+-=.
* `account_password_valid_time` - (Optional, Available since v1.265.0) The time when the password for the database account expires.
* `account_type` - (Optional, ForceNew) The account type. Default value:`Normal`. Valid values: `Normal`, `Super`.
* `db_cluster_id` - (Required, ForceNew) The cluster ID.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_cluster_id>:<account_name>`.
* `status` - (Available since v1.265.0) The status of the database account.

## Timeouts

-> **NOTE:** Available since v1.265.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Account.
* `delete` - (Defaults to 8 mins) Used when delete the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

Polar Db Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_account.example <db_cluster_id>:<account_name>
```
