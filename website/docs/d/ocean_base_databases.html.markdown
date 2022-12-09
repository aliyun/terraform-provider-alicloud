---
subcategory: "Ocean Base"
layout: "alicloud"
page_title: "Alicloud: alicloud_ocean_base_databases"
sidebar_current: "docs-alicloud-datasource-ocean-base-databases"
description: |-
  Provides a list of Ocean Base Databases to the user.
---

# alicloud\_ocean\_base\_databases

This data source provides the Ocean Base Databases of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ocean_base_databases" "ids" {
  tenant_id = "example_value"
  ids       = ["my-Database-1", "my-Database-2"]
}
output "ocean_base_database_id_1" {
  value = data.alicloud_ocean_base_databases.ids.databases.0.id
}

data "alicloud_ocean_base_databases" "nameRegex" {
  tenant_id  = "example_value"
  name_regex = "^my-Database"
}
output "ocean_base_database_id_2" {
  value = data.alicloud_ocean_base_databases.nameRegex.databases.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Database IDs. Its element value is same as Database Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Database name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of db. Valid Values: `ONLINE`and `DELETING`.
* `tenant_id` - (Required, ForceNew) Tenant ID.
* `with_tables` - (Optional, ForceNew) The with tables.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Database names.
* `databases` - A list of Ocean Base Databases. Each element contains the following attributes:
  * `collation` - Character order.
  * `database_name` - The name of the database. you cannot use some reserved keywords, such as test or mysql.
  * `description` - Database description information.
  * `encoding` - The encoding method of the database. for more information, see the Charset Field returned by DescribeCharset.
  * `id` - The ID of the Database. Its value is same as Queue Name.
  * `status` - The status of db.
  * `tenant_id` - Tenant ID.
  * `users` - User name and list of roles.
      * `user_type` - The account type of the users.
      * `role` - The role permissions given to the library by the account.
      * `user_name` - The account name of users.
  * `tables` - Database table information.
      * `table_name` - The database table name.