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

-> **NOTE:** The `account_password_wo` and `account_password_wo_version` arguments require Terraform v1.11.0 or later.

-> **NOTE:** The DynamoDB type account does not support deletion. When destroying the Terraform resource, the DynamoDB account will be removed from state but not deleted from the cloud.

-> **NOTE:** The `password_type` argument is only used when the account password is updated and is not read back by this resource, so it cannot be restored by import. The `db_name`, `account_privilege` and `priv_for_all_db` arguments are only used when the account is created and are not read back into the resource state, so they cannot be restored by import. The `DescribeAccounts` API does return database privileges in `DatabasePrivileges`, but this resource does not map them back to these arguments. For a PolarDB Search node account, `node_type` is part of the resource ID and is restored on import.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_account&exampleId=d304ccab-d255-da40-7c04-6026efedb0938bbe617a&activeTab=example&spm=docs.r.polardb_account.0.d304ccabd2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

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
  db_type            = "MySQL"
  db_version         = "8.0"
  db_node_class      = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type           = "PostPaid"
  vswitch_id         = alicloud_vswitch.default.id
  description        = "terraform-example"
  encrypt_new_tables = "OFF"
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
* `account_password` - (Optional) The account password. It conflicts with `kms_encrypted_password` and `account_password_wo`, and one of `account_password`, `kms_encrypted_password` and `account_password_wo` must be set. Must  meet the following requirements:
  - Contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters.
  - Be 8 to 32 characters in length.
  - Special characters include !@#$%^&*()_+-=.
* `account_password_wo` - (Optional, Available since v2.0.0-beta5) The write-only password of the account. It is only sent to the server when the resource is created or the `account_password_wo_version` changes and is never stored in the state or plan. It requires Terraform v1.11.0 or later and must be used together with `account_password_wo_version`. It conflicts with `account_password` and `kms_encrypted_password`, and one of `account_password`, `kms_encrypted_password` and `account_password_wo` must be set.
* `account_password_wo_version` - (Optional, Available since v2.0.0-beta5) The version of the write-only password. The write-only password is re-sent to the server only when this argument changes, so it must be changed whenever `account_password_wo` changes.
* `account_password_valid_time` - (Optional, Available since v1.265.0) The time when the password for the database account expires.
* `account_privilege` - (Optional, ForceNew, String, Available since v2.2.0) The permissions to grant the account on the databases specified by `db_name` when it is created. Valid values: `ReadWrite`, `ReadOnly`, `DMLOnly`, `DDLOnly`, `ReadIndex`. It must be set together with `db_name`; setting it alone is rejected at plan time. When `db_name` contains multiple database names, separate the permissions with commas (,) and keep the same order, for example `ReadWrite,ReadOnly`; each comma-separated value must be one of the valid values above and the length of `account_privilege` must not exceed 900 characters. Only standard accounts of PolarDB for MySQL clusters support this parameter. It is only sent when the account is created and cannot be restored by import; after the resource is imported, a configured `account_privilege` therefore plans a replacement. Changing this argument replaces the account, and the database privileges of the account should not be managed by both this argument and the `alicloud_polardb_account_privilege` resource at the same time.
* `account_type` - (Optional, ForceNew, String) The account type. Defaults to `Normal`. Valid values: `Normal`, `Super`, `DynamoDB`. The API reports a privilege-scoped type for some accounts, for example `ReadOnly` for an account created with a cluster-wide read-only privilege. The state follows the value reported by the API, and the diff from such a derived value back to the configured `Normal` is suppressed, so the plan converges after apply and the account is not replaced; changing `account_type` explicitly still replaces the account.
* `db_cluster_id` - (Required, ForceNew) The cluster ID.
* `db_name` - (Optional, ForceNew, String, Available since v2.2.0) The names of the databases configured for the account when it is created. Separate multiple database names with commas (,). Only standard accounts of PolarDB for MySQL clusters support this parameter. Set it together with `account_privilege` to grant the account permissions on these databases; the privileges are applied by the API at creation time and are not returned by this resource. It is only sent when the account is created and cannot be restored by import; after the resource is imported, a configured `db_name` therefore plans a replacement. Changing this argument replaces the account, and the database privileges of the account should not be managed by both this argument and the `alicloud_polardb_account_privilege` resource at the same time.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. It conflicts with `account_password` and `account_password_wo`, and one of `account_password`, `kms_encrypted_password` and `account_password_wo` must be set.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `node_type` - (Optional, ForceNew, String, Available since v2.2.0) The node type. Valid values: `Search`. This parameter is required only when the account is created for a PolarDB Search node. Both `CreateAccount` and `DescribeAccounts` send it: the query API requires it to locate a Search node account, so the resource ID of such an account is `<db_cluster_id>:<account_name>:Search`. Changing this argument replaces the account.
* `password_type` - (Optional, String, Available since v2.2.0) The password type of the account, for example, `Tair`. It is sent together with the new password when the account password is updated. The API only accepts it for Tair accounts; password updates for accounts of other types are rejected by the API.
* `priv_for_all_db` - (Optional, ForceNew, String, Available since v2.2.0) Specifies whether to grant the account permissions on all databases of the cluster and on all databases that are added later. Valid values: `0` and `1`. The API applies it only when the account is created together with `account_privilege`; setting it to `1` without `account_privilege` is rejected at plan time. When it is set to `1`, the CreateAccount API grants the permissions on all current and future databases of the cluster, and the query API does not list these privileges per database (verified for this account type), so the resource cannot verify them individually. It is only sent when the account is created and is not read back by the query API. Changing this argument replaces the account.
* `reset_permission_flag` - (Optional, Bool, Available since v2.2.0) Whether to reset the permissions of a high-privilege account to restore them when they were revoked by accident. Set it to `true` and change it to trigger the reset. The reset only applies to `Super` accounts of PolarDB for MySQL clusters and re-sends the configured password when one is set.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_cluster_id>:<account_name>`, or `<db_cluster_id>:<account_name>:<node_type>` for an account created with `node_type`.
* `status` - (Available since v1.265.0) The status of the database account.
* `dynamodb_auth_password` - (Sensitive, Available since v1.285.0) The DynamoDB authentication password. Only available for DynamoDB account type.

## Timeouts

-> **NOTE:** Available since v1.265.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Account.
* `delete` - (Defaults to 8 mins) Used when delete the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

Polar Db Account can be imported using the id, which consists of db_cluster_id and account_name, e.g.

```shell
$ terraform import alicloud_polardb_account.example <db_cluster_id>:<account_name>
```

For an account created with `node_type`, append the node type to the id, e.g.

```shell
$ terraform import alicloud_polardb_account.example <db_cluster_id>:<account_name>:Search
```
