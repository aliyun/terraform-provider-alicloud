---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_account"
sidebar_current: "docs-alicloud-resource-rds-account"
description: |-
  Provides a Alicloud RDS Account resource.
---

# alicloud_rds_account

Provides a RDS Account resource.

For information about RDS Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-createaccount).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_account&exampleId=4fb8c8f5-60e2-40b7-1de2-c569a0e970826bbb2cb8&activeTab=example&spm=docs.r.rds_account.0.4fb8c8f560&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_db_zones" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

data "alicloud_db_instance_classes" "default" {
  zone_id        = data.alicloud_db_zones.default.ids.0
  engine         = "MySQL"
  engine_version = "5.6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.1.instance_class
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = var.name
  account_password = "Example1234"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_rds_account&spm=docs.r.rds_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Optional, ForceNew) The ID of the instance.
* `account_description` - (Optional) The description of the account. The value must be 2 to 256 characters in length. The value can contain letters, digits, underscores (_), and hyphens (-), and must start with a letter.

-> **NOTE:** The name cannot start with http:// or https://.
* `account_name` - (Optional, ForceNew) The name of the database account.
    * The name must be unique.
    * The name can contain lowercase letters, digits, and underscores (_). For MySQL databases, the name can contain uppercase letters.
    * The name must start with a letter and end with a letter or digit.
    * For MySQL databases, the name of the privileged account cannot be the same as that of the standard account. For example, if the name of the privileged account is Test1, the name of the standard account cannot be test1.
    * The length of the value must meet the following requirements:
        * If the instance runs MySQL 5.7 or MySQL 8.0, the value must be 2 to 32 characters in length.
        * If the instance runs MySQL 5.6, the value must be 2 to 16 characters in length.
        * If the instance runs SQL Server, the value must be 2 to 64 characters in length.
        * If the instance runs PostgreSQL with cloud disks, the value must be 2 to 63 characters in length.
        * If the instance runs PostgreSQL with local disks, the value must be 2 to 16 characters in length.
        * If the instance runs MariaDB, the value must be 2 to 16 characters in length.
        * For more information about invalid characters, See [Forbidden keywords](https://help.aliyun.com/zh/rds/developer-reference/forbidden-keywords?spm=api-workbench.API%20Document.0.0.529e2defHKoZ3o).

* `account_password` - (Optional, Sensitive) The password of the account.
    * The value must be 8 to 32 characters in length.
    * The password must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters.
    * Special characters include ! @ # $ % ^ & * ( ) _ + - =
* `account_type` - (Optional, ForceNew) The account type. Valid values:
    * Normal: standard account (default).
    * Super: privileged account.
    * Sysadmin: system admin account. The account type is available only for ApsaraDB RDS for SQL Server instances.

-> **NOTE:** Before you create a system admin account, check whether the RDS instance meets all prerequisites. For more information, See [Create a system admin account](https://help.aliyun.com/zh/rds/apsaradb-rds-for-sql-server/create-a-system-admin-account-for-an-apsaradb-rds-for-sql-server-instance?spm=api-workbench.API%20Document.0.0.529e2defHKoZ3o).
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `description` - (Optional, Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_description` instead.
* `instance_id` - (Optional, ForceNew, Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `db_instance_id` instead.
* `name` - (Optional, ForceNew, Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_name` instead.
* `password` - (Optional, Sensitive, Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_password` instead.
* `type` - (Optional, ForceNew, Deprecated from v1.120.0) The attribute has been deprecated from 1.120.0 and using `account_type` instead.

-> **NOTE**: Only MySQL engine is supported resets permissions of the privileged account.
* `reset_permission_flag` - (Optional, Available in v1.198.0+) Resets permissions flag of the privileged account. Default to `false`. Set it to `true` can resets permissions of the privileged account.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value is formatted `<db_instance_id>:<account_name>`.
* `status` - The status of the resource. Valid values: `Available`, `Unavailable`.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when create the Account.
* `update` - (Defaults to 6 mins) Used when update the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.

## Import

RDS Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_account.example <db_instance_id>:<account_name>
```