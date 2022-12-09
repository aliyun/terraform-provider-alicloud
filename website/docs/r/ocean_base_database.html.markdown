---
subcategory: "Ocean Base"
layout: "alicloud"
page_title: "Alicloud: alicloud_ocean_base_database"
sidebar_current: "docs-alicloud-resource-ocean-base-database"
description: |-
  Provides a Alicloud Ocean Base Database resource.
---

# alicloud\_ocean\_base\_database

Provides a Ocean Base Database resource.

For information about Ocean Base Database and how to use it, see [What is Database](https://www.alibabacloud.com/help/en/apsaradb-for-oceanbase).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ocean_base_database" "example" {
  database_name = "MyDatabase"
  encoding      = "example_value"
  tenant_id     = alicloud_ocean_base_tenant.example.tenant_id
}

```

## Argument Reference

The following arguments are supported:

* `collation` - (Optional, ForceNew) Character order.
* `database_name` - (Required, ForceNew) The name of the database. you cannot use some reserved keywords, such as test or mysql.
* `description` - (Optional) Database description information.
* `encoding` - (Required, ForceNew) The encoding method of the database. for more information, see the Charset Field returned by DescribeCharset.
* `instance_id` - (Required) OceanBase cluster ID.
* `tenant_id` - (Required, ForceNew) Tenant ID.
* `users` - (Optional) User name and list of roles.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Database. The value formats as `<tenant_id>:<database_name>`.
* `status` - The status of db. Valid Values: `ONLINE`and `DELETING`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 10 mins) Used when delete the Database.

## Import

Ocean Base Database can be imported using the id, e.g.

```shell
$ terraform import alicloud_ocean_base_database.example <tenant_id>:<database_name>
```