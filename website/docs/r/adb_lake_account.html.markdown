---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_lake_account"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) Lake Account resource.
---

# alicloud_adb_lake_account

Provides a AnalyticDB for MySQL (ADB) Lake Account resource.

Account of the DBClusterLakeVesion.

For information about AnalyticDB for MySQL (ADB) Lake Account and how to use it, see [What is Lake Account](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/developer-reference/api-adb-2021-12-01-modifyaccountprivileges/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_lake_account&exampleId=83d73ad9-9ed6-3bbc-31fd-a4f1d2a9eb4948b0743f&activeTab=example&spm=docs.r.adb_lake_account.0.83d73ad99e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "VPCID" {
  vpc_name = var.name

  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "VSWITCHID" {
  vpc_id       = alicloud_vpc.VPCID.id
  zone_id      = "cn-hangzhou-k"
  vswitch_name = var.name

  cidr_block = "172.16.0.0/24"
}

resource "alicloud_adb_db_cluster_lake_version" "CreateInstance" {
  storage_resource       = "0ACU"
  zone_id                = "cn-hangzhou-k"
  vpc_id                 = alicloud_vpc.VPCID.id
  vswitch_id             = alicloud_vswitch.VSWITCHID.id
  db_cluster_description = var.name
  compute_resource       = "16ACU"
  db_cluster_version     = "5.0"
  payment_type           = "PayAsYouGo"
  security_ips           = "127.0.0.1"
}


resource "alicloud_adb_lake_account" "default" {
  db_cluster_id    = alicloud_adb_db_cluster_lake_version.CreateInstance.id
  account_type     = "Super"
  account_name     = "tfnormal"
  account_password = "normal@2023"
  account_privileges {
    privilege_type = "Database"
    privilege_object {
      database = "MYSQL"
    }

    privileges = [
      "select",
      "update"
    ]
  }
  account_privileges {
    privilege_type = "Table"
    privilege_object {
      database = "INFORMATION_SCHEMA"
      table    = "ENGINES"
    }

    privileges = [
      "update"
    ]
  }
  account_privileges {
    privilege_type = "Column"
    privilege_object {
      table    = "COLUMNS"
      column   = "PRIVILEGES"
      database = "INFORMATION_SCHEMA"
    }

    privileges = [
      "update"
    ]
  }

  account_description = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_adb_lake_account&spm=docs.r.adb_lake_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `account_description` - (Optional) The description of the account.
* `account_name` - (Required, ForceNew) The name of the account.
* `account_password` - (Required) AccountPassword.
* `account_privileges` - (Optional) List of permissions granted. See [`account_privileges`](#account_privileges) below.
* `account_type` - (Optional, ForceNew) The type of the account.
* `db_cluster_id` - (Required, ForceNew) The DBCluster ID.
* `engine` - (Optional, ForceNew) The engine type of the account. Valid values: `AnalyticDB` (default): the AnalyticDB MySQL engine. `Clickhouse`: the wide-table engine. The engine is immutable after creation and is forwarded to the account lifecycle APIs (Describe, Modify, Reset password, Delete) that accept it. `DescribeAccounts` does not return `Engine` for default-engine (`AnalyticDB`) accounts, so the value configured at create time is preserved on refresh rather than overwritten; for such accounts the attribute may not be read back during a legacy two-segment import. `DescribeAccounts` filters by `Engine` and defaults to `AnalyticDB`, so `Clickhouse` accounts cannot be located by the default query; import a `Clickhouse` account with the three-segment id `<db_cluster_id>:<account_name>:Clickhouse` so the engine is recorded before read.
* `ram_user_list` - (Optional, List, Available since v1.274.0) List of Alibaba Cloud RAM user IDs to bind.

### `account_privileges`

The account_privileges supports the following:
* `privilege_object` - (Optional) Object associated to privileges. See [`privilege_object`](#account_privileges-privilege_object) below.
* `privilege_type` - (Optional) The type of privileges.
* `privileges` - (Optional) privilege list.

### `account_privileges-privilege_object`

The privilege_object supports the following:
* `column` - (Optional) The name of column.
* `database` - (Optional) The name of database.
* `table` - (Optional) The name of table.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_cluster_id>:<account_name>`. Import supports the three-segment form `<db_cluster_id>:<account_name>:<engine>`; the `<engine>` segment is required for `Clickhouse` accounts and optional (defaults to `AnalyticDB`) for the legacy two-segment form.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Lake Account.
* `delete` - (Defaults to 5 mins) Used when delete the Lake Account.
* `update` - (Defaults to 5 mins) Used when update the Lake Account.

## Import

ADB Lake Account can be imported using the id, which consists of db_cluster_id and account_name, e.g.

```shell
$ terraform import alicloud_adb_lake_account.example <db_cluster_id>:<account_name>
```
