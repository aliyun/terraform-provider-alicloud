---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_databases"
sidebar_current: "docs-alicloud-datasource-rds-databases"
description: |-
  Provides a list of databases of an ApsaraDB RDS instance.
---

# alicloud_rds_databases

The `alicloud_rds_databases` data source provides a collection of databases of an ApsaraDB RDS instance. You can filter the results by the database name and status on the server side, and by a regular expression on the provider side.

For information on the API and its fields, see [DescribeDatabases](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-describedatabases).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** The data source reads every page returned by `DescribeDatabases` before applying `name_regex`, so the result is complete regardless of the number of databases on the instance. Engine-specific fields are not exported in this version: the SQL Server `AdvancedInfo`, `BasicInfo` and `RuntimeInfo`, the PostgreSQL `Collate`, `Ctype`, `ConnLimit` and `Tablespace`, and `DuckDBEnabled`.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
}

resource "alicloud_db_database" "default" {
  count          = 2
  instance_id    = alicloud_db_instance.default.id
  data_base_name = "terraform_example_${count.index}"
  character_set  = "utf8"
  description    = "terraform example database ${count.index}"
}

data "alicloud_rds_databases" "default" {
  # The splat reads the instance id from every database instance, so the data
  # source waits for all databases to be created before it reads.
  instance_id = element(alicloud_db_database.default[*].instance_id, 0)
  name_regex  = "^terraform_example_"
}

output "database_names" {
  value = data.alicloud_rds_databases.default.names
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of the RDS instance.
* `db_name` - (Optional, String) The name of the database. The value is sent to the server as the `DBName` filter.
* `db_status` - (Optional, String) The status of the database to filter. Valid values: `Creating`, `Running` and `Deleting`. The value is sent to the server as the `DBStatus` filter and is also applied on the provider side.
* `name_regex` - (Optional, String) A regex string to filter results by database name. The filter is applied on the provider side after all pages are read. An invalid regular expression fails at plan time.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Read-Only, String) The ID of the data source, derived from `instance_id` and the filters so that it stays stable across refreshes.
* `names` - (Read-Only, List) A list of database names that match the configured filters.
* `databases` - (Read-Only, List) A list of databases that match the configured filters. Each element contains the following attributes:
  * `db_name` - (Read-Only, String) The name of the database.
  * `db_status` - (Read-Only, String) The status of the database. Valid values: `Creating`, `Running`, `Deleting` and `Cold`.
  * `character_set_name` - (Read-Only, String) The character set of the database.
  * `db_description` - (Read-Only, String) The description of the database.
  * `engine` - (Read-Only, String) The database engine of the instance.
  * `instance_id` - (Read-Only, String) The ID of the instance that the database belongs to.
  * `accounts` - (Read-Only, List) A list of accounts that have privileges on the database. Each element contains the following attributes:
    * `account_name` - (Read-Only, String) The name of the account.
    * `account_privilege` - (Read-Only, String) The privilege that the account has on the database. Valid values: `ReadWrite`, `ReadOnly`, `DMLOnly` and `DDLOnly`.
    * `account_privilege_detail` - (Read-Only, String) The detailed privilege that the account has on the database.
