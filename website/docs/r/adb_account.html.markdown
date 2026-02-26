---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_account"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) Account resource.
---

# alicloud_adb_account

Provides a AnalyticDB for MySQL (ADB) Account resource.



For information about AnalyticDB for MySQL (ADB) Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/api-doc-adb-2019-03-15-api-doc-createaccount).

-> **NOTE:** Available since v1.71.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_account&exampleId=cf607244-cb42-08ac-56a7-3f3cf863b6ae3bca8348&activeTab=example&spm=docs.r.adb_account.0.cf607244cb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}

data "alicloud_adb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}

resource "alicloud_adb_db_cluster" "cluster" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "8Core32GB"
  vswitch_id          = data.alicloud_vswitches.default.ids.0
  description         = var.name
}

resource "alicloud_adb_account" "default" {
  db_cluster_id       = alicloud_adb_db_cluster.cluster.id
  account_name        = var.name
  account_password    = "tf_example123"
  account_description = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_adb_account&spm=docs.r.adb_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `account_description` - (Optional) The description of the account.
* `account_name` - (Required, ForceNew) The name of the database account. The name must meet the following requirements:
  - Start with a lowercase letter and end with a lowercase letter or a digit.
  - Contain only lowercase letters, digits, and underscores (_).
  - Its length must be between 2 and 16 characters.
  - Cannot be a reserved name, such as root, admin, or opsadmin.
* `account_password` - (Optional) The password of the database account. The password must meet the following requirements:
  - It must consist of uppercase letters, lowercase letters, digits, and special characters.
  - The allowed special characters are: (!), (@), (#), ($), (%), (^), (&), (*), (()), (_), (+), (-), (=).
  - Its length must be between 8 and 32 characters.
* `account_type` - (Optional, ForceNew, Available since v1.272.0) The type of the account. Valid values:
  - `Normal`: A standard account. You can create up to 256 standard accounts for a cluster.
  - `Super` (default): A privileged account. You can create only one privileged account for a cluster.
* `db_cluster_id` - (Required, ForceNew) The cluster ID of the data warehouse edition.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `tags` - (Optional, ForceNew, Map, Available since v1.272.0) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<db_cluster_id>:<account_name>`.
* `status` - (Available since v1.272.0) The status of the account.

## Timeouts

-> **NOTE:** Available since v1.272.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

AnalyticDB for MySQL (ADB) Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_account.example <db_cluster_id>:<account_name>
```
