---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_logic_database"
sidebar_current: "docs-alicloud-resource-dms-enterprise-logic-database"
description: |-
  Provides a Alicloud DMS Enterprise Logic Database resource.
---

# alicloud_dms_enterprise_logic_database

Provides a DMS Enterprise Logic Database resource.

For information about DMS Enterprise Logic Database and how to use it, see [What is Logic Database](https://www.alibabacloud.com/help/en/dms/developer-reference/api-dms-enterprise-2018-11-01-createlogicdatabase).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dms_enterprise_logic_database&exampleId=e5cbccf8-1e4a-c603-64b2-a1738724408af3ec9f8d&activeTab=example&spm=docs.r.dms_enterprise_logic_database.0.e5cbccf81e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_dms_enterprise_logic_database" "default" {
  alias = "TF_logic_db_test"
  database_ids = [
    "35617919", "35617920"
  ]
}
```

## Argument Reference

The following arguments are supported:
* `alias` - (Required) Logical Library alias.
* `database_ids` - (Required) Sub-Database ID

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `db_type` - Database type.
* `env_type` - Environment type, return value is as follows:-product: production environment-dev: development environment-pre: Advance Environment-test: test environment-sit:SIT environment-uat:UAT environment-pet: Pressure measurement environment-stag:STAG environment
* `logic` - Whether it is a logical Library, the return value is true.
* `logic_database_id` - The ID of the logical Library.
* `owner_id_list` - The user ID list of the logical library Owner.
* `owner_name_list` - The nickname list of the logical library Owner.
* `schema_name` - Logical Library name.
* `search_name` - Logical library search name.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Logic Database.
* `delete` - (Defaults to 5 mins) Used when delete the Logic Database.
* `update` - (Defaults to 5 mins) Used when update the Logic Database.

## Import

DMS Enterprise Logic Database can be imported using the id, e.g.

```shell
$terraform import alicloud_dms_enterprise_logic_database.example <id>
```