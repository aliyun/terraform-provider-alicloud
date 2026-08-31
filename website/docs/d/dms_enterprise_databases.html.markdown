---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_databases"
sidebar_current: "docs-alicloud-datasource-dms-enterprise-databases"
description: |-
  Provides a list of DMS Enterprise Database owned by an Alibaba Cloud account.
---

# alicloud_dms_enterprise_databases

This data source provides DMS Enterprise Database available to the user. [What is Database](https://www.alibabacloud.com/help/en/dms/developer-reference/api-dms-enterprise-2018-11-01-listdatabases).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_dms_enterprise_databases" "default" {
  name_regex  = "test2"
  instance_id = "2195118"
}

output "alicloud_dms_enterprise_database_example_id" {
  value = data.alicloud_dms_enterprise_databases.default.databases.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Database IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter the results by the database Schema Name.
* `instance_id` - (Required, ForceNew) The instance ID of the target database.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Database IDs.
* `databases` - A list of Database Entries. Each element contains the following attributes:
  * `id` - The Database ID, same value as `database_id`.
  * `catalog_name` - The name of the Database Directory.> PG Series databases will display the database name.
  * `database_id` - The ID of the physical library.
  * `db_type` - Database type.
  * `dba_id` - The DBA user ID of the target database.
  * `dba_name` - The DBA nickname of the target Library.
  * `encoding` - Database encoding.
  * `env_type` - The environment type of the database.
  * `host` - The database connection address.
  * `instance_id` - The instance ID of the target database.
  * `owner_id_list` - Library Owner User ID list.
  * `owner_name_list` - Library Owner nickname list.
  * `port` - The connection port of the database.
  * `schema_name` - The name of the database.> PG Series databases will display schema names.
  * `search_name` - Library search name.
  * `sid` - Database SID.> only Oracle Database Display.
  * `state` - Library status, value description:-**NORMAL**: NORMAL-**DISABLE**: Disabled-**OFFLINE**: OFFLINE-**NOT_EXIST**: does not exist
