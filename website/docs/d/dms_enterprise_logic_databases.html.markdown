---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_logic_databases"
sidebar_current: "docs-alicloud-datasource-dms-enterprise-logic-databases"
description: |-
  Provides a list of DMS Enterprise Logic Database owned by an Alibaba Cloud account.
---

# alicloud_dms_enterprise_logic_databases

This data source provides DMS Enterprise Logic Database available to the user. [What is Logic Database](https://www.alibabacloud.com/help/en/dms/developer-reference/api-dms-enterprise-2018-11-01-createlogicdatabase).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_dms_enterprise_instances" "dms_enterprise_instances_ds" {
  instance_type = "mysql"
  search_key    = "tf-test-no-deleting"
}

data "alicloud_dms_enterprise_logic_databases" "default" {
  instance_id = data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id
}

output "alicloud_dms_enterprise_logic_database_example_id" {
  value = data.alicloud_dms_enterprise_logic_databases.default.databases.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Logic Database IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Logic Database IDs.
* `databases` - A list of Logic Database Entries. Each element contains the following attributes:
  * `id` - Logic Database ID.
  * `alias` - Logical Library alias.
  * `database_ids` - Sub-Database ID.
  * `db_type` - Database type.
  * `env_type` - Environment type, return value is as follows:-product: production environment-dev: development environment-pre: Advance Environment-test: test environment-sit:SIT environment-uat:UAT environment-pet: Pressure measurement environment-stag:STAG environment
  * `logic` - Whether it is a logical Library, the return value is true.
  * `logic_database_id` - The ID of the logical Library.
  * `owner_id_list` - The user ID list of the logical library Owner.
  * `owner_name_list` - The nickname list of the logical library Owner.
  * `schema_name` - Logical Library name.
  * `search_name` - Logical library search name.
