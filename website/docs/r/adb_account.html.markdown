---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_account"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) Account resource.
---

# alicloud_adb_account

Provides a AnalyticDB for MySQL (ADB) Account resource.



For information about AnalyticDB for MySQL (ADB) Account and how to use it, see [What is Account](https://next.api.alibabacloud.com/document/adb/2019-03-15/CreateAccount).

-> **NOTE:** Available since v1.71.0.

## Example Usage

Basic Usage

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
* `account_description` - (Optional) The description of the account
* `account_name` - (Required, ForceNew) The name of the account
* `account_password` - (Required) The password of the ADB account.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `account_type` - (Optional) The type of the account
* `db_cluster_id` - (Required, ForceNew) The DBCluster ID
* `tag` - (Optional, ForceNew, List, Available since v1.270.0) The tag of the resource See [`tag`](#tag) below.

### `tag`

The tag supports the following:
* `key` - (Optional, ForceNew, Available since v1.270.0) The key of the tags
* `value` - (Optional, ForceNew, Available since v1.270.0) The value of the tags

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<db_cluster_id>:<account_name>`.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `delete` - (Defaults to 5 mins) Used when delete the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

AnalyticDB for MySQL (ADB) Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_account.example <db_cluster_id>:<account_name>
```