---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_accounts"
sidebar_current: "docs-alicloud-datasource-click-house-accounts"
description: |-
  Provides a list of Click House Accounts to the user.
---

# alicloud_click_house_accounts

This data source provides the Click House Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "oneaccountname"
}
variable "pwd" {
  default = "Tf-onepwd"
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
  zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "22.8.5.29"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = alicloud_vswitch.default.id
  vpc_id                  = alicloud_vpc.default.id
}

resource "alicloud_click_house_account" "default" {
  db_cluster_id       = alicloud_click_house_db_cluster.default.id
  account_description = "your_description"
  account_name        = var.name
  account_password    = var.pwd
  type                = var.type
}

data "alicloud_click_house_accounts" "default" {
  ids           = [alicloud_click_house_account.default.id]
  db_cluster_id = alicloud_click_house_db_cluster.default.id
}

output "account_id" {
  value = data.alicloud_click_house_accounts.default.ids.0
}

```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The DBCluster id.
* `ids` - (Optional, ForceNew, Computed)  A list of Account IDs. Its element value is same as Account Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Account name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid Status: `Creating`,`Available`,`Deleting`.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Account names.
* `accounts` - A list of Click House Accounts. Each element contains the following attributes:
	* `account_description` - In Chinese, English letter. May contain Chinese and English characters, lowercase letters, numbers, and underscores (_), the dash (-). Cannot start with http:// and https:// at the beginning. Length is from 2 to 256 characters.
	* `account_name` - Account name: lowercase letters, numbers, underscores, lowercase letter; length no more than 16 characters.
	* `account_type` - The Valid Account type: `Normal`, `Super`.
	* `db_cluster_id` - The DBCluster id.
	* `id` - The ID of the Account. Its value is same as Queue Name.
	* `status` - The status of the resource.
