---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_enterprise_db_cluster_account"
description: |-
  Provides a Alicloud Click House Enterprise Db Cluster Account resource.
---

# alicloud_click_house_enterprise_db_cluster_account

Provides a Click House Enterprise Db Cluster Account resource.

Clickhouse enterprise instance account.

For information about Click House Enterprise Db Cluster Account and how to use it, see [What is Enterprise Db Cluster Account](https://next.api.alibabacloud.com/document/clickhouse/2023-05-22/CreateAccount).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_click_house_enterprise_db_cluster_account&exampleId=76541049-20f5-fd0d-3784-2722cd0b7be390911d2c&activeTab=example&spm=docs.r.click_house_enterprise_db_cluster_account.0.7654104920&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}

resource "alicloud_click_house_enterprise_db_cluster" "defaultWrovOd" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


resource "alicloud_click_house_enterprise_db_cluster_account" "default" {
  account        = "abc"
  description    = "example_desc"
  db_instance_id = alicloud_click_house_enterprise_db_cluster.defaultWrovOd.id
  account_type   = "NormalAccount"
  password       = "abc123456!"
  dml_auth_setting {
    dml_authority      = "0"
    ddl_authority      = true
    allow_dictionaries = ["*"]
    allow_databases    = ["*"]
  }
}
```

## Argument Reference

The following arguments are supported:
* `account` - (Required, ForceNew) The name of the database account.
* `account_type` - (Required, ForceNew) The type of the database account. Valid values:
  - `NormalAccount`: Normal account number.
  - `SuperAccount`: The privileged account.
* `db_instance_id` - (Required, ForceNew) The cluster ID.
* `description` - (Optional) Note information.
* `dml_auth_setting` - (Optional, List) Authorization information. See [`dml_auth_setting`](#dml_auth_setting) below.
* `password` - (Required) Database account password. Set the following rules.
  - Consists of at least three of uppercase letters, lowercase letters, numbers, and special characters.
  - Oh-! @#$%^& *()_+-= is a special character.
  - Length is 8~32 characters.

### `dml_auth_setting`

The dml_auth_setting supports the following:
* `allow_databases` - (Optional, List) The list of databases that require authorization. If there are more than one, separate them with commas (,).
* `allow_dictionaries` - (Optional, List) List of dictionaries that require authorization. If there are more than one, separate them with commas (,).
* `ddl_authority` - (Required) Whether to grant the DDL permission to the database account. Value description:
  - `true`: allows DDL.
  - `false`: DDL is disabled.
* `dml_authority` - (Required, Int) Whether to grant the DML permission to the database account. The values are as follows:
  - `0`: Queries that allow reading, writing, and changing settings
  - `1`: Only queries for reading data are allowed.
  - `2`: allows queries to read data and change settings.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<account>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Db Cluster Account.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Db Cluster Account.
* `update` - (Defaults to 5 mins) Used when update the Enterprise Db Cluster Account.

## Import

Click House Enterprise Db Cluster Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_enterprise_db_cluster_account.example <db_instance_id>:<account>
```