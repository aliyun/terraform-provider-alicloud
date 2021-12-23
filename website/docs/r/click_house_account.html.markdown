---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_account"
sidebar_current: "docs-alicloud-resource-click-house-account"
description: |-
  Provides a Alicloud Click House Account resource.
---

# alicloud\_click\_house\_account

Provides a Click House Account resource.

For information about Click House Account and how to use it, see [What is Account](https://www.alibabacloud.com/product/clickhouse).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "testaccountname"
}
variable "pwd" {
  default = "Tf-testpwd"
}
resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = "your_vswitch_id"
}
resource "alicloud_click_house_account" "default" {
  db_cluster_id       = alicloud_click_house_db_cluster.default.id
  account_description = "your_description"
  account_name        = var.name
  account_password    = var.pwd
}

```

## Argument Reference

The following arguments are supported:

* `account_description` - (Optional) In Chinese, English letter. May contain Chinese and English characters, lowercase letters, numbers, and underscores (_), the dash (-). Cannot start with http:// and https:// at the beginning. Length is from 2 to 256 characters.
* `account_name` - (Required, ForceNew) Account name: lowercase letters, numbers, underscores, lowercase letter; length no more than 16 characters.
* `account_password` - (Required) The account password: uppercase letters, lowercase letters, lowercase letters, numbers, and special characters (special character! #$%^& author (s):_+-=) in a length of 8-32 bit.
* `db_cluster_id` - (Required, ForceNew) The db cluster id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Account. The value formats as `<db_cluster_id>:<account_name>`.
* `status` - The status of the resource. Valid Status: `Creating`,`Available`,`Deleting`.

## Import

Click House Account can be imported using the id, e.g.

```
$ terraform import alicloud_click_house_account.example <db_cluster_id>:<account_name>
```
