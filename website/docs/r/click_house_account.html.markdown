---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_account"
sidebar_current: "docs-alicloud-resource-click-house-account"
description: |-
  Provides a Alicloud Click House Account resource.
---

# alicloud_click_house_account

Provides a Click House Account resource.

For information about Click House Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/en/clickhouse/latest/api-clickhouse-2019-11-11-createaccount).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_click_house_regions" "default" {
  current = true
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "22.8.5.29"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = alicloud_vswitch.default.id
  vpc_id                  = alicloud_vpc.default.id
}

resource "alicloud_click_house_account" "default" {
  db_cluster_id       = alicloud_click_house_db_cluster.default.id
  account_description = "tf-example-description"
  account_name        = "examplename"
  account_password    = "Example1234"
}
```

## Argument Reference

The following arguments are supported:

* `account_description` - (Optional) In Chinese, English letter. May contain Chinese and English characters, lowercase letters, numbers, and underscores (_), the dash (-). Cannot start with http:// and https:// at the beginning. Length is from 2 to 256 characters.
* `account_name` - (Required, ForceNew) Account name: lowercase letters, numbers, underscores, lowercase letter; length no more than 16 characters.
* `account_password` - (Required) The account password: uppercase letters, lowercase letters, lowercase letters, numbers, and special characters (special character! #$%^& author (s):_+-=) in a length of 8-32 bit.
* `db_cluster_id` - (Required, ForceNew) The db cluster id.
* `dml_authority` - (Optional, Available since v1.163.0) Specifies whether to grant DML permissions to the database account. Valid values: `all` and `readOnly,modify`.
* `ddl_authority` - (Optional, Available since v1.163.0) Specifies whether to grant DDL permissions to the database account. Valid values: `true` and `false`.
  -`true`: grants DDL permissions to the database account.
  -`false`: does not grant DDL permissions to the database account.
* `allow_databases` - (Optional, Available since v1.163.0) The list of databases to which you want to grant permissions. Separate databases with commas (,).
* `total_databases` - (Optional, Available since v1.163.0) The list of all databases. Separate databases with commas (,).
* `allow_dictionaries` - (Optional, Available since v1.163.0) The list of dictionaries to which you want to grant permissions. Separate dictionaries with commas (,).
* `total_dictionaries` - (Optional, Available since v1.163.0) The list of all dictionaries. Separate dictionaries with commas (,).


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value formats as `<db_cluster_id>:<account_name>`.
* `status` - The status of the resource. Valid Status: `Creating`,`Available`,`Deleting`.
* `type` - The type of the database account. Valid values: `Normal` or `Super`.

## Timeouts

-> **NOTE:** Available since v1.163.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Click House Account.
* `update` - (Defaults to 1 mins) Used when update the Click House Account.
* `delete` - (Defaults to 1 mins) Used when delete the Click House Account.

## Import

Click House Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_account.example <db_cluster_id>:<account_name>
```
