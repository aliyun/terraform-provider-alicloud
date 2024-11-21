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

For information about Click House Account and how to use it, see [What is Account](https://www.alibabacloud.com/help/zh/clickhouse/latest/api-clickhouse-2019-11-11-createaccount).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_click_house_account&exampleId=ac9776ab-f3a9-0dc6-7408-f41e60c156e19408086a&activeTab=example&spm=docs.r.click_house_account.0.ac9776abf3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
variable "type" {
  default = "Normal"
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
  zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.1.zone_id
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
  type                = var.type
}
```

## Argument Reference

The following arguments are supported:

* `account_description` - (Optional) In Chinese, English letter. May contain Chinese and English characters, lowercase letters, numbers, and underscores (_), the dash (-). Cannot start with http:// and https:// at the beginning. Length is from 2 to 256 characters.
* `account_name` - (Required, ForceNew) Account name: lowercase letters, numbers, underscores, lowercase letter; length no more than 16 characters.
* `account_password` - (Required) The account password: uppercase letters, lowercase letters, lowercase letters, numbers, and special characters (special character! #$%^& author (s):_+-=) in a length of 8-32 bit.
* `db_cluster_id` - (Required, ForceNew) The db cluster id.
* `type` - (Optional, ForceNew) The type of the database account. Valid values: `Normal` or `Super`.
* `dml_authority` - (Optional, Available since v1.163.0) Specifies whether to grant DML permissions to the database account. Valid values: `all` and `readOnly,modify`.
* `ddl_authority` - (Optional, Available since v1.163.0) Specifies whether to grant DDL permissions to the database account. Valid values: `true` and `false`.
  -`true`: grants DDL permissions to the database account.
  -`false`: does not grant DDL permissions to the database account.
* `allow_databases` - (Optional, Available since v1.163.0) The list of databases to which you want to grant permissions. Separate databases with commas (,).
* `total_databases` - (Optional, Deprecated since v1.223.1) The list of all databases. Separate databases with commas (,). Field 'total_databases' has been deprecated from provider version 1.223.1.
* `allow_dictionaries` - (Optional, Available since v1.163.0) The list of dictionaries to which you want to grant permissions. Separate dictionaries with commas (,).
* `total_dictionaries` - (Optional, Deprecated since v1.223.1) The list of all dictionaries. Separate dictionaries with commas (,). Field 'total_dictionaries' has been deprecated from provider version 1.223.1.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value formats as `<db_cluster_id>:<account_name>`.
* `status` - The status of the resource. Valid Status: `Creating`,`Available`,`Deleting`.

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
