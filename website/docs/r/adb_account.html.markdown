---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_account"
sidebar_current: "docs-alicloud-resource-adb-account"
description: |-
  Provides a ADB account resource.
---

# alicloud_adb_account

Provides a [ADB](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/api-doc-adb-2019-03-15-api-doc-createaccount) account resource and used to manage databases.

-> **NOTE:** Available since v1.71.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_account&exampleId=cf607244-cb42-08ac-56a7-3f3cf863b6ae3bca8348&activeTab=example&spm=docs.r.adb_account.0.cf607244cb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "creation" {
  default = "ADB"
}

variable "name" {
  default = "tfexample"
}

data "alicloud_adb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_adb_db_cluster" "cluster" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "8Core32GB"
  vswitch_id          = local.vswitch_id
  description         = var.name
}

resource "alicloud_adb_account" "default" {
  db_cluster_id       = alicloud_adb_db_cluster.cluster.id
  account_name        = var.name
  account_password    = "tf_example123"
  account_description = var.name
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster in which account belongs.
* `account_name` - (Required, ForceNew) Operation account requiring a uniqueness check. It may consist of lower case letters, numbers, and underlines, and must start with a letter and have no more than 16 characters.
* `account_password` - (Optional) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters. You have to specify one of `account_password` and `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. If the `account_password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `account_description` - (Optional) Account description. It cannot begin with https://. It must start with a Chinese character or English letter. It can include Chinese and English characters, underlines (_), hyphens (-), and numbers. The length may be 2-256 characters.
* `account_type` - (Optional, ForceNew, Removed since provider version 1.81.0) The type of the database account. Default Value: `Super`. Valid values:
  - `Normal`: standard account. Up to 256 standard accounts can be created for a cluster.
  - `Super`: privileged account. Only a single privileged account can be created for a cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<name>`.

## Import

ADB account can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_account.example am-12345:tf_account
```
